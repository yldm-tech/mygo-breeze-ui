package breeze

import (
	"fmt"
	"image"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func init() {
	for _, n := range []string{"Tooltip", "Toasts", "Meter", "Avatar", "Rating"} {
		drawnWidgets[n] = true
	}
}

// TestTooltipKeepsItsBehavior checks the tooltip on every platform: it shows once the pointer rests on its anchor, where the platform puts it (below on macOS and GNOME, above on Windows), goes with Escape, and leaves the anchor clickable.
func TestTooltipKeepsItsBehavior(t *testing.T) {
	for _, p := range Platforms {
		clicks := 0
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Column(c).Padding(120, 40).Children(func() {
					if Tooltip(c, Button(c, "Save"), "Saves the document").Clicked() {
						clicks++
					}
				})
			})
		}, 400, 300)
		tester.Frame()
		b, ok := tester.Find("Save")
		if !ok {
			t.Fatalf("%v: the anchor is missing", p)
		}
		tester.Move(b.X+b.W/2, b.Y+b.H/2)
		tester.Frame()
		if tester.HasText("Saves the document") {
			t.Fatalf("%v: the tip shows before the pointer rested", p)
		}
		tester.Advance(700 * time.Millisecond)
		tester.Frame()
		tip, ok := tester.Find("Saves the document")
		if !ok {
			t.Fatalf("%v: no tip once the pointer rested", p)
		}
		if above := tip.Y+tip.H <= b.Y; above != (p == Windows) {
			t.Errorf("%v: the tip is at %v, the anchor at %v", p, tip, b)
		}
		tester.Key(0, ui.KeyEscape)
		tester.Frame()
		if tester.HasText("Saves the document") {
			t.Errorf("%v: Escape left the tip showing", p)
		}
		if err := tester.Click("Save"); err != nil {
			t.Fatalf("%v: click the anchor: %v", p, err)
		}
		tester.Frame()
		if clicks != 1 {
			t.Errorf("%v: the anchor clicked %d times", p, clicks)
		}
	}
}

// TestToastsKeepTheirBehavior checks the toasts on every platform: Toast shows its message once, in place of MyGo's own look, the action runs and closes the toast, the close button closes one on Windows and GNOME, and a toast goes once its time is over.
func TestToastsKeepTheirBehavior(t *testing.T) {
	for _, p := range Platforms {
		undone := 0
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Column(c).Padding(16).Gap(8).Children(func() {
					if Button(c, "Save").Clicked() {
						c.Toast("Saved")
					}
					if Button(c, "Delete").Clicked() {
						c.ToastAction("Note deleted", "Undo", func() { undone++ })
					}
				})
				Toasts(c)
			})
		}, 600, 400)
		tester.Frame()
		if err := tester.Click("Save"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		n := 0
		for _, s := range tester.Texts() {
			if s == "Saved" {
				n++
			}
		}
		if n != 1 {
			t.Fatalf("%v: the toast shows %d times: %q", p, n, tester.Texts())
		}
		if err := tester.Click("Delete"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		if err := tester.Click("Undo"); err != nil {
			t.Fatalf("%v: click the action: %v", p, err)
		}
		tester.Frame()
		tester.Advance(300 * time.Millisecond)
		tester.Frame()
		if undone != 1 || tester.HasText("Note deleted") {
			t.Fatalf("%v: the action ran %d times, texts %q", p, undone, tester.Texts())
		}
		if p != Mac {
			close := map[Platform]string{Windows: "Close", Linux: "Dismiss"}[p]
			if err := tester.Click(close); err != nil {
				t.Fatalf("%v: click the close button: %v", p, err)
			}
			tester.Frame()
			tester.Advance(300 * time.Millisecond)
			tester.Frame()
			if tester.HasText("Saved") {
				t.Fatalf("%v: the close button left the toast", p)
			}
			if err := tester.Click("Save"); err != nil {
				t.Fatalf("%v: %v", p, err)
			}
			tester.Frame()
		}
		tester.Move(590, 10)
		tester.Advance(5 * time.Second)
		tester.Frame()
		tester.Advance(300 * time.Millisecond)
		tester.Frame()
		if tester.HasText("Saved") {
			t.Fatalf("%v: the toast outlived its time", p)
		}
	}
}

// TestMeterAvatarRatingKeepTheirBehavior checks the indicators on every platform: the meter and the avatar are found by their names at the platform's sizes, and the rating's stars set, clear and step with the keys, from the right in a right-to-left window.
func TestMeterAvatarRatingKeepTheirBehavior(t *testing.T) {
	heights := map[Platform]float32{Mac: 16, Windows: 3, Linux: 8}
	for _, p := range Platforms {
		for _, rtl := range []bool{false, true} {
			stars, changes := 2, 0
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					ui.Column(c).Padding(16).Gap(12).AlignItems(ui.Start).Children(func() {
						Meter(c, 90, 0, 100, &ui.MeterLevels{Warning: 70, Critical: 85}).Width(200).Label("Disk")
						Avatar(c, "Ada Lovelace", nil)
						if Rating(c, &stars, 5).Label("Stars").Changed() {
							changes++
						}
					})
				})
			}, 400, 300)
			if rtl {
				tester.SetDirection(ui.RightToLeft)
			}
			tester.Frame()
			m, ok := tester.Find("Disk")
			if !ok || m.H != heights[p] || m.W != 200 {
				t.Fatalf("%v: meter %v, %v; want 200 by %v", p, m, ok, heights[p])
			}
			a, ok := tester.Find("Ada Lovelace")
			if !ok || a.W-a.H > 0.01 || a.H-a.W > 0.01 || a.W < 20 {
				t.Fatalf("%v: avatar %v, %v", p, a, ok)
			}
			r, ok := tester.Find("Stars")
			if !ok {
				t.Fatalf("%v: rating not found", p)
			}
			// The fourth star from the start of the reading direction.
			star := (r.W + ratingGap(p)) / 5
			x := r.X + star*3.5
			if rtl {
				x = r.X + r.W - star*3.5
			}
			tester.ClickAt(x, r.Y+r.H/2)
			tester.Frame()
			if stars != 4 {
				t.Fatalf("%v rtl %v: a click on the fourth star set %d", p, rtl, stars)
			}
			tester.ClickAt(x, r.Y+r.H/2)
			tester.Frame()
			if stars != 0 {
				t.Fatalf("%v rtl %v: a click on the star set left %d", p, rtl, stars)
			}
			tester.Key(0, ui.KeyEnd)
			tester.Frame()
			tester.Key(0, ui.KeyDown)
			tester.Frame()
			if stars != 4 || changes < 4 {
				t.Fatalf("%v rtl %v: End and Down set %d, %d changes", p, rtl, stars, changes)
			}
		}
	}
}

// ratingGap is the gap between the stars of p's rating, which the stars' row leaves after none.
func ratingGap(p Platform) float32 {
	return map[Platform]float32{Mac: 1, Windows: 8, Linux: 2}[p]
}

func TestAvatarInitialsAndColors(t *testing.T) {
	for name, want := range map[string]string{"Ada Lovelace": "AL", "grace": "G", "  Alan  M. Turing ": "AT", "": ""} {
		if got := avatarInitials(name); got != want {
			t.Errorf("avatarInitials(%q) = %q, want %q", name, got, want)
		}
	}
	for name, want := range map[string]string{"Ada Lovelace": "AL", "grace": "G", "Alan M. Turing": "AT", "  ": ""} {
		if got := adwAvatarInitials(name); got != want {
			t.Errorf("adwAvatarInitials(%q) = %q, want %q", name, got, want)
		}
	}
	// libadwaita picks a random color for an empty name; Breeze UI the first, as its doc says.
	if got := adwAvatarColor(""); got != 0 {
		t.Errorf("adwAvatarColor(\"\") = %d, want 0", got)
	}
	// g_str_hash("a") is 5381*33+97 = 177670, which is color 1 (cyan) of 14.
	if got := adwAvatarColor("a"); got != 177670%14 {
		t.Errorf("adwAvatarColor(\"a\") = %d", got)
	}
	// Bytes past ASCII count as negative, as GLib's signed char does: "é" is C3 A9.
	h := uint32(5381)
	h = h*33 + uint32(0xC3) - 256
	h = h*33 + uint32(0xA9) - 256
	if got := adwAvatarColor("é"); got != int(h%14) {
		t.Errorf("adwAvatarColor(\"é\") = %d, want %d", got, h%14)
	}
}

// TestFindBarKeepsItsBehavior checks the find bar on every platform: its field takes the focus and the query, Enter goes to the next match round the end, Escape closes it.
func TestFindBarKeepsItsBehavior(t *testing.T) {
	for _, p := range Platforms {
		open, query, current := true, "", 0
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Column(c).Fill().Children(func() {
					FindBar(c, &open, &query, 3, &current)
					ui.Text(c, "Body")
				})
			})
		}, 600, 300)
		tester.Frame()
		tester.Frame()
		tester.Type("cat")
		tester.Frame()
		if query != "cat" {
			t.Fatalf("%v: typing gave %q", p, query)
		}
		for range 3 {
			tester.Key(0, ui.KeyEnter)
			tester.Frame()
		}
		if current != 0 {
			t.Fatalf("%v: three Enters over three matches left match %d", p, current)
		}
		tester.Key(0, ui.KeyEnter)
		tester.Frame()
		if current != 1 {
			t.Fatalf("%v: Enter went to match %d", p, current)
		}
		tester.Key(0, ui.KeyEscape)
		tester.Frame()
		if open {
			t.Fatalf("%v: Escape left the find bar open", p)
		}
	}
}

// feedbackView lays out the feedback widgets: a button whose tooltip shows, meters at their three levels, avatars, ratings, a find bar and two toasts.
func feedbackView(p Platform, s *feedbackState) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			if !s.toasted {
				s.toasted = true
				c.AddToast(ui.Toast{ID: "copied", Title: "Copied to the clipboard", Timeout: -1})
				c.AddToast(ui.Toast{ID: "deleted", Title: "Note deleted", Description: "It is in the trash.", Type: "success", Action: "Undo", OnAction: func() {}, Timeout: -1})
			}
			ui.Column(c).Fill().Children(func() {
				FindBar(c, &s.find, &s.query, 3, &s.current)
				ui.Box(c).Grow(1).Children(func() {
					at := func(x, y float32, build func() ui.Element) {
						ui.Box(c).Absolute().Left(x).Top(y).Children(func() { build() })
					}
					at(16, 16, func() ui.Element { return Tooltip(c, Button(c, "Share"), "Share the document") })
					at(300, 16, func() ui.Element { return Meter(c, 40, 0, 100, s.levels).Width(220).Label("Fine") })
					at(300, 44, func() ui.Element { return Meter(c, 75, 0, 100, s.levels).Width(220).Label("Warning") })
					at(300, 72, func() ui.Element { return Meter(c, 95, 0, 100, s.levels).Width(220).Label("Critical") })
					at(16, 110, func() ui.Element { return Avatar(c, "Ada Lovelace", nil) })
					at(64, 110, func() ui.Element { return Avatar(c, "Grace Hopper", nil) })
					at(112, 110, func() ui.Element { return Avatar(c, "Linus", nil) })
					at(300, 110, func() ui.Element { return Rating(c, &s.stars, 5).Label("Stars") })
				})
			})
			Toasts(c)
		})
	}
}

// feedbackState is what feedbackView's widgets edit.
type feedbackState struct {
	find, toasted  bool
	query          string
	current, stars int
	levels         *ui.MeterLevels
}

// TestScreenshotsFeedback writes feedbackView on every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names, with the tooltip showing. Without it the test only builds the frames.
func TestScreenshotsFeedback(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			s := &feedbackState{find: true, query: "cat", stars: 3, levels: &ui.MeterLevels{Warning: 70, Critical: 90}}
			tester := ui.NewTester(feedbackView(p, s), 560, 400)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.Frame()
			b, ok := tester.Find("Share")
			if !ok {
				t.Fatalf("%v: the tooltip's anchor is missing", p)
			}
			tester.Move(b.X+b.W/2, b.Y+b.H/2)
			tester.Advance(700 * time.Millisecond)
			tester.Frame()
			tester.Advance(300 * time.Millisecond)
			tester.Frame()
			for _, s := range []string{"Share the document", "Copied to the clipboard", "Note deleted", "Undo"} {
				if !tester.HasText(s) {
					t.Fatalf("%v: %q is missing: %q", p, s, tester.Texts())
				}
			}
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			writeShot(t, dir, fmt.Sprintf("breeze-feedback-%s-%s.png", strings.ToLower(p.String()), mode), tester)
		}
	}
}

// TestTooltipFollowsTheFocusAndDisabled checks on every platform what the pointer test does not: Tab to the anchor shows its tip at once, a disabled anchor shows none however long the pointer rests on it, text becomes the anchor's Description, and an empty anchor or text changes nothing.
func TestTooltipFollowsTheFocusAndDisabled(t *testing.T) {
	for _, p := range Platforms {
		off := false
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Column(c).Padding(120, 40).Gap(40).AlignItems(ui.Start).Children(func() {
					Tooltip(c, Button(c, "Save"), "Saves the document")
					Tooltip(c, Button(c, "Print"), "Prints the document").Disabled(true)
					if Tooltip(c, ui.Element{}, "Nothing").Valid() || Tooltip(c, Button(c, "Plain"), "").Clicked() {
						off = true
					}
				})
			})
		}, 400, 400)
		tester.Frame()
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		if !tester.Focused("Save") {
			t.Fatalf("%v: Tab did not focus the anchor", p)
		}
		if !tester.HasText("Saves the document") {
			t.Fatalf("%v: the focus from the keyboard showed no tip: %q", p, tester.Texts())
		}
		tester.Key(0, ui.KeyEscape)
		tester.Frame()
		if tester.HasText("Saves the document") {
			t.Errorf("%v: Escape left the focused anchor's tip", p)
		}
		b, ok := tester.Find("Print")
		if !ok {
			t.Fatalf("%v: the disabled anchor is missing", p)
		}
		tester.Move(b.X+b.W/2, b.Y+b.H/2)
		tester.Frame()
		tester.Advance(time.Second)
		tester.Frame()
		if tester.HasText("Prints the document") {
			t.Errorf("%v: a disabled anchor showed its tip", p)
		}
		if off || tester.HasText("Nothing") {
			t.Errorf("%v: an empty anchor or an empty text did something", p)
		}
	}
}

// TestToastEscapeAndTypes checks on every platform that every Type, known or not, shows the toast's title and description, and that Escape closes the toast the keyboard focus is in without running its action.
func TestToastEscapeAndTypes(t *testing.T) {
	for _, p := range Platforms {
		for _, typ := range []string{"", "info", "success", "warning", "error", "danger", "other"} {
			added := false
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					if !added {
						added = true
						c.AddToast(ui.Toast{Title: "Uploaded", Description: "3 files", Type: typ, Timeout: -1})
					}
					Toasts(c)
				})
			}, 600, 400)
			tester.Frame()
			tester.Frame()
			if !tester.HasText("Uploaded") || !tester.HasText("3 files") {
				t.Fatalf("%v, type %q: texts %q", p, typ, tester.Texts())
			}
		}
		ran, closed, added := 0, 0, false
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				if !added {
					added = true
					c.AddToast(ui.Toast{Title: "Note deleted", Action: "Undo", OnAction: func() { ran++ }, OnClose: func() { closed++ }, Timeout: -1})
				}
				Toasts(c)
			})
		}, 600, 400)
		tester.Frame()
		tester.Frame()
		for range 4 {
			if tester.Focused("Undo") {
				break
			}
			tester.Key(0, ui.KeyTab)
			tester.Frame()
		}
		if !tester.Focused("Undo") {
			t.Fatalf("%v: Tab does not reach the toast's action", p)
		}
		tester.Key(0, ui.KeyEscape)
		tester.Frame()
		tester.Advance(300 * time.Millisecond)
		tester.Frame()
		if tester.HasText("Note deleted") || ran != 0 || closed != 1 {
			t.Fatalf("%v: after Escape the action ran %d times, OnClose %d, texts %q", p, ran, closed, tester.Texts())
		}
	}
}

// TestMeterLevelsAsMyGo checks meterLevel against ui.Meter's rule: none without levels or with equal ones, the critical level winning over the warning one, and a Critical below Warning making low values the bad ones.
func TestMeterLevelsAsMyGo(t *testing.T) {
	high := &ui.MeterLevels{Warning: 70, Critical: 90}
	low := &ui.MeterLevels{Warning: 20, Critical: 10}
	for _, tc := range []struct {
		value  float64
		levels *ui.MeterLevels
		want   int
	}{
		{50, nil, 0}, {95, &ui.MeterLevels{}, 0}, {95, &ui.MeterLevels{Warning: 80, Critical: 80}, 0},
		{69, high, 0}, {70, high, 1}, {89, high, 1}, {90, high, 2}, {100, high, 2},
		{50, low, 0}, {20, low, 1}, {11, low, 1}, {10, low, 2}, {0, low, 2},
	} {
		if got := meterLevel(tc.value, tc.levels); got != tc.want {
			t.Errorf("meterLevel(%v, %+v) = %d, want %d", tc.value, tc.levels, got, tc.want)
		}
	}
	for _, tc := range []struct {
		value, lo, hi float64
		want          float32
	}{{50, 0, 100, 0.5}, {-5, 0, 100, 0}, {150, 0, 100, 1}, {5, 10, 10, 0}, {5, 10, 0, 0}} {
		if got := meterFraction(tc.value, tc.lo, tc.hi); got != tc.want {
			t.Errorf("meterFraction(%v, %v, %v) = %v, want %v", tc.value, tc.lo, tc.hi, got, tc.want)
		}
	}
}

// TestRatingDisabled checks on every platform that a disabled rating takes neither clicks nor keys, and on Windows, where the stars grow toward the pointer, that the pointer changes the stars of an enabled rating but not of one in a disabled column.
func TestRatingDisabled(t *testing.T) {
	for _, p := range Platforms {
		for _, inColumn := range []bool{false, true} {
			stars := 2
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					col := ui.Column(c).Padding(16).AlignItems(ui.Start)
					if inColumn {
						col.Disabled(true)
					}
					col.Children(func() {
						r := Rating(c, &stars, 5).Label("Stars")
						if !inColumn {
							r.Disabled(true)
						}
					})
				})
			}, 300, 120)
			tester.Frame()
			r, ok := tester.Find("Stars")
			if !ok {
				t.Fatalf("%v: rating not found", p)
			}
			tester.ClickAt(r.X+r.W-2, r.Y+r.H/2)
			tester.Frame()
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			tester.Key(0, ui.KeyEnd)
			tester.Frame()
			if stars != 2 {
				t.Errorf("%v, disabled column %v: a disabled rating changed to %d", p, inColumn, stars)
			}
		}
	}
	// The pointer over an enabled rating, then a disabled one, on Windows.
	for _, disabled := range []bool{false, true} {
		stars := 1
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Windows, func() {
				Use(c)
				ui.Column(c).Padding(16).AlignItems(ui.Start).Disabled(disabled).Children(func() {
					Rating(c, &stars, 5).Label("Stars")
				})
			})
		}, 300, 120)
		tester.Frame()
		r, _ := tester.Find("Stars")
		still := append([]uint8(nil), tester.Image().Pix...)
		tester.Move(r.X+r.W*0.7, r.Y+r.H/2)
		tester.Frame()
		tester.Frame()
		moved := string(tester.Image().Pix) != string(still)
		if moved == disabled {
			t.Errorf("Windows, disabled %v: the pointer changed the stars: %v", disabled, moved)
		}
	}
}

// TestFindBarButtons checks on every platform the find bar's buttons and its closed state: Next and Previous step round the matches, Done closes the bar, which then shows nothing, and opening it again gives its field the focus.
func TestFindBarButtons(t *testing.T) {
	for _, p := range Platforms {
		open, query, current := true, "cat", 0
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Column(c).Fill().Children(func() {
					FindBar(c, &open, &query, 3, &current)
					ui.Text(c, "Body")
				})
			})
		}, 600, 300)
		tester.Frame()
		tester.Frame()
		if err := tester.Click("Previous"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		if current != 2 {
			t.Fatalf("%v: Previous from the first match went to %d", p, current)
		}
		if err := tester.Click("Next"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		if current != 0 {
			t.Fatalf("%v: Next from the last match went to %d", p, current)
		}
		if err := tester.Click("Done"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		if open || tester.HasText("Done") || !tester.HasText("Body") {
			t.Fatalf("%v: Done left open %v, texts %q", p, open, tester.Texts())
		}
		open = true
		tester.Frame()
		tester.Frame()
		tester.Type("s")
		tester.Frame()
		if query != "s" {
			t.Fatalf("%v: typing after opening again gave %q", p, query)
		}
	}
}

// TestAvatarWithAnImage checks on every platform that an avatar with an image is found by its name at the size of one without, and that an empty name builds.
func TestAvatarWithAnImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	bm := ui.NewBitmap(img)
	for _, p := range Platforms {
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Row(c).Padding(16).Gap(8).Children(func() {
					Avatar(c, "Ada Lovelace", bm)
					Avatar(c, "Grace Hopper", nil)
					Avatar(c, "", nil)
				})
			})
		}, 300, 100)
		tester.Frame()
		a, ok := tester.Find("Ada Lovelace")
		g, ok2 := tester.Find("Grace Hopper")
		if !ok || !ok2 || a.W != g.W || a.H != g.H || a.W != a.H {
			t.Fatalf("%v: with an image %v, %v; without %v, %v", p, a, ok, g, ok2)
		}
	}
}
