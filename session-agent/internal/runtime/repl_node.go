package runtime

const nodeReplScript = `const readline = require("readline");
const vm = require("vm");

const context = vm.createContext({
  console,
  require,
  process,
  Buffer,
  setTimeout,
  setInterval,
  clearTimeout,
  clearInterval,
});

const rl = readline.createInterface({
  input: process.stdin,
  crlfDelay: Infinity,
});

rl.on("line", (line) => {
  if (!line) {
    return;
  }
  let req;
  try {
    req = JSON.parse(line);
  } catch (err) {
    return;
  }
  let stdout = "";
  let stderr = "";
  let error = "";

  const originalStdoutWrite = process.stdout.write.bind(process.stdout);
  const originalStderrWrite = process.stderr.write.bind(process.stderr);
  const originalConsoleLog = console.log;
  const originalConsoleError = console.error;

  process.stdout.write = (chunk, encoding, cb) => {
    stdout += chunk instanceof Buffer ? chunk.toString() : chunk;
    if (typeof cb === "function") {
      cb();
    }
    return true;
  };
  process.stderr.write = (chunk, encoding, cb) => {
    stderr += chunk instanceof Buffer ? chunk.toString() : chunk;
    if (typeof cb === "function") {
      cb();
    }
    return true;
  };
  console.log = (...args) => {
    stdout += args.join(" ") + "\n";
  };
  console.error = (...args) => {
    stderr += args.join(" ") + "\n";
  };

  try {
    vm.runInContext(req.code || "", context);
  } catch (err) {
    error = err && err.stack ? err.stack : String(err);
  }

  process.stdout.write = originalStdoutWrite;
  process.stderr.write = originalStderrWrite;
  console.log = originalConsoleLog;
  console.error = originalConsoleError;

  const resp = JSON.stringify({ stdout, stderr, error });
  originalStdoutWrite(resp + "\n");
});
`
