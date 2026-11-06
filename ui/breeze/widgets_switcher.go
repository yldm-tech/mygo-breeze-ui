package breeze

import "github.com/egoist/mygo/ui"

// ---- view switchers: the control that switches an app's top-level views ----
//
// macOS: an NSToolbarItemGroup of selectionMode selectOne in a window's unified toolbar, as Music and Finder switch their views there, and an NSSegmentedControl of segments holding an image and a label outside a toolbar, which tools/appkit-reference/reference.swift's switcher scene (renderSwitcher) renders at 2x in Aqua and Dark Aqua on macOS 27.2 where TestScreenshotsViewSwitcher places Breeze UI's, offscreen and captured live (-live), printing the group's and the control's frames and the values across them. The live capture sets the chosen segment's pill, the offscreen render the group's and the segments' sizes; the group's segments for other images and for titles were measured with NSToolbarItemGroups of 8- to 32-point template images and of titles laid out in the same toolbar.
//
// Windows 11: WinUI's SelectorBar of microsoft/microsoft-ui-xaml main (controls/dev/SelectorBar): SelectorBar.xaml (the SelectorBarItem template: the icon and the text in a horizontal StackPanel spaced SelectorBarItemSpacing, the icon's ContentPresenter scaled by SelectorBarItemIconScale with SelectorBarItemIconVisualMargin, the pill in a row of its own under them), SelectorBar_themeresources.xaml (SelectorBarPadding 0,4, SelectorBarItemPadding 12,10,12,7, the foregrounds of each state, SelectorBarItemPillHeight 3 and Width 4, the transparent backgrounds of every state) and SelectorBarItemAutomationPeer.cpp (an item named by its Text).
//
// GNOME: libadwaita 1.10.0's AdwViewSwitcher: adw-view-switcher.c (a homogeneous GtkBoxLayout, the tab-list role, the policy, narrow by default), adw-view-switcher-button.ui and adw-view-switcher-button.c (a flat GtkToggleButton with the tab role, the icon beside the label in a .wide box, MIN_NAT_BUTTON_WIDTH_WIDE 120), _view-switcher.scss (viewswitcher's border-spacing of 3 and min-height of 34, the buttons' padding of 0, the .wide box's padding 2px 12px and border-spacing 6px) and _buttons.scss (%button_basic_flat's fills, the buttons' bold label and corners of 9).

// ViewItem is a view a ViewSwitcher switches to: its icon, nil for none, and its label, which names the item for assistive technology on every platform, also where only the icon shows.
type ViewItem struct {
	Icon  *ui.SVG
	Label string
}

// switcherMacImageSegment is the width of an image segment of an NSToolbarItemGroup in a unified toolbar: max(37, the image's width + 17) points, measured on macOS 27.2 with template images 8 to 32 points wide (a group of three 16-point images is 111 points wide, of three 24-point ones 123); Breeze UI's icons are 16 points.
const switcherMacImageSegment = 37

// switcherMacTitlePadding and switcherMacTitleLeast size a title segment of an NSToolbarItemGroup: the label's advance plus 23.25 points, at least 35, the segments as wide as the widest, measured on macOS 27.2 (three Browse segments, 44.38 points of 13-point text, make a group of 202.5; three A segments one of 105). AppKit gives all the segments the widest one's width when it is the last; when it is not, they come out a few points narrower (Radio, Browse, Radio make 196.5, not 202.5), which is not reproduced.
const (
	switcherMacTitlePadding = 23.25
	switcherMacTitleLeast   = 35
)

// switcherAdwLeast is MIN_NAT_BUTTON_WIDTH_WIDE (adw-view-switcher-button.c): a button of a wide AdwViewSwitcher asks for at least 120 pixels, and the homogeneous box gives every button the widest one's width.
const switcherAdwLeast = 120

// ViewSwitcher creates the control that switches an app's top-level views, one of items, chosen by *selected, as the platform draws its own; build the chosen view after it. Put it in a HeaderBar's end controls or in a Toolbar to switch the window's views from there, or anywhere in the view.
//
// It is MyGo's SegmentedBase on macOS, a radio group of radio buttons as AppKit reports an NSToolbarItemGroup and an NSSegmentedControl (AXRadioGroup, measured on macOS 27.2), and MyGo's TabsBase on Windows and GNOME, a tab list of tabs as AdwViewSwitcher is (GTK_ACCESSIBLE_ROLE_TAB_LIST and TAB) and the nearest MyGo has to a SelectorBar's list of selectable items. Either way Tab stops once in it, at the chosen item, the arrows choose among the items, mirrored right to left, and each item is named by its label, as SelectorBarItemAutomationPeer names an item by its Text and AppKit a segment by its label.
//   - macOS: in a Toolbar or a HeaderBar, an NSToolbarItemGroup of selectionMode selectOne as a key window's unified toolbar draws it: one 36-point platter of a bordered toolbar item (idiomMacPlatter) holding the segments, each 37 points wide around its 16-point icon in black (white in Dark Aqua); the chosen segment a pill 28 points high 4 inside the platter's top and bottom, 4 in from the platter's start in the first segment and from its end in the last, otherwise 3 after the boundary before the segment and 2 before the one after it (switcherMacInsets), each segment's icon centered on its pill, filled as a pressed toolbar item is (black at 7.8% in Aqua, (89,88,88) in Dark Aqua: the live capture reads 232 on a 252 platter and (91,90,90) on (71,70,70)); a separator 1 by 20 points at the boundary between two segments neither of which is chosen, black at 20% (white at 28% in Dark Aqua) over the platter, read from the live capture. The toolbar shows the icons only, as AppKit's does in its icon-only display mode, the labels naming the segments; when an item has no icon the segments show the labels instead, in 13 points, as wide as the widest. A pressed segment shows the pill too and the focus ring goes around the focused segment's pill (inferred). Outside a toolbar, Segmented's NSSegmentedControl, each segment's 16-point icon 2 points before its label (measured: three image and label segments are 18 points wider each than label segments alone).
//   - Windows: WinUI's SelectorBar, padded 4 above and below: items 40 DIPs high, padded 12,10,12,7 above a 3-DIP row for the pill, the 16-DIP icon (an IconElement of 20, scaled 0.8 with a margin of -2 at its sides) 8 before the text in ControlContentThemeFontSize (14), in TextFillColorPrimary, TextFillColorSecondary under the pointer and TextFillColorTertiary pressed (the chosen item TextFillColorSecondary under the pointer and pressed); the chosen item's accent pill, 16 by 3, centered at its bottom, growing from 4 wide as WinUI animates it. WinUI's SelectorBar fills no item under the pointer or pressed (its SelectorBarItemBackground brushes are transparent in every state), so only the text changes.
//   - GNOME: AdwViewSwitcher with the wide policy, as in a header bar: buttons 34 pixels high, 3 apart, all as wide as the widest and at least 120, each a flat button with corners of 9 holding the 16-pixel icon 6 before its bold label, padded 2 12 and centered; the chosen one currentColor at 10% (13% under the pointer, 19% pressed), another one 7% under the pointer and 16% pressed, with the focus ring inside its edge.
//
// AdwViewSwitcher's narrow policy (the icon over the label), AdwViewSwitcherBar and the needs-attention dot and badge are not drawn, and a HeaderBar holds the switcher among its end controls rather than in place of its title, as libadwaita apps put it: see the documentation's Limits. It returns the radio group or the tab list, whose Changed reports nothing: read *selected.
func ViewSwitcher(c *ui.Context, selected *int, items ...ViewItem) ui.Element {
	p := PlatformOf(c)
	n := len(items)
	var root ui.Element
	var item func(i int) ui.Element
	if p == Mac {
		parts := ui.SegmentedBase(c, selected, n)
		root, item = parts.Track, parts.Segment
	} else {
		parts := ui.TabsBase(c, selected, n)
		root, item = parts.List, parts.Tab
	}
	// The group is an item of a toolbar it is built in, never joining the image items beside it, as AppKit's item groups do not.
	tb, _ := idiomToolbarItem(c, root, idiomOther)
	sel := *selected
	switch {
	case p == Windows:
		switcherFluent(c, root, item, items, sel)
	case p == Linux:
		switcherAdw(c, root, item, items, sel)
	case tb != nil:
		switcherMacBar(c, root, item, items, sel)
	default:
		switcherMac(c, root, item, items, sel)
	}
	return root
}

// switcherContent builds what an item shows in a row, gap apart: its icon, size square, unless icon is false or the item has none, and its label unless label is false, both built in every frame and hidden when not shown, so that the item's children keep their places. It returns the row, the icon and the label.
func switcherContent(c *ui.Context, it ViewItem, icon, label bool, size, gap float32) (row, img, text ui.Element) {
	row = ui.Row(c).Shrink(0).AlignItems(ui.Center).Gap(gap)
	row.Children(func() {
		img = ui.Icon(c, it.Icon).Size(size, size).Shrink(0).Hidden(!icon || it.Icon == nil)
		text = ui.Text(c, it.Label).SingleLine().Hidden(!label || it.Label == "")
	})
	return row, img, text
}

// switcherEqualWidths gives every item the width of the widest content of the previous frame plus pad, and at least least, as GtkBoxLayout's homogeneous layout and an NSToolbarItemGroup make their buttons equal, and asks for another frame while that width changes, as a width is known only once laid out.
func switcherEqualWidths(c *ui.Context, root ui.Element, items, contents []ui.Element, pad, least float32) {
	widest := float32(0)
	for _, e := range contents {
		widest = max(widest, e.LayoutBounds().W)
	}
	w := max(least, widest+pad)
	for _, e := range items {
		e.MinWidth(w)
	}
	last := ui.Local(root, "breeze-switcher-width", func() float32 { return -1 })
	if *last != w {
		*last = w
		c.Invalidate()
	}
}

// ---- macOS ----

// switcherMac draws the switcher outside a toolbar as Segmented draws an NSSegmentedControl, each segment's icon 2 points before its label.
func switcherMac(c *ui.Context, root ui.Element, item func(int) ui.Element, items []ViewItem, sel int) {
	macSegmentedTrack(c, root)
	root.Children(func() {
		for i, it := range items {
			seg := item(i).Label(it.Label)
			// A separator shows between two segments neither of which is chosen, as AppKit draws them.
			macSegment(c, seg, i == sel, i > 0 && i != sel && i-1 != sel)
			seg.Children(func() { switcherContent(c, it, true, true, 16, 2) })
		}
	})
}

// switcherMacBar draws the switcher as an NSToolbarItemGroup of selectionMode selectOne in a key window's unified toolbar (see ViewSwitcher): one bordered item's platter holding the segments, the chosen one's pill, and separators between two segments neither of which is chosen.
func switcherMacBar(c *ui.Context, root ui.Element, item func(int) ui.Element, items []ViewItem, sel int) {
	t := c.Theme()
	n := len(items)
	images := n > 0
	for _, it := range items {
		images = images && it.Icon != nil
	}
	root.Height(36).Radius(18).Padding(0).Gap(0).AlignItems(ui.Stretch).Shrink(0)
	pressed, focused := -1, -1
	segs := make([]ui.Element, n)
	labels := make([]ui.Element, n)
	root.Children(func() {
		for i, it := range items {
			seg := item(i).Label(it.Label).Padding(0).Justify(ui.Center).AlignItems(ui.Center).FocusRing(false).TextColor(idiomMacInk(t, false))
			if images {
				seg.Width(switcherMacImageSegment)
			}
			if seg.Pressed() {
				pressed = i
			}
			if seg.FocusVisible() {
				focused = i
			}
			segs[i] = seg
			seg.Children(func() {
				var row ui.Element
				row, _, labels[i] = switcherContent(c, it, images, !images, 16, 0)
				d := switcherMacShift(n, i)
				row.MarginStart(d).MarginEnd(-d)
			})
		}
	})
	if !images {
		switcherEqualWidths(c, root, segs, labels, switcherMacTitlePadding, switcherMacTitleLeast)
	}
	rtl := c.Direction() == ui.RightToLeft
	disabled := disabledOnce(root)
	face := idiomMacPressedFace(t)
	line := pick(t.Dark, ui.RGBA(0, 0, 0, 0.2), ui.RGBA(255, 255, 255, 0.28))
	root.Draw(func(p *ui.Painter, r ui.Rect) {
		idiomMacPlatter(p, t, r, idiomMacRest, face, false, false, rtl, disabled())
		for i := range n {
			if i == sel || i == pressed {
				p.Fill(switcherMacPill(r, n, i, rtl), face, (r.H-8)/2)
			}
			if i > 0 && i != sel && i-1 != sel {
				// The boundary between segments i-1 and i, the line on segment i's side of it.
				x := r.X + float32(i)*r.W/float32(n)
				if rtl {
					x = r.X + r.W - float32(i)*r.W/float32(n) - 1
				}
				p.Fill(ui.Rect{X: x, Y: r.Y + (r.H-20)/2, W: 1, H: 20}, line, 0)
			}
		}
	})
	root.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if focused >= 0 {
			macFocusRing(p, switcherMacPill(r, n, focused, rtl), (r.H-8)/2)
		}
	})
}

// switcherMacInsets returns how far the pill of segment i of n stands in from the segment's start and from its end, as AppKit places it: 4 points from the platter's start in the first segment and from its end in the last, otherwise 3 after the boundary before the segment (the separator's point and 2) and 2 before the boundary after it. Measured on macOS 27.2 with NSToolbarItemGroups of two to five 16-point images (37-point segments), of three 24-point ones (41) and of three titles, chosen in turn, offscreen and live: in a 111-point group of three the pills span 4 to 35, 40 to 72 and 77 to 107 points, in a 148-point group of four 4 to 35, 40 to 72, 77 to 109 and 114 to 144, in a 74-point group of two 4 to 35 and 40 to 70. A single segment's pill is 4 in from both ends (inferred).
func switcherMacInsets(n, i int) (start, end float32) {
	start, end = 3, 2
	if i == 0 {
		start = 4
	}
	if i == n-1 {
		end = 4
	}
	return start, end
}

// switcherMacPill returns the pill of segment i of n in the group's box r (see ViewSwitcher and switcherMacInsets): 4 points inside the platter's top and bottom, the segments sharing the box equally; right to left it is mirrored (inferred).
func switcherMacPill(r ui.Rect, n, i int, rtl bool) ui.Rect {
	seg := r.W / float32(n)
	start, end := switcherMacInsets(n, i)
	x0, x1 := float32(i)*seg+start, float32(i+1)*seg-end
	if rtl {
		x0, x1 = r.W-x1, r.W-x0
	}
	return ui.Rect{X: r.X + x0, Y: r.Y + 4, W: x1 - x0, H: r.H - 8}
}

// switcherMacShift is how far toward the group's end AppKit moves the image or the title of segment i of n from the segment's middle, centering it on the segment's pill (switcherMacPill) whether or not the segment is chosen: half the difference of the pill's insets, a point toward the end in the first segment, half a point in the ones between and half a point back in the last (measured on macOS 27.2 with the groups switcherMacInsets lists: the images' middles at 19.5, 56 and 92 points into a group of three, 19.5, 56, 93 and 129 into one of four).
func switcherMacShift(n, i int) float32 {
	start, end := switcherMacInsets(n, i)
	return (start - end) / 2
}

// ---- Windows ----

// switcherFluent draws the switcher as WinUI's SelectorBar (see ViewSwitcher), its items Segmented's SelectorBarItems holding an icon 8 DIPs before the text: the icon's box 20 DIPs high, as an IconElement of 20 scaled 0.8 lays out, and 16 wide, its margin of -2 at either side taking back what the scale leaves.
func switcherFluent(c *ui.Context, root ui.Element, item func(int) ui.Element, items []ViewItem, sel int) {
	fluentSegmentedTrack(c, root)
	root.Children(func() {
		for i, it := range items {
			tab := item(i).Label(it.Label)
			fluentSegment(c, tab, i == sel)
			tab.Children(func() {
				_, img, _ := switcherContent(c, it, true, true, 16, 8)
				img.MarginY(2)
			})
		}
	})
}

// ---- GNOME ----

// switcherAdw draws the switcher as AdwViewSwitcher with the wide policy (see ViewSwitcher): flat toggle buttons 3 apart, at least 34 high, as wide as the widest and at least 120, each the icon 6 before the bold label, padded 2 12.
func switcherAdw(c *ui.Context, root ui.Element, item func(int) ui.Element, items []ViewItem, sel int) {
	root.Gap(3).MinHeight(34).AlignItems(ui.Stretch)
	n := len(items)
	tabs := make([]ui.Element, n)
	rows := make([]ui.Element, n)
	root.Children(func() {
		for i, it := range items {
			tab := item(i).Label(it.Label).Radius(9).Padding(2, 12).Justify(ui.Center).AlignItems(ui.Center).Border(0, ui.Color{}).FontWeight(700)
			idiomAdwFlatFace(c, tab, i == sel)
			tabs[i] = tab
			tab.Children(func() { rows[i], _, _ = switcherContent(c, it, true, true, 16, 6) })
		}
	})
	switcherEqualWidths(c, root, tabs, rows, 24, switcherAdwLeast)
}
