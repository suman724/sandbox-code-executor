.PHONY: build test lint format run run-local run-local-control-plane run-local-data-plane run-local-session-agent \
	build-control-plane build-data-plane build-session-agent build-runtime-images build-runtime-python build-runtime-node \
	test-control-plane test-data-plane test-session-agent \
	run-control-plane run-data-plane run-session-agent

build: build-control-plane build-data-plane build-session-agent

build-control-plane:
	@cd control-plane && go build ./...

build-data-plane:
	@cd data-plane && go build ./...

build-session-agent:
	@cd session-agent && go build ./...

build-runtime-images: build-runtime-python build-runtime-node

build-runtime-python:
	@docker build -f deploy/runtime/python/Dockerfile -t runtime-python:dev .

build-runtime-node:
	@docker build -f deploy/runtime/node/Dockerfile -t runtime-node:dev .

test: test-control-plane test-data-plane test-session-agent

test-control-plane:
	@cd control-plane && go test ./...

test-data-plane:
	@cd data-plane && go test ./...

test-session-agent:
	@cd session-agent && go test ./...

lint:
	@golangci-lint run ./...

format:
	@gofmt -w .

run: run-control-plane run-data-plane run-session-agent

run-local: run-local-control-plane run-local-data-plane run-local-session-agent

run-local-control-plane:
	@ENV=dev DATA_PLANE_URL=http://localhost:8081 DATABASE_DRIVER=sqlite DATABASE_URL='file:control-plane.db?cache=shared&mode=rwc' MCP_ADDR=:8090 AUTHZ_BYPASS=true \
		$(MAKE) run-control-plane

run-local-data-plane:
	@ENV=dev RUNTIME_NAMESPACE=default RUNTIME_CLASS=gvisor SESSION_RUNTIME_BACKEND=local SESSION_REGISTRY_BACKEND=memory SESSION_REGISTRY_PATH=/tmp/session-registry.json SESSION_AGENT_AUTH_MODE=bypass SESSION_AGENT_PREFER=true AUTHZ_BYPASS=true \
		$(MAKE) run-data-plane

run-local-session-agent:
	@ENV=dev SESSION_AGENT_ADDR=:9000 SESSION_AGENT_AUTH_BYPASS=true \
		$(MAKE) run-session-agent

run-control-plane:
	@cd control-plane && go run ./cmd/control-plane

run-data-plane:
	@cd data-plane && go run ./cmd/sandbox-runner

run-session-agent:
	@cd session-agent && go run ./cmd/session-agent
