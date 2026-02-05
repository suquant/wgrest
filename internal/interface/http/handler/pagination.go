package handler

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// setLinkHeader sets the Link header for pagination following GitHub's pattern.
func setLinkHeader(c *fiber.Ctx, page, perPage, total int) {
	if total == 0 {
		return
	}

	totalPages := (total + perPage - 1) / perPage
	baseURL := c.OriginalURL()

	// Parse and remove old pagination params
	u, err := url.Parse(baseURL)
	if err != nil {
		return
	}

	query := u.Query()
	query.Del("page")
	query.Del("per_page")

	var links []string

	// First page
	firstQuery := query
	firstQuery.Set("page", "0")
	firstQuery.Set("per_page", strconv.Itoa(perPage))
	u.RawQuery = firstQuery.Encode()
	links = append(links, fmt.Sprintf("<%s>; rel=\"first\"", u.String()))

	// Last page
	lastQuery := query
	lastQuery.Set("page", strconv.Itoa(totalPages-1))
	lastQuery.Set("per_page", strconv.Itoa(perPage))
	u.RawQuery = lastQuery.Encode()
	links = append(links, fmt.Sprintf("<%s>; rel=\"last\"", u.String()))

	// Previous page
	if page > 0 {
		prevQuery := query
		prevQuery.Set("page", strconv.Itoa(page-1))
		prevQuery.Set("per_page", strconv.Itoa(perPage))
		u.RawQuery = prevQuery.Encode()
		links = append(links, fmt.Sprintf("<%s>; rel=\"prev\"", u.String()))
	}

	// Next page
	if page < totalPages-1 {
		nextQuery := query
		nextQuery.Set("page", strconv.Itoa(page+1))
		nextQuery.Set("per_page", strconv.Itoa(perPage))
		u.RawQuery = nextQuery.Encode()
		links = append(links, fmt.Sprintf("<%s>; rel=\"next\"", u.String()))
	}

	// Set Link header
	linkHeader := ""
	for i, link := range links {
		if i > 0 {
			linkHeader += ", "
		}
		linkHeader += link
	}

	if linkHeader != "" {
		c.Set("Link", linkHeader)
	}
}
