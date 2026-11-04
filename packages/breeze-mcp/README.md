# Breeze UI MCP

This package exposes the Breeze UI catalog to MCP clients. Breeze UI draws MyGo's native UI widgets as macOS, Windows 11 and GNOME draw their controls. The server is read-only: tools search components, list the platforms and generate Go snippets. The `breeze-ui://catalog` resource contains the complete catalog.

Run it from a checkout with:

```sh
bun run packages/breeze-mcp/src/index.ts
```

For an MCP client that supports command configuration, use `breeze-ui-mcp` as the command after installing the package. The server communicates over stdio and never starts a network listener.

Available tools:

- `search_components` — search components by name, category, or description
- `get_component` — inspect one component, with its look (`drawn` or `themed`)
- `list_platforms` — list the platforms and the version each look follows
- `generate_snippet` — generate a Go component snippet
