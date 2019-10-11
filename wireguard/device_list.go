package wireguard

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/suquant/wgrest"
	"github.com/suquant/wgrest/models"
	"github.com/suquant/wgrest/restapi/operations/wireguard"
)

// DeviceListHandler wireguard device list
func DeviceListHandler(
	params wireguard.DeviceListParams,
	principal interface{},
) middleware.Responder {
	devices, err := GetDevices()
	if err != nil {
		wgrest.Logger.Println(err.Error())

		return wireguard.NewDeviceListDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{Detail: err.Error()},
		)
	}

	return wireguard.NewDeviceListOK().WithPayload(devices)
}
