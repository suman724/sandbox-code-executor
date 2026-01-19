.PHONY: build test lint format run run-local \
	build-control-plane build-data-plane \
	test-control-plane test-data-plane \
	run-control-plane run-data-plane

build: build-control-plane build-data-plane

build-control-plane:
	@cd control-plane && go build ./...

build-data-plane:
	@cd data-plane && go build ./...

test: test-control-plane test-data-plane

test-control-plane:
	@cd control-plane && go test ./...

test-data-plane:
	@cd data-plane && go test ./...

lint:
	@golangci-lint run ./...

format:
	@gofmt -w .

run: run-control-plane run-data-plane

run-local:
	@ENV=dev DATA_PLANE_URL=http://localhost:8081 DATABASE_DRIVER=sqlite DATABASE_URL='file:control-plane.db?cache=shared&mode=rwc' AUTHZ_BYPASS=true \
		$(MAKE) run-control-plane &
	@ENV=dev RUNTIME_NAMESPACE=default RUNTIME_CLASS=gvisor AUTHZ_BYPASS=true \
		$(MAKE) run-data-plane &
	@wait

run-control-plane:
	@cd control-plane && go run ./cmd/control-plane

run-data-plane:
	@cd data-plane && go run ./cmd/sandbox-runner
