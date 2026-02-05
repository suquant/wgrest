package entity

import "time"

// Peer represents a WireGuard peer.
type Peer struct {
	// PublicKey is the base64 encoded public key
	PublicKey string `json:"public_key"`

	// URLSafePublicKey is the URL-safe base64 encoded public key
	URLSafePublicKey string `json:"url_safe_public_key"`

	// PrivateKey is the base64 encoded private key (stored separately)
	PrivateKey string `json:"private_key,omitempty"`

	// PresharedKey is the base64 encoded preshared key
	PresharedKey string `json:"preshared_key,omitempty"`

	// AllowedIPs are the peer's allowed IPs in CIDR notation
	AllowedIPs []string `json:"allowed_ips"`

	// LastHandshakeTime is the peer's last handshake time (RFC3339)
	LastHandshakeTime time.Time `json:"last_handshake_time"`

	// PersistentKeepaliveInterval is the peer's keepalive interval
	PersistentKeepaliveInterval string `json:"persistent_keepalive_interval,omitempty"`

	// Endpoint is the peer's endpoint in host:port format
	Endpoint string `json:"endpoint,omitempty"`

	// ReceiveBytes is the number of bytes received from this peer
	ReceiveBytes int64 `json:"receive_bytes"`

	// TransmitBytes is the number of bytes transmitted to this peer
	TransmitBytes int64 `json:"transmit_bytes"`
}

// PeerCreateOrUpdateRequest represents parameters for creating or updating a peer.
type PeerCreateOrUpdateRequest struct {
	PrivateKey                  *string  `json:"private_key,omitempty"`
	PublicKey                   *string  `json:"public_key,omitempty"`
	PresharedKey                *string  `json:"preshared_key,omitempty"`
	AllowedIPs                  []string `json:"allowed_ips,omitempty"`
	PersistentKeepaliveInterval *string  `json:"persistent_keepalive_interval,omitempty"`
	Endpoint                    *string  `json:"endpoint,omitempty"`
}

