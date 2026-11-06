package breeze

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func init() {
	for _, n := range []string{"Sheet", "DestructiveButton", "DestructiveAlertDialog"} {
		drawnWidgets[n] = true
	}
}

// dialogsAt returns the color of the last frame at (x, y), in DIPs at a scale of 1.
func dialogsAt(tester *ui.Tester, x, y float32) ui.Color {
	px := tester.Image().RGBAAt(int(x), int(y))
	return ui.RGB(px.R, px.G, px.B)
}

// dialogsNear reports whether a and b differ by at most tol in every channel.
func dialogsNear(a, b ui.Color, tol int) bool {
	d := func(x, y uint8) int { return max(int(x)-int(y), int(y)-int(x)) }
	return d(a.R, b.R) <= tol && d(a.G, b.G) <= tol && d(a.B, b.B) <= tol
}

// dialogsSheetView is a window holding a button and a sheet with a text field, named "Rename" for the tests to find it, open while *open; behind counts the button's clicks, and dismissed notes a frame whose sheet reports Dismissed.
func dialogsSheetView(p Platform, open *bool, name *string, behind *int, dismissed *bool) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Column(c).Fill().Padding(16).AlignItems(ui.Start).Children(func() {
				if Button(c, "Behind").Clicked() {
					*behind++
				}
			})
			s := Sheet(c, open, func() {
				Text(c, "Rename", TextHeading)
				TextInput(c, name).Label("Name").Width(220)
				ui.Row(c).Gap(8).Justify(ui.End).Children(func() {
					if Button(c, "Cancel").Clicked() {
						*open = false
					}
					if PrimaryButton(c, "Rename").Clicked() {
						*open = false
					}
				})
			})
			if s.Valid() {
				s.Label("Rename sheet")
				if s.Dismissed() {
					*dismissed = true
				}
			}
		})
	}
}

// TestSheetIgnoresClicksOutside checks the sheet on every platform: a click outside its panel, on the backdrop, leaves it open and reaches nothing behind it, a click inside it reaches its controls, typing goes to its field, and Escape closes it, which the panel reports with Dismissed.
func TestSheetIgnoresClicksOutside(t *testing.T) {
	for _, p := range Platforms {
		open, name, behind, dismissed := true, "", 0, false
		tester := ui.NewTester(dialogsSheetView(p, &open, &name, &behind, &dismissed), 600, 400)
		tester.Frame()
		panel, ok := tester.Find("Rename sheet")
		if !ok {
			t.Fatalf("%v: the sheet is missing: %q", p, tester.Texts())
		}
		// The window's button, under the backdrop, and the corners outside the panel.
		for _, at := range [][2]float32{{30, 26}, {8, 8}, {590, 390}, {panel.X - 4, panel.Y + panel.H/2}, {panel.X + panel.W/2, panel.Y + panel.H + 4}} {
			tester.ClickAt(at[0], at[1])
			tester.Frame()
			if !open || dismissed || behind != 0 {
				t.Fatalf("%v: a click at %v outside the sheet: open %v, dismissed %v, the button behind clicked %d times", p, at, open, dismissed, behind)
			}
		}
		field, ok := tester.Find("Name")
		if !ok {
			t.Fatalf("%v: the sheet's field is missing", p)
		}
		tester.ClickAt(field.X+field.W/2, field.Y+field.H/2)
		tester.Frame()
		tester.Type("Notes")
		tester.Frame()
		if name != "Notes" {
			t.Fatalf("%v: typing in the sheet's field gave %q", p, name)
		}
		tester.Key(0, ui.KeyEscape)
		tester.Frame()
		if open || !dismissed {
			t.Fatalf("%v: Escape left the sheet open %v, dismissed %v", p, open, dismissed)
		}
		if behind != 0 {
			t.Fatalf("%v: the button behind the sheet was clicked", p)
		}
		tester.Frame()
		if err := tester.Click("Behind"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		if behind != 1 {
			t.Fatalf("%v: once the sheet closed, the button behind it took %d clicks", p, behind)
		}
	}
}

// TestSheetCloses checks that the sheet's own buttons close it without Dismissed, which reports Escape only.
func TestSheetCloses(t *testing.T) {
	for _, p := range Platforms {
		open, name, behind, dismissed := true, "", 0, false
		tester := ui.NewTester(dialogsSheetView(p, &open, &name, &behind, &dismissed), 600, 400)
		tester.Frame()
		if err := tester.Click("Cancel"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		if open || dismissed {
			t.Fatalf("%v: Cancel left the sheet open %v, dismissed %v", p, open, dismissed)
		}
	}
}

// TestSheetLooks checks where each platform puts the sheet and what it lays over the window: macOS 27 centers it on the whole window, the title bar's room included, and dims the window with black at 20% (23% in Dark Aqua), filling the sheet with the window's background; Windows lays the smoke layer, black at 30%, and GNOME libadwaita's dimming, the shade at 14% (50% dark), both centering the dialog.
func TestSheetLooks(t *testing.T) {
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			open, name, behind, dismissed := true, "", 0, false
			var th *ui.Theme
			view := dialogsSheetView(p, &open, &name, &behind, &dismissed)
			tester := ui.NewTester(func(c *ui.Context) {
				view(c)
				th = Theme(c, p)
			}, 600, 400)
			tester.SetDark(dark)
			tester.SetTitleBar(ui.TitleBar{Height: 28, Left: 80})
			tester.Frame()
			panel, ok := tester.Find("Rename sheet")
			if !ok {
				t.Fatalf("%v: the sheet is missing", p)
			}
			if cx, cy := panel.X+panel.W/2, panel.Y+panel.H/2; cx < 299 || cx > 301 || cy < 199 || cy > 201 {
				t.Errorf("%v dark %v: the sheet's center is at (%v, %v), want the window's (300, 200)", p, dark, cx, cy)
			}
			var dim ui.Color
			switch p {
			case Windows:
				dim = th.Background.Mix(black, 0.3)
			case Linux:
				dim = th.Background.Mix(adwInk, pick(dark, float32(0.14), 0.5))
			default:
				dim = th.Background.Mix(black, pick(dark, float32(0.20), 0.23))
			}
			if got := dialogsAt(tester, 590, 390); !dialogsNear(got, dim, 2) {
				t.Errorf("%v dark %v: the window behind the sheet reads %v, want %v", p, dark, got, dim)
			}
			if p == Mac {
				if got := dialogsAt(tester, panel.X+panel.W/2, panel.Y+6); !dialogsNear(got, th.Background, 1) {
					t.Errorf("dark %v: the macOS sheet is filled with %v, want the window's %v", dark, got, th.Background)
				}
			}
		}
	}
}

// TestMacSheetCorner checks the macOS sheet's corner against the capture of macOS 27.2's (tools/appkit-reference, appkit-dialogs-light-sheet-live.png): how far in from the panel's left edge, at 2x, each second row below its top the sheet starts, its white past the window dimmed to 204 and the shadow, which the capture leaves out, darker still.
func TestMacSheetCorner(t *testing.T) {
	open, name, behind, dismissed := true, "", 0, false
	tester := ui.NewTester(dialogsSheetView(Mac, &open, &name, &behind, &dismissed), 600, 400)
	tester.SetScale(2)
	tester.Frame()
	panel, ok := tester.Find("Rename sheet")
	if !ok {
		t.Fatal("the sheet is missing")
	}
	// The panel is centered, so its top may fall inside a pixel: the first row is the one it covers more than half of.
	x0, y0 := int(panel.X*2+0.5), int(panel.Y*2+0.5)
	want := []int{27, 19, 15, 12, 9, 8, 6, 5, 3, 3, 2, 1, 1, 1, 0}
	for i, w := range want {
		y := y0 + 2*i
		x := x0
		for x < x0+60 && dialogsAt(tester, float32(x), float32(y)).R <= 206 {
			x++
		}
		if d := x - x0 - w; d < -2 || d > 2 {
			t.Errorf("%d pixels below its top the sheet starts %d pixels in, macOS's %d", 2*i, x-x0, w)
		}
	}
}

// dialogsButtonsView lays out a destructive button beside a plain one and a disabled destructive one; clicks counts the destructive button's clicks.
func dialogsButtonsView(p Platform, clicks *int) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Row(c).Fill().Padding(16).Gap(12).AlignItems(ui.Start).Children(func() {
				Button(c, "Keep")
				if DestructiveButton(c, "Delete").Clicked() {
					*clicks++
				}
				DestructiveButton(c, "Erase").Disabled(true)
			})
		})
	}
}

// dialogsInk returns the color farthest from bg inside r, the label's ink, scanning its middle row. Find gives a button's label, whose text is inside the button's padding, so the face is read 4 points before it.
func dialogsInk(tester *ui.Tester, r ui.Rect, bg ui.Color) ui.Color {
	far, ink := -1, bg
	for x := r.X; x < r.X+r.W; x += 0.5 {
		px := dialogsAt(tester, x, r.Y+r.H/2)
		d := func(a, b uint8) int { return max(int(a)-int(b), int(b)-int(a)) }
		if s := d(px.R, bg.R) + d(px.G, bg.G) + d(px.B, bg.B); s > far {
			far, ink = s, px
		}
	}
	return ink
}

// TestDestructiveButtonLooks checks the destructive button's face and label on every platform: on macOS the window tinted with systemRed at NSAlert's share under a systemRed label, on GNOME destructive_bg_color under white, and on Windows a plain button's face and text; and that it clicks as a button does, and not while disabled.
func TestDestructiveButtonLooks(t *testing.T) {
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			clicks := 0
			var th *ui.Theme
			view := dialogsButtonsView(p, &clicks)
			tester := ui.NewTester(func(c *ui.Context) {
				view(c)
				th = Theme(c, p)
			}, 400, 120)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.Frame()
			del, ok := tester.Find("Delete")
			keep, _ := tester.Find("Keep")
			if !ok {
				t.Fatalf("%v: the destructive button is missing", p)
			}
			face := dialogsAt(tester, 2*(del.X-4), 2*(del.Y+del.H/2))
			var wantFace, wantInk ui.Color
			switch p {
			case Windows:
				wantFace, wantInk = dialogsAt(tester, 2*(keep.X-4), 2*(keep.Y+keep.H/2)), th.Text
			case Linux:
				wantFace, wantInk = th.Danger, white
			default:
				wantFace, wantInk = th.Background.Mix(th.Danger, macDestructiveShare(dark)), th.Danger
			}
			if !dialogsNear(face, wantFace, 2) {
				t.Errorf("%v dark %v: the destructive button's face is %v, want %v", p, dark, face, wantFace)
			}
			scaled := ui.Rect{X: 2 * del.X, Y: 2 * del.Y, W: 2 * del.W, H: 2 * del.H}
			if ink := dialogsInk(tester, scaled, face); !dialogsNear(ink, wantInk, 40) {
				t.Errorf("%v dark %v: the destructive button's label is %v, want %v", p, dark, ink, wantInk)
			}
			if err := tester.Click("Delete"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			if err := tester.Click("Erase"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			if clicks != 1 {
				t.Errorf("%v: the destructive button took %d clicks, want 1", p, clicks)
			}
		}
	}
}

// TestDestructiveAlertDialogKeepsTheDefaultSafe checks the destructive alert on every platform with the destructive button last, first and in the middle: Enter clicks the last button that is not destructive, never the destructive one, a click returns its button's index whatever the order the platform places the buttons in, and Escape dismisses.
func TestDestructiveAlertDialogKeepsTheDefaultSafe(t *testing.T) {
	for _, p := range Platforms {
		for _, tc := range []struct {
			destructive int
			buttons     []string
			def         int
		}{
			{1, []string{"Cancel", "Delete"}, 0},
			{0, []string{"Delete", "Cancel"}, 1},
			{1, []string{"Cancel", "Delete", "Keep"}, 2},
		} {
			for _, act := range []struct {
				name string
				do   func(tester *ui.Tester)
				want int
			}{
				{"Enter", func(tester *ui.Tester) { tester.Key(0, ui.KeyEnter) }, tc.def},
				{"Escape", func(tester *ui.Tester) { tester.Key(0, ui.KeyEscape) }, ui.AlertDismissed},
				{"click Delete", func(tester *ui.Tester) { _ = tester.Click("Delete") }, tc.destructive},
				{"click Cancel", func(tester *ui.Tester) { _ = tester.Click("Cancel") }, indexOf(tc.buttons, "Cancel")},
			} {
				open, got := true, -1
				tester := ui.NewTester(func(c *ui.Context) {
					Preview(c, p, func() {
						if r := DestructiveAlertDialog(c, &open, "Delete Notes?", "You can't undo this.", tc.destructive, tc.buttons...); r != -1 {
							got = r
						}
					})
				}, 600, 400)
				tester.Frame()
				if !tester.HasText("Delete Notes?") || !tester.HasText("You can't undo this.") {
					t.Fatalf("%v: the alert's text is missing", p)
				}
				if !tester.Focused(tc.buttons[tc.def]) {
					t.Fatalf("%v %v: %q does not have the focus", p, tc.buttons, tc.buttons[tc.def])
				}
				act.do(tester)
				tester.Frame()
				if got != act.want || open {
					t.Fatalf("%v %v %s: got %d, open %v; want %d, closed", p, tc.buttons, act.name, got, open, act.want)
				}
			}
		}
	}
}

func indexOf(s []string, v string) int {
	for i, x := range s {
		if x == v {
			return i
		}
	}
	return -1
}

// TestDestructiveAlertDialogLooks checks the buttons of the destructive alert: the destructive one drawn as DestructiveButton is on the platform, systemRed on macOS and destructive_bg_color on GNOME; the default in the accent only on Windows, where ContentDialog's DefaultButton is, and plain on macOS, whose NSAlert then has no default, and on GNOME.
func TestDestructiveAlertDialogLooks(t *testing.T) {
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			open := true
			var th *ui.Theme
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					th = c.Theme()
					DestructiveAlertDialog(c, &open, "Delete “Notes”?", "You can’t undo this action.", 1, "Cancel", "Delete")
				})
			}, 560, 340)
			tester.SetDark(dark)
			tester.Frame()
			del, ok := tester.Find("Delete")
			cancel, _ := tester.Find("Cancel")
			if !ok {
				t.Fatalf("%v: the alert's destructive button is missing", p)
			}
			delFace := dialogsAt(tester, del.X-4, del.Y+del.H/2)
			cancelFace := dialogsAt(tester, cancel.X-4, cancel.Y+cancel.H/2)
			switch p {
			case Linux:
				if !dialogsNear(delFace, th.Danger, 2) {
					t.Errorf("dark %v: GNOME's destructive response is %v, want %v", dark, delFace, th.Danger)
				}
			case Mac:
				if delFace.R <= delFace.G+10 {
					t.Errorf("dark %v: macOS's destructive button's face %v is not tinted red", dark, delFace)
				}
			case Windows:
				if !dialogsNear(cancelFace, th.Accent, 30) {
					t.Errorf("dark %v: Windows's default, Cancel, is %v, want the accent %v", dark, cancelFace, th.Accent)
				}
				if dialogsNear(delFace, th.Accent, 30) {
					t.Errorf("dark %v: Windows's destructive button is in the accent", dark)
				}
			}
			if p != Windows && dialogsNear(cancelFace, th.Accent, 30) {
				t.Errorf("%v dark %v: the default, Cancel, is in the accent %v", p, dark, cancelFace)
			}
		}
	}
}

// dialogsAlertView is the alert of TestScreenshots made destructive, over the same window content.
func dialogsAlertView(p Platform) func(c *ui.Context) {
	open := true
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Column(c).Fill().Padding(16).Children(func() { ui.Text(c, "Window content") })
			DestructiveAlertDialog(c, &open, "Delete “Notes”?", "You can’t undo this action.", 1, "Cancel", "Delete")
		})
	}
}

// TestScreenshotsDialogs writes the sheet, the destructive alert and the destructive buttons on every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names. Without it the test only builds the frames.
func TestScreenshotsDialogs(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			plat := strings.ToLower(p.String())
			open, name, behind, dismissed := true, "Notes", 0, false
			sheet := ui.NewTester(dialogsSheetView(p, &open, &name, &behind, &dismissed), 560, 340)
			sheet.SetScale(2)
			sheet.SetDark(dark)
			sheet.Frame()
			if !sheet.HasText("Rename") || !sheet.HasText("Cancel") {
				t.Fatalf("%v: the sheet is missing its texts: %q", p, sheet.Texts())
			}
			writeShot(t, dir, fmt.Sprintf("breeze-sheet-%s-%s.png", plat, mode), sheet)
			alert := ui.NewTester(dialogsAlertView(p), 560, 340)
			alert.SetScale(2)
			alert.SetDark(dark)
			alert.Frame()
			if !alert.HasText("Delete") || !alert.HasText("Cancel") {
				t.Fatalf("%v: the destructive alert is missing its buttons: %q", p, alert.Texts())
			}
			writeShot(t, dir, fmt.Sprintf("breeze-destructive-alert-%s-%s.png", plat, mode), alert)
			clicks := 0
			buttons := ui.NewTester(dialogsButtonsView(p, &clicks), 360, 80)
			buttons.SetScale(2)
			buttons.SetDark(dark)
			buttons.Frame()
			writeShot(t, dir, fmt.Sprintf("breeze-destructive-button-%s-%s.png", plat, mode), buttons)
		}
	}
}

// TestSheetGivesTheFocusBack checks that on every platform the sheet takes the focus as it opens and gives it back, as it closes on Escape, to the button that opened it.
func TestSheetGivesTheFocusBack(t *testing.T) {
	for _, p := range Platforms {
		open, name := false, ""
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				ui.Column(c).Fill().Padding(16).AlignItems(ui.Start).Children(func() {
					if Button(c, "Rename…").Clicked() {
						open = true
					}
				})
				Sheet(c, &open, func() { TextInput(c, &name).Label("Name").AutoFocus() })
			})
		}, 600, 400)
		tester.Frame()
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		if !tester.Focused("Rename…") {
			t.Fatalf("%v: Tab did not focus the button", p)
		}
		tester.Key(0, ui.KeySpace)
		tester.Frame()
		tester.Frame()
		if !open || !tester.Focused("Name") {
			t.Fatalf("%v: the sheet did not open with the focus in its field: open %v", p, open)
		}
		tester.Key(0, ui.KeyEscape)
		tester.Frame()
		tester.Frame()
		if open || !tester.Focused("Rename…") {
			t.Fatalf("%v: closing the sheet left open %v and the focus away from the button that opened it", p, open)
		}
	}
}
