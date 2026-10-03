# Breeze UI for MyGo

[![CI](https://github.com/yldm-tech/mygo-breeze-ui/actions/workflows/ci.yml/badge.svg)](https://github.com/yldm-tech/mygo-breeze-ui/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/yldm-tech/mygo-breeze-ui.svg)](https://pkg.go.dev/github.com/yldm-tech/mygo-breeze-ui)
[![License](https://img.shields.io/github/license/yldm-tech/mygo-breeze-ui)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/yldm-tech/mygo-breeze-ui)](https://github.com/yldm-tech/mygo-breeze-ui/stargazers)

[Website](https://yldm-tech.github.io/mygo-breeze-ui/) · [English](README.md) · [简体中文](README.zh-CN.md)

Breeze UI is a Tailwind-inspired typed design system for [MyGo](https://github.com/egoist/mygo)'s native UI. It provides Go utility styles, design tokens, accessible components, a CLI, and an MCP catalog without HTML, CSS, or JavaScript at runtime.

## What is included

- Typed utility styles for layout, spacing, sizing, color, typography, borders, and shadows.
- Native components such as buttons, cards, fields, badges, progress bars, tabs, and dialogs.
- `breeze-ui`, a CLI for discovering the catalog and generating Go snippets.
- `breeze-ui-mcp`, a read-only MCP server for catalog search and snippet generation.

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
bun run build
```

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=yldm-tech/mygo-breeze-ui&type=Date)](https://star-history.com/#yldm-tech/mygo-breeze-ui&Date)

## License

[MIT](LICENSE)
