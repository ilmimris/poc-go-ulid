package group

import (
	"github.com/gofiber/fiber/v2"
	handlerv1 "github.com/ilmimris/poc-go-ulid/internal/adapter/inbound/rest/handler/v1"
)

func InitV1Group(root fiber.Router, h handlerv1.Handler) fiber.Router {
	v1Group := root.Group("/v1")

	return v1Group
}
