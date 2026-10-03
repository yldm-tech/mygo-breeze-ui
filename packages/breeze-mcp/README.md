# Breeze UI MCP

This package exposes the Breeze UI catalog to MCP clients. Breeze UI is a third-party Tailwind-style native UI system for MyGo. The server is read-only: tools search components and utilities, inspect tokens, and generate Go snippets. The `breeze-ui://catalog` resource contains the complete catalog.

Run it from a checkout with:

```sh
bun run packages/breeze-mcp/src/index.ts
```

For an MCP client that supports command configuration, use `breeze-ui-mcp` as the command after installing the package. The server communicates over stdio and never starts a network listener.

Available tools:

- `search_components` — search components by name, category, or description
- `get_component` — inspect one component
- `search_utilities` — search layout and visual utilities
- `get_token` — inspect one design token
- `generate_snippet` — generate a Go component snippet
