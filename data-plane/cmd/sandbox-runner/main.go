package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"data-plane/internal/config"
	"data-plane/internal/execution"
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
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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

	runHandler := runtime.RunHandler{
		Runner: execution.Runner{
			Registry: runtime.DefaultRegistry(),
			Deps:     runtime.DependencyPolicy{},
		},
	}
	sessionRegistry, err := buildSessionRegistry(cfg)
	if err != nil {
		log.Fatalf("session registry error: %v", err)
	}
	sessionRuntime, err := buildSessionRuntime(cfg)
	if err != nil {
		log.Fatalf("session runtime error: %v", err)
	}
	sessionHandler := runtime.SessionHandler{
		Runtime:  sessionRuntime,
		Registry: sessionRegistry,
		Agent:    runtime.NewAgentClient(),
	}
	apiHandler := runtime.RouterWithDependencies(runtime.Dependencies{
		RunHandler:     runHandler,
		SessionHandler: sessionHandler,
	})
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

func buildSessionRuntime(cfg config.Config) (runtime.SessionRuntime, error) {
	switch cfg.SessionRuntime {
	case "k8s":
		restConfig, clientset, err := buildKubeClient()
		if err != nil {
			return nil, err
		}
		return runtime.KubernetesSessionRuntime{
			Client:         clientset,
			Config:         restConfig,
			Namespace:      cfg.RuntimeNamespace,
			RuntimeClass:   cfg.RuntimeClass,
			Image:          cfg.SessionImage,
			PythonImage:    cfg.SessionImagePython,
			NodeImage:      cfg.SessionImageNode,
			Env:            cfg.Env,
			AgentAddr:      getenv("SESSION_AGENT_ADDR", ":9000"),
			AgentAuthMode:  cfg.AgentAuthMode,
			AgentAuthToken: cfg.AgentAuthToken,
		}, nil
	default:
		return runtime.NewLocalSessionRuntime(), nil
	}
}

func buildSessionRegistry(cfg config.Config) (runtime.SessionRegistry, error) {
	switch cfg.SessionRegistry {
	case "file":
		return runtime.NewFileSessionRegistry(cfg.SessionRegistryPath)
	default:
		return runtime.NewInMemorySessionRegistry(), nil
	}
}

func buildKubeClient() (*rest.Config, kubernetes.Interface, error) {
	if cfg, err := rest.InClusterConfig(); err == nil {
		clientset, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			return nil, nil, err
		}
		return cfg, clientset, nil
	}
	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, nil, err
		}
		clientset, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			return nil, nil, err
		}
		return cfg, clientset, nil
	}
	return nil, nil, errors.New("missing kubernetes config")
}
