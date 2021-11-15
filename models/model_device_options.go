package models

// DeviceOptions - Device options
type DeviceOptions struct {

	// Device's allowed ips, it might be any of IPv4 or IPv6 addresses in CIDR notation. It might be owervrite in peer and device config.
	AllowedIps []string `json:"allowed_ips"`

	// Interface's DNS servers.
	DnsServers []string `json:"dns_servers"`

	// Device host, it might be domain name or IPv4/IPv6 address. It is used for external/internal connection
	Host string `json:"host"`
}
