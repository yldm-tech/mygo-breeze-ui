package breeze

import (
	"math"
	"sync"

	"github.com/egoist/mygo/ui"
)

// ---- layout: disclosures, accordions, fieldsets, forms, check box groups and splits ----
//
// macOS: AppKit's disclosure button (NSButton of bezel style disclosure), NSBox and NSSplitView, and SwiftUI's grouped Form with DisclosureGroup rows, rendered at 2x in Aqua and Dark Aqua on macOS 27 by a scratch script modeled on tools/appkit-reference/reference.swift. Windows 11: Expander_themeresources.xaml and Expander.xaml of microsoft/microsoft-ui-xaml, the card colors of Common_themeresources_any.xaml, the spacing of Windows's "Content layout and spacing" guidance and the settings page of the WinUI Gallery. GNOME: libadwaita 1.7.6's _expanders.scss, _lists.scss, _misc.scss (.card), _preferences.scss and _defaults.scss, adwaita-icon-theme 48's pan-end-symbolic and GTK 4.18's gtkexpander.c.

// Shares of black (Aqua) or white (Dark Aqua) over the window background, measured from AppKit's NSBox and SwiftUI's grouped Form.
const (
	macSectionFillAlpha      = 0.031 // an NSBox's fill and a grouped form section's, the same in both: 247 over 255, 37 over 30
	macSectionSeparatorAlpha = 0.078 // the line between two rows of a grouped section: 235 over 255, 47 over 30
	macDisclosureLightAlpha  = 0.25  // a disclosure button's chevron in Aqua: 191 over 255
	macDisclosureDarkAlpha   = 0.436 // and in Dark Aqua: 128 over 30
)

// macFormSectionFill is the fill of an NSBox and of a section of a grouped form: corners of 12, no border.
func macFormSectionFill(t *ui.Theme) ui.Color { return tint(t.Background, t.Dark, macSectionFillAlpha) }

// fluentExpanderCard returns CardBackgroundFillColorDefault, CardBackgroundFillColorSecondary and CardStrokeColorDefault over the window: white at 70.2%, #f6f6f6 at 50.2% and black at 5.9% in the light theme; white at 5.1%, white at 3.1% and black at 9.8% in the dark one.
func fluentExpanderCard(t *ui.Theme) (fill, secondary, stroke ui.Color) {
	if t.Dark {
		return t.Background.Mix(white, 0.051), t.Background.Mix(white, 0.031), t.Background.Mix(black, 0.098)
	}
	return t.Background.Mix(white, 0.702), t.Background.Mix(ui.Hex("#f6f6f6"), 0.502), t.Background.Mix(black, 0.059)
}

// adwBoxedListCard returns card_bg_color over the window, #ffffff in the light theme and white at 8% in the dark one, and card_shade_color over it, the shade at 7% and 36%: the separator of a boxed list's rows.
func adwBoxedListCard(t *ui.Theme) (fill, shade ui.Color) {
	if t.Dark {
		fill = t.Background.Mix(white, 0.08)
		return fill, fill.Mix(adwInk, 0.36)
	}
	return white, white.Mix(adwInk, 0.07)
}

// adwBoxedListShadow casts .card's box-shadow: a 1-pixel ring of the shade at 3%, and two soft shadows at 7% and 3%.
func adwBoxedListShadow(e ui.Element) {
	e.Shadow(0, 0, 0, 1, ui.RGBA(0, 0, 6, 0.03)).Shadow(0, 1, 3, 1, ui.RGBA(0, 0, 6, 0.07)).Shadow(0, 2, 6, 2, ui.RGBA(0, 0, 6, 0.03))
}

// strokeDisclosureChevron strokes a chevron pointing right with its tip tip past cx, cy, turned about cx, cy by deg degrees clockwise; right to left it points left and turns the other way, as disclosures mirror. Its arms reach back from the tip by back and out to either side by half.
func strokeDisclosureChevron(p *ui.Painter, cx, cy, deg float32, rtl bool, tip, back, half, width float32, col ui.Color) {
	s := float32(1)
	if rtl {
		s, deg = -1, -deg
	}
	sin, cos := math.Sincos(float64(deg) * math.Pi / 180)
	at := func(x, y float32) (float32, float32) {
		return cx + x*float32(cos) - y*float32(sin), cy + x*float32(sin) + y*float32(cos)
	}
	var path ui.Path
	path.MoveTo(at(s*(tip-back), -half)).LineTo(at(s*tip, 0)).LineTo(at(s*(tip-back), half))
	p.StrokePath(&path, width, col)
}

// Collapsible creates a disclosure showing label, which a click opens and closes, as do Enter and Space, showing below it what fn builds while *open, as the platform draws a disclosure button beside its label (macOS), an Expander (Windows) or a GtkExpander (GNOME). The arrow turns and the content grows into view as it opens; Changed reports a click. It is ui.Collapsible's counterpart, built on ui.CollapsibleBase.
func Collapsible(c *ui.Context, label string, open *bool, fn func()) ui.Element {
	p := PlatformOf(c)
	root := ui.Column(c).Shrink(0)
	root.Children(func() {
		parts := ui.CollapsibleBase(c, open)
		if parts.Trigger.Changed() {
			root.SetChanged()
		}
		progress := parts.Progress()
		switch p {
		case Windows:
			fluentExpanderHeader(c, root, parts.Trigger, label, progress)
		case Linux:
			adwExpanderTitle(c, parts.Trigger, label, progress)
		default:
			macDisclosure(c, parts.Trigger, label, progress)
		}
		if !*open && progress == 0 {
			return
		}
		panel := parts.Panel(fn)
		switch p {
		case Windows:
			fluentExpanderContent(c, panel)
		case Linux:
			// GtkExpander puts its child right below the title, not indented.
			panel.Gap(6) // inferred: GtkExpander holds one child
		default:
			// Under the label, past the 13-point box and the 1.5 after it (inferred: SwiftUI's DisclosureGroup centers its content in the width of its stack rather than indenting it), 8 apart as in its VStack.
			panel.PaddingStart(14.5).Gap(8)
		}
	})
	return root
}

// macDisclosure draws an NSButton of bezel style disclosure beside its label: a 13-point box holding a chevron 7 points tall and 3.5 deep along its stroke's middle, stroked 1.5 points with round ends in black at 25% (white at 43.6% in Dark Aqua), pointing right and turning down as it opens, with the label 1.5 points after the box in a row 24 points high. Measured from AppKit's button and SwiftUI's DisclosureGroup rendered at 2x on macOS 27: the closed chevron's ink sits 0.5 points right of the box's center and 0.25 above it, the open one's 0.5 left and 0.25 below, which a turn about a point 0.25 left of and above the center gives within a quarter point; the label's ink starts 5.5 points after the chevron's, and the rows are 32 points apart in a VStack spaced 8. AppKit's focus ring (macFocusRing) goes around the box, as AppKit draws it around the button (inferred from the push button's, measured live).
func macDisclosure(c *ui.Context, tr ui.Element, label string, progress float32) {
	t := c.Theme()
	ink := pick(t.Dark, t.Background.Mix(black, macDisclosureLightAlpha), t.Background.Mix(white, macDisclosureDarkAlpha))
	rtl := c.Direction() == ui.RightToLeft
	tr.AlignSelf(ui.Start).AlignItems(ui.Center).Gap(1.5).MinHeight(24).FocusRing(false)
	shift := float32(0.25)
	if rtl {
		shift = -shift
	}
	tr.Children(func() {
		ui.Box(c).Size(13, 13).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
			strokeDisclosureChevron(p, r.X+r.W/2-shift, r.Y+r.H/2-0.25, 90*progress, rtl, 2.5, 3.5, 3.5, 1.5, ink)
			if tr.FocusVisible() {
				macFocusRing(p, r, 3)
			}
		})
		ui.Text(c, label).SingleLine()
	})
}

// fluentExpanderHeader draws an Expander's header: a card at least 48 DIPs high (ExpanderMinHeight) with CardBackgroundFillColorDefault inside a 1-DIP CardStrokeColorDefault border, corners of 4 (ControlCornerRadius), the bottom ones square while it is open; the header text 16 DIPs in (ExpanderHeaderPadding), and the chevron button at its end, 32 DIPs (ExpanderChevronButtonSize) with a margin of 20,0,8,0 (ExpanderChevronMargin), SubtleFillColorSecondary under the pointer and Tertiary pressed, around the 12-DIP chevron glyph (E70D, E70E), which turns from down to up.
func fluentExpanderHeader(c *ui.Context, root, tr ui.Element, label string, progress float32) {
	t := c.Theme()
	fill, _, stroke := fluentExpanderCard(t)
	bottom := float32(4)
	if progress > 0 {
		bottom = 0
	}
	root.AlignItems(ui.Stretch)
	// The padding on logical sides, as an Expander mirrors right to left with its FlowDirection.
	tr.MinHeight(48).Padding(0).PaddingStart(16).PaddingEnd(8).AlignItems(ui.Center).Radius(4, 4, bottom, bottom).Background(fill).Border(1, stroke)
	hover, pressed := fluentSubtle(t)
	chevron := ui.Transparent
	switch {
	case tr.Pressed():
		chevron = pressed
	case tr.Hovered():
		chevron = hover
	}
	tr.Children(func() {
		ui.Text(c, label).Grow(1)
		ui.Box(c).Size(32, 32).Shrink(0).MarginStart(20).Radius(4).Background(chevron).Draw(func(p *ui.Painter, r ui.Rect) {
			// Segoe Fluent Icons' chevron in a 12-DIP glyph: 10 wide, 5 deep, a 1-DIP stroke (inferred from the glyph's outline).
			strokeDisclosureChevron(p, r.X+r.W/2, r.Y+r.H/2, 90+180*progress, false, 2.5, 5, 5, 1, t.Text)
		})
	})
}

// fluentExpanderContent draws an Expander's content: ExpanderContentPadding of 16, CardBackgroundFillColorSecondary, a CardStrokeColorDefault border on every side but the header's (ExpanderContentDownBorderThickness 1,0,1,1) and the bottom corners rounded by 4. An Expander holds one content element; 12 between several is inferred from Windows's spacing guidance for content areas.
func fluentExpanderContent(c *ui.Context, panel ui.Element) {
	_, secondary, stroke := fluentExpanderCard(c.Theme())
	panel.Padding(16).Gap(12).Background(secondary).BorderWidth(0, 1, 1, 1).BorderColor(stroke).Radius(0, 0, 4, 4)
}

// adwExpanderTitle draws a GtkExpander's title: the 16-pixel expander arrow, pan-end-symbolic (a chevron stroked 2 pixels, its tip 3 pixels right of the center), turning to pan-down-symbolic as it opens, and the label right after it; the arrow is the text color at 70%, opaque under the pointer, and the title's corners are 6 for its focus ring.
func adwExpanderTitle(c *ui.Context, tr ui.Element, label string, progress float32) {
	t := c.Theme()
	share := float32(0.7)
	if tr.Hovered() || tr.Pressed() {
		share = 1
	}
	ink := t.Background.Mix(t.Text, share)
	rtl := c.Direction() == ui.RightToLeft
	tr.AlignSelf(ui.Start).AlignItems(ui.Center).MinHeight(16).Radius(6)
	tr.Children(func() {
		ui.Box(c).Size(16, 16).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
			strokeDisclosureChevron(p, r.X+r.W/2, r.Y+r.H/2, 90*progress, rtl, 3, 4, 4, 2, ink)
		})
		ui.Text(c, label).SingleLine()
	})
}

// accordions holds the windows building an Accordion, with how many of its items they have built, so that an item knows whether it is the first, which has no separator above it.
var accordions sync.Map // *ui.Context -> *layoutAccordion

type layoutAccordion struct{ built int }

// Accordion creates sections, one above the other, that fn builds with AccordionItem, each opening and closing on its own; Up and Down move the focus between their headers, as do Home and End. The box around them is the platform's: a section of a grouped form on macOS (as NSBox and System Settings draw it), Expander cards 4 DIPs apart on Windows, a boxed list on GNOME; the headers, their arrows and the keys are MyGo's. It is ui.Accordion's counterpart.
func Accordion(c *ui.Context, fn func()) ui.Element {
	p := PlatformOf(c)
	outer, nested := accordions.Load(c)
	accordions.Store(c, &layoutAccordion{})
	a := ui.Accordion(c, fn)
	if nested {
		accordions.Store(c, outer)
	} else {
		accordions.Delete(c)
	}
	t := c.Theme()
	switch p {
	case Windows:
		// Expanders stacked as the cards of a settings page, SettingsCardSpacing apart (WinUI Gallery's SettingsPage.xaml).
		a.Radius(0).Border(0, ui.Color{}).Background(ui.Transparent).Gap(4)
	case Linux:
		// list.boxed-list: a .card with corners of 12 ($card_radius) and its shadow.
		fill, _ := adwBoxedListCard(t)
		a.Radius(12).Border(0, ui.Color{}).Background(fill)
		adwBoxedListShadow(a)
	default:
		a.Radius(12).Border(0, ui.Color{}).Background(macFormSectionFill(t))
	}
	return a
}

// AccordionItem creates a section of an Accordion: a header showing title, which a click opens and closes, as do Enter and Space, and below it what fn builds while *open. Changed reports a click. In an Accordion it is a row of the platform's box: below a separator inset 10 points on macOS, an Expander's card on Windows, below the separator of a boxed list's rows on GNOME. Its header is MyGo's. It is ui.AccordionItem's counterpart.
func AccordionItem(c *ui.Context, title string, open *bool, fn func()) ui.Element {
	first, in := true, false
	if v, ok := accordions.Load(c); ok {
		// Counted before fn builds, as items of an accordion inside this one count in their own.
		g := v.(*layoutAccordion)
		g.built++
		first, in = g.built == 1, true
	}
	p := PlatformOf(c)
	item := ui.AccordionItem(c, title, open, fn)
	if !in {
		return item
	}
	t := c.Theme()
	switch p {
	case Windows:
		fill, _, stroke := fluentExpanderCard(t)
		item.Radius(4).Background(fill).Border(1, stroke).Clip()
	case Linux:
		if !first {
			// %boxed_list_row: a 1-pixel border of card_shade_color between rows.
			_, shade := adwBoxedListCard(t)
			item.BorderColor(shade)
		}
	default:
		item.Border(0, ui.Color{})
		if !first {
			line := tint(t.Background, t.Dark, macSectionSeparatorAlpha)
			item.DrawOver(func(p *ui.Painter, r ui.Rect) {
				// SwiftUI's grouped form insets its separators 10 points from the section's edges.
				p.Fill(ui.Rect{X: r.X + 10, Y: r.Y, W: r.W - 20, H: 1}, line, 0)
			})
		}
	}
	return item
}

// Fieldset creates a group of the fields that fn builds under legend, which names the group for assistive technology; Disabled disables them all, and in a Form their labels line up with those of the other fields. The fields sit in the platform's group: on macOS a section of a grouped form (NSBox's fill, corners of 12) under the legend, inset 10 points as SwiftUI insets a section's header; on Windows, which has no group box, under the legend as a settings page puts its cards under a Body Strong header, 10 DIPs below it; on GNOME in a card under the legend, as an AdwPreferencesGroup puts its boxed list 12 pixels below its title. In a Form the section or the card reaches past the fields by its padding on either side, so that the fields line up with the form's others. The legend's text is MyGo's, semibold, where GNOME's .heading is bold. It is ui.Fieldset's counterpart.
func Fieldset(c *ui.Context, legend string, fn func()) ui.Element {
	p := PlatformOf(c)
	t := c.Theme()
	_, inForm := layoutForms.Load(c)
	g := ui.Fieldset(c, legend, func() {
		// 12 between the fields, as MyGo spaces them; Windows's guidance puts 12 between content areas.
		box := ui.Column(c).Gap(12)
		pad := float32(0)
		switch p {
		case Linux:
			// A .card, padded as a row's header is inset 12 pixels (inferred for fields rather than rows).
			pad = 12
			fill, _ := adwBoxedListCard(t)
			box.Padding(pad).Radius(12).Background(fill)
			adwBoxedListShadow(box)
		case Mac:
			// The section reaches back under the legend's inset; its rows are inset 10 points (measured from SwiftUI's grouped form).
			pad = 10
			box.MarginStart(-pad).Padding(pad).Radius(12).Background(macFormSectionFill(t))
		}
		if inForm && pad > 0 {
			// The labels of a form are as wide as the widest, so its fields line up only where they start and end alike.
			box.MarginStart(-pad).MarginEnd(-pad)
		}
		box.Children(fn)
	})
	switch p {
	case Windows:
		// WinUI Gallery's SettingsPage.xaml: SettingsSectionHeaderTextBlockStyle (BodyStrongTextBlockStyle) has a bottom margin of 6, and SettingsCardSpacing puts 4 more before the first card.
		g.Gap(10)
	case Linux:
		g.Gap(12) // preferencesgroup > box: border-spacing 6px, and 6px below the header
	default:
		// The header's baseline is 13 points above the section: about 10 below its line. In a Form the section reaches back instead, keeping the legend 10 in from its edge.
		if !inForm {
			g.PaddingStart(10)
		}
		g.Gap(10)
	}
	return g
}

// layoutForms holds the windows building a Form on macOS or GNOME, with how many Forms each is in, so that a Fieldset knows its fields line up with a form's.
var layoutForms sync.Map // *ui.Context -> int

// Form creates a column of the fields that fn builds (Field, Fieldset). On macOS and GNOME it is ui.Form: the labels sit in a column of their own beside the controls, aligned to the end and as wide as the widest, as macOS's forms and GNOME's rows put a title before the control. On Windows, whose controls take their header above them, it is a column of fields 12 DIPs apart, each label 8 DIPs above its control, as Windows's spacing guidance has them. It is ui.Form's counterpart.
func Form(c *ui.Context, fn func()) ui.Element {
	if PlatformOf(c) == Windows {
		return ui.Column(c).Gap(12).Children(fn)
	}
	outer, _ := layoutForms.Load(c)
	depth, _ := outer.(int)
	layoutForms.Store(c, depth+1)
	f := ui.Form(c, fn)
	if depth > 0 {
		layoutForms.Store(c, depth)
	} else {
		layoutForms.Delete(c)
	}
	return f
}

// ---- widgets of package ui, in the platform's theme (Use) ----

// CheckboxGroup creates a check box over the check boxes fn builds, indented below it: it is checked when they all are, mixed when some are, and a click or Space checks them all, or none when they all are. Build the boxes with Checkbox to draw them as the platform does; the check box over them is MyGo's, a 16-DIP box in the platform's theme, smaller than the 20-DIP boxes of Windows and GNOME below it: package ui has no base for it. It is package ui's.
func CheckboxGroup(c *ui.Context, label string, fn func()) ui.Element {
	return ui.CheckboxGroup(c, label, fn)
}

// Split lays out first and second side by side with a divider between them that a drag or the arrows move: *size is the width of first. The divider is MyGo's, a 1-DIP line in the theme's Border color, as NSSplitView's thin divider and GtkPaned's handle are 1 wide: separatorColor on macOS (AppKit's thin divider measures black at 10.2% in Aqua, which it matches, and white at 12% in Dark Aqua, a little lighter than separatorColor's 9.8%), $border_color on GNOME as GtkPaned's separator, ControlStrokeColorDefault on Windows, which has no splitter of its own. It turns to the accent under the pointer, as none of the three do. It is package ui's.
func Split(c *ui.Context, size *float32, first func(), second func()) ui.Element {
	return ui.Split(c, size, first, second)
}

// SplitVertical is Split with first above second, *size its height. It is package ui's.
func SplitVertical(c *ui.Context, size *float32, first func(), second func()) ui.Element {
	return ui.SplitVertical(c, size, first, second)
}

// SplitPercent is Split with *p the percentage of the room first takes, from 0 to 100. It is package ui's.
func SplitPercent(c *ui.Context, p *float32, first func(), second func()) ui.Element {
	return ui.SplitPercent(c, p, first, second)
}

// SplitVerticalPercent is SplitPercent with first above second. It is package ui's.
func SplitVerticalPercent(c *ui.Context, p *float32, first func(), second func()) ui.Element {
	return ui.SplitVerticalPercent(c, p, first, second)
}
