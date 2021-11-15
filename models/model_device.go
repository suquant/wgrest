package models

// Device - Information about wireguard device.
type Device struct {

	// WireGuard device name. Usually it is network interface name
	Name string `json:"name"`

	// WireGuard device listen port.
	ListenPort int32 `json:"listen_port"`

	// WireGuard device public key encoded by base64.
	PublicKey string `json:"public_key"`

	// WireGuard device firewall mark.
	FirewallMark int32 `json:"firewall_mark"`

	// IPv4 or IPv6 addresses in CIDR notation
	Networks []string `json:"networks"`

	// WireGuard device's peers count
	PeersCount int32 `json:"peers_count"`

	// WireGuard device's peers total receive bytes
	TotalReceiveBytes int64 `json:"total_receive_bytes"`

	// WireGuard device's peers total transmit bytes
	TotalTransmitBytes int64 `json:"total_transmit_bytes"`
}
