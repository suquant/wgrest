package wireguard

import (
	"fmt"
	"net/http"
	"os"

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
	deviceName := params.Dev
	_, err := GetDeviceByName(deviceName)
	if err != nil {
		switch {
		case IsErrNotFound(err):
			return wireguard.NewDeviceGetNotFound()
		default:
			wgrest.Logger.Println(err.Error())

			return wireguard.NewDeviceGetDefault(http.StatusInternalServerError).WithPayload(
				&models.Error{
					Detail: err.Error(),
				},
			)
		}
	}

	la := netlink.NewLinkAttrs()
	la.Name = deviceName

	dev := &netlink.GenericLink{
		LinkAttrs: la,
		LinkType:  "wireguard",
	}

	err = netlink.LinkDel(dev)
	if err != nil {
		msg := fmt.Sprintf("netlink err: %s\n", err.Error())
		wgrest.Logger.Println(msg)

		return wireguard.NewDeviceDeleteDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{Detail: msg},
		)
	}

	wgConfPath := getWgConfPath(params.Dev)
	err = os.Remove(wgConfPath)
	if err != nil {
		msg := fmt.Sprintf("os err: %s\n", err.Error())
		wgrest.Logger.Println(msg)

		return wireguard.NewDeviceDeleteDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{Detail: msg},
		)
	}

	return wireguard.NewDeviceDeleteNoContent()
}
