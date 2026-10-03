#!/usr/bin/env node

import { spawn } from "node:child_process"
import { fileURLToPath } from "node:url"
import path from "node:path"

const entry = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "../dist/index.js")
const child = spawn(process.execPath, [entry], { stdio: "inherit", env: process.env })
child.on("exit", (code, signal) => {
  if (signal) process.kill(process.pid, signal)
  else process.exit(code ?? 1)
})
