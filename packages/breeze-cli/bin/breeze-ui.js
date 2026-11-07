#!/usr/bin/env node
import { spawnSync } from "node:child_process";
import { existsSync, readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const here = dirname(fileURLToPath(import.meta.url));
const cwd = process.cwd();
const checkout = existsSync(join(cwd, "go.mod")) ? cwd : join(here, "..", "..", "..");
const args = process.argv.slice(2);
const binary = process.env.BREEZE_UI_CLI_BINARY;
let result;
if (binary) {
  result = spawnSync(binary, args, { stdio: "inherit" });
} else if (existsSync(join(checkout, "go.mod")) && readFileSync(join(checkout, "go.mod"), "utf8").startsWith("module github.com/yldm-tech/mygo-breeze-ui\n")) {
  result = spawnSync("go", ["run", "./cmd/breeze-ui", ...args], { cwd: checkout, stdio: "inherit" });
} else {
  result = spawnSync("go", ["run", "github.com/yldm-tech/mygo-breeze-ui/cmd/breeze-ui@latest", ...args], { cwd, stdio: "inherit" });
}
if (result.error) {
  console.error(`breeze-ui: ${result.error.message}`);
  process.exit(1);
}
process.exit(result.status ?? 1);
