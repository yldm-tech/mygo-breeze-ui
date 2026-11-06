package breeze

import "github.com/egoist/mygo/ui"

// Buttons without a face at rest, which every platform has beside its push buttons: macOS's NSButton that shows its bezel only while the pointer is inside it (showsBorderOnlyWhileMouseInside, of bezel style accessoryBarAction, or accessoryBar for a toggle), WinUI 3's SubtleButtonStyle and libadwaita's .flat button. They are ui.ButtonBase's (ui.ToggleBase's for FlatToggle), so Clicked, Menu, PopoverMenuButton, the focus, the keys and Disabled are MyGo's, and they show what is behind them: their fills are translucent layers, as the platforms' are, so a flat button over a card, a sidebar or a dark player's controls takes its colors from what it is laid over.
//
// Inside a Toolbar or a HeaderBar they are the toolbar's own items, as Button, IconButton and Toggle draw them there.

// FlatButton creates a button showing label without a face at rest, as the platform draws a borderless text button: on macOS an NSButton of bezel style accessoryBarAction showing its bezel only while the pointer is inside it, 24 points high and padded 12 either side of its title, the title in secondaryLabelColor at rest and labelColor over the bezel; on Windows a Button of SubtleButtonStyle, a Button's size (32 DIPs high, ButtonPadding 11,5,11,6 inside its 1-DIP border) with no fill or border at rest; on GNOME a .flat text button, 34 pixels high, padded 5px 17px, bold (flatStyle). Inside a Toolbar it is Button. Ask Clicked whether it was clicked; Menu and PopoverMenuButton open a menu from it.
func FlatButton(c *ui.Context, label string) ui.Element {
	if _, bar := idiomToolbars.Load(c); bar {
		return Button(c, label)
	}
	b := ui.ButtonBase(c)
	flatStyle(c, b, false, false)
	b.Children(func() { ui.Text(c, label).SingleLine() })
	return b
}

// FlatIconButton creates a button showing only icon without a face at rest, as the platform draws a borderless image button: on macOS an NSButton of bezel style accessoryBarAction in a 24-point square that shows its bezel only while the pointer is inside it; on Windows a Button of SubtleButtonStyle holding a 16-DIP icon, 40 by 29 as IconButton; on GNOME a .flat image button, 34 pixels square (flatStyle). The icon is the caller's SVG drawn 16 by 16 in the color of the button's text, whatever colors the SVG has: secondaryLabelColor at rest and labelColor over the bezel on macOS, TextFillColorPrimary on Windows (Secondary pressed), window_fg_color on GNOME. Inside a Toolbar it is IconButton, the toolbar's image item.
//
// label names the button for assistive technology and shows as its tooltip, the platform's, as IconButton's does. Ask Clicked whether it was clicked; Menu and PopoverMenuButton open a menu from it, as from a ⋯ button.
func FlatIconButton(c *ui.Context, icon *ui.SVG, label string) ui.Element {
	if _, bar := idiomToolbars.Load(c); bar {
		return IconButton(c, icon, label)
	}
	b := ui.ButtonBase(c).Label(label)
	flatStyle(c, b, true, false)
	b.Children(func() { ui.Icon(c, icon).Size(16, 16).Shrink(0) })
	if label != "" {
		platformTooltip(c, b, label)
	}
	return b
}

// FlatToggle creates a button that stays pressed while *on, showing icon without a face while off, as the platform draws a borderless toggle: on macOS an NSButton of bezel style accessoryBar and type pushOnPushOff in a 24-point square, which shows its bezel only while the pointer is inside it and, on but away from the pointer, a gray fill (black at 12.5% in Aqua, white in Dark Aqua) under its image in labelColor, the accent under the pointer; on Windows a ToggleButton of SubtleButtonStyle's look while off and a checked ToggleButton's, the accent, while on (WinUI has no subtle toggle; inferred); on GNOME a .flat image toggle, currentColor at 10% while checked (flatStyle). A nil icon shows label as its text instead, at the size of FlatButton. Inside a Toolbar it is Toggle, the toolbar's image toggle.
//
// label names it for assistive technology and, with an icon, shows as its tooltip. A click or Space while it has the focus turns *on over and Changed reports it; assistive technology hears a toggle button, pressed while *on, as ui.ToggleBase tells it.
func FlatToggle(c *ui.Context, on *bool, icon *ui.SVG, label string) ui.Element {
	if _, bar := idiomToolbars.Load(c); bar {
		if icon == nil {
			return Toggle(c, on, label)
		}
		b := Toggle(c, on, "").Label(label)
		b.Children(func() { ui.Icon(c, icon).Size(16, 16).Shrink(0) })
		if label != "" {
			platformTooltip(c, b, label)
		}
		return b
	}
	b := ui.ToggleBase(c, on).Label(label)
	flatStyle(c, b, icon != nil, *on)
	if icon == nil {
		b.Children(func() { ui.Text(c, label).SingleLine() })
		return b
	}
	b.Children(func() { ui.Icon(c, icon).Size(16, 16).Shrink(0) })
	if label != "" {
		platformTooltip(c, b, label)
	}
	return b
}

// flatStyle draws b, a flat button showing an icon (image) or text, toggled on when on, as the platform draws it.
func flatStyle(c *ui.Context, b ui.Element, image, on bool) {
	switch PlatformOf(c) {
	case Windows:
		fluentSubtleButton(c, b, image, on)
	case Linux:
		adwFlatButton(c, b, image, on)
	default:
		macFlatButton(c, b, image, on)
	}
}

// Shares of black (Aqua) or white (Dark Aqua) of a flat button's bezel, measured from AppKit's rendering on macOS 27.2 at 2x by tools/appkit-reference (appkit-flat-*.png): the bezel of an accessoryBarAction button is a push button's face (macFaceAlpha, 235 over white and 48 over the window's 30 in Dark Aqua) and pressed a pressed push button's (macPressedAlpha, 212 and 68); an accessoryBar toggle on with its bezel hidden is filled with 12.5% (223 and 58).
const macFlatOnAlpha = 0.125

// macFlatLayer is black at alpha in Aqua and white in Dark Aqua, translucent: a flat button's fill, which shows what is behind it.
func macFlatLayer(t *ui.Theme, alpha float32) ui.Color {
	return pick(t.Dark, black, white).Alpha(alpha)
}

// macFlatButton draws an NSButton that shows its bezel only while the pointer is inside it (showsBorderOnlyWhileMouseInside), as AppKit renders one on macOS 27.2 at 2x (tools/appkit-reference, appkit-flat-*.png): of bezel style accessoryBarAction, 24 points high, a text button padded 12 points either side of its title as a push button (51 points around "Skip"), an image button in a 24-point square frame, which AppKit fills whole with the bezel, corners of 6 as a push button's. At rest it has no face, its title or image in secondaryLabelColor (128 over white, 154 over the window in Dark Aqua); under the pointer the bezel shows, a push button's face, black at 7.8% in Aqua and white in Dark Aqua (235 and 48), its content in labelColor over it (36 and 223), and pressed a pressed push button's, 16.9% (212 and 68) with labelColor over it (32 and 226). A toggle (accessoryBar, pushOnPushOff) on and away from the pointer is filled at 12.5% (223 and 58) under labelColor (34 and 225); under the pointer its bezel shows on, the accent with a white image (#0078f9 and #067eff, t.Accent here, as for Toggle), and pressed the accent 10% darker in Aqua and lighter in Dark Aqua (#006ce1 and #1e8bff). AppKit cannot put the pointer inside a button offscreen: the bezel under the pointer is the one it draws with showsBorderOnlyWhileMouseInside false (inferred). AppKit's focus ring goes around it while it has the focus from the keyboard (inferred, as on a push button). Disabled, MyGo halves its opacity, the image's darkest pixel at 191 over white where AppKit's disabled one reads 189.
func macFlatButton(c *ui.Context, b ui.Element, image, on bool) {
	t := c.Theme()
	b.MinHeight(24).Radius(6).Gap(6).Justify(ui.Center).AlignItems(ui.Center)
	if image {
		b.Size(24, 24).Padding(0)
	} else {
		b.PaddingX(12)
	}
	pressed, hovered := b.Pressed(), b.Hovered()
	ink := pick(t.Dark, black, white)
	face, fg := ui.Transparent, t.TextMuted
	// label is labelColor over a fill of alpha laid over the window.
	label := func(alpha float32) ui.Color { return tint(tint(t.Background, t.Dark, alpha), t.Dark, 0.847) }
	switch {
	case on && pressed:
		face, fg = t.Accent.Mix(ink, 0.1), t.AccentText
	case on && hovered:
		face, fg = t.Accent, t.AccentText
	case on:
		face, fg = macFlatLayer(t, macFlatOnAlpha), label(macFlatOnAlpha)
	case pressed:
		face, fg = macFlatLayer(t, macPressedAlpha), label(macPressedAlpha)
	case hovered:
		face, fg = macFlatLayer(t, macFaceAlpha), label(macFaceAlpha)
	}
	b.Background(face).TextColor(fg)
	macRing(b, 6)
}

// fluentSubtleFill returns SubtleFillColorSecondary and SubtleFillColorTertiary as the translucent colors they are, #09000000 and #06000000 in the light theme, #0FFFFFFF and #0AFFFFFF in the dark one (Common_themeresources_any.xaml), where fluentSubtle lays them over the window.
func fluentSubtleFill(t *ui.Theme) (hover, pressed ui.Color) {
	if t.Dark {
		return white.Alpha(0x0f / 255.0), white.Alpha(0x0a / 255.0)
	}
	return black.Alpha(0x09 / 255.0), black.Alpha(0x06 / 255.0)
}

// fluentSubtleButton draws a Button of SubtleButtonStyle (Button_themeresources.xaml of microsoft-ui-xaml main at 04a5ab5): DefaultButtonStyle's ButtonPadding 11,5,11,6 inside a 1-DIP border, ControlCornerRadius (4) and FocusVisualMargin -3, so as large as a Button, 32 DIPs high around a label (fluentButton) and 40 by 29 around a 16-DIP icon (IconButton, measured from WinUI 3 3.3); its background and border SubtleFillColorTransparent at rest, SubtleFillColorSecondary under the pointer and SubtleFillColorTertiary pressed (SubtleButtonBackground*, SubtleButtonBorderBrush*, the border the same brush over the background's InnerBorderEdge, so the face is one fill), the content TextFillColorPrimary, TextFillColorSecondary pressed (SubtleButtonForeground*). The faces are those brushes, translucent, so the button shows what is behind it. WinUI has no subtle toggle: on, it is drawn as a checked ToggleButton, the accent button's look (fluentToggle; inferred).
func fluentSubtleButton(c *ui.Context, b ui.Element, image, on bool) {
	t := c.Theme()
	if on {
		fluentButton(c, b, true)
	} else {
		b.MinHeight(32).Padding(5, 11, 6, 11).Gap(8).Radius(4).Justify(ui.Center).AlignItems(ui.Center).FocusRing(false)
		pressed := b.Pressed()
		hover, down := fluentSubtleFill(t)
		face, fg := ui.Transparent, t.Text
		switch {
		case pressed:
			face, fg = down, t.TextMuted
		case b.Hovered():
			face = hover
		}
		// The border keeps its DIP, filled as the face is, so the content does not move.
		b.Background(face).TextColor(fg).Border(1, ui.Color{})
		focused := b.FocusVisible()
		b.DrawOver(func(p *ui.Painter, r ui.Rect) {
			if focused {
				fluentFocusRing(p, t, r, 3, 3, 3, 3, 4)
			}
		})
	}
	if image {
		b.MinHeight(29)
	}
}

// adwFlatButton draws a GtkButton of libadwaita 1.10's .flat (%button_basic_flat in _buttons.scss): a .text-button 34 pixels high (min-height 24, padding 5px 17px), bold, or an .image-button 34 pixels square (min-width 24, padding 5), corners of $button_radius (9); transparent at rest, $hover_color (currentColor at 7%) under the pointer and $active_color (16%) pressed; checked, $selected_color (10%), $selected_hover_color (13%) under the pointer and $selected_active_color (19%) pressed (_colors.scss), as idiomAdwFlatFace fills a toolbar's flat buttons, here as translucent layers over what is behind the button; the text window_fg_color over the face; the focus ring 2 pixels inside its edge.
func adwFlatButton(c *ui.Context, b ui.Element, image, on bool) {
	t := c.Theme()
	if image {
		b.MinWidth(34).MinHeight(34).Padding(5)
	} else {
		b.MinHeight(34).Padding(5, 17).FontWeight(700)
	}
	b.Radius(9).Gap(4).Justify(ui.Center).AlignItems(ui.Center).Border(0, ui.Color{})
	share := float32(0)
	switch {
	case b.Pressed():
		share = pick(on, float32(0.16), 0.19)
	case b.Hovered():
		share = pick(on, float32(0.07), 0.13)
	case on:
		share = 0.10
	}
	face := ui.Transparent
	if share > 0 {
		face = adwInkLayer(t, share)
	}
	b.Background(face).TextColor(adwTextOn(t, adwInkAlpha(t, t.Background, share)))
	adwFocusRing(t, b, 9, nil)
}
