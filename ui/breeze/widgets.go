// Package breeze draws MyGo's native UI widgets as the desktop the app runs on draws its own: macOS's controls on a Mac, Windows 11's on Windows and GNOME's libadwaita controls on Linux.
//
// Its widgets have the names and signatures of package ui's, so a view swaps ui.Button for breeze.Button and keeps everything else: they return ui.Element, ask Clicked and take more styles as ui's do. The pointer, the keyboard, the focus and accessibility are MyGo's, from the widgets' bases (ui.ButtonBase, ui.CheckboxBase, ui.SwitchBase); Breeze UI draws them. Widgets it does not draw itself yet are package ui's, styled by the platform's theme (Use).
package breeze

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// Button creates a push button showing label, as the platform's push buttons look. Ask Clicked whether it was clicked; give it other content with Children and an empty label. It is ui.Button's counterpart.
func Button(c *ui.Context, label string) ui.Element { return button(c, label, false) }

// PrimaryButton creates a button in the accent color for the main action, as the platform draws its default button (macOS), accent button (Windows) or suggested action (GNOME). It is ui.PrimaryButton's counterpart.
func PrimaryButton(c *ui.Context, label string) ui.Element { return button(c, label, true) }

func button(c *ui.Context, label string, primary bool) ui.Element {
	b := ui.ButtonBase(c)
	switch PlatformOf(c) {
	case Windows:
		fluentButton(c, b, primary)
	case Linux:
		adwButton(c, b, primary)
	default:
		macButton(c, b, primary)
	}
	if label != "" {
		b.Children(func() { ui.Text(c, label).SingleLine() })
	}
	return b
}

// Checkbox creates a check box toggling *checked, with a label, as the platform's check boxes look. It is ui.Checkbox's counterpart.
func Checkbox(c *ui.Context, checked *bool, label string) ui.Element {
	row := ui.CheckboxBase(c, checked).AlignItems(ui.Center).FocusRing(false)
	on := *checked
	var box ui.Element
	row.Children(func() {
		box = ui.Box(c).Shrink(0)
		if label != "" {
			ui.Text(c, label)
		}
	})
	switch PlatformOf(c) {
	case Windows:
		fluentCheckbox(c, row, box, on)
	case Linux:
		adwCheckbox(c, row, box, on)
	default:
		macCheckbox(c, row, box, on)
	}
	return row
}

// Switch creates a switch toggling *on, as the platform's switches look. Name it with Label for assistive technology, or put it in a field. It is ui.Switch's counterpart.
func Switch(c *ui.Context, on *bool) ui.Element {
	sw := ui.SwitchBase(c, on)
	switch PlatformOf(c) {
	case Windows:
		fluentSwitch(c, sw, *on)
	case Linux:
		adwSwitch(c, sw, *on)
	default:
		macSwitch(c, sw, *on)
	}
	return sw
}

// knob returns where a switch's knob is, from 0 off to 1 on, easing over d as the switch is turned.
func knob(sw ui.Element, on bool, d time.Duration) float32 {
	target := float32(0)
	if on {
		target = 1
	}
	return sw.Animate("breeze-knob", target, d)
}

// along returns where along its track a knob at pos is drawn: from the left in a left-to-right window, from the right in a right-to-left one, as switches mirror there.
func along(c *ui.Context, pos float32) float32 {
	if c.Direction() == ui.RightToLeft {
		return 1 - pos
	}
	return pos
}

// radii returns r for every corner, as Painter.FocusRing takes them.
func radii(r float32) [4]float32 { return [4]float32{r, r, r, r} }

// ---- widgets of package ui, styled by the platform's theme ----

// TextInput creates a single-line text input editing *value. It is package ui's, in the platform's theme when the view calls Use.
func TextInput(c *ui.Context, value *string) ui.Element { return ui.TextInput(c, value) }

// TextArea creates a multi-line text input editing *value. It is package ui's, in the platform's theme when the view calls Use.
func TextArea(c *ui.Context, value *string) ui.Element { return ui.TextArea(c, value) }

// Select creates a drop-down choosing *selected among options. It is package ui's, in the platform's theme when the view calls Use.
func Select(c *ui.Context, selected *string, options []string) ui.Element {
	return ui.Select(c, selected, options)
}

// Slider creates a slider setting *value between lo and hi. It is package ui's, in the platform's theme when the view calls Use.
func Slider(c *ui.Context, value *float64, lo, hi float64) ui.Element {
	return ui.Slider(c, value, lo, hi)
}

// Progress creates a progress bar at value, from 0 to 1. It is package ui's, in the platform's theme when the view calls Use.
func Progress(c *ui.Context, value float64) ui.Element { return ui.Progress(c, value) }

// Tabs creates tabs choosing *selected among labels. It is package ui's, in the platform's theme when the view calls Use.
func Tabs(c *ui.Context, selected *int, labels ...string) ui.Element {
	return ui.Tabs(c, selected, labels...)
}

// Segmented creates a segmented control choosing *selected among labels. It is package ui's, in the platform's theme when the view calls Use.
func Segmented(c *ui.Context, selected *int, labels ...string) ui.Element {
	return ui.Segmented(c, selected, labels...)
}
