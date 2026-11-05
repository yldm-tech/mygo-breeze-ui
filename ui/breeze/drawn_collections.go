package breeze

import (
	"math"
	"sync"
	"time"

	"github.com/egoist/mygo/ui"
)

// ---- collections drawn over MyGo's rows and cells ----
//
// MyGo does not clip what an element's Draw callback paints to the element's box (ui/paint.go Painter.element calls p.draw with no clip of its own; only an element with Clip or a scroll clips its children, to its padding box), and an element paints after its earlier siblings and its ancestors. So an element Breeze UI builds inside one of MyGo's cells or rows, the cover, can paint over the cell's border and fill, over the row's accent fill and over the disclosure arrow MyGo built before it, then draw the platform's row there. The cover takes no clicks (Hovered only asks for the pointer's hover), so the pointer, the keys, the focus and what assistive technology sees stay MyGo's. What a cover cannot change is MyGo's layout: the rows' and cells' padding, the gaps, the arrow's box and the rows' least heights.

// collBg returns the background of the platform's collection views, which the covers paint over MyGo's fills: controlBackgroundColor on macOS, the window's own color there since macOS 26 (as styleList says), the window's SolidBackgroundFillColorBase behind WinUI's ListView and TreeView, which have none of their own, and view_bg_color on GNOME (_lists.scss: listview and columnview are .view).
func collBg(p Platform, t *ui.Theme) ui.Color {
	if p == Linux {
		return adwListViewBg(t)
	}
	return t.Background
}

// collText returns the text color of a row: alternateSelectedControlTextColor, white, on a macOS selection while the view has the focus, and the text's own color otherwise (TextFillColorPrimary on Windows, whatever the item's state, and view_fg_color on GNOME), rather than the accent's text MyGo gives a chosen row.
func collText(p Platform, t *ui.Theme, bg ui.Color, chosen, focused bool) ui.Color {
	switch {
	case p == Linux:
		return adwTextOn(t, bg)
	case p == Mac && chosen && focused:
		return white
	}
	return t.Text
}

// collIndent returns how far each level of a tree indents its rows: NSOutlineView's indentationPerLevel of 13 points (measured on an inset outline view on macOS 27.2), TreeViewItem's 16 DIPs a level (TreeViewItem.cpp, depth * 16) and, on GNOME, a GtkTreeExpander's indent node of 16 pixels and the 4 pixels of border-spacing after it (gtktreeexpander.c prepends one "indent" GtkBuiltinIcon a level; _column-view.scss: treeexpander { border-spacing: 4px }; the icon's 16 pixels are -gtk-icon-size's default, inferred).
func collIndent(p Platform) float32 {
	switch p {
	case Windows:
		return 16
	case Linux:
		return 20
	}
	return 13
}

// collRow is a row, or the slice of a row one of a table's cells holds, as collPaintRow paints it: whether it is chosen, under the pointer and in a view with the focus, whether the rows above and below it are chosen, which join AppKit's selections into one block, and whether the slice holds the row's start and its end, in the reading direction.
type collRow struct {
	chosen, hovered, focused bool
	above, below             bool
	start, end               bool
	rtl                      bool
}

// collRowLook is where a collection puts the platform's row in the rect it paints: macInset and fluentInset are how far in from the row's start and end AppKit's selection and WinUI's backplate begin (negative to go past it), and pillRadius rounds WinUI's selection indicator.
type collRowLook struct {
	macInset, fluentInset, pillRadius float32
}

// collPaintRow paints the platform's row in r: on macOS the selection of an inset NSTableView, corners of 8 in macListSelection's color (as macListRow draws it), square on the sides of chosen rows beside it; on Windows a ListViewItem's or TreeViewItem's backplate, 2 DIPs in from the row's top and bottom (s_backplateMargin, TreeViewItemPresenterMargin 4,2), corners of 4, SubtleFillColorSecondary while chosen or under the pointer and SubtleFillColorTertiary while both (ListViewItem_themeresources.xaml, TreeView_themeresources.xaml), with the 3 by 16 DIP pill of AccentFillColorDefault at the chosen row's start (s_selectionIndicatorSize; TreeViewItem.xaml's SelectionIndicator, HorizontalAlignment Left); on GNOME a square fill, $view_selected_color (the accent at 25%, 32% under the pointer) or $view_hover_color (currentColor at 4%) (_lists.scss row:selected and row.activatable:hover, _colors.scss of libadwaita 1.10.0). The pressed rows' looks are not drawn: the press is MyGo's row's, which a cover cannot ask without taking its clicks.
func collPaintRow(p *ui.Painter, pl Platform, t *ui.Theme, bg ui.Color, r ui.Rect, row collRow, look collRowLook) {
	switch pl {
	case Windows:
		secondary, tertiary := fluentSubtle(t)
		var fill ui.Color
		switch {
		case row.chosen && row.hovered:
			fill = tertiary
		case row.chosen || row.hovered:
			fill = secondary
		}
		plate := collInset(r, row, look.fluentInset)
		plate.Y, plate.H = plate.Y+2, plate.H-4
		if fill.A > 0 {
			collFillCorners(p, plate, fill, 4, collSides(row, false))
		}
		if row.chosen && row.start {
			h := min(plate.H, 16)
			x := plate.X
			if row.rtl {
				x = plate.X + plate.W - 3
			}
			p.Fill(ui.Rect{X: x, Y: plate.Y + (plate.H-h)/2, W: 3, H: h}, t.Accent, look.pillRadius)
		}
	case Linux:
		var fill ui.Color
		switch {
		case row.chosen && row.hovered:
			fill = bg.Mix(t.Accent, 0.32)
		case row.chosen:
			fill = bg.Mix(t.Accent, 0.25)
		case row.hovered:
			fill = adwInkAlpha(t, bg, 0.04)
		}
		if fill.A > 0 {
			p.Fill(r, fill, 0)
		}
	default:
		if row.chosen {
			collFillCorners(p, collInset(r, row, look.macInset), macListSelection(t, row.focused), 8, collSides(row, true))
		}
	}
}

// collInset returns r taken in by inset at the row's start and end, where the slice holds them.
func collInset(r ui.Rect, row collRow, inset float32) ui.Rect {
	left, right := row.start, row.end
	if row.rtl {
		left, right = right, left
	}
	if left {
		r.X, r.W = r.X+inset, r.W-inset
	}
	if right {
		r.W -= inset
	}
	return r
}

// collSides returns which sides of a row's slice round their corners, left, top, right and bottom: the row's start and end, and with block, the top and bottom unless a chosen row is beside them. A corner is rounded where both its sides are.
func collSides(row collRow, block bool) [4]bool {
	left, right := row.start, row.end
	if row.rtl {
		left, right = right, left
	}
	return [4]bool{left, !(block && row.above), right, !(block && row.below)}
}

// collFillCorners fills r with corners of radius where both their sides, left, top, right and bottom, are in round, and square elsewhere: the rounded fill is grown past r on the other sides and clipped to r, as MyGo's continuous corners curve further in than radius.
func collFillCorners(p *ui.Painter, r ui.Rect, c ui.Color, radius float32, round [4]bool) {
	if round == [4]bool{} {
		p.Fill(r, c, 0)
		return
	}
	g := 2 * radius
	grown := r
	if !round[0] {
		grown.X, grown.W = grown.X-g, grown.W+g
	}
	if !round[1] {
		grown.Y, grown.H = grown.Y-g, grown.H+g
	}
	if !round[2] {
		grown.W += g
	}
	if !round[3] {
		grown.H += g
	}
	if grown == r {
		p.Fill(r, c, radius)
		return
	}
	p.Clip(r, 0, func() { p.Fill(grown, c, radius) })
}

// collChevronColor returns the color of a row's disclosure chevron: secondaryLabelColor on macOS, white on a selection with the focus (inferred from AppKit's emphasized background style; the measured unemphasized chevron is secondaryLabelColor, 127 on 255 and 154 on 30); TreeViewItemForeground, TextFillColorPrimary, on Windows (TreeView_themeresources.xaml); currentColor, the row's text, on GNOME (_expanders.scss).
func collChevronColor(p Platform, t *ui.Theme, bg ui.Color, chosen, focused bool) ui.Color {
	switch {
	case p == Windows:
		return t.Text
	case p == Linux:
		return adwTextOn(t, bg)
	case chosen && focused:
		return white
	}
	return t.TextMuted
}

// collChevronTurn returns how far a row's chevron has turned, from 0 closed to 1 open: AppKit turns its disclosure chevron (over MyGo's 150 ms here, its duration inferred); WinUI swaps CollapsedGlyph for ExpandedGlyph and GTK the expander's pan-end icon for pan-down at once.
func collChevronTurn(p Platform, cover ui.Element, open bool) float32 {
	target := float32(0)
	if open {
		target = 1
	}
	if p != Mac {
		return target
	}
	return cover.Animate("breeze-disclosure", target, 150*time.Millisecond)
}

// collChevron draws the platform's disclosure chevron centered in box, MyGo's arrow box, whose clicks open and close the row, turned by turn from closed (0) to open (1):
//   - macOS: NSOutlineView's chevron, measured on an inset outline view on macOS 27.2: closed a '>' about 4.5 points wide and 6.5 high, open a 'v' about 8 wide and 4.5 high, stroked about 1.5, turning between them.
//   - Windows: TreeViewItem's CollapsedGlyph E76C and ExpandedGlyph E70D of Segoe Fluent Icons at GlyphSize 8 in a 12-DIP box (Generated/TreeViewItem.properties.cpp, TreeViewItem.xaml); the glyphs' outline is inferred as fluentChevronDown's, 7 DIPs across and 3.5 deep, stroked 1.
//   - GNOME: adwaita-icon-theme 51.0's pan-end-symbolic (pan-end-symbolic-rtl right to left) and pan-down-symbolic, 16 pixels, a 2-pixel round stroke along (7,4), (11,8), (7,12) and (4,7), (8,11), (12,7).
func collChevron(p *ui.Painter, pl Platform, box ui.Rect, turn float32, rtl bool, col ui.Color) {
	cx, cy := box.X+box.W/2, box.Y+box.H/2
	var path ui.Path
	switch pl {
	case Windows:
		if turn >= 0.5 {
			fluentChevronDown(p, cx, cy, col)
			return
		}
		s := float32(1)
		if rtl {
			s = -1
		}
		path.MoveTo(cx-s*1.75, cy-3.5).LineTo(cx+s*1.75, cy).LineTo(cx-s*1.75, cy+3.5)
		p.StrokePath(&path, 1, col)
	case Linux:
		x, y := box.X+(box.W-16)/2, box.Y+(box.H-16)/2
		switch {
		case turn >= 0.5:
			path.MoveTo(x+4, y+7).LineTo(x+8, y+11).LineTo(x+12, y+7)
		case rtl:
			path.MoveTo(x+9, y+4).LineTo(x+5, y+8).LineTo(x+9, y+12)
		default:
			path.MoveTo(x+7, y+4).LineTo(x+11, y+8).LineTo(x+7, y+12)
		}
		p.StrokePath(&path, 2, col)
	default:
		// Pointing right, 3 points deep, its arms 2.5 points from its axis closed and 3.25 open; turned clockwise, or right to left pointing left and turned the other way.
		a, b := float32(1.5), 2.5+0.75*turn
		s, deg := float32(1), 90*turn
		if rtl {
			s, deg = -1, -deg
		}
		sin, cos := math.Sincos(float64(deg) * math.Pi / 180)
		at := func(x, y float32) (float32, float32) {
			return cx + x*float32(cos) - y*float32(sin), cy + x*float32(sin) + y*float32(cos)
		}
		path.MoveTo(at(-s*a, -b)).LineTo(at(s*a, 0)).LineTo(at(-s*a, b))
		p.StrokePath(&path, 1.5, col)
	}
}

// ---- rows of Outline, Table and OutlineTable ----

// collRows is what the covers of a list's rows share as they are built and painted: the platform and its theme, the list's background and state, whether the list chooses rows and has the focus, the rows chosen and under the pointer, and the first and last column each row built, in the reading direction (a table builds a row's cells in the order it shows them).
type collRows struct {
	p          Platform
	t          *ui.Theme
	bg         ui.Color
	s          *ui.ListState
	choosable  bool
	focusKnown bool
	focused    bool
	chosen     map[int]bool
	hovered    map[int]bool
	first      map[int]int
	last       map[int]int
}

func collNewRows(c *ui.Context, p Platform, s *ui.ListState) *collRows {
	t := c.Theme()
	return &collRows{
		p: p, t: t, bg: collBg(p, t), s: s,
		choosable: s != nil && (s.Selected != nil || s.Selection != nil),
		chosen:    map[int]bool{}, hovered: map[int]bool{}, first: map[int]int{}, last: map[int]int{},
	}
}

// focus reports whether the list has the focus, asked as its first row builds, as ListState.FocusWithin answers while the list is built.
func (rs *collRows) focus(c *ui.Context) bool {
	if !rs.focusKnown {
		rs.focused, rs.focusKnown = rs.s != nil && rs.s.FocusWithin(c), true
	}
	return rs.focused
}

// header reports whether row i heads a section, which MyGo neither chooses nor lights up.
func (rs *collRows) header(i int) bool { return rs.s != nil && rs.s.Header != nil && rs.s.Header(i) }

// isChosen reports whether the cell of row i and col is chosen, and keeps the row's choice for the rows beside it.
func (rs *collRows) isChosen(i, col int) bool {
	on := rs.choosable && rs.s.CellSelected(i, col)
	if on {
		rs.chosen[i] = true
	}
	return on
}

// outlineRow builds the cover of row i of an Outline, depth levels deep, a branch open or not, and build, the app's content, inside it. MyGo's row (ui/outline.go coreOutline and prefix) is padded Space(1) and Space(2) with gaps of Space(1), and holds an indentation of depth*indent (with its gap) and a Space(4) arrow box before the cover, so the row's rect is the cover's grown by those at its start and by the padding elsewhere, which the cover's negative margins reach (and only there: over the arrow it would take its clicks).
func (rs *collRows) outlineRow(c *ui.Context, i, depth int, indent float32, branch, open bool, build func()) {
	p, t, bg := rs.p, rs.t, rs.bg
	lead := t.Space(2) + t.Space(4) + t.Space(1)
	if depth > 0 {
		lead += float32(depth)*indent + t.Space(1)
	}
	chosen := rs.isChosen(i, 0)
	focused := rs.focus(c)
	cover := ui.Row(c).Grow(1).AlignSelf(ui.Stretch).AlignItems(ui.Center).Role(ui.RoleNone).
		MarginY(-t.Space(1)).PaddingY(t.Space(1)).MarginEnd(-t.Space(2)).PaddingEnd(t.Space(2))
	hovered := rs.choosable && !rs.header(i) && cover.Hovered()
	cover.TextColor(collText(p, t, bg, chosen, focused))
	if p == Windows {
		// TreeViewItem: 32 DIPs from one item to the next, a backplate 28 high (TreeViewItemPresenterPadding 0,3,0,5 around TreeViewItemContentHeight 20, TreeViewItemPresenterMargin 4,2), and the content 40 DIPs into the backplate, past the chevron's column (14, the 12-DIP glyph and 14), and 16 more a level (TreeViewItem.xaml). The outline's padding puts MyGo's arrow box on the glyph's center; the content then needs 8 more at the top level and 4 more deeper, where MyGo's indentation adds a gap of its own.
		start := float32(8)
		if depth > 0 {
			start = 4
		}
		cover.MinHeight(32).PaddingStart(start)
	}
	turn := collChevronTurn(p, cover, open)
	rtl := c.Direction() == ui.RightToLeft
	arrow := t.Space(4)
	cover.Draw(func(pt *ui.Painter, r ui.Rect) {
		row := ui.Rect{X: r.X - lead, Y: r.Y, W: r.W + lead, H: r.H}
		if rtl {
			row.X = r.X
		}
		pt.Fill(row, bg, 0)
		collPaintRow(pt, p, t, bg, row, collRow{chosen: chosen, hovered: hovered, focused: focused, above: rs.chosen[i-1], below: rs.chosen[i+1], start: true, end: true, rtl: rtl}, collRowLook{macInset: 8, fluentInset: -4, pillRadius: 2})
		if branch {
			box := ui.Rect{X: r.X - t.Space(1) - arrow, Y: r.Y + (r.H-arrow)/2, W: arrow, H: arrow}
			if rtl {
				box.X = r.X + r.W + t.Space(1)
			}
			collChevron(pt, p, box, turn, rtl, collChevronColor(p, t, bg, chosen, focused))
		}
	})
	cover.Children(build)
}

// cell builds the cover of the cell of row i and col of a Table or an OutlineTable, aligned as align, and build, the app's content, inside it. MyGo's cell box (ui/table.go table: Padding(Space(1.5), Space(2.5)), Clip()) clips the cover to it, so the cover reaches the box's edges with negative margins and paints its slice of the row: the selection's rounded start in the row's first cell, its end in the last. In an OutlineTable's first column lead is the width of the indentation and the arrow box before the cover, which MyGo builds there without a gap; the cover does not reach over them, which would take the arrow's clicks, and paints them from its own rect. A row heading a section keeps MyGo's look.
func (rs *collRows) cell(c *ui.Context, i, col int, align ui.Align, lead float32, branch, open bool, build func()) {
	if rs.header(i) {
		build()
		return
	}
	if _, ok := rs.first[i]; !ok {
		rs.first[i] = col
	}
	rs.last[i] = col
	p, t, bg := rs.p, rs.t, rs.bg
	chosen := rs.isChosen(i, col)
	focused := rs.focus(c)
	padY, padX := t.Space(1.5), t.Space(2.5)
	inY := padY
	if p == Linux {
		// GtkColumnView's cells: padding 8px 6px (_column-view.scss). The 6 pixels at their sides stay MyGo's 10, which its header's titles keep as well.
		inY = 8
	}
	cover := ui.Row(c).Grow(1).AlignSelf(ui.Stretch).AlignItems(ui.Center).Role(ui.RoleNone).
		MarginY(-padY).PaddingY(inY).MarginEnd(-padX).PaddingEnd(padX)
	if lead == 0 {
		cover.MarginStart(-padX).PaddingStart(padX)
	}
	switch align {
	case ui.Center:
		cover.Justify(ui.Center)
	case ui.End:
		cover.Justify(ui.End)
	}
	if p == Windows {
		// ListViewItemMinHeight 40: WinUI 3 has no table, so its rows are ListView's (inferred).
		cover.MinHeight(40)
	}
	if rs.choosable && cover.Hovered() {
		rs.hovered[i] = true
	}
	cover.TextColor(collText(p, t, bg, chosen, focused))
	turn := collChevronTurn(p, cover, open)
	rtl := c.Direction() == ui.RightToLeft
	arrow := t.Space(4)
	cover.Draw(func(pt *ui.Painter, r ui.Rect) {
		box := r
		if lead > 0 {
			box.W += lead + padX
			if !rtl {
				box.X -= lead + padX
			}
		}
		pt.Fill(box, bg, 0)
		row := collRow{chosen: chosen, hovered: rs.hovered[i], focused: focused, above: rs.chosen[i-1], below: rs.chosen[i+1], start: rs.first[i] == col, end: rs.last[i] == col, rtl: rtl}
		collPaintRow(pt, p, t, bg, box, row, collRowLook{macInset: 4, pillRadius: 1.5})
		if branch {
			ab := ui.Rect{X: r.X - arrow, Y: r.Y + (r.H-arrow)/2, W: arrow, H: arrow}
			if rtl {
				ab.X = r.X + r.W
			}
			collChevron(pt, p, ab, turn, rtl, collChevronColor(p, t, bg, chosen, focused))
		}
	})
	cover.Children(build)
}

// collOutlineIndex returns the row of an outline's item, from the rows MyGo flattened for the frame before it builds them (OutlineState.Rows and Item), indexed as the first row builds.
func collOutlineIndex[K comparable](s *ui.OutlineState[K]) func(K) int {
	var at map[K]int
	return func(item K) int {
		if at == nil {
			at = make(map[K]int, s.Rows())
			for i := range s.Rows() {
				at[s.Item(i)] = i
			}
		}
		return at[item]
	}
}

// collStyleOutline gives an Outline the platform's background and the room around its rows: on macOS an inset outline view's 10 points above and below, and 2 at the sides, so that MyGo's arrow box, Space(2) into the row, centers 8 points into the selection drawn 10 points in from the edges (AppKit's chevron centers about 8.5 points into it, measured; its content starts 15 in, MyGo's 20); on Windows 8 at the sides, so that the arrow box centers on TreeViewItem's glyph, 20 DIPs into a backplate 4 in from the edges; on GNOME none, as a GtkListView's rows are its full width (MyGo's arrow then centers 16 pixels in, GTK's expander 10, row padding 2px and 16-pixel expander).
func collStyleOutline(p Platform, o ui.Element, bg ui.Color) {
	o.Background(bg)
	switch p {
	case Windows:
		o.PaddingX(8)
	case Mac:
		o.Padding(10, 2)
	}
}

// collStyleTable gives a Table or an OutlineTable the platform's background and room at its sides: 6 points on macOS, so that the first column's text, 10 into its cell, starts 16 points in, as an inset NSTableView's does (its selection 10 in, measured), the header's titles with it; 4 DIPs on Windows, so that the backplate and the pill start 4 in, as ListViewItem's do (s_backplateMargin, s_selectionIndicatorMargin).
func collStyleTable(p Platform, tb ui.Element, bg ui.Color) {
	tb.Background(bg)
	switch p {
	case Mac:
		tb.PaddingX(6)
	case Windows:
		tb.PaddingX(4)
	}
}

// ---- GridView ----

// collGridChosen reports whether item i of a grid is chosen. GridState's Selector keeps its methods unexported, so a grid choosing several items is read through the exported Has of the Selection it holds when its keys are the items' indexes or strings; with keys of another type only the item last chosen, *Selected, is drawn chosen (MyGo needs a GridState.Chosen to tell them all without reflection).
func collGridChosen(s *ui.GridState, i int) bool {
	if s == nil || s.Selected == nil {
		return false
	}
	if s.Selection == nil {
		return *s.Selected == i
	}
	if s.Key == nil {
		if sel, ok := s.Selection.(*ui.Selection[int]); ok {
			return sel.Has(i)
		}
	} else if k, ok := s.Key(i).(string); ok {
		if sel, ok := s.Selection.(*ui.Selection[string]); ok {
			return sel.Has(k)
		}
	}
	return *s.Selected == i
}

// collGridCell gives the cover of a grid's item the platform's item. MyGo's cell (ui/gridview.go GridState.cell) is a column with a 2-DIP border, clear unless chosen, filled with the accent at 20% inside an accent border while chosen in a grid with the focus, SurfacePressed inside a Border-colored one while chosen otherwise, and SurfaceHover under the pointer; the cover reaches over its border with negative margins and paints over all of them:
//   - macOS: NSCollectionView draws no selection of its own, so its items are drawn as an inset table's rows are (inferred): macListSelection's color with corners of 8 while chosen, white text while the grid has the focus, nothing under the pointer.
//   - Windows: GridViewItem (GridViewItem_themeresources.xaml, ListViewBaseItemChrome.cpp): its content centered (HorizontalContentAlignment and VerticalContentAlignment Center), corners of 4 (GridViewItemCornerRadius); under the pointer SubtleFillColorSecondary inside a 1-DIP border of ControlStrokeColorOnAccentTertiary (#37000000) and the text in TextFillColorSecondary (GridViewItemForegroundPointerOver); chosen, SubtleFillColorTertiary 1 DIP in (SetBackplateMargin, as the border is opaque) inside a 2-DIP border of AccentFillColorDefault (GridViewItemSelectedBorderThickness, AccentFillColorSecondary under the pointer) and a 1-DIP inner border of ControlSolidFillColorDefault (#ffffff, #454545) with corners of 3 (s_innerBorderCornerRadius).
//   - GNOME: gridview > child (_views.scss): padded 3px with corners of 9 ($button_radius), $view_selected_color while chosen ($view_selected_hover_color under the pointer) and $view_hover_color under the pointer, as its items are activatable.
//
// MyGo's keyboard cursor ring and the line where dragged items would drop are drawn over the cover, and its gap of Space(2) between items stays (WinUI's items are 4 apart, GridViewItem's Margin 0,0,4,4).
func collGridCell(p Platform, t *ui.Theme, cover ui.Element, bg ui.Color, chosen, hovered, focused bool) {
	text := collText(p, t, bg, chosen, focused)
	switch p {
	case Windows:
		cover.Padding(2).Justify(ui.Center).AlignItems(ui.Center)
		if hovered && !chosen {
			text = t.TextMuted
		}
	case Linux:
		cover.Padding(3)
	default:
		cover.Padding(2)
	}
	cover.TextColor(text)
	cover.Draw(func(pt *ui.Painter, r ui.Rect) {
		pt.Fill(r, bg, 0)
		switch p {
		case Windows:
			secondary, tertiary := fluentSubtle(t)
			switch {
			case chosen:
				border := t.Accent
				if hovered {
					border = t.AccentHover
				}
				pt.Fill(ui.Rect{X: r.X + 1, Y: r.Y + 1, W: r.W - 2, H: r.H - 2}, tertiary, 3)
				pt.Stroke(r, border, 4, 2)
				pt.Stroke(ui.Rect{X: r.X + 2, Y: r.Y + 2, W: r.W - 4, H: r.H - 4}, pick(t.Dark, white, ui.Hex("#454545")), 3, 1)
			case hovered:
				pt.Fill(r, secondary, 4)
				pt.Stroke(r, bg.Mix(black, 0x37/255.0), 4, 1)
			}
		case Linux:
			switch {
			case chosen && hovered:
				pt.Fill(r, bg.Mix(t.Accent, 0.32), 9)
			case chosen:
				pt.Fill(r, bg.Mix(t.Accent, 0.25), 9)
			case hovered:
				pt.Fill(r, adwInkAlpha(t, bg, 0.04), 9)
			}
		default:
			if chosen {
				pt.Fill(r, macListSelection(t, focused), 8)
			}
		}
	})
}

// ---- Tree and TreeItem ----

// collTreeDepths holds how deep the TreeItem being built in each window is, which MyGo keeps to itself.
var collTreeDepths sync.Map // *ui.Context -> *int

// collTreeDepth returns how deep the TreeItem built now in c's window is: 0 at the top of its Tree.
func collTreeDepth(c *ui.Context) int {
	if v, ok := collTreeDepths.Load(c); ok {
		return *v.(*int)
	}
	return 0
}

// collTreeChildren returns children building one level deeper, or nil for nil, which MyGo takes for a leaf.
func collTreeChildren(c *ui.Context, children func()) func() {
	if children == nil {
		return nil
	}
	return func() {
		d := collTreeDepth(c)
		deeper := d + 1
		collTreeDepths.Store(c, &deeper)
		children()
		if d == 0 {
			collTreeDepths.Delete(c)
		} else {
			collTreeDepths.Store(c, &d)
		}
	}
}

// collTreeTheme sets th, the theme a Tree's items are built in, to the platform's rows. MyGo's TreeItem (ui/tree.go) fills the item its caller marks Selected with the Accent of the theme around it and gives it AccentText, fills an item under the pointer with SurfaceHover, draws its arrow in TextMuted (AccentText while chosen) and its focus ring with the theme's Radius:
//   - macOS: macListSelection's color while the tree has the focus, unemphasized gray otherwise, white text on the first and the text's own on the second, nothing under the pointer, and the selection's corners of 8.
//   - Windows: TreeViewItemBackgroundSelected and TreeViewItemBackgroundPointerOver, SubtleFillColorSecondary, TreeViewItemForeground, TextFillColorPrimary, for the text and the glyph whatever the state, and ControlCornerRadius 4 (TreeView_themeresources.xaml, TreeViewItem.xaml).
//   - GNOME: $view_selected_color, the accent at 25% over view_bg_color, $view_hover_color, view_fg_color for the text and the expander (currentColor), and square rows (_lists.scss).
func collTreeTheme(p Platform, t, th *ui.Theme, focused bool) {
	switch p {
	case Windows:
		hover, _ := fluentSubtle(t)
		th.Accent, th.AccentText, th.SurfaceHover, th.TextMuted, th.Radius = hover, t.Text, hover, t.Text, 4
	case Linux:
		bg := adwListViewBg(t)
		fill := bg.Mix(t.Accent, 0.25)
		th.Text = adwTextOn(t, bg)
		th.Accent, th.AccentText, th.SurfaceHover, th.TextMuted, th.Radius = fill, adwTextOn(t, fill), adwInkAlpha(t, bg, 0.04), th.Text, 0
	default:
		th.Accent, th.AccentText, th.SurfaceHover, th.Radius = macListSelection(t, focused), collText(p, t, t.Background, true, focused), ui.Transparent, 8
	}
}

// collTreeItem gives an item of a Tree, depth levels deep, the platform's geometry. MyGo's item is Space(7) high, its Space(4) arrow box Space(1+4*depth) from its start and Space(0.5) from its label:
//   - macOS: a row of an inset NSOutlineView, measured on macOS 27.2: 24 points high, its selection 10 points in from the tree's sides with corners of 8, the chevron centered 8.5 points into it and 13 more a level (indentationPerLevel), the label 15 in; MyGo's arrow box centers 8 points in and the label lands 16 in.
//   - Windows: TreeViewItem: 32 DIPs from one item to the next around a backplate 28 high, 4 in from the sides (TreeViewItemPresenterMargin 4,2, TreeViewItemPresenterPadding 0,3,0,5, TreeViewItemContentHeight 20), the glyph centered 20 DIPs into the backplate and the content 40 in, 16 more a level (ExpandCollapseChevron's Padding 14,0 around the 12-DIP glyph; TreeViewItem.cpp's depth * 16).
//   - GNOME: a GtkListView row holding a GtkTreeExpander: padding 2px (_lists.scss listview > row), the 16-pixel expander and the label 4 apart (treeexpander's border-spacing), 20 pixels more a level (collIndent), as high as the label and the padding (inferred: no min-height is set).
//
// The 3 by 16 accent pill WinUI shows at the start of the chosen item is not drawn: the caller marks the item Selected after it returns, and Element has nothing that tells a callback so as it paints. MyGo's arrow stays, a 4 by 8 chevron stroked 1.5, as it is a child of the item whose fill a cover could not know without the choice.
func collTreeItem(c *ui.Context, p Platform, item ui.Element, depth int, label string) {
	d := float32(depth)
	switch p {
	case Windows:
		item.Height(28).MarginX(4).MarginY(2).Radius(4).PaddingStart(12 + 16*d).Gap(12)
	case Linux:
		t := c.Theme()
		if label == "" {
			label = " "
		}
		_, h := c.MeasureText(0, ui.Span{Text: label, Size: t.FontSize})
		item.Height(max(h, 16) + 4).Radius(0).PaddingStart(2 + 20*d).Gap(4)
	default:
		item.Height(24).MarginX(10).Radius(8).PaddingStart(13 * d).Gap(0)
	}
}

// ---- EditableText ----

// collEditFrame draws the platform's frame around the input MyGo builds inside an EditableText while it is edited (ui/editable.go coreEditableText): the input's box is the text's grown by Space(0.5)+1 above and below and Space(1)+1 at the sides (its negative margins), with corners of Space(1) and a 1-DIP border of the accent, which a stroke of the input's own fill, the theme's Background, hides. Then, as the platform's field editor is framed (inferred: AppKit's editing a cell and GTK's editable label were not rendered with the focus):
//   - macOS: macFocusRing around it, as AppKit rings a field being edited.
//   - Windows: a focused TextBox's 1-DIP ControlStrokeColorDefault edge and its 2-DIP accent line along the bottom, as fluentField draws them.
//   - GNOME: nothing, as editablelabel > stack > text is only filled with view_bg_color (_labels.scss); the 2-pixel accent outline libadwaita draws around one in a column view's row that is not chosen (_column-view.scss) is left out, since a list edits the text of its chosen row.
//
// The input's fill stays the theme's Background and its caret the accent, MyGo's.
func collEditFrame(p Platform, t *ui.Theme, b ui.Element) {
	my, mx, radius := t.Space(0.5)+1, t.Space(1)+1, t.Space(1)
	b.DrawOver(func(pt *ui.Painter, r ui.Rect) {
		in := ui.Rect{X: r.X - mx, Y: r.Y - my, W: r.W + 2*mx, H: r.H + 2*my}
		pt.Stroke(in, t.Background, radius, 1)
		switch p {
		case Windows:
			pt.Stroke(in, t.Border, radius, 1)
			pt.Clip(in, radius, func() {
				pt.Fill(ui.Rect{X: in.X, Y: in.Y + in.H - 2, W: in.W, H: 2}, t.Accent, 0)
			})
		case Mac:
			macFocusRing(pt, in, radius)
		}
	})
}
