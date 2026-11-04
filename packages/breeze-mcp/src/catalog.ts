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

type CatalogFile = {
  version: number
  platforms: Omit<CatalogItem, "kind">[]
  components: Omit<CatalogItem, "kind">[]
}

const sourceCatalogFile = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "../../../ui/breeze/catalog.json")
const packagedCatalogFile = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "../catalog.json")

function catalogPath() {
  return fs.existsSync(sourceCatalogFile) ? sourceCatalogFile : packagedCatalogFile
}

export function loadCatalog(): CatalogItem[] {
  const source = fs.readFileSync(catalogPath(), "utf8")
  const parsed = JSON.parse(source) as CatalogFile
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
