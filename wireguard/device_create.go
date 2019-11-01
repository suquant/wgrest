package wireguard

import (
	"fmt"
	"net/http"
	"os"

	"github.com/suquant/wgrest/storage"

	"github.com/suquant/wgrest"

	"github.com/vishvananda/netlink"

	"github.com/suquant/wgrest/models"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"github.com/go-openapi/runtime/middleware"
	"github.com/suquant/wgrest/restapi/operations/wireguard"
	"golang.zx2c4.com/wireguard/wgctrl"
)

// DeviceCreateHandler wireguard device create handler
func DeviceCreateHandler(
	params wireguard.DeviceCreateParams,
	principal interface{},
) middleware.Responder {
	name := *params.Device.Name
	wgConfPath := getWgConfPath(name)

	// check if interface already exist
	_, err := os.Stat(wgConfPath)
	if err != nil && !os.IsNotExist(err) {
		msg := fmt.Sprintf("os err: %s", err.Error())
		wgrest.Logger.Println(msg)

		return wireguard.NewDeviceCreateDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{Detail: msg},
		)
	} else if err == nil {
		msg := fmt.Sprintf("file %s exists", wgConfPath)
		return wireguard.NewDeviceCreateConflict().WithPayload(
			&models.Error{Detail: msg},
		)
	}

	client, err := wgctrl.New()
	if err != nil {
		msg := fmt.Sprintf("wgctrl err: %s", err.Error())
		wgrest.Logger.Println(msg)

		return wireguard.NewDeviceCreateDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{Detail: msg},
		)
	}

	privateKey, err := wgtypes.ParseKey(*params.Device.PrivateKey)
	if err != nil {
		msg := fmt.Sprintf("wgctrl err: %s", err.Error())
		wgrest.Logger.Println(msg)

		return wireguard.NewDeviceCreateDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{Detail: msg},
		)
	}

	listenPort := int(*params.Device.ListenPort)

	la := netlink.NewLinkAttrs()
	la.Name = name

	wgDev := &netlink.GenericLink{
		LinkAttrs: la,
		LinkType:  "wireguard",
	}

	err = netlink.LinkAdd(wgDev)
	if err != nil {
		msg := fmt.Sprintf("netlink err: %s", err.Error())
		wgrest.Logger.Println(msg)

		return wireguard.NewDeviceCreateDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{Detail: msg},
		)
	}

	net, err := netlink.ParseAddr(*params.Device.Network)
	if err != nil {
		msg := fmt.Sprintf("netlink err: %s", err.Error())
		wgrest.Logger.Println(msg)

		return wireguard.NewDeviceCreateDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{Detail: msg},
		)
	}

	err = netlink.AddrAdd(wgDev, net)
	if err != nil {
		msg := fmt.Sprintf("netlink err: %s", err.Error())
		wgrest.Logger.Println(msg)

		return wireguard.NewDeviceCreateDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{Detail: msg},
		)
	}

	cfg := wgtypes.Config{
		PrivateKey: &privateKey,
		ListenPort: &listenPort,
	}

	err = client.ConfigureDevice(name, cfg)
	if err != nil {
		msg := fmt.Sprintf("wgctrl err: %s", err.Error())
		wgrest.Logger.Println(msg)

		return wireguard.NewDeviceCreateDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{Detail: msg},
		)
	}

	err = netlink.LinkSetUp(wgDev)
	if err != nil {
		msg := fmt.Sprintf("netlink err: %s", err.Error())
		wgrest.Logger.Println(msg)

		return wireguard.NewDeviceCreateDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{Detail: msg},
		)
	}

	st := storage.NewStorage(storage.DiskStorage)
	rwc, err := st.Open(wgConfPath)
	if err != nil {
		msg := fmt.Sprintf("storage err: %s", err.Error())
		wgrest.Logger.Println(msg)

		return wireguard.NewDeviceCreateDefault(http.StatusInternalServerError).WithPayload(
			&models.Error{Detail: msg},
		)
	}
	defer rwc.Close()

	fmt.Fprintf(rwc, "[Interface]\n")
	fmt.Fprintf(rwc, "Address = %s\n", *params.Device.Network)
	fmt.Fprintf(rwc, "ListenPort = %v\n", *params.Device.ListenPort)
	fmt.Fprintf(rwc, "PrivateKey = %s\n", *params.Device.PrivateKey)
	fmt.Fprintf(rwc, "SaveConfig = true\n\n")

	scheme := "https"
	if params.HTTPRequest.URL.Scheme != "" {
		scheme = params.HTTPRequest.URL.Scheme
	}

	location := fmt.Sprintf("%s://%s%s%s", scheme, params.HTTPRequest.Host, params.HTTPRequest.RequestURI, name)
	return wireguard.NewDeviceCreateCreated().WithLocation(location)
}
