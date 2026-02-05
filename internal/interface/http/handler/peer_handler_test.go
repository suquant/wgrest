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
)

func TestPeerHandler_CreatePeer_InvalidJSON(t *testing.T) {
	app := fiber.New()
	handler := &PeerHandler{}
	app.Post("/devices/:name/peers/", handler.CreatePeer)

	req := httptest.NewRequest(http.MethodPost, "/devices/wg0/peers/", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var errResp entity.Error
	json.Unmarshal(body, &errResp)
	assert.Equal(t, entity.ErrCodeInvalidRequest, errResp.Code)
}

func TestPeerHandler_UpdatePeer_InvalidJSON(t *testing.T) {
	app := fiber.New()
	handler := &PeerHandler{}
	app.Patch("/devices/:name/peers/:urlSafePubKey/", handler.UpdatePeer)

	req := httptest.NewRequest(http.MethodPatch, "/devices/wg0/peers/abc123/", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var errResp entity.Error
	json.Unmarshal(body, &errResp)
	assert.Equal(t, entity.ErrCodeInvalidRequest, errResp.Code)
}

func TestPeerHandler_ListPeers_QueryParams(t *testing.T) {
	// Test that query parameters are correctly parsed
	app := fiber.New()
	
	var capturedPage, capturedPerPage int
	var capturedQuery, capturedSort string

	app.Get("/devices/:name/peers/", func(c *fiber.Ctx) error {
		capturedPage = c.QueryInt("page", 0)
		capturedPerPage = c.QueryInt("per_page", 100)
		capturedQuery = c.Query("q")
		capturedSort = c.Query("sort")
		return c.JSON([]entity.Peer{})
	})

	req := httptest.NewRequest(http.MethodGet, "/devices/wg0/peers/?page=2&per_page=50&q=test&sort=-receive_bytes", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Equal(t, 2, capturedPage)
	assert.Equal(t, 50, capturedPerPage)
	assert.Equal(t, "test", capturedQuery)
	assert.Equal(t, "-receive_bytes", capturedSort)
}

func TestPeerHandler_GetQuickConfig_RouteParams(t *testing.T) {
	// Test that route parameters are correctly extracted
	app := fiber.New()

	var capturedDeviceName, capturedPubKey string

	app.Get("/devices/:name/peers/:urlSafePubKey/quick.conf", func(c *fiber.Ctx) error {
		capturedDeviceName = c.Params("name")
		capturedPubKey = c.Params("urlSafePubKey")
		return c.SendString("[Interface]\n# Config here")
	})

	req := httptest.NewRequest(http.MethodGet, "/devices/wg0/peers/dGVzdFB1YmxpY0tleQ==/quick.conf", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Equal(t, "wg0", capturedDeviceName)
	assert.Equal(t, "dGVzdFB1YmxpY0tleQ==", capturedPubKey)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "[Interface]")
}
