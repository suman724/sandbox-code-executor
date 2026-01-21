# Quickstart: Session Runtime Routing

## Scenario: Execute steps without remote exec

1. Create a session via control-plane.
2. Data-plane provisions a sandbox pod and waits for the session-agent health check.
3. Data-plane stores a SessionRoute entry mapping the session to the pod endpoint and token.
4. Submit a step via control-plane.
5. Data-plane forwards the step to the session-agent API and returns stdout/stderr.

## Validation

- Run two steps in the same session and verify in-memory state persists.
- Create two sessions and verify steps are routed to their own runtimes.
- Disable remote exec permissions and verify step execution still works.

## Local Non-Production Run

1. Start the session-agent locally with auth bypass enabled.
2. Start data-plane configured to route steps to the local agent endpoint.
3. Create a session and submit a step; verify output is returned without auth headers.
