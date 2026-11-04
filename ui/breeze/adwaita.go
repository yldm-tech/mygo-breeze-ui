package breeze

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// GNOME: libadwaita 1.7 (GNOME 48), from its stylesheet (src/stylesheet/_colors.scss, _common.scss, widgets/_buttons.scss, _checks.scss and _switch.scss) and, for the accent and the switch's size, its C code and GTK's. libadwaita colors controls with currentColor at an opacity; those are resolved here over window_bg_color.

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
	t.FontSize = 14.67                        // Adwaita Sans 11 at 96 DPI
	if dark {
		t.Background = ui.Hex("#222226")
		t.Text = white
		t.TextMuted = t.Background.Mix(white, 0.55) // --dim-opacity
		t.Border = t.Background.Mix(white, 0.15)    // $border_color
		t.Surface = t.Background.Mix(white, 0.10)   // $button_color
		t.SurfaceHover = t.Background.Mix(white, 0.15)
		t.SurfacePressed = t.Background.Mix(white, 0.30)
		t.Danger, t.Warning, t.Success = ui.Hex("#c01c28"), ui.Hex("#cd9309"), ui.Hex("#26a269")
		t.Focus = ui.RGBA(129, 208, 255, 0.5) // accent_color at 50%
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

// adwButton draws a GtkButton with a label (.text-button): 34 pixels high with padding 5px 17px, corners of 9 and bold text; .suggested-action when primary.
func adwButton(c *ui.Context, b ui.Element, primary bool) {
	t := c.Theme()
	b.MinHeight(34).Padding(5, 17).Gap(4).Radius(9).FontWeight(700).Justify(ui.Center).AlignItems(ui.Center)
	pressed, hovered := b.Pressed(), b.Hovered()
	if primary {
		b.Background(adwAccentFill(t, hovered, pressed)).TextColor(t.AccentText)
		return
	}
	face := t.Surface
	switch {
	case pressed:
		face = t.SurfacePressed
	case hovered:
		face = t.SurfaceHover
	}
	b.Background(face).TextColor(t.Text)
}

// adwCheckbox draws a GtkCheckButton: a 20-pixel indicator with corners of 6 and an inset 2-pixel trough border, 4 pixels from its label, inside a row padded by 4 that shows the focus ring; checked, the accent with libadwaita's check mark in accent_fg_color.
func adwCheckbox(c *ui.Context, row, box ui.Element, on bool) {
	t := c.Theme()
	row.Gap(4).Padding(4).Radius(9).FocusRing(true)
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
		// check.svg's 14-pixel mark, 3 pixels in.
		s := r.W / 20
		var path ui.Path
		path.MoveTo(r.X+5*s, r.Y+10*s).LineTo(r.X+8.9*s, r.Y+13.9*s).LineTo(r.X+16.3*s, r.Y+5.9*s)
		p.StrokePath(&path, 2*s, t.AccentText)
	})
}

// adwSwitch draws a GtkSwitch: 46 by 26 pixels, a capsule track padded by 3 around a 20-pixel round slider with a soft shadow, moving in 100 milliseconds.
func adwSwitch(c *ui.Context, sw ui.Element, on bool) {
	t := c.Theme()
	hovered, pressed := sw.Hovered(), sw.Pressed()
	pos := knob(sw, on, 100*time.Millisecond)
	at := along(c, pos)
	sw.Size(46, 26).Radius(13).Background(adwTrough(t, hovered, pressed).Mix(adwAccentFill(t, hovered, pressed), pos))
	slider := white
	if t.Dark && !on && !hovered && !pressed {
		slider = t.Background.Mix(white, 0.8) // $slider_color, white over view_bg_color
	}
	sw.Draw(func(p *ui.Painter, r ui.Rect) {
		const in, d = 3, 20
		k := ui.Rect{X: r.X + in + at*(r.W-d-2*in), Y: r.Y + (r.H-d)/2, W: d, H: d}
		p.Shadow(k, d/2, 0, 2, 4, 0, ui.RGBA(0, 0, 6, 0.2))
		p.Fill(k, slider, d/2)
	})
}
