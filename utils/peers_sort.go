package utils

import (
	"bytes"
	"fmt"
	"sort"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type sortPeerByPubKey []wgtypes.Peer

func (a sortPeerByPubKey) Len() int { return len(a) }
func (a sortPeerByPubKey) Less(i, j int) bool {
	return bytes.Compare(a[i].PublicKey[:], a[j].PublicKey[:]) > 0
}
func (a sortPeerByPubKey) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type sortPeerByReceiveBytes []wgtypes.Peer

func (a sortPeerByReceiveBytes) Len() int           { return len(a) }
func (a sortPeerByReceiveBytes) Less(i, j int) bool { return a[i].ReceiveBytes < a[j].ReceiveBytes }
func (a sortPeerByReceiveBytes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type sortPeerByReceiveBytesDesc []wgtypes.Peer

func (a sortPeerByReceiveBytesDesc) Len() int           { return len(a) }
func (a sortPeerByReceiveBytesDesc) Less(i, j int) bool { return a[i].ReceiveBytes > a[j].ReceiveBytes }
func (a sortPeerByReceiveBytesDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type sortPeerByTransmitBytes []wgtypes.Peer

func (a sortPeerByTransmitBytes) Len() int           { return len(a) }
func (a sortPeerByTransmitBytes) Less(i, j int) bool { return a[i].ReceiveBytes < a[j].ReceiveBytes }
func (a sortPeerByTransmitBytes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type sortPeerByTransmitBytesDesc []wgtypes.Peer

func (a sortPeerByTransmitBytesDesc) Len() int { return len(a) }
func (a sortPeerByTransmitBytesDesc) Less(i, j int) bool {
	return a[i].ReceiveBytes < a[j].ReceiveBytes
}
func (a sortPeerByTransmitBytesDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type sortPeerByTotalBytes []wgtypes.Peer

func (a sortPeerByTotalBytes) Len() int { return len(a) }
func (a sortPeerByTotalBytes) Less(i, j int) bool {
	return a[i].ReceiveBytes+a[i].TransmitBytes < a[j].ReceiveBytes+a[j].TransmitBytes
}
func (a sortPeerByTotalBytes) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type sortPeerByTotalBytesDesc []wgtypes.Peer

func (a sortPeerByTotalBytesDesc) Len() int { return len(a) }
func (a sortPeerByTotalBytesDesc) Less(i, j int) bool {
	return a[i].ReceiveBytes+a[i].TransmitBytes > a[j].ReceiveBytes+a[j].TransmitBytes
}
func (a sortPeerByTotalBytesDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type sortPeerByLastHandshakeTime []wgtypes.Peer

func (a sortPeerByLastHandshakeTime) Len() int { return len(a) }
func (a sortPeerByLastHandshakeTime) Less(i, j int) bool {
	return a[i].LastHandshakeTime.Before(a[j].LastHandshakeTime)
}
func (a sortPeerByLastHandshakeTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type sortPeerByLastHandshakeTimeDesc []wgtypes.Peer

func (a sortPeerByLastHandshakeTimeDesc) Len() int { return len(a) }
func (a sortPeerByLastHandshakeTimeDesc) Less(i, j int) bool {
	return a[i].LastHandshakeTime.After(a[j].LastHandshakeTime)
}
func (a sortPeerByLastHandshakeTimeDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func SortPeersByField(field string, peers []wgtypes.Peer) error {

	switch field {
	case "pub_key":
		sort.Sort(sortPeerByPubKey(peers))
		break
	case "receive_bytes":
		sort.Sort(sortPeerByReceiveBytes(peers))
		break
	case "-receive_bytes":
		sort.Sort(sortPeerByReceiveBytesDesc(peers))
		break
	case "transmit_bytes":
		sort.Sort(sortPeerByTransmitBytes(peers))
		break
	case "-transmit_bytes":
		sort.Sort(sortPeerByTransmitBytesDesc(peers))
		break
	case "total_bytes":
		sort.Sort(sortPeerByTotalBytes(peers))
		break
	case "-total_bytes":
		sort.Sort(sortPeerByTotalBytesDesc(peers))
		break
	case "last_handshake_time":
		sort.Sort(sortPeerByLastHandshakeTime(peers))
		break
	case "-last_handshake_time":
		sort.Sort(sortPeerByLastHandshakeTimeDesc(peers))
		break
	default:
		return fmt.Errorf("wrong sort field: %s", field)
	}

	return nil
}
