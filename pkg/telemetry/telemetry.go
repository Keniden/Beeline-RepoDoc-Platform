package telemetry

import (
	"context"
	"fmt"
	"net/url"

	"github.com/beeline/repodoc/configs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func InitTracer(ctx context.Context, cfg configs.TelemetryConfig) (*sdktrace.TracerProvider, error) {
	serviceName := cfg.ServiceName
	if serviceName == "" {
		serviceName = "repodoc"
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(attribute.String("service.name", serviceName)),
	)
	if err != nil {
		return nil, fmt.Errorf("resource: %w", err)
	}

	if cfg.OTelEndpoint == "" {
		return sdktrace.NewTracerProvider(sdktrace.WithResource(res)), nil
	}

	exporterOpts := []otlptracehttp.Option{}
	if parsed, err := url.Parse(cfg.OTelEndpoint); err == nil && parsed.Host != "" {
		exporterOpts = append(exporterOpts, otlptracehttp.WithEndpoint(parsed.Host))
		if parsed.Scheme == "http" {
			exporterOpts = append(exporterOpts, otlptracehttp.WithInsecure())
		}
	} else {
		exporterOpts = append(exporterOpts, otlptracehttp.WithEndpoint(cfg.OTelEndpoint))
	}

	exporter, err := otlptracehttp.New(ctx, exporterOpts...)
	if err != nil {
		return nil, fmt.Errorf("otlp exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	return tp, nil
}

func Shutdown(ctx context.Context, tp *sdktrace.TracerProvider) {
	if tp == nil {
		return
	}
	_ = tp.Shutdown(ctx)
}
