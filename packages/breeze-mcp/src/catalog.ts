import fs from "node:fs"
import path from "node:path"
import { fileURLToPath } from "node:url"

export type CatalogKind = "component" | "platform"

export type CatalogItem = {
  name: string
  kind: CatalogKind
  category?: string
  /** "drawn" when Breeze UI draws the component as each platform does, "themed" when it is MyGo's widget in the platform's theme. */
  look?: "drawn" | "themed"
  description: string
  value?: string
  example?: string
}

export type CatalogFile = {
  version: number
  platforms: Omit<CatalogItem, "kind">[]
  components: Omit<CatalogItem, "kind">[]
}

// The catalog's files in a checkout (ui/breeze/catalog/*.json), and the catalog assembled from them that the package carries (build.ts writes it).
const sourceCatalogDir = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "../../../ui/breeze/catalog")
const packagedCatalogFile = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "../catalog.json")

/** assembleCatalog merges the files of dir as the Go package does: in the order of their names, platforms.json giving the version and the platforms and each other file the components of one group, the components in the byte order of their names. The version must be given once, and no two components or platforms may share a name. */
export function assembleCatalog(dir = sourceCatalogDir): CatalogFile {
  let version: number | undefined
  const platforms: CatalogFile["platforms"] = []
  const components: CatalogFile["components"] = []
  const seen = new Set<string>()
  const files = fs.readdirSync(dir).filter((name) => name.endsWith(".json")).sort((a, b) => (a < b ? -1 : a > b ? 1 : 0))
  for (const name of files) {
    const file = JSON.parse(fs.readFileSync(path.join(dir, name), "utf8")) as Partial<CatalogFile>
    if (file.version !== undefined) {
      if (version !== undefined) throw new Error(`${name}: the catalog's version is given twice`)
      version = file.version
    }
    for (const [kind, items] of [["platform", file.platforms ?? []], ["component", file.components ?? []]] as const) {
      for (const item of items) {
        if (seen.has(`${kind} ${item.name}`)) throw new Error(`${name}: ${kind} ${item.name} is listed twice`)
        seen.add(`${kind} ${item.name}`)
      }
    }
    platforms.push(...(file.platforms ?? []))
    components.push(...(file.components ?? []))
  }
  if (version === undefined) throw new Error(`${dir}: no file gives the catalog's version`)
  components.sort((a, b) => (a.name < b.name ? -1 : a.name > b.name ? 1 : 0))
  return { version, platforms, components }
}

function readCatalog(): CatalogFile {
  if (fs.existsSync(sourceCatalogDir)) return assembleCatalog(sourceCatalogDir)
  return JSON.parse(fs.readFileSync(packagedCatalogFile, "utf8")) as CatalogFile
}

export function loadCatalog(): CatalogItem[] {
  const parsed = readCatalog()
  return [
    ...parsed.components.map((item) => ({ ...item, kind: "component" as const })),
    ...parsed.platforms.map((item) => ({ ...item, kind: "platform" as const })),
  ].sort((a, b) => a.name.localeCompare(b.name) || a.kind.localeCompare(b.kind))
}

export function searchCatalog(kind: CatalogKind, query = "", catalog = loadCatalog()): CatalogItem[] {
  const needle = query.trim().toLowerCase()
  return catalog.filter((item) => {
    if (item.kind !== kind) return false
    if (!needle) return true
    return [item.name, item.category, item.description, item.example].filter(Boolean).some((value) => value!.toLowerCase().includes(needle))
  })
}

export function findCatalog(kind: CatalogKind, name: string, catalog = loadCatalog()): CatalogItem | undefined {
  return catalog.find((item) => item.kind === kind && item.name.toLowerCase() === name.trim().toLowerCase())
}
