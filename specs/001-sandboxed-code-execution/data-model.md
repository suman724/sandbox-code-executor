# Data Model: Sandboxed Code Execution Service

## Tenant
- Fields: id, name, status, created_at
- Relationships: has many Agents, Policies, Jobs, Sessions, Workflows
- Validation: status in {active, suspended}

## Agent
- Fields: id, tenant_id, name, type, status, created_at
- Relationships: belongs to Tenant; submits Jobs, Sessions, Workflows
- Validation: status in {active, disabled}

## Policy
- Fields: id, tenant_id, name, version, ruleset, created_at, updated_at
- Relationships: belongs to Tenant; applies to Jobs, Sessions, Services
- Validation: version monotonic per policy name

## Job
- Fields: id, tenant_id, agent_id, policy_id, language, status, requested_at,
  started_at, completed_at, exit_status, resource_usage, output_refs
- Relationships: belongs to Tenant, Agent, Policy; has many Artifacts
- State: queued -> running -> succeeded | failed | canceled
 - Integration: control plane creates Job and invokes data plane to provision a
   SandboxInstance for execution.

## Session
- Fields: id, tenant_id, agent_id, policy_id, status, ttl_seconds, created_at,
  expires_at, last_activity_at
- Relationships: belongs to Tenant, Agent, Policy; has many SessionSteps
- State: active -> expired | terminated
 - Integration: control plane creates Session and requests the data plane to
   allocate a SandboxInstance for the session TTL.

## SessionStep
- Fields: id, session_id, sequence, command, status, started_at, completed_at
- Relationships: belongs to Session; produces Artifacts
- Validation: sequence increments strictly per session

## Artifact
- Fields: id, tenant_id, owner_type, owner_id, name, size_bytes, checksum,
  storage_uri, created_at
- Relationships: belongs to Tenant; owned by Job, SessionStep, or WorkflowStep
- Validation: size_bytes <= tenant policy limits

## Workflow
- Fields: id, tenant_id, status, created_at, started_at, completed_at
- Relationships: belongs to Tenant; has many WorkflowSteps; references Agents
- State: queued -> running -> succeeded | failed | canceled

## WorkflowStep
- Fields: id, workflow_id, sequence, agent_id, job_id, status, started_at,
  completed_at
- Relationships: belongs to Workflow; references Job; produces Artifacts
- Validation: sequence increments strictly per workflow

## Service
- Fields: id, tenant_id, policy_id, status, started_at, stopped_at, proxy_url
- Relationships: belongs to Tenant, Policy; produces AuditEvents
- State: starting -> running -> stopped | terminated

## RunnerNode
- Fields: id, status, zone, capacity, last_heartbeat_at
- Relationships: hosts SandboxInstances
- Validation: status in {healthy, degraded, offline}

## SandboxInstance
- Fields: id, job_id, session_id, runner_id, status, created_at, terminated_at
- Relationships: belongs to Job or Session; hosted on RunnerNode
- State: provisioning -> running -> terminated | failed

## AuditEvent
- Fields: id, tenant_id, actor_id, action, resource_type, resource_id,
  timestamp, outcome, details
- Relationships: belongs to Tenant; references Job, Session, Service, Workflow
- Validation: immutable after write
