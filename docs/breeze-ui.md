# Breeze UI for MyGo

`github.com/yldm-tech/mygo-breeze-ui/ui/breeze` draws [MyGo](https://github.com/egoist/mygo)'s native UI widgets as the desktop the app runs on draws its own controls: macOS's on a Mac, Windows 11's on Windows and GNOME's libadwaita controls on Linux.

Its widgets have the names and signatures of package `ui`'s. A view swaps `ui.Button` for `breeze.Button` and keeps everything else: the widgets return `ui.Element`, answer `Clicked`, and take more styles after the call as `ui`'s do.

```go
import (
	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func (app *settings) view(c *ui.Context) {
	breeze.Use(c) // the platform's colors and metrics, for package ui's widgets too
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

## Install

```sh
go get github.com/yldm-tech/mygo-breeze-ui
```

It needs MyGo v0.7.4 or later. The gallery shows the three platforms side by side:

```sh
go run ./examples/breeze-ui
```

## Widgets

| Widget | macOS | Windows 11 | GNOME |
|---|---|---|---|
| `Button` | push button, 24 points high, no border | button with its elevation border, 32 high | text button, 34 high, bold |
| `PrimaryButton` | default button in the accent color | accent button | suggested action |
| `Checkbox` | 16-point rounded box | 20-pixel box with a stroke | 20-pixel box with a trough border |
| `Switch` | 54 by 24 track, capsule knob | 40 by 20 track, knob that grows under the pointer | 46 by 26 track, round slider |

These four are drawn by Breeze UI. Their sizes and colors come from the platforms themselves:

- **macOS:** AppKit's semantic colors (`NSColor`) resolved on macOS 27, and `NSButton` and `NSSwitch` that AppKit rendered at 2x. A push button's face is `secondarySystemFill`; an unchecked box and an off switch are `systemFill`.
- **Windows 11:** the theme resources of WinUI 3 (`microsoft/microsoft-ui-xaml`): `Common_themeresources_any.xaml` and the Button, CheckBox and ToggleSwitch theme resources.
- **GNOME:** the libadwaita 1.7 stylesheet (GNOME 48), and the C code of libadwaita and GTK where the stylesheet leaves a value to them (the accent, the switch's size).

The pointer, the keyboard, the focus and accessibility are MyGo's: each widget is built on MyGo's base without a look (`ui.ButtonBase`, `ui.CheckboxBase`, `ui.SwitchBase`), so screen readers, Tab and Space work as they do for `ui`'s widgets. `Disabled(true)` grays a widget out: MyGo draws disabled elements at half opacity, which is how AppKit and libadwaita draw them. WinUI has disabled colors of its own instead (a lighter fill and `TextFillColorDisabled`), which Breeze UI does not draw yet; half opacity comes close to them.

`TextInput`, `TextArea`, `Select`, `Slider`, `Segmented`, `Tabs` and `Progress` are package `ui`'s widgets under the same names. They take the platform's colors, corner radius and font size from `Use`, but keep MyGo's shapes: a text field on macOS is filled rather than outlined, for example. `breeze-ui list` shows which widgets are drawn and which are themed.

## Platforms

`breeze.Current()` is the platform the app runs on: `Mac` on macOS, `Windows` on Windows and `Linux` elsewhere. Every widget looks it up when it is built.

- `breeze.SetPlatform(breeze.Windows)` makes every window draw another platform's controls, as for screenshots; `breeze.SetPlatform(breeze.Auto)` follows the system again.
- `breeze.Preview(c, p, body)` draws only what `body` builds as platform `p`, in this window. With `ui.Themed(c, breeze.Theme(c, p))` around it, a gallery shows the platforms side by side:

```go
for _, p := range breeze.Platforms {
	breeze.Preview(c, p, func() {
		ui.Themed(c, breeze.Theme(c, p)).Grow(1).Padding(16).Gap(12).Children(func() {
			breeze.Button(c, "Cancel")
			breeze.Switch(c, &app.sync).Label("Sync")
		})
	})
}
```

## Themes

`breeze.Theme(c, p)` returns platform `p`'s theme for the window's appearance: the window background, text, separators and accent of macOS, Windows 11 or GNOME, light or dark. `breeze.Use(c)` sets it for the window with `c.SetTheme`, so package `ui`'s widgets and elements of your own that read `c.Theme()` follow it too.

Like MyGo's default theme it follows the desktop's preferences:

- **Accent color.** The accent the user chose replaces the platform's default (`#007aff`, `#0078d4`, `#3584e4`). On Windows the controls are filled with the shade WinUI uses, darker in the light theme and lighter in the dark one.
- **Higher contrast.** Separators and secondary text are darker in the light theme and lighter in the dark one, and the focus ring is opaque.
- **Text size.** The font size grows with Windows's and GNOME's text scale.

## Limits

- **Inactive windows.** AppKit draws an inactive window's checked boxes, on switches and default button in gray. MyGo v0.7.4 does not tell a view whether its window is active, so Breeze UI always draws the active look.
- **Sizes.** The platforms' controls differ in size: a macOS push button is 24 points high, Windows's 32 and GNOME's 34. The same layout is denser on a Mac; size elements around the widgets from the content rather than with fixed heights.
- **Windows and GNOME fidelity.** The macOS look is compared with AppKit's rendering pixel by pixel. The Windows and GNOME looks follow the values of WinUI and libadwaita's sources but have not been compared with screenshots of those systems yet; two values (the padding of an accent button and the line height of 14-pixel text on Windows) are taken from the default button.
- **Icons and fonts.** Breeze UI ships no SF Symbols, Segoe Fluent Icons or Adwaita icons: check marks are drawn as paths. Text uses the system's font (`Font: ""`).
