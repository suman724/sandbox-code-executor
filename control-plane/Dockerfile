FROM golang:1.23 AS build
WORKDIR /src
COPY control-plane/go.mod ./control-plane/
WORKDIR /src/control-plane
ENV GOFLAGS="-mod=mod"
RUN go mod download
COPY control-plane/ .
RUN go build -o /out/control-plane ./cmd/control-plane

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/control-plane /app/control-plane
USER nonroot:nonroot
ENTRYPOINT ["/app/control-plane"]
