package models

import (
	"time"
)

// Peer - Information about wireguard peer.
type Peer struct {

	// Base64 encoded public key
	PublicKey string `json:"public_key"`

	// URL safe base64 encoded public key. It is usefull to use in peers api endpoint.
	UrlSafePublicKey string `json:"url_safe_public_key"`

	// Base64 encoded preshared key
	PresharedKey string `json:"preshared_key,omitempty"`

	// Peer's allowed ips, it might be any of IPv4 or IPv6 addresses in CIDR notation
	AllowedIps []string `json:"allowed_ips"`

	// Peer's last handshake time formated in RFC3339
	LastHandshakeTime time.Time `json:"last_handshake_time"`

	// Peer's persistend keepalive interval in
	PersistentKeepaliveInterval string `json:"persistent_keepalive_interval"`

	// Peer's endpoint in host:port format
	Endpoint string `json:"endpoint"`

	// Peer's receive bytes
	ReceiveBytes int64 `json:"receive_bytes"`

	// Peer's transmit bytes
	TransmitBytes int64 `json:"transmit_bytes"`
}
