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

// referenceView lays out the controls as tools/appkit-reference/reference.swift places AppKit's first set, so their screenshots compare control by control: three buttons, three check boxes and three switches, the last of each disabled.
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

// referenceView2 lays out more controls as tools/appkit-reference/reference.swift places AppKit's second set: radio buttons, fields, a pop-up button, a toggle, a slider, a progress bar, a segmented control, a spinner and a text view.
func referenceView2(p Platform, s *referenceState) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Box(c).Fill().Children(func() {
				at := func(x, y float32, build func() ui.Element) {
					ui.Box(c).Absolute().Left(x).Top(y).Children(func() { build() })
				}
				at(16, 16, func() ui.Element { return Radio(c, &s.size, "S", "Off") })
				at(100, 16, func() ui.Element { return Radio(c, &s.size, "M", "On") })
				at(16, 50, func() ui.Element { return TextInput(c, &s.name).Placeholder("Name").Width(180) })
				at(220, 50, func() ui.Element { return SearchField(c, &s.query).Width(180) })
				at(16, 90, func() ui.Element { return Select(c, &s.choice, []string{"Medium", "Small", "Large"}) })
				at(220, 90, func() ui.Element { return Toggle(c, &s.bold, "Bold") })
				at(16, 130, func() ui.Element { return Slider(c, &s.volume, 0, 100).Width(180) })
				at(220, 142, func() ui.Element { return Progress(c, 0.4).Width(180) })
				at(16, 170, func() ui.Element { return Segmented(c, &s.view, "Day", "Week", "Month") })
				at(220, 170, func() ui.Element { return Spinner(c) })
				at(16, 220, func() ui.Element { return TextArea(c, &s.notes).Size(180, 70).MinHeight(70) })
			})
		})
	}
}

// referenceState is what referenceView2's controls edit.
type referenceState struct {
	size, name, query, choice, notes string
	bold                             bool
	volume                           float64
	view                             int
}

func newReferenceState() *referenceState {
	return &referenceState{size: "M", choice: "Medium", notes: "Notes", bold: true, volume: 40}
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
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			writeShot(t, dir, fmt.Sprintf("breeze-%s-%s.png", strings.ToLower(p.String()), mode), tester)
			second := ui.NewTester(referenceView2(p, newReferenceState()), 520, 330)
			second.SetScale(2)
			second.SetDark(dark)
			second.Frame()
			for _, s := range []string{"Off", "On", "Medium", "Bold", "Week"} {
				if !second.HasText(s) {
					t.Fatalf("%v: %q is missing from the second set", p, s)
				}
			}
			writeShot(t, dir, fmt.Sprintf("breeze2-%s-%s.png", strings.ToLower(p.String()), mode), second)
		}
	}
}

// writeShot writes tester's frame to dir as name, when dir is set.
func writeShot(t *testing.T, dir, name string, tester *ui.Tester) {
	t.Helper()
	if dir == "" {
		return
	}
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
