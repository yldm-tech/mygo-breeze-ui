package breeze

import "github.com/egoist/mygo/ui"

// ButtonVariant controls a button's semantic color.
type ButtonVariant uint8

const (
	ButtonPrimary ButtonVariant = iota
	ButtonSecondary
	ButtonDanger
	ButtonGhost
)

// ButtonSize controls a button's padding and text size.
type ButtonSize uint8

const (
	ButtonMD ButtonSize = iota
	ButtonSM
	ButtonLG
)

// Button creates a themed, accessible push button. Its Clicked method can be
// queried in the same frame as the underlying MyGo button.
func Button(c *ui.Context, label string, variant ButtonVariant) *ui.Element {
	return ButtonWith(c, label, ButtonOptions{Variant: variant})
}

// ButtonOptions configures ButtonWith.
type ButtonOptions struct {
	Variant  ButtonVariant
	Size     ButtonSize
	Disabled bool
}

// ButtonWith creates a button from typed options.
func ButtonWith(c *ui.Context, label string, opts ButtonOptions) *ui.Element {
	p := Palette(c)
	base, hover, pressed, foreground, border := p.Surface, p.Hover, p.Pressed, p.Text, p.Border
	switch opts.Variant {
	case ButtonPrimary:
		base, hover, pressed, foreground, border = p.Accent, p.Accent.Mix(p.AccentText, 0.16), p.Accent.Mix(p.Text, 0.18), p.AccentText, ui.Transparent
	case ButtonDanger:
		base, hover, pressed, foreground, border = p.Danger, p.Danger.Mix(p.AccentText, 0.16), p.Danger.Mix(p.Text, 0.18), p.AccentText, ui.Transparent
	case ButtonGhost:
		base, hover, pressed, foreground, border = ui.Transparent, p.Hover, p.Pressed, p.Text, ui.Transparent
	}
	b := ui.ButtonBase(c).Disabled(opts.Disabled).Radius(float32(RadiusMD)).Background(base).TextColor(foreground)
	switch opts.Size {
	case ButtonSM:
		b.Padding(Space(1), Space(2)).FontSize(float32(TextSM))
	case ButtonLG:
		b.Padding(Space(2.5), Space(4)).FontSize(float32(TextLG))
	default:
		b.Padding(Space(1.5), Space(3)).FontSize(float32(TextBase))
	}
	if border.A != 0 {
		b.Border(1, border)
	}
	if opts.Disabled {
		b.Opacity(0.55)
	} else if b.Pressed() {
		b.Background(pressed)
	} else if b.Hovered() {
		b.Background(hover)
	}
	b.Children(func() { ui.Text(c, label).SingleLine() })
	return b
}

// Card creates a bordered surface container.
func Card(c *ui.Context, body func()) *ui.Element {
	p := Palette(c)
	return Column(c, P(4), Gap(3), Rounded(RadiusLG), Bg(p.Background), Border(p.Border), ShadowSM(ui.RGBA(0, 0, 0, 0.08))).Children(body)
}

// Input creates a themed single-line text input.
func Input(c *ui.Context, value *string, placeholder string) *ui.Element {
	e := ui.TextInput(c, value)
	if placeholder != "" {
		e.Placeholder(placeholder)
	}
	return e.Radius(float32(RadiusMD))
}

// TextArea creates a themed multi-line text input.
func TextArea(c *ui.Context, value *string, placeholder string) *ui.Element {
	e := ui.TextArea(c, value)
	if placeholder != "" {
		e.Placeholder(placeholder)
	}
	return e.Radius(float32(RadiusMD))
}

// Checkbox delegates behavior and accessibility to MyGo's checkbox widget.
func Checkbox(c *ui.Context, checked *bool, label string) *ui.Element {
	return ui.Checkbox(c, checked, label)
}

// Switch delegates behavior and accessibility to MyGo's switch widget.
func Switch(c *ui.Context, on *bool) *ui.Element { return ui.Switch(c, on) }

// Select delegates behavior and accessibility to MyGo's select widget.
func Select(c *ui.Context, selected *string, options []string) *ui.Element {
	return ui.Select(c, selected, options)
}

// BadgeVariant controls a badge's semantic color.
type BadgeVariant uint8

const (
	BadgeNeutral BadgeVariant = iota
	BadgeSuccess
	BadgeWarning
	BadgeDanger
)

// Badge creates a compact status label.
func Badge(c *ui.Context, label string, variant BadgeVariant) *ui.Element {
	p := Palette(c)
	background, foreground := p.Surface, p.Text
	switch variant {
	case BadgeSuccess:
		background, foreground = ui.Hex("#dcfce7"), ui.Hex("#166534")
	case BadgeWarning:
		background, foreground = ui.Hex("#fef3c7"), ui.Hex("#92400e")
	case BadgeDanger:
		background, foreground = ui.Hex("#fee2e2"), ui.Hex("#991b1b")
	}
	return Row(c, P(1), Px(2), Gap(1), Rounded(RadiusFull), Bg(background), TextColor(foreground), Shrink(0)).Children(func() {
		Text(c, label, TextSize(TextXS), FontWeight(600))
	})
}

// Progress creates a themed progress bar.
func Progress(c *ui.Context, value float64) *ui.Element { return ui.Progress(c, value) }

// Tabs delegates behavior and keyboard interaction to MyGo's tabs widget.
func Tabs(c *ui.Context, selected *int, labels ...string) *ui.Element {
	return ui.Tabs(c, selected, labels...)
}

// Dialog creates a modal dialog with a title, body, and optional footer.
func Dialog(c *ui.Context, open *bool, title string, body, footer func()) *ui.Element {
	return ui.Modal(c, open, func() {
		Text(c, title, TextSize(TextLG), FontWeight(600))
		if body != nil {
			body()
		}
		if footer != nil {
			Row(c, Gap(2), JustifyEnd).Children(footer)
		}
	})
}
