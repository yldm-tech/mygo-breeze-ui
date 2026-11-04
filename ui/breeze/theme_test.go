package breeze

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func lum(c ui.Color) float32 { return luminance(c) }

// themes returns every platform's theme, light and dark, as a window without preferences gets it.
func themes(t *testing.T) map[string]*ui.Theme {
	t.Helper()
	out := map[string]*ui.Theme{}
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			var got *ui.Theme
			tester := ui.NewTester(func(c *ui.Context) { got = Theme(c, p) }, 100, 100)
			tester.SetDark(dark)
			tester.Frame()
			if got.Dark != dark {
				t.Fatalf("%v theme Dark = %v, want %v", p, got.Dark, dark)
			}
			out[p.String()+map[bool]string{false: " light", true: " dark"}[dark]] = got
		}
	}
	return out
}

// TestThemeContrast checks that text and the accent stand out, darker on light themes and lighter on dark ones.
func TestThemeContrast(t *testing.T) {
	for name, th := range themes(t) {
		bg := lum(th.Background)
		if d := lum(th.Text) - bg; (th.Dark && d < 0.5) || (!th.Dark && d > -0.5) {
			t.Errorf("%s: text luminance %.2f against background %.2f", name, lum(th.Text), bg)
		}
		if d := lum(th.AccentText) - lum(th.Accent); d < 0.3 && d > -0.3 {
			t.Errorf("%s: accent text %.2f does not stand out on the accent %.2f", name, lum(th.AccentText), lum(th.Accent))
		}
	}
}

// TestTroughsStandOut checks the fills that show an unchecked box and an off switch: each must differ from the window background toward the text, as every platform draws them, or the control is invisible.
func TestTroughsStandOut(t *testing.T) {
	for name, th := range themes(t) {
		parts := map[string]ui.Color{
			"mac trough":      tint(th.Background, th.Dark, macTroughAlpha),
			"mac face":        macFace(th.Background, th.Dark),
			"adwaita trough":  adwTrough(th, false, false),
			"fluent stroke":   func() ui.Color { _, s := fluentUnchecked(th, false, false); return s }(),
			"fluent off knob": fluentOffKnob(th),
		}
		for part, c := range parts {
			d := lum(c) - lum(th.Background)
			if (th.Dark && d < 0.04) || (!th.Dark && d > -0.04) {
				t.Errorf("%s: %s luminance %.3f against background %.3f", name, part, lum(c), lum(th.Background))
			}
		}
	}
}

func TestPlatformFromSystem(t *testing.T) {
	for goos, want := range map[string]Platform{"darwin": Mac, "windows": Windows, "linux": Linux, "freebsd": Linux} {
		if got := system(goos); got != want {
			t.Errorf("system(%q) = %v, want %v", goos, got, want)
		}
	}
}

func TestPreviewOverridesAndRestores(t *testing.T) {
	var outer, inner, after Platform
	tester := ui.NewTester(func(c *ui.Context) {
		Preview(c, Windows, func() {
			outer = PlatformOf(c)
			Preview(c, Linux, func() { inner = PlatformOf(c) })
			after = PlatformOf(c)
		})
	}, 100, 100)
	tester.Frame()
	if outer != Windows || inner != Linux || after != Windows {
		t.Fatalf("Preview platforms = %v, %v, %v; want Windows, Linux, Windows", outer, inner, after)
	}
	SetPlatform(Linux)
	defer SetPlatform(Auto)
	if Current() != Linux {
		t.Fatalf("Current() = %v after SetPlatform(Linux)", Current())
	}
}

// TestWidgetsKeepTheirBehavior checks that the drawn widgets still click, toggle and refuse clicks when disabled, as package ui's do, on every platform.
func TestWidgetsKeepTheirBehavior(t *testing.T) {
	for _, p := range Platforms {
		clicks, checked, on := 0, false, false
		disabledClicks := 0
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				ui.Column(c).Gap(8).Children(func() {
					if Button(c, "Go").Clicked() {
						clicks++
					}
					if PrimaryButton(c, "Off").Disabled(true).Clicked() {
						disabledClicks++
					}
					Checkbox(c, &checked, "Remember")
					Switch(c, &on).Label("Sync")
				})
			})
		}, 400, 300)
		tester.Frame()
		for _, s := range []string{"Go", "Off", "Remember"} {
			if err := tester.Click(s); err != nil {
				t.Fatalf("%v: click %q: %v", p, s, err)
			}
			tester.Frame()
		}
		if clicks != 1 || disabledClicks != 0 || !checked {
			t.Fatalf("%v: clicks %d, disabled clicks %d, checked %v", p, clicks, disabledClicks, checked)
		}
		r, ok := tester.Find("Sync")
		if !ok {
			t.Fatalf("%v: switch named Sync not found", p)
		}
		tester.ClickAt(r.X+r.W/2, r.Y+r.H/2)
		tester.Frame()
		if !on {
			t.Fatalf("%v: switch did not turn on", p)
		}
	}
}
