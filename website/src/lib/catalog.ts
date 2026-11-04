import rawCatalog from "../../../ui/breeze/catalog.json";
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

const catalog = catalogSchema.parse(rawCatalog);

export async function getCatalog(): Promise<Catalog> {
  return catalog;
}

export const catalogCounts = {
  components: catalog.components.length,
  platforms: catalog.platforms.length,
  drawn: catalog.components.filter((item) => item.look === "drawn").length,
};
