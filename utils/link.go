package utils

import "github.com/vishvananda/netlink"

type WgLink struct {
    *netlink.LinkAttrs
}

func (WgLink) Type() string {
    return "wireguard"
}

func (wg *WgLink) Attrs() *netlink.LinkAttrs {
    return wg.LinkAttrs
}

func (wg *WgLink) Close() error {
    return netlink.LinkDel(wg)
}


