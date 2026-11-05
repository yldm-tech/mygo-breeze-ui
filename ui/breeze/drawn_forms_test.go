package breeze

import (
	"fmt"
	"image/color"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

func init() {
	for _, n := range []string{"CheckboxGroup", "NumberInput", "Split", "SplitVertical", "SplitPercent", "SplitVerticalPercent"} {
		drawnWidgets[n] = true
	}
}

// formsShotState is what formsShotView's widgets edit.
type formsShotState struct {
	mail, calendar, photos, all1, all2 bool
	copies, full                       float64
	name, address                      string
	alarm, due                         time.Time
	split                              float32
}

func newFormsShotState() *formsShotState {
	return &formsShotState{
		mail: true, all1: true, all2: true,
		copies: 5, full: 10,
		name:  "Ada",
		alarm: time.Date(2026, 1, 5, 15, 4, 0, 0, time.UTC), due: time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
		split: 120,
	}
}

// formsShotView lays out the widgets of the forms group for screenshots: a mixed and a checked check box group, two number inputs (one at its maximum), a field outside a form, a fieldset, a time and a date input, and a split.
func formsShotView(p Platform, s *formsShotState) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Row(c).Fill().Padding(16).Gap(24).AlignItems(ui.Start).Children(func() {
				ui.Column(c).Width(260).Gap(16).Children(func() {
					CheckboxGroup(c, "Notifications", func() {
						Checkbox(c, &s.mail, "Mail")
						Checkbox(c, &s.calendar, "Calendar")
						Checkbox(c, &s.photos, "Photos")
					})
					CheckboxGroup(c, "Sync", func() {
						Checkbox(c, &s.all1, "Contacts")
						Checkbox(c, &s.all2, "Notes")
					})
					NumberInput(c, &s.copies, 0, 10, 1)
					NumberInput(c, &s.full, 0, 10, 1)
				})
				ui.Column(c).Width(300).Gap(16).Children(func() {
					Field(c, "Name", func() { TextInput(c, &s.name) })
					Fieldset(c, "Shipping", func() {
						Field(c, "Address", func() { TextInput(c, &s.address).Placeholder("Street") })
					})
					ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
						TimeInput(c, &s.alarm).Label("Alarm")
						DateInput(c, &s.due).Label("Due")
					})
					Split(c, &s.split, func() {
						ui.Text(c, "Sidebar").Padding(8)
					}, func() {
						ui.Text(c, "Content").Padding(8)
					}).Height(60).Border(1, c.Theme().Border)
				})
			})
		})
	}
}

// TestScreenshotsWave2Forms writes the forms group of every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names, with the pointer over the split's divider. Without it the test only builds the frames.
func TestScreenshotsWave2Forms(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			s := newFormsShotState()
			tester := ui.NewTester(formsShotView(p, s), 640, 440)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tester.SetRegion(i18n.Region{Locale: "en-US", FirstWeekday: time.Sunday, HourCycle: 24})
			tester.Frame()
			tester.Frame()
			for _, want := range []string{"Notifications", "Mail", "Shipping", "Name", "Sidebar"} {
				if !tester.HasText(want) {
					t.Fatalf("%v: %q is missing", p, want)
				}
			}
			side, ok := tester.Find("Sidebar")
			if !ok {
				t.Fatalf("%v: no sidebar pane", p)
			}
			// The divider is s.split DIPs after the first pane's start, where its text starts.
			tester.Move(side.X+s.split+0.5, side.Y+side.H/2)
			tester.Frame()
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			writeShot(t, dir, fmt.Sprintf("breeze-forms-%s-%s.png", strings.ToLower(p.String()), mode), tester)
		}
	}
}

// formsTester builds view on platform p, in the platform's theme, in direction dir without motion, at scale 1, and lays out a first frame.
func formsTester(p Platform, dir ui.Direction, w, h int, view func(c *ui.Context)) *ui.Tester {
	return formsTesterAt(p, dir, w, h, 1, view)
}

// formsTesterAt is formsTester at scale: its Image has scale pixels a DIP.
func formsTesterAt(p Platform, dir ui.Direction, w, h int, scale float32, view func(c *ui.Context)) *ui.Tester {
	tester := ui.NewTester(func(c *ui.Context) { Preview(c, p, func() { Use(c); view(c) }) }, w, h)
	tester.SetScale(scale)
	tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tester.SetDirection(dir)
	tester.Frame()
	tester.Frame()
	return tester
}

// formsFind returns the box of the element showing s, failing the test without it.
func formsFind(t *testing.T, tester *ui.Tester, s string) ui.Rect {
	t.Helper()
	r, ok := tester.Find(s)
	if !ok {
		t.Fatalf("no %q among %q", s, tester.Texts())
	}
	return r
}

// formsScaled returns r in the pixels of an image of scale pixels a DIP.
func formsScaled(r ui.Rect, scale float32) ui.Rect {
	return ui.Rect{X: r.X * scale, Y: r.Y * scale, W: r.W * scale, H: r.H * scale}
}

// formsInkRows returns the first and last rows of r, in pixels, with a pixel a quarter away from bg.
func formsInkRows(tester *ui.Tester, r ui.Rect, bg ui.Color) (first, last int) {
	img := tester.Image()
	first, last = -1, -1
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if !formsNear(img.RGBAAt(x, y), bg, 64) {
				if first < 0 {
					first = y
				}
				last = y
				break
			}
		}
	}
	return first, last
}

// formsNear reports whether a pixel is within tol of c in every channel.
func formsNear(px color.RGBA, c ui.Color, tol int) bool {
	d := func(a, b uint8) bool { return int(a)-int(b) <= tol && int(b)-int(a) <= tol }
	return d(px.R, c.R) && d(px.G, c.G) && d(px.B, c.B)
}

// formsInk counts the pixels of r, in DIPs at scale 1, more than half way from bg to fg, and returns the farthest one's share of the way.
func formsInk(tester *ui.Tester, r ui.Rect, bg, fg ui.Color) (count int, darkest float32) {
	img := tester.Image()
	lum := func(r, g, b uint8) float32 { return 0.299*float32(r) + 0.587*float32(g) + 0.114*float32(b) }
	lb, lf := lum(bg.R, bg.G, bg.B), lum(fg.R, fg.G, fg.B)
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			px := img.RGBAAt(x, y)
			share := (lum(px.R, px.G, px.B) - lb) / (lf - lb)
			if share > 0.5 {
				count++
			}
			darkest = max(darkest, share)
		}
	}
	return count, darkest
}

// formsInkEdge returns where the ink of r, in DIPs at scale 1, starts in the reading direction: its first column from the left, or right to left its last one's right edge, a quarter away from bg.
func formsInkEdge(tester *ui.Tester, r ui.Rect, bg ui.Color, rtl bool) float32 {
	img := tester.Image()
	inked := func(x int) bool {
		for y := int(r.Y); y < int(r.Y+r.H); y++ {
			if px := img.RGBAAt(x, y); !formsNear(px, bg, 64) {
				return true
			}
		}
		return false
	}
	if rtl {
		for x := int(r.X + r.W); x >= int(r.X); x-- {
			if inked(x) {
				return float32(x + 1)
			}
		}
		return r.X
	}
	for x := int(r.X); x <= int(r.X+r.W); x++ {
		if inked(x) {
			return float32(x)
		}
	}
	return r.X + r.W
}

// formsCheckBoxAt returns the center of the box before the label at r of a check box on platform p: macCheckbox's 16-point box 6 before its label, fluentCheckbox's and adwCheckbox's 20 DIPs 8 and 4 before theirs, before it in the reading direction.
func formsCheckBoxAt(p Platform, r ui.Rect, rtl bool) (x, y float32) {
	size, gap := float32(16), float32(6)
	switch p {
	case Windows:
		size, gap = 20, 8
	case Linux:
		size, gap = 20, 4
	}
	if rtl {
		return r.X + r.W + gap + size/2, r.Y + r.H/2
	}
	return r.X - gap - size/2, r.Y + r.H/2
}

// formsIconRect returns where MyGo's 14-DIP calendar icon is in a date field r padded pad at its ends.
func formsIconRect(r ui.Rect, pad float32, rtl bool) ui.Rect {
	x := r.X + r.W - pad - 14
	if rtl {
		x = r.X + pad
	}
	return ui.Rect{X: x, Y: r.Y + (r.H-14)/2, W: 14, H: 14}
}

// TestCheckboxGroupKeepsItsBehavior checks the rebuilt check box group on every platform, left to right and right to left: mixed with one box of two checked, it draws the platform's dash in the accent's text color; its boxes are indented from the start by the platform's indent; a click checks them all and Space then none, each reported once by Changed; a click on a box makes it mixed again, which a click on the group checks all; Tab goes to the box over them and then to each box; and the window holds no group's state after the build.
func TestCheckboxGroupKeepsItsBehavior(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			mail, calendar, changes := true, false, 0
			var th *ui.Theme
			tester := formsTesterAt(p, dir, 400, 300, 2, func(c *ui.Context) {
				th = c.Theme()
				ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() {
					if CheckboxGroup(c, "Notifications", func() {
						Checkbox(c, &mail, "Mail")
						Checkbox(c, &calendar, "Calendar")
					}).Changed() {
						changes++
					}
				})
			})
			rtl := dir == ui.RightToLeft
			head := formsFind(t, tester, "Notifications")
			x, y := formsCheckBoxAt(p, head, rtl)
			if px := tester.Image().RGBAAt(int(2*x), int(2*y)); !formsNear(px, th.AccentText, 40) {
				t.Errorf("%v %v: the mixed box's middle is %v, not the dash's %v", p, dir, px, th.AccentText)
			}
			m := formsFind(t, tester, "Mail")
			in := m.X - head.X
			if rtl {
				in = head.X + head.W - (m.X + m.W)
			}
			indent, _ := formsCheckGroupSpacing(p)
			if d := in - indent; d < -0.5 || d > 0.5 {
				t.Errorf("%v %v: the boxes' labels are %v in from the head's, want %v", p, dir, in, indent)
			}
			if err := tester.Click("Notifications"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			if !mail || !calendar || changes != 1 {
				t.Fatalf("%v %v: a click on the mixed group left mail %v, calendar %v after %d changes", p, dir, mail, calendar, changes)
			}
			tester.Key(0, ui.KeySpace)
			tester.Frame()
			if mail || calendar || changes != 2 {
				t.Fatalf("%v %v: Space on the checked group left mail %v, calendar %v after %d changes", p, dir, mail, calendar, changes)
			}
			if err := tester.Click("Calendar"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			if mail || !calendar || changes != 2 {
				t.Fatalf("%v %v: a click on a box gave mail %v, calendar %v after %d group changes", p, dir, mail, calendar, changes)
			}
			if err := tester.Click("Notifications"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			if !mail || !calendar || changes != 3 {
				t.Fatalf("%v %v: a click on the group mixed again left mail %v, calendar %v after %d changes", p, dir, mail, calendar, changes)
			}
			if !layoutStateEmpty(&formsCheckGroups) {
				t.Errorf("%v %v: a CheckboxGroup left its state behind", p, dir)
			}
		}
		a, b := false, false
		tester := formsTester(p, ui.LeftToRight, 400, 300, func(c *ui.Context) {
			CheckboxGroup(c, "Sync", func() {
				Checkbox(c, &a, "Contacts")
				Checkbox(c, &b, "Notes")
			})
		})
		for _, want := range []string{"Sync", "Contacts", "Notes"} {
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			if !tester.Focused(want) {
				t.Fatalf("%v: Tab did not reach %q", p, want)
			}
		}
		tester.Key(0, ui.KeySpace)
		tester.Frame()
		if a || !b {
			t.Fatalf("%v: Space on Notes gave contacts %v, notes %v", p, a, b)
		}
	}
}

// TestCheckboxGroupDisabledAndNested checks on every platform that a disabled group refuses clicks on the box over its boxes and on its boxes and takes no focus from Tab, and that a group inside another counts its boxes in its own, as package ui's does.
func TestCheckboxGroupDisabledAndNested(t *testing.T) {
	for _, p := range Platforms {
		mail, calendar := false, false
		tester := formsTester(p, ui.LeftToRight, 400, 300, func(c *ui.Context) {
			CheckboxGroup(c, "Notifications", func() {
				Checkbox(c, &mail, "Mail")
				Checkbox(c, &calendar, "Calendar")
			}).Disabled(true)
		})
		for _, s := range []string{"Notifications", "Mail"} {
			if err := tester.Click(s); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
		}
		if mail || calendar {
			t.Fatalf("%v: a disabled group took a click: mail %v, calendar %v", p, mail, calendar)
		}
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		if tester.Focused("Notifications") || tester.Focused("Mail") {
			t.Fatalf("%v: Tab focused a disabled group", p)
		}
		outer, inner1, inner2 := false, false, false
		nested := formsTester(p, ui.LeftToRight, 400, 400, func(c *ui.Context) {
			CheckboxGroup(c, "All", func() {
				Checkbox(c, &outer, "Outer")
				CheckboxGroup(c, "Inner", func() {
					Checkbox(c, &inner1, "One")
					Checkbox(c, &inner2, "Two")
				})
			})
		})
		if err := nested.Click("All"); err != nil {
			t.Fatal(err)
		}
		nested.Frame()
		if !outer || inner1 || inner2 {
			t.Fatalf("%v: the outer group set outer %v, the inner boxes %v %v", p, outer, inner1, inner2)
		}
		if err := nested.Click("Inner"); err != nil {
			t.Fatal(err)
		}
		nested.Frame()
		if !inner1 || !inner2 {
			t.Fatalf("%v: the inner group set %v %v", p, inner1, inner2)
		}
	}
}

// TestNumberInputKeepsItsBehavior checks the rebuilt number input on every platform, left to right and right to left: the buttons sit where the platform puts them; Increase steps up, reported once by Changed; Up and Down step while the text has the focus; a number typed in range applies; Tab goes from the text to Decrease and then Increase, which Space clicks; Increase is disabled at the maximum, where its click does nothing; and a disabled input refuses its buttons.
func TestNumberInputKeepsItsBehavior(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			value, changes := 5.0, 0
			locked := 3.0
			tester := formsTester(p, dir, 400, 200, func(c *ui.Context) {
				ui.Column(c).Padding(16).Gap(16).AlignItems(ui.Start).Children(func() {
					if NumberInput(c, &value, 0, 10, 1).Label("Copies").Changed() {
						changes++
					}
					NumberInput(c, &locked, 0, 10, 1).Disabled(true)
				})
			})
			inc, dec := formsFind(t, tester, "Increase"), formsFind(t, tester, "Decrease")
			row := formsFind(t, tester, "Copies")
			rtl := dir == ui.RightToLeft
			// The text is at the row's start, 20 DIPs in, and the buttons after it in the reading direction.
			text := ui.Rect{X: row.X + 20, Y: row.Y + row.H/2, W: 1, H: 1}
			if rtl {
				text.X = row.X + row.W - 20
			}
			if after := inc.X > text.X; after == rtl {
				t.Errorf("%v %v: Increase at %v, the row at %v", p, dir, inc, row)
			}
			switch p {
			case Mac:
				if inc.Y >= dec.Y || inc.X != dec.X {
					t.Errorf("%v %v: Increase at %v is not above Decrease at %v", p, dir, inc, dec)
				}
			case Windows:
				// Up then Down, as NumberBox's columns, mirrored right to left.
				if (inc.X < dec.X) == rtl {
					t.Errorf("%v %v: Increase at %v, Decrease at %v", p, dir, inc, dec)
				}
			default:
				// − then +, as GtkSpinButton.
				if (dec.X < inc.X) == rtl {
					t.Errorf("%v %v: Decrease at %v, Increase at %v", p, dir, dec, inc)
				}
			}
			if err := tester.Click("Increase"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			if value != 6 || changes != 1 {
				t.Fatalf("%v %v: Increase gave %v after %d changes", p, dir, value, changes)
			}
			tester.ClickAt(text.X, text.Y)
			tester.Frame()
			tester.Key(0, ui.KeyUp)
			tester.Frame()
			tester.Key(0, ui.KeyUp)
			tester.Frame()
			tester.Key(0, ui.KeyDown)
			tester.Frame()
			if value != 7 || changes != 4 {
				t.Fatalf("%v %v: Up, Up and Down gave %v after %d changes", p, dir, value, changes)
			}
			tester.Key(0, ui.KeyEnd)
			tester.Key(0, ui.KeyBackspace)
			tester.Type("9")
			tester.Frame()
			if value != 9 {
				t.Fatalf("%v %v: typing 9 gave %v", p, dir, value)
			}
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			if !tester.Focused("Decrease") {
				t.Fatalf("%v %v: Tab from the text did not reach Decrease", p, dir)
			}
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			if !tester.Focused("Increase") {
				t.Fatalf("%v %v: Tab from Decrease did not reach Increase", p, dir)
			}
			tester.Key(0, ui.KeySpace)
			tester.Frame()
			if value != 10 {
				t.Fatalf("%v %v: Space on Increase gave %v", p, dir, value)
			}
			if err := tester.Click("Increase"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			if value != 10 || changes != 6 {
				t.Fatalf("%v %v: Increase at the maximum gave %v after %d changes", p, dir, value, changes)
			}
			// The disabled input's buttons: Tab passes them by, and a click does nothing.
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			tester.Key(0, ui.KeySpace)
			tester.Frame()
			if locked != 3 || value != 10 {
				t.Fatalf("%v %v: after Increase, Tab and Space the disabled input is at %v, the other at %v", p, dir, locked, value)
			}
		}
		// A click on a disabled input's button does nothing.
		locked := 3.0
		tester := formsTester(p, ui.LeftToRight, 400, 200, func(c *ui.Context) {
			ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() { NumberInput(c, &locked, 0, 10, 1).Disabled(true) })
		})
		if err := tester.Click("Increase"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		if locked != 3 {
			t.Fatalf("%v: a click on a disabled input's Increase gave %v", p, locked)
		}
	}
}

// TestNumberInputRoundsAndClamps checks on every platform that a click on the field just inside its edge focuses the text, as a TextBox's and an entry's padding do; that a number typed out of range does not apply; and that once the text loses the focus it shows the value rounded to step's decimals: Backspace at its end then removes the last of those decimals.
func TestNumberInputRoundsAndClamps(t *testing.T) {
	for _, p := range Platforms {
		value := 0.5
		tester := formsTester(p, ui.LeftToRight, 400, 200, func(c *ui.Context) {
			ui.Column(c).Padding(16).Gap(16).AlignItems(ui.Start).Children(func() {
				NumberInput(c, &value, 0, 2, 0.25).Label("Share")
				Button(c, "Other")
			})
		})
		row := formsFind(t, tester, "Share")
		tester.ClickAt(row.X+2, row.Y+row.H/2)
		tester.Frame()
		tester.Key(0, ui.KeyUp)
		tester.Frame()
		if value != 0.75 {
			t.Fatalf("%v: a click inside the field's edge, then Up, gave %v", p, value)
		}
		// 90.75, past hi.
		tester.Key(0, ui.KeyHome)
		tester.Type("9")
		tester.Frame()
		if value != 0.75 {
			t.Fatalf("%v: typing a number out of range applied %v", p, value)
		}
		if err := tester.Click("Other"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		tester.ClickAt(row.X+2, row.Y+row.H/2)
		tester.Frame()
		tester.Key(0, ui.KeyEnd)
		tester.Key(0, ui.KeyBackspace)
		tester.Frame()
		if value != 0.7 {
			t.Fatalf("%v: Backspace after the blur gave %v, as if the text had not shown 0.75 again", p, value)
		}
	}
}

// TestSplitDividerStaysBorderUnderThePointer checks the splits on every platform: under the pointer the divider stays the theme's Border, where ui.Split's turns to the accent; the panes build in the window's theme, with its accent, and lay out as ui.Split's; a drag and the arrows still move the divider, which Changed on what Split returns reports; and the vertical and percent splits keep their keys.
func TestSplitDividerStaysBorderUnderThePointer(t *testing.T) {
	for _, p := range Platforms {
		var editor [2]ui.Rect
		for i, ours := range []bool{true, false} {
			size, changes := float32(120), 0
			var th *ui.Theme
			var paneAccent ui.Color
			tester := formsTester(p, ui.LeftToRight, 600, 200, func(c *ui.Context) {
				th = c.Theme()
				first := func() {
					paneAccent = c.Theme().Accent
					ui.Text(c, "Files")
				}
				second := func() { ui.Text(c, "Editor") }
				var s ui.Element
				if ours {
					s = Split(c, &size, first, second).Fill()
				} else {
					s = ui.Split(c, &size, first, second).Fill()
				}
				if s.Changed() {
					changes++
				}
			})
			files := formsFind(t, tester, "Files")
			editor[i] = formsFind(t, tester, "Editor")
			if paneAccent != th.Accent {
				t.Errorf("%v (ours %v): a pane builds with accent %v, the window's is %v", p, ours, paneAccent, th.Accent)
			}
			x, y := files.X+size, float32(100)
			tester.Move(x+0.5, y)
			tester.Frame()
			want := th.Border
			if !ours {
				want = th.Accent
			}
			if px := tester.Image().RGBAAt(int(x), int(y)); !formsNear(px, want, 2) {
				t.Errorf("%v (ours %v): the divider under the pointer is %v, want %v", p, ours, px, want)
			}
			tester.Press(x+0.5, y)
			tester.Frame()
			tester.Move(x+40.5, y)
			tester.Frame()
			if px := tester.Image().RGBAAt(int(x+40), int(y)); ours && !formsNear(px, th.Border, 2) {
				t.Errorf("%v: the divider while dragged is %v, want %v", p, px, th.Border)
			}
			tester.Release(x+40.5, y)
			tester.Frame()
			if size != 160 {
				t.Fatalf("%v (ours %v): a drag of 40 left the split at %v", p, ours, size)
			}
			tester.Key(0, ui.KeyRight)
			tester.Frame()
			if size != 170 || changes != 2 {
				t.Fatalf("%v (ours %v): Right after the drag gave %v after %d changes", p, ours, size, changes)
			}
		}
		if editor[0] != editor[1] {
			t.Errorf("%v: the second pane's text is at %v, with ui.Split at %v", p, editor[0], editor[1])
		}
		height, share := float32(100), float32(40)
		tester := formsTester(p, ui.LeftToRight, 600, 400, func(c *ui.Context) {
			ui.Column(c).Fill().Gap(8).Children(func() {
				SplitVertical(c, &height, func() { ui.Text(c, "Top") }, func() { ui.Text(c, "Bottom") }).Height(240)
				SplitVerticalPercent(c, &share, func() { ui.Text(c, "Upper") }, func() { ui.Text(c, "Lower") }).Height(140)
			})
		})
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		tester.Key(0, ui.KeyDown)
		tester.Frame()
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		tester.Key(0, ui.KeyUp)
		tester.Frame()
		if height <= 100 || share >= 40 {
			t.Fatalf("%v: Down left the vertical split at %v, Up the vertical percent split at %v", p, height, share)
		}
	}
}

// TestFieldLabelAndLegendWeights checks, left to right and right to left, the label of a field on Windows, drawn at regular weight over MyGo's transparent medium one, and the legend of a fieldset on GNOME, bold over MyGo's semibold one: the label, the control and the description stay where ui.Field and ui.Fieldset put them, a click on the label still focuses the control, and the label's ink is lighter (Windows) or heavier (GNOME) than MyGo's, starting where MyGo's starts.
func TestFieldLabelAndLegendWeights(t *testing.T) {
	for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
		rtl := dir == ui.RightToLeft
		var ink [2]int
		var where [2][3]ui.Rect
		var edge [2]float32
		for i, ours := range []bool{true, false} {
			name := ""
			var th *ui.Theme
			tester := formsTester(Windows, dir, 400, 200, func(c *ui.Context) {
				th = c.Theme()
				ui.Column(c).Padding(16).Children(func() {
					build := func() { TextInput(c, &name).Label("Input") }
					if ours {
						Field(c, "Full name", build).Description("As on your card.")
					} else {
						ui.Field(c, "Full name", build).Gap(8).Description("As on your card.")
					}
				})
			})
			label := formsFind(t, tester, "Full name")
			where[i] = [3]ui.Rect{label, formsFind(t, tester, "Input"), formsFind(t, tester, "As on your card.")}
			box := ui.Rect{X: label.X, Y: label.Y, W: label.W, H: 20}
			ink[i], _ = formsInk(tester, box, th.Background, th.Text)
			edge[i] = formsInkEdge(tester, box, th.Background, rtl)
			tester.ClickAt(edge[i]+map[bool]float32{false: 3, true: -3}[rtl], label.Y+8)
			tester.Frame()
			tester.Type("Ada")
			tester.Frame()
			if name != "Ada" {
				t.Fatalf("Windows %v (ours %v): typing after a click on the label gave %q", dir, ours, name)
			}
		}
		if where[0] != where[1] {
			t.Errorf("Windows %v: label, input and description at %v, with ui.Field at %v", dir, where[0], where[1])
		}
		if ink[0] >= ink[1] || ink[0] == 0 {
			t.Errorf("Windows %v: the label has %d dark pixels, MyGo's %d", dir, ink[0], ink[1])
		}
		if d := edge[0] - edge[1]; d < -1.5 || d > 1.5 {
			t.Errorf("Windows %v: the label's ink starts at %v, MyGo's at %v", dir, edge[0], edge[1])
		}

		var legendInk [2]int
		var legendRows [2][2]int
		var legendEdge [2]float32
		for i, ours := range []bool{true, false} {
			address := ""
			var th *ui.Theme
			tester := formsTester(Linux, dir, 400, 300, func(c *ui.Context) {
				th = c.Theme()
				ui.Column(c).Padding(16).Children(func() {
					build := func() {
						Field(c, "Address", func() { TextInput(c, &address).Label("Address input") })
					}
					if ours {
						Fieldset(c, "Shipping", build)
					} else {
						ui.Fieldset(c, "Shipping", build).Gap(12)
					}
				})
			})
			legend := formsFind(t, tester, "Shipping")
			box := ui.Rect{X: legend.X, Y: legend.Y, W: legend.W, H: 20}
			legendRows[i][0], legendRows[i][1] = formsInkRows(tester, box, th.Background)
			legendInk[i], _ = formsInk(tester, box, th.Background, th.Text)
			legendEdge[i] = formsInkEdge(tester, box, th.Background, rtl)
			if ours {
				if err := tester.Click("Address"); err != nil {
					t.Fatal(err)
				}
				tester.Frame()
				tester.Type("Main St")
				tester.Frame()
				if address != "Main St" {
					t.Fatalf("GNOME %v: typing after a click on the label gave %q", dir, address)
				}
			}
		}
		if d0, d1 := legendRows[0][0]-legendRows[1][0], legendRows[0][1]-legendRows[1][1]; d0 < -1 || d0 > 1 || d1 < -1 || d1 > 1 {
			t.Errorf("GNOME %v: the legend's ink spans rows %v, MyGo's %v", dir, legendRows[0], legendRows[1])
		}
		if legendInk[0] <= legendInk[1] {
			t.Errorf("GNOME %v: the bold legend has %d dark pixels, MyGo's semibold %d", dir, legendInk[0], legendInk[1])
		}
		if d := legendEdge[0] - legendEdge[1]; d < -1.5 || d > 1.5 {
			t.Errorf("GNOME %v: the legend's ink starts at %v, MyGo's at %v", dir, legendEdge[0], legendEdge[1])
		}
	}
}

// TestDateAndTimeInputsOnMacAndGNOME checks, left to right and right to left, the parts this group draws of the date and time fields: on macOS the time's colon reaches labelColor (the theme's Text) and Up still steps the hours in the column the field is built in, and the date field shows no calendar icon, enabled or disabled, and still opens its calendar; on GNOME the icon is the text color at 70%.
func TestDateAndTimeInputsOnMacAndGNOME(t *testing.T) {
	for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
		rtl := dir == ui.RightToLeft
		alarm := time.Date(2026, 1, 5, 7, 30, 0, 0, time.UTC)
		due := time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC)
		var th *ui.Theme
		var face ui.Color
		mac := formsTesterAt(Mac, dir, 400, 500, 2, func(c *ui.Context) {
			th = c.Theme()
			face = pick(th.Dark, white, th.Background.Mix(black, macDateFillDark))
			ui.Column(c).Padding(16).Gap(16).AlignItems(ui.Start).Children(func() {
				TimeInput(c, &alarm).Label("Alarm")
				DateInput(c, &due).Label("Due")
				DateInput(c, &due).Label("Locked").Disabled(true)
			})
		})
		mac.SetRegion(i18n.Region{Locale: "en-GB", FirstWeekday: time.Monday, HourCycle: 24})
		mac.Frame()
		colon := formsFind(t, mac, ":")
		if _, darkest := formsInk(mac, formsScaled(colon, 2), face, th.Text); darkest < 0.9 {
			t.Errorf("macOS %v: the colon reaches %v of the text's contrast, want labelColor", dir, darkest)
		}
		for _, name := range []string{"Due", "Locked"} {
			r := formsFind(t, mac, name)
			icon := formsScaled(formsIconRect(r, 4, rtl), 2)
			bg := face
			if name == "Locked" {
				bg = th.Background.Mix(face, 0.5)
			}
			img := mac.Image()
			for y := int(icon.Y); y < int(icon.Y+icon.H); y++ {
				for x := int(icon.X); x < int(icon.X+icon.W); x++ {
					if px := img.RGBAAt(x, y); !formsNear(px, bg, 6) {
						t.Fatalf("macOS %v: %s shows %v at %d,%d where its icon was, its field is %v", dir, name, px, x, y, bg)
					}
				}
			}
		}
		if err := mac.Click("07"); err != nil {
			t.Fatal(err)
		}
		mac.Frame()
		mac.Key(0, ui.KeyUp)
		mac.Frame()
		if alarm.Hour() != 8 {
			t.Fatalf("macOS %v: Up on the hours set %v", dir, alarm)
		}
		if err := mac.Click("Due"); err != nil {
			t.Fatal(err)
		}
		mac.Frame()
		if !mac.HasText("January 2026") {
			t.Fatalf("macOS %v: no calendar after a click: %q", dir, mac.Texts())
		}

		var fill ui.Color
		gnome := formsTester(Linux, dir, 400, 400, func(c *ui.Context) {
			th = c.Theme()
			fill = adwInkAlpha(th, th.Background, 0.10)
			ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() { DateInput(c, &due).Label("Due") })
		})
		r := formsFind(t, gnome, "Due")
		_, darkest := formsInk(gnome, formsIconRect(r, 9, rtl), fill, adwInkAlpha(th, fill, 1))
		if darkest < 0.6 || darkest > 0.8 {
			t.Errorf("GNOME %v: the calendar icon reaches %v of the text color, want 0.7", dir, darkest)
		}
	}
}
