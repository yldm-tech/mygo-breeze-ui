package breeze

import (
	"fmt"
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
		if got := adwInitials(name); got != want {
			t.Errorf("adwInitials(%q) = %q, want %q", name, got, want)
		}
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
