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
	drawnWidgets["Text"] = true
}

// textSamples are the styles textView shows, each with the text it shows it in, top to bottom.
var textSamples = []struct {
	style TextStyle
	text  string
}{
	{TextLargeTitle, "Large Title"},
	{TextTitle, "Title"},
	{TextSubtitle, "Subtitle"},
	{TextHeading, "Heading"},
	{TextBody, "Body"},
	{TextCaption, "Caption"},
	{TextSecondary, "Secondary body"},
	{TextCaption | TextSecondary, "Secondary caption"},
}

// textTops are the tops of textSamples' texts in textView, 16 DIPs in from the left, where tools/appkit-reference, tools/winui-reference and tools/adwaita-reference place the platforms' own labels. They leave room for WinUI's 40-pixel title.
var textTops = []float32{16, 76, 122, 160, 188, 214, 240, 266}

// textView shows every step of p's type ramp, each at its place in textTops.
func textView(p Platform) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Box(c).Fill().Children(func() {
				for i, s := range textSamples {
					ui.Box(c).Absolute().Left(16).Top(textTops[i]).Children(func() { Text(c, s.style, s.text) })
				}
			})
		})
	}
}

// textScales are the sizes of the steps as multiples of the body's, as each platform's sources give them: AppKit's text styles on macOS 27.2 (largeTitle 26, title1 22, title2 17, headline 13, body 13, caption1 10), WinUI's TextBlock styles (40, 28, 20, 14, 14, 12) and libadwaita 1.10.0's classes (181%, 136%, 118%, 100%, 100%, 82%).
var textScales = map[Platform]map[TextStyle]float64{
	Mac:     {TextLargeTitle: 26.0 / 13, TextTitle: 22.0 / 13, TextSubtitle: 17.0 / 13, TextHeading: 1, TextBody: 1, TextCaption: 10.0 / 13},
	Windows: {TextLargeTitle: 40.0 / 14, TextTitle: 28.0 / 14, TextSubtitle: 20.0 / 14, TextHeading: 1, TextBody: 1, TextCaption: 12.0 / 14},
	Linux:   {TextLargeTitle: 1.81, TextTitle: 1.36, TextSubtitle: 1.18, TextHeading: 1, TextBody: 1, TextCaption: 0.82},
}

// textBounds builds one line of s in each of styles, on p, and returns their boxes, with the theme's font size.
func textBounds(t *testing.T, p Platform, s string, styles []TextStyle, setup func(*ui.Tester)) ([]ui.Rect, float32) {
	t.Helper()
	var size float32
	tester := ui.NewTester(func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			size = c.Theme().FontSize
			ui.Column(c).Padding(16).Gap(8).AlignItems(ui.Start).Children(func() {
				for i, st := range styles {
					Text(c, st, fmt.Sprintf("%s %d", s, i)).SingleLine()
				}
			})
		})
	}, 900, 600)
	if setup != nil {
		setup(tester)
	}
	tester.Frame()
	var out []ui.Rect
	for i := range styles {
		r, ok := tester.Find(fmt.Sprintf("%s %d", s, i))
		if !ok {
			t.Fatalf("%v: text %d is missing: %q", p, i, tester.Texts())
		}
		out = append(out, r)
	}
	return out, size
}

// TestTextFollowsTheTypeRamp checks every step's size on every platform against the body's, from the boxes of single lines, so that it holds in any font: a step whose line height is the font's own is as many times the body's height as its size is the body's size, and GNOME's caption, whose line height is 140% of its size, is 1.4 times its size high. The titles get smaller step by step, the heading is the body's height, and where the host's font draws a bold face apart from the regular one, the heading is wider than the body.
func TestTextFollowsTheTypeRamp(t *testing.T) {
	styles := []TextStyle{TextLargeTitle, TextTitle, TextSubtitle, TextHeading, TextBody, TextCaption}
	for _, p := range Platforms {
		// The heading's weight: 700 on macOS and GNOME, WinUI's SemiBold on Windows.
		bold := formsWeightsDiffer(400, map[Platform]int{Mac: 700, Windows: 600, Linux: 700}[p])
		r, size := textBounds(t, p, "Sphinx", styles, nil)
		body := float64(r[4].H)
		for i, st := range styles {
			want := textScales[p][st] * body
			if p == Linux && st == TextCaption {
				want = 1.4 * 0.82 * float64(size)
			}
			if got := float64(r[i].H); math.Abs(got-want) > max(1, want*0.04) {
				t.Errorf("%v: %s is %.2f high, want %.2f (body %.2f, font size %v)", p, textSamples[i].text, got, want, body, size)
			}
		}
		if !(r[0].H > r[1].H && r[1].H > r[2].H && r[2].H > r[4].H) {
			t.Errorf("%v: the titles do not get smaller step by step: %v", p, r)
		}
		if r[3].H != r[4].H {
			t.Errorf("%v: the heading is %v high, the body %v", p, r[3].H, r[4].H)
		}
		if bold && r[3].W <= r[4].W {
			t.Errorf("%v: the heading (%v wide) is not bolder than the body (%v)", p, r[3].W, r[4].W)
		}
	}
}

// TestTextBodyIsUIText checks that TextBody draws, on every platform, light and dark, exactly what ui.Text draws in the platform's theme, and that a style that is no step does too.
func TestTextBodyIsUIText(t *testing.T) {
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			var shots []*image.RGBA
			for _, build := range []func(c *ui.Context){
				func(c *ui.Context) { ui.Text(c, "The quick brown fox jumps over the lazy dog") },
				func(c *ui.Context) { Text(c, TextBody, "The quick brown fox jumps over the lazy dog") },
				func(c *ui.Context) { Text(c, TextStyle(42), "The quick brown fox jumps over the lazy dog") },
			} {
				tester := ui.NewTester(func(c *ui.Context) {
					Preview(c, p, func() {
						Use(c)
						ui.Column(c).Padding(16).Width(200).Children(func() { build(c) })
					})
				}, 240, 120)
				tester.SetDark(dark)
				tester.Frame()
				shots = append(shots, tester.Image())
			}
			for i := 1; i < len(shots); i++ {
				if string(shots[i].Pix) != string(shots[0].Pix) {
					t.Errorf("%v, dark %v: Text %d draws differently from ui.Text", p, dark, i)
				}
			}
		}
	}
}

// textInk returns how far the farthest pixel of r in img is from bg, in luminance from 0 to 1: how strongly the text in r is drawn.
func textInk(img *image.RGBA, r ui.Rect, bg ui.Color) float32 {
	far := float32(0)
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			o := img.PixOffset(x, y)
			px := ui.RGB(img.Pix[o], img.Pix[o+1], img.Pix[o+2])
			far = max(far, float32(math.Abs(float64(luminance(px)-luminance(bg)))))
		}
	}
	return far
}

// TestTextSecondaryDims checks TextSecondary on every platform, light and dark: the dimmed text is drawn fainter than the same step undimmed, at the platform's secondary text color on macOS and Windows (secondaryLabelColor, TextFillColorSecondary), and at libadwaita's 55% opacity on GNOME, 90% with higher contrast.
func TestTextSecondaryDims(t *testing.T) {
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			for _, contrast := range []bool{false, true} {
				var th *ui.Theme
				tester := ui.NewTester(func(c *ui.Context) {
					Preview(c, p, func() {
						Use(c)
						th = c.Theme()
						ui.Column(c).Padding(16).Gap(8).AlignItems(ui.Start).Children(func() {
							Text(c, TextHeading, "Wrote").SingleLine()
							Text(c, TextHeading|TextSecondary, "Faint").SingleLine()
						})
					})
				}, 200, 100)
				tester.SetDark(dark)
				tester.SetPreferences(ui.Preferences{HighContrast: contrast})
				tester.Frame()
				img := tester.Image()
				strong, _ := tester.Find("Wrote")
				faint, _ := tester.Find("Faint")
				full, dim := textInk(img, strong, th.Background), textInk(img, faint, th.Background)
				want := luminance(th.TextMuted) - luminance(th.Background)
				if p == Linux {
					want = pick(contrast, float32(0.55), 0.9) * (luminance(th.Text) - luminance(th.Background))
				}
				want = float32(math.Abs(float64(want)))
				if dim >= full || math.Abs(float64(dim-want)) > 0.06 {
					t.Errorf("%v, dark %v, contrast %v: the dimmed text's ink is %.3f, want %.3f (the plain text's %.3f)", p, dark, contrast, dim, want, full)
				}
			}
		}
	}
}

// TestTextKeepsItsBehavior checks on every platform that Text behaves as ui.Text: assistive technology and the tester find it by its text, it takes no focus (Tab goes to the button after it), it wraps at the width it gets, it is selected and copied whole in a Selectable container, a title made a heading with Role keeps its text, and in a right-to-left window it starts at the right.
func TestTextKeepsItsBehavior(t *testing.T) {
	const long = "Breeze UI draws each widget as the desktop the app runs on draws its own."
	for _, p := range Platforms {
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Column(c).Padding(16).Gap(8).AlignItems(ui.Start).Children(func() {
					Text(c, TextTitle, "Settings").Role(ui.RoleHeading).Level(1)
					ui.Column(c).Selectable().Width(160).Children(func() { Text(c, TextCaption|TextSecondary, long) })
					Text(c, TextBody, "One line").SingleLine()
					Button(c, "Done")
				})
			})
		}, 400, 400)
		tester.Frame()
		title, ok := tester.Find("Settings")
		if !ok {
			t.Fatalf("%v: the title is missing: %q", p, tester.Texts())
		}
		line, _ := tester.Find("One line")
		para, ok := tester.Find(long)
		if !ok {
			t.Fatalf("%v: the paragraph is missing", p)
		}
		if para.W > 160 || para.H < 2*line.H*0.7 {
			t.Errorf("%v: the caption did not wrap in 160 DIPs: %v (a line is %v high)", p, para, line.H)
		}
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		if !tester.Focused("Done") {
			t.Errorf("%v: Tab did not go past the texts to the button", p)
		}
		tester.ClickAt(para.X+4, para.Y+para.H/4)
		tester.Key(ui.Cmd, ui.KeyA)
		tester.Key(ui.Cmd, ui.KeyC)
		if got := tester.Clipboard(); got != long {
			t.Errorf("%v: a click, select all and copy took %q", p, got)
		}

		rtl := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() { Text(c, TextSubtitle, "Start") })
			})
		}, 400, 100)
		rtl.SetDirection(ui.RightToLeft)
		rtl.Frame()
		r, _ := rtl.Find("Start")
		if math.Abs(float64(r.X+r.W-384)) > 0.5 {
			t.Errorf("%v: right to left the subtitle is at %v, not 16 from the right (the title was at %v)", p, r, title)
		}
	}
}

// TestScreenshotsText writes textView on every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names, after checking that each text is where textTops puts it and that the titles get smaller step by step. Without it the test only builds the frames.
func TestScreenshotsText(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			tester := ui.NewTester(textView(p), 360, 300)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.Frame()
			var heights []float32
			for i, s := range textSamples {
				r, ok := tester.Find(s.text)
				if !ok {
					t.Fatalf("%v: %q is missing: %q", p, s.text, tester.Texts())
				}
				if r.X != 16 || r.Y != textTops[i] {
					t.Errorf("%v: %q is at %v, want 16, %v", p, s.text, r, textTops[i])
				}
				if i+1 < len(textTops) && r.Y+r.H > textTops[i+1] {
					t.Errorf("%v: %q (%v) runs into the next text at %v", p, s.text, r, textTops[i+1])
				}
				heights = append(heights, r.H)
			}
			if !(heights[0] > heights[1] && heights[1] > heights[2] && heights[2] > heights[4]) {
				t.Errorf("%v: the titles do not get smaller step by step: %v", p, heights)
			}
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			writeShot(t, dir, fmt.Sprintf("breeze-text-%s-%s.png", strings.ToLower(p.String()), mode), tester)
		}
	}
}
