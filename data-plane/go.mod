module data-plane

go 1.23

require (
	github.com/go-chi/chi/v5 v5.1.0
	github.com/golang-jwt/jwt/v5 v5.2.1
	go.opentelemetry.io/otel v1.26.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.26.0
	go.opentelemetry.io/otel/sdk v1.26.0
	k8s.io/api v0.30.1
	k8s.io/apimachinery v0.30.1
	k8s.io/client-go v0.30.1
)
