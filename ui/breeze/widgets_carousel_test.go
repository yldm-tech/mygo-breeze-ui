package breeze

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func init() {
	drawnWidgets["PipsPager"] = true
	drawnWidgets["Carousel"] = true
}

// carouselTestSlides are the texts of the test carousel's pages, which tell them apart from the pips' names.
var carouselTestSlides = []string{"Slide 1", "Slide 2", "Slide 3", "Slide 4"}

// carouselTestPageHeight is the height of the test carousel's pages.
const carouselTestPageHeight = 160

// carouselTestPage builds page i of the test carousel: its slide's text in the middle of a page 160 high, and counts the pages built in built, when not nil.
func carouselTestPage(c *ui.Context, built map[int]int) func(i int) {
	return func(i int) {
		if built != nil {
			built[i]++
		}
		ui.Column(c).Height(carouselTestPageHeight).Justify(ui.Center).AlignItems(ui.Center).Children(func() {
			ui.Text(c, carouselTestSlides[i]).SingleLine()
		})
	}
}

// carouselShotPlaces are where carouselShotView puts its three controls, in DIPs from the window's top left corner, from its top right corner right to left: the pager of 5, the pager of 9 and the carousel, 320 wide. tools/appkit-reference's carousel scene, tools/winui-reference's winui-carousel scene and tools/adwaita-reference's carousel set put the platforms' own there.
var carouselShotPlaces = [3]ui.Point{{X: 16, Y: 16}, {X: 16, Y: 56}, {X: 16, Y: 96}}

// carouselShotView lays out the screenshot test's controls on p: a PipsPager of 5 pages with the 2nd chosen, one of 9 with the 7th chosen, and a Carousel 320 wide of 4 pages with the 2nd shown, at carouselShotPlaces.
func carouselShotView(p Platform) func(c *ui.Context) {
	five, nine, shown := 1, 6, 1
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Box(c).Fill().Children(func() {
				at := func(k int, build func()) {
					ui.Box(c).Absolute().Start(carouselShotPlaces[k].X).Top(carouselShotPlaces[k].Y).Children(build)
				}
				at(0, func() { PipsPager(c, &five, 5) })
				at(1, func() { PipsPager(c, &nine, 9) })
				at(2, func() { Carousel(c, &shown, 4, carouselTestPage(c, nil)).Width(320) })
			})
		})
	}
}

// TestScreenshotsCarousel writes the pagers and the carousel of every platform, light and dark, and right to left in the light theme, at 2x to the directory BREEZE_SCREENSHOTS names. Without it the test only builds the frames.
func TestScreenshotsCarousel(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	type shot struct {
		dark bool
		rtl  bool
	}
	for _, p := range Platforms {
		for _, s := range []shot{{false, false}, {true, false}, {false, true}} {
			tester := ui.NewTester(carouselShotView(p), 352, 312)
			tester.SetScale(2)
			tester.SetDark(s.dark)
			tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			if s.rtl {
				tester.SetDirection(ui.RightToLeft)
			}
			tester.Frame()
			tester.Frame()
			tester.Frame()
			if !tester.HasText("Slide 2") {
				t.Fatalf("%v: the carousel does not show its second page", p)
			}
			mode := map[bool]string{false: "light", true: "dark"}[s.dark]
			if s.rtl {
				mode = "rtl-light"
			}
			writeShot(t, dir, fmt.Sprintf("breeze-carousel-%s-%s.png", strings.ToLower(p.String()), mode), tester)
		}
	}
}

// pagerTestState is what a test pager view keeps: the page chosen, the pages, whether the pager is disabled, and its box as last built.
type pagerTestState struct {
	selected, count int
	disabled        bool
	box             ui.Rect
}

// view builds a column of a button Before, the pager and a button After.
func (s *pagerTestState) view(c *ui.Context) {
	Use(c)
	ui.Column(c).Padding(16).Gap(12).AlignItems(ui.Start).Children(func() {
		Button(c, "Before")
		s.box = PipsPager(c, &s.selected, s.count).Disabled(s.disabled).Bounds()
		Button(c, "After")
	})
}

// carouselTestState is what a test carousel view keeps: the page shown, whether the carousel is disabled, the pages built, by index, and its box as last built.
type carouselTestState struct {
	shown    int
	disabled bool
	built    map[int]int
	box      ui.Rect
}

// view builds a column of a button Before, a carousel 320 wide of the 4 test pages and a button After.
func (s *carouselTestState) view(c *ui.Context) {
	Use(c)
	ui.Column(c).Padding(16).Gap(12).AlignItems(ui.Start).Children(func() {
		Button(c, "Before")
		s.box = Carousel(c, &s.shown, 4, carouselTestPage(c, s.built)).Width(320).Disabled(s.disabled).Bounds()
		Button(c, "After")
	})
}

// carouselPress presses a key and lets the frames that follow run.
func carouselPress(tester *ui.Tester, mods ui.Modifiers, k ui.Key) {
	tester.Key(mods, k)
	tester.Frame()
	tester.Frame()
}

// TestPipsPagerChoosesWithClicksAndKeys checks the pager on every platform, left to right and right to left: a click on a pip chooses its page; Tab stops at the pager once, named Pager; the arrow toward the start chooses the page before and the one toward the end the page after, Up and Down as WinUI's PipsPager takes them, stopping at the first and the last page rather than going round; and the next Tab leaves it.
func TestPipsPagerChoosesWithClicksAndKeys(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			name := fmt.Sprintf("%v %v", p, dir)
			s := &pagerTestState{selected: 1, count: 5}
			tester := layoutTester(p, 400, 200, s.view)
			tester.SetDirection(dir)
			tester.Frame()
			tester.Frame()
			if err := tester.Click("Page 4"); err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			tester.Frame()
			if s.selected != 3 {
				t.Fatalf("%s: a click on Page 4 chose %d", name, s.selected)
			}
			if err := tester.Click("Before"); err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			tester.Frame()
			carouselPress(tester, 0, ui.KeyTab)
			if !tester.Focused("Pager") {
				t.Fatalf("%s: Tab did not stop at the pager", name)
			}
			forth, back := ui.KeyRight, ui.KeyLeft
			if dir == ui.RightToLeft {
				forth, back = back, forth
			}
			carouselPress(tester, 0, forth)
			carouselPress(tester, 0, forth)
			if s.selected != 4 {
				t.Errorf("%s: the arrow toward the end twice from the 4th page chose %d, want the last, 4", name, s.selected)
			}
			for range 5 {
				carouselPress(tester, 0, back)
			}
			if s.selected != 0 {
				t.Errorf("%s: the arrow toward the start five times chose %d, want the first, 0", name, s.selected)
			}
			carouselPress(tester, 0, ui.KeyDown)
			carouselPress(tester, 0, ui.KeyDown)
			carouselPress(tester, 0, ui.KeyUp)
			if s.selected != 1 {
				t.Errorf("%s: Down twice and Up chose %d, want 1", name, s.selected)
			}
			carouselPress(tester, 0, ui.KeyTab)
			if !tester.Focused("After") {
				t.Errorf("%s: Tab from the pager did not leave it", name)
			}
		}
	}
}

// TestPipsPagerNamesItsPips checks on every platform that each pip is named by its page, Page 1 to Page 5, as WinUI names its pips, and the pager Pager.
func TestPipsPagerNamesItsPips(t *testing.T) {
	for _, p := range Platforms {
		s := &pagerTestState{count: 5}
		tester := layoutTester(p, 400, 200, s.view)
		tester.Frame()
		for _, label := range []string{"Pager", "Page 1", "Page 2", "Page 3", "Page 4", "Page 5"} {
			if _, ok := tester.Find(label); !ok {
				t.Errorf("%v: nothing is named %q", p, label)
			}
		}
		if got := pageName(nil, 6); got != "Page 7" {
			t.Errorf("the 7th page is named %q, want Page 7", got)
		}
	}
}

// TestPipsPagerDisabled checks on every platform that a disabled pager takes neither clicks nor keys.
func TestPipsPagerDisabled(t *testing.T) {
	for _, p := range Platforms {
		s := &pagerTestState{selected: 1, count: 5, disabled: true}
		tester := layoutTester(p, 400, 200, s.view)
		tester.Frame()
		if err := tester.Click("Page 4"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		if err := tester.Click("Before"); err != nil {
			t.Fatal(err)
		}
		carouselPress(tester, 0, ui.KeyTab)
		carouselPress(tester, 0, ui.KeyRight)
		if s.selected != 1 {
			t.Errorf("%v: a disabled pager chose %d", p, s.selected)
		}
	}
}

// TestPipsPagerSizes checks the pager's size on every platform, for 1, 5 and 9 pages: on macOS UIPageControl's 23 points high and 19 a dot with 8 at either end (111 for 5 and 187 for 9, as UIKit for Mac lays it out); on Windows 24 high and 12 a pip, at most 5 (MaxVisiblePips); on GNOME 20 high and 15 a dot with 6 at either end (AdwCarouselIndicatorDots's measure).
func TestPipsPagerSizes(t *testing.T) {
	for _, p := range Platforms {
		for _, n := range []int{1, 5, 9} {
			s := &pagerTestState{count: n}
			tester := layoutTester(p, 400, 200, s.view)
			tester.Frame()
			var w, h float32
			switch p {
			case Mac:
				w, h = float32(19*n+16), 23
			case Windows:
				w, h = float32(12*min(n, 5)), 24
			default:
				w, h = float32(15*n+12), 20
			}
			if s.box.W != w || s.box.H != h {
				t.Errorf("%v: a pager of %d is %v by %v, want %v by %v", p, n, s.box.W, s.box.H, w, h)
			}
		}
	}
}

// pagerPixel returns the color tester painted at x, y in DIPs, at scale.
func pagerPixel(tester *ui.Tester, scale, x, y float32) ui.Color {
	px := tester.Image().RGBAAt(int(x*scale), int(y*scale))
	return ui.RGB(px.R, px.G, px.B)
}

// TestPipsPagerMacDots checks UIPageControl's dots on macOS at 2x, light and dark, left to right and right to left: each 7 points across in the middle of its 19-point cell, the chosen one black (white in Dark Aqua) and the others at 45%, as UIKit for Mac draws them on macOS 27.2 (white at 45% reads 115 over black).
func TestPipsPagerMacDots(t *testing.T) {
	for _, dark := range []bool{false, true} {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			s := &pagerTestState{selected: 1, count: 5}
			var th ui.Theme
			tester := layoutTester(Mac, 400, 200, func(c *ui.Context) {
				s.view(c)
				th = *c.Theme()
			})
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.SetDirection(dir)
			tester.Frame()
			tester.Frame()
			ink := pick(dark, black, white)
			for i := range 5 {
				r, _ := tester.Find(fmt.Sprintf("Page %d", i+1))
				if r.W != 19 || r.H != 23 {
					t.Fatalf("dark %v %v: pip %d is %v by %v, want 19 by 23", dark, dir, i, r.W, r.H)
				}
				cx, cy := r.X+r.W/2, r.Y+r.H/2
				want := th.Background.Mix(ink, 0.45)
				if i == 1 {
					want = ink
				}
				// Half a point inside and outside the dot's edge, clear of the antialiased pixels on it.
				for _, d := range []float32{0, -3, 2.5} {
					if got := pagerPixel(tester, 2, cx+d, cy); !macNear(got, want, 3) {
						t.Errorf("dark %v %v: dot %d reads %v at %v from its middle, want %v", dark, dir, i, got, d, want)
					}
				}
				for _, d := range []float32{-4.5, 4} {
					if got := pagerPixel(tester, 2, cx+d, cy); !macNear(got, th.Background, 3) || !macNear(pagerPixel(tester, 2, cx, cy+d), th.Background, 3) {
						t.Errorf("dark %v %v: dot %d reaches %v from its middle (%v)", dark, dir, i, d, got)
					}
				}
			}
			if dir == ui.LeftToRight {
				if first, _ := tester.Find("Page 1"); first.X-s.box.X != 8 {
					t.Errorf("dark %v: the first cell is %v from the control's start, want 8", dark, first.X-s.box.X)
				}
			} else if first, _ := tester.Find("Page 1"); s.box.X+s.box.W-(first.X+first.W) != 8 {
				t.Errorf("dark %v: right to left, the first cell is %v from the control's right, want 8", dark, s.box.X+s.box.W-(first.X+first.W))
			}
		}
	}
}

// TestPipsPagerMacFocusRing checks that macOS's pager shows AppKit's focus ring around the control, round at its ends, when Tab brings the focus to it.
func TestPipsPagerMacFocusRing(t *testing.T) {
	s := &pagerTestState{selected: 1, count: 5}
	tester := layoutTester(Mac, 400, 200, s.view)
	tester.SetScale(2)
	tester.Frame()
	y := s.box.Y + s.box.H/2
	before := pagerPixel(tester, 2, s.box.X-2, y)
	if err := tester.Click("Before"); err != nil {
		t.Fatal(err)
	}
	carouselPress(tester, 0, ui.KeyTab)
	got := pagerPixel(tester, 2, s.box.X-2, y)
	if int(got.B)-int(got.R) < 40 || got == before {
		t.Errorf("2 points before the focused pager reads %v, want the focus ring", got)
	}
}

// TestPipsPagerFluentPips checks WinUI's pips on Windows at 2x, light and dark: the chosen pip's dot 6 DIPs across, another's 4, in ControlStrongFillColorDefault (#72000000, #8BFFFFFF) over the window; a pip under the pointer 6 DIPs across in TextFillColorSecondary (#9E000000, #C5FFFFFF); and the system focus visual on the chosen pip's edge once Tab brings the focus there.
func TestPipsPagerFluentPips(t *testing.T) {
	for _, dark := range []bool{false, true} {
		s := &pagerTestState{selected: 1, count: 5}
		var th ui.Theme
		tester := layoutTester(Windows, 400, 200, func(c *ui.Context) {
			s.view(c)
			th = *c.Theme()
		})
		tester.SetScale(2)
		tester.SetDark(dark)
		tester.Frame()
		tester.Frame()
		strong := th.Background.Mix(pick(dark, black, white), pick(dark, float32(0x72)/255, 0x8b/255.0))
		secondary := th.Background.Mix(pick(dark, black, white), pick(dark, float32(0x9e)/255, 0xc5/255.0))
		dot := func(label string, want ui.Color, d float32) {
			t.Helper()
			r, _ := tester.Find(label)
			if r.W != 12 || r.H != 24 {
				t.Fatalf("dark %v: %s is %v by %v, want 12 by 24", dark, label, r.W, r.H)
			}
			cx, cy := r.X+r.W/2, r.Y+r.H/2
			if got := pagerPixel(tester, 2, cx, cy); !macNear(got, want, 3) {
				t.Errorf("dark %v: %s's dot reads %v, want %v", dark, label, got, want)
			}
			// Half a DIP inside and outside the dot's edge, clear of the antialiased pixels on it.
			in, out := d/2-0.75, d/2+0.5
			if got := pagerPixel(tester, 2, cx+in, cy); !macNear(got, want, 3) {
				t.Errorf("dark %v: %s's dot is narrower than %v (%v at %v)", dark, label, d, got, in)
			}
			if got := pagerPixel(tester, 2, cx+out, cy); !macNear(got, th.Background, 3) {
				t.Errorf("dark %v: %s's dot is wider than %v (%v at %v)", dark, label, d, got, out)
			}
		}
		dot("Page 2", strong, 6)
		dot("Page 4", strong, 4)
		r, _ := tester.Find("Page 4")
		tester.Move(r.X+r.W/2, r.Y+r.H/2)
		tester.Frame()
		dot("Page 4", secondary, 6)
		tester.Move(300, 180)
		tester.Frame()
		if err := tester.Click("Before"); err != nil {
			t.Fatal(err)
		}
		carouselPress(tester, 0, ui.KeyTab)
		chosen, _ := tester.Find("Page 2")
		edge := pagerPixel(tester, 2, chosen.X+0.75, chosen.Y+chosen.H/2)
		if d := int(edge.G) - int(th.Background.G); (dark && d < 100) || (!dark && d > -100) {
			t.Errorf("dark %v: the chosen pip's edge reads %v, want the focus visual's outer stroke", dark, edge)
		}
	}
}

// TestPipsPagerFluentScrollsItsPips checks how WinUI's ScrollViewer scrolls the pips (ScrollToCenterOfViewport), as pipsFluentScroll gives it and as the pager shows it: 5 pips at most, the chosen one in the middle, held at the ends.
func TestPipsPagerFluentScrollsItsPips(t *testing.T) {
	for _, c := range []struct {
		n, sel int
		off    float32
	}{{5, 4, 0}, {9, 0, 0}, {9, 2, 0}, {9, 3, 12}, {9, 6, 48}, {9, 8, 48}, {3, 2, 0}} {
		if got := pipsFluentScroll(c.n, c.sel); got != c.off {
			t.Errorf("%d pips, the %dth chosen: scrolled %v, want %v", c.n, c.sel, got, c.off)
		}
	}
	for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
		s := &pagerTestState{selected: 6, count: 9}
		tester := layoutTester(Windows, 400, 200, s.view)
		tester.SetDirection(dir)
		tester.Frame()
		tester.Frame()
		chosen, _ := tester.Find("Page 7")
		if mid := chosen.X + chosen.W/2 - s.box.X; mid != 30 {
			t.Errorf("%v: the 7th of 9 pips is %v into the pager, want its middle, 30", dir, mid)
		}
		fifth, _ := tester.Find("Page 5")
		at := fifth.X - s.box.X
		if dir == ui.RightToLeft {
			at = s.box.X + s.box.W - (fifth.X + fifth.W)
		}
		if at != 0 {
			t.Errorf("%v: the 5th pip is %v from the pager's start, want the first shown", dir, at)
		}
	}
}

// TestPipsPagerAdwDots checks AdwCarouselIndicatorDots's geometry as snapshot_dots computes it: dots 15 pixels apart, the row centered on a width made odd, so that in its natural width the dots stand at 16.5, 31.5 and on from its left, and in a carousel 320 wide at 140.5 and on; and how a position between two pages shares the choice between their dots.
func TestPipsPagerAdwDots(t *testing.T) {
	for _, c := range []struct {
		w    float32
		n, s int
		want float32
	}{{87, 5, 0, 16.5}, {87, 5, 4, 76.5}, {147, 9, 6, 106.5}, {320, 4, 0, 140.5}, {320, 4, 3, 185.5}, {321, 4, 0, 141.5}} {
		if got := dotsAdwCenter(c.w, c.n, c.s); got != c.want {
			t.Errorf("dot %d of %d in %v pixels is at %v, want %v", c.s, c.n, c.w, got, c.want)
		}
	}
	for _, c := range []struct {
		pos  float32
		want [4]float32
	}{{1, [4]float32{0, 1, 0, 0}}, {1.25, [4]float32{0, 0.75, 0.25, 0}}, {0, [4]float32{1, 0, 0, 0}}, {3, [4]float32{0, 0, 0, 1}}} {
		for i := range 4 {
			if got := dotsAdwShare(4, i, c.pos); got != c.want[i] {
				t.Errorf("at %v, dot %d is %v chosen, want %v", c.pos, i, got, c.want[i])
			}
		}
	}
}

// TestPipsPagerAdwDotsDrawn checks GNOME's dots at 2x, light and dark, left to right and right to left: the chosen dot of radius 4 at 90% of the text's color, the others of radius 3 at 30%, as AdwCarouselIndicatorDots draws them, 10 pixels down the 20-pixel widget, the chosen one mirrored right to left.
func TestPipsPagerAdwDotsDrawn(t *testing.T) {
	for _, dark := range []bool{false, true} {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			s := &pagerTestState{selected: 1, count: 5}
			var th ui.Theme
			tester := layoutTester(Linux, 400, 200, func(c *ui.Context) {
				s.view(c)
				th = *c.Theme()
			})
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.SetDirection(dir)
			tester.Frame()
			tester.Frame()
			for i := range 5 {
				slot := i
				if dir == ui.RightToLeft {
					slot = 4 - i
				}
				cx, cy := s.box.X+16.5+15*float32(slot), s.box.Y+10
				share, radius := float32(0.3), float32(3)
				if i == 1 {
					share, radius = 0.9, 4
				}
				want := adwInkAlpha(&th, th.Background, share)
				if got := pagerPixel(tester, 2, cx, cy); !macNear(got, want, 3) {
					t.Errorf("dark %v %v: dot %d reads %v at its middle, want %v", dark, dir, i, got, want)
				}
				// Half a pixel inside and outside the dot's edge, clear of the antialiased pixels on it.
				if got := pagerPixel(tester, 2, cx+radius-0.75, cy); !macNear(got, want, 3) {
					t.Errorf("dark %v %v: dot %d is narrower than radius %v (%v)", dark, dir, i, radius, got)
				}
				if got := pagerPixel(tester, 2, cx+radius+0.5, cy); !macNear(got, th.Background, 3) {
					t.Errorf("dark %v %v: dot %d is wider than radius %v (%v)", dark, dir, i, radius, got)
				}
			}
		}
	}
}

// TestCarouselMovesWithKeysAndClicks checks the carousel on every platform, left to right and right to left: Tab stops at its view, named by the page shown; the arrow toward the end shows the next page and the one toward the start the page before, Page Down and Page Up the same, End the last and Home the first, without going round; Up and Down move only on Windows, as a horizontal FlipView takes them; a click on a page's dot shows it; and the next Tab leaves the carousel, its dots being no Tab stop.
func TestCarouselMovesWithKeysAndClicks(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			name := fmt.Sprintf("%v %v", p, dir)
			s := &carouselTestState{shown: 1}
			tester := layoutTester(p, 400, 320, s.view)
			tester.SetDirection(dir)
			tester.Frame()
			tester.Frame()
			if !tester.HasText("Slide 2") || tester.HasText("Slide 1") || tester.HasText("Slide 3") {
				t.Fatalf("%s: the carousel does not show the second page alone", name)
			}
			if err := tester.Click("Before"); err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			carouselPress(tester, 0, ui.KeyTab)
			if !tester.Focused("Page 2 of 4") {
				t.Fatalf("%s: Tab did not stop at the carousel's view named Page 2 of 4", name)
			}
			forth, back := ui.KeyRight, ui.KeyLeft
			if dir == ui.RightToLeft {
				forth, back = back, forth
			}
			steps := []struct {
				key  ui.Key
				want int
			}{{forth, 2}, {forth, 3}, {forth, 3}, {back, 2}, {ui.KeyPageUp, 1}, {ui.KeyPageDown, 2}, {ui.KeyHome, 0}, {back, 0}, {ui.KeyEnd, 3}}
			for _, step := range steps {
				carouselPress(tester, 0, step.key)
				if s.shown != step.want || !tester.HasText(carouselTestSlides[step.want]) {
					t.Fatalf("%s: %v showed page %d, want %d", name, step.key, s.shown, step.want)
				}
			}
			carouselPress(tester, 0, ui.KeyUp)
			want := 3
			if p == Windows {
				want = 2
			}
			if s.shown != want {
				t.Errorf("%s: Up showed page %d, want %d", name, s.shown, want)
			}
			if !tester.Focused(fmt.Sprintf("Page %d of 4", want+1)) {
				t.Errorf("%s: the view is not named after page %d", name, want+1)
			}
			if err := tester.Click("Page 1"); err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			tester.Frame()
			tester.Frame()
			if s.shown != 0 || !tester.HasText("Slide 1") {
				t.Errorf("%s: a click on the first dot showed page %d", name, s.shown)
			}
			if err := tester.Click("Before"); err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			carouselPress(tester, 0, ui.KeyTab)
			carouselPress(tester, 0, ui.KeyTab)
			if !tester.Focused("After") {
				t.Errorf("%s: Tab from the carousel did not go on to After", name)
			}
		}
	}
}

// TestCarouselBuildsTheShownPages checks that the carousel builds the page it shows alone at rest, and while it moves between two pages those two, on every platform.
func TestCarouselBuildsTheShownPages(t *testing.T) {
	for _, p := range Platforms {
		s := &carouselTestState{shown: 1, built: map[int]int{}}
		tester := ui.NewTester(func(c *ui.Context) { Preview(c, p, func() { s.view(c) }) }, 400, 320)
		tester.SetPreferences(ui.Preferences{TextScale: 1})
		tester.Frame()
		tester.Frame()
		clear(s.built)
		tester.Frame()
		if len(s.built) != 1 || s.built[1] == 0 {
			t.Errorf("%v: at rest the carousel built pages %v, want the second alone", p, s.built)
		}
		if err := tester.Click("Before"); err != nil {
			t.Fatal(err)
		}
		carouselPress(tester, 0, ui.KeyTab)
		tester.Key(0, ui.KeyRight)
		tester.Frame()
		tester.Advance(40 * 1e6)
		clear(s.built)
		tester.Frame()
		if len(s.built) != 2 || s.built[1] == 0 || s.built[2] == 0 {
			t.Errorf("%v: moving to the third page the carousel built pages %v, want the second and the third", p, s.built)
		}
		if !tester.HasText("Slide 2") || !tester.HasText("Slide 3") {
			t.Errorf("%v: moving to the third page the carousel does not show the second and the third", p)
		}
		tester.Advance(2e9)
		tester.Frame()
		tester.Frame()
		clear(s.built)
		tester.Frame()
		if len(s.built) != 1 || s.built[2] == 0 || tester.HasText("Slide 2") {
			t.Errorf("%v: come to rest the carousel built pages %v, want the third alone", p, s.built)
		}
	}
}
