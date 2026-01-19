---

description: "Task list template for feature implementation"
---

# Tasks: Sandboxed Code Execution Service

**Input**: Design documents from `/specs/001-sandboxed-code-execution/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: The examples below include test tasks. Unit and integration tests are MANDATORY for all user stories and affected components.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Multi-service**: `control-plane/`, `data-plane/`, `shared/`
- Paths shown below match the implementation plan structure

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Create project structure per implementation plan in control-plane/, data-plane/, shared/
- [X] T002 Initialize Go module for control plane in control-plane/go.mod
- [X] T003 Initialize Go module for data plane in data-plane/go.mod
- [X] T004 Initialize shared contracts module in shared/go.mod and shared/pkg/contracts/types.go
- [X] T005 [P] Create Go workspace for multi-module development in go.work
- [X] T006 [P] Configure linting in .golangci.yml
- [X] T007 Update Makefile targets to build/test/lint/format/run each plane in Makefile

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T008 Define shared error and ID types in shared/pkg/contracts/errors.go and shared/pkg/contracts/ids.go
- [X] T009 Implement control-plane config loader and validation in control-plane/internal/config/config.go
- [X] T010 Implement data-plane config loader and validation in data-plane/internal/config/config.go
- [X] T011 Implement control-plane auth middleware with AUTHZ_BYPASS flag in control-plane/internal/api/middleware/auth.go
- [X] T012 Implement data-plane auth middleware with AUTHZ_BYPASS flag in data-plane/internal/runtime/auth.go
- [X] T013 Implement control-plane audit logger in control-plane/internal/audit/logger.go
- [X] T014 Implement data-plane audit logger in data-plane/internal/telemetry/logger.go
- [X] T015 Implement policy engine adapter in control-plane/internal/policy/engine.go
- [X] T016 Implement control-plane to data-plane client in control-plane/pkg/client/data_plane_client.go
- [X] T017 Implement metadata storage interfaces in control-plane/internal/storage/store.go
- [X] T018 Implement PostgreSQL job/session metadata store in control-plane/internal/storage/postgres/job_store.go
- [X] T019 Implement PostgreSQL policy/audit metadata store in control-plane/internal/storage/postgres/policy_store.go
- [X] T020 Implement artifact storage interface in control-plane/internal/storage/artifacts.go
- [X] T021 Implement object storage adapter in control-plane/internal/storage/object/artifact_store.go
- [X] T022 Implement base router setup in control-plane/internal/api/router.go
- [X] T023 Implement base router setup in data-plane/internal/runtime/router.go
- [X] T024 Implement runtime registry in data-plane/internal/runtime/registry.go
- [X] T025 Implement network isolation enforcement in data-plane/internal/isolation/egress.go
- [X] T026 Implement resource limit enforcement in data-plane/internal/isolation/limits.go
- [X] T027 Implement workspace encryption utilities in data-plane/internal/workspace/encryption.go
- [X] T028 Implement secure deletion routines in data-plane/internal/workspace/secure_delete.go
- [X] T029 Implement secret materialization guard in data-plane/internal/workspace/secrets.go
- [X] T030 Implement idempotency key store in control-plane/internal/orchestration/idempotency_store.go
- [X] T031 Implement retry policy helper in control-plane/internal/orchestration/retry_policy.go
- [X] T032 Implement graceful degradation hooks in control-plane/internal/orchestration/degradation.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Run One-Shot Code Safely (Priority: P1) 🎯 MVP

**Goal**: Submit a one-shot job, run it in an isolated sandbox, and return outputs and artifacts.

**Independent Test**: Submit a job with known output and confirm policy enforcement, outputs, and audit records.

### Tests for User Story 1 (MANDATORY) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T033 [P] [US1] Unit test job lifecycle in control-plane/internal/orchestration/job_service_test.go
- [X] T034 [P] [US1] Unit test sandbox runner in data-plane/internal/execution/runner_test.go
- [X] T035 [P] [US1] Unit test dependency allowlist in data-plane/internal/runtime/deps_test.go
- [X] T036 [P] [US1] Unit test runtime registry in data-plane/internal/runtime/registry_test.go
- [X] T037 [P] [US1] Contract test /jobs in control-plane/tests/contract/jobs_contract_test.go
- [X] T038 [P] [US1] Integration test one-shot flow in control-plane/tests/integration/job_run_test.go

### Implementation for User Story 1

- [X] T039 [P] [US1] Create Job model in control-plane/internal/orchestration/job.go
- [X] T040 [US1] Implement Job service in control-plane/internal/orchestration/job_service.go
- [X] T041 [US1] Implement jobs API handler in control-plane/internal/api/handlers/jobs.go
- [X] T042 [US1] Implement data-plane run handler in data-plane/internal/runtime/handlers.go
- [X] T043 [US1] Implement sandbox execution runner in data-plane/internal/execution/runner.go
- [X] T044 [US1] Implement dependency allowlist enforcement in data-plane/internal/runtime/deps.go
- [X] T045 [US1] Implement runtime registry usage in data-plane/internal/runtime/registry.go
- [X] T046 [US1] Implement artifact capture in data-plane/internal/workspace/artifacts.go
- [X] T047 [US1] Add error handling/logging in control-plane/internal/api/handlers/jobs.go and data-plane/internal/execution/runner.go

**Checkpoint**: User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Use Stateful Sessions with Artifacts (Priority: P2)

**Goal**: Create sessions with TTL and run multiple steps in the same sandbox with artifact handling.

**Independent Test**: Create a session, run multiple steps, verify shared state, then confirm TTL cleanup.

### Tests for User Story 2 (MANDATORY) ⚠️

- [X] T048 [P] [US2] Unit test session lifecycle in control-plane/internal/sessions/session_service_test.go
- [X] T049 [P] [US2] Contract test /sessions in control-plane/tests/contract/sessions_contract_test.go
- [X] T050 [P] [US2] Integration test session flow in control-plane/tests/integration/session_flow_test.go

### Implementation for User Story 2

- [X] T051 [P] [US2] Create Session model in control-plane/internal/sessions/session.go
- [X] T052 [US2] Implement Session service in control-plane/internal/sessions/session_service.go
- [X] T053 [US2] Implement sessions API handler in control-plane/internal/api/handlers/sessions.go
- [X] T054 [US2] Implement session workspace in data-plane/internal/workspace/session.go
- [X] T055 [US2] Implement session step execution in data-plane/internal/execution/session_runner.go
- [X] T056 [US2] Implement TTL cleanup worker in control-plane/internal/sessions/ttl_worker.go
- [X] T057 [US2] Add error handling/logging in control-plane/internal/api/handlers/sessions.go

**Checkpoint**: User Stories 1 and 2 should both work independently

---

## Phase 5: User Story 3 - Configure Policies and Audit Execution (Priority: P2)

**Goal**: Manage policies and query audit trails for compliance and visibility.

**Independent Test**: Create a policy, enforce it on a job, and query audit events that record the decision.

### Tests for User Story 3 (MANDATORY) ⚠️

- [X] T058 [P] [US3] Unit test policy CRUD in control-plane/internal/policy/policy_service_test.go
- [X] T059 [P] [US3] Contract test /policies in control-plane/tests/contract/policies_contract_test.go
- [X] T060 [P] [US3] Integration test audit query in control-plane/tests/integration/audit_query_test.go

### Implementation for User Story 3

- [X] T061 [P] [US3] Implement policy store in control-plane/internal/policy/policy_store.go
- [X] T062 [US3] Implement policy handlers in control-plane/internal/api/handlers/policies.go
- [X] T063 [US3] Implement audit query handler in control-plane/internal/api/handlers/audit.go
- [X] T064 [US3] Implement audit persistence in control-plane/internal/audit/store.go
- [X] T065 [US3] Enforce policy checks in control-plane/internal/orchestration/job_service.go and control-plane/internal/sessions/session_service.go

**Checkpoint**: User Stories 1–3 should work independently

---

## Phase 6: User Story 4 - Run Optional Service Mode (Priority: P3)

**Goal**: Start a long-lived service in a sandbox and expose it through a controlled proxy.

**Independent Test**: Start a service, access it via proxy, then confirm termination when TTL or policy changes.

### Tests for User Story 4 (MANDATORY) ⚠️

- [X] T066 [P] [US4] Unit test service lifecycle in control-plane/internal/services/service_service_test.go
- [X] T067 [P] [US4] Contract test /services in control-plane/tests/contract/services_contract_test.go
- [X] T068 [P] [US4] Integration test service proxy flow in control-plane/tests/integration/service_proxy_test.go

### Implementation for User Story 4

- [X] T069 [P] [US4] Create Service model in control-plane/internal/services/service.go
- [X] T070 [US4] Implement services handler in control-plane/internal/api/handlers/services.go
- [X] T071 [US4] Implement service runtime in data-plane/internal/runtime/service_runner.go
- [X] T072 [US4] Implement proxy registration in data-plane/internal/runtime/proxy.go
- [X] T073 [US4] Add audit logging for service start/stop in control-plane/internal/audit/logger.go and data-plane/internal/telemetry/logger.go

**Checkpoint**: User Stories 1–4 should work independently

---

## Phase 7: User Story 5 - Orchestrate Multi-Agent Workflows (Priority: P3)

**Goal**: Orchestrate multi-agent workflows with shared memory and tracking.

**Independent Test**: Run a multi-step workflow and verify each step state and shared memory access.

### Tests for User Story 5 (MANDATORY) ⚠️

- [X] T074 [P] [US5] Unit test workflow orchestration in control-plane/internal/orchestration/workflow_service_test.go
- [X] T075 [P] [US5] Contract test /workflows in control-plane/tests/contract/workflows_contract_test.go
- [X] T076 [P] [US5] Integration test workflow run in control-plane/tests/integration/workflow_run_test.go

### Implementation for User Story 5

- [X] T077 [P] [US5] Create Workflow model in control-plane/internal/orchestration/workflow.go
- [X] T078 [US5] Implement workflow service in control-plane/internal/orchestration/workflow_service.go
- [X] T079 [US5] Implement workflows handler in control-plane/internal/api/handlers/workflows.go
- [X] T080 [US5] Implement shared memory store in control-plane/internal/orchestration/memory_store.go
- [X] T081 [US5] Add workflow audit events in control-plane/internal/audit/logger.go

**Checkpoint**: All user stories should now be independently functional

---

## Phase 8: MCP Tool Interface (Cross-Cutting)

**Purpose**: Expose agent-accessible tools via MCP for job/session/artifact/workflow operations

- [X] T082 [P] Implement MCP server core in control-plane/internal/mcp/server.go
- [X] T083 Implement MCP tool routing in control-plane/internal/mcp/router.go
- [X] T084 Implement MCP job tools in control-plane/internal/mcp/tools/jobs.go
- [X] T085 Implement MCP session tools in control-plane/internal/mcp/tools/sessions.go
- [X] T086 Implement MCP artifact tools in control-plane/internal/mcp/tools/artifacts.go
- [X] T087 Implement MCP workflow tools in control-plane/internal/mcp/tools/workflows.go
- [X] T088 Add MCP integration tests in control-plane/tests/integration/mcp_tools_test.go

---

## Phase 9: Deployment & Reliability (Cross-Cutting)

**Purpose**: Deployment artifacts and reliability features for both planes

- [X] T089 [P] Add control-plane Kubernetes deployment manifest in deploy/control-plane/deployment.yaml
- [X] T090 [P] Add data-plane Kubernetes deployment manifest in deploy/data-plane/deployment.yaml
- [X] T091 [P] Add control-plane service manifest in deploy/control-plane/service.yaml
- [X] T092 [P] Add data-plane service manifest in deploy/data-plane/service.yaml
- [X] T093 Add Helm chart for control plane in deploy/helm/control-plane/Chart.yaml
- [X] T094 Add Helm chart for data plane in deploy/helm/data-plane/Chart.yaml
- [X] T095 Add k6 integration test for jobs API in control-plane/tests/integration/k6/jobs.js
- [X] T096 Add k6 integration test for sessions API in control-plane/tests/integration/k6/sessions.js
- [X] T097 Add idempotency integration test in control-plane/tests/integration/idempotency_test.go
- [X] T098 Add retry/degradation integration test in control-plane/tests/integration/reliability_test.go

---

## Phase 10: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T099 [P] Add performance baseline test in control-plane/tests/integration/perf_baseline_test.go
- [X] T100 [P] Add data-plane load test in data-plane/tests/integration/load_test.go
- [X] T101 Update quickstart validation steps in specs/001-sandboxed-code-execution/quickstart.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **MCP Tools (Phase 8)**: Depends on Foundational and User Story 1 completion
- **Deployment & Reliability (Phase 9)**: Depends on Foundational and User Story 1 completion
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
Task: "Unit test job lifecycle in control-plane/internal/orchestration/job_service_test.go"
Task: "Unit test dependency allowlist in data-plane/internal/runtime/deps_test.go"
Task: "Integration test one-shot flow in control-plane/tests/integration/job_run_test.go"
```

---

## Parallel Example: User Story 2

```bash
Task: "Unit test session lifecycle in control-plane/internal/sessions/session_service_test.go"
Task: "Contract test /sessions in control-plane/tests/contract/sessions_contract_test.go"
Task: "Integration test session flow in control-plane/tests/integration/session_flow_test.go"
```

---

## Parallel Example: User Story 3

```bash
Task: "Unit test policy CRUD in control-plane/internal/policy/policy_service_test.go"
Task: "Contract test /policies in control-plane/tests/contract/policies_contract_test.go"
Task: "Integration test audit query in control-plane/tests/integration/audit_query_test.go"
```

---

## Parallel Example: MCP Tools

```bash
Task: "Implement MCP server core in control-plane/internal/mcp/server.go"
Task: "Implement MCP job tools in control-plane/internal/mcp/tools/jobs.go"
Task: "Add MCP integration tests in control-plane/tests/integration/mcp_tools_test.go"
```

---

## Parallel Example: Deployment & Reliability

```bash
Task: "Add control-plane Kubernetes deployment manifest in deploy/control-plane/deployment.yaml"
Task: "Add data-plane Kubernetes deployment manifest in deploy/data-plane/deployment.yaml"
Task: "Add k6 integration test for jobs API in control-plane/tests/integration/k6/jobs.js"
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
8. Add Deployment & Reliability → Test independently → Deploy/Demo

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
   - Developer G: Deployment & Reliability
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
