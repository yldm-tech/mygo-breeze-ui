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

// TestMoreWidgetsKeepTheirBehavior checks the widgets drawn on MyGo's other bases on every platform: radio buttons choose, toggles toggle, segments and tabs choose, a select opens and chooses with the keyboard, a slider follows a click, and the fields edit their text.
func TestMoreWidgetsKeepTheirBehavior(t *testing.T) {
	for _, p := range Platforms {
		size, bold, view, tab := "S", false, 0, 0
		choice, name, query, volume := "Small", "", "", 0.0
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				ui.Column(c).Gap(8).Padding(8).Children(func() {
					RadioGroup(c, func() {
						Radio(c, &size, "S", "Small size")
						Radio(c, &size, "L", "Large size")
					})
					Toggle(c, &bold, "Bold")
					Segmented(c, &view, "Day", "Week")
					Tabs(c, &tab, "General", "Advanced")
					Select(c, &choice, []string{"Small", "Medium", "Large"}).Label("Size")
					Slider(c, &volume, 0, 100).Width(200).Label("Volume")
					TextInput(c, &name).Label("Name")
					SearchField(c, &query).Label("Search")
					Progress(c, 0.5)
					Spinner(c)
				})
			})
		}, 480, 800)
		tester.Frame()
		for _, s := range []string{"Large size", "Bold", "Week", "Advanced"} {
			if err := tester.Click(s); err != nil {
				t.Fatalf("%v: click %q: %v", p, s, err)
			}
			tester.Frame()
		}
		if size != "L" || !bold || view != 1 || tab != 1 {
			t.Fatalf("%v: size %q, bold %v, view %d, tab %d", p, size, bold, view, tab)
		}
		if err := tester.Click("Small"); err != nil {
			t.Fatalf("%v: open the select: %v", p, err)
		}
		tester.Frame()
		tester.Key(0, ui.KeyDown)
		tester.Frame()
		tester.Key(0, ui.KeyEnter)
		tester.Frame()
		if choice != "Medium" {
			t.Fatalf("%v: select chose %q, want Medium", p, choice)
		}
		r, ok := tester.Find("Volume")
		if !ok {
			t.Fatalf("%v: slider named Volume not found", p)
		}
		tester.ClickAt(r.X+r.W-2, r.Y+r.H/2)
		tester.Frame()
		if volume < 90 {
			t.Fatalf("%v: a click at the slider's end set %v", p, volume)
		}
		for label, want := range map[string]*string{"Name": &name, "Search": &query} {
			r, ok := tester.Find(label)
			if !ok {
				t.Fatalf("%v: field %q not found", p, label)
			}
			tester.ClickAt(r.X+r.W/2, r.Y+r.H/2)
			tester.Frame()
			tester.Type("abc")
			tester.Frame()
			if *want != "abc" {
				t.Fatalf("%v: typing in %q gave %q", p, label, *want)
			}
		}
	}
}
