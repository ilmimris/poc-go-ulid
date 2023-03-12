package rest

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/ilmimris/poc-go-ulid/internal/adapter/inbound/rest/group"
	handlerv1 "github.com/ilmimris/poc-go-ulid/internal/adapter/inbound/rest/handler/v1"
	"github.com/ilmimris/poc-go-ulid/internal/adapter/inbound/rest/middleware"
	"github.com/ilmimris/poc-go-ulid/internal/adapter/inbound/rest/registry"

	log "github.com/sirupsen/logrus"
)

type Options struct {
	// Port is the port that the REST server will listen on.
	Port         string
	BodyLimit    int
	ReadTimeout  int
	WriteTimeout int
}

type Rest struct {
	// contains filtered or unexported fields
	router  *fiber.App
	options Options
}

func NewRest(o *Options, serviceRegistry *registry.ServiceRegistry) *Rest {
	app := fiber.New(fiber.Config{
		BodyLimit:        o.BodyLimit,
		DisableKeepalive: true,
		ReadTimeout:      time.Duration(o.ReadTimeout) * time.Second,
		WriteTimeout:     time.Duration(o.WriteTimeout) * time.Second,
	})

	r := &Rest{
		router:  app,
		options: *o,
	}

	v1 := handlerv1.New(handlerv1.OptHandler{
		ServiceRegistry: serviceRegistry,
	})

	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(csrf.New())
	app.Use(middleware.NewRequestId())
	app.Use(middleware.NewTelemetry(), middleware.SetHeaderTracerID())
	app.Use(middleware.NewLatency())
	app.Use(middleware.NewHTTPLog())
	app.Use(idempotency.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	group.InitRouterV1(r, *v1)

	return r

}

func (r *Rest) Serve() {
	log.Infof("REST server is listening on port %s", r.options.Port)
	if err := r.router.Listen(":" + r.options.Port); err != nil {
		log.Fatalf("REST server failed to listen on port %s: %v", r.options.Port, err)
	}
}

func (r *Rest) Shutdown() {
	log.Infof("REST server is shutting down")
	if err := r.router.Shutdown(); err != nil {
		log.Fatalf("REST server failed to shutdown: %v", err)
	}
}

func (r *Rest) GetRouter() *fiber.App {
	return r.router
}
