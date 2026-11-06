package breeze

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func init() {
	for _, n := range []string{"ComboRow", "EntryRow", "SpinRow", "ExpanderRow"} {
		drawnWidgets[n] = true
	}
}

// settingsView lays out the settings rows for screenshots, in two columns 300 wide as rowsView's: a card of a combo row with a subtitle, an entry row holding text, an empty one and a spin row with a subtitle, then a combo row alone and a disabled spin row alone; beside them an open expander row with a subtitle holding an action row and a switch row, and a closed one without a subtitle.
func settingsView(p Platform) func(c *ui.Context) {
	appearance, language, name, email := "Auto", "English", "Evan", ""
	size, count := 12.0, 3.0
	open, advanced, metered := true, false, false
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Row(c).Fill().Padding(16).Gap(24).AlignItems(ui.Start).Children(func() {
				ui.Column(c).Width(300).Gap(16).Children(func() {
					Card(c, func() {
						ComboRow(c, "Appearance", "Light or dark", &appearance, []string{"Light", "Dark", "Auto"})
						EntryRow(c, "Name", &name)
						EntryRow(c, "Email", &email)
						SpinRow(c, "Font size", "Points", &size, 6, 72, 1)
					})
					ComboRow(c, "Language", "", &language, []string{"English", "Français"})
					SpinRow(c, "Copies", "", &count, 1, 9, 1).Disabled(true)
				})
				ui.Column(c).Width(300).Gap(16).Children(func() {
					ExpanderRow(c, "Network", "Wi-Fi and proxies", &open, func() {
						ActionRow(c, "Proxy", "Automatic", nil)
						SwitchRow(c, "Metered", "", &metered)
					})
					ExpanderRow(c, "Advanced", "", &advanced, func() {
						ActionRow(c, "Logs", "", nil)
					})
				})
			})
		})
	}
}

// TestScreenshotsSettings writes the settings rows of every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names. Without it the test only builds the frames.
func TestScreenshotsSettings(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			h := 480
			if p == Windows {
				h = 640
			}
			tester := ui.NewTester(settingsView(p), 672, h)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tester.Frame()
			tester.Frame()
			for _, s := range []string{"Appearance", "Light or dark", "Auto", "Font size", "Points", "Language", "Network", "Wi-Fi and proxies", "Proxy", "Metered", "Advanced"} {
				if !tester.HasText(s) {
					t.Fatalf("%v: %q is missing", p, s)
				}
			}
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			writeShot(t, dir, fmt.Sprintf("breeze-settings-%s-%s.png", strings.ToLower(p.String()), mode), tester)
		}
	}
}
