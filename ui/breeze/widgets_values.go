package breeze

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// ---- values: sliders with steps, dates, times and colors ----

// StepSlider creates a slider setting *value to lo or a multiple of step from it, up to hi, with a tick mark at each, as the platform draws a slider with tick marks: AppKit's NSSlider with its dots below the track, WinUI's Slider with ticks below it (TickPlacement BottomRight), or GTK's GtkScale with marks after it. Dragging snaps to the ticks, and the arrows move by step. It is ui.StepSlider's counterpart: ui.StepSlider's own slider with Breeze UI's drawing in place of MyGo's, so it keeps going across under Vertical as ui.StepSlider does, which a SliderBase would not.
func StepSlider(c *ui.Context, value *float64, lo, hi, step float64) ui.Element {
	s := ui.StepSlider(c, value, lo, hi, step)
	frac := float32(0)
	if hi > lo {
		frac = float32(min(max((*value-lo)/(hi-lo), 0), 1))
	}
	ticks := stepSliderTicks(lo, hi, step)
	switch PlatformOf(c) {
	case Windows:
		fluentStepSlider(c, s, frac, ticks)
	case Linux:
		adwStepSlider(c, s, frac, ticks)
	default:
		macStepSlider(c, s, frac, ticks)
	}
	return s
}

// stepSliderTicks returns where along the track a StepSlider's tick marks go, from 0 at lo to 1 at hi: at lo and every multiple of step from it up to hi. Like MyGo's StepSlider it shows none for no step, or past 61 ticks, too many to tell apart.
func stepSliderTicks(lo, hi, step float64) []float32 {
	if step <= 0 || hi <= lo {
		return nil
	}
	n := int(math.Floor((hi-lo)/step+1e-9)) + 1
	if n > 61 {
		return nil
	}
	ticks := make([]float32, n)
	for i := range ticks {
		ticks[i] = float32(float64(i) * step / (hi - lo))
	}
	return ticks
}

// RangeSlider creates a slider of two knobs setting *low and *high between lo and hi, with tick marks at the steps when step is above 0. It is package ui's: its knobs are elements inside it, which take the focus, and only MyGo's drawing knows which of them has it, so Breeze UI cannot draw their focus rings itself. Its track and knobs take the platform's accent and colors from Use.
func RangeSlider(c *ui.Context, low, high *float64, lo, hi, step float64) ui.Element {
	return ui.RangeSlider(c, low, high, lo, hi, step)
}

// DecimalInput creates a field editing *value as an exact decimal between lo and hi in steps of step, with buttons stepping it beside it. It is package ui's: its field and buttons are built inside it, out of reach of a style applied to what it returns.
func DecimalInput(c *ui.Context, value *string, lo, hi, step string) ui.Element {
	return ui.DecimalInput(c, value, lo, hi, step)
}

// DateInput creates a field showing *date in the user's regional formats, which opens a calendar below it, as the platform draws a date field: AppKit's NSDatePicker of style textField, WinUI's CalendarDatePicker, or a GNOME entry, as GTK has no date field of its own. The text, the calendar icon, the calendar and the keys are MyGo's. It is ui.DateInput's counterpart.
func DateInput(c *ui.Context, date *time.Time) ui.Element {
	b := ui.DateInput(c, date)
	switch PlatformOf(c) {
	case Windows:
		fluentDateInput(c, b)
	case Linux:
		adwDateInput(c, b)
	default:
		macDatePicker(c, b, false)
		// NSDatePicker is as wide as its text.
		b.MinWidth(0).Padding(2, 4).Gap(6)
	}
	return b
}

// TimeInput creates a field of the time of day of *tm whose hours, minutes and half of the day each take the focus, in the box the platform draws around a time: AppKit's NSDatePicker of style textField, the button of WinUI's TimePicker, or a GNOME entry. The segments and their keys are MyGo's. It is ui.TimeInput's counterpart.
func TimeInput(c *ui.Context, tm *time.Time) ui.Element {
	f := ui.TimeInput(c, tm)
	switch PlatformOf(c) {
	case Windows:
		fluentTimeInput(c, f)
	case Linux:
		adwTimeInput(c, f)
	default:
		macDatePicker(c, f, f.FocusWithin())
		// MyGo's segments are padded 2 at their sides, which brings the digits 4 in, as AppKit's.
		f.Padding(2, 2)
	}
	return f
}

// Calendar creates a month's calendar choosing *date. It is package ui's, as SwiftUI's graphical date picker lays it out: its days and buttons are built inside it, out of reach of a style applied to what it returns. Its colors, corners and text size come from Use.
func Calendar(c *ui.Context, date *time.Time) ui.Element { return ui.Calendar(c, date) }

// ColorWell creates a swatch of *color that opens a color picker below it, as the platform draws a color well: AppKit's NSColorWell, a WinUI DropDownButton showing the color (WinUI has no color well; its ColorPicker is opened from such a button), or GTK's GtkColorDialogButton. The picker it opens is MyGo's. It is ui.ColorWell's counterpart.
func ColorWell(c *ui.Context, color *ui.Color) ui.Element {
	b := ui.ColorWell(c, color)
	now := *color
	switch PlatformOf(c) {
	case Windows:
		fluentColorWell(c, b, now)
	case Linux:
		adwColorWell(c, b, now)
	default:
		macColorWell(c, b, now)
	}
	return b
}

// ColorPicker creates a picker of *color: a square of saturation and brightness, sliders of hue and opacity, the color in hex and swatches. It is package ui's, laid out as AppKit's color panel; WinUI's ColorPicker and GTK's color chooser lay theirs out differently. Its parts are built inside it, so their colors and corners come from Use.
func ColorPicker(c *ui.Context, color *ui.Color) ui.Element { return ui.ColorPicker(c, color) }

// colorWellChecks paints a checkerboard of size squares over r, under a color that is not opaque, as WinUI's ColorPicker and GTK's color swatch show transparency: first the square at r's top left in first, then the next in second, and so on across and down. A transparent color leaves its squares as they are.
func colorWellChecks(p *ui.Painter, r ui.Rect, size float32, first, second ui.Color) {
	for y, row := r.Y, 0; y < r.Y+r.H; y, row = y+size, row+1 {
		for x, col := r.X, 0; x < r.X+r.W; x, col = x+size, col+1 {
			c := first
			if (row+col)%2 == 1 {
				c = second
			}
			if c.A > 0 {
				p.Fill(ui.Rect{X: x, Y: y, W: min(size, r.X+r.W-x), H: min(size, r.Y+r.H-y)}, c, 0)
			}
		}
	}
}

// hideWellSwatch sizes a color well w by h with corners of radius and keeps MyGo's swatch inside it from showing: a transparent border half as wide as the well is high leaves its children no room to paint in, so the platform's swatch, painted over the well, is the only one, also at the half opacity of a disabled well, through which a swatch underneath would show. The well's pointer, keys, focus ring and accessibility stay MyGo's.
func hideWellSwatch(b ui.Element, w, h, radius float32) {
	b.Size(w, h).Padding(0).Radius(radius).Background(ui.Transparent).Border(h/2, ui.Color{}).Clip()
}

// colorWellSwatch paints color over r with corners of radius, over a checkerboard of size squares in first and second (colorWellChecks) when it is not opaque.
func colorWellSwatch(p *ui.Painter, r ui.Rect, color ui.Color, radius, size float32, first, second ui.Color) {
	if color.A < 255 {
		p.Clip(r, radius, func() { colorWellChecks(p, r, size, first, second) })
	}
	p.Fill(r, color, radius)
}

// ---- macOS ----

// AppKit's controls of this group, measured as mac.go's are, from NSSlider, NSDatePicker and NSColorWell rendered at 2x on macOS 27 in Aqua and Dark Aqua.
const (
	macStepTickAlpha       = 0.259 // a slider's tick mark in Aqua, tertiaryLabelColor: black over the window, 189 over white
	macStepTickAlphaDark   = 0.247 // and in Dark Aqua, white: 86 over the window
	macDateBorder          = 0.078 // NSDatePicker's bezel in Aqua, black over the window
	macDateBorderDark      = 0.049 // and in Dark Aqua, white
	macDateFillDark        = 0.233 // its field in Dark Aqua, black over the window
	macWellEdgeShift       = 32    // a color well's swatch edge: each channel 32 darker in Aqua, lighter in Dark Aqua
	macDatePickerRadius    = 7     // the outer edge of NSDatePicker's bezel
	macDatePickerMinHeight = 23    // its 21-point frame with the bezel drawn a point outside it
)

// macStepSlider draws an NSSlider with tick marks below it, as Slider draws one without: a 6-point track of systemFill, filled with the accent up to a 20 by 16 point capsule knob with AppKit's shadow and focus ring, and a round tick of 2 points in tertiaryLabelColor (measured in a key window on macOS 27.2: 189 over white, 86 over the window in Dark Aqua), 3 points below the track at each step, where the knob's center comes to it. The control stays 16 points high, as AppKit's intrinsic size is; the ticks are at its bottom edge. The knob is white in Aqua and in Dark Aqua alike, as NSSlider's measured 255 in both.
func macStepSlider(c *ui.Context, s ui.Element, frac float32, ticks []float32) {
	t := c.Theme()
	const kw, kh = 20, 16
	s.Height(kh).PaddingX(kw / 2).MinWidth(60).FocusRing(false)
	trough := tint(t.Background, t.Dark, macTroughAlpha)
	dot := tint(t.Background, t.Dark, pick(t.Dark, float32(macStepTickAlpha), macStepTickAlphaDark))
	knobFill := white
	rtl := c.Direction() == ui.RightToLeft
	at := along(c, frac)
	focused := s.FocusVisible()
	s.Draw(func(p *ui.Painter, r ui.Rect) {
		track := ui.Rect{X: r.X, Y: r.Y + (r.H-6)/2, W: r.W, H: 6}
		for _, f := range ticks {
			if rtl {
				f = 1 - f
			}
			x := r.X + kw/2 + f*(r.W-kw)
			p.Fill(ui.Rect{X: x - 1, Y: track.Y + 6 + 3, W: 2, H: 2}, dot, 1)
		}
		p.Fill(track, trough, 3)
		x := r.X + kw/2 + at*(r.W-kw)
		fill := ui.Rect{X: r.X, Y: track.Y, W: x - r.X, H: 6}
		if rtl {
			fill = ui.Rect{X: x, Y: track.Y, W: r.X + r.W - x, H: 6}
		}
		p.Fill(fill, t.Accent, 3)
		k := ui.Rect{X: x - kw/2, Y: r.Y + (r.H-kh)/2, W: kw, H: kh}
		if focused {
			macFocusRing(p, k, kh/2)
		}
		macKnobShadow(p, k)
		p.Fill(k, knobFill, kh/2)
	})
}

// macDatePicker draws the bezel of an NSDatePicker of style textField: a 1-point line, lighter than a text field's, around a white field (darker than the window in Dark Aqua) with corners of 7, 23 points high with the line, which AppKit draws a point outside the control's 21-point frame (the bezel, fill and text measured again in a key window on macOS 27.2: 235 and 41, a 23 field in Dark Aqua, labelColor text). AppKit's focus ring goes around it while ringed or while it has the focus from the keyboard, as around a text field (macField; inferred for the date picker).
func macDatePicker(c *ui.Context, e ui.Element, ringed bool) {
	t := c.Theme()
	fill, border := white, t.Background.Mix(black, macDateBorder)
	if t.Dark {
		fill, border = t.Background.Mix(black, macDateFillDark), t.Background.Mix(white, macDateBorderDark)
	}
	e.MinHeight(macDatePickerMinHeight).Radius(macDatePickerRadius).Background(fill).Border(0, ui.Color{}).TextColor(t.Text).FocusRing(false)
	ringed = ringed || e.FocusVisible()
	e.DrawOver(func(p *ui.Painter, r ui.Rect) {
		p.Stroke(r, border, macDatePickerRadius, 1)
		if ringed {
			macFocusRing(p, ui.Rect{X: r.X + 1, Y: r.Y + 1, W: r.W - 2, H: r.H - 2}, macDatePickerRadius-1)
		}
	})
}

// macColorWell draws an NSColorWell of the default style: a 48 by 24 point capsule of systemFill, and inside it, 3 points in, the color as a capsule with an edge of the color 32 a channel darker (lighter in Dark Aqua). A color that is not opaque shows as AppKit shows it, with no checkerboard: split along a line through the well's center that rises a point for every two across, the color over black above it and over white below it, each half with its own edge. Pressed, the capsule darkens as a pressed check box does (inferred).
func macColorWell(c *ui.Context, b ui.Element, color ui.Color) {
	t := c.Theme()
	bezel := tint(t.Background, t.Dark, macTroughAlpha)
	if b.Pressed() {
		bezel = tint(t.Background, t.Dark, macTroughDownAlpha)
	}
	solid := color
	solid.A = 255
	alpha := float32(color.A) / 255
	// Opaque, both halves are the color itself.
	overBlack, overWhite := black.Mix(solid, alpha), white.Mix(solid, alpha)
	hideWellSwatch(b, 48, 24, 12)
	b.DrawOver(func(p *ui.Painter, r ui.Rect) {
		// The bezel fills the well, under the swatch too, as a disabled well at half opacity shows it there.
		p.Fill(r, bezel, r.H/2)
		s := ui.Rect{X: r.X + 3, Y: r.Y + 3, W: r.W - 6, H: r.H - 6}
		in := ui.Rect{X: s.X + 1, Y: s.Y + 1, W: s.W - 2, H: s.H - 2}
		// The edge is a ring around the swatch, not under it, so a disabled well at half opacity shows the bezel through the swatch once, as AppKit's does (measured: 159 over 242 in Aqua).
		swatch := func(fill ui.Color) {
			p.Stroke(s, macWellEdge(fill, t.Dark), s.H/2, 1)
			p.Fill(in, fill, in.H/2)
		}
		swatch(overWhite)
		if color.A == 255 {
			return
		}
		// The half over black, above the line: a clip turned with the line, the swatch drawn turned back. The line rises a point for every two across, atan(1/2), 26.565 degrees.
		cx, cy := r.X+r.W/2, r.Y+r.H/2
		const tilt = -26.565
		p.Transform(ui.Rotation(tilt), cx, cy, func() {
			p.Clip(ui.Rect{X: cx - r.W, Y: cy - r.W, W: 2 * r.W, H: r.W}, 0, func() {
				p.Transform(ui.Rotation(-tilt), cx, cy, func() { swatch(overBlack) })
			})
		})
	})
}

// macWellEdge returns the edge AppKit draws around a color well's swatch: the color with each channel shifted by macWellEdgeShift, darker in Aqua and lighter in Dark Aqua.
func macWellEdge(color ui.Color, dark bool) ui.Color {
	shift := func(v uint8) uint8 {
		if dark {
			return uint8(min(int(v)+macWellEdgeShift, 255))
		}
		return uint8(max(int(v)-macWellEdgeShift, 0))
	}
	return ui.Color{R: shift(color.R), G: shift(color.G), B: shift(color.B), A: 255}
}

// ---- Windows 11 ----

// fluentStepSlider draws a Slider with TickPlacement BottomRight: 32 DIPs high, a 4-DIP track of ControlStrongFillColorDefault filled with the accent, the thumb of Slider (an 18-DIP Thumb whose border reaches 2 beyond it, around an accent dot that grows under the pointer and shrinks pressed), and a TickBar below the track: ticks one device pixel wide (TickBar_Partial.cpp's singlePixelWidthScaled, 1 over the scale) and 4 DIPs high (SliderOutsideTickBarThemeHeight) in SliderTickBarFill (ControlStrongFillColorDefault), 4 below the track. The thumb's center and the ticks travel the track less the Thumb's 18 DIPs, as TickBar arranges them.
func fluentStepSlider(c *ui.Context, s ui.Element, frac float32, ticks []float32) {
	t := c.Theme()
	const thumb, travel = 22, 18
	s.Height(32).PaddingX(travel / 2).MinWidth(60).FocusRing(false)
	inner := float32(10.3)
	accent := t.Accent
	switch {
	case s.Pressed():
		inner, accent = 8.5, t.AccentPressed
	case s.Hovered():
		inner, accent = 14, t.AccentHover
	}
	strong := t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.447), 0.545))
	face := pick(t.Dark, white, ui.Hex("#454545"))
	rtl := c.Direction() == ui.RightToLeft
	at := along(c, frac)
	focused := s.FocusVisible()
	s.Draw(func(p *ui.Painter, r ui.Rect) {
		ty := r.Y + (r.H-4)/2
		px := 1 / p.Scale()
		for _, f := range ticks {
			if rtl {
				f = 1 - f
			}
			x := r.X + travel/2 + f*(r.W-travel)
			p.Fill(ui.Rect{X: x - px/2, Y: ty + 4 + 4, W: px, H: 4}, strong, 0)
		}
		p.Fill(ui.Rect{X: r.X, Y: ty, W: r.W, H: 4}, strong, 2)
		x := r.X + travel/2 + at*(r.W-travel)
		fill := ui.Rect{X: r.X, Y: ty, W: x - r.X, H: 4}
		if rtl {
			fill = ui.Rect{X: x, Y: ty, W: r.X + r.W - x, H: 4}
		}
		p.Fill(fill, t.Accent, 2)
		k := ui.Rect{X: x - thumb/2, Y: r.Y + (r.H-thumb)/2, W: thumb, H: thumb}
		p.Fill(k, face.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.161), 0.094)), thumb/2)
		p.Fill(ui.Rect{X: k.X + 1, Y: k.Y + 1, W: thumb - 2, H: thumb - 2}, face, thumb/2-1)
		p.Fill(ui.Rect{X: x - inner/2, Y: k.Y + (thumb-inner)/2, W: inner, H: inner}, accent, inner/2)
		// The keyboard's focus goes to Slider's FocusBorder, the whole 32-DIP slider with corners of ControlCornerRadius, widened 7 at each end by its FocusVisualMargin of -7,0,-7,0; the ring goes around the thumb only while focus is engaged, which takes a game pad.
		if focused {
			p.FocusRing(ui.Rect{X: r.X - 7, Y: r.Y, W: r.W + 14, H: r.H}, radii(4))
		}
	})
}

// fluentDateTimeFace gives a date or time picker the face of a button: ControlFillColorDefault, Secondary under the pointer and Tertiary pressed, inside ControlElevationBorderBrush (ControlStrokeColorDefault while pressed), corners of 4, at least 32 DIPs high, as CalendarDatePicker and TimePicker draw theirs.
func fluentDateTimeFace(c *ui.Context, e ui.Element, hovered, pressed bool) {
	t := c.Theme()
	face, fg := t.Surface, t.Text
	switch {
	case pressed:
		face, fg = t.SurfacePressed, t.TextMuted // CalendarDatePickerTextForegroundPressed, TextFillColorSecondary
	case hovered:
		face = t.SurfaceHover
	}
	e.MinHeight(32).Radius(4).Background(face).TextColor(fg)
	if pressed {
		e.Border(1, t.Border)
	} else {
		e.Border(1, ui.Color{}).BorderLinearGradient(fluentElevation(t, 32))
	}
}

// fluentDateInput draws a CalendarDatePicker: a picker's face as wide as its text, the date 12 DIPs in with 2 below it, and the calendar glyph centered in a 32-DIP column at the end, where MyGo's calendar icon is placed, in TextFillColorSecondary.
func fluentDateInput(c *ui.Context, b ui.Element) {
	fluentDateTimeFace(c, b, b.Hovered(), b.Pressed())
	// The icon is 14 DIPs wide: 9 from the end centers it in the last 32, and 9 after the text begins that column.
	b.MinWidth(0).Padding(0, 0, 2, 0).PaddingStart(12).PaddingEnd(9).Gap(9)
}

// fluentTimeInput draws the button of a TimePicker around MyGo's segments: a picker's face at least 242 DIPs wide (TimePickerThemeMinWidth), padded 3 above and 6 below as TimePickerHostPadding, with the segments spread evenly across it, as WinUI centers the hours, the minutes and the half of the day each in an equal column of its own. WinUI draws a 1-DIP divider of TimePickerSpacerFill between each two columns, which a style on the field cannot draw, as it cannot know where MyGo's segments end; MyGo's colon stands near the first (in the 12-hour field rendered here, 7 DIPs before it) and nothing in place of the second. Under the pointer it lightens as the button does.
func fluentTimeInput(c *ui.Context, f ui.Element) {
	fluentDateTimeFace(c, f, f.Hovered(), false)
	f.MinWidth(242).Justify(ui.SpaceEvenly).Padding(3, 0, 6, 0)
}

// fluentColorWell draws a DropDownButton showing the color, as WinUI's apps open a ColorPicker from one: a button's face and elevation border, 32 DIPs high; ButtonPadding 11 at the start; a swatch 36 by 20 (MyGo's size) with corners of 4 (ControlCornerRadius) and a 1-DIP ControlStrokeColorDefault edge (inferred), a color that is not opaque over ColorPicker's checkerboard; then the ChevronDownSmall glyph, 8 DIPs in a 12-DIP box 8 after the swatch, in TextFillColorSecondary (Tertiary under the pointer and pressed), and ButtonPadding 11 at the end: 80 DIPs wide.
func fluentColorWell(c *ui.Context, b ui.Element, color ui.Color) {
	t := c.Theme()
	pressed, hovered := b.Pressed(), b.Hovered()
	face := t.Surface
	chevron := t.TextMuted
	switch {
	case pressed:
		face = t.SurfacePressed
	case hovered:
		face = t.SurfaceHover
	}
	if pressed || hovered {
		chevron = t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.447), 0.529)) // TextFillColorTertiary: #72000000, #87FFFFFF in Dark
	}
	// ColorPicker's checkerboard (ColorHelpers.cpp): squares of CheckerSize 4, the first one clear, the next in CheckerColorBrush, SystemListLowColor: #19000000, #19FFFFFF in Dark.
	checker := pick(t.Dark, black, white).Alpha(0x19 / 255.0)
	hideWellSwatch(b, 1+11+36+8+12+11+1, 32, 4)
	rtl := c.Direction() == ui.RightToLeft
	b.DrawOver(func(p *ui.Painter, r ui.Rect) {
		fluentWellBorder(p, r, t, pressed)
		p.Fill(ui.Rect{X: r.X + 1, Y: r.Y + 1, W: r.W - 2, H: r.H - 2}, face, 3)
		x, cx := r.X+12, r.X+r.W-12-6
		if rtl {
			x, cx = r.X+r.W-12-36, r.X+12+6
		}
		s := ui.Rect{X: x, Y: r.Y + (r.H-20)/2, W: 36, H: 20}
		colorWellSwatch(p, s, color, 4, 4, ui.Transparent, checker)
		p.Stroke(s, t.Border, 4, 1)
		cy := r.Y + r.H/2
		var path ui.Path
		path.MoveTo(cx-4, cy-2).LineTo(cx, cy+2).LineTo(cx+4, cy-2)
		p.StrokePath(&path, 1, chevron)
	})
}

// fluentWellBorder strokes a button's 1-DIP border around r: ControlElevationBorderBrush (as fluentElevation lays it, its 3-DIP blend drawn as a step 2 DIPs from the edge, which keeps the strokes from overlapping), or ControlStrokeColorDefault while pressed.
func fluentWellBorder(p *ui.Painter, r ui.Rect, t *ui.Theme, pressed bool) {
	if pressed {
		p.Stroke(r, t.Border, 4, 1)
		return
	}
	def, sec := t.Background.Mix(black, 0.059), t.Background.Mix(black, 0.161)
	split := r.Y + r.H - 2
	top, bottom := def, sec
	if t.Dark {
		def, sec = t.Background.Mix(white, 0.071), t.Background.Mix(white, 0.094)
		split, top, bottom = r.Y+2, sec, def
	}
	p.Clip(ui.Rect{X: r.X, Y: r.Y, W: r.W, H: split - r.Y}, 0, func() { p.Stroke(r, top, 4, 1) })
	p.Clip(ui.Rect{X: r.X, Y: split, W: r.W, H: r.Y + r.H - split}, 0, func() { p.Stroke(r, bottom, 4, 1) })
}

// ---- GNOME ----

// adwStepSlider draws a GtkScale with marks after it (.marks-after): as Slider, a 4-pixel trough of $trough_color filled with the accent, padded 12 but below, where the marks are; the slider a 20-pixel circle whose bottom-right corner is square, turned 45 degrees so it points down at the marks; and the marks, 1-pixel indicators 6 high, 6 below the 20-pixel range, in currentColor at --dim-opacity. The scale stays 44 pixels high: 12 above, 20 of range and 12 of marks.
func adwStepSlider(c *ui.Context, s ui.Element, frac float32, ticks []float32) {
	t := c.Theme()
	const d = 20
	hovered, pressed := s.Hovered(), s.Pressed()
	s.Height(12 + d + 12).PaddingX(12).MinWidth(60).FocusRing(false)
	trough := adwTrough(t, false, false)
	fill := t.Accent
	knobFill := pick(t.Dark, white, t.Background.Mix(white, 0.8))
	if hovered || pressed {
		trough, fill, knobFill = adwTrough(t, true, false), t.AccentHover, white
	}
	rtl := c.Direction() == ui.RightToLeft
	at := along(c, frac)
	focused := s.FocusVisible()
	s.Draw(func(p *ui.Painter, r ui.Rect) {
		cy := r.Y + 12 + d/2
		left, width := r.X+12, r.W-24
		for _, f := range ticks {
			if rtl {
				f = 1 - f
			}
			x := left + f*width
			p.Fill(ui.Rect{X: x - 0.5, Y: cy + d/2 + 6, W: 1, H: 6}, t.TextMuted, 0)
		}
		p.Fill(ui.Rect{X: left, Y: cy - 2, W: width, H: 4}, trough, 2)
		x := left + at*width
		bar := ui.Rect{X: left, Y: cy - 2, W: x - left, H: 4}
		if rtl {
			bar = ui.Rect{X: x, Y: cy - 2, W: left + width - x, H: 4}
		}
		p.Fill(bar, fill, 2)
		// The shadow 1px 1px 4px, turned with the slider, falls straight down; it is cast by the circle alone.
		p.Shadow(ui.Rect{X: x - d/2, Y: cy - d/2, W: d, H: d}, d/2, 0, 1.4, 4, 0, ui.RGBA(0, 0, 6, 0.2))
		p.FillPath(adwMarkedSlider(x, cy, d/2+1), t.Background.Mix(adwInk, 0.1))
		p.FillPath(adwMarkedSlider(x, cy, d/2), knobFill)
		// _scale.scss: focus-ring("> trough > slider", $offset: 0), a 2-pixel outline right outside the slider's edge, over its hairline, following the square corner as Slider's follows the circle.
		if focused {
			p.StrokePath(adwMarkedSlider(x, cy, d/2+1), 2, t.Focus)
		}
	})
}

// adwMarkedSlider returns the outline of a scale's slider with marks after it, of radius rad around cx, cy: a circle but for its square corner, which points straight down, rad times the square root of 2 below the center.
func adwMarkedSlider(cx, cy, rad float32) *ui.Path {
	var path ui.Path
	point := func(a float64) (float32, float32) {
		return cx + rad*float32(math.Cos(a)), cy + rad*float32(math.Sin(a))
	}
	path.MoveTo(cx, cy+rad*math.Sqrt2)
	// The sides of the corner touch the circle at 45 degrees either side of straight down; the circle goes round over the top from one to the other, in three quarters drawn as cubic curves.
	a := math.Pi / 4
	x, y := point(a)
	path.LineTo(x, y)
	for range 3 {
		b := a - math.Pi/2
		k := float32(4.0 / 3 * math.Tan((b-a)/4))
		x0, y0 := point(a)
		x1, y1 := point(b)
		path.CubeTo(x0-k*rad*float32(math.Sin(a)), y0+k*rad*float32(math.Cos(a)),
			x1+k*rad*float32(math.Sin(b)), y1-k*rad*float32(math.Cos(b)), x1, y1)
		a = b
	}
	return path.Close()
}

// adwDateTimeEntry draws a GtkEntry's box around a field MyGo builds of its own: 34 pixels high, corners of 9, currentColor at 10% without a border, and the 2-pixel focus ring inside its edge while ringed, as the entry draws it while the focus is within it. The fill is painted rather than set as the background, so a button underneath does not change it under the pointer: entries stay as they are.
func adwDateTimeEntry(c *ui.Context, e ui.Element, ringed bool) {
	t := c.Theme()
	fill := adwInkAlpha(t, t.Background, 0.10)
	e.MinHeight(34).Radius(9).Background(ui.Transparent).Border(0, ui.Color{}).TextColor(t.Text).FocusRing(false)
	e.Draw(func(p *ui.Painter, r ui.Rect) { p.Fill(r, fill, 9) })
	e.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if ringed {
			p.Stroke(r, t.Focus, 9, 2)
		}
	})
}

// adwDateInput draws a date field as a GtkEntry: GTK has no date field, and GNOME's apps show a date in an entry or a button opening a GtkCalendar. The text 9 pixels in and MyGo's calendar icon 6 after it (border-spacing), at MyGo's width.
func adwDateInput(c *ui.Context, b ui.Element) {
	adwDateTimeEntry(c, b, b.Focused())
	b.Padding(0, 9).Gap(6)
}

// adwTimeInput draws a time field as a GtkEntry, the digits 9 pixels in (MyGo's segments are padded 2 at their sides), with the focus ring while one of its segments has the focus.
func adwTimeInput(c *ui.Context, f ui.Element) {
	adwDateTimeEntry(c, f, f.FocusWithin())
	f.Padding(0, 7)
}

// adwColorWell draws a GtkColorDialogButton: a button (button.color) padded 5 around its swatch, 34 pixels high as a button is at least, the swatch as wide as "Black" in the text (about MyGo's 36, inferred) and as high as the button's content, corners of 4.5, with a 1-pixel inner edge of view_fg_color at 10% on a light color (a color whose intensity is above a half, as gtkcolorswatch.c decides). A color that is not opaque is drawn over GTK's checkerboard, as the button's swatch keeps its alpha.
func adwColorWell(c *ui.Context, b ui.Element, color ui.Color) {
	t := c.Theme()
	face := t.Surface
	switch {
	case b.Pressed():
		face = t.SurfacePressed
	case b.Hovered():
		face = t.SurfaceHover
	}
	light := (0.30*float32(color.R)+0.59*float32(color.G)+0.11*float32(color.B))/255 > 0.5
	// _color-chooser.scss: an inset 1px of view_fg_color at 10%, which is rgb(0 0 6 / 80%) in light (8% in all) and white in dark.
	edge := pick(t.Dark, adwInk, white).Alpha(pick(t.Dark, float32(0.08), 0.10))
	hideWellSwatch(b, 36+10, 34, 9)
	b.DrawOver(func(p *ui.Painter, r ui.Rect) {
		// The face fills the button under the swatch, as GTK paints the button's background behind its child, so no sliver of the window shows between the two corners.
		p.Fill(r, face, 9)
		s := ui.Rect{X: r.X + 5, Y: r.Y + 5, W: r.W - 10, H: r.H - 10}
		// gtkcolorchooser.c's checkered pattern: 10-pixel squares, #a8a8a8 at the top left and #545454 beside it.
		colorWellSwatch(p, s, color, 4.5, 10, ui.Hex("#a8a8a8"), ui.Hex("#545454"))
		// The edge is the overlay's inset shadow, whose corners are of 4 ($button_radius - 5px).
		if light {
			p.Stroke(s, edge, 4, 1)
		}
	})
}
