package breeze

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// macOS 26 and later. Colors are AppKit's semantic colors (NSColor) resolved in the Aqua and Dark Aqua appearances on macOS 27; control fills and sizes are measured from AppKit's own controls rendered at 2x by tools/appkit-reference/reference.swift. A push button's face is secondarySystemFill and an unchecked box's or an off switch's systemFill, the same share of black or white over the window in both appearances.

// macAccent is controlAccentColor, the blue accent.
var macAccent = ui.Hex("#007aff")

func macTheme(dark bool, accent ui.Color) *ui.Theme {
	if accent.A == 0 {
		accent = macAccent
	}
	t := &ui.Theme{
		Dark:          dark,
		Accent:        accent,
		AccentHover:   accent,
		AccentPressed: accent.Mix(black, 0.18),
		AccentText:    textOn(accent),
		Radius:        6,
		Spacing:       4,
		FontSize:      13,
	}
	if dark {
		t.Background = ui.Hex("#1e1e1e")             // windowBackgroundColor
		t.Text = t.Background.Mix(white, 0.847)      // labelColor
		t.TextMuted = t.Background.Mix(white, 0.549) // secondaryLabelColor
		t.Border = t.Background.Mix(white, 0.098)    // separatorColor
		t.Danger, t.Warning, t.Success = ui.Hex("#ff4245"), ui.Hex("#ff9230"), ui.Hex("#30d158")
		t.Selection = ui.Hex("#3f638b")             // selectedTextBackgroundColor
		t.Focus = accent.Mix(white, 0.1).Alpha(0.5) // keyboardFocusIndicatorColor
		t.Scrollbar = ui.RGBA(255, 255, 255, 0.35)
	} else {
		t.Background = ui.Hex("#ffffff")
		t.Text = t.Background.Mix(black, 0.847)
		t.TextMuted = t.Background.Mix(black, 0.498)
		t.Border = t.Background.Mix(black, 0.098)
		t.Danger, t.Warning, t.Success = ui.Hex("#ff383c"), ui.Hex("#ff8d28"), ui.Hex("#34c759")
		t.Selection = ui.Hex("#b3d7ff")
		t.Focus = accent.Alpha(0.5)
		t.Scrollbar = ui.RGBA(0, 0, 0, 0.32)
	}
	t.Surface = macFace(t.Background, dark)
	t.SurfaceHover = t.Surface // AppKit's push buttons do not change under the pointer
	t.SurfacePressed = macPressed(t.Background, dark)
	return t
}

// Shares of black (Aqua) or white (Dark Aqua) over the window background, measured from AppKit's controls.
const (
	macFaceAlpha       = 0.078 // a push button, secondarySystemFill
	macPressedAlpha    = 0.169 // a push button while pressed
	macTroughAlpha     = 0.098 // an unchecked box or an off switch, systemFill
	macTroughDownAlpha = 0.188 // an unchecked box while pressed
)

// macFace is a push button's face.
func macFace(bg ui.Color, dark bool) ui.Color { return tint(bg, dark, macFaceAlpha) }

// macPressed is a push button's face while pressed.
func macPressed(bg ui.Color, dark bool) ui.Color { return tint(bg, dark, macPressedAlpha) }

// macButton draws an NSButton of bezel style push, regular size: 24 points high, a continuous rounded rectangle without a border, the default button in the accent color.
func macButton(c *ui.Context, b ui.Element, primary bool) {
	t := c.Theme()
	b.MinHeight(24).PaddingX(12).Gap(6).Radius(6).Justify(ui.Center).AlignItems(ui.Center)
	face, fg := macFace(t.Background, t.Dark), t.Text
	if b.Pressed() {
		face = macPressed(t.Background, t.Dark)
	}
	if primary {
		face, fg = t.Accent, t.AccentText
		if b.Pressed() {
			face = t.AccentPressed
		}
	}
	b.Background(face).TextColor(fg)
}

// macCheckbox draws an NSButton check box: a 16-point rounded square without a border, 6 points from its label, filled with the accent and a white check mark when checked.
func macCheckbox(c *ui.Context, row, box ui.Element, on bool) {
	t := c.Theme()
	row.Gap(6)
	alpha := float32(macTroughAlpha)
	if row.Pressed() {
		alpha = macTroughDownAlpha
	}
	fill := tint(t.Background, t.Dark, alpha)
	if on {
		fill = t.Accent
		if row.Pressed() {
			fill = t.AccentPressed
		}
	}
	box.Size(16, 16).Radius(4.5).Background(fill)
	box.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if on {
			var path ui.Path
			path.MoveTo(r.X+r.W*0.25, r.Y+r.H*0.53).LineTo(r.X+r.W*0.43, r.Y+r.H*0.72).LineTo(r.X+r.W*0.76, r.Y+r.H*0.28)
			p.StrokePath(&path, 2, t.AccentText)
		}
		if row.FocusVisible() {
			p.FocusRing(r, radii(4.5))
		}
	})
}

// macSwitch draws an NSSwitch of macOS 26: a 54 by 24 point track with a 32 by 20 point capsule knob, 2 points in from the track's edge.
func macSwitch(c *ui.Context, sw ui.Element, on bool) {
	t := c.Theme()
	pos := knob(sw, on, 200*time.Millisecond)
	at := along(c, pos)
	off := tint(t.Background, t.Dark, macTroughAlpha)
	knobFill := white
	if t.Dark {
		knobFill = ui.Hex("#e1e1e1")
	}
	sw.Size(54, 24).Radius(12).Background(off.Mix(t.Accent, pos))
	sw.Draw(func(p *ui.Painter, r ui.Rect) {
		const in, w = 2, 32
		k := ui.Rect{X: r.X + in + at*(r.W-w-2*in), Y: r.Y + in, W: w, H: r.H - 2*in}
		p.Shadow(k, k.H/2, 0, 0.5, 1.5, 0, ui.RGBA(0, 0, 0, 0.18))
		p.Fill(k, knobFill, k.H/2)
	})
}

// ---- more of AppKit's controls, measured as the ones above ----

const (
	macToggleOnAlpha   = 0.125 // a push-on-push-off button that is on
	macSegmentOnAlpha  = 0.196 // the chosen segment of an NSSegmentedControl
	macSeparatorAlpha  = 0.294 // the line between two segments that are not chosen
	macFieldBorder     = 0.149 // an NSTextField's bezel in Aqua, black
	macFieldBorderDark = 0.093 // and in Dark Aqua, white
	macFieldDarkFill   = 0.267 // a text field's fill in Dark Aqua: black over the window
	macViewBorder      = 0.102 // an NSScrollView's bezel border around a text view
)

// macRadio draws an NSButton radio button: a 16-point circle of systemFill, the accent with a white dot of 5 points when chosen, 6 points from its label.
func macRadio(c *ui.Context, row, dot ui.Element, on bool) {
	t := c.Theme()
	row.Gap(6)
	fill := tint(t.Background, t.Dark, macTroughAlpha)
	if row.Pressed() {
		fill = tint(t.Background, t.Dark, macTroughDownAlpha)
	}
	if on {
		fill = t.Accent
		if row.Pressed() {
			fill = t.AccentPressed
		}
	}
	dot.Size(16, 16).Radius(8).Background(fill)
	dot.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if on {
			const d = 5
			p.Fill(ui.Rect{X: r.X + (r.W-d)/2, Y: r.Y + (r.H-d)/2, W: d, H: d}, t.AccentText, d/2)
		}
		if row.FocusVisible() {
			p.FocusRing(r, radii(8))
		}
	})
}

// macToggle draws an NSButton of type push-on-push-off: a push button whose face is darker while it is on.
func macToggle(c *ui.Context, b ui.Element, on bool) {
	macButton(c, b, false)
	t := c.Theme()
	if on {
		face := tint(t.Background, t.Dark, macToggleOnAlpha)
		if b.Pressed() {
			face = macPressed(t.Background, t.Dark)
		}
		b.Background(face)
	}
}

// macField draws the bezel of an NSTextField (a 1-point line around a white field, darker than the window in Dark Aqua, corners of 6; AppKit draws it a point outside the field's 24-point frame, so it shows 26 points high), of an NSSearchField (the same as a capsule), or of an NSScrollView around a text view (square, the window's own background). The focus ring goes around it while it has the focus, as AppKit draws it around the field being edited.
func macField(c *ui.Context, f ui.Element, kind fieldKind) {
	t := c.Theme()
	fill := white
	border := t.Background.Mix(black, macFieldBorder)
	if t.Dark {
		fill = t.Background.Mix(black, macFieldDarkFill)
		border = t.Background.Mix(white, macFieldBorderDark)
	}
	radius := float32(6)
	switch kind {
	case multiLine:
		radius, fill, border = 0, pick(t.Dark, white, t.Background), tint(t.Background, t.Dark, macViewBorder)
		f.Padding(5, 6)
	case searchLine:
		radius = 13
		f.MinHeight(26).Padding(3, 10, 3, 8)
	default:
		f.MinHeight(26).Padding(3, 7)
	}
	f.Radius(radius).Background(fill).Border(0, ui.Color{})
	focused := f.Focused()
	f.DrawOver(func(p *ui.Painter, r ui.Rect) {
		p.Stroke(r, border, radius, 1)
		if focused {
			p.FocusRing(r, radii(radius))
		}
	})
}

// macSelect draws an NSPopUpButton: a push button, 24 points high, showing the choice and the up and down arrows at its end.
func macSelect(c *ui.Context, trigger ui.Element, selected string) {
	macButton(c, trigger, false)
	t := c.Theme()
	trigger.Padding(0, 8, 0, 12).Gap(10).Justify(ui.SpaceBetween).MinWidth(96)
	trigger.Children(func() {
		ui.Text(c, selected).SingleLine()
		ui.Box(c).Size(7, 12).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
			var up, down ui.Path
			up.MoveTo(r.X, r.Y+r.H*0.38).LineTo(r.X+r.W/2, r.Y+r.H*0.08).LineTo(r.X+r.W, r.Y+r.H*0.38)
			down.MoveTo(r.X, r.Y+r.H*0.62).LineTo(r.X+r.W/2, r.Y+r.H*0.92).LineTo(r.X+r.W, r.Y+r.H*0.62)
			p.StrokePath(&up, 1.5, t.Text)
			p.StrokePath(&down, 1.5, t.Text)
		})
	})
}

// macMenu draws the menu a pop-up button opens: a rounded panel with a hairline edge and a soft shadow.
func macMenu(c *ui.Context, panel ui.Element) {
	t := c.Theme()
	bg := pick(t.Dark, ui.Hex("#f6f6f6"), ui.Hex("#2b2b2b"))
	panel.Margin(4, 0, 0, 0).Padding(5).Radius(10).Background(bg).
		Border(1, tint(bg, t.Dark, 0.12)).Shadow(0, 8, 24, 0, ui.RGBA(0, 0, 0, pick(t.Dark, float32(0.18), 0.4)))
}

// macMenuItem draws a menu item: 22 points high, the highlighted one in the accent, a check mark before the choice.
func macMenuItem(c *ui.Context, item ui.Element, label string, chosen bool) {
	t := c.Theme()
	fg := t.Text
	item.MinHeight(22).Padding(0, 10, 0, 6).Gap(4).Radius(6).AlignItems(ui.Center)
	if item.Highlighted() {
		fg = t.AccentText
		item.Background(t.Accent).TextColor(fg)
	}
	item.Children(func() {
		mark := ui.Box(c).Size(12, 12).Shrink(0)
		if chosen {
			mark.Draw(func(p *ui.Painter, r ui.Rect) {
				var path ui.Path
				path.MoveTo(r.X+r.W*0.15, r.Y+r.H*0.55).LineTo(r.X+r.W*0.4, r.Y+r.H*0.8).LineTo(r.X+r.W*0.88, r.Y+r.H*0.2)
				p.StrokePath(&path, 1.6, fg)
			})
		}
		ui.Text(c, label).SingleLine()
	})
}

// macSlider draws an NSSlider: a 6-point track of systemFill, filled with the accent up to a 20 by 16 point capsule knob.
func macSlider(c *ui.Context, s ui.Element, frac float32) {
	t := c.Theme()
	const kw, kh = 20, 16
	// The knob's center travels the content box, inside the padding.
	s.Height(kh).PaddingX(kw / 2).MinWidth(60)
	trough := tint(t.Background, t.Dark, macTroughAlpha)
	// The knob is white in Dark Aqua too: tools/appkit-reference/reference.swift reads 255 across it in both appearances.
	knobFill := white
	at := along(c, frac)
	s.Draw(func(p *ui.Painter, r ui.Rect) {
		track := ui.Rect{X: r.X, Y: r.Y + (r.H-6)/2, W: r.W, H: 6}
		p.Fill(track, trough, 3)
		x := r.X + kw/2 + at*(r.W-kw)
		fill := ui.Rect{X: r.X, Y: track.Y, W: x - r.X, H: 6}
		if c.Direction() == ui.RightToLeft {
			fill = ui.Rect{X: x, Y: track.Y, W: r.X + r.W - x, H: 6}
		}
		p.Fill(fill, t.Accent, 3)
		k := ui.Rect{X: x - kw/2, Y: r.Y + (r.H-kh)/2, W: kw, H: kh}
		p.Shadow(k, kh/2, 0, 0.5, 2, 0, ui.RGBA(0, 0, 0, 0.25))
		p.Fill(k, knobFill, kh/2)
		if s.FocusVisible() {
			p.FocusRing(k, radii(kh/2))
		}
	})
	s.FocusRing(false)
}

// macProgress draws an NSProgressIndicator bar: 8 points high, a pale track with a hairline edge, filled with the accent.
func macProgress(c *ui.Context, e ui.Element, value float64) {
	t := c.Theme()
	e.Height(8).MinWidth(40).Radius(4).Background(tint(t.Background, t.Dark, 0.059)).Border(0.5, tint(t.Background, t.Dark, 0.125)).Clip()
	e.Draw(func(p *ui.Painter, r ui.Rect) { drawBar(c, p, r, value, t.Accent, 4) })
}

// drawBar fills a progress bar's track r up to value, from the start of the reading direction, or a moving segment for a negative value.
func drawBar(c *ui.Context, p *ui.Painter, r ui.Rect, value float64, fill ui.Color, radius float32) {
	rtl := c.Direction() == ui.RightToLeft
	if value >= 0 {
		w := r.W * float32(min(value, 1))
		x := r.X
		if rtl {
			x = r.X + r.W - w
		}
		p.Fill(ui.Rect{X: x, Y: r.Y, W: w, H: r.H}, fill, radius)
		return
	}
	p.AnimationFrame()
	phase := float32(p.Now().UnixMilli()%1500) / 1500
	w := r.W * 0.35
	x := r.X - w + (r.W+w)*phase
	if rtl {
		x = r.X + r.W - (x - r.X) - w
	}
	p.Fill(ui.Rect{X: x, Y: r.Y, W: w, H: r.H}, fill, radius)
}

// macSpinner draws a spinning NSProgressIndicator of regular size: 32 points, 8 round spokes from a quarter of the way out to the edge, fading behind the darkest one as it steps around once a second.
func macSpinner(c *ui.Context, e ui.Element) {
	t := c.Theme()
	ink := pick(t.Dark, black, white)
	e.Size(32, 32)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		const spokes, turn = 8, 1000
		ms := int(p.Now().UnixMilli() % turn)
		lead := ms * spokes / turn
		p.After(time.Duration((lead+1)*turn/spokes-ms) * time.Millisecond)
		cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, r.W/2
		w := r.W * 0.09
		in, out := rad*0.375+w/2, rad-w/2
		for i := range spokes {
			age := (lead - i + spokes) % spokes
			alpha := 0.275 - float32(age)/spokes*0.21
			a := float64(i) * 2 * math.Pi / spokes
			sin, cos := float32(math.Sin(a)), float32(math.Cos(a))
			var path ui.Path
			path.MoveTo(cx+sin*in, cy-cos*in).LineTo(cx+sin*out, cy-cos*out)
			p.StrokePath(&path, w, t.Background.Mix(ink, alpha))
		}
	})
}

// macSegmentedTrack draws an NSSegmentedControl's track: a push button's face, 24 points high, corners of 6.
func macSegmentedTrack(c *ui.Context, track ui.Element) {
	t := c.Theme()
	track.MinHeight(24).Radius(6).Background(macFace(t.Background, t.Dark)).AlignItems(ui.Stretch)
}

// macSegment draws a segment: 12 points either side of its label, the chosen one darker, a line between two that are not chosen.
func macSegment(c *ui.Context, seg ui.Element, chosen, separator bool) {
	t := c.Theme()
	seg.Padding(0, 12).Radius(6).Justify(ui.Center).AlignItems(ui.Center).TextColor(t.Text)
	switch {
	case chosen:
		seg.Background(tint(t.Background, t.Dark, macSegmentOnAlpha))
	case seg.Pressed():
		seg.Background(macPressed(t.Background, t.Dark))
	}
	if separator {
		line := tint(t.Background, t.Dark, macSeparatorAlpha)
		seg.Draw(func(p *ui.Painter, r ui.Rect) {
			x := r.X
			if c.Direction() == ui.RightToLeft {
				x = r.X + r.W - 1
			}
			p.Fill(ui.Rect{X: x, Y: r.Y + (r.H-13)/2, W: 1, H: 13}, line, 0)
		})
	}
}

// macTabList draws the tabs of an NSTabView, which macOS draws as a segmented control, centered over the page.
func macTabList(c *ui.Context, list ui.Element) {
	macSegmentedTrack(c, list)
	list.AlignSelf(ui.Center)
}

// ---- dialogs: NSAlert of macOS 26, measured as the controls above ----

// macSheet is an alert's background: 251 gray in Aqua, 28 in Dark Aqua.
func macSheet(dark bool) ui.Color { return pick(dark, ui.Hex("#fbfbfb"), ui.Hex("#1c1c1c")) }

// macDialog draws a dialog as macOS draws a sheet or an alert's window: corners of 24, a window's shadow, and nothing dimming the window behind, as AppKit leaves it.
func macDialog(c *ui.Context, back, panel ui.Element) {
	t := c.Theme()
	back.Background(ui.Transparent)
	panel.Padding(20).Gap(12).Radius(24).Background(macSheet(t.Dark)).Border(1, tint(macSheet(t.Dark), t.Dark, 0.1)).
		Shadow(0, 12, 36, 0, ui.RGBA(0, 0, 0, pick(t.Dark, float32(0.22), 0.5)))
}

// macAlert draws an NSAlert: 260 points wide, its bold title and its message centered, 10 points apart, above buttons 28 points high with round ends, 16 points in from the edges. Two buttons share a row, the default last; more stack, the default on top.
func macAlert(c *ui.Context, back, panel ui.Element, title, message string, n int, button func(int) ui.Element) {
	t := c.Theme()
	macDialog(c, back, panel)
	panel.Width(260).Padding(20, 16, 16, 16).Gap(0).AlignItems(ui.Stretch)
	panel.Children(func() {
		ui.Text(c, title).Bold().TextAlign(ui.Center).MarginX(4)
		if message != "" {
			ui.Text(c, message).TextAlign(ui.Center).MarginX(4).Margin(10, 4, 0, 4).TextColor(t.Text)
		}
		big := func(b ui.Element) { b.MinHeight(28).Radius(14).Grow(1).Basis(0) }
		if n <= 2 {
			ui.Row(c).Gap(8).Margin(16, 0, 0, 0).Children(func() {
				for i := range n {
					big(button(i))
				}
			})
			return
		}
		ui.Column(c).Gap(8).Margin(16, 0, 0, 0).AlignItems(ui.Stretch).Children(func() {
			big(button(n - 1))
			for i := range n - 1 {
				big(button(i))
			}
		})
	})
}

// macPopover draws an NSPopover's panel as macOS draws its menus and popovers: a light panel with a hairline edge and a soft shadow, corners of 10. NSPopover does not render offscreen, so this follows the menu above rather than a measurement.
func macPopover(c *ui.Context, panel ui.Element) {
	macMenu(c, panel)
	panel.Padding(12).Gap(8)
}
