package breeze

import "github.com/egoist/mygo/ui"

// ---- header bar: the title bar a view draws in a window whose title bar is hidden ----
//
// macOS: a titled window's unified toolbar showing its title (NSWindow.toolbarStyle .unified, titleVisibility visible), rendered at 2x in Aqua and Dark Aqua on macOS 27.2 by tools/appkit-reference/reference.swift (its header bar scene: renderHeaderBar), with a navigational back item before the title, a bordered text item and a bordered image item after a flexible space, with and without NSWindow.subtitle. Windows 11: WinUI's TitleBar control (controls/dev/TitleBar of microsoft/microsoft-ui-xaml main at cb1b8a6: TitleBar.xaml, TitleBar_themeresources.xaml, TitleBar.cpp and TitleBarAutomationPeer.cpp). GNOME: libadwaita 1.10.0's AdwHeaderBar and AdwWindowTitle (_header-bar.scss, _toolbars.scss and _colors.scss) and GTK 4.24.1's font-size: smaller (gtkcssenumvalue.c).

// headerMacStartGap is the room AppKit leaves after the traffic lights: the first item of a unified toolbar starts 17 points after the zoom button's right edge (its platter at 96, the zoom button ending at 79), and without an item before it the title starts 29 points after it (108), 12 more, measured on macOS 27.2.
const headerMacStartGap = 17

// headerMacTitleGap is the room before the title: 12 points after the platter of the item before it (the title's text at 144 after a platter ending at 132), or after where the first item would start, measured on macOS 27.2.
const headerMacTitleGap = 12

// headerFluentMinDragRegion is TitleBarMinDragRegionWidth (TitleBar_themeresources.xaml): a column 48 DIPs wide that the TitleBar keeps between its RightHeader and the caption buttons, so that there is always room to drag the window.
const headerFluentMinDragRegion = 48

// headerLayout is where a platform puts the parts of its title bar.
type headerLayout struct {
	height      float32 // the bar's height at least
	start, end  float32 // the room before the first control and after the last, from the bar's edges
	top, bottom float32 // the bar's padding above and below its content
	gap         float32 // the room between two controls
	titleGap    float32 // the room before the title
	// before and after are the room the window controls take at the bar's start and end, which GNOME's start and end boxes hold (headerAdwSide).
	before, after float32
}

// headerLayoutOf returns p's layout of a title bar with controls (items) or without, in a window whose controls take bar's room: before, the room at the start of the bar (Left, or Right right to left), and after, at its end.
//   - macOS: a unified toolbar is 52 points high, its items 8 apart, the first 17 after the traffic lights (8 from the window's edge without them, as the last is from the other edge: inferred from the trailing edge), and the title 12 after the item before it.
//   - Windows: TitleBarExpandedHeight 48 with a LeftHeader, Content or RightHeader, TitleBarCompactHeight 32 without (TitleBar::UpdateHeight); the LeftPaddingColumn as wide as the caption buttons at the start, none without them (TitleBar::UpdatePadding replaces TitleBarLeftPaddingWidth's 2 with the window's AppWindowTitleBar.LeftInset, RightInset right to left); the RightHeader before the 48-DIP TitleBarMinDragRegionWidth column and the caption buttons' RightPaddingColumn; the title 14 after the LeftHeader (TitleBarLeftHeaderPaddingWidth), so 16 after the visible edge of an image button there, which keeps the back button's Margin of 2 (headerFluentButton). WinUI spaces no controls inside a header, which the app's own panel does: 4 apart, as a CommandBar's buttons (inferred).
//   - GNOME: an AdwHeaderBar is at least 47 pixels high, its box padded 6 above, 7 below and 7 at either end, its start and end boxes 6 apart inside (border-spacing) and 6 from the title (_header-bar.scss); the bar's margins of -1px at its sides put it a pixel past the window's edges, so its controls start and end 6 from them (the back button at 6 and GTK's close button ending 6 from the edge in libadwaita 1.10.0's render, tools/adwaita-reference's headerbar set); the window controls are the first child of the start box or the last of the end box, so the title is centered across the whole bar.
func headerLayoutOf(p Platform, bar ui.TitleBar, before, after float32, items bool) headerLayout {
	var l headerLayout
	switch p {
	case Windows:
		l = headerLayout{height: 32, start: before, end: after + headerFluentMinDragRegion, gap: 4, titleGap: 14}
		if items {
			l.height = 48
		}
	case Linux:
		l = headerLayout{height: 47, start: 6, end: 6, top: 6, bottom: 7, gap: 6, titleGap: 6, before: before, after: after}
	default:
		l = headerLayout{height: 52, start: 8, end: 8, gap: idiomMacToolbarGap, titleGap: headerMacTitleGap}
		if before > 0 {
			l.start = before + headerMacStartGap
		}
		if after > 0 {
			l.end = after + headerMacStartGap
		}
	}
	l.height = max(l.height, bar.Height)
	return l
}

// HeaderBar creates the title bar of a window whose title bar is hidden (mygo.WindowOptions.TitleBarStyle), as the platform draws a window's own: the title, the subtitle under it or beside it (an empty subtitle shows none), and the controls start and end build before and after it (nil builds none). It keeps clear of the window controls, in the room c.TitleBar() reports, is at least as high as they are, and moves the window when dragged anywhere but over its controls (ui.Element.DragWindow), maximizing it on a double click. Put it at the top of the window's column.
//
// The controls of this package built in start and end take the platform's look of a title bar's controls, as they take a toolbar's in a Toolbar: on macOS bordered NSToolbarItems, 36-point capsules, image items side by side sharing a platter; on Windows the subtle 40-DIP buttons of a TitleBar for image buttons (IconButton, BackButton, ForwardButton, a Button with an empty label), the others drawn as outside a toolbar, as WinUI's TitleBar gives its headers' content no look of its own; on GNOME an AdwHeaderBar's, flat image and menu buttons and raised text buttons. Each control is a Tab stop of its own, as in a header bar of GTK and a TitleBar of WinUI.
//   - macOS: the title bar of a window with a unified toolbar on macOS 26 and 27, 52 points high on the window's background with a half-point line along its bottom (black at 10.2%, white at 7.1% in Dark Aqua, measured offscreen on macOS 27.2), the start controls 17 points after the traffic lights as navigational items, then the title 12 points after them in 15-point semibold labelColor, or with a subtitle in 13-point bold over the 11-point subtitle in secondaryLabelColor, then the end controls at the far end, 8 points from the edge. AppKit centers the traffic lights in a 52-point bar (19 points down and 19 in); MyGo leaves them where its hidden title bar puts them unless mygo.WindowOptions.TrafficLightPosition moves them, which {X: 19, Y: 19} does.
//   - Windows: WinUI's TitleBar, 48 DIPs high with controls and 32 without, transparent over the window's backdrop, the start controls as its LeftHeader from the window's edge (or the caption buttons, at the start right to left), an image button 2 DIPs in, then the title 14 DIPs after them in CaptionTextBlockStyle (12) and TextFillColorPrimary and the subtitle 8 after it, in the same style and TextFillColorSecondary, the end controls as its RightHeader, 48 DIPs before the caption buttons, an image button 2 DIPs inside it. MyGo's Windows host draws the caption buttons itself, as tall as mygo.WindowOptions.TitleBarHeight (32 by default; 48 matches a bar with controls, as WinUI apps ask with PreferredHeightOption Tall).
//   - GNOME: libadwaita's AdwHeaderBar, 47 pixels high at least, on headerbar_bg_color with headerbar_shade_color along its bottom, the title centered across the bar as AdwWindowTitle's, bold, with the subtitle under it smaller and dimmed, the start and end controls at the bar's ends 6 apart. GTK draws the window buttons themselves.
//
// The bar is a group that assistive technology names by title, as WinUI's TitleBarAutomationPeer names the TitleBar by its Title. MyGo v0.7.4 does not tell a view whether its window is active, so the bar is drawn as in the active window: the inactive window's dimmed title (AppKit's), the TitleBar's TextFillColorTertiary title and half-opaque headers, and the header bar's backdrop (window_bg_color, its content at 50%) are not drawn.
func HeaderBar(c *ui.Context, title, subtitle string, start, end func()) ui.Element {
	p := PlatformOf(c)
	t := c.Theme()
	room := c.TitleBar()
	before, after := room.Left, room.Right
	if c.Direction() == ui.RightToLeft {
		before, after = after, before
	}
	l := headerLayoutOf(p, room, before, after, start != nil || end != nil)
	left, right := l.start, l.end
	if c.Direction() == ui.RightToLeft {
		left, right = right, left
	}
	var bar ui.Element
	if p == Linux {
		// The controls of the bar are drawn over headerbar_bg_color, which they take their faces from as from a theme's background.
		th := *t
		th.Background = adwHeaderBarBackground(t)
		bar = ui.Themed(c, &th).Row()
	} else {
		bar = ui.Row(c)
	}
	bar.Shrink(0).FillWidth().MinHeight(l.height).AlignItems(ui.Center).Gap(0).Padding(l.top, right, l.bottom, left).DragWindow().Role(ui.RoleGroup).Label(title)
	bar.Children(func() {
		switch p {
		case Linux:
			headerAdwContent(c, l, title, subtitle, start, end)
		case Windows:
			headerFluentContent(c, l, title, subtitle, start, end)
		default:
			headerMacContent(c, l, title, subtitle, start, end)
		}
	})
	switch p {
	case Linux:
		adwHeaderBar(c, bar)
	case Windows:
		// The TitleBar's Background is Transparent: the window's backdrop shows through it.
		bar.Background(ui.Transparent)
	default:
		macHeaderBar(c, bar)
	}
	return bar
}

// headerItems builds fn's controls in a row of the bar, gap apart, as the controls of a toolbar's start or end (idiomWithToolbar), so that they take the platform's look of a title bar's controls. The row is built hidden when fn is nil, so that what follows it keeps its place among its siblings.
func headerItems(c *ui.Context, l headerLayout, fn func()) ui.Element {
	row := ui.Row(c).Shrink(0).AlignItems(ui.Center).Gap(l.gap).Hidden(fn == nil)
	if fn != nil {
		idiomWithToolbar(c, &idiomToolbar{header: l.height}, func() { row.Children(fn) })
	}
	return row
}

// ---- macOS ----

// headerMacContent lays out the parts of a unified toolbar showing the window's title: the start controls, the title 12 points after them (or where they would start), growing, and the end controls at the far end.
func headerMacContent(c *ui.Context, l headerLayout, title, subtitle string, start, end func()) {
	t := c.Theme()
	headerItems(c, l, start)
	ui.Column(c).Grow(1).Shrink(1).MinWidth(0).MarginStart(l.titleGap).MarginEnd(l.gap).Justify(ui.Center).Children(func() {
		if subtitle == "" {
			// _NSToolbarTitleField alone: .SFNS-Semibold at 15 points in labelColor (NSFont.fontName and pointSize read on macOS 27.2).
			ui.Text(c, title).SingleLine().FontSize(t.Rem(15.0 / 13)).FontWeight(600).TextColor(t.Text)
		} else {
			// With NSWindow.subtitle: the title in .SFNS-Bold at 13 points over the subtitle in .SFNS-Regular at 11 points and secondaryLabelColor, the two 30 points high together, centered.
			ui.Text(c, title).SingleLine().FontWeight(700).TextColor(t.Text)
		}
		ui.Text(c, subtitle).SingleLine().FontSize(t.Rem(11.0 / 13)).TextColor(t.TextMuted).Hidden(subtitle == "")
	})
	headerItems(c, l, end)
}

// macHeaderBar paints the unified toolbar's background, the window's, and the line AppKit draws along its bottom: half a point of black at 10.2% in Aqua (229 over the bar's 255) and of white at 7.1% in Dark Aqua (46 over 30), measured offscreen on macOS 27.2. AppKit draws it just below the 52-point bar, over the content's first half point; it is drawn inside the bar's bottom edge here, so that a background of the content below does not cover it.
func macHeaderBar(c *ui.Context, bar ui.Element) {
	t := c.Theme()
	bar.Background(t.Background)
	line := tint(t.Background, t.Dark, pick(t.Dark, float32(0.102), 0.071))
	bar.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 0.5, W: r.W, H: 0.5}, line, 0)
	})
}

// ---- Windows ----

// headerFluentContent lays out the parts of WinUI's TitleBar: the LeftHeader, the title 14 DIPs after it in CaptionTextBlockStyle and TitleBarForegroundBrush (TextFillColorPrimary) with TitleBarTitleMargin 0,0,8,0, the subtitle in the same style and TitleBarSubtitleForegroundBrush (TextFillColorSecondary) with TitleBarSubtitleMargin 0,0,16,0, both trimmed with an ellipsis and centered vertically, the Content column taking the rest, and the RightHeader at its end.
func headerFluentContent(c *ui.Context, l headerLayout, title, subtitle string, start, end func()) {
	t := c.Theme()
	caption := t.Rem(12.0 / 14)
	headerItems(c, l, start)
	ui.Text(c, title).SingleLine().FontSize(caption).TextColor(t.Text).MarginStart(l.titleGap).MarginEnd(8).Shrink(1).MinWidth(0)
	ui.Text(c, subtitle).SingleLine().FontSize(caption).TextColor(t.TextMuted).MarginEnd(16).Shrink(1).MinWidth(0).Hidden(subtitle == "")
	ui.Spacer(c)
	headerItems(c, l, end)
}

// headerFluentButton draws b, an image button in a HeaderBar h DIPs high, as TitleBarBackButtonStyle and TitleBarPaneToggleButtonStyle draw the TitleBar's own buttons (TitleBar_themeresources.xaml): TitleBarBackButtonWidth 40 with a Margin of 2 on every side, so as tall as the bar less 4 (VerticalAlignment Stretch) and 2 inside the header's edges, ControlCornerRadius (4), its 16-DIP glyph centered; SubtleFillColorTransparent, SubtleFillColorSecondary under the pointer and SubtleFillColorTertiary pressed, the glyph TextFillColorPrimary, TextFillColorSecondary pressed. The styles set no FocusVisualMargin, so the system focus visual is drawn at the button's edge (FrameworkElement's default 0).
func headerFluentButton(c *ui.Context, b ui.Element, h float32) {
	t := c.Theme()
	b.Width(40).Height(h-4).Margin(2).Padding(0).Radius(4).Justify(ui.Center).AlignItems(ui.Center).Border(0, ui.Color{}).FocusRing(false)
	hover, down := fluentSubtle(t)
	face, fg := ui.Transparent, t.Text
	switch {
	case b.Pressed():
		face, fg = down, t.TextMuted
	case b.Hovered():
		face = hover
	}
	b.Background(face).TextColor(fg)
	focused := b.FocusVisible()
	b.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if focused {
			fluentFocusRing(p, t, r, 0, 0, 0, 0, 4)
		}
	})
}

// ---- GNOME ----

// headerAdwContent lays out the parts of an AdwHeaderBar: the start box, the AdwWindowTitle centered across the bar, as GtkCenterBox centers its center child while the sides leave it room, and the end box, the boxes 6 pixels from the title. The title's .title is bold and the subtitle's .subtitle smaller (GTK's font-size: smaller, the size divided by 1.2) and .dimmed (55%, 90% with higher contrast), each padded 12 at either end and centered (_header-bar.scss, adw-window-title.ui).
func headerAdwContent(c *ui.Context, l headerLayout, title, subtitle string, start, end func()) {
	t := c.Theme()
	headerAdwSide(c, l, start, l.before, false).MarginEnd(l.titleGap)
	ui.Column(c).Shrink(1).MinWidth(0).PaddingX(12).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
		ui.Text(c, title).SingleLine().FontWeight(700)
		ui.Text(c, subtitle).SingleLine().FontSize(t.Rem(1 / 1.2)).Opacity(pick(c.Preferences().HighContrast, float32(0.55), 0.9)).Hidden(subtitle == "")
	})
	headerAdwSide(c, l, end, l.after, true).MarginStart(l.titleGap).Justify(ui.End)
}

// headerAdwSide builds the start box (or the end box, atEnd) of an AdwHeaderBar: fn's controls and, at the bar's edge, the room the window controls take there, which GTK draws in a GTK 3 header bar of its own, padded 6 at either end (MyGo's linux/titlebar.go): the room less the 6 pixels the bar leaves at its edge and the box's spacing of 6, so that the controls end where that header bar starts, 6 from the window buttons inside it, as they are from an AdwHeaderBar's windowcontrols (in libadwaita 1.10.0's render the send-to button ends at 514 and the close button spans 520 to 554 in a 560-pixel window, where a room of 46 ends Breeze UI's controls at 514). The two boxes share the room the title leaves equally, so that the title is centered across the bar, but neither is narrower than what it holds, as GtkCenterBox moves its center child aside for a side that needs more than half: each box is at least as wide as its controls were in the previous frame, and a frame is asked for while that width changes.
func headerAdwSide(c *ui.Context, l headerLayout, fn func(), room float32, atEnd bool) ui.Element {
	side := ui.Row(c).Grow(1).Basis(0).AlignItems(ui.Center).Gap(l.gap)
	// The window controls' room at the bar's edge, the same side in every frame, so that the siblings keep their places.
	w := max(room-l.start-l.gap, 0)
	var items ui.Element
	side.Children(func() {
		if !atEnd {
			ui.Box(c).Shrink(0).Width(w).Role(ui.RoleNone).Hidden(room == 0)
		}
		items = headerItems(c, l, fn)
		if atEnd {
			ui.Box(c).Shrink(0).Width(w).Role(ui.RoleNone).Hidden(room == 0)
		}
	})
	least := w
	if fn != nil {
		least += items.LayoutBounds().W
		if room > 0 {
			least += l.gap
		}
	}
	last := ui.Local(side, "breeze-header-side", func() float32 { return -1 })
	if *last != least {
		*last = least
		c.Invalidate()
	}
	return side.MinWidth(least)
}

// adwHeaderBarBackground is headerbar_bg_color: #ffffff in the light theme and #2e2e32 in the dark one (_colors.scss).
func adwHeaderBarBackground(t *ui.Theme) ui.Color {
	return pick(t.Dark, ui.Hex("#ffffff"), ui.Hex("#2e2e32"))
}

// adwHeaderBar paints an AdwHeaderBar outside a window's titlebar slot and outside an AdwToolbarView (where it is flat, 46 pixels high): headerbar_bg_color and its box-shadow's inset 0 -1px of headerbar_shade_color along the bottom, RGB(0 0 6 / 12%) in the light theme and RGB(0 0 6 / 36%) in the dark one, and with higher contrast $headerbar_border_color, headerbar_border_color at --border-opacity (50%). Its side insets of $headerbar_border_color fall outside the window, under the bar's margins of -1px.
func adwHeaderBar(c *ui.Context, bar ui.Element) {
	t := c.Theme()
	bg := adwHeaderBarBackground(t)
	shade := bg.Mix(adwInk, pick(t.Dark, float32(0.12), 0.36))
	if c.Preferences().HighContrast {
		// headerbar_border_color is RGB(0 0 6 / 80%) in the light theme and white in the dark one.
		shade = pick(t.Dark, bg.Mix(adwInk, 0.8*0.5), bg.Mix(white, 0.5))
	}
	bar.Background(bg)
	bar.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 1, W: r.W, H: 1}, shade, 0)
	})
}
