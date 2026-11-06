package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestShortcutStoriesFollowTheirControls checks that the shortcut stories draw each platform's shortcut and follow their controls: the ShortcutLabel story shows ⇧⌘T in one keycap on macOS and Ctrl, Shift and T on Windows and GNOME, whatever the host, the key and modifier controls change the keys and the Go shown, and the ShortcutText story writes the platform's text.
func TestShortcutStoriesFollowTheirControls(t *testing.T) {
	b := newBook()
	b.selected = "ShortcutLabel/Default"
	tester := ui.NewTester(b.view, 1280, 820)
	a := b.argsOf(b.story())
	keys := map[int][]string{0: {"⇧⌘T"}, 1: {"Ctrl", "Shift", "T"}, 2: {"Ctrl", "Shift", "T"}}
	for platform := range 3 {
		b.platform = platform
		tester.Frame()
		for _, k := range keys[platform] {
			if !tester.HasText(k) {
				t.Errorf("%s: the keycap %q is missing: %q", platformNames[platform], k, tester.Texts())
			}
		}
	}
	if !tester.HasText("breeze.ShortcutLabel(c, ui.Cmd|ui.Shift, ui.KeyT)") {
		t.Errorf("the code does not name the shortcut: %q", tester.Texts())
	}
	*a.on["shift"], *a.on["alt"], *a.text["key"] = false, true, "Left"
	b.platform = 0
	tester.Frame()
	if !tester.HasText("⌥⌘←") || !tester.HasText("breeze.ShortcutLabel(c, ui.Cmd|ui.Alt, ui.KeyLeft)") {
		t.Errorf("the controls did not change the shortcut: %q", tester.Texts())
	}

	b.selected = "ShortcutText/Menu row"
	tester.Frame()
	texts := map[int]string{0: "⇧⌘T", 1: "Ctrl+Shift+T", 2: "Shift+Ctrl+T"}
	for platform := range 3 {
		b.platform = platform
		tester.Frame()
		if !tester.HasText(texts[platform]) {
			t.Errorf("%s: the menu row does not write %q: %q", platformNames[platform], texts[platform], tester.Texts())
		}
	}

	b.selected = "Keycaps/Chord"
	tester.Frame()
	if !tester.HasText("[") || !tester.HasText("]") || !tester.HasText(`breeze.Keycaps(c, "[", "]")`) {
		t.Errorf("the chord story does not show its keys: %q", tester.Texts())
	}
}
