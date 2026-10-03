# The `breeze-ui` CLI

`breeze-ui` discovers and scaffolds the catalog in
`github.com/yldm-tech/mygo-breeze-ui/ui/breeze`. It is useful when starting a native UI screen,
and it reads the same catalog used by the website and MCP server.

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
breeze-ui list [-kind utility|component|token]
breeze-ui add Button [-o ui/button_snippet.go]
breeze-ui tokens [-o theme/tokens.go]
breeze-ui init [directory] [--force]
breeze-ui example
```

`list` prints catalog entries. `add` writes a Go snippet for a component;
options can appear before or after the component name. `tokens` writes a small
generated Go file containing the numeric design tokens. `init` creates a
`breeze-ui.json` marker with the catalog version and refuses to overwrite an
existing file unless `--force` is provided. `example` prints a minimal view.

The CLI is a scaffold and discovery tool. It does not parse CSS or add a
runtime dependency beyond the normal `ui/breeze` Go package.
