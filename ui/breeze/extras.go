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

// styleMenuButton gives b, a button opening a menu, the platform's look of its pull-down or menu button and its children: label, if any, and the platform's arrow. tb is the Toolbar b is item i of, nil outside one.
//   - macOS: outside a toolbar a push button (macButton) with the arrow 4 points after its label; in a toolbar an NSMenuToolbarItem, its pull-down button a 36-point capsule (idiomMacPlatter) padded 12 before its 13-point label and 13 after the 7.5-point arrow, 4 points after the label (Share's is 71.5 points wide: 12 + 34.91 + 4.1 + 7.5 + 13, measured on macOS 27.2). The arrow outside a toolbar is inferred from the toolbar's.
//   - Windows: outside a toolbar a DropDownButton, a button (fluentButton) with the 12-DIP ChevronDownSmall in TextFillColorSecondary 8 DIPs after its label; in a toolbar an AppBarButton with a flyout (idiomFluentAppBarButton), its SubItemChevron 5 DIPs after the label (AppBarButtonSubItemChevronLabelOnRightMargin -7,20,12,0 after the label's margin of 12) in TextFillColorSecondary, Tertiary pressed.
//   - GNOME: a GtkMenuButton, whose button carries .arrow-button beside .text-button when it shows a label and the arrow, or .image-button when it shows only the arrow (gtkmenubutton.c's update_style_classes, GTK 4.24.0): padded 5px 9px with 4 pixels between the label and the 16-pixel pan-down arrow, or 5 around the arrow alone (_buttons.scss); raised outside a toolbar (adwButton) and flat inside one, as %toolbar makes button.arrow-button and button.image-button (_toolbars.scss of libadwaita 1.10.0).
func styleMenuButton(c *ui.Context, b ui.Element, label string, tb *idiomToolbar, i int) ui.Element {
	t := c.Theme()
	var arrow func()
	switch PlatformOf(c) {
	case Windows:
		if tb != nil {
			idiomFluentAppBarButton(c, b, false, false)
			b.Gap(5)
			color := t.TextMuted
			if b.Pressed() {
				color = t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.447), 0.529))
			}
			arrow = func() { idiomFluentSubItemChevron(c, color) }
		} else {
			fluentButton(c, b, false)
			arrow = func() {
				ui.Box(c).Size(12, 12).Shrink(0).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
					fluentChevronDown(p, r.X+r.W/2, r.Y+r.H/2, t.TextMuted)
				})
			}
		}
	case Linux:
		if tb != nil {
			b.MinHeight(34).MinWidth(34).Radius(9).FontWeight(700).Justify(ui.Center).AlignItems(ui.Center).Border(0, ui.Color{})
			idiomAdwFlatFace(c, b, false)
		} else {
			adwButton(c, b, false)
		}
		if label == "" {
			b.Padding(5)
		} else {
			b.Padding(5, 9).Gap(4)
		}
		color := adwTextOn(t, t.Background)
		arrow = func() { idiomAdwArrow(c, color) }
	default:
		if tb != nil {
			idiomMacToolbarItem(c, b, tb, i, false)
			b.PaddingStart(12).PaddingEnd(13).Gap(4)
			color := idiomMacInk(t, false)
			arrow = func() { idiomMacChevron(c, color) }
		} else {
			macButton(c, b, false)
			b.Gap(4)
			face := macFace(t.Background, t.Dark)
			if b.Pressed() {
				face = macPressed(t.Background, t.Dark)
			}
			color := tint(face, t.Dark, 0.847)
			arrow = func() { idiomMacChevron(c, color) }
		}
	}
	b.Children(func() {
		if label != "" {
			ui.Text(c, label).SingleLine()
		}
		arrow()
	})
	return b
}

// MenuButton creates a button, drawn as the platform's pull-down button, opening the system's own menu below it, which build fills; inside a Toolbar it is drawn as the platform's toolbar menu item. It is ui.MenuButton's counterpart, built as it is of ui.ButtonBase with Element.Menu: the menu is native on every platform.
func MenuButton(c *ui.Context, label string, build func(m *ui.Menu)) ui.Element {
	return idiomMenuButton(c, label, build, false)
}

// PopoverMenuButton creates a button, drawn as the platform's pull-down button (inside a Toolbar as the platform's toolbar menu item), opening a menu that MyGo draws in the window in the platform's menu colors and corners: NSMenu's on macOS, MenuFlyout's on Windows, popover.menu's on GNOME. Its rows' sizes are MyGo's. It is ui.PopoverMenuButton's counterpart, built as it is of ui.ButtonBase with Element.PopoverMenuButton.
func PopoverMenuButton(c *ui.Context, label string, build func(m *ui.Menu)) ui.Element {
	return idiomMenuButton(c, label, build, true)
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

type toggleGroup struct {
	depth, built int
	// face is the fill of the toggle built last, over which GNOME's next linked toggle lays its own where the two overlap.
	face ui.Color
}

// adwToggleGroupFace returns the fill of the toggle c's ToggleGroup built before the one being built, and records face as the fill of the one being built.
func adwToggleGroupFace(c *ui.Context, face ui.Color) (prev ui.Color) {
	v, ok := toggleGroups.Load(c)
	if !ok {
		return face
	}
	g := v.(*toggleGroup)
	prev, g.face = g.face, face
	return prev
}

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

// Field creates a labeled control: the label names the control fn builds for assistive technology, and a click on it focuses the control. In a Form the labels line up in a column before the controls, as macOS lays out a form; elsewhere the label is above the control. On Windows it is 8 DIPs above it and of regular weight, as a TextBox's header is (TextBox_themeresources.xaml: HeaderContentPresenter FontWeight Normal, TextBoxTopHeaderMargin 0,0,0,8): MyGo's medium label is left transparent, keeping its place, its click and the name it gives, and the regular one drawn over it in the theme's text color, which does not dim with a control disabled on its own, as MyGo's label does (a Field disabled as a whole dims). It is ui.Field's counterpart.
func Field(c *ui.Context, label string, fn func()) ui.Element {
	if PlatformOf(c) != Windows {
		return ui.Field(c, label, fn)
	}
	t := c.Theme()
	// The column repeats the gap MyGo puts between the control and the texts below it, so they stay where they were, and gives the control back the text color the label gives up.
	f := ui.Field(c, label, func() { ui.Column(c).Gap(t.Space(1.5)).TextColor(t.Text).Children(fn) })
	rtl := c.Direction() == ui.RightToLeft
	f.Gap(8).TextColor(ui.Transparent).DrawOver(func(p *ui.Painter, r ui.Rect) {
		formsLabelOver(p, r, label, 400, t.FontSize, t.Text, rtl)
	})
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
