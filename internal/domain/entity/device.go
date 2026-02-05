package entity

// Device represents a WireGuard device/interface with wg-quick configuration.
type Device struct {
	// Name is the WireGuard interface name
	Name string `json:"name"`

	// ListenPort is the WireGuard listen port
	ListenPort int32 `json:"listen_port"`

	// PublicKey is the device public key (base64)
	PublicKey string `json:"public_key"`

	// PrivateKey is the device private key (base64) - only returned on create
	PrivateKey string `json:"private_key,omitempty"`

	// FirewallMark is the device firewall mark
	FirewallMark int32 `json:"firewall_mark,omitempty"`

	// --- wg-quick [Interface] options ---

	// Addresses are IP addresses to assign to the interface (CIDR notation)
	Addresses []string `json:"addresses,omitempty"`

	// DNS servers to configure when interface is up
	DNS []string `json:"dns,omitempty"`

	// MTU for the interface (auto-calculated if not set)
	MTU int32 `json:"mtu,omitempty"`

	// Table is the routing table (auto, off, or table number)
	Table string `json:"table,omitempty"`

	// PreUp commands to run before interface comes up
	PreUp []string `json:"pre_up,omitempty"`

	// PostUp commands to run after interface comes up
	PostUp []string `json:"post_up,omitempty"`

	// PreDown commands to run before interface goes down
	PreDown []string `json:"pre_down,omitempty"`

	// PostDown commands to run after interface goes down
	PostDown []string `json:"post_down,omitempty"`

	// --- Status ---

	// Running indicates if the interface is currently up
	Running bool `json:"running"`

	// --- Statistics ---

	// PeersCount is the number of connected peers
	PeersCount int32 `json:"peers_count"`

	// TotalReceiveBytes across all peers
	TotalReceiveBytes int64 `json:"total_receive_bytes"`

	// TotalTransmitBytes across all peers
	TotalTransmitBytes int64 `json:"total_transmit_bytes"`
}

// DeviceCreateOrUpdateRequest represents parameters for creating or updating a device.
type DeviceCreateOrUpdateRequest struct {
	// Name is required for creation
	Name *string `json:"name,omitempty"`

	// ListenPort for the WireGuard interface
	ListenPort *int32 `json:"listen_port,omitempty"`

	// PrivateKey (base64) - if not provided, one will be generated
	PrivateKey *string `json:"private_key,omitempty"`

	// FirewallMark for the interface
	FirewallMark *int32 `json:"firewall_mark,omitempty"`

	// --- wg-quick options ---

	// Addresses to assign (CIDR notation)
	Addresses []string `json:"addresses,omitempty"`

	// DNS servers
	DNS []string `json:"dns,omitempty"`

	// MTU for the interface
	MTU *int32 `json:"mtu,omitempty"`

	// Table routing table
	Table *string `json:"table,omitempty"`

	// PreUp commands
	PreUp []string `json:"pre_up,omitempty"`

	// PostUp commands
	PostUp []string `json:"post_up,omitempty"`

	// PreDown commands
	PreDown []string `json:"pre_down,omitempty"`

	// PostDown commands
	PostDown []string `json:"post_down,omitempty"`
}
