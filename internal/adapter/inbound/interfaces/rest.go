package interfaces

import "github.com/gofiber/fiber/v2"

type Rest interface {
	Serve()
	GetRouter() *fiber.App
	Shutdown()
}
