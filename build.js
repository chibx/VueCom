"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var node_child_process_1 = require("node:child_process");
var node_fs_1 = require("node:fs");
var promises_1 = require("node:fs/promises");
var node_os_1 = require("node:os");
var node_path_1 = require("node:path");
var OUTPUT_DIR = (0, node_path_1.join)(process.cwd(), "./.output");
var BACKEND_DIR = (0, node_path_1.join)(process.cwd(), "./backend");
var FRONTEND_DIR = (0, node_path_1.join)(process.cwd(), "./frontend");
var prev = "f"; // Switch to preserve grouping
var BINARY_NAME = "vuecom-server".concat((0, node_os_1.platform)() === "win32" ? ".exe" : "");
var backend = (0, node_child_process_1.spawn)("go", ["build", "-o", "./".concat(BINARY_NAME)], { cwd: BACKEND_DIR });
var frontend = (0, node_child_process_1.spawn)("npm", ["run", "build"], { cwd: FRONTEND_DIR });
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
await Promise.all([
    new Promise(function (resolve, reject) {
        backend.once("close", function (code) {
            if (code !== 0) {
                console.error("--------Server Build Failed---------\nAborting all active operations!");
                if (frontend.connected) {
                    frontend.kill("SIGKILL");
                }
                reject();
                return;
            }
            resolve();
        });
    }),
    new Promise(function (resolve, reject) {
        frontend.once("close", function (code) {
            if (code !== 0) {
                console.error("--------Client Build Failed---------\nAborting all active operations!");
                if (backend.connected) {
                    backend.kill("SIGKILL");
                }
                reject();
                return;
            }
            resolve();
        });
    }),
]).catch(function () { });
if ((0, node_fs_1.existsSync)(OUTPUT_DIR)) {
    (0, node_fs_1.rmSync)(OUTPUT_DIR, {
        force: true,
        recursive: true,
    });
}
(0, node_fs_1.mkdirSync)(OUTPUT_DIR);
(0, node_fs_1.copyFileSync)((0, node_path_1.join)(BACKEND_DIR, "./".concat(BINARY_NAME)), (0, node_path_1.join)(OUTPUT_DIR, "./".concat(BINARY_NAME)));
(0, node_fs_1.cpSync)((0, node_path_1.join)(FRONTEND_DIR, "./dist"), (0, node_path_1.join)(OUTPUT_DIR, "./dist"), { recursive: true });
await Promise.all([
    (0, promises_1.rm)((0, node_path_1.join)(BACKEND_DIR, "./".concat(BINARY_NAME))),
    (0, promises_1.rm)((0, node_path_1.join)(FRONTEND_DIR, "./dist"), { recursive: true, force: true }),
]);
