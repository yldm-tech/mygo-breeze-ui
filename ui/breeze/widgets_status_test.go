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
	for _, n := range []string{"InfoBar", "InfoBarActions", "ProgressRing", "StatusPage"} {
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
		// The close button's distance from the bar's end, left to right, which right to left must be its distance from the bar's start: the bar is 16 from the window's edges.
		var endGap float32
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
			// The close button is at the end of the reading direction, as far from it in either.
			if atEnd := b.X+b.W/2 > 300; atEnd == rtl {
				t.Errorf("%v rtl %v: the close button is at %v", p, rtl, b)
			}
			if !rtl {
				endGap = 600 - 16 - (b.X + b.W)
			} else if gap := b.X - 16; math.Abs(float64(gap-endGap)) > 0.01 {
				t.Errorf("%v: the close button is %v from the bar's start right to left, %v from its end left to right", p, gap, endGap)
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
		if err := tester.Click(infoBarCloseName(p)); err != nil {
			t.Fatalf("%v: a disabled bar has no close button: %v", p, err)
		}
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

// statusTitlesView lays out statusView's four message bars with their titles and no message, where tools/winui-reference places WinUI's InfoBars without a message (winui-status-titles).
func statusTitlesView(p Platform, open *[4]bool) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Box(c).Fill().Children(func() {
				for i, kind := range []string{"info", "success", "warning", "error"} {
					var o *bool
					if i == 0 || i == 3 {
						o = &open[i]
					}
					title := []string{"Informational", "Success", "Warning", "Error"}[i]
					ui.Column(c).Absolute().Left(16).Top(16 + float32(i)*68).Right(16).AlignItems(ui.Start).Children(func() { InfoBar(c, kind, title, "", o) })
				}
			})
		})
	}
}

// statusActionsView lays out message bars with actions where tools/appkit-reference's statusActionsSet places SwiftUI's: 16 from the window's edges and 68 apart, a closable warning with a button, an informational bar with a primary button, a closable success bar without a message with a button, and a closable error bar 280 wide with a button, which Windows lays out vertically and GNOME with its button below the text.
func statusActionsView(p Platform, open *[4]bool) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Box(c).Fill().Children(func() {
				bars := []struct {
					kind, title, message, action string
					primary, closable            bool
					width                        float32
				}{
					{"warning", "Warning", "Some message about it.", "Open Settings", false, true, 528},
					{"info", "Informational", "Some message about it.", "Restart", true, false, 528},
					{"success", "Success", "", "Undo", false, true, 528},
					{"error", "Error", "Some message about it.", "Retry", false, true, 280},
				}
				for i, b := range bars {
					var o *bool
					if b.closable {
						o = &open[i]
					}
					ui.Column(c).Absolute().Left(16).Top(16 + float32(i)*68).Width(b.width).AlignItems(ui.Start).Children(func() {
						InfoBarActions(c, b.kind, b.title, b.message, o, func() {
							if b.primary {
								PrimaryButton(c, b.action)
							} else {
								Button(c, b.action)
							}
						})
					})
				}
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

// TestScreenshotsStatus writes statusView, statusTitlesView, statusActionsView and statusPageView on every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names. Without it the test only builds the frames.
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
			titles := ui.NewTester(statusTitlesView(p, &open), 560, 300)
			titles.SetScale(2)
			titles.SetDark(dark)
			titles.Frame()
			titles.Frame()
			for _, s := range []string{"Informational", "Success", "Warning", "Error"} {
				if _, ok := titles.Find(s); !ok {
					t.Fatalf("%v: the bar %q without a message is missing: %q", p, s, titles.Texts())
				}
			}
			writeShot(t, dir, fmt.Sprintf("breeze-statustitles-%s-%s.png", strings.ToLower(p.String()), mode), titles)
			acts := ui.NewTester(statusActionsView(p, &open), 560, 380)
			acts.SetScale(2)
			acts.SetDark(dark)
			acts.Frame()
			acts.Frame()
			for _, s := range []string{"Warning", "Open Settings", "Restart", "Undo", "Retry"} {
				if _, ok := acts.Find(s); !ok {
					t.Fatalf("%v: %q of the bars with actions is missing: %q", p, s, acts.Texts())
				}
			}
			writeShot(t, dir, fmt.Sprintf("breeze-statusactions-%s-%s.png", strings.ToLower(p.String()), mode), acts)
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

// TestFluentInfoBarTitleOnly checks that a WinUI InfoBar without a message lays its title out as InfoBarPanel lays a single child, vertically, the panel padded 14 above and 18 below, where a bar whose title and message share a line is the 48 DIPs of InfoBarMinHeight inside its 1-DIP border.
func TestFluentInfoBarTitleOnly(t *testing.T) {
	tester := ui.NewTester(func(c *ui.Context) {
		Preview(c, Windows, func() {
			Use(c)
			ui.Column(c).Padding(16).Gap(16).Children(func() {
				InfoBar(c, "info", "Alone", "", nil)
				InfoBar(c, "info", "Paired", "Message", nil)
			})
		})
	}, 600, 300)
	tester.Frame()
	tester.Frame()
	// Find gives the bars, which their titles name, and the message's line.
	alone, _ := tester.Find("Alone")
	paired, _ := tester.Find("Paired")
	line, ok := tester.Find("Message")
	if !ok {
		t.Fatalf("no message: %q", tester.Texts())
	}
	if paired.H != 50 {
		t.Errorf("a bar with a title and a message on one line is %v high, want 50", paired.H)
	}
	if want := max(1+14+line.H+18+1, 50); alone.H != want {
		t.Errorf("a bar without a message is %v high, want %v", alone.H, want)
	}
}

// TestInfoBarArrangements checks the arrangements Windows and GNOME switch between by the bar's width: in a wide bar WinUI's message follows the title on its line and AdwBanner's button sits beside its text; in a narrow one the message goes below the title, and the banner's button below its text. The bar is narrower than the window, so the first frame, which knows only the window's width, must be laid out again at the bar's. Scaled by a transform, the bar keeps the arrangement of its width in the layout, though it is drawn twice as wide.
func TestInfoBarArrangements(t *testing.T) {
	const message = "Restart to finish installing."
	for _, p := range []Platform{Windows, Linux} {
		for _, scale := range []float32{1, 2} {
			for _, width := range []float32{560, 260} {
				open := true
				tester := ui.NewTester(func(c *ui.Context) {
					Preview(c, p, func() {
						Use(c)
						// The scale is about the column's middle: the space before it keeps it in the window.
						ui.Column(c).Padding(150, 0, 0, 300).Children(func() {
							ui.Column(c).Padding(16).Width(width+32).Scale(scale, scale).Children(func() {
								InfoBar(c, "info", "Update available", message, &open)
							})
						})
					})
				}, 1300, 500)
				tester.Frame()
				tester.Frame()
				narrow := width < 300
				if p == Windows {
					m, ok := tester.Find(message)
					if !ok {
						t.Fatalf("Windows: no message: %q", tester.Texts())
					}
					// Across, the message is 1+14 below the bar's top, which Find gives as the bar is named by its title; below the title, further down. Both are scaled.
					bar, _ := tester.Find("Update available")
					if below := m.Y-bar.Y > 25*scale; below != narrow {
						t.Errorf("Windows scale %v width %v: the message is at %v in the bar %v", scale, width, m, bar)
					}
					continue
				}
				text, ok := tester.Find("Update available " + message)
				if !ok {
					t.Fatalf("GNOME: no banner text: %q", tester.Texts())
				}
				b, _ := tester.Find(infoBarCloseName(p))
				if below := b.Y >= text.Y+text.H; below != narrow {
					t.Errorf("GNOME scale %v width %v: the button %v, the text %v", scale, width, b, text)
				}
			}
		}
	}
}

// indexOf returns the index of s in texts, -1 when it is not there.
func indexOf(texts []string, s string) int {
	for i, x := range texts {
		if x == s {
			return i
		}
	}
	return -1
}

// TestInfoBarActionsKeepTheirBehavior checks a bar with actions on every platform, left to right and right to left: the action is inside the bar and takes clicks and the keyboard, assistive technology reads the title, the message, the action and then the close button, the action is before the close button in the reading direction, and on GNOME the action takes the place of the banner's Dismiss button; a bar without open shows its action too.
func TestInfoBarActionsKeepTheirBehavior(t *testing.T) {
	const title, message = "Model missing", "Download a model to transcribe."
	for _, p := range Platforms {
		for _, rtl := range []bool{false, true} {
			open, opened, restarted := true, 0, 0
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					ui.Column(c).Padding(16).Gap(12).Children(func() {
						InfoBarActions(c, "warning", title, message, &open, func() {
							if Button(c, "Open Settings").Clicked() {
								opened++
							}
						})
						InfoBarActions(c, "info", "Update ready", "", nil, func() {
							if PrimaryButton(c, "Restart").Clicked() {
								restarted++
							}
						})
					})
				})
			}, 640, 360)
			if rtl {
				tester.SetDirection(ui.RightToLeft)
			}
			tester.Frame()
			tester.Frame()
			bar, _ := tester.Find(title)
			b, ok := tester.Find("Open Settings")
			if !ok {
				t.Fatalf("%v rtl %v: no action: %q", p, rtl, tester.Texts())
			}
			if b.X < bar.X || b.Y < bar.Y || b.X+b.W > bar.X+bar.W || b.Y+b.H > bar.Y+bar.H {
				t.Errorf("%v rtl %v: the action %v is outside the bar %v", p, rtl, b, bar)
			}
			if err := tester.Click("Open Settings"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			if opened != 1 || !open {
				t.Fatalf("%v rtl %v: a click on the action ran it %d times, open %v", p, rtl, opened, open)
			}
			if err := tester.Click("Restart"); err != nil || restarted != 1 {
				t.Fatalf("%v rtl %v: the action of a bar that cannot close: %v, ran %d times", p, rtl, err, restarted)
			}
			texts := tester.Texts()
			// GNOME's banner shows the message in its title's label.
			order := []string{title, message, "Open Settings", "Close"}
			if p == Linux {
				order = []string{title + " " + message, "Open Settings"}
			}
			for i := 1; i < len(order); i++ {
				if a, z := indexOf(texts, order[i-1]), indexOf(texts, order[i]); a < 0 || z < a {
					t.Errorf("%v rtl %v: %q is not read after %q: %q", p, rtl, order[i], order[i-1], texts)
				}
			}
			if p == Linux {
				if _, ok := tester.Find("Dismiss"); ok {
					t.Errorf("GNOME rtl %v: the banner shows its Dismiss button beside its action", rtl)
				}
				continue
			}
			x, ok := tester.Find("Close")
			if !ok {
				t.Fatalf("%v rtl %v: no close button: %q", p, rtl, texts)
			}
			if before := b.X+b.W <= x.X; before == rtl {
				t.Errorf("%v rtl %v: the action %v is not before the close button %v in the reading direction", p, rtl, b, x)
			}
			// The keyboard: Tab reaches the action, Space presses it, and the next Tab the close button.
			for range 6 {
				if tester.Focused("Open Settings") {
					break
				}
				tester.Key(0, ui.KeyTab)
				tester.Frame()
			}
			if !tester.Focused("Open Settings") {
				t.Fatalf("%v rtl %v: Tab does not reach the action", p, rtl)
			}
			tester.Key(0, ui.KeySpace)
			tester.Frame()
			if opened != 2 {
				t.Errorf("%v rtl %v: Space on the action ran it %d times in all", p, rtl, opened)
			}
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			if !tester.Focused("Close") {
				t.Errorf("%v rtl %v: Tab from the action does not reach the close button", p, rtl)
			}
		}
	}
}

// TestInfoBarActionsWithoutActions checks that InfoBar is InfoBarActions with nil actions: the same frame, pixel for pixel, on every platform.
func TestInfoBarActionsWithoutActions(t *testing.T) {
	for _, p := range Platforms {
		frame := func(bar func(c *ui.Context)) string {
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					ui.Column(c).Padding(16).Children(func() { bar(c) })
				})
			}, 500, 160)
			tester.Frame()
			tester.Frame()
			return string(tester.Image().Pix)
		}
		open := true
		a := frame(func(c *ui.Context) { InfoBar(c, "warning", "Title", "Message", &open) })
		b := frame(func(c *ui.Context) { InfoBarActions(c, "warning", "Title", "Message", &open, nil) })
		if a != b {
			t.Errorf("%v: InfoBar and InfoBarActions without actions draw different frames", p)
		}
	}
}

// TestMacInfoBarActionsAsMeasured checks the actions of a bar on macOS against SwiftUI's buttons in that row of a grouped form, as tools/appkit-reference's statusActionsSet renders them on macOS 27.2, light and dark: a 24-point push button 10 points below the bar's top, the bar 4 points taller than without it as the title's line moves down to the button's middle, the close button's middle on the button's, 6.25 points between the button and the close button's box, the button 10 points from the bar's end without a close button, and the button's face laid over the bar's fill (228 over 247, 54 over 37).
func TestMacInfoBarActionsAsMeasured(t *testing.T) {
	for _, dark := range []bool{false, true} {
		open := true
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Mac, func() {
				Use(c)
				ui.Column(c).Padding(16).Gap(16).Children(func() {
					InfoBar(c, "warning", "Plain", "Some message about it.", nil)
					// Labels other than the buttons' text let Find give the buttons rather than their text.
					InfoBarActions(c, "warning", "Warning", "Some message about it.", &open, func() { Button(c, "Open Settings").Label("Settings button") })
					InfoBarActions(c, "info", "Informational", "Some message about it.", nil, func() { Button(c, "Restart").Label("Restart button") })
				})
			})
		}, 560, 300)
		tester.SetDark(dark)
		tester.Frame()
		tester.Frame()
		plain, _ := tester.Find("Plain")
		bar, _ := tester.Find("Warning")
		b, _ := tester.Find("Settings button")
		x, _ := tester.Find("Close")
		if b.H != 24 || b.Y-bar.Y != 10 {
			t.Errorf("dark %v: the button %v in the bar %v, want 24 high and 10 below its top", dark, b, bar)
		}
		// Half of what the button is taller than the title's line: 4 points with SF's line of 16, a little more with MyGo's shorter one.
		if d := bar.H - plain.H; d < 4 || d > 5 {
			t.Errorf("dark %v: the bar with a button is %v taller than without, want about 4", dark, d)
		}
		if mid, xmid := b.Y+b.H/2, x.Y+x.H/2; math.Abs(float64(mid-xmid)) > 0.01 {
			t.Errorf("dark %v: the close button's middle %v, the button's %v", dark, xmid, mid)
		}
		if gap := x.X - (b.X + b.W); math.Abs(float64(gap-6.25)) > 0.01 {
			t.Errorf("dark %v: %v between the button and the close button, want 6.25", dark, gap)
		}
		info, _ := tester.Find("Informational")
		r, _ := tester.Find("Restart button")
		if end := info.X + info.W - (r.X + r.W); math.Abs(float64(end-10)) > 0.01 {
			t.Errorf("dark %v: the button is %v from the end of a bar without a close button, want 10", dark, end)
		}
		img := tester.Image()
		face := img.RGBAAt(int(b.X+3), int(b.Y+b.H/2))
		fill := img.RGBAAt(int(bar.X+bar.W/2), int(bar.Y+2))
		levels := map[bool][2]uint8{false: {228, 247}, true: {54, 37}}[dark]
		want, under := levels[0], levels[1]
		if !nearColor(face, color.RGBA{want, want, want, 255}, 1) || !nearColor(fill, color.RGBA{under, under, under, 255}, 1) {
			t.Errorf("dark %v: the button's face %v over the fill %v, want %d over %d", dark, face, fill, want, under)
		}
	}
}

// TestFluentInfoBarActions checks WinUI InfoBar's ActionButton: across a wide bar 16 DIPs after the message and 8 from the panel's top, the bar keeping its 48 DIPs inside the border; in a narrow bar 12 below the message at its start, the panel padded 18 below it; laid out vertically too when the action with its margin is taller than InfoBarMinHeight, and across beside a title without a message, as InfoBarPanel then has two items; and the button's face laid over the bar's color, as fluentOver lays it.
func TestFluentInfoBarActions(t *testing.T) {
	const message = "Some message about it."
	for _, width := range []float32{600, 300} {
		var theme *ui.Theme
		var accent ui.Color
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Windows, func() {
				Use(c)
				theme, accent = c.Theme(), fluentSystemAccent(c)
				ui.Column(c).Padding(16).Gap(16).Width(width + 32).Children(func() {
					InfoBarActions(c, "success", "Saved", message, nil, func() { Button(c, "Undo").Label("Undo button") })
					InfoBarActions(c, "info", "Tall", message, nil, func() { Button(c, "Tall action").MinHeight(44) })
					InfoBarActions(c, "info", "Alone", "", nil, func() { Button(c, "Retry").Label("Retry button") })
				})
			})
		}, 700, 600)
		tester.Frame()
		tester.Frame()
		bar, _ := tester.Find("Saved")
		m, _ := tester.Find(message)
		b, ok := tester.Find("Undo button")
		if !ok {
			t.Fatalf("width %v: no action: %q", width, tester.Texts())
		}
		if width > 500 {
			if bar.H != 50 || b.Y != bar.Y+1+8 || math.Abs(float64(b.X-(m.X+m.W+16))) > 0.01 || b.H != 32 {
				t.Errorf("width %v: across, the action %v after the message %v in the bar %v", width, b, m, bar)
			}
			tall, _ := tester.Find("Tall")
			if ta, _ := tester.Find("Tall action"); ta.Y < tall.Y+30 {
				t.Errorf("an action taller than InfoBarMinHeight less its margin is at %v in the bar %v, not below the message", ta, tall)
			}
			alone, _ := tester.Find("Alone")
			if r, _ := tester.Find("Retry button"); r.Y != alone.Y+1+8 || alone.H != 50 {
				t.Errorf("a title and an action are not across: the action %v in the bar %v", r, alone)
			}
		} else {
			if b.Y != m.Y+m.H+12 || b.X != m.X {
				t.Errorf("width %v: the action %v is not 12 below the message %v at its start", width, b, m)
			}
			if end := bar.Y + bar.H - (b.Y + b.H); end != 18+1 {
				t.Errorf("width %v: the panel ends %v below the action, want 18 and the border", width, end)
			}
		}
		bg, _ := fluentToastColors(theme, accent, fluentToastSuccess)
		want := fluentOver(theme, bg).Surface
		got := tester.Image().RGBAAt(int(b.X+3), int(b.Y+b.H/2))
		if !nearColor(got, color.RGBA{want.R, want.G, want.B, 255}, 1) {
			t.Errorf("width %v: the action's face is %v, want %v over the bar's %v", width, got, want, bg)
		}
	}
}

// TestAdwBannerActions checks a banner with actions on GNOME: the action takes the Dismiss button's place, beside the text in a wide banner and below it in a narrow one, as AdwBanner places its button; a PrimaryButton shows as the suggested action, in the accent; and the button's face is the Dismiss button's, laid over the strip.
func TestAdwBannerActions(t *testing.T) {
	const message = "Some message about it."
	for _, width := range []float32{600, 260} {
		open := true
		var theme *ui.Theme
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Linux, func() {
				Use(c)
				theme = c.Theme()
				ui.Column(c).Padding(16).Gap(16).Width(width + 32).Children(func() {
					InfoBarActions(c, "", "Offline", message, &open, func() { Button(c, "Retry").Label("Retry button") })
					InfoBarActions(c, "", "Update", message, nil, func() { PrimaryButton(c, "Restart").Label("Restart button") })
				})
			})
		}, 700, 500)
		tester.Frame()
		tester.Frame()
		if _, ok := tester.Find("Dismiss"); ok {
			t.Fatalf("width %v: the banner shows Dismiss beside its action", width)
		}
		text, _ := tester.Find("Offline " + message)
		b, ok := tester.Find("Retry button")
		if !ok {
			t.Fatalf("width %v: no action: %q", width, tester.Texts())
		}
		if below := b.Y >= text.Y+text.H; below != (width < 300) {
			t.Errorf("width %v: the action %v, the text %v", width, b, text)
		}
		// Beside a narrow action, narrower than the Dismiss button's 84, the title stays centered across the whole strip, as AdwBanner keeps it.
		if bar, _ := tester.Find("Offline"); width > 300 && math.Abs(float64(text.X+text.W/2-(bar.X+bar.W/2))) > 0.5 {
			t.Errorf("width %v: the text %v is not centered across the banner %v beside the action %v", width, text, bar, b)
		}
		img := tester.Image()
		bg := theme.Background.Mix(ui.Hex("#7d7d83"), 0.3)
		want := bg.Mix(pick(theme.Dark, adwInk, white), adwBannerShares[0][0])
		if got := img.RGBAAt(int(b.X+4), int(b.Y+b.H/2)); !nearColor(got, color.RGBA{want.R, want.G, want.B, 255}, 1) {
			t.Errorf("width %v: the action's face is %v, want the Dismiss button's %v", width, got, want)
		}
		r, _ := tester.Find("Restart button")
		if got, acc := img.RGBAAt(int(r.X+4), int(r.Y+r.H/2)), theme.Accent; !nearColor(got, color.RGBA{acc.R, acc.G, acc.B, 255}, 1) {
			t.Errorf("width %v: the primary action's face is %v, want the accent %v", width, got, acc)
		}
	}
}
