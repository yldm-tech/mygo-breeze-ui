package breeze

import (
	"sync"

	"github.com/egoist/mygo/ui"
)

// Link creates text that opens url, in the platform's link color: AppKit's linkColor, WinUI's AccentTextFillColorPrimary or libadwaita's accent_color, underlined on Windows and GNOME, and under the pointer on macOS. It is ui.Link's counterpart, with its routing.
func Link(c *ui.Context, label, url string) ui.Element {
	l := ui.Link(c, label, url)
	t := c.Theme()
	switch PlatformOf(c) {
	case Windows:
		// A Hyperlink in text: underlined, the underline gone under the pointer in WinUI though MyGo's link adds one there.
		l.TextColor(fluentLinkColor(c, l.Hovered(), l.Pressed())).Underline()
	case Linux:
		// libadwaita underlines links whether or not the pointer is over them.
		l.TextColor(adwLinkColor(t)).Underline()
	default:
		l.TextColor(pick(t.Dark, ui.Hex("#0068da"), ui.Hex("#419cff")))
	}
	return l
}

// Divider creates a line between parts of a view, in the platform's separator color. It is ui.Divider's counterpart.
func Divider(c *ui.Context) ui.Element {
	d := ui.Divider(c)
	t := c.Theme()
	switch PlatformOf(c) {
	case Windows:
		d.Background(fluentDivider(t))
	case Linux:
		d.Background(adwSeparator(t))
	default:
		d.Background(macSeparator(t.Background, t.Dark))
	}
	return d
}

// Badge creates a small label for a count or a status, as Windows draws an InfoBadge and GNOME the badge of a view switcher. macOS has no badge control of its own; there it is ui.Badge in the platform's theme. It is ui.Badge's counterpart.
func Badge(c *ui.Context, text string) ui.Element {
	b := ui.Badge(c, text)
	switch PlatformOf(c) {
	case Windows:
		fluentBadge(c, b)
	case Linux:
		adwBadge(c, b)
	}
	return b
}

// styleMenuButton gives a button opening a menu the platform's look of its pull-down or menu button.
func styleMenuButton(c *ui.Context, b ui.Element) ui.Element {
	switch PlatformOf(c) {
	case Windows:
		fluentButton(c, b, false)
	case Linux:
		// libadwaita's buttons have no border; MyGo's menu button has one.
		adwButton(c, b, false)
		b.Border(0, ui.Color{})
	default:
		macButton(c, b, false)
		b.Border(0, ui.Color{})
	}
	return b
}

// MenuButton creates a button, drawn as the platform's pull-down button, opening the system's own menu below it, which build fills. It is ui.MenuButton's counterpart: the menu is native on every platform.
func MenuButton(c *ui.Context, label string, build func(m *ui.Menu)) ui.Element {
	return styleMenuButton(c, ui.MenuButton(c, label, build))
}

// PopoverMenuButton creates a button, drawn as the platform's pull-down button, opening a menu that MyGo draws in the window in the platform's theme. It is ui.PopoverMenuButton's counterpart.
func PopoverMenuButton(c *ui.Context, label string, build func(m *ui.Menu)) ui.Element {
	return styleMenuButton(c, ui.PopoverMenuButton(c, label, build))
}

// Popover shows fn's elements in a panel beside anchor while *open, as the platform draws a popover or a flyout; clicking outside it or pressing Escape sets *open to false. It is ui.Popover's counterpart.
func Popover(c *ui.Context, anchor ui.Element, open *bool, fn func()) ui.Element {
	if !*open {
		return ui.Element{}
	}
	p := PlatformOf(c)
	return ui.PopoverBase(c, anchor, open, func(panel ui.Element) {
		panel.MinWidth(anchor.Bounds().W)
		Preview(c, p, func() {
			switch p {
			case Windows:
				fluentFlyout(c, panel, anchor.Bounds().W)
			case Linux:
				adwPopover(c, panel)
			default:
				macPopover(c, panel)
			}
			fn()
		})
	})
}

// toggleGroups holds the windows building a ToggleGroup's toggles, which Toggle draws as segments of it, with how many it has built.
var toggleGroups sync.Map // *ui.Context -> *toggleGroup

type toggleGroup struct{ depth, built int }

// nextInToggleGroup reports whether a toggle being built is in a ToggleGroup, and whether it is the group's first.
func nextInToggleGroup(c *ui.Context) (grouped, first bool) {
	v, ok := toggleGroups.Load(c)
	if !ok {
		return false, false
	}
	g := v.(*toggleGroup)
	g.built++
	return true, g.built == 1
}

// ToggleGroup creates a row of the toggles fn builds, joined as the platform joins them: a segmented control on macOS, linked buttons on GNOME, and on Windows, which has no joined toggles, ToggleButtons side by side. The arrows move among them, as in ui.ToggleGroup, whose counterpart it is.
func ToggleGroup(c *ui.Context, fn func()) ui.Element {
	p := PlatformOf(c)
	outer, nested := toggleGroups.Load(c)
	toggleGroups.Store(c, &toggleGroup{})
	g := ui.ToggleGroup(c, fn)
	if nested {
		toggleGroups.Store(c, outer)
	} else {
		toggleGroups.Delete(c)
	}
	switch p {
	case Windows:
		fluentToggleGroup(c, g)
	case Linux:
		adwToggleGroup(c, g)
	default:
		macSegmentedTrack(c, g)
		g.Padding(0).Gap(0).Border(0, ui.Color{})
	}
	return g
}

// Field creates a labeled control: the label names the control fn builds for assistive technology. In a Form the labels line up in a column before the controls, as macOS lays out a form; elsewhere the label is above the control, 8 DIPs above it on Windows, as a TextBox's header is. It is ui.Field's counterpart.
func Field(c *ui.Context, label string, fn func()) ui.Element {
	f := ui.Field(c, label, fn)
	if PlatformOf(c) == Windows {
		f.Gap(8)
	}
	return f
}

// Stepper creates a pair of arrows changing *value by step between lo and hi: a click steps, holding an arrow keeps stepping, and Up and Down step while it has the focus. Its box is the platform's: AppKit's NSStepper on macOS, a button's face elsewhere; the arrows and their behavior are MyGo's. It is ui.Stepper's counterpart.
func Stepper(c *ui.Context, value *float64, lo, hi, step float64) ui.Element {
	s := ui.Stepper(c, value, lo, hi, step)
	t := c.Theme()
	switch PlatformOf(c) {
	case Windows:
		s.Width(32).Radius(4).Background(t.Surface).Border(1, ui.Color{}).BorderLinearGradient(fluentElevation(t, 32))
	case Linux:
		s.Width(34).Radius(9).Background(t.Surface).Border(0, ui.Color{})
	default:
		macStepper(c, s)
	}
	return s
}
