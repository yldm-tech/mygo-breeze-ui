import fs from "node:fs"
import path from "node:path"

import { assembleCatalog } from "./src/catalog"

// The package carries the catalog assembled from ui/breeze/catalog/, which an installed package has no checkout to read; it is built here, never committed.
const target = path.resolve(import.meta.dir, "catalog.json")
fs.writeFileSync(target, JSON.stringify(assembleCatalog(path.resolve(import.meta.dir, "../../ui/breeze/catalog")), null, 2) + "\n")
const result = await Bun.build({
  entrypoints: [path.resolve(import.meta.dir, "src/index.ts")],
  outdir: path.resolve(import.meta.dir, "dist"),
  target: "node",
  format: "esm",
  minify: false,
})
if (!result.success) throw new Error("failed to build breeze-ui-mcp")
console.log(`built ${target} and dist/index.js`)
