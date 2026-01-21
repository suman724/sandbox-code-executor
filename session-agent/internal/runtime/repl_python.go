package runtime

const pythonReplScript = `import contextlib
import io
import json
import sys
import traceback

globals_ns = {"__name__": "__main__"}

for line in sys.stdin:
    line = line.rstrip("\n")
    if not line:
        continue
    try:
        req = json.loads(line)
    except Exception:
        continue
    code = req.get("code", "")
    out = io.StringIO()
    err = io.StringIO()
    failure = ""
    try:
        with contextlib.redirect_stdout(out), contextlib.redirect_stderr(err):
            exec(code, globals_ns)
    except Exception:
        failure = traceback.format_exc()
    resp = {"stdout": out.getvalue(), "stderr": err.getvalue(), "error": failure}
    sys.stdout.write(json.dumps(resp) + "\n")
    sys.stdout.flush()
`
