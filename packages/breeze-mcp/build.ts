import fs from "node:fs"
import path from "node:path"

const source = path.resolve(import.meta.dir, "../../ui/breeze/catalog.json")
const target = path.resolve(import.meta.dir, "catalog.json")
fs.copyFileSync(source, target)
const result = await Bun.build({
  entrypoints: [path.resolve(import.meta.dir, "src/index.ts")],
  outdir: path.resolve(import.meta.dir, "dist"),
  target: "node",
  format: "esm",
  minify: false,
})
if (!result.success) throw new Error("failed to build breeze-ui-mcp")
console.log(`built ${target} and dist/index.js`)
