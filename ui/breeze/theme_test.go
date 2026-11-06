package breeze

import (
	"fmt"
	"image"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

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

// themeOfAccent is the app accent the ThemeOf tests give: a purple no platform defaults to.
var themeOfAccent = ui.Hex("#8e44ad")

// themeBeforeThemeOf is the theme Theme returned before ThemeOf existed, written out from that code: the platform's theme of the window's appearance from the user's accent made opaque, with the higher contrast and text size adjustments laid on it.
func themeBeforeThemeOf(p Platform, dark bool, prefs ui.Preferences) *ui.Theme {
	accent := prefs.Accent
	if accent.A > 0 {
		accent.A = 255
	}
	var t *ui.Theme
	switch p {
	case Windows:
		t = fluentTheme(dark, accent)
	case Linux:
		t = adwTheme(dark, accent)
	default:
		t = macTheme(dark, accent)
	}
	if prefs.HighContrast {
		t.Border = t.Border.Mix(t.Text, 0.45)
		t.TextMuted = t.TextMuted.Mix(t.Text, 0.4)
		t.Focus.A = 255
	}
	if s := prefs.TextScale; s > 0 && s != 1 {
		t.FontSize *= s
	}
	return t
}

// TestThemeOfMatchesTheme checks that Theme, and ThemeOf with the window's appearance and the zero accent, return what Theme returned before ThemeOf existed, on every platform, light and dark, without preferences, with a translucent user accent, higher contrast and text size, and with an opaque user accent.
func TestThemeOfMatchesTheme(t *testing.T) {
	for _, prefs := range []ui.Preferences{{TextScale: 1}, {Accent: ui.RGBA(255, 45, 85, 0.5), HighContrast: true, TextScale: 1.25}, {Accent: ui.Hex("#34c759"), TextScale: 0.9}} {
		for _, p := range Platforms {
			for _, dark := range []bool{false, true} {
				var theme, of *ui.Theme
				tester := ui.NewTester(func(c *ui.Context) {
					theme = Theme(c, p)
					of = ThemeOf(c, p, c.Theme().Dark, ui.Color{})
				}, 100, 100)
				tester.SetPreferences(prefs)
				tester.SetDark(dark)
				tester.Frame()
				want := themeBeforeThemeOf(p, dark, prefs)
				if !reflect.DeepEqual(theme, want) {
					t.Errorf("%v dark %v, preferences %+v: Theme %+v, want %+v", p, dark, prefs, *theme, *want)
				}
				if !reflect.DeepEqual(of, want) {
					t.Errorf("%v dark %v, preferences %+v: ThemeOf %+v, want %+v", p, dark, prefs, *of, *want)
				}
			}
		}
	}
}

// TestThemeOfAppearanceAndAccent checks that ThemeOf gives the theme of the appearance asked for whatever the window's is, the one a window of that appearance gets from Theme, and that an accent given to it replaces the user's, made opaque, with the platform's derivations: macOS and GNOME fill the accent controls with it, Windows with its AccentFillColorDefault shade.
func TestThemeOfAppearanceAndAccent(t *testing.T) {
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			var want *ui.Theme
			probe := ui.NewTester(func(c *ui.Context) { want = Theme(c, p) }, 100, 100)
			probe.SetDark(dark)
			probe.Frame()
			var other, custom *ui.Theme
			tester := ui.NewTester(func(c *ui.Context) {
				other = ThemeOf(c, p, dark, ui.Color{})
				custom = ThemeOf(c, p, dark, themeOfAccent.Alpha(0.5))
			}, 100, 100)
			tester.SetPreferences(ui.Preferences{Accent: ui.Hex("#ff2d55"), TextScale: 1})
			tester.SetDark(!dark)
			tester.Frame()
			if other.Dark != dark || other.Background != want.Background || other.Text != want.Text || other.Surface != want.Surface {
				t.Errorf("%v dark %v in a window of the other appearance: background %v text %v surface %v, want %v %v %v", p, dark, other.Background, other.Text, other.Surface, want.Background, want.Text, want.Surface)
			}
			fill := themeOfAccent
			if p == Windows {
				fill = fluentAccentFill(themeOfAccent, dark)
			}
			if custom.Accent != fill {
				t.Errorf("%v dark %v: accent %v, want %v from the app's accent", p, dark, custom.Accent, fill)
			}
			if custom.AccentPressed == want.AccentPressed {
				t.Errorf("%v dark %v: the pressed fill %v did not follow the app's accent", p, dark, custom.AccentPressed)
			}
		}
	}
}

// themeOfControls builds a button and a segmented control named after where they are, so that a test finds them inside and outside a subtree.
func themeOfControls(c *ui.Context, where string) {
	view := 1
	ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
		Button(c, where)
		Segmented(c, &view, "Day", "Week").Label(where + " view")
	})
}

// themeOfView builds, in the platform's window theme, controls outside and controls inside ui.Themed with ThemeOf's dark theme.
func themeOfView(p Platform) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Column(c).Fill().Padding(16).Gap(16).Children(func() {
				themeOfControls(c, "Outside")
				ui.Themed(c, ThemeOf(c, p, true, ui.Color{})).Padding(16).Children(func() { themeOfControls(c, "Inside") })
			})
		})
	}
}

// crop returns the pixels of img under r, a box in DIPs at scale 1.
func crop(img *image.RGBA, r ui.Rect) []uint8 {
	var out []uint8
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			px := img.RGBAAt(x, y)
			out = append(out, px.R, px.G, px.B, px.A)
		}
	}
	return out
}

// TestThemeOfDarkSubtree checks ThemeOf inside ui.Themed on every platform: in a light window the button and the segmented control inside draw pixel for pixel as in a dark window, and the ones outside keep the light window's faces.
func TestThemeOfDarkSubtree(t *testing.T) {
	for _, p := range Platforms {
		shots := map[bool]*ui.Tester{}
		for _, dark := range []bool{false, true} {
			tester := ui.NewTester(themeOfView(p), 480, 200)
			tester.SetDark(dark)
			tester.Frame()
			shots[dark] = tester
		}
		light, dark := shots[false], shots[true]
		for _, name := range []string{"Inside", "Inside view", "Outside", "Outside view"} {
			r, ok := light.Find(name)
			if !ok {
				t.Fatalf("%v: %q is missing", p, name)
			}
			same := string(crop(light.Image(), r)) == string(crop(dark.Image(), r))
			if inside := strings.HasPrefix(name, "Inside"); same != inside {
				t.Errorf("%v: %q in a light window draws as in a dark one: %v, want %v", p, name, same, inside)
			}
		}
		// The button's face just before its label: dark inside, light outside.
		face := func(name string) float32 {
			r, _ := light.Find(name)
			px := light.Image().RGBAAt(int(r.X-3), int(r.Y+r.H/2))
			return lum(ui.RGB(px.R, px.G, px.B))
		}
		if in, out := face("Inside"), face("Outside"); in > 0.35 || out < 0.65 {
			t.Errorf("%v: the button's face reads %.2f inside the dark subtree and %.2f outside it", p, in, out)
		}
	}
}

// hasColor reports whether any pixel of img inside r, a box in DIPs at scale 1, is col.
func hasColor(img *image.RGBA, r ui.Rect, col ui.Color) bool {
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if px := img.RGBAAt(x, y); px.R == col.R && px.G == col.G && px.B == col.B && px.A == 255 {
				return true
			}
		}
	}
	return false
}

// TestFluentThemeOfAccent checks that on Windows an InfoBar, a toast and a link inside ThemeOf with an app's accent take that accent, not the system's: the informational icon is SystemAccentColor in the light theme and its Light2 shade in the dark one, and a link's text AccentTextFillColorPrimary, the Dark2 or Light3 shade.
func TestFluentThemeOfAccent(t *testing.T) {
	for _, dark := range []bool{false, true} {
		saved := false
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Windows, func() {
				Use(c)
				ui.Themed(c, ThemeOf(c, Windows, c.Theme().Dark, themeOfAccent)).Fill().Padding(16).Gap(12).Children(func() {
					InfoBar(c, "info", "Update ready", "Restart to finish.", nil)
					Link(c, "Release notes", "https://example.com")
					if Button(c, "Save").Clicked() {
						saved = true
						c.Toast("Saved")
					}
					Toasts(c)
				})
			})
		}, 600, 400)
		tester.SetDark(dark)
		tester.Frame()
		if err := tester.Click("Save"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		tester.Advance(time.Second)
		tester.Frame()
		if !saved || !tester.HasText("Saved") {
			t.Fatalf("dark %v: the toast did not show", dark)
		}
		icon := pick(dark, themeOfAccent, fluentShade(themeOfAccent, 2))
		link := pick(dark, fluentShade(themeOfAccent, -2), fluentShade(themeOfAccent, 3))
		img := tester.Image()
		for _, tc := range []struct {
			name string
			want ui.Color
		}{{"Update ready", icon}, {"Saved", icon}, {"Release notes", link}} {
			r, ok := tester.Find(tc.name)
			if !ok {
				t.Fatalf("dark %v: %q is missing", dark, tc.name)
			}
			// The bar's icon is before its title, on the same row; the link's color is its text's.
			row := ui.Rect{X: 0, Y: r.Y - 16, W: r.X + r.W, H: r.H + 32}
			if tc.name == "Release notes" {
				row = r
			}
			if !hasColor(img, row, tc.want) {
				t.Errorf("dark %v: %q shows no %v from the app's accent", dark, tc.name, tc.want)
			}
		}
	}
}

// themeOfShotView builds the controls ThemeOf themes in a window of p's theme: a row in the window's theme, a row in the other appearance and a row in the window's appearance with an app's purple accent, the last with a link and an informational bar, which on Windows take the accent's shades.
func themeOfShotView(p Platform) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			on, view := true, 1
			row := func(name string) {
				ui.Text(c, name).FontSize(12).TextColor(c.Theme().TextMuted)
				ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
					Button(c, "Cancel")
					PrimaryButton(c, "Save")
					Segmented(c, &view, "Day", "Week")
					Switch(c, &on).Label(name + " switch")
				})
			}
			dark := c.Theme().Dark
			ui.Column(c).Fill().Padding(16).Gap(12).Children(func() {
				row("Window theme")
				ui.Themed(c, ThemeOf(c, p, !dark, ui.Color{})).Padding(16).Gap(8).Radius(8).Children(func() {
					row(pick(dark, "ThemeOf dark", "ThemeOf light"))
				})
				ui.Themed(c, ThemeOf(c, p, dark, themeOfAccent)).Padding(16).Gap(8).Radius(8).Border(1, c.Theme().Border).Children(func() {
					row("ThemeOf with an app accent")
					Link(c, "Release notes", "https://example.com")
					InfoBar(c, "info", "Update ready", "Restart to finish.", nil)
				})
			})
		})
	}
}

// TestScreenshotsThemeOf writes themeOfShotView on every platform, in a light and a dark window, at 2x to the directory BREEZE_SCREENSHOTS names, as breeze-themeof-<platform>-<appearance>.png. Without it the test only builds the frames.
func TestScreenshotsThemeOf(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			tester := ui.NewTester(themeOfShotView(p), 560, 400)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.Frame()
			tester.Frame()
			for _, s := range []string{"Window theme", "ThemeOf with an app accent", "Release notes", "Update ready"} {
				if !tester.HasText(s) {
					t.Fatalf("%v: %q is missing", p, s)
				}
			}
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			writeShot(t, dir, fmt.Sprintf("breeze-themeof-%s-%s.png", strings.ToLower(p.String()), mode), tester)
		}
	}
}

// TestFluentThemeAccent checks how fluentSystemAccent tells a Windows theme's SystemAccentColor: fluentTheme's themes, light and dark, with the default accent or another, give theirs back, while package ui's themes and a Windows theme whose Accent an app replaced do not, so that they keep the user's accent.
func TestFluentThemeAccent(t *testing.T) {
	for _, dark := range []bool{false, true} {
		for _, accent := range []ui.Color{fluentAccent, themeOfAccent} {
			if got, ok := fluentThemeAccent(fluentTheme(dark, accent)); !ok || got != accent {
				t.Errorf("dark %v: fluentTheme with %v gives %v, %v", dark, accent, got, ok)
			}
		}
		replaced := fluentTheme(dark, fluentAccent)
		replaced.Accent = themeOfAccent
		for name, th := range map[string]*ui.Theme{"package ui's": pick(dark, ui.LightTheme(), ui.DarkTheme()), "replaced": replaced} {
			if got, ok := fluentThemeAccent(th); ok {
				t.Errorf("dark %v: %s theme gives the accent %v", dark, name, got)
			}
		}
	}
}
