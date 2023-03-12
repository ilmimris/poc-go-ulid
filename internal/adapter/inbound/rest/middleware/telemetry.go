package middleware

import (
	"fmt"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/ilmimris/poc-go-ulid/shared/constant"
	"go.opentelemetry.io/otel/trace"
)

func NewTelemetry() fiber.Handler {
	return otelfiber.Middleware(
		otelfiber.WithSpanNameFormatter(func(ctx *fiber.Ctx) string {
			return fmt.Sprintf("%s %s", ctx.Method(), ctx.Route().Path)
		}),
		otelfiber.WithServerName(constant.AppName),
	)
}

func SetHeaderTracerID() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		XTraceID := ctx.Get(constant.HeaderXTraceID)
		if XTraceID != "" {
			ctx.Set(constant.HeaderXTraceID, XTraceID)
			return ctx.Next()
		}

		span := trace.SpanFromContext(ctx.UserContext())
		defer span.End()

		XTraceID = span.SpanContext().TraceID().String()
		ctx.Set(constant.HeaderXTraceID, XTraceID)
		return ctx.Next()
	}
}
