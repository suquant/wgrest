package wgquick

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/suquant/wgrest/internal/domain/entity"
)

func TestParseConfig_Basic(t *testing.T) {
	configStr := `[Interface]
PrivateKey = cGVla2Fib28K
ListenPort = 51820

[Peer]
PublicKey = cHVibGljS2V5Cg==
AllowedIPs = 10.0.0.2/32
Endpoint = 192.168.1.1:51820
`

	cfg, err := ParseConfig(strings.NewReader(configStr))
	require.NoError(t, err)

	assert.Equal(t, "cGVla2Fib28K", cfg.PrivateKey)
	assert.Equal(t, 51820, cfg.ListenPort)
	assert.Len(t, cfg.Peers, 1)
	assert.Equal(t, "cHVibGljS2V5Cg==", cfg.Peers[0].PublicKey)
	assert.Equal(t, []string{"10.0.0.2/32"}, cfg.Peers[0].AllowedIPs)
	assert.Equal(t, "192.168.1.1:51820", cfg.Peers[0].Endpoint)
}

func TestParseConfig_AllInterfaceOptions(t *testing.T) {
	configStr := `[Interface]
PrivateKey = cGVla2Fib28K
ListenPort = 51820
Address = 10.0.0.1/24, fd00::1/64
DNS = 1.1.1.1, 8.8.8.8
MTU = 1420
Table = auto
FwMark = 0x1234
PreUp = echo "pre-up"
PostUp = iptables -A FORWARD -i wg0 -j ACCEPT
PreDown = echo "pre-down"
PostDown = iptables -D FORWARD -i wg0 -j ACCEPT
SaveConfig = true
`

	cfg, err := ParseConfig(strings.NewReader(configStr))
	require.NoError(t, err)

	assert.Equal(t, "cGVla2Fib28K", cfg.PrivateKey)
	assert.Equal(t, 51820, cfg.ListenPort)
	assert.Equal(t, []string{"10.0.0.1/24", "fd00::1/64"}, cfg.Addresses)
	assert.Equal(t, []string{"1.1.1.1", "8.8.8.8"}, cfg.DNS)
	assert.Equal(t, 1420, cfg.MTU)
	assert.Equal(t, "auto", cfg.Table)
	assert.Equal(t, []string{`echo "pre-up"`}, cfg.PreUp)
	assert.Equal(t, []string{"iptables -A FORWARD -i wg0 -j ACCEPT"}, cfg.PostUp)
	assert.Equal(t, []string{`echo "pre-down"`}, cfg.PreDown)
	assert.Equal(t, []string{"iptables -D FORWARD -i wg0 -j ACCEPT"}, cfg.PostDown)
	assert.True(t, cfg.SaveConfig)
}

func TestParseConfig_MultiplePeers(t *testing.T) {
	configStr := `[Interface]
PrivateKey = cGVla2Fib28K
ListenPort = 51820

[Peer]
PublicKey = peer1PublicKey
AllowedIPs = 10.0.0.2/32
PersistentKeepalive = 25

[Peer]
PublicKey = peer2PublicKey
AllowedIPs = 10.0.0.3/32, 10.0.0.4/32
Endpoint = vpn.example.com:51820
PresharedKey = cHJlc2hhcmVkS2V5
`

	cfg, err := ParseConfig(strings.NewReader(configStr))
	require.NoError(t, err)

	assert.Len(t, cfg.Peers, 2)

	// First peer
	assert.Equal(t, "peer1PublicKey", cfg.Peers[0].PublicKey)
	assert.Equal(t, []string{"10.0.0.2/32"}, cfg.Peers[0].AllowedIPs)
	assert.Equal(t, 25, cfg.Peers[0].PersistentKeepaliveInterval)

	// Second peer
	assert.Equal(t, "peer2PublicKey", cfg.Peers[1].PublicKey)
	assert.Equal(t, []string{"10.0.0.3/32", "10.0.0.4/32"}, cfg.Peers[1].AllowedIPs)
	assert.Equal(t, "vpn.example.com:51820", cfg.Peers[1].Endpoint)
	assert.Equal(t, "cHJlc2hhcmVkS2V5", cfg.Peers[1].PresharedKey)
}

func TestParseConfig_Comments(t *testing.T) {
	configStr := `# This is a comment
[Interface]
PrivateKey = cGVla2Fib28K
# Another comment
ListenPort = 51820
`

	cfg, err := ParseConfig(strings.NewReader(configStr))
	require.NoError(t, err)

	assert.Equal(t, "cGVla2Fib28K", cfg.PrivateKey)
	assert.Equal(t, 51820, cfg.ListenPort)
}

func TestParseConfig_EmptyFile(t *testing.T) {
	cfg, err := ParseConfig(strings.NewReader(""))
	require.NoError(t, err)
	assert.Empty(t, cfg.PrivateKey)
	assert.Empty(t, cfg.Peers)
}

func TestBuildConfig(t *testing.T) {
	svc := &Service{configDirs: []string{"/tmp"}}

	device := &entity.Device{
		Name:       "wg0",
		PrivateKey: "privateKeyBase64",
		ListenPort: 51820,
		Addresses:  []string{"10.0.0.1/24"},
		DNS:        []string{"1.1.1.1", "8.8.8.8"},
		MTU:        1420,
		Table:      "auto",
		PreUp:      []string{"echo pre-up"},
		PostUp:     []string{"echo post-up"},
		PreDown:    []string{"echo pre-down"},
		PostDown:   []string{"echo post-down"},
	}

	peers := []entity.Peer{
		{
			PublicKey:                   "peerPublicKey",
			AllowedIPs:                  []string{"10.0.0.2/32"},
			Endpoint:                    "192.168.1.1:51820",
			PersistentKeepaliveInterval: "25s",
		},
	}

	config := svc.buildConfig(device, peers)

	assert.Contains(t, config, "[Interface]")
	assert.Contains(t, config, "PrivateKey = privateKeyBase64")
	assert.Contains(t, config, "ListenPort = 51820")
	assert.Contains(t, config, "Address = 10.0.0.1/24")
	assert.Contains(t, config, "DNS = 1.1.1.1, 8.8.8.8")
	assert.Contains(t, config, "MTU = 1420")
	assert.Contains(t, config, "Table = auto")
	assert.Contains(t, config, "PreUp = echo pre-up")
	assert.Contains(t, config, "PostUp = echo post-up")
	assert.Contains(t, config, "PreDown = echo pre-down")
	assert.Contains(t, config, "PostDown = echo post-down")

	assert.Contains(t, config, "[Peer]")
	assert.Contains(t, config, "PublicKey = peerPublicKey")
	assert.Contains(t, config, "AllowedIPs = 10.0.0.2/32")
	assert.Contains(t, config, "Endpoint = 192.168.1.1:51820")
	assert.Contains(t, config, "PersistentKeepalive = 25s")
}

func TestGetPlatform(t *testing.T) {
	platform := GetPlatform()
	// Should be one of the supported platforms
	validPlatforms := []string{"linux", "darwin", "freebsd", "windows"}
	assert.Contains(t, validPlatforms, platform)
}

func TestNewService(t *testing.T) {
	// Use a temp directory
	svc, err := NewService([]string{t.TempDir()})
	require.NoError(t, err)
	assert.NotNil(t, svc)
}

func TestFindConfigPath(t *testing.T) {
	// Use temp dirs that don't exist anywhere
	svc := &Service{configDirs: []string{"/nonexistent/first", "/nonexistent/second"}}
	path := svc.FindConfigPath("wg0")
	// Returns first dir path when not found
	assert.Equal(t, "/nonexistent/first/wg0.conf", path)
}

func TestFindConfigPath_ExistingFile(t *testing.T) {
	// Create temp dir with a config file
	tmpDir := t.TempDir()
	configPath := tmpDir + "/wg0.conf"
	err := os.WriteFile(configPath, []byte("[Interface]\n"), 0600)
	require.NoError(t, err)

	svc := &Service{configDirs: []string{"/nonexistent", tmpDir}}
	path := svc.FindConfigPath("wg0")
	// Should find the existing file in tmpDir
	assert.Equal(t, configPath, path)
}

func TestGetConfigDir(t *testing.T) {
	// Create temp dir with a config file
	tmpDir := t.TempDir()
	configPath := tmpDir + "/wg1.conf"
	err := os.WriteFile(configPath, []byte("[Interface]\n"), 0600)
	require.NoError(t, err)

	svc := &Service{configDirs: []string{"/nonexistent", tmpDir}}
	dir := svc.GetConfigDir("wg1")
	assert.Equal(t, tmpDir, dir)

	// Non-existent config returns first dir
	dir = svc.GetConfigDir("wg99")
	assert.Equal(t, "/nonexistent", dir)
}

func TestListConfigDevices(t *testing.T) {
	// Create temp dir with config files
	tmpDir := t.TempDir()
	err := os.WriteFile(tmpDir+"/wg0.conf", []byte("[Interface]\n"), 0600)
	require.NoError(t, err)
	err = os.WriteFile(tmpDir+"/wg1.conf", []byte("[Interface]\n"), 0600)
	require.NoError(t, err)
	// Create a non-config file (should be ignored)
	err = os.WriteFile(tmpDir+"/readme.txt", []byte("test"), 0600)
	require.NoError(t, err)

	svc := &Service{configDirs: []string{tmpDir}}
	devices := svc.ListConfigDevices()

	assert.Len(t, devices, 2)
	assert.Contains(t, devices, "wg0")
	assert.Contains(t, devices, "wg1")
}

func TestListConfigDevices_MultipleDirs(t *testing.T) {
	// Create two temp dirs with config files
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()

	err := os.WriteFile(tmpDir1+"/wg0.conf", []byte("[Interface]\n"), 0600)
	require.NoError(t, err)
	err = os.WriteFile(tmpDir2+"/wg1.conf", []byte("[Interface]\n"), 0600)
	require.NoError(t, err)
	// Same device in both dirs (should be deduplicated)
	err = os.WriteFile(tmpDir2+"/wg0.conf", []byte("[Interface]\n"), 0600)
	require.NoError(t, err)

	svc := &Service{configDirs: []string{tmpDir1, tmpDir2}}
	devices := svc.ListConfigDevices()

	// Should have 2 unique devices, not 3
	assert.Len(t, devices, 2)
	assert.Contains(t, devices, "wg0")
	assert.Contains(t, devices, "wg1")
}

func TestNewService_EmptyDirs(t *testing.T) {
	_, err := NewService([]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one config directory")
}

func TestConfigDirs(t *testing.T) {
	dirs := []string{"/a", "/b", "/c"}
	svc := &Service{configDirs: dirs}
	assert.Equal(t, dirs, svc.ConfigDirs())
}

