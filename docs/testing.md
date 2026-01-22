# Testing Guide (Local + Kubernetes)

This guide provides step-by-step commands and helper scripts to test the control plane, data plane, and session-agent locally and on a Kubernetes cluster.

## Prerequisites

- Go 1.23
- Docker (for runtime images)
- `kubectl` with cluster access (for Kubernetes testing)

## Local Testing

### 1) Start all services

Use the Makefile targets to run the stack locally:

```sh
make run-local
```

This starts:
- Control plane: `http://localhost:8080` (health: `/healthz`)
- Data plane: `http://localhost:8081` (health: `/healthz`)
- Session-agent: `http://localhost:9000` (health: `/v1/health`)

### 2) Run the session flow script

The script creates a session and submits multiple steps to verify REPL state.

```sh
BASE_URL=http://localhost:8080 RUNTIME=python bash scripts/run-session-step.sh
```

You can override defaults with:
- `BASE_URL` (control-plane URL)
- `TENANT_ID`, `AGENT_ID`, `POLICY_ID`
- `RUNTIME` (`python` or `node`)

### 3) Optional smoke tests

```sh
curl -s http://localhost:8080/healthz
curl -s http://localhost:8081/healthz
curl -s http://localhost:9000/v1/health
```

### 4) Local test helper script

Run this to check health endpoints and execute the session flow:

```sh
python3 scripts/test-local-stack.py
```

---

## Kubernetes Testing

### 1) Build and push runtime images (session-agent baked in)

```sh
docker build -f deploy/runtime/python/Dockerfile -t <registry>/runtime-python:dev .
docker build -f deploy/runtime/node/Dockerfile -t <registry>/runtime-node:dev .
docker push <registry>/runtime-python:dev
docker push <registry>/runtime-node:dev
```

### 2) Deploy control-plane and data-plane

Use the manifests or Helm charts under `deploy/`:

```sh
# Example (adjust to your environment):
kubectl apply -f deploy/control-plane
kubectl apply -f deploy/data-plane
```

### 3) Configure data-plane for Kubernetes sessions

Ensure these env vars are set on the data-plane deployment:

- `SESSION_RUNTIME_BACKEND=k8s`
- `RUNTIME_NAMESPACE=<namespace>`
- `RUNTIME_CLASS=<runtimeClass>` (optional)
- `SESSION_AGENT_AUTH_MODE=enforced` (or `bypass` for non-prod)
- `SESSION_AGENT_PREFER=true`
- `RUNTIME_PYTHON_IMAGE=<registry>/runtime-python:dev`
- `RUNTIME_NODE_IMAGE=<registry>/runtime-node:dev`

### 4) Port-forward services

Expose control-plane and data-plane locally:

```sh
kubectl -n <namespace> port-forward svc/<control-plane-service> 8080:8080
kubectl -n <namespace> port-forward svc/<data-plane-service> 8081:8081
```

### 5) Run the session flow script

```sh
BASE_URL=http://localhost:8080 RUNTIME=python bash scripts/run-session-step.sh
```

### 6) Kubernetes test helper script

This script port-forwards the services and runs the session flow:

```sh
NAMESPACE=<namespace> CONTROL_PLANE_SERVICE=<svc-name> DATA_PLANE_SERVICE=<svc-name> \
  python3 scripts/test-k8s-stack.py
```

---

## Notes

- Session runtime is selected at session creation time (`runtime` on POST `/sessions`).
- Step execution does not require `runtime` in the step payload.
- For local testing, `SESSION_AGENT_AUTH_BYPASS=true` is set in `make run-local-session-agent`.

## Persistence & Filesystem Locations

### Control plane

- **What is stored**: session metadata, job status, and session runtime IDs.
- **Mechanism**: SQLite for local dev (or Postgres in production).
- **Local file**: `control-plane.db` in the repo root (created by `make run-local-control-plane`).

### Data plane

- **What is stored**:
  - **Session registry** (optional): session → runtime route mappings (runtime ID, endpoint, auth mode, token, runtime).
  - **Workspaces**: files created/modified by code execution within a session.
- **Mechanism**:
  - **Registry**: in-memory or JSON file (controlled by `SESSION_REGISTRY_BACKEND` and `SESSION_REGISTRY_PATH`).
  - **Workspaces**: filesystem directory on the data-plane host.
- **Local file locations**:
  - **Registry** (if file-backed): `/tmp/session-registry.json` (default in `make run-local-data-plane`).
  - **Workspaces**: `/tmp/sessions/<sessionId>` or `/tmp/sessions/<workspaceRef>`.

### Session-agent

- **What is stored**: per-session REPL process state (in-memory only).
- **Mechanism**: in-memory map of active sessions and live interpreter processes.
- **Filesystem location**: none; state is lost when the process restarts.
