package tracing

import (
	"github.com/ilmimris/poc-go-ulid/config"
	"github.com/ilmimris/poc-go-ulid/pkg/tracing"
	"github.com/ilmimris/poc-go-ulid/shared/constant"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func InitTracing(cfg config.TracerConfig) {
	tp, err := tracing.GetOtelProvider(cfg.Tracer, constant.AppName, cfg)
	if err != nil {
		panic(err)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}, b3.New()))
}
