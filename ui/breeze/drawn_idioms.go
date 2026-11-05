package breeze

import (
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"
)

// Idioms that cross widgets: the look a control takes inside a Toolbar, and the popups of the fields and menu buttons drawn as the platform's menus.

// ---- toolbar items ----

// idiomItemKind tells apart the controls of a Toolbar, which the platforms draw differently: AppKit joins adjacent image-only items into one platter, WinUI draws an icon-only AppBarButton wider than its icon, libadwaita turns image buttons flat.
type idiomItemKind uint8

const (
	idiomText    idiomItemKind = iota // a button or a toggle showing a label
	idiomImage                        // a button with an empty label, as the documented Children(icon) use, or a history button
	idiomToggle                       // a toggle, which AppKit draws as a one-item NSToolbarItemGroup that never joins
	idiomMenu                         // a menu button
	idiomPrimary                      // the primary button
	idiomOther                        // a segment of a ToggleGroup, which ends a run of image items
)

// idiomToolbars holds the Toolbar each window is building, which the controls built inside it read, as toggleGroups does for ToggleGroup.
var idiomToolbars sync.Map // *ui.Context -> *idiomToolbar

// idiomToolbar is a Toolbar being built: the kind and element of each control this package built in it, in order, and whether each shares one platter with the one before it (macOS only).
type idiomToolbar struct {
	kinds  []idiomItemKind
	items  []ui.Element
	joined []bool
}

// idiomBuildToolbar builds ui.Toolbar with fn's controls marked as a toolbar's, restoring the toolbar around it, if any, after.
func idiomBuildToolbar(c *ui.Context, fn func()) ui.Element {
	outer, nested := idiomToolbars.Load(c)
	idiomToolbars.Store(c, &idiomToolbar{})
	tb := ui.Toolbar(c, fn)
	if nested {
		idiomToolbars.Store(c, outer)
	} else {
		idiomToolbars.Delete(c)
	}
	return tb
}

// idiomToolbarItem notes e, of kind, as the next control of the Toolbar being built, and returns the toolbar and e's place in it, or nil outside a toolbar.
func idiomToolbarItem(c *ui.Context, e ui.Element, kind idiomItemKind) (*idiomToolbar, int) {
	v, ok := idiomToolbars.Load(c)
	if !ok {
		return nil, 0
	}
	tb := v.(*idiomToolbar)
	tb.kinds = append(tb.kinds, kind)
	tb.items = append(tb.items, e)
	tb.joined = append(tb.joined, false)
	return tb, len(tb.kinds) - 1
}

// idiomMacToolbarGap is the room NSToolbar leaves between two items, which Toolbar's Gap gives on macOS: 8 points, measured on macOS 27.2.
const idiomMacToolbarGap = 8

// idiomMacJoin decides whether image item i of tb shares one platter with the image item before it, as AppKit draws adjacent image-only bordered NSToolbarItems on macOS 27.2 (a live capture of a key window: bold and italic items, 36 and 37 points wide, touching in one 73-point NSToolbarPlatterView, while text items, menu items and item groups never join). Only controls of this package are noted in a toolbar, so two image items noted one after the other may still have something between them, as a Spacer: they join only when the last layout put them side by side, Toolbar's gap apart, or touching once joined. In the first frame nothing is laid out yet, so they join in the next, which is asked for once.
func idiomMacJoin(c *ui.Context, tb *idiomToolbar, i int) bool {
	if i == 0 || tb.kinds[i] != idiomImage || tb.kinds[i-1] != idiomImage {
		return false
	}
	prev, cur := tb.items[i-1].Bounds(), tb.items[i].Bounds()
	if prev.W == 0 || cur.W == 0 {
		asked := ui.Local(tb.items[i], "breeze-toolbar-join", func() bool { return false })
		if !*asked {
			*asked = true
			c.After(0)
		}
		return false
	}
	// The room between them, whichever side the previous one is on (right to left it is on the right).
	gap := max(cur.X-(prev.X+prev.W), prev.X-(cur.X+cur.W))
	same := cur.Y-prev.Y < 1 && prev.Y-cur.Y < 1
	if !same || gap < -0.5 || gap > idiomMacToolbarGap+0.5 {
		return false
	}
	tb.joined[i] = true
	// The negative margin takes back the toolbar's gap, so that the two touch.
	tb.items[i].MarginStart(-idiomMacToolbarGap)
	return true
}

// idiomMacFill is what fills a toolbar item's platter.
type idiomMacFill uint8

const (
	idiomMacRest idiomMacFill = iota
	idiomMacPressed
	idiomMacOn // a selected toggle or a prominent item, in the accent
	idiomMacOnPressed
)

// idiomMacPlatter paints the platter of a bordered NSToolbarItem in r, as a live capture of a key window on macOS 27.2 (BreezeProbeIdioms, converted from Display P3 to sRGB) shows it: a capsule 36 points high, white in Aqua, and in Dark Aqua (69,68,68) with a lighter rim inside its top and bottom edges ((122,120,120), (105,104,104), then (81,80,80) a point in, and (104,103,103), (120,118,117) at the bottom), inside a half-point line outside the frame that is darker at the sides than above and below (black at 3.1% above, 24.3% at the sides and 4.3% below over white, 247, 193 and 244; in Dark Aqua black at 63% at the sides, 15 over the toolbar's 41, and nothing above and below). Pressed, the whole platter is flat black at 7.8% (235) in Aqua and (89,88,88) in Dark Aqua (an NSToolbarButton highlighted), or pressed is given, for the history buttons. A selected toggle (a one-item NSToolbarItemGroup) and a prominent item (NSToolbarItem.style .prominent) fill it with the accent, controlAccentColor, with a lighter rim inside its top and bottom edges ((28,185,255), (8,158,255) half a point in, for the blue accent). Pressed while selected it is the accent pressed (inferred). A disabled item keeps its platter. Items that share a platter (before, after) draw their slice of it, the platter running on under the neighbor's side, clipped to the item; the slice pressed is a pill 28 points high, 4 points inside the platter's edges, as AppKit highlights an item of a joined platter (measured: the italic slice of bold and italic, 253 to 286 points of 217 to 290). Right to left (rtl), the item before is on the right, so the platter runs on to that side.
func idiomMacPlatter(p *ui.Painter, t *ui.Theme, r ui.Rect, fill idiomMacFill, pressed ui.Color, before, after, rtl, disabled bool) {
	const radius = 18
	// left and right are the sides a neighbor shares the platter on.
	left, right := before, after
	if rtl {
		left, right = after, before
	}
	ext := r
	if left {
		ext.X -= radius
		ext.W += radius
	}
	if right {
		ext.W += radius
	}
	clip := ui.Rect{X: r.X, Y: r.Y - 0.5, W: r.W, H: r.H + 1}
	if !left {
		clip.X -= 0.5
		clip.W += 0.5
	}
	if !right {
		clip.W += 0.5
	}
	// MyGo draws a disabled element at half opacity, each fill on its own; AppKit dims a disabled item's label, not its platter (inferred: the platter of a disabled item in the capture is that of an enabled one), so the platter takes the colors that read as the enabled ones at half opacity, as macDisabledOn does, the label alone fading. It is also what keeps a joined platter whole beside a disabled history button.
	shade := func(col ui.Color) ui.Color {
		if !disabled {
			return col
		}
		if col.A < 255 {
			return ui.RGBA(col.R, col.G, col.B, min(float32(col.A)*2/255, 1))
		}
		up := func(c, bg uint8) uint8 { return uint8(min(max(2*int(c)-int(bg), 0), 255)) }
		return ui.RGB(up(col.R, t.Background.R), up(col.G, t.Background.G), up(col.B, t.Background.B))
	}
	grad := func(g ui.LinearGradient) ui.LinearGradient {
		for i := range g.Stops {
			g.Stops[i].Color = shade(g.Stops[i].Color)
		}
		return g
	}
	accent := fill == idiomMacOn || fill == idiomMacOnPressed
	p.Clip(clip, 0, func() {
		// The half-point line outside the frame, as a band, so that nothing of it shows under the face.
		ring := ui.Rect{X: ext.X - 0.5, Y: ext.Y - 0.5, W: ext.W + 1, H: ext.H + 1}
		var band ui.Path
		idiomCapsule(&band, ring, false)
		idiomCapsule(&band, ext, true)
		if t.Dark {
			p.FillPathGradient(&band, grad(idiomMacStops(ring.H, []idiomStop{{0, ui.RGBA(0, 0, 0, 0)}, {ring.H / 2, ui.RGBA(0, 0, 0, 0.63)}, {ring.H, ui.RGBA(0, 0, 0, 0)}})))
		} else {
			p.FillPathGradient(&band, grad(idiomMacStops(ring.H, []idiomStop{{0, ui.RGBA(0, 0, 0, 0.031)}, {ring.H / 2, ui.RGBA(0, 0, 0, 0.243)}, {ring.H, ui.RGBA(0, 0, 0, 0.043)}})))
		}
		switch {
		case accent:
			face := t.Accent
			if fill == idiomMacOnPressed {
				face = t.AccentPressed
			}
			top, near := face.Mix(white, 0.3), face.Mix(white, 0.15)
			if face == macAccent {
				top, near = ui.RGB(28, 185, 255), ui.RGB(8, 158, 255)
			}
			p.FillGradient(ext, grad(idiomMacStops(ext.H, []idiomStop{{0.25, top}, {0.75, near}, {2, face}, {ext.H - 2, face}, {ext.H - 0.75, near}, {ext.H - 0.25, top}})), radius)
		case fill == idiomMacPressed && (before || after):
			p.FillGradient(ext, grad(idiomMacFace(t, ext.H)), radius)
			pill := ui.Rect{X: r.X, Y: r.Y + 4, W: r.W, H: r.H - 8}
			if !left {
				pill.X += 4
				pill.W -= 4
			}
			if !right {
				pill.W -= 4
			}
			p.Fill(pill, shade(pressed), pill.H/2)
		case fill == idiomMacPressed:
			p.Fill(ext, shade(pressed), radius)
		default:
			p.FillGradient(ext, grad(idiomMacFace(t, ext.H)), radius)
		}
	})
}

// idiomCapsule adds to path a capsule filling r, at least as wide as it is high, its ends half circles, clockwise or, with reverse, counterclockwise, so that a capsule inside another drawn the other way leaves a hole (non-zero winding).
func idiomCapsule(path *ui.Path, r ui.Rect, reverse bool) {
	rad := r.H / 2
	k := 0.5523 * rad
	x0, y0, x1, y1 := r.X, r.Y, r.X+r.W, r.Y+r.H
	// Each segment is a cubic from its start through two control points to its end; a straight one has its control points at its ends.
	segs := [][4][2]float32{
		{{x0 + rad, y0}, {x0 + rad, y0}, {x1 - rad, y0}, {x1 - rad, y0}},
		{{x1 - rad, y0}, {x1 - rad + k, y0}, {x1, y0 + rad - k}, {x1, y0 + rad}},
		{{x1, y0 + rad}, {x1, y1 - rad + k}, {x1 - rad + k, y1}, {x1 - rad, y1}},
		{{x1 - rad, y1}, {x1 - rad, y1}, {x0 + rad, y1}, {x0 + rad, y1}},
		{{x0 + rad, y1}, {x0 + rad - k, y1}, {x0, y1 - rad + k}, {x0, y1 - rad}},
		{{x0, y0 + rad}, {x0, y0 + rad - k}, {x0 + rad - k, y0}, {x0 + rad, y0}},
	}
	if reverse {
		path.MoveTo(x0+rad, y0)
		for i := len(segs) - 1; i >= 0; i-- {
			sg := segs[i]
			path.CubeTo(sg[2][0], sg[2][1], sg[1][0], sg[1][1], sg[0][0], sg[0][1])
		}
	} else {
		path.MoveTo(x0+rad, y0)
		for _, sg := range segs {
			path.CubeTo(sg[1][0], sg[1][1], sg[2][0], sg[2][1], sg[3][0], sg[3][1])
		}
	}
	path.Close()
}

// idiomMacPressedFace is the face of a pressed toolbar item: black at 7.8% (235) in Aqua and (89,88,88) in Dark Aqua, measured as idiomMacPlatter says.
func idiomMacPressedFace(t *ui.Theme) ui.Color {
	return pick(t.Dark, ui.RGB(235, 235, 235), ui.RGB(89, 88, 88))
}

// idiomMacFace is the face of a toolbar item at rest, as idiomMacPlatter describes it, for a platter h points high.
func idiomMacFace(t *ui.Theme, h float32) ui.LinearGradient {
	if !t.Dark {
		return ui.LinearGradient{Angle: 180, Stops: []ui.GradientStop{{Color: white, At: 0}, {Color: white, At: 1}}}
	}
	mid := ui.RGB(69, 68, 68)
	return idiomMacStops(h, []idiomStop{
		{0.25, ui.RGB(122, 120, 120)}, {0.75, ui.RGB(105, 104, 104)}, {1.25, ui.RGB(81, 80, 80)}, {3, ui.RGB(74, 73, 73)}, {h / 2, mid},
		{h - 3, ui.RGB(74, 73, 73)}, {h - 1.25, ui.RGB(81, 80, 80)}, {h - 0.75, ui.RGB(104, 103, 103)}, {h - 0.25, ui.RGB(120, 118, 117)},
	})
}

// idiomStop is a color at a distance in points from the top of a platter.
type idiomStop struct {
	at    float32
	color ui.Color
}

// idiomMacStops turns stops in points down a platter h points high into a gradient running from its top to its bottom.
func idiomMacStops(h float32, stops []idiomStop) ui.LinearGradient {
	g := ui.LinearGradient{Angle: 180, Stops: make([]ui.GradientStop, len(stops))}
	for i, s := range stops {
		g.Stops[i] = ui.GradientStop{Color: s.color, At: min(max(s.at/h, 0), 1)}
	}
	return g
}

// idiomMacInk is the label or glyph of a toolbar item: AppKit draws it pure black in Aqua and white in Dark Aqua in a key window (0 and 255 across the label of a text item and the glyph of an image item, measured as idiomMacPlatter says), white over the accent (the selected toggle's label reads 255).
func idiomMacInk(t *ui.Theme, accent bool) ui.Color {
	if accent {
		return t.AccentText
	}
	return pick(t.Dark, black, white)
}

// idiomMacToolbarItem draws b, item i of tb, as a bordered NSToolbarItem of macOS 27 in a window's unified toolbar, measured from a live capture of a key window on macOS 27.2: a 36-point capsule (idiomMacPlatter), a text item padded 12 points either side of its 13-point regular label (the platters of Bold, Italic and Wrap are their labels' advance plus 24, 51.5, 53 and 56 points wide), an image item 36 points square, joined with an image item beside it (idiomMacJoin), and AppKit's focus ring around it. on is a selected toggle's and the primary button's accent fill.
func idiomMacToolbarItem(c *ui.Context, b ui.Element, tb *idiomToolbar, i int, on bool) {
	t := c.Theme()
	kind := tb.kinds[i]
	b.Height(36).Radius(18).Gap(6).Justify(ui.Center).AlignItems(ui.Center).Border(0, ui.Color{})
	if kind == idiomImage {
		b.MinWidth(36).Padding(0)
	} else {
		b.PaddingX(12)
	}
	before := idiomMacJoin(c, tb, i)
	fill := idiomMacRest
	switch {
	case on && b.Pressed():
		fill = idiomMacOnPressed
	case on:
		fill = idiomMacOn
	case b.Pressed():
		fill = idiomMacPressed
	}
	b.TextColor(idiomMacInk(t, on))
	pressed, disabled := idiomMacPressedFace(t), disabledOnce(b)
	rtl := c.Direction() == ui.RightToLeft
	b.Draw(func(p *ui.Painter, r ui.Rect) {
		after := i+1 < len(tb.joined) && tb.joined[i+1]
		idiomMacPlatter(p, t, r, fill, pressed, before, after, rtl, disabled())
	})
	macRing(b, 18)
}

// idiomFluentAppBarButton draws b as an AppBarButton in a CommandBar, or an AppBarToggleButton when toggle, with its label on the right (DefaultLabelPosition Right, which sizes the button to its content), from AppBarButton_themeresources.xaml and AppBarToggleButton_themeresources.xaml of microsoft-ui-xaml main at cb1b8a6: the face is AppBarButtonInnerBorder, 36 DIPs high inside the 48-DIP AppBarThemeCompactHeight row (AppBarButtonInnerBorderMargin 2,6,2,6, which Toolbar's padding and gap give), with ControlCornerRadius (4); the label is AppBarButtonLabelFontSize, 12, padded 18 at the start (AppBarButtonContentViewboxMargin 12 and AppBarButtonTextLabelOnRightMargin 8, less the inner border's 2, an empty icon taking no room: inferred) and 10 at the end (the label's margin 12 less 2); an icon-only one (LabelCollapsed) keeps the style's Width 68, a 64-DIP face around its 16-DIP icon. It is SubtleFillColorTransparent, SubtleFillColorSecondary under the pointer and SubtleFillColorTertiary pressed, with TextFillColorSecondary pressed; checked, AccentFillColorDefault with TextOnAccentFillColorPrimary inside AccentControlElevationBorderBrush, AccentFillColorSecondary under the pointer and AccentFillColorTertiary pressed, each with the HighLightOverlay of SubtleFillColorSecondary or Tertiary laid over it, pressed with TextOnAccentFillColorSecondary and no border. The system focus visual goes 3 DIPs around the face (FocusVisualMargin -3, the inner border being the template's focus target).
func idiomFluentAppBarButton(c *ui.Context, b ui.Element, iconOnly, on bool) {
	t := c.Theme()
	b.Height(36).Radius(4).FontSize(12).Gap(8).Justify(ui.Center).AlignItems(ui.Center).FocusRing(false)
	if iconOnly {
		b.Width(64).Padding(0)
	} else {
		b.PaddingY(0).PaddingStart(18).PaddingEnd(10)
	}
	pressed, hovered := b.Pressed(), b.Hovered()
	hover, down := fluentSubtle(t)
	face, fg := ui.Transparent, t.Text
	switch {
	case pressed:
		face, fg = down, t.TextMuted
	case hovered:
		face = hover
	}
	// The border is always there, transparent but while checked, so that the label does not move.
	b.Border(1, ui.Color{})
	if on {
		ink := pick(t.Dark, black, white)
		face, fg = t.Accent, t.AccentText
		switch {
		case pressed:
			face, fg = t.AccentPressed.Mix(ink, pick(t.Dark, float32(0.024), 0.039)), fg.Alpha(pick(t.Dark, float32(0.7), 0.5))
		case hovered:
			face = t.AccentHover.Mix(ink, pick(t.Dark, float32(0.035), 0.059))
		}
		if !pressed {
			sec := face.Mix(black, pick(t.Dark, float32(0.4), 0.14))
			hi := face.Mix(white, 0.08)
			b.BorderLinearGradient(ui.LinearGradient{Angle: 180, Stops: []ui.GradientStop{{Color: hi, At: 0}, {Color: hi, At: 33.0 / 36}, {Color: sec, At: 35.0 / 36}, {Color: sec, At: 1}}})
		}
	}
	b.Background(face).TextColor(fg)
	focused := b.FocusVisible()
	b.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if focused {
			fluentFocusRing(p, t, r, 3, 3, 3, 3, 4)
		}
	})
}

// idiomAdwFlatButton draws b as libadwaita 1.10 draws an image button, or an image toggle when toggle, in a header bar or a .toolbar (_toolbars.scss: %toolbar's button.image-button extends %button_basic_flat): 34 pixels square (min-height 24 and padding 5), corners of $button_radius (9), transparent, $hover_color (currentColor at 7%) under the pointer and $active_color (16%) pressed; checked, $selected_color (10%), $selected_hover_color (13%) under the pointer and $selected_active_color (19%) pressed (_colors.scss); the focus ring 2 pixels inside its edge. A flat button disabled is at $strong_disabled_opacity (30%), which MyGo's own half opacity takes the place of.
func idiomAdwFlatButton(c *ui.Context, b ui.Element, on bool) {
	b.MinWidth(34).MinHeight(34).Padding(5).Radius(9).Justify(ui.Center).AlignItems(ui.Center).Border(0, ui.Color{})
	idiomAdwFlatFace(c, b, on)
}

// idiomAdwFlatFace fills a flat button: transparent, currentColor at 7% under the pointer and 16% pressed, 10%, 13% and 19% while checked (%button_basic_flat in _buttons.scss and _colors.scss of libadwaita 1.10.0), with its focus ring inside its edge.
func idiomAdwFlatFace(c *ui.Context, b ui.Element, on bool) {
	t := c.Theme()
	share := float32(0)
	switch {
	case b.Pressed():
		share = pick(on, float32(0.16), 0.19)
	case b.Hovered():
		share = pick(on, float32(0.07), 0.13)
	case on:
		share = 0.10
	}
	face := ui.Transparent
	if share > 0 {
		face = adwInkAlpha(t, t.Background, share)
	}
	b.Background(face).TextColor(adwTextOn(t, face))
	adwFocusRing(t, b, 9, nil)
}

// ---- menu buttons ----

// idiomMenuButton builds a menu button of the platform showing label and the platform's arrow, from ui.ButtonBase as ui.MenuButton builds its own: Menu(build) opens the system's menu, or with popover, PopoverMenuButton(build) the menu MyGo draws, which takes the platform's menu colors (idiomMenuTheme).
func idiomMenuButton(c *ui.Context, label string, build func(m *ui.Menu), popover bool) ui.Element {
	b := ui.ButtonBase(c)
	tb, i := idiomToolbarItem(c, b, idiomMenu)
	p := PlatformOf(c)
	styleMenuButton(c, b, label, tb, i)
	if !popover {
		return b.Menu(build)
	}
	// MyGo draws the menu with the theme of the scope PopoverMenuButton is called in. A themed box inside the button, out of the layout and without size, gives it the platform's menu colors without wrapping the button.
	mt := idiomMenuTheme(c, p)
	b.Children(func() {
		ui.Themed(c, mt).Background(ui.Transparent).Absolute().Size(0, 0).Children(func() { b.PopoverMenuButton(build) })
	})
	return b
}

// idiomMenuTheme returns the theme MyGo's drawn menu takes the platform's menu colors from. MyGo's menu takes its panel's background and border from the theme's Background and Border, its corners from Radius+2 (its rows' from Radius), its highlighted row from Accent with AccentText, and the row leading to an open submenu from SurfaceHover; its paddings stay MyGo's, from Spacing.
//   - macOS: macMenu's panel, #f6f6f6 or #2b2b2b with separatorColor's 12% edge and corners of 10, the highlighted row in the accent, as NSMenu draws it.
//   - Windows: a MenuFlyoutPresenter, the in-app acrylic's fallback color (#f9f9f9, #2c2c2c) with SurfaceStrokeColorFlyout's edge and OverlayCornerRadius (8), the row under the pointer SubtleFillColorSecondary in TextFillColorPrimary (MenuFlyoutItemBackgroundPointerOver).
//   - GNOME: popover.menu, popover_bg_color (#ffffff, #36363a) with corners of $popover_radius (15), libadwaita 1.10's 1-pixel ring of black at 5% drawn as MyGo's border, a modelbutton under the pointer or selected currentColor at 10% ($selected_color, _menus.scss).
func idiomMenuTheme(c *ui.Context, p Platform) *ui.Theme {
	t := c.Theme()
	mt := *t
	switch p {
	case Windows:
		bg := pick(t.Dark, ui.Hex("#f9f9f9"), ui.Hex("#2c2c2c"))
		mt.Background, mt.Border, mt.Radius = bg, bg.Mix(black, pick(t.Dark, float32(0.059), 0.2)), 6
		hover := bg.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.035), 0.059))
		mt.Accent, mt.AccentText, mt.SurfaceHover = hover, t.Text, hover
	case Linux:
		bg := pick(t.Dark, white, ui.Hex("#36363a"))
		mt.Background, mt.Border, mt.Radius = bg, bg.Mix(black, 0.05), 13
		hover := adwInkAlpha(t, bg, 0.10)
		mt.Accent, mt.AccentText, mt.SurfaceHover = hover, adwTextOn(t, hover), hover
	default:
		bg := pick(t.Dark, ui.Hex("#f6f6f6"), ui.Hex("#2b2b2b"))
		mt.Background, mt.Border, mt.Radius = bg, tint(bg, t.Dark, 0.12), 8
	}
	return &mt
}

// idiomMacChevron draws the arrow of an NSMenuToolbarItem's pull-down button in r: a downward chevron 7.5 by 4 points of ink, its stroke about 1.5 points (measured from a live capture of a key window on macOS 27.2: the ink of Share's arrow spans 7.5 by 4 points, 4 points after its label's advance and 13 before the platter's end; the stroke is approximate).
func idiomMacChevron(c *ui.Context, color ui.Color) {
	ui.Box(c).Size(7.5, 4).Shrink(0).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
		var path ui.Path
		path.MoveTo(r.X+0.75, r.Y+0.75).LineTo(r.X+r.W/2, r.Y+r.H-0.75).LineTo(r.X+r.W-0.75, r.Y+0.75)
		p.StrokePath(&path, 1.5, color)
	})
}

// idiomFluentSubItemChevron draws AppBarButtonFlyoutGlyph, U+E974 of Segoe Fluent Icons, at AppBarButtonSubItemChevronFontSize 8, the chevron navCrumbSeparator draws at 12 scaled to 8 (the glyph's shape is inferred: the font is not open), pointing left right to left (MirroredWhenRightToLeft).
func idiomFluentSubItemChevron(c *ui.Context, color ui.Color) {
	rtl := c.Direction() == ui.RightToLeft
	ui.Box(c).Size(8, 8).Shrink(0).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
		x0, x1, y0, y1 := r.X+r.W*0.375, r.X+r.W*0.66, r.Y+r.H*0.21, r.Y+r.H*0.79
		if rtl {
			x0, x1 = r.X+r.W-(x0-r.X), r.X+r.W-(x1-r.X)
		}
		var path ui.Path
		path.MoveTo(x0, y0).LineTo(x1, (y0+y1)/2).LineTo(x0, y1)
		p.StrokePath(&path, 1, color)
	})
}

// idiomAdwArrow draws adwaita-icon-theme 51.0's pan-down-symbolic in 16 pixels, as adwSelect does.
func idiomAdwArrow(c *ui.Context, color ui.Color) {
	ui.Box(c).Size(16, 16).Shrink(0).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
		var path ui.Path
		path.MoveTo(r.X+4, r.Y+7).LineTo(r.X+8, r.Y+11).LineTo(r.X+12, r.Y+7)
		p.StrokePath(&path, 2, color)
	})
}

// ---- the popup of an autocomplete ----

// idiomAutocomplete builds an autocomplete on ui.ComboboxBase, as ui.Autocomplete is built, with its popup drawn as the platform's menu. It returns the input, as ui.Autocomplete does, so that Label, Placeholder, Focus, Changed and Submitted keep their meaning; the arrows, Enter, Escape, the focus closing the popup and what assistive technology sees are ComboboxBase's.
func idiomAutocomplete(c *ui.Context, value *string, suggestions []string) ui.Element {
	cb := ui.ComboboxBase(c, value)
	in := cb.Input
	// typed is set while the text is what the user typed since the popup opened, which the suggestions match: ComboboxBase's own (Filtering) is cleared below.
	typed := ui.Local(in, "breeze-typed", func() bool { return false })
	// end is set once a suggestion is taken, for the caret to go after it in the next build pass, once the input shows it.
	end := ui.Local(in, "breeze-caret-end", func() bool { return false })
	if *end {
		n := utf8.RuneCountInString(*value)
		in.SetTextSelection(n, n)
		*end = false
	}
	if in.Changed() {
		*typed = true
		// ComboboxBase highlights the first option as the user types, for Enter to take it; ui.Autocomplete highlights none, so that Enter without the arrows is Submitted. Closing the popup and opening it again leaves it open without a highlight.
		cb.SetOpen(false)
		cb.SetOpen(true)
	}
	if !cb.Open() {
		*typed = false
	}
	shown := idiomSuggestions(suggestions, *value, *typed)
	p := PlatformOf(c)
	t := c.Theme()
	cb.Popup(func(panel ui.Element) {
		// MyGo's own list scrolls past Space(70); WinUI's AutoSuggestListMaxHeight is 374 (generic.xaml).
		height := t.Space(70)
		switch p {
		case Windows:
			// AutoSuggestBox's SuggestionsContainer: AcrylicBackgroundFillColorDefault, SurfaceStrokeColorFlyout, OverlayCornerRadius, the list padded 2 above and below (AutoSuggestListMargin 0,2,0,2).
			fluentMenu(c, panel)
			panel.Padding(2, 0)
			height = 374
		case Linux:
			adwMenu(c, panel)
		default:
			macMenu(c, panel)
		}
		// The scroll area is no more than a frame of the list to assistive technology.
		ui.Scroll(c).MaxHeight(height).Role(ui.RoleNone).Children(func() {
			for _, s := range shown {
				item := cb.Item(s)
				switch p {
				case Windows:
					idiomFluentSuggestion(c, item, s)
				case Linux:
					adwMenuItem(c, item, s, false)
				default:
					macMenuItem(c, item, s, false)
				}
			}
		})
	})
	if v, ok := cb.Chosen(); ok {
		if *value != v {
			*value = v
			in.SetChanged()
		}
		*typed, *end = false, true
	}
	return in
}

// idiomSuggestions returns those of suggestions containing what was typed, ignoring case, those starting with it first, but the one equal to it, as ui.Autocomplete shows them; none until the user typed, or for a blank text.
func idiomSuggestions(suggestions []string, value string, typed bool) []string {
	q := strings.ToLower(strings.TrimSpace(value))
	if !typed || q == "" {
		return nil
	}
	var starts, contains []string
	for _, s := range suggestions {
		if strings.EqualFold(s, value) {
			continue
		}
		switch l := strings.ToLower(s); {
		case strings.HasPrefix(l, q):
			starts = append(starts, s)
		case strings.Contains(l, q):
			contains = append(contains, s)
		}
	}
	return append(starts, contains...)
}

// idiomFluentSuggestion draws a suggestion of an AutoSuggestBox, whose SuggestionsList is a ListView: a ListViewItem as fluentListRow draws one (ListViewItemMinHeight 40, Padding 16,0,12,0 after the selection indicator's 7 DIPs, a backplate 4,2 in with ListViewItemCornerRadius 4), SubtleFillColorSecondary under the pointer and while selected, which the arrows make it, SubtleFillColorTertiary pressed and selected under the pointer, and the selected one's 3 by 16 accent pill (ListViewItem_themeresources.xaml). MyGo highlights the option under the pointer too; the one highlighted while the pointer is elsewhere is taken for the selected one.
func idiomFluentSuggestion(c *ui.Context, item ui.Element, label string) {
	t := c.Theme()
	item.MinHeight(40).PaddingStart(16 + 7).PaddingEnd(12).AlignItems(ui.Center).TextColor(t.Text)
	// The subtle fills are translucent, laid here over the list's own background, fluentMenu's.
	panel := pick(t.Dark, ui.Hex("#f9f9f9"), ui.Hex("#2c2c2c"))
	ink := pick(t.Dark, black, white)
	secondary, tertiary := panel.Mix(ink, pick(t.Dark, float32(0.035), 0.059)), panel.Mix(ink, pick(t.Dark, float32(0.024), 0.039))
	hovered, pressed := item.Hovered(), item.Pressed()
	selected := item.Highlighted() && !hovered
	var fill ui.Color
	switch {
	case pressed:
		fill = tertiary
	case hovered || selected:
		fill = secondary
	}
	rtl := c.Direction() == ui.RightToLeft
	item.Draw(func(p *ui.Painter, r ui.Rect) {
		if fill.A > 0 {
			p.Fill(ui.Rect{X: r.X + 4, Y: r.Y + 2, W: r.W - 8, H: r.H - 4}, fill, 4)
		}
		if !selected {
			return
		}
		x := r.X + 4
		if rtl {
			x = r.X + r.W - 4 - 3
		}
		p.Fill(ui.Rect{X: x, Y: r.Y + (r.H-16)/2, W: 3, H: 16}, t.Accent, 1.5)
	})
	item.Children(func() { ui.Text(c, label).SingleLine() })
}
