# Feedback

The widgets of [`ui/breeze/widgets_feedback.go`](../../ui/breeze/widgets_feedback.go), whose catalog entries are [`ui/breeze/catalog/feedback.json`](../../ui/breeze/catalog/feedback.json). What every widget shares, its platforms, themes and limits, is in [Breeze UI for MyGo](../breeze-ui.md).

## Looks

| Widget | macOS | Windows 11 | GNOME |
|---|---|---|---|
| `Tooltip` | 11-point text (`NSFont.toolTipsFont`) in a light or dark panel below the element; panel inferred | ToolTip: in-app acrylic fallback (#f9f9f9, #2c2c2c), 1-pixel flyout stroke, 4-pixel corners, padding 9,6,9,8, 12-pixel text, above the element | GTK tooltip: the shade at 80%, white text, 1-pixel white 10% border, 9-pixel corners, padding 6 10, below the element |
| `Toasts` | panel like a notification banner (macOS has no toast; inferred) | InfoBar with a severity icon, colored by the toast's Type, 38-pixel close button | AdwToast: #505053 pill, heading title, action button (426 wide with it), circular flat close button |
| `Meter` | NSLevelIndicator: 16-point bar, systemGreen, systemYellow or systemRed with an edge 13% darker | ProgressBar look: 3-pixel bar over a 1-pixel track, caution and critical colors | GtkLevelBar: 8-pixel round trough, accent, warning, error or success |
| `Avatar` | gray monogram, white initials, a silhouette without a name (inferred) | PersonPicture: gray ellipse, semibold initials at 42% by InitialsGenerator's rules, Contact glyph without them | AdwAvatar: one of 14 gradients by the name's hash, avatar-default-symbolic without a name, white initials below 25 pixels |
| `Rating` | NSLevelIndicator rating: 12-point stars, secondaryLabelColor and systemFill | RatingControl: 16-pixel accent stars 8 apart, growing under the pointer; faint fill hovering an unset rating | GNOME Software's stars: 32 pixels, yellow (GNOME has no rating control) |

`Tooltip(c, anchor, text)` shows a tip drawn as the platform draws its own, built on `ui.TooltipBase`: 11-point text in a panel below the element on macOS, WinUI's ToolTip above it, GNOME's tooltip below it. `Element.Tooltip` keeps MyGo's look in the theme's colors. `Toasts(c)` draws the window's toasts, those of `c.Toast`, `c.ToastAction` and `c.AddToast`, as the platform would, on `ui.ToastViewportBase` and `ui.ToastBase`: GNOME's AdwToast, an InfoBar on Windows, and on macOS, which has no toast, a panel like a notification banner. Call it once in the view; without it MyGo draws its own toasts.

`FindBar` stays package `ui`'s in the platform's theme, for the reason [values](values.md#looks) gives for its themed widgets:

- `FindBar` names its bar, field and buttons and says its count in MyGo's translations ("Find", "No matches", "2 of 5", "Previous", "Next", "Done"), from the same package. Breeze UI draws the bar itself (NSTextFinder's bar, a strip on Windows, which has none, or GtkSearchBar's); the field and buttons inside it stay MyGo's.

## Sources

They were measured with the other widgets: see [Sources](../breeze-ui.md#sources).

## Limits

Besides [those of every widget](../breeze-ui.md#limits):

- **Rating's layout.** `Rating` returns `ui.Rating`'s row inside a column of its own, which takes the place in the parent's layout: give a margin or a grow to an element around it rather than to the row.
- **Tooltip and Description.** `Tooltip` sets the anchor's Description to its text, replacing one given before it, as package `ui` offers no way to read it: give a Description of your own after `Tooltip`.
