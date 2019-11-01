package wireguard

import (
	"io"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// DeviceEncoder device encoder
type DeviceEncoder struct {
	w io.Writer
}

// NewDeviceEncoder new device encoder
func NewDeviceEncoder(w io.Writer) *DeviceEncoder {
	return &DeviceEncoder{
		w: w,
	}
}

// Encode encode device
func (enc *DeviceEncoder) Encode(v wgtypes.Device) error {
	return nil
}
