import { describe, expect, test } from "bun:test"
import { Client } from "@modelcontextprotocol/sdk/client/index.js"
import { InMemoryTransport } from "@modelcontextprotocol/sdk/inMemory.js"

import { findCatalog, loadCatalog, searchCatalog } from "./catalog"
import { createServer } from "./index"

describe("Breeze UI MCP catalog", () => {
  test("loads the shared catalog", () => {
    const catalog = loadCatalog()
    expect(catalog.some((item) => item.kind === "component" && item.name === "Button")).toBe(true)
    expect(catalog.some((item) => item.kind === "platform" && item.name === "Windows")).toBe(true)
    expect(catalog.find((item) => item.name === "Switch")?.look).toBe("drawn")
  })

  test("searches by name and description", () => {
    expect(searchCatalog("component", "drop-down").map((item) => item.name)).toContain("Select")
    expect(searchCatalog("platform", "libadwaita").map((item) => item.name)).toEqual(["Linux"])
    expect(findCatalog("component", "button")?.name).toBe("Button")
  })

  test("serves tools through MCP", async () => {
    const server = createServer()
    const [clientTransport, serverTransport] = InMemoryTransport.createLinkedPair()
    await server.connect(serverTransport)
    const client = new Client({ name: "test-client", version: "0.0.0" })
    await client.connect(clientTransport)
    const tools = await client.listTools()
    expect(tools.tools.map((tool) => tool.name)).toEqual(expect.arrayContaining(["generate_snippet", "list_platforms"]))
    expect(tools.tools.map((tool) => tool.name)).not.toContain("search_utilities")
    const result = await client.callTool({ name: "generate_snippet", arguments: { component: "Button" } })
    expect(JSON.stringify(result)).toContain("breeze.Button")
    await client.close()
    await server.close()
  })
})
