package breeze

import (
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
