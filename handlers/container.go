package handlers

import (
	"github.com/suquant/wgrest/storage"
)

type WireGuardContainerOptions struct {
	Storage              storage.Storage
	DefaultDeviceOptions storage.StoreDeviceOptions
}

// WireGuardContainer will hold all dependencies for your application.
type WireGuardContainer struct {
	storage              storage.Storage
	defaultDeviceOptions storage.StoreDeviceOptions
}

// NewWireGuardContainer returns an empty or an initialized container for your handlers.
func NewWireGuardContainer(options WireGuardContainerOptions) (WireGuardContainer, error) {
	c := WireGuardContainer{
		storage:              options.Storage,
		defaultDeviceOptions: options.DefaultDeviceOptions,
	}

	return c, nil
}
