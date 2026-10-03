import fs from "node:fs"
import path from "node:path"
import { fileURLToPath } from "node:url"

export type CatalogKind = "utility" | "component" | "token"

export type CatalogItem = {
  name: string
  kind: CatalogKind
  category?: string
  description: string
  value?: string
  example?: string
}

type CatalogFile = {
  version: number
  utilities: Omit<CatalogItem, "kind">[]
  components: Omit<CatalogItem, "kind">[]
  tokens: Omit<CatalogItem, "kind">[]
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
    ...parsed.utilities.map((item) => ({ ...item, kind: "utility" as const })),
    ...parsed.components.map((item) => ({ ...item, kind: "component" as const })),
    ...parsed.tokens.map((item) => ({ ...item, kind: "token" as const })),
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
