package breeze

import (
	"sync"
	"time"

	"github.com/egoist/mygo/ui"
)

// ---- rows: content cards, settings rows and switch rows ----
//
// macOS: SwiftUI's GroupBox and a grouped Form with a Section of rows (LabeledContent with a button, Toggles of style switch with and without a second Text, LabeledContent with a value), rendered at 2x in Aqua and Dark Aqua on macOS 27.2 by tools/appkit-reference/reference.swift (its rows scene). Windows 11: CommunityToolkit/Windows's SettingsCard.xaml and SettingsExpander.xaml (main at a784bad), the card colors of Common_themeresources_any.xaml, ToggleSwitch_Partial.cpp of microsoft/microsoft-ui-xaml for the On and Off content, and the WinUI Gallery's settings page for the space between cards. GNOME: libadwaita 1.10.0's _lists.scss (AdwActionRow, %boxed_list_row, list.boxed-list, row.expander's nested list), _misc.scss (.card), _labels.scss, _preferences.scss and _colors.scss, adw-action-row.ui and adw-switch-row.c.

// rowsIn tells where a row is built: alone, or in a Card, a Fieldset or an AccordionItem, whose rows join as the platform joins rows in a group.
type rowsIn uint8

const (
	rowsAlone rowsIn = iota
	rowsInCard
	rowsInFieldset
	rowsInAccordion
)

// rowsItem is a row of a group, with the corners its container gives it once the group is built: the ends of a boxed list are rounded.
type rowsItem struct {
	e      ui.Element
	radius [4]float32
}

// rowsGroup holds the rows built directly in a Card, a Fieldset or an AccordionItem.
type rowsGroup struct {
	in   rowsIn
	rows []*rowsItem
}

// rowGroups holds the group each window is building rows in.
var rowGroups sync.Map // *ui.Context -> *rowsGroup

// rowsCollect builds fn with the rows built in it counted in a group of kind in, which it returns; the group around it, if any, is the window's again after.
func rowsCollect(c *ui.Context, in rowsIn, fn func()) *rowsGroup {
	outer, nested := rowGroups.Load(c)
	g := &rowsGroup{in: in}
	rowGroups.Store(c, g)
	fn()
	if nested {
		rowGroups.Store(c, outer)
	} else {
		rowGroups.Delete(c)
	}
	return g
}

// rowsIsolate builds fn, the trailing controls of a row, outside the row's group, so that rows built in it are rows of their own.
func rowsIsolate(c *ui.Context, fn func()) {
	outer, nested := rowGroups.Load(c)
	rowGroups.Delete(c)
	fn()
	if nested {
		rowGroups.Store(c, outer)
	}
}

// rowsJoin adds row to the group being built, if any, and tells where it is and whether it is the group's first.
func rowsJoin(c *ui.Context, row ui.Element) (item *rowsItem, in rowsIn, first bool) {
	item = &rowsItem{e: row}
	v, ok := rowGroups.Load(c)
	if !ok {
		return item, rowsAlone, true
	}
	g := v.(*rowsGroup)
	g.rows = append(g.rows, item)
	return item, g.in, len(g.rows) == 1
}

// rowsFinish lays out box, a Card or a Fieldset's box, around the rows of g once they are built, as the platform joins rows: without padding, the rows touching, and the ends of the list rounded with the box (macOS's grouped section and GNOME's boxed list, corners of 12; Windows's card, whose corners of 4 leave 3 inside its border). In a Fieldset on Windows the rows are SettingsCards of their own, 4 DIPs apart, as a settings page stacks them (SettingsCardSpacing in the WinUI Gallery's SettingsPage.xaml).
func rowsFinish(c *ui.Context, box ui.Element, g *rowsGroup) {
	if len(g.rows) == 0 {
		return
	}
	p := PlatformOf(c)
	if p == Windows && g.in == rowsInFieldset {
		box.Gap(4)
		return
	}
	box.Padding(0).Gap(0)
	r := float32(12)
	if p == Windows {
		r = 3
	}
	first, last := g.rows[0], g.rows[len(g.rows)-1]
	first.radius[0], first.radius[1] = r, r
	last.radius[2], last.radius[3] = r, r
	for _, it := range g.rows {
		it.e.Radius(it.radius[0], it.radius[1], it.radius[2], it.radius[3])
	}
}

// Card creates a card holding what fn builds, as the platform draws a box around content that belongs together: on macOS a SwiftUI GroupBox (the fill of a grouped section, corners of 12, padded 5 points, its content 8 apart as in SwiftUI's VStack), on Windows the card of WinUI's Expander and CommunityToolkit's SettingsCard (CardBackgroundFillColorDefault inside a 1-DIP CardStrokeColorDefault border, corners of ControlCornerRadius, 4; padded 16 as ExpanderContentPadding and SettingsCardPadding, its content 12 apart as Windows's spacing guidance puts content areas), on GNOME a .card (card_bg_color, corners of 12, the card shadow; padded 12 and spaced 12, inferred, as .card sets neither). Rows built in fn (ActionRow, SwitchRow) join into the platform's list instead: a grouped section with inset separators on macOS, a boxed list on GNOME, and on Windows the items of a SettingsExpander, a line between each. It has no role of its own, as a GtkBox with .card and a Border have none.
func Card(c *ui.Context, fn func()) ui.Element {
	p := PlatformOf(c)
	t := c.Theme()
	card := ui.Column(c).Shrink(0)
	g := rowsCollect(c, rowsInCard, func() { card.Children(fn) })
	switch p {
	case Windows:
		fill, _, stroke := fluentExpanderCard(t)
		card.Radius(4).Background(fill).Border(1, stroke).Padding(16).Gap(12)
	case Linux:
		fill, _ := adwBoxedListCard(t)
		card.Radius(12).Background(fill).Padding(12).Gap(12)
		adwBoxedListShadow(card)
	default:
		// Measured from SwiftUI's GroupBox on macOS 27.2: a red Rectangle inside one is 5 points in from every edge, and the box's corners and fill are a grouped section's.
		card.Radius(12).Background(macFormSectionFill(t)).Padding(5).Gap(8)
	}
	rowsFinish(c, card, g)
	return card
}

// ActionRow creates a row of settings: title, subtitle below it in the platform's secondary style (an empty subtitle shows none), and at its end the controls fn builds (nil for none). Alone it is drawn as a group of one row; rows one after another in a Card, a Fieldset or an AccordionItem join as the platform joins them. macOS: a row of a SwiftUI grouped Form as System Settings draws it, 10 points in from the section's edges, the 13-point title in labelColor over the 11-point subtitle in secondaryLabelColor, the controls at the trailing edge, a separator inset 10 points between rows (measured on macOS 27.2). Windows: CommunityToolkit's SettingsCard, 68 DIPs high at least, padded 16, the 14-DIP header and the 12-DIP description in TextFillColorSecondary, the content 24 after them; in a Card or an AccordionItem a SettingsExpander's item, 52 high, padded 16,8, a line of CardStrokeColorDefault above it. GNOME: AdwActionRow, 54 pixels high at least (a 50-pixel header in the row's padding of 2), its header 12 in, the title box padded 6 with 3 between the title and the subtitle (smaller, dimmed), its suffixes 6 apart, the rows of a boxed list a line of card_shade_color apart. The row is a group that assistive technology names by title and describes by subtitle, as AdwActionRow is labelled by its title and described by its subtitle and SettingsCard's peer is a group.
func ActionRow(c *ui.Context, title, subtitle string, fn func()) ui.Element {
	row := ui.Row(c).Shrink(0).Role(ui.RoleGroup).Label(title)
	if subtitle != "" {
		row.Description(subtitle)
	}
	var trailing func(under ui.Color, hovered, pressed bool)
	if fn != nil {
		trailing = func(ui.Color, bool, bool) { rowsIsolate(c, fn) }
	}
	rowsBuild(c, row, title, subtitle, false, trailing)
	return row
}

// SwitchRow creates an ActionRow with the platform's switch at its end toggling *on, the whole row toggling it on a click, Space or Enter, as AdwSwitchRow activates on a click anywhere: it is one switch for assistive technology, named by title and described by subtitle, the switch inside it not a stop of its own (adw-switch-row.c sets the role switch on the row and makes its GtkSwitch unfocusable). macOS: a SwiftUI Toggle of style switch in a grouped Form, a mini switch, 36 by 16 points with a 21 by 12 knob, 2 points in from its top and bottom and 1.5 from its ends (measured on macOS 27.2), AppKit's focus ring around it. Windows: a SettingsCard holding a ToggleSwitch of the card's RightAlignedCompactToggleSwitchStyle, 36 high, its On or Off content 12 DIPs before the track, as Windows Settings shows them (ToggleSwitch's default OnContent and OffContent, TEXT_TOGGLESWITCH_ON and TEXT_TOGGLESWITCH_OFF; Breeze UI shows them in English, without access to MyGo's translations); the card under the pointer and pressed as a clickable SettingsCard (ControlFillColorSecondary and ControlElevationBorderBrush, ControlFillColorTertiary and ControlStrokeColorDefault with the header in TextFillColorSecondary), its focus visual 3 DIPs around it. GNOME: AdwSwitchRow, its GtkSwitch at the end, the row's background currentColor at 3% under the pointer and 8% pressed (%boxed_list_row), its focus ring 2 pixels across its edge (outline-offset -1px). Changed reports a toggle. It is a counterpart of ui.Switch with ui.Field's label.
func SwitchRow(c *ui.Context, title, subtitle string, on *bool) ui.Element {
	row := ui.SwitchBase(c, on).Label(title).FocusRing(false)
	if subtitle != "" {
		row.Description(subtitle)
	}
	value := *on
	focused := row.FocusVisible()
	p := PlatformOf(c)
	rowsBuild(c, row, title, subtitle, true, func(under ui.Color, hovered, pressed bool) {
		u := *c.Theme()
		u.Background = under
		switch p {
		case Windows:
			// The OffContentPresenter or OnContentPresenter in the first column, the 12-DIP column, then the 40 by 20 track (SettingsCard.xaml's RightAlignedCompactToggleSwitchStyle), the content in ToggleSwitchContentForeground, TextFillColorPrimary.
			ui.Row(c).Height(36).Shrink(0).AlignItems(ui.Center).Gap(12).Children(func() {
				state := "Off"
				if value {
					state = "On"
				}
				ui.Text(c, state).SingleLine().TextColor(u.Text)
				track := ui.Box(c).Shrink(0)
				// The ToggleSwitch is a control of its own in the card: under the pointer only over it.
				over := track.Hovered()
				fluentSwitchTrack(c, &u, track, over, over && pressed, value)
			})
		case Linux:
			track := ui.Box(c).Shrink(0)
			// GTK lights the switch only under the pointer over it, not over the row; a press on the row does not press it.
			adwSwitchTrack(c, &u, track, track.Hovered(), false, value)
		default:
			track := ui.Box(c).Shrink(0)
			macMiniSwitch(c, &u, track, value, focused)
		}
	})
	return row
}

// macMiniSwitch draws, on sw, the switch a SwiftUI Toggle shows in a grouped Form on macOS 27.2: a mini NSSwitch, 36 by 16 points, its knob 21 by 12 with corners of 6, 2 points in from the track's top and bottom and 1.5 from its ends (a mini Toggle measures the same), white in both appearances, the track systemFill over what is under it while off and the accent while on, with the regular switch's knob shadow and AppKit's focus ring around the track while focused.
func macMiniSwitch(c *ui.Context, t *ui.Theme, sw ui.Element, on, focused bool) {
	pos := knob(sw, on, 200*time.Millisecond)
	at := along(c, pos)
	off := tint(t.Background, t.Dark, macTroughAlpha)
	sw.Size(36, 16).Radius(8).Background(off.Mix(t.Accent, pos))
	sw.Draw(func(p *ui.Painter, r ui.Rect) {
		const in, w, h = 1.5, 21, 12
		k := ui.Rect{X: r.X + in + at*(r.W-w-2*in), Y: r.Y + (r.H-h)/2, W: w, H: h}
		p.Shadow(k, h/2, 0, 0.5, 1.5, 0, ui.RGBA(0, 0, 0, 0.18))
		p.Fill(k, white, h/2)
	})
	sw.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if focused {
			macFocusRing(p, r, 8)
		}
	})
}

// rowsBuild styles row as the platform's settings row where it is built (rowsJoin) and builds in it the title and the subtitle, in a column that takes the room left, and after them what trailing builds, given the color under the row's end and whether the row is under the pointer and pressed; toggle says the row toggles a switch, reacting to the pointer and showing the focus.
func rowsBuild(c *ui.Context, row ui.Element, title, subtitle string, toggle bool, trailing func(under ui.Color, hovered, pressed bool)) {
	t := c.Theme()
	p := PlatformOf(c)
	item, in, first := rowsJoin(c, row)
	hovered, pressed, focused := false, false, false
	if toggle {
		hovered, pressed, focused = row.Hovered(), row.Pressed(), row.FocusVisible()
	}
	if p == Mac {
		row.AlignItems(ui.Start)
	} else {
		// GTK's suffixes and SettingsCard's content are centered on the row (valign center, VerticalAlignment Center).
		row.AlignItems(ui.Center)
	}
	if in == rowsInAccordion {
		// Under an AccordionItem's header, reaching across MyGo's panel (padded Space(3) at the sides and the bottom, its children Space(2) apart): every row gives back the panel's bottom padding, and the next takes the panel's gap back from it, so the rows touch one another and the item's edges.
		side, gap := t.Space(3), t.Space(2)
		// On macOS a DisclosureGroup in a grouped form leaves 5 points below its last row (measured on macOS 27.2 by reference.swift: the last row's title 23.5 above the section's end).
		keep := float32(0)
		if p == Mac {
			keep = 5
		}
		top := side - keep - gap
		if first {
			top = 0
		}
		row.Margin(top, -side, keep-side, -side)
	}
	var under ui.Color
	var titleColor, subColor ui.Color
	var subSize, lineGap, endGap, suffixGap float32
	switch p {
	case Windows:
		fill, _, stroke := fluentExpanderCard(t)
		under = fill
		joined := in == rowsInCard || in == rowsInAccordion
		if joined {
			// A SettingsExpander's item (DefaultSettingsExpanderItemStyle): MinHeight 52, BorderThickness 0,1,0,0, CornerRadius 0; Padding 58,8,44,8 puts its header past the expander's icon and its content before the chevron, which a Card has neither of, so 16 at both ends as ClickableSettingsExpanderItemPadding's end (inferred).
			row.MinHeight(52).Padding(8, 16)
			if !first || in == rowsInAccordion {
				row.BorderWidth(1, 0, 0, 0).BorderColor(stroke)
			}
		} else {
			row.MinHeight(68).Padding(16).Radius(4).Background(fill).Border(1, stroke)
		}
		titleColor, subColor = t.Text, t.TextMuted
		if toggle && (hovered || pressed) {
			under = fluentFace(t, pressed, hovered)
			row.Background(under)
			if !joined {
				if pressed {
					row.Border(1, t.Border)
				} else {
					h := max(row.Bounds().H, 68)
					row.Border(1, ui.Color{}).BorderLinearGradient(fluentElevation(t, h))
				}
			}
			if pressed {
				titleColor = t.TextMuted
			}
		}
		subSize, lineGap, endGap, suffixGap = 12, 0, 24, 8
		if toggle {
			radius := float32(4)
			if joined {
				radius = 0
			}
			row.DrawOver(func(p *ui.Painter, r ui.Rect) {
				if focused {
					// FocusVisualMargin -3 (SettingsCard.xaml).
					fluentFocusRing(p, t, r, 3, 3, 3, 3, radius)
				}
			})
		}
	case Linux:
		fill, shade := adwBoxedListCard(t)
		under = fill
		// list > row: padding 2px; row > box.header: margin 0 12px, min-height 50px, border-spacing 6px.
		row.MinHeight(54).Padding(2, 14)
		switch in {
		case rowsAlone:
			// A boxed list of one row.
			row.Radius(12).Background(fill)
			adwBoxedListShadow(row)
			item.radius = [4]float32{12, 12, 12, 12}
		case rowsInAccordion:
			// row.expander list.nested: card_shade_color at 50% over the card, a line of card_shade_color above every row, the expander's header above the first.
			under = fill.Mix(adwInk, pick(t.Dark, float32(0.035), 0.18))
			row.Background(under).BorderWidth(1, 0, 0, 0).BorderColor(shade)
		default:
			if !first {
				row.BorderWidth(1, 0, 0, 0).BorderColor(shade)
			}
		}
		if toggle && (hovered || pressed) {
			// %boxed_list_row: color-mix(in srgb, currentColor 3%, transparent) under the pointer, 8% pressed; currentColor is card_fg_color, the shade at 80% in the light theme.
			share := pick(pressed, float32(0.03), 0.08)
			ink := white
			if !t.Dark {
				ink, share = adwInk, share*0.8
			}
			row.Background(under.Mix(ink, share))
		}
		// card_fg_color, and the subtitle .dimmed at --dim-opacity 55%.
		titleColor = pick(t.Dark, fill.Mix(adwInk, 0.8), white)
		subColor = pick(t.Dark, fill.Mix(adwInk, 0.8*0.55), fill.Mix(white, 0.55))
		// label.subtitle: font-size smaller, which GTK takes as the parent's size divided by 1.2.
		subSize, lineGap, endGap, suffixGap = t.FontSize/1.2, 3, 6, 6
		if toggle {
			row.DrawOver(func(p *ui.Painter, r ui.Rect) {
				if focused {
					// The focus-ring mixin with outline-offset -1px: 2 pixels from 1 inside the row's edge to 1 outside it (MyGo's FocusRing draws 1 to 3 outside the box it is given).
					rad := item.radius
					for i := range rad {
						rad[i] = max(rad[i]-2, 0)
					}
					p.FocusRing(ui.Rect{X: r.X + 2, Y: r.Y + 2, W: r.W - 4, H: r.H - 4}, rad)
				}
			})
		}
	default:
		fill := macFormSectionFill(t)
		under = fill
		// 10 points in from the section's sides; 10.5 above and 11 below with 2.5 between the title and the subtitle put their ink where AppKit's is in a 52-point row (the title's 13 points and the subtitle's 30.5 below the row's top, the subtitle's 11.5 above its bottom) and make a row of one line 36.5 points high, as measured on macOS 27.2 with the system font.
		row.Padding(10.5, 10, 11, 10)
		switch in {
		case rowsAlone:
			// A grouped section of one row.
			row.Radius(12).Background(fill)
		case rowsInAccordion:
			// The rows of a DisclosureGroup in a grouped form are closer and have no separators: a row of a title and a subtitle 39.5 points from the next, 12.5 less than in a section (measured on macOS 27.2 by reference.swift).
			row.Padding(4.25, 10, 4.75, 10)
		}
		if in != rowsAlone && in != rowsInAccordion && !first {
			line := tint(t.Background, t.Dark, macSectionSeparatorAlpha)
			row.DrawOver(func(p *ui.Painter, r ui.Rect) {
				p.Fill(ui.Rect{X: r.X + 10, Y: r.Y, W: r.W - 20, H: 1}, line, 0)
			})
		}
		// labelColor and secondaryLabelColor over the section, as AppKit composites them: 38 and 124 over 247 in Aqua, 222 and 157 over 37 in Dark Aqua.
		titleColor = pick(t.Dark, fill.Mix(black, 0.847), fill.Mix(white, 0.847))
		subColor = pick(t.Dark, fill.Mix(black, 0.498), fill.Mix(white, 0.549))
		// The spacing of the controls is inferred.
		subSize, lineGap, endGap, suffixGap = 11, 2.5, 8, 8
	}
	row.Gap(endGap)
	row.Children(func() {
		labels := ui.Column(c).Grow(1).Shrink(1).MinWidth(0).Gap(lineGap)
		if p == Linux {
			// row > box.header > box.title: margin-top and margin-bottom 6px.
			labels.Padding(6, 0)
		}
		labels.Children(func() {
			ui.Text(c, title).TextColor(titleColor)
			if subtitle != "" {
				ui.Text(c, subtitle).FontSize(subSize).TextColor(subColor)
			}
		})
		if trailing != nil {
			end := ui.Row(c).Shrink(0).AlignItems(ui.Center).Gap(suffixGap)
			if p == Mac {
				// A grouped form centers its controls on the title's line, not on the row: a 16-point switch's top and a 24-point button's middle line up with the title's (measured on macOS 27.2), the button reaching into the row's padding.
				_, line := c.MeasureText(0, ui.Span{Text: title, Size: t.FontSize})
				end.Height(line)
			}
			end.Children(func() { trailing(under, hovered, pressed) })
		}
	})
}
