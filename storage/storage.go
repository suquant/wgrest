package storage

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type StoreDeviceOptions struct {
	DNSServers []string
	AllowedIPs []string
	Host       string
}

func (o *StoreDeviceOptions) Dump(w io.Writer) error {
	fmt.Fprintf(w, "DNS = %s\n", strings.Join(o.DNSServers, ", "))
	fmt.Fprintf(w, "AllowedIPs = %s\n", strings.Join(o.AllowedIPs, ", "))
	fmt.Fprintf(w, "Host = %s\n", o.Host)

	return nil
}

func (o *StoreDeviceOptions) Restore(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		terms := strings.SplitN(line, "=", 2)
		if len(terms) != 2 {
			return fmt.Errorf("failed to parse line: %s", line)
		}

		left := strings.ToLower(strings.TrimSpace(terms[0]))
		switch left {
		case "dns":
			right := strings.Split(terms[1], ",")
			dnsServers := make([]string, len(right))
			for i, v := range right {
				dnsServers[i] = strings.TrimSpace(v)
			}
			o.DNSServers = dnsServers
			break
		case "allowedips":
			right := strings.Split(terms[1], ",")
			for _, v := range right {
				o.AllowedIPs = append(o.AllowedIPs, strings.TrimSpace(v))
			}
			break
		case "host":
			o.Host = strings.TrimSpace(terms[1])
			break
		default:
			break
		}
	}

	return nil
}

type StorePeerOptions struct {
	PrivateKey string
}

func (o *StorePeerOptions) Dump(w io.Writer) error {
	fmt.Fprintf(w, "PrivateKey = %s\n", o.PrivateKey)

	return nil
}

func (o *StorePeerOptions) Restore(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		terms := strings.SplitN(line, "=", 2)
		if len(terms) != 2 {
			return fmt.Errorf("failed to parse line: %s", line)
		}

		left := strings.ToLower(strings.TrimSpace(terms[0]))
		switch left {
		case "privatekey":
			o.PrivateKey = strings.TrimSpace(terms[1])
			break
		default:
			return fmt.Errorf("invalid option: %s", left)
		}
	}

	return nil
}

type Storage interface {
	WriteDeviceOptions(name string, options StoreDeviceOptions) error
	WritePeerOptions(pubKey wgtypes.Key, options StorePeerOptions) error

	ReadDeviceOptions(name string) (*StoreDeviceOptions, error)
	ReadPeerOptions(pubKey wgtypes.Key) (*StorePeerOptions, error)
}
