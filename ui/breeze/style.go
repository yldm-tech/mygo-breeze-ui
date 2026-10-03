package breeze

import "github.com/egoist/mygo/ui"

// Style is a typed utility that styles one MyGo element.
type Style func(*ui.Element) *ui.Element

// Apply applies styles from left to right and returns e.
func Apply(e *ui.Element, styles ...Style) *ui.Element {
	for _, style := range styles {
		if style != nil {
			e = style(e)
		}
	}
	return e
}

// Merge combines styles into one reusable style.
func Merge(styles ...Style) Style {
	return func(e *ui.Element) *ui.Element { return Apply(e, styles...) }
}

// Row creates a horizontal flex container and applies styles.
func Row(c *ui.Context, styles ...Style) *ui.Element { return Apply(ui.Row(c), styles...) }

// Column creates a vertical flex container and applies styles.
func Column(c *ui.Context, styles ...Style) *ui.Element { return Apply(ui.Column(c), styles...) }

// Box creates a column container and applies styles.
func Box(c *ui.Context, styles ...Style) *ui.Element { return Apply(ui.Box(c), styles...) }

// Scroll creates a vertical scroll container and applies styles.
func Scroll(c *ui.Context, styles ...Style) *ui.Element { return Apply(ui.Scroll(c), styles...) }

// Wrap enables flex wrapping on a row.
func Wrap(e *ui.Element) *ui.Element { return e.Wrap() }

// Gap sets the gap in default spacing units.
func Gap(n float32) Style { return func(e *ui.Element) *ui.Element { return e.Gap(Space(n)) } }

// P sets padding on every side in default spacing units.
func P(n float32) Style { return func(e *ui.Element) *ui.Element { return e.Padding(Space(n)) } }

// Px sets horizontal padding in default spacing units.
func Px(n float32) Style { return func(e *ui.Element) *ui.Element { return e.PaddingX(Space(n)) } }

// Py sets vertical padding in default spacing units.
func Py(n float32) Style { return func(e *ui.Element) *ui.Element { return e.PaddingY(Space(n)) } }

// M sets margin on every side in default spacing units.
func M(n float32) Style { return func(e *ui.Element) *ui.Element { return e.Margin(Space(n)) } }

// Size sets a fixed size in DIPs.
func Size(width, height float32) Style {
	return func(e *ui.Element) *ui.Element { return e.Size(width, height) }
}

// W sets a fixed width in DIPs.
func W(width float32) Style { return func(e *ui.Element) *ui.Element { return e.Width(width) } }

// H sets a fixed height in DIPs.
func H(height float32) Style { return func(e *ui.Element) *ui.Element { return e.Height(height) } }

// Fill makes an element consume its parent's available size.
func Fill(e *ui.Element) *ui.Element { return e.Fill() }

// Grow gives an element free space in its parent's main axis.
func Grow(f float32) Style { return func(e *ui.Element) *ui.Element { return e.Grow(f) } }

// Shrink sets the flex shrink factor.
func Shrink(f float32) Style { return func(e *ui.Element) *ui.Element { return e.Shrink(f) } }

// ItemsCenter centers children across the cross axis.
func ItemsCenter(e *ui.Element) *ui.Element { return e.AlignItems(ui.Center) }

// ItemsStart aligns children at the start of the cross axis.
func ItemsStart(e *ui.Element) *ui.Element { return e.AlignItems(ui.Start) }

// ItemsEnd aligns children at the end of the cross axis.
func ItemsEnd(e *ui.Element) *ui.Element { return e.AlignItems(ui.End) }

// JustifyCenter centers children along the main axis.
func JustifyCenter(e *ui.Element) *ui.Element { return e.Justify(ui.Center) }

// JustifyBetween distributes children with space between them.
func JustifyBetween(e *ui.Element) *ui.Element { return e.Justify(ui.SpaceBetween) }

// JustifyEnd places children at the end of the main axis.
func JustifyEnd(e *ui.Element) *ui.Element { return e.Justify(ui.End) }

// Rounded applies a corner radius token.
func Rounded(radius RadiusToken) Style {
	return func(e *ui.Element) *ui.Element { return e.Radius(float32(radius)) }
}

// Bg sets a background color.
func Bg(color ui.Color) Style { return func(e *ui.Element) *ui.Element { return e.Background(color) } }

// Border sets a one-DIP border.
func Border(color ui.Color) Style {
	return func(e *ui.Element) *ui.Element { return e.Border(1, color) }
}

// BorderWidth sets a border with an explicit width.
func BorderWidth(width float32, color ui.Color) Style {
	return func(e *ui.Element) *ui.Element { return e.Border(width, color) }
}

// ShadowSM applies a small surface shadow.
func ShadowSM(color ui.Color) Style {
	return func(e *ui.Element) *ui.Element { return e.Shadow(0, 1, 3, 0, color) }
}

// Alpha changes the opacity of an element and its descendants.
func Alpha(value float32) Style { return func(e *ui.Element) *ui.Element { return e.Opacity(value) } }

// TextSize sets the text size token.
func TextSize(size TextToken) Style {
	return func(e *ui.Element) *ui.Element { return e.FontSize(float32(size)) }
}

// FontWeight sets the text weight.
func FontWeight(weight int) Style {
	return func(e *ui.Element) *ui.Element { return e.FontWeight(weight) }
}

// TextColor sets the inherited text color.
func TextColor(color ui.Color) Style {
	return func(e *ui.Element) *ui.Element { return e.TextColor(color) }
}

// Muted applies the active theme's muted text color.
func Muted(c *ui.Context) Style { return TextColor(Palette(c).Muted) }

// Text creates a text element with styles.
func Text(c *ui.Context, value string, styles ...Style) *ui.Element {
	return Apply(ui.Text(c, value), styles...)
}
