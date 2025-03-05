package tracer

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type Interface interface {
	Start(ctx context.Context, spanName string) (context.Context, trace.Span)
	Shutdown(ctx context.Context) error
}

type Config struct {
	TracerExporterEndpoint string
	ServiceName            string
}

type oteltracer struct {
	conf Config
	tp   *sdktrace.TracerProvider
}

func Init(cfg Config) Interface {
	ot := &oteltracer{
		conf: cfg,
	}
	ot.initTracer()
	return ot
}

func (o *oteltracer) initTracer() {
	// Create a Jaeger exporter
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(o.conf.TracerExporterEndpoint)))
	if err != nil {
		panic(err)
	}

	// Create a resource with the service name
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			// Add the service name
			semconv.ServiceNameKey.String(o.conf.ServiceName),
		),
	)
	if err != nil {
		panic(err)
	}

	// Create a trace provider with the Jaeger exporter
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Set the global trace provider
	otel.SetTracerProvider(tp)

	// Set the global propagator to trace context
	otel.SetTextMapPropagator(propagation.TraceContext{})

	o.tp = tp
	log.Printf("OTEL-TRACER-JAEGER: success init with exporter %s", o.conf.TracerExporterEndpoint)
}

func (o *oteltracer) Start(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return o.tp.Tracer("").Start(ctx, spanName)
}

func (o *oteltracer) Shutdown(ctx context.Context) error {
	return o.tp.Shutdown(ctx)
}
