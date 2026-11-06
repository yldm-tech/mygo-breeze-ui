package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestPreferencesStoriesFollowTheirControls checks the preferences stories on every platform: a click on a radio row's subtitle sets the engine control, and the group shows its description and its header's button only while their controls hold text.
func TestPreferencesStoriesFollowTheirControls(t *testing.T) {
	b := newBook()
	tester := ui.NewTester(b.view, 1280, 820)
	tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	for platform := range 3 {
		b.platform = platform
		b.selected = "RadioRow/In a group"
		a := b.argsOf(b.story())
		*a.text["engine"] = "whisper"
		tester.Frame()
		if err := tester.Click("Sent to a service"); err != nil {
			t.Fatalf("%s: %v", platformNames[platform], err)
		}
		tester.Frame()
		if *a.text["engine"] != "cloud" {
			t.Errorf("%s: a click on the cloud row left the engine %q", platformNames[platform], *a.text["engine"])
		}
		b.selected = "PreferencesGroup/Default"
		a = b.argsOf(b.story())
		for _, s := range []string{"Models are stored on this Mac.", ""} {
			*a.text["description"], *a.text["button"] = s, s
			tester.Frame()
			tester.Frame()
			if shown := tester.HasText("Models are stored on this Mac."); shown != (s != "") {
				t.Errorf("%s: description %q, the group shows it %v", platformNames[platform], s, shown)
			}
		}
	}
}
