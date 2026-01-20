---

description: "Task list template for feature implementation"
---

# Tasks: Sandboxed Code Execution Service

**Input**: Design documents from `/specs/001-sandboxed-code-execution/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Unit and integration tests are MANDATORY for all user stories and affected components.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Multi-service**: `control-plane/`, `data-plane/`, `shared/`
- Paths shown below match the implementation plan structure

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and dependency scaffolding

- [X] T001 Add control-plane dependencies (chi, otel, opa, jwt, postgres driver) in `control-plane/go.mod`
- [X] T002 Add data-plane dependencies (otel, client-go, sandbox runtime libs) in `data-plane/go.mod`
- [X] T003 [P] Expand shared API contracts to match OpenAPI in `shared/pkg/contracts/types.go` and `shared/pkg/contracts/ids.go`
- [X] T004 [P] Extend control-plane config for storage/auth/otel/mcp in `control-plane/internal/config/config.go`
- [X] T005 [P] Extend data-plane config for runtime/auth/otel in `data-plane/internal/config/config.go`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T006 Implement chi routing and register all OpenAPI paths in `control-plane/internal/api/router.go`
- [X] T007 Implement data-plane routing for /runs create/get/terminate in `data-plane/internal/runtime/router.go`
- [X] T008 Implement JWT auth, tenant/agent context, and bypass audit logging in `control-plane/internal/api/middleware/auth.go`
- [X] T009 Implement service-to-service auth and bypass audit logging in `data-plane/internal/runtime/auth.go`
- [X] T010 Implement OPA policy evaluator loader and Evaluate logic in `control-plane/internal/policy/engine.go`
- [X] T011 Implement persistent idempotency store and wire usage in `control-plane/internal/orchestration/idempotency_store.go`
- [X] T012 Implement retry/backoff handling + degradation gates in `control-plane/internal/orchestration/retry_policy.go` and `control-plane/internal/orchestration/degradation.go`
- [X] T013 Implement PostgreSQL job/session stores (create/update/get) in `control-plane/internal/storage/postgres/job_store.go`
- [X] T014 Implement PostgreSQL policy/audit stores + audit queries in `control-plane/internal/storage/postgres/policy_store.go`
- [X] T015 Implement artifact storage adapter with signed URLs in `control-plane/internal/storage/object/artifact_store.go`
- [X] T016 Implement audit logger + persistence wiring in `control-plane/internal/audit/logger.go` and `control-plane/internal/audit/store.go`
- [X] T017 Implement runtime adapters registry (python/node) in `data-plane/internal/runtime/registry.go`
- [X] T018 Implement egress/resource limit enforcement in `data-plane/internal/isolation/egress.go` and `data-plane/internal/isolation/limits.go`
- [X] T019 Implement workspace encryption/secure delete/secret materialization in `data-plane/internal/workspace/encryption.go`, `data-plane/internal/workspace/secure_delete.go`, `data-plane/internal/workspace/secrets.go`
- [X] T020 Implement control-plane data-plane HTTP client (auth headers, timeouts, retries) in `control-plane/pkg/client/data_plane_client.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Run One-Shot Code Safely (Priority: P1) 🎯 MVP

**Goal**: Submit a one-shot job, run it in an isolated sandbox, and return outputs and artifacts.

**Independent Test**: Submit a job with known output and confirm policy enforcement, outputs, and audit records.

### Tests for User Story 1 (MANDATORY) ⚠️

- [X] T021 [P] [US1] Replace placeholder job lifecycle tests in `control-plane/internal/orchestration/job_service_test.go`
- [X] T022 [P] [US1] Replace placeholder contract tests for /jobs in `control-plane/tests/contract/jobs_contract_test.go`
- [X] T023 [P] [US1] Replace placeholder integration test for one-shot flow in `control-plane/tests/integration/job_run_test.go`
- [X] T024 [P] [US1] Expand sandbox runner tests to include artifacts and failures in `data-plane/internal/execution/runner_test.go`
- [X] T025 [P] [US1] Add allowlist enforcement tests for denied deps in `data-plane/internal/runtime/deps_test.go`

### Implementation for User Story 1

- [X] T026 [P] [US1] Expand Job model fields (agent, policy, outputs, artifacts) in `control-plane/internal/orchestration/job.go`
- [X] T027 [US1] Implement job state transitions + persistence in `control-plane/internal/orchestration/job_service.go`
- [X] T028 [US1] Align /jobs POST/GET handlers with OpenAPI schemas in `control-plane/internal/api/handlers/jobs.go`
- [X] T029 [US1] Implement run create/get/terminate handlers in `data-plane/internal/runtime/handlers.go`
- [X] T030 [US1] Implement sandbox execution runner (workspace, stdout/stderr, artifacts) in `data-plane/internal/execution/runner.go`
- [X] T031 [US1] Implement artifact metadata + upload/download wiring in `data-plane/internal/workspace/artifacts.go` and `control-plane/internal/storage/artifacts.go`
- [X] T032 [US1] Add job audit events (accepted/running/finished) in `control-plane/internal/audit/logger.go`

**Checkpoint**: User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Use Stateful Sessions with Artifacts (Priority: P2)

**Goal**: Create sessions with TTL and run multiple steps in the same sandbox with artifact handling.

**Independent Test**: Create a session, run multiple steps, verify shared state, then confirm TTL cleanup.

### Tests for User Story 2 (MANDATORY) ⚠️

- [X] T033 [P] [US2] Replace placeholder session lifecycle tests in `control-plane/internal/sessions/session_service_test.go`
- [X] T034 [P] [US2] Replace placeholder contract tests for /sessions in `control-plane/tests/contract/sessions_contract_test.go`
- [X] T035 [P] [US2] Replace placeholder integration test for session flow in `control-plane/tests/integration/session_flow_test.go`

### Implementation for User Story 2

- [X] T036 [P] [US2] Expand Session and SessionStep models in `control-plane/internal/sessions/session.go`
- [X] T037 [US2] Implement session create + step execution logic in `control-plane/internal/sessions/session_service.go`
- [X] T038 [US2] Align /sessions POST and /sessions/{id}/steps handlers with OpenAPI in `control-plane/internal/api/handlers/sessions.go`
- [X] T039 [US2] Implement session runner for step execution in `data-plane/internal/execution/session_runner.go`
- [X] T040 [US2] Implement session workspace state handling in `data-plane/internal/workspace/session.go`
- [X] T041 [US2] Implement TTL cleanup worker (expiry + secure delete) in `control-plane/internal/sessions/ttl_worker.go`
- [X] T042 [US2] Add session audit events in `control-plane/internal/audit/logger.go`

**Checkpoint**: User Stories 1 and 2 should both work independently

---

## Phase 5: User Story 3 - Configure Policies and Audit Execution (Priority: P2)

**Goal**: Manage policies and query audit trails for compliance and visibility.

**Independent Test**: Create a policy, enforce it on a job, and query audit events that record the decision.

### Tests for User Story 3 (MANDATORY) ⚠️

- [X] T043 [P] [US3] Replace placeholder policy CRUD tests in `control-plane/internal/policy/policy_service_test.go`
- [X] T044 [P] [US3] Replace placeholder contract tests for /policies in `control-plane/tests/contract/policies_contract_test.go`
- [X] T045 [P] [US3] Replace placeholder integration test for audit queries in `control-plane/tests/integration/audit_query_test.go`

### Implementation for User Story 3

- [X] T046 [US3] Implement policy store rules + versioning in `control-plane/internal/policy/policy_store.go`
- [X] T047 [US3] Align /policies handler validation with OpenAPI in `control-plane/internal/api/handlers/policies.go`
- [X] T048 [US3] Implement audit query filters (tenant, time range) in `control-plane/internal/api/handlers/audit.go`
- [X] T049 [US3] Wire policy enforcement inputs in `control-plane/internal/orchestration/job_service.go` and `control-plane/internal/sessions/session_service.go`
- [X] T050 [US3] Implement audit persistence wiring in `control-plane/internal/storage/postgres/policy_store.go`

**Checkpoint**: User Stories 1–3 should work independently

---

## Phase 6: User Story 4 - Run Optional Service Mode (Priority: P3)

**Goal**: Start a long-lived service in a sandbox and expose it through a controlled proxy.

**Independent Test**: Start a service, access it via proxy, then confirm termination when TTL or policy changes.

### Tests for User Story 4 (MANDATORY) ⚠️

- [X] T051 [P] [US4] Replace placeholder service lifecycle tests in `control-plane/internal/services/service_service_test.go`
- [X] T052 [P] [US4] Replace placeholder contract tests for /services in `control-plane/tests/contract/services_contract_test.go`
- [X] T053 [P] [US4] Replace placeholder integration test for service proxy flow in `control-plane/tests/integration/service_proxy_test.go`

### Implementation for User Story 4

- [X] T054 [P] [US4] Expand Service model fields in `control-plane/internal/services/service.go`
- [X] T055 [US4] Align /services handler response + proxy URL with OpenAPI in `control-plane/internal/api/handlers/services.go`
- [X] T056 [US4] Implement sandbox service runner lifecycle in `data-plane/internal/runtime/service_runner.go`
- [X] T057 [US4] Implement proxy registration/teardown in `data-plane/internal/runtime/proxy.go`
- [X] T058 [US4] Add service audit events in `control-plane/internal/audit/logger.go` and `data-plane/internal/telemetry/logger.go`

**Checkpoint**: User Stories 1–4 should work independently

---

## Phase 7: User Story 5 - Orchestrate Multi-Agent Workflows (Priority: P3)

**Goal**: Orchestrate multi-agent workflows with shared memory and tracking.

**Independent Test**: Run a multi-step workflow and verify each step state and shared memory access.

### Tests for User Story 5 (MANDATORY) ⚠️

- [X] T059 [P] [US5] Replace placeholder workflow tests in `control-plane/internal/orchestration/workflow_service_test.go`
- [X] T060 [P] [US5] Replace placeholder contract tests for /workflows in `control-plane/tests/contract/workflows_contract_test.go`
- [X] T061 [P] [US5] Replace placeholder integration test for workflow run in `control-plane/tests/integration/workflow_run_test.go`

### Implementation for User Story 5

- [X] T062 [P] [US5] Expand Workflow and WorkflowStep models in `control-plane/internal/orchestration/workflow.go`
- [X] T063 [US5] Implement workflow orchestration logic in `control-plane/internal/orchestration/workflow_service.go`
- [X] T064 [US5] Align /workflows handler with OpenAPI in `control-plane/internal/api/handlers/workflows.go`
- [X] T065 [US5] Implement shared memory store in `control-plane/internal/orchestration/memory_store.go`
- [X] T066 [US5] Add workflow audit events in `control-plane/internal/audit/logger.go`

**Checkpoint**: All user stories should now be independently functional

---

## Phase 8: Database Abstraction & SQLite Support (Cross-Cutting)

**Purpose**: Support Postgres and SQLite for non-production testing without changing service logic.

- [X] T067 [P] Add database driver configuration and validation in `control-plane/internal/config/config.go`
- [X] T068 [P] Add SQLite schema + migration helper in `control-plane/internal/storage/sqlite/schema.sql` and `control-plane/internal/storage/sqlite/migrate.go`
- [X] T069 [P] Implement SQLite job/session stores in `control-plane/internal/storage/sqlite/job_store.go` and `control-plane/internal/storage/sqlite/session_store.go`
- [X] T070 [P] Implement SQLite policy/audit stores in `control-plane/internal/storage/sqlite/policy_store.go` and `control-plane/internal/storage/sqlite/audit_store.go`
- [X] T071 [P] Implement SQLite idempotency store in `control-plane/internal/orchestration/sqlite_idempotency_store.go`
- [X] T072 Add store factory + wiring for Postgres/SQLite in `control-plane/internal/storage/factory.go`, `control-plane/internal/api/router.go`, and `control-plane/cmd/control-plane/main.go`
- [X] T073 [P] Add SQLite-backed tests for non-production scenarios in `control-plane/tests/integration/sqlite_store_test.go`

---

## Phase 9: MCP Tool Interface (Cross-Cutting)

**Purpose**: Expose agent-accessible tools via MCP for job/session/artifact/workflow operations

- [X] T074 [P] Implement MCP server wiring in `control-plane/internal/mcp/server.go` and `control-plane/internal/mcp/router.go`
- [X] T075 Implement MCP jobs tool in `control-plane/internal/mcp/tools/jobs.go`
- [X] T076 Implement MCP sessions tool in `control-plane/internal/mcp/tools/sessions.go`
- [X] T077 Implement MCP artifacts tool in `control-plane/internal/mcp/tools/artifacts.go`
- [X] T078 Implement MCP workflows tool in `control-plane/internal/mcp/tools/workflows.go`
- [X] T079 Replace placeholder MCP integration tests in `control-plane/tests/integration/mcp_tools_test.go`

---

## Phase 10: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T080 [P] Add OpenTelemetry init and propagation in `control-plane/cmd/control-plane/main.go` and `data-plane/cmd/sandbox-runner/main.go`
- [X] T081 Add metrics for latency/queue depth/denials in `control-plane/internal/orchestration/job_service.go` and `data-plane/internal/execution/runner.go`
- [X] T082 [P] Flesh out k6 scenarios in `control-plane/tests/integration/k6/jobs.js` and `control-plane/tests/integration/k6/sessions.js`
- [X] T083 Update validation steps in `specs/001-sandboxed-code-execution/quickstart.md`

---

## Deferred Follow-ups

**Purpose**: Items identified for later hardening.

- [ ] T084 Implement artifact persistence beyond no-op validation in `control-plane/internal/storage/object/artifact_store.go`
- [ ] T085 Wire audit store to persistent backing (Postgres/SQLite) in `control-plane/cmd/control-plane/main.go` and `control-plane/internal/audit/store.go`

---

## Session Persistence & Execution Model (Follow-up)

**Purpose**: Ensure session steps execute in a persistent runtime (local child process for dev, Kubernetes resource for prod).

- [X] T086 Define session runtime contract and lifecycle (start/step/terminate) in `specs/001-sandboxed-code-execution/contracts/data-plane-openapi.yaml`
- [X] T087 Add data-plane session runtime interface and implementations in `data-plane/internal/runtime/session_runtime.go`
- [X] T088 Implement local session runtime using child processes in `data-plane/internal/runtime/session_local.go`
- [X] T089 Implement Kubernetes-backed session runtime (Pod/Job with RuntimeClass) in `data-plane/internal/runtime/session_k8s.go`
- [X] T090 Implement data-plane session runtime registry for session_id -> runtime_id in `data-plane/internal/runtime/session_registry.go`
- [X] T091 Wire session step handling to data-plane runtime (HTTP endpoint + handler) in `data-plane/internal/runtime/handlers.go` and `data-plane/internal/runtime/router.go`
- [X] T092 Update control-plane session step execution to call data-plane endpoint in `control-plane/internal/sessions/session_service.go` and `control-plane/pkg/client/data_plane_client.go`
- [X] T093 Add session step persistence and state tracking for running runtime IDs in `control-plane/internal/sessions/session.go` and `control-plane/internal/storage/postgres/job_store.go` (or new session store file)
- [X] T094 Add integration tests for session step persistence (local) in `control-plane/tests/integration/session_flow_test.go` and `data-plane/tests/integration/session_runner_test.go`
- [X] T095 Decide and implement persistence backend for session registry (in-memory for local, persistent for prod) in `data-plane/internal/runtime/session_registry.go` and configuration wiring in `data-plane/internal/config/config.go`
- [ ] T096 Define shared session registry and routing strategy for multi-replica data-plane (e.g., DB/Redis-backed registry + sticky routing) in `specs/001-sandboxed-code-execution/tasks.md`
- [X] T097 Extend session step API contract to return stdout/stderr (add schema) in `specs/001-sandboxed-code-execution/contracts/data-plane-openapi.yaml` and `data-plane/internal/runtime/openapi.yaml`
- [X] T098 Capture session step stdout/stderr in local runtime and return in response in `data-plane/internal/runtime/session_local.go` and `data-plane/internal/runtime/handlers.go`
- [X] T099 Capture session step stdout/stderr from k8s exec and return in response in `data-plane/internal/runtime/session_k8s.go` and `data-plane/internal/runtime/handlers.go`
- [X] T100 Propagate session step output through control-plane client and API response in `control-plane/pkg/client/data_plane_client.go`, `control-plane/internal/sessions/session_service.go`, and `control-plane/internal/api/handlers/sessions.go`
- [X] T101 Add tests for session step output propagation in `data-plane/tests/integration/session_runner_test.go` and `control-plane/tests/integration/session_flow_test.go`

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Database Abstraction (Phase 8)**: Depends on Foundational completion
- **MCP Tools (Phase 9)**: Depends on Foundational and User Story 1 completion
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - May integrate with US1 but should be independently testable
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - May integrate with US1 but should be independently testable
- **User Story 4 (P3)**: Can start after Foundational (Phase 2) - May integrate with US1/US2 but should be independently testable
- **User Story 5 (P3)**: Can start after Foundational (Phase 2) - May integrate with prior stories but should be independently testable

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Models before services
- Services before endpoints
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- All tests for a user story marked [P] can run in parallel
- Models within a story marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
Task: "Replace placeholder job lifecycle tests in control-plane/internal/orchestration/job_service_test.go"
Task: "Replace placeholder contract tests for /jobs in control-plane/tests/contract/jobs_contract_test.go"
Task: "Implement artifact metadata + upload/download wiring in data-plane/internal/workspace/artifacts.go"
```

---

## Parallel Example: User Story 2

```bash
Task: "Replace placeholder session lifecycle tests in control-plane/internal/sessions/session_service_test.go"
Task: "Implement session runner for step execution in data-plane/internal/execution/session_runner.go"
Task: "Implement TTL cleanup worker (expiry + secure delete) in control-plane/internal/sessions/ttl_worker.go"
```

---

## Parallel Example: User Story 3

```bash
Task: "Replace placeholder policy CRUD tests in control-plane/internal/policy/policy_service_test.go"
Task: "Align /policies handler validation with OpenAPI in control-plane/internal/api/handlers/policies.go"
Task: "Implement audit query filters (tenant, time range) in control-plane/internal/api/handlers/audit.go"
```

---

## Parallel Example: User Story 4

```bash
Task: "Replace placeholder service lifecycle tests in control-plane/internal/services/service_service_test.go"
Task: "Implement sandbox service runner lifecycle in data-plane/internal/runtime/service_runner.go"
Task: "Implement proxy registration/teardown in data-plane/internal/runtime/proxy.go"
```

---

## Parallel Example: User Story 5

```bash
Task: "Replace placeholder workflow tests in control-plane/internal/orchestration/workflow_service_test.go"
Task: "Implement workflow orchestration logic in control-plane/internal/orchestration/workflow_service.go"
Task: "Implement shared memory store in control-plane/internal/orchestration/memory_store.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Add User Story 4 → Test independently → Deploy/Demo
6. Add User Story 5 → Test independently → Deploy/Demo
7. Add MCP tools → Test independently → Deploy/Demo
8. Add Polish items

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1
   - Developer B: User Story 2
   - Developer C: User Story 3
   - Developer D: User Story 4
   - Developer E: User Story 5
   - Developer F: MCP tools

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
