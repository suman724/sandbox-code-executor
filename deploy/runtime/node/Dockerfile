FROM golang:1.23 AS builder

WORKDIR /src/session-agent

COPY session-agent/go.mod ./
RUN go mod download

COPY session-agent/ ./
RUN CGO_ENABLED=0 go build -o /out/session-agent ./cmd/session-agent

FROM node:20-slim

WORKDIR /app

COPY --from=builder /out/session-agent /app/session-agent

ENV SESSION_AGENT_ADDR=:9000
ENV SESSION_AGENT_AUTH_BYPASS=false

EXPOSE 9000

CMD ["/app/session-agent"]
