package models

import "github.com/suquant/wgrest/storage"

func NewDeviceOptions(options storage.StoreDeviceOptions) DeviceOptions {
	return DeviceOptions{
		Host:       options.Host,
		AllowedIps: options.AllowedIPs,
		DnsServers: options.DNSServers,
	}
}

func (r *DeviceOptionsUpdateRequest) Apply(options *storage.StoreDeviceOptions) error {
	if r.Host != nil {
		options.Host = *r.Host
	}

	if r.DnsServers != nil {
		options.DNSServers = *r.DnsServers
	}

	if r.AllowedIps != nil {
		options.AllowedIPs = *r.AllowedIps
	}

	return nil
}
