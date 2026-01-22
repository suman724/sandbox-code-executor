#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
TENANT_ID="${TENANT_ID:-tenant-1}"
AGENT_ID="${AGENT_ID:-agent-1}"
POLICY_ID="${POLICY_ID:-policy-1}"
RUNTIME="${RUNTIME:-python}"

create_payload=$(cat <<EOF
{
  "tenantId": "${TENANT_ID}",
  "agentId": "${AGENT_ID}",
  "policyId": "${POLICY_ID}",
  "ttlSeconds": 3600,
  "runtime": "${RUNTIME}"
}
EOF
)

create_resp=$(curl -s -X POST "${BASE_URL}/sessions" \
  -H "accept: application/json" \
  -H "Content-Type: application/json" \
  -d "${create_payload}")

session_id=$(printf '%s' "${create_resp}" | python3 -c 'import json,sys; print(json.load(sys.stdin).get("id",""))')
if [ -z "${session_id}" ]; then
  echo "Failed to parse session id: ${create_resp}" >&2
  exit 1
fi

step_payload='{"command":"print(\"hello\")"}'

curl -s -X POST "${BASE_URL}/sessions/${session_id}/steps" \
  -H "accept: application/json" \
  -H "Content-Type: application/json" \
  -d "${step_payload}"

step_payload='{"command":"x=5*10"}'

curl -s -X POST "${BASE_URL}/sessions/${session_id}/steps" \
  -H "accept: application/json" \
  -H "Content-Type: application/json" \
  -d "${step_payload}"

step_payload='{"command":"print(x)"}'

curl -s -X POST "${BASE_URL}/sessions/${session_id}/steps" \
  -H "accept: application/json" \
  -H "Content-Type: application/json" \
  -d "${step_payload}"
echo
