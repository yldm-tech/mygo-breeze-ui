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
	var top float32
	switch p {
	case Windows:
		// NavigationViewItemHeaderTextStyle: 14 pixels, semibold, TextFillColorSecondary, 16 from the pane's edge. WinUI's header row is 40 high; MyGo's padding of the title, a third of it above and below, and 4 above the section (inferred) come to about 36.
		head.FontSize, head.Spacing, head.Radius = 15, 6, 4
		top = 4
	case Linux:
		// libadwaita 1.7 has no sidebar headings: .heading .dim-label, the text at 55% (inferred), lined up with the rows' labels, 12 pixels below the row above (inferred).
		head.FontSize, head.Spacing, head.Radius = t.FontSize+1, 4, 9
		head.TextMuted = adwInkAlpha(t, adwSidebarBg(t), 0.55)
		top = 12
	default:
		// AppKit's group row: 11-point text in a row 19 points high (MyGo's padding of the title makes it about 17), 2 points inside the selection's inset and 13 points below the row above, measured, in tertiaryLabelColor (black at 26% in Aqua, white at 25% in Dark Aqua). Rendered offscreen, without the sidebar's vibrancy, AppKit draws it lighter still (20% and 12%).
		head.FontSize, head.Spacing, head.Radius = 12, 2, 8
		head.TextMuted = tint(macSidebarBg(t), t.Dark, pick(t.Dark, float32(0.259), 0.247))
		top = 13
	}
	st := navCurrentSidebar(c)
	wrap := ui.Themed(c, &head).Background(ui.Transparent)
	if st != nil {
		if st.built > 0 {
			wrap.Margin(top, 0, 0, 0)
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
			ui.Themed(c, t).Background(ui.Transparent).Gap(gap).Children(fn)
		})
	})
	return sec
}

// SidebarItem creates an item of a Sidebar choosing id, showing label after icon, nil for none, as the platform draws a row of its sidebar: on macOS 32 points high with a rounded selection, in the accent color while the sidebar has the focus and in gray otherwise; on Windows a NavigationViewItem, 36 high with a 3 by 16 accent pill at the start of the chosen one; on GNOME a row 36 high, the chosen one darker. Add to it with Children, as a Badge. Build it in this package's Sidebar, which tells it whether it is chosen and whether the sidebar has the focus. It is ui.SidebarItem's counterpart.
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
		adwSidebarItem(c, item, chosen)
	default:
		macSidebarItem(c, item, chosen, focused)
	}
	return item
}

// macSidebarBg is a source list's background, measured from an NSOutlineView of style sourceList rendered by AppKit on macOS 27.
func macSidebarBg(t *ui.Theme) ui.Color { return pick(t.Dark, ui.Hex("#f0f0f0"), ui.Hex("#282828")) }

// macSidebarItem draws a row of an NSOutlineView of style sourceList at the default (medium) row size, measured on macOS 27: 32 points high, an icon or the text 6 points inside the selection, the text 43 points in after a 16-point icon, as AppKit lays out a cell's image and text; the selection has corners of 8 and fills the row, and while the outline view has the focus it is the emphasized selection in the accent color with white text, otherwise unemphasizedSelectedContentBackgroundColor, AppKit's semantic color (offscreen, AppKit draws that selection without its material, in black).
func macSidebarItem(c *ui.Context, item ui.Element, chosen, focused bool) {
	t := c.Theme()
	item.MinHeight(32).Padding(0, 6).Gap(11).Radius(8)
	switch {
	case chosen && focused:
		item.Background(macSidebarSelection(t)).TextColor(white)
	case chosen:
		item.Background(pick(t.Dark, ui.Hex("#dcdcdc"), ui.Hex("#464646"))).TextColor(t.Text)
	}
}

// macSidebarSelection is the emphasized selection of a source list: #0070f5 in Aqua and #007aff, controlAccentColor, in Dark Aqua for the blue accent, measured from a row view drawn emphasized (AppKit renders it in Display P3, converted here to sRGB). For another accent it is the accent, a little darker in Aqua (inferred).
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
		// The chevron of Segoe Fluent Icons, BreadcrumbBarChevronLeftToRight U+E974 (U+E973 right to left, mirrored here), BreadcrumbBarChevronFontSize 12 with BreadcrumbBarChevronPadding 2,0, in TextFillColorPrimary.
		sep.Size(12, 12).MarginX(2)
		color, width = t.Text, 1
	} else {
		// NSPathControl's chevron: 4 by 6 points in secondaryLabelColor, 7.5 points after an item's text and 5.5 before the next, measured.
		sep.Size(4, 6).MarginStart(5).MarginEnd(3)
	}
	sep.Draw(func(pt *ui.Painter, r ui.Rect) {
		x0, x1, y0, y1 := r.X, r.X+r.W, r.Y, r.Y+r.H
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

// BackButton creates a button going back in r's history, disabled at its start, as the platform draws it: a toolbar button on macOS, NavigationView's back button on Windows and a flat image button, as in a header bar, on GNOME. A right click lists the pages before, to go back several. Its arrow is MyGo's. It is ui.BackButton's counterpart.
func BackButton(c *ui.Context, r *ui.Router) ui.Element {
	return navHistoryButton(c, ui.BackButton(c, r))
}

// ForwardButton creates a button going forward in r's history, as BackButton goes back, drawn as the platform draws it. It is ui.ForwardButton's counterpart.
func ForwardButton(c *ui.Context, r *ui.Router) ui.Element {
	return navHistoryButton(c, ui.ForwardButton(c, r))
}

// navHistoryButton gives a back or forward button of package ui the platform's look.
func navHistoryButton(c *ui.Context, b ui.Element) ui.Element {
	t := c.Theme()
	switch PlatformOf(c) {
	case Windows:
		// NavigationBackButtonNormalStyle: 40 by 36, ControlCornerRadius, SubtleFillColorTransparent, SubtleFillColorSecondary under the pointer and Tertiary pressed, the same as an AppBarButton in a CommandBar.
		hover, pressed := fluentSubtle(t)
		face := ui.Transparent
		switch {
		case b.Pressed():
			face = pressed
		case b.Hovered():
			face = hover
		}
		b.Size(40, 36).Padding(0).Radius(4).Border(0, ui.Color{}).Background(face)
	case Linux:
		// A flat .image-button: 34 pixels square (24 and padding 5), corners of 9, $hover_color (currentColor at 7%) under the pointer and $active_color (16%) pressed, as buttons with an icon are in a header bar or a .toolbar.
		share := float32(0)
		switch {
		case b.Pressed():
			share = 0.16
		case b.Hovered():
			share = 0.07
		}
		face := ui.Transparent
		if share > 0 {
			face = adwInkAlpha(t, t.Background, share)
		}
		b.Size(34, 34).Padding(0).Radius(9).Border(0, ui.Color{}).Background(face)
	default:
		// An NSButton of bezel style toolbar with an image, as AppKit renders it on macOS 27 outside a window's toolbar: 22 by 23 points, the bezel a point inside, corners of 4 and a 1-point line of black (Aqua) or white (Dark Aqua) at 5%, filled at 10% while pressed, measured.
		// The clear face is white at 0% rather than ui.Transparent, which MyGo's Toolbar takes for a face of its own and fills under the pointer: an NSButton does not track the pointer, so the bezel stays as it is until pressed (inferred).
		face, edge := white.Alpha(0), tint(t.Background, t.Dark, 0.051)
		if b.Pressed() {
			face = tint(t.Background, t.Dark, 0.102)
			edge = face
		}
		b.Size(20, 21).Margin(1).Padding(0).Radius(4).Border(1, edge).Background(face)
	}
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
		// NSToolbar's items, 8 points apart (inferred: AppKit does not render a toolbar offscreen).
		tb.Padding(4, 8).Gap(8)
	}
	return tb
}
