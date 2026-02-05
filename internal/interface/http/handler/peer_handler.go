package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/suquant/wgrest/internal/domain/entity"
	"github.com/suquant/wgrest/internal/usecase"
)

// PeerHandler handles HTTP requests for peer operations.
type PeerHandler struct {
	useCase *usecase.PeerUseCase
}

// NewPeerHandler creates a new peer handler.
func NewPeerHandler(uc *usecase.PeerUseCase) *PeerHandler {
	return &PeerHandler{useCase: uc}
}

// ListPeers godoc
// @Summary List all peers for a device
// @Tags Peers
// @Accept json
// @Produce json
// @Param name path string true "Device name"
// @Param page query int false "Page number" default(0)
// @Param per_page query int false "Items per page" default(100)
// @Param q query string false "Search by allowed IPs"
// @Param sort query string false "Sort field (prefix with - for desc)" Enums(pub_key, -pub_key, receive_bytes, -receive_bytes, transmit_bytes, -transmit_bytes, total_bytes, -total_bytes, last_handshake_time, -last_handshake_time)
// @Success 200 {array} entity.Peer
// @Failure 404 {object} entity.Error
// @Failure 500 {object} entity.Error
// @Security BearerAuth
// @Router /devices/{name}/peers/ [get]
func (h *PeerHandler) ListPeers(c *fiber.Ctx) error {
	deviceName := c.Params("name")
	page := c.QueryInt("page", 0)
	perPage := c.QueryInt("per_page", 100)
	query := c.Query("q")
	sort := c.Query("sort")

	peers, total, err := h.useCase.ListPeers(deviceName, page, perPage, query, sort)
	if err != nil {
		if isDeviceNotFoundError(err) {
			return c.Status(fiber.StatusNotFound).JSON(entity.Error{
				Code:    entity.ErrCodeDeviceNotFound,
				Message: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(entity.Error{
			Code:    entity.ErrCodeInternalError,
			Message: err.Error(),
		})
	}

	// Set Link header for pagination
	setLinkHeader(c, page, perPage, total)

	return c.JSON(peers)
}

// GetPeer godoc
// @Summary Get a specific peer
// @Tags Peers
// @Accept json
// @Produce json
// @Param name path string true "Device name"
// @Param urlSafePubKey path string true "URL-safe base64 encoded public key"
// @Success 200 {object} entity.Peer
// @Failure 404 {object} entity.Error
// @Security BearerAuth
// @Router /devices/{name}/peers/{urlSafePubKey}/ [get]
func (h *PeerHandler) GetPeer(c *fiber.Ctx) error {
	deviceName := c.Params("name")
	urlSafePubKey := c.Params("urlSafePubKey")

	peer, err := h.useCase.GetPeer(deviceName, urlSafePubKey)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(entity.Error{
			Code:    entity.ErrCodePeerNotFound,
			Message: err.Error(),
		})
	}

	return c.JSON(peer)
}

// CreatePeer godoc
// @Summary Create a new peer
// @Tags Peers
// @Accept json
// @Produce json
// @Param name path string true "Device name"
// @Param request body entity.PeerCreateOrUpdateRequest true "Peer creation request"
// @Success 201 {object} entity.Peer
// @Failure 400 {object} entity.Error
// @Failure 404 {object} entity.Error
// @Failure 500 {object} entity.Error
// @Security BearerAuth
// @Router /devices/{name}/peers/ [post]
func (h *PeerHandler) CreatePeer(c *fiber.Ctx) error {
	deviceName := c.Params("name")

	var req entity.PeerCreateOrUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.Error{
			Code:    entity.ErrCodeInvalidRequest,
			Message: err.Error(),
		})
	}

	peer, err := h.useCase.CreatePeer(deviceName, req)
	if err != nil {
		if isDeviceNotFoundError(err) {
			return c.Status(fiber.StatusNotFound).JSON(entity.Error{
				Code:    entity.ErrCodeDeviceNotFound,
				Message: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(entity.Error{
			Code:    entity.ErrCodeInternalError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(peer)
}

// UpdatePeer godoc
// @Summary Update a peer
// @Tags Peers
// @Accept json
// @Produce json
// @Param name path string true "Device name"
// @Param urlSafePubKey path string true "URL-safe base64 encoded public key"
// @Param request body entity.PeerCreateOrUpdateRequest true "Peer update request"
// @Success 200 {object} entity.Peer
// @Failure 400 {object} entity.Error
// @Failure 404 {object} entity.Error
// @Security BearerAuth
// @Router /devices/{name}/peers/{urlSafePubKey}/ [patch]
func (h *PeerHandler) UpdatePeer(c *fiber.Ctx) error {
	deviceName := c.Params("name")
	urlSafePubKey := c.Params("urlSafePubKey")

	var req entity.PeerCreateOrUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.Error{
			Code:    entity.ErrCodeInvalidRequest,
			Message: err.Error(),
		})
	}

	peer, err := h.useCase.UpdatePeer(deviceName, urlSafePubKey, req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(entity.Error{
			Code:    entity.ErrCodePeerNotFound,
			Message: err.Error(),
		})
	}

	return c.JSON(peer)
}

// DeletePeer godoc
// @Summary Delete a peer
// @Tags Peers
// @Param name path string true "Device name"
// @Param urlSafePubKey path string true "URL-safe base64 encoded public key"
// @Success 200 {object} entity.Peer
// @Failure 404 {object} entity.Error
// @Security BearerAuth
// @Router /devices/{name}/peers/{urlSafePubKey}/ [delete]
func (h *PeerHandler) DeletePeer(c *fiber.Ctx) error {
	deviceName := c.Params("name")
	urlSafePubKey := c.Params("urlSafePubKey")

	peer, err := h.useCase.DeletePeer(deviceName, urlSafePubKey)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(entity.Error{
			Code:    entity.ErrCodePeerNotFound,
			Message: err.Error(),
		})
	}

	return c.JSON(peer)
}

func isDeviceNotFoundError(err error) bool {
	return err != nil && contains(err.Error(), "not found")
}


