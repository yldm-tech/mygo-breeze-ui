import { cp, mkdir, rm } from "node:fs/promises"
import { join } from "node:path"

const root = import.meta.dir
const out = join(root, "dist")
await rm(out, { recursive: true, force: true })
await mkdir(join(out, "assets"), { recursive: true })
const result = await Bun.build({
  entrypoints: [join(root, "src/main.ts")],
  outdir: join(out, "assets"),
  target: "browser",
  minify: true,
  sourcemap: "none",
})
if (!result.success) throw new Error("website build failed")
await cp(join(root, "index.html"), join(out, "index.html"))
await cp(join(root, "src/styles.css"), join(out, "styles.css"))
await cp(join(root, "public"), join(out), { recursive: true })
console.log(`built ${out}`)
