# Feature Specification: Session Runtime Routing

**Feature Branch**: `001-session-runtime-routing`  
**Created**: 2026-01-20  
**Status**: Draft  
**Input**: User description: "The current implementation makes assumption that Kubernetes allows exec into pods. It is not allowed in our enviornment due to security reasons. We need to do the following - Preserve session state across steps (REPL-like behavior). - Route steps to the correct session runtime. - Avoid `kubectl exec`/SPDY attach in production."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Preserve Session State Across Steps (Priority: P1)

As a platform user, I want to submit multiple steps to the same session and have in-memory state preserved between steps so that iterative workflows behave consistently.

**Why this priority**: State preservation is the core value of sessions; without it, multi-step workflows are unreliable.

**Independent Test**: Create a session, run a step that defines a variable, then run a second step that uses that variable and verify the expected output is returned.

**Test Requirements**: Must include unit tests for session-agent runner behavior and integration tests for multi-step session execution.

**Acceptance Scenarios**:

1. **Given** an active session, **When** a step defines an in-memory variable, **Then** a subsequent step can reference it without redefinition.
2. **Given** an active session, **When** a step writes a file to the session workspace, **Then** a subsequent step can read that file.

---

### User Story 2 - Route Steps to the Correct Runtime (Priority: P2)

As a platform operator, I need steps to be routed to the correct session runtime so outputs and state belong to the intended session.

**Why this priority**: Incorrect routing breaks isolation and state consistency across sessions.

**Independent Test**: Create two sessions, run distinct steps in each, and verify outputs and state do not cross between sessions.

**Test Requirements**: Must include unit tests for routing/registry logic and integration tests that verify cross-session isolation.

**Acceptance Scenarios**:

1. **Given** two active sessions, **When** steps are sent to each session, **Then** each step executes against its own session runtime.
2. **Given** a session runtime restart, **When** a new step is sent, **Then** the system detects the mismatch and returns a clear error.

---

### User Story 3 - Operate Without Remote Exec Access (Priority: P3)

As a security-conscious operator, I need session execution to work without requiring remote shell/exec access into runtime instances.

**Why this priority**: The production environment prohibits remote exec access, so this is required for deployment.

**Independent Test**: Configure production mode and run steps without any remote exec permissions; verify execution succeeds using the supported channel.

**Test Requirements**: Must include unit tests for auth bypass and integration tests proving production mode avoids exec/attach.

**Acceptance Scenarios**:

1. **Given** production mode with remote exec access disabled, **When** a step is submitted, **Then** it executes successfully using the approved execution channel.
2. **Given** production mode, **When** remote exec access is unavailable, **Then** the system does not attempt to use it.
3. **Given** non-production mode, **When** a step is submitted, **Then** session-agent authentication can be bypassed using an explicit configuration flag.

---

### Edge Cases

- What happens when a session runtime crashes between steps?
- How does the system handle steps submitted to an expired session?
- If a session runtime is healthy but unreachable, the system retries up to 3 times with exponential backoff; if still unreachable, it returns `runtime_unreachable` with `retryable=false`.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST preserve in-memory session state across steps within an active session.
- **FR-002**: System MUST preserve file-based workspace state across steps within an active session.
- **FR-003**: System MUST route each step to the correct session runtime based on the session identifier.
- **FR-004**: System MUST prevent steps from being routed to a different session runtime.
- **FR-005**: System MUST return a deterministic error when a session runtime is unavailable.
  - Error includes: status code, error code `runtime_unreachable`, and a retryable flag.
  - Retry behavior: up to 3 retries with exponential backoff, then fail.
- **FR-006**: System MUST operate without requiring remote exec or attach access into runtime instances in production.
- **FR-007**: System MUST log routing decisions with fields: session_id, runtime_id, target_endpoint, auth_mode, step_id, and outcome status.
- **FR-008**: System MUST host a session-agent application at the repository root that is packaged into Python and Node runtime pods.
- **FR-009**: System MUST support local, non-production execution where the session-agent can bypass authentication via an explicit configuration flag.
- **FR-010**: Each user story MUST be validated by both unit and integration tests.

### Key Entities *(include if feature involves data)*

- **Session Runtime**: The execution context that holds in-memory and workspace state for a session.
- **Session Route**: The mapping between a session identifier and its runtime target.
- **Step Execution**: A single command or code submission executed against a session runtime.

## Assumptions & Dependencies

- Session identifiers are unique and stable for the lifetime of a session.
- The platform provides a reliable way to address a session runtime without remote exec access.
- Operators can access logs or telemetry for troubleshooting routing decisions.
- Runtime images can include additional binaries required by the session-agent.
- Auth bypass is only enabled when an explicit non-production environment flag is set; it must be disabled by default.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 95% of multi-step sessions complete with correct state preservation on the first attempt.
- **SC-002**: Steps for concurrent sessions are routed correctly with zero cross-session state leakage in validation tests.
- **SC-003**: Sessions can execute steps in production without requiring remote exec access.
- **SC-004**: 90% of operator troubleshooting cases can identify the session runtime target from logs or telemetry.
