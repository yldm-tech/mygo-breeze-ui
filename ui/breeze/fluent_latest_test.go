package breeze

import (
	"math"
	"testing"
	"time"

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

// fluentNear reports whether a and b differ by at most 1 in each channel, as rounding leaves them.
func fluentNear(a, b ui.Color) bool {
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
		if want := bg.Mix(th.Accent, 0.9); !fluentNear(th.AccentHover, want) {
			t.Errorf("dark %v: AccentHover %v, want the fill at 90%% over the window, %v", dark, th.AccentHover, want)
		}
		if want := bg.Mix(th.Accent, 0.8); !fluentNear(th.AccentPressed, want) {
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
		if want := pick(dark, ui.RGBA(0, 0, 0, 0x72/255.0), ui.RGBA(255, 255, 255, 0x8b/255.0)); !fluentNear(th.Scrollbar, want) {
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

// TestFluentSelectFocus checks that a ComboBox with the keyboard's focus shows its template's focus look, not the system focus visual: the accent Pill 1 DIP in from the box's start, mirrored right to left, which is not there without the focus.
func TestFluentSelectFocus(t *testing.T) {
	for _, rtl := range []bool{false, true} {
		sel := "Medium"
		var box ui.Rect
		var accent ui.Color
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Windows, func() {
				Use(c)
				accent = c.Theme().Accent
				ui.Column(c).Fill().Padding(24).AlignItems(ui.Start).Children(func() {
					box = Select(c, &sel, []string{"Small", "Medium"}).Bounds()
				})
			})
		}, 240, 90)
		if rtl {
			tester.SetDirection(ui.RightToLeft)
		}
		tester.Frame()
		tester.Frame()
		x := int(box.X + 2)
		if rtl {
			x = int(box.X + box.W - 3)
		}
		y := int(box.Y + box.H/2)
		pill := func() ui.Color {
			r, g, b, a := tester.Image().At(x, y).RGBA()
			return ui.Color{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
		}
		if fluentNear(pill(), accent) {
			t.Fatalf("rtl %v: the pill shows without the focus", rtl)
		}
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		if got := pill(); !fluentNear(got, accent) {
			t.Errorf("rtl %v: with the keyboard's focus the box's start is %v, want the accent pill %v", rtl, got, accent)
		}
	}
}

// fluentPixel returns the pixel of tester's 2x frame at x, y in DIPs.
func fluentPixel(tester *ui.Tester, x, y float32) ui.Color {
	r, g, b, a := tester.Image().At(int(x*2), int(y*2)).RGBA()
	return ui.Color{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}

// fluentWithin reports whether a and b differ by at most d in each color channel.
func fluentWithin(a, b ui.Color, d int) bool {
	ok := func(x, y uint8) bool { return int(x)-int(y) <= d && int(y)-int(x) <= d }
	return ok(a.R, b.R) && ok(a.G, b.G) && ok(a.B, b.B)
}

// fluentShot builds view on Windows at 2x, light or dark, and returns the tester after two frames, so that bounds read while building are the laid-out ones.
func fluentShot(dark bool, w, h int, view func(c *ui.Context)) *ui.Tester {
	tester := ui.NewTester(func(c *ui.Context) {
		Preview(c, Windows, func() {
			Use(c)
			view(c)
		})
	}, w, h)
	tester.SetScale(2)
	tester.SetDark(dark)
	tester.Frame()
	tester.Frame()
	return tester
}

// TestFluentAccentAsMeasured checks the accent fills and link colors against WinUI 3's rendering with Windows's default accent (tools/winui-reference, WinUI 3 3.3 on Windows 11 26H2): the accent button #005a9e and #76b9ed, under the pointer #1869a6 and #6da9d8, pressed #3179af and #649ac4, and a hyperlink #004275 and #a6d8ff.
func TestFluentAccentAsMeasured(t *testing.T) {
	for _, dark := range []bool{false, true} {
		th := windowsTheme(t, dark)
		want := map[string][2]ui.Color{
			"Accent":        {th.Accent, pick(dark, ui.Hex("#005a9e"), ui.Hex("#76b9ed"))},
			"AccentHover":   {th.AccentHover, pick(dark, ui.Hex("#1869a6"), ui.Hex("#6da9d8"))},
			"AccentPressed": {th.AccentPressed, pick(dark, ui.Hex("#3179af"), ui.Hex("#649ac4"))},
		}
		for name, c := range want {
			if !fluentNear(c[0], c[1]) {
				t.Errorf("dark %v: %s %v, WinUI's reads %v", dark, name, c[0], c[1])
			}
		}
		var link ui.Color
		fluentShot(dark, 100, 40, func(c *ui.Context) { link = fluentLinkColor(c, false, false) })
		if wantLink := pick(dark, ui.Hex("#004275"), ui.Hex("#a6d8ff")); link != wantLink {
			t.Errorf("dark %v: link %v, WinUI's reads %v", dark, link, wantLink)
		}
	}
}

// TestFluentChoicesAsMeasured checks a check box, a radio button and a switch against WinUI 3's rendering (tools/winui-reference, WinUI 3 3.3 on Windows 11 26H2): the label's text box 5 DIPs below the check box's top and 6 below the radio button's, the box 6 below; the unchecked box's stroke reading 134 then 131 into the box (154 and 152 dark), 191 and 173 pressed (67 and 80); the off switch's knob 90 on its track (204 dark) and 87 under the pointer (207); the chosen radio button's dot ring over the accent below the dot, #00365f in the light theme; and a disabled checked box gray, 191 (67 dark), as MyGo's half opacity lets it read.
func TestFluentChoicesAsMeasured(t *testing.T) {
	for _, dark := range []bool{false, true} {
		off, on, sw, size := false, true, false, "M"
		var boxRow, onRow, disRow, radioRow, swBox ui.Rect
		view := func(c *ui.Context) {
			ui.Box(c).Fill().Children(func() {
				at := func(x, y float32, build func() ui.Element) ui.Rect {
					var r ui.Rect
					ui.Box(c).Absolute().Left(x).Top(y).Children(func() { r = build().Bounds() })
					return r
				}
				boxRow = at(16, 16, func() ui.Element { return Checkbox(c, &off, "Unchecked") })
				onRow = at(160, 16, func() ui.Element { return Checkbox(c, &on, "Checked") })
				disRow = at(300, 16, func() ui.Element { return Checkbox(c, &on, "Disabled").Disabled(true) })
				radioRow = at(16, 60, func() ui.Element { return Radio(c, &size, "M", "Chosen") })
				swBox = at(160, 66, func() ui.Element { return Switch(c, &sw) })
			})
		}
		tester := fluentShot(dark, 440, 120, view)
		if r, ok := tester.Find("Unchecked"); !ok || r.Y-boxRow.Y != 5 {
			t.Errorf("dark %v: the check box's label is %v DIPs below its top, want 5", dark, r.Y-boxRow.Y)
		}
		if r, ok := tester.Find("Chosen"); !ok || r.Y-radioRow.Y != 6 {
			t.Errorf("dark %v: the radio button's label is %v DIPs below its top, want 6", dark, r.Y-radioRow.Y)
		}
		box := ui.Rect{X: boxRow.X, Y: boxRow.Y + 6, W: 20, H: 20}
		outer, inner := fluentPixel(tester, box.X+10, box.Y+0.25), fluentPixel(tester, box.X+10, box.Y+0.75)
		if want := pick(dark, [2]uint8{134, 131}, [2]uint8{154, 152}); !fluentWithin(outer, ui.RGBA(want[0], want[0], want[0], 1), 2) || !fluentWithin(inner, ui.RGBA(want[1], want[1], want[1], 1), 2) {
			t.Errorf("dark %v: the box's stroke reads %v then %v, WinUI's %v", dark, outer, inner, want)
		}
		cx, cy := radioRow.X+10, radioRow.Y+6+10
		if got := fluentPixel(tester, cx, cy+6-0.25); !dark && !fluentWithin(got, ui.Hex("#00365f"), 4) {
			t.Errorf("the dot's ring below the dot reads %v over the accent, WinUI's #00365f", got)
		}
		knob := func() ui.Color { return fluentPixel(tester, swBox.X+10, swBox.Y+10) }
		if got, want := knob(), pick(dark, uint8(90), 204); !fluentWithin(got, ui.RGBA(want, want, want, 1), 2) {
			t.Errorf("dark %v: the off knob reads %v, WinUI's %d", dark, got, want)
		}
		if got, want := fluentPixel(tester, disRow.X+2.5, disRow.Y+6+3), pick(dark, uint8(191), 67); !fluentWithin(got, ui.RGBA(want, want, want, 1), 3) {
			t.Errorf("dark %v: the disabled checked box reads %v, WinUI's %d", dark, got, want)
		}
		tester.Move(swBox.X+20, swBox.Y+10)
		tester.Advance(time.Second)
		tester.Frame()
		if got, want := knob(), pick(dark, uint8(87), 207); !fluentWithin(got, ui.RGBA(want, want, want, 1), 2) {
			t.Errorf("dark %v: the off knob under the pointer reads %v, WinUI's %d", dark, got, want)
		}
		tester.Move(box.X+10, box.Y+10)
		tester.Press(box.X+10, box.Y+10)
		tester.Advance(time.Second)
		tester.Frame()
		outer, inner = fluentPixel(tester, box.X+10, box.Y+0.25), fluentPixel(tester, box.X+10, box.Y+0.75)
		if want := pick(dark, [2]uint8{191, 173}, [2]uint8{67, 80}); !fluentWithin(outer, ui.RGBA(want[0], want[0], want[0], 1), 2) || !fluentWithin(inner, ui.RGBA(want[1], want[1], want[1], 1), 2) {
			t.Errorf("dark %v: the pressed box's stroke reads %v then %v, WinUI's %v", dark, outer, inner, want)
		}
		if onRow.W == 0 {
			t.Fatalf("dark %v: the checked box has no size", dark)
		}
	}
}

// TestFluentSelectAndSliderAsMeasured checks a ComboBox and a slider against WinUI 3's rendering (tools/winui-reference, WinUI 3 3.3 on Windows 11 26H2): the ComboBox's text 12 DIPs in, its chevron's point 1.75 DIPs below the box's middle (centered on the box, not on its padded content) and 14 DIPs from its end, and the slider thumb's edge translucent over the accent track, #005595 in the light theme.
func TestFluentSelectAndSliderAsMeasured(t *testing.T) {
	for _, dark := range []bool{false, true} {
		sel, volume := "Medium", 40.0
		var box, slider ui.Rect
		tester := fluentShot(dark, 300, 120, func(c *ui.Context) {
			ui.Box(c).Fill().Children(func() {
				ui.Box(c).Absolute().Left(16).Top(16).Children(func() { box = Select(c, &sel, []string{"Medium", "Small"}).Bounds() })
				ui.Box(c).Absolute().Left(16).Top(64).Children(func() { slider = Slider(c, &volume, 0, 100).Width(180).Bounds() })
			})
		})
		if r, ok := tester.Find("Medium"); !ok || r.X-box.X != 12 {
			t.Errorf("dark %v: the ComboBox's text is %v DIPs in, want 12", dark, r.X-box.X)
		}
		// The chevron is 7 DIPs wide, its point in the middle of its 12-DIP box, which ends 14 DIPs from the end.
		px := int((box.X + box.W - 14 - 6) * 2)
		face := fluentPixel(tester, box.X+3, box.Y+3)
		best, darkest := -1, 0
		for y := int(box.Y * 2); y < int((box.Y+box.H)*2); y++ {
			r, _, _, _ := tester.Image().At(px, y).RGBA()
			d := int(face.R) - int(r>>8)
			if dark {
				d = -d
			}
			if d > darkest {
				best, darkest = y, d
			}
		}
		if point := float32(best)/2 + 0.25; point < box.Y+box.H/2+1.25 || point > box.Y+box.H/2+2.25 {
			t.Errorf("dark %v: the chevron's point is at %v, want 1.75 below the box's middle, %v", dark, point, box.Y+box.H/2+1.75)
		}
		thumb := slider.X + 9 + 0.4*(slider.W-18)
		if got := fluentPixel(tester, thumb-11+0.25, slider.Y+slider.H/2); !dark && !fluentWithin(got, ui.Hex("#005595"), 3) {
			t.Errorf("the thumb's edge over the track reads %v, WinUI's #005595", got)
		}
	}
}

// TestFluentOverlayStrokes checks the edges of a flyout and of an alert, which WinUI draws over what is behind them, not over their fill (tools/winui-reference, WinUI 3 3.3 on Windows 11 26H2): the flyout's top edge reads 229 in the light theme and 26 in the dark one, and the alert's edge is SurfaceStrokeColorDefault over the smoke over the window.
func TestFluentOverlayStrokes(t *testing.T) {
	for _, dark := range []bool{false, true} {
		third := ui.NewTester(referenceView3(Windows), 400, 240)
		third.SetScale(2)
		third.SetDark(dark)
		third.Frame()
		third.Frame()
		d, ok := third.Find("Details")
		if !ok {
			t.Fatal("the popover is not open")
		}
		// The text is inside the 1-DIP border and the padding of 15 above it.
		if got, want := fluentPixel(third, d.X+d.W/2, d.Y-16+0.25), pick(dark, uint8(229), 26); !fluentWithin(got, ui.RGBA(want, want, want, 1), 1) {
			t.Errorf("dark %v: the flyout's top edge reads %v, WinUI's %d", dark, got, want)
		}
		open := true
		var th *ui.Theme
		alert := fluentShot(dark, 560, 340, func(c *ui.Context) {
			th = c.Theme()
			AlertDialog(c, &open, "Delete?", "You can’t undo this action.", "Cancel", "Delete")
		})
		panel, ok := alert.Find("Delete?")
		if !ok {
			t.Fatal("the alert is not open")
		}
		if got := fluentPixel(alert, panel.X+panel.W/2, panel.Y+0.25); !fluentWithin(got, fluentDialogStroke(th), 1) {
			t.Errorf("dark %v: the alert's edge reads %v, want %v", dark, got, fluentDialogStroke(th))
		}
	}
}
