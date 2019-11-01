package wireguard

import (
	"io"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type DeviceEncoder struct {
	w io.Writer
}

func NewDeviceEncoder(w io.Writer) *DeviceEncoder {
	return &DeviceEncoder{
		w: w,
	}
}

func (enc *DeviceEncoder) Encode(v wgtypes.Device) error {
	return nil
}
