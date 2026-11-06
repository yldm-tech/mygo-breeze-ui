package breeze

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// GNOME: libadwaita 1.10.0 (GNOME 51, on GTK 4.24.0, with adwaita-icon-theme 51.0's symbolic icons), from its stylesheet (src/stylesheet/_colors.scss, which holds what _defaults.scss did before 1.10, _common.scss, _drawing.scss, widgets/_buttons.scss, _checks.scss and _switch.scss) and, for the accent and the switch's size, its C code and GTK's. Its colors are the same as libadwaita 1.7.6's. libadwaita colors controls with currentColor at an opacity; those are resolved here over window_bg_color. The controls were then compared pixel by pixel with libadwaita 1.10.0 itself, rendered with GTK 4.24.1 and Adwaita Sans by tools/adwaita-reference, whose values are cited as measured from it.

// adwAccent is libadwaita's default blue accent, accent_bg_color.
var adwAccent = ui.Hex("#3584e4")

// adwInk is the light theme's window_fg_color without its 80% opacity.
var adwInk = ui.RGB(0, 0, 6)

func adwTheme(dark bool, accent ui.Color) *ui.Theme {
	custom := accent.A != 0
	if !custom {
		accent = adwAccent
	}
	t := pick(dark, ui.LightTheme(), ui.DarkTheme())
	t.Accent, t.AccentText = accent, white
	t.AccentHover = accent.Mix(white, 0.1)    // a 10% currentColor layer
	t.AccentPressed = accent.Mix(adwInk, 0.2) // a 20% shade layer
	t.Radius, t.Spacing = 9, 4                // $button_radius
	t.FontSize = 14.67                        // Adwaita Sans 11 at 96 DPI: font-name in gsettings-desktop-schemas 51.0's org.gnome.desktop.interface
	if dark {
		t.Background = ui.Hex("#222226")
		t.Text = white
		t.TextMuted = t.Background.Mix(white, 0.55) // --dim-opacity
		t.Border = t.Background.Mix(white, 0.15)    // $border_color
		t.Surface = t.Background.Mix(white, 0.10)   // $button_color
		t.SurfaceHover = t.Background.Mix(white, 0.15)
		t.SurfacePressed = t.Background.Mix(white, 0.30)
		t.Danger, t.Warning, t.Success = ui.Hex("#c01c28"), ui.Hex("#cd9309"), ui.Hex("#26a269")
		// accent_color at 50%. The dark theme's accent_color is the accent with its Oklab lightness raised to 0.85, which for #3584e4 is outside sRGB (128.6, 208.4, 310.8: its blue past 255), and GTK 4.24 mixes it at 50% unclamped: libadwaita 1.10.0's focus ring reads (92, 132, 185) over an entry's (57, 57, 60) (measured by tools/adwaita-reference), where the clamped #81d0ff at 50% gives (93, 132, 158). The blue's 185 needs 64.3% of 255; these values give GTK's ring over the entry's and a button's faces, and come within 3 levels of it over the window.
		t.Focus = ui.RGBA(112, 175, 255, 0.643)
	} else {
		t.Background = ui.Hex("#fafafb")
		t.Text = t.Background.Mix(adwInk, 0.8)
		t.TextMuted = t.Background.Mix(adwInk, 0.44)
		t.Border = t.Background.Mix(adwInk, 0.12)
		t.Surface = t.Background.Mix(adwInk, 0.08)
		t.SurfaceHover = t.Background.Mix(adwInk, 0.12)
		t.SurfacePressed = t.Background.Mix(adwInk, 0.24)
		t.Danger, t.Warning, t.Success = ui.Hex("#e01b24"), ui.Hex("#e5a50a"), ui.Hex("#2ec27e")
		t.Focus = ui.RGBA(4, 97, 190, 0.5)
	}
	if custom {
		// accent_color is the accent clamped in Oklab lightness, darker in the light theme and lighter in the dark one; a mix comes close.
		t.Focus = pick(dark, accent.Mix(black, 0.2), accent.Mix(white, 0.45)).Alpha(0.5)
	}
	t.Selection = accent.Alpha(pick(dark, float32(0.25), 0.4))
	// An overlay scroll bar's slider while the pointer is not on the bar (_scrolling.scss): currentColor at 20%, 3 pixels wide. MyGo shows the bar while the pointer is over the view, not only over the bar, so the 8-pixel slider over its 10% trough that GTK widens it to under the pointer, and its 1-pixel outline, are not drawn.
	t.Scrollbar = pick(dark, ui.RGBA(0, 0, 6, 0.16), ui.RGBA(255, 255, 255, 0.2))
	t.ScrollbarWidth = 3
	// tooltip.background: the shade at 80% with white text, in light and dark.
	t.Inverse, t.InverseText = ui.RGBA(0, 0, 6, 0.8), white
	return t
}

// adwTrough returns $trough_color and its hover and active shades: currentColor at 15%, 20% and 25%, the unchecked box of a check button and a switch that is off.
func adwTrough(t *ui.Theme, hovered, pressed bool) ui.Color {
	ink := pick(t.Dark, adwInk, white)
	alpha := pick(t.Dark, float32(0.12), 0.15)
	switch {
	case pressed:
		alpha = pick(t.Dark, float32(0.20), 0.25)
	case hovered:
		alpha = pick(t.Dark, float32(0.16), 0.20)
	}
	return t.Background.Mix(ink, alpha)
}

// adwAccentFill returns the accent background of a checked control under the pointer and pressed.
func adwAccentFill(t *ui.Theme, hovered, pressed bool) ui.Color {
	switch {
	case pressed:
		return t.AccentPressed
	case hovered:
		return t.AccentHover
	}
	return t.Accent
}

// adwFocusRing draws libadwaita's focus ring on e in place of MyGo's: while e has the visible keyboard focus, 2 pixels of $focus_border_color (the theme's Focus) inside its edge, as the focus-ring mixin's outline-width 2px and outline-offset -2px draw it (_drawing.scss) on buttons and check buttons, around corners of radius (at most half e's height). over, nil for nothing, paints what e shows over its children, under the ring.
func adwFocusRing(t *ui.Theme, e ui.Element, radius float32, over func(p *ui.Painter, r ui.Rect)) {
	e.FocusRing(false).DrawOver(func(p *ui.Painter, r ui.Rect) {
		if over != nil {
			over(p, r)
		}
		if e.FocusVisible() {
			p.Stroke(r, t.Focus, min(radius, r.H/2), 2)
		}
	})
}

// adwButton draws a GtkButton with a label (.text-button): 34 pixels high with padding 5px 17px, corners of 9 and bold text, its focus ring inside its edge; .suggested-action when primary, whose ring (%opaque_button's, outside with a 1-pixel gap) is MyGo's.
func adwButton(c *ui.Context, b ui.Element, primary bool) {
	t := c.Theme()
	b.MinHeight(34).Padding(5, 17).Gap(4).Radius(9).FontWeight(700).Justify(ui.Center).AlignItems(ui.Center)
	pressed, hovered := b.Pressed(), b.Hovered()
	if primary {
		b.Background(adwAccentFill(t, hovered, pressed)).TextColor(t.AccentText)
		return
	}
	adwFocusRing(t, b, 9, nil)
	// $button_color, $button_hover_color and $button_active_color are currentColor at 10%, 15% and 30%, translucent, so a button shows what is behind it: over a boxed list's card its face is the card's color mixed, 235 in the light theme and 72 in the dark one where over the window it is 230 and 57 (measured from libadwaita 1.10.0 by tools/adwaita-reference, an AdwActionRow's suffix and the inputs set's buttons). The text takes its color over the face as it is over the window.
	share := float32(0.10)
	switch {
	case pressed:
		share = 0.30
	case hovered:
		share = 0.15
	}
	b.Background(adwInkLayer(t, share)).TextColor(adwTextOn(t, adwInkAlpha(t, t.Background, share)))
}

// adwInkLayer is currentColor at share as a translucent color, what adwInkAlpha lays over a background: white in the dark theme, and in the light one window_fg_color, the shade at 80%.
func adwInkLayer(t *ui.Theme, share float32) ui.Color {
	if t.Dark {
		return white.Alpha(share)
	}
	return adwInk.Alpha(share * 0.8)
}

// adwCheckbox draws a GtkCheckButton: a 20-pixel indicator (14 and padding 3) with corners of 6 and an inset 2-pixel trough border, 4 pixels from its label, inside a row padded by 4 (.text-button) with corners of 9 that shows the focus ring inside its edge; checked, the accent with libadwaita's check mark in accent_fg_color.
func adwCheckbox(c *ui.Context, row, box ui.Element, on bool) {
	t := c.Theme()
	row.Gap(4).Padding(4).Radius(9)
	adwFocusRing(t, row, 9, nil)
	hovered, pressed := row.Hovered(), row.Pressed()
	box.Size(20, 20).Radius(6)
	switch {
	case on:
		box.Background(adwAccentFill(t, hovered, pressed))
	case pressed:
		box.Background(adwTrough(t, false, true))
	default:
		box.Border(2, adwTrough(t, hovered, false))
	}
	box.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if !on {
			return
		}
		// assets/check.svg, which libadwaita 1.10 draws in place of the check-symbolic images: a 14-pixel icon, 3 pixels in, filled along a 2-pixel stroke with round ends through (2,7), (6,11) and (12,3), the centers of its arcs.
		s := r.W / 20
		var path ui.Path
		path.MoveTo(r.X+5*s, r.Y+10*s).LineTo(r.X+9*s, r.Y+14*s).LineTo(r.X+15*s, r.Y+6*s)
		p.StrokePath(&path, 2*s, t.AccentText)
	})
}

// adwSwitch draws a GtkSwitch: 46 by 26 pixels, a capsule track padded by 3 around a 20-pixel round slider with a soft shadow, moving in 100 milliseconds; the slider is adwSliderFill off and white when checked or under the pointer. Disabled, the slider casts no shadow (switch > slider:disabled), as measured from libadwaita 1.10.0 (GTK 4.24.1) by tools/adwaita-reference.
func adwSwitch(c *ui.Context, sw ui.Element, on bool) {
	adwSwitchTrack(c, c.Theme(), sw, sw.Hovered(), sw.Pressed(), on)
}

// adwSwitchTrack draws a GtkSwitch on sw, its trough laid over t's Background, under the pointer if hovered and pressed if pressed: adwSwitch's look for a switch whose state is another element's, as AdwSwitchRow's.
func adwSwitchTrack(c *ui.Context, t *ui.Theme, sw ui.Element, hovered, pressed, on bool) {
	pos := knob(sw, on, 100*time.Millisecond)
	at := along(c, pos)
	sw.Size(46, 26).Radius(13).Background(adwTrough(t, hovered, pressed).Mix(adwAccentFill(t, hovered, pressed), pos))
	slider := adwSliderFill(t)
	if on || hovered || pressed {
		// switch:checked > slider and switch:hover > slider: white, as libadwaita 1.10.0 (GTK 4.24.1) renders a checked switch's slider in the dark theme too (measured by tools/adwaita-reference).
		slider = white
	}
	sw.Draw(func(p *ui.Painter, r ui.Rect) {
		const in, d = 3, 20
		k := ui.Rect{X: r.X + in + at*(r.W-d-2*in), Y: r.Y + (r.H-d)/2, W: d, H: d}
		if !sw.IsDisabled() {
			p.Shadow(k, d/2, 0, 2, 4, 0, ui.RGBA(0, 0, 6, 0.2))
		}
		p.Fill(k, slider, d/2)
	})
}

// adwSliderFill is the slider of a switch and a scale away from the pointer, color-mix(in srgb, white 80%, var(--view-bg-color)): white in the light theme and #d2d2d2 over the dark view_bg_color #1d1d20, as measured from libadwaita 1.10.0 (GTK 4.24.1) by tools/adwaita-reference; under the pointer it is white.
func adwSliderFill(t *ui.Theme) ui.Color {
	return adwListViewBg(t).Mix(white, 0.8)
}

// ---- more of libadwaita's controls, from the same stylesheet ----

// adwInkAlpha lays currentColor at share over bg, as libadwaita's color-mix(in srgb, currentColor share, transparent) does: window_fg is 80% opaque in the light theme.
func adwInkAlpha(t *ui.Theme, bg ui.Color, share float32) ui.Color {
	if t.Dark {
		return bg.Mix(white, share)
	}
	return bg.Mix(adwInk, share*0.8)
}

// adwRadio draws a GtkCheckButton in a group: a 20-pixel circle with an inset 2-pixel trough ring, 4 pixels from its label, inside a padded row with round ends (.grouped) that shows the focus ring inside its edge; chosen, the accent with libadwaita's 8-pixel white bullet (assets/bullet.svg).
func adwRadio(c *ui.Context, row, dot ui.Element, on bool) {
	t := c.Theme()
	row.Gap(4).Padding(4).Radius(999)
	adwFocusRing(t, row, 999, nil)
	adwRadioDot(c, dot, on, row.Hovered(), row.Pressed())
}

// adwRadioDot draws a grouped GtkCheckButton's circle on dot, as adwRadio describes, under the pointer while hovered and pressed while pressed.
func adwRadioDot(c *ui.Context, dot ui.Element, on, hovered, pressed bool) {
	t := c.Theme()
	dot.Size(20, 20).Radius(10)
	switch {
	case on:
		dot.Background(adwAccentFill(t, hovered, pressed))
	case pressed:
		dot.Background(adwTrough(t, false, true))
	default:
		dot.Border(2, adwTrough(t, hovered, false))
	}
	if on {
		dot.DrawOver(func(p *ui.Painter, r ui.Rect) {
			const d = 8
			p.Fill(ui.Rect{X: r.X + (r.W-d)/2, Y: r.Y + (r.H-d)/2, W: d, H: d}, t.AccentText, d/2)
		})
	}
}

// adwToggle draws a GtkToggleButton: a button, darker while checked (currentColor at 30%, 35% under the pointer and 40% pressed).
func adwToggle(c *ui.Context, b ui.Element, on bool) {
	adwButton(c, b, false)
	if !on {
		return
	}
	t := c.Theme()
	face := adwInkAlpha(t, t.Background, adwToggleShare(b, true))
	b.Background(face).TextColor(adwTextOn(t, face))
}

// adwToggleShare returns the share of currentColor a toggle button lays over what is behind it: 10%, 15% under the pointer and 30% pressed ($button_color and its hover and active shades), and checked 30%, 35% and 40%.
func adwToggleShare(b ui.Element, on bool) float32 {
	switch {
	case on && b.Pressed():
		return 0.40
	case on && b.Hovered():
		return 0.35
	case on:
		return 0.30
	case b.Pressed():
		return 0.30
	case b.Hovered():
		return 0.15
	}
	return 0.10
}

// adwField draws a GtkEntry: 34 pixels high with padding 0 9px, corners of 9, filled with currentColor at 10% and no border, and a 2-pixel focus ring inside its edge while it has the focus; a search entry is the same. A multi-line text view sits in a frame: view_bg_color inside a 1-pixel border, corners of 12. The fill is translucent, as libadwaita's is, so an entry shows what is behind it: over the window 230 in the light theme and 57 in the dark one, in an AdwActionRow over a boxed list's card 235 and 72 (measured from libadwaita 1.10.0 (GTK 4.24.1) by tools/adwaita-reference, its inputs and overcard sets). An entry's placeholder is .dimmed, window_fg_color at 55% over the fill: 129 and 166 over the window, 132 and 172 over the card, where MyGo's TextMuted, opaque over the window, read 140 and 156.
func adwField(c *ui.Context, f ui.Element, kind fieldKind) {
	t := c.Theme()
	f.Border(0, ui.Color{})
	radius := float32(9)
	if kind == multiLine {
		radius = 12
		f.Padding(8, 9).Background(pick(t.Dark, white, ui.Hex("#1d1d20")))
	} else {
		f.MinHeight(34).Padding(0, 9).Background(adwInkLayer(t, 0.10)).PlaceholderColor(adwInkLayer(t, 0.55))
	}
	f.Radius(radius)
	focused := f.Focused()
	f.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if kind == multiLine {
			p.Stroke(r, t.Border, radius, 1)
		}
		if focused {
			p.Stroke(r, t.Focus, radius, 2)
		}
	})
}

// adwSelect draws a GtkDropDown's button: a regular button, darker while its list is open, with its focus ring inside its edge, showing the choice and the 16-pixel pan-down arrow 6 pixels after it, as wide as these need and no wider: GtkDropDown sets no least width (a "Medium" drop-down is 101 pixels wide in Adwaita Sans, measured from libadwaita 1.10.0 (GTK 4.24.1) by tools/adwaita-reference).
func adwSelect(c *ui.Context, trigger ui.Element, selected string) {
	t := c.Theme()
	b := trigger
	b.MinHeight(34).Padding(5, 10).Gap(6).Radius(9).FontWeight(700).Justify(ui.SpaceBetween).AlignItems(ui.Center).TextColor(t.Text)
	face := t.Surface
	switch {
	case b.Pressed():
		face = t.SurfacePressed
	case b.Hovered():
		face = t.SurfaceHover
	}
	fg := adwTextOn(t, face)
	b.Background(face).TextColor(fg)
	adwFocusRing(t, b, 9, nil)
	b.Children(func() {
		ui.Text(c, selected).SingleLine()
		ui.Box(c).Size(16, 16).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
			// adwaita-icon-theme 51.0's pan-down-symbolic: a chevron filled along a 2-pixel stroke with round ends from (4,7) down to its tip at (8,11) and up to (12,7).
			var path ui.Path
			path.MoveTo(r.X+4, r.Y+7).LineTo(r.X+8, r.Y+11).LineTo(r.X+12, r.Y+7)
			p.StrokePath(&path, 2, fg)
		})
	})
}

// adwMenu draws a dropdown's popover: popover_bg_color, corners of 15 and popover > contents's box-shadow (_popovers.scss; libadwaita 1.10 replaced the 1-pixel border of black at 14% with its first shadow): a 1-pixel ring of black at 5% around it and two soft shadows; its list padded 6 pixels.
func adwMenu(c *ui.Context, panel ui.Element) {
	t := c.Theme()
	panel.Margin(6, 0, 0, 0).Padding(6).Radius(15).Background(pick(t.Dark, white, ui.Hex("#36363a"))).Border(0, ui.Color{}).
		Shadow(0, 0, 0, 1, ui.RGBA(0, 0, 0, 0.05)).Shadow(0, 1, 5, 1, ui.RGBA(0, 0, 0, 0.09)).Shadow(0, 2, 14, 3, ui.RGBA(0, 0, 0, 0.05))
}

// adwMenuItem draws a row of the list: padding 9px 12px, corners of 9, currentColor at 10% under the pointer and 19% pressed (_dropdowns.scss and _menus.scss), and object-select-symbolic 6 pixels after the chosen option's label.
func adwMenuItem(c *ui.Context, item ui.Element, label string, chosen bool) {
	t := c.Theme()
	panel := pick(t.Dark, white, ui.Hex("#36363a"))
	item.Padding(9, 12).Gap(6).Radius(9).AlignItems(ui.Center).MinWidth(100)
	switch {
	case item.Pressed():
		item.Background(adwInkAlpha(t, panel, 0.19))
	case item.Highlighted():
		item.Background(adwInkAlpha(t, panel, 0.10))
	}
	item.Children(func() {
		ui.Text(c, label).SingleLine()
		mark := ui.Box(c).Size(16, 16).Shrink(0)
		if chosen {
			mark.Draw(func(p *ui.Painter, r ui.Rect) {
				// adwaita-icon-theme 51.0's object-select-symbolic: a check filled along a 2-pixel stroke with round ends through (3,9), (6,12) and (13,4), as GTK 4.24's own fallback strokes it.
				var path ui.Path
				path.MoveTo(r.X+3, r.Y+9).LineTo(r.X+6, r.Y+12).LineTo(r.X+13, r.Y+4)
				p.StrokePath(&path, 2, t.Text)
			})
		}
	})
}

// adwScaleTravel is how far in from the trough's ends a scale's slider center stops: GtkRange moves the slider's allocation, 20 pixels with margins of -8 (_scale.scss), so 4 wide, along the trough, and its center stays 2 pixels in from the ends, where a scale's first and last marks are drawn, as measured from libadwaita 1.10.0 (GTK 4.24.1) by tools/adwaita-reference.
const adwScaleTravel = 2

// adwScaleSlider returns the center of a scale's slider at along (from 0 to 1, from the trough's start in the reading direction) over a trough from left, width wide: GtkRange allocates the slider at a whole pixel, the travel times along rounded down, so a 156-pixel trough puts it 62 pixels in at 40% and 78 at 50%, as measured from libadwaita 1.10.0 (GTK 4.24.1) by tools/adwaita-reference.
func adwScaleSlider(left, width, along float32) float32 {
	return left + adwScaleTravel + float32(math.Floor(float64(along*(width-2*adwScaleTravel))))
}

// adwSlider draws a GtkScale: 34 pixels high, padded 12 pixels around a 10-pixel range (min-height 10px), a 4-pixel trough of $trough_color centered in it and filled with the accent, and a 20-pixel round slider with a hairline (box-shadow 0 0 0 1px rgb(0 0 6 / 10%)) and a soft shadow, which turns white under the pointer; disabled, the slider casts no shadow. The sizes were measured from libadwaita 1.10.0 (GTK 4.24.1) by tools/adwaita-reference: GtkScale with draw-value off is 34 pixels high, the trough 15 pixels from its top.
func adwSlider(c *ui.Context, s ui.Element, frac float32) {
	t := c.Theme()
	const d = 20
	hovered, pressed := s.Hovered(), s.Pressed()
	// The pointer and the slider's center both travel the trough but 2 pixels at each end (adwScaleTravel), inside the scale's padding; the slider overhangs its ends, as GtkRange draws it.
	s.Height(12 + 10 + 12).PaddingX(12 + adwScaleTravel).MinWidth(60).FocusRing(false)
	trough := adwTrough(t, false, false)
	fill := t.Accent
	knobFill := adwSliderFill(t)
	if hovered || pressed {
		trough, fill, knobFill = adwTrough(t, true, false), t.AccentHover, white
	}
	at := along(c, frac)
	s.Draw(func(p *ui.Painter, r ui.Rect) {
		ty := r.Y + (r.H-4)/2
		left := r.X + 12
		width := r.W - 24
		p.Fill(ui.Rect{X: left, Y: ty, W: width, H: 4}, trough, 2)
		x := adwScaleSlider(left, width, at)
		bar := ui.Rect{X: left, Y: ty, W: x - left, H: 4}
		if c.Direction() == ui.RightToLeft {
			bar = ui.Rect{X: x, Y: ty, W: left + width - x, H: 4}
		}
		p.Fill(bar, fill, 2)
		k := ui.Rect{X: x - d/2, Y: r.Y + (r.H-d)/2, W: d, H: d}
		if !s.IsDisabled() {
			p.Shadow(k, d/2, 0, 2, 4, 0, ui.RGBA(0, 0, 6, 0.2))
		}
		// The hairline is ---slider-border-color, rgb(0 0 6 / 10%) in both themes, over what lies under it: the accent fill, the trough or the window.
		p.Fill(ui.Rect{X: k.X - 1, Y: k.Y - 1, W: d + 2, H: d + 2}, ui.RGBA(0, 0, 6, 0.1), d/2+1)
		p.Fill(k, knobFill, d/2)
		if s.FocusVisible() {
			p.Stroke(ui.Rect{X: k.X - 2, Y: k.Y - 2, W: d + 4, H: d + 4}, t.Focus, d/2+2, 2)
		}
	})
}

// adwProgress draws a GtkProgressBar: an 8-pixel trough of $trough_color with round ends, filled with the accent.
func adwProgress(c *ui.Context, e ui.Element, value float64) {
	t := c.Theme()
	e.Height(8).MinWidth(40).Radius(4).Background(adwTrough(t, false, false)).Clip()
	e.Draw(func(p *ui.Painter, r ui.Rect) { drawBar(c, p, r, value, t.Accent, 4) })
}

// adwSpinner draws an AdwSpinner at its natural 16 pixels: a 2-pixel ring of the text color at 15%, and an arc of it that turns once in 1.2 seconds while it extends and contracts, as adw-spinner-paintable.c draws it. libadwaita 1.10 strokes it an eighth of its diameter wide at the widget's icon weight of 400 (calculate_line_width), where 1.7 stroked 2.5 pixels.
func adwSpinner(c *ui.Context, e ui.Element) {
	t := c.Theme()
	e.Size(16, 16)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		p.AnimationFrame()
		const lw = 2
		rad := min(r.W, r.H)/2 - lw/2
		cx, cy := r.X+r.W/2, r.Y+r.H/2
		var ring ui.Path
		ring.Circle(cx, cy, rad)
		p.StrokePath(&ring, lw, t.Background.Mix(t.Text, 0.15))
		start, end := adwSpinnerArc(float64(p.Now().UnixMilli()))
		// The paintable's angles run clockwise from 3 o'clock; strokeArc's from 12.
		strokeArc(p, cx, cy, rad, end+math.Pi/2, start-end, lw, t.Text)
	})
}

// adwSpinnerArc returns the angles of the spinner's arc ms milliseconds in, as adw-spinner-paintable.c computes them.
func adwSpinnerArc(ms float64) (start, end float64) {
	const (
		minArc, maxArc = 0.015 * math.Pi, 0.9 * math.Pi
		startAngle     = 0.35 * math.Pi
		cycle          = 2.65 * math.Pi
	)
	base := ms / 1200 * 2 * math.Pi
	a := math.Mod(base, cycle)
	ease := func(x float64) float64 { return -(math.Cos(math.Pi*x) - 1) / 2 }
	ts := 1.0
	if a < 1.1*math.Pi {
		ts = ease(a / (1.1 * math.Pi))
	}
	te := 0.0
	switch {
	case a >= 1.75*math.Pi:
		te = 1
	case a >= 0.4*math.Pi:
		te = ease((a - 0.4*math.Pi) / (1.35 * math.Pi))
	}
	shift := base - a*maxArc/cycle + startAngle
	return shift + minArc + (maxArc-minArc)*ts, shift + (maxArc-minArc)*te
}

// adwSegmentedTrack draws an AdwToggleGroup: padded 3 pixels, currentColor at 10%, corners of 9.
func adwSegmentedTrack(c *ui.Context, track ui.Element) {
	t := c.Theme()
	track.Padding(3).Radius(9).Background(adwInkAlpha(t, t.Background, 0.10)).AlignItems(ui.Stretch)
}

// adwSegment draws a toggle of the group (.text-button): 28 pixels high, padded 8 at the sides, bold, corners of 6, a light fill with a shadow while active, currentColor at 7% under the pointer and 16% pressed, and a 1-pixel separator between toggles that are not active.
func adwSegment(c *ui.Context, seg ui.Element, chosen, separator bool) {
	t := c.Theme()
	track := adwInkAlpha(t, t.Background, 0.10)
	seg.MinHeight(28).MinWidth(34).Padding(0, 8).Radius(6).FontWeight(700).Justify(ui.Center).AlignItems(ui.Center).TextColor(t.Text)
	switch {
	case chosen:
		seg.Background(pick(t.Dark, white, track.Mix(white, 0.2))).Shadow(0, 1, 3, 1, ui.RGBA(0, 0, 6, 0.07)).Shadow(0, 2, 6, 2, ui.RGBA(0, 0, 6, 0.03))
	case seg.Pressed():
		seg.Background(adwInkAlpha(t, track, 0.16))
	case seg.Hovered():
		seg.Background(adwInkAlpha(t, track, 0.07))
	}
	if separator && !seg.Hovered() {
		line := t.Border
		seg.Draw(func(p *ui.Painter, r ui.Rect) {
			x := r.X - 1.5
			if c.Direction() == ui.RightToLeft {
				x = r.X + r.W + 0.5
			}
			p.Fill(ui.Rect{X: x, Y: r.Y + 3, W: 1, H: r.H - 6}, line, 0)
		})
	}
}

// adwTabList draws a GtkNotebook's header: tabs in a row padded 4 pixels, over a 1-pixel border on the page's side.
func adwTabList(c *ui.Context, list ui.Element) {
	t := c.Theme()
	list.PaddingX(4).AlignItems(ui.Stretch).BorderWidth(0, 0, 1, 0).BorderColor(t.Border)
}

// adwTab draws a notebook tab: at least 30 pixels, padding 3px 12px 4px, 4 pixels apart, currentColor at 7% under the pointer, and a 4-pixel accent line along its bottom while chosen.
func adwTab(c *ui.Context, tab ui.Element, chosen bool) {
	t := c.Theme()
	tab.MinHeight(30).MinWidth(30).Padding(3, 12, 4, 12).MarginX(4).Justify(ui.Center).AlignItems(ui.Center).TextColor(t.Text)
	if tab.Hovered() || tab.Pressed() {
		tab.Background(adwInkAlpha(t, t.Background, 0.07))
	}
	if chosen || tab.Hovered() {
		line := t.Accent
		if !chosen {
			line = t.Border
		}
		tab.Draw(func(p *ui.Painter, r ui.Rect) { p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 4, W: r.W, H: 4}, line, 0) })
	}
}

// ---- dialogs: AdwDialog and AdwAlertDialog ----

// adwSheetShadow casts floating-sheet > sheet's shadows (_dialogs.scss; libadwaita 1.10 made them smaller and centered), the same in light and dark: two soft ones and a 1-pixel ring of black at 5%.
func adwSheetShadow(panel ui.Element) {
	panel.Shadow(0, 0, 14, 2, ui.RGBA(0, 0, 6, 0.03)).Shadow(0, 0, 5, 2, ui.RGBA(0, 0, 6, 0.10)).Shadow(0, 0, 0, 1, ui.RGBA(0, 0, 0, 0.05))
}

// adwDimming fills the window behind a dialog with the shade color at twice its alpha: 14% in light, 50% in dark.
func adwDimming(t *ui.Theme, back ui.Element) {
	back.Background(ui.RGBA(0, 0, 6, pick(t.Dark, float32(0.14), 0.5)))
}

// adwDialog draws a floating AdwDialog: window_bg_color, corners of 15, a hairline inside its edge and the sheet's shadows, over the dimmed window.
func adwDialog(c *ui.Context, back, panel ui.Element) {
	t := c.Theme()
	adwDimming(t, back)
	panel.Padding(24).Gap(12).Radius(15).Background(t.Background).Border(1, t.Background.Mix(white, 0.07))
	adwSheetShadow(panel)
}

// adwAlert draws an AdwAlertDialog: dialog_bg_color, corners of 18, 300 to 372 pixels wide, with the sheet's 1-pixel $window_outline_color (white at 7%) inside its edge; the heading in .title-2 (weight 800, 136%) and the body (.body, line-height 140%) 10 pixels below it, centered, padded 32 above and 24 at the sides; and the responses in a row below, padded 12 above and 24 at the sides and below, equal widths 12 apart, at least 60 pixels wide past their padding while they all fit so, 44 pixels high with corners of 12, in the order they were added, so the default comes last; every response but button styled, the suggested or the destructive one, which keeps its own fill, takes the response face. In a window at most 360 pixels high the dialog is .short, as adw-dialog.c's breakpoint makes it: the heading in .title-4 (weight 700, 118%), 6 pixels above the body, padded 24 above and 18 at the sides, and the responses padded 9 above and 18 at the sides and below, 40 pixels high, with no least width (_message-dialog.scss).
func adwAlert(c *ui.Context, back, panel ui.Element, title, message string, n, styled int, button func(int) ui.Element) {
	t := c.Theme()
	adwDimming(t, back)
	w, h := c.Size()
	top, side, gap, size, weight, padY, padX := float32(32), float32(24), float32(10), float32(1.36), 800, float32(10), float32(20)
	above, least := float32(12), float32(60)
	if h <= 360 {
		top, side, gap, size, weight, padY, padX = 24, 18, 6, 1.18, 700, 8, 16
		above, least = 9, 0
	}
	// libadwaita stacks the responses when they are wider side by side than the dialog can be (adw-alert-dialog.c, measure_responses_do's wide_nat against DIALOG_MAX_WIDTH); this row never stacks them, so the least width applies only while the responses fit in the widest alert the window allows, and otherwise they share the row as before rather than overflow it.
	if float32(n)*(least+2*padX)+12*float32(n-1)+2*side > min(372, w-40) {
		least = 0
	}
	bg := pick(t.Dark, ui.Hex("#fafafb"), ui.Hex("#36363a"))
	panel.MinWidth(300).MaxWidth(372).Radius(18).Background(bg).AlignItems(ui.Stretch)
	adwSheetShadow(panel)
	// libadwaita 1.10.0 keeps the sheet's outline on an alert in normal contrast and removes it in high contrast, the other way round from 1.7.6: cb3219bc ("stylesheet: Stop using $contrast") turned `@if $contrast != 'high' { outline: none }` into `@media (prefers-contrast: more) { outline: none }` (_message-dialog.scss), as its doc/images/alert-dialog-dark.png shows. GTK draws an outline over the content and outside the layout, as DrawOver does.
	outline := bg.Mix(white, 0.07)
	panel.DrawOver(func(p *ui.Painter, r ui.Rect) { p.Stroke(r, outline, 18, 1) })
	panel.Children(func() {
		ui.Column(c).Padding(top, side, 9, side).Gap(gap).AlignItems(ui.Stretch).Children(func() {
			ui.Text(c, title).FontSize(t.FontSize * size).FontWeight(weight).TextAlign(ui.Center)
			if message != "" {
				ui.Text(c, message).TextAlign(ui.Center).LineHeight(1.4)
			}
		})
		// The response area is padded as the message area at the sides and below (RESPONSE_HORZ_PADDING_SHORT while short); a response's content is at least 24 pixels high and, unless short, 60 wide.
		ui.Row(c).Padding(above, side, side, side).Gap(12).Children(func() {
			for i := range n {
				b := button(i).MinHeight(24+2*padY).MinWidth(least+2*padX).Padding(padY, padX).Radius(12).Grow(1).Basis(0)
				if i != styled {
					// currentColor at 10%, 15% and 30%, over the alert's own background rather than the window's.
					share := float32(0.10)
					switch {
					case b.Pressed():
						share = 0.30
					case b.Hovered():
						share = 0.15
					}
					face := adwInkAlpha(t, bg, share)
					b.Background(face).TextColor(adwTextOn(t, face))
				}
			}
		})
	})
}

// adwTextOn returns window_fg_color over fill: libadwaita's text is 80% opaque in the light theme, so it is a little lighter over a darker button than over the window.
func adwTextOn(t *ui.Theme, fill ui.Color) ui.Color {
	if t.Dark {
		return white
	}
	return fill.Mix(adwInk, 0.8)
}

// ---- links, separators, badges, popovers and linked toggles ----

// adwLinkColor is libadwaita's link color, accent_color: the accent darkened in the light theme and lightened in the dark one. For the default blue, the colors are measured from links libadwaita 1.10.0 (GTK 4.24.1) rendered (tools/adwaita-reference), #357fca and #81d0ff, as 1.7.6 rendered them.
func adwLinkColor(t *ui.Theme) ui.Color {
	if t.Accent == adwAccent {
		return pick(t.Dark, ui.Hex("#357fca"), ui.Hex("#81d0ff"))
	}
	return pick(t.Dark, t.Accent.Mix(black, 0.2), t.Accent.Mix(white, 0.45))
}

// adwPopover draws a popover's contents: popover_bg_color, corners of 15, padding 8 and adwMenu's ring and two soft shadows, its contents right at the anchor's edge, as GtkPopover places them without its arrow (has-arrow false, as MyGo draws popovers): in libadwaita 1.10.0's render (tools/adwaita-reference's third set) the contents of the Details popover start at 104, where the Info button above it ends, in the light and the dark theme. adwMenu's margin of 6 is a dropdown's.
func adwPopover(c *ui.Context, panel ui.Element) {
	adwMenu(c, panel)
	panel.Margin(0).Padding(8).Gap(6)
}

// adwSeparator is a separator's color, $border_color: currentColor at 15%.
func adwSeparator(t *ui.Theme) ui.Color { return t.Border }

// adwBadge draws an indicator badge (indicatorbin.badge, _view-switcher.scss): a capsule at least 13 pixels high, currentColor at 40% with bold white text at 0.6rem, padded 4 at the sides. In a Sidebar it is the count AdwViewSwitcherSidebar (libadwaita 1.9+) puts at the end of an AdwSidebar item (view-switcher-sidebar .indicator, _sidebars.scss): a capsule at least 1.6em square for its text at font-size smaller (GTK's factor of 1.2, inferred), bold and white, padded 6 at the sides, over currentColor at 40% over sidebar_bg_color. MyGo's badge sizes its text itself, 2 below the theme's font size, so the text is not at 0.6rem or smaller, the sizes the capsules take.
func adwBadge(c *ui.Context, b ui.Element) {
	t := c.Theme()
	if navCurrentSidebar(c) != nil {
		size := t.FontSize / 1.2
		b.MinHeight(size*1.6).MinWidth(size*1.6).Padding(0, 6).Radius(999).Background(adwInkAlpha(t, adwSidebarBg(t), 0.4)).
			TextColor(white).FontSize(size).FontWeight(700).Border(0, ui.Color{})
		return
	}
	b.MinHeight(13).MinWidth(13).Padding(0, 4).Radius(7).Background(adwInkAlpha(t, t.Background, 0.4)).
		TextColor(white).FontSize(t.FontSize*0.6).FontWeight(700).Border(0, ui.Color{})
}

// adwToggleGroup draws a box.linked of toggle buttons: no gap between them and only the group's outer corners rounded, 9 pixels.
func adwToggleGroup(c *ui.Context, g ui.Element) {
	g.Gap(0).Padding(0).Radius(9).Background(ui.Transparent).Border(0, ui.Color{}).Clip()
}

// adwLinkedToggle draws a toggle button inside a box.linked: a toggle button whose corners the group rounds. Linked buttons overlap by a pixel (margin-left: -1px), so a line shows where two meet: the button's translucent fill laid over the one before it, drawn along the start of every button but the first. Measured from libadwaita 1.10.0 (GTK 4.24.1) by tools/adwaita-reference: after a checked Bold, Italic's seam is #afafb1 in the light theme, darker than both, and #747476 in the dark one, lighter than both.
func adwLinkedToggle(c *ui.Context, b ui.Element, on bool, first bool) {
	adwToggle(c, b, on)
	b.Radius(0)
	t := c.Theme()
	share := adwToggleShare(b, on)
	prev := adwToggleGroupFace(c, adwInkAlpha(t, t.Background, share))
	if first {
		return
	}
	seam := adwInkAlpha(t, prev, share)
	// The focus ring, inside the button's edge as on a toggle alone, keeps corners of 9 so that the group's rounded clip does not cut it at the group's corners; between two toggles libadwaita's are square (inferred: a toggle does not know whether it is the last).
	adwFocusRing(t, b, 9, func(p *ui.Painter, r ui.Rect) {
		x := r.X
		if c.Direction() == ui.RightToLeft {
			x = r.X + r.W - 1
		}
		p.Fill(ui.Rect{X: x, Y: r.Y, W: 1, H: r.H}, seam, 0)
	})
}
