import "dotenv/config";
import { spawn } from "node:child_process";
import { join } from "node:path";
import process from "node:process";

const backendDir = join(process.cwd(), "./backend");
const services = ["analytics", "catalog", "email", "gateway", "inventory", "orders", "payment"];

const golangPackagesDirs = [...services.map((service) => `services/${service}`), "shared"].map(
    (dir) => join(backendDir, dir)
);

const childProcesses = golangPackagesDirs.map((dir) => {
    return new Promise<void>((resolve) => {
        const child = spawn("go", ["mod", "download"], {
            cwd: dir,
            env: { ...process.env, FORCE_COLOR: "3" },
            stdio: ["ignore", "pipe", "pipe"],
        });

        child.once("exit", (code) => {
            if (code !== null && code !== 0) {
                console.log(`Package dependencies install for ${dir} failed with code ${code}`);
            } else {
                console.log(`Package dependencies install for ${dir} completed successfully`);
            }
            resolve();
        });
    });
});

await Promise.all(childProcesses);
