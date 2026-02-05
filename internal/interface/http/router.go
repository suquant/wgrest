package http

import (
	"net/http"
	"os/exec"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	"github.com/suquant/wgrest/internal/interface/http/handler"
	"github.com/suquant/wgrest/internal/interface/http/middleware"
)

// RouterConfig contains configuration for the router.
type RouterConfig struct {
	DeviceHandler *handler.DeviceHandler
	PeerHandler   *handler.PeerHandler
	AuthToken     string
	Version       string
	OpenAPISpec   []byte
}

// SetupRouter configures all routes for the Fiber application.
func SetupRouter(app *fiber.App, cfg RouterConfig) {
	// Global middleware
	app.Use(logger.New(logger.Config{
		Format: "${time} ${status} - ${method} ${path}\n",
	}))
	app.Use(recover.New())

	// Version endpoint (no auth required)
	app.Get("/version", func(c *fiber.Ctx) error {
		wgVersion := getWireGuardVersion()
		return c.JSON(fiber.Map{
			"wgrest":    cfg.Version,
			"wireguard": wgVersion,
		})
	})

	// OpenAPI spec endpoint (no auth required)
	app.Get("/openapi.json", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Send(cfg.OpenAPISpec)
	})

	// Swagger UI (no auth required)
	app.Get("/docs/*", swagger.New(swagger.Config{
		URL:         "/openapi.json",
		Title:       "WGRest API Documentation",
		DeepLinking: true,
	}))

	// Rewrite middleware for backward compatibility
	// Redirect /devices to /v1/devices/
	app.Use(func(c *fiber.Ctx) error {
		path := c.Path()
		if path == "/devices" || path == "/devices/" {
			return c.Redirect("/v1/devices/", http.StatusTemporaryRedirect)
		}
		return c.Next()
	})

	// API v1 group
	v1 := app.Group("/v1")

	// CORS middleware for v1
	v1.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,HEAD,PUT,PATCH,POST,DELETE",
		AllowHeaders: "Content-Type,Accept,Accept-Language,Link,Authorization",
	}))

	// Auth middleware if token configured
	if cfg.AuthToken != "" {
		v1.Use(middleware.BearerAuth(cfg.AuthToken))
	}

	// Device routes
	v1.Get("/devices/", cfg.DeviceHandler.ListDevices)
	v1.Post("/devices/", cfg.DeviceHandler.CreateDevice)
	v1.Get("/devices/:name/", cfg.DeviceHandler.GetDevice)
	v1.Patch("/devices/:name/", cfg.DeviceHandler.UpdateDevice)
	v1.Delete("/devices/:name/", cfg.DeviceHandler.DeleteDevice)

	// wg-quick operations
	v1.Post("/devices/:name/up/", cfg.DeviceHandler.Up)
	v1.Post("/devices/:name/down/", cfg.DeviceHandler.Down)

	// Peer routes
	v1.Get("/devices/:name/peers/", cfg.PeerHandler.ListPeers)
	v1.Post("/devices/:name/peers/", cfg.PeerHandler.CreatePeer)
	v1.Get("/devices/:name/peers/:urlSafePubKey/", cfg.PeerHandler.GetPeer)
	v1.Patch("/devices/:name/peers/:urlSafePubKey/", cfg.PeerHandler.UpdatePeer)
	v1.Delete("/devices/:name/peers/:urlSafePubKey/", cfg.PeerHandler.DeletePeer)
}

// getWireGuardVersion executes wg --version and returns the version string.
func getWireGuardVersion() string {
	out, err := exec.Command("wg", "--version").Output()
	if err != nil {
		return "unknown"
	}
	// Output format: "wireguard-tools v1.0.20210914 - https://..."
	version := strings.TrimSpace(string(out))
	// Extract just the version part
	parts := strings.Fields(version)
	if len(parts) >= 2 {
		return parts[1] // e.g., "v1.0.20210914"
	}
	return version
}
