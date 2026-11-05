package breeze

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// macAt lays out build's element with its top left corner at x, y in a window of macOS's controls.
func macAt(x, y float32, build func(c *ui.Context)) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, Mac, func() {
			Use(c)
			ui.Box(c).Fill().Children(func() {
				ui.Box(c).Absolute().Left(x).Top(y).Children(func() { build(c) })
			})
		})
	}
}

// macPixel returns the color the tester painted at x, y.
func macPixel(tester *ui.Tester, x, y int) ui.Color {
	px := tester.Image().RGBAAt(x, y)
	return ui.RGB(px.R, px.G, px.B)
}

// macNear reports whether a and b differ by at most d in every channel.
func macNear(a, b ui.Color, d int) bool {
	diff := func(x, y uint8) int { return max(int(x)-int(y), int(y)-int(x)) }
	return diff(a.R, b.R) <= d && diff(a.G, b.G) <= d && diff(a.B, b.B) <= d
}

// TestMacThemeOfMacOS27 checks the theme's colors measured on macOS 27.2: separatorColor, the focus ring of the blue accent as a key window shows it, and the overlay scroller's knob; another accent keeps its own ring.
func TestMacThemeOfMacOS27(t *testing.T) {
	for _, dark := range []bool{false, true} {
		th := macTheme(dark, ui.Color{})
		border := pick(dark, ui.RGB(230, 230, 230), ui.RGB(57, 57, 57))
		if !macNear(th.Border, border, 1) {
			t.Errorf("dark %v: Border %v, want separatorColor %v", dark, th.Border, border)
		}
		focus := pick(dark, ui.RGBA(27, 105, 249, 0.5), ui.RGBA(46, 169, 255, 0.5))
		if th.Focus != focus {
			t.Errorf("dark %v: Focus %v, want %v", dark, th.Focus, focus)
		}
		if th.Scrollbar.A != ui.RGBA(0, 0, 0, 0.5).A {
			t.Errorf("dark %v: Scrollbar %v, want the knob at half opacity", dark, th.Scrollbar)
		}
		red := ui.Hex("#ff383c")
		if got := macTheme(dark, red).Focus; got.A != focus.A || (!dark && got != red.Alpha(0.5)) {
			t.Errorf("dark %v: Focus for a red accent %v", dark, got)
		}
	}
}

// TestMacAccentControlsDrawTheKeyWindow checks that the on toggle and the chosen segment are drawn in the accent, as AppKit draws them in a key window, while the chosen tab is a neutral platter, not the accent.
func TestMacAccentControlsDrawTheKeyWindow(t *testing.T) {
	for _, dark := range []bool{false, true} {
		on, seg, tab := true, 0, 0
		tester := ui.NewTester(macAt(40, 40, func(c *ui.Context) {
			ui.Column(c).Gap(36).AlignItems(ui.Start).Children(func() {
				Toggle(c, &on, "Bold")
				Segmented(c, &seg, "Day", "Week")
				Tabs(c, &tab, "General", "Advanced")
			})
		}), 400, 240)
		tester.SetDark(dark)
		tester.Frame()
		accent := macTheme(dark, ui.Color{}).Accent
		if got := macPixel(tester, 44, 52); got != accent {
			t.Errorf("dark %v: on toggle %v, want the accent %v", dark, got, accent)
		}
		if got := macPixel(tester, 44, 112); got != accent {
			t.Errorf("dark %v: chosen segment %v, want the accent %v", dark, got, accent)
		}
		platter := pick(dark, white, tint(ui.Hex("#1e1e1e"), true, 0.2))
		if got := macPixel(tester, 46, 172); !macNear(got, platter, 2) {
			t.Errorf("dark %v: chosen tab %v, want the platter %v", dark, got, platter)
		}
	}
}

// TestMacDisabledCheckedBoxIsGray checks that a checked box made disabled after it is built is drawn gray, as AppKit draws it, not in the accent: at MyGo's half opacity it reads 236 over white and 46 over the window in Dark Aqua.
func TestMacDisabledCheckedBoxIsGray(t *testing.T) {
	for _, dark := range []bool{false, true} {
		checked := true
		tester := ui.NewTester(macAt(40, 40, func(c *ui.Context) { Checkbox(c, &checked, "Disabled").Disabled(true) }), 300, 120)
		tester.SetDark(dark)
		tester.Frame()
		want := pick(dark, ui.RGB(236, 236, 236), ui.RGB(46, 46, 46))
		if got := macPixel(tester, 42, 44); !macNear(got, want, 2) {
			t.Errorf("dark %v: disabled checked box %v, want %v", dark, got, want)
		}
	}
}

// TestMacFocusRing checks AppKit's focus ring around a push button that has the focus from the keyboard: 3.5 points wide, from 3 points outside the button to half a point over its edge, in the theme's Focus color, which reads #8db4fc over white.
func TestMacFocusRing(t *testing.T) {
	tester := ui.NewTester(macAt(40, 40, func(c *ui.Context) { Button(c, "Cancel") }), 300, 120)
	tester.Frame()
	tester.Key(0, ui.KeyTab)
	tester.Frame()
	if got := macPixel(tester, 38, 52); !macNear(got, ui.Hex("#8db4fc"), 2) {
		t.Errorf("focus ring %v, want #8db4fc", got)
	}
	if got := macPixel(tester, 35, 52); got != white {
		t.Errorf("past the focus ring %v, want the window", got)
	}
}

// TestMacDisabledCheckedMarkIsGray checks the mark of a disabled checked box and the dot of a disabled chosen radio button as AppKit draws them in a key window on macOS 27.2: 175 over the 236 box in Aqua and 98 over the 46 box in Dark Aqua. MyGo halves the opacity of each fill on its own, so the mark is composited over the box as shown.
func TestMacDisabledCheckedMarkIsGray(t *testing.T) {
	for _, dark := range []bool{false, true} {
		checked, choice := true, "a"
		tester := ui.NewTester(macAt(40, 40, func(c *ui.Context) {
			ui.Column(c).Gap(24).AlignItems(ui.Start).Children(func() {
				Checkbox(c, &checked, "Disabled").Disabled(true)
				Radio(c, &choice, "a", "Disabled").Disabled(true)
			})
		}), 300, 140)
		tester.SetScale(2)
		tester.SetDark(dark)
		tester.Frame()
		want := pick(dark, ui.RGB(175, 175, 175), ui.RGB(98, 98, 98))
		// The middle row of the box crosses the mark's first stroke 3.5 points in; the radio button's dot is at its center.
		if got := macPixel(tester, 2*40+7, 2*48); !macNear(got, want, 3) {
			t.Errorf("dark %v: disabled check mark %v, want %v", dark, got, want)
		}
		if got := macPixel(tester, 2*48, 2*(40+16+24+8)); !macNear(got, want, 3) {
			t.Errorf("dark %v: disabled radio dot %v, want %v", dark, got, want)
		}
	}
}

// TestMacSliderRingOverKnob checks that the focus ring of a slider's knob is drawn over the knob's edge, as a live capture of a key window on macOS 27.2 shows it in Dark Aqua: the knob's outermost half point reads the ring over white, not white.
func TestMacSliderRingOverKnob(t *testing.T) {
	v := 40.0
	tester := ui.NewTester(macAt(40, 40, func(c *ui.Context) { Slider(c, &v, 0, 100).Width(180) }), 300, 120)
	tester.SetScale(2)
	tester.SetDark(true)
	tester.Frame()
	tester.Key(0, ui.KeyTab)
	tester.Frame()
	// The knob's center is 10 + 0.4 of 160 points in, its edge 10 points after it.
	edge := 2*(40+10+64+10) - 1
	if got := macPixel(tester, edge, 2*48); got == white || !macNear(got, ui.RGB(150, 212, 255), 6) {
		t.Errorf("knob edge under the ring %v, want the ring over white", got)
	}
}
