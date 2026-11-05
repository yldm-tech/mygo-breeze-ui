package breeze

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// Windows 11: WinUI 3's Fluent 2 styles, from the theme dictionaries of microsoft/microsoft-ui-xaml main at cb1b8a6 (Common_themeresources_any.xaml and the controls' theme resources), its control templates, the AnimatedIcon sources it draws its check mark and chevrons with (AnimatedAcceptVisualSource, AnimatedChevronDownSmallVisualSource, AnimatedChevronUpDownSmallVisualSource), the ProgressRing's Lottie (ProgressRingIndeterminate.cpp), ProgressBar.cpp and the focus visual's FocusRectManager.cpp. Their colors are translucent layers over the window's SolidBackgroundFillColorBase, resolved here to opaque ones.

// fluentAccent is Windows's default accent color, SystemAccentColor.
var fluentAccent = ui.Hex("#0078d4")

// fluentSystemAccent returns SystemAccentColor: the accent the user chose, opaque, else Windows's default.
func fluentSystemAccent(c *ui.Context) ui.Color {
	if a := c.Preferences().Accent; a.A > 0 {
		a.A = 255
		return a
	}
	return fluentAccent
}

// fluentDefaultShades are the shades Windows derives from its default accent, #0078d4: SystemAccentColorLight2 and Light3 (2, 3) and Dark1 to Dark3 (-1 to -3). tools/winui-reference sets them so to render WinUI's controls, and its accent button reads #005a9e in the light theme and #76b9ed in the dark one, its hyperlink #004275 and #a6d8ff (measured from WinUI 3 3.3 in Windows App SDK 2.6.5 on Windows 11 26H2 by tools/winui-reference).
var fluentDefaultShades = map[int]ui.Color{2: ui.Hex("#76b9ed"), 3: ui.Hex("#a6d8ff"), -1: ui.Hex("#005a9e"), -2: ui.Hex("#004275"), -3: ui.Hex("#002642")}

// fluentShadeMix is how far fluentShade mixes another accent with white (positive n) or black (negative n): mixes that give the default accent's shades within a few levels, as Windows's own derivation is not public.
var fluentShadeMix = map[int]float32{2: 0.48, 3: 0.68, -1: 0.25, -2: 0.45, -3: 0.68}

// fluentShade returns SystemAccentColorLight<n> for n of 2 or 3, or SystemAccentColorDark<-n> for n of -1 to -3: Windows's own shades for its default accent, mixes with white or black for another one.
func fluentShade(accent ui.Color, n int) ui.Color {
	accent.A = 255
	if accent == fluentAccent {
		return fluentDefaultShades[n]
	}
	if n > 0 {
		return accent.Mix(white, fluentShadeMix[n])
	}
	return accent.Mix(black, fluentShadeMix[n])
}

// fluentAccentFill returns AccentFillColorDefault for the system accent: SystemAccentColorDark1 in the light theme and SystemAccentColorLight2 in the dark one.
func fluentAccentFill(accent ui.Color, dark bool) ui.Color {
	if dark {
		return fluentShade(accent, 2)
	}
	return fluentShade(accent, -1)
}

// fluentSecondaryOver returns TextFillColorSecondary, #9E000000 in the light theme and #C5FFFFFF in the dark one, laid over under: the color of a ComboBox's or a DropDownButton's chevron on its face and of an off switch's knob on its track, which WinUI draws translucent over them rather than over the window (measured from WinUI 3 3.3 on Windows 11 26H2 by tools/winui-reference: the chevron reads 95 on the face of 251 where the window's TextFillColorSecondary is 92, 207 on 45 in the dark theme; the knob 90 on the track of 237, 87 under the pointer and 84 pressed).
func fluentSecondaryOver(t *ui.Theme, under ui.Color) ui.Color {
	return under.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0x9e)/255, 0xc5/255.0))
}

// fluentInnerStroke draws the inner half of the 1-DIP stroke of a WinUI Rectangle or Ellipse with a corner radius of radius along r, in c. Such a shape's fill reaches the middle of its stroke, so a translucent stroke reads as itself over the window on its outer half and over the fill on its inner half: an unchecked box's ControlStrongStrokeColorDefault reads 134 and 131 in the light theme, 191 and 173 pressed, and a checked box's AccentFillColorSecondary stroke under the pointer reads the fill, #1869a6, then #025b9e, the accent at 90% laid twice (measured from WinUI 3 3.3 on Windows 11 26H2 by tools/winui-reference). The element's border draws the outer half.
func fluentInnerStroke(p *ui.Painter, r ui.Rect, radius float32, c ui.Color) {
	p.Stroke(ui.Rect{X: r.X + 0.5, Y: r.Y + 0.5, W: r.W - 1, H: r.H - 1}, c, max(radius-0.5, 0), 0.5)
}

// fluentTwice returns a fill of the accent at alpha a laid over the window twice, as the inner half of the stroke of a checked box or a chosen radio button under the pointer (a 0.9) or pressed (a 0.8) reads, its stroke the same translucent brush as its fill.
func fluentTwice(t *ui.Theme, a float32) ui.Color {
	return t.Background.Mix(t.Accent, 1-(1-a)*(1-a))
}

func fluentTheme(dark bool, accent ui.Color) *ui.Theme {
	if accent.A == 0 {
		accent = fluentAccent
	}
	t := pick(dark, ui.LightTheme(), ui.DarkTheme())
	fill := fluentAccentFill(accent, dark)
	t.Accent = fill
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
	// AccentFillColorSecondaryBrush and AccentFillColorTertiaryBrush are the fill at Opacity 0.9 and 0.8 (Common_themeresources_any.xaml), laid over the window's SolidBackgroundFillColorBase set above.
	t.AccentHover = t.Background.Mix(fill, 0.9)
	t.AccentPressed = t.Background.Mix(fill, 0.8)
	// SystemFillColorCritical, SystemFillColorCaution and SystemFillColorSuccess (Common_themeresources_any.xaml).
	t.Danger = pick(dark, ui.Hex("#c42b1c"), ui.Hex("#ff99a4"))
	t.Warning = pick(dark, ui.Hex("#9d5d00"), ui.Hex("#fce100"))
	t.Success = pick(dark, ui.Hex("#0f7b0f"), ui.Hex("#6ccb5f"))
	// ScrollBarThumbFill is ControlStrongFillColorDefault, #72000000 and #8BFFFFFF (ScrollBar_themeresources.xaml). The thumb keeps MyGo's 6 DIPs, the fill of WinUI's expanded thumb (a 12-DIP Rectangle inside a 6-DIP transparent ScrollBarThumbStrokeThickness), rather than the 2-DIP indicator WinUI shows while the pointer is away from the bar, as MyGo shows one thumb that is dragged too; which width a Windows render shows over content is not settled.
	t.Scrollbar = pick(dark, ui.RGBA(0, 0, 0, 0x72/255.0), ui.RGBA(255, 255, 255, 0x8b/255.0))
	// TextControlSelectionHighlightColor is AccentFillColorSelectedTextBackgroundBrush, the opaque SystemAccentColor, under TextOnAccentFillColorSelectedText, white (TextBox_themeresources.xaml, Common_themeresources_any.xaml). MyGo keeps the text's own color over the highlight: in the dark theme that is white, so the opaque accent is WinUI's; in the light theme the text is near black, so the accent stays translucent there to keep it readable.
	t.Selection = pick(dark, accent.Alpha(0.4), accent)
	// ToolTip: the in-app acrylic's fallback color, not the theme turned over.
	t.Inverse, t.InverseText = pick(dark, ui.Hex("#f9f9f9"), ui.Hex("#2c2c2c")), t.Text
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

// fluentElevationOver is ControlElevationBorderBrush as fluentElevation lays it, left translucent, for a border drawn over something other than the window: a slider's thumb over its track, where WinUI's thumb edge reads #005595 over the light theme's accent and 126 over the gray track, 161 over the dark theme's (measured from WinUI 3 3.3 on Windows 11 26H2 by tools/winui-reference).
func fluentElevationOver(t *ui.Theme, h float32) ui.LinearGradient {
	if t.Dark {
		def, sec := white.Alpha(0x12/255.0), white.Alpha(0x18/255.0)
		return ui.LinearGradient{Angle: 180, Stops: []ui.GradientStop{{Color: sec, At: 0}, {Color: sec, At: 1 / h}, {Color: def, At: 3 / h}, {Color: def, At: 1}}}
	}
	def, sec := black.Alpha(0x0f/255.0), black.Alpha(0x29/255.0)
	return ui.LinearGradient{Angle: 180, Stops: []ui.GradientStop{{Color: def, At: 0}, {Color: def, At: (h - 3) / h}, {Color: sec, At: (h - 1) / h}, {Color: sec, At: 1}}}
}

// fluentFocusRing draws WinUI's system focus visual whose outer edge is bounds, the focus target grown by the negative of its FocusVisualMargin: the 2-DIP primary stroke in FocusStrokeColorOuter (the theme's Focus) and inside it the 1-DIP secondary stroke in FocusStrokeColorInner, #B3FFFFFF in the light theme and #B3000000 in the dark one (Common_themeresources_any.xaml). FocusRectManager.cpp's GetFocusOptionsForElement takes the thicknesses from the element's FocusVisualPrimaryThickness and FocusVisualSecondaryThickness, whose defaults, 2 and 1, are FrameworkElement's documented ones and not in the sources read (inferred; the same file's hyperlink focus options also use 2 and 1), and ConfigureFocusElement puts the inner stroke inside the outer one. radius is the target's corner radius; AdjustCornerRadius in the same file adds the mean of the two margins at each corner for the outer stroke and takes the primary's 2 off again for the inner one. MyGo's FocusRing draws its 2-DIP ring 1 to 3 DIPs outside the rect it is given, so it is given bounds less 3.
func fluentFocusRing(p *ui.Painter, t *ui.Theme, target ui.Rect, start, top, end, bottom, radius float32) {
	bounds := ui.Rect{X: target.X - start, Y: target.Y - top, W: target.W + start + end, H: target.H + top + bottom}
	outer := radius + (top+start)/2
	p.FocusRing(ui.Rect{X: bounds.X + 3, Y: bounds.Y + 3, W: bounds.W - 6, H: bounds.H - 6}, radii(max(outer-3, 0)))
	inner := pick(t.Dark, white, black).Alpha(0.7)
	p.Stroke(ui.Rect{X: bounds.X + 2, Y: bounds.Y + 2, W: bounds.W - 4, H: bounds.H - 4}, inner, max(outer-2, 0), 1)
}

// fluentButton draws DefaultButtonStyle, or AccentButtonStyle when primary: 32 DIPs high with ButtonPadding 11,5,11,6 inside a 1-DIP elevation border, corners of 4, and the system focus visual 3 DIPs around it (FocusVisualMargin -3, Button_themeresources.xaml and ToggleButton_themeresources.xaml).
func fluentButton(c *ui.Context, b ui.Element, primary bool) {
	t := c.Theme()
	b.MinHeight(32).Padding(5, 11, 6, 11).Gap(8).Radius(4).Justify(ui.Center).AlignItems(ui.Center).FocusRing(false)
	focused := b.FocusVisible()
	b.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if focused {
			fluentFocusRing(p, t, r, 3, 3, 3, 3, 4)
		}
	})
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
	face, fg := fluentFace(t, pressed, hovered), t.Text
	if pressed {
		fg = t.TextMuted
	}
	b.Background(face).TextColor(fg)
	if pressed {
		b.Border(1, t.Border)
	} else {
		b.Border(1, ui.Color{}).BorderLinearGradient(fluentElevation(t, 32))
	}
}

// fluentFace returns a button's ControlFillColorDefault, Secondary under the pointer and Tertiary pressed, laid over the window.
func fluentFace(t *ui.Theme, pressed, hovered bool) ui.Color {
	switch {
	case pressed:
		return t.SurfacePressed
	case hovered:
		return t.SurfaceHover
	}
	return t.Surface
}

// fluentUnchecked returns the fill and the stroke of an unchecked check box, or of a switch that is off: CheckBoxCheckBackgroundFillUnchecked and CheckBoxCheckBackgroundStrokeUnchecked (ToggleSwitchFillOff and ToggleSwitchStrokeOff are the same), the stroke, ControlStrongStrokeColorDefault, laid over the window for its outer half and over the fill for its inner half (fluentInnerStroke).
func fluentUnchecked(t *ui.Theme, hovered, pressed bool) (fill, stroke, inner ui.Color) {
	bg := t.Background
	k, a := black, float32(0.447)
	if t.Dark {
		k, a = white, 0.545
		fill = bg.Mix(black, 0.098)
		switch {
		case pressed:
			fill = bg.Mix(white, 0.071)
		case hovered:
			fill = bg.Mix(white, 0.043)
		}
	} else {
		fill = bg.Mix(black, 0.024)
		switch {
		case pressed:
			fill = bg.Mix(black, 0.094)
		case hovered:
			fill = bg.Mix(black, 0.059)
		}
	}
	return fill, bg.Mix(k, a), fill.Mix(k, a)
}

// fluentPressedStroke returns the stroke of an unchecked check box or radio button that is pressed, ControlStrongStrokeColorDisabled (#37000000, #28FFFFFF), laid over the window and over fill.
func fluentPressedStroke(t *ui.Theme, fill ui.Color) (stroke, inner ui.Color) {
	k, a := black, float32(0x37)/255
	if t.Dark {
		k, a = white, 0x28/255.0
	}
	return t.Background.Mix(k, a), fill.Mix(k, a)
}

// fluentCheckbox draws a CheckBox: a 20-DIP box with corners of 4 and a 1-DIP stroke, 8 DIPs from its label, in a row at least 32 high and 120 wide (CheckBoxHeight, CheckBoxMinWidth); checked, the accent fill with the check mark of AnimatedAcceptVisualSource in TextOnAccentFillColorPrimary, drawing itself in as WinUI's does, or with mixed the indeterminate dash (formsFluentDash), and gray when disabled. The focus visual goes around the whole row, grown by CheckBoxFocusVisualMargin -7,-3,-7,-3 (CheckBox_themeresources.xaml).
func fluentCheckbox(c *ui.Context, row, box ui.Element, on, mixed bool) {
	t := c.Theme()
	// WinUI does not center the box and the label in the row: the box is centered in a grid CheckBoxHeight (32) high at the top, so 6 DIPs down, and the label is at the top of the ContentPresenter inside the style's Padding 8,5,0,0, so 5 DIPs down (measured from WinUI 3 3.3 on Windows 11 26H2 by tools/winui-reference: the box's top at 6 and the label's text box at 5, which put its baseline 1.5 DIPs above where centering the label put Breeze UI's).
	row.Gap(8).MinHeight(32).MinWidth(120).AlignItems(ui.Start).Padding(5, 0, 0, 0)
	hovered, pressed := row.Hovered(), row.Pressed()
	fill, stroke, inner := fluentUnchecked(t, hovered, pressed)
	if !on && pressed {
		stroke, inner = fluentPressedStroke(t, fill)
	}
	mark := t.AccentText
	if on {
		fill, inner = t.Accent, t.Accent
		switch {
		case pressed:
			fill, inner, mark = t.AccentPressed, fluentTwice(t, 0.8), mark.Alpha(pick(t.Dark, float32(0.7), 0.5))
		case hovered:
			fill, inner = t.AccentHover, fluentTwice(t, 0.9)
		}
	}
	box.Size(20, 20).Radius(4).Margin(1, 0, 0, 0)
	if on {
		// Disabled is known only once the row is built, so a checked box is filled as it is painted, its stroke the fill.
		box.Border(1, ui.Color{})
	} else {
		box.Background(fill).Border(1, stroke)
	}
	disabled := disabledOnce(row)
	// The check mark draws in by its TrimEnd from 0 to 1 over frames 15 to 34 of the 60-frame-a-second Lottie (NormalOffToNormalOn, 317 ms) along cubic-bezier(0.55,0,0,1); unchecked, it is gone within 4 frames, drawn here at once.
	target, dur := float32(0), time.Duration(0)
	if on {
		target, dur = 1, 317*time.Millisecond
	}
	drawn := box.AnimateWith("breeze-check", target, dur, ui.CubicBezier(0.55, 0, 0, 1))
	box.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if !on {
			fluentInnerStroke(p, r, 4, inner)
			return
		}
		face, ring, mark := fill, inner, mark
		if disabled() {
			face, ring, mark = fluentDisabledOn(t)
		}
		p.Fill(r, face, 4)
		if ring != face {
			fluentInnerStroke(p, r, 4, ring)
		}
		if mixed {
			formsFluentDash(p, r, mark)
			return
		}
		// AnimatedAcceptVisualSource.cpp's check, (-15.172,0.016) to (-5,10.188) to (15.337,-10.337) stroked 4 with round caps and join, in a shape scaled 0.7 at 24,23 in a container scaled 1.05 about 24,24, its 48-unit visual fit to the 20-DIP box moved down 1 by the AnimatedIcon's Margin 0,1,0,-1: x = 10 + 0.30625·px, y = 10.5625 + 0.30625·py, and a stroke of 4·0.7·1.05·20/48 = 1.225 DIPs. No state scales the check, so the stroke is that at every size.
		var path ui.Path
		path.MoveTo(r.X+5.354, r.Y+10.567).LineTo(r.X+8.469, r.Y+13.683).LineTo(r.X+14.697, r.Y+7.397).Trim(0, drawn)
		p.StrokePath(&path, 1.225, mark)
	})
	row.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if row.FocusVisible() {
			fluentFocusRing(p, t, r, 7, 3, 7, 3, 4)
		}
	})
}

// fluentDisabledOn returns the box, the inner half of its stroke and the mark of a checked check box that is disabled: WinUI draws the box in AccentFillColorDisabled (#37000000, #28FFFFFF), its stroke the same brush laid twice on its inner half, and the mark in TextOnAccentFillColorDisabled (white, #87FFFFFF), reading 191, 150 and 255 in the light theme and 67, 96 and 167 in the dark one (measured from WinUI 3 3.3 on Windows 11 26H2 by tools/winui-reference). MyGo draws a disabled element at half opacity, each fill on its own, so each is the color that reads so at half opacity over what is under it (fluentUnderHalf); the mark can only come to 223 and 161 so, white being the lightest there is.
func fluentDisabledOn(t *ui.Theme) (box, ring, mark ui.Color) {
	k, a, onAccent := black, float32(0x37)/255, white
	if t.Dark {
		k, a = white, 0x28/255.0
	}
	face, twice := t.Background.Mix(k, a), t.Background.Mix(k, 1-(1-a)*(1-a))
	if t.Dark {
		onAccent = face.Mix(white, 0x87/255.0)
	}
	return fluentUnderHalf(face, t.Background), fluentUnderHalf(twice, face), fluentUnderHalf(onAccent, face)
}

// fluentUnderHalf returns the color that reads as want when drawn at half opacity over under, as far as one can.
func fluentUnderHalf(want, under ui.Color) ui.Color {
	ch := func(w, u uint8) uint8 { return uint8(min(max(2*int(w)-int(u), 0), 255)) }
	return ui.Color{R: ch(want.R, under.R), G: ch(want.G, under.G), B: ch(want.B, under.B), A: 255}
}

// fluentSwitch draws a ToggleSwitch: a 40 by 20 track, a 1-DIP stroke while off, and a round knob of 12 DIPs that grows to 14 under the pointer and stretches to 17 wide while pressed, moving in 83 milliseconds.
func fluentSwitch(c *ui.Context, sw ui.Element, on bool) {
	t := c.Theme()
	fluentSwitchTrack(c, t, sw, sw.Hovered(), sw.Pressed(), on)
	// The focus target is SwitchAreaGrid: every column of the switch's grid, the 40-DIP track, the 12-DIP gap and the On or Off content, which Breeze's switch does not show, so 52 DIPs, and its rows, 10 + 20 + 10 less its Margin 0,5, so 30 DIPs around the track; grown by FocusVisualMargin -7,-3,-7,-3 with corners of ControlCornerRadius (ToggleSwitch_themeresources.xaml).
	sw.FocusRing(false)
	rtl := c.Direction() == ui.RightToLeft
	sw.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if !sw.FocusVisible() {
			return
		}
		area := ui.Rect{X: r.X, Y: r.Y - 5, W: r.W + 12, H: r.H + 10}
		if rtl {
			area.X -= 12
		}
		fluentFocusRing(p, t, area, 7, 3, 7, 3, 4)
	})
}

// fluentSwitchTrack draws a ToggleSwitch's track and knob on sw, laid over t's Background, under the pointer if hovered and pressed if pressed: the look of fluentSwitch without its focus visual, for a switch whose state is another element's, as SwitchRow's.
func fluentSwitchTrack(c *ui.Context, t *ui.Theme, sw ui.Element, hovered, pressed, on bool) {
	pos := knob(sw, on, 83*time.Millisecond)
	at := along(c, pos)
	offFill, offStroke, offInner := fluentUnchecked(t, hovered, pressed)
	onFill := t.Accent
	switch {
	case pressed:
		onFill = t.AccentPressed
	case hovered:
		onFill = t.AccentHover
	}
	knobFill := fluentOffKnob(t, offFill).Mix(t.AccentText, pos)
	inner := offInner.Mix(onFill, pos)
	sw.Size(40, 20).Radius(10).Background(offFill.Mix(onFill, pos)).Border(1, offStroke.Mix(onFill, pos))
	w, h := float32(12), float32(12)
	switch {
	case pressed:
		w, h = 17, 14
	case hovered:
		w, h = 14, 14
	}
	sw.Draw(func(p *ui.Painter, r ui.Rect) {
		fluentInnerStroke(p, r, 10, inner)
		// The knob's center goes from 10 to 30 DIPs; pressed, it keeps its outer edge and stretches toward the other end.
		cx := r.X + r.H/2 + at*(r.W-r.H)
		x := cx - h/2 - at*(w-h)
		p.Fill(ui.Rect{X: x, Y: r.Y + (r.H-h)/2, W: w, H: h}, knobFill, h/2)
	})
}

// fluentOffKnob is ToggleSwitchKnobFillOff, TextFillColorSecondary, a gray knob laid over the track's fill, track.
func fluentOffKnob(t *ui.Theme, track ui.Color) ui.Color { return fluentSecondaryOver(t, track) }

// ---- more of WinUI 3's controls, from the same theme resources ----

// fluentSubtle returns SubtleFillColorSecondary and SubtleFillColorTertiary, the backgrounds of a list item under the pointer and pressed.
func fluentSubtle(t *ui.Theme) (hover, pressed ui.Color) {
	if t.Dark {
		return t.Background.Mix(white, 0.059), t.Background.Mix(white, 0.039)
	}
	return t.Background.Mix(black, 0.035), t.Background.Mix(black, 0.024)
}

// fluentRadio draws a RadioButton: a 20-DIP circle with a 1-DIP stroke, 8 DIPs from its label in a row at least 32 high and 120 wide (the style's MinWidth); chosen, the accent with a dot of TextOnAccentFillColorPrimary, 12 DIPs across, growing to 14 under the pointer and shrinking to 10 pressed over ControlNormalAnimationDuration (250 ms) along ControlFastOutSlowInKeySpline 0,0,0,1, back to 12 at once, inside a 1-DIP AccentControlElevationBorderBrush stroke (RadioButtonCheckGlyphStrokeChecked); pressed, PressedCheckGlyph, a dot of the same fill growing from 4 to 10 DIPs over ControlFastAnimationDuration (167 ms), over the chosen dot's ring too. The focus visual goes around the whole row, grown by FocusVisualMargin -7,-3,-7,-3 (RadioButton_themeresources.xaml).
func fluentRadio(c *ui.Context, row, dot ui.Element, on bool) {
	t := c.Theme()
	// The circle is centered in a grid 32 high at the top and the label at the top inside the style's Padding 8,6,0,0, so both 6 DIPs down (measured from WinUI 3 3.3 on Windows 11 26H2 by tools/winui-reference).
	row.Gap(8).MinHeight(32).MinWidth(120).AlignItems(ui.Start).Padding(6, 0, 0, 0)
	hovered, pressed := row.Hovered(), row.Pressed()
	fill, stroke, inner := fluentUnchecked(t, hovered, pressed)
	if pressed && !on {
		stroke, inner = fluentPressedStroke(t, fill)
	}
	ease := ui.CubicBezier(0, 0, 0, 1)
	size, dur := float32(12), time.Duration(0)
	if on {
		fill, inner = t.Accent, t.Accent
		switch {
		case pressed:
			fill, inner, size, dur = t.AccentPressed, fluentTwice(t, 0.8), 10, 250*time.Millisecond
		case hovered:
			fill, inner, size, dur = t.AccentHover, fluentTwice(t, 0.9), 14, 250*time.Millisecond
		}
		stroke = fill
	}
	d := dot.AnimateWith("breeze-dot", size, dur, ease)
	down := pressed && !on
	pressSize, pressDur := float32(4), time.Duration(0)
	if pressed {
		pressSize, pressDur = 10, 167*time.Millisecond
	}
	press := dot.AnimateWith("breeze-press-dot", pressSize, pressDur, ease)
	// AccentControlElevationBorderBrush over the dot: ControlStrokeColorOnAccentDefault (#14FFFFFF) and, in the bottom DIP blending into it by 3 DIPs from the bottom, ControlStrokeColorOnAccentSecondary (#66000000 in the light theme, #23000000 in the dark one), the brush's ScaleY -1 putting the secondary at the bottom in both (Common_themeresources_any.xaml). The dot's fill reaches the middle of that stroke, so its outer half lies over the accent around the dot, translucent: at the dot's bottom WinUI's reads #00365f over the light theme's #005a9e, and a 12-DIP dot's fill is 11 across (measured from WinUI 3 3.3 on Windows 11 26H2 by tools/winui-reference).
	on1 := white.Alpha(0x14 / 255.0)
	sec := black.Alpha(pick(t.Dark, float32(0x66)/255, 0x23/255.0))
	// Under the pointer and pressed, the CommonStates set the dot's Stroke to RadioButtonCheckGlyphStrokePointerOver and Pressed, CircleElevationBorderBrush: ControlStrokeColorDefault down to half the dot's height, blending into ControlStrokeColorSecondary by 70% (#0F000000 and #29000000, #12FFFFFF and #18FFFFFF in the dark theme; RadioButton_themeresources.xaml, Common_themeresources_any.xaml). Under the pointer WinUI's reads so, #17639d above the dot and #14598d below it over the light theme's accent (measured from WinUI 3 3.3 on Windows 11 26H2 by tools/winui-reference). Pressed, the Pressed state also shows PressedCheckGlyph, a Border of RadioButtonCheckGlyphFillPressed (TextOnAccentFillColorPrimary) without a border thickness, growing from 4 to 10 DIPs over ControlFastAnimationDuration above the dot, which shrinks to 10 over ControlNormalAnimationDuration; the Checked state hides it with a To animation that the Pressed state's, started later, overrides, so it covers the dot's ring: WinUI's pressed render reads a 10-DIP dot of TextOnAccentFillColorPrimary without a ring (RadioButton_themeresources.xaml; measured from WinUI 3 3.3 on Windows 11 26H2 by tools/winui-reference).
	cdef, csec := black.Alpha(0x0f/255.0), black.Alpha(0x29/255.0)
	if t.Dark {
		cdef, csec = white.Alpha(0x12/255.0), white.Alpha(0x18/255.0)
	}
	circle := ui.LinearGradient{Angle: 180, Stops: []ui.GradientStop{{Color: cdef, At: 0}, {Color: cdef, At: 0.5}, {Color: csec, At: 0.7}, {Color: csec, At: 1}}}
	dot.Size(20, 20).Radius(10).Background(fill).Border(1, stroke)
	dot.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if inner != fill {
			fluentInnerStroke(p, r, 10, inner)
		}
		cx, cy := r.X+r.W/2, r.Y+r.H/2
		switch {
		case on:
			p.Fill(ui.Rect{X: cx - d/2 + 0.5, Y: cy - d/2 + 0.5, W: d - 1, H: d - 1}, t.AccentText, d/2-0.5)
			var ring ui.Path
			ring.Circle(cx, cy, d/2-0.5)
			g := ui.LinearGradient{Angle: 180, Stops: []ui.GradientStop{{Color: on1, At: 0}, {Color: on1, At: (d - 3) / d}, {Color: sec, At: (d - 1) / d}, {Color: sec, At: 1}}}
			if hovered || pressed {
				g = circle
			}
			p.StrokePathGradient(&ring, 1, g)
			if pressed {
				p.Fill(ui.Rect{X: cx - press/2, Y: cy - press/2, W: press, H: press}, t.AccentText, press/2)
			}
		case down:
			p.Fill(ui.Rect{X: cx - press/2, Y: cy - press/2, W: press, H: press}, t.AccentText, press/2)
		}
	})
	row.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if row.FocusVisible() {
			fluentFocusRing(p, t, r, 7, 3, 7, 3, 4)
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
	// The 1-DIP border, drawn over the box here, is inside the template's BorderElement, and TextControlThemePadding 10,5,6,6 inside it: the text starts 11 DIPs in and 6 down (measured from WinUI 3 3.3 on Windows 11 26H2 by tools/winui-reference: the text box at 27,56 in a TextBox at 16,50). The template mirrors right to left, so the sides are the start and the end.
	f.Padding(6, 0, 7, 0).PaddingStart(11).PaddingEnd(7).Radius(4).Background(fill).Border(0, ui.Color{})
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

// fluentSelect draws a closed ComboBox: a button's face and elevation border, 32 DIPs high and at least ComboBoxThemeMinWidth (64) wide with padding 12,5,0,7, a 12-DIP chevron 14 DIPs from its end in TextFillColorSecondary. Its keyboard focus is drawn by its template, not by the system focus visual, which DefaultComboBoxStyle turns off (UseSystemFocusVisuals is IsApplicationFocusVisualKindReveal, False): in the Focused and FocusedPressed states, which ComboBox_Partial.cpp picks for focus other than the pointer's, HighlightBackground shows 4 DIPs around the box (Margin -4) with a 2-DIP ComboBoxBackgroundBorderBrushFocused edge (FocusStrokeColorOuter, the theme's Focus), ComboBoxHiglightBorderCornerRadius 7 and ComboBoxBackgroundFocused (ControlFillColorDefault) behind the box, and the 3 by 16 accent Pill (ComboBoxItemPill, AccentFillColorDefault, corners of 1.5) shows 1 DIP in from the box's start (ComboBox_themeresources.xaml).
func fluentSelect(c *ui.Context, trigger ui.Element, selected string) {
	t := c.Theme()
	pressed, hovered := trigger.Pressed(), trigger.Hovered()
	face, fg := fluentFace(t, pressed, hovered), t.Text
	if pressed {
		fg = t.TextMuted
	}
	// ComboBoxPadding 12,5,0,7 and the DropDownGlyph's Margin 0,0,14,0 are from the box's outer edge, the 1-DIP border drawn inside them, where MyGo's padding starts inside the border: the text starts 12 DIPs in and the chevron ends 14 from the end (measured from WinUI 3 3.3 on Windows 11 26H2 by tools/winui-reference: the text at 28 and the 12-DIP glyph at 92 in a ComboBox from 16 to 118), at the start and the end as the template mirrors right to left.
	trigger.MinHeight(32).MinWidth(64).Padding(5, 0, 7, 0).PaddingStart(11).PaddingEnd(13).Gap(12).Radius(4).Justify(ui.SpaceBetween).AlignItems(ui.Center).
		Background(face).TextColor(fg).FocusRing(false)
	if pressed {
		trigger.Border(1, t.Border)
	} else {
		trigger.Border(1, ui.Color{}).BorderLinearGradient(fluentElevation(t, 32))
	}
	focused := trigger.FocusVisible()
	rtl := c.Direction() == ui.RightToLeft
	trigger.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if !focused {
			return
		}
		// HighlightBackground's fill shows only in the 2 DIPs between its edge and the box; its edge is MyGo's ring, which lies 1 to 3 DIPs outside the rect it is given.
		p.Stroke(ui.Rect{X: r.X - 2, Y: r.Y - 2, W: r.W + 4, H: r.H + 4}, t.Surface, 5, 2)
		p.FocusRing(ui.Rect{X: r.X - 1, Y: r.Y - 1, W: r.W + 2, H: r.H + 2}, radii(4))
		x := r.X + 1
		if rtl {
			x = r.X + r.W - 1 - 3
		}
		p.Fill(ui.Rect{X: x, Y: r.Y + (r.H-16)/2, W: 3, H: 16}, t.Accent, 1.5)
	})
	nudge := fluentChevronNudge(trigger, pressed)
	trigger.Children(func() {
		ui.Text(c, selected).SingleLine()
		// The DropDownGlyph is as high as the box, so the chevron is centered on the box, 1 DIP below the middle of the padded content (5 above, 7 below), in TextFillColorSecondary over the face (fluentSecondaryOver).
		ui.Box(c).Size(12, 12).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
			fluentChevronDown(p, r.X+r.W/2, r.Y+r.H/2+1+nudge, fluentSecondaryOver(t, face))
		})
	})
}

// fluentChevronDown strokes AnimatedChevronDownSmallVisualSource's chevron centered on cx, cy: its outline is a 1-unit stroke with round ends and corner along (-3.5,-1.75), (0,1.75), (3.5,-1.75), scaled 4 in a 48-unit visual that a 12-DIP AnimatedIcon shows at a quarter, so 1 unit is 1 DIP (AnimatedChevronDownSmallVisualSource.cpp). ComboBox and DropDownButton show it in 12 DIPs.
func fluentChevronDown(p *ui.Painter, cx, cy float32, col ui.Color) {
	var path ui.Path
	path.MoveTo(cx-3.5, cy-1.75).LineTo(cx, cy+1.75).LineTo(cx+3.5, cy-1.75)
	p.StrokePath(&path, 1, col)
}

// fluentChevronNudge returns how far the ChevronDownSmall glyph of e has moved down: pressed, its offset goes from 24 to 31.5 of the 48 units, 1.875 DIPs, over frames 0 to 9 (151 ms) along cubic-bezier(0.167,0.167,0.65,1); released it comes back over frames 10 to 29 (about 300 ms), overshooting 0.75 DIP above rest on the way, which is left out here: it eases back along the last segment's cubic-bezier(0.35,0,0,1) (AnimatedChevronDownSmallVisualSource.cpp, OffsetVector2Animation_0 and _1, c_durationTicks 4833333).
func fluentChevronNudge(e ui.Element, pressed bool) float32 {
	if pressed {
		return e.AnimateWith("breeze-chevron", 1.875, 151*time.Millisecond, ui.CubicBezier(0.167, 0.167, 0.65, 1))
	}
	return e.AnimateWith("breeze-chevron", 0, 300*time.Millisecond, ui.CubicBezier(0.35, 0, 0, 1))
}

// fluentMenu draws a ComboBox's drop-down: the in-app acrylic's fallback color, a 1-DIP SurfaceStrokeColorFlyout edge, OverlayCornerRadius and a flyout's shadow.
func fluentMenu(c *ui.Context, panel ui.Element) {
	t := c.Theme()
	panel.Margin(4, 0, 0, 0).Padding(4, 0).Radius(8).Background(pick(t.Dark, ui.Hex("#f9f9f9"), ui.Hex("#2c2c2c"))).
		Border(1, fluentFlyoutStroke(t)).Shadow(0, 8, 16, 0, ui.RGBA(0, 0, 0, pick(t.Dark, float32(0.14), 0.26)))
}

// fluentFlyoutStroke is SurfaceStrokeColorFlyout, #0F000000 in the light theme and #33000000 in the dark one, laid over the window: a flyout's background stops inside its border, so the translucent edge reads over what is behind the flyout, not over the flyout's own fill, which MyGo's border is drawn over (measured from WinUI 3 3.3 on Windows 11 26H2 by tools/winui-reference: the edge reads 229 and 26 above a flyout over the window, where over its fill it would read 234 and 35).
func fluentFlyoutStroke(t *ui.Theme) ui.Color {
	return t.Background.Mix(black, pick(t.Dark, float32(0x0f)/255, 0x33/255.0))
}

// fluentMenuItem draws a ComboBoxItem: padding 11,5,11,7 inside a margin of 5,2, corners of 3, and the chosen one marked by a 3 by 16 accent pill at its start. Its background is SubtleFillColorSecondary under the pointer and while chosen, Tertiary pressed and while chosen under the pointer, and Secondary chosen and pressed (ComboBoxItemBackground*); pressed, its text is TextFillColorSecondary, and the chosen one's pill shrinks to ComboBoxItemPillMinScale, 0.625, over ComboBoxItemScaleAnimationDuration, 167 ms, along KeySpline 0,0,0,1 (ComboBox_themeresources.xaml).
func fluentMenuItem(c *ui.Context, item ui.Element, label string, chosen bool) {
	t := c.Theme()
	secondary, tertiary := fluentSubtle(t)
	item.Margin(2, 5).Padding(5, 11, 7, 11).Radius(3).MinWidth(100)
	pressed, over := item.Pressed(), item.Highlighted()
	switch {
	case pressed && chosen:
		item.Background(secondary).TextColor(t.TextMuted)
	case pressed:
		item.Background(tertiary).TextColor(t.TextMuted)
	case chosen && item.Hovered():
		// SelectedPointerOver takes the pointer; MyGo's highlight by the arrows leaves the chosen item Selected.
		item.Background(tertiary)
	case chosen || over:
		item.Background(secondary)
	}
	if chosen {
		// The Selected state animates no ScaleY, so the pill is whole again at once on release.
		scale, dur := float32(1), time.Duration(0)
		if pressed {
			scale, dur = 0.625, 167*time.Millisecond
		}
		scale = item.AnimateWith("breeze-pill", scale, dur, ui.CubicBezier(0, 0, 0, 1))
		item.Draw(func(p *ui.Painter, r ui.Rect) {
			x := r.X
			if c.Direction() == ui.RightToLeft {
				x = r.X + r.W - 3
			}
			h := 16 * scale
			p.Fill(ui.Rect{X: x, Y: r.Y + (r.H-h)/2, W: 3, H: h}, t.Accent, 1.5)
		})
	}
	item.Children(func() { ui.Text(c, label).SingleLine() })
}

// fluentSlider draws a Slider: 32 DIPs high, a 4-DIP track of ControlStrongFillColorDefault filled with the accent (AccentFillColorSecondary under the pointer and Tertiary pressed, SliderTrackValueFillPointerOver and Pressed), and a 22-DIP thumb of ControlSolidFillColorDefault with an elevation border around an accent dot that grows under the pointer and shrinks pressed.
func fluentSlider(c *ui.Context, s ui.Element, frac float32) {
	t := c.Theme()
	// The thumb is SliderHorizontalThumbWidth, 18, whose border reaches 2 beyond it on each side; its center travels the track less 18.
	const travel = 18
	s.Height(32).PaddingX(travel / 2).MinWidth(60).FocusRing(false)
	inner, accent := fluentSliderState(t, s)
	track := t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.447), 0.545))
	at := along(c, frac)
	s.Draw(func(p *ui.Painter, r ui.Rect) {
		ty := r.Y + (r.H-4)/2
		p.Fill(ui.Rect{X: r.X, Y: ty, W: r.W, H: 4}, track, 2)
		x := r.X + travel/2 + at*(r.W-travel)
		fill := ui.Rect{X: r.X, Y: ty, W: x - r.X, H: 4}
		if c.Direction() == ui.RightToLeft {
			fill = ui.Rect{X: x, Y: ty, W: r.X + r.W - x, H: 4}
		}
		p.Fill(fill, accent, 2)
		fluentSliderThumb(p, t, x, r.Y+r.H/2, inner, accent)
		// The keyboard's focus goes to Slider's FocusBorder, the whole 32-DIP slider with corners of ControlCornerRadius, widened 7 at each end by its FocusVisualMargin of -7,0,-7,0.
		if s.FocusVisible() {
			fluentFocusRing(p, t, r, 7, 0, 7, 0, 4)
		}
	})
}

// fluentSliderState returns the size of a slider's inner thumb, SliderInnerThumbWidth 12 scaled 0.86, 1.167 under the pointer and 0.71 pressed, and the accent of its value fill and inner thumb: AccentFillColorDefault, Secondary under the pointer and Tertiary pressed (Slider_themeresources.xaml).
func fluentSliderState(t *ui.Theme, s ui.Element) (inner float32, accent ui.Color) {
	switch {
	case s.Pressed():
		return 8.5, t.AccentPressed
	case s.Hovered():
		return 14, t.AccentHover
	}
	return 10.3, t.Accent
}

// fluentSliderThumb draws a slider's thumb centered on cx, cy: the 18-DIP Thumb's Border with Margin -2, so 22 DIPs, with SliderThumbCornerRadius 10, a 1-DIP SliderThumbBorderBrush (ControlElevationBorderBrush, translucent over the track as fluentElevationOver lays it) around SliderOuterThumbBackground (ControlSolidFillColorDefault, #FFFFFF and #454545), and the inner dot of the given size in accent (Slider_themeresources.xaml).
func fluentSliderThumb(p *ui.Painter, t *ui.Theme, cx, cy, inner float32, accent ui.Color) {
	const thumb = 22
	k := ui.Rect{X: cx - thumb/2, Y: cy - thumb/2, W: thumb, H: thumb}
	p.FillGradient(k, fluentElevationOver(t, thumb), 10)
	p.Fill(ui.Rect{X: k.X + 1, Y: k.Y + 1, W: thumb - 2, H: thumb - 2}, pick(t.Dark, white, ui.Hex("#454545")), 9)
	p.Fill(ui.Rect{X: cx - inner/2, Y: cy - inner/2, W: inner, H: inner}, accent, inner/2)
}

// fluentProgress draws a ProgressBar: a 3-DIP accent bar with rounded ends over a 1-DIP line of ControlStrongStrokeColorDefault, which hides while two bars move for progress of unknown length (fluentIndeterminateBars).
func fluentProgress(c *ui.Context, e ui.Element, value float64) {
	t := c.Theme()
	line := t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.447), 0.545))
	e.Height(3).MinWidth(40)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		if value < 0 {
			fluentIndeterminateBars(c, p, r, t.Accent)
			return
		}
		p.Fill(ui.Rect{X: r.X, Y: r.Y + 1, W: r.W, H: 1}, line, 0.5)
		drawBar(c, p, r, value, t.Accent, 1.5)
	})
}

// fluentIndeterminateBars draws the Indeterminate state of a ProgressBar in r, clipped to it, repeating every 2 seconds (ProgressBar.xaml's storyboard, ProgressBar.cpp's UpdateWidthBasedTemplateSettings): a bar 40% of the width moving from -100% to 300% of its width over 1.5 s along KeySpline 0.4,0,0.6,1 and held there, and a bar 60% of the width held at -150% of its width until 0.75 s, then moving to 166% of it by 2 s along the same spline. Right to left they move the other way.
func fluentIndeterminateBars(c *ui.Context, p *ui.Painter, r ui.Rect, fill ui.Color) {
	p.AnimationFrame()
	bars := fluentIndeterminateAt(float32(p.Now().UnixMilli()%2000)/1000, r.W)
	rtl := c.Direction() == ui.RightToLeft
	p.Clip(r, 0, func() {
		for _, b := range bars {
			x := r.X + b[0]
			if rtl {
				x = r.X + r.W - b[0] - b[1]
			}
			p.Fill(ui.Rect{X: x, Y: r.Y, W: b[1], H: r.H}, fill, 1.5)
		}
	})
}

// fluentIndeterminateAt returns where the two bars of an indeterminate ProgressBar w DIPs wide are s seconds into the 2-second loop: each one's start from the bar's start and its width.
func fluentIndeterminateAt(s, w float32) [2][2]float32 {
	ease := ui.CubicBezier(0.4, 0, 0.6, 1)
	w1, w2 := w*0.4, w*0.6
	x1 := -w1 + 4*w1*ease(min(s/1.5, 1))
	x2 := -1.5 * w2
	if s > 0.75 {
		x2 += (1.5 + 1.66) * w2 * ease(min((s-0.75)/1.25, 1))
	}
	return [2][2]float32{{x1, w1}, {x2, w2}}
}

// fluentSpinner draws an indeterminate ProgressRing as its Lottie, ProgressRingIndeterminate.cpp, does: an ellipse of radius 7 scaled 5 on an 80-unit canvas, stretched to the ring's 32 DIPs, so radius 14 with a 1.5·5·0.4 = 3-DIP accent stroke and round ends. Over each 2 s it turns from 0 to 450 degrees in the first second and on to 900 in the next, along cubic-bezier(0.167,0.167,0.833,0.833); in the first second its end grows from its start to half the circle (TrimEnd 0 to 0.5), in the next its start catches up (TrimStart 0 to 0.5), on the same curve. Where on the ellipse the trim starts is not in the source: 12 o'clock here (inferred).
func fluentSpinner(c *ui.Context, e ui.Element) {
	t := c.Theme()
	e.Size(32, 32)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		p.AnimationFrame()
		start, sweep := fluentRingArcAt(float32(p.Now().UnixMilli()%2000) / 1000)
		strokeArc(p, r.X+r.W/2, r.Y+r.H/2, 14, start, sweep, 3, t.Accent)
	})
}

// fluentRingArcAt returns where the arc of an indeterminate ProgressRing starts and how far it sweeps, in radians clockwise, s seconds into its 2-second loop.
func fluentRingArcAt(s float32) (start, sweep float64) {
	ease := ui.CubicBezier(0.167, 0.167, 0.833, 0.833)
	var turn, from, to float32
	if s < 1 {
		u := ease(s)
		turn, from, to = 450*u, 0, 0.5*u
	} else {
		u := ease(min(s-1, 1))
		turn, from, to = 450+450*u, 0.5*u, 0.5
	}
	return float64(turn)*math.Pi/180 + float64(from)*2*math.Pi, float64(to-from) * 2 * math.Pi
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

// fluentSegment draws a SelectorBarItem (SelectorBar.xaml, SelectorBar_themeresources.xaml): padding 12,10,12,7 above its own 3-DIP row for the pill, text in TextFillColorPrimary, Secondary under the pointer and Tertiary pressed, while the chosen one is Secondary both under the pointer and pressed (SelectedPointerOver and SelectedPressed take SelectorBarItemForegroundPointerOver); chosen, the 4 by 3 accent pill at the bottom of the item, centered, scaled 4 wide to 16 and fading in over 167 ms along KeySpline 0,0,0,1. Its RadiusX 0.5 and RadiusY 1, scaled with it, round its corners 2 by 1, drawn here with corners of 1 as MyGo's rectangles have circular corners. The focus visual is 2 DIPs around it (SelectorBarItemFocusVisualMargin -2).
func fluentSegment(c *ui.Context, seg ui.Element, chosen bool) {
	t := c.Theme()
	fg := t.Text
	switch {
	case chosen && (seg.Pressed() || seg.Hovered()):
		fg = t.TextMuted
	case seg.Pressed():
		fg = t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.447), 0.529))
	case seg.Hovered():
		fg = t.TextMuted
	}
	seg.Padding(10, 12, 10, 12).Radius(4).Justify(ui.Center).AlignItems(ui.Center).TextColor(fg).FocusRing(false)
	// The Selected states take the pill's ScaleX to 4 and its Opacity to 1 over ComboBoxItemScaleAnimationDuration along the same spline; the Unselected ones animate neither, so it goes at once.
	grown, dur := float32(0), time.Duration(0)
	if chosen {
		grown, dur = 1, 167*time.Millisecond
	}
	grown = seg.AnimateWith("breeze-pill", grown, dur, ui.CubicBezier(0, 0, 0, 1))
	if chosen {
		seg.Draw(func(p *ui.Painter, r ui.Rect) {
			w := 4 + 12*grown
			p.Fill(ui.Rect{X: r.X + (r.W-w)/2, Y: r.Y + r.H - 3, W: w, H: 3}, t.Accent.Alpha(grown), 1)
		})
	}
	focused := seg.FocusVisible()
	seg.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if focused {
			fluentFocusRing(p, t, r, 2, 2, 2, 2, 4)
		}
	})
}

// fluentTabList draws tabs as a SelectorBar, which WinUI offers for switching between pages of a view.
func fluentTabList(c *ui.Context, list ui.Element) { fluentSegmentedTrack(c, list) }

// fluentTab draws a tab as a SelectorBarItem.
func fluentTab(c *ui.Context, tab ui.Element, chosen bool) { fluentSegment(c, tab, chosen) }

// ---- dialogs: WinUI 3's ContentDialog ----

// fluentDialogShadow casts a ContentDialog's elevation of 128: an ambient and a directional shadow, as DropShadowRecipe.h computes them.
func fluentDialogShadow(t *ui.Theme, panel ui.Element) {
	panel.Shadow(0, 2, 21, 0, ui.RGBA(0, 0, 0, pick(t.Dark, float32(0.15), 0.37))).
		Shadow(0, 32, 64, 0, ui.RGBA(0, 0, 0, pick(t.Dark, float32(0.19), 0.37)))
}

// fluentDialogStroke is a ContentDialog's border, SurfaceStrokeColorDefault (#66757575), laid over SmokeFillColorDefault over the window: the dialog's background stops inside its border, so the edge reads over what is behind the dialog, not over its own fill (measured from WinUI 3 3.3 on Windows 11 26H2 by tools/winui-reference: shown in place, without the smoke, its sides read 177 and 178 and its top 184 where over the dialog's fill they would read 193, the shadow behind darkening them).
func fluentDialogStroke(t *ui.Theme) ui.Color {
	return t.Background.Mix(black, 0.3).Mix(ui.Hex("#757575"), 0x66/255.0)
}

// fluentDialog draws a ContentDialog: 320 to 548 DIPs wide, OverlayCornerRadius, a 1-DIP SurfaceStrokeColorDefault border, LayerFillColorAlt inside its 24-DIP padding, over SmokeFillColorDefault.
func fluentDialog(c *ui.Context, back, panel ui.Element) {
	t := c.Theme()
	back.Background(ui.RGBA(0, 0, 0, 0.3))
	panel.MinWidth(320).MaxWidth(548).Padding(24).Gap(12).Radius(8).Border(1, fluentDialogStroke(t)).
		Background(pick(t.Dark, white, t.Background.Mix(white, 0.051)))
	fluentDialogShadow(t, panel)
}

// fluentAlert draws a ContentDialog with a title, a message and buttons: the title 20 DIPs and semibold, 12 above the message, in the padded top area with a line under it, and the buttons in the command area below, stretched to equal widths 8 apart. The default button comes first, as the primary button does, then the others in order; a single button takes the right half, as ContentDialog places its close button.
func fluentAlert(c *ui.Context, back, panel ui.Element, title, message string, n int, button func(int) ui.Element) {
	t := c.Theme()
	back.Background(ui.RGBA(0, 0, 0, 0.3))
	panel.MinWidth(320).MaxWidth(548).Radius(8).Border(1, fluentDialogStroke(t)).Background(t.Background).Clip().AlignItems(ui.Stretch)
	fluentDialogShadow(t, panel)
	panel.Children(func() {
		ui.Column(c).Padding(24).Background(pick(t.Dark, white, t.Background.Mix(white, 0.051))).
			BorderWidth(0, 0, 1, 0).BorderColor(t.Background.Mix(black, pick(t.Dark, float32(0.059), 0.098))).Children(func() {
			ui.Text(c, title).FontSize(20).FontWeight(600).Margin(0, 0, 12, 0)
			if message != "" {
				ui.Text(c, message)
			}
		})
		ui.Row(c).Padding(24).Gap(8).Children(func() {
			wide := func(b ui.Element) { b.Grow(1).Basis(0) }
			if n == 1 {
				ui.Box(c).Grow(1).Basis(0)
			}
			if n > 0 {
				wide(button(n - 1))
			}
			for i := range n - 1 {
				wide(button(i))
			}
		})
	})
}

// ---- links, dividers, badges, flyouts and grouped toggles ----

// fluentLinkColor returns a hyperlink's foreground: AccentTextFillColorPrimary (SystemAccentColorDark2 in the light theme, Light3 in the dark one), Secondary under the pointer (Dark3, Light3) and Tertiary pressed (Dark1, Light2), the shades of fluentShade.
func fluentLinkColor(c *ui.Context, hovered, pressed bool) ui.Color {
	base := fluentSystemAccent(c)
	if c.Theme().Dark {
		if pressed {
			return fluentShade(base, 2)
		}
		return fluentShade(base, 3)
	}
	switch {
	case pressed:
		return fluentShade(base, -1)
	case hovered:
		return fluentShade(base, -3)
	}
	return fluentShade(base, -2)
}

// fluentDivider is DividerStrokeColorDefault.
func fluentDivider(t *ui.Theme) ui.Color {
	return t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.059), 0.082))
}

// fluentBadge draws an InfoBadge with a value: AccentFillColorDefault with TextOnAccentFillColorPrimary, 16 DIPs high and at least as wide, round, the number 11 DIPs with a margin of 4,0,4,2.
func fluentBadge(c *ui.Context, b ui.Element) {
	t := c.Theme()
	b.Height(16).MinWidth(16).Padding(0, 4, 2, 4).Radius(8).Background(t.Accent).TextColor(t.AccentText).FontSize(11).FontWeight(400).Border(0, ui.Color{}).Justify(ui.Center).AlignItems(ui.Center)
}

// fluentFlyout draws a FlyoutPresenter: the in-app acrylic's fallback color, a 1-DIP SurfaceStrokeColorFlyout edge, OverlayCornerRadius, padding 16,15,16,17, at least FlyoutThemeMinWidth (96, generic.xaml) wide or as wide as its anchor, anchorW, as Breeze's popovers are on every platform, at most 456 DIPs wide, and the directional shadow of an elevation of 32. FlyoutThemeMaxHeight, 758, is left out: WinUI scrolls what is taller inside the presenter, which a style on the panel cannot do.
func fluentFlyout(c *ui.Context, panel ui.Element, anchorW float32) {
	t := c.Theme()
	panel.Margin(4, 0, 0, 0).Padding(15, 16, 17, 16).Gap(8).MinWidth(max(96, anchorW)).MinHeight(40).MaxWidth(456).Radius(8).
		Background(pick(t.Dark, ui.Hex("#f9f9f9"), ui.Hex("#2c2c2c"))).
		Border(1, fluentFlyoutStroke(t)).
		Shadow(0, 8, 17, 0, ui.RGBA(0, 0, 0, pick(t.Dark, float32(0.14), 0.26)))
}

// fluentToggleGroup lays out toggles on Windows, which has no control joining them: ToggleButtons in a row, 8 DIPs apart.
func fluentToggleGroup(c *ui.Context, g ui.Element) {
	g.Gap(8).Padding(0).Radius(0).Background(ui.Transparent).Border(0, ui.Color{})
}
