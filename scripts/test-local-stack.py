#!/usr/bin/env python3
import os
import subprocess
import sys
import urllib.request


def get_env(name, default):
    return os.environ.get(name, default)


def check_health(url, label):
    try:
        with urllib.request.urlopen(url, timeout=5) as resp:
            if resp.status != 200:
                raise RuntimeError(f"{label} returned status {resp.status}")
    except Exception as exc:
        raise RuntimeError(f"{label} health check failed: {exc}") from exc


def main():
    base_url = get_env("BASE_URL", "http://localhost:8080")
    data_plane_url = get_env("DATA_PLANE_URL", "http://localhost:8081")
    session_agent_url = get_env("SESSION_AGENT_URL", "http://localhost:9000")

    print("Checking control-plane health...")
    check_health(f"{base_url}/healthz", "control-plane")

    print("Checking data-plane health...")
    check_health(f"{data_plane_url}/healthz", "data-plane")

    print("Checking session-agent health...")
    check_health(f"{session_agent_url}/v1/health", "session-agent")

    print("Running session flow...")
    result = subprocess.run(
        ["bash", "run-session-step.sh"],
        env={**os.environ, "BASE_URL": base_url},
        check=False,
    )
    if result.returncode != 0:
        raise RuntimeError("session flow script failed")

    print("Local stack test complete.")


if __name__ == "__main__":
    try:
        main()
    except Exception as exc:
        print(str(exc), file=sys.stderr)
        sys.exit(1)
