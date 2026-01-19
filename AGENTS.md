# Agent Guidance

## Purpose

This repository supports feature specs and planning. Agents should follow the
constitution and templates in `.specify/`.

## Workflow

- Use the `/speckit.specify` flow to generate feature specs.
- Use the `/speckit.plan` flow to create implementation plans and design
  artifacts.
- Keep documentation deterministic and avoid implementation details in specs.

## Quality Gates

- Ensure `Makefile` targets run without errors.
- Require unit and integration test plans for each feature.
- Capture observability and performance budgets in plans.

## Conventions

- Use workspace-relative paths in documentation.
- Keep language and frameworks in plan technical context only.
- Prefer explicit, testable statements in requirements.
- Treat control plane and data plane as separate deployables with independent
  build, test, and release steps.

## Active Technologies
- Go 1.22 + go-chi/chi (HTTP API), OpenTelemetry (tracing), (001-sandboxed-code-execution)
- PostgreSQL (metadata, audit index), object storage (artifacts), (001-sandboxed-code-execution)
- Go 1.22 + go-chi/chi (HTTP API), OpenTelemetry (tracing), Open Policy Agent (policy evaluation), jwt-go (authn), Kubernetes client-go (runtime coordination) (001-sandboxed-code-execution)
- PostgreSQL (metadata, audit index), object storage (artifacts), Redis (session state), append-only log store (audit events) (001-sandboxed-code-execution)
- Go 1.23 + go-chi/chi (HTTP API), OpenTelemetry (tracing), Open Policy Agent (policy evaluation), jwt-go (authn), Kubernetes client-go (runtime coordination) (001-sandboxed-code-execution)

## Recent Changes
- 001-sandboxed-code-execution: Added Go 1.22 + go-chi/chi (HTTP API), OpenTelemetry (tracing),
