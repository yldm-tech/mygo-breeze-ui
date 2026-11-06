package breeze

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func init() {
	for _, n := range []string{"PreferencesGroup", "PrefixRow", "ContentRow", "RadioRow"} {
		drawnWidgets[n] = true
	}
}

// prefsTile builds the tile the preferences tests put in a PrefixRow: a 30-point rounded square of systemBlue, as an app's tile, 20 DIPs on Windows, as SettingsCard's HeaderIcon is.
func prefsTile(c *ui.Context) {
	s := float32(30)
	if PlatformOf(c) == Windows {
		s = 20
	}
	ui.Box(c).Size(s, s).Shrink(0).Radius(s * 7 / 30).Background(ui.Hex("#0088ff"))
}

// prefsView lays out the preferences widgets for screenshots, in two columns 300 wide as rowsView's: a group with a header button and a description holding a prefix row with a subtitle and a button, a row without a prefix and a prefix row of one line, then a card of three radio rows; beside them a group without a description holding a row and a content row with a text area.
func prefsView(p Platform) func(c *ui.Context) {
	engine, prompt := "whisper", "Translate the subtitles."
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Row(c).Fill().Padding(16).Gap(24).AlignItems(ui.Start).Children(func() {
				ui.Column(c).Width(300).Gap(24).Children(func() {
					PreferencesGroup(c, "Speech", "Models are stored on this Mac.", func() { Button(c, "Detect") }, func() {
						PrefixRow(c, "Whisper", "Base · 142 MB", func() { prefsTile(c) }, func() { Button(c, "Remove") })
						ActionRow(c, "Language", "", func() { ui.Text(c, "Auto").TextColor(c.Theme().TextMuted) })
						PrefixRow(c, "Tiny", "", func() { prefsTile(c) }, nil)
					})
					Card(c, func() {
						RadioRow(c, &engine, "whisper", "Whisper", "On this Mac", nil)
						RadioRow(c, &engine, "cloud", "Cloud", "Sent to a service", nil)
						RadioRow(c, &engine, "off", "Off", "", nil)
					})
				})
				ui.Column(c).Width(300).Gap(24).Children(func() {
					PreferencesGroup(c, "Translation", "", nil, func() {
						ActionRow(c, "Model", "", func() { ui.Text(c, "GPT").TextColor(c.Theme().TextMuted) })
						ContentRow(c, "Prompt", func() { TextArea(c, &prompt).Height(60) })
					})
				})
			})
		})
	}
}

// TestScreenshotsPreferences writes the preferences widgets of every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names. Without it the test only builds the frames.
func TestScreenshotsPreferences(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			h := 480
			if p == Windows {
				h = 600
			}
			tester := ui.NewTester(prefsView(p), 672, h)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tester.Frame()
			tester.Frame()
			for _, s := range []string{"Speech", "Detect", "Models are stored on this Mac.", "Whisper", "Base · 142 MB", "Remove", "Tiny", "On this Mac", "Off", "Translation", "Prompt"} {
				if !tester.HasText(s) {
					t.Fatalf("%v: %q is missing", p, s)
				}
			}
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			writeShot(t, dir, fmt.Sprintf("breeze-preferences-%s-%s.png", strings.ToLower(p.String()), mode), tester)
		}
	}
}

// TestRadioRowChooses checks the radio rows on every platform, left to right and right to left: a click on a row's subtitle and one on its end choose it, Changed reporting each new choice once and a click on the chosen row none; Space and Enter choose the focused row; Down and Up choose the next and the previous row of the card sharing the value, round from the last to the first, passing over a disabled one, and move the focus there; Tab goes from the last row to the button after the card, the radio button inside the row no stop of its own.
func TestRadioRowChooses(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			engine, changes, disabled := "whisper", map[string]int{}, false
			tester := layoutTester(p, 560, 400, func(c *ui.Context) {
				ui.Column(c).Padding(16).Gap(8).Children(func() {
					Card(c, func() {
						for _, r := range [][3]string{{"whisper", "Whisper", "On this Mac"}, {"cloud", "Cloud", "Sent to a service"}, {"off", "Off", ""}} {
							row := RadioRow(c, &engine, r[0], r[1], r[2], nil)
							if r[0] == "cloud" {
								row.Disabled(disabled)
							}
							if row.Changed() {
								changes[r[0]]++
							}
						}
					})
					Button(c, "After")
				})
			})
			tester.SetDirection(dir)
			tester.Frame()
			step := func(what string, act func(), want string, wantChanges map[string]int) {
				t.Helper()
				act()
				tester.Frame()
				tester.Frame()
				if engine != want {
					t.Fatalf("%v %v: %s chose %q, want %q", p, dir, what, engine, want)
				}
				for k, n := range wantChanges {
					if changes[k] != n {
						t.Fatalf("%v %v: after %s, %q changed %d times, want %d (%v)", p, dir, what, k, changes[k], n, changes)
					}
				}
			}
			click := func(s string) func() {
				return func() {
					if err := tester.Click(s); err != nil {
						t.Fatalf("%v %v: %v", p, dir, err)
					}
				}
			}
			step("a click on the subtitle", click("Sent to a service"), "cloud", map[string]int{"cloud": 1, "whisper": 0})
			step("a click on the chosen row", click("Cloud"), "cloud", map[string]int{"cloud": 1})
			off, ok := tester.Find("Off")
			if !ok {
				t.Fatalf("%v %v: no row is named Off", p, dir)
			}
			// The row's end: the right left to right, the left right to left.
			x := off.X + off.W - 12
			if dir == ui.RightToLeft {
				x = off.X + 12
			}
			step("a click on the row's end", func() { tester.ClickAt(x, off.Y+off.H/2) }, "off", map[string]int{"off": 1})
			if !tester.Focused("Off") {
				t.Fatalf("%v %v: the click did not focus the row", p, dir)
			}
			step("Down from the last row", func() { tester.Key(0, ui.KeyDown) }, "whisper", map[string]int{"whisper": 1})
			if !tester.Focused("Whisper") {
				t.Fatalf("%v %v: Down did not move the focus to the first row", p, dir)
			}
			step("Down", func() { tester.Key(0, ui.KeyDown) }, "cloud", map[string]int{"cloud": 2})
			step("Up", func() { tester.Key(0, ui.KeyUp) }, "whisper", map[string]int{"whisper": 2})
			step("Up from the first row", func() { tester.Key(0, ui.KeyUp) }, "off", map[string]int{"off": 2})
			if !tester.Focused("Off") {
				t.Fatalf("%v %v: Up did not move the focus round to the last row", p, dir)
			}
			disabled = true
			tester.Frame()
			step("Up past a disabled row", func() { tester.Key(0, ui.KeyUp) }, "whisper", map[string]int{"whisper": 3, "cloud": 2})
			step("Down past a disabled row", func() { tester.Key(0, ui.KeyDown) }, "off", map[string]int{"off": 3, "cloud": 2})
			engine = "cloud"
			disabled = false
			tester.Frame()
			step("Space", func() { tester.Key(0, ui.KeySpace) }, "off", map[string]int{"off": 4})
			engine = "cloud"
			tester.Frame()
			step("Enter", func() { tester.Key(0, ui.KeyEnter) }, "off", map[string]int{"off": 5})
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			if !tester.Focused("After") {
				t.Fatalf("%v %v: Tab from the last row did not go to the button after the card", p, dir)
			}
			if !layoutStateEmpty(&rowGroups) || !layoutStateEmpty(&radioFocus) {
				t.Fatalf("%v %v: a group of rows or a focus to move outlived the build", p, dir)
			}
		}
	}
}

// TestRadioRowAloneTakesNoArrows checks on every platform that radio rows built outside a group, which share no list, leave Down to what is around them: the choice stays.
func TestRadioRowAloneTakesNoArrows(t *testing.T) {
	for _, p := range Platforms {
		v := 1
		tester := layoutTester(p, 560, 300, func(c *ui.Context) {
			ui.Column(c).Padding(16).Gap(8).Children(func() {
				RadioRow(c, &v, 1, "One", "", nil)
				RadioRow(c, &v, 2, "Two", "", nil)
			})
		})
		if err := tester.Click("One"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		tester.Key(0, ui.KeyDown)
		tester.Frame()
		tester.Frame()
		if v != 1 {
			t.Errorf("%v: Down on a row alone chose %d", p, v)
		}
	}
}

// TestPrefixRowsJoin checks on every platform, left to right and right to left, that prefix rows join a card's list with rows without a prefix, touching, the line between rows drawn under the prefix's column too, and that the prefix and the title sit where the platform puts them: the prefix 10 points in on macOS and the title 10 after it, 18 DIPs in on Windows (the card's padding of 16 and the HeaderIcon's margin of 2) and the title 20 after it, 14 pixels in on GNOME (the row's padding of 2 and the header's margin of 12) and the title 12 after it.
func TestPrefixRowsJoin(t *testing.T) {
	near := func(a, b float32) bool { return a-b > -0.5 && a-b < 0.5 }
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			var rows [3]ui.Rect
			var tile ui.Rect
			tester := layoutTester(p, 600, 400, func(c *ui.Context) {
				Use(c)
				ui.Column(c).Padding(16).Children(func() {
					Card(c, func() {
						rows[0] = PrefixRow(c, "Whisper", "Base model", func() {
							tile = ui.Box(c).Size(20, 20).Background(ui.Hex("#0088ff")).Bounds()
						}, func() { Button(c, "Remove") }).Bounds()
						rows[1] = ActionRow(c, "Language", "", nil).Bounds()
						rows[2] = PrefixRow(c, "Tiny", "Small model", func() { ui.Box(c).Size(20, 20).Background(ui.Hex("#0088ff")) }, nil).Bounds()
					})
				})
			})
			tester.SetDirection(dir)
			tester.SetScale(2)
			tester.Frame()
			tester.Frame()
			for i := 1; i < len(rows); i++ {
				if gap := rows[i].Y - (rows[i-1].Y + rows[i-1].H); gap < -0.5 || gap > 0.5 {
					t.Errorf("%v %v: rows %d and %d are %v apart", p, dir, i-1, i, gap)
				}
			}
			in, after := map[Platform][2]float32{Mac: {10, 10}, Windows: {18, 20}, Linux: {14, 12}}[p][0], map[Platform][2]float32{Mac: {10, 10}, Windows: {18, 20}, Linux: {14, 12}}[p][1]
			sub, _ := tester.Find("Base model")
			tileIn, subIn := tile.X-rows[0].X, sub.X-(tile.X+tile.W)
			if dir == ui.RightToLeft {
				tileIn, subIn = rows[0].X+rows[0].W-(tile.X+tile.W), tile.X-(sub.X+sub.W)
			}
			if !near(tileIn, in) || subIn < after-0.5 || subIn > after+1.5 {
				t.Errorf("%v %v: the prefix is %v in and the subtitle %v after it, want %v and %v", p, dir, tileIn, subIn, in, after)
			}
			if !near(tile.Y+tile.H/2, rows[0].Y+rows[0].H/2) && p != Mac {
				t.Errorf("%v %v: the prefix %v is not centered on its row %v", p, dir, tile, rows[0])
			}
			// The line above the second and the third row, under the first row's prefix column: unlike the row's fill below it.
			img := tester.Image()
			for _, r := range rows[1:] {
				x := int(2 * (tile.X + tile.W/2))
				line, fill := img.RGBAAt(x, int(2*r.Y+1)), img.RGBAAt(x, int(2*(r.Y+4)))
				if line == fill {
					t.Errorf("%v %v: no line above the row at %v under the prefix's column (%v)", p, dir, r, line)
				}
			}
		}
	}
}

// TestPrefixRowDropsItsPrefixWhenNarrow checks that a prefix row narrower than 286 DIPs drops its prefix on Windows, as SettingsCard collapses its HeaderIcon (RightWrappedNoIcon), and keeps it elsewhere; a radio row keeps its radio button on every platform.
func TestPrefixRowDropsItsPrefixWhenNarrow(t *testing.T) {
	for _, p := range Platforms {
		v := 0
		var tile ui.Rect
		tester := layoutTester(p, 260, 400, func(c *ui.Context) {
			ui.Column(c).Padding(8).Gap(8).Children(func() {
				PrefixRow(c, "Whisper", "", func() { tile = ui.Box(c).Size(20, 20).Bounds() }, nil)
				RadioRow(c, &v, 0, "Zero", "", nil)
			})
		})
		tester.Frame()
		tester.Frame()
		if shown := tile.W > 0; shown != (p != Windows) {
			t.Errorf("%v: in a row 244 wide the prefix shows %v (%v)", p, shown, tile)
		}
	}
}

// TestPreferencesGroupHeader checks the group on every platform, left to right and right to left: the suffix's button takes clicks, at the end of the header's line; the description is below the rows on macOS and Windows and under the title, above the rows, on GNOME; the rows join as in a Fieldset; and the list starts where the platform puts it below the header: 10 below the header's line on macOS (the button's 24 points) and Windows, and on GNOME 6 below a single-line header of 34 pixels and 12 below one with a description.
func TestPreferencesGroupHeader(t *testing.T) {
	near := func(a, b float32) bool { return a-b > -0.75 && a-b < 0.75 }
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			clicks := 0
			var groups [2]ui.Rect
			var rows [3]ui.Rect
			var button ui.Rect
			tester := layoutTester(p, 600, 600, func(c *ui.Context) {
				ui.Column(c).Padding(16).Gap(24).Children(func() {
					groups[0] = PreferencesGroup(c, "Speech", "Models are stored here.", func() {
						b := Button(c, "Detect")
						if b.Clicked() {
							clicks++
						}
						button = b.Bounds()
					}, func() {
						rows[0] = ActionRow(c, "One", "", nil).Bounds()
						rows[1] = ActionRow(c, "Two", "", nil).Bounds()
					}).Bounds()
					groups[1] = PreferencesGroup(c, "Plain", "", nil, func() {
						rows[2] = ActionRow(c, "Three", "", nil).Bounds()
					}).Bounds()
				})
			})
			tester.SetDirection(dir)
			tester.Frame()
			tester.Frame()
			if err := tester.Click("Detect"); err != nil {
				t.Fatalf("%v %v: %v", p, dir, err)
			}
			tester.Frame()
			if clicks != 1 {
				t.Fatalf("%v %v: the suffix's button took %d clicks", p, dir, clicks)
			}
			desc, _ := tester.Find("Models are stored here.")
			inset := map[Platform]float32{Mac: 10}[p]
			end := groups[0].X + groups[0].W - (button.X + button.W)
			if dir == ui.RightToLeft {
				end = button.X - groups[0].X
			}
			if !near(end, inset) {
				t.Errorf("%v %v: the button ends %v from the group's end, want %v", p, dir, end, inset)
			}
			want := map[Platform]float32{Mac: 0, Windows: 4, Linux: 0}[p]
			if gap := rows[1].Y - (rows[0].Y + rows[0].H); !near(gap, want) {
				t.Errorf("%v %v: the rows are %v apart, want %v", p, dir, gap, want)
			}
			if below := desc.Y > rows[1].Y+rows[1].H; below != (p != Linux) {
				t.Errorf("%v %v: the description %v is below the rows %v: %v", p, dir, desc, rows[1], below)
			}
			switch p {
			case Linux:
				if !near(rows[0].Y-(desc.Y+desc.H), 12) {
					t.Errorf("%v %v: the list is %v below the description, want 12", p, dir, rows[0].Y-(desc.Y+desc.H))
				}
				if !near(rows[2].Y-groups[1].Y, 40) {
					t.Errorf("%v %v: a single-line group's list is %v below its top, want 40", p, dir, rows[2].Y-groups[1].Y)
				}
			default:
				if !near(rows[0].Y-(button.Y+button.H), 10) {
					t.Errorf("%v %v: the list is %v below the header's button, want 10", p, dir, rows[0].Y-(button.Y+button.H))
				}
			}
			if !layoutStateEmpty(&rowGroups) {
				t.Errorf("%v %v: a group of rows outlived the build", p, dir)
			}
		}
	}
}

// TestContentRowHoldsItsContent checks the content row on every platform: what fn builds is below the label, the platform's gap apart (9.5 on macOS, 8 on Windows, 6 on GNOME), reaches across the row inside its inset (10 on macOS, 16 on Windows, 14 on GNOME), takes clicks and the focus, and the row joins a card's list.
func TestContentRowHoldsItsContent(t *testing.T) {
	near := func(a, b float32) bool { return a-b > -0.5 && a-b < 0.5 }
	for _, p := range Platforms {
		text := ""
		var first, row, area ui.Rect
		var line float32
		tester := layoutTester(p, 600, 400, func(c *ui.Context) {
			_, line = c.MeasureText(0, ui.Span{Text: "Prompt", Size: c.Theme().FontSize})
			ui.Column(c).Padding(16).Children(func() {
				Card(c, func() {
					first = ActionRow(c, "Model", "", nil).Bounds()
					row = ContentRow(c, "Prompt", func() { area = TextArea(c, &text).Label("Prompt text").Height(60).Bounds() }).Bounds()
				})
			})
		})
		tester.Frame()
		if !near(row.Y, first.Y+first.H) {
			t.Errorf("%v: the content row %v does not touch the row above it %v", p, row, first)
		}
		gap := map[Platform]float32{Mac: 9.5, Windows: 8, Linux: 6}[p]
		inset := map[Platform]float32{Mac: 10, Windows: 16, Linux: 14}[p]
		if !near(area.X-row.X, inset) || !near(row.X+row.W-(area.X+area.W), inset) {
			t.Errorf("%v: row %v, content %v: want it %v in from both sides", p, row, area, inset)
		}
		// The label's line is above the content: at least the gap and the label's line below the row's top.
		if area.Y-row.Y < gap+line {
			t.Errorf("%v: the content %v starts %v below the row's top", p, area, area.Y-row.Y)
		}
		if err := tester.Click("Prompt text"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		tester.Type("Hi")
		tester.Frame()
		if text != "Hi" {
			t.Errorf("%v: typing in the content gave %q", p, text)
		}
	}
}

// TestMacPreferencesMatchAppKit checks prefsView on macOS at 2x, light and dark, against the values read from reference.swift's preferences scene (appkit-preferences-*.png, macOS 27.2): the section of the group with a header button 50 points from the window's top (the button's 24 points from 16, and 10 below it); the tile of the one-line prefix row 10 points in from the section's edge and its row 50 points high between the separator and the section's end; the description's ink in secondaryLabelColor, 128 over the window in Aqua and 154 in Dark Aqua, its top 13 points below the section; and an unchosen radio button's trough over the section, 223 over 247 and 58 over 37.
func TestMacPreferencesMatchAppKit(t *testing.T) {
	for _, dark := range []bool{false, true} {
		tester := ui.NewTester(prefsView(Mac), 672, 480)
		tester.SetScale(2)
		tester.SetDark(dark)
		tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		tester.Frame()
		tester.Frame()
		img := tester.Image()
		gray := func(x, y int) int { return int(img.RGBAAt(x, y).R) }
		window := gray(80, 10)
		// Down the tile's column, below the header: the section's top, the tile of the last row, the section's end.
		blue := func(y int) bool { v := img.RGBAAt(80, y); return v.B > 200 && v.R < 50 }
		top, tileTop, tileEnd, end := 0, 0, 0, 0
		for y := 85; y < 460 && end == 0; y++ {
			switch {
			case top == 0:
				if gray(80, y) != window {
					top = y
				}
			case blue(y) && !blue(y-1):
				tileTop, tileEnd = y, 0
			case !blue(y) && blue(y-1):
				tileEnd = y
			case tileEnd != 0 && gray(80, y) == window:
				end = y
			}
		}
		if top != 100 {
			t.Errorf("dark %v: the section starts at %v pixels, want 100", dark, top)
		}
		if tileEnd-tileTop != 60 || end-tileEnd != 20 {
			t.Errorf("dark %v: the last tile spans %v to %v and the section ends at %v, want 60 high and 20 above the end", dark, tileTop, tileEnd, end)
		}
		// The tile 10 points in: its first column 52 pixels from the window's edge, the section's 32.
		if v := img.RGBAAt(51, tileTop+30); v.B > 200 && v.R < 50 {
			t.Errorf("dark %v: the tile reaches the column at 51 pixels", dark)
		}
		if v := img.RGBAAt(52, tileTop+30); !(v.B > 200 && v.R < 50) {
			t.Errorf("dark %v: the tile does not start at 52 pixels (%v)", dark, v)
		}
		// The description: its ink's extreme and its top.
		ink, inkTop := window, 0
		for y := end; y < end+50; y++ {
			for x := 52; x < 80; x++ {
				v := gray(x, y)
				if (!dark && v < ink) || (dark && v > ink) {
					ink = v
				}
				if inkTop == 0 && v != window {
					inkTop = y
				}
			}
		}
		want := map[bool]int{false: 128, true: 154}[dark]
		if ink < want-2 || ink > want+2 || inkTop-end != 26 {
			t.Errorf("dark %v: the description's ink reads %v from %v pixels below the section, want %v from 26", dark, ink, inkTop-end, want)
		}
		// The unchosen radio button of the second radio row, at the trailing edge: its middle.
		row, ok := tester.Find("Cloud")
		if !ok {
			t.Fatal("no row is named Cloud")
		}
		var trough int
		for y := int(2 * row.Y); y < int(2*(row.Y+row.H)); y++ {
			if v := gray(int(2*(row.X+row.W))-28, y); v != gray(int(2*(row.X+row.W))-60, y) {
				trough = gray(int(2*(row.X+row.W))-28, y+10)
				break
			}
		}
		if want := map[bool]int{false: 223, true: 58}[dark]; trough != want {
			t.Errorf("dark %v: the unchosen radio button's trough reads %v, want %v", dark, trough, want)
		}
	}
}
