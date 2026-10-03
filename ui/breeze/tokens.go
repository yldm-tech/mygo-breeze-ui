package breeze

import "github.com/egoist/mygo/ui"

// RadiusToken names the standard corner radius tokens.
type RadiusToken float32

const (
	RadiusNone RadiusToken = 0
	RadiusSM   RadiusToken = 4
	RadiusMD   RadiusToken = 6
	RadiusLG   RadiusToken = 10
	RadiusXL   RadiusToken = 14
	RadiusFull RadiusToken = 999
)

// TextToken names the standard type scale in DIPs.
type TextToken float32

const (
	TextXS   TextToken = 11
	TextSM   TextToken = 13
	TextBase TextToken = 14
	TextLG   TextToken = 18
	TextXL   TextToken = 24
	Text2XL  TextToken = 30
)

// Colors are semantic colors used by components. They are derived from the
// active MyGo theme, so components follow light/dark appearance changes.
type Colors struct {
	Background ui.Color
	Surface    ui.Color
	Hover      ui.Color
	Pressed    ui.Color
	Border     ui.Color
	Text       ui.Color
	Muted      ui.Color
	Accent     ui.Color
	AccentText ui.Color
	Danger     ui.Color
	Focus      ui.Color
}

// Palette returns the semantic colors of the current frame.
func Palette(c *ui.Context) Colors {
	t := c.Theme()
	return Colors{
		Background: t.Background,
		Surface:    t.Surface,
		Hover:      t.SurfaceHover,
		Pressed:    t.SurfacePressed,
		Border:     t.Border,
		Text:       t.Text,
		Muted:      t.TextMuted,
		Accent:     t.Accent,
		AccentText: t.AccentText,
		Danger:     t.Danger,
		Focus:      t.Focus,
	}
}

// SetTheme applies a design-system theme to the current MyGo frame.
func SetTheme(c *ui.Context, dark bool) {
	if dark {
		t := ui.DarkTheme()
		c.SetTheme(t)
		return
	}
	c.SetTheme(ui.LightTheme())
}

// Space converts a Tailwind-style spacing unit to DIPs. The default MyGo
// spacing unit is four DIPs.
func Space(n float32) float32 { return n * 4 }
