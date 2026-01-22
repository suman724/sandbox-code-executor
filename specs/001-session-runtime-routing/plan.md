# Implementation Plan: Session Runtime Routing

**Branch**: `001-session-runtime-routing` | **Date**: 2026-01-20 | **Spec**: `/Users/suman/Projects/specify-playground/specs/001-session-runtime-routing/spec.md`
**Input**: Feature specification from `/specs/001-session-runtime-routing/spec.md`

## Summary

Introduce a root-level `session-agent` app packaged into Python and Node runtime pods so the data-plane can submit steps over HTTP without remote exec, while preserving REPL-like session state and supporting a non-production auth bypass.

## Technical Context

**Language/Version**: Go 1.23  
**Primary Dependencies**: chi router, Kubernetes client (existing), standard net/http  
**Storage**: Data-plane session route registry (current local registry + future persistent store)  
**Testing**: Go test (unit + integration)  
**Target Platform**: Linux server, Kubernetes  
**Project Type**: Multi-service (control-plane, data-plane, session-agent)  
**Performance Goals**: Preserve session state with step routing under 200ms p95 in-cluster  
**Constraints**: No remote exec/attach in production; per-session isolation required; auth bypass only for non-prod  
**Scale/Scope**: Long-running sessions with multiple steps, multi-tenant usage

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- Makefile targets `build`, `test`, `lint`, `format`, `run` defined and runnable
- `AGENTS.md` present with workflow and conventions
- Unit + integration tests planned for all changes
- Error handling and observability (logs/metrics) designed
- Performance/latency budgets identified for critical paths

## Project Structure

### Documentation (this feature)

```text
specs/001-session-runtime-routing/
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
  cmd/
  internal/
  tests/

data-plane/
  cmd/
  internal/
  tests/

session-agent/
  cmd/
  internal/
  tests/
```

**Structure Decision**: Multi-service layout with dedicated control-plane, data-plane, and session-agent apps.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | N/A | N/A |
