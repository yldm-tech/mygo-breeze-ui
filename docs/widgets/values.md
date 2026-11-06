# Values

The widgets of [`ui/breeze/widgets_values.go`](../../ui/breeze/widgets_values.go), whose catalog entries are [`ui/breeze/catalog/values.json`](../../ui/breeze/catalog/values.json). What every widget shares, its platforms, themes and limits, is in [Breeze UI for MyGo](../breeze-ui.md).

## Looks

| Widget | macOS | Windows 11 | GNOME |
|---|---|---|---|
| `StepSlider` | slider with 2-point dots 3 below the track, white knob in both modes, 16 high | Slider with 1-pixel ticks 4 below (TickPlacement BottomRight), focus ring around the whole slider | scale with 6-pixel marks after it, slider pointing down |
| `DateInput` | NSDatePicker textField bezel, 23 high, corners of 7 | CalendarDatePicker: button face, 32 high, glyph column at the end | entry, 34 high (GTK has no date field) |
| `TimeInput` | NSDatePicker textField bezel around MyGo's segments | TimePicker's button face, 242 wide, 32 high, segments spread across it | entry with the focus ring while a segment has the focus |
| `ColorWell` | NSColorWell: 48 by 24 capsule of systemFill around a capsule swatch; a translucent color split over black and white | DropDownButton with a 36 by 20 swatch and a chevron, 32 high; ColorPicker's checkerboard | GtkColorDialogButton: button padded 5 around a swatch, 34 high; GTK's 10-pixel checkerboard |

These stay package `ui`'s in the platform's theme. Their parts are built inside them, out of reach of a style applied to what they return, and each was checked against MyGo v0.7.4's source for a rebuild from its exported bases, as `NumberInput` is rebuilt; each needs something only package `ui` has:

- `DecimalInput` marks what is typed wrong with an error in MyGo's translations ("Enter a number", "Use at most 2 decimal places", "Enter 0 or more"), which come from MyGo's internal `l10n` package. A rebuild could read an app's own `mygo.*` overrides through a catalog of its own in the namespace `mygo` (`i18n.NewLibraryCatalog`), but not MyGo's translations, which the packages `mygo/locales/<tag>` an app imports add to that internal catalog only: it would have to carry copies of them, which would show languages the app did not import and fall out of step with MyGo's. `NumberInput` shows no error, which is why it could be rebuilt.
- `RangeSlider` could have its track and knobs drawn in place of MyGo's, as `StepSlider` has, since `Draw` on what it returns replaces MyGo's painting; but its knobs are elements inside it, and only MyGo's painting can ask which of them has the focus from the keyboard, so the focus ring would be lost. Knobs of its own, focusable elements with `Role(ui.RoleSlider)` and `Range`, could take the keys and the roles, but MyGo names its knobs after the slider's `Label`, given after the call ("Price minimum"), through a link between elements that is not exported.
- `Calendar` names its month buttons in MyGo's translations ("Previous month", "Next month"), and `Browse`, which an app calls on what it returns, only acts on MyGo's own calendar. Telling assistive technology which day the keys are on would not block a rebuild: `Element.ActiveDescendant` is exported. `CalendarView` draws the platforms' calendars without MyGo's translations, its names and dates from the exported `i18n` package and its buttons' names from Breeze UI's own catalog; `Calendar` stays MyGo's for `Browse`.
- `ColorPicker` names its square and its hex field, describes the square's value and names its swatches' colors in MyGo's translations ("Saturation", "Hex", the color names).

## Sources

They were measured with the other widgets: see [Sources](../breeze-ui.md#sources).

## Limits

Only [those of every widget](../breeze-ui.md#limits).
