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
			"fluent stroke":   func() ui.Color { _, s, _ := fluentUnchecked(th, false, false); return s }(),
			"fluent off knob": func() ui.Color { f, _, _ := fluentUnchecked(th, false, false); return fluentOffKnob(th, f) }(),
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

// TestDialogsKeepTheirBehavior checks the alerts and dialogs on every platform: a click returns its button's index whatever the order the platform places the buttons in, Enter clicks the default, Escape dismisses, and a modal closes on Escape.
func TestDialogsKeepTheirBehavior(t *testing.T) {
	for _, p := range Platforms {
		for _, tc := range []struct {
			name string
			act  func(tester *ui.Tester)
			want int
		}{
			{"click Cancel", func(tester *ui.Tester) { _ = tester.Click("Cancel") }, 0},
			{"click Delete", func(tester *ui.Tester) { _ = tester.Click("Delete") }, 1},
			{"Enter", func(tester *ui.Tester) { tester.Key(0, ui.KeyEnter) }, 1},
			{"Escape", func(tester *ui.Tester) { tester.Key(0, ui.KeyEscape) }, ui.AlertDismissed},
		} {
			open, got := true, -1
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					if r := AlertDialog(c, &open, "Delete Notes?", "You can't undo this.", "Cancel", "Delete"); r != -1 {
						got = r
					}
				})
			}, 600, 400)
			tester.Frame()
			if !tester.HasText("Delete Notes?") || !tester.HasText("You can't undo this.") {
				t.Fatalf("%v: the alert's text is missing", p)
			}
			tc.act(tester)
			tester.Frame()
			if got != tc.want || open {
				t.Fatalf("%v %s: got %d, open %v; want %d, closed", p, tc.name, got, open, tc.want)
			}
		}
		open, font := true, ""
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Modal(c, &open, func() { Combobox(c, &font, []string{"Serif", "Sans"}).Label("Font") })
			})
		}, 600, 400)
		tester.Frame()
		r, ok := tester.Find("Font")
		if !ok {
			t.Fatalf("%v: the modal's combobox is missing", p)
		}
		tester.ClickAt(r.X+r.W/2, r.Y+r.H/2)
		tester.Frame()
		tester.Type("Mono")
		tester.Frame()
		tester.Key(0, ui.KeyEscape)
		tester.Frame()
		tester.Key(0, ui.KeyEscape)
		tester.Frame()
		if open {
			t.Fatalf("%v: Escape left the modal open", p)
		}
	}
}

// TestRemainingWidgetsKeepTheirBehavior checks the widgets restyled from package ui's on every platform: toggles in a group still toggle, a stepper steps with Up, a link opens its URL, a menu button opens its menu, and a popover opens beside its anchor and closes on Escape.
func TestRemainingWidgetsKeepTheirBehavior(t *testing.T) {
	for _, p := range Platforms {
		bold, italic, count, open := false, false, 5.0, false
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				ui.Column(c).Gap(8).Padding(8).Children(func() {
					ToggleGroup(c, func() {
						Toggle(c, &bold, "B")
						Toggle(c, &italic, "I")
					}).Label("Style")
					Stepper(c, &count, 0, 10, 1).Label("Count")
					Link(c, "Docs", "https://example.com/docs")
					Divider(c)
					Badge(c, "3")
					MenuButton(c, "Sort", func(m *ui.Menu) { m.Item("Name") })
					anchor := Button(c, "Info")
					if anchor.Clicked() {
						open = true
					}
					Popover(c, anchor, &open, func() { ui.Text(c, "Details") })
					Field(c, "Name", func() { ui.Text(c, "Ada") })
				})
			})
		}, 480, 600)
		tester.Frame()
		for _, s := range []string{"B", "I"} {
			if err := tester.Click(s); err != nil {
				t.Fatalf("%v: click %q: %v", p, s, err)
			}
			tester.Frame()
		}
		if !bold || !italic {
			t.Fatalf("%v: grouped toggles bold %v italic %v", p, bold, italic)
		}
		r, ok := tester.Find("Count")
		if !ok {
			t.Fatalf("%v: stepper not found", p)
		}
		tester.ClickAt(r.X+r.W/2, r.Y+2)
		tester.Frame()
		tester.Key(0, ui.KeyUp)
		tester.Frame()
		if count < 6 {
			t.Fatalf("%v: stepper value %v after a click on its up arrow and Up", p, count)
		}
		if err := tester.Click("Docs"); err != nil {
			t.Fatalf("%v: click the link: %v", p, err)
		}
		tester.Frame()
		if urls := tester.OpenedURLs(); len(urls) != 1 || urls[0] != "https://example.com/docs" {
			t.Fatalf("%v: the link opened %v", p, urls)
		}
		if err := tester.Click("Sort"); err != nil {
			t.Fatalf("%v: click the menu button: %v", p, err)
		}
		tester.Frame()
		if m := tester.Menu(); len(m) == 0 || m[0] != "Name" {
			t.Fatalf("%v: the menu shows %v", p, m)
		}
		tester.CloseMenu()
		tester.Frame()
		if err := tester.Click("Info"); err != nil {
			t.Fatalf("%v: click the popover's anchor: %v", p, err)
		}
		tester.Frame()
		if !open || !tester.HasText("Details") {
			t.Fatalf("%v: the popover did not open", p)
		}
		tester.Key(0, ui.KeyEscape)
		tester.Frame()
		if open {
			t.Fatalf("%v: Escape left the popover open", p)
		}
		if !tester.HasText("Ada") || !tester.HasText("3") {
			t.Fatalf("%v: the field's or the badge's text is missing", p)
		}
	}
}
