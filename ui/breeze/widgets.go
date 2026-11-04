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

// Radio creates a radio button that selects value into *selected, with a label, as the platform's radio buttons look. Put radio buttons in a RadioGroup. It is ui.Radio's counterpart.
func Radio[T comparable](c *ui.Context, selected *T, value T, label string) ui.Element {
	row := ui.RadioBase(c, selected, value).AlignItems(ui.Center).FocusRing(false)
	on := *selected == value
	var dot ui.Element
	row.Children(func() {
		dot = ui.Box(c).Shrink(0)
		if label != "" {
			ui.Text(c, label)
		}
	})
	switch PlatformOf(c) {
	case Windows:
		fluentRadio(c, row, dot, on)
	case Linux:
		adwRadio(c, row, dot, on)
	default:
		macRadio(c, row, dot, on)
	}
	return row
}

// Toggle creates a button that stays pressed while *on, showing label, as the platform draws a toggle button. It is ui.Toggle's counterpart.
func Toggle(c *ui.Context, on *bool, label string) ui.Element {
	b := ui.ToggleBase(c, on)
	switch PlatformOf(c) {
	case Windows:
		fluentToggle(c, b, *on)
	case Linux:
		adwToggle(c, b, *on)
	default:
		macToggle(c, b, *on)
	}
	if label != "" {
		b.Children(func() { ui.Text(c, label).SingleLine() })
	}
	return b
}

// fieldKind tells the text fields apart, which platforms draw differently.
type fieldKind uint8

const (
	singleLine fieldKind = iota
	multiLine
	searchLine
)

// styleField gives a text field of package ui, the box around its text, the platform's look. The field keeps MyGo's own behavior: its border turns to the accent while it has the focus, and to the theme's Danger color with an Error.
func styleField(c *ui.Context, f ui.Element, kind fieldKind) ui.Element {
	switch PlatformOf(c) {
	case Windows:
		fluentField(c, f, kind)
	case Linux:
		adwField(c, f, kind)
	default:
		macField(c, f, kind)
	}
	return f
}

// TextInput creates a single-line text input editing *value, in the box the platform draws around its text fields. It is ui.TextInput's counterpart.
func TextInput(c *ui.Context, value *string) ui.Element {
	return styleField(c, ui.TextInput(c, value), singleLine)
}

// TextArea creates a multi-line text input editing *value, in the box the platform draws around its text views. It is ui.TextArea's counterpart.
func TextArea(c *ui.Context, value *string) ui.Element {
	return styleField(c, ui.TextArea(c, value), multiLine)
}

// SearchField creates a field for a search query with a magnifying glass, and a button clearing it while it holds text, as the platform's search fields look. It is ui.SearchField's counterpart.
func SearchField(c *ui.Context, query *string) ui.Element {
	return styleField(c, ui.SearchField(c, query), searchLine)
}

// Select creates a drop-down choosing *selected among options, as the platform's pop-up button, combo box or drop-down looks, opening a list of the options. The arrows move among them and Enter chooses. It is ui.Select's counterpart.
func Select(c *ui.Context, selected *string, options []string) ui.Element {
	sel := ui.SelectBase(c, selected)
	p := PlatformOf(c)
	switch p {
	case Windows:
		fluentSelect(c, sel.Trigger, *selected)
	case Linux:
		adwSelect(c, sel.Trigger, *selected)
	default:
		macSelect(c, sel.Trigger, *selected)
	}
	sel.Popup(func(panel ui.Element) {
		// Options past the room the window has scroll.
		panel.ClipY()
		switch p {
		case Windows:
			fluentMenu(c, panel)
		case Linux:
			adwMenu(c, panel)
		default:
			macMenu(c, panel)
		}
		for _, opt := range options {
			item := sel.Item(opt)
			chosen := opt == *selected
			switch p {
			case Windows:
				fluentMenuItem(c, item, opt, chosen)
			case Linux:
				adwMenuItem(c, item, opt, chosen)
			default:
				macMenuItem(c, item, opt, chosen)
			}
		}
	})
	return sel.Trigger
}

// Slider creates a slider setting *value between lo and hi, as the platform's sliders look. It is ui.Slider's counterpart.
func Slider(c *ui.Context, value *float64, lo, hi float64) ui.Element {
	s := ui.SliderBase(c, value, lo, hi)
	frac := float32(0)
	if hi > lo {
		frac = float32(min(max((*value-lo)/(hi-lo), 0), 1))
	}
	switch PlatformOf(c) {
	case Windows:
		fluentSlider(c, s, frac)
	case Linux:
		adwSlider(c, s, frac)
	default:
		macSlider(c, s, frac)
	}
	return s
}

// Progress creates a progress bar at value, from 0 to 1, as the platform's progress bars look; a negative value shows progress of unknown length, moving. It is ui.Progress's counterpart.
func Progress(c *ui.Context, value float64) ui.Element {
	e := ui.Box(c).Role(ui.RoleProgress).Range(0, 1, value)
	switch PlatformOf(c) {
	case Windows:
		fluentProgress(c, e, value)
	case Linux:
		adwProgress(c, e, value)
	default:
		macProgress(c, e, value)
	}
	return e
}

// Spinner creates an indicator of work of unknown length, turning, as the platform's spinning progress indicator, progress ring or spinner. It is ui.Spinner's counterpart.
func Spinner(c *ui.Context) ui.Element {
	e := ui.Box(c).Shrink(0).Role(ui.RoleProgress).Range(0, 1, -1)
	switch PlatformOf(c) {
	case Windows:
		fluentSpinner(c, e)
	case Linux:
		adwSpinner(c, e)
	default:
		macSpinner(c, e)
	}
	return e
}

// Segmented creates a segmented control choosing *selected among labels, as the platform's segmented control, selector bar or toggle group looks. It is ui.Segmented's counterpart.
func Segmented(c *ui.Context, selected *int, labels ...string) ui.Element {
	parts := ui.SegmentedBase(c, selected, len(labels))
	p := PlatformOf(c)
	switch p {
	case Windows:
		fluentSegmentedTrack(c, parts.Track)
	case Linux:
		adwSegmentedTrack(c, parts.Track)
	default:
		macSegmentedTrack(c, parts.Track)
	}
	parts.Track.Children(func() {
		for i, label := range labels {
			seg := parts.Segment(i)
			// A separator shows between two segments neither of which is chosen, as AppKit draws them.
			sep := i > 0 && i != *selected && i-1 != *selected
			switch p {
			case Windows:
				fluentSegment(c, seg, i == *selected)
			case Linux:
				adwSegment(c, seg, i == *selected, sep)
			default:
				macSegment(c, seg, i == *selected, sep)
			}
			seg.Children(func() { ui.Text(c, label).SingleLine() })
		}
	})
	return parts.Track
}

// Tabs creates tabs choosing *selected among labels, as the platform's tab view, selector bar or view switcher looks. Build the chosen page after it. It is ui.Tabs's counterpart.
func Tabs(c *ui.Context, selected *int, labels ...string) ui.Element {
	parts := ui.TabsBase(c, selected, len(labels))
	p := PlatformOf(c)
	switch p {
	case Windows:
		fluentTabList(c, parts.List)
	case Linux:
		adwTabList(c, parts.List)
	default:
		macTabList(c, parts.List)
	}
	parts.List.Children(func() {
		for i, label := range labels {
			tab := parts.Tab(i)
			sep := i > 0 && i != *selected && i-1 != *selected
			switch p {
			case Windows:
				fluentTab(c, tab, i == *selected)
			case Linux:
				adwTab(c, tab, i == *selected)
			default:
				macSegment(c, tab, i == *selected, sep)
			}
			tab.Children(func() { ui.Text(c, label).SingleLine() })
		}
	})
	return parts.List
}

// ---- widgets of package ui, in the platform's theme (Use) ----

// Link creates text that opens url, in the accent color. It is package ui's.
func Link(c *ui.Context, label, url string) ui.Element { return ui.Link(c, label, url) }

// Divider creates a line between parts of a view. It is package ui's.
func Divider(c *ui.Context) ui.Element { return ui.Divider(c) }

// Badge creates a small label for a count or a status. It is package ui's.
func Badge(c *ui.Context, text string) ui.Element { return ui.Badge(c, text) }

// MenuButton creates a button opening the system's own menu below it, which build fills. It is package ui's: the menu is the platform's native one.
func MenuButton(c *ui.Context, label string, build func(m *ui.Menu)) ui.Element {
	return ui.MenuButton(c, label, build)
}

// PopoverMenuButton creates a button opening a menu that MyGo draws in the window. It is package ui's.
func PopoverMenuButton(c *ui.Context, label string, build func(m *ui.Menu)) ui.Element {
	return ui.PopoverMenuButton(c, label, build)
}

// Combobox creates a field choosing *selected among options, whose text filters them, in the box the platform draws around its text fields. It is ui.Combobox's counterpart, with MyGo's own list of options.
func Combobox(c *ui.Context, selected *string, options []string) ui.Element {
	return styleField(c, ui.Combobox(c, selected, options), singleLine)
}

// NumberInput creates a field editing *value between lo and hi in steps of step. It is package ui's.
func NumberInput(c *ui.Context, value *float64, lo, hi, step float64) ui.Element {
	return ui.NumberInput(c, value, lo, hi, step)
}

// Stepper creates buttons stepping *value between lo and hi by step. It is package ui's.
func Stepper(c *ui.Context, value *float64, lo, hi, step float64) ui.Element {
	return ui.Stepper(c, value, lo, hi, step)
}

// Field creates a labeled row around the control fn builds, naming it for assistive technology. It is package ui's.
func Field(c *ui.Context, label string, fn func()) ui.Element { return ui.Field(c, label, fn) }

// RadioGroup creates the group that radio buttons built in fn choose within, which the arrows move through. It is package ui's.
func RadioGroup(c *ui.Context, fn func()) ui.Element { return ui.RadioGroup(c, fn) }

// ToggleGroup creates a row of toggles joined together. It is package ui's.
func ToggleGroup(c *ui.Context, fn func()) ui.Element { return ui.ToggleGroup(c, fn) }

// Modal shows a dialog built by fn over the window while *open, as the platform draws its dialogs; clicking outside it or pressing Escape sets *open to false, which Dismissed reports. It is ui.Modal's counterpart.
func Modal(c *ui.Context, open *bool, fn func()) ui.Element {
	if !*open {
		return ui.Element{}
	}
	p := PlatformOf(c)
	w, h := c.Size()
	return ui.DialogBase(c, open, func(back, panel ui.Element) {
		panel.MaxWidth(w - 40).MaxHeight(h - 40)
		Preview(c, p, func() {
			switch p {
			case Windows:
				fluentDialog(c, back, panel)
			case Linux:
				adwDialog(c, back, panel)
			default:
				macDialog(c, back, panel)
			}
			fn()
		})
	})
}

// AlertDialog shows an alert over the window while *open, as the platform draws its alerts: title, message, and buttons, the last of which is the default, in the accent color, with the focus, which Enter clicks. Each platform places the buttons its own way: macOS and GNOME put the default last, Windows first, while the index returned is the same everywhere. Escape closes the alert choosing no button; a click outside it does nothing. It returns the index of the button clicked, in the frame it is, which closes the alert, ui.AlertDismissed in the frame Escape closes it, and -1 otherwise. It is ui.AlertDialog's counterpart.
func AlertDialog(c *ui.Context, open *bool, title, message string, buttons ...string) int {
	if !*open {
		return -1
	}
	p := PlatformOf(c)
	w, _ := c.Size()
	chosen, dismissed := -1, false
	ui.Overlay(c, func() {
		back := ui.Box(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).Center().Modal()
		dismissed = back.OverlayShortcut(0, ui.KeyEscape)
		back.Children(func() {
			panel := ui.Column(c).MaxWidth(w - 40).Role(ui.RoleAlertDialog).Label(title).Description(message)
			// button builds button i of the alert, the last as the default.
			button := func(i int) ui.Element {
				var b ui.Element
				if i == len(buttons)-1 {
					b = PrimaryButton(c, buttons[i]).AutoFocus()
				} else {
					b = Button(c, buttons[i])
				}
				if b.Clicked() {
					chosen = i
				}
				return b
			}
			Preview(c, p, func() {
				switch p {
				case Windows:
					fluentAlert(c, back, panel, title, message, len(buttons), button)
				case Linux:
					adwAlert(c, back, panel, title, message, len(buttons), button)
				default:
					macAlert(c, back, panel, title, message, len(buttons), button)
				}
			})
		})
	})
	if dismissed && chosen < 0 {
		chosen = ui.AlertDismissed
	}
	if chosen != -1 {
		*open = false
	}
	return chosen
}

// Popover creates a panel beside anchor while *open. It is package ui's.
func Popover(c *ui.Context, anchor ui.Element, open *bool, fn func()) ui.Element {
	return ui.Popover(c, anchor, open, fn)
}

// Table creates a table of n rows in columns, whose cells cell builds. It is package ui's.
func Table(c *ui.Context, s *ui.ListState, columns []ui.TableColumn, n int, cell func(row, col int)) ui.Element {
	return ui.Table(c, s, columns, n, cell)
}

// Tree creates a tree of the items fn builds with TreeItem. It is package ui's.
func Tree(c *ui.Context, fn func()) ui.Element { return ui.Tree(c, fn) }

// TreeItem creates an item of a Tree, open while *open, with the items children builds. It is package ui's.
func TreeItem(c *ui.Context, label string, open *bool, children func()) ui.Element {
	return ui.TreeItem(c, label, open, children)
}
