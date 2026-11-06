package breeze

import (
	"sync"

	"github.com/egoist/mygo/ui"
)

// ---- preferences: groups with a header and a footer, rows with a prefix, rows holding content below their title, and radio rows ----
//
// The pieces a settings page builds around ActionRow and SwitchRow (widgets_rows.go), on their row look (rowsJoin, rowsStyle), so that they join a Card, a Fieldset, an AccordionItem or a PreferencesGroup as the platform joins rows there. macOS: SwiftUI's grouped Form with a Section whose header holds a Button after its title and whose footer is a Text, rows whose label is a Label with an image, a VStack row of a title over a TextEditor and a Picker of the radio-group style, rendered at 2x in Aqua and Dark Aqua on macOS 27.2 by tools/appkit-reference/Sources/appkit-reference/Scenes/preferences.swift (preferencesSet). Windows 11: CommunityToolkit/Windows's SettingsCard.xaml and SettingsCard.cs at main a784bad (SettingsCardHeaderIconMaxSize 20, SettingsCardHeaderIconMargin 2,0,20,0, SettingsCardWrapNoIconThreshold 286, the Vertical ContentAlignment state and SettingsCardVerticalHeaderContentSpacing 8), and the WinUI Gallery's settings page for a section's header; not rendered. GNOME: libadwaita 1.10.0's adw-preferences-group.ui and .c (the header box, its labels and the header suffix, single-line), _preferences.scss, adw-action-row.ui and _lists.scss (box.header > .prefixes); not rendered.

// PreferencesGroup creates a group of settings rows under a header, as the platform's settings pages group them: the title, at the end of the header's line the controls suffix builds (nil for none, such as a button that adds an item or checks again), and the description (an empty one shows none). Rows built in fn join as they join a Fieldset: one section on macOS, a boxed list on GNOME, SettingsCards 4 DIPs apart on Windows; other content sits in the section padded as a Fieldset pads it. Assistive technology sees a group named by the title and described by the description, as AdwPreferencesGroup is labelled by its title and described by its description.
//   - macOS: a SwiftUI Section of a grouped Form with a header and a footer: the header the title in a Fieldset's legend (13-point semibold, 10 points in from the section's edge) with the suffix at the end of its line, 10 points before the section's end, the title and the suffix centered on each other and the section 10 points below them (a 24-point push button makes the header line 24 points high); the description below the section as its footer, 11 points in secondaryLabelColor, 10 points in and 10 below it (measured on macOS 27.2).
//   - Windows: the WinUI Gallery's settings section: the title in BodyStrong (semibold), the suffix right-aligned on its line, the cards 10 DIPs below (as Fieldset); the description in the caption size and TextFillColorSecondary 10 DIPs under the cards (inferred: SettingsCard has no footer).
//   - GNOME: AdwPreferencesGroup: the header box holds the labels, the bold title over the description (body, dimmed) 6 pixels apart, and the suffix at its end, centered; a header without a description is single-line, at least 34 pixels high with its labels centered in it, and the list 6 pixels below it; with a description, 12 below it (the header's margin-bottom of 6 and the group's border-spacing of 6).
func PreferencesGroup(c *ui.Context, title, description string, suffix, fn func()) ui.Element {
	p := PlatformOf(c)
	t := c.Theme()
	g := ui.Column(c).Shrink(0).Role(ui.RoleGroup).Label(title)
	if description != "" {
		g.Description(description)
	}
	footer := description != "" && p != Linux
	header := title != "" || suffix != nil || (description != "" && p == Linux)
	g.Children(func() {
		if header {
			prefsHeader(c, p, title, description, suffix)
		}
		box := ui.Column(c).Gap(12)
		switch p {
		case Linux:
			// A boxed list, or a .card padded as Fieldset's for other content.
			fill, _ := adwBoxedListCard(t)
			box.Padding(12).Radius(12).Background(fill).TextColor(t.Text)
			adwBoxedListShadow(box)
		case Mac:
			box.Padding(10).Radius(12).Background(macFormSectionFill(t))
		}
		rowsFinish(c, box, rowsCollect(c, rowsInFieldset, func() { box.Children(fn) }))
		if footer {
			prefsFooter(c, p, description)
		}
	})
	switch p {
	case Linux:
		// preferencesgroup > box: border-spacing 6px; > box.header:not(.single-line): margin-bottom 6px.
		g.Gap(pick(description != "", float32(6), 12))
	default:
		g.Gap(10)
	}
	return g
}

// prefsHeader builds a PreferencesGroup's header as PreferencesGroup describes: the title, the description under it on GNOME, and the suffix at the end of its line.
func prefsHeader(c *ui.Context, p Platform, title, description string, suffix func()) {
	t := c.Theme()
	h := ui.Row(c).Shrink(0).AlignItems(ui.Center).Gap(8)
	switch p {
	case Mac:
		// SwiftUI insets a section's header 10 points from the section's edges, the suffix's end too (measured).
		h.PaddingX(10)
	case Linux:
		// box.single-line: min-height 34px, where only the title and the suffix show.
		if description == "" {
			h.MinHeight(adwGroupHeader)
		}
		h.Gap(0)
	}
	h.Children(func() {
		labels := ui.Column(c).Grow(1).Shrink(1).MinWidth(0)
		labels.Children(func() {
			if title != "" {
				// The legend's weight: semibold on macOS and Windows, bold on GNOME (.heading).
				ui.Text(c, title).FontWeight(pick(p == Linux, 600, 700)).TextColor(t.Text).Role(ui.RoleNone)
			}
			if p == Linux && description != "" {
				// box.labels: border-spacing 6px; the description .body .dimmed, window_fg_color at 55%.
				labels.Gap(6)
				ui.Text(c, description).TextColor(pick(t.Dark, t.Background.Mix(adwInk, 0.8*0.55), t.Background.Mix(white, 0.55))).Role(ui.RoleNone)
			}
		})
		if suffix != nil {
			ui.Row(c).Shrink(0).AlignItems(ui.Center).Gap(8).Children(func() { rowsIsolate(c, suffix) })
		}
	})
}

// prefsFooter builds a PreferencesGroup's description below its rows on macOS and Windows, as PreferencesGroup describes.
func prefsFooter(c *ui.Context, p Platform, description string) {
	t := c.Theme()
	text := ui.Text(c, description).Role(ui.RoleNone)
	if p == Windows {
		// CaptionTextBlockStyle, 12 where the body is 14, in TextFillColorSecondary.
		text.FontSize(t.FontSize * 12 / 14).TextColor(t.TextMuted)
		return
	}
	// secondaryLabelColor over the window, as AppKit composites it: 128 over 255 in Aqua (measured on macOS 27.2), white at 55% in Dark Aqua as a row's subtitle.
	text.FontSize(t.FontSize * 11 / 13).TextColor(pick(t.Dark, t.Background.Mix(black, 0.498), t.Background.Mix(white, 0.549))).PaddingX(10)
}

// PrefixRow creates an ActionRow with what prefix builds before its title (nil for nothing), such as an icon, an avatar or an app's tile, and at its end the controls suffix builds (nil for none). It joins a Card, a Fieldset, an AccordionItem or a PreferencesGroup as ActionRow does, its separators where the platform puts them between rows with and without a prefix. It is a group that assistive technology names by title and describes by subtitle, as ActionRow is; what prefix builds keeps its own roles.
//   - macOS: a row of a grouped Form whose label holds an image (a Label with an icon): the prefix 10 points in from the section's edge, the title 10 after it, the prefix and the labels centered on each other and padded 10 above and below (a row of a 30-point image and one line is 50 points high), the separators still inset 10 points from the section's edges, under the prefix too (measured on macOS 27.2).
//   - Windows: SettingsCard's HeaderIcon: a column 2 DIPs in from the card's padding and 20 before the header, at least 20 DIPs square, centered on the card; narrower than 286 DIPs the card drops it (RightWrappedNoIcon). WinUI scales the HeaderIcon into a 20-DIP Viewbox; Breeze UI lays the prefix at its own size, so build it 20 DIPs.
//   - GNOME: AdwActionRow's prefixes: the prefix after the header's 12-pixel inset, centered, 12 pixels before the title (the prefixes' margin of 6 and the header's border-spacing of 6).
func PrefixRow(c *ui.Context, title, subtitle string, prefix, suffix func()) ui.Element {
	row := ui.Row(c).Shrink(0).Role(ui.RoleGroup).Label(title)
	if subtitle != "" {
		row.Description(subtitle)
	}
	item, in, first := rowsJoin(c, row)
	look := rowsLook{collapse: true}
	if prefix != nil {
		look.prefix = func(rowsInk) { prefix() }
	}
	var trailing func(under ui.Color, hovered, pressed bool)
	if suffix != nil {
		trailing = func(ui.Color, bool, bool) { rowsIsolate(c, suffix) }
	}
	rowsStyle(c, row, item, in, first, title, subtitle, look, trailing)
	return row
}

// ContentRow creates a settings row holding what fn builds below its label, such as a TextArea or a list of models, as the platforms lay a control too wide for the end of a row under its title. It joins a Card, a Fieldset, an AccordionItem or a PreferencesGroup as ActionRow does, and is a group that assistive technology names by label.
//   - macOS: a VStack row of a grouped Form: the label in labelColor as a row's title, 13.5 points from the row's top to its ink, then 9.5 points below its line the content, filling the row's width inside its inset of 10, and 10 points below it (measured on macOS 27.2 from a VStack of a Text and a TextEditor).
//   - Windows: SettingsCard with ContentAlignment Vertical: the header, then 8 DIPs below it (SettingsCardVerticalHeaderContentSpacing) the content, stretched across the card.
//   - GNOME: a row of a boxed list built of the action row's header: the label as its title, the content 6 pixels below it (the header's border-spacing) inside the title box's margins of 6 (inferred: libadwaita has no such row).
func ContentRow(c *ui.Context, label string, fn func()) ui.Element {
	p := PlatformOf(c)
	row := ui.Row(c).Shrink(0).Role(ui.RoleGroup).Label(label)
	item, in, first := rowsJoin(c, row)
	look := rowsLook{content: func(ink rowsInk) {
		var col ui.Element
		if p == Windows {
			// What a SettingsCard holds takes its translucent colors over the card's fill, as a row's trailing controls do.
			col = ui.Themed(c, fluentOver(c.Theme(), ink.under)).Background(ui.Transparent).Column()
		} else {
			col = ui.Column(c)
		}
		col.Grow(1).Shrink(1).MinWidth(0).AlignItems(ui.Stretch)
		switch p {
		case Mac:
			col.Gap(9.5)
		case Windows:
			col.Gap(8)
		default:
			col.Padding(6, 0).Gap(6)
		}
		col.Children(func() {
			if label != "" {
				ui.Text(c, label).TextColor(ink.title).Role(ui.RoleNone)
			}
			if fn != nil {
				rowsIsolate(c, fn)
			}
		})
	}}
	rowsStyle(c, row, item, in, first, label, "", look, nil)
	if p == Mac {
		// The label as far below the row's top as a row's title, the content 10 points above the row's end; in a DisclosureGroup 6 less, as rowsStyle takes it there (inferred).
		top, bottom := float32(10.5), float32(10)
		switch {
		case in == rowsInAccordion:
			top, bottom = 4.5, 4
		case in != rowsAlone && !first:
			top = 11.5
		}
		row.Padding(top, 10, bottom, 10)
	}
	return row
}

// RadioRow creates a settings row that chooses value into *selected, one of the RadioRows sharing selected: its title, the subtitle below it (an empty subtitle shows none), the platform's radio button, and the controls suffix builds at its end (nil for none). A click anywhere on the row, Space or Enter choose it, and Changed reports a new choice; while it has the focus, Up and Down choose the RadioRow before or after it among those of the same Card, Fieldset, AccordionItem or PreferencesGroup sharing selected, round from the last to the first, passing over disabled ones, and move the focus there. It joins a group as ActionRow does. Assistive technology sees one radio button named by the title, described by the subtitle and checked while chosen, the radio button inside it not a stop of its own.
//   - macOS: a row of a SwiftUI Picker of the radio-group style in a grouped Form, as macOS 27.2 draws one: every option a row of its own between the section's separators, its radio button (Radio's 16-point circle, its trough laid over the section: 223 over 247 in Aqua) at the row's trailing edge, centered on the title's line, after the suffix; AppKit's focus ring around the circle (measured on macOS 27.2).
//   - Windows: a clickable SettingsCard (IsClickEnabled) holding a RadioButton as its HeaderIcon: the 20-DIP circle 2 DIPs in and 20 before the header, the card lit under the pointer and pressed as SwitchRow's, the circle lit only under the pointer over it, the card's focus visual (inferred: WinUI has no radio card). Breeze UI keeps the circle where SettingsCard would collapse a HeaderIcon in a card narrower than 286 DIPs.
//   - GNOME: an activatable AdwActionRow with a GtkCheckButton in group mode as its prefix (the documented radio row): the 20-pixel circle 12 pixels before the title, the row lit at 3% under the pointer and 8% pressed, the circle only under the pointer over it, the row's focus ring across its edge.
func RadioRow[T comparable](c *ui.Context, selected *T, value T, title, subtitle string, suffix func()) ui.Element {
	p := PlatformOf(c)
	row := ui.RadioBase(c, selected, value).Shrink(0).Label(title).FocusRing(false)
	if subtitle != "" {
		row.Description(subtitle)
	}
	group := rowsGroupOf(c)
	item, in, first := rowsJoin(c, row)
	key := radioKey{c: c, selected: selected}
	if v, ok := radioFocus.Load(key); ok && v.(T) == value {
		// An arrow on the RadioRow before chose this one: it takes the focus.
		radioFocus.Delete(key)
		row.Focus().SetChanged()
	}
	if group != nil {
		radioRowKeys(group, row, key, selected, value)
	}
	on := *selected == value
	pressed, focused := row.Pressed(), row.FocusVisible()
	look := rowsLook{lit: true, hovered: row.Hovered(), pressed: pressed, focused: focused}
	disabled := disabledOnce(row)
	dot := func(under ui.Color) {
		d := ui.Box(c).Shrink(0).Role(ui.RoleNone)
		switch p {
		case Windows:
			over := d.Hovered()
			fluentRadioDot(c, d, on, over, over && pressed)
		case Linux:
			// GTK lights the check button only under the pointer over it; a press on the row does not press it.
			adwRadioDot(c, d, on, d.Hovered(), false)
		default:
			macRadioDot(c, d, under, on, pressed, focused, disabled)
		}
	}
	var trailing func(under ui.Color, hovered, pressed bool)
	if p == Mac {
		trailing = func(under ui.Color, _, _ bool) {
			if suffix != nil {
				rowsIsolate(c, suffix)
			}
			dot(under)
		}
	} else {
		look.prefix = func(ink rowsInk) { dot(ink.under) }
		if suffix != nil {
			trailing = func(ui.Color, bool, bool) { rowsIsolate(c, suffix) }
		}
	}
	rowsStyle(c, row, item, in, first, title, subtitle, look, trailing)
	return row
}

// radioKey names the RadioRows of a window that choose into the same value.
type radioKey struct {
	c        *ui.Context
	selected any
}

// radioFocus holds, for a radioKey, the value whose RadioRow takes the focus as it is built next, chosen by an arrow on another.
var radioFocus sync.Map // radioKey -> T

// radioRun is the RadioRows of a group that choose into the same value, in the order they are built.
type radioRun[T comparable] struct {
	rows []*radioEntry[T]
}

// radioEntry is a RadioRow of a run: its value and whether it is disabled, which its row learns as the frame is painted.
type radioEntry[T comparable] struct {
	value    T
	disabled func() bool
}

// rowsGroupOf returns the group of rows being built in c, nil for none.
func rowsGroupOf(c *ui.Context) *rowsGroup {
	if v, ok := rowGroups.Load(c); ok {
		return v.(*rowsGroup)
	}
	return nil
}

// radioRowKeys adds row, choosing value into *selected, to its run in group g, and has Up and Down on it choose the RadioRow before or after it in the run that is not disabled, round from the last to the first, and give that row the focus as it is built next.
func radioRowKeys[T comparable](g *rowsGroup, row ui.Element, key radioKey, selected *T, value T) {
	if g.radios == nil {
		g.radios = map[any]any{}
	}
	run, ok := g.radios[key.selected].(*radioRun[T])
	if !ok {
		run = &radioRun[T]{}
		g.radios[key.selected] = run
	}
	at := len(run.rows)
	entry := &radioEntry[T]{value: value, disabled: disabledOnce(row)}
	run.rows = append(run.rows, entry)
	// Disabled may come after RadioRow returns: the row asks itself as it is painted, before any key comes.
	row.Draw(func(*ui.Painter, ui.Rect) { entry.disabled() })
	row.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind != ui.InputKeyDown || ev.Mods != 0 {
			return false
		}
		step := 0
		switch ev.Key {
		case ui.KeyDown:
			step = 1
		case ui.KeyUp:
			step = -1
		default:
			return false
		}
		n := len(run.rows)
		for k := 1; k < n; k++ {
			to := run.rows[((at+step*k)%n+n)%n]
			if !to.disabled() {
				*selected = to.value
				radioFocus.Store(key, to.value)
				return true
			}
		}
		return false
	})
}
