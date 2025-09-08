"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
require("dotenv/config");
var node_child_process_1 = require("node:child_process");
var node_path_1 = require("node:path");
var prev = "f"; // Switch to preserve grouping
var backend = (0, node_child_process_1.spawn)("gow", ["run", "."], {
    cwd: (0, node_path_1.join)(process.cwd(), "./backend"),
});
var frontend = (0, node_child_process_1.spawn)("npm", ["run", "dev"], {
    cwd: (0, node_path_1.join)(process.cwd(), "./frontend"),
});
backend.once("spawn", function () {
    console.log("Backend Started Successfully");
});
backend.stderr.on("data", function (err) {
    err = String(err);
    console.error("\n", "--------------Error - Backend-------------------");
    console.log(err.substring(0, err.length - 1), "\n");
});
backend.stdout.on("data", function (d) {
    d = String(d);
    if (prev !== "b") {
        console.log("--------------Backend-------------------");
    }
    console.log(d.substring(0, d.length - 1));
    prev = "b";
});
frontend.once("spawn", function () {
    console.log("Frontend Started Successfully");
});
frontend.stdout.on("data", function (data) {
    data = data + "";
    if (prev !== "f") {
        console.log("--------------Frontend-------------------");
    }
    console.log(data.substring(0, data.length - 1));
    prev = "f";
});
frontend.stderr.on("data", function (err) {
    err = String(err);
    console.log("\n", "--------------Error - Frontend-------------------");
    console.error(err.substring(0, err.length - 1), "\n");
});
/** Self Terminate */
backend.once("close", function () {
    if (frontend.connected) {
        frontend.kill("SIGKILL");
    }
});
frontend.once("close", function () {
    if (backend.connected) {
        backend.kill("SIGKILL");
    }
});
