---

description: "Task list for session runtime routing"
---

# Tasks: Session Runtime Routing

**Input**: Design documents from `/specs/001-session-runtime-routing/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Unit and integration tests are mandatory for all user stories and affected components.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Create session-agent module skeleton in `session-agent/cmd/session-agent/main.go` and `session-agent/go.mod`
- [X] T002 Add session-agent to workspace in `go.work`
- [X] T003 Update build/test/run targets for session-agent in `Makefile`
- [X] T004 [P] Add session-agent package documentation stub in `session-agent/README.md`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure required before any user story

- [X] T005 Define shared agent API models in `shared/sessionagent/types.go`
- [X] T006 Implement session-agent configuration (auth token, auth mode, listen address) in `session-agent/internal/config/config.go`
- [X] T007 Implement session-agent HTTP router (health + steps) in `session-agent/internal/api/router.go`
- [X] T008 Implement session-agent auth middleware with non-prod bypass flag in `session-agent/internal/api/middleware/auth.go`
- [X] T009 Implement session-agent runtime runner (persistent interpreter) in `session-agent/internal/runtime/runner.go`
- [X] T010 Add structured logging helpers in `session-agent/internal/telemetry/logging.go`
- [X] T011 Implement agent client for data-plane in `data-plane/internal/runtime/agent_client.go`
- [X] T012 Extend session registry to store endpoint/token/auth mode in `data-plane/internal/runtime/session_registry.go`
- [X] T013 Update data-plane runtime API to return route metadata in `data-plane/internal/runtime/session_runtime.go`
- [X] T014 Update data-plane handlers to store session routes in `data-plane/internal/runtime/handlers.go`
- [X] T015 Add data-plane config for agent endpoint/auth mode in `data-plane/internal/config/config.go`
- [X] T016 Add runtime images with session-agent baked in at `deploy/runtime/python/Dockerfile`
- [X] T017 Add runtime images with session-agent baked in at `deploy/runtime/node/Dockerfile`
- [X] T018 Wire Kubernetes runtime to use runtime-specific images and agent env in `data-plane/internal/runtime/session_k8s.go`
- [X] T019 Add local session-agent process launcher for non-prod in `data-plane/internal/runtime/session_local.go`

---

## Phase 3: User Story 1 - Preserve Session State Across Steps (Priority: P1) 🎯 MVP

**Goal**: Steps in a session run in the same long-lived interpreter with stdout/stderr returned.

**Independent Test**: Create a session, set a variable in one step, use it in a second step, and verify output.

### Tests for User Story 1 (MANDATORY) ⚠️

- [X] T020 [P] [US1] Unit test for runner state persistence in `session-agent/tests/unit/runner_test.go`
- [X] T021 [P] [US1] Integration test for session step persistence in `data-plane/tests/integration/session_agent_local_test.go`
- [X] T022 [P] [US1] Contract test for agent API in `session-agent/tests/contract/agent_api_test.go`
- [X] T042 [P] [US1] Integration test for file-based workspace persistence in `data-plane/tests/integration/session_workspace_persistence_test.go`

### Implementation for User Story 1

- [X] T023 [US1] Implement step execution and stdout/stderr capture in `session-agent/internal/runtime/runner.go`
- [X] T024 [US1] Implement step handler wiring to runner in `session-agent/internal/api/handlers/steps.go`
- [X] T025 [US1] Return step output to data-plane via client in `data-plane/internal/runtime/agent_client.go`
- [X] T026 [US1] Update data-plane step handling to return stdout/stderr in `data-plane/internal/runtime/handlers.go`

---

## Phase 4: User Story 2 - Route Steps to the Correct Runtime (Priority: P2)

**Goal**: Steps are routed to the correct session runtime using registry entries.

**Independent Test**: Create two sessions, run steps, and verify outputs do not cross.

### Tests for User Story 2 (MANDATORY) ⚠️

- [X] T027 [P] [US2] Unit test for session registry route lookups in `data-plane/tests/unit/session_registry_test.go`
- [X] T028 [P] [US2] Integration test for multi-session routing in `data-plane/tests/integration/session_routing_test.go`

### Implementation for User Story 2

- [X] T029 [US2] Store session route metadata on create in `data-plane/internal/runtime/handlers.go`
- [X] T030 [US2] Resolve session route and invoke agent client in `data-plane/internal/runtime/handlers.go`
- [X] T031 [US2] Update runtime OpenAPI to include route metadata in `data-plane/internal/runtime/openapi.yaml`

---

## Phase 5: User Story 3 - Operate Without Remote Exec Access (Priority: P3)

**Goal**: Production mode avoids exec/attach and supports auth bypass in non-prod.

**Independent Test**: Run a session in non-prod with auth bypass, and ensure production mode enforces auth without exec.

### Tests for User Story 3 (MANDATORY) ⚠️

- [X] T032 [P] [US3] Unit test for auth bypass flag in `session-agent/tests/unit/auth_test.go`
- [X] T033 [P] [US3] Integration test for non-prod auth bypass in `data-plane/tests/integration/session_agent_auth_bypass_test.go`

### Implementation for User Story 3

- [X] T034 [US3] Enforce token auth when auth mode is enforced in `session-agent/internal/api/middleware/auth.go`
- [X] T035 [US3] Disable exec/attach in Kubernetes runtime path in `data-plane/internal/runtime/session_k8s.go`
- [X] T036 [US3] Add observability for routed runtime targets in `data-plane/internal/runtime/handlers.go`

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T037 [P] Update documentation for session-agent usage in `README.md`
- [X] T038 [P] Update data-plane deployment notes for runtime images in `deploy/data-plane/deployment.yaml`
- [X] T039 [P] Validate quickstart steps in `specs/001-session-runtime-routing/quickstart.md`
- [X] T040 [P] Add concurrency validation for multi-step sessions in `data-plane/tests/integration/session_concurrency_test.go`
- [X] T041 [P] Add cross-session isolation validation in `data-plane/tests/integration/session_isolation_test.go`

---

## Phase 7: Hardening for Session Agent & Routing

**Goal**: Ensure session-agent provides long-lived per-session REPL state, enforce per-session auth, and stabilize K8s readiness + workspace handling.

### Tests for Phase 7 (MANDATORY) ⚠️

- [X] T052 [P] Unit test per-session token enforcement in `session-agent/tests/unit/auth_test.go`
- [X] T053 [P] Contract test update for session lifecycle endpoints in `session-agent/tests/contract/agent_api_test.go`
- [X] T054 [P] Unit test for K8s readiness wait/health probe in `data-plane/tests/unit/session_k8s_ready_test.go`

### Implementation for Phase 7

- [X] T043 [P] Convert runtime images to multi-stage builds to compile session-agent in `deploy/runtime/python/Dockerfile` and `deploy/runtime/node/Dockerfile`
- [X] T044 Implement per-session REPL process registry and IO locking in `session-agent/internal/runtime/runner.go`
- [X] T045 Add Python/Node REPL helpers in `session-agent/internal/runtime/repl_python.go` and `session-agent/internal/runtime/repl_node.go`
- [X] T046 Add session registration endpoint and per-session token storage in `session-agent/internal/api/router.go`, `session-agent/internal/api/handlers/sessions.go`, `session-agent/internal/runtime/runner.go`, and `shared/sessionagent/types.go`
- [X] T047 Register sessions with agent on start in `data-plane/internal/runtime/session_local.go`, `data-plane/internal/runtime/session_k8s.go`, and `data-plane/internal/runtime/agent_client.go`
- [X] T048 Wait for pod readiness and agent `/v1/health` before returning route in `data-plane/internal/runtime/session_k8s.go` and `data-plane/internal/runtime/agent_client.go`
- [X] T049 Wire `workspaceRef` into session working directory and pod volume mount in `data-plane/internal/runtime/session_local.go`, `data-plane/internal/runtime/session_k8s.go`, `data-plane/internal/workspace/session.go`, and `session-agent/internal/runtime/runner.go`
- [X] T050 Prefer agent routing for session steps when configured in `data-plane/internal/runtime/handlers.go` and `data-plane/internal/config/config.go`
- [X] T051 Add session termination endpoint and client wiring in `session-agent/internal/api/router.go`, `session-agent/internal/api/handlers/sessions.go`, `session-agent/internal/runtime/runner.go`, `data-plane/internal/runtime/agent_client.go`, and `data-plane/internal/runtime/session_local.go`
- [X] T055 [P] Update session-agent, data-plane, and runtime image docs with local and Kubernetes flow diagrams in `README.md` and `architecture/sandbox-executor.md`
- [X] T056 [P] Verify and update build/run targets for session-agent and runtime images in `Makefile`
- [X] T057 [P] Update runtime Dockerfiles for session-agent packaging in `deploy/runtime/python/Dockerfile` and `deploy/runtime/node/Dockerfile`
- [X] T058 [P] Update GitHub Actions workflows for runtime image builds and session-agent changes in `.github/workflows`
- [X] T059 Add retry with exponential backoff and `runtime_unreachable` error payload for agent step calls in `data-plane/internal/runtime/agent_client.go` and `data-plane/internal/runtime/handlers.go`
- [X] T060 Validate session runtime identity before routing steps and return a deterministic error on mismatch in `data-plane/internal/runtime/handlers.go` and `data-plane/internal/runtime/session_registry.go`
- [X] T061 Enforce header/body token consistency for session registration in `session-agent/internal/api/handlers/sessions.go` and `session-agent/internal/api/middleware/auth.go`
- [X] T062 Add explicit errors for expired sessions/runtime crashes in `data-plane/internal/runtime/handlers.go` and `data-plane/internal/runtime/session_registry.go`
- [X] T063 Cleanup k8s pod on readiness/health timeout in `data-plane/internal/runtime/session_k8s.go`
- [X] T064 Reject re-registering a session with a different runtime in `session-agent/internal/runtime/runner.go`
- [X] T065 Remove runtime from data-plane step API: update `data-plane/internal/runtime/handlers.go`, `data-plane/internal/runtime/openapi.yaml`, and `control-plane/pkg/client/data_plane_client.go` to stop requiring runtime on step requests
- [X] T066 Route agent step calls using the session route runtime (not request payload) in `data-plane/internal/runtime/handlers.go` and `data-plane/internal/runtime/session_registry.go`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - blocks all user stories
- **User Stories (Phase 3+)**: Depend on Foundational completion
- **Polish (Final Phase)**: Depends on desired user stories being complete

### User Story Dependencies

- **US1 (P1)**: Can start after Foundational
- **US2 (P2)**: Can start after Foundational
- **US3 (P3)**: Can start after Foundational

### Parallel Opportunities

- T004, T016, T017, T020, T021, T022, T027, T028, T032, T033, T037, T038, T039, T043, T052, T053, T054 can run in parallel as marked

---

## Parallel Example: User Story 1

```bash
Task: "Unit test for runner state persistence in session-agent/tests/unit/runner_test.go"
Task: "Integration test for session step persistence in data-plane/tests/integration/session_agent_local_test.go"
Task: "Contract test for agent API in session-agent/tests/contract/agent_api_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1
4. Validate state persistence and step output

### Incremental Delivery

1. Complete Setup + Foundational
2. Deliver US1 (stateful steps)
3. Deliver US2 (routing correctness)
4. Deliver US3 (no exec, auth bypass in non-prod)
