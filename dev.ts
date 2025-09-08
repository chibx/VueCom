import "dotenv/config";
import { spawn } from "node:child_process";
import { join } from "node:path";

let prev: "f" | "b" = "f"; // Switch to preserve grouping

const backend = spawn("gow", ["run", "."], {
  cwd: join(process.cwd(), "./backend"),
});
const frontend = spawn("npm", ["run", "dev"], {
  cwd: join(process.cwd(), "./frontend"),
});

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
backend.once("close", () => {
  if (frontend.connected) {
    frontend.kill("SIGKILL");
  }
});

frontend.once("close", () => {
  if (backend.connected) {
    backend.kill("SIGKILL");
  }
});
