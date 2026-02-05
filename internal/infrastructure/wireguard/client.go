package wireguard

import (
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"github.com/suquant/wgrest/internal/domain/entity"
)

// EmptyKey represents an empty WireGuard key (all zeros).
var EmptyKey = wgtypes.Key{}

// Client provides access to WireGuard devices via wgctrl.
type Client struct {
	ctrl *wgctrl.Client
}

// NewClient creates a new WireGuard client.
func NewClient() (*Client, error) {
	ctrl, err := wgctrl.New()
	if err != nil {
		return nil, err
	}
	return &Client{ctrl: ctrl}, nil
}

// Close closes the WireGuard client.
func (c *Client) Close() error {
	return c.ctrl.Close()
}

// resolveInterfaceName resolves the actual interface name for a logical name.
// On macOS/BSD, wg-quick maps wg0 -> utunX and stores mapping in /var/run/wireguard/wg0.name
func resolveInterfaceName(name string) string {
	if runtime.GOOS != "darwin" && runtime.GOOS != "freebsd" {
		return name
	}

	// Check if there's a name mapping file
	nameFile := filepath.Join("/var/run/wireguard", name+".name")
	data, err := os.ReadFile(nameFile)
	if err != nil {
		return name // No mapping, use original name
	}

	realName := strings.TrimSpace(string(data))
	if realName == "" {
		return name
	}

	return realName
}

// List returns all WireGuard devices.
func (c *Client) List() ([]entity.Device, error) {
	devices, err := c.ctrl.Devices()
	if err != nil {
		return nil, err
	}

	result := make([]entity.Device, len(devices))
	for i, d := range devices {
		result[i] = deviceToEntity(d)
	}

	return result, nil
}

// Get returns a specific device by name.
func (c *Client) Get(name string) (*entity.Device, error) {
	// Resolve actual interface name (macOS uses utunX)
	realName := resolveInterfaceName(name)

	d, err := c.ctrl.Device(realName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("device %s not found", name)
		}
		return nil, err
	}

	device := deviceToEntity(d)
	// Preserve logical name (wg0) instead of real name (utun7)
	device.Name = name
	return &device, nil
}

// Create creates a new WireGuard device.
func (c *Client) Create(req entity.DeviceCreateOrUpdateRequest) (*entity.Device, error) {
	if req.Name == nil || *req.Name == "" {
		return nil, fmt.Errorf("device name is required")
	}

	name := *req.Name

	// Check if device already exists
	_, err := c.ctrl.Device(name)
	if err == nil {
		return nil, fmt.Errorf("device %s already exists", name)
	}

	// Build configuration
	cfg := wgtypes.Config{}

	if req.PrivateKey != nil {
		key, err := wgtypes.ParseKey(*req.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("invalid private key: %w", err)
		}
		cfg.PrivateKey = &key
	}

	if req.ListenPort != nil {
		port := int(*req.ListenPort)
		cfg.ListenPort = &port
	}

	if req.FirewallMark != nil {
		mark := int(*req.FirewallMark)
		cfg.FirewallMark = &mark
	}

	if err := c.ctrl.ConfigureDevice(name, cfg); err != nil {
		return nil, err
	}

	return c.Get(name)
}

// Update updates an existing WireGuard device.
func (c *Client) Update(name string, req entity.DeviceCreateOrUpdateRequest) (*entity.Device, error) {
	// Check if device exists
	_, err := c.ctrl.Device(name)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("device %s not found", name)
		}
		return nil, err
	}

	cfg := wgtypes.Config{}

	if req.PrivateKey != nil {
		key, err := wgtypes.ParseKey(*req.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("invalid private key: %w", err)
		}
		cfg.PrivateKey = &key
	}

	if req.ListenPort != nil {
		port := int(*req.ListenPort)
		cfg.ListenPort = &port
	}

	if req.FirewallMark != nil {
		mark := int(*req.FirewallMark)
		cfg.FirewallMark = &mark
	}

	if err := c.ctrl.ConfigureDevice(name, cfg); err != nil {
		return nil, err
	}

	return c.Get(name)
}

// Delete removes a WireGuard device.
func (c *Client) Delete(name string) error {
	// Note: wgctrl doesn't support device deletion directly
	// This would typically be handled by removing the network interface
	return fmt.Errorf("device deletion not supported via wgctrl")
}

// ListPeers returns all peers for a device.
func (c *Client) ListPeers(deviceName string) ([]entity.Peer, error) {
	// Resolve actual interface name (macOS uses utunX)
	realName := resolveInterfaceName(deviceName)

	d, err := c.ctrl.Device(realName)
	if err != nil {
		// Handle various "not found" errors across platforms
		if os.IsNotExist(err) || strings.Contains(err.Error(), "does not exist") ||
			strings.Contains(err.Error(), "no such file") ||
			strings.Contains(err.Error(), "not found") {
			return nil, fmt.Errorf("device %s not found", deviceName)
		}
		// Return actual error for debugging
		return nil, fmt.Errorf("failed to get device %s: %w", deviceName, err)
	}

	result := make([]entity.Peer, len(d.Peers))
	for i, p := range d.Peers {
		result[i] = peerToEntity(p)
	}

	return result, nil
}

// GetPeer returns a specific peer by URL-safe public key.
func (c *Client) GetPeer(deviceName string, urlSafePubKey string) (*entity.Peer, error) {
	pubKey, err := decodeURLSafeKey(urlSafePubKey)
	if err != nil {
		return nil, err
	}

	// Resolve actual interface name (macOS uses utunX)
	realName := resolveInterfaceName(deviceName)

	d, err := c.ctrl.Device(realName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("device %s not found", deviceName)
		}
		return nil, err
	}

	for _, p := range d.Peers {
		if p.PublicKey.String() == pubKey.String() {
			peer := peerToEntity(p)
			return &peer, nil
		}
	}

	return nil, fmt.Errorf("peer not found")
}

// CreatePeer creates a new peer for a device.
func (c *Client) CreatePeer(deviceName string, req entity.PeerCreateOrUpdateRequest) (*entity.Peer, error) {
	d, err := c.ctrl.Device(deviceName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("device %s not found", deviceName)
		}
		return nil, err
	}

	peerCfg := wgtypes.PeerConfig{}

	// Generate or use provided keys
	if req.PublicKey != nil {
		key, err := wgtypes.ParseKey(*req.PublicKey)
		if err != nil {
			return nil, fmt.Errorf("invalid public key: %w", err)
		}
		peerCfg.PublicKey = key
	} else {
		// Generate new key pair
		privateKey, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			return nil, err
		}
		peerCfg.PublicKey = privateKey.PublicKey()
		// Store private key in request for caller to persist
		pk := privateKey.String()
		req.PrivateKey = &pk
	}

	if req.PresharedKey != nil {
		key, err := wgtypes.ParseKey(*req.PresharedKey)
		if err != nil {
			return nil, fmt.Errorf("invalid preshared key: %w", err)
		}
		peerCfg.PresharedKey = &key
	}

	if len(req.AllowedIPs) > 0 {
		allowedIPs := make([]net.IPNet, 0, len(req.AllowedIPs))
		for _, ip := range req.AllowedIPs {
			_, ipNet, err := net.ParseCIDR(ip)
			if err != nil {
				return nil, fmt.Errorf("invalid allowed IP %s: %w", ip, err)
			}
			allowedIPs = append(allowedIPs, *ipNet)
		}
		peerCfg.AllowedIPs = allowedIPs
	}

	if req.PersistentKeepaliveInterval != nil {
		duration, err := time.ParseDuration(*req.PersistentKeepaliveInterval)
		if err != nil {
			return nil, fmt.Errorf("invalid keepalive interval: %w", err)
		}
		peerCfg.PersistentKeepaliveInterval = &duration
	}

	if req.Endpoint != nil && *req.Endpoint != "" {
		addr, err := net.ResolveUDPAddr("udp", *req.Endpoint)
		if err != nil {
			return nil, fmt.Errorf("invalid endpoint: %w", err)
		}
		peerCfg.Endpoint = addr
	}

	cfg := wgtypes.Config{
		Peers: []wgtypes.PeerConfig{peerCfg},
	}

	if err := c.ctrl.ConfigureDevice(d.Name, cfg); err != nil {
		return nil, err
	}

	peer := entity.Peer{
		PublicKey:                   peerCfg.PublicKey.String(),
		URLSafePublicKey:            base64.URLEncoding.EncodeToString(peerCfg.PublicKey[:]),
		AllowedIPs:                  req.AllowedIPs,
		PersistentKeepaliveInterval: "0s",
		LastHandshakeTime:           time.Time{},
	}

	if req.PrivateKey != nil {
		peer.PrivateKey = *req.PrivateKey
	}

	if req.PresharedKey != nil {
		peer.PresharedKey = *req.PresharedKey
	}

	if req.PersistentKeepaliveInterval != nil {
		peer.PersistentKeepaliveInterval = *req.PersistentKeepaliveInterval
	}

	if req.Endpoint != nil {
		peer.Endpoint = *req.Endpoint
	}

	return &peer, nil
}

// UpdatePeer updates an existing peer.
func (c *Client) UpdatePeer(deviceName string, urlSafePubKey string, req entity.PeerCreateOrUpdateRequest) (*entity.Peer, error) {
	pubKey, err := decodeURLSafeKey(urlSafePubKey)
	if err != nil {
		return nil, err
	}

	d, err := c.ctrl.Device(deviceName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("device %s not found", deviceName)
		}
		return nil, err
	}

	// Find existing peer
	var existingPeer *wgtypes.Peer
	for _, p := range d.Peers {
		if p.PublicKey.String() == pubKey.String() {
			existingPeer = &p
			break
		}
	}

	if existingPeer == nil {
		return nil, fmt.Errorf("peer not found")
	}

	peerCfg := wgtypes.PeerConfig{
		PublicKey:         *pubKey,
		UpdateOnly:        true,
		ReplaceAllowedIPs: len(req.AllowedIPs) > 0,
	}

	if req.PresharedKey != nil {
		key, err := wgtypes.ParseKey(*req.PresharedKey)
		if err != nil {
			return nil, fmt.Errorf("invalid preshared key: %w", err)
		}
		peerCfg.PresharedKey = &key
	}

	if len(req.AllowedIPs) > 0 {
		allowedIPs := make([]net.IPNet, 0, len(req.AllowedIPs))
		for _, ip := range req.AllowedIPs {
			_, ipNet, err := net.ParseCIDR(ip)
			if err != nil {
				return nil, fmt.Errorf("invalid allowed IP %s: %w", ip, err)
			}
			allowedIPs = append(allowedIPs, *ipNet)
		}
		peerCfg.AllowedIPs = allowedIPs
	}

	if req.PersistentKeepaliveInterval != nil {
		duration, err := time.ParseDuration(*req.PersistentKeepaliveInterval)
		if err != nil {
			return nil, fmt.Errorf("invalid keepalive interval: %w", err)
		}
		peerCfg.PersistentKeepaliveInterval = &duration
	}

	if req.Endpoint != nil && *req.Endpoint != "" {
		addr, err := net.ResolveUDPAddr("udp", *req.Endpoint)
		if err != nil {
			return nil, fmt.Errorf("invalid endpoint: %w", err)
		}
		peerCfg.Endpoint = addr
	}

	cfg := wgtypes.Config{
		Peers: []wgtypes.PeerConfig{peerCfg},
	}

	if err := c.ctrl.ConfigureDevice(d.Name, cfg); err != nil {
		return nil, err
	}

	return c.GetPeer(deviceName, urlSafePubKey)
}

// DeletePeer removes a peer from a device.
func (c *Client) DeletePeer(deviceName string, urlSafePubKey string) (*entity.Peer, error) {
	pubKey, err := decodeURLSafeKey(urlSafePubKey)
	if err != nil {
		return nil, err
	}

	// Get peer info before deletion
	peer, err := c.GetPeer(deviceName, urlSafePubKey)
	if err != nil {
		return nil, err
	}

	peerCfg := wgtypes.PeerConfig{
		PublicKey: *pubKey,
		Remove:    true,
	}

	cfg := wgtypes.Config{
		Peers: []wgtypes.PeerConfig{peerCfg},
	}

	if err := c.ctrl.ConfigureDevice(deviceName, cfg); err != nil {
		return nil, err
	}

	return peer, nil
}

// GetDevice returns the raw wgtypes.Device for advanced operations.
func (c *Client) GetDevice(name string) (*wgtypes.Device, error) {
	return c.ctrl.Device(name)
}

func deviceToEntity(d *wgtypes.Device) entity.Device {
	var totalReceive, totalTransmit int64
	for _, p := range d.Peers {
		totalReceive += p.ReceiveBytes
		totalTransmit += p.TransmitBytes
	}

	networks := getDeviceNetworks(d.Name)

	return entity.Device{
		Name:               d.Name,
		ListenPort:         int32(d.ListenPort),
		PublicKey:          d.PublicKey.String(),
		PrivateKey:         d.PrivateKey.String(),
		FirewallMark:       int32(d.FirewallMark),
		Addresses:          networks,
		PeersCount:         int32(len(d.Peers)),
		TotalReceiveBytes:  totalReceive,
		TotalTransmitBytes: totalTransmit,
	}
}

func peerToEntity(p wgtypes.Peer) entity.Peer {
	allowedIPs := make([]string, len(p.AllowedIPs))
	for i, ip := range p.AllowedIPs {
		allowedIPs[i] = ip.String()
	}

	endpoint := ""
	if p.Endpoint != nil {
		endpoint = p.Endpoint.String()
	}

	presharedKey := ""
	if p.PresharedKey != EmptyKey {
		presharedKey = p.PresharedKey.String()
	}

	return entity.Peer{
		PublicKey:                   p.PublicKey.String(),
		URLSafePublicKey:            base64.URLEncoding.EncodeToString(p.PublicKey[:]),
		PresharedKey:                presharedKey,
		AllowedIPs:                  allowedIPs,
		LastHandshakeTime:           p.LastHandshakeTime,
		PersistentKeepaliveInterval: p.PersistentKeepaliveInterval.String(),
		Endpoint:                    endpoint,
		ReceiveBytes:                p.ReceiveBytes,
		TransmitBytes:               p.TransmitBytes,
	}
}

func decodeURLSafeKey(urlSafeKey string) (*wgtypes.Key, error) {
	// Handle both standard and URL-safe base64
	keyBytes, err := base64.URLEncoding.DecodeString(urlSafeKey)
	if err != nil {
		// Try standard base64 with URL-safe replacement
		keyBytes, err = base64.StdEncoding.DecodeString(
			strings.ReplaceAll(strings.ReplaceAll(urlSafeKey, "-", "+"), "_", "/"),
		)
		if err != nil {
			return nil, fmt.Errorf("invalid public key encoding: %w", err)
		}
	}

	if len(keyBytes) != wgtypes.KeyLen {
		return nil, fmt.Errorf("invalid key length")
	}

	var key wgtypes.Key
	copy(key[:], keyBytes)
	return &key, nil
}

func getDeviceNetworks(name string) []string {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil
	}

	networks := make([]string, len(addrs))
	for i, addr := range addrs {
		networks[i] = addr.String()
	}

	return networks
}
