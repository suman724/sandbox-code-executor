# Quickstart: Sandboxed Code Execution Service

## Prerequisites

- Go 1.23 installed
- Docker installed for local container builds
- Access to an internal Kubernetes cluster
- Tenant identity provider and audit storage configured

## Local Development

1. Run `make build-control-plane` to compile the control plane.
2. Run `make build-data-plane` to compile the data plane.
3. Run `make test-control-plane` and `make test-data-plane` for unit and
   integration tests.
4. Run `make run-control-plane` and `make run-data-plane` to start each service
   locally.

### SQLite (non-production only)

Set these environment variables to use SQLite for local testing:

```bash
export DATABASE_DRIVER=sqlite
export DATABASE_URL='file:control-plane.db?cache=shared&mode=rwc'
```

### MCP server (separate port)

The MCP tool interface runs on its own port. Set `MCP_ADDR` before starting
the control plane, for example:

```bash
export MCP_ADDR=':8090'
```

## Container Builds

- Build control plane: `docker build -f control-plane/Dockerfile -t control-plane:dev .`
- Build data plane: `docker build -f data-plane/Dockerfile -t data-plane:dev .`

## Authorization Bypass (Non-Production Only)

- Use the `AUTHZ_BYPASS=true` environment flag for local or test environments.
- When enabled, both control plane and data plane log audit events indicating
  authorization bypass was active.

## Basic Workflow

1. Create a policy for a tenant (control plane).
2. Submit a one-shot job; the control plane validates policy and calls the data
   plane to create a sandbox and execute the run.
3. Create a session; the control plane asks the data plane to allocate a
   sandbox for the session TTL, then run steps inside it.
4. Review audit events for all executions (control plane).

## Validation

- Run unit tests: `make test-control-plane` and `make test-data-plane`.
- Run integration tests: `go test ./control-plane/tests/integration` and
  `go test ./data-plane/tests/integration`.
- Run k6 checks: `k6 run control-plane/tests/integration/k6/jobs.js` and
  `k6 run control-plane/tests/integration/k6/sessions.js`.
- Start locally with defaults: `make run-local` (starts REST on `:8080` and MCP on `:8090`).
- Validate MCP health: `curl http://localhost:8090/healthz`.
- Validate Prometheus metrics (control plane): `curl http://localhost:8080/metrics`.
- Validate Prometheus metrics (data plane): `curl http://localhost:8081/metrics`.

## Deployment Notes

- Deploy control plane and data plane as separate services with independent
  scaling and rollout schedules.
- Service mode requires explicit policy enablement and proxy access.
- Network egress is denied by default; allowlist as needed per policy.
