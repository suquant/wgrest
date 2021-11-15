package utils

import (
	"github.com/suquant/wgrest/models"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"strings"
)

func FilterPeersByQuery(q string, peers []wgtypes.Peer) []wgtypes.Peer {
	var filteredPeers []wgtypes.Peer
	for _, peer := range peers {
		var terms = []string{
			peer.PublicKey.String(),
		}

		if peer.PresharedKey != models.EmptyKey {
			terms = append(terms, peer.PresharedKey.String())
		}

		for _, v := range peer.AllowedIPs {
			terms = append(terms, v.String())
		}

		if peer.Endpoint != nil {
			terms = append(terms, peer.Endpoint.String())
		}

		if strings.Contains(strings.Join(terms, " "), q) {
			filteredPeers = append(filteredPeers, peer)
		}
	}

	return filteredPeers
}
