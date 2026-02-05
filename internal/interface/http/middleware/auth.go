package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/suquant/wgrest/internal/domain/entity"
)

// BearerAuth creates a middleware that validates Bearer tokens.
func BearerAuth(token string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")

		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(entity.Error{
				Code:    entity.ErrCodeUnauthorized,
				Message: "missing authorization header",
			})
		}

		// Check for Bearer prefix
		const prefix = "Bearer "
		if !strings.HasPrefix(auth, prefix) {
			return c.Status(fiber.StatusUnauthorized).JSON(entity.Error{
				Code:    entity.ErrCodeUnauthorized,
				Message: "invalid authorization format",
			})
		}

		// Extract and validate token
		providedToken := auth[len(prefix):]
		if providedToken != token {
			return c.Status(fiber.StatusUnauthorized).JSON(entity.Error{
				Code:    entity.ErrCodeUnauthorized,
				Message: "invalid token",
			})
		}

		return c.Next()
	}
}
