#!/usr/bin/env python3
import os
import subprocess
import sys
import time


def get_env(name, default=""):
    return os.environ.get(name, default)


def require_env(name):
    value = os.environ.get(name, "")
    if not value:
        raise RuntimeError(f"Missing {name}. Set {name} environment variable.")
    return value


def main():
    namespace = get_env("NAMESPACE", "default")
    control_plane_service = require_env("CONTROL_PLANE_SERVICE")
    data_plane_service = require_env("DATA_PLANE_SERVICE")
    local_control_plane_port = get_env("LOCAL_CONTROL_PLANE_PORT", "8080")
    local_data_plane_port = get_env("LOCAL_DATA_PLANE_PORT", "8081")

    cp_cmd = [
        "kubectl",
        "-n",
        namespace,
        "port-forward",
        f"svc/{control_plane_service}",
        f"{local_control_plane_port}:8080",
    ]
    dp_cmd = [
        "kubectl",
        "-n",
        namespace,
        "port-forward",
        f"svc/{data_plane_service}",
        f"{local_data_plane_port}:8081",
    ]

    cp_proc = subprocess.Popen(cp_cmd, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    dp_proc = subprocess.Popen(dp_cmd, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    try:
        time.sleep(2)
        base_url = f"http://localhost:{local_control_plane_port}"
        result = subprocess.run(
            ["bash", "scripts/run-session-step.sh"],
            env={**os.environ, "BASE_URL": base_url},
            check=False,
        )
        if result.returncode != 0:
            raise RuntimeError("session flow script failed")
    finally:
        for proc in (cp_proc, dp_proc):
            proc.terminate()
            try:
                proc.wait(timeout=5)
            except subprocess.TimeoutExpired:
                proc.kill()

    print("Kubernetes stack test complete.")


if __name__ == "__main__":
    try:
        main()
    except Exception as exc:
        print(str(exc), file=sys.stderr)
        sys.exit(1)
