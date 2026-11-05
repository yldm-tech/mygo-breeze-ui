package breeze

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

// referenceView3 lays out the controls restyled from package ui's as tools/adwaita-reference/reference.py places libadwaita's third set: a link, a separator, joined toggles, a menu button and a popover open below its anchor.
func referenceView3(p Platform) func(c *ui.Context) {
	bold, italic, open := true, false, true
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Box(c).Fill().Children(func() {
				at := func(x, y float32, build func() ui.Element) {
					ui.Box(c).Absolute().Left(x).Top(y).Children(func() { build() })
				}
				at(16, 16, func() ui.Element { return Link(c, "Docs", "https://example.com") })
				at(16, 50, func() ui.Element { return Divider(c).Width(200) })
				at(16, 70, func() ui.Element {
					return ToggleGroup(c, func() {
						Toggle(c, &bold, "Bold")
						Toggle(c, &italic, "Italic")
					})
				})
				at(16, 120, func() ui.Element { return MenuButton(c, "Sort", func(m *ui.Menu) { m.Item("Name") }) })
				at(200, 70, func() ui.Element {
					info := Button(c, "Info")
					Popover(c, info, &open, func() { ui.Text(c, "Details") })
					return info
				})
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
			open := true
			alert := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					ui.Column(c).Fill().Padding(16).Children(func() { ui.Text(c, "Window content") })
					AlertDialog(c, &open, "Delete “Notes”?", "You can’t undo this action.", "Cancel", "Delete")
				})
			}, 560, 340)
			alert.SetScale(2)
			alert.SetDark(dark)
			alert.Frame()
			writeShot(t, dir, fmt.Sprintf("breeze-alert-%s-%s.png", strings.ToLower(p.String()), mode), alert)
			third := ui.NewTester(referenceView3(p), 400, 240)
			third.SetScale(2)
			third.SetDark(dark)
			third.Frame()
			third.Frame()
			if !third.HasText("Details") || !third.HasText("Italic") {
				t.Fatalf("%v: the third set is missing its popover or toggles", p)
			}
			writeShot(t, dir, fmt.Sprintf("breeze3-%s-%s.png", strings.ToLower(p.String()), mode), third)
		}
	}
}

// windowsState is a control of tools/winui-reference's fifth set, which holds WinUI's controls in their PointerOver and Pressed states: where it is, the region around it that its shot is cut from, whether the pointer presses it or only moves over it, and what builds it.
type windowsState struct {
	x, y    float32
	region  image.Rectangle // in DIPs
	pressed bool
	build   func(c *ui.Context) ui.Element
}

func windowsStates() []windowsState {
	off, on := false, true
	size, chosen := "S", "M"
	volume := 40.0
	button := func(label string, primary bool) func(c *ui.Context) ui.Element {
		return func(c *ui.Context) ui.Element {
			if primary {
				return PrimaryButton(c, label)
			}
			return Button(c, label)
		}
	}
	check := func(v *bool, label string) func(c *ui.Context) ui.Element {
		return func(c *ui.Context) ui.Element { return Checkbox(c, v, label) }
	}
	sw := func(v *bool) func(c *ui.Context) ui.Element {
		return func(c *ui.Context) ui.Element { return Switch(c, v) }
	}
	radio := func(v *string, label string) func(c *ui.Context) ui.Element {
		return func(c *ui.Context) ui.Element { return Radio(c, v, "M", label) }
	}
	slider := func(c *ui.Context) ui.Element { return Slider(c, &volume, 0, 100).Width(180) }
	r := func(x0, y0, x1, y1 int) image.Rectangle { return image.Rect(x0, y0, x1, y1) }
	return []windowsState{
		{16, 16, r(0, 0, 106, 56), false, button("Cancel", false)},
		{110, 16, r(106, 0, 200, 56), true, button("Cancel", false)},
		{204, 16, r(200, 0, 280, 56), false, button("Save", true)},
		{284, 16, r(280, 0, 380, 56), true, button("Save", true)},
		{16, 64, r(0, 56, 136, 100), false, check(&off, "Off")},
		{140, 64, r(136, 56, 260, 100), true, check(&off, "Off")},
		{264, 64, r(260, 56, 384, 100), false, check(&on, "On")},
		{388, 64, r(384, 56, 520, 100), true, check(&on, "On")},
		{16, 106, r(0, 100, 136, 146), false, sw(&off)},
		{140, 106, r(136, 100, 260, 146), true, sw(&off)},
		{264, 106, r(260, 100, 384, 146), false, sw(&on)},
		{388, 106, r(384, 100, 520, 146), true, sw(&on)},
		{16, 150, r(0, 146, 136, 188), true, radio(&size, "Off")},
		{140, 150, r(136, 146, 260, 188), false, radio(&chosen, "On")},
		{264, 150, r(260, 146, 384, 188), true, radio(&chosen, "On")},
		{16, 194, r(0, 188, 210, 240), false, slider},
		{220, 194, r(210, 188, 420, 240), true, slider},
	}
}

// TestScreenshotsWindowsStates writes the controls of tools/winui-reference's fifth set, each under the pointer or pressed by it, as breeze-states-windows-light.png and -dark.png, cut together from one shot per control as the pointer is over one control at a time.
func TestScreenshotsWindowsStates(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, dark := range []bool{false, true} {
		out := image.NewRGBA(image.Rect(0, 0, 520*2, 260*2))
		for i, s := range windowsStates() {
			var box ui.Rect
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, Windows, func() {
					Use(c)
					ui.Box(c).Fill().Children(func() {
						ui.Box(c).Absolute().Left(s.x).Top(s.y).Children(func() { box = s.build(c).Bounds() })
					})
				})
			}, 520, 260)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.Frame()
			tester.Frame()
			if box.W == 0 {
				t.Fatalf("state %d: the control has no size", i)
			}
			// The slider is pressed on its thumb, at 40, so that its value stays.
			px, py := box.X+box.W/2, box.Y+box.H/2
			if box.W == 180 {
				px = box.X + 9 + 0.4*(box.W-18)
			}
			tester.Move(px, py)
			if s.pressed {
				tester.Press(px, py)
			}
			tester.Advance(time.Second)
			tester.Frame()
			if i == 0 {
				draw.Draw(out, out.Bounds(), tester.Image(), image.Point{}, draw.Src)
			}
			reg := image.Rect(s.region.Min.X*2, s.region.Min.Y*2, s.region.Max.X*2, s.region.Max.Y*2)
			draw.Draw(out, reg, tester.Image(), reg.Min, draw.Src)
		}
		if dir == "" {
			continue
		}
		f, err := os.Create(filepath.Join(dir, fmt.Sprintf("breeze-states-windows-%s.png", map[bool]string{false: "light", true: "dark"}[dark])))
		if err != nil {
			t.Fatal(err)
		}
		if err := png.Encode(f, out); err != nil {
			t.Fatal(err)
		}
		if err := f.Close(); err != nil {
			t.Fatal(err)
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
