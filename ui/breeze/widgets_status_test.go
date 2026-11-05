package breeze

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func init() {
	for _, n := range []string{"InfoBar", "ProgressRing", "StatusPage"} {
		drawnWidgets[n] = true
	}
}

// statusSearchIcon is a magnifying glass, as the status page of a search without results shows.
var statusSearchIcon = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><circle cx="6.5" cy="6.5" r="4.75" fill="none" stroke="#000" stroke-width="1.5"/><path d="M10 10l4.5 4.5" stroke="#000" stroke-width="1.5" stroke-linecap="round"/></svg>`))

// infoBarCloseName is the accessible name of an InfoBar's close button on p.
func infoBarCloseName(p Platform) string {
	return map[Platform]string{Mac: "Close", Windows: "Close", Linux: "Dismiss"}[p]
}

// TestInfoBarKeepsItsBehavior checks the bar on every platform, left to right and right to left: it shows its title and message, its close button, at the end of the reading direction, closes it with a click and with the keyboard, a disabled bar does not close, and a bar without open has no close button.
func TestInfoBarKeepsItsBehavior(t *testing.T) {
	for _, p := range Platforms {
		for _, rtl := range []bool{false, true} {
			open, fixed := true, false
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					ui.Column(c).Padding(16).Gap(12).Children(func() {
						InfoBar(c, "warning", "Update available", "Restart to finish installing.", &open)
						InfoBar(c, "error", "Disk full", "Free some space.", nil)
						if Button(c, "Fix").Clicked() {
							fixed = true
						}
					})
				})
			}, 600, 300)
			if rtl {
				tester.SetDirection(ui.RightToLeft)
			}
			tester.Frame()
			for _, s := range []string{"Update available", "Restart to finish installing.", "Disk full", "Free some space."} {
				if !tester.HasText(s) {
					t.Fatalf("%v rtl %v: %q is missing: %q", p, rtl, s, tester.Texts())
				}
			}
			name := infoBarCloseName(p)
			b, ok := tester.Find(name)
			if !ok {
				t.Fatalf("%v rtl %v: no close button named %q", p, rtl, name)
			}
			if b.Y > 16+60 {
				t.Errorf("%v rtl %v: the close button at %v belongs to the second bar, which has none", p, rtl, b)
			}
			// The close button is at the end of the reading direction.
			if atEnd := b.X+b.W/2 > 300; atEnd == rtl {
				t.Errorf("%v rtl %v: the close button is at %v", p, rtl, b)
			}
			if err := tester.Click(name); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			if open || tester.HasText("Update available") {
				t.Fatalf("%v rtl %v: a click on the close button left the bar open", p, rtl)
			}
			if !tester.HasText("Disk full") || fixed {
				t.Fatalf("%v rtl %v: closing one bar changed the rest", p, rtl)
			}
			// The keyboard: Tab reaches the close button, Space presses it.
			open = true
			tester.Frame()
			for range 4 {
				if tester.Focused(name) {
					break
				}
				tester.Key(0, ui.KeyTab)
				tester.Frame()
			}
			if !tester.Focused(name) {
				t.Fatalf("%v rtl %v: Tab does not reach the close button", p, rtl)
			}
			tester.Key(0, ui.KeySpace)
			tester.Frame()
			if open {
				t.Fatalf("%v rtl %v: Space on the close button left the bar open", p, rtl)
			}
		}
		// A bar in a disabled column does not close.
		open := true
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Column(c).Padding(16).Disabled(true).Children(func() {
					InfoBar(c, "", "Offline", "", &open)
				})
			})
		}, 600, 200)
		tester.Frame()
		_ = tester.Click(infoBarCloseName(p))
		tester.Frame()
		if !open || !tester.HasText("Offline") {
			t.Errorf("%v: a disabled bar closed", p)
		}
	}
}

// TestInfoBarSeverities checks the bar's kinds: Windows colors each severity's bar (the informational one, success, warning and error, as a toast's Type names them), macOS draws them on one fill with a symbol of the severity's color, and GNOME draws every kind the same, as libadwaita's banners have no severity.
func TestInfoBarSeverities(t *testing.T) {
	kinds := []string{"info", "success", "warning", "error"}
	for _, p := range Platforms {
		fills := map[string]color.Color{}
		symbols := map[string]color.Color{}
		for _, kind := range kinds {
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					ui.Column(c).Padding(16).Children(func() {
						InfoBar(c, kind, "Title", "Message", nil)
					})
				})
			}, 500, 160)
			tester.Frame()
			if !tester.HasText("Title") || !tester.HasText("Message") {
				t.Fatalf("%v %s: texts %q", p, kind, tester.Texts())
			}
			title, _ := tester.Find("Title")
			img := tester.Image()
			fills[kind] = img.At(22, int(title.Y+title.H/2))
			// The macOS symbol's middle, 11.5 points in from the bar's start, 13 points across.
			symbols[kind] = img.At(16+11.5+6.5, int(title.Y+title.H/2)-4)
		}
		distinct := func(m map[string]color.Color) int {
			seen := map[color.Color]bool{}
			for _, c := range m {
				seen[c] = true
			}
			return len(seen)
		}
		switch p {
		case Windows:
			if n := distinct(fills); n != 4 {
				t.Errorf("Windows: %d distinct bar colors for four severities: %v", n, fills)
			}
		case Linux:
			if n := distinct(fills); n != 1 {
				t.Errorf("GNOME: the banner's color changes with the kind: %v", fills)
			}
		default:
			if n := distinct(fills); n != 1 {
				t.Errorf("macOS: the section's fill changes with the kind: %v", fills)
			}
			if n := distinct(symbols); n != 4 {
				t.Errorf("macOS: %d distinct symbol colors for four severities: %v", n, symbols)
			}
		}
	}
	for typ, want := range map[string]int{"": fluentToastInformational, "info": fluentToastInformational, "success": fluentToastSuccess, "Warning": fluentToastWarning, "danger": fluentToastError, "critical": fluentToastError} {
		if got := fluentToastSeverity(typ); got != want {
			t.Errorf("fluentToastSeverity(%q) = %d, want %d", typ, got, want)
		}
	}
}

// ringPixel returns the color of img on the ring of radius rad around cx, cy at angle degrees clockwise from 12 o'clock.
func ringPixel(img interface{ At(x, y int) color.Color }, cx, cy, rad, degrees float64) color.Color {
	a := degrees * math.Pi / 180
	return img.At(int(cx+rad*math.Sin(a)), int(cy-rad*math.Cos(a)))
}

// cropPixels returns the pixels of img inside r, to compare two frames there.
func cropPixels(img *image.RGBA, r ui.Rect) string {
	return string(img.SubImage(image.Rect(int(r.X), int(r.Y), int(r.X+r.W), int(r.Y+r.H))).(*image.RGBA).Pix)
}

// nearColor reports whether a and b differ by at most tol in every channel, as the antialiased edge of a thin stroke lets two points of one fill differ.
func nearColor(a, b color.Color, tol int) bool {
	ar, ag, ab, _ := a.RGBA()
	br, bg, bb, _ := b.RGBA()
	d := func(x, y uint32) bool { return math.Abs(float64(int(x>>8)-int(y>>8))) <= float64(tol) }
	return d(ar, br) && d(ag, bg) && d(ab, bb)
}

// TestProgressRingKeepsItsBehavior checks the ring on every platform: it is found by its name at the platform's size, Size gives another, it fills clockwise from 12 o'clock in either reading direction, a full ring is the same all round, and a negative value shows the platform's spinner.
func TestProgressRingKeepsItsBehavior(t *testing.T) {
	for _, p := range Platforms {
		for _, rtl := range []bool{false, true} {
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					ui.Row(c).Padding(16).Gap(16).AlignItems(ui.Start).Children(func() {
						ProgressRing(c, 0.25).Label("Quarter")
						ProgressRing(c, 1).Label("Full")
						ProgressRing(c, 0.6).Size(16, 16).Label("Small")
						ProgressRing(c, -1).Label("Busy")
					})
				})
			}, 400, 100)
			if rtl {
				tester.SetDirection(ui.RightToLeft)
			}
			tester.Frame()
			q, ok := tester.Find("Quarter")
			if !ok || q.W != 32 || q.H != 32 {
				t.Fatalf("%v: ring %v, %v; want 32 by 32", p, q, ok)
			}
			if s, ok := tester.Find("Small"); !ok || s.W != 16 || s.H != 16 {
				t.Fatalf("%v: small ring %v, %v", p, s, ok)
			}
			if _, ok := tester.Find("Busy"); !ok {
				t.Fatalf("%v: no spinner for a negative value", p)
			}
			img := tester.Image()
			busy, _ := tester.Find("Busy")
			before := cropPixels(img, busy)
			tester.Advance(250 * time.Millisecond)
			tester.Frame()
			if cropPixels(tester.Image(), busy) == before {
				t.Errorf("%v: the spinner of a negative value does not move", p)
			}
			// The middle of each platform's stroke: 13.5 from the center on macOS, 14.16 on Windows, 14 on GNOME.
			rad := map[Platform]float64{Mac: 13.5, Windows: 14.16, Linux: 14}[p]
			cx, cy := float64(q.X+q.W/2), float64(q.Y+q.H/2)
			filled, empty := ringPixel(img, cx, cy, rad, 45), ringPixel(img, cx, cy, rad, 225)
			if filled == empty {
				t.Errorf("%v rtl %v: a quarter ring is the same at 45 and 225 degrees: %v", p, rtl, filled)
			}
			if back := ringPixel(img, cx, cy, rad, 315); back != empty {
				t.Errorf("%v rtl %v: a quarter ring is filled at 315 degrees: %v, %v", p, rtl, back, empty)
			}
			f, _ := tester.Find("Full")
			fx, fy := float64(f.X+f.W/2), float64(f.Y+f.H/2)
			for _, deg := range []float64{90, 180, 270} {
				if got := ringPixel(img, fx, fy, rad, deg); !nearColor(got, filled, 8) {
					t.Errorf("%v rtl %v: a full ring at %v degrees is %v, at 45 on a quarter %v", p, rtl, deg, got, filled)
				}
			}
		}
	}
}

// TestFluentRingTrim checks the determinate ProgressRing's trim against its Lottie's keyframes: nothing before 1/120, the value itself at the quarters, 1 from 0.9917.
func TestFluentRingTrim(t *testing.T) {
	for _, tc := range []struct{ value, want float64 }{{0, 0}, {0.005, 0}, {0.25, 0.25}, {0.5, 0.5}, {0.75, 0.75}, {0.991666675, 1}, {1, 1}} {
		if got := fluentRingTrim(tc.value); math.Abs(got-tc.want) > 1e-4 {
			t.Errorf("fluentRingTrim(%v) = %v, want %v", tc.value, got, tc.want)
		}
	}
	for v := 0.01; v < 1; v += 0.01 {
		if fluentRingTrim(v) < fluentRingTrim(v-0.01) {
			t.Fatalf("fluentRingTrim goes back at %v", v)
		}
	}
}

// TestStatusPageKeepsItsBehavior checks the page on every platform: its texts show centered across it, the actions below them take clicks and the keyboard, and a page without an icon, a description or actions still shows its title.
func TestStatusPageKeepsItsBehavior(t *testing.T) {
	for _, p := range Platforms {
		cleared := 0
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Column(c).Fill().Children(func() {
					StatusPage(c, statusSearchIcon, "No Results", "Check the spelling or try a new search.", func() {
						if Button(c, "Clear Search").Clicked() {
							cleared++
						}
					})
				})
			})
		}, 600, 500)
		tester.Frame()
		d, ok := tester.Find("Check the spelling or try a new search.")
		if !ok {
			t.Fatalf("%v: no description: %q", p, tester.Texts())
		}
		if mid := d.X + d.W/2; mid < 299 || mid > 301 {
			t.Errorf("%v: the description is at %v, not centered across the page", p, d)
		}
		b, ok := tester.Find("Clear Search")
		if !ok || b.Y < d.Y+d.H {
			t.Fatalf("%v: the action %v, %v is not below the description %v", p, b, ok, d)
		}
		if err := tester.Click("Clear Search"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		if !tester.Focused("Clear Search") {
			t.Fatalf("%v: Tab does not reach the action", p)
		}
		tester.Key(0, ui.KeySpace)
		tester.Frame()
		if cleared != 2 {
			t.Errorf("%v: the action ran %d times", p, cleared)
		}
		bare := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Column(c).Fill().Children(func() { StatusPage(c, nil, "Empty", "", nil) })
			})
		}, 400, 300)
		bare.Frame()
		if !bare.HasText("Empty") {
			t.Errorf("%v: a bare page lost its title: %q", p, bare.Texts())
		}
	}
}

// statusView lays out the status widgets where tools/appkit-reference/reference.swift, tools/winui-reference and tools/adwaita-reference place their counterparts: four message bars, one of each severity, the first and the last closable, across the window 16 from its edges and 68 apart, and below them rings of 0.25, 0.6, 1 and 0 at the regular size and one of 0.6 at 16.
func statusView(p Platform, open *[4]bool) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Box(c).Fill().Children(func() {
				at := func(x, y float32, build func()) {
					ui.Column(c).Absolute().Left(x).Top(y).Right(16).AlignItems(ui.Start).Children(build)
				}
				bars := []struct{ kind, title, message string }{
					{"info", "Informational", "Some message about it."},
					{"success", "Success", "Some message about it."},
					{"warning", "Warning", "Some message about it."},
					{"error", "Error", "Some message about it."},
				}
				for i, b := range bars {
					var o *bool
					if i == 0 || i == 3 {
						o = &open[i]
					}
					at(16, 16+float32(i)*68, func() { InfoBar(c, b.kind, b.title, b.message, o) })
				}
				for i, v := range []float64{0.25, 0.6, 1, 0} {
					ui.Box(c).Absolute().Left(16 + float32(i)*48).Top(300).Children(func() { ProgressRing(c, v).Label(fmt.Sprint("Ring ", i)) })
				}
				ui.Box(c).Absolute().Left(208).Top(308).Children(func() { ProgressRing(c, 0.6).Size(16, 16).Label("Small ring") })
			})
		})
	}
}

// statusPageView lays out a status page filling the window, with an icon, a title, a description and an action, as reference.swift places SwiftUI's ContentUnavailableView.
func statusPageView(p Platform) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Column(c).Fill().Children(func() {
				StatusPage(c, statusSearchIcon, "No Results", "Check the spelling or try a new search.", func() { Button(c, "Clear Search") })
			})
		})
	}
}

// TestScreenshotsStatus writes statusView and statusPageView on every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names. Without it the test only builds the frames.
func TestScreenshotsStatus(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			open := [4]bool{true, true, true, true}
			tester := ui.NewTester(statusView(p, &open), 560, 360)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.Frame()
			for _, s := range []string{"Informational", "Success", "Warning", "Error", "Ring 0", "Small ring"} {
				if _, ok := tester.Find(s); !ok {
					t.Fatalf("%v: %q is missing: %q", p, s, tester.Texts())
				}
			}
			writeShot(t, dir, fmt.Sprintf("breeze-status-%s-%s.png", strings.ToLower(p.String()), mode), tester)
			page := ui.NewTester(statusPageView(p), 560, 420)
			page.SetScale(2)
			page.SetDark(dark)
			page.Frame()
			if !page.HasText("No Results") || !page.HasText("Clear Search") {
				t.Fatalf("%v: the status page is missing its texts: %q", p, page.Texts())
			}
			writeShot(t, dir, fmt.Sprintf("breeze-statuspage-%s-%s.png", strings.ToLower(p.String()), mode), page)
		}
	}
}

// TestInfoBarArrangements checks the arrangements Windows and GNOME switch between by the bar's width: in a wide bar WinUI's message follows the title on its line and AdwBanner's button sits beside its text; in a narrow one the message goes below the title, and the banner's button below its text. The bar is narrower than the window, so the first frame, which knows only the window's width, must be laid out again at the bar's.
func TestInfoBarArrangements(t *testing.T) {
	const message = "Restart to finish installing."
	for _, p := range []Platform{Windows, Linux} {
		for _, width := range []float32{560, 220} {
			open := true
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					ui.Column(c).Padding(16).Width(width + 32).Children(func() {
						InfoBar(c, "info", "Update available", message, &open)
					})
				})
			}, 800, 400)
			tester.Frame()
			tester.Frame()
			narrow := width < 300
			if p == Windows {
				m, ok := tester.Find(message)
				if !ok {
					t.Fatalf("Windows: no message: %q", tester.Texts())
				}
				// Across, the message is 14 from the bar's top, at 16; below the title, further down.
				if below := m.Y > 16+24; below != narrow {
					t.Errorf("Windows width %v: the message is at %v", width, m)
				}
				continue
			}
			text, ok := tester.Find("Update available " + message)
			if !ok {
				t.Fatalf("GNOME: no banner text: %q", tester.Texts())
			}
			b, _ := tester.Find(infoBarCloseName(p))
			if below := b.Y >= text.Y+text.H; below != narrow {
				t.Errorf("GNOME width %v: the button %v, the text %v", width, b, text)
			}
		}
	}
}
