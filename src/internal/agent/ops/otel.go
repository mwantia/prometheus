package ops

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/mwantia/queueverse/internal/config"
	"github.com/mwantia/queueverse/internal/registry"
	"github.com/mwantia/queueverse/pkg/log"
)

type OpenTelemetry struct {
	Operation

	provider *sdktrace.TracerProvider
	exporter *otlptrace.Exporter

	Log log.Logger
}

func (m *OpenTelemetry) Create(cfg *config.Config, reg *registry.PluginRegistry) (Cleanup, error) {
	m.Log = *log.New("telemetry")

	prop := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(prop)

	exp, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint(cfg.Telemetry.Endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	m.exporter = exp
	m.provider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp,
			sdktrace.WithBatchTimeout(time.Second*5),
			sdktrace.WithMaxExportBatchSize(10),
		),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.Telemetry.ServiceName),
			semconv.TelemetrySDKNameKey.String("opentelemetry"),
			semconv.TelemetrySDKLanguageKey.String("go"),
		)),
	)

	return func(ctx context.Context) error {
		m.provider.Shutdown(ctx)
		m.exporter.Shutdown(ctx)
		return nil
	}, nil
}

func (m *OpenTelemetry) Serve(ctx context.Context) error {
	otel.SetTracerProvider(m.provider)
	m.exporter.Start(ctx)
	return nil
}
