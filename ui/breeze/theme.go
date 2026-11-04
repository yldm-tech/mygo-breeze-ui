package breeze

import "github.com/egoist/mygo/ui"

var (
	black = ui.RGB(0, 0, 0)
	white = ui.RGB(255, 255, 255)
)

// tint lays the platform's neutral fill over bg: black in a light theme and white in a dark one, alpha strong. Platforms draw control faces as such translucent layers, so a face keeps its contrast over any background.
func tint(bg ui.Color, dark bool, alpha float32) ui.Color {
	if dark {
		return bg.Mix(white, alpha)
	}
	return bg.Mix(black, alpha)
}

// Theme returns the theme of p's desktop, light or dark as c's window is, for the widgets of package ui and for elements of your own. Like MyGo's default theme it follows the desktop's preferences: the accent color the user chose, higher contrast and the text size. Use sets it for the window.
func Theme(c *ui.Context, p Platform) *ui.Theme {
	if p == Auto {
		p = Current()
	}
	prefs := c.Preferences()
	dark := c.Theme().Dark
	accent := prefs.Accent
	if accent.A > 0 {
		accent.A = 255
	}
	var t *ui.Theme
	switch p {
	case Windows:
		t = fluentTheme(dark, accent)
	case Linux:
		t = adwTheme(dark, accent)
	default:
		t = macTheme(dark, accent)
	}
	if prefs.HighContrast {
		t.Border = t.Border.Mix(t.Text, 0.45)
		t.TextMuted = t.TextMuted.Mix(t.Text, 0.4)
		t.Focus.A = 255
	}
	if s := prefs.TextScale; s > 0 && s != 1 {
		t.FontSize *= s
	}
	return t
}

// Use gives c's window the theme of the platform whose controls are drawn (PlatformOf), for package ui's widgets too. Call it at the start of the view, before the widgets it should style.
func Use(c *ui.Context) { c.SetTheme(Theme(c, PlatformOf(c))) }

// luminance returns how light c is, from 0 to 1, weighing its channels as the eye does.
func luminance(c ui.Color) float32 {
	return (0.2126*float32(c.R) + 0.7152*float32(c.G) + 0.0722*float32(c.B)) / 255
}

// textOn returns the text color that stands out on fill: white, or near black on a light fill such as yellow.
func textOn(fill ui.Color) ui.Color {
	if luminance(fill) > 0.6 {
		return ui.RGB(24, 24, 27)
	}
	return white
}

// pick returns light in a light theme and dark in a dark one.
func pick[T any](isDark bool, light, dark T) T {
	if isDark {
		return dark
	}
	return light
}
