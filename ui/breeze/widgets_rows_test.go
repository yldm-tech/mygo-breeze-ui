package breeze

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func init() {
	for _, n := range []string{"Card", "ActionRow", "SwitchRow"} {
		drawnWidgets[n] = true
	}
}

// TestSwitchRowToggles checks the switch row on every platform, left to right and right to left: a click on its subtitle, in its middle and at its far end toggles it and Changed reports each, Shift-Tab from the button after it focuses the row as one stop named by its title, Space and Enter toggle it, and the switch inside it is no stop of its own.
func TestSwitchRowToggles(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			on, changes, after := false, 0, 0
			tester := layoutTester(p, 420, 300, func(c *ui.Context) {
				ui.Column(c).Padding(16).Gap(8).Children(func() {
					if SwitchRow(c, "Dark mode", "Use a dark appearance", &on).Changed() {
						changes++
					}
					if Button(c, "After").Clicked() {
						after++
					}
				})
			})
			tester.SetDirection(dir)
			tester.Frame()
			// The subtitle, then the row's middle (Find returns the row, named by its title).
			for i, s := range []string{"Use a dark appearance", "Dark mode"} {
				if err := tester.Click(s); err != nil {
					t.Fatalf("%v %v: %v", p, dir, err)
				}
				tester.Frame()
				if on != (i == 0) || changes != i+1 {
					t.Fatalf("%v %v: a click on %q left it on %v after %d changes", p, dir, s, on, changes)
				}
			}
			row, ok := tester.Find("Dark mode")
			if !ok {
				t.Fatalf("%v %v: no element is named Dark mode", p, dir)
			}
			// The switch's end of the row: the right left to right, the left right to left.
			x := row.X + row.W - 12
			if dir == ui.RightToLeft {
				x = row.X + 12
			}
			tester.ClickAt(x, row.Y+row.H/2)
			tester.Frame()
			if !on || changes != 3 {
				t.Fatalf("%v %v: a click on the switch's end left it on %v after %d changes", p, dir, on, changes)
			}
			// From the button after it, Shift-Tab comes back to the row as one stop.
			if err := tester.Click("After"); err != nil {
				t.Fatalf("%v %v: %v", p, dir, err)
			}
			tester.Frame()
			tester.Key(ui.Shift, ui.KeyTab)
			tester.Frame()
			if !tester.Focused("Dark mode") {
				t.Fatalf("%v %v: Shift-Tab from the button did not focus the row", p, dir)
			}
			tester.Key(0, ui.KeySpace)
			tester.Frame()
			if on || changes != 4 {
				t.Fatalf("%v %v: Space left it on %v after %d changes", p, dir, on, changes)
			}
			tester.Key(0, ui.KeyEnter)
			tester.Frame()
			if !on || changes != 5 {
				t.Fatalf("%v %v: Enter left it on %v after %d changes", p, dir, on, changes)
			}
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			if !tester.Focused("After") {
				t.Fatalf("%v %v: a second Tab did not go past the row to the button: the switch inside is a stop of its own", p, dir)
			}
		}
	}
}

// TestSwitchRowDisabled checks on every platform that a disabled switch row takes neither a click nor the focus.
func TestSwitchRowDisabled(t *testing.T) {
	for _, p := range Platforms {
		on := false
		tester := layoutTester(p, 420, 200, func(c *ui.Context) {
			ui.Column(c).Padding(16).Children(func() {
				SwitchRow(c, "Bluetooth", "", &on).Disabled(true)
			})
		})
		if err := tester.Click("Bluetooth"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		if on || tester.Focused("Bluetooth") {
			t.Fatalf("%v: a disabled row turned on %v, focused %v", p, on, tester.Focused("Bluetooth"))
		}
	}
}

// TestActionRowHoldsItsControls checks the action row on every platform: its title and subtitle show, a button at its end takes clicks and the focus as it would anywhere, the button sits after the title, at the row's end, and the texts the platform's inset from the row's start, mirrored right to left.
func TestActionRowHoldsItsControls(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			clicks := 0
			tester := layoutTester(p, 420, 200, func(c *ui.Context) {
				ui.Column(c).Padding(16).Children(func() {
					ActionRow(c, "Wallpaper", "Choose a picture", func() {
						if Button(c, "Open").Clicked() {
							clicks++
						}
					})
				})
			})
			tester.SetDirection(dir)
			tester.Frame()
			if !tester.HasText("Choose a picture") {
				t.Fatalf("%v %v: the subtitle is missing", p, dir)
			}
			if err := tester.Click("Open"); err != nil {
				t.Fatalf("%v %v: %v", p, dir, err)
			}
			tester.Frame()
			if clicks != 1 {
				t.Fatalf("%v %v: the button took %d clicks", p, dir, clicks)
			}
			if !tester.Focused("Open") {
				t.Fatalf("%v %v: Tab did not go to the row's button", p, dir)
			}
			row, _ := tester.Find("Wallpaper")
			open, _ := tester.Find("Open")
			sub, _ := tester.Find("Choose a picture")
			after := open.X > sub.X+sub.W && row.X+row.W-(open.X+open.W) < 40
			if dir == ui.RightToLeft {
				after = open.X+open.W < sub.X && open.X-row.X < 40
			}
			if !after {
				t.Errorf("%v %v: row %v, subtitle %v, button %v: the button is not at the row's end", p, dir, row, sub, open)
			}
			// The texts are inset from the row's start as the platform insets them: 10 points on macOS, 16 DIPs on Windows, 14 pixels on GNOME.
			inset := sub.X - row.X
			if dir == ui.RightToLeft {
				inset = row.X + row.W - (sub.X + sub.W)
			}
			if want := map[Platform]float32{Mac: 10, Windows: 17, Linux: 14}[p]; inset < want-0.5 || inset > want+1.5 {
				t.Errorf("%v %v: the subtitle is %v in from the row's start, want %v", p, dir, inset, want)
			}
		}
	}
}

// TestRowsJoinInGroups checks how rows one after another join on every platform: in a Card and a Fieldset on macOS and GNOME they touch, one list; on Windows they touch in a Card, as a SettingsExpander's items, and stand 4 DIPs apart in a Fieldset, as SettingsCards; alone, two rows are apart as the column around them sets them. It also checks that no group outlives the build.
func TestRowsJoinInGroups(t *testing.T) {
	for _, p := range Platforms {
		a, b := false, true
		var rows [3][2]ui.Rect
		tester := layoutTester(p, 420, 700, func(c *ui.Context) {
			ui.Column(c).Padding(16).Gap(20).Children(func() {
				Card(c, func() {
					rows[0][0] = ActionRow(c, "One", "", nil).Bounds()
					rows[0][1] = SwitchRow(c, "Two", "", &a).Bounds()
				})
				Fieldset(c, "Group", func() {
					rows[1][0] = ActionRow(c, "Three", "", nil).Bounds()
					rows[1][1] = SwitchRow(c, "Four", "", &b).Bounds()
				})
				ui.Column(c).Gap(10).Children(func() {
					rows[2][0] = ActionRow(c, "Five", "", nil).Bounds()
					rows[2][1] = SwitchRow(c, "Six", "", &a).Bounds()
				})
			})
		})
		tester.Frame()
		want := [3]float32{0, 0, 10}
		if p == Windows {
			want[1] = 4
		}
		for i, pair := range rows {
			gap := pair[1].Y - (pair[0].Y + pair[0].H)
			if d := gap - want[i]; d < -0.5 || d > 0.5 {
				t.Errorf("%v: rows %d are %v apart, want %v", p, i, gap, want[i])
			}
		}
		if !layoutStateEmpty(&rowGroups) {
			t.Errorf("%v: a group of rows outlived the build", p)
		}
	}
}

// TestRowsFillAnAccordionItem checks on every platform that rows in an accordion's open section reach its edges and touch one another and the section's end, as the nested rows of AdwExpanderRow and the items of a SettingsExpander do, or on macOS stop 5 points above it, as a DisclosureGroup's rows in a grouped form.
func TestRowsFillAnAccordionItem(t *testing.T) {
	for _, p := range Platforms {
		open, on := true, false
		var ib, fb, sb ui.Rect
		tester := layoutTester(p, 420, 500, func(c *ui.Context) {
			ui.Column(c).Padding(16).Children(func() {
				Accordion(c, func() {
					ib = AccordionItem(c, "Network", &open, func() {
						fb = ActionRow(c, "Proxy", "Automatic", nil).Bounds()
						sb = SwitchRow(c, "Metered", "", &on).Bounds()
					}).Bounds()
				})
			})
		})
		tester.Frame()
		near := func(a, b float32) bool { return a-b > -0.5 && a-b < 0.5 }
		inset := float32(0)
		if p == Windows {
			inset = 1 // the item's border
		}
		if !near(fb.X, ib.X+inset) || !near(fb.X+fb.W, ib.X+ib.W-inset) {
			t.Errorf("%v: item %v, row %v: the row does not reach the item's sides", p, ib, fb)
		}
		if !near(sb.Y, fb.Y+fb.H) {
			t.Errorf("%v: rows %v and %v do not touch", p, fb, sb)
		}
		// A DisclosureGroup of a grouped form leaves 5 points below its last row.
		below := inset
		if p == Mac {
			below = 5
		}
		if !near(sb.Y+sb.H, ib.Y+ib.H-below) {
			t.Errorf("%v: item %v, last row %v: the row does not reach the item's end", p, ib, sb)
		}
	}
}

// rowsView lays out the widgets of the rows group for screenshots: a card of text, a card of rows (an action row with a button, a switch row on with a subtitle, one off without), a fieldset of rows, a row alone and a disabled switch row, and an accordion section of rows.
func rowsView(p Platform) func(c *ui.Context) {
	dark, notify, location, metered, bt := true, false, true, false, true
	open := true
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Row(c).Fill().Padding(16).Gap(24).AlignItems(ui.Start).Children(func() {
				ui.Column(c).Width(300).Gap(16).Children(func() {
					Card(c, func() {
						ui.Text(c, "A card holds what belongs together.")
					})
					Card(c, func() {
						ActionRow(c, "Wallpaper", "Choose a picture", func() { Button(c, "Open") })
						SwitchRow(c, "Dark mode", "Use a dark appearance", &dark)
						SwitchRow(c, "Notifications", "", &notify)
					})
					ActionRow(c, "Storage", "", func() { ui.Text(c, "120 GB").TextColor(c.Theme().TextMuted) })
					SwitchRow(c, "Bluetooth", "", &bt).Disabled(true)
				})
				ui.Column(c).Width(300).Gap(16).Children(func() {
					Fieldset(c, "Privacy", func() {
						SwitchRow(c, "Location", "Apps may ask for it", &location)
						ActionRow(c, "Camera", "2 apps", nil)
					})
					Accordion(c, func() {
						AccordionItem(c, "Network", &open, func() {
							ActionRow(c, "Proxy", "Automatic", nil)
							SwitchRow(c, "Metered", "", &metered)
						})
					})
				})
			})
		})
	}
}

// TestScreenshotsRows writes the rows group of every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names. Without it the test only builds the frames.
func TestScreenshotsRows(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			tester := ui.NewTester(rowsView(p), 672, 520)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tester.Frame()
			tester.Frame()
			for _, s := range []string{"A card holds what belongs together.", "Choose a picture", "Notifications", "Privacy", "Apps may ask for it", "Metered", "Bluetooth"} {
				if !tester.HasText(s) {
					t.Fatalf("%v: %q is missing", p, s)
				}
			}
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			writeShot(t, dir, fmt.Sprintf("breeze-rows-%s-%s.png", strings.ToLower(p.String()), mode), tester)
		}
	}
}

// TestSwitchRowShowsPointerAndFocus checks the switch row's states on every platform: under the pointer its background changes on Windows and GNOME, as a clickable SettingsCard and an activatable row of a boxed list do, and stays on macOS, whose form rows do not light up; with the focus from the keyboard it shows a ring on every platform.
func TestSwitchRowShowsPointerAndFocus(t *testing.T) {
	for _, p := range Platforms {
		on := false
		tester := layoutTester(p, 420, 200, func(c *ui.Context) {
			Use(c)
			ui.Column(c).Padding(24).Children(func() {
				SwitchRow(c, "Dark mode", "", &on)
			})
		})
		tester.Frame()
		row, ok := tester.Find("Dark mode")
		if !ok {
			t.Fatalf("%v: no row", p)
		}
		before := tester.Image()
		// Inside the row before its text, which every platform insets by 10 or more.
		x, y := int(row.X+4), int(row.Y+row.H/2)
		rest := before.RGBAAt(x, y)
		tester.Move(row.X+row.W/2, row.Y+row.H/2)
		tester.Frame()
		hover := tester.Image().RGBAAt(x, y)
		if lit := hover != rest; lit != (p != Mac) {
			t.Errorf("%v: under the pointer the row reads %v, at rest %v", p, hover, rest)
		}
		tester.Move(0, 0)
		tester.Frame()
		rest2 := tester.Image()
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		if !tester.Focused("Dark mode") {
			t.Fatalf("%v: Tab did not focus the row", p)
		}
		focused := tester.Image()
		changed := 0
		for i := range focused.Pix {
			if focused.Pix[i] != rest2.Pix[i] {
				changed++
			}
		}
		if changed == 0 {
			t.Errorf("%v: the focus from the keyboard shows no ring", p)
		}
	}
}
