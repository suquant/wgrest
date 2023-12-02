package models

import (
	"encoding/base64"
	"fmt"
	"net"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var (
	EmptyKey = wgtypes.Key{}
)

func NewPeer(peer wgtypes.Peer) Peer {
	allowedIPs := make([]string, len(peer.AllowedIPs))
	for i, v := range peer.AllowedIPs {
		allowedIPs[i] = v.String()
	}

	p := Peer{
		PublicKey:                   peer.PublicKey.String(),
		UrlSafePublicKey:            base64.URLEncoding.EncodeToString(peer.PublicKey[:]),
		AllowedIps:                  allowedIPs,
		LastHandshakeTime:           peer.LastHandshakeTime,
		ReceiveBytes:                peer.ReceiveBytes,
		TransmitBytes:               peer.TransmitBytes,
		PersistentKeepaliveInterval: peer.PersistentKeepaliveInterval.String(),
	}

	if peer.PresharedKey != EmptyKey {
		p.PresharedKey = peer.PresharedKey.String()
	}

	if peer.Endpoint != nil {
		p.Endpoint = peer.Endpoint.String()
	}

	return p
}

func (r *PeerCreateOrUpdateRequest) Apply(conf *wgtypes.PeerConfig) error {
	if r.Endpoint != "" {
		endpoint, err := net.ResolveUDPAddr("udp", r.Endpoint)
		if err != nil {
			return err
		}

		conf.Endpoint = endpoint
	}

	if r.PersistentKeepaliveInterval != "" {
		keepaliveInterval, err := time.ParseDuration(r.PersistentKeepaliveInterval)
		if err != nil {
			return err
		}

		conf.PersistentKeepaliveInterval = &keepaliveInterval
	}

	if r.AllowedIps != nil {
		allowedIPs := make([]net.IPNet, len(*r.AllowedIps))
		for i, v := range *r.AllowedIps {
			_, ipNet, err := net.ParseCIDR(v)
			if err != nil {
				return err
			}

			if ipNet == nil {
				return fmt.Errorf("failed to parse CIDR: %s", v)
			}

			allowedIPs[i] = *ipNet
		}

		conf.AllowedIPs = allowedIPs
	}

	if r.PresharedKey != nil {
		psKey, err := wgtypes.ParseKey(*r.PresharedKey)
		if err != nil {
			return err
		}

		conf.PresharedKey = &psKey
	}

	if r.PrivateKey != nil {
		privKey, err := wgtypes.ParseKey(*r.PrivateKey)
		if err != nil {
			return err
		}

		if privKey.PublicKey() != conf.PublicKey {
			return fmt.Errorf("wrong private key")
		}
	}

	return nil
}
