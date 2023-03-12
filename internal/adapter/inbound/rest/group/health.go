package group

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ilmimris/poc-go-ulid/internal/adapter/inbound/rest/handler/healthcheck"
)

func InitHealthCheck(root fiber.Router) {
	healthGroup := root.Group("/health")
	healthGroup.Get("/", healthcheck.Check())
}
