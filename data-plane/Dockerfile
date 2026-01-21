FROM golang:1.23 AS build
WORKDIR /src
COPY data-plane/go.mod ./data-plane/
COPY shared/go.mod shared/go.sum ./shared/
WORKDIR /src/data-plane
ENV GOFLAGS="-mod=mod"
RUN go mod download
COPY data-plane/ .
COPY shared/ /src/shared/
RUN go build -o /out/data-plane ./cmd/sandbox-runner

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/data-plane /app/data-plane
USER nonroot:nonroot
ENTRYPOINT ["/app/data-plane"]
