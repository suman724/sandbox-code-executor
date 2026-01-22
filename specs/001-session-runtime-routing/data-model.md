# Data Model: Session Runtime Routing

## Entities

### SessionRuntime

Represents the live runtime bound to a session.

**Fields**:
- session_id (string, unique)
- runtime_id (string)
- namespace (string)
- pod_name (string)
- service_name (string, optional)
- endpoint (string, pod IP or service DNS)
- auth_token (string, per-session token)
- auth_mode (string: enforced|bypass)
- status (string: provisioning|ready|failed|terminated)
- created_at (timestamp)
- last_seen_at (timestamp)

### SessionRoute

Lookup record for routing steps to a runtime.

**Fields**:
- session_id (string, unique)
- target_endpoint (string)
- target_port (number)
- auth_token (string)
- auth_mode (string: enforced|bypass)
- runtime_status (string)
- last_updated_at (timestamp)

### StepExecution

Represents a submitted step and its outcome.

**Fields**:
- step_id (string)
- session_id (string)
- submitted_at (timestamp)
- started_at (timestamp)
- completed_at (timestamp)
- exit_code (number)
- stdout (string)
- stderr (string)
- status (string: accepted|running|completed|failed)
