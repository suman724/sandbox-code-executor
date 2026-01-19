package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"data-plane/internal/config"
	"data-plane/internal/runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	addr := ":" + getenv("PORT", "8081")
	log.Printf("data-plane starting env=%s addr=%s", cfg.Env, addr)

	serviceName := cfg.OtelService
	if serviceName == "" {
		serviceName = "data-plane"
	}
	shutdown, err := initTelemetry(context.Background(), serviceName, cfg.OtelEndpoint)
	if err != nil {
		log.Fatalf("otel init error: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdown(ctx); err != nil {
			log.Printf("otel shutdown error: %v", err)
		}
	}()

	if err := http.ListenAndServe(addr, runtime.Router()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func initTelemetry(ctx context.Context, serviceName string, endpoint string) (func(context.Context) error, error) {
	otel.SetTextMapPropagator(propagation.TraceContext{})
	res, err := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceNameKey.String(serviceName),
	))
	if err != nil {
		return nil, err
	}
	var tracerProvider *sdktrace.TracerProvider
	if endpoint != "" {
		exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(endpoint), otlptracegrpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		tracerProvider = sdktrace.NewTracerProvider(
			sdktrace.WithResource(res),
			sdktrace.WithBatcher(exporter),
		)
	} else {
		tracerProvider = sdktrace.NewTracerProvider(
			sdktrace.WithResource(res),
		)
	}
	otel.SetTracerProvider(tracerProvider)
	return func(ctx context.Context) error {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	}, nil
}
