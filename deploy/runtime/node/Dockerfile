FROM node:20-slim

WORKDIR /app

COPY session-agent /app/session-agent

WORKDIR /app/session-agent
RUN go mod download
RUN go build -o /app/session-agent/bin/session-agent ./cmd/session-agent

ENV SESSION_AGENT_ADDR=:9000
ENV SESSION_AGENT_AUTH_BYPASS=false

EXPOSE 9000

CMD ["/app/session-agent/bin/session-agent"]
