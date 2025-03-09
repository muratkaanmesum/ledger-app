package monitoring

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
)

func InitTracer() (*trace.TracerProvider, error) {
	ctx := context.Background()

	conn, err := grpc.NewClient("localhost:4317")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Create OTLP gRPC exporter
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	// Create Tracer Provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("echo-app"),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}
