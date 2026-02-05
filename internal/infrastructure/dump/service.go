package dump

import (
	"context"
	"log"
	"time"

	"github.com/suquant/wgrest/internal/infrastructure/wgquick"
	"github.com/suquant/wgrest/internal/infrastructure/wireguard"
)

// Service provides periodic config dump functionality.
type Service struct {
	interval   time.Duration
	wgClient   *wireguard.Client
	wgquickSvc *wgquick.Service
}

// NewService creates a new dump service.
func NewService(
	interval time.Duration,
	wgClient *wireguard.Client,
	wgquickSvc *wgquick.Service,
) *Service {
	return &Service{
		interval:   interval,
		wgClient:   wgClient,
		wgquickSvc: wgquickSvc,
	}
}

// Start begins the periodic dump loop.
func (s *Service) Start(ctx context.Context) {
	// Do an initial dump
	if err := s.SaveAll(); err != nil {
		log.Printf("Initial config dump failed: %v", err)
	}

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Final save on shutdown
			log.Println("Performing final config dump before shutdown...")
			if err := s.SaveAll(); err != nil {
				log.Printf("Final config dump failed: %v", err)
			}
			return
		case <-ticker.C:
			if err := s.SaveAll(); err != nil {
				log.Printf("Periodic config dump failed: %v", err)
			}
		}
	}
}

// SaveAll saves configs for all WireGuard devices using wg showconf.
func (s *Service) SaveAll() error {
	devices, err := s.wgClient.List()
	if err != nil {
		return err
	}

	var lastErr error
	for _, device := range devices {
		if err := s.wgquickSvc.SaveFromShowconf(device.Name); err != nil {
			log.Printf("Failed to save config for %s: %v", device.Name, err)
			lastErr = err
		} else {
			log.Printf("Saved config for interface: %s", device.Name)
		}
	}

	return lastErr
}
