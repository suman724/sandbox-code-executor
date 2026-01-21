FROM golang:1.23 AS builder

WORKDIR /src/session-agent

COPY session-agent/go.mod session-agent/go.sum ./
COPY shared/go.mod ./shared/
RUN go mod download

COPY session-agent/ ./
COPY shared/ /src/shared/
RUN CGO_ENABLED=0 go build -o /out/session-agent ./cmd/session-agent

FROM node:20-slim

WORKDIR /app

COPY --from=builder /out/session-agent /app/session-agent

ENV SESSION_AGENT_ADDR=:9000
ENV SESSION_AGENT_AUTH_BYPASS=false

EXPOSE 9000

CMD ["/app/session-agent"]
