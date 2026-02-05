package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBearerAuth_ValidToken(t *testing.T) {
	app := fiber.New()
	app.Use(BearerAuth("test-token"))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer test-token")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "OK", string(body))
}

func TestBearerAuth_InvalidToken(t *testing.T) {
	app := fiber.New()
	app.Use(BearerAuth("test-token"))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer wrong-token")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestBearerAuth_MissingToken(t *testing.T) {
	app := fiber.New()
	app.Use(BearerAuth("test-token"))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestBearerAuth_MalformedHeader(t *testing.T) {
	app := fiber.New()
	app.Use(BearerAuth("test-token"))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	testCases := []struct {
		name   string
		header string
	}{
		{"Basic auth", "Basic dXNlcjpwYXNz"},
		{"No Bearer prefix", "test-token"},
		{"Empty", ""},
		{"Bearer only", "Bearer"},
		{"Bearer with space only", "Bearer "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	}
}

func TestBearerAuth_EmptyTokenConfig(t *testing.T) {
	// When token is empty string, middleware still checks for Bearer format
	// but any token will fail since "" != providedToken
	app := fiber.New()
	app.Use(BearerAuth(""))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	// Empty token config still requires auth header
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
