package wireguard

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/suquant/wgrest"
	"github.com/suquant/wgrest/models"
	"github.com/suquant/wgrest/restapi/operations/wireguard"
)

// PeerGetHandler get device peer by public key
func PeerGetHandler(
	params wireguard.PeerGetParams,
	pricipal interface{},
) middleware.Responder {
	wgPeer, err := getWGPeerByPeerID(params.Dev, params.PeerID)
	if err != nil {
		switch {
		case IsErrNotFound(err):
			return wireguard.NewPeerGetNotFound()
		default:
			msg := fmt.Sprintf(err.Error())
			wgrest.Logger.Println(msg)

			return wireguard.NewDeviceGetDefault(http.StatusInternalServerError).WithPayload(
				&models.Error{
					Detail: msg,
				},
			)
		}
	}

	peer, err := getPeerByWGPeer(*wgPeer)
	if err != nil {
		msg := fmt.Sprintf(err.Error())
		wgrest.Logger.Println(msg)

		return wireguard.NewDeviceGetDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{
				Detail: msg,
			},
		)
	}

	return wireguard.NewPeerGetOK().WithPayload(peer)
}
