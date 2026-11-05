// Package breeze draws MyGo's native UI widgets as the desktop the app runs on draws its own: macOS's controls on a Mac, Windows 11's on Windows and GNOME's libadwaita controls on Linux.
//
// Its widgets have the names and signatures of package ui's, so a view swaps ui.Button for breeze.Button and keeps everything else: they return ui.Element, ask Clicked and take more styles as ui's do. The pointer, the keyboard, the focus and accessibility are MyGo's, from the widgets' bases (ui.ButtonBase, ui.CheckboxBase, ui.SwitchBase); Breeze UI draws them. Widgets it does not draw itself yet are package ui's, styled by the platform's theme (Use).
package breeze

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// Button creates a push button showing label, as the platform's push buttons look; inside a Toolbar, as the platform's toolbar items: a bordered NSToolbarItem on macOS, whose image-only items side by side share one platter, an AppBarButton on Windows, and on GNOME a raised text button or, without a label, a flat image button. Ask Clicked whether it was clicked; give it other content with Children and an empty label. It is ui.Button's counterpart.
func Button(c *ui.Context, label string) ui.Element { return button(c, label, false) }

// PrimaryButton creates a button in the accent color for the main action, as the platform draws its default button (macOS), accent button (Windows) or suggested action (GNOME); inside a Toolbar on macOS, a prominent toolbar item, an accent capsule. It is ui.PrimaryButton's counterpart.
func PrimaryButton(c *ui.Context, label string) ui.Element { return button(c, label, true) }

func button(c *ui.Context, label string, primary bool) ui.Element {
	b := ui.ButtonBase(c)
	kind := idiomText
	switch {
	case primary:
		kind = idiomPrimary
	case label == "":
		kind = idiomImage
	}
	tb, i := idiomToolbarItem(c, b, kind)
	switch p := PlatformOf(c); {
	case p == Windows && tb != nil && !primary:
		idiomFluentAppBarButton(c, b, label == "", false)
	case p == Windows:
		fluentButton(c, b, primary)
	case p == Linux && tb != nil && kind == idiomImage:
		idiomAdwFlatButton(c, b, false)
	case p == Linux:
		// Text buttons and the suggested action stay raised in a toolbar (_toolbars.scss).
		adwButton(c, b, primary)
	case tb != nil:
		idiomMacToolbarItem(c, b, tb, i, primary)
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
	formsJoinCheckGroup(c, checked)
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

// Toggle creates a button that stays pressed while *on, showing label, as the platform draws a toggle button; inside a ToggleGroup, as a segment of the group; inside a Toolbar, as the platform's toolbar toggle: a one-item NSToolbarItemGroup on macOS, its platter in the accent while on, an AppBarToggleButton on Windows, and on GNOME a raised toggle button or, without a label, a flat image toggle. It is ui.Toggle's counterpart.
func Toggle(c *ui.Context, on *bool, label string) ui.Element {
	b := ui.ToggleBase(c, on)
	grouped, first := nextInToggleGroup(c)
	kind := idiomToggle
	if grouped {
		kind = idiomOther
	}
	tb, i := idiomToolbarItem(c, b, kind)
	bar := tb != nil && !grouped
	switch p := PlatformOf(c); {
	case p == Windows && bar:
		idiomFluentAppBarButton(c, b, label == "", *on)
	case p == Windows:
		fluentToggle(c, b, *on)
	case p == Linux && grouped:
		adwLinkedToggle(c, b, *on, first)
	case p == Linux && bar && label == "":
		idiomAdwFlatButton(c, b, *on)
	case p == Linux:
		adwToggle(c, b, *on)
	case grouped:
		// A segment of an NSSegmentedControl that selects any number of segments.
		macSegment(c, b, *on, false)
	case bar:
		idiomMacToolbarItem(c, b, tb, i, *on)
		if label == "" {
			// An image toggle is a group's 36-point square item (measured: a one-image NSToolbarItemGroup is 36 by 36).
			b.MinWidth(36).Padding(0)
		}
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
			switch p {
			case Windows:
				fluentTab(c, tab, i == *selected)
			case Linux:
				adwTab(c, tab, i == *selected)
			default:
				macTab(c, tab, i == *selected)
			}
			tab.Children(func() { ui.Text(c, label).SingleLine() })
		}
	})
	return parts.List
}

// ---- widgets of package ui, in the platform's theme (Use) ----

// Combobox creates a field choosing *selected among options, whose text filters them, in the box the platform draws around its text fields. It is ui.Combobox's counterpart, with MyGo's own list of options.
func Combobox(c *ui.Context, selected *string, options []string) ui.Element {
	return styleField(c, ui.Combobox(c, selected, options), singleLine)
}

// NumberInput creates a field editing *value between lo and hi in steps of step, with buttons stepping it, as the platform draws a number field: an NSTextField with an NSStepper beside it on macOS, a NumberBox with its spin buttons inline on Windows, a GtkSpinButton with its − and + buttons after the text on GNOME. What is typed applies once it is a number in range, Up and Down step it, the buttons are disabled at the limits and named Decrease and Increase, and the text shows the value rounded to step's decimals once the input loses the focus; Changed reports a new value. It is ui.NumberInput's counterpart, rebuilt from ui.TextInputBase and ui.ButtonBase as ui.NumberInput is built, so its keys, focus order (the text, Decrease, Increase) and names are MyGo's: the buttons stay tab stops, as WinUI's and GTK's are not, and step once a click, as AppKit's, WinUI's and GTK's repeat while held. On Windows and GNOME the box holds the buttons and fills the width a column gives it, as a NumberBox and a GtkSpinButton do; on macOS the field keeps MyGo's 80 DIPs, the stepper beside it.
func NumberInput(c *ui.Context, value *float64, lo, hi, step float64) ui.Element {
	switch PlatformOf(c) {
	case Windows:
		return formsFluentNumberInput(c, value, lo, hi, step)
	case Linux:
		return formsAdwNumberInput(c, value, lo, hi, step)
	}
	return formsMacNumberInput(c, value, lo, hi, step)
}

// RadioGroup creates the group that radio buttons built in fn choose within, which the arrows move through. It is package ui's.
func RadioGroup(c *ui.Context, fn func()) ui.Element { return ui.RadioGroup(c, fn) }

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

// Table creates a table of n rows in columns, whose cells cell builds, as package ui's Table does, its rows drawn as the platform draws them: an inset NSTableView's rounded selection, accent while the table has the focus and gray otherwise (macOS); ListView's items, with their backplate and the accent pill on the chosen one, as WinUI 3 has no table (Windows); a GtkColumnView's rows with the accent-tinted row (GNOME). The choice, the sort, the columns' widths and order, the keys, the focus and accessibility are MyGo's, as are the header, in the platform's theme, the cells' padding and the rows' least height of Space(8). A table choosing cells (ListState.CellSelect) keeps MyGo's cells in the platform's theme, and so does the room a row leaves after its last column when every column is narrower than the table. It is ui.Table's counterpart.
func Table(c *ui.Context, s *ui.ListState, columns []ui.TableColumn, n int, cell func(row, col int)) ui.Element {
	p := PlatformOf(c)
	if s != nil && s.CellSelect {
		tb := ui.Table(c, s, columns, n, cell)
		collStyleTable(p, tb, collBg(p, c.Theme()))
		return tb
	}
	rows := collNewRows(c, p, s)
	tb := ui.Table(c, s, columns, n, func(row, col int) {
		rows.cell(c, row, col, columns[col].Align, 0, false, false, func() { cell(row, col) })
	})
	collStyleTable(p, tb, rows.bg)
	return tb
}

// Tree creates a tree of the items fn builds with TreeItem, as package ui's Tree does, drawn as the platform draws a tree's rows: an inset NSOutlineView's, its selection accent while the tree has the focus and gray otherwise (macOS), TreeView's items with their backplate (Windows), or a GtkListView's rows in a view's background (GNOME). The items are built in a theme of the platform's row colors, which MyGo's TreeItem takes as its caller marks it Selected; the keys, the focus and accessibility are MyGo's. It is ui.Tree's counterpart.
func Tree(c *ui.Context, fn func()) ui.Element {
	p := PlatformOf(c)
	t := c.Theme()
	th := *t
	tree := ui.Tree(c, func() {
		items := ui.Themed(c, &th).Background(ui.Transparent).Role(ui.RoleNone)
		// The focus of the frame before, as the column holding the items has it before they build.
		collTreeTheme(p, t, &th, items.FocusWithin())
		// Themed took the text's color from th before collTreeTheme set it.
		items.TextColor(th.Text)
		items.Children(fn)
	})
	if p == Linux {
		tree.Background(adwListViewBg(t))
	}
	return tree
}

// TreeItem creates an item of a Tree, open while *open, with the items children builds, as package ui's TreeItem does, at the platform's height, insets, corners and indentation of its rows (collTreeItem): Selected shows it chosen in the platform's colors inside this package's Tree. MyGo's disclosure arrow stays, and WinUI's accent pill on the chosen item is not drawn. It is ui.TreeItem's counterpart.
func TreeItem(c *ui.Context, label string, open *bool, children func()) ui.Element {
	p := PlatformOf(c)
	depth := collTreeDepth(c)
	item := ui.TreeItem(c, label, open, collTreeChildren(c, children))
	collTreeItem(c, p, item, depth, label)
	return item
}
