# Breeze UI MCP

This package exposes the Breeze UI catalog to MCP clients. Breeze UI draws MyGo's native UI widgets as macOS, Windows 11 and GNOME draw their controls. The server is read-only: tools search components, list the platforms and generate Go snippets. The `breeze-ui://catalog` resource contains the complete catalog.

Run it from a checkout with:

```sh
bun run packages/breeze-mcp/src/index.ts
```

For Claude Code, add the package as a project-scoped stdio server:

```sh
claude mcp add --transport stdio --scope project breeze-ui -- npx -y breeze-ui-mcp
claude mcp list
```

For Codex CLI:

```sh
codex mcp add breeze-ui -- npx -y breeze-ui-mcp
codex mcp list
```

For Cursor Agent:

```sh
agent mcp add breeze-ui -- npx -y breeze-ui-mcp
agent mcp list
```

The server communicates over stdio and never starts a network listener.

Available tools:

- `search_components` — search components by name, category, or description
- `get_component` — inspect one component, with its look (`drawn` or `themed`)
- `list_platforms` — list the platforms and the version each look follows
- `generate_snippet` — generate a Go component snippet
