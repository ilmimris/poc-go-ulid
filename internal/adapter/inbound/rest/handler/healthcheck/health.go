package healthcheck

import "github.com/gofiber/fiber/v2"

func Check() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	}
}
