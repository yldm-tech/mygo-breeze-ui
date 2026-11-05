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
					ui.Box(c).Absolute().Left(16).Top(textTops[i]).Children(func() { Text(c, s.text, s.style) })
				}
			})
		})
	}
}

// textScales are the sizes of the steps as multiples of the body's, as each platform's sources give them: AppKit's text styles on macOS 27.2 (largeTitle 26, title1 22, title2 17, headline 13, body 13, caption1 10), WinUI's TextBlock styles (40, 28, 20, 14, 14, 12) and libadwaita 1.10.0's classes (181%, 136%, 118%, 100%, 100%, 82%).
var textScales = map[Platform]map[TextStyle]float32{
	Mac:     {TextLargeTitle: 26.0 / 13, TextTitle: 22.0 / 13, TextSubtitle: 17.0 / 13, TextHeading: 1, TextBody: 1, TextCaption: 10.0 / 13},
	Windows: {TextLargeTitle: 40.0 / 14, TextTitle: 28.0 / 14, TextSubtitle: 20.0 / 14, TextHeading: 1, TextBody: 1, TextCaption: 12.0 / 14},
	Linux:   {TextLargeTitle: 1.81, TextTitle: 1.36, TextSubtitle: 1.18, TextHeading: 1, TextBody: 1, TextCaption: 0.82},
}

// textWeights are the steps' weights: AppKit's headline in bold and the others regular; WinUI's SemiBold from BaseTextBlockStyle but on Body and Caption; libadwaita's 800 on .title-1 and .title-2, bold on .title-4 and .heading.
var textWeights = map[Platform]map[TextStyle]int{
	Mac:     {TextLargeTitle: 400, TextTitle: 400, TextSubtitle: 400, TextHeading: 700, TextBody: 400, TextCaption: 400},
	Windows: {TextLargeTitle: 600, TextTitle: 600, TextSubtitle: 600, TextHeading: 600, TextBody: 400, TextCaption: 400},
	Linux:   {TextLargeTitle: 800, TextTitle: 800, TextSubtitle: 700, TextHeading: 700, TextBody: 400, TextCaption: 400},
}

// textRampStyles are the steps in the order of the ramp, largest first.
var textRampStyles = []TextStyle{TextLargeTitle, TextTitle, TextSubtitle, TextHeading, TextBody, TextCaption}

// textBoxes builds one line for each of textRampStyles on p, with build, and returns their boxes and the theme's font size.
func textBoxes(t *testing.T, p Platform, prefs ui.Preferences, build func(c *ui.Context, st TextStyle, s string)) ([]ui.Rect, float32) {
	t.Helper()
	var size float32
	tester := ui.NewTester(func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			size = c.Theme().FontSize
			ui.Column(c).Padding(16).Gap(8).AlignItems(ui.Start).Children(func() {
				for i, st := range textRampStyles {
					build(c, st, fmt.Sprintf("Sphinx %d", i))
				}
			})
		})
	}, 900, 600)
	tester.SetPreferences(prefs)
	tester.Frame()
	var out []ui.Rect
	for i := range textRampStyles {
		r, ok := tester.Find(fmt.Sprintf("Sphinx %d", i))
		if !ok {
			t.Fatalf("%v: text %d is missing: %q", p, i, tester.Texts())
		}
		out = append(out, r)
	}
	return out, size
}

// TestTextFollowsTheTypeRamp checks every step on every platform, at the usual text size and at the desktop's 150%, against a ui.Text given the platform's size (the theme's FontSize times textScales), weight (textWeights) and line height (140% on GNOME's caption, the font's own elsewhere) by hand: their boxes must be the same, whatever font the host has. Then, from the boxes, the titles get smaller step by step, the heading is the body's height, and where the host's font draws the heading's weight apart from the regular one, the heading is wider than the body.
func TestTextFollowsTheTypeRamp(t *testing.T) {
	for _, p := range Platforms {
		bold := formsWeightsDiffer(400, textWeights[p][TextHeading])
		var usual float32
		for _, scale := range []float32{1, 1.5} {
			prefs := ui.Preferences{TextScale: scale}
			r, size := textBoxes(t, p, prefs, func(c *ui.Context, st TextStyle, s string) { Text(c, s, st).SingleLine() })
			want, _ := textBoxes(t, p, prefs, func(c *ui.Context, st TextStyle, s string) {
				e := ui.Text(c, s).SingleLine()
				if st == TextBody {
					return
				}
				e.FontSize(c.Theme().FontSize * textScales[p][st]).FontWeight(textWeights[p][st])
				if p == Linux && st == TextCaption {
					e.LineHeight(1.4)
				}
			})
			if scale == 1 {
				usual = size
			} else if size != usual*scale {
				t.Errorf("%v: the theme's font size is %v at a text scale of %v, %v usually", p, size, scale, usual)
			}
			for i := range textRampStyles {
				if r[i] != want[i] {
					t.Errorf("%v at %v: %s is %v, want %v (font size %v)", p, scale, textSamples[i].text, r[i], want[i], size)
				}
			}
			if !(r[0].H > r[1].H && r[1].H > r[2].H && r[2].H > r[4].H) {
				t.Errorf("%v at %v: the titles do not get smaller step by step: %v", p, scale, r)
			}
			if math.Abs(float64(r[3].H-r[4].H)) > 0.01 {
				t.Errorf("%v at %v: the heading is %v high, the body %v", p, scale, r[3].H, r[4].H)
			}
			if bold && r[3].W <= r[4].W {
				t.Errorf("%v at %v: the heading (%v wide) is not bolder than the body (%v)", p, scale, r[3].W, r[4].W)
			}
		}
	}
}

// TestTextBodyIsUIText checks that Text without a style, and with TextBody, draws, on every platform, light and dark, exactly what ui.Text draws in the platform's theme, so that a view can swap ui.Text for it, and that a style that is no step does too.
func TestTextBodyIsUIText(t *testing.T) {
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			var shots []*image.RGBA
			for _, build := range []func(c *ui.Context){
				func(c *ui.Context) { ui.Text(c, "The quick brown fox jumps over the lazy dog") },
				func(c *ui.Context) { Text(c, "The quick brown fox jumps over the lazy dog") },
				func(c *ui.Context) { Text(c, "The quick brown fox jumps over the lazy dog", TextBody) },
				func(c *ui.Context) { Text(c, "The quick brown fox jumps over the lazy dog", TextStyle(42)) },
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

// textInk returns how far the farthest pixel of r, in DIPs, in img, drawn at scale, is from bg, in luminance from 0 to 1: how strongly the text in r is drawn.
func textInk(img *image.RGBA, r ui.Rect, scale float32, bg ui.Color) float32 {
	far := float32(0)
	for y := int(r.Y * scale); y < int((r.Y+r.H)*scale); y++ {
		for x := int(r.X * scale); x < int((r.X+r.W)*scale); x++ {
			o := img.PixOffset(x, y)
			px := ui.RGB(img.Pix[o], img.Pix[o+1], img.Pix[o+2])
			far = max(far, float32(math.Abs(float64(luminance(px)-luminance(bg)))))
		}
	}
	return far
}

// TestTextSecondaryDims checks TextSecondary on every platform, light and dark, at 2x so that the strokes of the host's font cover whole pixels: the dimmed text is drawn fainter than the same step undimmed, at the platform's secondary text color on macOS and Windows (secondaryLabelColor, TextFillColorSecondary), and at libadwaita's 55% opacity on GNOME, 90% with higher contrast.
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
							Text(c, "Wrote", TextHeading).SingleLine()
							Text(c, "Faint", TextHeading|TextSecondary).SingleLine()
						})
					})
				}, 200, 100)
				tester.SetScale(2)
				tester.SetDark(dark)
				tester.SetPreferences(ui.Preferences{HighContrast: contrast})
				tester.Frame()
				img := tester.Image()
				strong, _ := tester.Find("Wrote")
				faint, _ := tester.Find("Faint")
				full, dim := textInk(img, strong, 2, th.Background), textInk(img, faint, 2, th.Background)
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
					Text(c, "Settings", TextTitle).Role(ui.RoleHeading).Level(1)
					ui.Column(c).Selectable().Width(160).Children(func() { Text(c, long, TextCaption, TextSecondary) })
					Text(c, "One line").SingleLine()
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
				ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() { Text(c, "Start", TextSubtitle) })
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
