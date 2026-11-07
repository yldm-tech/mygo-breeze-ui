# Extras

The widgets of [`ui/breeze/extras.go`](../../ui/breeze/extras.go), whose catalog entries are [`ui/breeze/catalog/extras.json`](../../ui/breeze/catalog/extras.json). What every widget shares, its platforms, themes and limits, is in [Breeze UI for MyGo](../breeze-ui.md).

## Looks

| Widget | macOS | Windows 11 | GNOME |
|---|---|---|---|
| `Link` | `linkColor`, underlined under the pointer | HyperlinkButton's accent text colors, underlined | `accent_color`, underlined |
| `Divider` | `separatorColor` | `DividerStrokeColorDefault` | `$border_color` |
| `Badge` | MyGo's (macOS has no badge control) | InfoBadge: 16 high, accent | indicator badge: 13 high, 40% gray, white bold text |
| `MenuButton`, `PopoverMenuButton` | push button with a 7.5 by 4 chevron after the label | DropDownButton with ChevronDownSmall | `.arrow-button` with pan-down |
| `Popover` | panel like a menu | FlyoutPresenter, at least 96 wide | popover contents with a 1-pixel 5% ring |
| `ToggleGroup` | segmented control selecting any segments | ToggleButtons side by side (WinUI joins none) | `box.linked`, with the seam where two buttons meet |
| `Field` | form labels in a column before the controls (in a `Form`) | header at regular weight 8 above the control | label above the control |
| `Stepper` | NSStepper's box, 20 points wide, a line between its halves | a button's face and elevation border | a button's face |

## Sources

They were measured with the other widgets: see [Sources](../breeze-ui.md#sources).

## Limits

Only [those of every widget](../breeze-ui.md#limits).
