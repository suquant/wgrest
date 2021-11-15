package models

// PeerCreateOrUpdateRequest - Peer params that might be used due to creation or updation process
type PeerCreateOrUpdateRequest struct {

	// Base64 encoded private key. If present it will be stored in persistent storage.
	PrivateKey *string `json:"private_key,omitempty"`

	// Base64 encoded public key
	PublicKey *string `json:"public_key,omitempty"`

	// Base64 encoded preshared key
	PresharedKey *string `json:"preshared_key,omitempty"`

	// Peer's allowed ips, it might be any of IPv4 or IPv6 addresses in CIDR notation
	AllowedIps *[]string `json:"allowed_ips,omitempty"`

	// Peer's persistend keepalive interval. Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".
	PersistentKeepaliveInterval string `json:"persistent_keepalive_interval,omitempty"`

	// Peer's endpoint in host:port format
	Endpoint string `json:"endpoint,omitempty"`
}
