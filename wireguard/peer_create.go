package wireguard

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/suquant/wgrest"
	"github.com/suquant/wgrest/models"
	"github.com/suquant/wgrest/restapi/operations/wireguard"

	"encoding/base64"
)

// PeerCreateHandler add/modify peere in device
func PeerCreateHandler(
	params wireguard.PeerCreateParams,
	principal interface{},
) middleware.Responder {
	dev := params.Dev
	wgDev, err := getWGDeviceByName(dev)
	if err != nil {
		switch {
		case IsErrNotFound(err):
			return wireguard.NewPeerCreateNotFound()
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

	peer, err := getWGPeerByPeer(params.Peer)
	if err != nil {
		msg := fmt.Sprintf(err.Error())
		wgrest.Logger.Println(msg)

		return wireguard.NewDeviceGetDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{
				Detail: msg,
			},
		)
	}

	err = addOrRemovePeer(wgDev.Name, peer, false)
	if err != nil {
		msg := fmt.Sprintf(err.Error())
		wgrest.Logger.Println(msg)

		return wireguard.NewDeviceGetDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{
				Detail: msg,
			},
		)
	}

	scheme := "https"
	if params.HTTPRequest.URL.Scheme != "" {
		scheme = params.HTTPRequest.URL.Scheme
	}

	publicKey := peer.PublicKey.String()
	resourcePath := base64.URLEncoding.EncodeToString([]byte(publicKey))

	location := fmt.Sprintf("%s://%s%s%s", scheme, params.HTTPRequest.Host, params.HTTPRequest.RequestURI, resourcePath)
	return wireguard.NewPeerCreateCreated().WithLocation(location)
}
