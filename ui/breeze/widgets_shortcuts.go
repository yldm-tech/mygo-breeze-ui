package breeze

import (
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"
)

// ---- shortcuts: keys as the platforms write and draw them ----
//
// macOS writes a shortcut as its menus show key equivalents: the modifiers' glyphs in the order ⌃⌥⇧⌘, then the key's glyph or name (NSMenuItem's keyEquivalentModifierMask, as MyGo's popover menus and GTK's macOS backend write it, gtkaccelgroup.c of GTK 4.24.0). Windows writes it as WinUI's KeyboardAccelerator does for a MenuFlyoutItem's KeyboardAcceleratorTextLabel: Ctrl, Alt, Windows and Shift in that order, then the key, joined by KEYBOARD_ACCELERATOR_TEXT_JOIN, "%s+%s", the keys named by the TEXT_VK_ strings of Microsoft.UI.Xaml.Common.rc (KeyboardAccelerator_Partial.cpp of microsoft-ui-xaml main). GNOME writes it as gtk_accelerator_get_label of GTK 4.24.0 does for a menu's accelerator: Shift, Ctrl, Alt and Super in that order, joined by +, the key a character in upper case, Space or Backslash, or else GDK's keyval name without its underscores. Its keycaps are AdwShortcutLabel's of libadwaita 1.10.0 (adw-shortcut-label.c), which orders the modifiers Super, Ctrl, Alt and Shift and draws arrows, Space (␣) and Return (⏎) as symbols.
//
// AppKit has no keycap control, so the macOS keycap is inferred: one rounded plate holding the whole shortcut, as the cheat sheets and settings of macOS apps draw it. The Windows keycap is the KeyVisual of PowerToys (Common.UI.Controls), as its Keyboard Manager draws keys with OriginalKeyVisualStyle: a plate per key, 4 DIPs apart.

// shortcutCmd returns mods with the modifier ui.Cmd stands for on the host, Super on macOS and Ctrl elsewhere, replaced by p's: a shortcut written with ui.Cmd shows as the platform drawn writes it, so a Preview of Windows on a Mac shows Ctrl+S where macOS shows ⌘S.
func shortcutCmd(p Platform, mods ui.Modifiers) ui.Modifiers {
	if mods&ui.Cmd == 0 {
		return mods
	}
	return mods&^ui.Cmd | pick(p == Mac, ui.Ctrl, ui.Super)
}

// shortcutKeyNames is how each platform names the keys that are not a character: macOS's glyphs (NSMenu's and Apple's style guide's, Space and the function keys spelled out), WinUI's TEXT_VK_ strings, GTK's keyval names in its menus (gtk_accelerator_get_label, without underscores) and in AdwShortcutLabel's keycaps (get_labels, which draws arrows, Space and Return as symbols and spells Page_Up and Page_Down with an underscore; Breeze UI spells them with a space as GTK's menus do).
var shortcutKeyNames = map[ui.Key][4]string{
	//                 macOS  Windows      GNOME text   GNOME keycap
	ui.KeyEnter:       {"↩", "Enter", "Return", "⏎"},
	ui.KeyEscape:      {"⎋", "Esc", "Escape", "Escape"},
	ui.KeyBackspace:   {"⌫", "Back", "BackSpace", "BackSpace"},
	ui.KeyTab:         {"⇥", "Tab", "Tab", "Tab"},
	ui.KeySpace:       {"Space", "Space", "Space", "␣"},
	ui.KeyDelete:      {"⌦", "Delete", "Delete", "Delete"},
	ui.KeyInsert:      {"Insert", "Insert", "Insert", "Insert"},
	ui.KeyHome:        {"↖", "Home", "Home", "Home"},
	ui.KeyEnd:         {"↘", "End", "End", "End"},
	ui.KeyPageUp:      {"⇞", "PageUp", "Page Up", "Page Up"},
	ui.KeyPageDown:    {"⇟", "PageDown", "Page Down", "Page Down"},
	ui.KeyLeft:        {"←", "Left", "Left", "←"},
	ui.KeyRight:       {"→", "Right", "Right", "→"},
	ui.KeyUp:          {"↑", "Up", "Up", "↑"},
	ui.KeyDown:        {"↓", "Down", "Down", "↓"},
	ui.KeyContextMenu: {"Menu", "Application", "Menu", "Menu"},
	ui.KeyBackslash:   {"\\", "\\", "Backslash", "Backslash"},
}

// shortcutChars are the characters of the punctuation keys, as a US layout types them: WinUI names these keys by the character the user's layout types (MapVirtualKeyW), and macOS and GTK by the character of the key equivalent or keyval.
var shortcutChars = map[ui.Key]string{
	ui.KeyMinus: "-", ui.KeyEqual: "=", ui.KeyComma: ",", ui.KeyPeriod: ".", ui.KeySlash: "/",
	ui.KeySemicolon: ";", ui.KeyQuote: "'", ui.KeyBracketLeft: "[", ui.KeyBracketRight: "]", ui.KeyBackquote: "`",
}

// shortcutKey returns how p names key, in its text when keycap is false and on its keycaps when true (which differ on GNOME only), or "" for a key it does not name, as KeyBack and KeyForward.
func shortcutKey(p Platform, key ui.Key, keycap bool) string {
	switch {
	case key >= ui.KeyA && key <= ui.KeyZ:
		return string(rune('A' + key - ui.KeyA))
	case key >= ui.Key0 && key <= ui.Key9:
		return string(rune('0' + key - ui.Key0))
	case key >= ui.KeyF1 && key <= ui.KeyF12:
		return "F" + strconv.Itoa(int(key-ui.KeyF1)+1)
	}
	if s, ok := shortcutChars[key]; ok {
		return s
	}
	names, ok := shortcutKeyNames[key]
	if !ok {
		return ""
	}
	switch p {
	case Windows:
		return names[1]
	case Linux:
		return names[pick(keycap, 2, 3)]
	}
	return names[0]
}

// shortcutModifier is a modifier as each platform names it.
type shortcutModifier struct {
	mod                            ui.Modifiers
	mac, windows, gnome, spokenMac string
}

var (
	modCtrl  = shortcutModifier{ui.Ctrl, "⌃", "Ctrl", "Ctrl", "Control"}
	modAlt   = shortcutModifier{ui.Alt, "⌥", "Alt", "Alt", "Option"}
	modShift = shortcutModifier{ui.Shift, "⇧", "Shift", "Shift", "Shift"}
	modSuper = shortcutModifier{ui.Super, "⌘", "Windows", "Super", "Command"}
)

// shortcutOrders are the orders the platforms put the modifiers in: macOS's menus ⌃⌥⇧⌘, WinUI's Ctrl, Alt, Windows and Shift (KeyboardAccelerator::GetStringRepresentation), GTK's menus Shift, Ctrl, Alt and Super (gtk_accelerator_print_label), AdwShortcutLabel's keycaps Super, Ctrl, Alt and Shift (get_labels).
var (
	macOrder         = []shortcutModifier{modCtrl, modAlt, modShift, modSuper}
	windowsOrder     = []shortcutModifier{modCtrl, modAlt, modSuper, modShift}
	gnomeTextOrder   = []shortcutModifier{modShift, modCtrl, modAlt, modSuper}
	gnomeKeycapOrder = []shortcutModifier{modSuper, modCtrl, modAlt, modShift}
)

// shortcutParts returns the keys of mods+key as p shows them on keycaps when keycap is true, or in its text when false, the modifiers first; nil when p does not name key.
func shortcutParts(p Platform, mods ui.Modifiers, key ui.Key, keycap bool) []string {
	name := shortcutKey(p, key, keycap)
	if name == "" {
		return nil
	}
	mods = shortcutCmd(p, mods)
	order := macOrder
	switch {
	case p == Windows:
		order = windowsOrder
	case p == Linux && keycap:
		order = gnomeKeycapOrder
	case p == Linux:
		order = gnomeTextOrder
	}
	var parts []string
	for _, m := range order {
		if mods&m.mod == 0 {
			continue
		}
		switch p {
		case Windows:
			parts = append(parts, m.windows)
		case Linux:
			parts = append(parts, m.gnome)
		default:
			parts = append(parts, m.mac)
		}
	}
	return append(parts, name)
}

// shortcutText writes mods+key as p's menus do: macOS's glyphs run together, the other platforms' names joined by +.
func shortcutText(p Platform, mods ui.Modifiers, key ui.Key) string {
	parts := shortcutParts(p, mods, key, false)
	if p == Mac {
		return strings.Join(parts, "")
	}
	return strings.Join(parts, "+")
}

// macSpokenKeys are the names VoiceOver speaks for the keys macOS writes as glyphs.
var macSpokenKeys = map[ui.Key]string{
	ui.KeyEnter: "Return", ui.KeyEscape: "Escape", ui.KeyBackspace: "Delete", ui.KeyTab: "Tab", ui.KeyDelete: "Forward Delete",
	ui.KeyHome: "Home", ui.KeyEnd: "End", ui.KeyPageUp: "Page Up", ui.KeyPageDown: "Page Down",
	ui.KeyLeft: "Left Arrow", ui.KeyRight: "Right Arrow", ui.KeyUp: "Up Arrow", ui.KeyDown: "Down Arrow",
}

// shortcutSpoken returns the name assistive technology gives mods+key on p: on macOS the modifiers and the key in words, in the order of the glyphs, as VoiceOver reads a menu's key equivalent ("Shift Command L"); elsewhere the text, as WinUI's KeyboardAccelerator and AdwShortcutLabel (gtk_accelerator_get_label) name it.
func shortcutSpoken(p Platform, mods ui.Modifiers, key ui.Key) string {
	if p != Mac {
		return shortcutText(p, mods, key)
	}
	name := shortcutKey(p, key, false)
	if name == "" {
		return ""
	}
	if s, ok := macSpokenKeys[key]; ok {
		name = s
	}
	mods = shortcutCmd(p, mods)
	var words []string
	for _, m := range macOrder {
		if mods&m.mod != 0 {
			words = append(words, m.spokenMac)
		}
	}
	return strings.Join(append(words, name), " ")
}

// ShortcutText returns the shortcut of mods and key as the platform's menus write it, for text after a command such as a menu-like row's or a tooltip's: ⇧⌘L on macOS, Ctrl+Shift+L on Windows (WinUI's KeyboardAccelerator text) and Shift+Ctrl+L on GNOME (GTK's accelerator label) for ui.Cmd|ui.Shift with ui.KeyL. ui.Cmd stands for the command key of the platform drawn, Command on macOS and Ctrl elsewhere, so a Preview of another platform writes that platform's shortcut. A key the platform does not name, as ui.KeyBack, gives "".
func ShortcutText(c *ui.Context, mods ui.Modifiers, key ui.Key) string {
	return shortcutText(PlatformOf(c), mods, key)
}

// ShortcutLabel creates the keycaps of the shortcut of mods and key, as the platform shows a shortcut in a settings page or a cheat sheet: on macOS one keycap holding the shortcut as its menus write it (⇧⌘L; inferred, as AppKit has no keycap control), on Windows a keycap per key 4 DIPs apart (Ctrl, Shift, L; the KeyVisual of PowerToys), on GNOME AdwShortcutLabel's keycaps 6 pixels apart (Ctrl, Shift, L). ui.Cmd stands for the platform's command key, as in ShortcutText. Assistive technology sees one text named by the shortcut, Shift Command L on macOS and the text elsewhere; the keycaps stay left to right in a right-to-left window, as AdwShortcutLabel keeps them. Place several side by side for alternatives.
func ShortcutLabel(c *ui.Context, mods ui.Modifiers, key ui.Key) ui.Element {
	p := PlatformOf(c)
	keys := shortcutParts(p, mods, key, true)
	if p == Mac && keys != nil {
		keys = []string{strings.Join(keys, "")}
	}
	return keycaps(c, p, keys).Label(shortcutSpoken(p, mods, key))
}

// Keycaps creates a keycap for each of keys, as ShortcutLabel draws them on the platform, for keys that are not one shortcut: a chord such as Keycaps(c, "Z", "X"), a range, or keys of one's own. The keys are shown as given, each on its own keycap on every platform, macOS's too. Assistive technology sees one text named by the keys, separated by spaces.
func Keycaps(c *ui.Context, keys ...string) ui.Element {
	return keycaps(c, PlatformOf(c), keys).Label(strings.Join(keys, " "))
}

// keycaps builds a row of a keycap for each of keys, drawn as p draws them:
//   - macOS (inferred): 22 points high and at least as wide, padded 8 either side, corners of 6, the text 12 points, on a plate of black at 5% (white at 12% in the dark theme) with a 1-point line inside its bottom of black at 8% (30%), keycaps 4 points apart.
//   - Windows: PowerToys's KeyVisual in Keyboard Manager's OriginalKeyVisualStyle: padded 6, a 1-DIP border of CardStrokeColorDefault, ControlFillColorDefault, corners of ControlCornerRadius (4), the text 12 DIPs at normal weight in TextFillColorPrimary, keycaps 4 DIPs apart (its StackPanel's Spacing).
//   - GNOME: AdwShortcutLabel's .keycap (_shortcuts-dialog.scss of libadwaita 1.10.0): padded 6, at least 20 wide, corners of 6, currentColor at 10% with an inset 2-pixel shade along its bottom of --card-shade-color (the shade at 7%, 36% in the dark theme) and under prefers-contrast: more a 1-pixel ring of --border-color (currentColor at 50%), the text the label's, keycaps 6 pixels apart (border-spacing).
func keycaps(c *ui.Context, p Platform, keys []string) ui.Element {
	t := c.Theme()
	row := ui.Row(c).Direction(ui.LeftToRight).Role(ui.RoleText).AlignItems(ui.Center).Shrink(0)
	row.Gap(pick(p == Linux, float32(4), 6))
	row.Children(func() {
		for _, k := range keys {
			plate := ui.Row(c).Justify(ui.Center).AlignItems(ui.Center).Shrink(0)
			switch p {
			case Windows:
				_, _, stroke := fluentExpanderCard(t)
				plate.Padding(6).Radius(4).Background(t.Surface).Border(1, stroke)
			case Linux:
				plate.Padding(6).MinWidth(20).Radius(6).Background(adwInkLayer(t, 0.10)).InsetShadow(0, -2, 0, 0, ui.RGBA(0, 0, 6, pick(t.Dark, float32(0.07), 0.36)))
				if c.Preferences().HighContrast {
					plate.Shadow(0, 0, 0, 1, adwInkLayer(t, 0.5))
				}
			default:
				plate.MinHeight(22).MinWidth(22).Padding(0, 8).Radius(6).Background(pick(t.Dark, black.Alpha(0.05), white.Alpha(0.12))).InsetShadow(0, -1, 0, 0, black.Alpha(pick(t.Dark, float32(0.08), 0.3)))
			}
			plate.Children(func() {
				text := ui.Text(c, k).SingleLine().TextColor(t.Text).Role(ui.RoleNone)
				switch p {
				case Windows:
					text.FontSize(t.Rem(12.0 / 14)).FontWeight(400)
				case Mac:
					text.FontSize(t.Rem(12.0 / 13))
				}
			})
		}
	})
	return row
}
