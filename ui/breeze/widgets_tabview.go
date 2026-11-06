package breeze

import (
	"github.com/egoist/mygo/ui"
)

// ---- tab views: document tabs that close, and a button adding one ----
//
// macOS: the tab bar of NSWindow's tabbing (two or more NSWindows joined with addTabbedWindow, a delegate answering newWindowForTab: for the + button), laid out offscreen and captured live at 2x in Aqua and Dark Aqua on macOS 27.2 by tools/appkit-reference/reference.swift (its tabs scene): an NSTabBar 28 points high under the title bar, its NSTabBarTrackView 8 points in from the window's edges, NSTabButtons 24 high 2 points inside the track and a point apart, NSTabBarTrackInversionView separators 1 by 16 in those points, and an NSTabBarNewTabButton 28 points square 4 after the track. The capture's window is not the active app's, so its titles are drawn dimmed; the titles' colors are AppKit's label colors (inferred), the fills are measured. The close button, which shows under the pointer, never shows in a capture: its place, 18 points square 4 in from the tab's leading edge, is its accessibility frame, and its look is inferred.
//
// Windows 11: WinUI's TabView of microsoft/microsoft-ui-xaml main at 5ed8dd3: TabView_themeresources.xaml, TabView.xaml (its TabViewItem, TabViewCloseButtonStyle and TabViewButtonStyle templates), TabViewItem.cpp (UpdateTabGeometry's selected shape, UpdateCloseButton) and TabView.cpp (UpdateTabWidths with TabWidthMode Equal, SetTabSeparatorOpacity, its KeyboardAccelerators, the choice after a tab is removed), with the brushes of Common_themeresources_any.xaml.
//
// GNOME: libadwaita 1.10.0's AdwTabBar: _tab-view.scss, adw-tab-bar.ui, adw-tab.ui and adw-tab.c (the close button's place and when it shows), adw-tab-box.c (SPACING, the tabs sharing the width with expand-tabs on, its default, the separators' places and update_separators) and adw-tab-view.c (its shortcuts), with the end action holding the new tab button, src/icons/adw-tab-new.svg, and adwaita-icon-theme 51's window-close-symbolic.

// tabViewClose is what the platform names a tab's close button for assistive technology: AppKit's NSTabButton's close button reads "Close" (macOS 27.2), WinUI's TabViewCloseButtonName and adw-tab.ui's tooltip "Close Tab".
func tabViewClose(p Platform) string {
	if p == Mac {
		return "Close"
	}
	return "Close Tab"
}

// tabViewAdd is what the platform names the button adding a tab: WinUI's TabViewAddButtonName "Add New Tab"; AppKit names its + button only by its symbol ("add"), so macOS takes the name of the File menu's item, "New Tab", as GNOME takes the libadwaita demo's tooltip (inferred).
func tabViewAdd(p Platform) string {
	if p == Windows {
		return "Add New Tab"
	}
	return "New Tab"
}

// tabViewAfterClose returns the tab chosen once tab i of n closes, as WinUI's TabView chooses it (TabView.cpp's OnItemsChanged) and AdwTabView does: the tab after the closed one, now at its index, or the last when it was the last; a tab before the chosen one moves the choice one back.
func tabViewAfterClose(selected, i, n int) int {
	switch {
	case n <= 1:
		return 0
	case i < selected:
		return selected - 1
	case i == selected && i == n-1:
		return i - 1
	}
	return selected
}

// TabView creates the tab strip of a window's documents: tabs choosing *selected among labels, each with a close button when close is not nil, and after them a button adding a tab when add is not nil. Build the chosen tab's content after it, as with Tabs. A click on tab i's close button calls close(i) once the strip is built; TabView has moved *selected first, as the platforms move the choice (to the tab after a closed chosen one, or before it when it was the last), so close removes label i and whatever it keeps for the tab, and a caller keeping the tab sets *selected back. A click on the add button calls add, which appends a label and usually chooses it.
//
// The tabs are MyGo's TabsBase: a tab list for assistive technology, of tabs named by their labels, Tab stopping once at the chosen tab, the arrows, Home and End choosing among them. The close buttons are buttons named as the platform names them, which the pointer and assistive technology press but Tab does not stop at, as WinUI's (IsTabStop False) and GTK's (can-focus False) are not; the add button is a Tab stop after the tabs. The platform's keys switch tabs while the focus is in the strip, or anywhere in the window when no element around the focus takes them: Ctrl+Tab and Ctrl+Shift+Tab on every platform (AppKit's Window menu, WinUI's KeyboardAccelerators, AdwTabView's shortcuts), wrapping around; Ctrl+F4 closes the chosen tab on Windows; on GNOME Ctrl+Page Up and Ctrl+Page Down switch too, Alt+1 to Alt+9 and Alt+0 choose the first ten tabs, and Ctrl+Home and Ctrl+End the first and the last while the focus is in the strip (in the window they would take those keys from text fields). Closing the tab with Ctrl+W or Cmd+W is the app's, as GNOME and macOS leave it to the app's menus.
//   - macOS: the tab bar of a window's tabs (NSWindow tabbing): a 28-point capsule track 8 points in from the strip's sides, 9% of the label color over the window, holding the tabs, equal in width, 24 points high 2 points inside it and a point apart, their titles centered in 13 points; the chosen tab a white capsule (white at 11.7% in Dark Aqua) in labelColor, the others in secondaryLabelColor; a separator 1 by 16 points, the label color at 12% over the track, between two tabs neither of which is chosen or under the pointer; the close button under the pointer at the tab's leading edge; the + button, a 28-point circle 4 points after the track; a separator line under the strip as under a title bar.
//   - Windows: WinUI's TabView with TabWidthMode Equal and CloseButtonOverlayMode Auto: tabs 32 DIPs high, 8 below the strip's top, sharing the width up to 240 each, their 12-DIP headers in TextFillColorSecondary, semibold in TextFillColorPrimary on the chosen tab, which SolidBackgroundFillColorTertiary fills with corners of 8 at the top and 4-DIP curves joining CardStrokeColorDefault's line along the strip's bottom; DividerStrokeColorDefault separators 16 high at the ends of tabs neither chosen nor beside the chosen one, nor under the pointer; a 32 by 24 close button with the Cancel glyph on every closable tab, as Auto shows them (it maps to Always); the 32 by 24 add button 3 DIPs after the tabs and above the strip's bottom.
//   - GNOME: AdwTabBar with expand-tabs: a bar of headerbar_bg_color with headerbar_shade_color along its bottom, the tabs 34 pixels high, 6 below its top, 5 apart and 5 from its ends, sharing its width, their titles centered; the chosen tab currentColor at 10% (13% under the pointer, 19% pressed), another one 7% under the pointer and 16% pressed, corners of 9, none of it with a single tab; a 1-pixel separator 28 high centered in the space before every tab neither it nor the tab before is chosen or under the pointer; a round 24-pixel flat close button with window-close-symbolic at the tab's end while it is chosen or under the pointer; the add button, a flat 34-pixel image button with adw-tab-new-symbolic, at the bar's end.
//
// Dragging tabs to reorder them, GNOME's Ctrl+Shift+Page Up and Page Down that move the chosen tab, and the middle click closing a tab on Windows and GNOME are not provided: see the documentation's Limits.
func TabView(c *ui.Context, selected *int, labels []string, close func(i int), add func()) ui.Element {
	p := PlatformOf(c)
	n := len(labels)
	strip := ui.Row(c).Shrink(0)
	// Whether the strip was disabled as it was last painted: the caller disables what TabView returns after it, so its window-wide keys learn of it a frame later.
	off := ui.Local(strip, "breeze-tabview-off", func() bool { return false })
	sel := *selected
	hovered := make([]bool, n)
	closing, adding := -1, false
	rtl := c.Direction() == ui.RightToLeft
	strip.DrawOver(func(*ui.Painter, ui.Rect) { *off = strip.IsDisabled() })
	strip.Children(func() {
		parts := ui.TabsBase(c, selected, n)
		list := parts.List.Grow(1).Shrink(1).MinWidth(0)
		switch p {
		case Windows:
			fluentTabViewStrip(c, strip, list, n)
		case Linux:
			adwTabViewStrip(c, strip, list)
		default:
			macTabViewStrip(c, strip, list)
		}
		list.Children(func() {
			for i, label := range labels {
				tab := parts.Tab(i).Label(label).Grow(1).Shrink(1).MinWidth(0).FocusRing(false)
				hovered[i] = tab.Hovered()
				chosen := i == sel
				var x ui.Element
				switch p {
				case Windows:
					x = fluentTabViewTab(c, tab, label, chosen, close != nil, i, hovered, sel, rtl)
				case Linux:
					x = adwTabViewTab(c, tab, label, chosen, close != nil, i, hovered, sel, n, rtl)
				default:
					x = macTabViewTab(c, tab, label, chosen, close != nil, i, hovered, sel, rtl)
				}
				x.Role(ui.RoleButton).Label(tabViewClose(p)).Hidden(close == nil)
				if close != nil && x.Clicked() {
					closing = i
				}
			}
		})
		b := ui.ButtonBase(c).Shrink(0).Label(tabViewAdd(p)).Hidden(add == nil).FocusRing(false)
		switch p {
		case Windows:
			fluentTabViewAdd(c, b)
		case Linux:
			adwTabViewAdd(c, b)
		default:
			macTabViewAdd(c, b)
		}
		if add != nil && b.Clicked() {
			adding = true
		}
	})
	if n > 0 {
		choose := func(i int) { *selected = (i%n + n) % n }
		key := func(mods ui.Modifiers, k ui.Key, window bool) bool {
			// Both register their shortcut every frame, the strip's first: with the focus in the strip it takes the key, else the window does.
			in := strip.Shortcut(mods, k)
			out := window && c.Shortcut(mods, k) && !*off
			return in || out
		}
		if key(ui.Ctrl, ui.KeyTab, true) {
			choose(*selected + 1)
		}
		if key(ui.Ctrl|ui.Shift, ui.KeyTab, true) {
			choose(*selected - 1)
		}
		switch p {
		case Windows:
			if key(ui.Ctrl, ui.KeyF4, true) && close != nil && closing < 0 {
				closing = *selected
			}
		case Linux:
			if key(ui.Ctrl, ui.KeyPageDown, true) {
				choose(*selected + 1)
			}
			if key(ui.Ctrl, ui.KeyPageUp, true) {
				choose(*selected - 1)
			}
			if key(ui.Ctrl, ui.KeyHome, false) {
				choose(0)
			}
			if key(ui.Ctrl, ui.KeyEnd, false) {
				choose(n - 1)
			}
			for i, k := range []ui.Key{ui.Key1, ui.Key2, ui.Key3, ui.Key4, ui.Key5, ui.Key6, ui.Key7, ui.Key8, ui.Key9, ui.Key0} {
				// AdwTabView's select_nth_page_cb lets a key pass for a page that is not there.
				if key(ui.Alt, k, true) && i < n {
					choose(i)
				}
			}
		}
	}
	if closing >= 0 && closing < n {
		*selected = tabViewAfterClose(*selected, closing, n)
		close(closing)
	}
	if adding {
		add()
	}
	return strip
}

// ---- macOS: the tab bar of NSWindow's tabbing ----

// macTabViewStrip lays out the window tab bar: the strip padded 8 points at the sides and 7 below the 28-point track, with a separator line along its bottom point (the title bar's, 67 points below a tabbed window's top where the bar spans 32 to 60), the track a capsule of the label color at 9% over the window (232 over white, 62 over the title bar's 42 in Dark Aqua) with a half-point rim (220 and 77), its tabs 2 points inside it and a point apart.
func macTabViewStrip(c *ui.Context, strip, list ui.Element) {
	t := c.Theme()
	strip.Padding(0, 8, 8, 8).Gap(4).AlignItems(ui.Center)
	line := macSeparator(t.Background, t.Dark)
	strip.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 1, W: r.W, H: 1}, line, 0)
	})
	track := tint(t.Background, t.Dark, macTabTrackAlpha)
	rim := tint(t.Background, t.Dark, pick(t.Dark, float32(0.137), 0.164))
	list.Height(28).Padding(2).Gap(1).Radius(14).Background(track).AlignItems(ui.Stretch)
	list.DrawOver(func(p *ui.Painter, r ui.Rect) { p.Stroke(r, rim, 14, 0.5) })
}

// macTabTrackAlpha is the share of the label color the tab bar's track lays over the window: 232 over white, 62 over 42 in Dark Aqua (measured on macOS 27.2).
const macTabTrackAlpha = 0.09

// macTabViewTab draws tab i of the window tab bar and builds its children, returning its close button: the chosen tab a capsule, white with a faint rim (254 inside 248 over the track in Aqua; white at 11.7% with a rim at 16% in Dark Aqua, 67 and 76 over 42), the title centered in labelColor, the others in secondaryLabelColor (inferred: the capture's window was not the active one, which dims them); under the pointer an unchosen tab shows a capsule of 5% and the close button, which takes 18 points 4 in from the leading edge, an xmark in secondaryLabelColor, darker under the pointer (inferred). The separator before tab i, 1 by 16 points 4 below the tab's top in the point between it and tab i-1, is the label color at 12% over the track (204 over 232, 85 over 62), hidden beside the chosen tab and one under the pointer. AppKit's focus ring goes around a tab with the focus (inferred from the segments').
func macTabViewTab(c *ui.Context, tab ui.Element, label string, chosen, closable bool, i int, hovered []bool, sel int, rtl bool) ui.Element {
	t := c.Theme()
	track := tint(t.Background, t.Dark, macTabTrackAlpha)
	under := track
	over := hovered[i]
	pressed := tab.Pressed()
	focused := tab.FocusVisible()
	tab.Radius(12).AlignItems(ui.Center).Padding(0, 4).TextColor(tint(track, t.Dark, pick(t.Dark, float32(0.498), 0.549)))
	switch {
	case chosen:
		under = pick(t.Dark, white, t.Background.Mix(white, 0.117))
		edge := pick(t.Dark, track.Mix(black, 0.03), t.Background.Mix(white, 0.164))
		tab.Background(under).Border(0.5, edge).TextColor(tint(under, t.Dark, 0.847))
	case pressed:
		under = tint(track, t.Dark, 0.08)
		tab.Background(under)
	case over:
		under = tint(track, t.Dark, 0.05)
		tab.Background(under)
	}
	sep := tint(track, t.Dark, 0.12)
	tab.Draw(func(p *ui.Painter, r ui.Rect) {
		if i == 0 || chosen || i-1 == sel || hovered[i] || hovered[i-1] {
			return
		}
		x := r.X - 1
		if rtl {
			x = r.X + r.W
		}
		p.Fill(ui.Rect{X: x, Y: r.Y + 4, W: 1, H: 16}, sep, 0)
	})
	tab.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if focused {
			macFocusRing(p, r, 12)
		}
	})
	var x ui.Element
	tab.Children(func() {
		x = ui.Box(c).Size(18, 18).Shrink(0).Radius(9).Center()
		shown := over && closable
		x.IgnoreInput(!shown)
		if !shown {
			x.Opacity(0)
		}
		ink := tint(under, t.Dark, pick(t.Dark, float32(0.498), 0.549))
		if x.Hovered() {
			x.Background(tint(under, t.Dark, 0.1))
			ink = tint(under, t.Dark, 0.847)
		}
		x.Draw(func(p *ui.Painter, r ui.Rect) {
			toastCross(p, ui.Rect{X: r.X + 5.5, Y: r.Y + 5.5, W: 7, H: 7}, 0, 1.25, ink)
		})
		ui.Text(c, label).SingleLine().Grow(1).Shrink(1).MinWidth(0).TextAlign(ui.Center)
		// As wide as the close button, so that the title stays centered on the tab.
		ui.Box(c).Size(18, 18).Shrink(0).Role(ui.RoleNone)
	})
	return x
}

// macTabViewAdd draws the tab bar's + button (NSTabBarNewTabButton): a 28-point circle of glass, near white with a half-point rim in Aqua (250 inside 226), white at 13.6% with a lighter rim in Dark Aqua (71 inside 126 over 42), and the plus, 10.5 points across in 1.3-point strokes, black or white (measured on macOS 27.2); darker pressed (inferred), with AppKit's focus ring around it.
func macTabViewAdd(c *ui.Context, b ui.Element) {
	t := c.Theme()
	face := pick(t.Dark, t.Background.Mix(black, 0.02), t.Background.Mix(white, 0.136))
	if b.Pressed() {
		face = tint(face, t.Dark, 0.08)
	}
	rim := pick(t.Dark, t.Background.Mix(black, 0.11), t.Background.Mix(white, 0.3))
	ink := pick(t.Dark, black, white)
	b.Size(28, 28).Radius(14).Background(face).Border(0.5, rim)
	macRing(b, 14)
	b.Draw(func(p *ui.Painter, r ui.Rect) {
		cx, cy := r.X+r.W/2, r.Y+r.H/2
		p.Line(cx-5.25, cy, cx+5.25, cy, 1.3, ink)
		p.Line(cx, cy-5.25, cx, cy+5.25, 1.3, ink)
	})
}

// ---- Windows: WinUI's TabView ----

// fluentTabViewStrip lays out a TabView's strip: TabViewHeaderPadding 0,8,0,0 above the 32-DIP tabs; the tabs 9 DIPs from the start (LeftContentColumn's MinWidth 2, the ScrollViewer's first column's MinWidth 2, its ScrollContentPresenter's Padding 1,0,0,0 and the ItemsPresenter's 4-DIP header; inferred from the templates, as no render was compared) and 4 before the add button's container (the footer), the tab column as wide as TabWidthMode Equal makes the tabs, at most TabViewItemMaxWidth 240 each; and TabViewBorderBrush, CardStrokeColorDefault, 1 DIP along the strip's bottom, which the chosen tab covers. TabViewItemMinWidth (100) is not kept: where n tabs do not fit, they shrink rather than scroll.
func fluentTabViewStrip(c *ui.Context, strip, list ui.Element, n int) {
	t := c.Theme()
	strip.AlignItems(ui.End)
	_, _, line := fluentExpanderCard(t)
	strip.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 1, W: r.W, H: 1}, line, 0)
	})
	list.Padding(8, 0, 0, 0).PaddingStart(9).PaddingEnd(4).MaxWidth(13 + 240*float32(max(n, 1))).AlignItems(ui.Stretch)
}

// fluentTabViewTab draws TabViewItem i and builds its children, returning its close button. A tab is at least 32 DIPs high (TabViewItemMinHeight), its TabContainer padded 8,3,4,3 with the close button and 8,3,8,3 without (TabViewItemHeaderPaddingWithCloseButton and WithoutCloseButton), corners of OverlayCornerRadius (8) at the top, its 12-DIP header (TabViewItemHeaderFontSize, following the text size) at the start in TextFillColorSecondary, TextFillColorTertiary pressed; under the pointer LayerOnMicaBaseAltFillColorSecondary (black at 3.9%, white at 5.9% in the dark theme) and pressed LayerOnMicaBaseAltFillColorDefault (white at 70.2%, #3a3a3a at 45.1%). The chosen tab draws TabGeometry, SolidBackgroundFillColorTertiary (#f9f9f9, #282828) from 4 DIPs before its start to 4 after its end along the bottom, with 4-DIP concave curves there and corners of 8 at the top, and around it TabViewSelectedItemBorderBrush, CardStrokeColorDefault 1 DIP outside its sides and along its top (TabViewSelectedItemHeaderMargin -1,0,-1,1), which LeftRadiusRenderArc and RightRadiusRenderArc carry down into the strip's line; its header semibold in TextFillColorPrimary. TabSeparator, DividerStrokeColorDefault 1 DIP at the tab's end, 8 from its top and bottom (TabViewItemSeparatorMargin), shows unless the tab is chosen, before the chosen one, or it or the next is under the pointer (SetTabSeparatorOpacity, HideLeftAdjacentTabSeparator). The focus visual goes around the TabContainer (IsTemplateFocusTarget).
func fluentTabViewTab(c *ui.Context, tab ui.Element, label string, chosen, closable bool, i int, hovered []bool, sel int, rtl bool) ui.Element {
	t := c.Theme()
	pressed := tab.Pressed()
	over := hovered[i]
	focused := tab.FocusVisible()
	tab.MinHeight(32).Padding(3, 4, 3, 8).PaddingStart(8).PaddingEnd(pick(closable, float32(4), 8)).AlignItems(ui.Center).
		Radius(8, 8, 0, 0).FontSize(t.FontSize * 12 / 14).TextColor(t.TextMuted)
	under := t.Background
	switch {
	case chosen:
		under = pick(t.Dark, ui.Hex("#f9f9f9"), ui.Hex("#282828"))
		tab.TextColor(t.Text).FontWeight(600)
	case pressed:
		under = pick(t.Dark, t.Background.Mix(white, 0.702), t.Background.Mix(ui.Hex("#3a3a3a"), 0.451))
		tab.Background(under).TextColor(t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0x72)/255, 0x87/255.0)))
	case over:
		under = pick(t.Dark, t.Background.Mix(black, 0.039), t.Background.Mix(white, 0.059))
		tab.Background(under)
	}
	_, _, stroke := fluentExpanderCard(t)
	divider := fluentDivider(t)
	tab.Draw(func(p *ui.Painter, r ui.Rect) {
		if chosen {
			fill, outline := fluentTabGeometry(r)
			p.FillPath(fill, under)
			p.StrokePath(outline, 1, stroke)
			return
		}
		next := i+1 < len(hovered) && hovered[i+1]
		if i+1 == sel || hovered[i] || next {
			return
		}
		x := r.X + r.W - 1
		if rtl {
			x = r.X
		}
		p.Fill(ui.Rect{X: x, Y: r.Y + 8, W: 1, H: r.H - 16}, divider, 0)
	})
	tab.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if focused {
			fluentFocusRing(p, t, r, 0, 0, 0, 0, 8)
		}
	})
	var x ui.Element
	tab.Children(func() {
		ui.Text(c, label).SingleLine().Grow(1).Shrink(1).MinWidth(0)
		x = ui.Box(c).Size(32, 24).Shrink(0).MarginStart(4).Radius(4).Center()
		fluentTabViewButtonFace(t, x, under)
		ink := t.Text
		if x.Pressed() {
			ink = fluentSecondaryOver(t, under)
		}
		x.Draw(func(p *ui.Painter, r ui.Rect) {
			// The Cancel glyph at TabViewItemHeaderCloseFontSize 12: InfoBar's at 16 is an × 11 DIPs square in 1-DIP strokes (measured from WinUI 3 3.3), here scaled by 12/16 (inferred).
			toastCross(p, ui.Rect{X: r.X + r.W/2 - 4.125, Y: r.Y + r.H/2 - 4.125, W: 8.25, H: 8.25}, 0, 1, ink)
		})
	})
	return x
}

// fluentTabViewButtonFace fills a TabView's close or add button, b, over under: transparent, SubtleFillColorSecondary under the pointer and SubtleFillColorTertiary pressed (TabViewItemHeaderCloseButtonBackground and TabViewButtonBackground, PointerOver and Pressed; black at 3.5% and 2.4%, white at 5.9% and 3.9% in the dark theme).
func fluentTabViewButtonFace(t *ui.Theme, b ui.Element, under ui.Color) {
	ink := pick(t.Dark, black, white)
	switch {
	case b.Pressed():
		b.Background(under.Mix(ink, pick(t.Dark, float32(0.024), 0.039)))
	case b.Hovered():
		b.Background(under.Mix(ink, pick(t.Dark, float32(0.035), 0.059)))
	}
}

// fluentTabGeometry returns the chosen tab's shape over its box r, as TabViewItem.cpp's UpdateTabGeometry draws it with the SelectedBackgroundPath's Margin -4,0, and the line around it, the centers of the TabContainer's 1-DIP border (corners of 8 at its outer edge, 1 DIP outside the shape's sides) and of the radius arcs' quarter rings (radii 3 to 4, centered 4 outside the tab and 4 above its bottom), which end on the strip's line.
func fluentTabGeometry(r ui.Rect) (fill, outline *ui.Path) {
	const k = 0.5523 // a quarter circle as a cubic curve
	x0, x1, y0, y1 := r.X, r.X+r.W, r.Y, r.Y+r.H
	fill = new(ui.Path)
	fill.MoveTo(x0-4, y1).
		CubeTo(x0-4+4*k, y1, x0, y1-4+4*k, x0, y1-4).
		LineTo(x0, y0+8).
		CubeTo(x0, y0+8-8*k, x0+8-8*k, y0, x0+8, y0).
		LineTo(x1-8, y0).
		CubeTo(x1-8+8*k, y0, x1, y0+8-8*k, x1, y0+8).
		LineTo(x1, y1-4).
		CubeTo(x1, y1-4+4*k, x1+4-4*k, y1, x1+4, y1).
		Close()
	outline = new(ui.Path)
	outline.MoveTo(x0-4, y1-0.5).
		CubeTo(x0-4+3.5*k, y1-0.5, x0-0.5, y1-4+3.5*k, x0-0.5, y1-4).
		LineTo(x0-0.5, y0+8).
		CubeTo(x0-0.5, y0+8-7.5*k, x0+7-7.5*k, y0+0.5, x0+7, y0+0.5).
		LineTo(x1-7, y0+0.5).
		CubeTo(x1-7+7.5*k, y0+0.5, x1+0.5, y0+8-7.5*k, x1+0.5, y0+8).
		LineTo(x1+0.5, y1-4).
		CubeTo(x1+0.5, y1-4+3.5*k, x1+4-3.5*k, y1-0.5, x1+4, y1-0.5)
	outline.Cap(ui.ButtCap)
	return fill, outline
}

// fluentTabViewAdd draws the TabView's AddButton (TabViewButtonStyle): 32 by 24 (TabViewItemAddButtonWidth and Height) in a container padded 3,0,0,3 at the strip's bottom (TabViewItemAddButtonContainerPadding), corners of ControlCornerRadius, the Add glyph at 12 (TabViewItemAddButtonFontSize) in TextFillColorPrimary, TextFillColorSecondary pressed, its face as the close button's, and the system focus visual 3 DIPs around it (FocusVisualMargin -3). The glyph's ink, a plus 9 DIPs across in 1-DIP strokes, is inferred from the Cancel glyph's.
func fluentTabViewAdd(c *ui.Context, b ui.Element) {
	t := c.Theme()
	b.Size(32, 24).Margin(0, 0, 3, 0).MarginStart(3).Radius(4)
	fluentTabViewButtonFace(t, b, t.Background)
	ink := t.Text
	if b.Pressed() {
		ink = t.TextMuted
	}
	focused := b.FocusVisible()
	b.Draw(func(p *ui.Painter, r ui.Rect) {
		cx, cy := r.X+r.W/2, r.Y+r.H/2
		p.Fill(ui.Rect{X: cx - 4.5, Y: cy - 0.5, W: 9, H: 1}, ink, 0)
		p.Fill(ui.Rect{X: cx - 0.5, Y: cy - 4.5, W: 1, H: 9}, ink, 0)
	})
	b.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if focused {
			fluentFocusRing(p, t, r, 3, 3, 3, 3, 4)
		}
	})
}

// ---- GNOME: libadwaita's AdwTabBar ----

// adwTabViewBar returns AdwTabBar's .box colors: headerbar_bg_color (white, #2e2e32 in the dark theme) and headerbar_shade_color over it (the shade at 12%, 36%).
func adwTabViewBar(t *ui.Theme) (bg, shade ui.Color) {
	bg = pick(t.Dark, white, ui.Hex("#2e2e32"))
	return bg, bg.Mix(adwInk, pick(t.Dark, float32(0.12), 0.36))
}

// adwTabViewStrip lays out an AdwTabBar: .box padded 1 pixel at the sides and the bottom, none at the top, its inset shade along the bottom over that pixel; the tabbox padded 6 above and below, its tabs SPACING (5) apart and 5 from its ends, sharing its width (expand-tabs, AdwTabBar's default).
func adwTabViewStrip(c *ui.Context, strip, list ui.Element) {
	t := c.Theme()
	bg, shade := adwTabViewBar(t)
	strip.Background(bg).TextColor(adwTextOn(t, bg)).Padding(0, 1, 1, 1).AlignItems(ui.Center).BorderWidth(0, 0, 1, 0).BorderColor(shade)
	list.Padding(6, 5).Gap(5).AlignItems(ui.Stretch)
}

// adwTabViewTab draws AdwTab i and builds its children, returning its close button: 34 pixels high (min-height 26px, padding 4px), corners of $button_radius 9, its title in a .tab-title label with margins of 4, centered on the tab (adw_tab_size_allocate centers it across the whole tab, so the close button is matched by as much room at the start); currentColor at 10%, 13% under the pointer and 19% pressed while chosen, 7% under the pointer and 16% pressed otherwise, and none in a bar of one tab (.single-tab); the close button (window-close-symbolic in a .flat button 24 pixels round, currentColor at 7% under the pointer and 16% pressed) at the end, shown while the tab is chosen or under the pointer (update_state) and taking no clicks otherwise (can-target follows its opacity). The separator before tab i, a GtkSeparator of $border_color (currentColor at 15%) 1 pixel wide and 3 in from the tabbox's content top and bottom, is centered in the 5-pixel space before the tab, hidden before the first tab and where it or the tab before is chosen, under the pointer or pressed (update_separators). The focus ring is inside the tab's edge (tabboxchild's focus-ring mixin).
func adwTabViewTab(c *ui.Context, tab ui.Element, label string, chosen, closable bool, i int, hovered []bool, sel, n int, rtl bool) ui.Element {
	t := c.Theme()
	bg, _ := adwTabViewBar(t)
	pressed := tab.Pressed()
	over := hovered[i]
	tab.MinHeight(34).Padding(4).Radius(9).AlignItems(ui.Center)
	share := float32(0)
	if n > 1 {
		switch {
		case chosen && pressed:
			share = 0.19
		case chosen && over:
			share = 0.13
		case chosen:
			share = 0.10
		case pressed:
			share = 0.16
		case over:
			share = 0.07
		}
	}
	under := bg
	if share > 0 {
		under = adwInkAlpha(t, bg, share)
		tab.Background(under)
	}
	sep := adwInkAlpha(t, bg, 0.15)
	tab.Draw(func(p *ui.Painter, r ui.Rect) {
		if i == 0 || chosen || i-1 == sel || hovered[i] || hovered[i-1] {
			return
		}
		x := r.X - 3
		if rtl {
			x = r.X + r.W + 2
		}
		p.Fill(ui.Rect{X: x, Y: r.Y + 3, W: 1, H: r.H - 6}, sep, 0)
	})
	adwFocusRing(t, tab, 9, nil)
	var x ui.Element
	tab.Children(func() {
		// As wide as the close button, so that the title stays centered on the tab.
		ui.Box(c).Size(24, 24).Shrink(0).Role(ui.RoleNone)
		ui.Text(c, label).SingleLine().Grow(1).Shrink(1).MinWidth(0).MarginX(4).TextAlign(ui.Center)
		x = ui.Box(c).Size(24, 24).Shrink(0).Radius(12).Center()
		shown := closable && (chosen || over)
		x.IgnoreInput(!shown)
		if !shown {
			x.Opacity(0)
		}
		switch {
		case x.Pressed():
			x.Background(adwInkAlpha(t, under, 0.16))
		case x.Hovered():
			x.Background(adwInkAlpha(t, under, 0.07))
		}
		ink := adwTextOn(t, under)
		x.Draw(func(p *ui.Painter, r ui.Rect) {
			// window-close-symbolic, a 16-pixel icon: the cross stroked 2 wide past the square its arms end in, clipped to it.
			o := ui.Rect{X: r.X + r.W/2 - 8, Y: r.Y + r.H/2 - 8, W: 16, H: 16}
			p.Clip(ui.Rect{X: o.X + 4, Y: o.Y + 4, W: 8, H: 8}, 0, func() { toastCross(p, o, 0.1875, 2, ink) })
		})
	})
	return x
}

// adwTabViewAdd draws the button in AdwTabBar's end action (.end-action, padded 6px 5px with no padding toward the tabs): a 34-pixel .image-button, flat as a header bar's (%button_basic_flat: currentColor at 7% under the pointer and 16% pressed; GTK leaves a button in the end action raised unless it is .flat, and libadwaita's demo puts its New Tab button in the header bar), showing adw-tab-new-symbolic: a 14-pixel square of corners about 3 and a 6-pixel plus inside it, stroked 2 wide with round ends.
func adwTabViewAdd(c *ui.Context, b ui.Element) {
	t := c.Theme()
	bg, _ := adwTabViewBar(t)
	b.Size(34, 34).Radius(9).Margin(6, 0).MarginEnd(5)
	adwFocusRing(t, b, 9, nil)
	face := bg
	switch {
	case b.Pressed():
		face = adwInkAlpha(t, bg, 0.16)
	case b.Hovered():
		face = adwInkAlpha(t, bg, 0.07)
	}
	b.Background(face)
	ink := adwTextOn(t, face)
	b.Draw(func(p *ui.Painter, r ui.Rect) {
		const k = 0.5523
		x, y := r.X+r.W/2-8, r.Y+r.H/2-8
		var path ui.Path
		path.MoveTo(x+8, y+5).LineTo(x+8, y+11).MoveTo(x+5, y+8).LineTo(x+11, y+8)
		path.MoveTo(x+3.965, y+1).LineTo(x+12.035, y+1).
			CubeTo(x+12.035+2.965*k, y+1, x+15, y+3.965-2.965*k, x+15, y+3.965).
			LineTo(x+15, y+12.035).
			CubeTo(x+15, y+12.035+2.965*k, x+12.035+2.965*k, y+15, x+12.035, y+15).
			LineTo(x+3.965, y+15).
			CubeTo(x+3.965-2.965*k, y+15, x+1, y+12.035+2.965*k, x+1, y+12.035).
			LineTo(x+1, y+3.965).
			CubeTo(x+1, y+3.965-2.965*k, x+3.965-2.965*k, y+1, x+3.965, y+1).Close()
		p.StrokePath(&path, 2, ink)
	})
}
