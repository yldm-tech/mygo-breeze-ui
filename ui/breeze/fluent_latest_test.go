package breeze

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"
)

// windowsTheme returns the Windows theme a window without preferences gets, light or dark.
func windowsTheme(t *testing.T, dark bool) *ui.Theme {
	t.Helper()
	var got *ui.Theme
	tester := ui.NewTester(func(c *ui.Context) { got = Theme(c, Windows) }, 100, 100)
	tester.SetDark(dark)
	tester.Frame()
	return got
}

// near reports whether a and b differ by at most 1 in each channel, as rounding leaves them.
func near(a, b ui.Color) bool {
	d := func(x, y uint8) bool { return int(x)-int(y) <= 1 && int(y)-int(x) <= 1 }
	return d(a.R, b.R) && d(a.G, b.G) && d(a.B, b.B) && d(a.A, b.A)
}

// TestFluentThemeFromWinUI checks the theme values taken from WinUI's Common_themeresources_any.xaml, ScrollBar_themeresources.xaml and TextBox_themeresources.xaml: the accent's hover and pressed fills at 90% and 80% over the window's own background, the SystemFillColor status colors, ScrollBarThumbFill and the text selection.
func TestFluentThemeFromWinUI(t *testing.T) {
	for _, dark := range []bool{false, true} {
		th := windowsTheme(t, dark)
		bg := pick(dark, ui.Hex("#f3f3f3"), ui.Hex("#202020"))
		if th.Background != bg {
			t.Fatalf("dark %v: background %v, want %v", dark, th.Background, bg)
		}
		if want := bg.Mix(th.Accent, 0.9); !near(th.AccentHover, want) {
			t.Errorf("dark %v: AccentHover %v, want the fill at 90%% over the window, %v", dark, th.AccentHover, want)
		}
		if want := bg.Mix(th.Accent, 0.8); !near(th.AccentPressed, want) {
			t.Errorf("dark %v: AccentPressed %v, want the fill at 80%% over the window, %v", dark, th.AccentPressed, want)
		}
		status := map[string][2]ui.Color{
			"Danger":  {th.Danger, pick(dark, ui.Hex("#c42b1c"), ui.Hex("#ff99a4"))},
			"Warning": {th.Warning, pick(dark, ui.Hex("#9d5d00"), ui.Hex("#fce100"))},
			"Success": {th.Success, pick(dark, ui.Hex("#0f7b0f"), ui.Hex("#6ccb5f"))},
		}
		for name, c := range status {
			if c[0] != c[1] {
				t.Errorf("dark %v: %s %v, want %v", dark, name, c[0], c[1])
			}
		}
		if want := pick(dark, ui.RGBA(0, 0, 0, 0x72/255.0), ui.RGBA(255, 255, 255, 0x8b/255.0)); !near(th.Scrollbar, want) {
			t.Errorf("dark %v: Scrollbar %v, want ControlStrongFillColorDefault %v", dark, th.Scrollbar, want)
		}
		if dark && th.Selection != fluentAccent {
			t.Errorf("dark: Selection %v, want the opaque SystemAccentColor %v", th.Selection, fluentAccent)
		}
		if !dark && th.Selection.A == 255 {
			t.Errorf("light: Selection %v is opaque under near-black text", th.Selection)
		}
	}
}

// TestFluentRingArc checks the ProgressRing's Lottie motion: the arc never sweeps more than half the circle, grows to it in the first second while turning 450 degrees, and shrinks to nothing in the next while turning on to 900.
func TestFluentRingArc(t *testing.T) {
	for i := 0; i < 2000; i++ {
		if _, sweep := fluentRingArcAt(float32(i) / 1000); sweep < 0 || sweep > math.Pi+1e-6 {
			t.Fatalf("at %d ms the arc sweeps %.3f radians, outside 0 to π", i, sweep)
		}
	}
	checks := []struct {
		s            float32
		start, sweep float64
	}{
		{0, 0, 0},
		{1, 2.5 * math.Pi, math.Pi},
		{2, 6 * math.Pi, 0},
	}
	for _, c := range checks {
		start, sweep := fluentRingArcAt(c.s)
		if math.Abs(start-c.start) > 1e-3 || math.Abs(sweep-c.sweep) > 1e-3 {
			t.Errorf("at %v s: start %.4f, sweep %.4f; want %.4f, %.4f", c.s, start, sweep, c.start, c.sweep)
		}
	}
}

// TestFluentIndeterminateBars checks the two bars of ProgressBar.xaml's Indeterminate storyboard for a 200-DIP bar: 40% and 60% wide, the first from -100% to 300% of its width by 1.5 s, the second held at -150% of its width until 0.75 s and at 166% of it by 2 s.
func TestFluentIndeterminateBars(t *testing.T) {
	const w = 200
	at := func(s float32) [2][2]float32 { return fluentIndeterminateAt(s, w) }
	for _, s := range []float32{0, 1, 1.999} {
		if b := at(s); math.Abs(float64(b[0][1]-80)) > 0.01 || math.Abs(float64(b[1][1]-120)) > 0.01 {
			t.Fatalf("widths %v and %v, want 80 and 120", b[0][1], b[1][1])
		}
	}
	checks := []struct {
		name      string
		got, want float32
	}{
		{"first bar at 0 s", at(0)[0][0], -80},
		{"first bar at 1.5 s", at(1.5)[0][0], 240},
		{"first bar at 1.9 s", at(1.9)[0][0], 240},
		{"second bar at 0 s", at(0)[1][0], -180},
		{"second bar at 0.75 s", at(0.75)[1][0], -180},
		{"second bar at 2 s", at(2)[1][0], 199.2},
	}
	for _, c := range checks {
		if math.Abs(float64(c.got-c.want)) > 0.01 {
			t.Errorf("%s: %v, want %v", c.name, c.got, c.want)
		}
	}
}

// TestFluentCheckAndRadioMinWidth checks that a check box and a radio button with a short label take WinUI's MinWidth of 120 on Windows, so that their focus visual, which goes around the whole control, is as wide as WinUI's.
func TestFluentCheckAndRadioMinWidth(t *testing.T) {
	on, size := false, "S"
	var boxW, radioW float32
	tester := ui.NewTester(func(c *ui.Context) {
		Preview(c, Windows, func() {
			Use(c)
			ui.Column(c).AlignItems(ui.Start).Children(func() {
				// Bounds are the previous frame's, read while building.
				boxW = Checkbox(c, &on, "A").Bounds().W
				radioW = Radio(c, &size, "S", "B").Bounds().W
			})
		})
	}, 400, 200)
	tester.Frame()
	tester.Frame()
	if boxW < 120 {
		t.Errorf("check box %v DIPs wide, want at least 120 (CheckBoxMinWidth)", boxW)
	}
	if radioW < 120 {
		t.Errorf("radio button %v DIPs wide, want at least 120", radioW)
	}
}
