import "dotenv/config";
import { spawn } from "node:child_process";
import { platform } from "node:os";
import { join } from "node:path";

process.env.SERVER_MODE = "production";

const backend = spawn(`./vuecom-server${platform() === "win32" ? ".exe" : ""}`, {
    cwd: join(process.cwd(), "./.output"),
});

backend.once("spawn", () => {
    console.log("Server Started Successfully");
});
backend.stderr.on("data", (err) => {
    err = String(err);
    console.error("\n", "--------------Error - Backend-------------------");
    console.log(err.substring(0, err.length - 1), "\n");
});

backend.stdout.on("data", (d) => {
    d = String(d);
    console.log(d.substring(0, d.length - 1));
});

/** Self Terminate */
backend.once("close", () => {
    process.exit();
});
