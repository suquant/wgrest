package models

// DeviceCreateOrUpdateRequest - Device params that might be used due to creation or updation process
type DeviceCreateOrUpdateRequest struct {

	// WireGuard device name. Usually it is network interface name
	Name *string `json:"name,omitempty"`

	// WireGuard device listen port.
	ListenPort *int32 `json:"listen_port,omitempty"`

	// WireGuard device private key encoded by base64.
	PrivateKey *string `json:"private_key,omitempty"`

	// WireGuard device firewall mark.
	FirewallMark *int32 `json:"firewall_mark,omitempty"`

	// IPv4 or IPv6 addresses in CIDR notation
	Networks *[]string `json:"networks,omitempty"`
}
