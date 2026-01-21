# Research: Session Runtime Routing

## Decision 1: Use in-pod session agent with HTTP API

**Decision**: Run a lightweight session-agent app inside each sandbox pod that exposes an authenticated HTTP API for step execution.

**Rationale**:
- Avoids `exec`/SPDY attach while preserving REPL-like state.
- Uses standard network paths that can be audited and secured.
- Keeps interpreter state in-process and persistent for the session lifetime.

**Alternatives considered**:
- gRPC streaming API: richer typing/streaming but higher complexity and schema tooling overhead.
- External broker per session: adds extra hop and infrastructure overhead.

## Decision 2: One session per sandbox pod (default)

**Decision**: Map each session to a dedicated sandbox pod by default.

**Rationale**:
- Simplifies routing and isolation guarantees.
- Reduces risk of state leakage between sessions.

**Alternatives considered**:
- Shared pod with multiple agent instances: reduces resource usage but adds multiplexing complexity and stronger isolation requirements.

## Decision 3: Session routing via explicit registry

**Decision**: Persist a session route record in the data-plane to map session IDs to pod endpoints and auth tokens.

**Rationale**:
- Ensures steps route deterministically to the correct session runtime.
- Enables recovery and troubleshooting when a runtime is unavailable.

**Alternatives considered**:
- Derive routing from Kubernetes labels only: simpler but less reliable under restarts and requires more cluster queries.

## Decision 4: Auth via short-lived token between data-plane and agent

**Decision**: Use a per-session token issued by data-plane and validated by the agent for step calls.

**Rationale**:
- Enforces least-privilege access to the session runtime.
- Compatible with service mesh or network policy enforcement.

**Alternatives considered**:
- mTLS-only: strong but depends on mesh configuration; can be layered later.

## Decision 5: Session-agent as a top-level application

**Decision**: Add `session-agent/` as a root-level app and bake it into runtime images for Python and Node.

**Rationale**:
- Keeps agent code isolated and reusable across runtimes.
- Ensures the pod always has the same agent binary regardless of runtime.

**Alternatives considered**:
- Embedding agent logic inside data-plane: increases coupling and complicates runtime packaging.

## Decision 6: Auth bypass for non-production

**Decision**: Support an explicit non-production flag to bypass session-agent authentication when running locally.

**Rationale**:
- Simplifies local testing without requiring token provisioning.
- Keeps production secure by default.

**Alternatives considered**:
- Hardcoded dev tokens: less flexible and harder to audit.
