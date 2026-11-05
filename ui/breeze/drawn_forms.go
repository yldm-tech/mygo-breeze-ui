package breeze

import (
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/egoist/mygo/ui"
)

// ---- forms: check box groups, number inputs, splits, and the weights of field labels and legends ----
//
// macOS: NSButton check boxes (mixed state), NSTextField beside an NSStepper and NSDatePicker of style textField, rendered at 2x in Aqua and Dark Aqua on macOS 27.2 into an sRGB bitmap in a window answering hasKeyAppearance, as tools/appkit-reference/reference.swift renders. Windows 11: microsoft/microsoft-ui-xaml main at cb1b8a6 (CheckBox_themeresources.xaml, AnimatedAcceptVisualSource.cpp, NumberBox.xaml, NumberBox_themeresources.xaml, TextBox_themeresources.xaml) and WinUI-Gallery main at 7360ac3 (Samples/CheckBox/CheckBoxPage.xaml). GNOME: libadwaita 1.10.0's _checks.scss, assets/dash.svg, _spin-button.scss, _entries.scss, _buttons.scss, _colors.scss, _paned.scss and _labels.scss, and adwaita-icon-theme 51.0's value-increase-symbolic and value-decrease-symbolic.

// formsCheckState is what a check box shows: off, on, or mixed, as a CheckboxGroup's box over boxes that differ.
type formsCheckState uint8

const (
	formsCheckOff formsCheckState = iota
	formsCheckOn
	formsCheckMixed
)

// formsCheckGroups holds the windows building a CheckboxGroup, with the values of the check boxes Checkbox builds inside it, which the group counts and sets as ui.CheckboxGroup does package ui's (its registry, the context's checks, is MyGo's own).
var formsCheckGroups sync.Map // *ui.Context -> *formsCheckGroup

type formsCheckGroup struct{ boxes []*bool }

// formsJoinCheckGroup adds the value of a check box being built to the CheckboxGroup being built, if any.
func formsJoinCheckGroup(c *ui.Context, checked *bool) {
	if v, ok := formsCheckGroups.Load(c); ok {
		g := v.(*formsCheckGroup)
		g.boxes = append(g.boxes, checked)
	}
}

// formsCheckGroupSpacing returns how far a CheckboxGroup indents its boxes and how far apart it stacks its rows. macOS has no group of check boxes: 22 puts the boxes under the label of the box over them (macCheckbox's 16-point box and the 6 after it) and 8 is MyGo's gap, both inferred. Windows: WinUI Gallery's three-state CheckBox sample indents its options by Margin 24,0,0,0 in a StackPanel without spacing (CheckBoxPage.xaml at 7360ac3, lines 60-87), its rows already 32 high (CheckBoxHeight). GNOME has none either: 24 puts each box under the label above it (adwCheckbox's row padding 4, 20-pixel indicator and 4 of spacing, less the next row's own padding of 4), and no gap, as the rows are padded already (inferred).
func formsCheckGroupSpacing(p Platform) (indent, gap float32) {
	switch p {
	case Windows, Linux:
		return 24, 0
	}
	return 22, 8
}

// formsDrawCheckbox draws a check box of the platform in state, as Checkbox draws one on or off; mixed, the box is drawn as checked and its check mark replaced with the platform's dash.
func formsDrawCheckbox(c *ui.Context, p Platform, row, box ui.Element, state formsCheckState) {
	on := state != formsCheckOff
	switch p {
	case Windows:
		fluentCheckbox(c, row, box, on)
	case Linux:
		adwCheckbox(c, row, box, on)
	default:
		macCheckbox(c, row, box, on)
	}
	if state != formsCheckMixed {
		return
	}
	switch p {
	case Windows:
		formsFluentDash(c, row, box)
	case Linux:
		formsAdwDash(c, box)
	default:
		formsMacDash(c, row, box)
	}
}

// formsMacDash draws the mark of an NSButton check box in the mixed state over the box macCheckbox fills as checked: the accent (pressed, AccentPressed; disabled, macDisabledOn's gray) and a white capsule 2 points high and 6.6 wide centered in the 16-point box, with round ends. Measured in a key window at 2x on macOS 27.2: the dash covers rows 14 to 17 of the 32 pixels, the two middle ones over 12 full pixels and two at 62% (13.2 pixels), the outer two over the middle 8 only.
func formsMacDash(c *ui.Context, row, box ui.Element) {
	t := c.Theme()
	fill := t.Accent
	if row.Pressed() {
		fill = t.AccentPressed
	}
	ring, disabled := row.FocusVisible(), disabledOnce(row)
	box.DrawOver(func(p *ui.Painter, r ui.Rect) {
		face, mark := fill, t.AccentText
		if disabled() {
			face, mark = macDisabledOn(t)
		}
		p.Fill(r, face, macCheckboxRadius)
		w, h := 6.6*r.W/16, 2*r.H/16
		p.Fill(ui.Rect{X: r.X + (r.W-w)/2, Y: r.Y + (r.H-h)/2, W: w, H: h}, mark, h/2)
		if ring {
			macFocusRing(p, r, macCheckboxRadius)
		}
	})
}

// formsFluentDash draws a CheckBox's indeterminate glyph over the box fluentCheckbox fills as checked (CheckBoxCheckBackgroundFillIndeterminate and its stroke are AccentFillColorDefault, Secondary under the pointer and Tertiary pressed, as checked): AnimatedAcceptVisualSource.cpp's Geometry_0, a line from (-11.75,-0.125) to (11.875,-0.125) stroked 4 with round caps in a shape scaled 0.7 at 24,24, in a container scaled 1.05 about 24,24, its 48-unit visual fit to the 20-DIP box without the checked glyph's Margin (IndeterminateNormal sets CheckGlyph.Margin 0): x = 10 + 0.30625·px, y = 10 + 0.30625·py, so from (6.402, 9.962) to (13.637, 9.962), 1.225 DIPs wide. The glyph is TextOnAccentFillColorPrimary, Secondary pressed (CheckBoxCheckGlyphForegroundIndeterminatePressed), as fluentCheckbox's check mark.
func formsFluentDash(c *ui.Context, row, box ui.Element) {
	t := c.Theme()
	mark := t.AccentText
	if row.Pressed() {
		mark = mark.Alpha(pick(t.Dark, float32(0.7), 0.5))
	}
	box.DrawOver(func(p *ui.Painter, r ui.Rect) {
		s := r.W / 20
		var path ui.Path
		path.MoveTo(r.X+6.402*s, r.Y+9.962*s).LineTo(r.X+13.637*s, r.Y+9.962*s)
		p.StrokePath(&path, 1.225*s, mark)
	})
}

// formsAdwDash draws a GtkCheckButton's :indeterminate indicator over the box adwCheckbox fills as checked (_checks.scss styles :checked and :indeterminate alike, accent_bg_color under accent_fg_color): assets/dash.svg, a capsule from 1 to 13 across and 5 to 9 down with ends of radius 2 in the 14-pixel icon, which the indicator's padding of 3 puts at 4 to 16 and 8 to 12 of the 20-pixel box.
func formsAdwDash(c *ui.Context, box ui.Element) {
	t := c.Theme()
	box.DrawOver(func(p *ui.Painter, r ui.Rect) {
		s := r.W / 20
		p.Fill(ui.Rect{X: r.X + 4*s, Y: r.Y + 8*s, W: 12 * s, H: 4 * s}, t.AccentText, 2*s)
	})
}

// formsSplit builds a split of package ui with build, keeping its divider in the theme's Border color under the pointer and while held, where MyGo turns it to the accent: the split builds in a sub-theme whose Accent is its Border, which the divider's drawing takes as it is built, and each pane builds in the window's theme again, so what the panes hold keeps the accent. libadwaita 1.10.0's _paned.scss has no :hover rule for a GtkPaned's separator, NSSplitView's thin divider does not change under the pointer, and WinUI's only resize handle (ResizeGripper, private to TableView) shows no line under the pointer. The split's drag, keys, focus ring, Changed and sizes stay MyGo's; the column holding it takes its place in the parent's layout and the styles given to it, Fill and Height among them.
func formsSplit(c *ui.Context, build func(first, second func()) ui.Element, first, second func()) ui.Element {
	orig := c.Theme()
	th := *orig
	th.Accent = th.Border
	wrap := ui.Themed(c, &th).Background(ui.Transparent).AlignItems(ui.Stretch)
	pane := func(fn func()) func() {
		return func() { ui.Themed(c, orig).Background(ui.Transparent).Grow(1).AlignSelf(ui.Stretch).Children(fn) }
	}
	var s ui.Element
	wrap.Children(func() { s = build(pane(first), pane(second)).Grow(1) })
	if s.Changed() {
		wrap.SetChanged()
	}
	return wrap
}

// formsNumberDecimals returns the decimals a NumberInput shows, those of step, as ui.NumberInput counts them.
func formsNumberDecimals(step float64) int {
	s := strconv.FormatFloat(step, 'f', -1, 64)
	if i := strings.IndexByte(s, '.'); i >= 0 {
		return len(s) - i - 1
	}
	return 0
}

// formsNumberInput builds NumberInput as ui.NumberInput does (number.go of MyGo v0.7.4), from package ui's exported bases: a row holding a text input and two buttons labeled "Decrease" and "Increase", which MyGo names in English in its code too; what is typed applies once it is a number in range, Up and Down step, the buttons step and are disabled at the limits, and the text shows the value rounded to step's decimals once the input loses the focus. The row's Changed reports a new value. field gives the input its look and returns it; buttons builds the two buttons in the row, Decrease first, as MyGo orders them for Tab, with step, which builds a button stepping by sign and returns it to be drawn.
func formsNumberInput(c *ui.Context, value *float64, lo, hi, step float64, field func(text *string) ui.Element, buttons func(step func(sign float64) ui.Element)) ui.Element {
	row := ui.Row(c).Shrink(0).AlignItems(ui.Center)
	decimals := formsNumberDecimals(step)
	format := func(v float64) string { return strconv.FormatFloat(v, 'f', decimals, 64) }
	text := ui.Local(row, "breeze-number-text", func() string { return format(*value) })
	set := func(v float64) {
		v = math.Max(lo, math.Min(hi, v))
		p := math.Pow(10, float64(decimals))
		v = math.Round(v*p) / p
		if v != *value {
			*value = v
			row.SetChanged()
		}
	}
	row.Children(func() {
		in := field(text)
		if in.Changed() {
			if v, err := strconv.ParseFloat(strings.TrimSpace(*text), 64); err == nil && v >= lo && v <= hi {
				set(v)
			}
		}
		if in.Shortcut(0, ui.KeyUp) {
			set(*value + step)
			*text = format(*value)
		}
		if in.Shortcut(0, ui.KeyDown) {
			set(*value - step)
			*text = format(*value)
		}
		if !in.Focused() {
			// What the app set, or what was typed, shown in full.
			*text = format(*value)
		}
		in.Range(lo, hi, *value).Step(step)
		buttons(func(sign float64) ui.Element {
			b := ui.ButtonBase(c).Label(map[bool]string{true: "Increase", false: "Decrease"}[sign > 0]).FocusRing(false)
			b.Disabled(sign < 0 && *value <= lo || sign > 0 && *value >= hi)
			if b.Clicked() {
				set(*value + sign*step)
				*text = format(*value)
			}
			return b
		})
	})
	return row
}

// formsMacNumberInput draws NumberInput as an NSTextField with an NSStepper beside it: the text field (macField), as wide as MyGo's (80), and 4 after it (MyGo's gap, inferred) a stepper of regular size, measured in a key window at 2x on macOS 27.2: 20 by 24 points, a push button's face (235 over white, 48 over the window in Dark Aqua) with corners of 6, a line 13 points wide across its middle (macStepper's), and in each half a chevron 9.2 points wide and 4.7 high along its stroke's middle, stroked 1.75 with round ends (its arms measure 4.5 pixels across a row at 45 degrees), centered in the half, in labelColor over the face (36 over 235, 223 over 48), pointing up in the upper half, which increases, and down in the lower one. A half darkens to a pressed push button's face while pressed (inferred), its chevron dims at its limit as MyGo's buttons are disabled there (AppKit's does not), and AppKit's focus ring goes around the whole stepper while either half has the focus from the keyboard, as around an NSStepper.
func formsMacNumberInput(c *ui.Context, value *float64, lo, hi, step float64) ui.Element {
	t := c.Theme()
	face, down := macFace(t.Background, t.Dark), macPressed(t.Background, t.Dark)
	ink := tint(face, t.Dark, 0.847)
	line := tint(t.Background, t.Dark, pick(t.Dark, float32(0.204), 0.222))
	row := formsNumberInput(c, value, lo, hi, step, func(text *string) ui.Element {
		return TextInput(c, text).Width(80)
	}, func(stepBy func(sign float64) ui.Element) {
		var ringed [2]bool
		// Built Decrease first, for MyGo's Tab order; reversed, Increase shows on top.
		col := ui.Column(c).Reverse().Size(20, 24).Shrink(0).MarginStart(4).Radius(6).Background(face)
		col.Children(func() {
			for i, sign := range []float64{-1, 1} {
				b := stepBy(sign).Height(12).Width(20)
				if b.Pressed() {
					if sign > 0 {
						b.Background(down).Radius(6, 6, 0, 0)
					} else {
						b.Background(down).Radius(0, 0, 6, 6)
					}
				}
				ringed[i] = b.FocusVisible()
				b.Draw(func(p *ui.Painter, r ui.Rect) {
					formsStrokeChevron(p, r.X+r.W/2, r.Y+r.H/2, 4.6, 2.35*float32(-sign), 1.75, ink)
				})
			}
		})
		col.DrawOver(func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X + (r.W-13)/2, Y: r.Y + r.H/2 - 0.5, W: 13, H: 1}, line, 0)
			if ringed[0] || ringed[1] {
				macFocusRing(p, r, 6)
			}
		})
	})
	return row
}

// formsStrokeChevron strokes a chevron centered on cx, cy, its ends half to either side and its tip tip below the center (above it for a negative tip), the ends as far the other way.
func formsStrokeChevron(p *ui.Painter, cx, cy, half, tip, width float32, col ui.Color) {
	var path ui.Path
	path.MoveTo(cx-half, cy-tip).LineTo(cx, cy+tip).LineTo(cx+half, cy-tip)
	p.StrokePath(&path, width, col)
}

// formsFluentNumberInput draws NumberBox with SpinButtonPlacementMode Inline: the row is the TextBox, drawn as fluentField draws one (formsFluentTextBox) while the input inside has the focus, its text padded 10,5,6,6 (TextControlThemePadding); at its end the spin buttons (NumberBox.xaml): UpSpinButton then DownSpinButton, each a RepeatButton of NumberBoxSpinButtonStyle at least 32 DIPs wide and stretched to the box less its margin, Margin 4 for Up and 0,4,4,4 for Down, corners of 4 (ControlCornerRadius), with the RepeatButton brushes NumberBox sets from TextBox's inner button: no fill (TextControlButtonBackground, from generic.xaml, not read: inferred), SubtleFillColorSecondary under the pointer and SubtleFillColorTertiary pressed, a transparent border, and glyphs E70E and E70D of FontSize 12 in TextFillColorSecondary, TextFillColorTertiary pressed, drawn as fluentChevronDown's 7 by 3.5 chevron (the glyphs' outline in Segoe Fluent Icons is inferred). The text is at least as wide as MyGo's field, 80, so the box is at least 152 DIPs wide, past NumberBoxMinWidth's 120; it fills the width a column gives it, as a NumberBox stretches by default (HorizontalAlignment Stretch). WinUI's spin buttons are not tab stops; MyGo's are, so the system focus visual goes 3 DIPs around one while it has the focus from the keyboard (FocusVisualMargin -3, as a button's).
func formsFluentNumberInput(c *ui.Context, value *float64, lo, hi, step float64) ui.Element {
	t := c.Theme()
	var in ui.Element
	hover, press := fluentSubtle(t)
	tertiary := t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.447), 0.529))
	row := formsNumberInput(c, value, lo, hi, step, func(text *string) ui.Element {
		in = ui.TextInputBase(c, text).Width(80).Grow(1).Padding(5, 6, 6, 10)
		return in
	}, func(stepBy func(sign float64) ui.Element) {
		// Built Decrease first, for MyGo's Tab order; reversed, Up shows first, as NumberBox's columns put it.
		ui.Row(c).Reverse().AlignSelf(ui.Stretch).Shrink(0).Children(func() {
			for _, sign := range []float64{-1, 1} {
				b := stepBy(sign).Width(32).AlignSelf(ui.Stretch).MarginY(4).MarginEnd(4).Radius(4)
				if sign > 0 {
					b.MarginStart(4)
				}
				fg := t.TextMuted
				switch {
				case b.Pressed():
					b.Background(press)
					fg = tertiary
				case b.Hovered():
					b.Background(hover)
				}
				focused := b.FocusVisible()
				b.Draw(func(p *ui.Painter, r ui.Rect) {
					cx, cy := r.X+r.W/2, r.Y+r.H/2
					if sign > 0 {
						p.Transform(ui.Rotation(180), cx, cy, func() { fluentChevronDown(p, cx, cy, fg) })
					} else {
						fluentChevronDown(p, cx, cy, fg)
					}
				})
				b.DrawOver(func(p *ui.Painter, r ui.Rect) {
					if focused {
						fluentFocusRing(p, t, r, 3, 3, 3, 3, 4)
					}
				})
			}
		})
	})
	formsFluentTextBox(c, row, in.Focused())
	return row
}

// formsFluentTextBox draws the box of a TextBox around a row holding its input and its buttons, as fluentField draws a text field's, with the focus look while focused: the input inside has the focus, which fluentField reads from the field itself.
func formsFluentTextBox(c *ui.Context, row ui.Element, focused bool) {
	t := c.Theme()
	fill := t.Surface
	switch {
	case focused:
		fill = pick(t.Dark, white, t.Background.Mix(ui.Hex("#1e1e1e"), 0.7))
	case row.Hovered():
		fill = t.SurfaceHover
	}
	row.MinHeight(32).Radius(4).Background(fill).Border(0, ui.Color{})
	strong := t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.447), 0.545))
	row.DrawOver(func(p *ui.Painter, r ui.Rect) {
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

// formsAdwNumberInput draws NumberInput as a horizontal GtkSpinButton (_spin-button.scss): an entry (%entry: 34 pixels high, corners of 9, currentColor at 10%, its 2-pixel focus ring inside its edge while the text has the focus) without padding, its text padded 6 (at least 28 wide; at least MyGo's 80 here, filling the width a column gives the entry, as a GtkSpinButton fills a vertical GtkBox), then the − and + buttons, in that order (GTK 4's spin button, inferred from GTK, not libadwaita), each a flat image button (%button_basic_flat, .image-button: min-width 22 under _spin-button.scss and 5 of padding at either side, 32 with the 16-pixel icon centered) as high as the entry, with a 1-pixel border of currentColor at 10% on its start side, $hover_color (currentColor at 7%) under the pointer and $active_color (16%) pressed, the last one's end corners of 9 ($button_radius), which the entry's clip gives. Their icons are adwaita-icon-theme 51.0's value-decrease-symbolic (a bar from 3 to 13 across, 7 to 9 down) and value-increase-symbolic (that bar crossed by one from 3 to 13 down, 7 to 9 across) in the text color. GTK's spin buttons are not focusable; MyGo's are, so a button shows libadwaita's focus ring inside its edge while it has the focus from the keyboard.
func formsAdwNumberInput(c *ui.Context, value *float64, lo, hi, step float64) ui.Element {
	t := c.Theme()
	fill := adwInkAlpha(t, t.Background, 0.10)
	edge, ink := adwInkAlpha(t, fill, 0.10), adwInkAlpha(t, fill, 1)
	var in ui.Element
	row := formsNumberInput(c, value, lo, hi, step, func(text *string) ui.Element {
		in = ui.TextInputBase(c, text).Width(80).Grow(1).Padding(6)
		return in
	}, func(stepBy func(sign float64) ui.Element) {
		for _, sign := range []float64{-1, 1} {
			b := stepBy(sign).Width(33).AlignSelf(ui.Stretch).Radius(0).BorderStartWidth(1).BorderColor(edge)
			switch {
			case b.Pressed():
				b.Background(adwInkAlpha(t, fill, 0.16))
			case b.Hovered():
				b.Background(adwInkAlpha(t, fill, 0.07))
			}
			b.Draw(func(p *ui.Painter, r ui.Rect) {
				// The 16-pixel icon, centered in the button past its 1-pixel border.
				x, y := r.X+(r.W-16)/2, r.Y+(r.H-16)/2
				p.Fill(ui.Rect{X: x + 3, Y: y + 7, W: 10, H: 2}, ink, 0)
				if sign > 0 {
					p.Fill(ui.Rect{X: x + 7, Y: y + 3, W: 2, H: 10}, ink, 0)
				}
			})
			adwFocusRing(t, b, 0, nil)
		}
	})
	row.MinHeight(34).Radius(9).Background(fill).Border(0, ui.Color{}).AlignItems(ui.Center).Clip()
	focused := in.Focused()
	row.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if focused {
			p.Stroke(r, t.Focus, 9, 2)
		}
	})
	return row
}

// formsLabelOver draws label at weight in color at the top of r, where MyGo lays out a field's label or a fieldset's legend above what it holds, from r's start: from its right edge right to left. MyGo's own text there is left transparent, so it keeps the layout, the clicks and the names it gives.
func formsLabelOver(p *ui.Painter, r ui.Rect, label string, weight int, size float32, color ui.Color, rtl bool) {
	sp := ui.Span{Text: label, Size: size, Weight: weight, Color: color}
	x := r.X
	if rtl {
		w, _ := p.MeasureText(r.W, sp)
		x = r.X + r.W - w
	}
	p.RichText(x, r.Y, r.W, sp)
}

// formsCoverIcon paints fill over MyGo's 14-DIP calendar icon at the end of a date field r, padEnd in from its end (its start right to left), a DIP past it on every side, and returns where the icon is. Disabled, the field paints at half opacity, which one cover would let the icon show through; the cover is then the field as it shows at half opacity over the window, painted four times, which leaves a sixteenth of the icon's half.
func formsCoverIcon(p *ui.Painter, r ui.Rect, padEnd float32, fill, window ui.Color, rtl, disabled bool) ui.Rect {
	const size = 14
	icon := ui.Rect{X: r.X + r.W - padEnd - size, Y: r.Y + (r.H-size)/2, W: size, H: size}
	if rtl {
		icon.X = r.X + padEnd
	}
	cover := ui.Rect{X: icon.X - 1, Y: icon.Y - 1, W: icon.W + 2, H: icon.H + 2}
	times := 1
	if disabled {
		fill, times = window.Mix(fill, 0.5), 4
	}
	for range times {
		p.Fill(cover, fill, 0)
	}
	return icon
}

// formsDateIconInk returns the opacity of the icon of a GNOME entry (_entries.scss: entry > image at 70%, 100% under the pointer, 80% pressed) for a date field b, its icon padEnd in from its end (its start right to left): the share of currentColor it shows in. The pointer is read against the field as the last frame laid it out.
func formsDateIconInk(b ui.Element, padEnd float32, rtl bool) float32 {
	const size = 14
	x, y, over := b.PointerPosition()
	r := b.Bounds()
	left := r.W - padEnd - size
	if rtl {
		left = padEnd
	}
	over = over && x >= left && x < left+size && y >= (r.H-size)/2 && y < (r.H+size)/2
	switch {
	case over && b.Pressed():
		return 0.8
	case over:
		return 1
	}
	return 0.7
}
