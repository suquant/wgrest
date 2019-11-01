package wireguard

import (
	"fmt"
	"io"

	"gopkg.in/ini.v1"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type DeviceDecoder struct {
	r io.Reader
}

func NewDeviceDecoder(r io.Reader) *DeviceDecoder {
	return &DeviceDecoder{
		r: r,
	}
}

func (dec *DeviceDecoder) Decode(v *wgtypes.Device) error {
	cfg, err := ini.Load(dec.r)
	if err != nil {
		return fmt.Errorf("load ini config error: %w", err)
	}

	_, err = cfg.GetSection("interface")
	if err != nil {
		return fmt.Errorf("section interface not found: %w", err)
	}

	return nil
}
