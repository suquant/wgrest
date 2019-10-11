package wireguard

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/suquant/wgrest"
	"github.com/suquant/wgrest/models"
	"github.com/suquant/wgrest/restapi/operations/wireguard"
)

// PeerListHandler device peer list
func PeerListHandler(
	params wireguard.PeerListParams,
	pricipal interface{},
) middleware.Responder {
	dev, err := getWGDeviceByName(params.Dev)
	if err != nil {
		switch {
		case IsErrNotFound(err):
			return wireguard.NewPeerListNotFound()
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

	peers := make([]*models.WireguardPeer, len(dev.Peers))
	for i, p := range dev.Peers {
		peer, err := getPeerByWGPeer(p)
		if err != nil {
			msg := fmt.Sprintf(err.Error())
			wgrest.Logger.Println(msg)

			return wireguard.NewDeviceGetDefault(http.StatusInternalServerError).WithPayload(
				&models.Error{
					Detail: msg,
				},
			)
		}

		peers[i] = peer
	}

	return wireguard.NewPeerListOK().WithPayload(peers)
}
