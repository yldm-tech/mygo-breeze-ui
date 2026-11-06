package breeze

import (
	"fmt"
	"image"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func init() {
	for _, n := range []string{"ShortcutLabel", "Keycaps", "ShortcutText"} {
		drawnWidgets[n] = true
	}
}

// shortcutCase is a shortcut and how each platform writes it: macOS's menus, WinUI's KeyboardAccelerator and GTK's accelerator label.
type shortcutCase struct {
	mods                ui.Modifiers
	key                 ui.Key
	mac, windows, gnome string
}

// shortcutCases are shortcuts written with ui.Cmd, Alt and Shift, which read the same on every host, and keys each platform names its own way.
var shortcutCases = []shortcutCase{
	{ui.Cmd | ui.Shift, ui.KeyL, "⇧⌘L", "Ctrl+Shift+L", "Shift+Ctrl+L"},
	{ui.Cmd, ui.KeyK, "⌘K", "Ctrl+K", "Ctrl+K"},
	{ui.Cmd | ui.Alt | ui.Shift, ui.KeyS, "⌥⇧⌘S", "Ctrl+Alt+Shift+S", "Shift+Ctrl+Alt+S"},
	{ui.Alt | ui.Cmd, ui.KeyLeft, "⌥⌘←", "Ctrl+Alt+Left", "Ctrl+Alt+Left"},
	{0, ui.KeySpace, "Space", "Space", "Space"},
	{ui.Cmd, ui.KeyEnter, "⌘↩", "Ctrl+Enter", "Ctrl+Return"},
	{0, ui.KeyEscape, "⎋", "Esc", "Escape"},
	{ui.Shift, ui.KeyBackspace, "⇧⌫", "Shift+Back", "Shift+BackSpace"},
	{0, ui.KeyDelete, "⌦", "Delete", "Delete"},
	{0, ui.KeyTab, "⇥", "Tab", "Tab"},
	{0, ui.KeyPageUp, "⇞", "PageUp", "Page Up"},
	{0, ui.KeyHome, "↖", "Home", "Home"},
	{0, ui.KeyF5, "F5", "F5", "F5"},
	{ui.Cmd, ui.Key0, "⌘0", "Ctrl+0", "Ctrl+0"},
	{ui.Cmd, ui.KeyComma, "⌘,", "Ctrl+,", "Ctrl+,"},
	{ui.Cmd, ui.KeyBackslash, "⌘\\", "Ctrl+\\", "Ctrl+Backslash"},
	{ui.Cmd, ui.KeyBracketLeft, "⌘[", "Ctrl+[", "Ctrl+["},
	{0, ui.KeyContextMenu, "Menu", "Application", "Menu"},
}

// shortcutTextOn returns what ShortcutText gives mods and key in a Preview of p.
func shortcutTextOn(p Platform, mods ui.Modifiers, key ui.Key) string {
	var got string
	tester := ui.NewTester(func(c *ui.Context) {
		Preview(c, p, func() { got = ShortcutText(c, mods, key) })
	}, 100, 100)
	tester.Frame()
	return got
}

// TestShortcutTextPerPlatform checks ShortcutText against each platform's own writing of shortcutCases: macOS's glyphs in the order ⌃⌥⇧⌘, WinUI's names in the order Ctrl, Alt, Windows, Shift, GTK's in the order Shift, Ctrl, Alt, Super. It builds them in a Preview of each platform, so on any host it checks that the platform drawn, not runtime.GOOS, decides, and that ui.Cmd stands for that platform's command key.
func TestShortcutTextPerPlatform(t *testing.T) {
	for _, s := range shortcutCases {
		for p, want := range map[Platform]string{Mac: s.mac, Windows: s.windows, Linux: s.gnome} {
			if got := shortcutTextOn(p, s.mods, s.key); got != want {
				t.Errorf("%v: ShortcutText(%v, %v) = %q, want %q", p, s.mods, s.key, got, want)
			}
		}
	}
}

// TestShortcutTextOtherModifiers checks the modifier that is not ui.Cmd on the host, which keeps its own name: Control on a Mac, whose Cmd is Super, and Super on Windows and Linux, whose Cmd is Ctrl. WinUI names Super Windows and puts it before Shift; GTK names it Super and puts it last; AdwShortcutLabel puts it first.
func TestShortcutTextOtherModifiers(t *testing.T) {
	other := ui.Ctrl
	want := map[Platform]string{Mac: "⌃⇧⇥", Windows: "Ctrl+Shift+Tab", Linux: "Shift+Ctrl+Tab"}
	caps := map[Platform][]string{Mac: {"⌃⇧⇥"}, Windows: {"Ctrl", "Shift", "Tab"}, Linux: {"Ctrl", "Shift", "Tab"}}
	if ui.Cmd == ui.Ctrl {
		other = ui.Super
		want = map[Platform]string{Mac: "⇧⌘⇥", Windows: "Windows+Shift+Tab", Linux: "Shift+Super+Tab"}
		caps = map[Platform][]string{Mac: {"⇧⌘⇥"}, Windows: {"Windows", "Shift", "Tab"}, Linux: {"Super", "Shift", "Tab"}}
	}
	for _, p := range Platforms {
		if got := shortcutTextOn(p, other|ui.Shift, ui.KeyTab); got != want[p] {
			t.Errorf("%v: ShortcutText(%v|Shift, Tab) = %q, want %q", p, other, got, want[p])
		}
		got := shortcutParts(p, other|ui.Shift, ui.KeyTab, true)
		if p == Mac {
			got = []string{strings.Join(got, "")}
		}
		if strings.Join(got, " ") != strings.Join(caps[p], " ") {
			t.Errorf("%v: the keycaps of %v|Shift+Tab are %q, want %q", p, other, got, caps[p])
		}
	}
}

// TestShortcutTextUnnamedKey checks that a key the platforms do not name, as KeyBack, gives no text and ShortcutLabel no keycap.
func TestShortcutTextUnnamedKey(t *testing.T) {
	for _, p := range Platforms {
		if got := shortcutTextOn(p, ui.Cmd, ui.KeyBack); got != "" {
			t.Errorf("%v: ShortcutText(Cmd, KeyBack) = %q, want nothing", p, got)
		}
		var box ui.Rect
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() { box = ShortcutLabel(c, ui.Cmd, ui.KeyBack).Bounds() })
			})
		}, 200, 100)
		tester.Frame()
		tester.Frame()
		if box.W != 0 || len(tester.Texts()) != 0 {
			t.Errorf("%v: ShortcutLabel(Cmd, KeyBack) shows %v, %q", p, box, tester.Texts())
		}
	}
}

// shortcutLabelCases are what ShortcutLabel(ui.Cmd|ui.Shift, ui.KeyL) shows on a platform: the text of each keycap and the name assistive technology hears.
var shortcutLabelCases = map[Platform]struct {
	caps   []string
	spoken string
}{
	Mac:     {[]string{"⇧⌘L"}, "Shift Command L"},
	Windows: {[]string{"Ctrl", "Shift", "L"}, "Ctrl+Shift+L"},
	Linux:   {[]string{"Ctrl", "Shift", "L"}, "Shift+Ctrl+L"},
}

// keycapsApart returns the DIPs p puts between the end of a keycap's text, prev, and the start of the next one's, next: the padding and border after the one, the gap and the padding and border before the other, a text narrower than the keycap's minimum width (22 on macOS, 20 on GNOME, none on Windows) being centered in it.
func keycapsApart(p Platform, prev, next ui.Rect) float32 {
	side := func(w float32) float32 {
		switch p {
		case Windows:
			return 6 + 1
		case Linux:
			return max(6, (20-w)/2)
		}
		return max(8, (22-w)/2)
	}
	return side(prev.W) + pick(p == Linux, float32(4), 6) + side(next.W)
}

// TestShortcutLabelKeycaps checks ShortcutLabel(Cmd|Shift, L) on every platform: one keycap holding ⇧⌘L on macOS, three on Windows and GNOME (Ctrl, Shift, L), their texts as far apart as the platform's padding, border and gap put them, in a row of the platform's keycap height that assistive technology finds by the spoken shortcut, that takes no focus (Tab goes to the button after it), and that stays left to right in a right-to-left window.
func TestShortcutLabelKeycaps(t *testing.T) {
	for _, p := range Platforms {
		want := shortcutLabelCases[p]
		for _, rtl := range []bool{false, true} {
			var row ui.Rect
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					ui.Column(c).Padding(16).Gap(8).AlignItems(ui.Start).Children(func() {
						row = ShortcutLabel(c, ui.Cmd|ui.Shift, ui.KeyL).Bounds()
						Button(c, "Done")
					})
				})
			}, 300, 160)
			if rtl {
				tester.SetDirection(ui.RightToLeft)
			}
			tester.Frame()
			tester.Frame()
			if r, ok := tester.Find(want.spoken); !ok || !rectNear(r, row) {
				t.Errorf("%v, rtl %v: the shortcut is not named %q (%v, the row %v)", p, rtl, want.spoken, r, row)
			}
			var boxes []ui.Rect
			for _, k := range want.caps {
				r, ok := tester.Find(k)
				if !ok {
					t.Fatalf("%v, rtl %v: keycap %q is missing: %q", p, rtl, k, tester.Texts())
				}
				boxes = append(boxes, r)
			}
			if n := len(tester.Texts()); n != len(want.caps)+2 {
				t.Errorf("%v, rtl %v: %d texts, want the %d keycaps, the row's name and the button: %q", p, rtl, n, len(want.caps), tester.Texts())
			}
			for i := 1; i < len(boxes); i++ {
				if d, apart := boxes[i].X-(boxes[i-1].X+boxes[i-1].W), keycapsApart(p, boxes[i-1], boxes[i]); math.Abs(float64(d-apart)) > 0.01 {
					t.Errorf("%v, rtl %v: %q starts %v after %q ends, want %v", p, rtl, want.caps[i], d, want.caps[i-1], apart)
				}
			}
			height := map[Platform]float32{Mac: 22, Windows: boxes[0].H + 14, Linux: boxes[0].H + 12}[p]
			if math.Abs(float64(row.H-height)) > 0.01 {
				t.Errorf("%v, rtl %v: the keycaps are %v high, want %v", p, rtl, row.H, height)
			}
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			if !tester.Focused("Done") {
				t.Errorf("%v, rtl %v: Tab did not go past the keycaps to the button", p, rtl)
			}
		}
	}
}

// TestKeycapsShowsKeysAsGiven checks Keycaps on every platform: a keycap for each key, macOS's too, in the order given, and a name of the keys separated by spaces.
func TestKeycapsShowsKeysAsGiven(t *testing.T) {
	for _, p := range Platforms {
		var row ui.Rect
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() { row = Keycaps(c, "Z", "X", "⌘K").Bounds() })
			})
		}, 300, 100)
		tester.Frame()
		tester.Frame()
		if r, ok := tester.Find("Z X ⌘K"); !ok || !rectNear(r, row) {
			t.Errorf("%v: Keycaps is not named \"Z X ⌘K\": %v, the row %v", p, r, row)
		}
		z, okZ := tester.Find("Z")
		x, okX := tester.Find("X")
		k, okK := tester.Find("⌘K")
		if !okZ || !okX || !okK || !(z.X < x.X && x.X < k.X) {
			t.Errorf("%v: the keycaps are not Z, X, ⌘K in order: %v %v %v (%q)", p, z, x, k, tester.Texts())
		}
		if d, apart := x.X-(z.X+z.W), keycapsApart(p, z, x); math.Abs(float64(d-apart)) > 0.01 {
			t.Errorf("%v: X starts %v after Z ends, want %v", p, d, apart)
		}
	}
}

// keycapPixel returns the pixel of img, drawn at scale, at (x, y) in DIPs.
func keycapPixel(img *image.RGBA, scale, x, y float32) ui.Color {
	o := img.PixOffset(int(x*scale), int(y*scale))
	return ui.RGB(img.Pix[o], img.Pix[o+1], img.Pix[o+2])
}

// rectNear reports whether a and b are the same box but for rounding.
func rectNear(a, b ui.Rect) bool {
	near := func(x, y float32) bool { return math.Abs(float64(x-y)) < 0.01 }
	return near(a.X, b.X) && near(a.Y, b.Y) && near(a.W, b.W) && near(a.H, b.H)
}

// colorNear reports whether a and b differ by at most tol in every channel.
func colorNear(a, b ui.Color, tol int) bool {
	d := func(x, y uint8) bool { return int(x)-int(y) <= tol && int(y)-int(x) <= tol }
	return d(a.R, b.R) && d(a.G, b.G) && d(a.B, b.B)
}

// TestShortcutLabelPixels checks the keycap's colors on every platform, light and dark, at 2x, in the padding before the first key's text, and the shade along its bottom: macOS's black at 5% (white at 12% in the dark theme) with a 1-point line of black at 8% (30%); Windows's ControlFillColorDefault inside a 1-DIP CardStrokeColorDefault border; GNOME's currentColor at 10% with a 2-pixel shade of the card's shade (7%, 36%), and with higher contrast a 1-pixel ring of currentColor at 50% outside it.
func TestShortcutLabelPixels(t *testing.T) {
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			for _, contrast := range []bool{false, true} {
				if contrast && p != Linux {
					continue
				}
				var th *ui.Theme
				var row ui.Rect
				tester := ui.NewTester(func(c *ui.Context) {
					Preview(c, p, func() {
						Use(c)
						th = c.Theme()
						ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() { row = ShortcutLabel(c, ui.Cmd|ui.Shift, ui.KeyL).Bounds() })
					})
				}, 300, 100)
				tester.SetScale(2)
				tester.SetDark(dark)
				tester.SetPreferences(ui.Preferences{HighContrast: contrast})
				tester.Frame()
				tester.Frame()
				img := tester.Image()
				px := func(x, y float32) ui.Color { return keycapPixel(img, 2, x, y) }
				bg := th.Background
				mid := row.Y + row.H/2
				var face, bottom, edge, edgeWant ui.Color
				var faceWant, bottomWant ui.Color
				switch p {
				case Mac:
					face, bottom = px(row.X+4, mid), px(row.X+12, row.Y+row.H-0.25)
					faceWant = pick(dark, bg.Mix(black, 0.05), bg.Mix(white, 0.12))
					bottomWant = faceWant.Mix(black, pick(dark, float32(0.08), 0.3))
				case Windows:
					_, _, stroke := fluentExpanderCard(th)
					face, edge = px(row.X+3, mid), px(row.X+0.25, mid)
					faceWant, edgeWant = th.Surface, stroke
					bottom, bottomWant = face, faceWant
				case Linux:
					face, bottom = px(row.X+3, mid), px(row.X+10, row.Y+row.H-1)
					faceWant = adwInkAlpha(th, bg, 0.10)
					bottomWant = faceWant.Mix(ui.RGB(0, 0, 6), pick(dark, float32(0.07), 0.36))
					if contrast {
						edge, edgeWant = px(row.X-0.25, mid), adwInkAlpha(th, bg, 0.5)
					}
				}
				name := fmt.Sprintf("%v, dark %v, contrast %v", p, dark, contrast)
				if !colorNear(face, faceWant, 2) {
					t.Errorf("%s: the keycap's face reads %v, want %v", name, face, faceWant)
				}
				if !colorNear(bottom, bottomWant, 2) {
					t.Errorf("%s: the keycap's bottom reads %v, want %v", name, bottom, bottomWant)
				}
				if edgeWant != (ui.Color{}) && !colorNear(edge, edgeWant, 2) {
					t.Errorf("%s: the keycap's edge reads %v, want %v", name, edge, edgeWant)
				}
				if p == Linux && !contrast && !colorNear(px(row.X-0.25, mid), bg, 1) {
					t.Errorf("%s: a ring shows outside the keycap without higher contrast: %v", name, px(row.X-0.25, mid))
				}
			}
		}
	}
}

// shortcutSamples are the shortcuts shortcutView shows, one to a line, as LingoCut's cheat sheet lists them.
var shortcutSamples = []struct {
	label string
	mods  ui.Modifiers
	key   ui.Key
}{
	{"Play or pause", 0, ui.KeySpace},
	{"Command palette", ui.Cmd, ui.KeyK},
	{"Show translation", ui.Cmd | ui.Shift, ui.KeyT},
	{"Previous sentence", ui.Alt | ui.Cmd, ui.KeyLeft},
	{"Export", ui.Cmd, ui.KeyEnter},
	{"Close", 0, ui.KeyEscape},
}

// shortcutView shows shortcutSamples on p, each label before its ShortcutLabel and its ShortcutText, then a chord of Keycaps.
func shortcutView(p Platform) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Column(c).Fill().Padding(16).Gap(10).Children(func() {
				for _, s := range shortcutSamples {
					ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
						Text(c, s.label).Width(140)
						ui.Box(c).Width(170).Children(func() { ShortcutLabel(c, s.mods, s.key) })
						Text(c, ShortcutText(c, s.mods, s.key), TextSecondary)
					})
				}
				ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
					Text(c, "Speed").Width(140)
					Keycaps(c, "[", "]")
				})
			})
		})
	}
}

// TestScreenshotsShortcuts writes shortcutView on every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names, after checking that every shortcut's text and keycaps show. Without it the test only builds the frames.
func TestScreenshotsShortcuts(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			tester := ui.NewTester(shortcutView(p), 480, 330)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.Frame()
			tester.Frame()
			for _, s := range shortcutSamples {
				text := shortcutText(p, s.mods, s.key)
				if !tester.HasText(s.label) || !tester.HasText(text) {
					t.Fatalf("%v: %q or %q is missing: %q", p, s.label, text, tester.Texts())
				}
				for _, k := range shortcutParts(p, s.mods, s.key, true) {
					if p != Mac && !tester.HasText(k) {
						t.Errorf("%v: the keycap %q of %q is missing", p, k, s.label)
					}
				}
			}
			if !tester.HasText("[ ]") {
				t.Errorf("%v: the chord's keycaps are missing: %q", p, tester.Texts())
			}
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			writeShot(t, dir, fmt.Sprintf("breeze-shortcuts-%s-%s.png", strings.ToLower(p.String()), mode), tester)
		}
	}
}
