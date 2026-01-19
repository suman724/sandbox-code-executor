# Phase 0 Research: Sandboxed Code Execution Service

## Decision 1: Control plane vs data plane split
- Decision: Separate control plane and data plane into independently deployable
  services with distinct Go modules.
- Rationale: Enables independent scaling, failure isolation, and deployment
  cadence, while aligning with security boundaries between orchestration and
  sandbox execution.
- Alternatives considered: Single monolith, shared binary with feature flags.

## Decision 2: Control plane to data plane integration
- Decision: Control plane calls the data plane over internal APIs to provision
  sandboxes for jobs and sessions, then tracks execution status and artifacts.
- Rationale: Centralized policy enforcement and audit tracking in the control
  plane while keeping execution isolated in the data plane.
- Alternatives considered: Direct client calls to data plane, asynchronous queue
  without synchronous control plane coordination.

## Decision 3: MCP tool interface
- Decision: Provide an MCP server in the control plane to expose job, session,
  artifact, and workflow tools to agents.
- Rationale: Standardizes agent integration and keeps tool access governed by
  control-plane policy.
- Alternatives considered: Direct REST-only access, embedded agent plugins.

## Decision 4: Dependency allowlist enforcement
- Decision: Enforce dependency allowlists in the data plane before runtime
  installs and record decisions in audit logs.
- Rationale: Reduces supply-chain risk and prevents unvetted packages.
- Alternatives considered: Allow-all with monitoring, global denylist only.

## Decision 5: Runtime registry for multi-language support
- Decision: Implement a runtime registry in the data plane to map languages to
  sandbox images and execution adapters.
- Rationale: Enables controlled expansion to new languages without breaking
  clients.
- Alternatives considered: Hard-coded language handling.

## Decision 6: Isolation and egress controls
- Decision: Enforce default-deny egress and resource limits at the data plane
  using runtime and policy constraints.
- Rationale: Prevents credential exfiltration and cross-tenant leakage.
- Alternatives considered: Network monitoring only.

## Decision 7: Workspace encryption and secure deletion
- Decision: Encrypt workspace storage and perform secure deletion on teardown.
- Rationale: Meets data control requirements for sensitive artifacts.
- Alternatives considered: Best-effort deletion without encryption.

## Decision 8: Secret injection controls
- Decision: Validate secret injection policy in the control plane and gate
  secret materialization in the data plane.
- Rationale: Prevents secrets from being persisted or exfiltrated.
- Alternatives considered: Global secret injection without policy checks.

## Decision 9: Authorization bypass for testing
- Decision: Add an explicit feature flag that bypasses authorization checks in
  both planes, restricted to non-production environments and always audited.
- Rationale: Enables integration testing without production credentials while
  preserving traceability and preventing accidental production use.
- Alternatives considered: Mock auth only, separate test-only endpoints.

## Decision 10: Container packaging
- Decision: Maintain separate Dockerfiles for control plane and data plane.
- Rationale: Enables independent builds, release pipelines, and runtime
  constraints for each service.
- Alternatives considered: Single multi-binary image, shared base image only.

## Decision 11: API framework
- Decision: Use a lightweight HTTP router (go-chi/chi) for the control plane.
- Rationale: Minimal middleware overhead, strong community use, fits internal
  service patterns, and keeps request handling explicit.
- Alternatives considered: Gin, Echo, net/http with custom routing.

## Decision 12: Sandbox runtime integration
- Decision: Use Kubernetes RuntimeClass with gVisor for sandboxed execution
  pods, coordinated via client-go.
- Rationale: Strong isolation within Kubernetes, operational alignment with
  enterprise clusters, and avoids managing custom hypervisors initially.
- Alternatives considered: Firecracker-based microVMs, Kata Containers,
  user-namespace isolation.

## Decision 13: Policy evaluation
- Decision: Use Open Policy Agent for policy evaluation and enforcement.
- Rationale: Expressive policies, versioned rule sets, and well-established
  enterprise usage for multi-tenant authorization and controls.
- Alternatives considered: Custom policy engine, Cedar, Casbin.

## Decision 14: Authentication and authorization
- Decision: Validate JWTs issued by the internal identity provider and map to
  tenant/app/agent identities.
- Rationale: Aligns with existing internal auth, supports service-to-service
  calls, and provides consistent identity claims.
- Alternatives considered: mTLS-only auth, custom API keys.

## Decision 15: Storage layout
- Decision: PostgreSQL for metadata and audit indexing, object storage for
  artifacts, Redis for session state, append-only log store for audit events.
- Rationale: Separates hot session state from durable audit records and keeps
  artifacts in scalable object storage.
- Alternatives considered: Single database for all data, document store for
  audit data, in-cluster file storage.

## Decision 16: Observability
- Decision: Use OpenTelemetry for traces and structured logs with correlation
  ids, plus metrics for queue depth, execution latency, and policy denials.
- Rationale: Standardized telemetry makes auditing and incident response
  consistent across tenants.
- Alternatives considered: Proprietary logging formats, tracing-only approach.

## Decision 17: Controlled service mode
- Decision: Expose services through a controlled proxy with explicit allowlist
  and per-tenant rate limits.
- Rationale: Keeps long-lived services isolated while enabling safe access.
- Alternatives considered: Direct pod networking access, shared ingress without
  policy controls.

## Decision 18: Reliability strategy
- Decision: Idempotent job submission with request ids and at-least-once
  execution guarded by job de-duplication.
- Rationale: Supports retries without duplicate side effects or double billing.
- Alternatives considered: Best-effort retries without idempotency.
