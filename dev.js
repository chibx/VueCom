import "dotenv/config";
import { spawn } from "node:child_process";
import { join } from "node:path";

let prev = null;

const backend = spawn("air", [], {
    cwd: join(process.cwd(), "./backend/gateway_service"),
    env: { ...process.env, FORCE_COLOR: "1" },
    stdio: ["ignore", "pipe", "pipe"],
});

const frontend = spawn("npm", ["run", "dev"], {
    cwd: join(process.cwd(), "./frontend"),
    stdio: ["ignore", "pipe", "pipe"],
});

// Helper to process lines from a stream
function pipeOutput(stream, sourceName, isError = false) {
    stream.setEncoding("utf8");
    stream.on("data", (data) => {
        const lines = data.split(/\r?\n/);
        lines.forEach((line) => {
            if (line === "") return; // Skip empty lines from split

            if (prev !== sourceName) {
                console.log(`--------------${sourceName}-------------------`);
                prev = sourceName;
            }

            if (isError) {
                console.error(line);
            } else {
                console.log(line); // Prints with colors preserved
            }
        });
    });
}

backend.once("spawn", () => {
    console.log("Backend Started Successfully");
});
pipeOutput(backend.stdout, "Backend");
pipeOutput(backend.stderr, "Backend", true);

frontend.once("spawn", () => {
    console.log("Frontend Started Successfully");
});
pipeOutput(frontend.stdout, "Frontend");
pipeOutput(frontend.stderr, "Frontend", true);


backend.once("close", (code) => {
    console.log(`Backend exited with code ${code}`);
    if (frontend.connected) frontend.kill("SIGKILL");
});

frontend.once("close", (code) => {
    console.log(`Frontend exited with code ${code}`);
    if (backend.connected) backend.kill("SIGKILL");
});
