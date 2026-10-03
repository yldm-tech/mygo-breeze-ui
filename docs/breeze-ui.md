# Breeze UI for MyGo

`github.com/yldm-tech/mygo-breeze-ui/ui/breeze` is a third-party typed design system for MyGo's native UI. It
uses Tailwind's utility-first vocabulary and token model while keeping MyGo's
Go rendering model: there is no DOM, CSS parser, browser, or JavaScript.

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

## Install

Install the package alongside MyGo:

```sh
go get github.com/yldm-tech/mygo-breeze-ui/ui/breeze
```

A native-only app still needs only Go. Start from the gallery with:

```sh
go run ./examples/breeze-ui
```

## Utilities

Utility helpers return typed functions that can be composed with `breeze.Apply` or
passed to `breeze.Row`, `breeze.Column`, and `breeze.Box`:

| Tailwind idea | MyGo helper |
|---|---|
| `flex flex-row` | `breeze.Row(c)` |
| `flex flex-col` | `breeze.Column(c)` |
| `gap-4` | `breeze.Gap(4)` |
| `p-4`, `px-4`, `py-4` | `breeze.P(4)`, `breeze.Px(4)`, `breeze.Py(4)` |
| `w-*`, `h-*`, `grow` | `breeze.W`, `breeze.H`, `breeze.Grow` |
| `items-center` | `breeze.ItemsCenter` |
| `justify-between` | `breeze.JustifyBetween` |
| `rounded-lg` | `breeze.Rounded(breeze.RadiusLG)` |
| `bg-*`, `text-*`, `border` | `breeze.Bg`, `breeze.TextColor`, `breeze.Border` |
| `text-lg`, `font-semibold` | `breeze.TextSize(breeze.TextLG)`, `breeze.FontWeight(600)` |

The default spacing unit is four DIPs, so `breeze.P(4)` means 16 DIPs. Use MyGo's
`Context.SetTheme` for a custom native theme; `breeze.Palette(c)` reads semantic
colors from the current frame.

## Components

The first component set includes `Button`, `Card`, `Input`, `TextArea`,
`Checkbox`, `Switch`, `Select`, `Badge`, `Progress`, `Tabs`, and `Dialog`.
Components delegate keyboard behavior, focus, text editing, and accessibility
to the existing MyGo widgets. Buttons expose variants such as
`ButtonPrimary`, `ButtonSecondary`, `ButtonDanger`, and `ButtonGhost`.

For a component with a custom visual design, start with MyGo's headless bases
such as `ui.ButtonBase`, `ui.SelectBase`, `ui.DialogBase`, and
`ui.TextInputBase`, then apply `breeze` utilities and `ui.Painter` drawing.

## State and responsive layouts

CSS selector variants become explicit Go state. Use `element.Hovered()`,
`Pressed()`, `Focused()`, `Clicked()`, and `Disabled()` while building the
frame. Use `c.Theme().Dark` for appearance-specific styling and `c.Size()` for
window-size-dependent layout decisions.

CSS-specific features such as arbitrary selectors, pseudo-elements, CSS Grid,
container queries, and media-query breakpoints are outside this package. Use
MyGo layout primitives or custom drawing for those cases.
