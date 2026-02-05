package wgquick

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/suquant/wgrest/internal/domain/entity"
)

// Config represents a parsed wg-quick configuration.
type Config struct {
	// Interface section
	PrivateKey   string
	ListenPort   int
	FirewallMark int
	Addresses    []string
	DNS          []string
	MTU          int
	Table        string
	PreUp        []string
	PostUp       []string
	PreDown      []string
	PostDown     []string
	SaveConfig   bool

	// Peers section (raw WireGuard config)
	Peers []PeerConfig
}

// PeerConfig represents a peer in the config file.
type PeerConfig struct {
	PublicKey                   string
	PresharedKey                string
	AllowedIPs                  []string
	Endpoint                    string
	PersistentKeepaliveInterval int
}

// Service manages wg-quick configuration files.
type Service struct {
	// Config directories to search (in order)
	configDirs []string
	mu         sync.Mutex
}

// NewService creates a new wg-quick config service.
func NewService(configDirs []string) (*Service, error) {
	if len(configDirs) == 0 {
		return nil, fmt.Errorf("at least one config directory must be specified")
	}
	return &Service{configDirs: configDirs}, nil
}

// ensureDir creates a directory if it doesn't exist.
func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0700)
}

// FindConfigPath searches all directories and returns the first existing config file.
// If not found, returns path in the first directory (for new configs).
func (s *Service) FindConfigPath(name string) string {
	filename := name + ".conf"

	// Search all paths in order
	for _, dir := range s.configDirs {
		path := filepath.Join(dir, filename)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Not found - return path in first directory (for new configs)
	return filepath.Join(s.configDirs[0], filename)
}

// GetConfigDir returns the directory where a config file exists, or first dir for new configs.
func (s *Service) GetConfigDir(name string) string {
	filename := name + ".conf"

	// Search all paths in order
	for _, dir := range s.configDirs {
		path := filepath.Join(dir, filename)
		if _, err := os.Stat(path); err == nil {
			return dir
		}
	}

	// Not found - return first directory (for new configs)
	return s.configDirs[0]
}

// ConfigDirs returns the configured search directories.
func (s *Service) ConfigDirs() []string {
	return s.configDirs
}

// ListConfigDevices returns names of all devices that have config files.
func (s *Service) ListConfigDevices() []string {
	seen := make(map[string]bool)
	var devices []string

	for _, dir := range s.configDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue // Directory might not exist
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if strings.HasSuffix(name, ".conf") {
				deviceName := strings.TrimSuffix(name, ".conf")
				if !seen[deviceName] {
					seen[deviceName] = true
					devices = append(devices, deviceName)
				}
			}
		}
	}

	return devices
}

// SaveConfig writes a device configuration to file (in the same directory where found).
func (s *Service) SaveConfig(device *entity.Device, peers []entity.Peer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Write to the same dir where config was found, or first dir for new configs
	configDir := s.GetConfigDir(device.Name)
	if err := ensureDir(configDir); err != nil {
		return err
	}

	config := s.buildConfig(device, peers)
	configPath := filepath.Join(configDir, device.Name+".conf")
	tmpPath := configPath + ".tmp"

	if err := os.WriteFile(tmpPath, []byte(config), 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	if err := os.Rename(tmpPath, configPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to rename config: %w", err)
	}

	return nil
}

// SaveFromShowconf saves config using `wg showconf` command output.
func (s *Service) SaveFromShowconf(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cmd := exec.Command("wg", "showconf", name)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wg showconf failed: %w: %s", err, stderr.String())
	}

	// Write to the same dir where config was found, or first dir for new configs
	configDir := s.GetConfigDir(name)
	if err := ensureDir(configDir); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, name+".conf")
	tmpPath := configPath + ".tmp"

	if err := os.WriteFile(tmpPath, stdout.Bytes(), 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	if err := os.Rename(tmpPath, configPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to rename config: %w", err)
	}

	return nil
}

// LoadConfig parses a wg-quick config file, searching all paths.
func (s *Service) LoadConfig(name string) (*Config, error) {
	configPath := s.FindConfigPath(name)
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ParseConfig(f)
}

// Up brings up a WireGuard interface using wg-quick.
// Note: wg-quick typically requires root/sudo.
func (s *Service) Up(name string) error {
	// Use timeout to prevent indefinite waiting for sudo password
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "wg-quick", "up", name)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	// Prevent any stdin interaction (sudo password prompts) - use explicit devnull
	devNull, err := os.Open(os.DevNull)
	if err == nil {
		cmd.Stdin = devNull
		defer devNull.Close()
	}
	// Set environment to prevent sudo from asking for password
	cmd.Env = append(os.Environ(), "SUDO_ASKPASS=/bin/false", "SSH_ASKPASS=/bin/false")

	if err := cmd.Run(); err != nil {
		errMsg := stderr.String()
		if errMsg == "" {
			errMsg = stdout.String()
		}
		// Check for context timeout (usually means waiting for sudo)
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("wg-quick up timed out (likely waiting for sudo password): run wgrest as root")
		}
		// Check for permission errors
		if strings.Contains(errMsg, "must be run as root") ||
			strings.Contains(errMsg, "Permission denied") ||
			strings.Contains(errMsg, "Operation not permitted") {
			return fmt.Errorf("wg-quick up requires root privileges: %s", errMsg)
		}
		return fmt.Errorf("wg-quick up failed: %w: %s", err, errMsg)
	}

	return nil
}

// Down brings down a WireGuard interface using wg-quick.
// Note: wg-quick typically requires root/sudo.
func (s *Service) Down(name string) error {
	// Use timeout to prevent indefinite waiting for sudo password
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "wg-quick", "down", name)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	// Prevent any stdin interaction (sudo password prompts) - use explicit devnull
	devNull, err := os.Open(os.DevNull)
	if err == nil {
		cmd.Stdin = devNull
		defer devNull.Close()
	}
	// Set environment to prevent sudo from asking for password
	cmd.Env = append(os.Environ(), "SUDO_ASKPASS=/bin/false", "SSH_ASKPASS=/bin/false")

	if err := cmd.Run(); err != nil {
		errMsg := stderr.String()
		if errMsg == "" {
			errMsg = stdout.String()
		}
		// Check for context timeout (usually means waiting for sudo)
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("wg-quick down timed out (likely waiting for sudo password): run wgrest as root")
		}
		// Check for permission errors
		if strings.Contains(errMsg, "must be run as root") ||
			strings.Contains(errMsg, "Permission denied") ||
			strings.Contains(errMsg, "Operation not permitted") {
			return fmt.Errorf("wg-quick down requires root privileges: %s", errMsg)
		}
		return fmt.Errorf("wg-quick down failed: %w: %s", err, errMsg)
	}

	return nil
}

// GetPlatform returns the current platform (linux, darwin, freebsd).
func GetPlatform() string {
	return runtime.GOOS
}

func (s *Service) buildConfig(device *entity.Device, peers []entity.Peer) string {
	var b strings.Builder

	b.WriteString("[Interface]\n")

	if device.PrivateKey != "" {
		b.WriteString(fmt.Sprintf("PrivateKey = %s\n", device.PrivateKey))
	}

	if device.ListenPort > 0 {
		b.WriteString(fmt.Sprintf("ListenPort = %d\n", device.ListenPort))
	}

	if device.FirewallMark > 0 {
		b.WriteString(fmt.Sprintf("FwMark = %d\n", device.FirewallMark))
	}

	for _, addr := range device.Addresses {
		b.WriteString(fmt.Sprintf("Address = %s\n", addr))
	}

	if len(device.DNS) > 0 {
		b.WriteString(fmt.Sprintf("DNS = %s\n", strings.Join(device.DNS, ", ")))
	}

	if device.MTU > 0 {
		b.WriteString(fmt.Sprintf("MTU = %d\n", device.MTU))
	}

	if device.Table != "" {
		b.WriteString(fmt.Sprintf("Table = %s\n", device.Table))
	}

	for _, cmd := range device.PreUp {
		b.WriteString(fmt.Sprintf("PreUp = %s\n", cmd))
	}

	for _, cmd := range device.PostUp {
		b.WriteString(fmt.Sprintf("PostUp = %s\n", cmd))
	}

	for _, cmd := range device.PreDown {
		b.WriteString(fmt.Sprintf("PreDown = %s\n", cmd))
	}

	for _, cmd := range device.PostDown {
		b.WriteString(fmt.Sprintf("PostDown = %s\n", cmd))
	}

	// Add peers
	for _, peer := range peers {
		b.WriteString("\n[Peer]\n")
		b.WriteString(fmt.Sprintf("PublicKey = %s\n", peer.PublicKey))

		if peer.PresharedKey != "" {
			b.WriteString(fmt.Sprintf("PresharedKey = %s\n", peer.PresharedKey))
		}

		if len(peer.AllowedIPs) > 0 {
			b.WriteString(fmt.Sprintf("AllowedIPs = %s\n", strings.Join(peer.AllowedIPs, ", ")))
		}

		if peer.Endpoint != "" {
			b.WriteString(fmt.Sprintf("Endpoint = %s\n", peer.Endpoint))
		}

		if peer.PersistentKeepaliveInterval != "" && peer.PersistentKeepaliveInterval != "0s" {
			b.WriteString(fmt.Sprintf("PersistentKeepalive = %s\n", peer.PersistentKeepaliveInterval))
		}
	}

	return b.String()
}

// ParseConfig parses a wg-quick configuration file.
func ParseConfig(r io.Reader) (*Config, error) {
	cfg := &Config{}
	var currentPeer *PeerConfig

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Section headers
		if line == "[Interface]" {
			currentPeer = nil
			continue
		}
		if line == "[Peer]" {
			currentPeer = &PeerConfig{}
			cfg.Peers = append(cfg.Peers, *currentPeer)
			continue
		}

		// Key-value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(strings.ToLower(parts[0]))
		value := strings.TrimSpace(parts[1])

		if currentPeer != nil {
			// Peer section
			switch key {
			case "publickey":
				cfg.Peers[len(cfg.Peers)-1].PublicKey = value
			case "presharedkey":
				cfg.Peers[len(cfg.Peers)-1].PresharedKey = value
			case "allowedips":
				ips := strings.Split(value, ",")
				for _, ip := range ips {
					cfg.Peers[len(cfg.Peers)-1].AllowedIPs = append(cfg.Peers[len(cfg.Peers)-1].AllowedIPs, strings.TrimSpace(ip))
				}
			case "endpoint":
				cfg.Peers[len(cfg.Peers)-1].Endpoint = value
			case "persistentkeepalive":
				fmt.Sscanf(value, "%d", &cfg.Peers[len(cfg.Peers)-1].PersistentKeepaliveInterval)
			}
		} else {
			// Interface section
			switch key {
			case "privatekey":
				cfg.PrivateKey = value
			case "listenport":
				fmt.Sscanf(value, "%d", &cfg.ListenPort)
			case "fwmark":
				fmt.Sscanf(value, "%d", &cfg.FirewallMark)
			case "address":
				addrs := strings.Split(value, ",")
				for _, addr := range addrs {
					cfg.Addresses = append(cfg.Addresses, strings.TrimSpace(addr))
				}
			case "dns":
				servers := strings.Split(value, ",")
				for _, s := range servers {
					cfg.DNS = append(cfg.DNS, strings.TrimSpace(s))
				}
			case "mtu":
				fmt.Sscanf(value, "%d", &cfg.MTU)
			case "table":
				cfg.Table = value
			case "preup":
				cfg.PreUp = append(cfg.PreUp, value)
			case "postup":
				cfg.PostUp = append(cfg.PostUp, value)
			case "predown":
				cfg.PreDown = append(cfg.PreDown, value)
			case "postdown":
				cfg.PostDown = append(cfg.PostDown, value)
			case "saveconfig":
				cfg.SaveConfig = strings.ToLower(value) == "true"
			}
		}
	}

	return cfg, scanner.Err()
}
