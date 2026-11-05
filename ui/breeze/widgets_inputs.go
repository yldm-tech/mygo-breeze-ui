package breeze

import (
	"github.com/egoist/mygo/ui"
)

// Inputs that MyGo has no widget for, built from its bases and drawn as each platform draws them: a password field, a split button and a button showing only an icon.

// ---- password fields ----

// PasswordField creates a single-line field editing *value that hides its text, MyGo's TextInput with Password, in the box the platform draws around its text fields (TextInput's), with the platform's own ways to see the text:
//   - macOS: an NSSecureTextField, the text field's bezel with nothing to reveal the text (AppKit offers none). AppKit 14 and later show a caps-lock indicator at its end while Caps Lock is on, which Breeze UI cannot draw: MyGo does not tell a view whether Caps Lock is on.
//   - Windows: a PasswordBox whose PasswordRevealMode is Peek, its default: the RevealButton, named "Show Password", shows at its end while it has the focus and holds text typed since it got it, and reveals the text while it is pressed; Alt+F8 reveals it while held, as WinUI's accelerator does. The button is not a Tab stop, as in WinUI (IsTabStop False), so the focus stays in the field.
//   - GNOME: a GtkPasswordEntry, an entry whose context menu ends with a "Show Text" check item that shows or hides the text, from the keyboard too (the menu key or Shift+F10). GTK leaves the peek icon off by default (show-peek-icon), and so does Breeze UI: AdwPasswordEntryRow shows one, but it is a row of a boxed list, not an entry. Its caps-lock warning icon is left out, as on macOS.
//
// Editing, the selection, Submitted on Enter, the focus and what assistive technology hears (a secure text field, whose text it is not told) are MyGo's; the bullets are MyGo's too (•, where AppKit and WinUI draw larger dots). Copy and Cut do nothing while the text is hidden.
func PasswordField(c *ui.Context, value *string) ui.Element {
	in := ui.TextInput(c, value)
	switch PlatformOf(c) {
	case Windows:
		fluentPasswordBox(c, in, value)
	case Linux:
		adwPasswordEntry(c, in)
	default:
		in.Password()
		macField(c, in, singleLine)
	}
	return in
}

// fluentPeek is what a PasswordBox remembers between frames to show its RevealButton as PasswordBox.cpp does: the text it last saw, whether it had the focus, whether the text was typed into it since it got the focus (m_fCanShowRevealButton), and whether Alt+F8 is held.
type fluentPeek struct {
	last                  string
	focused, canShow, key bool
}

// fluentPeekKeys returns the input handler of a PasswordBox: Alt+F8 down reveals the text and F8 up hides it again (PasswordBox_Partial.cpp's OnKeyDown, which checks IsMenuKeyDown, and OnKeyUp). It takes no other input.
func fluentPeekKeys(st *fluentPeek) func(ev ui.InputEvent) bool {
	return func(ev ui.InputEvent) bool {
		switch {
		case ev.Kind == ui.InputKeyDown && ev.Key == ui.KeyF8 && ev.Mods&ui.Alt != 0:
			st.key = true
			return true
		case ev.Kind == ui.InputKeyUp && ev.Key == ui.KeyF8 && st.key:
			st.key = false
			return true
		}
		return false
	}
}

// fluentPasswordBox draws in as a PasswordBox (PasswordBox_themeresources.xaml of microsoft-ui-xaml main at ace5f33): the TextBox look (fluentField), and in its template's second column the RevealButton, a ToggleButton 30 DIPs wide stretched to the box's height whose ButtonLayoutGrid sits inside TextBoxInnerButtonMargin 0,4,4,4 (26 by 24 in a 32-DIP box) with ControlCornerRadius, transparent at rest (TextControlButtonBackground, which TextBox_themeresources.xaml leaves to generic.xaml: inferred, as for NumberBox), SubtleFillColorSecondary under the pointer and SubtleFillColorTertiary pressed, its glyph U+F78D of Segoe Fluent Icons at PasswordBoxIconFontSize 12 in TextFillColorSecondary, TextFillColorTertiary pressed (TextControlButtonForeground*). The text stops before that column while the button shows: ContentElement is in the first column with TextControlThemePadding 10,5,6,6 inside the 1-DIP border.
//
// When it shows is PasswordBox.cpp's CanInvokeRevealButton: the box has the focus, it is wider than 5 em of its FontSize 14 (ArrangeOverride), and its text went from empty to not since it got the focus (OnGotFocus clears m_fCanShowRevealButton, and emptying the box clears it again), so that a password filled in before cannot be peeked at. The text shows while the button is pressed or Alt+F8 is held, and the button's IsChecked follows (RevealPassword), which assistive technology sees.
func fluentPasswordBox(c *ui.Context, in ui.Element, value *string) {
	t := c.Theme()
	fluentField(c, in, singleLine)
	st := ui.Local(in, "breeze-peek", func() fluentPeek { return fluentPeek{last: *value} })
	focused := in.Focused()
	if focused && !st.focused {
		st.canShow = false
	}
	if !focused {
		// F8's release goes to whatever has the focus then.
		st.key = false
	}
	if *value != st.last {
		if st.last == "" {
			st.canShow = true
		}
		if *value == "" {
			st.canShow = false
		}
		st.last = *value
	}
	st.focused = focused
	in.HandleInput(fluentPeekKeys(st))
	show := st.canShow && focused && in.Bounds().W > 5*14
	// The button is built whether it shows or not, Hidden while not, so that the elements built after the field keep their places, by which MyGo knows them (the focus, moved to the next of them by Tab as the button hides, would otherwise land on the one after).
	b := ui.Box(c).Role(ui.RoleToggleButton).Label("Show Password").KeepFocus().AttachTo(in, ui.AnchorEnd, ui.AnchorEnd).MarginEnd(4).Size(26, 24).Radius(4).Hidden(!show)
	pressed := false
	if show {
		// The text stops 6 DIPs before the button's 30-DIP column, past the 1-DIP border's room the column takes.
		in.PaddingEnd(37)
		pressed = b.Pressed()
		hover, down := fluentSubtle(t)
		glyph := t.TextMuted
		switch {
		case pressed:
			b.Background(down)
			glyph = t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.447), 0.529))
		case b.Hovered():
			b.Background(hover)
		}
		b.Checked(pressed || st.key)
		b.Draw(func(p *ui.Painter, r ui.Rect) { fluentRevealGlyph(p, r.X+r.W/2, r.Y+r.H/2, glyph) })
	}
	if !show || !pressed && !st.key {
		in.Password()
	}
}

// fluentRevealGlyph strokes U+F78D of Segoe Fluent Icons at 12 DIPs centered on cx, cy: an eye, an almond 11 DIPs wide and 7 high around a pupil 4 across, in 1-DIP strokes. The font is not open: the shape is inferred from the glyph as Windows shows it.
func fluentRevealGlyph(p *ui.Painter, cx, cy float32, col ui.Color) {
	var eye ui.Path
	eye.MoveTo(cx-5.5, cy).CubeTo(cx-3, cy-4.5, cx+3, cy-4.5, cx+5.5, cy).CubeTo(cx+3, cy+4.5, cx-3, cy+4.5, cx-5.5, cy).Close()
	p.StrokePath(&eye, 1, col)
	p.Stroke(ui.Rect{X: cx - 2, Y: cy - 2, W: 4, H: 4}, col, 2, 1)
}

// adwPasswordEntry draws in as a GtkPasswordEntry (gtkpasswordentry.c of GTK 4.24.0): the entry's look (adwField; libadwaita 1.10.0's entry.password only dims the caps-lock indicator, _entries.scss), the text hidden until the "_Show Text" item GTK adds after the editing items of the context menu (gtk_password_entry_set_extra_menu, the stateful action misc.toggle-visibility, shown checked while the text shows) turns it on; the same item hides it again.
func adwPasswordEntry(c *ui.Context, in ui.Element) {
	shown := ui.Local(in, "breeze-show-text", func() bool { return false })
	if !*shown {
		in.Password()
	}
	in.ContextMenu(func(m *ui.Menu) {
		m.EditItems()
		m.Separator()
		if m.Item("Show Text").Checked(*shown).Chosen() {
			*shown = !*shown
		}
	})
	adwField(c, in, singleLine)
}

// ---- split buttons ----

// inputsMoreOptions names the part of a SplitButton that opens its menu: libadwaita's default tooltip of AdwSplitButton's dropdown (adw-split-button.c), which Breeze UI gives as the part's accessible name on every platform: WinUI's secondary button is not in the accessibility tree (AccessibilityView Raw), the SplitButton itself expanding, and AppKit names no part of an NSComboButton.
const inputsMoreOptions = "More Options"

// SplitButton creates a button of two parts side by side: the returned one shows label and does the main action, which Clicked reports, and the other, after it (before it right to left), shows the platform's arrow and opens the system's own menu below it, which build fills, as MenuButton's does. It is drawn as macOS's NSComboButton of style split (a push button with a divider before the arrow), Windows's SplitButton (SplitButton.xaml) and GNOME's AdwSplitButton; inside a Toolbar, as the platform's toolbar split button: an NSComboButton in a toolbar item's platter, the SplitButtonCommandBarStyle of a CommandBar, and a flat AdwSplitButton.
//
// The two parts are ui.ButtonBase's in a row that assistive technology sees as a group, as AdwSplitButton's: the main part a button named by its label, the other a menu button named "More Options". Both are Tab stops, as in GTK, where WinUI's and AppKit's controls take one; Enter and Space click the main part, and Enter, Space and Down open the menu from the other, which the arrows reach in a Toolbar. Styles given to what SplitButton returns apply to the main part only, but for Disabled, which the other part takes from it in the next frame: give margins and sizes to an element around it.
func SplitButton(c *ui.Context, label string, build func(m *ui.Menu)) ui.Element {
	p := PlatformOf(c)
	t := c.Theme()
	rtl := c.Direction() == ui.RightToLeft
	var primary, arrow ui.Element
	group := ui.Row(c).Shrink(0).AlignItems(ui.Stretch).Role(ui.RoleGroup)
	tb, _ := idiomToolbarItem(c, group, idiomOther)
	bar := tb != nil
	group.Children(func() {
		primary = ui.ButtonBase(c).FocusRing(false).Justify(ui.Center).AlignItems(ui.Center)
		primary.Children(func() { ui.Text(c, label).SingleLine() })
		if p == Windows && !bar {
			// DividerBackgroundGrid, built between the two so that it lays out there; drawn over the main part's background, which spans its column.
			face := fluentFace(t, primary.Pressed(), primary.Hovered())
			ui.Box(c).Width(1).Shrink(0).Role(ui.RoleNone).Background(face.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.059), 0.071)))
		}
		arrow = ui.ButtonBase(c).Label(inputsMoreOptions).FocusRing(false).Shrink(0).AlignItems(ui.Center)
	})
	// The other part follows the main part's Disabled, which the caller gives after SplitButton returns: it is read as the group paints and given to the other part in the frame it asks for.
	off := ui.Local(primary, "breeze-split-disabled", func() bool { return false })
	arrow.Disabled(*off)
	disabled := disabledOnce(primary)
	group.DrawOver(func(pt *ui.Painter, _ ui.Rect) {
		if d := disabled(); d != *off {
			*off = d
			pt.After(0)
		}
	})
	start, end := inputsCorners(rtl)
	switch {
	case p == Windows && bar:
		fluentSplitAppBar(c, primary, arrow, start, end)
	case p == Windows:
		fluentSplitButton(c, primary, arrow, start, end)
	case p == Linux && bar:
		adwSplitButton(c, group, primary, arrow, start, end, true)
	case p == Linux:
		adwSplitButton(c, group, primary, arrow, start, end, false)
	case bar:
		macSplitToolbarItem(c, group, primary, arrow, disabled)
	default:
		macComboButton(c, group, primary, arrow, start, end)
	}
	arrow.Menu(build)
	return primary
}

// inputsCorners returns which physical corners are a part's at the start and at the end of a row (ui.Element.Radius takes top-left, top-right, bottom-right and bottom-left), as a row mirrors right to left.
func inputsCorners(rtl bool) (start, end func(r float32) [4]float32) {
	left := func(r float32) [4]float32 { return [4]float32{r, 0, 0, r} }
	right := func(r float32) [4]float32 { return [4]float32{0, r, r, 0} }
	if rtl {
		return right, left
	}
	return left, right
}

// inputsRadius gives e the corners k: top-left, top-right, bottom-right and bottom-left.
func inputsRadius(e ui.Element, k [4]float32) ui.Element { return e.Radius(k[0], k[1], k[2], k[3]) }

// macComboButton draws an NSComboButton of style split (macOS 13 and later), measured from AppKit's rendering on macOS 27.2 at 2x by tools/appkit-reference (appkit-inputs-*.png): a push button's 24-point face (macFace) with corners of 6, its title in labelColor over the face (36 over 235, 223 over 48, as macButton's) centered in the main part, which is the title's advance and 12.5 points either side ("Save", 29.1 points, in a 54-point part), then a divider a point wide and 14 high, centered, black at 23.4% over the face in Aqua (180 over 235) and white at 25.1% in Dark Aqua (100 over 48), and a 16-point part holding the chevron.down of its arrow, 9 by 5 points of ink in the title's color, its stroke about a point (macComboChevron). A part pressed takes a pressed push button's face (macPressed, inferred: AppKit does not render a press offscreen), and AppKit's focus ring goes around the part with the keyboard's focus (inferred: an NSComboButton takes the focus once).
func macComboButton(c *ui.Context, group, primary, arrow ui.Element, start, end func(float32) [4]float32) {
	t := c.Theme()
	face := macFace(t.Background, t.Dark)
	ink := tint(face, t.Dark, 0.847)
	group.Height(24).Radius(6).Background(face)
	primary.PaddingX(12.5).TextColor(ink)
	if primary.Pressed() {
		inputsRadius(primary.Background(macPressed(t.Background, t.Dark)), start(6))
	}
	macRing(primary, 6)
	line := tint(face, t.Dark, pick(t.Dark, float32(0.234), 0.251))
	arrow.Width(17).Justify(ui.Center).PaddingStart(1)
	if arrow.Pressed() {
		inputsRadius(arrow.Background(macPressed(t.Background, t.Dark)), end(6))
	}
	rtl := c.Direction() == ui.RightToLeft
	arrow.Draw(func(p *ui.Painter, r ui.Rect) {
		x := r.X
		if rtl {
			x = r.X + r.W - 1
		}
		p.Fill(ui.Rect{X: x, Y: r.Y + (r.H-14)/2, W: 1, H: 14}, line, 0)
	})
	macRing(arrow, 6)
	arrow.Children(func() { macComboChevron(c, 9, 5, 1, ink) })
}

// macComboChevron draws the chevron.down of an NSComboButton's arrow: w by h points of ink, its stroke width wide.
func macComboChevron(c *ui.Context, w, h, width float32, color ui.Color) {
	ui.Box(c).Size(w, h).Shrink(0).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
		in := width / 2
		var path ui.Path
		path.MoveTo(r.X+in, r.Y+in).LineTo(r.X+r.W/2, r.Y+r.H-in).LineTo(r.X+r.W-in, r.Y+in)
		p.StrokePath(&path, width, color)
	})
}

// macSplitToolbarItem draws an NSComboButton as the view of an NSToolbarItem, as AppKit renders a window's unified toolbar on macOS 27.2 at 2x (tools/appkit-reference): the item's 36-point platter (idiomMacPlatter) holding the main part as a pill 4 points in from the platter's start, its title padded 7.5 points either side ("New", 26.9 points, in a 42-point pill), filled with black at 9.2% over the platter in Aqua and white at 9.3% in Dark Aqua (228 over 251, 49 over 28 in the render, which draws the platter lighter than a live window does), and after it a 40-point part holding the arrow, a chevron 9 by 5 points stroked about 1.2 points in the toolbar's ink (idiomMacInk). The pill is 28 points high, 4 inside the platter as the pressed slice of a joined platter is (inferred: the offscreen render draws the platter 32 points high). A pill pressed takes a pressed push button's share (16.9%), and the arrow part pressed a pill of its own (both inferred).
func macSplitToolbarItem(c *ui.Context, group, primary, arrow ui.Element, disabled func() bool) {
	t := c.Theme()
	rtl := c.Direction() == ui.RightToLeft
	group.Height(36).Radius(18).Gap(0).PaddingStart(4).Border(0, ui.Color{})
	pressed := idiomMacPressedFace(t)
	group.Draw(func(p *ui.Painter, r ui.Rect) {
		idiomMacPlatter(p, t, r, idiomMacRest, pressed, false, false, rtl, disabled())
	})
	platter := pick(t.Dark, white, ui.RGB(69, 68, 68))
	ink := idiomMacInk(t, false)
	share := pick(t.Dark, float32(0.092), 0.093)
	if primary.Pressed() {
		share = macPressedAlpha
	}
	primary.Height(28).AlignSelf(ui.Center).PaddingX(7.5).Radius(14).Background(tint(platter, t.Dark, share)).TextColor(ink)
	macRing(primary, 14)
	arrow.Width(40).Height(28).AlignSelf(ui.Center).Justify(ui.Center).Radius(14)
	if arrow.Pressed() {
		arrow.Background(tint(platter, t.Dark, macPressedAlpha))
	}
	macRing(arrow, 14)
	arrow.Children(func() { macComboChevron(c, 9, 5, 1.2, ink) })
}

// fluentSplitButton draws a SplitButton (SplitButton.xaml and SplitButton_themeresources.xaml of microsoft-ui-xaml main at ace5f33): its columns are the main part, at least SplitButtonPrimaryButtonSize (35) wide, padded by SplitButtonPadding 11,6,11,7, a 1-DIP DividerBackgroundGrid in SplitButtonBorderBrushDivider (ControlStrokeColorDefault over the main part's background, which spans that column), and the secondary part, SplitButtonSecondaryButtonSize (35) wide, its 12-DIP AnimatedChevronDownSmall at its end padded 12 (Padding 0,0,12,0), in SplitButtonForegroundSecondary (TextFillColorSecondary over the face), TextFillColorTertiary pressed. Each part is a button's face, ControlFillColorDefault, Secondary under the pointer and Tertiary pressed, each state its own part's (PrimaryPointerOver, SecondaryPressed); the border is drawn by PrimaryButtonBorder (1,1,0,1, corners 4,0,0,4) and SecondaryButtonBorder (0,1,1,1, corners 0,4,4,0), ControlElevationBorderBrush, a part's ControlStrokeColorDefault while pressed, leaving the divider's column without a top or bottom edge. The label is TextFillColorPrimary, Secondary pressed. The system focus visual goes 1 DIP around (FocusVisualMargin -1) the part with the keyboard's focus, where WinUI's control takes the focus once.
func fluentSplitButton(c *ui.Context, primary, arrow ui.Element, start, end func(float32) [4]float32) {
	t := c.Theme()
	part := func(b ui.Element, k [4]float32) (face ui.Color, pressed bool) {
		pressed = b.Pressed()
		face = fluentFace(t, pressed, b.Hovered())
		inputsRadius(b.MinHeight(32).Background(face), k)
		b.Border(1, ui.Color{})
		if pressed {
			b.BorderColor(t.Border)
		} else {
			b.BorderLinearGradient(fluentElevation(t, 32))
		}
		focused := b.FocusVisible()
		b.DrawOver(func(p *ui.Painter, r ui.Rect) {
			if focused {
				fluentFocusRing(p, t, r, 1, 1, 1, 1, 4)
			}
		})
		return face, pressed
	}
	_, pressed := part(primary, start(4))
	// SplitButtonPadding is from the control's outer edge, the border drawn inside it, where MyGo's padding starts inside the border: the label starts 11 DIPs in.
	primary.BorderEndWidth(0).MinWidth(35).Padding(5, 0, 6, 0).PaddingStart(10).PaddingEnd(11)
	if pressed {
		primary.TextColor(t.TextMuted)
	} else {
		primary.TextColor(t.Text)
	}
	face, pressed := part(arrow, end(4))
	arrow.BorderStartWidth(0).Width(35).Justify(ui.End).PaddingEnd(11)
	chevron := fluentSecondaryOver(t, face)
	if pressed {
		chevron = t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.447), 0.529))
	}
	nudge := fluentChevronNudge(arrow, pressed)
	arrow.Children(func() {
		ui.Box(c).Size(12, 12).Shrink(0).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			fluentChevronDown(p, r.X+r.W/2, r.Y+r.H/2+nudge, chevron)
		})
	})
}

// fluentSplitAppBar draws a SplitButton of SplitButtonCommandBarStyle, as a CommandBar holds one (SplitButton_themeresources.xaml): its face 36 DIPs high inside AppBarButtonInnerBorderMargin 2,6,2,6, which Toolbar's padding and gap give; the main part padded 12,8,4,8 around its label at FontSize 12, and the secondary part 23 DIPs wide, padded 2,0,8,0, holding U+E972 of Segoe Fluent Icons at FontSize 8, right-aligned with its TextBlock's padding 2,1,2,0; no divider. Both are transparent (AppBarButtonBackground), the part under the pointer SubtleFillColorSecondary (AppBarButtonBackgroundPointerOver) with the other at SplitButtonInAppBarUnfocusedPointerOver (SubtleFillColorTertiary), a part pressed SubtleFillColorTertiary, corners of 4 at the control's ends; the label and glyph TextFillColorPrimary, Secondary pressed (AppBarButtonForegroundPressed). The system focus visual goes 3 DIPs around (FocusVisualMargin -3) the part with the keyboard's focus.
func fluentSplitAppBar(c *ui.Context, primary, arrow ui.Element, start, end func(float32) [4]float32) {
	t := c.Theme()
	hover, down := fluentSubtle(t)
	ph, pp, ah, ap := primary.Hovered(), primary.Pressed(), arrow.Hovered(), arrow.Pressed()
	part := func(b ui.Element, k [4]float32, hovered, pressed, otherHovered bool) ui.Color {
		face, fg := ui.Transparent, t.Text
		switch {
		case pressed:
			face, fg = down, t.TextMuted
		case hovered:
			face = hover
		case otherHovered:
			face = down
		}
		inputsRadius(b.Height(36).Background(face).TextColor(fg), k)
		focused := b.FocusVisible()
		b.DrawOver(func(p *ui.Painter, r ui.Rect) {
			if focused {
				fluentFocusRing(p, t, r, 3, 3, 3, 3, 4)
			}
		})
		return fg
	}
	part(primary, start(4), ph, pp, ah)
	primary.FontSize(12).PaddingY(0).PaddingStart(12).PaddingEnd(4)
	fg := part(arrow, end(4), ah, ap, ph)
	arrow.Width(23).PaddingStart(2).PaddingEnd(8).Justify(ui.End)
	arrow.Children(func() {
		// The TextBlock: the 8-DIP glyph inside its padding of 2 at the sides and 1 above.
		ui.Box(c).Size(12, 9).Shrink(0).Padding(1, 2, 0, 2).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			fluentChevronDown8(p, r.X+r.W/2, r.Y+1+4, fg)
		})
	})
}

// fluentChevronDown8 strokes U+E972 of Segoe Fluent Icons, ChevronDown, at FontSize 8 centered on cx, cy: idiomFluentSubItemChevron's chevron (U+E974 at 8) turned to point down, 4.6 DIPs wide and 2.3 high in a 1-DIP stroke (the shape is inferred: the font is not open).
func fluentChevronDown8(p *ui.Painter, cx, cy float32, col ui.Color) {
	var path ui.Path
	path.MoveTo(cx-2.3, cy-1.15).LineTo(cx, cy+1.15).LineTo(cx+2.3, cy-1.15)
	p.StrokePath(&path, 1, col)
}

// adwSplitButton draws an AdwSplitButton (adw-split-button.c and _buttons.scss of libadwaita 1.10.0): its button and its menu button's button side by side, the button padded 5px 10px (the inner button carries neither .text-button nor .image-button), at least 34 pixels high, bold, and the arrow's button padded 5px 4px around the 16-pixel pan-down-symbolic, 24 pixels wide; each has its own corners at the control's ends ($button_radius, 9), its own state and its focus ring inside its edge. The two overlap by a pixel (margin -1px on either side of the separator): raised, the separator's background is none, and the arrow's translucent face laid over the button's draws the seam where they meet, as linked buttons do (adwLinkedToggle): $button_color, $button_hover_color and $button_active_color (currentColor at 10%, 15% and 30%). Flat, as in a .toolbar (%flat_split_button, _toolbars.scss), the buttons are %button_basic_flat (transparent, $hover_color 7% and $active_color 16%) over the split button, which is $hover_color as a whole while the pointer is over it or a part is pressed, and the separator shows between them, currentColor at $dimmer_opacity (30%) with margins of 6px above and below, but while the pointer is over the split button or a part has the focus from the keyboard.
func adwSplitButton(c *ui.Context, group, primary, arrow ui.Element, start, end func(float32) [4]float32, flat bool) {
	t := c.Theme()
	rtl := c.Direction() == ui.RightToLeft
	shares := func(b ui.Element) float32 {
		switch {
		case b.Pressed():
			return pick(flat, float32(0.30), 0.16)
		case b.Hovered():
			return pick(flat, float32(0.15), 0.07)
		}
		return pick(flat, float32(0.10), 0)
	}
	under := t.Background
	hot := primary.Hovered() || arrow.Hovered() || primary.Pressed() || arrow.Pressed()
	if flat && hot {
		under = adwInkAlpha(t, t.Background, 0.07)
		group.Radius(9).Background(under)
	}
	face := func(share float32, base ui.Color) ui.Color {
		if share == 0 {
			return ui.Transparent
		}
		return adwInkAlpha(t, base, share)
	}
	ps, as := shares(primary), shares(arrow)
	pf := face(ps, under)
	for _, b := range []ui.Element{primary, arrow} {
		b.MinHeight(34).FontWeight(700).Border(0, ui.Color{})
	}
	inputsRadius(primary.Padding(5, 10).Background(pf), start(9)).TextColor(adwTextOn(t, pick(ps == 0, pf, under)))
	adwFocusRing(t, primary, 9, nil)
	af := face(as, under)
	inputsRadius(arrow.MinWidth(24).Padding(5, 4).MarginStart(-1).Background(af), end(9))
	ink := adwTextOn(t, pick(as == 0, af, under))
	focus := primary.FocusVisible() || arrow.FocusVisible()
	edge := func(r ui.Rect) float32 {
		if rtl {
			return r.X + r.W - 1
		}
		return r.X
	}
	adwFocusRing(t, arrow, 9, func(p *ui.Painter, r ui.Rect) {
		switch {
		case !flat:
			// The arrow's face over the button's, a pixel wide.
			p.Fill(ui.Rect{X: edge(r), Y: r.Y, W: 1, H: r.H}, adwInkAlpha(t, pf, as), 0)
		case !hot && !focus:
			p.Fill(ui.Rect{X: edge(r), Y: r.Y + 6, W: 1, H: r.H - 12}, adwInkAlpha(t, under, 0.30), 0)
		}
	})
	arrow.Children(func() { idiomAdwArrow(c, ink) })
}

// ---- icon buttons ----

// IconButton creates a button showing only icon, as the platform draws a button with an image and no title: a bordered push button of macOS 24 points high, 12 points either side of the image (measured from AppKit's rendering on macOS 27.2 by tools/appkit-reference: an NSButton showing a 15-point SF Symbol is 39 points wide); WinUI's Button with a 16-DIP icon as its content, its ButtonPadding 11,5,11,6 (32 high as with a label: inferred, as WinUI sizes a Button to its content); GNOME's button.image-button, 34 pixels square, raised, padded 5 around the 16-pixel icon (_buttons.scss). Inside a Toolbar it is the toolbar's image item, as Button with an empty label draws it: a 36-point image item of an NSToolbar, which joins one beside it into one platter, an AppBarButton with its label collapsed (64 DIPs wide), and a flat image button. The icon is the caller's SVG drawn 16 by 16 in the color of the button's text, which ui.Icon takes from its ancestors whatever colors the SVG has: labelColor over the face on macOS, TextFillColorPrimary on Windows (Secondary pressed), window_fg_color on GNOME.
//
// label names the button for assistive technology and shows as its tooltip, the platform's (Tooltip), without becoming its Description too, which would repeat the name; in a Toolbar's overflow menu it is the item's text. Clicked, the focus, the keys and Disabled are ui.ButtonBase's.
func IconButton(c *ui.Context, icon *ui.SVG, label string) ui.Element {
	_, bar := idiomToolbars.Load(c)
	b := button(c, "", false).Label(label)
	if !bar && PlatformOf(c) == Linux {
		b.MinWidth(34).Padding(5)
	}
	b.Children(func() { ui.Icon(c, icon).Size(16, 16).Shrink(0) })
	inputsTooltip(c, b, label)
	return b
}

// inputsTooltip shows text as anchor's tooltip, the platform's (Tooltip's look), without making it anchor's Description: an icon button's tooltip is its name.
func inputsTooltip(c *ui.Context, anchor ui.Element, text string) {
	if text == "" {
		return
	}
	p := PlatformOf(c)
	ui.TooltipBase(c, anchor, func(tip ui.Element) {
		switch p {
		case Windows:
			fluentTooltip(c, anchor, tip)
		case Linux:
			adwTooltip(c, anchor, tip)
		default:
			macTooltip(c, anchor, tip)
		}
		tip.Children(func() { ui.Text(c, text) })
	})
}
