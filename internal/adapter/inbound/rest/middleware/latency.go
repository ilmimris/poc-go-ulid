package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ilmimris/poc-go-ulid/shared/constant"
)

func NewLatency() func(ctx *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		ctx = context.WithValue(ctx, constant.ContextKeyLatency, time.Now().UTC().Format(time.RFC3339Nano))
		c.SetUserContext(ctx)
		return c.Next()
	}
}
