import { describe, expect, it } from "vite-plus/test";
import { apiRequestSchema } from "./api";

describe("catalog API request validation", () => {
  it("fills defaults for a catalog query", () => {
    expect(apiRequestSchema.parse({})).toEqual({ query: "", kind: "all" });
  });

  it("accepts a supported filter", () => {
    expect(apiRequestSchema.parse({ query: "button", kind: "component" })).toEqual({
      query: "button",
      kind: "component",
    });
  });

  it("accepts the platform filter", () => {
    expect(apiRequestSchema.parse({ kind: "platform" })).toEqual({ query: "", kind: "platform" });
  });

  it("rejects kinds the catalog no longer has", () => {
    expect(apiRequestSchema.safeParse({ kind: "utility" }).success).toBe(false);
    expect(apiRequestSchema.safeParse({ kind: "token" }).success).toBe(false);
  });

  it("rejects oversized search input", () => {
    expect(apiRequestSchema.safeParse({ query: "x".repeat(81) }).success).toBe(false);
  });
});
