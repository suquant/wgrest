package utils

import (
	externalip "github.com/glendc/go-external-ip"
)

func GetExternalIP() (string, error) {
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	if err != nil {
		return "", err
	}

	return ip.String(), nil
}
