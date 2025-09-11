import { spawn } from "node:child_process";
import { copyFileSync, cpSync, existsSync, mkdirSync, rmSync } from "node:fs";
import { platform } from "node:os";
import { join } from "node:path";
const OUTPUT_DIR = join(process.cwd(), "./.output");
const BACKEND_DIR = join(process.cwd(), "./backend");
const FRONTEND_DIR = join(process.cwd(), "./frontend");
let prev = "f"; // Switch to preserve grouping
const BINARY_NAME = `vuecom-server${platform() === "win32" ? ".exe" : ""}`;
const backend = spawn("go", ["build", "-o", `./${BINARY_NAME}`], { cwd: BACKEND_DIR });
const frontend = spawn("npm", ["run", "generate"], { cwd: FRONTEND_DIR });
backend.once("spawn", () => {
    console.log("Backend Started Successfully");
});
backend.stderr.on("data", (err) => {
    err = String(err);
    console.error("\n", "--------------Error - Backend-------------------");
    console.log(err.substring(0, err.length - 1), "\n");
});
backend.stdout.on("data", (d) => {
    d = String(d);
    if (prev !== "b") {
        console.log("--------------Backend-------------------");
    }
    console.log(d.substring(0, d.length - 1));
    prev = "b";
});
frontend.once("spawn", () => {
    console.log("Frontend Started Successfully");
});
frontend.stdout.on("data", (data) => {
    data = data + "";
    if (prev !== "f") {
        console.log("--------------Frontend-------------------");
    }
    console.log(data.substring(0, data.length - 1));
    prev = "f";
});
frontend.stderr.on("data", (err) => {
    err = String(err);
    console.log("\n", "--------------Error - Frontend-------------------");
    console.error(err.substring(0, err.length - 1), "\n");
});
/** Self Terminate */
await Promise.all([
    new Promise((resolve, reject) => {
        backend.once("close", (code) => {
            if (code !== 0) {
                console.error(
                    "--------Server Build Failed---------\nAborting all active operations!"
                );
                if (frontend.connected) {
                    frontend.kill("SIGKILL");
                }
                reject();
                return;
            }
            resolve();
        });
    }),
    new Promise((resolve, reject) => {
        frontend.once("close", (code) => {
            if (code !== 0) {
                console.error(
                    "--------Client Build Failed---------\nAborting all active operations!"
                );
                if (backend.connected) {
                    backend.kill("SIGKILL");
                }
                reject();
                return;
            }
            resolve();
        });
    }),
]).catch(() => {});

if (existsSync(OUTPUT_DIR)) {
    rmSync(OUTPUT_DIR, {
        force: true,
        recursive: true,
    });
}
mkdirSync(OUTPUT_DIR);
copyFileSync(join(BACKEND_DIR, `./${BINARY_NAME}`), join(OUTPUT_DIR, `./${BINARY_NAME}`));
cpSync(join(FRONTEND_DIR, "./.output/public"), join(OUTPUT_DIR, "./dist"), { recursive: true });
// await Promise.all([
//     rm(join(BACKEND_DIR, `./${BINARY_NAME}`)),
//     rm(join(FRONTEND_DIR, "./dist"), { recursive: true, force: true }),
// ]);
