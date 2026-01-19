FROM golang:1.23 AS build
WORKDIR /src
COPY data-plane/ ./data-plane/
WORKDIR /src/data-plane
RUN go mod download
RUN go build -o /out/data-plane ./cmd/sandbox-runner

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/data-plane /app/data-plane
USER nonroot:nonroot
ENTRYPOINT ["/app/data-plane"]
