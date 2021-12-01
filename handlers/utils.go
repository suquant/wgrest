package handlers

import (
	"encoding/base64"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/suquant/wgrest/models"
	"github.com/suquant/wgrest/utils"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"net/http"
	"strconv"
)

func getPaginator(ctx echo.Context, nums int) (*utils.Paginator, error) {
	perPageParam := ctx.QueryParam("per_page")
	var perPage int = 100

	if perPageParam != "" {
		parsedPerPage, err := strconv.Atoi(perPageParam)
		if err != nil {
			return nil, &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  "failed to parse per_page param",
				Internal: err,
			}
		}

		perPage = parsedPerPage
	}

	return utils.NewPaginator(ctx.Request(), perPage, nums), nil
}

func parseUrlSafeKey(encodedKey string) (wgtypes.Key, error) {
	decodedKey, err := base64.URLEncoding.DecodeString(encodedKey)
	if err != nil {
		return wgtypes.Key{}, fmt.Errorf("failed to parse key: %s", err)
	}

	if len(decodedKey) != wgtypes.KeyLen {
		return wgtypes.Key{}, fmt.Errorf("failed to parse key: wrong length")
	}
	var key wgtypes.Key
	copy(key[:32], decodedKey[:])

	return key, nil
}

func applyNetworks(device *models.Device) error {
	addresses, err := utils.GetInterfaceIPs(device.Name)
	if err != nil {
		return err
	}

	device.Networks = addresses
	return nil
}
