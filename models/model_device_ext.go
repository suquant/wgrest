package models

import "golang.zx2c4.com/wireguard/wgctrl/wgtypes"

func NewDevice(device *wgtypes.Device) Device {
	var totalReceiveBytes int64
	var totalTransmitBytes int64

	for _, peer := range device.Peers {
		totalReceiveBytes += peer.ReceiveBytes
		totalTransmitBytes += peer.TransmitBytes
	}

	return Device{
		Name:               device.Name,
		ListenPort:         int32(device.ListenPort),
		PublicKey:          device.PublicKey.String(),
		PeersCount:         int32(len(device.Peers)),
		FirewallMark:       int32(device.FirewallMark),
		TotalReceiveBytes:  totalReceiveBytes,
		TotalTransmitBytes: totalTransmitBytes,
	}
}

func (r *DeviceCreateOrUpdateRequest) Apply(conf *wgtypes.Config) error {
	if r.FirewallMark != nil {
		fwMark := int(*r.FirewallMark)
		conf.FirewallMark = &fwMark
	}

	if r.PrivateKey != nil {
		privKey, err := wgtypes.ParseKey(*r.PrivateKey)
		if err != nil {
			return err
		}

		conf.PrivateKey = &privKey
	}

	if r.ListenPort != nil {
		listenPort := int(*r.ListenPort)
		conf.ListenPort = &listenPort
	}

	return nil
}
