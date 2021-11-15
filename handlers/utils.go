package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/suquant/wgrest/utils"
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
