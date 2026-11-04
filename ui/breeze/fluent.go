package breeze

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// Windows 11: WinUI 3's Fluent 2 styles, from the theme dictionaries of microsoft/microsoft-ui-xaml (Common_themeresources_any.xaml and the Button, CheckBox and ToggleSwitch theme resources). Their colors are translucent layers over the window's SolidBackgroundFillColorBase, resolved here to opaque ones.

// fluentAccent is Windows's default accent color, SystemAccentColor.
var fluentAccent = ui.Hex("#0078d4")

// fluentAccentFill returns AccentFillColorDefault for the system accent: SystemAccentColorDark1 in the light theme and SystemAccentColorLight2 in the dark one. Windows computes those shades; these mixes give its default accent's exactly (#005a9e and about #76b9ed).
func fluentAccentFill(accent ui.Color, dark bool) ui.Color {
	if dark {
		return accent.Mix(white, 0.48)
	}
	return accent.Mix(black, 0.25)
}

func fluentTheme(dark bool, accent ui.Color) *ui.Theme {
	if accent.A == 0 {
		accent = fluentAccent
	}
	t := pick(dark, ui.LightTheme(), ui.DarkTheme())
	fill := fluentAccentFill(accent, dark)
	t.Accent = fill
	t.AccentHover = t.Background.Mix(fill, 0.9) // AccentFillColorSecondary: the fill at 90%
	t.AccentPressed = t.Background.Mix(fill, 0.8)
	t.Radius, t.Spacing, t.FontSize = 4, 4, 14 // ControlCornerRadius, ControlContentThemeFontSize
	if dark {
		t.Background = ui.Hex("#202020") // SolidBackgroundFillColorBase
		t.Text = white                   // TextFillColorPrimary
		t.TextMuted = t.Background.Mix(white, 0.773)
		t.Border = t.Background.Mix(white, 0.071) // ControlStrokeColorDefault
		t.Surface = t.Background.Mix(white, 0.059)
		t.SurfaceHover = t.Background.Mix(white, 0.082)
		t.SurfacePressed = t.Background.Mix(white, 0.031)
		t.AccentText = black // TextOnAccentFillColorPrimary
		t.Focus = white      // FocusStrokeColorOuter
	} else {
		t.Background = ui.Hex("#f3f3f3")
		t.Text = t.Background.Mix(black, 0.894)
		t.TextMuted = t.Background.Mix(black, 0.62)
		t.Border = t.Background.Mix(black, 0.059)
		t.Surface = t.Background.Mix(white, 0.702)
		t.SurfaceHover = t.Background.Mix(ui.Hex("#f9f9f9"), 0.5)
		t.SurfacePressed = t.Background.Mix(ui.Hex("#f9f9f9"), 0.302)
		t.AccentText = white
		t.Focus = t.Background.Mix(black, 0.894)
	}
	t.Selection = accent.Alpha(0.4)
	return t
}

// fluentElevation is ControlElevationBorderBrush for a control h DIPs high: ControlStrokeColorDefault, and ControlStrokeColorSecondary along the bottom edge in the light theme and along the top edge in the dark one, blending over 3 DIPs.
func fluentElevation(t *ui.Theme, h float32) ui.LinearGradient {
	if t.Dark {
		def, sec := t.Background.Mix(white, 0.071), t.Background.Mix(white, 0.094)
		return ui.LinearGradient{Angle: 180, Stops: []ui.GradientStop{{Color: sec, At: 0}, {Color: sec, At: 1 / h}, {Color: def, At: 3 / h}, {Color: def, At: 1}}}
	}
	def, sec := t.Background.Mix(black, 0.059), t.Background.Mix(black, 0.161)
	return ui.LinearGradient{Angle: 180, Stops: []ui.GradientStop{{Color: def, At: 0}, {Color: def, At: (h - 3) / h}, {Color: sec, At: (h - 1) / h}, {Color: sec, At: 1}}}
}

// fluentButton draws DefaultButtonStyle, or AccentButtonStyle when primary: 32 DIPs high with ButtonPadding 11,5,11,6 inside a 1-DIP elevation border, corners of 4.
func fluentButton(c *ui.Context, b ui.Element, primary bool) {
	t := c.Theme()
	b.MinHeight(32).Padding(5, 11, 6, 11).Gap(8).Radius(4).Justify(ui.Center).AlignItems(ui.Center)
	pressed, hovered := b.Pressed(), b.Hovered()
	if primary {
		face, fg := t.Accent, t.AccentText
		switch {
		case pressed:
			face, fg = t.AccentPressed, fg.Alpha(pick(t.Dark, float32(0.7), 0.5))
		case hovered:
			face = t.AccentHover
		}
		b.Background(face).TextColor(fg)
		// The border stays while pressed, transparent, so the label does not move.
		b.Border(1, ui.Color{})
		if !pressed {
			// AccentControlElevationBorderBrush: ControlStrokeColorOnAccentSecondary along the bottom.
			sec := face.Mix(black, pick(t.Dark, float32(0.4), 0.14))
			on := face.Mix(white, 0.08)
			b.BorderLinearGradient(ui.LinearGradient{Angle: 180, Stops: []ui.GradientStop{{Color: on, At: 0}, {Color: on, At: 29.0 / 32}, {Color: sec, At: 31.0 / 32}, {Color: sec, At: 1}}})
		}
		return
	}
	face, fg := t.Surface, t.Text
	switch {
	case pressed:
		face, fg = t.SurfacePressed, t.TextMuted
	case hovered:
		face = t.SurfaceHover
	}
	b.Background(face).TextColor(fg)
	if pressed {
		b.Border(1, t.Border)
	} else {
		b.Border(1, ui.Color{}).BorderLinearGradient(fluentElevation(t, 32))
	}
}

// fluentUnchecked returns the fill and the stroke of an unchecked check box, or of a switch that is off: CheckBoxCheckBackgroundFillUnchecked and CheckBoxCheckBackgroundStrokeUnchecked (ToggleSwitchFillOff and ToggleSwitchStrokeOff are the same).
func fluentUnchecked(t *ui.Theme, hovered, pressed bool) (fill, stroke ui.Color) {
	bg := t.Background
	if t.Dark {
		fill, stroke = bg.Mix(black, 0.098), bg.Mix(white, 0.545)
		switch {
		case pressed:
			fill = bg.Mix(white, 0.071)
		case hovered:
			fill = bg.Mix(white, 0.043)
		}
		return fill, stroke
	}
	fill, stroke = bg.Mix(black, 0.024), bg.Mix(black, 0.447)
	switch {
	case pressed:
		fill = bg.Mix(black, 0.094)
	case hovered:
		fill = bg.Mix(black, 0.059)
	}
	return fill, stroke
}

// fluentCheckbox draws a CheckBox: a 20-DIP box with corners of 4 and a 1-DIP stroke, 8 DIPs from its label, in a row at least 32 high; checked, the accent fill with the check mark of Segoe Fluent Icons (U+E73E) in TextOnAccentFillColorPrimary.
func fluentCheckbox(c *ui.Context, row, box ui.Element, on bool) {
	t := c.Theme()
	row.Gap(8).MinHeight(32)
	hovered, pressed := row.Hovered(), row.Pressed()
	fill, stroke := fluentUnchecked(t, hovered, pressed)
	mark := t.AccentText
	if on {
		fill = t.Accent
		switch {
		case pressed:
			fill, mark = t.AccentPressed, mark.Alpha(pick(t.Dark, float32(0.7), 0.5))
		case hovered:
			fill = t.AccentHover
		}
		stroke = fill
	} else if pressed {
		stroke = t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.216), 0.157))
	}
	box.Size(20, 20).Radius(4).Background(fill).Border(1, stroke)
	box.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if on {
			var path ui.Path
			path.MoveTo(r.X+r.W*0.27, r.Y+r.H*0.52).LineTo(r.X+r.W*0.43, r.Y+r.H*0.68).LineTo(r.X+r.W*0.74, r.Y+r.H*0.36)
			p.StrokePath(&path, 1.5, mark)
		}
		if row.FocusVisible() {
			p.FocusRing(ui.Rect{X: r.X - 2, Y: r.Y - 2, W: r.W + 4, H: r.H + 4}, radii(6))
		}
	})
}

// fluentSwitch draws a ToggleSwitch: a 40 by 20 track, a 1-DIP stroke while off, and a round knob of 12 DIPs that grows to 14 under the pointer and stretches to 17 wide while pressed, moving in 83 milliseconds.
func fluentSwitch(c *ui.Context, sw ui.Element, on bool) {
	t := c.Theme()
	hovered, pressed := sw.Hovered(), sw.Pressed()
	pos := knob(sw, on, 83*time.Millisecond)
	at := along(c, pos)
	offFill, offStroke := fluentUnchecked(t, hovered, pressed)
	onFill := t.Accent
	switch {
	case pressed:
		onFill = t.AccentPressed
	case hovered:
		onFill = t.AccentHover
	}
	knobFill := fluentOffKnob(t).Mix(t.AccentText, pos)
	sw.Size(40, 20).Radius(10).Background(offFill.Mix(onFill, pos)).Border(1, offStroke.Mix(onFill, pos))
	w, h := float32(12), float32(12)
	switch {
	case pressed:
		w, h = 17, 14
	case hovered:
		w, h = 14, 14
	}
	sw.Draw(func(p *ui.Painter, r ui.Rect) {
		// The knob's center goes from 10 to 30 DIPs; pressed, it keeps its outer edge and stretches toward the other end.
		cx := r.X + r.H/2 + at*(r.W-r.H)
		x := cx - h/2 - at*(w-h)
		p.Fill(ui.Rect{X: x, Y: r.Y + (r.H-h)/2, W: w, H: h}, knobFill, h/2)
	})
}

// fluentOffKnob is ToggleSwitchKnobFillOff: TextFillColorSecondary, a gray knob on the track.
func fluentOffKnob(t *ui.Theme) ui.Color {
	return t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.62), 0.773))
}

// ---- more of WinUI 3's controls, from the same theme resources ----

// fluentSubtle returns SubtleFillColorSecondary and SubtleFillColorTertiary, the backgrounds of a list item under the pointer and pressed.
func fluentSubtle(t *ui.Theme) (hover, pressed ui.Color) {
	if t.Dark {
		return t.Background.Mix(white, 0.059), t.Background.Mix(white, 0.039)
	}
	return t.Background.Mix(black, 0.035), t.Background.Mix(black, 0.024)
}

// fluentRadio draws a RadioButton: a 20-DIP circle with a 1-DIP stroke, 8 DIPs from its label in a row at least 32 high; chosen, the accent with a dot of TextOnAccentFillColorPrimary, 12 DIPs across, 14 under the pointer and 10 pressed.
func fluentRadio(c *ui.Context, row, dot ui.Element, on bool) {
	t := c.Theme()
	row.Gap(8).MinHeight(32)
	hovered, pressed := row.Hovered(), row.Pressed()
	fill, stroke := fluentUnchecked(t, hovered, pressed)
	if pressed && !on {
		stroke = t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.216), 0.157))
	}
	d := float32(12)
	if on {
		fill = t.Accent
		switch {
		case pressed:
			fill, d = t.AccentPressed, 10
		case hovered:
			fill, d = t.AccentHover, 14
		}
		stroke = fill
	}
	dot.Size(20, 20).Radius(10).Background(fill).Border(1, stroke)
	dot.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if on {
			p.Fill(ui.Rect{X: r.X + (r.W-d)/2, Y: r.Y + (r.H-d)/2, W: d, H: d}, t.AccentText, d/2)
		}
		if row.FocusVisible() {
			p.FocusRing(ui.Rect{X: r.X - 2, Y: r.Y - 2, W: r.W + 4, H: r.H + 4}, radii(12))
		}
	})
}

// fluentToggle draws a ToggleButton: a button, and while on the accent button's look, as ToggleButtonBackgroundChecked is AccentFillColorDefault.
func fluentToggle(c *ui.Context, b ui.Element, on bool) { fluentButton(c, b, on) }

// fluentField draws a TextBox: 32 DIPs high with padding 10,5,6,6, ControlFillColorDefault, a 1-DIP ControlStrokeColorDefault border whose bottom edge is ControlStrongStrokeColorDefault, and while it has the focus a white fill (darker in the dark theme) above a 2-DIP accent line. A multi-line one has the same box; a search box is an AutoSuggestBox, drawn as a TextBox.
func fluentField(c *ui.Context, f ui.Element, kind fieldKind) {
	t := c.Theme()
	focused := f.Focused()
	fill := t.Surface
	switch {
	case focused:
		fill = pick(t.Dark, white, t.Background.Mix(ui.Hex("#1e1e1e"), 0.7))
	case f.Hovered():
		fill = t.SurfaceHover
	}
	f.Padding(5, 6, 6, 10).Radius(4).Background(fill).Border(0, ui.Color{})
	if kind != multiLine {
		f.MinHeight(32)
	}
	strong := t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.447), 0.545))
	f.DrawOver(func(p *ui.Painter, r ui.Rect) {
		p.Stroke(r, t.Border, 4, 1)
		p.Clip(r, 4, func() {
			if focused {
				p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 2, W: r.W, H: 2}, t.Accent, 0)
			} else {
				p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 1, W: r.W, H: 1}, strong, 0)
			}
		})
	})
}

// fluentSelect draws a closed ComboBox: a button's face and elevation border, 32 DIPs high with padding 12,5,0,7, and a 12-DIP chevron 14 DIPs from its end in TextFillColorSecondary.
func fluentSelect(c *ui.Context, trigger ui.Element, selected string) {
	t := c.Theme()
	pressed, hovered := trigger.Pressed(), trigger.Hovered()
	face, fg := t.Surface, t.Text
	switch {
	case pressed:
		face, fg = t.SurfacePressed, t.TextMuted
	case hovered:
		face = t.SurfaceHover
	}
	trigger.MinHeight(32).MinWidth(120).Padding(5, 14, 7, 12).Gap(12).Radius(4).Justify(ui.SpaceBetween).AlignItems(ui.Center).
		Background(face).TextColor(fg)
	if pressed {
		trigger.Border(1, t.Border)
	} else {
		trigger.Border(1, ui.Color{}).BorderLinearGradient(fluentElevation(t, 32))
	}
	trigger.Children(func() {
		ui.Text(c, selected).SingleLine()
		ui.Box(c).Size(12, 12).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
			var path ui.Path
			path.MoveTo(r.X+r.W*0.12, r.Y+r.H*0.35).LineTo(r.X+r.W*0.5, r.Y+r.H*0.7).LineTo(r.X+r.W*0.88, r.Y+r.H*0.35)
			p.StrokePath(&path, 1, t.TextMuted)
		})
	})
}

// fluentMenu draws a ComboBox's drop-down: the in-app acrylic's fallback color, a 1-DIP SurfaceStrokeColorFlyout edge, OverlayCornerRadius and a flyout's shadow.
func fluentMenu(c *ui.Context, panel ui.Element) {
	t := c.Theme()
	panel.Margin(4, 0, 0, 0).Padding(4, 0).Radius(8).Background(pick(t.Dark, ui.Hex("#f9f9f9"), ui.Hex("#2c2c2c"))).
		Border(1, pick(t.Dark, ui.RGBA(0, 0, 0, 0.059), ui.RGBA(0, 0, 0, 0.2))).Shadow(0, 8, 16, 0, ui.RGBA(0, 0, 0, pick(t.Dark, float32(0.14), 0.26)))
}

// fluentMenuItem draws a ComboBoxItem: padding 11,5,11,7 inside a margin of 5,2, corners of 3, a subtle fill under the pointer, and the chosen one marked by a 3 by 16 accent pill at its start.
func fluentMenuItem(c *ui.Context, item ui.Element, label string, chosen bool) {
	t := c.Theme()
	hover, pressed := fluentSubtle(t)
	item.Margin(2, 5).Padding(5, 11, 7, 11).Radius(3).MinWidth(100)
	switch {
	case item.Pressed():
		item.Background(pressed).TextColor(t.TextMuted)
	case item.Highlighted() || chosen:
		item.Background(hover)
	}
	if chosen {
		item.Draw(func(p *ui.Painter, r ui.Rect) {
			x := r.X
			if c.Direction() == ui.RightToLeft {
				x = r.X + r.W - 3
			}
			p.Fill(ui.Rect{X: x, Y: r.Y + (r.H-16)/2, W: 3, H: 16}, t.Accent, 1.5)
		})
	}
	item.Children(func() { ui.Text(c, label).SingleLine() })
}

// fluentSlider draws a Slider: 32 DIPs high, a 4-DIP track of ControlStrongFillColorDefault filled with the accent, and a 22-DIP thumb of ControlSolidFillColorDefault with an elevation border around an accent dot that grows under the pointer and shrinks pressed.
func fluentSlider(c *ui.Context, s ui.Element, frac float32) {
	t := c.Theme()
	const thumb = 22
	s.Height(32).PaddingX(thumb / 2).MinWidth(60).FocusRing(false)
	inner := float32(10.3)
	accent := t.Accent
	switch {
	case s.Pressed():
		inner, accent = 8.5, t.AccentPressed
	case s.Hovered():
		inner, accent = 14, t.AccentHover
	}
	track := t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.447), 0.545))
	face := pick(t.Dark, white, ui.Hex("#454545"))
	at := along(c, frac)
	s.Draw(func(p *ui.Painter, r ui.Rect) {
		ty := r.Y + (r.H-4)/2
		p.Fill(ui.Rect{X: r.X, Y: ty, W: r.W, H: 4}, track, 2)
		x := r.X + thumb/2 + at*(r.W-thumb)
		fill := ui.Rect{X: r.X, Y: ty, W: x - r.X, H: 4}
		if c.Direction() == ui.RightToLeft {
			fill = ui.Rect{X: x, Y: ty, W: r.X + r.W - x, H: 4}
		}
		p.Fill(fill, t.Accent, 2)
		k := ui.Rect{X: x - thumb/2, Y: r.Y + (r.H-thumb)/2, W: thumb, H: thumb}
		p.Fill(k, face.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.161), 0.094)), thumb/2)
		p.Fill(ui.Rect{X: k.X + 1, Y: k.Y + 1, W: thumb - 2, H: thumb - 2}, face, thumb/2-1)
		p.Fill(ui.Rect{X: x - inner/2, Y: k.Y + (thumb-inner)/2, W: inner, H: inner}, accent, inner/2)
		if s.FocusVisible() {
			p.FocusRing(k, radii(thumb/2))
		}
	})
}

// fluentProgress draws a ProgressBar: a 3-DIP accent bar with rounded ends over a 1-DIP line of ControlStrongStrokeColorDefault, which hides while the bar moves for progress of unknown length.
func fluentProgress(c *ui.Context, e ui.Element, value float64) {
	t := c.Theme()
	line := t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.447), 0.545))
	e.Height(3).MinWidth(40)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		if value >= 0 {
			p.Fill(ui.Rect{X: r.X, Y: r.Y + 1, W: r.W, H: 1}, line, 0.5)
		}
		drawBar(c, p, r, value, t.Accent, 1.5)
	})
}

// fluentSpinner draws an indeterminate ProgressRing: 32 DIPs, a 3-DIP accent arc with round ends turning two and a half times every 2 seconds, growing and shrinking as it goes.
func fluentSpinner(c *ui.Context, e ui.Element) {
	t := c.Theme()
	e.Size(32, 32)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		p.AnimationFrame()
		phase := float64(p.Now().UnixMilli()%2000) / 2000
		start := phase * 2.5 * 2 * math.Pi
		sweep := 0.15*math.Pi + 1.35*math.Pi*math.Sin(phase*math.Pi)
		strokeArc(p, r.X+r.W/2, r.Y+r.H/2, 14.5, start, sweep, 3, t.Accent)
	})
}

// strokeArc strokes the arc of radius rad around cx, cy from angle start, in radians clockwise from 12 o'clock, sweeping sweep, with round ends.
func strokeArc(p *ui.Painter, cx, cy, rad float32, start, sweep float64, width float32, col ui.Color) {
	var path ui.Path
	steps := max(int(sweep/(math.Pi/32)), 2)
	for i := 0; i <= steps; i++ {
		a := start + sweep*float64(i)/float64(steps)
		x, y := cx+rad*float32(math.Sin(a)), cy-rad*float32(math.Cos(a))
		if i == 0 {
			path.MoveTo(x, y)
		} else {
			path.LineTo(x, y)
		}
	}
	p.StrokePath(&path, width, col)
}

// fluentSegmentedTrack draws a SelectorBar: items in a row, padded 4 above and below, without a background.
func fluentSegmentedTrack(c *ui.Context, track ui.Element) {
	track.PaddingY(4).AlignItems(ui.Stretch)
}

// fluentSegment draws a SelectorBarItem: padding 12,10,12,7, text in TextFillColorPrimary (Secondary under the pointer, Tertiary pressed) and, chosen, a 16 by 3 accent pill centered below the text.
func fluentSegment(c *ui.Context, seg ui.Element, chosen bool) {
	t := c.Theme()
	fg := t.Text
	switch {
	case seg.Pressed():
		fg = t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.447), 0.529))
	case seg.Hovered():
		fg = t.TextMuted
	}
	seg.Padding(10, 12, 10, 12).Radius(4).Justify(ui.Center).AlignItems(ui.Center).TextColor(fg)
	if chosen {
		seg.Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X + (r.W-16)/2, Y: r.Y + r.H - 4, W: 16, H: 3}, t.Accent, 1.5)
		})
	}
}

// fluentTabList draws tabs as a SelectorBar, which WinUI offers for switching between pages of a view.
func fluentTabList(c *ui.Context, list ui.Element) { fluentSegmentedTrack(c, list) }

// fluentTab draws a tab as a SelectorBarItem.
func fluentTab(c *ui.Context, tab ui.Element, chosen bool) { fluentSegment(c, tab, chosen) }
