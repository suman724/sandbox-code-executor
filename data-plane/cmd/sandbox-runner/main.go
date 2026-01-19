package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"data-plane/internal/config"
	"data-plane/internal/runtime"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
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
	telemetry, err := initTelemetry(context.Background(), serviceName, cfg.OtelEndpoint)
	if err != nil {
		log.Fatalf("otel init error: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := telemetry.Shutdown(ctx); err != nil {
			log.Printf("otel shutdown error: %v", err)
		}
	}()

	apiHandler := runtime.Router()
	router := chi.NewRouter()
	if telemetry.MetricsHandler != nil {
		router.Handle("/metrics", telemetry.MetricsHandler)
	}
	router.Mount("/", apiHandler)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

type telemetryInit struct {
	Shutdown       func(context.Context) error
	MetricsHandler http.Handler
}

func initTelemetry(ctx context.Context, serviceName string, endpoint string) (telemetryInit, error) {
	otel.SetTextMapPropagator(propagation.TraceContext{})
	res, err := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceNameKey.String(serviceName),
	))
	if err != nil {
		return telemetryInit{}, err
	}
	var tracerProvider *sdktrace.TracerProvider
	if endpoint != "" {
		exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(endpoint), otlptracegrpc.WithInsecure())
		if err != nil {
			return telemetryInit{}, err
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
	metricsExporter, err := prometheus.New()
	if err != nil {
		return telemetryInit{}, err
	}
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(metricsExporter),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(meterProvider)
	return telemetryInit{
		Shutdown: func(ctx context.Context) error {
			if err := tracerProvider.Shutdown(ctx); err != nil {
				return err
			}
			return meterProvider.Shutdown(ctx)
		},
		MetricsHandler: promhttp.Handler(),
	}, nil
}
