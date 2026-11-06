# Breeze UI for MyGo

[![CI](https://github.com/yldm-tech/mygo-breeze-ui/actions/workflows/ci.yml/badge.svg)](https://github.com/yldm-tech/mygo-breeze-ui/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/yldm-tech/mygo-breeze-ui.svg)](https://pkg.go.dev/github.com/yldm-tech/mygo-breeze-ui)
[![License](https://img.shields.io/github/license/yldm-tech/mygo-breeze-ui)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/yldm-tech/mygo-breeze-ui)](https://github.com/yldm-tech/mygo-breeze-ui/stargazers)

[Website](https://yldm-tech.github.io/mygo-breeze-ui/) · [English](README.md) · [简体中文](README.zh-CN.md)

Breeze UI draws [MyGo](https://github.com/egoist/mygo)'s native UI widgets as the desktop the app runs on draws its own: macOS's controls on a Mac, Windows 11's on Windows and GNOME's libadwaita controls on Linux. The widgets have the names and signatures of MyGo's package `ui`, so a view swaps `ui.Button` for `breeze.Button` and keeps everything else.

![The same view drawn as macOS, Windows 11 and GNOME](docs/assets/gallery-light.png)

## What is included

- `ui/breeze`: the common controls of macOS 27, Windows 11 (WinUI 3) and GNOME 51 (libadwaita 1.10), each drawn as its platform draws it: buttons and toggles, text fields and pickers, sliders and progress, tabs, navigation and header bars, lists, tables and trees, settings rows, dialogs, menus and messages. The [widget table](docs/breeze-ui.md) lists every one with its counterpart on each platform, and the [storybook](https://yldm-tech.github.io/mygo-breeze-ui/storybook/) shows them.
- Themes of the three platforms, light and dark, that follow the accent color, contrast and text size the user chose.
- `breeze-ui`, a CLI that lists the widgets and platforms and writes snippets.
- `breeze-ui-mcp`, a read-only MCP server over the same catalog.

## Go package

```sh
go get github.com/yldm-tech/mygo-breeze-ui
```

```go
import (
    "github.com/egoist/mygo/ui"
    breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func (app *settings) view(c *ui.Context) {
    breeze.Use(c) // the platform's colors and metrics
    ui.Column(c).Padding(16).Gap(12).Children(func() {
        breeze.Checkbox(c, &app.remember, "Remember me")
        breeze.Switch(c, &app.sync).Label("Sync")
        ui.Row(c).Gap(8).Children(func() {
            breeze.Button(c, "Cancel")
            if breeze.PrimaryButton(c, "Save").Clicked() {
                app.save()
            }
        })
    })
}
```

The platform is the system's unless `breeze.SetPlatform` picks another; `breeze.Preview(c, p, body)` draws part of a window as another platform, as the gallery does:

```sh
go run ./examples/breeze-ui
```

The pointer, the keyboard, the focus and accessibility come from MyGo's widget bases (`ui.ButtonBase`, `ui.CheckboxBase`, `ui.SwitchBase`). [docs/breeze-ui.md](docs/breeze-ui.md) lists what each widget looks like on each platform, where the values come from, and the current limits.

## Storybook

Every widget, on every platform, light and dark, with controls that change its arguments and the Go that builds it: **[the storybook](https://yldm-tech.github.io/mygo-breeze-ui/storybook/)** runs in the browser, compiled to WebAssembly, and the same program opens as a window:

```sh
go run ./examples/storybook
```

## CLI

```sh
npx breeze-ui list
npx breeze-ui add Switch -o ui/switch_snippet.go
```

Or install the Go command:

```sh
go install github.com/yldm-tech/mygo-breeze-ui/cmd/breeze-ui@latest
```

or download it for Linux, macOS or Windows from the [releases](https://github.com/yldm-tech/mygo-breeze-ui/releases). `breeze-ui version` prints the version it was built as.

## MCP server

```sh
bun install
bun run build
```

The `breeze-ui-mcp` package exposes read-only tools to search components, list the platforms and generate snippets.

## Development

```sh
go test ./...
go run ./examples/breeze-ui
BREEZE_SCREENSHOTS=/tmp/shots go test ./ui/breeze -run TestScreenshots   # write each platform's controls as PNG
bun install
bun run test
bun run typecheck
bun run build
```

Pull request titles and commit messages are [Conventional Commits](https://www.conventionalcommits.org/) (`feat(breeze): ...`, `fix: ...`); the Labels workflow checks titles and labels each pull request by type, size, area and platform, and the Conflicts workflow labels those that conflict with their base branch. Once CI passes on `main`, the Release workflow releases it when the commits since the last tag call for a new version: it writes the version into the npm packages and the MCP server, commits that to `main` as `chore(release)`, tags the commit and publishes a GitHub release with the `breeze-ui` command built for Linux, macOS and Windows. For the version, `feat` bumps the minor one, `fix`, `perf` and `revert` the patch one, a breaking change (`feat!:` or a `BREAKING CHANGE:` footer) the major one, or the minor one while the version is 0.x. Run the Release workflow by hand to force a bump.

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=yldm-tech/mygo-breeze-ui&type=Date)](https://star-history.com/#yldm-tech/mygo-breeze-ui&Date)

## License

[MIT](LICENSE)
