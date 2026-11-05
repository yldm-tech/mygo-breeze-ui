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

## Storybook

[The storybook](https://yldm-tech.github.io/mygo-breeze-ui/storybook/) shows every widget as Storybook shows components: stories in a sidebar, each drawn as macOS, Windows 11 or GNOME draws it, or as the three side by side, light or dark, with controls for its arguments and the Go that builds it. It is [`examples/storybook`](../examples/storybook), which runs as a window with `go run ./examples/storybook`; the site's deployment compiles it to WebAssembly (`GOOS=js GOARCH=wasm`), which MyGo draws in the page with WebGL 2. To try the web build locally, build the site, then the storybook into it, and serve `website/dist` under `/mygo-breeze-ui/`:

```sh
bun run site:build
GOOS=js GOARCH=wasm go build -o website/dist/storybook/app.wasm ./examples/storybook
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" website/dist/storybook/
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
| `Collapsible` | disclosure chevron in a 13-point box, drawn where AppKit draws it, label 1.5 after the box, rows 24 high | Expander: 48-high header card, 32-pixel chevron button, content card padded 16 | GtkExpander: 16-pixel pan-end arrow at 70%, label right after it |
| `Accordion` | grouped form section: corners of 12, 3.1% fill | Expander cards 4 apart (MyGo's headers) | boxed list: card, corners of 12, card shadow |
| `AccordionItem` | row below a separator inset 10 points | card with CardStrokeColorDefault, corners of 4 | row below a `card_shade_color` separator |
| `Fieldset` | legend over a grouped section (NSBox's fill), inset 10; in a Form the section reaches 10 past the fields | legend over the fields, 10 above them (settings section header; no group box) | legend over a card, 12 above it (AdwPreferencesGroup); in a Form the card reaches 12 past the fields |
| `Form` | labels in a column before the controls (ui.Form) | labels above the controls, 8 above, fields 12 apart | labels in a column before the controls (ui.Form) |
| `Sidebar` | source list, 10-point inset, its own gray background (#f0f0f0, #282828) | NavigationView's left pane over the window | `.navigation-sidebar` on `sidebar_bg_color`, a line at its end |
| `SidebarSection` | group row: 11-point semibold title in `tertiaryLabelColor`, text placed as AppKit's (13 points from the top, 13 below the row above) | NavigationViewItemHeader: a 40-pixel row, 14 pixels, semibold, secondary text, 16 from the pane's edge | dimmed heading (libadwaita has none) |
| `SidebarItem` | 32-point row, corners of 8, accent selection while focused, gray otherwise, label semibold when chosen | NavigationViewItem: 36 high, subtle fill (tertiary pressed), 3 by 16 accent pill | 36-pixel row, corners of 9, 10% when chosen, 7% hovered, 16% pressed |
| `Breadcrumbs` | NSPathControl: secondary gray items and 4 by 6 chevrons | BreadcrumbBar: 14-pixel items, 12-pixel chevrons | Files' path bar: bold buttons, slashes, in a 10% trough |
| `BackButton`, `ForwardButton` | toolbar item: 36-point round platter with a hairline, SF Symbols' chevron | NavigationView back button: 40 by 36, subtle, WinUI's back arrow | flat image button, 34 pixels, `go-previous-symbolic` |
| `Toolbar` | 52 points high, items 8 from its edges and 8 apart | CommandBar, 48 high | `.toolbar`: padding and spacing 6 |
| `Tooltip` | 11-point text (`NSFont.toolTipsFont`) in a light or dark panel below the element; panel inferred | ToolTip: in-app acrylic fallback (#f9f9f9, #2c2c2c), 1-pixel flyout stroke, 4-pixel corners, padding 9,6,9,8, 12-pixel text, above the element | GTK tooltip: the shade at 80%, white text, 1-pixel white 10% border, 9-pixel corners, padding 6 10, below the element |
| `Toasts` | panel like a notification banner (macOS has no toast; inferred) | InfoBar with a severity icon, colored by the toast's Type, 38-pixel close button | AdwToast: #505053 pill, heading title, action button, circular flat close button |
| `Meter` | NSLevelIndicator: 16-point bar, systemGreen, systemYellow or systemRed with an edge 13% darker | ProgressBar look: 3-pixel bar over a 1-pixel track, caution and critical colors | GtkLevelBar: 8-pixel round trough, accent, warning, error or success |
| `Avatar` | gray monogram, white initials, a silhouette without a name (inferred) | PersonPicture: gray ellipse, semibold initials at 42% by InitialsGenerator's rules, Contact glyph without them | AdwAvatar: one of 14 gradients by the name's hash, avatar-default-symbolic without a name, white initials below 25 pixels |
| `Rating` | NSLevelIndicator rating: 12-point stars, secondaryLabelColor and systemFill | RatingControl: 16-pixel accent stars 8 apart, growing under the pointer; faint fill hovering an unset rating | GNOME Software's stars: 32 pixels, yellow (GNOME has no rating control) |
| `StepSlider` | slider with 2-point dots 3 below the track, white knob in both modes, 16 high | Slider with 1-pixel ticks 4 below (TickPlacement BottomRight), focus ring around the whole slider | scale with 6-pixel marks after it, slider pointing down |
| `DateInput` | NSDatePicker textField bezel, 23 high, corners of 7 | CalendarDatePicker: button face, 32 high, glyph column at the end | entry, 34 high (GTK has no date field) |
| `TimeInput` | NSDatePicker textField bezel around MyGo's segments | TimePicker's button face, 242 wide, 32 high, segments spread across it | entry with the focus ring while a segment has the focus |
| `ColorWell` | NSColorWell: 48 by 24 capsule of systemFill around a capsule swatch; a translucent color split over black and white | DropDownButton with a 36 by 20 swatch and a chevron, 32 high; ColorPicker's checkerboard | GtkColorDialogButton: button padded 5 around a swatch, 34 high; GTK's 10-pixel checkerboard |
| `List` | inset table view: 24-point rows, 10 in from the edges, the text 8 in from the selection; the chosen row rounded 8 in `selectedContentBackgroundColor` with white text, `unemphasizedSelectedContentBackgroundColor` without the focus; rows chosen together make one block rounded at its ends; no hover | ListView: 40-pixel items, `SubtleFillColorSecondary` backplate inset 4,2 with corners of 4 (`SubtleFillColorTertiary` chosen and under the pointer), 3 by 16 `AccentFillColorDefault` pill on the chosen one; pressed states not drawn | GtkListView `.rich-list`: 48-pixel rows on `view_bg_color`, the chosen one the accent at 25% (32% under the pointer), the hover `currentColor` at 4%; pressed states not drawn |
| `Autocomplete` | bezeled field around MyGo's list | TextBox around MyGo's list | entry around MyGo's list |
| `TokenField` | bezeled field around MyGo's gray chips (not NSTokenField's blue tokens) | TextBox around MyGo's chips (not TokenizingTextBox's) | entry around MyGo's chips (GNOME has no token field) |

These are drawn by Breeze UI. Their sizes and colors come from the platforms themselves:

- **macOS:** AppKit's semantic colors (`NSColor`) resolved on macOS 27, and AppKit's own controls rendered at 2x by [`tools/appkit-reference/reference.swift`](../tools/appkit-reference/reference.swift), which writes them as PNGs and prints the gray values across each control. A push button's face is `secondarySystemFill`; an unchecked box, a radio button and an off switch are `systemFill`. Run `BREEZE_SCREENSHOTS=<dir> go test ./ui/breeze -run TestScreenshots` for Breeze UI's controls in the same places.
- **Windows 11:** the theme resources of WinUI 3 (`microsoft/microsoft-ui-xaml`): `Common_themeresources_any.xaml` and the theme resources of Button, ToggleButton, CheckBox, RadioButton, ToggleSwitch, TextBox, ComboBox, Slider, ProgressBar, ProgressRing, SelectorBar, ContentDialog and ToolTip; the dialog's shadow from `DropShadowRecipe.h`.
- **GNOME:** the libadwaita 1.7 stylesheet (GNOME 48), and the C code of libadwaita and GTK where the stylesheet leaves a value to them: the accent, the switch's and the scale's sizes, the drop-down's check mark, the spinner's arc and the alert's width. The look is compared with libadwaita 1.7.6 itself: [`tools/adwaita-reference/reference.py`](../tools/adwaita-reference/reference.py) renders the same controls with GTK 4.18 in the same places as the screenshot test, and button sizes and colors, check boxes, switches, radio buttons, entries, toggles, scales, progress bars, text views, alerts and links match it to the pixel but for the font.

The widgets of the `widgets_*.go` files (layout, navigation, feedback, values and collections) were measured the same way, each group with a test that writes its screenshots (`BREEZE_SCREENSHOTS=<dir> go test ./ui/breeze -run 'TestScreenshots'`):

- **macOS:** AppKit's controls rendered at 2x on macOS 27.2 and compared pixel by pixel with Breeze UI's: NSButton's disclosure bezel, NSBox, NSSplitView's thin divider and a SwiftUI grouped `Form` with `DisclosureGroup`s; NSOutlineView of style source list, NSPathControl and a window's unified NSToolbar (which renders offscreen with its item platters, though not the Liquid Glass material); NSLevelIndicator in its capacity and rating styles and NSTextFinder's bar; NSSlider with tick marks, NSDatePicker of style text field, NSColorWell in its styles; an inset NSTableView, NSTokenField and NSComboBox. Colored fills were read from a render into an sRGB bitmap or resolved through `NSColor`, because the display profile that `reference.swift` reads shifts colors though not grays (NSLevelIndicator's green reads `#65c466` there, `#34c759` in sRGB). Tooltips, notification banners and contact monograms draw outside a view, so the macOS tooltip, the toast's panel and the avatar's monogram are inferred.
- **Windows 11:** microsoft-ui-xaml's main branch: the theme resources and templates of Expander, NavigationView (its items, headers, selection pill and back button, and `AnimatedBackVisualSource.cpp` for the arrow), BreadcrumbBar, CommandBar and AppBarButton, ToolTip, InfoBar, ProgressBar, RatingControl, PersonPicture (with `InitialsGenerator.cpp`), Slider and `TickBar_Partial.cpp`, CalendarDatePicker, TimePicker, DropDownButton, ColorPicker's checkerboard and ListViewItem (with `ListViewBaseItemChrome.cpp`), and WinUI Gallery's settings page for the card spacing. No Windows machine was at hand, so nothing was compared with a rendered WinUI control, and the glyphs of Segoe Fluent Icons (the chevrons of Expander and BreadcrumbBar) are drawn as paths inferred from their shapes.
- **GNOME:** libadwaita 1.7.6's `_expanders.scss`, `_lists.scss`, `_preferences.scss`, `_paned.scss`, `_sidebars.scss`, `_toolbars.scss`, `_tooltip.scss`, `_misc.scss` (the toast), `_level-bar.scss`, `_avatar.scss` with `adw-avatar.c` (its 14 gradients chosen by the name's hash), `_scale.scss`, `_entries.scss`, `_buttons.scss` and `_color-chooser.scss`; GTK 4.18.6's `gtkexpander.c`, `gtkcolordialogbutton.c` and `gtkcolorswatch.c`; adwaita-icon-theme 48's `pan-end-symbolic` and `go-previous-symbolic`, whose strokes Breeze UI draws; Nautilus 48's path bar for `Breadcrumbs` and GNOME Software 48's stars for `Rating`, as libadwaita has neither. These were not rendered with GTK for comparison.

The pointer, the keyboard, the focus and accessibility are MyGo's. The buttons, choices, slider, segments and tabs are built on MyGo's bases without a look (`ui.ButtonBase`, `ui.ToggleBase`, `ui.CheckboxBase`, `ui.RadioBase`, `ui.SwitchBase`, `ui.SelectBase`, `ui.SliderBase`, `ui.SegmentedBase`, `ui.TabsBase`), so screen readers, Tab, Space and the arrows work as they do for `ui`'s widgets. The text fields are `ui`'s own, with the box around their text drawn by the platform, so editing, input methods, the search field's clear button and Escape are MyGo's. `Disabled(true)` grays a widget out: MyGo draws disabled elements at half opacity, which is how AppKit and libadwaita draw them. WinUI has disabled colors of its own instead (a lighter fill and `TextFillColorDisabled`), which Breeze UI does not draw yet; half opacity comes close to them.

The alerts keep `ui.AlertDialog`'s contract: the last button is the default, which has the focus and which Enter clicks, Escape returns `ui.AlertDismissed`, and the index returned names the same button on every platform, though Windows places the default first.

The other widgets under `ui`'s names (`NumberInput`, `RadioGroup`, `Table`, `Tree`, `TreeItem`) are package `ui`'s in the platform's theme. MyGo offers no base without a look for them, even on its main branch, so drawing the platforms' rows and fields would mean writing their keyboard handling, scrolling and accessibility again; their colors, corner radius and font size come from `Use`. `RadioGroup` is a container around the platform's radio buttons. `MenuButton`'s and `PopoverMenuButton`'s buttons are the platform's; the menu of `MenuButton` is the system's own, and `PopoverMenuButton`'s is MyGo's in the platform's theme until MyGo releases its `MenuBase`. `Stepper` keeps MyGo's arrows and behavior inside the platform's box. `breeze-ui list` shows which widgets are drawn and which are themed.

These are themed too, for the same reason: `CheckboxGroup` (the platform's check boxes under MyGo's 16-DIP group box, which has no base), `Split`, `SplitVertical`, `SplitPercent` and `SplitVerticalPercent` (MyGo's 1-DIP divider in the platform's separator color), `FindBar` (MyGo's bar restyled as NSTextFinder's, a strip on Windows, which has none, or GtkSearchBar's), `RangeSlider`, `DecimalInput`, `Calendar`, `ColorPicker`, `GridView`, `Outline`, `OutlineTable` and `EditableText`. Some drawn widgets keep MyGo's parts that no base reaches: the headers of `Accordion`, the legend of `Fieldset`, the suggestion lists of `Autocomplete` and `Combobox`, and the chips of `TokenField`.

`Tooltip(c, anchor, text)` shows a tip drawn as the platform draws its own, built on `ui.TooltipBase`: 11-point text in a panel below the element on macOS, WinUI's ToolTip above it, GNOME's tooltip below it. `Element.Tooltip` keeps MyGo's look in the theme's colors. `Toasts(c)` draws the window's toasts, those of `c.Toast`, `c.ToastAction` and `c.AddToast`, as the platform would, on `ui.ToastViewportBase` and `ui.ToastBase`: GNOME's AdwToast, an InfoBar on Windows, and on macOS, which has no toast, a panel like a notification banner. Call it once in the view; without it MyGo draws its own toasts.

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
- **Buttons in a toolbar.** `Toolbar` and the history buttons are the platform's, but a `Button` or `Toggle` inside a `Toolbar` keeps its usual look: on macOS a 24-point push button beside 36-point toolbar platters, on Windows a Button rather than an AppBarButton. AppKit also joins adjacent toolbar items into one platter, which Breeze UI does not.
- **Parts MyGo draws inside a widget.** Without a base for them, these stay MyGo's in the platform's theme: the accordion's headers (SwiftUI puts the chevron before a regular-weight title, an Expander's header is 48 high, an AdwExpanderRow 50 high with an accent arrow when open), the fieldset's legend (weight 600, where GNOME's heading is 700), the check box group's 16-DIP box, the split's divider, which turns to the accent under the pointer, the suggestion lists of `Autocomplete` and `Combobox` (highlighted in the accent, not the platform's menu), the chips of `TokenField` (gray, not NSTokenField's flat blue tokens), the rows of `GridView` and `Outline`, the calendar icon of `DateInput` and the colon of `TimeInput`. A label of `Field` outside a `Form` is MyGo's weight 500, where WinUI's headers are regular.
- **Pressed rows.** `List` draws chosen and hovered rows as each platform does but not pressed ones, as MyGo's list rows do not report a press.
- **Rating's layout.** `Rating` returns `ui.Rating`'s row inside a column of its own, which takes the place in the parent's layout: give a margin or a grow to an element around it rather than to the row.
- **Tooltip and Description.** `Tooltip` sets the anchor's Description to its text, replacing one given before it, as package `ui` offers no way to read it: give a Description of your own after `Tooltip`.
