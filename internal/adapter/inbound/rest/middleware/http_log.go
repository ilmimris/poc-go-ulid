package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ilmimris/poc-go-ulid/shared/constant"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

var contentTypeAllowed = map[string]bool{
	fiber.MIMEApplicationJSON:      true,
	fiber.MIMEApplicationForm:      true,
	fiber.MIMETextPlain:            true,
	fiber.MIMETextPlainCharsetUTF8: true,
}

// NewHTTPLog is middleware fiber for http logging.
// Every you request, the middleware will execute and logging.
// Need custom a http log, because we need to put request_id and tracer_id in fields that can search easier when debug in staging/production.
// If you need a advance log like payload request, you can set a settings "DEBUG" in your config apps.
// Don't recommendation to use "DEBUG" in production, because payload will print in every request that success (if the http_status is failed, the payload keeps print a payload).
// Instead of that, use "INFO" in production. the request it will print but the payload is not print.
func NewHTTPLog() fiber.Handler {
	return func(ctx *fiber.Ctx) (err error) {
		var (
			methodName  = string(ctx.Request().Header.Method())
			url         = ctx.OriginalURL()
			getLevel    = log.GetLevel().String()
			debugStatus = "DEBUG"
			requestID   string
		)

		if err = ctx.Next(); err != nil {
			return
		}

		if ctx.Locals(constant.ContextKeyRequestID.ToString()) != nil {
			requestID = fmt.Sprintf("%v", ctx.Locals(constant.ContextKeyRequestID.ToString()))
		}

		span := trace.SpanFromContext(ctx.UserContext())
		defer span.End()

		httpCode := ctx.Response().StatusCode()
		entry := log.WithFields(log.Fields{
			"request_id":  requestID,
			"trace_id":    span.SpanContext().TraceID().String(),
			"method":      methodName,
			"url":         url,
			"http_status": httpCode,
		})

		if httpCode >= fiber.StatusBadRequest {
			errorResponse := constant.ContextKeyErrorResp.Error(ctx.UserContext())
			entry.WithField("payload", getCleanBodyRequest(ctx)).WithError(errorResponse).Error()
			return
		}

		if httpCode >= fiber.StatusOK && httpCode < fiber.StatusBadRequest && getLevel == debugStatus {
			entry.WithField("payload", getCleanBodyRequest(ctx)).Debug()
			return
		}

		entry.Info()
		return
	}
}

func getCleanBodyRequest(ctx *fiber.Ctx) string {
	var (
		contentType = string(ctx.Request().Header.ContentType())
		body        = ctx.Body()
		bb          = &bytes.Buffer{}
	)

	if contentType == "" {
		return string(body)
	}

	if !contentTypeAllowed[contentType] {
		return constant.ErrNotSupportedDisplayPayload.Error()
	}

	if !contentTypeAllowed[fiber.MIMEApplicationJSON] {
		return string(body)
	}

	err := json.Compact(bb, body)
	if err != nil {
		return string(body)
	}

	return bb.String()
}
