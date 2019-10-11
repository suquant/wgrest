package wireguard

import (
	"github.com/suquant/wgrest/restapi/operations"
	"github.com/suquant/wgrest/restapi/operations/wireguard"
)

// ApplyAPI apply to api
func ApplyAPI(api *operations.WgrestAPI) {
	// device
	api.WireguardDeviceCreateHandler = wireguard.DeviceCreateHandlerFunc(DeviceCreateHandler)
	api.WireguardDeviceListHandler = wireguard.DeviceListHandlerFunc(DeviceListHandler)
	api.WireguardDeviceGetHandler = wireguard.DeviceGetHandlerFunc(DeviceGetHandler)
	api.WireguardDeviceDeleteHandler = wireguard.DeviceDeleteHandlerFunc(DeviceDeleteHandler)

	// peer
	api.WireguardPeerCreateHandler = wireguard.PeerCreateHandlerFunc(PeerCreateHandler)
	api.WireguardPeerListHandler = wireguard.PeerListHandlerFunc(PeerListHandler)
	api.WireguardPeerGetHandler = wireguard.PeerGetHandlerFunc(PeerGetHandler)
	api.WireguardPeerDeleteHandler = wireguard.PeerDeleteHandlerFunc(PeerDeleteHandler)
}
