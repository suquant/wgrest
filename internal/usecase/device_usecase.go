package usecase

import (
	"github.com/suquant/wgrest/internal/domain/entity"
	"github.com/suquant/wgrest/internal/infrastructure/wgquick"
	"github.com/suquant/wgrest/internal/infrastructure/wireguard"
)

// DeviceUseCase handles business logic for device operations.
type DeviceUseCase struct {
	wgClient      *wireguard.Client
	wgquickSvc    *wgquick.Service
}

// NewDeviceUseCase creates a new device use case.
func NewDeviceUseCase(
	wgClient *wireguard.Client,
	wgquickSvc *wgquick.Service,
) *DeviceUseCase {
	return &DeviceUseCase{
		wgClient:   wgClient,
		wgquickSvc: wgquickSvc,
	}
}

// ListDevices returns all devices (running + config-only) with pagination.
func (uc *DeviceUseCase) ListDevices(page, perPage int) ([]entity.Device, int, error) {
	// Get running devices
	runningDevices, err := uc.wgClient.List()
	if err != nil {
		return nil, 0, err
	}

	// Track running device names
	runningNames := make(map[string]bool)
	for _, d := range runningDevices {
		runningNames[d.Name] = true
	}

	// Enrich running devices with wg-quick config
	for i := range runningDevices {
		runningDevices[i].Running = true
		uc.enrichDeviceWithConfig(&runningDevices[i])
	}

	// Get devices from config files that aren't running
	configDeviceNames := uc.wgquickSvc.ListConfigDevices()
	var configOnlyDevices []entity.Device
	for _, name := range configDeviceNames {
		if !runningNames[name] {
			device := entity.Device{
				Name:    name,
				Running: false,
			}
			uc.enrichDeviceWithConfig(&device)
			configOnlyDevices = append(configOnlyDevices, device)
		}
	}

	// Merge: running first, then config-only
	devices := append(runningDevices, configOnlyDevices...)
	total := len(devices)

	// Apply pagination
	if perPage <= 0 {
		perPage = 100
	}
	if page < 0 {
		page = 0
	}

	start := page * perPage
	if start >= len(devices) {
		return []entity.Device{}, total, nil
	}

	end := start + perPage
	if end > len(devices) {
		end = len(devices)
	}

	return devices[start:end], total, nil
}

// GetDevice returns a device by name (running or config-only).
func (uc *DeviceUseCase) GetDevice(name string) (*entity.Device, error) {
	// Try to get running device first
	device, err := uc.wgClient.Get(name)
	if err == nil {
		device.Running = true
		uc.enrichDeviceWithConfig(device)
		return device, nil
	}

	// Not running - check if config exists
	cfg, cfgErr := uc.wgquickSvc.LoadConfig(name)
	if cfgErr != nil {
		// Neither running nor config exists
		return nil, err
	}

	// Build device from config
	device = &entity.Device{
		Name:       name,
		Running:    false,
		PrivateKey: cfg.PrivateKey,
		ListenPort: int32(cfg.ListenPort),
		Addresses:  cfg.Addresses,
		DNS:        cfg.DNS,
		MTU:        int32(cfg.MTU),
		Table:      cfg.Table,
		PreUp:      cfg.PreUp,
		PostUp:     cfg.PostUp,
		PreDown:    cfg.PreDown,
		PostDown:   cfg.PostDown,
	}

	return device, nil
}

// CreateDevice creates a new device and writes wg-quick config.
func (uc *DeviceUseCase) CreateDevice(req entity.DeviceCreateOrUpdateRequest) (*entity.Device, error) {
	device, err := uc.wgClient.Create(req)
	if err != nil {
		return nil, err
	}

	// Apply wg-quick options
	device.Addresses = req.Addresses
	device.DNS = req.DNS
	if req.MTU != nil {
		device.MTU = *req.MTU
	}
	if req.Table != nil {
		device.Table = *req.Table
	}
	device.PreUp = req.PreUp
	device.PostUp = req.PostUp
	device.PreDown = req.PreDown
	device.PostDown = req.PostDown

	// Save wg-quick config
	peers, _ := uc.wgClient.ListPeers(device.Name)
	if err := uc.wgquickSvc.SaveConfig(device, peers); err != nil {
		// Log but don't fail - device was created
	}

	return device, nil
}

// UpdateDevice updates a device and writes wg-quick config.
func (uc *DeviceUseCase) UpdateDevice(name string, req entity.DeviceCreateOrUpdateRequest) (*entity.Device, error) {
	device, err := uc.wgClient.Update(name, req)
	if err != nil {
		return nil, err
	}

	// Apply wg-quick options
	if len(req.Addresses) > 0 {
		device.Addresses = req.Addresses
	}
	if len(req.DNS) > 0 {
		device.DNS = req.DNS
	}
	if req.MTU != nil {
		device.MTU = *req.MTU
	}
	if req.Table != nil {
		device.Table = *req.Table
	}
	if len(req.PreUp) > 0 {
		device.PreUp = req.PreUp
	}
	if len(req.PostUp) > 0 {
		device.PostUp = req.PostUp
	}
	if len(req.PreDown) > 0 {
		device.PreDown = req.PreDown
	}
	if len(req.PostDown) > 0 {
		device.PostDown = req.PostDown
	}

	// Enrich with existing config values
	uc.enrichDeviceWithConfig(device)

	// Save wg-quick config
	peers, _ := uc.wgClient.ListPeers(device.Name)
	if err := uc.wgquickSvc.SaveConfig(device, peers); err != nil {
		// Log but don't fail
	}

	return device, nil
}

// DeleteDevice deletes a device.
func (uc *DeviceUseCase) DeleteDevice(name string) error {
	return uc.wgClient.Delete(name)
}

// Up brings up a WireGuard interface using wg-quick.
func (uc *DeviceUseCase) Up(name string) error {
	return uc.wgquickSvc.Up(name)
}

// Down brings down a WireGuard interface using wg-quick.
func (uc *DeviceUseCase) Down(name string) error {
	return uc.wgquickSvc.Down(name)
}

func (uc *DeviceUseCase) enrichDeviceWithConfig(device *entity.Device) {
	cfg, err := uc.wgquickSvc.LoadConfig(device.Name)
	if err != nil {
		return
	}

	device.Addresses = cfg.Addresses
	device.DNS = cfg.DNS
	device.MTU = int32(cfg.MTU)
	device.Table = cfg.Table
	device.PreUp = cfg.PreUp
	device.PostUp = cfg.PostUp
	device.PreDown = cfg.PreDown
	device.PostDown = cfg.PostDown
}
