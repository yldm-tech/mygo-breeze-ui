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
| `Toggle` | push-on-push-off button, darker while on | ToggleButton, the accent while on | toggle button, darker while checked |
| `Checkbox` | 16-point rounded box | 20-pixel box with a stroke | 20-pixel box with a trough border |
| `Radio` | 16-point circle, 5-point dot | 20-pixel circle, dot that grows under the pointer | 20-pixel circle, 8-pixel bullet |
| `Switch` | 54 by 24 track, capsule knob | 40 by 20 track, knob that grows under the pointer | 46 by 26 track, round slider |
| `TextInput` | bezeled field, 26 high, focus ring around it | TextBox, accent line along the bottom while focused | filled entry, 34 high, inner focus ring |
| `TextArea` | text view in a bezel border | TextBox | text view in a frame |
| `SearchField` | capsule search field | AutoSuggestBox | search entry |
| `Combobox` | bezeled field around MyGo's list | TextBox around MyGo's list | entry around MyGo's list |
| `Select` | pop-up button and menu with a check mark | ComboBox, accent pill on the chosen item | drop-down and popover with a check mark |
| `Slider` | 6-point track, 20 by 16 capsule knob | 4-pixel track, thumb with an accent dot | 4-pixel trough, 20-pixel slider |
| `Progress` | 8-point bar | 3-pixel bar over a 1-pixel line | 8-pixel bar |
| `Spinner` | 8 spokes, 32 points | ProgressRing, 32 pixels | AdwSpinner, 16 pixels |
| `Segmented` | NSSegmentedControl | SelectorBar | AdwToggleGroup |
| `Tabs` | tab view, drawn as a segmented control | SelectorBar | notebook tabs with the accent line |
| `Modal` | sheet with a window's shadow, the window behind not dimmed | ContentDialog over the smoke layer | floating AdwDialog over the dimmed window |
| `AlertDialog` | NSAlert: 260 points, centered text, capsule buttons, the default last | ContentDialog: title and message above a command area, the default first | AdwAlertDialog: title-2 heading, 44-pixel buttons, the default last |
| `Link` | `linkColor`, underlined under the pointer | HyperlinkButton's accent text colors, underlined | `accent_color`, underlined |
| `Divider` | `separatorColor` | `DividerStrokeColorDefault` | `$border_color` |
| `Badge` | MyGo's (macOS has no badge control) | InfoBadge: 16 high, accent | indicator badge: 13 high, 40% gray, white bold text |
| `MenuButton`, `PopoverMenuButton` | push button | button | button without a border |
| `Popover` | panel like a menu | FlyoutPresenter | popover contents |
| `ToggleGroup` | segmented control selecting any segments | ToggleButtons side by side (WinUI joins none) | `box.linked`, with the seam where two buttons meet |
| `Field` | form labels in a column before the controls (in a `Form`) | header 8 above the control | label above the control |
| `Stepper` | NSStepper's box, 20 points wide | a button's face and elevation border | a button's face |

These are drawn by Breeze UI. Their sizes and colors come from the platforms themselves:

- **macOS:** AppKit's semantic colors (`NSColor`) resolved on macOS 27, and AppKit's own controls rendered at 2x by [`tools/appkit-reference/reference.swift`](../tools/appkit-reference/reference.swift), which writes them as PNGs and prints the gray values across each control. A push button's face is `secondarySystemFill`; an unchecked box, a radio button and an off switch are `systemFill`. Run `BREEZE_SCREENSHOTS=<dir> go test ./ui/breeze -run TestScreenshots` for Breeze UI's controls in the same places.
- **Windows 11:** the theme resources of WinUI 3 (`microsoft/microsoft-ui-xaml`): `Common_themeresources_any.xaml` and the theme resources of Button, ToggleButton, CheckBox, RadioButton, ToggleSwitch, TextBox, ComboBox, Slider, ProgressBar, ProgressRing, SelectorBar, ContentDialog and ToolTip; the dialog's shadow from `DropShadowRecipe.h`.
- **GNOME:** the libadwaita 1.7 stylesheet (GNOME 48), and the C code of libadwaita and GTK where the stylesheet leaves a value to them: the accent, the switch's and the scale's sizes, the drop-down's check mark, the spinner's arc and the alert's width. The look is compared with libadwaita 1.7.6 itself: [`tools/adwaita-reference/reference.py`](../tools/adwaita-reference/reference.py) renders the same controls with GTK 4.18 in the same places as the screenshot test, and button sizes and colors, check boxes, switches, radio buttons, entries, toggles, scales, progress bars, text views, alerts and links match it to the pixel but for the font.

The pointer, the keyboard, the focus and accessibility are MyGo's. The buttons, choices, slider, segments and tabs are built on MyGo's bases without a look (`ui.ButtonBase`, `ui.ToggleBase`, `ui.CheckboxBase`, `ui.RadioBase`, `ui.SwitchBase`, `ui.SelectBase`, `ui.SliderBase`, `ui.SegmentedBase`, `ui.TabsBase`), so screen readers, Tab, Space and the arrows work as they do for `ui`'s widgets. The text fields are `ui`'s own, with the box around their text drawn by the platform, so editing, input methods, the search field's clear button and Escape are MyGo's. `Disabled(true)` grays a widget out: MyGo draws disabled elements at half opacity, which is how AppKit and libadwaita draw them. WinUI has disabled colors of its own instead (a lighter fill and `TextFillColorDisabled`), which Breeze UI does not draw yet; half opacity comes close to them.

The alerts keep `ui.AlertDialog`'s contract: the last button is the default, which has the focus and which Enter clicks, Escape returns `ui.AlertDismissed`, and the index returned names the same button on every platform, though Windows places the default first.

The other widgets under `ui`'s names (`NumberInput`, `RadioGroup`, `Table`, `Tree`, `TreeItem`) are package `ui`'s in the platform's theme. MyGo offers no base without a look for them, even on its main branch, so drawing the platforms' rows and fields would mean writing their keyboard handling, scrolling and accessibility again; their colors, corner radius and font size come from `Use`. `RadioGroup` is a container around the platform's radio buttons. `MenuButton`'s and `PopoverMenuButton`'s buttons are the platform's; the menu of `MenuButton` is the system's own, and `PopoverMenuButton`'s is MyGo's in the platform's theme until MyGo releases its `MenuBase`. `Stepper` keeps MyGo's arrows and behavior inside the platform's box. `breeze-ui list` shows which widgets are drawn and which are themed.

Tooltips (`Tooltip` on any element) take their colors from the theme: on Windows the ToolTip's light or dark acrylic fallback with the text color, on GNOME libadwaita's tooltip, the shade at 80% with white text. On macOS they keep MyGo's own look, which is not measured from AppKit's.

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

- **Inactive windows.** AppKit draws an inactive window's checked boxes, chosen radio buttons, on switches and default button in gray. MyGo v0.7.4 does not tell a view whether its window is active, so Breeze UI always draws the active look.
- **What AppKit would not render offscreen.** The reference tool renders AppKit's controls in an inactive window, so the chosen segment of a segmented control is measured in its inactive gray, which Breeze UI draws; an active window may draw it differently. Menus are not rendered offscreen at all: the pop-up button's menu (its corners, padding and shadow) is drawn to resemble macOS's, not measured.
- **Sizes.** The platforms' controls differ in size: a macOS push button is 24 points high, Windows's 32 and GNOME's 34, and a GNOME scale is 44 pixels high with its padding. The same layout is denser on a Mac; size elements around the widgets from the content rather than with fixed heights.
- **Windows fidelity.** The macOS look is compared with AppKit's rendering and the GNOME look with libadwaita's, pixel by pixel. The Windows look follows the values of WinUI's sources but has not been compared with screenshots of Windows yet. Some values are inferred rather than read: the padding of an accent button and the line height of 14-pixel text on Windows, the shape of the ProgressRing's and the indeterminate ProgressBar's motion (Lottie animations and code in WinUI), the ComboBox's item height, and the spacing of toggles in a `ToggleGroup`, as WinUI has no joined toggles. Accent shades for a custom accent (link colors and fills) are mixes that give Windows's default accent's shades exactly.
- **What only MyGo draws.** The arrows of `Stepper`, the chevron of `MenuButton` and the menus of `PopoverMenuButton` are MyGo's: GNOME's menu button shows a solid triangle, and GNOME's spin button puts − and + beside the text rather than a stepper.
- **Dialogs.** An alert shows no icon: NSAlert shows the app's icon above its title, which a view cannot reach. GNOME's alert stacks its buttons when they do not fit side by side; Breeze UI keeps them in a row. On macOS the corner radius and the alert's colors are measured from NSAlert's window; the hairline around a dialog and its shadow are drawn to resemble AppKit's, not measured.
- **Errors on text fields.** The platforms' text fields draw their own borders and focus, so a field's `Error` does not turn its border red as `ui`'s fields do.
- **Icons and fonts.** Breeze UI ships no SF Symbols, Segoe Fluent Icons or Adwaita icons: check marks, chevrons and arrows are drawn as paths, and the search field's magnifying glass is MyGo's. Text uses the system's font (`Font: ""`).
