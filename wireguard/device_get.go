package wireguard

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/suquant/wgrest"
	"github.com/suquant/wgrest/models"
	"github.com/suquant/wgrest/restapi/operations/wireguard"
)

// DeviceGetHandler wireguard device get details
func DeviceGetHandler(
	params wireguard.DeviceGetParams,
	principal interface{},
) middleware.Responder {
	device, err := GetDeviceByName(params.Dev)
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

	return wireguard.NewDeviceGetOK().WithPayload(device)
}
