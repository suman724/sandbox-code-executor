import http from "k6/http";
import { check, sleep } from "k6";

export const options = { vus: 1, duration: "5s" };

export default function () {
  const res = http.post("http://localhost:8080/sessions", JSON.stringify({ id: "session-1" }), {
    headers: { "Content-Type": "application/json" },
  });
  check(res, { "status is 201": (r) => r.status === 201 });
  sleep(1);
}
