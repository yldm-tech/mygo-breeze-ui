package breeze

import "github.com/egoist/mygo/ui"

// Sheet shows a dialog built by fn over the window while *open, as the platform draws a dialog holding a form: an NSWindow sheet on macOS, a ContentDialog on Windows and a floating AdwDialog on GNOME. Unlike Modal, a click outside it does nothing, as none of the three platforms closes such a dialog on one, so what was typed in it is not lost to a stray click; Escape sets *open to false, which Dismissed reports. It returns the panel, which takes more styles as Modal's does.
func Sheet(c *ui.Context, open *bool, fn func()) ui.Element {
	if !*open {
		return ui.Element{}
	}
	p := PlatformOf(c)
	w, h := c.Size()
	// A layer under the dialog takes the clicks outside its panel and does nothing with them: the dialog's backdrop lets the pointer through to it, so the backdrop is never clicked and only Escape dismisses the dialog.
	ui.Overlay(c, func() { ui.Box(c).Absolute().Left(0).Top(0).Right(0).Bottom(0) })
	return ui.DialogBase(c, open, func(back, panel ui.Element) {
		back.PassThrough()
		panel.MaxWidth(w - 40).MaxHeight(h - 40)
		Preview(c, p, func() {
			switch p {
			case Windows:
				fluentDialog(c, back, panel)
			case Linux:
				adwDialog(c, back, panel)
			default:
				macSheetPanel(c, back, panel)
			}
			fn()
		})
	})
}

// DestructiveButton creates a push button for an action that destroys data, as the platform marks one: on macOS NSAlert's destructive button (hasDestructiveAction), a face tinted with systemRed under a systemRed label; on GNOME a .destructive-action button, filled with destructive_bg_color under a white label; on Windows, which has no destructive style, a plain Button. Ask Clicked whether it was clicked.
func DestructiveButton(c *ui.Context, label string) ui.Element {
	b := ui.ButtonBase(c)
	switch PlatformOf(c) {
	case Windows:
		fluentButton(c, b, false)
	case Linux:
		adwDestructiveButton(c, b)
	default:
		macDestructiveButton(c, b)
	}
	if label != "" {
		b.Children(func() { ui.Text(c, label).SingleLine() })
	}
	return b
}

// DestructiveAlertDialog shows an alert as AlertDialog does, with buttons[destructive] drawn as the platform draws a destructive button and never the default: the default, which has the focus and which Enter clicks, is the last of the other buttons, so Enter never destroys anything. With the destructive button alone there is no default: the alert itself has the focus and Enter clicks nothing. The buttons keep AlertDialog's places. macOS's NSAlert, which then has no default button, draws the default in the plain face, as do GNOME's AdwAlertDialog, where it is only the default response; Windows's ContentDialog draws it in the accent, as its DefaultButton, and the destructive button plain. A destructive that is not the index of a button draws none destructive, and the alert is AlertDialog's. It returns what AlertDialog returns: the index of the button clicked, ui.AlertDismissed when Escape closes the alert, and -1 otherwise.
func DestructiveAlertDialog(c *ui.Context, open *bool, title, message string, destructive int, buttons ...string) int {
	if destructive >= len(buttons) {
		destructive = -1
	}
	return alertDialog(c, open, title, message, destructive, buttons)
}

// alertDialog shows the alert of AlertDialog, or of DestructiveAlertDialog with buttons[destructive] drawn destructive; destructive is -1 for none.
func alertDialog(c *ui.Context, open *bool, title, message string, destructive int, buttons []string) int {
	if !*open {
		return -1
	}
	p := PlatformOf(c)
	w, _ := c.Size()
	n := len(buttons)
	// The default is the last button that is not destructive; the button drawn in its own style, which the GNOME alert does not give its response face, is the destructive one, or else the default.
	def := n - 1
	if def == destructive {
		def--
	}
	styled := destructive
	if destructive < 0 {
		styled = def
	}
	chosen, dismissed := -1, false
	ui.Overlay(c, func() {
		back := ui.Box(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).Center().Modal()
		dismissed = back.OverlayShortcut(0, ui.KeyEscape)
		back.Children(func() {
			panel := ui.Column(c).MaxWidth(w - 40).Role(ui.RoleAlertDialog).Label(title).Description(message)
			if def < 0 && n > 0 {
				// The destructive button is alone, so no button is the default: the panel takes the focus the dialog would give its first button, the destructive one, so Enter clicks nothing.
				panel.Focusable().FocusRing(false).AutoFocus()
			}
			// button builds button i of the alert.
			button := func(i int) ui.Element {
				var b ui.Element
				switch {
				case i == destructive:
					b = DestructiveButton(c, buttons[i])
				case i == def && (destructive < 0 || p == Windows):
					b = PrimaryButton(c, buttons[i])
				default:
					b = Button(c, buttons[i])
				}
				if i == def {
					b.AutoFocus()
				}
				if b.Clicked() {
					chosen = i
				}
				return b
			}
			Preview(c, p, func() {
				switch p {
				case Windows:
					fluentAlert(c, back, panel, title, message, n, button)
				case Linux:
					adwAlert(c, back, panel, title, message, n, styled, button)
				default:
					macAlert(c, back, panel, title, message, n, button)
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

// macSheetPanel draws an NSWindow sheet as macOS 27 shows one, measured from a live capture on macOS 27.2 converted to sRGB (a sheet's window does not render offscreen): centered on its window, which it dims with black at 20% in Aqua and 23% in Dark Aqua (255 to 204, and 35 to 27 over the window), filled with the window's background, its corners those of MyGo's radius of 16, the closest to the capture's 15-point arc. The capture leaves the window's shadow out, which is macDialog's. The padding of 20 and the gap of 12 are macDialog's, for the content: a sheet's window has none of its own.
func macSheetPanel(c *ui.Context, back, panel ui.Element) {
	t := c.Theme()
	back.Background(black.Alpha(pick(t.Dark, float32(0.20), 0.23)))
	panel.Padding(20).Gap(12).Radius(16).Background(t.Background).
		Shadow(0, 12, 36, 0, ui.RGBA(0, 0, 0, pick(t.Dark, float32(0.22), 0.5)))
}

// macDestructiveShare is how much systemRed tints the face of NSAlert's destructive button, measured from a live capture of an alert on macOS 27.2 converted to sRGB: (232, 182, 181) over the alert's 220 in Aqua and (100, 54, 54) over its 47 in Dark Aqua, which the theme's Danger at these shares gives within 4 levels.
func macDestructiveShare(dark bool) float32 { return pick(dark, float32(0.24), 0.255) }

// macDestructiveButton draws NSAlert's destructive button (hasDestructiveAction), the only destructive look AppKit draws (a push button with hasDestructiveAction elsewhere, and SwiftUI's bordered Button(role: .destructive), stay plain, rendered on macOS 27.2): a push button's shape, its face the window tinted with systemRed (the theme's Danger) and its label systemRed. Pressed, the tint doubles (inferred: the capture shows no pressed button).
func macDestructiveButton(c *ui.Context, b ui.Element) {
	t := c.Theme()
	b.MinHeight(24).PaddingX(12).Gap(6).Radius(6).Justify(ui.Center).AlignItems(ui.Center)
	share := macDestructiveShare(t.Dark)
	if b.Pressed() {
		share *= 2
	}
	b.Background(t.Background.Mix(t.Danger, share)).TextColor(t.Danger)
	macRing(b, 6)
}

// adwDestructiveButton draws a GtkButton with a label and .destructive-action: adwButton's shape, destructive_bg_color (the theme's Danger, #e01b24 in the light theme and #c01c28 in the dark one) with destructive_fg_color, white, and under the pointer and pressed the 10% currentColor and 20% shade layers %opaque_button lays over it, as over the suggested action (inferred from _buttons.scss: libadwaita 1.10 does not run here). Its focus ring, outside with a gap, is MyGo's, as the suggested action's.
func adwDestructiveButton(c *ui.Context, b ui.Element) {
	t := c.Theme()
	b.MinHeight(34).Padding(5, 17).Gap(4).Radius(9).FontWeight(700).Justify(ui.Center).AlignItems(ui.Center)
	fill := t.Danger
	switch {
	case b.Pressed():
		fill = fill.Mix(adwInk, 0.2)
	case b.Hovered():
		fill = fill.Mix(white, 0.1)
	}
	b.Background(fill).TextColor(white)
}
