package breeze

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func init() {
	drawnWidgets["HeaderBar"] = true
}

// headerTestBackIcon is a chevron pointing back, for the image button at the start of the test bars.
var headerTestBackIcon = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><path d="M10.5 2.5 5 8l5.5 5.5" fill="none" stroke="#000" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>`))

// headerTestRoom is the room the window controls take in a window with a hidden title bar on p, as MyGo reports it (c.TitleBar), left to right or mirrored right to left: on macOS the traffic lights moved to 19, 19 (mygo.WindowOptions.TrafficLightPosition), as AppKit places them in a unified toolbar, the zoom button ending 79 points in and the band 52 high; on Windows MyGo's three 46-DIP caption buttons in a 48-DIP title bar (TitleBarHeight 48); on GNOME GTK's close button in a GTK 3 header bar of its own, padded 6 either side of the 34-pixel button and 46 high (inferred from MyGo's title bar; the end controls it leaves end at 514 in a 560-pixel window, where libadwaita 1.10.0 ends them, 6 before its close button at 520 to 554, in tools/adwaita-reference's headerbar set).
func headerTestRoom(p Platform, dir ui.Direction) ui.TitleBar {
	var bar ui.TitleBar
	switch p {
	case Windows:
		bar = ui.TitleBar{Height: 48, Right: 138}
	case Linux:
		bar = ui.TitleBar{Height: 46, Right: 46}
	default:
		bar = ui.TitleBar{Height: 52, Left: 79}
	}
	if dir == ui.RightToLeft {
		bar.Left, bar.Right = bar.Right, bar.Left
	}
	return bar
}

// headerTestState is what the test bar's controls change.
type headerTestState struct {
	back, edit, share int
	subtitle          string
	noStart, noEnd    bool
	shareOff          bool
	name              string // a label given after HeaderBar, naming the bar instead of its title
	editName          string // a label given to the Edit button, to find the button rather than its text
}

// headerView builds a window with a header bar at its top on platform p: the title Documents, s's subtitle, an image button Back at its start and a text button Edit and an image button Share at its end, then the window's content.
func headerView(p Platform, s *headerTestState) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Column(c).Fill().Children(func() {
				start := func() {
					if IconButton(c, headerTestBackIcon, "Back").Clicked() {
						s.back++
					}
				}
				end := func() {
					edit := Button(c, "Edit")
					if s.editName != "" {
						edit.Label(s.editName)
					}
					if edit.Clicked() {
						s.edit++
					}
					if IconButton(c, idiomTestShareIcon, "Share").Disabled(s.shareOff).Clicked() {
						s.share++
					}
				}
				if s.noStart {
					start = nil
				}
				if s.noEnd {
					end = nil
				}
				bar := HeaderBar(c, "Documents", s.subtitle, start, end)
				if s.name != "" {
					bar.Label(s.name)
				}
				ui.Column(c).Grow(1).Padding(16).Children(func() { ui.Text(c, "Content") })
			})
		})
	}
}

// headerTester builds headerView on p in a w by h window whose controls take headerTestRoom, laid out in dir.
func headerTester(p Platform, dir ui.Direction, w, h int, s *headerTestState) *ui.Tester {
	tester := layoutTester(p, w, h, headerView(p, s))
	tester.SetDirection(dir)
	tester.SetTitleBar(headerTestRoom(p, dir))
	tester.Frame()
	return tester
}

// headerFind returns the box of the element showing or named s, failing the test without one.
func headerFind(t *testing.T, tester *ui.Tester, s string) ui.Rect {
	t.Helper()
	r, ok := tester.Find(s)
	if !ok {
		t.Fatalf("no element shows %q", s)
	}
	return r
}

// overlaps reports whether a and b share any area.
func overlaps(a, b ui.Rect) bool {
	return a.X < b.X+b.W && b.X < a.X+a.W && a.Y < b.Y+b.H && b.Y < a.Y+a.H
}

// TestHeaderBarKeepsClearOfWindowControls checks the bar on every platform, left to right and right to left, with the window controls' room on their side (the left on macOS, the right on Windows and GNOME, mirrored right to left): it runs along the top of the window, as wide as it and at least as high as the room and the platform's bar, nothing of it overlaps the room, and the title is placed as the platform places it: 12 points after the start controls on macOS, 14 DIPs after the LeftHeader on Windows, which holds the image button 2 DIPs in (the back button's Margin), so 16 after the button, centered across the bar on GNOME, with the end controls at the far end.
func TestHeaderBarKeepsClearOfWindowControls(t *testing.T) {
	const w, h = 560, 160
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			s := &headerTestState{name: "Window"}
			tester := headerTester(p, dir, w, h, s)
			room := headerTestRoom(p, dir)
			bar := headerFind(t, tester, "Window")
			least := map[Platform]float32{Mac: 52, Windows: 48, Linux: 47}[p]
			if bar.X != 0 || bar.Y != 0 || bar.W != w || bar.H < room.Height || bar.H < least {
				t.Errorf("%v %v: the bar is %+v, want the window's width along its top, at least %v and %v high", p, dir, bar, room.Height, least)
			}
			left := ui.Rect{X: 0, Y: 0, W: room.Left, H: room.Height}
			right := ui.Rect{X: w - room.Right, Y: 0, W: room.Right, H: room.Height}
			parts := map[string]ui.Rect{}
			for _, name := range []string{"Back", "Documents", "Edit", "Share"} {
				r := headerFind(t, tester, name)
				parts[name] = r
				if overlaps(r, left) || overlaps(r, right) {
					t.Errorf("%v %v: %s at %+v overlaps the window controls' room %+v %+v", p, dir, name, r, left, right)
				}
				if r.Y < bar.Y || r.Y+r.H > bar.Y+bar.H {
					t.Errorf("%v %v: %s at %+v is not inside the bar %+v", p, dir, name, r, bar)
				}
			}
			back, title, edit, share := parts["Back"], parts["Documents"], parts["Edit"], parts["Share"]
			// after is how far a's start is past b's end along the line: from b's right edge left to right, from b's left edge right to left.
			after := func(a, b ui.Rect) float32 {
				if dir == ui.RightToLeft {
					return b.X - (a.X + a.W)
				}
				return a.X - (b.X + b.W)
			}
			switch p {
			case Linux:
				if mid := title.X + title.W/2; mid < w/2-1 || mid > w/2+1 {
					t.Errorf("%v %v: the title's middle is at %v, not centered across the %v-pixel bar", p, dir, mid, w)
				}
			default:
				want := map[Platform]float32{Mac: 12, Windows: 16}[p]
				if got := after(title, back); got < want-0.5 || got > want+0.5 {
					t.Errorf("%v %v: the title starts %v after the start controls, want %v", p, dir, got, want)
				}
			}
			// The start controls begin where the platform begins them: 17 points after the traffic lights on macOS, 2 DIPs in on Windows (the LeftPaddingColumn is as wide as the caption buttons at the start, none here, and the back button keeps its Margin of 2), 6 pixels in on GNOME (the bar's padding of 7 less its margin of -1px, the back button at 6 in libadwaita 1.10.0's render, tools/adwaita-reference's headerbar set).
			begin := back.X
			if dir == ui.RightToLeft {
				begin = w - (back.X + back.W)
			}
			if wantBegin := map[Platform]float32{Mac: 79 + 17, Windows: 2, Linux: 6}[p]; begin < wantBegin-0.5 || begin > wantBegin+0.5 {
				t.Errorf("%v %v: the start controls begin %v from the window's edge, want %v", p, dir, begin, wantBegin)
			}
			if after(edit, title) <= 0 || after(share, edit) <= 0 {
				t.Errorf("%v %v: the end controls do not follow the title: title %+v, Edit %+v, Share %+v", p, dir, title, edit, share)
			}
			// The end controls are at the far end: the room after them is the platform's, from the window controls' room or the edge.
			end := w - (share.X + share.W)
			if dir == ui.RightToLeft {
				end = share.X
			}
			want := map[Platform]float32{Mac: 8, Windows: 138 + 48 + 2, Linux: 46}[p]
			if end < want-0.5 || end > want+0.5 {
				t.Errorf("%v %v: the last control ends %v from the window's edge, want %v", p, dir, end, want)
			}
		}
	}
}

// TestHeaderBarControlsTakeClicks checks the bar on every platform: its controls take clicks and are Tab stops of their own in order, a disabled one takes neither, and a press on the bar between the title and the end controls clicks nothing and focuses nothing, as it moves the window (DragWindow, which ui.Tester cannot observe: MyGo's headless host does nothing on a window drag). The bar is a group named by the title.
func TestHeaderBarControlsTakeClicks(t *testing.T) {
	for _, p := range Platforms {
		s := &headerTestState{}
		tester := headerTester(p, ui.LeftToRight, 560, 160, s)
		named := headerFind(t, tester, "Documents")
		if named.Y != 0 || named.W != 560 {
			t.Errorf("%v: Documents names %+v, not the bar", p, named)
		}
		for _, name := range []string{"Back", "Edit", "Share"} {
			if err := tester.Click(name); err != nil {
				t.Fatalf("%v: %v", p, err)
			}
			tester.Frame()
		}
		if s.back != 1 || s.edit != 1 || s.share != 1 {
			t.Fatalf("%v: clicks reached Back %d, Edit %d, Share %d times", p, s.back, s.edit, s.share)
		}
		// Tab from the start: each control is a stop of its own.
		edit := headerFind(t, tester, "Edit")
		tester.ClickAt(edit.X-30, edit.Y+edit.H/2)
		tester.Frame()
		if tester.Focused("Back") || tester.Focused("Edit") || tester.Focused("Share") {
			t.Fatalf("%v: a press on the bar's empty part focused a control", p)
		}
		if s.back != 1 || s.edit != 1 || s.share != 1 {
			t.Fatalf("%v: a press on the bar's empty part clicked a control", p)
		}
		for _, name := range []string{"Back", "Edit", "Share"} {
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			if !tester.Focused(name) {
				t.Fatalf("%v: Tab did not reach %s", p, name)
			}
		}
		tester.Key(0, ui.KeySpace)
		tester.Frame()
		if s.share != 2 {
			t.Fatalf("%v: Space on Share clicked it %d times", p, s.share)
		}
		s.shareOff = true
		tester.Frame()
		tester.Frame()
		if err := tester.Click("Share"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		if s.share != 2 {
			t.Fatalf("%v: a disabled control took a click", p)
		}
	}
}

// TestHeaderBarKeepsItsParts checks on every platform that a subtitle shows only when given, that nil start and end controls build nothing, and that the focus stays on a control while the subtitle and the start controls come and go, as the bar builds its parts in every frame and hides those not shown. On Windows the bar is 32 DIPs high without controls and 48 with them, as WinUI's TitleBar.
func TestHeaderBarKeepsItsParts(t *testing.T) {
	for _, p := range Platforms {
		s := &headerTestState{name: "Window", subtitle: "12 items"}
		tester := headerTester(p, ui.LeftToRight, 560, 160, s)
		if !tester.HasText("12 items") {
			t.Fatalf("%v: the subtitle is missing", p)
		}
		title, sub := headerFind(t, tester, "Documents"), headerFind(t, tester, "12 items")
		if p == Windows {
			if sub.X <= title.X+title.W {
				t.Errorf("%v: the subtitle %+v is not beside the title %+v", p, sub, title)
			}
		} else if sub.Y < title.Y+title.H-0.5 {
			t.Errorf("%v: the subtitle %+v is not under the title %+v", p, sub, title)
		}
		if err := tester.Click("Edit"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		s.subtitle, s.noStart = "", true
		tester.Frame()
		if tester.HasText("12 items") {
			t.Errorf("%v: an empty subtitle still shows", p)
		}
		if _, ok := tester.Find("Back"); ok {
			t.Errorf("%v: a nil start still builds its controls", p)
		}
		if !tester.Focused("Edit") {
			t.Errorf("%v: the focus left Edit as the subtitle and the start controls went", p)
		}
		s.subtitle, s.noStart = "3 items", false
		tester.Frame()
		if !tester.Focused("Edit") || !tester.HasText("3 items") {
			t.Errorf("%v: the focus left Edit as the subtitle and the start controls came back", p)
		}
		s.noStart, s.noEnd = true, true
		tester.SetTitleBar(ui.TitleBar{})
		bar := headerFind(t, tester, "Window")
		want := map[Platform]float32{Mac: 52, Windows: 32, Linux: 47}[p]
		if bar.H != want {
			t.Errorf("%v: without controls the bar is %v high, want %v", p, bar.H, want)
		}
	}
}

// TestHeaderBarControlsLook checks the look the controls take in a header bar on every platform, as they take a toolbar's: on macOS 36-point toolbar items, on Windows the TitleBar's 40-DIP image button as tall as the bar less 4 beside a push button 32 high, outside a toolbar's look, and on GNOME 34-pixel buttons, the image button flat, its face the bar's background.
func TestHeaderBarControlsLook(t *testing.T) {
	for _, p := range Platforms {
		s := &headerTestState{editName: "Edit button"}
		tester := headerTester(p, ui.LeftToRight, 560, 160, s)
		edit, share := headerFind(t, tester, "Edit button"), headerFind(t, tester, "Share")
		var editH, shareW, shareH float32
		switch p {
		case Windows:
			editH, shareW, shareH = 32, 40, 44
		case Linux:
			editH, shareW, shareH = 34, 34, 34
		default:
			editH, shareW, shareH = 36, 36, 36
		}
		if edit.H != editH || share.W != shareW || share.H != shareH {
			t.Errorf("%v: Edit is %v high and Share %v by %v, want %v and %v by %v", p, edit.H, share.W, share.H, editH, shareW, shareH)
		}
		if p == Linux {
			img := tester.Image()
			bg := adwHeaderBarBackground(ui.LightTheme())
			at := img.RGBAAt(int(share.X+2), int(share.Y+share.H/2))
			if at.R != bg.R || at.G != bg.G || at.B != bg.B {
				t.Errorf("%v: the flat image button's face reads %v, want the bar's %v", p, at, bg)
			}
		}
	}
}

// TestScreenshotsHeaderBar writes a window with a header bar on every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names, once without a subtitle and once with one, the window controls' room as headerTestRoom gives it, as tools/appkit-reference's header bar scenes and the headerbar scenes of tools/winui-reference and tools/adwaita-reference draw the platforms' own. Without it the test only builds the frames.
func TestScreenshotsHeaderBar(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			for _, subtitle := range []string{"", "12 items"} {
				tester := ui.NewTester(headerView(p, &headerTestState{subtitle: subtitle}), 560, 120)
				tester.SetScale(2)
				tester.SetDark(dark)
				tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
				tester.SetTitleBar(headerTestRoom(p, ui.LeftToRight))
				tester.Frame()
				tester.Frame()
				if !tester.HasText("Documents") || !tester.HasText("Edit") {
					t.Fatalf("%v: the bar's title or controls are missing", p)
				}
				mode := map[bool]string{false: "light", true: "dark"}[dark]
				name := "headerbar"
				if subtitle != "" {
					name += "-subtitle"
				}
				writeShot(t, dir, fmt.Sprintf("breeze-%s-%s-%s.png", name, strings.ToLower(p.String()), mode), tester)
			}
		}
	}
}

// TestHeaderBarHistoryButtons checks BackButton and ForwardButton at the start of a header bar on every platform: they go back and forward as anywhere, and take the bar's look: on macOS image items sharing a 36-point platter, on Windows the TitleBar's back button, 40 DIPs wide and as tall as the 48-DIP bar less 4, and on GNOME flat 34-pixel buttons.
func TestHeaderBarHistoryButtons(t *testing.T) {
	for _, p := range Platforms {
		r := ui.NewRouter("/notes")
		r.Push("/notes/42")
		tester := layoutTester(p, 560, 160, func(c *ui.Context) {
			Use(c)
			ui.Column(c).Fill().Children(func() {
				HeaderBar(c, "Notes", "", func() {
					BackButton(c, r)
					ForwardButton(c, r)
				}, nil)
			})
		})
		tester.SetTitleBar(headerTestRoom(p, ui.LeftToRight))
		tester.Frame()
		back := headerFind(t, tester, "Back")
		want := map[Platform][2]float32{Mac: {36, 36}, Windows: {40, 44}, Linux: {34, 34}}[p]
		if back.W < want[0]-0.5 || back.W > want[0]+0.5 || back.H != want[1] {
			t.Errorf("%v: Back is %v by %v, want %v by %v", p, back.W, back.H, want[0], want[1])
		}
		if begin := map[Platform]float32{Mac: 79 + 17, Windows: 2, Linux: 6}[p]; back.X < begin-0.5 || back.X > begin+0.5 {
			t.Errorf("%v: Back begins at %v, want %v", p, back.X, begin)
		}
		if err := tester.Click("Back"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		if r.Path() != "/notes" {
			t.Fatalf("%v: Back left the router at %q", p, r.Path())
		}
		if err := tester.Click("Forward"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		if r.Path() != "/notes/42" {
			t.Fatalf("%v: Forward left the router at %q", p, r.Path())
		}
	}
}

// TestHeaderBarAdwSidesHoldTheirControls checks GNOME's bar in a window too narrow to center the title between its sides: the end controls, wider than half the room the title leaves, stay inside the bar and clear of the window controls, and the title moves toward the start rather than under them, as GtkCenterBox moves its center child aside.
func TestHeaderBarAdwSidesHoldTheirControls(t *testing.T) {
	const w = 330
	s := &headerTestState{name: "Window"}
	tester := headerTester(Linux, ui.LeftToRight, w, 160, s)
	tester.Frame()
	title, edit, share := headerFind(t, tester, "Documents"), headerFind(t, tester, "Edit"), headerFind(t, tester, "Share")
	room := headerTestRoom(Linux, ui.LeftToRight)
	if share.X+share.W > w-room.Right+0.5 {
		t.Errorf("Share ends at %v, past the window controls' room at %v", share.X+share.W, w-room.Right)
	}
	if title.X+title.W > edit.X {
		t.Errorf("the title, ending at %v, runs under Edit at %v", title.X+title.W, edit.X)
	}
}
