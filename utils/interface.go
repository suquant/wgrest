package utils

import "net"

func GetInterfaceIPs(name string) (addresses []string, err error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		addresses = append(addresses, addr.String())
	}

	return
}
