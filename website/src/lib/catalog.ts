import { z } from "zod";

// "drawn": Breeze UI draws the widget as the platform does; "themed": MyGo's own widget in the platform's theme.
export const lookSchema = z.enum(["drawn", "themed"]);

const platformSchema = z.object({
  name: z.string(),
  description: z.string(),
  value: z.string(),
  example: z.string(),
});

const componentSchema = z.object({
  name: z.string(),
  category: z.string(),
  look: lookSchema,
  description: z.string(),
  example: z.string(),
});

export const catalogSchema = z.object({
  version: z.number(),
  platforms: z.array(platformSchema),
  components: z.array(componentSchema),
});

export type CatalogLook = z.infer<typeof lookSchema>;
export type CatalogPlatform = z.infer<typeof platformSchema>;
export type CatalogComponent = z.infer<typeof componentSchema>;
export type Catalog = z.infer<typeof catalogSchema>;
export type CatalogKind = "component" | "platform";

// One row of the catalog regardless of kind, as the catalog grid and table list them.
export type CatalogEntry = {
  kind: CatalogKind;
  name: string;
  description: string;
  example: string;
  category?: string;
  look?: CatalogLook;
  value?: string;
};

export function catalogEntries(catalog: Catalog): CatalogEntry[] {
  return [
    ...catalog.components.map((item) => ({ ...item, kind: "component" as const })),
    ...catalog.platforms.map((item) => ({ ...item, kind: "platform" as const })),
  ];
}

type CatalogFile = { version?: number; platforms?: unknown[]; components?: { name: string }[] };

// The catalog's files, ui/breeze/catalog/*.json: platforms.json gives the version and the platforms, each other file the components of one group of widgets.
const catalogFiles = import.meta.glob<CatalogFile>("../../../ui/breeze/catalog/*.json", { eager: true, import: "default" });

// assembleCatalog merges the files as the Go package does: in the order of their paths, the components in the byte order of their names, which is the order the grid and the table list them in. The version must be given once, and no two components may share a name.
export function assembleCatalog(files: Record<string, CatalogFile>) {
  let version: number | undefined;
  const platforms: unknown[] = [];
  const components: { name: string }[] = [];
  for (const path of Object.keys(files).sort((a, b) => (a < b ? -1 : a > b ? 1 : 0))) {
    const file = files[path]!;
    if (file.version !== undefined) {
      if (version !== undefined) throw new Error(`${path}: the catalog's version is given twice`);
      version = file.version;
    }
    platforms.push(...(file.platforms ?? []));
    components.push(...(file.components ?? []));
  }
  components.sort((a, b) => (a.name < b.name ? -1 : a.name > b.name ? 1 : 0));
  for (let i = 1; i < components.length; i++) {
    if (components[i]!.name === components[i - 1]!.name) throw new Error(`component ${components[i]!.name} is listed twice`);
  }
  return { version, platforms, components };
}

const catalog = catalogSchema.parse(assembleCatalog(catalogFiles));

export async function getCatalog(): Promise<Catalog> {
  return catalog;
}

export const catalogCounts = {
  components: catalog.components.length,
  platforms: catalog.platforms.length,
  drawn: catalog.components.filter((item) => item.look === "drawn").length,
};
