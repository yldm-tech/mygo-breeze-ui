# The `breeze-ui` CLI

`breeze-ui` lists the widgets of `github.com/yldm-tech/mygo-breeze-ui/ui/breeze` and the platforms they follow, and writes snippets that use them. It reads the same catalog as the website and the MCP server.

Install the Go command:

```sh
go install github.com/yldm-tech/mygo-breeze-ui/cmd/breeze-ui@latest
```

Or use the npm wrapper in a Bun or npm workspace:

```sh
npx breeze-ui init
```

## Commands

```sh
breeze-ui list [-kind component|platform]
breeze-ui add Switch [-o ui/switch_snippet.go]
breeze-ui init [directory] [--force]
breeze-ui example
```

- `list` prints the catalog: each component with its look (`drawn` or `themed`), and each platform with the version its look follows.
- `add` writes a Go snippet for a component; options can come before or after the component's name.
- `init` creates a `breeze-ui.json` marker with the catalog's version, and does not replace an existing one without `--force`.
- `example` prints a small view that uses the widgets.

The CLI only discovers and scaffolds: an app needs nothing beyond the `ui/breeze` Go package.
