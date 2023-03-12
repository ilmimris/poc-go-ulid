package middleware

import (
	"crypto/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/ilmimris/poc-go-ulid/shared/constant"
	"github.com/oklog/ulid"
)

func NewRequestId() fiber.Handler {
	return requestid.New(requestid.Config{
		Generator: func() string {
			return ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String()
		},
		ContextKey: constant.ContextKeyRequestID.ToString(),
	})
}
