import { z } from "zod";

export const apiRequestSchema = z.object({
  query: z.string().trim().max(80, "Keep search under 80 characters").optional().default(""),
  kind: z.enum(["all", "component", "platform"]).default("all"),
});

export type ApiRequest = z.infer<typeof apiRequestSchema>;

export const apiEndpoints = [
  {
    method: "GET",
    path: "/api/catalog",
    label: "Full catalog",
    description:
      "Return the versioned catalog: the components, each drawn or themed, and the platforms they are drawn for.",
    response: '{ "version": 2, "platforms": [...], "components": [...] }',
  },
  {
    method: "GET",
    path: "/api/catalog?kind=component",
    label: "Filter entries",
    description: "Filter the catalog by component or platform.",
    response: '{ "items": [{ "name": "Button", "kind": "component", "look": "drawn" }] }',
  },
  {
    method: "GET",
    path: "/api/health",
    label: "Health check",
    description: "Confirm the Breeze UI catalog is available.",
    response: '{ "status": "ok", "version": 2 }',
  },
] as const;
