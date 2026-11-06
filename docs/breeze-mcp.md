# Breeze UI MCP server

`packages/breeze-mcp` exposes the Breeze UI catalog to MCP clients over stdio: the widgets, how each looks on macOS, Windows 11 and GNOME, and snippets that use them. It is read-only, so it can be enabled in an agent without granting access to modify an app or the repository.

Run it from a MyGo checkout:

```sh
bun run packages/breeze-mcp/src/index.ts
```

An installed package provides the `breeze-ui-mcp` command. Add that command to the MCP client configuration used by your editor or agent. The server does not open a network port.

## Tools

| Tool | Purpose |
|---|---|
| `search_components` | Find components by name, category, or description. |
| `get_component` | Read one component's catalog entry, with its look (`drawn` or `themed`). |
| `list_platforms` | List the platforms Breeze UI draws and the version each look follows. |
| `generate_snippet` | Generate a Go snippet for a component. |

The `breeze-ui://catalog` resource returns the complete JSON catalog. The source of truth is the files of `ui/breeze/catalog/`, one per group of widgets beside `platforms.json`; the Go package, CLI, MCP server, and website all merge those files into the same catalog so names and examples stay aligned. The published package carries the merged catalog, which its build writes to `packages/breeze-mcp/catalog.json` (not committed).
