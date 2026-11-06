# Navigation

The widgets of [`ui/breeze/widgets_navigation.go`](../../ui/breeze/widgets_navigation.go), whose catalog entries are [`ui/breeze/catalog/navigation.json`](../../ui/breeze/catalog/navigation.json). What every widget shares, its platforms, themes and limits, is in [Breeze UI for MyGo](../breeze-ui.md).

## Looks

| Widget | macOS | Windows 11 | GNOME |
|---|---|---|---|
| `Sidebar` | source list, 10-point inset, its own gray background (#f0f0f0, #282828) | NavigationView's left pane over the window | AdwSidebar on `sidebar_bg_color`, a line at its end |
| `SidebarSection` | group row: 11-point semibold title in `tertiaryLabelColor`, text placed as AppKit's (13 points from the top, 13 below the row above) | NavigationViewItemHeader: a 40-pixel row, 14 pixels, semibold, secondary text, 16 from the pane's edge | AdwSidebar's heading: bold, in the sidebar's text color, 12 in from the edge |
| `SidebarItem` | 32-point row, corners of 8, accent selection while focused, gray otherwise, label semibold when chosen; the icon in the accent (#0070f5, #007aff), white on the focused selection | NavigationViewItem: 36 high, subtle fill (tertiary pressed), 3 by 16 accent pill; the icon in the text color | AdwSidebar row: 42 pixels, padding 3 14, 10% when chosen, 7% hovered, 16% pressed; the icon in the text color |
| `SidebarItemImage` | `SidebarItem`'s row with a 16-point image in its own colors, 9.5 in, kept on the accent selection | `SidebarItem`'s row with a 16-DIP image in its own colors in the 40-DIP column | `SidebarItem`'s row with a 16-pixel image in its own colors, 12 before the label |
| `SidebarSectionWithActions` | `SidebarSection` with controls at the end of the title's row, 2 points inside the selection's inset, centered on the title, before the disclosure arrow of an open section | `SidebarSection` with controls 12 from the pane's edge (inferred) | `SidebarSection` with controls at the heading's end, 12 in (inferred) |
| `Breadcrumbs` | NSPathControl: secondary gray items and 4 by 6 chevrons | BreadcrumbBar: 14-pixel items, 12-pixel chevrons | Files' path bar: bold buttons, slashes, in a 10% trough |
| `BackButton`, `ForwardButton` | toolbar item: 36-point platter with a hairline, shared by the two in a toolbar, SF Symbols' chevron | NavigationView back button: 40 by 36, subtle, WinUI's back arrow | flat image button, 34 pixels, `go-previous-symbolic` |
| `Toolbar` | 52 points high, items 8 from its edges and 8 apart | CommandBar, 48 high | `.toolbar`: padding and spacing 6 |
| `Button`, `PrimaryButton`, `Toggle`, `MenuButton`, `PopoverMenuButton` in a `Toolbar` | 36-point capsules padded 12, adjacent image items sharing one platter | AppBarButton: transparent, subtle fill under the pointer, the accent when checked | text buttons stay raised; image buttons flat, 34 pixels, with the focus ring inside |

## Sources

They were measured with the other widgets: see [Sources](../breeze-ui.md#sources).

`appkit-reference`'s sidebar scene (`renderSidebar`, `appkit-sidebar-*.png`) renders an NSOutlineView of style sourceList holding template SF Symbols and a colored image, chosen with and without the focus: AppKit draws a template image in the accent (#0070f5 in Aqua, #007aff in Dark Aqua, the emphasized selection's fill), on a row not chosen, on the selection without the focus and under the pointer, and white on the emphasized selection, while an image that is not a template keeps its colors there too; `SidebarItem` and `SidebarItemImage` draw them so.

## Limits

Besides [those of every widget](../breeze-ui.md#limits):

- **Toolbars.** Buttons built inside a `Toolbar` take its look, and so do buttons in a `Popover` or `Modal` opened from inside the toolbar's build, as MyGo's own toolbar look reaches them too. On macOS adjacent image items join into one platter from the second frame, as joining reads the first frame's layout. The overflow button MyGo adds when items do not fit is MyGo's.
- **Sidebar sections with actions.** `SidebarSectionWithActions` lays its controls over the end of MyGo's title row, which has no slot for them: a title long enough to reach them runs under them, and controls taller than the row grow it as much above as below from the next frame, as the row's height is read from the last layout (a new section shows them overflowing it for one frame). They are drawn as anywhere else, so an `IconButton` is a bordered button where a SwiftUI section header holds a borderless one; their places on Windows and GNOME, which have no such control, are inferred. They are Tab stops after the sidebar, not inside it, and assistive technology finds them after the section's items.
