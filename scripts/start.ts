import "dotenv/config";
import { spawn } from "node:child_process";
import { platform } from "node:os";
import { join } from "node:path";

process.env.SERVER_MODE = "production";

const BINARY_NAME = `./gateway${platform() === "win32" ? ".exe" : ""}`;

const backend = spawn(BINARY_NAME, [], {
    cwd: join(process.cwd(), "./.output/gateway"),
    env: { ...process.env, FORCE_COLOR: "3" },
    stdio: ["ignore", "pipe", "pipe"],
});

backend.once("spawn", () => {
    console.log("Server Started Successfully");
});

let buffer = "";

backend.stdout.setEncoding("utf8");
backend.stdout.on("data", (chunk: string) => {
    buffer += chunk;
    const lines = buffer.split(/\r?\n/);
    buffer = lines.pop() || "";

    for (const line of lines) {
        if (line !== "") {
            console.log(line);
        }
    }
});

backend.stdout.on("end", () => {
    if (buffer.trim() !== "") {
        console.log(buffer.trim());
    }
});

backend.stderr.setEncoding("utf8");
backend.stderr.on("data", (chunk: string) => {
    buffer += chunk;
    const lines = buffer.split(/\r?\n/);
    buffer = lines.pop() || "";

    for (const line of lines) {
        if (line !== "") {
            console.error(line);
        }
    }
});

backend.stderr.on("end", () => {
    if (buffer.trim() !== "") {
        console.error(buffer.trim());
    }
});

backend.once("close", (code) => {
    if (code !== null) {
        console.log(`Server exited with code ${code}`);
    }
    process.exit(code ?? 0);
});
