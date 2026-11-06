package breeze

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func init() {
	drawnWidgets["ViewSwitcher"] = true
}

// switcherTestIcons are the icons of the test switchers, Library, Browse and Radio, 16 points square; tools/appkit-reference's switcher scene reads the same SVGs.
var switcherTestIcons = []*ui.SVG{
	ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><path d="M2.5 2.5h3v11h-3zM6.5 2.5h3v11h-3zM10.4 3.4l2.4-.6 2.4 10.2-2.4.6z" fill="none" stroke="#000" stroke-width="1.3" stroke-linejoin="round"/></svg>`)),
	ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><path d="M2.5 2.5h4v4h-4zM9.5 2.5h4v4h-4zM2.5 9.5h4v4h-4zM9.5 9.5h4v4h-4z" fill="none" stroke="#000" stroke-width="1.3" stroke-linejoin="round"/></svg>`)),
	ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><circle cx="8" cy="8" r="1.6"/><path d="M5.3 5.3a3.8 3.8 0 0 0 0 5.4M10.7 5.3a3.8 3.8 0 0 1 0 5.4M3.2 3.2a6.8 6.8 0 0 0 0 9.6M12.8 3.2a6.8 6.8 0 0 1 0 9.6" fill="none" stroke="#000" stroke-width="1.3" stroke-linecap="round"/></svg>`)),
}

// switcherTestLabels are the labels of the test switchers' items.
var switcherTestLabels = []string{"Library", "Browse", "Radio"}

// switcherTestItems returns the test switchers' items, with their icons when icons is true.
func switcherTestItems(icons bool) []ViewItem {
	items := make([]ViewItem, len(switcherTestLabels))
	for i, s := range switcherTestLabels {
		items[i].Label = s
		if icons {
			items[i].Icon = switcherTestIcons[i]
		}
	}
	return items
}

// switcherTestState is what a test switcher view keeps: the chosen item, whether the switcher is in the header bar or in the content, and the switcher's box as last built.
type switcherTestState struct {
	selected int
	inBar    bool
	disabled bool
	box      ui.Rect
}

// view builds a window with a header bar titled Music over its content, a button Before, then a button After; the switcher of Library, Browse and Radio with icons is among the header bar's end controls when inBar, or between the two buttons otherwise.
func (s *switcherTestState) view(c *ui.Context) {
	switcher := func() {
		v := ViewSwitcher(c, &s.selected, switcherTestItems(true)...).Disabled(s.disabled)
		s.box = v.Bounds()
	}
	ui.Column(c).Fill().Children(func() {
		var end func()
		if s.inBar {
			end = switcher
		}
		HeaderBar(c, "Music", "", nil, end)
		ui.Column(c).Padding(16).Gap(12).AlignItems(ui.Start).Children(func() {
			Button(c, "Before")
			if !s.inBar {
				switcher()
			}
			Button(c, "After")
		})
	})
}

// TestViewSwitcherChoosesWithClicksAndKeys checks the switcher on every platform, left to right and right to left, alone and in a header bar: a click chooses an item; Tab reaches the chosen item once, the arrow toward the start chooses the previous item and the one toward the end the next one, the focus following, as MyGo's radio groups and tabs move it; and the next Tab leaves the switcher.
func TestViewSwitcherChoosesWithClicksAndKeys(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			for _, inBar := range []bool{false, true} {
				name := fmt.Sprintf("%v %v in bar %v", p, dir, inBar)
				s := &switcherTestState{inBar: inBar}
				tester := layoutTester(p, 640, 300, s.view)
				tester.SetDirection(dir)
				tester.Frame()
				tester.Frame()
				if err := tester.Click("Radio"); err != nil {
					t.Fatalf("%s: %v", name, err)
				}
				tester.Frame()
				if s.selected != 2 {
					t.Fatalf("%s: a click on Radio chose %d", name, s.selected)
				}
				if err := tester.Click("Before"); err != nil {
					t.Fatalf("%s: %v", name, err)
				}
				tester.Frame()
				press := func(mods ui.Modifiers, k ui.Key) {
					tester.Key(mods, k)
					tester.Frame()
					tester.Frame()
				}
				// In the header bar the switcher comes before the content's buttons.
				into, out := ui.Modifiers(0), "After"
				if inBar {
					into, out = ui.Shift, "Before"
				}
				press(into, ui.KeyTab)
				if !tester.Focused("Radio") {
					t.Fatalf("%s: Tab did not stop at the chosen item", name)
				}
				forth, back := ui.KeyRight, ui.KeyLeft
				if dir == ui.RightToLeft {
					forth, back = back, forth
				}
				press(0, back)
				press(0, back)
				if s.selected != 0 || !tester.Focused("Library") {
					t.Errorf("%s: the arrow toward the start twice chose %d", name, s.selected)
				}
				press(0, forth)
				if s.selected != 1 || !tester.Focused("Browse") {
					t.Errorf("%s: the arrow toward the end chose %d", name, s.selected)
				}
				press(0, ui.KeyTab)
				if !tester.Focused(out) {
					t.Errorf("%s: Tab from the switcher did not leave it for %s", name, out)
				}
			}
		}
	}
}

// TestViewSwitcherNamesItsItems checks on every platform that each item is named by its label, also on macOS in a header bar, whose toolbar group shows the icons only, as AppKit's icon-only toolbar does.
func TestViewSwitcherNamesItsItems(t *testing.T) {
	for _, p := range Platforms {
		for _, inBar := range []bool{false, true} {
			s := &switcherTestState{inBar: inBar}
			tester := layoutTester(p, 640, 300, s.view)
			tester.Frame()
			for _, label := range switcherTestLabels {
				if r, ok := tester.Find(label); !ok || r.W < 30 {
					t.Errorf("%v in bar %v: no item is named %q", p, inBar, label)
				}
			}
		}
	}
}

// TestViewSwitcherDisabled checks on every platform that a disabled switcher takes neither clicks nor the keys.
func TestViewSwitcherDisabled(t *testing.T) {
	for _, p := range Platforms {
		s := &switcherTestState{selected: 1, disabled: true}
		tester := layoutTester(p, 640, 300, s.view)
		tester.Frame()
		if err := tester.Click("Radio"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		if err := tester.Click("Before"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		tester.Key(0, ui.KeyRight)
		tester.Frame()
		if s.selected != 1 {
			t.Errorf("%v: a disabled switcher chose %d", p, s.selected)
		}
	}
}

// switcherSizes builds a switcher with icons and one without on p, in a header bar when inBar, and returns the items' boxes of each, Library, Browse and Radio, and the switchers' boxes.
func switcherSizes(p Platform, inBar bool, labels []string) (with, without [3]ui.Rect, boxes [2]ui.Rect) {
	items := func(icons bool) []ViewItem {
		its := make([]ViewItem, len(labels))
		for i, s := range labels {
			its[i].Label = s
			if icons {
				its[i].Icon = switcherTestIcons[i]
			}
		}
		return its
	}
	a, b := 1, 1
	// The items without icons are named apart, so that the two switchers' items can be told apart.
	plain := items(false)
	for i := range plain {
		plain[i].Label += " "
	}
	tester := layoutTester(p, 900, 300, func(c *ui.Context) {
		build := func() {
			boxes[0] = ViewSwitcher(c, &a, items(true)...).Bounds()
			boxes[1] = ViewSwitcher(c, &b, plain...).Bounds()
		}
		if inBar {
			HeaderBar(c, "Music", "", nil, build)
			return
		}
		ui.Column(c).Padding(16).Gap(12).AlignItems(ui.Start).Children(build)
	})
	tester.Frame()
	tester.Frame()
	tester.Frame()
	for i, s := range labels {
		with[i], _ = tester.Find(s)
		without[i], _ = tester.Find(s + " ")
	}
	return with, without, boxes
}

// TestViewSwitcherSizes checks the switcher's sizes on every platform, without reading the host font's widths: on macOS in a header bar, the toolbar group's 36-point platter of 37-point image segments (111 for three, as AppKit lays out three 16-point images), or of title segments as wide as the widest when an item has no icon; on macOS alone, NSSegmentedControl's 24-point segments, an icon making a segment 18 points wider (16 and 2); on Windows a SelectorBar 48 DIPs high of 40-DIP items, an icon making an item 24 wider (16 and SelectorBarItemSpacing 8); on GNOME buttons 34 pixels high, 3 apart, as wide as the widest and at least 120.
func TestViewSwitcherSizes(t *testing.T) {
	near := func(a, b float32) bool { return a-b < 0.01 && b-a < 0.01 }
	for _, p := range Platforms {
		for _, inBar := range []bool{false, true} {
			name := fmt.Sprintf("%v in bar %v", p, inBar)
			with, without, boxes := switcherSizes(p, inBar, switcherTestLabels)
			switch {
			case p == Mac && inBar:
				if !near(boxes[0].W, 111) || !near(boxes[0].H, 36) {
					t.Errorf("%s: the image group is %v by %v, want 111 by 36", name, boxes[0].W, boxes[0].H)
				}
				for i := range 3 {
					if !near(with[i].W, 37) || !near(without[i].W, without[0].W) || without[i].W < 35 {
						t.Errorf("%s: segment %d is %v wide with an icon and %v with a title, want 37 and the widest's", name, i, with[i].W, without[i].W)
					}
				}
			case p == Mac:
				for i := range 3 {
					if !near(with[i].H, 24) || !near(with[i].W-without[i].W, 18) {
						t.Errorf("%s: segment %d is %v high and %v wider with its icon, want 24 and 18", name, i, with[i].H, with[i].W-without[i].W)
					}
				}
			case p == Windows:
				if !near(boxes[0].H, 48) {
					t.Errorf("%s: the SelectorBar is %v high, want 48", name, boxes[0].H)
				}
				for i := range 3 {
					if !near(with[i].H, 40) || !near(with[i].W-without[i].W, 24) {
						t.Errorf("%s: item %d is %v high and %v wider with its icon, want 40 and 24", name, i, with[i].H, with[i].W-without[i].W)
					}
				}
			default:
				for i := range 3 {
					if !near(with[i].H, 34) || !near(with[i].W, with[0].W) || with[i].W < 120 {
						t.Errorf("%s: button %d is %v by %v, want 34 high, as wide as the first and at least 120", name, i, with[i].W, with[i].H)
					}
				}
				if gap := with[1].X - (with[0].X + with[0].W); !near(gap, 3) {
					t.Errorf("%s: the buttons are %v apart, want 3", name, gap)
				}
			}
		}
	}
	// Equal widths follow the widest label: GNOME's buttons and macOS's title segments in a toolbar.
	for _, p := range []Platform{Mac, Linux} {
		labels := []string{"A", "A much longer label", "B"}
		with, without, _ := switcherSizes(p, true, labels)
		row := with
		if p == Mac {
			row = without
		}
		if !near(row[0].W, row[1].W) || !near(row[2].W, row[1].W) || row[1].W <= 120 {
			t.Errorf("%v: buttons %v, %v and %v wide, want all as wide as the widest", p, row[0].W, row[1].W, row[2].W)
		}
	}
}

// TestViewSwitcherMacToolbarPill checks the macOS toolbar group's pixels, left to right and right to left: the chosen segment's pill is a pressed toolbar item's face (235 in Aqua) beside its icon, an unchosen segment shows the platter there, and a 1-point separator darkens the boundary between two unchosen segments, on the later segment's side.
func TestViewSwitcherMacToolbarPill(t *testing.T) {
	for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
		s := &switcherTestState{selected: 1, inBar: true}
		tester := layoutTester(Mac, 640, 300, func(c *ui.Context) {
			Use(c)
			s.view(c)
		})
		tester.SetDirection(dir)
		tester.Frame()
		tester.Frame()
		img := tester.Image()
		gray := func(x, y float32) int { return int(img.RGBAAt(int(x), int(y)).G) }
		browse, _ := tester.Find("Browse")
		library, _ := tester.Find("Library")
		y := browse.Y + browse.H/2
		// 6 points in from a side of the segment is inside its pill and clear of its centered 16-point icon.
		if v := gray(browse.X+6, y); v < 233 || v > 237 {
			t.Errorf("%v: the chosen segment's pill reads %d, want 235", dir, v)
		}
		if v := gray(library.X+library.W/2-11, y); v < 250 {
			t.Errorf("%v: an unchosen segment reads %d, want the white platter", dir, v)
		}
		s.selected = 0
		tester.Frame()
		img = tester.Image()
		radio, _ := tester.Find("Radio")
		x := radio.X
		if dir == ui.RightToLeft {
			x = radio.X + radio.W - 1
		}
		if v := gray(x, y); v > 215 {
			t.Errorf("%v: the separator before Radio reads %d, want the platter darkened by 20%%", dir, v)
		}
	}
}

// TestViewSwitcherChosenLooks checks the chosen item's look on Windows and GNOME, reading pixels clear of the text: WinUI's accent pill, 16 by 3, centered at the bottom of the chosen SelectorBarItem and on no other; GNOME's chosen button currentColor at 10% over the window, another one transparent.
func TestViewSwitcherChosenLooks(t *testing.T) {
	for _, p := range []Platform{Windows, Linux} {
		var th ui.Theme
		s := &switcherTestState{selected: 1}
		tester := layoutTester(p, 640, 300, func(c *ui.Context) {
			Use(c)
			th = *c.Theme()
			s.view(c)
		})
		tester.Frame()
		tester.Frame()
		img := tester.Image()
		at := func(x, y float32) ui.Color {
			c := img.RGBAAt(int(x), int(y))
			return ui.RGB(c.R, c.G, c.B)
		}
		same := func(a, b ui.Color) bool {
			d := func(x, y uint8) bool { return int(x)-int(y) <= 2 && int(y)-int(x) <= 2 }
			return d(a.R, b.R) && d(a.G, b.G) && d(a.B, b.B)
		}
		browse, _ := tester.Find("Browse")
		radio, _ := tester.Find("Radio")
		if p == Windows {
			for _, r := range []ui.Rect{browse, radio} {
				got := at(r.X+r.W/2, r.Y+r.H-2)
				if chosen := r == browse; same(got, th.Accent) != chosen {
					t.Errorf("%v: the pill's place under %v reads %v, chosen %v", p, r, got, chosen)
				}
				if got := at(r.X+r.W/2-9, r.Y+r.H-2); same(got, th.Accent) {
					t.Errorf("%v: the pill reaches 9 DIPs from the item's middle", p)
				}
			}
			continue
		}
		face := adwInkAlpha(&th, th.Background, 0.10)
		if got := at(browse.X+4, browse.Y+browse.H/2); !same(got, face) {
			t.Errorf("%v: the chosen button reads %v, want %v", p, got, face)
		}
		if got := at(radio.X+4, radio.Y+radio.H/2); !same(got, th.Background) {
			t.Errorf("%v: an unchosen button reads %v, want the window's %v", p, got, th.Background)
		}
	}
}

// switcherView lays out the switchers for screenshots: a header bar titled Music with a switcher of Library, Browse and Radio among its end controls, and the same switcher alone 16 from the window's side and 16 below the bar, Browse chosen in both. tools/appkit-reference's switcher scene, tools/winui-reference's winui-switcher scene and tools/adwaita-reference's switcher set place the platforms' own there.
func switcherView(p Platform) func(c *ui.Context) {
	bar, alone := 1, 1
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Column(c).Fill().Children(func() {
				HeaderBar(c, "Music", "", nil, func() { ViewSwitcher(c, &bar, switcherTestItems(true)...) })
				ui.Column(c).Grow(1).Padding(16).AlignItems(ui.Start).Children(func() {
					ViewSwitcher(c, &alone, switcherTestItems(true)...)
				})
			})
		})
	}
}

// TestScreenshotsViewSwitcher writes the switchers of every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names, the window controls' room as headerTestRoom gives it. Without it the test only builds the frames.
func TestScreenshotsViewSwitcher(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			tester := ui.NewTester(switcherView(p), 560, 140)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tester.SetTitleBar(headerTestRoom(p, ui.LeftToRight))
			tester.Frame()
			tester.Frame()
			tester.Frame()
			if !tester.HasText("Music") || !tester.HasText("Browse") {
				t.Fatalf("%v: the bar's title or the switcher is missing", p)
			}
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			writeShot(t, dir, fmt.Sprintf("breeze-switcher-%s-%s.png", strings.ToLower(p.String()), mode), tester)
		}
	}
}
