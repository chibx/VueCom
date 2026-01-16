import "dotenv/config";
import { spawn } from "node:child_process";
import { join } from "node:path";

let prev: string | null = null;

const backend = spawn("air", [], {
    cwd: join(process.cwd(), "./backend/services/gateway_service"),
    env: { ...process.env, FORCE_COLOR: "1" },
    stdio: ["ignore", "pipe", "pipe"],
});

const frontend = spawn("npm", ["run", "dev"], {
    cwd: join(process.cwd(), "./frontend"),
    stdio: ["ignore", "pipe", "pipe"],
});

// Helper to process lines from a stream
function pipeOutput(stream: NodeJS.ReadableStream, sourceName: string, isError = false) {
    stream.setEncoding("utf8");
    stream.on("data", (data: string) => {
        const lines = data.split(/\r?\n/);
        lines.forEach((line: string) => {
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
