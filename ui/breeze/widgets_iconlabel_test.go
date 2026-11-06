package breeze

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func init() {
	for _, n := range []string{"IconLabelButton", "PrimaryIconLabelButton"} {
		drawnWidgets[n] = true
	}
}

// iconLabelTestSquare is a solid 16-point square, as the template image tools/appkit-reference's probes measure, whose ink tells exactly where the icon is drawn.
var iconLabelTestSquare = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><rect width="16" height="16" fill="#000"/></svg>`))

// iconLabelGeometry is what an icon-label button adds to the button showing only its label on a platform, and where it draws its icon.
type iconLabelGeometry struct {
	h      float32 // its height
	wider  float32 // how much wider it is than the button showing the label alone
	inset  float32 // from its start to the icon
	nudged float32 // how much further from its start the label is
}

// iconLabelWant is the geometry of IconLabelButton (or PrimaryIconLabelButton when primary) on p, in a Toolbar when bar.
func iconLabelWant(p Platform, bar, primary bool) iconLabelGeometry {
	switch {
	case p == Mac && bar:
		// A bordered NSToolbarItem with an image and a title: 82 points wide where the title alone is 64, the image 12 in.
		return iconLabelGeometry{36, 18, 12, 18}
	case p == Mac:
		// An NSButton push button with imagePosition .imageLeading: 82 points wide where the title alone is 64, the image 12 in, the title 2 after it.
		return iconLabelGeometry{24, 18, 12, 18}
	case p == Windows && bar && !primary:
		// An AppBarButton with its label on the right: the icon 10 inside the face and the label 8 after it, where the label alone is 19 in.
		return iconLabelGeometry{36, 14, 10, 15}
	case p == Windows:
		// A Button holding a StackPanel of a FontIcon and a TextBlock 8 apart, 12 in past its border.
		return iconLabelGeometry{32, 24, 12, 24}
	}
	// .image-text-button: 9 in, the icon 6 from the label, which keeps 2 at its end, where a .text-button is padded 17.
	return iconLabelGeometry{34, 8, 9, 14}
}

// iconLabelInk is the color the platform gives the icon and label of an icon-label button at rest: the label's.
func iconLabelInk(c *ui.Context, p Platform, bar, primary bool) ui.Color {
	if primary {
		return c.Theme().AccentText
	}
	return inputsIconInk(c, p, bar)
}

// inkStart returns the first pixel of row y from x0 towards x1 (either way) that differs from face by more than 60 in any channel, or -1.
func inkStart(img *image.RGBA, y, x0, x1 int, face color.RGBA) int {
	step := 1
	if x1 < x0 {
		step = -1
	}
	for x := x0; x != x1; x += step {
		px := img.RGBAAt(x, y)
		if max(absDiff(px.R, face.R), absDiff(px.G, face.G), absDiff(px.B, face.B)) > 60 {
			return x
		}
	}
	return -1
}

// TestIconLabelButtonGeometry checks IconLabelButton and PrimaryIconLabelButton on every platform, alone and in a Toolbar, left to right and right to left, against the button showing the label alone beside it: its height, how much wider it is, the icon drawn where the platform puts its image, in the label's color, and the label moved past it by the platform's spacing.
func TestIconLabelButtonGeometry(t *testing.T) {
	for _, p := range Platforms {
		for _, bar := range []bool{false, true} {
			for _, primary := range []bool{false, true} {
				for _, rtl := range []bool{false, true} {
					name := fmt.Sprintf("%v toolbar %v primary %v rtl %v", p, bar, primary, rtl)
					var a, b ui.Rect
					var ink ui.Color
					tester := ui.NewTester(func(c *ui.Context) {
						Preview(c, p, func() {
							Use(c)
							ui.Column(c).Padding(40).AlignItems(ui.Start).Children(func() {
								build := func() {
									if primary {
										a = PrimaryButton(c, "Export").Bounds()
										b = PrimaryIconLabelButton(c, iconLabelTestSquare, "Export").Bounds()
									} else {
										a = Button(c, "Export").Bounds()
										b = IconLabelButton(c, iconLabelTestSquare, "Export").Bounds()
									}
								}
								if bar {
									Toolbar(c, build)
								} else {
									ui.Row(c).Gap(16).Children(build)
								}
							})
							ink = iconLabelInk(c, p, bar, primary)
						})
					}, 480, 200)
					tester.SetScale(2)
					if rtl {
						tester.SetDirection(ui.RightToLeft)
					}
					tester.Frame()
					tester.Frame()
					tester.Frame()
					want := iconLabelWant(p, bar, primary)
					if b.H != want.h || b.W-a.W != want.wider {
						t.Errorf("%s: the button is %v by %v and the label's alone %v wide, want %v high and %v wider", name, b.W, b.H, a.W, want.h, want.wider)
					}
					img := tester.Image()
					// Along the middle of each button, from its start inward, past its edge and its face's rim. A button after one of a fractional width starts between pixels, so where ink starts is read to the pixel, half a point.
					edge := func(r ui.Rect, in float32) int {
						if rtl {
							return int(2*(r.X+r.W-in)) - 1
						}
						return int(2 * (r.X + in))
					}
					inward := func(r ui.Rect) int {
						if rtl {
							return int(2 * r.X)
						}
						return int(2 * (r.X + r.W))
					}
					// from is how far into r from its start the pixel x is.
					from := func(r ui.Rect, x int) float32 {
						if rtl {
							return r.X + r.W - float32(x+1)/2
						}
						return float32(x)/2 - r.X
					}
					near := func(got, want float32) bool { return got-want <= 0.5 && want-got <= 0.5 }
					y := int(2 * (b.Y + b.H/2))
					face := img.RGBAAt(edge(b, 4), y)
					if x := inkStart(img, y, edge(b, 4), inward(b), face); x < 0 || !near(from(b, x), want.inset) {
						t.Errorf("%s: the icon's ink starts %v in, want %v", name, from(b, x), want.inset)
					}
					for _, in := range []float32{want.inset + 1, want.inset + 15} {
						px := img.RGBAAt(edge(b, in), y)
						if d := max(absDiff(px.R, ink.R), absDiff(px.G, ink.G), absDiff(px.B, ink.B)); d > 24 {
							t.Errorf("%s: the icon reads %v %v in, want the label's color %v", name, px, in, ink)
						}
					}
					// The label's first ink, past the icon, against the label's alone: the same glyph, moved by the platform's spacing.
					ya := int(2 * (a.Y + a.H/2))
					label := inkStart(img, y-2, edge(b, want.inset+17), inward(b), face)
					alone := inkStart(img, ya-2, edge(a, 4), inward(a), img.RGBAAt(edge(a, 4), ya))
					if label < 0 || alone < 0 || !near(from(b, label)-from(a, alone), want.nudged) {
						t.Errorf("%s: the label's ink starts %v in, %v in the label alone, want %v further", name, from(b, label), from(a, alone), want.nudged)
					}
				}
			}
		}
	}
}

// TestIconLabelButtonBehavior checks on every platform, alone, in a Toolbar and in a HeaderBar, that the button is as high as the platform's there (on Windows a HeaderBar's TitleBar gives it the look it has outside a toolbar), that the label names the button itself for assistive technology and the tester, that a click and Space on the focused button click it, and that a disabled one neither clicks nor takes the focus; and on macOS that it never joins the image items beside it into one platter, as AppKit joins only image-only items.
func TestIconLabelButtonBehavior(t *testing.T) {
	heights := map[Platform]map[string]float32{
		Mac:     {"row": 24, "toolbar": 36, "header": 36},
		Windows: {"row": 32, "toolbar": 36, "header": 32},
		Linux:   {"row": 34, "toolbar": 34, "header": 34},
	}
	for _, p := range Platforms {
		for _, in := range []string{"row", "toolbar", "header"} {
			name := fmt.Sprintf("%v in a %s", p, in)
			clicks, deletes := 0, 0
			var share, export, trash ui.Rect
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					ui.Column(c).Padding(40).AlignItems(ui.Start).Children(func() {
						build := func() {
							share = IconButton(c, idiomTestShareIcon, "Share").Bounds()
							e := IconLabelButton(c, idiomTestShareIcon, "Export")
							if e.Clicked() {
								clicks++
							}
							export = e.Bounds()
							d := IconLabelButton(c, inputsTestTrashIcon, "Delete").Disabled(true)
							if d.Clicked() {
								deletes++
							}
							trash = d.Bounds()
						}
						switch in {
						case "toolbar":
							Toolbar(c, build)
						case "header":
							HeaderBar(c, "Files", "", nil, build).Width(400)
						default:
							ui.Row(c).Gap(8).Children(build)
						}
					})
				})
			}, 480, 200)
			tester.SetScale(2)
			tester.Frame()
			tester.Frame()
			tester.Frame()
			if export.H != heights[p][in] {
				t.Errorf("%s: the button is %v high, want %v", name, export.H, heights[p][in])
			}
			// Find scales the box by the zoom, which may move its last digit.
			same := func(a, b ui.Rect) bool {
				return max(a.X-b.X, b.X-a.X, a.Y-b.Y, b.Y-a.Y, a.W-b.W, b.W-a.W, a.H-b.H, b.H-a.H) < 0.01
			}
			if r, ok := tester.Find("Export"); !ok || !same(r, export) {
				t.Errorf("%s: the element named Export is %v, %v, want the button %v", name, r, ok, export)
			}
			if strings.Count(strings.Join(tester.Texts(), "\n"), "Export") != 2 {
				t.Errorf("%s: want the button's name and its label, got %q", name, tester.Texts())
			}
			if err := tester.Click("Export"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			if clicks != 1 {
				t.Fatalf("%s: a click clicked %d times", name, clicks)
			}
			tester.Key(0, ui.KeySpace)
			tester.Frame()
			if clicks != 2 {
				t.Errorf("%s: Space on the focused button clicked %d times in all", name, clicks)
			}
			if err := tester.Click("Delete"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			if deletes != 0 || tester.Focused("Delete") {
				t.Errorf("%s: the disabled button clicked %d times, focused %v", name, deletes, tester.Focused("Delete"))
			}
			if p == Mac && in == "toolbar" {
				s, e, d := share, export, trash
				if e.X-(s.X+s.W) != idiomMacToolbarGap || d.X-(e.X+e.W) != idiomMacToolbarGap {
					t.Errorf("%s: the items are at %v, %v and %v, want them %v apart, each on its own platter", name, s, e, d, idiomMacToolbarGap)
				}
			}
		}
	}
}

// iconLabelView lays out the icon-label buttons on platform p as tools/appkit-reference/reference.swift places AppKit's (iconLabelSet): a push button showing an image and a title, the default one and a disabled one; below them, a toolbar holding an item with an image and a title, an image item and a prominent item with an image and a title.
func iconLabelView(p Platform) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Box(c).Fill().Children(func() {
				at := func(x, y float32, build func()) {
					ui.Box(c).Absolute().Left(x).Top(y).Children(build)
				}
				at(16, 16, func() { IconLabelButton(c, idiomTestShareIcon, "Share") })
				at(130, 16, func() { PrimaryIconLabelButton(c, idiomTestShareIcon, "Export") })
				at(250, 16, func() { IconLabelButton(c, inputsTestTrashIcon, "Delete").Disabled(true) })
				ui.Box(c).Absolute().Left(0).Top(64).Width(520).Children(func() {
					Toolbar(c, func() {
						IconLabelButton(c, idiomTestShareIcon, "Share")
						IconButton(c, inputsTestTrashIcon, "Delete")
						PrimaryIconLabelButton(c, idiomTestShareIcon, "Export")
					})
				})
			})
		})
	}
}

// TestScreenshotsIconLabel writes iconLabelView on every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names. Without it the test only builds the frames.
func TestScreenshotsIconLabel(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			tester := ui.NewTester(iconLabelView(p), 520, 130)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.Frame()
			tester.Frame()
			for _, text := range []string{"Share", "Export", "Delete"} {
				if _, ok := tester.Find(text); !ok {
					t.Fatalf("%v: %q is missing", p, text)
				}
			}
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			writeShot(t, dir, fmt.Sprintf("breeze-iconlabel-%s-%s.png", strings.ToLower(p.String()), mode), tester)
		}
	}
}
