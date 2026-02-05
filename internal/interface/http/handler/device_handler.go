package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/suquant/wgrest/internal/domain/entity"
	"github.com/suquant/wgrest/internal/usecase"
)

// DeviceHandler handles HTTP requests for device operations.
type DeviceHandler struct {
	useCase *usecase.DeviceUseCase
}

// NewDeviceHandler creates a new device handler.
func NewDeviceHandler(uc *usecase.DeviceUseCase) *DeviceHandler {
	return &DeviceHandler{useCase: uc}
}

// ListDevices godoc
// @Summary List all WireGuard devices
// @Tags Devices
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(0)
// @Param per_page query int false "Items per page" default(100)
// @Success 200 {array} entity.Device
// @Failure 500 {object} entity.Error
// @Security BearerAuth
// @Router /devices/ [get]
func (h *DeviceHandler) ListDevices(c *fiber.Ctx) error {
	page := c.QueryInt("page", 0)
	perPage := c.QueryInt("per_page", 100)

	devices, total, err := h.useCase.ListDevices(page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entity.Error{
			Code:    entity.ErrCodeInternalError,
			Message: err.Error(),
		})
	}

	// Set Link header for pagination
	setLinkHeader(c, page, perPage, total)

	return c.JSON(devices)
}

// GetDevice godoc
// @Summary Get a specific device
// @Tags Devices
// @Accept json
// @Produce json
// @Param name path string true "Device name"
// @Success 200 {object} entity.Device
// @Failure 404 {object} entity.Error
// @Security BearerAuth
// @Router /devices/{name}/ [get]
func (h *DeviceHandler) GetDevice(c *fiber.Ctx) error {
	name := c.Params("name")

	device, err := h.useCase.GetDevice(name)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(entity.Error{
			Code:    entity.ErrCodeDeviceNotFound,
			Message: err.Error(),
		})
	}

	return c.JSON(device)
}

// CreateDevice godoc
// @Summary Create a new WireGuard device
// @Tags Devices
// @Accept json
// @Produce json
// @Param request body entity.DeviceCreateOrUpdateRequest true "Device creation request"
// @Success 201 {object} entity.Device
// @Failure 400 {object} entity.Error
// @Failure 409 {object} entity.Error
// @Failure 500 {object} entity.Error
// @Security BearerAuth
// @Router /devices/ [post]
func (h *DeviceHandler) CreateDevice(c *fiber.Ctx) error {
	var req entity.DeviceCreateOrUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.Error{
			Code:    entity.ErrCodeInvalidRequest,
			Message: err.Error(),
		})
	}

	device, err := h.useCase.CreateDevice(req)
	if err != nil {
		// Check if device already exists
		if isDeviceExistsError(err) {
			return c.Status(fiber.StatusConflict).JSON(entity.Error{
				Code:    entity.ErrCodeDeviceExists,
				Message: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(entity.Error{
			Code:    entity.ErrCodeInternalError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(device)
}

// UpdateDevice godoc
// @Summary Update a device
// @Tags Devices
// @Accept json
// @Produce json
// @Param name path string true "Device name"
// @Param request body entity.DeviceCreateOrUpdateRequest true "Device update request"
// @Success 200 {object} entity.Device
// @Failure 400 {object} entity.Error
// @Failure 404 {object} entity.Error
// @Security BearerAuth
// @Router /devices/{name}/ [patch]
func (h *DeviceHandler) UpdateDevice(c *fiber.Ctx) error {
	name := c.Params("name")

	var req entity.DeviceCreateOrUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.Error{
			Code:    entity.ErrCodeInvalidRequest,
			Message: err.Error(),
		})
	}

	device, err := h.useCase.UpdateDevice(name, req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(entity.Error{
			Code:    entity.ErrCodeDeviceNotFound,
			Message: err.Error(),
		})
	}

	return c.JSON(device)
}

// DeleteDevice godoc
// @Summary Delete a device
// @Tags Devices
// @Param name path string true "Device name"
// @Success 204 "No Content"
// @Failure 404 {object} entity.Error
// @Security BearerAuth
// @Router /devices/{name}/ [delete]
func (h *DeviceHandler) DeleteDevice(c *fiber.Ctx) error {
	name := c.Params("name")

	if err := h.useCase.DeleteDevice(name); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(entity.Error{
			Code:    entity.ErrCodeDeviceNotFound,
			Message: err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Up godoc
// @Summary Bring interface up (wg-quick up)
// @Tags Devices
// @Produce json
// @Param name path string true "Device name"
// @Success 200 {object} map[string]string
// @Failure 500 {object} entity.Error
// @Security BearerAuth
// @Router /devices/{name}/up/ [post]
func (h *DeviceHandler) Up(c *fiber.Ctx) error {
	name := c.Params("name")

	if err := h.useCase.Up(name); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entity.Error{
			Code:    entity.ErrCodeInternalError,
			Message: err.Error(),
		})
	}

	return c.JSON(fiber.Map{"status": "up", "interface": name})
}

// Down godoc
// @Summary Bring interface down (wg-quick down)
// @Tags Devices
// @Produce json
// @Param name path string true "Device name"
// @Success 200 {object} map[string]string
// @Failure 500 {object} entity.Error
// @Security BearerAuth
// @Router /devices/{name}/down/ [post]
func (h *DeviceHandler) Down(c *fiber.Ctx) error {
	name := c.Params("name")

	if err := h.useCase.Down(name); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entity.Error{
			Code:    entity.ErrCodeInternalError,
			Message: err.Error(),
		})
	}

	return c.JSON(fiber.Map{"status": "down", "interface": name})
}

func isDeviceExistsError(err error) bool {
	return err != nil && (err.Error() == "device already exists" ||
		contains(err.Error(), "already exists"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
