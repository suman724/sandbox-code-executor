# Implementation Plan: Sandboxed Code Execution Service

**Branch**: `001-sandboxed-code-execution` | **Date**: 2026-01-19 | **Spec**: specs/001-sandboxed-code-execution/spec.md
**Input**: Feature specification from `/specs/001-sandboxed-code-execution/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Build an internal service that executes untrusted code in ephemeral sandboxes
with strong tenant isolation, policy enforcement, and auditable observability.
Separate control plane and data plane into independently deployable services to
allow isolated scaling, deployment, and operational boundaries. The control
plane receives external requests and calls the data plane to provision sandboxes
for jobs and sessions, while the data plane runs the actual execution runtime.
Expose a tool interface via an MCP server, enforce dependency allowlists,
network isolation, workspace encryption/secure deletion, and support multiple
language runtimes. Add a non-production feature flag to bypass authorization
for testing control plane and data plane flows end-to-end.

## Technical Context

**Language/Version**: Go 1.23
**Primary Dependencies**: go-chi/chi (HTTP API), OpenTelemetry (tracing), Open Policy Agent (policy evaluation), jwt-go (authn), Kubernetes client-go (runtime coordination)
**Storage**: PostgreSQL (metadata, audit index), object storage (artifacts), Redis (session state), append-only log store (audit events)
**Testing**: go test (unit), k6 (integration/load tests)
**Target Platform**: Kubernetes (enterprise standard)
**Project Type**: multi-service (control plane + data plane)
**Performance Goals**: p95 job start latency < 3s for warm pools, p95 one-shot execution completion < 30s for standard workloads, 500 concurrent jobs per cluster baseline
**Constraints**: Default-deny network egress, strict resource limits per job, workspace encryption and secure deletion, secret injection controls, audit trail completeness, and authorization bypass flag restricted to non-production environments with audit logging
**Scale/Scope**: Support 10x burst workload over baseline without isolation violations

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- Makefile targets `build`, `test`, `lint`, `format`, `run` defined and runnable
- `AGENTS.md` present with workflow and conventions
- Unit + integration tests planned for all changes
- Error handling and observability (logs/metrics) designed
- Performance/latency budgets identified for critical paths

Status: PASS (Makefile and AGENTS present; test, observability, and latency
budgets captured in this plan.)

## Project Structure

### Documentation (this feature)

```text
specs/001-sandboxed-code-execution/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
control-plane/
├── cmd/
│   └── control-plane/
├── internal/
│   ├── api/
│   ├── audit/
│   ├── mcp/
│   ├── orchestration/
│   ├── policy/
│   ├── sessions/
│   ├── storage/
│   └── tenancy/
├── pkg/
│   └── client/
├── tests/
│   ├── contract/
│   ├── integration/
│   └── unit/
├── Dockerfile
└── go.mod

data-plane/
├── cmd/
│   └── sandbox-runner/
├── internal/
│   ├── execution/
│   ├── isolation/
│   ├── runtime/
│   ├── telemetry/
│   └── workspace/
├── tests/
│   ├── integration/
│   └── unit/
├── Dockerfile
└── go.mod

shared/
└── pkg/                 # Shared contracts/types versioned for both planes
```

**Structure Decision**: Separate control-plane and data-plane services with
independent Go modules and tests so each can be built, tested, and deployed
independently. The control plane calls the data plane over internal APIs to
create sandboxes for jobs and sessions, then streams back execution results.
Shared contracts and types live in `shared/` with versioned interfaces to avoid
tight coupling.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
