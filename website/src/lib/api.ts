import { z } from "zod";

export const apiRequestSchema = z.object({
  query: z.string().trim().max(80, "Keep search under 80 characters").optional().default(""),
  kind: z.enum(["all", "component", "utility", "token"]).default("all"),
});

export type ApiRequest = z.infer<typeof apiRequestSchema>;

export const apiEndpoints = [
  {
    method: "GET",
    path: "/api/catalog",
    label: "Full catalog",
    description: "Return the versioned components, utilities, and design tokens catalog.",
    response: '{ "version": 1, "components": [...], "utilities": [...], "tokens": [...] }',
  },
  {
    method: "GET",
    path: "/api/catalog?kind=component",
    label: "Filter entries",
    description: "Filter the catalog by component, utility, or token.",
    response: '{ "items": [{ "name": "Button", "kind": "component" }] }',
  },
  {
    method: "GET",
    path: "/api/health",
    label: "Health check",
    description: "Confirm the Breeze UI catalog is available.",
    response: '{ "status": "ok", "version": 1 }',
  },
] as const;
