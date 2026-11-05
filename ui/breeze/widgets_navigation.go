package breeze

import (
	"sync"

	"github.com/egoist/mygo/ui"
)

// ---- sidebars: AppKit's source list, WinUI 3's NavigationView pane, libadwaita's .navigation-sidebar ----

// navSidebars holds the Sidebar each window is building, which its sections and items read.
var navSidebars sync.Map // *ui.Context -> *navSidebar

// navSidebar is a Sidebar being built: the ID of the item chosen, whether the sidebar has the keyboard focus in an active window, as AppKit draws the selection then in the accent color, and how many sections and items it holds so far.
type navSidebar struct {
	selected *string
	focused  bool
	built    int
}

// navCurrentSidebar returns the Sidebar being built in c's window, or nil outside one.
func navCurrentSidebar(c *ui.Context) *navSidebar {
	if v, ok := navSidebars.Load(c); ok {
		return v.(*navSidebar)
	}
	return nil
}

// Sidebar creates a sidebar of the sections and items fn builds (SidebarSection, SidebarItem), as the platform draws its source list (macOS), NavigationView's pane (Windows) or navigation sidebar (GNOME): a click chooses an item, whose ID *selected holds, and while the sidebar has the focus Up, Down, Home, End and the first letters of an item choose among them, as in ui.Sidebar. Give it a width. It is ui.Sidebar's counterpart, and its items are drawn as the platform draws them only when built with SidebarItem.
func Sidebar(c *ui.Context, selected *string, fn func()) ui.Element {
	p := PlatformOf(c)
	var focus *ui.Handle
	sb := ui.Sidebar(c, selected, func() {
		body := ui.Column(c)
		// The sidebar is bound to a handle kept by its first child, so that its items know while they build whether it has the focus.
		focus = ui.Local(body, "breeze-sidebar", func() ui.Handle { return ui.Handle{} })
		st := &navSidebar{selected: selected, focused: focus.Focused(c)}
		outer, nested := navSidebars.Load(c)
		navSidebars.Store(c, st)
		switch p {
		case Windows:
			body.Gap(4) // NavigationViewItemButtonMargin 4,2 between two items
		case Linux:
			body.Gap(2) // a row's margin: 0 6px 2px
		}
		body.Children(fn)
		if nested {
			navSidebars.Store(c, outer)
		} else {
			navSidebars.Delete(c)
		}
	})
	sb.Bind(focus)
	t := c.Theme()
	switch p {
	case Windows:
		// NavigationViewExpandedPaneBackground is SolidBackgroundFillColorTransparent: the window shows through the pane.
		sb.Background(ui.Transparent).Padding(4)
	case Linux:
		// .sidebar-pane: sidebar_bg_color with a 1-pixel sidebar_border_color line along its end; the list is padded 6 pixels ($menu_margin).
		sb.Background(adwSidebarBg(t)).Padding(6).BorderEndWidth(1).BorderColor(adwSidebarBorder(t))
	default:
		// The selection is 10 points in from the source list's edges, measured.
		sb.Background(macSidebarBg(t)).Padding(10)
	}
	return sb
}

// SidebarSection creates a section of a Sidebar under title, holding the items fn builds while *open is true; a click on the title shows and hides them, as in ui.SidebarSection. The title is drawn as the platform heads a group of its sidebar: AppKit's group row, NavigationViewItemHeader or a dimmed heading on GNOME. A nil open keeps the section open. It is ui.SidebarSection's counterpart.
func SidebarSection(c *ui.Context, title string, open *bool, fn func()) ui.Element {
	p := PlatformOf(c)
	t := c.Theme()
	// MyGo draws the title FontSize-1 large in TextMuted, padded Space(1) and Space(2), so a theme of those around the section draws the platform's heading; fn's items build with the theme around the sidebar again.
	head := *t
	// first is the space above the first section, top the space above another, below the space between the title and the items.
	var first, top, below float32
	switch p {
	case Windows:
		// NavigationViewItemHeader: a 40-pixel row, its text (NavigationViewItemHeaderTextStyle: 14 pixels, semibold, TextFillColorSecondary) centered, 16 from the pane's edge, between items whose NavigationViewItemButtonMargin keeps 2 above and below. MyGo's title is padded Space(1) above and below: 36 high with its disclosure arrow (Space(4) square), 7.5 less without one, around the 16.5-pixel line of its text (measured in the tester), so the space it lacks to 40 goes half above and half below it.
		head.FontSize, head.Spacing, head.Radius = 15, 6, 4
		short := float32(0)
		if open == nil {
			short = 7.5
		}
		first, top, below = 2+short/2, short/2, 1+short/2
	case Linux:
		// libadwaita 1.7 has no sidebar headings: .heading .dim-label, the text at 55% (inferred), lined up with the rows' labels, 12 pixels below the row above (inferred).
		head.FontSize, head.Spacing, head.Radius = t.FontSize+1, 4, 9
		head.TextMuted = adwInkAlpha(t, adwSidebarBg(t), 0.55)
		top = 12
	default:
		// AppKit's group row, measured on macOS 27: 11-point semibold text in a row 19 points high, 2 points inside the selection's inset, the first 13 points from the top (an item's row would be 10) and the others 13 points below the row above, in tertiaryLabelColor (black at 26% in Aqua, white at 25% in Dark Aqua; rendered offscreen, without the sidebar's vibrancy, AppKit draws it lighter still, 20% and 12%). MyGo's title row is 18 points high with the space below it, so the margins above put its text where AppKit's is, the items following half a point higher.
		head.FontSize, head.Spacing, head.Radius = 12, 2, 8
		head.TextMuted = tint(macSidebarBg(t), t.Dark, pick(t.Dark, float32(0.259), 0.247))
		first, top = 3.5, 14
	}
	st := navCurrentSidebar(c)
	wrap := ui.Themed(c, &head).Background(ui.Transparent)
	if st != nil {
		if st.built > 0 {
			wrap.Margin(top, 0, 0, 0)
		} else {
			wrap.Margin(first, 0, 0, 0)
		}
		st.built++
	}
	var sec ui.Element
	wrap.Children(func() {
		sec = ui.SidebarSection(c, title, open, func() {
			gap := float32(0)
			switch p {
			case Windows:
				gap = 4
			case Linux:
				gap = 2
			}
			ui.Themed(c, t).Background(ui.Transparent).Gap(gap).Padding(below, 0, 0, 0).Children(fn)
		})
	})
	return sec
}

// SidebarItem creates an item of a Sidebar choosing id, showing label after icon, nil for none, as the platform draws a row of its sidebar: on macOS 32 points high with a rounded selection, in the accent color while the sidebar has the focus and in gray otherwise, its label semibold; on Windows a NavigationViewItem, 36 high with a 3 by 16 accent pill at the start of the chosen one; on GNOME a row 36 high, the chosen one darker. Add to it with Children, as a Badge. Build it in this package's Sidebar, which tells it whether it is chosen and whether the sidebar has the focus. It is ui.SidebarItem's counterpart.
func SidebarItem(c *ui.Context, id string, icon *ui.SVG, label string) ui.Element {
	p := PlatformOf(c)
	t := c.Theme()
	st := navCurrentSidebar(c)
	chosen, focused := false, false
	if st != nil {
		chosen, focused = *st.selected == id, st.focused
		st.built++
	}
	// MyGo draws an item's icon in the theme's Accent, its text while chosen in the sidebar with the focus in AccentText, and the item under the pointer in SurfaceHover; a theme of the platform's colors around it draws them so.
	row := *t
	switch p {
	case Windows:
		hover, _ := fluentSubtle(t)
		row.Accent, row.AccentText, row.SurfaceHover = t.Text, t.Text, hover
	case Linux:
		bg := adwSidebarBg(t)
		row.Text = adwSidebarFg(t)
		row.Accent, row.AccentText, row.SurfaceHover = row.Text, row.Text, adwInkAlpha(t, bg, 0.07)
	default:
		// AppKit's source list draws a template image in the text color and does not mark the row under the pointer.
		row.Accent, row.AccentText, row.SurfaceHover = t.Text, white, ui.Transparent
	}
	var item ui.Element
	c.Key(id)
	ui.Themed(c, &row).Background(ui.Transparent).Children(func() { item = ui.SidebarItem(c, id, icon, label) })
	switch p {
	case Windows:
		if !chosen && item.Pressed() {
			// NavigationViewItemBackgroundPressed is SubtleFillColorTertiary. MyGo paints an item that is not chosen in the SurfaceHover of the theme around it as it paints, so the theme takes that color now.
			_, row.SurfaceHover = fluentSubtle(t)
		}
		fluentSidebarItem(c, item, chosen, icon != nil)
	case Linux:
		if !chosen && item.Pressed() {
			// $active_color, currentColor at 16%, as Windows takes its pressed fill above.
			row.SurfaceHover = adwInkAlpha(t, adwSidebarBg(t), 0.16)
		}
		adwSidebarItem(c, item, chosen)
	default:
		macSidebarItem(c, item, chosen, focused, icon != nil)
	}
	return item
}

// macSidebarBg is a source list's background, measured from an NSOutlineView of style sourceList rendered by AppKit on macOS 27.
func macSidebarBg(t *ui.Theme) ui.Color { return pick(t.Dark, ui.Hex("#f0f0f0"), ui.Hex("#282828")) }

// macSidebarItem draws a row of an NSOutlineView of style sourceList at the default (medium) row size, measured on macOS 27: 32 points high, the selection 10 points in from the list's edges with corners of 8, filling the row. A cell's image view is 23 points wide, 6 inside the selection, so the 16-point icon is centered in it, 9.5 points in, and the text field 3 points after it, its text 2 points further in (the text ink starts at 46, measured); without an icon the text starts 6 points inside the selection. AppKit draws the chosen row's label semibold (.SFNS-Semibold in its attributed string, measured); while the outline view has the focus the selection is the emphasized one in the accent color with white text, otherwise unemphasizedSelectedContentBackgroundColor (#dcdcdc, #464646), AppKit's semantic color (offscreen, AppKit draws that selection without its material, in black).
func macSidebarItem(c *ui.Context, item ui.Element, chosen, focused, hasIcon bool) {
	t := c.Theme()
	item.MinHeight(32).Padding(0, 6).Radius(8)
	if hasIcon {
		item.PaddingStart(9.5).Gap(9.5)
	}
	if !chosen {
		return
	}
	item.FontWeight(600)
	if focused {
		item.Background(macSidebarSelection(t)).TextColor(white)
	} else {
		item.Background(pick(t.Dark, ui.Hex("#dcdcdc"), ui.Hex("#464646"))).TextColor(t.Text)
	}
}

// macSidebarSelection is the emphasized selection of a source list: #0070f5 in Aqua and #007aff, controlAccentColor, in Dark Aqua for the blue accent, measured from a row view drawn emphasized and again from a live capture of a key window's sidebar on macOS 27.2 (#0070f5 and #007aff, converted from the display's profile to sRGB). For another accent it is the accent, a little darker in Aqua (inferred).
func macSidebarSelection(t *ui.Theme) ui.Color {
	if t.Accent == macAccent {
		return pick(t.Dark, ui.Hex("#0070f5"), ui.Hex("#007aff"))
	}
	return pick(t.Dark, t.Accent.Mix(black, 0.06), t.Accent)
}

// fluentSidebarItem draws a NavigationViewItem in the left pane: NavigationViewItemOnLeftMinHeight 36, ControlCornerRadius 4, the 16-DIP icon centered in a 40-DIP column and the text 4 after it (8 and 4 without an icon), 14 before the end; SubtleFillColorSecondary under the pointer and while chosen (Tertiary chosen under the pointer), TextFillColorSecondary pressed, and the chosen one's selection indicator, AccentFillColorDefault, 3 by 16 with corners of 2, at its start.
func fluentSidebarItem(c *ui.Context, item ui.Element, chosen, hasIcon bool) {
	t := c.Theme()
	secondary, tertiary := fluentSubtle(t)
	item.MinHeight(36).Radius(4).PaddingY(0).PaddingStart(12).PaddingEnd(14)
	if hasIcon {
		item.Gap(16)
	} else {
		item.Gap(8)
	}
	fg := t.Text
	if item.Pressed() {
		fg = t.TextMuted
	}
	item.TextColor(fg)
	if !chosen {
		return
	}
	if item.Hovered() && !item.Pressed() {
		item.Background(tertiary)
	} else {
		item.Background(secondary)
	}
	item.Draw(func(p *ui.Painter, r ui.Rect) {
		x := r.X
		if c.Direction() == ui.RightToLeft {
			x = r.X + r.W - 3
		}
		p.Fill(ui.Rect{X: x, Y: r.Y + (r.H-16)/2, W: 3, H: 16}, t.Accent, 1.5)
	})
}

// adwSidebarBg is sidebar_bg_color.
func adwSidebarBg(t *ui.Theme) ui.Color { return pick(t.Dark, ui.Hex("#ebebed"), ui.Hex("#2e2e32")) }

// adwSidebarFg is sidebar_fg_color over the sidebar: RGB(0 0 6 / 80%) in the light theme, white in the dark one.
func adwSidebarFg(t *ui.Theme) ui.Color {
	if t.Dark {
		return white
	}
	return adwSidebarBg(t).Mix(adwInk, 0.8)
}

// adwSidebarBorder is sidebar_border_color: RGB(0 0 6 / 7%) in the light theme, 36% in the dark one.
func adwSidebarBorder(t *ui.Theme) ui.Color {
	return adwSidebarBg(t).Mix(adwInk, pick(t.Dark, float32(0.07), 0.36))
}

// adwSidebarItem draws a row of a .navigation-sidebar list: at least 36 pixels high, padding 0 8px, corners of $menu_radius (9), currentColor at 7% under the pointer, and while selected at 10%, 13% under the pointer and 19% pressed. The 12 pixels between a row's icon and its label are inferred, as libadwaita leaves a row's content to the app.
func adwSidebarItem(c *ui.Context, item ui.Element, chosen bool) {
	t := c.Theme()
	bg := adwSidebarBg(t)
	item.MinHeight(36).Padding(0, 8).Gap(12).Radius(9)
	if !chosen {
		return
	}
	share := float32(0.10)
	switch {
	case item.Pressed():
		share = 0.19
	case item.Hovered():
		share = 0.13
	}
	item.Background(adwInkAlpha(t, bg, share))
}

// ---- breadcrumbs: NSPathControl, WinUI 3's BreadcrumbBar, Nautilus's path bar ----

// Breadcrumbs creates a path of items, as the platform's path control (macOS), BreadcrumbBar (Windows) or Files' path bar (GNOME) draws it, separated by chevrons, or slashes on GNOME: the items but the last are links, which a click chooses, as Enter does, setting *chosen to its index, which Changed reports. Long items shorten with an ellipsis as room runs short. Assistive technology sees links in a group; name it with Label. It is ui.Breadcrumbs's counterpart, built as it is of a row of buttons without a look (ui.ButtonBase) that assistive technology sees as links.
func Breadcrumbs(c *ui.Context, items []string, chosen *int) ui.Element {
	p := PlatformOf(c)
	t := c.Theme()
	row := ui.Row(c).AlignItems(ui.Center).MinWidth(0).Role(ui.RoleGroup)
	if p == Linux {
		// #NautilusPathBar: currentColor at 10%, corners of 9.
		row.Background(adwInkAlpha(t, t.Background, 0.10)).Radius(9)
	} else if p != Windows {
		// An NSPathControl of style standard is 22 points high, its first item's text 5 points in, measured.
		row.MinHeight(22).PaddingX(2.5)
	}
	row.Children(func() {
		for i, item := range items {
			if i > 0 {
				navCrumbSeparator(c, p)
			}
			if i == len(items)-1 {
				// Where the path is: no link.
				cur := ui.Row(c).Key(i).AlignItems(ui.Center).MinWidth(24).Shrink(0.5)
				navCrumb(c, p, cur, false)
				cur.Children(func() { ui.Text(c, item).SingleLine() })
				continue
			}
			c.Key(i)
			b := ui.ButtonBase(c).Role(ui.RoleLink).MinWidth(24).Shrink(1)
			if b.Clicked() {
				*chosen = i
				row.SetChanged()
			}
			navCrumb(c, p, b, true)
			b.Children(func() { ui.Text(c, item).SingleLine() })
		}
	})
	return row
}

// navCrumb draws an item of a path: a link, or the last item when link is false.
func navCrumb(c *ui.Context, p Platform, e ui.Element, link bool) {
	t := c.Theme()
	switch p {
	case Windows:
		// A BreadcrumbBarItem: ControlContentThemeFontSize, padding 1,3 around a 20-DIP line, TextFillColorPrimary, Secondary under the pointer and Tertiary pressed; the current item is TextFillColorPrimary too.
		fg := t.Text
		if link {
			switch {
			case e.Pressed():
				fg = t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.447), 0.529))
			case e.Hovered():
				fg = t.TextMuted
			}
		}
		e.MinHeight(26).Padding(3, 1).Radius(4).TextColor(fg)
	case Linux:
		// #NautilusPathButton: a button 24 pixels high without vertical padding, margin 3, corners of 7, its label bold; the others than the current one dimmed, without a background but under the pointer, where they take a button's hover color (currentColor at 15%) and their full color, 30% pressed.
		trough := adwInkAlpha(t, t.Background, 0.10)
		e.Margin(3).MinHeight(24).PaddingX(10).Radius(7).FontWeight(700).TextColor(t.Text)
		if !link {
			return
		}
		switch {
		case e.Pressed():
			e.Background(adwInkAlpha(t, trough, 0.30))
		case e.Hovered():
			e.Background(adwInkAlpha(t, trough, 0.15))
		default:
			e.TextColor(t.TextMuted)
		}
	default:
		// NSPathControl draws every item in secondaryLabelColor, without a background; pressed, a link darkens to labelColor (inferred).
		fg := t.TextMuted
		if link && e.Pressed() {
			fg = t.Text
		}
		e.PaddingX(2).Radius(4).TextColor(fg)
	}
}

// navCrumbSeparator draws what separates two items of a path, pointing along it: left in a right-to-left window.
func navCrumbSeparator(c *ui.Context, p Platform) {
	t := c.Theme()
	if p == Linux {
		// G_DIR_SEPARATOR_S in a .dim-label, in the weight of the text around it.
		ui.Text(c, "/").Role(ui.RoleNone).FontWeight(400).TextColor(t.TextMuted)
		return
	}
	rtl := c.Direction() == ui.RightToLeft
	sep := ui.Box(c).Shrink(0).Role(ui.RoleNone)
	color, width := t.TextMuted, float32(1.2)
	if p == Windows {
		// The chevron of Segoe Fluent Icons, BreadcrumbBarChevronLeftToRight U+E974 (U+E973 right to left, mirrored here), BreadcrumbBarChevronFontSize 12 with BreadcrumbBarChevronPadding 2,0, in BreadcrumbBarNormalForegroundBrush (TextFillColorPrimary). The glyph's shape within its 12-pixel em, about 3.5 by 7 pixels with a 1-pixel stroke, is inferred: the font is not open.
		sep.Size(12, 12).MarginX(2)
		color, width = t.Text, 1
	} else {
		// NSPathControl's chevron: 4 by 6 points of ink in secondaryLabelColor, 7.5 points after an item's text and 5.5 before the next, measured; the stroke, about 1.2 points, keeps inside that box.
		sep.Size(4, 6).MarginStart(5).MarginEnd(3)
	}
	sep.Draw(func(pt *ui.Painter, r ui.Rect) {
		in := width / 2
		x0, x1, y0, y1 := r.X+in, r.X+r.W-in, r.Y+in, r.Y+r.H-in
		if p == Windows {
			x0, x1, y0, y1 = r.X+r.W*0.375, r.X+r.W*0.66, r.Y+r.H*0.21, r.Y+r.H*0.79
		}
		if rtl {
			x0, x1 = x1, x0
		}
		var path ui.Path
		path.MoveTo(x0, y0).LineTo(x1, (y0+y1)/2).LineTo(x0, y1)
		pt.StrokePath(&path, width, color)
	})
}

// ---- history buttons and toolbars ----

// BackButton creates a button going back in r's history, disabled at its start, as the platform draws it: a toolbar item on macOS, a glass capsule around a chevron, NavigationView's back button on Windows, with its arrow, and a flat image button, as in a header bar, on GNOME. A right click lists the pages before, to go back several. It is ui.BackButton's counterpart, whose chevron gives way to the platform's glyph.
func BackButton(c *ui.Context, r *ui.Router) ui.Element {
	return navHistoryButton(c, r, -1)
}

// ForwardButton creates a button going forward in r's history, as BackButton goes back, drawn as the platform draws it, its glyph pointing the other way. It is ui.ForwardButton's counterpart.
func ForwardButton(c *ui.Context, r *ui.Router) ui.Element {
	return navHistoryButton(c, r, 1)
}

// navHistoryButton builds the back (d < 0) or forward button of package ui with the platform's look and glyph.
func navHistoryButton(c *ui.Context, r *ui.Router, d int) ui.Element {
	t := c.Theme()
	p := PlatformOf(c)
	// MyGo draws its chevron in the theme's Text, or TextMuted while disabled; a theme whose text is clear hides it, so that the platform's glyph is drawn over the button instead. Its tooltip keeps the colors it takes from the theme around it.
	clear := *t
	clear.Inverse, clear.InverseText = t.Inverse, t.InverseText
	if clear.Inverse == ui.Transparent {
		clear.Inverse = t.Text
	}
	if clear.InverseText == ui.Transparent {
		clear.InverseText = t.Background
	}
	clear.Text, clear.TextMuted = ui.Transparent, ui.Transparent
	var b ui.Element
	ui.Themed(c, &clear).Background(ui.Transparent).Children(func() {
		if d < 0 {
			b = ui.BackButton(c, r)
		} else {
			b = ui.ForwardButton(c, r)
		}
	})
	pressed, hovered, disabled := b.Pressed(), b.Hovered(), b.IsDisabled()
	// dir is 1 while the glyph points left, the way back goes left to right, and -1 while it points right: forward, or back right to left.
	dir := float32(-d)
	if c.Direction() == ui.RightToLeft {
		dir = -dir
	}
	var glyph func(pt *ui.Painter, rc ui.Rect)
	switch p {
	case Windows:
		// NavigationBackButtonNormalStyle: 40 by 36, ControlCornerRadius, NavigationViewBackButtonBackground SubtleFillColorTransparent, NavigationViewButtonBackgroundPointerOver SubtleFillColorSecondary and Pressed Tertiary; its glyph in NavigationViewItemForeground (TextFillColorPrimary), TextFillColorSecondary pressed and TextFillColorDisabled disabled, without a face then.
		hover, down := fluentSubtle(t)
		face, fg := ui.Transparent, t.Text
		switch {
		case disabled:
			fg = t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.361), 0.365))
		case pressed:
			face, fg = down, t.TextMuted
		case hovered:
			face = hover
		}
		b.Size(40, 36).Padding(0).Radius(4).Border(0, ui.Color{}).Background(face)
		glyph = func(pt *ui.Painter, rc ui.Rect) {
			// AnimatedBackVisualSource, the back arrow of a 16-DIP AnimatedIcon (its 48-unit visual scaled by 2.25 and a third): a shaft 12 DIPs long ending round and a head reaching 5.1 back and 4.9 to either side, strokes 0.75 DIP wide; ForwardButton mirrors it (inferred: NavigationView has no forward button).
			cx, cy := rc.X+rc.W/2, rc.Y+rc.H/2
			var arrow ui.Path
			arrow.MoveTo(cx+5.63*dir, cy).LineTo(cx-5.48*dir, cy)
			arrow.MoveTo(cx-0.37*dir, cy-4.88).LineTo(cx-5.48*dir, cy).LineTo(cx-0.37*dir, cy+4.88)
			pt.StrokePath(&arrow, 0.75, fg)
		}
	case Linux:
		// A flat .image-button: 34 pixels square (min-height 24 and padding 5), corners of $button_radius (9), $hover_color (currentColor at 7%) under the pointer and $active_color (16%) pressed, as buttons with an icon are in a header bar or a .toolbar (%toolbar extends %button_basic_flat); disabled, the button at $strong_disabled_opacity (30%).
		face, fg := ui.Transparent, t.Text
		switch {
		case disabled:
			fg = t.Background.Mix(t.Text, 0.3)
		case pressed:
			face = adwInkAlpha(t, t.Background, 0.16)
		case hovered:
			face = adwInkAlpha(t, t.Background, 0.07)
		}
		b.Size(34, 34).Padding(0).Radius(9).Border(0, ui.Color{}).Background(face)
		glyph = func(pt *ui.Painter, rc ui.Rect) {
			// Adwaita's go-previous-symbolic (go-next-symbolic mirrored), adwaita-icon-theme 48: a chevron in a 16-pixel icon from 11,2 to 5,8 to 11,14, a 2-pixel round stroke.
			cx, cy := rc.X+rc.W/2, rc.Y+rc.H/2
			var chevron ui.Path
			chevron.MoveTo(cx+3*dir, cy-6).LineTo(cx-3*dir, cy).LineTo(cx+3*dir, cy+6)
			pt.StrokePath(&chevron, 2, fg)
		}
	default:
		// A bordered NSToolbarItem with an image in a window's toolbar, as AppKit renders it on macOS 27 offscreen, without the Liquid Glass material: a platter 36 points square, round, filled black at 1.6% (Aqua) or 6.7% (Dark Aqua) over the window, with a half-point line of black at 21% or 57% around it; pressed, black at 20% or white at 18%, measured. Adjacent items share one platter in AppKit; here each button has its own.
		// The face is never ui.Transparent, which MyGo's Toolbar takes for a face of its own and fills under the pointer: a toolbar item does not mark the pointer offscreen (inferred for a live window).
		face := t.Background.Mix(black, pick(t.Dark, float32(0.016), 0.067))
		edge := t.Background.Mix(black, pick(t.Dark, float32(0.21), 0.57))
		fg := t.Text
		switch {
		case disabled:
			// tertiaryLabelColor (inferred: offscreen, in an inactive window, AppKit draws the glyph at 46% of the label's contrast, and disabled at 43% of that).
			fg = tint(t.Background, t.Dark, pick(t.Dark, float32(0.259), 0.247))
		case pressed:
			face = pick(t.Dark, t.Background.Mix(black, 0.2), t.Background.Mix(white, 0.178))
			fg = pick(t.Dark, black, white)
		}
		b.Size(36, 36).Padding(0).Radius(18).Border(0.5, edge).Background(face)
		glyph = func(pt *ui.Painter, rc ui.Rect) {
			// SF Symbols' chevron.left (chevron.right mirrored) in the item, measured: about 8 by 14 points of ink with a stroke near 1.7, its point 4 points before the platter's middle.
			cx, cy := rc.X+rc.W/2, rc.Y+rc.H/2
			var chevron ui.Path
			chevron.MoveTo(cx+2.4*dir, cy-6.2).LineTo(cx-3.9*dir, cy).LineTo(cx+2.4*dir, cy+6.2)
			pt.StrokePath(&chevron, 1.7, fg)
		}
	}
	b.DrawOver(glyph)
	return b
}

// Toolbar creates a row of the controls fn builds, as along the top of a window, spaced as the platform's toolbar: NSToolbar's items on macOS, a CommandBar on Windows and a .toolbar box on GNOME. It is one stop of Tab, and Left and Right move the focus among its controls; controls that do not fit go into a menu at its end, as in ui.Toolbar. Package ui's buttons inside it take MyGo's toolbar look, without a face until hovered; this package's Button and Toggle keep the platform's button look, and BackButton and ForwardButton are drawn as the platform's. It is ui.Toolbar's counterpart.
func Toolbar(c *ui.Context, fn func()) ui.Element {
	tb := ui.Toolbar(c, fn)
	switch PlatformOf(c) {
	case Windows:
		// CommandBar: AppBarThemeCompactHeight 48, CommandBarBackground ControlFillColorTransparent, padding 4,0,0,0, and its buttons' faces 36 high inside AppBarButtonInnerBorderMargin 2,6,2,6, 4 apart.
		tb.MinHeight(48).Padding(6).Gap(4).Radius(4).Background(ui.Transparent)
	case Linux:
		// .toolbar: padding 6px, border-spacing 6px.
		tb.Padding(6).Gap(6)
	default:
		// NSToolbarView of a window's unified toolbar, rendered by AppKit on macOS 27: 52 points high, its items' 36-point platters 8 points from its edges and 8 apart, measured.
		tb.MinHeight(52).Padding(8).Gap(8)
	}
	return tb
}
