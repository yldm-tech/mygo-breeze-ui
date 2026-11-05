package breeze

import "github.com/egoist/mygo/ui"

// Collections: lists, grids, outlines and the fields that hold several values or edit a value in place.
//
// MyGo builds a List's row inside an element of its own, which it fills with the theme's accent while the row is chosen; the row builder runs inside it. Breeze UI's List covers that element with a cell of its own, in the list's background, and draws the platform's chosen and hovered row in it: AppKit's inset table view, WinUI's ListViewItem and libadwaita's GtkListView row. The choice, the keys, the focus and accessibility stay MyGo's. MyGo's other collections (GridView, Outline, OutlineTable) build their rows' and cells' frames around the content, out of reach of a cover, so they are package ui's in the platform's theme.

// List creates a scrolling list of n rows, as the platform's table view (macOS), list view (Windows) or list view (GNOME) looks: row builds row i inside the platform's row, which shows the platform's chosen and hovered row while ListState's Selected or Selection lets the user choose rows. It builds only the rows in view, as ui.List does, and keeps its keys, its choice and its accessibility. Without row, configure ItemKey and Selection and build the rows with Rows, as with ui.List: those rows are MyGo's, in the platform's theme, inside the platform's list. It is ui.List's counterpart.
func List(c *ui.Context, s *ui.ListState, n int, row ...func(int)) ui.Element {
	p := PlatformOf(c)
	if len(row) != 1 {
		// ui.List builds the rows given (and refuses more than one), or waits for Rows.
		return styleList(c, p, ui.List(c, s, n, row...))
	}
	build := row[0]
	choosable := s != nil && (s.Selected != nil || s.Selection != nil)
	// The rows chosen in this frame, which a chosen row's selection joins on macOS as it is painted, once the rows in view are all built.
	chosen := map[int]bool{}
	l := styleList(c, p, ui.List(c, s, n))
	l.Rows(func(r ui.ListRow) {
		// The cell fills the row MyGo builds, which has no padding of its own, and covers what MyGo paints there.
		cell := ui.Box(c).Grow(1)
		// A header of a section is not chosen, and does not light up under the pointer, as MyGo's does not.
		hovers := choosable && (s.Header == nil || !s.Header(r.Index))
		switch p {
		case Windows:
			fluentListRow(c, cell, r, hovers)
		case Linux:
			adwListRow(c, cell, r, hovers)
		default:
			macListRow(c, cell, r, chosen)
		}
		cell.Children(func() { build(r.Index) })
	})
	return l
}

// styleList gives a list the platform's background and the room around its rows.
func styleList(c *ui.Context, p Platform, l ui.Element) ui.Element {
	t := c.Theme()
	switch p {
	case Windows:
		// A ListView has no background of its own; the window's shows through, which the rows' cells cover with.
		l.Background(t.Background)
	case Linux:
		l.Background(adwListViewBg(t))
	default:
		// NSTableView of style inset: controlBackgroundColor, the window's own color on macOS 26, and its rows 10 points in from every edge (measured: the first row's rect starts 10 points down, its selection 10 points in from each side).
		l.Background(t.Background).Padding(10)
	}
	return l
}

// ---- macOS: NSTableView of style inset ----

// macListRow draws a row of an inset NSTableView: 24 points high (the rows of rowSizeStyle default, measured), the content 8 points in from the selection (measured: the cell view starts 16 points into the row, 6 into its selection, and a label's text 2 into its NSTextField; the text field is taken flush with the cell view, as it was in the rendered table, inferred for Interface Builder's template, which may set it further in), and while chosen a rectangle with corners of 8 (measured from the drawSelection of a row view in an inset table at 2x, matching a rounded rectangle of radius 8) in selectedContentBackgroundColor with white text while the list has the focus, or unemphasizedSelectedContentBackgroundColor with the text's own color while it has not. Rows chosen together make one block, rounded at its ends only: NSTableRowView squares the corners on the side of a chosen previous or next row (measured from drawSelection with isPreviousRowSelected and isNextRowSelected). AppKit's rows do not change under the pointer. chosen gathers the rows chosen in the frame, for the block.
func macListRow(c *ui.Context, cell ui.Element, r ui.ListRow, chosen map[int]bool) {
	t := c.Theme()
	// The text's own color, rather than the accent's text MyGo gives the chosen row.
	cell.MinHeight(24).PaddingX(8).Justify(ui.Center).Background(t.Background).TextColor(t.Text)
	if !r.Selected() {
		return
	}
	i := r.Index
	chosen[i] = true
	focused := r.ListFocused()
	fill := macListSelection(t, focused)
	if focused {
		// alternateSelectedControlTextColor
		cell.TextColor(white)
	}
	cell.Draw(func(p *ui.Painter, rect ui.Rect) {
		p.Fill(rect, fill, 8)
		// The rows beside it are built by now: square the corners on the side of a chosen one.
		half := rect.H / 2
		if chosen[i-1] {
			p.Fill(ui.Rect{X: rect.X, Y: rect.Y, W: rect.W, H: half}, fill, 0)
		}
		if chosen[i+1] {
			p.Fill(ui.Rect{X: rect.X, Y: rect.Y + half, W: rect.W, H: rect.H - half}, fill, 0)
		}
	})
}

// macListSelection returns the fill of a chosen row: selectedContentBackgroundColor while the list has the focus, #0064e1 in Aqua and #0059d1 in Dark Aqua for the blue accent (resolved in AppKit), and for another accent that accent mixed with 16% black in Aqua and 24% in Dark Aqua, the mix closest to what AppKit resolves for the red, orange, green, purple and pink accents (within 13 and 18 of 255 per channel; AppKit darkens the accent's Oklab lightness and chroma rather than mixing in black); unemphasizedSelectedContentBackgroundColor while it has not, #dcdcdc and #464646 whatever the accent (resolved in AppKit).
func macListSelection(t *ui.Theme, focused bool) ui.Color {
	switch {
	case !focused:
		return pick(t.Dark, ui.Hex("#dcdcdc"), ui.Hex("#464646"))
	case t.Accent == macAccent:
		return pick(t.Dark, ui.Hex("#0064e1"), ui.Hex("#0059d1"))
	}
	return t.Accent.Mix(black, pick(t.Dark, float32(0.16), 0.24))
}

// ---- Windows 11: WinUI 3's ListView ----

// fluentListRow draws a ListViewItem (ListViewItem_themeresources.xaml and ListViewBaseItemChrome.cpp): at least 40 DIPs high (ListViewItemMinHeight), its content after the inline selection indicator's 7 DIPs and the item's Padding 16,0,12,0, in TextFillColorPrimary whatever its state; a backplate 4,2 in from its edges (s_backplateMargin) with corners of 4 (ListViewItemCornerRadius), SubtleFillColorSecondary while chosen or under the pointer and SubtleFillColorTertiary while chosen and under the pointer; and on the chosen one a 3 by 16 DIP pill of AccentFillColorDefault, 4 DIPs from the item's start, centered, with corners of 1.5 (s_selectionIndicatorSize, s_selectionIndicatorMargin, ListViewItemSelectionIndicatorCornerRadius), growing with an item taller than 56. The backplate shows under the pointer only when hovers, for a row that can be chosen. The pressed item's backplates (SubtleFillColorTertiary, or SubtleFillColorSecondary while chosen) and its pill 6 DIPs shorter (s_selectionIndicatorHeightShrinkage) are not drawn: the press belongs to MyGo's row, which ListRow does not report, and asking the cell for it would take the row's clicks.
func fluentListRow(c *ui.Context, cell ui.Element, r ui.ListRow, hovers bool) {
	t := c.Theme()
	cell.MinHeight(40).PaddingStart(16 + 7).PaddingEnd(12).Justify(ui.Center).Background(t.Background).TextColor(t.Text)
	secondary, tertiary := fluentSubtle(t)
	chosen, hovered := r.Selected(), hovers && cell.Hovered()
	var fill ui.Color
	switch {
	case chosen && hovered:
		fill = tertiary
	case chosen || hovered:
		fill = secondary
	}
	rtl := c.Direction() == ui.RightToLeft
	cell.Draw(func(p *ui.Painter, rect ui.Rect) {
		if fill.A > 0 {
			p.Fill(ui.Rect{X: rect.X + 4, Y: rect.Y + 2, W: rect.W - 8, H: rect.H - 4}, fill, 4)
		}
		if !chosen {
			return
		}
		h := min(rect.H, max(16, rect.H-40))
		x := rect.X + 4
		if rtl {
			x = rect.X + rect.W - 4 - 3
		}
		p.Fill(ui.Rect{X: x, Y: rect.Y + (rect.H-h)/2, W: 3, H: h}, t.Accent, 1.5)
	})
}

// ---- GNOME: GtkListView of class rich-list ----

// adwListViewBg is view_bg_color, the background of lists and text views: #ffffff and #1d1d20 (_colors.scss of libadwaita 1.10.0).
func adwListViewBg(t *ui.Theme) ui.Color { return pick(t.Dark, white, ui.Hex("#1d1d20")) }

// adwListRow draws a row of a GtkListView of class rich-list (_lists.scss): padded 8px 12px with a min-height of 32px, so 48 pixels high, square, in view_fg_color over view_bg_color; while chosen $view_selected_color, the accent at 25% ($view_selected_hover_color, 32%, under the pointer), and under the pointer $view_hover_color, currentColor at 4%, as rows that can be chosen are activatable: the hover shows only when hovers. The pressed row's $view_active_color (8%) and $view_selected_active_color (39%) are not drawn, the press being MyGo's row's, as on Windows.
func adwListRow(c *ui.Context, cell ui.Element, r ui.ListRow, hovers bool) {
	t := c.Theme()
	bg := adwListViewBg(t)
	cell.MinHeight(48).Padding(8, 12).Justify(ui.Center).Background(bg).TextColor(adwTextOn(t, bg))
	chosen, hovered := r.Selected(), hovers && cell.Hovered()
	var fill ui.Color
	switch {
	case chosen && hovered:
		fill = bg.Mix(t.Accent, 0.32)
	case chosen:
		fill = bg.Mix(t.Accent, 0.25)
	case hovered:
		fill = adwInkAlpha(t, bg, 0.04)
	}
	if fill.A > 0 {
		cell.Background(fill)
	}
}

// ---- package ui's collections, in the platform's theme (Use) ----

// GridView creates a scrolling grid of n items at least minWidth wide and height high, built with item(i) only as they show, as package ui's GridView does: the items chosen, their frames and the keys are MyGo's, in the platform's theme. MyGo builds each item's frame around what item builds, so Breeze UI cannot draw NSCollectionView's, GridView's or GtkGridView's cells there. It is ui.GridView's counterpart.
func GridView(c *ui.Context, s *ui.GridState, n int, minWidth, height float32, item func(i int)) ui.Element {
	return ui.GridView(c, s, n, minWidth, height, item)
}

// Outline creates a list of items in a tree, roots at the top and below each item open its children, which children returns, nil for an item without any; row builds an item's content after its arrow. The arrows, the indentation and the chosen row are MyGo's, in the platform's theme: MyGo builds them around what row builds, out of reach of NSOutlineView's, TreeView's or GtkTreeExpander's look. It is ui.Outline's counterpart.
func Outline[K comparable](c *ui.Context, s *ui.OutlineState[K], roots []K, children func(K) []K, row func(item K)) ui.Element {
	return ui.Outline(c, s, roots, children, row)
}

// OutlineTable creates an outline whose rows are a Table's: the first of columns shows the items' arrows, indented as they are deep, before what cell builds. It is package ui's, in the platform's theme, as Table and Outline are. It is ui.OutlineTable's counterpart.
func OutlineTable[K comparable](c *ui.Context, s *ui.OutlineState[K], columns []ui.TableColumn, roots []K, children func(K) []K, cell func(item K, col int)) ui.Element {
	return ui.OutlineTable(c, s, columns, roots, children, cell)
}

// EditableText shows *value as text that a double click, or Enter while it has the focus, turns into a text input editing it, as a file's name in Finder; Enter keeps what was typed and Escape goes back. It is package ui's, in the platform's theme: the input appears inside it as it is edited, out of reach of a style applied to what it returns. It is ui.EditableText's counterpart.
func EditableText(c *ui.Context, value *string) ui.Element {
	return ui.EditableText(c, value)
}

// ---- fields of package ui, in the platform's text field box ----

// Autocomplete creates a text input editing *value that suggests below it those of suggestions containing what was typed, those starting with it first, in the box the platform draws around its text fields: Up and Down move among them and Enter or a click takes one, the caret going after it, while Enter on none is Submitted, and Escape closes the list. The list is drawn as the platform's menu: NSMenu's panel and rows on macOS, an AutoSuggestBox's suggestion list on Windows, a dropdown's popover on GNOME. It is ui.Autocomplete's counterpart, built as it is on ui.ComboboxBase, which keeps its keys, its focus and what assistive technology sees.
func Autocomplete(c *ui.Context, value *string, suggestions []string) ui.Element {
	return styleField(c, idiomAutocomplete(c, value, suggestions), singleLine)
}

// TokenField creates a field of tokens, as of tags or the recipients of a mail, in the box the platform draws around its text fields: *tokens shows as chips before an input where Enter or a comma adds what was typed, Backspace in it while empty takes out the last, and the suggestions containing what was typed show below it. The chips and the list are MyGo's, in the platform's theme: macOS's NSTokenField and Windows's TokenizingTextBox draw their own tokens, GNOME has no token field. It is ui.TokenField's counterpart.
func TokenField(c *ui.Context, tokens *[]string, suggestions []string) ui.Element {
	f := styleField(c, ui.TokenField(c, tokens, suggestions), singleLine)
	if PlatformOf(c) == Linux {
		// A GtkEntry has no vertical padding; a field whose chips wrap onto more lines keeps them off its edges as its first line is (inferred: (34 - a chip's height) / 2).
		f.PaddingY(6)
	}
	return f
}
