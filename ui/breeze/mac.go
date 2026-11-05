package breeze

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// macOS 26 and later, measured on macOS 27.2. Colors are AppKit's semantic colors (NSColor) resolved in sRGB in the Aqua and Dark Aqua appearances; control fills and sizes are measured from AppKit's own controls rendered at 2x into an sRGB bitmap by tools/appkit-reference/reference.swift, in a window drawn as the key window: MyGo cannot tell a view whether its window is active, so Breeze UI draws the active look, the accent controls in the accent. A push button's face is secondarySystemFill and an unchecked box's or an off switch's systemFill, the same share of black or white over the window in both appearances. What AppKit does not draw offscreen, its focus rings and the selected tab, was measured from a live capture of a key window on macOS 27.2, converted from the display's profile to sRGB.

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
		t.Danger, t.Warning, t.Success = ui.Hex("#ff4245"), ui.Hex("#ff9230"), ui.Hex("#30d158")
		t.Selection = ui.Hex("#3f638b") // selectedTextBackgroundColor
		t.Focus = accent.Mix(white, 0.1).Alpha(0.5)
	} else {
		t.Background = ui.Hex("#ffffff")
		t.Text = t.Background.Mix(black, 0.847)
		t.TextMuted = t.Background.Mix(black, 0.498)
		t.Danger, t.Warning, t.Success = ui.Hex("#ff383c"), ui.Hex("#ff8d28"), ui.Hex("#34c759")
		t.Selection = ui.Hex("#b3d7ff")
		t.Focus = accent.Alpha(0.5)
	}
	t.Border = macSeparator(t.Background, dark)
	if accent == macAccent {
		// The focus ring of a key window as it shows on screen, measured from a live capture on macOS 27.2: one layer 3.5 points wide that reads #8db4fc over white and #286590 over the window in Dark Aqua, which this color at half opacity gives. keyboardFocusIndicatorColor resolves to #0067f4 and #1aa9ff at 0.498 in sRGB; the window server blends it in the display's color space, which reads redder once converted. For another accent the ring is that accent at half opacity (inferred).
		t.Focus = pick(dark, ui.RGBA(27, 105, 249, 0.5), ui.RGBA(46, 169, 255, 0.5))
	}
	// An overlay scroller's knob, drawn by NSScroller.drawKnob into an sRGB bitmap on macOS 27.2: black at 50.2% over white, white at 50.2% over the window in Dark Aqua.
	t.Scrollbar = pick(dark, ui.RGBA(0, 0, 0, 0.5), ui.RGBA(255, 255, 255, 0.5))
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

// macSeparator is separatorColor over bg: black at 10% in Aqua and white at 12% in Dark Aqua (NSColor.separatorColor resolved in sRGB on macOS 27.2).
func macSeparator(bg ui.Color, dark bool) ui.Color {
	return tint(bg, dark, pick(dark, float32(0.100), 0.120))
}

// macFocusRing draws AppKit's focus ring around a shape r of corners radius, as a live capture of a key window on macOS 27.2 shows it around a push button, a check box, a radio button, a segment and a switch: 3.5 points of the theme's Focus color, 3 points outside the shape and half a point over its edge. The ring's corners follow the shape's, 3 points out (inferred).
func macFocusRing(p *ui.Painter, r ui.Rect, radius float32) {
	p.Stroke(ui.Rect{X: r.X - 3, Y: r.Y - 3, W: r.W + 6, H: r.H + 6}, p.Theme().Focus, radius+3, 3.5)
}

// macFace is a push button's face.
func macFace(bg ui.Color, dark bool) ui.Color { return tint(bg, dark, macFaceAlpha) }

// macPressed is a push button's face while pressed.
func macPressed(bg ui.Color, dark bool) ui.Color { return tint(bg, dark, macPressedAlpha) }

// macButton draws an NSButton of bezel style push, regular size: 24 points high, a continuous rounded rectangle without a border, the default button in the accent color, and AppKit's focus ring around it.
func macButton(c *ui.Context, b ui.Element, primary bool) {
	t := c.Theme()
	b.MinHeight(24).PaddingX(12).Gap(6).Radius(6).Justify(ui.Center).AlignItems(ui.Center)
	face := macFace(t.Background, t.Dark)
	if b.Pressed() {
		face = macPressed(t.Background, t.Dark)
	}
	// AppKit composites labelColor over the face: 36 on 235 in Aqua and 223 on 48 in Dark Aqua, measured, where t.Text over the window would read 39 and 221.
	fg := tint(face, t.Dark, 0.847)
	if primary {
		face, fg = t.Accent, t.AccentText
		if b.Pressed() {
			face = t.AccentPressed
		}
	}
	b.Background(face).TextColor(fg)
	macRing(b, 6)
}

// macRing gives e AppKit's focus ring in place of MyGo's, around its box with corners of radius, while it has the focus from the keyboard. The focus is read as e is built: a frame painted again without a build, as while a toast moves, cannot ask e.
func macRing(e ui.Element, radius float32) {
	ring := e.FocusVisible()
	e.FocusRing(false).DrawOver(func(p *ui.Painter, r ui.Rect) {
		if ring {
			macFocusRing(p, r, radius)
		}
	})
}

// disabledOnce reports whether e is disabled, which Disabled may make it only after a widget is built: it asks e while e belongs to the frame being painted, and keeps the answer for the frames painted again without a build, which cannot ask it.
func disabledOnce(e ui.Element) func() bool {
	disabled := false
	return func() bool {
		if e.Valid() {
			disabled = e.IsDisabled()
		}
		return disabled
	}
}

// macCheckboxRadius is the corner of a check box: the corner a checked box leaves uncovered in AppKit's render (6.55 square points of its 16-point box) is the one MyGo's continuous corner of 5 leaves (6.41; 4.5 leaves 5.31 and 5.5 leaves 7.56), measured at 2x on macOS 27.2.
const macCheckboxRadius = 5

// macCheckbox draws an NSButton check box: a 16-point rounded square without a border, 6 points from its label, filled with the accent and a white check mark when checked, gray when checked and disabled.
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
	box.Size(16, 16).Radius(macCheckboxRadius)
	if !on {
		box.Background(fill)
	}
	ring, disabled := row.FocusVisible(), disabledOnce(row)
	box.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if on {
			// Disabled is known only once the row is built, so a checked box is filled as it is painted, gray if disabled.
			face, mark := fill, t.AccentText
			if disabled() {
				face, mark = macDisabledOn(t)
			}
			p.Fill(r, face, macCheckboxRadius)
			var path ui.Path
			path.MoveTo(r.X+r.W*0.25, r.Y+r.H*0.53).LineTo(r.X+r.W*0.43, r.Y+r.H*0.72).LineTo(r.X+r.W*0.76, r.Y+r.H*0.28)
			p.StrokePath(&path, 2, mark)
		}
		if ring {
			macFocusRing(p, r, macCheckboxRadius)
		}
	})
}

// macDisabledOn returns the box and the mark of a checked check box or a chosen radio button that is disabled: AppKit draws them gray, not in the accent, reading 236 and 175 in Aqua and 46 and 98 in Dark Aqua on screen (measured, in a key window). MyGo draws a disabled element at half opacity, so these are the colors that give those readings at half opacity over the window: black at 14.9% and 62.7% in Aqua, white at 14.2% and 60.4% in Dark Aqua (computed from the measurements).
func macDisabledOn(t *ui.Theme) (box, mark ui.Color) {
	return tint(t.Background, t.Dark, pick(t.Dark, float32(0.149), 0.142)), tint(t.Background, t.Dark, pick(t.Dark, float32(0.627), 0.604))
}

// macSwitch draws an NSSwitch of macOS 26: a 54 by 24 point track with a 32 by 20 point capsule knob, 2 points in from the track's edge, and AppKit's focus ring around the track.
func macSwitch(c *ui.Context, sw ui.Element, on bool) {
	t := c.Theme()
	pos := knob(sw, on, 200*time.Millisecond)
	at := along(c, pos)
	off := tint(t.Background, t.Dark, macTroughAlpha)
	// The knob is white in Dark Aqua too: AppKit's reads 255 across it, off and on (reference.swift, macOS 27.2).
	knobFill := white
	sw.Size(54, 24).Radius(12).Background(off.Mix(t.Accent, pos))
	macRing(sw, 12)
	sw.Draw(func(p *ui.Painter, r ui.Rect) {
		const in, w = 2, 32
		k := ui.Rect{X: r.X + in + at*(r.W-w-2*in), Y: r.Y + in, W: w, H: r.H - 2*in}
		p.Shadow(k, k.H/2, 0, 0.5, 1.5, 0, ui.RGBA(0, 0, 0, 0.18))
		p.Fill(k, knobFill, k.H/2)
	})
}

// ---- more of AppKit's controls, measured as the ones above ----

const (
	macSeparatorAlpha  = 0.294 // the line between two segments that are not chosen
	macFieldBorder     = 0.149 // an NSTextField's bezel in Aqua, black
	macFieldBorderDark = 0.093 // and in Dark Aqua, white
	macFieldDarkFill   = 0.267 // a text field's fill in Dark Aqua: black over the window
	macViewBorder      = 0.102 // an NSScrollView's bezel border around a text view
)

// macRadio draws an NSButton radio button: a 16-point circle of systemFill, the accent with a white dot of 5 points when chosen, gray when chosen and disabled, as a check box, 6 points from its label.
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
	dot.Size(16, 16).Radius(8)
	if !on {
		dot.Background(fill)
	}
	ring, disabled := row.FocusVisible(), disabledOnce(row)
	dot.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if on {
			// A disabled chosen radio button reads as a disabled checked box does: 236 and 175 in Aqua, 46 and 98 in Dark Aqua (measured in a key window on macOS 27.2).
			face, mark := fill, t.AccentText
			if disabled() {
				face, mark = macDisabledOn(t)
			}
			p.Fill(r, face, r.H/2)
			const d = 5
			p.Fill(ui.Rect{X: r.X + (r.W-d)/2, Y: r.Y + (r.H-d)/2, W: d, H: d}, mark, d/2)
		}
		if ring {
			macFocusRing(p, r, 8)
		}
	})
}

// macToggle draws an NSButton of type push-on-push-off: a push button that is the accent with a white label while it is on, as AppKit draws it in a key window (#0078f9 with a 255 label in Aqua, #067eff in Dark Aqua, measured by reference.swift on macOS 27.2; t.Accent here). Pressed while on, it is the accent pressed (inferred).
func macToggle(c *ui.Context, b ui.Element, on bool) {
	macButton(c, b, false)
	t := c.Theme()
	if on {
		face := t.Accent
		if b.Pressed() {
			face = t.AccentPressed
		}
		b.Background(face).TextColor(t.AccentText)
	}
}

// macField draws the bezel of an NSTextField (a 1-point line around a white field, darker than the window in Dark Aqua, corners of 7; AppKit draws it a point outside the field's 24-point frame, so it shows 26 points high), of an NSSearchField (the same as a capsule), or of an NSScrollView around a text view (square, the window's own background). The focus ring goes around it while it has the focus, as AppKit draws it around the field being edited: a live capture of a key window on macOS 27.2 shows it 3.5 points wide from 2 points outside the bezel to 1.5 inside it, which is macFocusRing around the field's frame, a point inside the bezel (around a text view and a search field, inferred to be the same).
func macField(c *ui.Context, f ui.Element, kind fieldKind) {
	t := c.Theme()
	fill := white
	border := t.Background.Mix(black, macFieldBorder)
	if t.Dark {
		fill = t.Background.Mix(black, macFieldDarkFill)
		border = t.Background.Mix(white, macFieldBorderDark)
	}
	// The bezel's outer corner, compared with AppKit's text field rendered at 2x on macOS 27.2: row by row from its top, AppKit's edge sits 1.2 to 1.5 pixels further in than MyGo's corner of 6 draws it, which a corner of 7 gives, as NSDatePicker's (macDatePickerRadius).
	radius := float32(7)
	switch kind {
	case multiLine:
		radius, fill, border = 0, pick(t.Dark, white, t.Background), tint(t.Background, t.Dark, macViewBorder)
		f.Padding(5, 6)
	case searchLine:
		radius = 13
		// AppKit's magnifier ink starts 8 points into the bezel and the placeholder 26.5, measured from reference.swift's render; MyGo's padding of 8 and gap of 6 put them at 10 and 30.5. MyGo's glyph stays its own, 8 points wide where AppKit's is 9.5.
		f.MinHeight(26).Padding(3, 10, 3, 6).Gap(4)
	default:
		f.MinHeight(26).Padding(3, 7)
	}
	f.Radius(radius).Background(fill).Border(0, ui.Color{})
	focused := f.Focused()
	f.DrawOver(func(p *ui.Painter, r ui.Rect) {
		p.Stroke(r, border, radius, 1)
		if focused {
			macFocusRing(p, ui.Rect{X: r.X + 1, Y: r.Y + 1, W: r.W - 2, H: r.H - 2}, max(radius-1, 0))
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

// macKnobShadow casts the shadow of a slider's knob. AppKit's, measured by reference.swift on macOS 27.2 over white, is soft and reaches further above the knob than below it: 238, 241, 244, 246, 249 and on over the 4.5 points above, 249, 251, 252, 253, 254 below, and 221 beside it over the 230 track. A shadow of black at 9%, 4 points of blur, a point and a quarter up gives that profile within 2 (inferred from it).
func macKnobShadow(p *ui.Painter, k ui.Rect) {
	p.Shadow(k, k.H/2, 0, -1.25, 4, 0, ui.RGBA(0, 0, 0, 0.09))
}

// macSlider draws an NSSlider: a 6-point track of systemFill, filled with the accent up to a 20 by 16 point capsule knob, with AppKit's focus ring around the knob, under it.
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
		if s.FocusVisible() {
			// AppKit draws the knob over the ring: the ring reads flush outside it, 3.5 points wide (live capture).
			macFocusRing(p, k, kh/2)
		}
		macKnobShadow(p, k)
		p.Fill(k, knobFill, kh/2)
	})
	s.FocusRing(false)
}

// macProgress draws an NSProgressIndicator bar: 8 points high, a pale track with a hairline edge, filled with the accent. The track is black at 5.9% in Aqua and white at 7.8% in Dark Aqua, the edge black at 12.5% and white at 11.6% (240 and 223 over white, 48 and 56 over the window, reference.swift on macOS 27.2).
func macProgress(c *ui.Context, e ui.Element, value float64) {
	t := c.Theme()
	track, edge := tint(t.Background, t.Dark, pick(t.Dark, float32(0.059), 0.078)), tint(t.Background, t.Dark, pick(t.Dark, float32(0.125), 0.116))
	e.Height(8).MinWidth(40).Radius(4).Background(track).Border(0.5, edge).Clip()
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

// macSpinner draws a spinning NSProgressIndicator of regular size as a key window draws it: 32 points, 8 spokes 4 points wide with round ends from 6 points out to the edge, the darkest black (white in Dark Aqua) at 54.9% and each one behind it 7.06% lighter, stepping around once a second. Measured by reference.swift on macOS 27.2: 115, 133, 150, 168, 186, 204, 222 and 240 over white, 154 down to 43 over the window in Dark Aqua, and 8 pixels across a spoke at 2x, its sides straight and its ends reaching 6 and 16 points out along its middle; the turn's period is not measured.
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
		w := r.W / 8
		for i := range spokes {
			age := (lead - i + spokes) % spokes
			alpha := 0.549 - float32(age)*0.0706
			// The spoke pointing up, from 6 points out to the edge, turned to its place.
			spoke := ui.Rect{X: cx - w/2, Y: cy - rad, W: w, H: rad * 0.625}
			p.Transform(ui.Rotation(float32(i)*360/spokes), cx, cy, func() {
				p.Fill(spoke, t.Background.Mix(ink, alpha), w/2)
			})
		}
	})
}

// macSegmentedTrack draws an NSSegmentedControl's track: a push button's face, 24 points high, corners of 6.
func macSegmentedTrack(c *ui.Context, track ui.Element) {
	t := c.Theme()
	track.MinHeight(24).Radius(6).Background(macFace(t.Background, t.Dark)).AlignItems(ui.Stretch)
}

// macSegment draws a segment: 12 points either side of its label, a line between two that are not chosen, and the chosen one in the accent with a white label, as a key window draws it whether the control selects one segment or any (#0078f9 with a 255 label in Aqua, #067eff in Dark Aqua, measured by reference.swift on macOS 27.2; t.Accent here). AppKit's focus ring goes around the segment that has the focus.
func macSegment(c *ui.Context, seg ui.Element, chosen, separator bool) {
	t := c.Theme()
	seg.Padding(0, 12).Radius(6).Justify(ui.Center).AlignItems(ui.Center).TextColor(t.Text)
	switch {
	case chosen && seg.Pressed():
		seg.Background(t.AccentPressed).TextColor(t.AccentText)
	case chosen:
		seg.Background(t.Accent).TextColor(t.AccentText)
	case seg.Pressed():
		seg.Background(macPressed(t.Background, t.Dark))
	}
	macRing(seg, 6)
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

// macTabList draws the tabs of an NSTabView, which macOS draws as a segmented control, centered over the page: a live capture of a key window on macOS 27.2 shows a 24-point strip of black at 7.5% (236 over white, white at 7.3% in Dark Aqua), a push button's face.
func macTabList(c *ui.Context, list ui.Element) {
	macSegmentedTrack(c, list)
	list.AlignSelf(ui.Center)
}

// macTab draws a tab of an NSTabView as a live capture of a key window on macOS 27.2 shows it: the chosen one a raised platter 2 points inside the strip, white in Aqua and white at 20% over the window in Dark Aqua (the strip reads 236 and the platter 255 with a 227 edge over white; 52, 81 and a 135 edge in Dark Aqua), not the accent, its label in labelColor; no lines between the others. The platter's corners, concentric with the strip's, and its shadow are inferred from the capture. AppKit's focus ring goes around the tab that has the focus (inferred from the segments').
func macTab(c *ui.Context, tab ui.Element, chosen bool) {
	t := c.Theme()
	tab.Padding(0, 12).Radius(6).Justify(ui.Center).AlignItems(ui.Center).TextColor(t.Text)
	if !chosen && tab.Pressed() {
		tab.Background(macPressed(t.Background, t.Dark))
	}
	macRing(tab, 6)
	if !chosen {
		return
	}
	fill := pick(t.Dark, white, tint(t.Background, true, 0.2))
	edge := tint(fill, t.Dark, pick(t.Dark, float32(0.11), 0.3))
	tab.Draw(func(p *ui.Painter, r ui.Rect) {
		k := ui.Rect{X: r.X + 2, Y: r.Y + 2, W: r.W - 4, H: r.H - 4}
		p.Shadow(k, 4, 0, 0.5, 1.5, 0, ui.RGBA(0, 0, 0, pick(t.Dark, float32(0.12), 0.4)))
		p.Fill(k, fill, 4)
		p.Stroke(k, edge, 4, 0.5)
	})
}

// macStepper draws an NSStepper of regular size: 20 points wide, a push button's face with corners of 6, and across its middle a line 13 points wide, black at 20.4% in Aqua and white at 22.2% in Dark Aqua (203 and 80, measured from AppKit's stepper rendered in a key window at 2x on macOS 27.2), with AppKit's focus ring around it. The arrows are MyGo's.
func macStepper(c *ui.Context, s ui.Element) {
	t := c.Theme()
	line := tint(t.Background, t.Dark, pick(t.Dark, float32(0.204), 0.222))
	s.Width(20).Radius(6).Background(macFace(t.Background, t.Dark)).Border(0, ui.Color{})
	// MyGo's stepper draws its own ring in DrawOver, which this takes the place of.
	ring := s.FocusVisible()
	s.DrawOver(func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X + (r.W-13)/2, Y: r.Y + r.H/2 - 0.5, W: 13, H: 1}, line, 0)
		if ring {
			macFocusRing(p, r, 6)
		}
	})
}

// ---- dialogs: NSAlert of macOS 26, measured as the controls above. On macOS 27 its panel is a material that does not render offscreen, and the window capture used for the focus rings does not reach an alert's or a popover's window, so the colors here stay those measured on macOS 26 ----

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
