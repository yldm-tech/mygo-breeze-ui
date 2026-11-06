package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestThemeStoryFollowsItsControls checks the ThemeOf story on every platform: it shows its controls twice, the canvas's and the subtree's, and its code names the appearance and the accent its controls choose.
func TestThemeStoryFollowsItsControls(t *testing.T) {
	b := newBook()
	b.selected = "ThemeOf/Dark subtree"
	tester := ui.NewTester(b.view, 1280, 820)
	a := b.argsOf(b.story())
	for platform := range 3 {
		b.platform = platform
		tester.Frame()
		n := 0
		for _, s := range tester.Texts() {
			if s == "Play" {
				n++
			}
		}
		if n != 2 {
			t.Errorf("%s: the story shows Play %d times, want 2", platformNames[platform], n)
		}
	}
	if !tester.HasText("ui.Themed(c, breeze.ThemeOf(c, breeze.PlatformOf(c), true, ui.Color{})).Padding(16).Gap(12).Children(func() {\n\tbreeze.PrimaryButton(c, \"Play\")\n\tbreeze.Segmented(c, &app.view, \"Day\", \"Week\")\n\tbreeze.Switch(c, &app.subtitles).Label(\"Subtitles\")\n})") {
		t.Errorf("the code of the dark subtree in the system accent is missing: %q", tester.Texts())
	}
	*a.on["dark"] = false
	*a.text["accent"] = "Purple"
	tester.Frame()
	if !tester.HasText("ui.Themed(c, breeze.ThemeOf(c, breeze.PlatformOf(c), false, ui.Hex(\"#8e44ad\"))).Padding(16).Gap(12).Children(func() {\n\tbreeze.PrimaryButton(c, \"Play\")\n\tbreeze.Segmented(c, &app.view, \"Day\", \"Week\")\n\tbreeze.Switch(c, &app.subtitles).Label(\"Subtitles\")\n})") {
		t.Errorf("the code does not follow the dark and accent controls: %q", tester.Texts())
	}
}
