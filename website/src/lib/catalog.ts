import rawCatalog from "../../../ui/breeze/catalog.json";
import { z } from "zod";

const catalogItemSchema = z.object({
  name: z.string(),
  description: z.string().optional(),
  example: z.string().optional(),
  category: z.string().optional(),
  kind: z.string().optional(),
  value: z.string().optional(),
});

export const catalogSchema = z.object({
  version: z.number(),
  utilities: z.array(catalogItemSchema),
  components: z.array(catalogItemSchema),
  tokens: z.array(catalogItemSchema),
});

export type CatalogItem = z.infer<typeof catalogItemSchema>;
export type Catalog = z.infer<typeof catalogSchema>;

const catalog = catalogSchema.parse(rawCatalog);

export async function getCatalog(): Promise<Catalog> {
  return catalog;
}

export const catalogCounts = {
  components: catalog.components.length,
  utilities: catalog.utilities.length,
  tokens: catalog.tokens.length,
};
