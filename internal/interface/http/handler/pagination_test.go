package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSetLinkHeader_FirstPage(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		setLinkHeader(c, 0, 10, 100)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	link := resp.Header.Get("Link")
	assert.Contains(t, link, "rel=\"next\"")
	assert.Contains(t, link, "rel=\"last\"")
	assert.Contains(t, link, "rel=\"first\"")
}

func TestSetLinkHeader_MiddlePage(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		setLinkHeader(c, 5, 10, 100)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	link := resp.Header.Get("Link")
	assert.Contains(t, link, "rel=\"next\"")
	assert.Contains(t, link, "rel=\"prev\"")
	assert.Contains(t, link, "rel=\"first\"")
	assert.Contains(t, link, "rel=\"last\"")
}

func TestSetLinkHeader_LastPage(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		setLinkHeader(c, 9, 10, 100) // Page 9 is the last page for 100 items with 10 per page
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	link := resp.Header.Get("Link")
	assert.NotContains(t, link, "rel=\"next\"")
	assert.Contains(t, link, "rel=\"prev\"")
	assert.Contains(t, link, "rel=\"first\"")
}

func TestSetLinkHeader_SinglePage(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		setLinkHeader(c, 0, 10, 5) // Only 5 items, fits on one page
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	link := resp.Header.Get("Link")
	// On single page, no next/prev needed
	assert.NotContains(t, link, "rel=\"next\"")
	assert.NotContains(t, link, "rel=\"prev\"")
}

func TestSetLinkHeader_EmptyResults(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		setLinkHeader(c, 0, 10, 0)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	link := resp.Header.Get("Link")
	assert.Empty(t, link)
}

func TestSetLinkHeader_SecondPage(t *testing.T) {
	app := fiber.New()
	app.Get("/devices", func(c *fiber.Ctx) error {
		setLinkHeader(c, 1, 10, 50)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/devices?sort=name", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	link := resp.Header.Get("Link")
	assert.Contains(t, link, "rel=\"next\"")
	assert.Contains(t, link, "rel=\"prev\"")
}

func TestSetLinkHeader_LargeDataset(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		setLinkHeader(c, 50, 25, 10000)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	link := resp.Header.Get("Link")
	assert.Contains(t, link, "page=0")   // first
	assert.Contains(t, link, "page=399") // last page for 10000 items / 25 per page
	assert.Contains(t, link, "page=49")  // prev
	assert.Contains(t, link, "page=51")  // next
}
