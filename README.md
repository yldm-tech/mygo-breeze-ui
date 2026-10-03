# Breeze UI for MyGo

Breeze UI is a Tailwind-inspired typed design system for [MyGo](https://github.com/egoist/mygo)'s native UI. It provides Go utility styles, design tokens, accessible components, a CLI, and an MCP catalog without HTML, CSS, or JavaScript at runtime.

## Go package

```sh
go get github.com/yldm-tech/mygo-breeze-ui
```

```go
import (
    "github.com/egoist/mygo/ui"
    breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func view(c *ui.Context) {
    breeze.Card(c, func() {
        breeze.Text(c, "Account", breeze.TextSize(breeze.TextLG), breeze.FontWeight(600))
        breeze.Input(c, &name, "Your name")
        if breeze.Button(c, "Save", breeze.ButtonPrimary).Clicked() {
            saved = true
        }
    })
}
```

The native renderer comes from the MyGo dependency; Breeze UI adds the Tailwind-style vocabulary and components on top of it.

## CLI

```sh
npx breeze-ui init
npx breeze-ui list -kind component
npx breeze-ui add Button -o ui/button_snippet.go
```

Or install the Go command:

```sh
go install github.com/yldm-tech/mygo-breeze-ui/cmd/breeze-ui@latest
```

## MCP server

```sh
bun install
bun run build
```

The `breeze-ui-mcp` package exposes read-only tools for searching components, utilities, tokens, and generated snippets.

## Development

```sh
go test ./...
bun install
bun run test
bun run typecheck
```
