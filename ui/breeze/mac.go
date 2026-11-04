package breeze

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// macOS 26 and later. Colors are AppKit's semantic colors (NSColor) resolved in the Aqua and Dark Aqua appearances on macOS 27; control fills and sizes are measured from NSButton and NSSwitch that AppKit rendered at 2x. A push button's face is secondarySystemFill and an unchecked box's or an off switch's systemFill, the same share of black or white over the window in both appearances.

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
