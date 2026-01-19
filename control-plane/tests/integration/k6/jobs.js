import http from "k6/http";
import { check, sleep } from "k6";

export const options = { vus: 2, duration: "10s" };

const baseUrl = __ENV.BASE_URL || "http://localhost:8080";
const authToken = __ENV.AUTH_TOKEN;

export default function () {
  const payload = {
    tenantId: "tenant-1",
    agentId: "agent-1",
    policyId: "policy-1",
    language: "python",
    code: "print('ok')",
  };
  const headers = { "Content-Type": "application/json" };
  if (authToken) {
    headers.Authorization = `Bearer ${authToken}`;
  }
  const res = http.post(`${baseUrl}/jobs`, JSON.stringify(payload), { headers });
  check(res, { "status is 202": (r) => r.status === 202 });
  sleep(1);
}
