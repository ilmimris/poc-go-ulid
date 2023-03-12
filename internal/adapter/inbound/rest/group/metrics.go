package group

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/ilmimris/poc-go-ulid/shared/constant"
)

func InitMetrics(app *fiber.App) {
	prometheus := fiberprometheus.New(constant.AppName)
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)
}
