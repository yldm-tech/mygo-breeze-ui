package breeze

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

// referenceView lays out the controls as cmd/appkit-reference places AppKit's, so their screenshots compare control by control: three buttons, three check boxes and three switches, the last of each disabled.
func referenceView(p Platform, state *[6]bool) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Box(c).Fill().Children(func() {
				at := func(x, y float32, build func() ui.Element) {
					ui.Box(c).Absolute().Left(x).Top(y).Children(func() { build() })
				}
				at(16, 16, func() ui.Element { return Button(c, "Cancel") })
				at(110, 16, func() ui.Element { return PrimaryButton(c, "Save") })
				at(200, 16, func() ui.Element { return Button(c, "Disabled").Disabled(true) })
				at(16, 64, func() ui.Element { return Checkbox(c, &state[0], "Unchecked") })
				at(130, 64, func() ui.Element { return Checkbox(c, &state[1], "Checked") })
				at(240, 64, func() ui.Element { return Checkbox(c, &state[2], "Disabled").Disabled(true) })
				at(16, 106, func() ui.Element { return Switch(c, &state[3]).Label("Off") })
				at(100, 106, func() ui.Element { return Switch(c, &state[4]).Label("On") })
				at(184, 106, func() ui.Element { return Switch(c, &state[5]).Label("Disabled").Disabled(true) })
			})
		})
	}
}

// TestScreenshots writes the reference layout of every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names, to compare with the platforms' own controls. Without it the test only builds the frames.
func TestScreenshots(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			state := [6]bool{false, true, true, false, true, false}
			tester := ui.NewTester(referenceView(p, &state), 360, 200)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.Frame()
			for _, s := range []string{"Cancel", "Save", "Unchecked", "Checked"} {
				if !tester.HasText(s) {
					t.Fatalf("%v: %q is missing", p, s)
				}
			}
			if dir == "" {
				continue
			}
			name := fmt.Sprintf("breeze-%s-%s.png", strings.ToLower(p.String()), map[bool]string{false: "light", true: "dark"}[dark])
			f, err := os.Create(filepath.Join(dir, name))
			if err != nil {
				t.Fatal(err)
			}
			if err := png.Encode(f, tester.Image()); err != nil {
				t.Fatal(err)
			}
			if err := f.Close(); err != nil {
				t.Fatal(err)
			}
		}
	}
}
