# Breeze UI for MyGo

[![CI](https://github.com/yldm-tech/mygo-breeze-ui/actions/workflows/ci.yml/badge.svg)](https://github.com/yldm-tech/mygo-breeze-ui/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/yldm-tech/mygo-breeze-ui.svg)](https://pkg.go.dev/github.com/yldm-tech/mygo-breeze-ui)
[![License](https://img.shields.io/github/license/yldm-tech/mygo-breeze-ui)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/yldm-tech/mygo-breeze-ui)](https://github.com/yldm-tech/mygo-breeze-ui/stargazers)

[官网](https://yldm-tech.github.io/mygo-breeze-ui/) · [English](README.md) · [简体中文](README.zh-CN.md)

官网提供[目录 API 文档](https://yldm-tech.github.io/mygo-breeze-ui/api/)，展示版本化的只读目录接口。

Breeze UI 是面向 [MyGo](https://github.com/egoist/mygo) 原生 UI 的 Tailwind 风格类型化设计系统。它提供 Go 工具样式、设计令牌、原生组件、CLI 和 MCP 目录服务，运行时不需要 HTML、CSS 或 JavaScript。

## 包含内容

- 覆盖布局、间距、尺寸、颜色、排版、边框和阴影的类型化工具样式。
- Button、Card、Input、Badge、Progress、Tabs、Dialog 等原生组件。
- `breeze-ui` CLI，用于浏览目录和生成 Go 代码片段。
- `breeze-ui-mcp` 只读 MCP 服务，用于搜索目录和生成代码片段。

## Go 包

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
        breeze.Text(c, "账户", breeze.TextSize(breeze.TextLG), breeze.FontWeight(600))
        breeze.Input(c, &name, "你的名字")
        if breeze.Button(c, "保存", breeze.ButtonPrimary).Clicked() {
            saved = true
        }
    })
}
```

原生渲染器由 MyGo 依赖提供，Breeze UI 在其上增加 Tailwind 风格的设计词汇和组件。

## CLI

```sh
npx breeze-ui init
npx breeze-ui list -kind component
npx breeze-ui add Button -o ui/button_snippet.go
```

也可以直接安装 Go 命令：

```sh
go install github.com/yldm-tech/mygo-breeze-ui/cmd/breeze-ui@latest
```

## MCP 服务

```sh
bun install
bun run build
```

`breeze-ui-mcp` 提供只读工具，用于搜索组件、工具样式、设计令牌并生成代码片段。

## 开发

```sh
go test ./...
bun install
bun run test
bun run typecheck
bun run build
```

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=yldm-tech/mygo-breeze-ui&type=Date)](https://star-history.com/#yldm-tech/mygo-breeze-ui&Date)

## 许可证

[MIT](LICENSE)
