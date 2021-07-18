package wireguard

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"os"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"github.com/suquant/wgrest/models"
	"golang.zx2c4.com/wireguard/wgctrl"
)

var (
	// ErrNotFound not found error
	ErrNotFound = errors.New("not found")
)

// IsErrNotFound return true if error is ErrNotFound
func IsErrNotFound(err error) bool {
	return err == ErrNotFound
}

func getWGDeviceByName(name string) (*wgtypes.Device, error) {
	client, err := wgctrl.New()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	wgDevice, err := client.Device(name)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return wgDevice, nil
}

func getDeviceByWGDevice(wgdev *wgtypes.Device) *models.WireguardDevice {
	listenPort := int64(wgdev.ListenPort)
	name := wgdev.Name
	privateKey := wgdev.PrivateKey.String()
	publicKey := wgdev.PublicKey.String()
	var network string

	addrs, err := getInterfaceNetworks(name)
	if err != nil {
		// logging.Logger.Error(
		// 	fmt.Sprintf("iface err: %s", err.Error()),
		// )
	} else {
		for _, a := range addrs {
			networkString := a.String()
			network = networkString
		}
	}

	return &models.WireguardDevice{
		ListenPort: &listenPort,
		Name:       &name,
		PrivateKey: &privateKey,
		PublicKey:  publicKey,
		Network:    &network,
	}
}

// getInterfaceNetworks get interface network
func getInterfaceNetworks(name string) ([]net.Addr, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}

	return iface.Addrs()
}

// GetDevices get all devices
func GetDevices() ([]*models.WireguardDevice, error) {
	client, err := wgctrl.New()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	wgDevices, err := client.Devices()
	if err != nil {
		return nil, err
	}

	devices := make([]*models.WireguardDevice, len(wgDevices))
	for i, wgDevice := range wgDevices {
		devices[i] = getDeviceByWGDevice(wgDevice)
	}

	return devices, nil
}

// GetDeviceByName get device by name
func GetDeviceByName(name string) (*models.WireguardDevice, error) {
	wgDevice, err := getWGDeviceByName(name)
	if err != nil {
		return nil, err
	}

	return getDeviceByWGDevice(wgDevice), nil
}

func addOrRemovePeer(dev string, peer *wgtypes.Peer, remove bool) error {
	client, err := wgctrl.New()
	if err != nil {
		return err
	}
	defer client.Close()

	var presharedKey *wgtypes.Key
	if len(peer.PresharedKey) > 0 {
		presharedKey = &peer.PresharedKey
	}

	peerCfg := wgtypes.PeerConfig{
		PublicKey:                   peer.PublicKey,
		PresharedKey:                presharedKey,
		Endpoint:                    peer.Endpoint,
		PersistentKeepaliveInterval: &peer.PersistentKeepaliveInterval,

		ReplaceAllowedIPs: true,
		AllowedIPs:        peer.AllowedIPs,

		Remove: remove,
	}

	cfg := wgtypes.Config{
		ReplacePeers: false,
		Peers:        []wgtypes.PeerConfig{peerCfg},
	}

	return client.ConfigureDevice(dev, cfg)
}

func getWGPeerByPeer(peer *models.WireguardPeer) (*wgtypes.Peer, error) {
	publicKey, err := wgtypes.ParseKey(*peer.PublicKey)
	if err != nil {
		return nil, err
	}

	var presharedKey wgtypes.Key
	if peer.PresharedKey != "" {
		presharedKey, err = wgtypes.ParseKey(peer.PresharedKey)
		if err != nil {
			return nil, err
		}
	}

	var allowedIPs []net.IPNet
	for _, cidr := range peer.AllowedIps {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}

		allowedIPs = append(allowedIPs, *ipnet)
	}

	return &wgtypes.Peer{
		PublicKey:    publicKey,
		PresharedKey: presharedKey,
		AllowedIPs:   allowedIPs,
	}, nil
}

func getPeerByWGPeer(peer wgtypes.Peer) (*models.WireguardPeer, error) {
	publicKey := peer.PublicKey.String()
	allowedIPs := make([]string, len(peer.AllowedIPs))
	for i, ip := range peer.AllowedIPs {
		allowedIPs[i] = ip.String()
	}
	presharedKey := ""
	var emptyKey wgtypes.Key
	if peer.PresharedKey != emptyKey {
		presharedKey = peer.PresharedKey.String()
	}

	peerID, err := getPeerID(peer)
	if err != nil {
		return nil, err
	}

	return &models.WireguardPeer{
		PublicKey:    &publicKey,
		PresharedKey: presharedKey,
		AllowedIps:   allowedIPs,
		PeerID:       peerID,
	}, nil
}

func getWGPeerByPeerID(dev, peerID string) (*wgtypes.Peer, error) {
	wg, err := getWGDeviceByName(dev)
	if err != nil {
		return nil, err
	}

	decodedPublicKey, err := base64.URLEncoding.DecodeString(peerID)
	if err != nil {
		return nil, err
	}

	publicKey, err := wgtypes.NewKey(decodedPublicKey)
	if err != nil {
		return nil, err
	}

	for _, p := range wg.Peers {
		if p.PublicKey.String() == publicKey.String() {
			return &p, nil
		}
	}

	return nil, ErrNotFound
}

func getPeerID(peer wgtypes.Peer) (string, error) {
	return base64.URLEncoding.EncodeToString(peer.PublicKey[:]), nil
}

func getWgConfPath(dev string) string {
	return fmt.Sprintf("/etc/wireguard/%s.conf", dev)
}
