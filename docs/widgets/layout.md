# Layout

The widgets of [`ui/breeze/widgets_layout.go`](../../ui/breeze/widgets_layout.go), whose catalog entries are [`ui/breeze/catalog/layout.json`](../../ui/breeze/catalog/layout.json). What every widget shares, its platforms, themes and limits, is in [Breeze UI for MyGo](../breeze-ui.md).

## Looks

| Widget | macOS | Windows 11 | GNOME |
|---|---|---|---|
| `Collapsible` | disclosure chevron in a 13-point box, drawn where AppKit draws it, label 1.5 after the box, rows 24 high | Expander: 48-high header card, 32-pixel chevron button, content card padded 16 | GtkExpander: 16-pixel pan-end arrow at 70%, label right after it |
| `Accordion` | grouped form section: corners of 12, 3.1% fill | Expander cards 4 apart (MyGo's headers) | boxed list: card, corners of 12, card shadow |
| `AccordionItem` | row below a separator inset 10 points; rows in it as a DisclosureGroup's, closer and without separators | card with CardStrokeColorDefault, corners of 4; rows in it as a SettingsExpander's items | row below a `card_shade_color` separator; rows in it as AdwExpanderRow's nested rows on `card_shade_color` at 50% |
| `Fieldset` | legend over a grouped section (NSBox's fill), inset 10; in a Form the section reaches 10 past the fields; rows in it one section | legend over the fields, 10 above them (settings section header; no group box); rows in it SettingsCards 4 apart | bold legend centered in a 34-pixel header, the card 6 below it (AdwPreferencesGroup); in a Form the card reaches 12 past the fields; rows in it a boxed list |
| `Form` | labels in a column before the controls (ui.Form) | labels above the controls, 8 above, fields 12 apart | labels in a column before the controls (ui.Form) |
| `CheckboxGroup` | the head a check box with the accent dash while mixed, the boxes 22 in | the head a CheckBox with its indeterminate dash, the boxes 24 in | the head a check button with `dash.svg` while mixed, the boxes 24 in |
| `Split`, `SplitVertical`, `SplitPercent`, `SplitVerticalPercent` | 1-point divider in `separatorColor`, the same under the pointer | 1-pixel divider in the separator color | GtkPaned's 1-pixel separator |

`Scroll`, `ScrollBoth` and `ScrollHorizontal` are package `ui`'s too, under breeze's name so that a view can use breeze everywhere: MyGo paints a scroll container's bars itself, after its children, from the theme's `Scrollbar` color and `ScrollbarWidth` alone (`Painter.scrollbars` in `paint.go`), so with `Use` the thumb takes the platform's color and width (the overlay knob's black or white at 50% on macOS, WinUI's `ScrollBarThumbFill`, GNOME's 3-pixel slider at 20%), while its place, length, track and growth are MyGo's on every platform. No element can restyle that bar, and drawing another one would mean hiding MyGo's with `Scrollbars(false)`, which takes away its pointer targets, and writing the thumb's drag and the track's paging again.

## Sources

They were measured with the other widgets: see [Sources](../breeze-ui.md#sources).

## Limits

Besides [those of every widget](../breeze-ui.md#limits):

- **Text inside a sub-theme.** The panes of the `Split` family are built in a theme of their own, so their text takes the theme's color, size and font rather than those around them.
