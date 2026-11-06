package breeze

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func init() {
	for _, n := range []string{"ComboRow", "EntryRow", "SpinRow", "ExpanderRow"} {
		drawnWidgets[n] = true
	}
}

// settingsView lays out the settings rows for screenshots, in two columns 300 wide as rowsView's: a card of a combo row with a subtitle, an entry row holding text, an empty one and a spin row with a subtitle, then a combo row alone and a disabled spin row alone; beside them an open expander row with a subtitle holding an action row and a switch row, and a closed one without a subtitle.
func settingsView(p Platform) func(c *ui.Context) {
	appearance, language, name, email := "Auto", "English", "Evan", ""
	size, count := 12.0, 3.0
	open, advanced, metered := true, false, false
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Row(c).Fill().Padding(16).Gap(24).AlignItems(ui.Start).Children(func() {
				ui.Column(c).Width(300).Gap(16).Children(func() {
					Card(c, func() {
						ComboRow(c, "Appearance", "Light or dark", &appearance, []string{"Light", "Dark", "Auto"})
						EntryRow(c, "Name", &name)
						EntryRow(c, "Email", &email)
						SpinRow(c, "Font size", "Points", &size, 6, 72, 1)
					})
					ComboRow(c, "Language", "", &language, []string{"English", "Français"})
					SpinRow(c, "Copies", "", &count, 1, 9, 1).Disabled(true)
				})
				ui.Column(c).Width(300).Gap(16).Children(func() {
					ExpanderRow(c, "Network", "Wi-Fi and proxies", &open, func() {
						ActionRow(c, "Proxy", "Automatic", nil)
						SwitchRow(c, "Metered", "", &metered)
					})
					ExpanderRow(c, "Advanced", "", &advanced, func() {
						ActionRow(c, "Logs", "", nil)
					})
				})
			})
		})
	}
}

// TestScreenshotsSettings writes the settings rows of every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names. Without it the test only builds the frames.
func TestScreenshotsSettings(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			h := 480
			if p == Windows {
				h = 640
			}
			tester := ui.NewTester(settingsView(p), 672, h)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tester.Frame()
			tester.Frame()
			for _, s := range []string{"Appearance", "Light or dark", "Auto", "Font size", "Points", "Language", "Network", "Wi-Fi and proxies", "Proxy", "Metered", "Advanced"} {
				if !tester.HasText(s) {
					t.Fatalf("%v: %q is missing", p, s)
				}
			}
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			writeShot(t, dir, fmt.Sprintf("breeze-settings-%s-%s.png", strings.ToLower(p.String()), mode), tester)
		}
	}
}

// settingsBefore builds a button before the row fn builds and one after it, in a column, for the tests that walk the focus with Tab.
func settingsAround(c *ui.Context, fn func()) {
	ui.Column(c).Padding(16).Gap(8).Children(func() {
		Button(c, "Before")
		fn()
		Button(c, "After")
	})
}

// TestComboRowChooses checks the combo row on every platform, left to right and right to left: Tab from the button before it stops once, at the control named by the title (the row on GNOME, the pop-up button or the ComboBox elsewhere), Down opens the list and Enter chooses, a click on the choice opens it and a click on an option chooses it, Changed reports each choice, the choice sits after the title at the row's end, and on GNOME a click on the title opens the list too, as AdwComboRow activates on the whole row.
func TestComboRowChooses(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			choice, changes := "Light", 0
			tester := layoutTester(p, 560, 400, func(c *ui.Context) {
				settingsAround(c, func() {
					if ComboRow(c, "Appearance", "Light or dark", &choice, []string{"Light", "Dark", "Auto"}).Changed() {
						changes++
					}
				})
			})
			tester.SetDirection(dir)
			tester.Frame()
			title, _ := tester.Find("Light or dark")
			value, ok := tester.Find("Light")
			if !ok {
				t.Fatalf("%v %v: the choice is missing", p, dir)
			}
			after := value.X > title.X+title.W
			if dir == ui.RightToLeft {
				after = value.X+value.W < title.X
			}
			if !after {
				t.Errorf("%v %v: subtitle %v, choice %v: the choice is not after the title", p, dir, title, value)
			}
			if err := tester.Click("Before"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			if !tester.Focused("Appearance") {
				t.Fatalf("%v %v: Tab did not reach the control named Appearance", p, dir)
			}
			tester.Key(0, ui.KeyDown)
			tester.Frame()
			tester.Key(0, ui.KeyDown)
			tester.Frame()
			tester.Key(0, ui.KeyEnter)
			tester.Frame()
			if choice != "Dark" || changes != 1 {
				t.Fatalf("%v %v: Down, Down and Enter chose %q after %d changes", p, dir, choice, changes)
			}
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			if !tester.Focused("After") {
				t.Fatalf("%v %v: a second Tab did not go past the row to the button", p, dir)
			}
			if err := tester.Click("Dark"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			if err := tester.Click("Auto"); err != nil {
				t.Fatalf("%v %v: a click on the choice opened no list: %v", p, dir, err)
			}
			tester.Frame()
			if choice != "Auto" || changes != 2 {
				t.Fatalf("%v %v: a click on Auto chose %q after %d changes", p, dir, choice, changes)
			}
			if p == Linux {
				if err := tester.Click("Light or dark"); err != nil {
					t.Fatal(err)
				}
				tester.Frame()
				if !tester.HasText("Light") {
					t.Errorf("%v: a click on the row's subtitle opened no list", dir)
				}
			}
		}
	}
}

// TestComboRowDisabled checks on every platform that a disabled combo row opens no list and takes no focus.
func TestComboRowDisabled(t *testing.T) {
	for _, p := range Platforms {
		choice := "Light"
		tester := layoutTester(p, 560, 300, func(c *ui.Context) {
			ui.Column(c).Padding(16).Children(func() {
				ComboRow(c, "Appearance", "", &choice, []string{"Light", "Dark"}).Disabled(true)
			})
		})
		if err := tester.Click("Light"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		if tester.HasText("Dark") || tester.Focused("Appearance") {
			t.Errorf("%v: a disabled row opened its list %v or took the focus %v", p, tester.HasText("Dark"), tester.Focused("Appearance"))
		}
	}
}

// TestEntryRowEdits checks the entry row on every platform, left to right and right to left: Tab from the button before it focuses the field named by the title, typing edits the value and Changed reports it, a second Tab leaves the row, and on GNOME a click on the row's end, the pencil, focuses the text, as AdwEntryRow's does, and the pencil hides while the text has the focus.
func TestEntryRowEdits(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			name, changes := "", 0
			var row ui.Rect
			tester := layoutTester(p, 560, 300, func(c *ui.Context) {
				settingsAround(c, func() {
					r := EntryRow(c, "Name", &name)
					if r.Changed() {
						changes++
					}
					row = r.Bounds()
				})
			})
			tester.SetDirection(dir)
			tester.Frame()
			if err := tester.Click("Before"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			if !tester.Focused("Name") {
				t.Fatalf("%v %v: Tab did not focus the field named Name", p, dir)
			}
			tester.Type("Ada")
			tester.Frame()
			if name != "Ada" || changes == 0 {
				t.Fatalf("%v %v: typing left %q after %d changes", p, dir, name, changes)
			}
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			if !tester.Focused("After") {
				t.Fatalf("%v %v: a second Tab did not go past the row", p, dir)
			}
			if p != Linux {
				continue
			}
			// The pencil's middle: its 34-pixel box ends 14 in from the row's end.
			x := row.X + row.W - 31
			if dir == ui.RightToLeft {
				x = row.X + 31
			}
			pencil := func() int {
				img, ink := tester.Image(), 0
				for y := int(row.Y + row.H/2 - 8); y < int(row.Y+row.H/2+8); y++ {
					for px := int(x - 8); px < int(x+8); px++ {
						if img.RGBAAt(px, y).R < 128 {
							ink++
						}
					}
				}
				return ink
			}
			if pencil() == 0 {
				t.Errorf("%v: no pencil at the row's end", dir)
			}
			tester.ClickAt(x, row.Y+row.H/2)
			tester.Frame()
			if !tester.Focused("Name") {
				t.Fatalf("%v: a click on the row's end did not focus the text", dir)
			}
			if pencil() != 0 {
				t.Errorf("%v: the pencil still shows while the text has the focus", dir)
			}
		}
	}
}

// TestEntryRowTitleShrinks checks on GNOME that the row's title is drawn at the row's middle while the row is empty and unfocused, and above the text, smaller, once it holds text, as AdwEntryRow's: the ink of the title moves up.
func TestEntryRowTitleShrinks(t *testing.T) {
	for _, filled := range []bool{false, true} {
		value := ""
		if filled {
			value = "x"
		}
		var row ui.Rect
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Linux, func() {
				Use(c)
				ui.Column(c).Padding(16).Children(func() { row = EntryRow(c, "Name", &value).Bounds() })
			})
		}, 400, 120)
		tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		tester.Frame()
		tester.Frame()
		img := tester.Image()
		// The topmost ink in the first 40 pixels of the text's column, left of where the text "x" would be drawn below the title.
		top := -1
		for y := int(row.Y); y < int(row.Y+row.H) && top < 0; y++ {
			for x := int(row.X + 14); x < int(row.X+40); x++ {
				if img.RGBAAt(x, y).R < 200 {
					top = y
					break
				}
			}
		}
		mid := int(row.Y + row.H/2)
		if filled && (top < 0 || top > mid-8) || !filled && (top < mid-12 || top > mid) {
			t.Errorf("filled %v: the title's ink starts at %d in a row from %v to %v", filled, top, row.Y, row.Y+row.H)
		}
	}
}

// TestSpinRowSteps checks the spin row on every platform: Tab stops once, at the text named by the title, Up steps the value, a click on Increase steps it and leaves the focus in the text, the buttons are no stops of their own, Increase is disabled at the upper limit, and on GNOME a click on the title focuses the text, as AdwSpinRow's grab_focus does, and so does a click on a button while the focus is elsewhere, as GtkSpinButton's button_pressed_cb grabs it (gtkspinbutton.c of GTK 4.24.1).
func TestSpinRowSteps(t *testing.T) {
	for _, p := range Platforms {
		size, changes := 10.0, 0
		tester := layoutTester(p, 560, 300, func(c *ui.Context) {
			settingsAround(c, func() {
				if SpinRow(c, "Font size", "Points", &size, 6, 12, 1).Changed() {
					changes++
				}
			})
		})
		if err := tester.Click("Before"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		if !tester.Focused("Font size") {
			t.Fatalf("%v: Tab did not focus the text named Font size", p)
		}
		tester.Key(0, ui.KeyUp)
		tester.Frame()
		if size != 11 || changes != 1 {
			t.Fatalf("%v: Up left %v after %d changes", p, size, changes)
		}
		if err := tester.Click("Increase"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		if size != 12 || changes != 2 {
			t.Fatalf("%v: Increase left %v after %d changes", p, size, changes)
		}
		if !tester.Focused("Font size") {
			t.Errorf("%v: a click on Increase took the focus from the text", p)
		}
		if err := tester.Click("Increase"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		if size != 12 {
			t.Errorf("%v: Increase went past the limit to %v", p, size)
		}
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		if !tester.Focused("After") {
			t.Fatalf("%v: Tab from the text did not go past the buttons to the button after the row", p)
		}
		if p == Linux {
			if err := tester.Click("Points"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			if !tester.Focused("Font size") {
				t.Error("a click on the row did not focus the text")
			}
			if err := tester.Click("After"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			if err := tester.Click("Decrease"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			if size != 11 || !tester.Focused("Font size") {
				t.Errorf("a click on Decrease with the focus on a button left %v, the text focused %v", size, tester.Focused("Font size"))
			}
		}
	}
}

// TestExpanderRowOpens checks the expander row on every platform, left to right and right to left: closed it hides its rows, a click on its title opens it and Changed reports it, its rows reach its sides and touch one another, Tab stops once at the header, named by the title, Space closes it and Enter opens it again, and the header is as high as the platform's at least: 68 DIPs on Windows (SettingsCardMinHeight), 54 pixels on GNOME (the row's 50 and its padding of 2).
func TestExpanderRowOpens(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			open, changes, on := false, 0, false
			var root, a, b ui.Rect
			tester := layoutTester(p, 560, 500, func(c *ui.Context) {
				settingsAround(c, func() {
					e := ExpanderRow(c, "Network", "Wi-Fi and proxies", &open, func() {
						a = ActionRow(c, "Proxy", "Automatic", nil).Bounds()
						b = SwitchRow(c, "Metered", "", &on).Bounds()
					})
					if e.Changed() {
						changes++
					}
					root = e.Bounds()
				})
			})
			tester.SetDirection(dir)
			tester.Frame()
			if tester.HasText("Proxy") {
				t.Fatalf("%v %v: a closed row shows its rows", p, dir)
			}
			header, ok := tester.Find("Network")
			if !ok {
				t.Fatalf("%v %v: no element is named Network", p, dir)
			}
			if want := map[Platform]float32{Mac: 40, Windows: 68, Linux: 54}[p]; header.H < want-0.5 {
				t.Errorf("%v %v: the header is %v high, want at least %v", p, dir, header.H, want)
			}
			if err := tester.Click("Wi-Fi and proxies"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			tester.Frame()
			if !open || changes != 1 || !tester.HasText("Proxy") {
				t.Fatalf("%v %v: a click left it open %v after %d changes", p, dir, open, changes)
			}
			inset := float32(0)
			if p == Windows {
				inset = 1 // the content's border
			}
			near := func(x, y float32) bool { return x-y > -0.5 && x-y < 0.5 }
			if !near(a.X, root.X+inset) || !near(a.X+a.W, root.X+root.W-inset) || !near(b.Y, a.Y+a.H) {
				t.Errorf("%v %v: row %v, rows %v and %v: the rows do not reach the sides or do not touch", p, dir, root, a, b)
			}
			if err := tester.Click("Before"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			if !tester.Focused("Network") {
				t.Fatalf("%v %v: Tab did not focus the header", p, dir)
			}
			tester.Key(0, ui.KeySpace)
			tester.Frame()
			if open || changes != 2 {
				t.Fatalf("%v %v: Space left it open %v after %d changes", p, dir, open, changes)
			}
			tester.Key(0, ui.KeyEnter)
			tester.Frame()
			if !open || changes != 3 {
				t.Fatalf("%v %v: Enter left it open %v after %d changes", p, dir, open, changes)
			}
			if !layoutStateEmpty(&rowGroups) {
				t.Errorf("%v %v: a group of rows outlived the build", p, dir)
			}
		}
	}
}

// TestSettingsRowsJoinInACard checks on every platform that the settings rows join a Card's list as ActionRow does: on macOS and GNOME they touch, and so they do on Windows, as a SettingsExpander's items; alone each is a group of its own.
func TestSettingsRowsJoinInACard(t *testing.T) {
	for _, p := range Platforms {
		choice, name, size, open := "A", "", 1.0, false
		var rows [4]ui.Rect
		tester := layoutTester(p, 600, 600, func(c *ui.Context) {
			ui.Column(c).Padding(16).Children(func() {
				Card(c, func() {
					rows[0] = ComboRow(c, "Combo", "", &choice, []string{"A", "B"}).Bounds()
					rows[1] = EntryRow(c, "Entry", &name).Bounds()
					rows[2] = SpinRow(c, "Spin", "", &size, 0, 9, 1).Bounds()
					rows[3] = ExpanderRow(c, "Expander", "", &open, func() {}).Bounds()
				})
			})
		})
		tester.Frame()
		for i := 1; i < len(rows); i++ {
			if gap := rows[i].Y - (rows[i-1].Y + rows[i-1].H); gap < -0.5 || gap > 0.5 {
				t.Errorf("%v: rows %d and %d are %v apart", p, i-1, i, gap)
			}
		}
	}
}

// TestSettingsRowsDisabled checks on every platform that the four rows, disabled, are no Tab stops and that clicks on them open, focus, step and type nothing.
func TestSettingsRowsDisabled(t *testing.T) {
	for _, p := range Platforms {
		choice, name, size, open := "Light", "", 3.0, false
		tester := layoutTester(p, 560, 500, func(c *ui.Context) {
			settingsAround(c, func() {
				ComboRow(c, "Appearance", "", &choice, []string{"Light", "Dark"}).Disabled(true)
				EntryRow(c, "Name", &name).Disabled(true)
				SpinRow(c, "Copies", "", &size, 1, 9, 1).Disabled(true)
				ExpanderRow(c, "Network", "", &open, func() { ActionRow(c, "Proxy", "", nil) }).Disabled(true)
			})
		})
		if err := tester.Click("Before"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		if !tester.Focused("After") {
			t.Errorf("%v: Tab from the button before the rows stopped in a disabled row", p)
		}
		for _, s := range []string{"Light", "Name", "Copies", "Increase", "Network"} {
			if err := tester.Click(s); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
		}
		tester.Type("abc")
		tester.Frame()
		if tester.HasText("Dark") || name != "" || size != 3 || open || tester.Focused("Name") || tester.Focused("Copies") {
			t.Errorf("%v: clicks on disabled rows opened the list %v, typed %q, stepped to %v, opened %v, focused the entry %v or the spin row %v", p, tester.HasText("Dark"), name, size, open, tester.Focused("Name"), tester.Focused("Copies"))
		}
	}
}

// TestMacFormRowHeights checks on macOS the heights SwiftUI gives a grouped form's rows, measured on macOS 27.2, against a row of text so that they hold whatever the host's font: alone, a picker's row of one line is at least 37.5 points high and a stepper's 40.5, taller for their controls than a row of text; an entry row is half a point shorter than a row of text; a closed expander row is its title's line and 14.5 points above and below it.
func TestMacFormRowHeights(t *testing.T) {
	choice, name, size, open := "A", "", 1.0, false
	var text, combo, entry, spin, expander ui.Rect
	tester := layoutTester(Mac, 400, 500, func(c *ui.Context) {
		ui.Column(c).Width(300).Gap(16).Children(func() {
			text = ActionRow(c, "Text", "", nil).Bounds()
			combo = ComboRow(c, "Combo", "", &choice, []string{"A", "B"}).Bounds()
			entry = EntryRow(c, "Entry", &name).Bounds()
			spin = SpinRow(c, "Spin", "", &size, 0, 9, 1).Bounds()
			expander = ExpanderRow(c, "Expander", "", &open, func() {}).Bounds()
		})
	})
	tester.Frame()
	near := func(a, b float32) bool { return a-b > -0.01 && a-b < 0.01 }
	// A row of text is 10.5 above its title's line and 11 below it.
	line := text.H - 21.5
	if combo.H < 37.5-0.01 || spin.H < 40.5-0.01 {
		t.Errorf("a picker's row is %v high and a stepper's %v, want 37.5 and 40.5 at least", combo.H, spin.H)
	}
	if !near(entry.H, text.H-0.5) {
		t.Errorf("an entry row is %v high where a row of text is %v, want half a point less", entry.H, text.H)
	}
	if !near(expander.H, line+29) {
		t.Errorf("a closed expander row is %v high, want its title's %v and 29", expander.H, line)
	}
}

// TestMacStepperFieldInASection checks on macOS a spin row's field at 2x, light and dark: its bezel's line is laid over the section as AppKit lays it (210 over 247, 57 over 37, measured on macOS 27.2), and the bezel ends 9 points from the row's end, a point past the stepper's.
func TestMacStepperFieldInASection(t *testing.T) {
	for _, dark := range []bool{false, true} {
		size := 3.0
		var row ui.Rect
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Mac, func() {
				Use(c)
				ui.Column(c).Padding(16).Width(332).Children(func() { row = SpinRow(c, "Copies", "", &size, 1, 9, 1).Bounds() })
			})
		}, 400, 120)
		tester.SetScale(2)
		tester.SetDark(dark)
		tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		tester.Frame()
		tester.Frame()
		inc, ok := tester.Find("Increase")
		if !ok {
			t.Fatal("no stepper half is named Increase")
		}
		if end := row.X + row.W - 10; inc.X+20 < end-0.01 || inc.X+20 > end+0.01 {
			t.Errorf("dark %v: the stepper ends at %v, want %v, a point before the bezel's end", dark, inc.X+20, end)
		}
		want, section := map[bool]uint8{false: 210, true: 57}[dark], map[bool]uint8{false: 247, true: 37}[dark]
		img := tester.Image()
		near := func(got, want uint8) bool { return got+2 >= want && got <= want+2 }
		// The bezel's top line over the field's middle, and its end line across the stepper's middle: its last pixel and the section's first after it.
		top := img.RGBAAt(int((inc.X-30)*2), int((inc.Y-1)*2)).R
		y := int((inc.Y + 12) * 2)
		last, after := img.RGBAAt(int((row.X+row.W-9)*2)-1, y).R, img.RGBAAt(int((row.X+row.W-9)*2), y).R
		if !near(top, want) || !near(last, want) || !near(after, section) {
			t.Errorf("dark %v: the bezel's line reads %d at the top and %d at the end, %d after it; want %d and the section's %d after it", dark, top, last, after, want, section)
		}
	}
}

// TestFluentExpanderRowHeader checks on Windows, light and dark, the SettingsExpander's header as WinUI renders it (tools/winui-reference's winui-rows): 70 DIPs high closed, 68 of fill inside its border; open, the first item's line right below 68 DIPs of fill, a single line, the header on ControlFillColorDefault (251, 45) and the items on CardBackgroundFillColorDefault (251, 43).
func TestFluentExpanderRowHeader(t *testing.T) {
	for _, dark := range []bool{false, true} {
		open := false
		var item ui.Rect
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Windows, func() {
				Use(c)
				ui.Column(c).Padding(16).Width(332).Children(func() {
					ExpanderRow(c, "Network", "", &open, func() { item = ActionRow(c, "Proxy", "", nil).Bounds() })
				})
			})
		}, 400, 300)
		tester.SetDark(dark)
		tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		tester.Frame()
		header, ok := tester.Find("Network")
		if !ok {
			t.Fatal("no element is named Network")
		}
		if header.H < 70-0.01 || header.H > 70+0.01 {
			t.Errorf("dark %v: the closed header is %v high, want 70", dark, header.H)
		}
		open = true
		tester.Frame()
		tester.Frame()
		header, _ = tester.Find("Network")
		if top := header.Y + 69; item.Y < top-0.01 || item.Y > top+0.01 {
			t.Errorf("dark %v: the first item starts at %v, want right below the header's border and 68 DIPs of fill, %v", dark, item.Y, top)
		}
		img := tester.Image()
		x := int(item.X + 8)
		face, line, above, fill := img.RGBAAt(x, int(header.Y+header.H/2)).R, img.RGBAAt(x, int(item.Y)).R, img.RGBAAt(x, int(item.Y)-1).R, img.RGBAAt(x, int(item.Y+item.H/2)).R
		want := map[bool][3]uint8{false: {251, 229, 251}, true: {45, 29, 43}}[dark]
		if face != want[0] || line != want[1] || above != face || fill != want[2] {
			t.Errorf("dark %v: header %d, line %d with %d above it, item %d; want %d, a single line of %d, %d", dark, face, line, above, fill, want[0], want[1], want[2])
		}
	}
}

// TestAdwExpanderRowArrow checks on GNOME, light and dark, that an open expander row's arrow is accent_color as libadwaita 1.10.0 draws it, #0461be and #81d0ff for the default blue (measured by tools/adwaita-reference).
func TestAdwExpanderRowArrow(t *testing.T) {
	for _, dark := range []bool{false, true} {
		open := true
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Linux, func() {
				Use(c)
				ui.Column(c).Padding(16).Width(332).Children(func() { ExpanderRow(c, "Network", "", &open, func() {}) })
			})
		}, 400, 200)
		tester.SetScale(2)
		tester.SetDark(dark)
		tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		tester.Frame()
		tester.Frame()
		header, ok := tester.Find("Network")
		if !ok {
			t.Fatal("no element is named Network")
		}
		want := map[bool][3]int{false: {4, 97, 190}, true: {129, 208, 255}}[dark]
		abs := func(v int) int { return max(v, -v) }
		img, found := tester.Image(), false
		// The 16-pixel arrow ends 14 pixels before the row's end: 2 of the row's padding and 12 of the header's margin.
		end := header.X + header.W - 14
		for y := int(header.Y * 2); y < int((header.Y+header.H)*2) && !found; y++ {
			for x := int((end - 16) * 2); x < int(end*2); x++ {
				c := img.RGBAAt(x, y)
				if d := abs(int(c.R)-want[0]) + abs(int(c.G)-want[1]) + abs(int(c.B)-want[2]); d <= 6 {
					found = true
					break
				}
			}
		}
		if !found {
			t.Errorf("dark %v: no pixel of the open arrow is accent_color %v", dark, want)
		}
	}
}
