package utils

import (
	"bytes"
	"fmt"
	"github.com/suquant/wgrest/models"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"io"
	"strings"
)

type PeerQuickConfigOptions struct {
	PrivateKey *string
	DNSServers *[]string
	AllowedIPs *[]string
	Host       *string
}

func GetPeerQuickConfig(device wgtypes.Device, peer wgtypes.Peer, options PeerQuickConfigOptions) (io.Reader, error) {
	b := &bytes.Buffer{}
	fmt.Fprintln(b, "[Interface]")
	if options.PrivateKey != nil {
		fmt.Fprintln(b, "PrivateKey =", *options.PrivateKey)
	}

	addresses := make([]string, len(peer.AllowedIPs))
	for i, v := range peer.AllowedIPs {
		addresses[i] = v.String()
	}

	fmt.Fprintf(b, "Address = %s\n", strings.Join(addresses, ","))
	if options.DNSServers != nil && len(*options.DNSServers) > 0 {
		fmt.Fprintf(b, "DNS = %s\n", strings.Join(*options.DNSServers, ","))
	}

	fmt.Fprintln(b, "")
	fmt.Fprintln(b, "[Peer]")

	fmt.Fprintf(b, "PublicKey = %s\n", device.PublicKey.String())
	if peer.PresharedKey != models.EmptyKey {
		fmt.Fprintf(b, "PresharedKey = %s\n", peer.PresharedKey.String())
	}
	if options.Host != nil {
		fmt.Fprintf(b, "Endpoint = %s:%v\n", *options.Host, device.ListenPort)
	}
	if options.AllowedIPs != nil && len(*options.AllowedIPs) > 0 {
		fmt.Fprintf(b, "AllowedIPs = %s\n", strings.Join(*options.AllowedIPs, ","))
	}

	return b, nil
}
