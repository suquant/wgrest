package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"golang.org/x/crypto/acme/autocert"

	"github.com/suquant/wgrest/api/docs"
	"github.com/suquant/wgrest/internal/infrastructure/dump"
	"github.com/suquant/wgrest/internal/infrastructure/wgquick"
	"github.com/suquant/wgrest/internal/infrastructure/wireguard"
	httpInterface "github.com/suquant/wgrest/internal/interface/http"
	"github.com/suquant/wgrest/internal/interface/http/handler"
	"github.com/suquant/wgrest/internal/usecase"
)

var (
	appVersion string // Populated during build time
)

// defaultConfigDirs returns the platform-specific default WireGuard config directories.
// Matches wg-quick search order.
func defaultConfigDirs() []string {
	switch runtime.GOOS {
	case "darwin", "freebsd":
		return []string{
			"/etc/wireguard",
			"/usr/local/etc/wireguard",
			"/opt/homebrew/etc/wireguard",
		}
	default:
		return []string{"/etc/wireguard"}
	}
}

// @title WGRest API
// @version 1.0
// @description REST API for managing WireGuard interfaces and peers
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8000
// @basePath /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @schemes http https
func main() {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "conf",
			Value:   "/etc/wgrest/wgrest.conf",
			Usage:   "wgrest config file path",
			EnvVars: []string{"WGREST_CONF"},
		},
		&cli.BoolFlag{
			Name:  "version",
			Value: false,
			Usage: "Print version and exit",
		},
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "listen",
			Value:   "127.0.0.1:8000",
			Usage:   "Listen address",
			EnvVars: []string{"WGREST_LISTEN"},
		}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name:    "config-dir",
			Value:   cli.NewStringSlice(defaultConfigDirs()...),
			Usage:   "WireGuard config directories (wg-quick style, can specify multiple)",
			EnvVars: []string{"WGREST_CONFIG_DIR"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "certs-dir",
			Value:   "/var/lib/wgrest/certs",
			Usage:   "ACME TLS certificates cache directory",
			EnvVars: []string{"WGREST_CERTS_DIR"},
		}),
		altsrc.NewDurationFlag(&cli.DurationFlag{
			Name:    "dump-interval",
			Value:   10 * time.Minute,
			Usage:   "Config dump interval",
			EnvVars: []string{"WGREST_DUMP_INTERVAL"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "static-auth-token",
			Value:   "",
			Usage:   "Bearer token for authorization",
			EnvVars: []string{"WGREST_STATIC_AUTH_TOKEN"},
		}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name:    "tls-domain",
			Value:   cli.NewStringSlice(),
			Usage:   "TLS Domains for ACME (Let's Encrypt)",
			EnvVars: []string{"WGREST_TLS_DOMAIN"},
		}),
	}

	app := &cli.App{
		Name:   "wgrest",
		Usage:  "wgrest - REST API for WireGuard",
		Flags:  flags,
		Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewTomlSourceFromFlagFunc("conf")),
		Action: func(c *cli.Context) error {
			if c.Bool("version") {
				fmt.Printf("wgrest version: %s\n", appVersion)
				return nil
			}

			// Create context for graceful shutdown
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Initialize Fiber
			fiberApp := fiber.New(fiber.Config{
				DisableStartupMessage: false,
				AppName:               "wgrest",
			})

			// Initialize WireGuard client
			wgClient, err := wireguard.NewClient()
			if err != nil {
				return fmt.Errorf("failed to create wireguard client: %w", err)
			}
			defer wgClient.Close()

			// Initialize wg-quick config service
			configDirs := c.StringSlice("config-dir")
			wgquickSvc, err := wgquick.NewService(configDirs)
			if err != nil {
				return fmt.Errorf("failed to create wgquick service: %w", err)
			}

			// Initialize dump service
			dumpInterval := c.Duration("dump-interval")
			dumpService := dump.NewService(dumpInterval, wgClient, wgquickSvc)

			// Start dump service in background
			go dumpService.Start(ctx)
			log.Printf("Config dump service started (interval: %s, dirs: %v)", dumpInterval, configDirs)

			// Initialize use cases
			deviceUC := usecase.NewDeviceUseCase(wgClient, wgquickSvc)
			peerUC := usecase.NewPeerUseCase(wgClient, wgquickSvc)

			// Initialize handlers
			deviceHandler := handler.NewDeviceHandler(deviceUC)
			peerHandler := handler.NewPeerHandler(peerUC)

			// Setup routes
			httpInterface.SetupRouter(fiberApp, httpInterface.RouterConfig{
				DeviceHandler: deviceHandler,
				PeerHandler:   peerHandler,
				AuthToken:     c.String("static-auth-token"),
				Version:       appVersion,
				OpenAPISpec:   docs.OpenAPISpec,
			})

			// Handle graceful shutdown
			go func() {
				sigCh := make(chan os.Signal, 1)
				signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
				<-sigCh

				log.Println("Received shutdown signal...")
				cancel() // Triggers final config dump

				time.Sleep(500 * time.Millisecond)

				if err := fiberApp.Shutdown(); err != nil {
					log.Printf("Error shutting down server: %v", err)
				}
			}()

			// Start server
			listen := c.String("listen")
			tlsDomains := c.StringSlice("tls-domain")

			if len(tlsDomains) > 0 {
				// ACME TLS with Let's Encrypt
				certsDir := c.String("certs-dir")
				certManager := &autocert.Manager{
					Prompt:     autocert.AcceptTOS,
					HostPolicy: autocert.HostWhitelist(tlsDomains...),
					Cache:      autocert.DirCache(certsDir),
				}

				// Create TLS listener with autocert
				tlsListener := certManager.Listener()

				log.Printf("Starting wgrest server with ACME TLS for domains: %v", tlsDomains)
				return fiberApp.Listener(tlsListener)
			}

			log.Printf("Starting wgrest server on %s", listen)
			return fiberApp.Listen(listen)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
