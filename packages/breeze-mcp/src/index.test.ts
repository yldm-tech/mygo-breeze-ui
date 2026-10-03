import { describe, expect, test } from "bun:test"
import { Client } from "@modelcontextprotocol/sdk/client/index.js"
import { InMemoryTransport } from "@modelcontextprotocol/sdk/inMemory.js"

import { findCatalog, loadCatalog, searchCatalog } from "./catalog"
import { createServer } from "./index"

describe("Breeze UI MCP catalog", () => {
  test("loads the shared catalog", () => {
    const catalog = loadCatalog()
    expect(catalog.some((item) => item.kind === "component" && item.name === "Button")).toBe(true)
    expect(catalog.some((item) => item.kind === "token" && item.name === "space-1")).toBe(true)
  })

  test("searches by name and description", () => {
    expect(searchCatalog("utility", "flex").map((item) => item.name)).toContain("row")
    expect(findCatalog("component", "button")?.name).toBe("Button")
  })

  test("serves tools through MCP", async () => {
    const server = createServer()
    const [clientTransport, serverTransport] = InMemoryTransport.createLinkedPair()
    await server.connect(serverTransport)
    const client = new Client({ name: "test-client", version: "0.0.0" })
    await client.connect(clientTransport)
    const tools = await client.listTools()
    expect(tools.tools.map((tool) => tool.name)).toContain("generate_snippet")
    const result = await client.callTool({ name: "generate_snippet", arguments: { component: "Button" } })
    expect(JSON.stringify(result)).toContain("breeze.Button")
    await client.close()
    await server.close()
  })
})
