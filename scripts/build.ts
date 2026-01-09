import { spawn } from "node:child_process";
import { copyFileSync, cpSync, existsSync, mkdirSync, rmSync } from "node:fs";
import { platform } from "node:os";
import { join } from "node:path";

const OUTPUT_DIR = join(process.cwd(), "./.output");
const BACKEND_DIR = join(process.cwd(), "./backend/services/gateway_service");
const FRONTEND_DIR = join(process.cwd(), "./frontend");

const BINARY_NAME = `gateway${platform() === "win32" ? ".exe" : ""}`;

let prev: string | null = null;

const colorEnv = { ...process.env, FORCE_COLOR: "3" }; // 3 = full 24-bit color

const backend = spawn("make", ["build"], {
    cwd: BACKEND_DIR,
    env: colorEnv,
    stdio: ["ignore", "pipe", "pipe"],
});

const frontend = spawn("npm", ["run", "build"], {
    cwd: FRONTEND_DIR,
    env: colorEnv,
    stdio: ["ignore", "pipe", "pipe"],
});

// Reusable function to pipe and prefix output with grouping
function pipeOutput(stream: NodeJS.ReadableStream, sourceName: string, isError = false) {
    stream.setEncoding("utf8");

    let buffer = "";
    stream.on("data", (chunk) => {
        buffer += chunk;
        const lines = buffer.split(/\r?\n/);
        buffer = lines.pop() || "";

        for (const line of lines) {
            if (line === "") continue;

            if (prev !== sourceName) {
                console.log(`--------------${sourceName}-------------------`);
                prev = sourceName;
            }

            if (isError) {
                console.error(line);
            } else {
                console.log(line);
            }
        }
    });

    // Flush any remaining buffered data on close
    stream.on("end", () => {
        if (buffer.trim() !== "") {
            if (prev !== sourceName) {
                console.log(`--------------${sourceName}-------------------`);
                prev = sourceName;
            }
            if (isError) console.error(buffer.trim());
            else console.log(buffer.trim());
        }
    });
}

backend.once("spawn", () => console.log("Backend Build Started"));
frontend.once("spawn", () => console.log("Frontend Build Started"));

pipeOutput(backend.stdout, "Backend");
pipeOutput(backend.stderr, "Backend", true);

pipeOutput(frontend.stdout, "Frontend");
pipeOutput(frontend.stderr, "Frontend", true);

try {
    await Promise.all([
        new Promise<void>((resolve, reject) => {
            backend.on("close", (code) => {
                if (code !== 0) {
                    console.error("--------Server Build Failed---------");
                    reject(new Error(`Go build exited with code ${code}`));
                } else {
                    console.log("---------Server Build Completed----------");
                    resolve();
                }
            });
        }),
        new Promise<void>((resolve, reject) => {
            frontend.on("close", (code) => {
                if (code !== 0) {
                    console.error("--------Client Build Failed---------");
                    reject(new Error(`npm build exited with code ${code}`));
                } else {
                    console.log("---------Client Build Completed----------");
                    resolve();
                }
            });
        }),
    ]);

    console.log("Both builds completed successfully!");

    if (existsSync(OUTPUT_DIR)) {
        rmSync(OUTPUT_DIR, { recursive: true, force: true });
    }
    mkdirSync(OUTPUT_DIR, { recursive: true });

    const binarySrc = join(BACKEND_DIR, "bin", BINARY_NAME);
    const binaryDest = join(OUTPUT_DIR, BINARY_NAME);
    copyFileSync(binarySrc, binaryDest);

    const frontendDistSrc = join(FRONTEND_DIR, "dist");
    const frontendDistDest = join(OUTPUT_DIR, "dist");
    cpSync(frontendDistSrc, frontendDistDest, { recursive: true });

    console.log(`Build artifacts copied to ${OUTPUT_DIR}`);

    // rmSync(binarySrc, { force: true });
    // rmSync(frontendDistSrc, { recursive: true, force: true });
} catch (err) {
    console.error("\nBuild process aborted due to failure:", err);
    process.exit(1);
}
