package wireguard

import (
	"net/http"

	"github.com/suquant/wgrest"

	"github.com/go-openapi/runtime/middleware"
	"github.com/suquant/wgrest/models"
	"github.com/suquant/wgrest/restapi/operations/wireguard"
	"github.com/vishvananda/netlink"
)

// DeviceDeleteHandler wireguard device delete handler
func DeviceDeleteHandler(
	params wireguard.DeviceDeleteParams,
	principal interface{},
) middleware.Responder {
	la := netlink.NewLinkAttrs()
	la.Name = params.Dev

	dev := &netlink.GenericLink{
		LinkAttrs: la,
		LinkType:  "wireguard",
	}

	err := netlink.LinkDel(dev)
	if err != nil {
		wgrest.Logger.Printf("netlink err: %s\n", err.Error())

		return wireguard.NewDeviceDeleteDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{Detail: err.Error()},
		)
	}

	return wireguard.NewDeviceDeleteNoContent()
}
