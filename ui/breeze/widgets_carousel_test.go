package breeze

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

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

// carouselShotView lays out the screenshot test's controls on p: a PipsPager of 5 pages with the 2nd chosen, one of 9 with the 7th chosen, and a Carousel 320 wide of 4 pages with the 2nd shown, at carouselShotPlaces, all three disabled when disabled is.
func carouselShotView(p Platform, disabled bool) func(c *ui.Context) {
	five, nine, shown := 1, 6, 1
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Box(c).Fill().Children(func() {
				at := func(k int, build func()) {
					ui.Box(c).Absolute().Start(carouselShotPlaces[k].X).Top(carouselShotPlaces[k].Y).Children(build)
				}
				at(0, func() { PipsPager(c, &five, 5).Disabled(disabled) })
				at(1, func() { PipsPager(c, &nine, 9).Disabled(disabled) })
				at(2, func() { Carousel(c, &shown, 4, carouselTestPage(c, nil)).Width(320).Disabled(disabled) })
			})
		})
	}
}

// TestScreenshotsCarousel writes the pagers and the carousel of every platform, light and dark, and right to left in the light theme, at 2x to the directory BREEZE_SCREENSHOTS names, and on Windows the same three again with FlipView's buttons shown by the mouse over the view, and the pagers and the carousel disabled, light and dark. Without it the test only builds the frames.
func TestScreenshotsCarousel(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	type shot struct {
		dark bool
		rtl  bool
	}
	for _, p := range Platforms {
		for _, s := range []shot{{false, false}, {true, false}, {false, true}} {
			tester := ui.NewTester(carouselShotView(p, false), 352, 312)
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
			if p == Windows {
				// FlipView's buttons, which a move of the mouse over the view shows, as tools/winui-reference's winui-carousel-buttons and winui-carousel-buttons-rtl scenes show WinUI's.
				tester.Move(carouselShotPlaces[2].X+160, carouselShotPlaces[2].Y+carouselTestPageHeight/2)
				tester.Frame()
				tester.Frame()
				if _, ok := tester.Find("Next"); !ok {
					t.Fatalf("%v: the mouse over the carousel does not show its buttons", p)
				}
				writeShot(t, dir, fmt.Sprintf("breeze-carousel-buttons-windows-%s.png", mode), tester)
			}
		}
	}
	// Disabled on Windows, light and dark, as tools/winui-reference's winui-carousel-disabled scene shows WinUI's.
	for _, dark := range []bool{false, true} {
		tester := ui.NewTester(carouselShotView(Windows, true), 352, 312)
		tester.SetScale(2)
		tester.SetDark(dark)
		tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		tester.Frame()
		tester.Frame()
		writeShot(t, dir, fmt.Sprintf("breeze-carousel-disabled-windows-%s.png", map[bool]string{false: "light", true: "dark"}[dark]), tester)
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

// TestPipsPagerSizes checks the pager's size on every platform, for 0, 1, 5 and 9 pages: on macOS UIPageControl's 23 points high and 19 a dot with 8 at either end (35 for 1, 111 for 5 and 187 for 9, as UIKit for Mac lays it out), and 0 by 0 without pages, as sizeThatFits gives it; on Windows 24 high and 12 a pip, at most 5 (MaxVisiblePips), and 0 by 0 without pages, as the empty ItemsRepeater is; on GNOME 20 high and 15 a dot with 6 at either end (AdwCarouselIndicatorDots's measure), 12 by 20 without pages.
func TestPipsPagerSizes(t *testing.T) {
	for _, p := range Platforms {
		for _, n := range []int{0, 1, 5, 9} {
			s := &pagerTestState{count: n}
			tester := layoutTester(p, 400, 200, s.view)
			tester.Frame()
			var w, h float32
			switch p {
			case Mac:
				w, h = float32(19*n+16), 23
				if n == 0 {
					w, h = 0, 0
				}
			case Windows:
				w, h = float32(12*min(n, 5)), 24
				if n == 0 {
					h = 0
				}
			default:
				w, h = float32(15*n+12), 20
			}
			if s.box.W != w || s.box.H != h {
				t.Errorf("%v: a pager of %d is %v by %v, want %v by %v", p, n, s.box.W, s.box.H, w, h)
			}
		}
	}
}

// inkArea returns how much ink of color ink over bg tester painted in r, in DIPs, at scale, in square DIPs: each pixel's share of the way from bg to ink, by its green, summed.
func inkArea(tester *ui.Tester, scale float32, r ui.Rect, bg, ink ui.Color) float32 {
	img := tester.Image()
	sum := float32(0)
	for y := int(r.Y * scale); y < int((r.Y+r.H)*scale); y++ {
		for x := int(r.X * scale); x < int((r.X+r.W)*scale); x++ {
			share := (float32(bg.G) - float32(img.RGBAAt(x, y).G)) / (float32(bg.G) - float32(ink.G))
			sum += max(0, min(share, 1))
		}
	}
	return sum / (scale * scale)
}

// inkDiameter returns the diameter of a circle as large as the ink inkArea reads in r.
func inkDiameter(tester *ui.Tester, scale float32, r ui.Rect, bg, ink ui.Color) float32 {
	return 2 * float32(math.Sqrt(float64(inkArea(tester, scale, r, bg, ink))/math.Pi))
}

// pagerPixel returns the color tester painted at x, y in DIPs, at scale.
func pagerPixel(tester *ui.Tester, scale, x, y float32) ui.Color {
	px := tester.Image().RGBAAt(int(x*scale), int(y*scale))
	return ui.RGB(px.R, px.G, px.B)
}

// TestPipsPagerMacDots checks UIPageControl's dots on macOS at 2x, light and dark, left to right and right to left: each 6 points across in the middle of its 19-point cell, the chosen one black (white in Dark Aqua) and the others at 45%, as UIKit for Mac draws them on macOS 27.2 (white at 45% reads 131 over the dark window's 30).
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
				for _, d := range []float32{0, -2.5, 2} {
					if got := pagerPixel(tester, 2, cx+d, cy); !macNear(got, want, 3) {
						t.Errorf("dark %v %v: dot %d reads %v at %v from its middle, want %v", dark, dir, i, got, d, want)
					}
				}
				for _, d := range []float32{-3.5, 3.5} {
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

// TestPipsPagerFluentPips checks WinUI's pips on Windows at 2x, light and dark: the chosen pip's dot 5.625 DIPs across, another's 3.75 (the EA3B glyph's circle at 6 and 4), in ControlStrongFillColorDefault (#72000000, #8BFFFFFF) over the window; a pip under the pointer 5.625 DIPs across and a pressed one 3.75 in TextFillColorSecondary (#9E000000, #C5FFFFFF), as tools/winui-reference's winui-carousel-states scene shows them; and the system focus visual on the chosen pip's edge once Tab brings the focus there.
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
			// The ink is as large as a circle d across, its center in the pip's middle: the glyph's circle fills 15/16 of its font size, which winui-carousel measured as 3.74 and 5.63 DIPs across in the light theme.
			if got := inkDiameter(tester, 2, r, th.Background, want); got < d-0.08 || got > d+0.08 {
				t.Errorf("dark %v: %s's dot is %v DIPs across by its ink, want %v", dark, label, got, d)
			}
		}
		dot("Page 2", strong, 5.625)
		dot("Page 4", strong, 3.75)
		r, _ := tester.Find("Page 4")
		tester.Move(r.X+r.W/2, r.Y+r.H/2)
		tester.Frame()
		dot("Page 4", secondary, 5.625)
		tester.Press(r.X+r.W/2, r.Y+r.H/2)
		tester.Frame()
		dot("Page 4", secondary, 3.75)
		tester.Move(300, 180)
		tester.Release(300, 180)
		tester.Frame()
		if s.selected != 1 {
			t.Fatalf("dark %v: a press let go away from the pip chose page %d", dark, s.selected+1)
		}
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

// TestPipsPagerFluentDisabled checks that a disabled pager's pips and a disabled carousel's on Windows read ControlStrongFillColorDisabled (#51000000, #3FFFFFFF) over the window, chosen or not, at their sizes, as tools/winui-reference's winui-carousel-disabled scene measured WinUI's: 166 in the light theme and 87 in the dark one, where MyGo's half opacity alone would make them 189 and 93.
func TestPipsPagerFluentDisabled(t *testing.T) {
	for _, dark := range []bool{false, true} {
		for _, inCarousel := range []bool{false, true} {
			var th ui.Theme
			pager, carousel := &pagerTestState{selected: 1, count: 5, disabled: true}, &carouselTestState{shown: 1, disabled: true}
			tester := layoutTester(Windows, 400, 320, func(c *ui.Context) {
				if inCarousel {
					carousel.view(c)
				} else {
					pager.view(c)
				}
				th = *c.Theme()
			})
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.Frame()
			tester.Frame()
			want := th.Background.Mix(pick(dark, black, white), pick(dark, float32(0x51)/255, 0x3f/255.0))
			for _, c := range []struct {
				label string
				d     float32
			}{{"Page 2", 5.625}, {"Page 4", 3.75}} {
				r, ok := tester.Find(c.label)
				if !ok {
					t.Fatalf("dark %v, in a carousel %v: no %s", dark, inCarousel, c.label)
				}
				if got := pagerPixel(tester, 2, r.X+r.W/2, r.Y+r.H/2); !macNear(got, want, 2) {
					t.Errorf("dark %v, in a carousel %v: %s's dot reads %v, want %v", dark, inCarousel, c.label, got, want)
				}
				if got := inkDiameter(tester, 2, r, th.Background, want); got < c.d-0.08 || got > c.d+0.08 {
					t.Errorf("dark %v, in a carousel %v: %s's dot is %v DIPs across by its ink, want %v", dark, inCarousel, c.label, got, c.d)
				}
			}
		}
	}
}

// TestPipsPagerFluentScrollsItsPips checks how WinUI's ScrollViewer scrolls the pips (ScrollToCenterOfViewport), as pipsFluentScroll gives it and as the pager shows it: 5 pips at most, the chosen one in the middle, held at the ends; tools/winui-reference's winui-carousel scene measured the pager of 9 with the 7th chosen scrolled by 48 DIPs.
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

// carouselMotionTester returns a tester of the test carousel on p without ReduceMotion, so that it animates, laid out dir.
func carouselMotionTester(p Platform, s *carouselTestState, dir ui.Direction) *ui.Tester {
	tester := ui.NewTester(func(c *ui.Context) { Preview(c, p, func() { s.view(c) }) }, 400, 320)
	tester.SetPreferences(ui.Preferences{TextScale: 1})
	tester.SetDirection(dir)
	tester.Frame()
	tester.Frame()
	return tester
}

// TestCarouselSlidesTowardTheNextPage checks on every platform that the next page comes in from the end side as the carousel moves to it, from the right left to right and from the left right to left, the page shown going out the other way.
func TestCarouselSlidesTowardTheNextPage(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			s := &carouselTestState{shown: 1}
			tester := carouselMotionTester(p, s, dir)
			if err := tester.Click("Before"); err != nil {
				t.Fatal(err)
			}
			carouselPress(tester, 0, ui.KeyTab)
			forth := ui.KeyRight
			if dir == ui.RightToLeft {
				forth = ui.KeyLeft
			}
			tester.Key(0, forth)
			tester.Frame()
			tester.Advance(60 * 1e6)
			tester.Frame()
			out, ok1 := tester.Find("Slide 2")
			in, ok2 := tester.Find("Slide 3")
			if !ok1 || !ok2 {
				t.Fatalf("%v %v: moving, the carousel does not show both pages", p, dir)
			}
			if (dir == ui.LeftToRight) != (in.X > out.X) {
				t.Errorf("%v %v: the next page is at %v, the page shown at %v", p, dir, in.X, out.X)
			}
			view := s.box
			if mid := out.X + out.W/2; mid > view.X && mid < view.X+view.W && (mid-view.X-view.W/2)*(in.X-out.X) > 0 {
				t.Errorf("%v %v: the page shown went toward the next page's side", p, dir)
			}
		}
	}
}

// TestCarouselMotion checks how each platform moves between pages: GNOME's spring takes a page over 0.3 seconds or so and passes the pages between on a longer move; Windows animates a move to the adjacent page but jumps to another (FlipView::OnSelectedIndexChanged); macOS slides over 0.25 seconds (inferred).
func TestCarouselMotion(t *testing.T) {
	for _, p := range Platforms {
		s := &carouselTestState{shown: 0}
		tester := carouselMotionTester(p, s, ui.LeftToRight)
		if err := tester.Click("Before"); err != nil {
			t.Fatal(err)
		}
		carouselPress(tester, 0, ui.KeyTab)
		tester.Key(0, ui.KeyEnd)
		tester.Frame()
		tester.Advance(30 * 1e6)
		tester.Frame()
		jumped := tester.HasText("Slide 4") && !tester.HasText("Slide 1") && !tester.HasText("Slide 2")
		if jumped != (p == Windows) {
			t.Errorf("%v: End from the first page jumped %v, want %v", p, jumped, p == Windows)
		}
		tester.Advance(time1s)
		tester.Frame()
		tester.Frame()
		if !tester.HasText("Slide 4") || tester.HasText("Slide 3") {
			t.Errorf("%v: a second after End the carousel does not rest on the last page", p)
		}
		tester.Key(0, ui.KeyLeft)
		tester.Frame()
		tester.Advance(30 * 1e6)
		tester.Frame()
		if !tester.HasText("Slide 3") || !tester.HasText("Slide 4") {
			t.Errorf("%v: 30 ms after Left the carousel is not between the last two pages", p)
		}
		tester.Advance(240 * 1e6)
		tester.Frame()
		tester.Frame()
		moving := tester.HasText("Slide 4")
		if moving != (p == Linux) {
			t.Errorf("%v: 270 ms after Left the carousel still moves %v, want %v", p, moving, p == Linux)
		}
	}
}

// time1s is a second.
const time1s = 1e9

// carouselMiddle returns the middle of the test carousel's view, which the pages fill.
func carouselMiddle(tester *ui.Tester) (x, y float32) {
	r, _ := tester.Find("Slide 2")
	if r.W == 0 {
		for _, s := range carouselTestSlides {
			if b, ok := tester.Find(s); ok {
				r = b
				break
			}
		}
	}
	return r.X + r.W/2, r.Y + r.H/2
}

// TestCarouselSwipesWithATouchpad checks the touchpad's swipe on GNOME and macOS, left to right and right to left: the pages follow the fingers, on GNOME 400 pixels to a page (TOUCHPAD_BASE_DISTANCE_H), on macOS the view's width (inferred); a slow swipe ending less than half a page on goes back, past half a page on; a quick one goes on to the next page (VELOCITY_THRESHOLD_TOUCHPAD); and a scroll less than 1.2 times as sideways as up or down leaves the pages, as SCROLL_ORIENTATION_FACTOR rejects it.
func TestCarouselSwipesWithATouchpad(t *testing.T) {
	for _, p := range []Platform{Linux, Mac} {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			name := fmt.Sprintf("%v %v", p, dir)
			page := float32(400)
			if p == Mac {
				page = 320
			}
			sign := float32(1)
			if dir == ui.RightToLeft {
				sign = -1
			}
			s := &carouselTestState{shown: 1}
			tester := carouselMotionTester(p, s, dir)
			x, y := carouselMiddle(tester)
			// Slowly 0.45 of a page on: it goes back.
			for range 9 {
				tester.Scroll(x, y, sign*page/20, 0)
				tester.Advance(100 * 1e6)
				tester.Frame()
			}
			if !tester.HasText("Slide 3") {
				t.Errorf("%s: 0.45 of a page into a swipe the next page does not show", name)
			}
			tester.Advance(time1s)
			tester.Frame()
			tester.Frame()
			if s.shown != 1 || tester.HasText("Slide 3") {
				t.Errorf("%s: a slow swipe of 0.45 of a page showed page %d, want it back on 1", name, s.shown)
			}
			// Slowly two thirds of a page on: it goes on.
			for range 8 {
				tester.Scroll(x, y, sign*page/12, 0)
				tester.Advance(100 * 1e6)
				tester.Frame()
			}
			tester.Advance(time1s)
			tester.Frame()
			tester.Frame()
			if s.shown != 2 {
				t.Errorf("%s: a slow swipe of two thirds of a page showed page %d, want 2", name, s.shown)
			}
			// Quickly a tenth of a page back: it goes on to the page before.
			for range 4 {
				tester.Scroll(x, y, -sign*page/40, 0)
				tester.Advance(10 * 1e6)
			}
			tester.Frame()
			tester.Advance(time1s)
			tester.Frame()
			tester.Frame()
			if s.shown != 1 {
				t.Errorf("%s: a quick swipe of a tenth of a page back showed page %d, want 1", name, s.shown)
			}
			// 300 sideways and 260 down is not 1.2 times as sideways as down.
			tester.Scroll(x, y, sign*300, 260)
			tester.Advance(time1s)
			tester.Frame()
			tester.Frame()
			if s.shown != 1 {
				t.Errorf("%s: a scroll less than 1.2 times as sideways as down moved to page %d", name, s.shown)
			}
		}
	}
}

// TestCarouselWheel checks a mouse's wheel: on Windows a notch flips a page, down or to the right the next, but another the same way within 200 ms flips none, while one the other way or after the pause does (FlipView::OnPointerWheelChanged); on GNOME a notch moves a page and another within 150 ms none (AdwCarousel's scroll_cb); on macOS the wheel moves no page (inferred, as NSPageController follows only a touchpad).
func TestCarouselWheel(t *testing.T) {
	for _, p := range Platforms {
		s := &carouselTestState{shown: 1}
		tester := carouselMotionTester(p, s, ui.LeftToRight)
		x, y := carouselMiddle(tester)
		notch := func() {
			tester.Wheel(x, y, 0, 40)
			tester.Frame()
		}
		notch()
		want := 2
		if p == Mac {
			want = 1
		}
		if s.shown != want {
			t.Fatalf("%v: a notch down showed page %d, want %d", p, s.shown, want)
		}
		if p == Mac {
			continue
		}
		tester.Advance(100 * 1e6)
		notch()
		if s.shown != 2 {
			t.Errorf("%v: a notch 100 ms after the last showed page %d, want none flipped", p, s.shown)
		}
		tester.Advance(250 * 1e6)
		notch()
		if s.shown != 3 {
			t.Errorf("%v: a notch 250 ms after the last showed page %d, want 3", p, s.shown)
		}
		if p == Windows {
			tester.Wheel(x, y, 0, -40)
			tester.Frame()
			if s.shown != 2 {
				t.Errorf("%v: a notch up right after showed page %d, want 2, as it goes the other way", p, s.shown)
			}
		}
	}
}

// TestCarouselDragsOnGNOME checks that on GNOME the pointer drags the pages, as AdwCarousel lets a mouse, once it went 16 pixels, and a slow drag past half a page shows the next one, left to right and right to left; on Windows and macOS a drag moves none.
func TestCarouselDragsOnGNOME(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			s := &carouselTestState{shown: 1}
			tester := carouselMotionTester(p, s, dir)
			x, y := carouselMiddle(tester)
			step := float32(-20)
			if dir == ui.RightToLeft {
				step = 20
			}
			tester.Press(x, y)
			tester.Frame()
			moved := false
			for k := 1; k <= 10; k++ {
				tester.Advance(100 * 1e6)
				tester.Move(x+step*float32(k), y)
				tester.Frame()
				moved = moved || tester.HasText("Slide 3")
			}
			tester.Release(x+step*10, y)
			tester.Frame()
			tester.Advance(time1s)
			tester.Frame()
			tester.Frame()
			want := 1
			if p == Linux {
				want = 2
			}
			if s.shown != want || moved != (p == Linux) {
				t.Errorf("%v %v: a slow drag of 200 DIPs showed page %d (the pages moved %v), want %d", p, dir, s.shown, moved, want)
			}
		}
	}
}

// TestCarouselFlipButtons checks FlipView's previous and next buttons on Windows, left to right and right to left: hidden until the mouse moves over the view; then 16 by 38 DIPs, 1 DIP in from the view's sides and centered down it, the previous one at the start; a click on Next shows the next page; hidden again 3 seconds after the mouse last moved over the view, not over a button; on the first page no previous button and on the last no next one.
func TestCarouselFlipButtons(t *testing.T) {
	for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
		s := &carouselTestState{shown: 1}
		tester := carouselMotionTester(Windows, s, dir)
		if _, ok := tester.Find("Next"); ok {
			t.Fatalf("%v: the buttons show before the mouse moved over the view", dir)
		}
		x, y := carouselMiddle(tester)
		tester.Move(x, y)
		tester.Frame()
		prev, ok1 := tester.Find("Previous")
		next, ok2 := tester.Find("Next")
		if !ok1 || !ok2 {
			t.Fatalf("%v: the mouse over the view does not show both buttons", dir)
		}
		view := s.box
		view.H -= 34 // the PipsPager below, 10 and 24
		for _, b := range []ui.Rect{prev, next} {
			if b.W != 16 || b.H != 38 || b.Y+b.H/2 != view.Y+view.H/2 {
				t.Errorf("%v: a button is %v, want 16 by 38 centered down the view %v", dir, b, view)
			}
		}
		start, end := prev, next
		if dir == ui.RightToLeft {
			start, end = next, prev
		}
		if start.X != view.X+1 || end.X+end.W != view.X+view.W-1 {
			t.Errorf("%v: the buttons are at %v and %v in the view %v, want 1 in from its sides, the previous one at the start", dir, prev, next, view)
		}
		if err := tester.Click("Next"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		if s.shown != 2 {
			t.Errorf("%v: a click on Next showed page %d", dir, s.shown)
		}
		// Over a button they stay, as OnPointerEnteredNavigationButtons keeps them; the mouse moves back over the page.
		tester.Move(x, y)
		tester.Frame()
		tester.Advance(2900 * 1e6)
		tester.Frame()
		if _, ok := tester.Find("Next"); !ok {
			t.Errorf("%v: the buttons are gone before 3 seconds", dir)
		}
		tester.Advance(200 * 1e6)
		tester.Frame()
		if _, ok := tester.Find("Next"); ok {
			t.Errorf("%v: the buttons stay 3 seconds after the mouse last moved", dir)
		}
		s.shown = 0
		tester.Move(x+1, y)
		tester.Frame()
		tester.Frame()
		if _, ok := tester.Find("Previous"); ok {
			t.Errorf("%v: the first page shows a previous button", dir)
		}
		s.shown = 3
		tester.Move(x, y)
		tester.Frame()
		tester.Frame()
		if _, ok := tester.Find("Next"); ok {
			t.Errorf("%v: the last page shows a next button", dir)
		}
	}
}

// TestCarouselFlipButtonsDrawn checks FlipView's buttons at 2x, light and dark, left to right and right to left: the in-app acrylic's fallback color (#f9f9f9, #2c2c2c) and a caret 6 DIPs wide and 8 high in ControlStrongFillColorDefault in the middle, the previous button's pointing to the start, mirrored right to left as its glyph is (MirroredWhenRightToLeft), as tools/winui-reference's winui-carousel-buttons and winui-carousel-buttons-rtl scenes show them; and the caret under the pointer and pressed.
func TestCarouselFlipButtonsDrawn(t *testing.T) {
	for _, dark := range []bool{false, true} {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			name := fmt.Sprintf("dark %v %v", dark, dir)
			s := &carouselTestState{shown: 1}
			tester := ui.NewTester(func(c *ui.Context) { Preview(c, Windows, func() { s.view(c) }) }, 400, 320)
			tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.SetDirection(dir)
			tester.Frame()
			x, y := carouselMiddle(tester)
			tester.Move(x, y)
			tester.Frame()
			face := pick(dark, ui.Hex("#f9f9f9"), ui.Hex("#2c2c2c"))
			strong := face.Mix(pick(dark, black, white), pick(dark, float32(0x72)/255, 0x8b/255.0))
			prev, _ := tester.Find("Previous")
			cx, cy := prev.X+prev.W/2, prev.Y+prev.H/2
			if got := pagerPixel(tester, 2, prev.X+2, prev.Y+4); !macNear(got, face, 2) {
				t.Errorf("%s: the button's face reads %v, want %v", name, got, face)
			}
			// The caret's flat side, toward the end, and its tip, toward the start: right and left of the middle left to right, mirrored right to left.
			flat := float32(2.25)
			if dir == ui.RightToLeft {
				flat = -flat
			}
			if got := pagerPixel(tester, 2, cx+flat, cy); !macNear(got, strong, 3) {
				t.Errorf("%s: the caret reads %v near its flat side, want %v", name, got, strong)
			}
			if got := pagerPixel(tester, 2, cx+flat, cy-2.75); !macNear(got, strong, 3) {
				t.Errorf("%s: the caret's flat side is not 8 DIPs high (%v)", name, got)
			}
			if got := pagerPixel(tester, 2, cx-flat, cy-2); !macNear(got, face, 3) {
				t.Errorf("%s: the caret does not point toward the start (%v)", name, got)
			}
			if got := pagerPixel(tester, 2, cx+flat*3.5/2.25, cy); !macNear(got, face, 3) {
				t.Errorf("%s: the caret is wider than 6 DIPs (%v)", name, got)
			}
			// CaretLeftSolid8's outline, rounded at its corners and its tip, covers 32.2 square DIPs, where a sharp triangle 6 by 8 would cover 24; winui-carousel-buttons measured its ink as 32.3.
			if got := inkArea(tester, 2, ui.Rect{X: cx - 4, Y: cy - 5, W: 8, H: 10}, face, strong); got < 31.2 || got > 33.2 {
				t.Errorf("%s: the caret covers %v square DIPs, want 32.2", name, got)
			}
			// Its rounded tip, blunt where a sharp one comes to a point: on the middle row, the pixel of the first half DIP of the 6-DIP box reads mostly ink (WinUI's 144 of 138 to 249 in the light theme, a sharp tip's 212).
			tip := cx - 2.75
			if dir == ui.RightToLeft {
				tip = cx + 2.75
			}
			if got := pagerPixel(tester, 2, tip, cy); int(got.G) > (int(face.G)+int(strong.G))/2 == !dark {
				t.Errorf("%s: the caret's tip reads %v in its first half DIP, want it blunt", name, got)
			}
			// Under the pointer the caret is TextFillColorSecondary at its size, and pressed shrunk to 0.875 about its em's middle, as winui-carousel-states shows them: 95 in the light theme and 207 in the dark one, covering 32.2 and 24.7 square DIPs (WinUI's pressed glyph, hinted to whole pixels at 7 DIPs, measured 23.9 and 24.3).
			secondary := face.Mix(pick(dark, black, white), pick(dark, float32(0x9e)/255, 0xc5/255.0))
			box := ui.Rect{X: cx - 4, Y: cy - 5, W: 8, H: 10}
			tester.Move(cx, cy)
			tester.Frame()
			if got := pagerPixel(tester, 2, cx+flat, cy); !macNear(got, secondary, 3) {
				t.Errorf("%s: the caret under the pointer reads %v, want %v", name, got, secondary)
			}
			if got := inkArea(tester, 2, box, face, secondary); got < 31.2 || got > 33.2 {
				t.Errorf("%s: the caret under the pointer covers %v square DIPs, want 32.2", name, got)
			}
			// halves returns the ink left of the middle over that right of it, and the ink above the middle less that below it.
			halves := func() (float32, float32) {
				ink := func(x, y, w, h float32) float32 {
					return inkArea(tester, 2, ui.Rect{X: x, Y: y, W: w, H: h}, face, secondary)
				}
				return ink(cx-4, cy-5, 4, 10) / ink(cx, cy-5, 4, 10), ink(cx-4, cy-5, 8, 5) - ink(cx-4, cy, 8, 5)
			}
			across, _ := halves()
			tester.Press(cx, cy)
			tester.Frame()
			if got := pagerPixel(tester, 2, cx+flat, cy); !macNear(got, secondary, 3) {
				t.Errorf("%s: the pressed caret reads %v, want %v", name, got, secondary)
			}
			if got := inkArea(tester, 2, box, face, secondary); got < 23.7 || got > 25.7 {
				t.Errorf("%s: the pressed caret covers %v square DIPs, want 24.7", name, got)
			}
			// Shrunk about the middle, the ink keeps its share on either side of it and is as much above as below.
			if got, down := halves(); got < across-0.05 || got > across+0.05 || down < -0.3 || down > 0.3 {
				t.Errorf("%s: the pressed caret is not shrunk about its middle: %v of its ink left of it for every 1 right, where unpressed %v, and %v more above than below", name, got, across, down)
			}
			tester.Move(x, y)
			tester.Release(x, y)
			tester.Frame()
			if s.shown != 1 {
				t.Errorf("%s: a press let go away from the button showed page %d", name, s.shown+1)
			}
		}
	}
}

// TestCarouselLayout checks where each platform puts the dots below the pages, from the view's bottom, without reading the host font's widths: WinUI's PipsPager centered 10 DIPs below the FlipView; AdwCarouselIndicatorDots 6 below with a margin of 6 under it, as wide as the carousel; UIPageControl centered right below (inferred).
func TestCarouselLayout(t *testing.T) {
	for _, p := range Platforms {
		s := &carouselTestState{shown: 1}
		tester := layoutTester(p, 400, 320, s.view)
		tester.Frame()
		first, _ := tester.Find("Page 1")
		last, _ := tester.Find("Page 4")
		pageBottom := s.box.Y + carouselTestPageHeight
		var gap, below float32
		switch p {
		case Windows:
			gap, below = 10, 0
		case Linux:
			gap, below = 6, 6
		}
		if first.Y-pageBottom != gap || s.box.Y+s.box.H-(first.Y+first.H) != below {
			t.Errorf("%v: the dots are %v below the pages and %v above the carousel's bottom, want %v and %v", p, first.Y-pageBottom, s.box.Y+s.box.H-(first.Y+first.H), gap, below)
		}
		mid := (first.X + last.X + last.W) / 2
		want := s.box.X + s.box.W/2
		if p == Linux {
			// The dots' widget is 308 wide inside its margins of 6; snapshot_dots centers the dots on it made a pixel narrower, 307, and they stand half a spacing toward its end: their middle, and that of the 15-pixel cells around them, is 6 + 153.5 + 3.5 = 163 pixels into the carousel.
			want = s.box.X + 163
		}
		if mid != want {
			t.Errorf("%v: the dots' middle is at %v, want %v", p, mid, want)
		}
	}
}

// TestCarouselFocusRing checks the focus ring around the carousel's view once Tab brings the focus there: WinUI's focus visual on the view's edge, libadwaita's focus color 2 pixels inside it (inferred), AppKit's ring 3 points outside it (inferred).
func TestCarouselFocusRing(t *testing.T) {
	for _, p := range Platforms {
		s := &carouselTestState{shown: 1}
		var th ui.Theme
		tester := layoutTester(p, 400, 320, func(c *ui.Context) {
			s.view(c)
			th = *c.Theme()
		})
		tester.SetScale(2)
		tester.Frame()
		view := s.box
		y := view.Y + carouselTestPageHeight/2
		x := view.X + 0.75
		if p == Mac {
			x = view.X - 2
		}
		before := pagerPixel(tester, 2, x, y)
		if err := tester.Click("Before"); err != nil {
			t.Fatal(err)
		}
		carouselPress(tester, 0, ui.KeyTab)
		got := pagerPixel(tester, 2, x, y)
		switch p {
		case Linux:
			want := th.Background.Mix(ui.RGB(th.Focus.R, th.Focus.G, th.Focus.B), float32(th.Focus.A)/255)
			if !macNear(got, want, 3) {
				t.Errorf("%v: the view's edge reads %v focused, want %v", p, got, want)
			}
		default:
			if got == before {
				t.Errorf("%v: the view's edge reads %v focused as unfocused", p, got)
			}
		}
	}
}

// TestCarouselDisabled checks on every platform that a disabled carousel takes neither keys, clicks on its dots, nor the wheel.
func TestCarouselDisabled(t *testing.T) {
	for _, p := range Platforms {
		s := &carouselTestState{shown: 1, disabled: true}
		tester := layoutTester(p, 400, 320, s.view)
		tester.Frame()
		if err := tester.Click("Page 4"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		x, y := carouselMiddle(tester)
		tester.Wheel(x, y, 0, 40)
		tester.Frame()
		if err := tester.Click("Before"); err != nil {
			t.Fatal(err)
		}
		carouselPress(tester, 0, ui.KeyTab)
		carouselPress(tester, 0, ui.KeyRight)
		if s.shown != 1 {
			t.Errorf("%v: a disabled carousel showed page %d", p, s.shown)
		}
	}
}

// TestPagersWithoutPagesOrOutOfRange checks on every platform that a pager or a carousel of no pages is no Tab stop, as a PipsPager without pips and a FlipView or an AdwCarousel without pages have nothing to focus, and builds no page; that a page chosen before the first or past the last is held to the pages, as PipsPager::OnSelectedPageIndexChanged holds SelectedPageIndex; and that a carousel whose pages become fewer than the page shown shows the last.
func TestPagersWithoutPagesOrOutOfRange(t *testing.T) {
	for _, p := range Platforms {
		pager, shown, count := 3, 2, 0
		built := map[int]int{}
		tester := layoutTester(p, 400, 320, func(c *ui.Context) {
			Use(c)
			ui.Column(c).Padding(16).Gap(12).AlignItems(ui.Start).Children(func() {
				Button(c, "Before")
				PipsPager(c, &pager, count)
				Carousel(c, &shown, count, carouselTestPage(c, built)).Width(320)
				Button(c, "After")
			})
		})
		tester.Frame()
		tester.Frame()
		if len(built) != 0 {
			t.Errorf("%v: a carousel of no pages built %v", p, built)
		}
		if err := tester.Click("Before"); err != nil {
			t.Fatal(err)
		}
		carouselPress(tester, 0, ui.KeyTab)
		if !tester.Focused("After") {
			t.Errorf("%v: Tab from Before stopped before After with no pages", p)
		}
		pager, shown, count = 7, -2, 4
		tester.Frame()
		tester.Frame()
		if pager != 3 || shown != 0 || !tester.HasText("Slide 1") {
			t.Errorf("%v: page 7 of a pager of 4 became %d and page -2 of a carousel %d", p, pager, shown)
		}
		shown = 3
		tester.Frame()
		tester.Frame()
		count = 2
		tester.Advance(time1s)
		tester.Frame()
		tester.Frame()
		if shown != 1 || !tester.HasText("Slide 2") || tester.HasText("Slide 4") {
			t.Errorf("%v: page 4 of a carousel whose pages became 2 showed page %d", p, shown)
		}
	}
}

// TestCarouselKeysFromWhereThePagesShow checks on GNOME that the arrows move from the page nearest to where the pages show, as AdwCarousel's navigate_to_direction takes round (self->position), so that an arrow pressed again just after the first, the pages still near the page they leave, moves to the same page; on Windows and macOS from the page they go to, as FlipView moves from its SelectedIndex, two pages on; and on GNOME a notch of the wheel just after an arrow the same.
func TestCarouselKeysFromWhereThePagesShow(t *testing.T) {
	for _, p := range Platforms {
		s := &carouselTestState{shown: 0}
		tester := carouselMotionTester(p, s, ui.LeftToRight)
		if err := tester.Click("Before"); err != nil {
			t.Fatal(err)
		}
		carouselPress(tester, 0, ui.KeyTab)
		tester.Key(0, ui.KeyRight)
		tester.Frame()
		tester.Advance(10 * 1e6)
		tester.Frame()
		tester.Key(0, ui.KeyRight)
		tester.Frame()
		want := 2
		if p == Linux {
			want = 1
		}
		if s.shown != want {
			t.Errorf("%v: Right twice 10 ms apart showed page %d, want %d", p, s.shown, want)
		}
		if p != Linux {
			continue
		}
		// A notch of the wheel just after a key moves from the page nearest to where the pages show too (scroll_cb's get_page_at_position).
		tester.Advance(time1s)
		tester.Frame()
		tester.Key(0, ui.KeyRight)
		tester.Frame()
		tester.Advance(10 * 1e6)
		tester.Frame()
		x, y := carouselMiddle(tester)
		tester.Wheel(x, y, 0, 40)
		tester.Frame()
		if s.shown != 2 {
			t.Errorf("%v: Right and a notch down 10 ms after from page 1 showed page %d, want 2", p, s.shown)
		}
	}
}

// TestCarouselSwipeEnd checks where a swipe ends, as adw-swipe-tracker.c's get_end_progress chooses without long swipes: a swipe from a page reaches the pages on either side, one from between two pages the page before the first and after the second (get_bounds, within EPSILON 0.005 of a page counting as at it); a slow one goes to the nearest page, the first of two as near; a quick one glides on, a touchpad's 0.997 and a drag's 0.998 a millisecond (0.3323 and 0.499 pages for each DIP a millisecond, along a parabola past 2), to the page beyond when the page before is the one it began at, else to the nearest, so that a swipe back flicked on lands on the page before rather than past the page it began at, and one on flicked back stays on the page after; and the velocity goes on to the spring as AdwCarousel hands it on, DIPs a millisecond taken as pages a second, and none at a bound.
func TestCarouselSwipeEnd(t *testing.T) {
	bounds := []struct {
		from   float32
		lo, hi float32
	}{{1, 0, 2}, {1.5, 0, 3}, {1.004, 0, 2}, {1.006, 0, 3}, {0, 0, 1}, {3, 2, 3}}
	for _, b := range bounds {
		if lo, hi := carouselSwipeBounds(b.from, 4); lo != b.lo || hi != b.hi {
			t.Errorf("a swipe from %v reaches %v to %v, want %v to %v", b.from, lo, hi, b.lo, b.hi)
		}
	}
	glides := []struct {
		speed    float64
		touchpad bool
		want     float64
	}{{1, true, 0.33233}, {1, false, 0.499}, {3, true, 1.34698}, {3, false, 1.84700}}
	for _, g := range glides {
		if got := carouselSwipeGlide(g.speed, g.touchpad); math.Abs(got-g.want) > 1e-4 {
			t.Errorf("a swipe at %v (touchpad %v) glides %v, want %v", g.speed, g.touchpad, got, g.want)
		}
	}
	ends := []struct {
		pos, v float32
		want   int
	}{
		{1.45, 0.2, 1}, {1.5, 0.2, 1}, {1.55, 0.2, 2}, {0.6, -0.5, 1},
		{1.2, 0.7, 2}, {1, -0.7, 0}, {0.2, 0.7, 0}, {0.6, 0.7, 1}, {1.9, -0.7, 2},
	}
	for _, e := range ends {
		st := &carouselState{shown: 1, dist: 400}
		carouselSwipeBegin(st, true, 400, 4)
		st.pos = e.pos
		st.hasEnded, st.ended = true, e.v
		if got := carouselSwipeEnd(st, 4); got != e.want {
			t.Errorf("a touchpad's swipe from 1 ending at %v at %v DIPs a millisecond goes to %d, want %d", e.pos, e.v, got, e.want)
		}
		if want := e.v; st.velocity != want {
			t.Errorf("a swipe ending at %v DIPs a millisecond hands on %v, want %v", e.v, st.velocity, want)
		}
	}
	st := &carouselState{shown: 1}
	carouselSwipeBegin(st, false, 320, 4)
	st.pos, st.hasEnded, st.ended = 2, true, 1
	if got := carouselSwipeEnd(st, 4); got != 2 || st.velocity != 0 {
		t.Errorf("a drag ending at its upper bound going on went to %d with %v, want 2 with none", got, st.velocity)
	}
}

// TestCarouselDragStart checks on GNOME how a drag starts, as drag_update_cb does: the pages follow the pointer from the move that takes it 16 DIPs from where it pressed, the way before it lost (4-DIP moves to 40 DIPs move the pages 28), and a drag that goes more down than sideways at its first move leaves them.
func TestCarouselDragStart(t *testing.T) {
	for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
		s := &carouselTestState{shown: 1}
		tester := carouselMotionTester(Linux, s, dir)
		x, y := carouselMiddle(tester)
		at, _ := tester.Find("Slide 2")
		step := float32(-4)
		if dir == ui.RightToLeft {
			step = 4
		}
		tester.Press(x, y)
		tester.Frame()
		for k := 1; k <= 10; k++ {
			tester.Advance(20 * 1e6)
			tester.Move(x+step*float32(k), y)
			tester.Frame()
		}
		now, _ := tester.Find("Slide 2")
		if d := now.X - at.X; math.Abs(float64(d-7*step)) > 0.01 {
			t.Errorf("%v: 10 moves of 4 DIPs moved the page %v, want %v", dir, d, 7*step)
		}
		tester.Release(x+step*10, y)
		tester.Advance(time1s)
		tester.Frame()
		tester.Frame()
		tester.Press(x, y)
		tester.Frame()
		// 4 sideways and 8 down at first, then only sideways, to 40 sideways.
		for k := 1; k <= 10; k++ {
			tester.Advance(20 * 1e6)
			tester.Move(x+step*float32(k), y+8)
			tester.Frame()
		}
		if now, _ := tester.Find("Slide 2"); now.X != at.X || tester.HasText("Slide 3") || tester.HasText("Slide 1") {
			t.Errorf("%v: a drag first going down moved the pages", dir)
		}
		tester.Release(x+step*10, y+8)
		tester.Frame()
	}
}

// TestAdwCarouselDragAsRendered drags Carousel on GNOME as tools/adwaita-reference's carouseldrag and carouseldragrtl sets drag AdwCarousel, and checks the pages and the dots against libadwaita 1.10.0's render: the mouse pressed in the view's middle and moved 10 times by 10 DIPs toward the next page, 100 ms apart, the first move lost below the 16-DIP threshold and the others applied, so the pages stand 90 DIPs on, at 1.28125 pages, as libadwaita prints; the dots of the page shown, the 2nd, and of the next one (centered at 171.5 and 186.5, 272 down, mirrored right to left) at radii 3.72 and 3.28 and opacities 73% and 47%, their centers (104 and 156, 195 and 138 in the dark theme) as libadwaita's, and the pixel their rims cross on the side away from each other partly covered, as GTK draws them, where dots rounded to whole pixels left the 2nd one's background and filled the 3rd one's. Before the drag, the chosen dot, of radius 4, is drawn as GTK draws it to the level: its rim's pixel 3.5 DIPs right of its middle and half a DIP below it reads (97, 97, 101), (203, 203, 204) in the dark theme, where the dot drawn as a path read (109, 109, 113) and (190, 190, 191); this is checked off macOS, where MyGo rounds a rectangle's corners as Apple does (continuousCorners), as GNOME's dots are drawn on Linux. It does so light and dark, left to right and right to left, and writes carouselShotView after the drag to the directory BREEZE_SCREENSHOTS names.
func TestAdwCarouselDragAsRendered(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, mode := range []string{"light", "dark", "rtl-light", "rtl-dark"} {
		tester := ui.NewTester(carouselShotView(Linux, false), 352, 312)
		tester.SetScale(2)
		dark, rtl := strings.HasSuffix(mode, "dark"), strings.HasPrefix(mode, "rtl")
		tester.SetDark(dark)
		step := float32(-10)
		// The dot of the page shown, the 2nd, and the 3rd's, and the pixels their outer rims cross, in DIPs.
		first, second, firstRim, secondRim := float32(171.5), float32(186.5), float32(175), float32(189.5)
		if rtl {
			tester.SetDirection(ui.RightToLeft)
			step = 10
			first, second, firstRim, secondRim = 186.5, 171.5, 182.5, 168
		}
		tester.Frame()
		tester.Frame()
		rest, rim := [3]int{97, 97, 101}, first+3.5
		if dark {
			rest = [3]int{203, 203, 204}
		}
		if got := shotPixel(tester, rim, 272.5); runtime.GOOS != "darwin" && !samePixel(got, rest, 2) {
			t.Errorf("%s: the chosen dot's rim at rest is %v, want libadwaita's %v", mode, got, rest)
		}
		at, _ := tester.Find("Slide 2")
		x, y := float32(352)/2, carouselShotPlaces[2].Y+carouselTestPageHeight/2
		tester.Press(x, y)
		tester.Frame()
		for k := 1; k <= 10; k++ {
			tester.Advance(100 * time.Millisecond)
			tester.Move(x+step*float32(k), y)
			tester.Frame()
		}
		now, _ := tester.Find("Slide 2")
		if d := now.X - at.X; math.Abs(float64(d-9*step)) > 0.01 {
			t.Errorf("%s: 10 moves of 10 DIPs moved the page %v, want %v", mode, d, 9*step)
		}
		bg, inks := 250, [2][3]int{{104, 104, 108}, {156, 156, 159}}
		if dark {
			bg, inks = 34, [2][3]int{{195, 195, 196}, {138, 138, 140}}
		}
		for k, cx := range []float32{first, second} {
			if got := shotPixel(tester, cx, 272); !samePixel(got, inks[k], 2) {
				t.Errorf("%s: page %d's dot is %v in its middle, want libadwaita's %v", mode, k+2, got, inks[k])
			}
		}
		for k, rx := range []float32{firstRim, secondRim} {
			got, lo, hi := shotPixel(tester, rx, 271.5)[0], min(bg, inks[k][0])+15, max(bg, inks[k][0])-15
			if got < lo || got > hi {
				t.Errorf("%s: page %d's dot's rim at %v is %d, want it partly covered, between %d and %d", mode, k+2, rx, got, lo, hi)
			}
		}
		writeShot(t, dir, fmt.Sprintf("breeze-carousel-drag-linux-%s.png", mode), tester)
		tester.Release(x+step*10, y)
		tester.Frame()
	}
}

// TestAdwCarouselDotsWidth checks that a GNOME Carousel of many small pages is as wide as its dots, as tools/adwaita-reference's carouselmany set renders libadwaita 1.10.0's: a vertical GtkBox with no width of its own, of an AdwCarousel of 25 pages, each a label of its number centered, over AdwCarouselIndicatorDots with a margin of 6, is 399 pixels wide, the dots 387 (adw_carousel_indicator_dots_measure gives their natural width, 15 a page and 6 at either end, as their minimum), and the first, chosen, dot is centered 16.5 into them, where Breeze UI's dots, which asked for no width, left the carousel as wide as its margins. It writes the view, light and dark, to the directory BREEZE_SCREENSHOTS names.
func TestAdwCarouselDotsWidth(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, dark := range []bool{false, true} {
		shown := 0
		var box ui.Rect
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Linux, func() {
				Use(c)
				ui.Box(c).Fill().Children(func() {
					ui.Box(c).Absolute().Left(16).Top(16).Children(func() {
						box = Carousel(c, &shown, 25, func(i int) {
							ui.Column(c).AlignItems(ui.Center).Children(func() { ui.Text(c, fmt.Sprint(i+1)).SingleLine() })
						}).Bounds()
					})
				})
			})
		}, 432, 96)
		tester.SetScale(2)
		tester.SetDark(dark)
		tester.Frame()
		tester.Frame()
		mode := map[bool]string{false: "light", true: "dark"}[dark]
		dots, ok := tester.Find("Pager")
		if !ok {
			t.Fatalf("%s: no pager", mode)
		}
		if box.W != 399 || dots.X != 22 || dots.W != 387 {
			t.Errorf("%s: the carousel is %v wide and its dots %v from %v, want libadwaita's 399, and 387 from 22", mode, box.W, dots.W, dots.X)
		}
		ink := [3]int{70, 70, 74}
		if dark {
			ink = [3]int{233, 233, 234}
		}
		if got := shotPixel(tester, dots.X+16.5, dots.Y+10); !samePixel(got, ink, 2) {
			t.Errorf("%s: the chosen dot's middle is %v, want libadwaita's %v", mode, got, ink)
		}
		writeShot(t, dir, fmt.Sprintf("breeze-carousel-many-linux-%s.png", mode), tester)
	}
}
