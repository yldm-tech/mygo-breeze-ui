# Collections

The widgets of [`ui/breeze/widgets_collections.go`](../../ui/breeze/widgets_collections.go), whose catalog entries are [`ui/breeze/catalog/collections.json`](../../ui/breeze/catalog/collections.json). What every widget shares, its platforms, themes and limits, is in [Breeze UI for MyGo](../breeze-ui.md).

## Looks

| Widget | macOS | Windows 11 | GNOME |
|---|---|---|---|
| `List` | inset table view: 24-point rows, 10 in from the edges, the text 8 in from the selection; the chosen row rounded 8 in `selectedContentBackgroundColor` with white text, `unemphasizedSelectedContentBackgroundColor` without the focus; rows chosen together make one block rounded at its ends; no hover | ListView: 40-pixel items, `SubtleFillColorSecondary` backplate inset 4,2 with corners of 4 (`SubtleFillColorTertiary` chosen and under the pointer), 3 by 16 `AccentFillColorDefault` pill on the chosen one; pressed states not drawn | GtkListView `.rich-list`: 48-pixel rows on `view_bg_color`, the chosen one the accent at 25% (32% under the pointer), the hover `currentColor` at 4%; pressed states not drawn |
| `Autocomplete` | bezeled field and the pop-up button's menu | TextBox and AutoSuggestBox's list of ListViewItems | entry and a popover menu |
| `TokenField` | bezeled field around MyGo's gray chips (not NSTokenField's blue tokens) | TextBox around MyGo's chips (not TokenizingTextBox's) | entry around MyGo's chips (GNOME has no token field) |
| `Outline` | inset rows with a rounded selection and a gray or white chevron | TreeViewItem's backplate and accent pill | tinted rows with pan chevrons |
| `GridView` | the chosen item a rounded fill, the accent with the focus, gray without | GridViewItem's accent double border | tinted rounded item |
| `EditableText` | a text field's frame while editing (inferred) | TextBox's frame (inferred) | an entry's frame |
| `OutlineTable` | inset rows, the chosen one rounded in the accent (gray without the focus), 10 in from the edges; the turning chevron | ListViewItem's backplate and 3 by 16 accent pill (WinUI has no table) | full-width rows tinted with the accent on the view's background; pan-end chevrons |

## Sources

They were measured with the other widgets: see [Sources](../breeze-ui.md#sources).

## Limits

Besides [those of every widget](../breeze-ui.md#limits):

- **Shared outline states.** `Outline` and `OutlineTable` set an `OutlineState`'s `Indent` to the platform's when it is 0; a state shared by several `Preview`s keeps the first platform's.
