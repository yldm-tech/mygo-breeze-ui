# Core controls

The widgets of [`ui/breeze/widgets.go`](../../ui/breeze/widgets.go), whose catalog entries are [`ui/breeze/catalog/widgets.json`](../../ui/breeze/catalog/widgets.json). What every widget shares, its platforms, themes and limits, is in [Breeze UI for MyGo](../breeze-ui.md).

## Looks

| Widget | macOS | Windows 11 | GNOME |
|---|---|---|---|
| `Button` | push button, 24 points high, no border | button with its elevation border, 32 high | text button, 34 high, bold |
| `PrimaryButton` | default button in the accent color | accent button | suggested action |
| `Toggle` | push-on-push-off button, the accent with white text while on | ToggleButton, the accent while on | toggle button, darker while checked |
| `Checkbox` | 16-point box with corners of 5 | 20-pixel box with a stroke | 20-pixel box with a trough border |
| `Radio` | 16-point circle, 5-point dot | 20-pixel circle, dot that grows under the pointer | 20-pixel circle, 8-pixel bullet |
| `RadioGroup` | radio buttons 6 points apart down a column, 8 across a row (SwiftUI's radio-group picker) | RadioButtons: rows 8 apart, columns 7 | check buttons next to each other (GNOME has no radio group) |
| `Switch` | 54 by 24 track, capsule knob | 40 by 20 track, knob that grows under the pointer | 46 by 26 track, round slider |
| `TextInput`, `TextInputBuffer` | bezeled field, 26 high, focus ring around it | TextBox, accent line along the bottom while focused | filled entry, 34 high, inner focus ring |
| `TextArea`, `TextAreaBuffer` | text view in a bezel border | TextBox | text view in a frame |
| `SearchField` | capsule search field | AutoSuggestBox | search entry |
| `Combobox` | bezeled field around MyGo's list | TextBox around MyGo's list | entry around MyGo's list |
| `Select` | pop-up button and menu with a check mark | ComboBox, accent pill on the chosen item | drop-down with the pan-down chevron, popover with a check mark |
| `Slider` | 6-point track, 20 by 16 capsule knob | 4-pixel track, thumb with an accent dot | 4-pixel trough, 20-pixel slider |
| `Progress` | 8-point bar | 3-pixel bar over a 1-pixel line | 8-pixel bar |
| `Spinner` | 8 spokes 4 points wide, 32 points | ProgressRing, 32 pixels, an arc of at most half the circle | AdwSpinner, 16 pixels, 2-pixel stroke |
| `Segmented` | NSSegmentedControl, the chosen segment in the accent | SelectorBar | AdwToggleGroup |
| `Tabs` | tab view, the chosen tab on a raised neutral platter | SelectorBar | notebook tabs with the accent line |
| `Modal` | sheet with a window's shadow, the window behind not dimmed | ContentDialog over the smoke layer | floating AdwDialog over the dimmed window, with libadwaita 1.10's lighter shadows |
| `AlertDialog` | NSAlert: 260 points, centered text, capsule buttons, the default last | ContentDialog: title and message above a command area, the default first | AdwAlertDialog: title-2 heading, body at 140% line height, 44-pixel buttons, the default last, a 1-pixel hairline inside its edge |
| `NumberInput` | text field with an NSStepper beside it | NumberBox with inline spin buttons | GtkSpinButton: entry with flat − and + buttons after it |
| `Tree`, `TreeItem` | the platform's row height, inset and corners, the selection in its color | the same | the same |
| `Table` | inset rows, the chosen one rounded in the accent (gray without the focus), 10 in from the edges; the turning chevron | ListViewItem's backplate and 3 by 16 accent pill (WinUI has no table) | full-width rows tinted with the accent on the view's background; pan-end chevrons |

The alerts keep `ui.AlertDialog`'s contract: the last button is the default, which has the focus and which Enter clicks, Escape returns `ui.AlertDismissed`, and the index returned names the same button on every platform, though Windows places the default first.

## Sources

They were measured with the other widgets: see [Sources](../breeze-ui.md#sources).

## Limits

Besides [those of every widget](../breeze-ui.md#limits):

- **Dialogs.** An alert shows no icon: NSAlert shows the app's icon above its title, which a view cannot reach. GNOME's alert stacks its buttons when they do not fit side by side; Breeze UI keeps them in a row. On macOS the corner radius and the alert's colors are measured from NSAlert's window; the hairline around a dialog and its shadow are drawn to resemble AppKit's, not measured.
- **Errors on text fields.** The platforms' text fields draw their own borders and focus, so a field's `Error` does not turn its border red as `ui`'s fields do.
- **Radio groups.** `RadioGroup` sets the platform's spacing but shows no header: its `Label` names it for assistive technology only, where WinUI's RadioButtons shows its Header above the buttons. Across a row, WinUI's RadioButtons gives every button the width of the widest; Breeze UI's keep their own widths, at least 120 as a RadioButton's.
