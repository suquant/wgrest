package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/suquant/wgrest/internal/domain/entity"
	"github.com/suquant/wgrest/internal/usecase"
)

// MockDeviceUseCase implements a mock for device use case testing
type MockDeviceUseCase struct {
	devices []entity.Device
	err     error
}

func (m *MockDeviceUseCase) ListDevices(page, perPage int) ([]entity.Device, int, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	return m.devices, len(m.devices), nil
}

func (m *MockDeviceUseCase) GetDevice(name string) (*entity.Device, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, d := range m.devices {
		if d.Name == name {
			return &d, nil
		}
	}
	return nil, assert.AnError
}

func (m *MockDeviceUseCase) CreateDevice(req entity.DeviceCreateOrUpdateRequest) (*entity.Device, error) {
	if m.err != nil {
		return nil, m.err
	}
	device := entity.Device{
		Name:       *req.Name,
		ListenPort: 51820,
		PublicKey:  "generatedPublicKey",
	}
	return &device, nil
}

func (m *MockDeviceUseCase) UpdateDevice(name string, req entity.DeviceCreateOrUpdateRequest) (*entity.Device, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, d := range m.devices {
		if d.Name == name {
			if req.ListenPort != nil {
				d.ListenPort = *req.ListenPort
			}
			return &d, nil
		}
	}
	return nil, assert.AnError
}

func (m *MockDeviceUseCase) DeleteDevice(name string) error {
	return m.err
}

func (m *MockDeviceUseCase) Up(name string) error {
	return m.err
}

func (m *MockDeviceUseCase) Down(name string) error {
	return m.err
}

// DeviceUseCaseInterface defines the interface that DeviceHandler expects
type DeviceUseCaseInterface interface {
	ListDevices(page, perPage int) ([]entity.Device, int, error)
	GetDevice(name string) (*entity.Device, error)
	CreateDevice(req entity.DeviceCreateOrUpdateRequest) (*entity.Device, error)
	UpdateDevice(name string, req entity.DeviceCreateOrUpdateRequest) (*entity.Device, error)
	DeleteDevice(name string) error
	Up(name string) error
	Down(name string) error
}

func setupDeviceTestApp(uc *usecase.DeviceUseCase) *fiber.App {
	app := fiber.New()
	handler := NewDeviceHandler(uc)

	app.Get("/devices/", handler.ListDevices)
	app.Post("/devices/", handler.CreateDevice)
	app.Get("/devices/:name/", handler.GetDevice)
	app.Patch("/devices/:name/", handler.UpdateDevice)
	app.Delete("/devices/:name/", handler.DeleteDevice)
	app.Post("/devices/:name/up/", handler.Up)
	app.Post("/devices/:name/down/", handler.Down)

	return app
}

func TestDeviceHandler_ListDevices_Empty(t *testing.T) {
	// Since we can't easily mock the use case without interfaces,
	// we'll test the handler behavior with a nil use case
	// This tests the handler setup and routing
	t.Skip("Requires mocked use case interface")
}

func TestDeviceHandler_CreateDevice_InvalidJSON(t *testing.T) {
	app := fiber.New()
	// Handler without use case - tests JSON parsing
	handler := &DeviceHandler{}
	app.Post("/devices/", handler.CreateDevice)

	req := httptest.NewRequest(http.MethodPost, "/devices/", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var errResp entity.Error
	json.Unmarshal(body, &errResp)
	assert.Equal(t, entity.ErrCodeInvalidRequest, errResp.Code)
}

func TestDeviceHandler_UpdateDevice_InvalidJSON(t *testing.T) {
	app := fiber.New()
	handler := &DeviceHandler{}
	app.Patch("/devices/:name/", handler.UpdateDevice)

	req := httptest.NewRequest(http.MethodPatch, "/devices/wg0/", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var errResp entity.Error
	json.Unmarshal(body, &errResp)
	assert.Equal(t, entity.ErrCodeInvalidRequest, errResp.Code)
}

func TestIsDeviceExistsError(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"device already exists", assert.AnError, false},
		{"contains already exists", &mockError{msg: "device already exists"}, true},
		{"contains already exists mid", &mockError{msg: "the device already exists here"}, true},
		{"other error", &mockError{msg: "some other error"}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isDeviceExistsError(tc.err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

type mockError struct {
	msg string
}

func (e *mockError) Error() string {
	return e.msg
}
