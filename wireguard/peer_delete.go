package wireguard

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/suquant/wgrest"
	"github.com/suquant/wgrest/models"
	"github.com/suquant/wgrest/restapi/operations/wireguard"
)

// PeerDeleteHandler delete peer from device
func PeerDeleteHandler(
	params wireguard.PeerDeleteParams,
	principal interface{},
) middleware.Responder {
	dev := params.Dev
	wgPeer, err := getWGPeerByPeerID(dev, params.PeerID)
	if err != nil {
		switch {
		case IsErrNotFound(err):
			return wireguard.NewPeerDeleteNotFound()
		default:
			msg := fmt.Sprintf(err.Error())
			wgrest.Logger.Println(msg)

			return wireguard.NewDeviceDeleteDefault(http.StatusInternalServerError).WithPayload(
				&models.Error{
					Detail: msg,
				},
			)
		}
	}

	err = addOrRemovePeer(dev, wgPeer, true)
	if err != nil {
		msg := fmt.Sprintf(err.Error())
		wgrest.Logger.Println(msg)

		return wireguard.NewDeviceDeleteDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{
				Detail: msg,
			},
		)
	}

	return wireguard.NewPeerDeleteNoContent()
}
