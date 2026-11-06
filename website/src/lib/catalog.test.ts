import { describe, expect, it } from "vite-plus/test";
import { assembleCatalog, catalogEntries, getCatalog } from "./catalog";

describe("catalog assembled from ui/breeze/catalog", () => {
  it("lists the components in the byte order of their names, once each", async () => {
    const names = (await getCatalog()).components.map((item) => item.name);
    expect(names).toEqual([...names].sort((a, b) => (a < b ? -1 : a > b ? 1 : 0)));
    expect(new Set(names).size).toBe(names.length);
  });

  it("finds a component and a platform", async () => {
    const entries = catalogEntries(await getCatalog());
    expect(entries.some((item) => item.kind === "component" && item.name === "Button")).toBe(true);
    expect(entries.some((item) => item.kind === "platform" && item.name === "Windows")).toBe(true);
  });

  it("rejects a second version and a component listed twice", () => {
    expect(() => assembleCatalog({ "a.json": { version: 2 }, "b.json": { version: 2 } })).toThrow();
    expect(() => assembleCatalog({ "a.json": { version: 2, components: [{ name: "Button" }] }, "b.json": { components: [{ name: "Button" }] } })).toThrow();
  });
});
