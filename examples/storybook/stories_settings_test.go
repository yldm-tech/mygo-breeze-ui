package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestSettingsStoriesFollowTheirControls checks the settings rows' stories on every platform: the combo row shows the choice its control holds, and the expander row shows its rows while its open control is on and hides them while it is off.
func TestSettingsStoriesFollowTheirControls(t *testing.T) {
	b := newBook()
	tester := ui.NewTester(b.view, 1280, 820)
	tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	for platform := range 3 {
		b.platform = platform
		b.selected = "ComboRow/Default"
		a := b.argsOf(b.story())
		*a.text["selected"] = "Dark"
		tester.Frame()
		if !tester.HasText("Dark") || !tester.HasText(`breeze.ComboRow(c, "Appearance", "Light or dark", &app.appearance, []string{"Light", "Dark", "Auto"})`) {
			t.Errorf("%s: the combo row story does not show its choice or its code", platformNames[platform])
		}
		b.selected = "ExpanderRow/Default"
		a = b.argsOf(b.story())
		for _, open := range []bool{true, false} {
			*a.on["open"] = open
			tester.Frame()
			tester.Frame()
			// The code panel names the row too: what shows is a row with room of its own.
			r, _ := tester.Find("Proxy")
			if shown := r.H > 0; shown != open {
				t.Errorf("%s: open %v, the expander row shows its rows %v", platformNames[platform], open, shown)
			}
		}
	}
}
