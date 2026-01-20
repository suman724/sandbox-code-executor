# Kubernetes Options Without Exec

This document outlines alternatives for session-based execution when Kubernetes `exec` into a running pod is disallowed.

## Goals

- Preserve session state across steps (REPL-like behavior).
- Route steps to the correct session runtime.
- Avoid `kubectl exec`/SPDY attach in production.

## Option A: In-Pod Session Agent (HTTP/gRPC)

Run a lightweight "session agent" inside each sandbox pod that exposes an authenticated HTTP/gRPC API. The data-plane calls this API to submit code and retrieve output. Each session maps to a dedicated pod (or a shared pod with per-session agent instances).

How it works:
- Control-plane creates a sandbox pod via data-plane (Deployment/Pod/Job).
- Pod starts an agent process listening on an internal port (e.g., `127.0.0.1:9000`).
- Data-plane sends step requests to the agent using the pod IP or a ClusterIP service.
- Agent runs a persistent interpreter for the session and returns stdout/stderr.

Pros:
- No `exec` usage; all traffic is over standard network protocols.
- Works with REPL-like state via a long-lived interpreter in the pod.
- Easier to audit and secure via mTLS or service mesh policies.

Cons:
- Requires an agent binary/library baked into runtime images.
- Needs network policy and service discovery wiring.
- Increases the surface area for authN/authZ between data-plane and agent.

Best fit:
- Long-running sessions with multiple steps.
- Environments that allow pod-to-pod traffic with strict policies.

---

## Option B: Sidecar Runner + Shared Volume Queue

Use a sidecar container per session pod. The data-plane writes step requests to a shared volume (e.g., emptyDir). The sidecar watches for new files, executes them in the runtime container, and writes results back to the volume for the data-plane to read (or exposes a local HTTP API to read results).

How it works:
- Pod contains runtime container + sidecar.
- Data-plane writes `steps/<step-id>.json` into the pod volume via an API or object store-backed volume.
- Sidecar reads steps, executes in runtime container context (or directly if it embeds runtime), and writes `results/<step-id>.json`.
- Data-plane fetches results from the volume API or object storage.

Pros:
- Avoids `exec` and direct shell access.
- Can be built around file-based IPC, which is simple and auditable.
- Sidecar can enforce policy checks and resource limits.

Cons:
- Requires a secure way for data-plane to write/read the shared volume (usually via a dedicated API).
- More moving parts (two containers, coordination, retries).
- Latency can increase due to file polling.

Best fit:
- Clusters with strict network policies where pod exec and direct pod networking are limited.

---

## Option C: Queue-Driven Session Workers (Message Broker)

Make each session pod connect to a message broker (NATS, Kafka, Redis Streams). Data-plane publishes step requests to a session-specific topic/stream; the session worker consumes, executes, and publishes results.

How it works:
- Session worker starts with session ID and subscribes to `sessions.<id>.steps`.
- Data-plane publishes steps to broker and listens on `sessions.<id>.results`.
- Worker processes steps in order and maintains in-memory state.

Pros:
- No direct pod communication required; no `exec`.
- Strong ordering guarantees with per-session topics.
- Scales well and supports retries/dead-letter flows.

Cons:
- Requires managing a broker and credentials.
- Higher operational complexity.
- Latency depends on broker and consumer configuration.

Best fit:
- Production-grade deployments that already run a message broker.

---

## Option D: Session Pod with HTTP Ingress (Service per Session)

Expose each session pod as a dedicated ClusterIP service (or headless service). The data-plane calls the service endpoint to submit steps and retrieve output.

How it works:
- Data-plane creates a pod and an accompanying service with selector labels.
- Session runtime listens on a known port; data-plane calls the service DNS.

Pros:
- Simple routing and discovery using Kubernetes services.
- No exec, standard HTTP flow.

Cons:
- Service-per-session can be expensive at scale (IP usage, API server load).
- Requires strong auth/mTLS to prevent cross-tenant calls.

Best fit:
- Lower session volume or short-lived sessions.

---

## Option E: Build and Run a Job per Step (Stateless)

If strict REPL-like state is optional, you can run each step as a separate Kubernetes Job or Pod. State persists only via shared volumes or object storage.

Pros:
- Simplest security posture (no long-lived session pods).
- Easy to enforce per-step resource limits.

Cons:
- No in-memory state across steps (only file-based state).
- Higher latency per step and more scheduling overhead.

Best fit:
- Stateless workloads or workflows with explicit file-based checkpoints.

---

## Recommendation Summary

- If you need REPL-like state and no exec: **Option A (In-Pod Session Agent)** is the most direct and scalable.
- If your platform already uses a message broker: **Option C (Queue-Driven Session Workers)** is robust and auditable.
- For lower scale and simpler ops: **Option D (Service per Session)** works well.

## Next Steps

1. Choose the runtime pattern (agent vs broker vs sidecar).
2. Define the session routing contract (session ID -> pod/worker endpoint).
3. Add authN/authZ between data-plane and session workers (mTLS, JWT, or SPIFFE).
4. Update data-plane to use the chosen transport instead of `exec`.
