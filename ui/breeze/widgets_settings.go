package breeze

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// ---- settings rows: combo rows, entry rows, spin rows and expander rows ----
//
// The rows that complete ActionRow and SwitchRow (widgets_rows.go), built on its row look (rowsJoin, rowsStyle) so that they join a Card, a Fieldset or an AccordionItem as the platform joins rows there. macOS: SwiftUI's grouped Form with a Picker of the menu style, TextFields with a label, a Stepper with a format (its value in a text field) and DisclosureGroups whose label is a title over a subtitle, rendered at 2x in Aqua and Dark Aqua on macOS 27.2 by tools/appkit-reference/reference.swift (its settings scene). Windows 11: CommunityToolkit/Windows's SettingsCard.xaml (the MinWidth of 120 it gives a ComboBox and a TextBox in a card, SettingsCardContentMinWidth) and SettingsExpander.xaml (main at a784bad: SettingsExpanderHeaderPadding 16,16,4,16, the 32-DIP chevron button with a Margin of 0,0,8,0 and a 16-DIP AnimatedChevronUpDownSmallVisualSource, the header's ToggleButton states, the items' 58,8,44,8), microsoft-ui-xaml's Expander_themeresources.xaml (ExpanderContentDownBorderThickness 1,0,1,1) and the SettingsExpander of tools/winui-reference's winui-rows, rendered by WinUI 3 3.3 on Windows 11 26H2, with breeze's Select, TextInput and NumberInput for the ComboBox, the TextBox and the NumberBox. GNOME: libadwaita 1.10.0's adw-combo-row.ui and .c, adw-entry-row.ui and .c, adw-spin-row.ui and .c, adw-expander-row.ui, _lists.scss (row.entry, row.spin, row.combo, row.expander, %boxed_list_row), _spin-button.scss, _buttons.scss (%circular_button) and its icons adw-entry-edit.svg and adw-expander-arrow.svg, and the AdwExpanderRow of tools/adwaita-reference's rows set for the open arrow's color.

// ComboRow creates a settings row choosing *selected among options: the title, the subtitle below it (an empty subtitle shows none) and the choice at its end, opening a list of the options as the platform does. It joins a Card, a Fieldset or an AccordionItem as ActionRow does. Changed reports a new choice.
//   - macOS: a SwiftUI Picker of the menu style in a grouped Form, as System Settings shows one: at the row's trailing edge a borderless pop-up button, 24 points high, showing the choice in labelColor 12 points before a 20-point circle of black at 7.8% (white in Dark Aqua: 228 over the section's 247, 54 over 37) holding up and down chevrons in labelColor over it, the circle 4 points in from the button's end, a row of one line 37.5 points high, a point taller than a row of text for the button (measured on macOS 27.2). The button is the row's focus stop, with AppKit's focus ring around it; its menu is Select's.
//   - Windows: a SettingsCard holding a ComboBox (breeze's Select) at its end, at least 120 DIPs wide as SettingsCard makes a ComboBox (SettingsCardContentMinWidth); the ComboBox is the focus stop, as WinUI focuses the control in a card that is not clickable.
//   - GNOME: AdwComboRow, the whole row a combo box that a click, Space or Enter opens (adw_combo_row_activate): the choice in card_fg_color, padded 2 as an inline list view's row, and pan-down-symbolic 6 pixels after it, the row lit under the pointer and while its popover is open (has-open-popup), its focus ring across its edge; the popover is a menu of the options with a check mark by the chosen one, at least 120 pixels wide (row.combo popover > contents), aligned with the row's end, where GTK centers it on the arrow (inferred: MyGo places a popup at an edge of what it is attached to).
//
// Assistive technology hears the title as the name of the choosing control and the subtitle as its description.
func ComboRow(c *ui.Context, title, subtitle string, selected *string, options []string) ui.Element {
	p := PlatformOf(c)
	if p == Linux {
		return adwComboRow(c, title, subtitle, selected, options)
	}
	var trigger ui.Element
	build := func() {
		if p == Windows {
			trigger = Select(c, selected, options).MinWidth(120)
		} else {
			trigger = macFormPicker(c, selected, options)
		}
		trigger.Label(title)
		if subtitle != "" {
			trigger.Description(subtitle)
		}
	}
	var row ui.Element
	if p == Windows {
		row = ActionRow(c, title, subtitle, build)
	} else {
		// SwiftUI makes a picker's row of one line 37.5 points high for its 24-point button (measured on macOS 27.2), where a row of text is 36.5.
		row = macFormRow(c, title, subtitle, 37.5, build)
	}
	if trigger.Changed() {
		row.SetChanged()
	}
	return row
}

// macFormPicker draws the pop-up button of a SwiftUI Picker of the menu style in a grouped Form on macOS 27.2: no bezel, the choice in labelColor over the section, then, 12 points after it, a 20-point circle of the section's separator share (black at 7.8% in Aqua, white in Dark Aqua) holding two chevrons stroked 1.5 points with round ends, 5 points wide and 2.5 deep along the stroke, their tips 4.4 points above and below the circle's middle, in labelColor over the circle (35 over 228, 224 over 54), the circle 4 points before the button's end, the button 24 points high (measured from the render of reference.swift's settings scene). Its menu is Select's, macMenu and macMenuItem; its focus ring is AppKit's around the button (macRing, corners of 6 inferred from the push button's).
func macFormPicker(c *ui.Context, selected *string, options []string) ui.Element {
	t := c.Theme()
	sel := ui.SelectBase(c, selected)
	b := sel.Trigger
	fill := macFormSectionFill(t)
	circle := tint(fill, t.Dark, macSectionSeparatorAlpha)
	ink, chevron := tint(fill, t.Dark, 0.847), tint(circle, t.Dark, 0.847)
	// The start's padding is inferred: the button draws nothing there.
	b.Height(24).Shrink(0).PaddingStart(12).PaddingEnd(4).Gap(12).AlignItems(ui.Center).Radius(6)
	macRing(b, 6)
	value := *selected
	b.Children(func() {
		ui.Text(c, value).SingleLine().TextColor(ink)
		ui.Box(c).Size(20, 20).Shrink(0).Radius(10).Background(circle).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			cx, cy := r.X+r.W/2, r.Y+r.H/2
			formsStrokeChevron(p, cx, cy-3.15, 2.5, -1.25, 1.5, chevron)
			formsStrokeChevron(p, cx, cy+3.15, 2.5, 1.25, 1.5, chevron)
		})
	})
	sel.Popup(func(panel ui.Element) {
		panel.ClipY()
		macMenu(c, panel)
		for _, opt := range options {
			macMenuItem(c, sel.Item(opt), opt, opt == value)
		}
	})
	return b
}

// adwComboRow draws ComboRow on GNOME: the row is ui.SelectBase's trigger, so the whole row opens the list and is one focus stop, as AdwComboRow is.
func adwComboRow(c *ui.Context, title, subtitle string, selected *string, options []string) ui.Element {
	sel := ui.SelectBase(c, selected)
	row := sel.Trigger.Shrink(0).Label(title).FocusRing(false)
	if subtitle != "" {
		row.Description(subtitle)
	}
	item, in, first := rowsJoin(c, row)
	look := rowsLook{lit: true, hovered: row.Hovered() || sel.Open(), pressed: row.Pressed(), focused: row.FocusVisible()}
	value := *selected
	rowsStyle(c, row, item, in, first, title, subtitle, look, func(under ui.Color, _, _ bool) {
		// The value is a suffix of the row, a label in an inline GtkListView whose row is padded 2 (list > row); the arrow box is the next suffix, 6 after it (border-spacing of .suffixes).
		fg := adwTextOn(c.Theme(), under)
		ui.Text(c, value).SingleLine().MarginX(2).TextColor(fg)
		idiomAdwArrow(c, fg)
	})
	sel.Popup(func(panel ui.Element) {
		panel.ClipY()
		adwMenu(c, panel)
		panel.MinWidth(120).AttachTo(row, ui.AnchorBottomEnd, ui.AnchorTopEnd)
		for _, opt := range options {
			adwMenuItem(c, sel.Item(opt), opt, opt == value)
		}
	})
	return row
}

// EntryRow creates a settings row editing *value, named title, which joins a Card, a Fieldset or an AccordionItem as ActionRow does; Changed reports an edit. The text field is the row's focus stop and assistive technology hears the title as its name.
//   - macOS: a SwiftUI TextField with a label in a grouped Form: the title in labelColor at the row's leading edge and, 8 points after it, the field taking the rest of the row without a bezel, its text at the trailing edge (measured on macOS 27.2; an empty field shows nothing). AppKit's focus ring goes around the field while it is edited (inferred).
//   - Windows: a SettingsCard holding a TextBox (breeze's TextInput) at its end, 120 DIPs wide, the MinWidth SettingsCard gives a TextBox (SettingsCardContentMinWidth), where an empty TextBox takes no more (WinUI's grows with its text, Breeze UI's scrolls it).
//   - GNOME: AdwEntryRow: the title inside the row, dimmed, centered while the row is empty and unfocused, and shrinking to the subtitle's size above the text, 3 pixels above it (TITLE_SPACING), as the row is focused or filled, over 150 milliseconds (EMPTY_ANIMATION_DURATION); the text 12 pixels in from the row's edge (a header margin of 6 and the editable area's padding of 6, row.entry), and at its end the 16-pixel adw-entry-edit pencil in a 34-pixel box 6 after the text (.edit-icon: min-width and min-height 24, padding 5, margin 6) while the row is not being edited. A click anywhere on the row focuses the text (pressed_cb); the row draws its focus ring while the text has the focus, from the pointer too (row.entry's focus ring on .focused), and is not lit under the pointer then.
//
// AdwPasswordEntryRow is left out: it adds a peek button that reveals the text, which needs PasswordField's per-platform reveal behavior inside the row rather than an option on this one; Password on the field is not reachable through the returned row.
func EntryRow(c *ui.Context, title string, value *string) ui.Element {
	switch PlatformOf(c) {
	case Windows:
		var in ui.Element
		row := ActionRow(c, title, "", func() { in = TextInput(c, value).Width(120).Label(title) })
		if in.Changed() {
			row.SetChanged()
		}
		return row
	case Linux:
		return adwEntryRow(c, title, value)
	}
	return macEntryRow(c, title, value)
}

// macEntryRow draws EntryRow on macOS as a row of a grouped Form holding a TextField with a label.
func macEntryRow(c *ui.Context, title string, value *string) ui.Element {
	row := ui.Row(c).Shrink(0).Role(ui.RoleGroup).Label(title)
	item, in, first := rowsJoin(c, row)
	var field ui.Element
	look := rowsLook{content: func(ink rowsInk) {
		ui.Text(c, title).Shrink(0).TextColor(ink.title)
		field = ui.TextInputBase(c, value).Grow(1).MinWidth(0).TextAlign(ui.End).TextColor(ink.title).Label(title)
		focused := field.Focused()
		field.DrawOver(func(p *ui.Painter, r ui.Rect) {
			if focused {
				macFocusRing(p, r, 4)
			}
		})
	}}
	rowsStyle(c, row, item, in, first, title, "", look, nil)
	// The title and the field share the title's line, 8 apart (measured: the field starts 8 points after the title's advance). SwiftUI lays a TextField's row half a point shorter than a row of text, 36 points between two separators, the title as far below the row's top (measured on macOS 27.2): half a point less below the text. In a DisclosureGroup the same half point is taken from rowsStyle's padding there (inferred).
	top, bottom := float32(10.5), float32(10.5)
	switch {
	case in == rowsInAccordion:
		top, bottom = 4.5, 4.5
	case in != rowsAlone && !first:
		top = 11.5
	}
	row.Gap(8).Padding(top, 10, bottom, 10)
	if field.Changed() {
		row.SetChanged()
	}
	return row
}

// adwEntryRow draws EntryRow on GNOME as AdwEntryRow.
func adwEntryRow(c *ui.Context, title string, value *string) ui.Element {
	t := c.Theme()
	row := ui.Row(c).Shrink(0)
	item, in, first := rowsJoin(c, row)
	within := row.FocusWithin()
	look := rowsLook{lit: true, hovered: row.Hovered() && !within, pressed: row.Pressed() && !within, focused: within}
	var input ui.Element
	rtl := c.Direction() == ui.RightToLeft
	look.content = func(ink rowsInk) {
		// The editable area: the header's 50 pixels, the title where it shrinks to and the text 3 below it, centered as allocate_editable_area centers them.
		area := ui.Column(c).Grow(1).Shrink(1).MinWidth(0).MinHeight(50).Justify(ui.Center).Gap(3)
		var small ui.Element
		area.Children(func() {
			small = ui.Text(c, title).SingleLine().FontSize(ink.subSize).TextColor(ui.Transparent).Role(ui.RoleNone)
			input = ui.TextInputBase(c, value).TextColor(ink.title).Label(title)
		})
		focused := input.Focused()
		target := float32(0)
		if *value != "" || focused {
			target = 1
		}
		filled := area.Animate("breeze-entry-row-filled", target, 150*time.Millisecond)
		dim, big, sub, family := ink.sub, t.FontSize, ink.subSize, t.Font
		area.DrawOver(func(p *ui.Painter, r ui.Rect) {
			size := big + (sub-big)*filled
			// The title is shaped at its own size, so that its line is that size's, as the small text's is: a span of RichText sits on the baseline of a line of the theme's size, which drew the filled row's title about 2 pixels below the small text's place (the ink of "Name" 3.5 pixels below libadwaita 1.10.0's in tools/adwaita-reference's settings set, the text below it 1.5).
			l := ui.ShapeText(title, ui.Font{Family: family, Size: size}, 0)
			w, _ := l.Size()
			_, hb := ui.ShapeText(title, ui.Font{Family: family, Size: big}, 0).Size()
			y := r.Y + (r.H-hb)/2
			y += (small.Bounds().Y - y) * filled
			x := r.X
			if rtl {
				x = r.X + r.W - w
			}
			p.TextLayout(l, x, y, dim)
		})
		// The pencil while the row is not being edited, its room given back to the text while it is (gtk_widget_set_child_visible in update_empty).
		pen := ink.title
		ui.Box(c).Size(34, 34).Shrink(0).MarginStart(6).Role(ui.RoleNone).Hidden(focused).Draw(func(p *ui.Painter, r ui.Rect) {
			adwEditIcon(p, r.X+(r.W-16)/2, r.Y+(r.H-16)/2, pen)
		})
	}
	rowsStyle(c, row, item, in, first, title, "", look, nil)
	if row.Clicked() {
		input.Focus()
	}
	if input.Changed() {
		row.SetChanged()
	}
	return row
}

// adwEditIcon strokes libadwaita 1.10.0's adw-entry-edit icon, a pencil, in the 16-pixel square at x, y: its outline from (2,14) up to (2,11), along the body to the rounded end at (13.5,2.5) and back to (5,14), and the strokes across it, 2 pixels wide with round ends and corners. It is not mirrored right to left, as the icon has no -rtl variant.
func adwEditIcon(p *ui.Painter, x, y float32, col ui.Color) {
	var path ui.Path
	path.MoveTo(x+2, y+13).LineTo(x+3, y+14)
	path.MoveTo(x+2, y+11).LineTo(x+5, y+14)
	path.MoveTo(x+10, y+3).LineTo(x+13, y+6)
	path.MoveTo(x+2, y+14).LineTo(x+2, y+11).LineTo(x+10.458621, y+2.54137921).
		CubeTo(x+11.3099012, y+1.69009948, x+12.6900988, y+1.69009948, x+13.541379, y+2.54137921).
		CubeTo(x+14.3469534, y+3.3926599, x+14.3099003, y+4.73580599, x+13.4586191, y+5.54137945).
		LineTo(x+5, y+14).Close()
	p.StrokePath(&path, 2, col)
}

// SpinRow creates a settings row setting *value between lo and hi in steps of step: the title, the subtitle below it (an empty subtitle shows none) and at its end the value with the platform's stepping buttons, which joins a Card, a Fieldset or an AccordionItem as ActionRow does. Typing a number in range sets it, Up and Down step it, and the buttons step it and are disabled at the limits, as NumberInput's; Changed reports a new value. The text is the row's one focus stop, named by the title for assistive technology (on GNOME described by the subtitle too, as AdwSpinRow's spin button is; on macOS and Windows the row, a group as ActionRow's, is); the buttons are not stops, as WinUI's NumberBox and GTK's spin buttons are not; on macOS and Windows they leave the focus where it is while pressed, and on GNOME a press on them moves it into the text, as GtkSpinButton's does.
//   - macOS: a SwiftUI Stepper with a format in a grouped Form: at the row's trailing edge, 9 points from the section's, a text field's bezel (macField's, its line laid over the section), 80 points wide, its value at the trailing edge 10 points before an NSStepper without a face inside the field's end, 20 by 24 points, its chevrons those of NumberInput's stepper (formsMacNumberInput) in labelColor over the field and the line between them black at 12.5% (white at 13.8% in Dark Aqua: 223 over 255, 55 over 23); a row of one line is 40.5 points high, taller than a row of text for the field (measured on macOS 27.2).
//   - Windows: a SettingsCard holding a NumberBox with its spin buttons inline (NumberInput's look, formsFluentNumberInput).
//   - GNOME: AdwSpinRow, an action row with a GtkSpinButton at its end (row.spin): no entry around the text, which ends at the buttons, 6 pixels before them (the text's padding of 6), then the − and + buttons as circular buttons, 30 pixels, margin 10px 2px and 6 apart (border-spacing), filled as libadwaita's buttons are (currentColor at 10%, 15% under the pointer and 30% pressed) with value-decrease-symbolic and value-increase-symbolic in the button's color over that face. A click on the row focuses the text (adw_spin_row_grab_focus), and so does a press on a button, wherever the focus was (button_pressed_cb in gtkspinbutton.c grabs it for the spin button); the row draws its focus ring while the text has the focus and is not lit under the pointer then (row.spin's focus ring on .focused). The text is 64 pixels wide, its value at its end, where GtkSpinButton fills the room the title leaves (inferred: what shows is the same).
func SpinRow(c *ui.Context, title, subtitle string, value *float64, lo, hi, step float64) ui.Element {
	p := PlatformOf(c)
	if p == Linux {
		return adwSpinRow(c, title, subtitle, value, lo, hi, step)
	}
	var spin ui.Element
	var row ui.Element
	if p == Windows {
		row = ActionRow(c, title, subtitle, func() { spin = formsFluentNumberInput(c, value, lo, hi, step, title, false) })
	} else {
		// SwiftUI makes a Stepper's row of one line 40.5 points high, the bezel 5.5 points below its top and 9 above its bottom (measured on macOS 27.2), where a row of text is 36.5.
		row = macFormRow(c, title, subtitle, 40.5, func() { spin = macFormStepper(c, value, lo, hi, step, title) })
	}
	if spin.Changed() {
		row.SetChanged()
	}
	return row
}

// macFormRow builds a row of a grouped Form on macOS as ActionRow does, a group named by the title and described by the subtitle, holding at its end what fn builds, and at least least points high where it is alone: SwiftUI makes a row of one line taller than the title's line for a control taller than it. A row with a separator above it takes the separator's point too, and a row of a DisclosureGroup is 12 points shorter, as rowsStyle pads it (both inferred).
func macFormRow(c *ui.Context, title, subtitle string, least float32, fn func()) ui.Element {
	row := ui.Row(c).Shrink(0).Role(ui.RoleGroup).Label(title)
	if subtitle != "" {
		row.Description(subtitle)
	}
	item, in, first := rowsJoin(c, row)
	rowsStyle(c, row, item, in, first, title, subtitle, rowsLook{}, func(ui.Color, bool, bool) { rowsIsolate(c, fn) })
	switch {
	case in == rowsInAccordion:
		least -= 12
	case in != rowsAlone && !first:
		least++
	}
	return row.MinHeight(least)
}

// macFormStepper draws the text field and the stepper of a SwiftUI Stepper with a format in a grouped Form, as SpinRow describes: the field's bezel 82 by 26 points, a point outside its 80 by 24 frame as macField draws it, and the stepper over the frame's last 20 points.
func macFormStepper(c *ui.Context, value *float64, lo, hi, step float64, name string) ui.Element {
	t := c.Theme()
	fill := white
	if t.Dark {
		fill = t.Background.Mix(black, macFieldDarkFill)
	}
	ink := tint(fill, t.Dark, 0.847)
	line := tint(fill, t.Dark, pick(t.Dark, float32(0.125), 0.138))
	// macField lays the bezel's line over the window; AppKit lays it over what is under the field, here the section: 210 over 247 and 57 over 37 (measured on macOS 27.2), macField's black at 14.9% and white at 9.3% over the section's fill.
	border := tint(macFormSectionFill(t), t.Dark, pick(t.Dark, float32(macFieldBorder), macFieldBorderDark))
	row, _ := formsNumberInput(c, value, lo, hi, step, name, false, func(text *string) ui.Element {
		// The value ends 10 points before the stepper: 30 from the bezel's end, which is 9 points from the section's (measured).
		in := TextInput(c, text).Width(82).MarginEnd(-1).TextAlign(ui.End).PaddingEnd(30)
		focused := in.Focused()
		in.DrawOver(func(p *ui.Painter, r ui.Rect) {
			// macField's bezel line and focus ring, with its corner of 7.
			p.Stroke(r, border, 7, 1)
			if focused {
				macFocusRing(p, ui.Rect{X: r.X + 1, Y: r.Y + 1, W: r.W - 2, H: r.H - 2}, 6)
			}
		})
		return in
	}, func(stepBy func(sign float64) ui.Element) {
		// Built Decrease first, for the order assistive technology reads; reversed, Increase shows on top. It covers the last 20 points of the field's 80-point frame, which ends a point inside the bezel.
		col := ui.Column(c).Reverse().Size(20, 24).Shrink(0).MarginStart(-20)
		col.Children(func() {
			for _, sign := range []float64{-1, 1} {
				stepBy(sign).Height(12).Width(20).Draw(func(p *ui.Painter, r ui.Rect) {
					formsStrokeChevron(p, r.X+r.W/2, r.Y+r.H/2, 4.6, 2.35*float32(-sign), 1.75, ink)
				})
			}
		})
		col.DrawOver(func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X + (r.W-13)/2, Y: r.Y + r.H/2 - 0.5, W: 13, H: 1}, line, 0)
		})
	})
	return row
}

// adwSpinRow draws SpinRow on GNOME as AdwSpinRow.
func adwSpinRow(c *ui.Context, title, subtitle string, value *float64, lo, hi, step float64) ui.Element {
	t := c.Theme()
	row := ui.Row(c).Shrink(0)
	item, in, first := rowsJoin(c, row)
	within := row.FocusWithin()
	look := rowsLook{lit: true, hovered: row.Hovered() && !within, pressed: row.Pressed() && !within, focused: within}
	var spin, input ui.Element
	rowsStyle(c, row, item, in, first, title, subtitle, look, func(under ui.Color, _, _ bool) {
		var box ui.Element
		spin, box = formsNumberInput(c, value, lo, hi, step, title, false, func(text *string) ui.Element {
			input = ui.TextInputBase(c, text).Width(64).Padding(6).TextAlign(ui.End).TextColor(adwTextOn(t, under))
			if subtitle != "" {
				input.Description(subtitle)
			}
			return input
		}, func(stepBy func(sign float64) ui.Element) {
			for _, sign := range []float64{-1, 1} {
				b := stepBy(sign).Size(30, 30).Shrink(0).MarginX(2).Radius(15)
				share := float32(0.10)
				switch {
				case b.Pressed():
					share = 0.30
				case b.Hovered():
					share = 0.15
				}
				face := adwInkAlpha(t, under, share)
				b.Background(face)
				// The icon is the button's color, window_fg_color (black at 80% in the light theme), over its face rather than over the row: 47 over the light face of 235, where over the card it would read 51, as libadwaita 1.10.0 renders it (tools/adwaita-reference's settings set).
				icon := adwTextOn(t, face)
				if b.Pressed() {
					// GtkSpinButton's button_pressed_cb grabs the focus for the spin button, which gives it to its text, whatever had it (gtkspinbutton.c of GTK 4.24.1).
					input.Focus()
				}
				b.Draw(func(p *ui.Painter, r ui.Rect) {
					// value-decrease-symbolic and value-increase-symbolic of adwaita-icon-theme 51.0, centered (formsAdwNumberInput's).
					x, y := r.X+(r.W-16)/2, r.Y+(r.H-16)/2
					p.Fill(ui.Rect{X: x + 3, Y: y + 7, W: 10, H: 2}, icon, 0)
					if sign > 0 {
						p.Fill(ui.Rect{X: x + 7, Y: y + 3, W: 2, H: 10}, icon, 0)
					}
				})
			}
		})
		box.Gap(6)
	})
	if row.Clicked() {
		input.Focus()
	}
	if spin.Changed() {
		row.SetChanged()
	}
	return row
}

// ExpanderRow creates a settings row that opens and closes, showing below it what fn builds while *open: its header shows the title, the subtitle below it (an empty subtitle shows none) and the platform's arrow, and a click on it, Space or Enter open and close it; Changed reports a click. The header is the row's one focus stop, a disclosure that assistive technology names by the title, describes by the subtitle and hears as expanded or collapsed. Rows built in fn (ActionRow, SwitchRow and the other rows) join its nested list as they join an AccordionItem; alone it is drawn as a group of its own, and in a Card, a Fieldset or an AccordionItem it joins as a row does. The arrow turns and the nested rows grow into view as it opens, as Collapsible's do.
//   - macOS: a SwiftUI DisclosureGroup in a grouped Form whose label is a title over a subtitle: the disclosure chevron (macDisclosure's, its ink as an NSButton of bezel style disclosure draws it) in a 13-point box 7 points in from the row's leading edge, the regular-weight title 1.5 after it, both centered on the labels, the header's text 14.5 points below its top and above its bottom while closed (a one-line header 44 points high), 1 above its bottom while open; the nested rows those of a DisclosureGroup, without separators (measured on macOS 27.2). AppKit's focus ring goes around the chevron, where SwiftUI puts the focus.
//   - Windows: CommunityToolkit's SettingsExpander: its header the Expander's ToggleButton holding a SettingsCard (70 DIPs high, 68 inside its 1-DIP CardStrokeColorDefault border, padded 16 and 4 before the chevron) on the card's fill while closed and ControlFillColorDefault while open, lit under the pointer and pressed as a clickable SettingsCard, with the 16-DIP AnimatedChevronUpDownSmall chevron (fluentExpanderHeader's, grown from 12 DIPs to 16) centered in a 32-DIP button 8 DIPs from the end, pointing down closed and up open; the content below inside the stroke but on the header's side, corners of 4 at the bottom, its rows a SettingsExpander's items on the card's fill, padded 58,8,44,8 with a line above each, the first's right below the header's fill (measured from WinUI 3 3.3 by tools/winui-reference's winui-rows).
//   - GNOME: AdwExpanderRow: its header an AdwActionRow (54 pixels, 55 with a line above it) whose last suffix is adw-expander-arrow, the up chevron from (3,10) to (8,5) to (13,10) stroked 2 pixels, turned half a turn while closed, 3 pixels after the other suffixes' 6 (margin-start of image.expander-row-arrow), dimmed while closed and opaque in accent_color while open (#0461be and #81d0ff for the default blue, measured); the nested list below on card_shade_color at 50%, a line of card_shade_color above each nested row.
func ExpanderRow(c *ui.Context, title, subtitle string, open *bool, fn func()) ui.Element {
	p := PlatformOf(c)
	t := c.Theme()
	root := ui.Column(c).Shrink(0)
	item, in, first := rowsJoin(c, root)
	hdrIn, hdrFirst := in, first
	if p != Windows {
		// The nested rows' backgrounds and the lit header keep to the group's corners.
		root.Clip()
		if in == rowsAlone {
			// A boxed list or a grouped section of its own, which the header is the first row of.
			hdrIn, hdrFirst = rowsInCard, true
			item.radius = [4]float32{12, 12, 12, 12}
			if p == Linux {
				fill, _ := adwBoxedListCard(t)
				root.Radius(12).Background(fill)
				adwBoxedListShadow(root)
			} else {
				root.Radius(12).Background(macFormSectionFill(t))
			}
		}
	}
	header := &rowsItem{}
	var progress float32
	root.Children(func() {
		parts := ui.CollapsibleBase(c, open)
		tr := parts.Trigger.Shrink(0).Label(title).FocusRing(false)
		if subtitle != "" {
			tr.Description(subtitle)
		}
		if tr.Changed() {
			root.SetChanged()
		}
		progress = parts.Progress()
		header.e = tr
		look := rowsLook{lit: true, hovered: tr.Hovered(), pressed: tr.Pressed(), focused: tr.FocusVisible()}
		switch p {
		case Windows:
			fluentExpanderRowHeader(c, tr, header, hdrIn, hdrFirst, title, subtitle, look, progress)
		case Linux:
			arrow := *open
			rowsStyle(c, tr, header, hdrIn, hdrFirst, title, subtitle, look, func(under ui.Color, _, _ bool) {
				adwExpanderArrow(c, progress, arrow)
			})
		default:
			macExpanderRowHeader(c, tr, header, hdrIn, hdrFirst, title, subtitle, look, progress)
		}
		if !*open && progress == 0 {
			return
		}
		panel := parts.Panel(func() { rowsCollect(c, rowsInAccordion, fn) })
		// Padded and spaced as MyGo pads an AccordionItem's panel, which rows built in an AccordionItem reach across (rowsStyle's rowsInAccordion).
		side := t.Space(3)
		panel.Padding(0, side, side, side).Gap(t.Space(2))
		if p == Windows && (in == rowsAlone || in == rowsInFieldset) {
			// ExpanderContent: ExpanderContentDownBorderThickness 1,0,1,1 of ExpanderContentBorderBrush (CardStrokeColorDefault) and the bottom corners of 4, without a fill of its own; the items are SettingsCards on SettingsCardBackground, CardBackgroundFillColorDefault, which they cover it with (measured: 251 and 43, as the header's card). In a Card or an AccordionItem the rows are the group's items and the panel draws nothing.
			fill, _, stroke := fluentExpanderCard(t)
			panel.Background(fill).BorderWidth(0, 1, 1, 1).BorderColor(stroke).Radius(0, 0, 4, 4)
		}
	})
	// The header's corners, for its focus ring: the group's at the top, and at the bottom too while closed.
	root.Draw(func(*ui.Painter, ui.Rect) {
		header.radius = item.radius
		if progress > 0 {
			header.radius[2], header.radius[3] = 0, 0
		}
	})
	return root
}

// adwExpanderArrow draws AdwExpanderRow's arrow, adw-expander-arrow (libadwaita 1.10.0's adw-expander-arrow.svg: a chevron from (3,10) up to (8,5) and down to (13,10), stroked 2 pixels with round ends), 3 pixels after the suffixes before it, turned half a turn while closed and back as it opens (-gtk-icon-transform rotate(0.5turn) to none), dimmed while closed and accent_color while open (_lists.scss, image.expander-row-arrow).
func adwExpanderArrow(c *ui.Context, progress float32, open bool) {
	t := c.Theme()
	fill, _ := adwBoxedListCard(t)
	col := pick(t.Dark, fill.Mix(adwInk, 0.8*0.55), fill.Mix(white, 0.55))
	if open {
		// accent_color: for the default blue #0461be and #81d0ff, which libadwaita 1.10.0 draws the open arrow in (measured from tools/adwaita-reference's rows set; adwLinkColor's lighter #357fca is the color of a link's text as it renders); for another accent, the accent darkened or lightened as the theme's focus ring takes it.
		col = pick(t.Dark, t.Accent.Mix(black, 0.2), t.Accent.Mix(white, 0.45))
		if t.Accent == adwAccent {
			col = pick(t.Dark, ui.Hex("#0461be"), ui.Hex("#81d0ff"))
		}
	}
	ui.Box(c).Size(16, 16).Shrink(0).MarginStart(3).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
		// strokeDisclosureChevron points right, turned: -90 degrees points up, 90 down.
		strokeDisclosureChevron(p, r.X+8, r.Y+8, 90-180*progress, false, 3, 5, 5, 2, col)
	})
}

// fluentExpanderRowHeader draws a SettingsExpander's header on tr, as ExpanderRow describes.
func fluentExpanderRowHeader(c *ui.Context, tr ui.Element, header *rowsItem, in rowsIn, first bool, title, subtitle string, look rowsLook, progress float32) {
	t := c.Theme()
	rtl := c.Direction() == ui.RightToLeft
	look.content = func(ink rowsInk) {
		ui.Column(c).Grow(1).Shrink(1).MinWidth(0).Children(func() {
			ui.Text(c, title).TextColor(ink.title)
			if subtitle != "" {
				ui.Text(c, subtitle).FontSize(ink.subSize).TextColor(ink.sub)
			}
		})
		ui.Box(c).Size(32, 32).Shrink(0).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			// AnimatedChevronUpDownSmallVisualSource in the 16-DIP AnimatedIcon of SettingsExpander's ExpandCollapseChevron, its 48-unit visual shown at a third: fluentExpanderHeader's chevron of a 12-DIP icon grown by 4/3, 9.3 DIPs wide and 4.7 deep with a stroke of 1.3, down while closed and up while open.
			strokeDisclosureChevron(p, r.X+r.W/2, r.Y+r.H/2, 90+180*progress, rtl, 1.75*4/3, 3.5*4/3, 3.5*4/3, 4.0/3, t.Text)
		})
	}
	rowsStyle(c, tr, header, in, first, title, subtitle, look, nil)
	// SettingsExpanderHeaderPadding 16,16,4,16 and the chevron's Margin 0,0,8,0; the bottom corners square while open, as the Expander's header is in its ExpandDown state.
	tr.PaddingEnd(8).Gap(4)
	if in == rowsAlone || in == rowsInFieldset {
		// The header is the Expander's ToggleButton, its MinHeight the SettingsExpander's 68, holding a SettingsCard of MinHeight 68 inside its 1-DIP border: 70 DIPs high, 68 of fill (measured from WinUI 3 3.3 with CommunityToolkit's SettingsControls 8.2.251219 by tools/winui-reference's winui-rows: the expander's border, 68 DIPs of fill, one line, then the first item). Open, the line below the fill is the first item's, which the content's -2 top margin lays over the header's own bottom border, so the header gives its bottom border up. Closed it is in the Normal state, its fill SettingsCardBackground (CardBackgroundFillColorDefault); open, the Checked state sets no fill, which leaves the ToggleButton's own, ControlFillColorDefault (45 in the dark theme where the card reads 43, as measured).
		bottom := float32(4)
		tr.MinHeight(70)
		if progress > 0 {
			bottom = 0
			tr.MinHeight(69).BorderWidth(1, 1, 0, 1)
			if !look.hovered && !look.pressed {
				tr.Background(t.Surface)
			}
		}
		tr.Radius(4, 4, bottom, bottom)
	}
}

// macExpanderRowHeader draws a DisclosureGroup's header in a grouped Form on tr, as ExpanderRow describes.
func macExpanderRowHeader(c *ui.Context, tr ui.Element, header *rowsItem, in rowsIn, first bool, title, subtitle string, look rowsLook, progress float32) {
	t := c.Theme()
	rtl := c.Direction() == ui.RightToLeft
	focused := look.focused
	look.content = func(ink rowsInk) {
		col := pick(t.Dark, ink.under.Mix(black, macDisclosureLightAlpha), ink.under.Mix(white, macDisclosureDarkAlpha))
		shift := float32(0.25)
		if rtl {
			shift = -shift
		}
		ui.Box(c).Size(13, 13).Shrink(0).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			strokeDisclosureChevron(p, r.X+r.W/2-shift, r.Y+r.H/2-0.25, 90*progress, rtl, 2.5, 3.5, 3.5, 1.5, col)
			if focused {
				macFocusRing(p, r, 3)
			}
		})
		ui.Column(c).Grow(1).Shrink(1).MinWidth(0).Gap(2.5).Children(func() {
			ui.Text(c, title).TextColor(ink.title)
			if subtitle != "" {
				ui.Text(c, subtitle).FontSize(ink.subSize).TextColor(ink.sub)
			}
		})
	}
	rowsStyle(c, tr, header, in, first, title, subtitle, look, nil)
	// The chevron's box 7 points in, the title 1.5 after it (22 from the section's edge), the labels 14.5 below the header's top and as far above its bottom, 1 while open, where the nested rows follow (measured on macOS 27.2 from reference.swift's settings scene: a closed one-line header 44 points high, the title's ink 20.5 points below its top and 18.5 above its bottom; the first nested row's title 73 points below the top of a header with a subtitle).
	tr.AlignItems(ui.Center).Gap(1.5).Padding(14.5, 10, 14.5-13.5*progress, 10).PaddingStart(7).PaddingEnd(10)
}
