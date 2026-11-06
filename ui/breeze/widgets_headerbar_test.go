package breeze

import (
	"fmt"
	"image"
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
	name              string  // a label given after HeaderBar, naming the bar instead of its title
	editName          string  // a label given to the Edit button, to find the button rather than its text
	sidebar           float32 // in headerPaneView, the width of the sidebar at the window's start, none at 0
	sidebarAtEnd      bool    // in headerPaneView, the sidebar is at the window's end
	above             float32 // in headerPaneView, the height of a box above the bar, none at 0
	backRef, editRef  ui.Handle
	shareRef          ui.Handle // bound to the Back, Edit and Share buttons
}

// headerView builds a window with a header bar at its top on platform p: the title Documents, s's subtitle, an image button Back at its start and a text button Edit and an image button Share at its end, then the window's content.
func headerView(p Platform, s *headerTestState) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Column(c).Fill().Children(func() { headerBarAndContent(c, s) })
		})
	}
}

// headerBarAndContent builds the bar of headerView and the window's content under it, in the column it is called in.
func headerBarAndContent(c *ui.Context, s *headerTestState) {
	start := func() {
		if IconButton(c, headerTestBackIcon, "Back").Bind(&s.backRef).Clicked() {
			s.back++
		}
	}
	end := func() {
		edit := Button(c, "Edit").Bind(&s.editRef)
		if s.editName != "" {
			edit.Label(s.editName)
		}
		if edit.Clicked() {
			s.edit++
		}
		if IconButton(c, idiomTestShareIcon, "Share").Disabled(s.shareOff).Bind(&s.shareRef).Clicked() {
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
}

// headerPaneView builds a window on p whose window controls take room, holding a Sidebar s.sidebar DIPs wide at its start, or at its end with s.sidebarAtEnd (none at 0), its items below the controls' band, and beside it a column holding a box s.above high (none at 0), then headerView's bar and content, as a window with a sidebar puts a header bar over its content pane.
func headerPaneView(p Platform, room ui.TitleBar, s *headerTestState) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Row(c).Fill().AlignItems(ui.Stretch).Children(func() {
				selected := "inbox"
				sidebar := func() {
					Sidebar(c, &selected, func() {
						// The window controls' band, over the sidebar.
						ui.Box(c).Shrink(0).Height(room.Height).Role(ui.RoleNone)
						SidebarItem(c, "inbox", navTestFolder, "Inbox")
						SidebarItem(c, "sent", navTestFolder, "Sent")
					}).Shrink(0).Width(s.sidebar).Hidden(s.sidebar == 0).Label("Folders")
				}
				if !s.sidebarAtEnd {
					sidebar()
				}
				ui.Column(c).Grow(1).Shrink(1).MinWidth(0).Children(func() {
					ui.Box(c).Shrink(0).Height(s.above).Hidden(s.above == 0)
					headerBarAndContent(c, s)
				})
				if s.sidebarAtEnd {
					sidebar()
				}
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

// headerPaneSidebar is the width of the sidebar beside the bar in the pane tests, as wide as the sidebar of tools/appkit-reference's split header bar scenes.
const headerPaneSidebar = 236

// headerPaneTester builds headerPaneView on p in a w by 160 window whose controls take headerTestRoom, laid out in dir.
func headerPaneTester(p Platform, dir ui.Direction, w int, s *headerTestState) *ui.Tester {
	room := headerTestRoom(p, dir)
	tester := layoutTester(p, w, 160, headerPaneView(p, room, s))
	tester.SetDirection(dir)
	tester.SetTitleBar(room)
	tester.Frame()
	return tester
}

// headerPaneParts returns the boxes of the bar, of Back, Documents, Edit and Share in tester, failing the test for one missing or overlapping the window controls' room.
func headerPaneParts(t *testing.T, tester *ui.Tester, p Platform, dir ui.Direction, w float32) (bar, back, title, edit, share ui.Rect) {
	t.Helper()
	room := headerTestRoom(p, dir)
	left := ui.Rect{X: 0, Y: 0, W: room.Left, H: room.Height}
	right := ui.Rect{X: w - room.Right, Y: 0, W: room.Right, H: room.Height}
	parts := map[string]ui.Rect{}
	for _, name := range []string{"Back", "Documents", "Edit", "Share"} {
		r := headerFind(t, tester, name)
		parts[name] = r
		if overlaps(r, left) || overlaps(r, right) {
			t.Errorf("%v %v: %s at %+v overlaps the window controls' room %+v %+v", p, dir, name, r, left, right)
		}
	}
	return headerFind(t, tester, "Window"), parts["Back"], parts["Documents"], parts["Edit"], parts["Share"]
}

// TestHeaderBarInPaneKeepsClearOfTheControlsOverIt checks the bar over a window's content pane, beside a 236-DIP sidebar at the window's start, on every platform, left to right and right to left: it keeps clear of only the window controls it lies under, so that its start controls begin where the platform begins them from the pane's edge rather than after controls that are over the sidebar, as AppKit lays out the items and the title of a unified toolbar after a tracking separator over a split view's divider (tools/appkit-reference's split header bar scenes on macOS 27.2: the back item's platter 8 points after the divider, the title 12 after it, and 20 after the divider without the item), WinUI lays out a TitleBar with no caption buttons at its start (no LeftPaddingColumn, the back button 2 DIPs in) and libadwaita a content pane's AdwHeaderBar without window buttons at its start (6 pixels in, the title centered across the bar). The end controls keep clear of the window controls at the window's end, as across the whole window.
func TestHeaderBarInPaneKeepsClearOfTheControlsOverIt(t *testing.T) {
	const w = 720
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			s := &headerTestState{name: "Window", sidebar: headerPaneSidebar}
			tester := headerPaneTester(p, dir, w, s)
			bar, back, title, edit, share := headerPaneParts(t, tester, p, dir, w)
			// The pane is at the window's end: on the right left to right, on the left right to left.
			pane := ui.Rect{X: headerPaneSidebar, W: w - headerPaneSidebar}
			if dir == ui.RightToLeft {
				pane.X = 0
			}
			if bar.X != pane.X || bar.W != pane.W || bar.Y != 0 {
				t.Errorf("%v %v: the bar is %+v, want the pane's width %v at %v along the window's top", p, dir, bar, pane.W, pane.X)
			}
			// begin is how far r starts from the pane's start edge, and after how far a starts past b's end.
			begin := func(r ui.Rect) float32 {
				if dir == ui.RightToLeft {
					return pane.X + pane.W - (r.X + r.W)
				}
				return r.X - pane.X
			}
			after := func(a, b ui.Rect) float32 {
				if dir == ui.RightToLeft {
					return b.X - (a.X + a.W)
				}
				return a.X - (b.X + b.W)
			}
			if got, want := begin(back), map[Platform]float32{Mac: 8, Windows: 2, Linux: 6}[p]; got < want-0.5 || got > want+0.5 {
				t.Errorf("%v %v: the start controls begin %v from the pane's edge, want %v", p, dir, got, want)
			}
			switch p {
			case Linux:
				if mid := title.X + title.W/2; mid < pane.X+pane.W/2-1 || mid > pane.X+pane.W/2+1 {
					t.Errorf("%v %v: the title's middle is at %v, not centered across the pane %+v", p, dir, mid, pane)
				}
			default:
				if got, want := after(title, back), map[Platform]float32{Mac: 12, Windows: 16}[p]; got < want-0.5 || got > want+0.5 {
					t.Errorf("%v %v: the title starts %v after the start controls, want %v", p, dir, got, want)
				}
			}
			if after(edit, title) <= 0 || after(share, edit) <= 0 {
				t.Errorf("%v %v: the end controls do not follow the title: title %+v, Edit %+v, Share %+v", p, dir, title, edit, share)
			}
			end := w - (share.X + share.W)
			if dir == ui.RightToLeft {
				end = share.X
			}
			if want := map[Platform]float32{Mac: 8, Windows: 138 + 48 + 2, Linux: 46}[p]; end < want-0.5 || end > want+0.5 {
				t.Errorf("%v %v: the last control ends %v from the window's edge, want %v", p, dir, end, want)
			}
			if p == Mac {
				s.noStart = true
				tester.Frame()
				if got := begin(headerFind(t, tester, "Documents")); got < 20-0.5 || got > 20+0.5 {
					t.Errorf("%v %v: without start controls the title begins %v from the pane's edge, want 20", p, dir, got)
				}
			}
		}
	}
}

// TestHeaderBarInPaneAtTheStart checks the bar over a pane at the window's start, beside a 236-DIP sidebar at its end, on every platform, left to right and right to left: it keeps clear of the window controls at the window's start as across the whole window (17 points after the traffic lights on macOS, 2 DIPs in on Windows and 6 pixels on GNOME, where the controls are at the end), and of none at its own end, which stops short of the controls over the sidebar, so that its last control ends where the platform ends it from the bar's edge: 8 points on macOS, the 48-DIP TitleBarMinDragRegionWidth and the image button's margin of 2 on Windows, and 6 pixels on GNOME.
func TestHeaderBarInPaneAtTheStart(t *testing.T) {
	const w = 720
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			s := &headerTestState{name: "Window", sidebar: headerPaneSidebar, sidebarAtEnd: true}
			tester := headerPaneTester(p, dir, w, s)
			bar, back, _, _, share := headerPaneParts(t, tester, p, dir, w)
			pane := ui.Rect{X: 0, W: w - headerPaneSidebar}
			if dir == ui.RightToLeft {
				pane.X = headerPaneSidebar
			}
			if bar.X != pane.X || bar.W != pane.W {
				t.Errorf("%v %v: the bar is %+v, want the pane's width %v at %v", p, dir, bar, pane.W, pane.X)
			}
			begin, end := back.X-pane.X, pane.X+pane.W-(share.X+share.W)
			if dir == ui.RightToLeft {
				begin, end = pane.X+pane.W-(back.X+back.W), share.X-pane.X
			}
			if want := map[Platform]float32{Mac: 79 + 17, Windows: 2, Linux: 6}[p]; begin < want-0.5 || begin > want+0.5 {
				t.Errorf("%v %v: the start controls begin %v from the pane's edge, want %v", p, dir, begin, want)
			}
			if want := map[Platform]float32{Mac: 8, Windows: 48 + 2, Linux: 6}[p]; end < want-0.5 || end > want+0.5 {
				t.Errorf("%v %v: the end controls end %v from the pane's edge, want %v", p, dir, end, want)
			}
		}
	}
}

// TestHeaderBarFollowsItsPane checks on every platform, left to right and right to left, that the bar takes the room of the window controls it lies under in the frame it moves in: as the sidebar beside it goes and the bar runs under the controls at the window's start, a single frame keeps its start controls clear of them as across the whole window, and as the sidebar comes back they begin at the pane's edge again; a bar below the controls' band, under a box as high as it, keeps clear of none of them; and as the window grows, the bar over the pane keeps clear of the controls at the window's end.
func TestHeaderBarFollowsItsPane(t *testing.T) {
	for _, p := range Platforms {
		w := float32(720)
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			s := &headerTestState{name: "Window", sidebar: headerPaneSidebar}
			tester := headerPaneTester(p, dir, int(w), s)
			// start is how far the start controls begin from the window's start edge, and end how far the end controls end from its end edge.
			start := func() float32 {
				r := headerFind(t, tester, "Back")
				if dir == ui.RightToLeft {
					return w - (r.X + r.W)
				}
				return r.X
			}
			end := func() float32 {
				r := headerFind(t, tester, "Share")
				if dir == ui.RightToLeft {
					return r.X
				}
				return w - (r.X + r.W)
			}
			inPane := map[Platform]float32{Mac: headerPaneSidebar + 8, Windows: headerPaneSidebar + 2, Linux: headerPaneSidebar + 6}[p]
			whole := map[Platform]float32{Mac: 79 + 17, Windows: 2, Linux: 6}[p]
			steps := []struct {
				sidebar, above float32
				start, end     float32
			}{
				{0, 0, whole, map[Platform]float32{Mac: 8, Windows: 138 + 48 + 2, Linux: 46}[p]},
				{headerPaneSidebar, 0, inPane, map[Platform]float32{Mac: 8, Windows: 138 + 48 + 2, Linux: 46}[p]},
				{0, 60, map[Platform]float32{Mac: 8, Windows: 2, Linux: 6}[p], map[Platform]float32{Mac: 8, Windows: 48 + 2, Linux: 6}[p]},
			}
			for _, step := range steps {
				s.sidebar, s.above = step.sidebar, step.above
				tester.Frame()
				if got := start(); got < step.start-0.5 || got > step.start+0.5 {
					t.Errorf("%v %v: with a sidebar %v wide and a box %v high above the bar, the start controls begin %v from the window's edge, want %v", p, dir, step.sidebar, step.above, got, step.start)
				}
				if got := end(); got < step.end-0.5 || got > step.end+0.5 {
					t.Errorf("%v %v: with a sidebar %v wide and a box %v high above the bar, the end controls end %v from the window's edge, want %v", p, dir, step.sidebar, step.above, got, step.end)
				}
			}
			s.sidebar, s.above, w = headerPaneSidebar, 0, 900
			tester.SetSize(int(w), 160)
			if got, want := start(), inPane; got < want-0.5 || got > want+0.5 {
				t.Errorf("%v %v: in a window grown to %v the start controls begin %v from its edge, want %v", p, dir, w, got, want)
			}
			if got, want := end(), map[Platform]float32{Mac: 8, Windows: 138 + 48 + 2, Linux: 46}[p]; got < want-0.5 || got > want+0.5 {
				t.Errorf("%v %v: in a window grown to %v the end controls end %v from its edge, want %v", p, dir, w, got, want)
			}
		}
	}
}

// TestHeaderBarFollowsTheWindowInAFrameBuiltTwice checks on every platform, left to right and right to left, that a bar across the window keeps clear of the window controls in the frame the window grows in when that frame is built more than once, as MyGo builds a frame again after an event handled while building it (here a focus request the view makes as it sees the new width): every frame laid out at the new width, the first included, keeps the start and end controls clear of the controls' room.
func TestHeaderBarFollowsTheWindowInAFrameBuiltTwice(t *testing.T) {
	const from, to = 720, 820
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			s := &headerTestState{name: "Window"}
			room := headerTestRoom(p, dir)
			seen := float32(from)
			var laid []ui.Rect
			view := headerPaneView(p, room, s)
			tester := layoutTester(p, from, 160, func(c *ui.Context) {
				w, _ := c.Size()
				if w != seen {
					// Focusing Edit is handled while this pass builds, so MyGo builds the frame a second time before laying it out.
					seen = w
					s.editRef.Focus(c)
				}
				if w == to {
					// What the frames laid out since the window grew put Back and Share at; the frame it grew in reads the last frame laid out at the old width, whose controls are farther from the controls' room at the new width.
					laid = append(laid, s.backRef.Bounds(c), s.shareRef.Bounds(c))
				}
				view(c)
			})
			tester.SetDirection(dir)
			tester.SetTitleBar(room)
			tester.Frame()
			tester.SetSize(to, 160)
			tester.Frame()
			if !tester.Focused("Edit") {
				t.Fatalf("%v %v: Edit did not take the focus, so the frame was not built twice", p, dir)
			}
			left := ui.Rect{X: 0, Y: 0, W: room.Left, H: room.Height}
			right := ui.Rect{X: to - room.Right, Y: 0, W: room.Right, H: room.Height}
			for _, r := range laid {
				// Half a DIP in from each edge, for the rounding of a control that ends where the room begins.
				in := ui.Rect{X: r.X + 0.5, Y: r.Y + 0.5, W: r.W - 1, H: r.H - 1}
				if overlaps(in, left) || overlaps(in, right) {
					t.Errorf("%v %v: a frame laid out at the width of %v put a control at %+v, over the window controls' room %+v %+v", p, dir, to, r, left, right)
				}
			}
		}
	}
}

// TestHeaderBarInPaneLook checks the pixels of the bar over a content pane on every platform, light and dark, at 2x: the bar is drawn pixel for pixel as a bar across a window as wide as the pane whose controls take only the room at its end, as AppKit draws a content pane's toolbar after a tracking separator as a window's own from the divider on (tools/appkit-reference's split header bar scenes); and its line along its bottom runs under the pane alone, the sidebar beside it drawn as above the line's row.
func TestHeaderBarInPaneLook(t *testing.T) {
	const w, h, scale = 720, 120, 2
	const pane = w - headerPaneSidebar
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			shot := func(width int, room ui.TitleBar, sidebar float32) *image.RGBA {
				tester := ui.NewTester(headerPaneView(p, room, &headerTestState{sidebar: sidebar}), width, h)
				tester.SetScale(scale)
				tester.SetDark(dark)
				tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
				tester.SetTitleBar(room)
				tester.Frame()
				return tester.Image()
			}
			room := headerTestRoom(p, ui.LeftToRight)
			inPane := shot(w, room, headerPaneSidebar)
			alone := shot(pane, ui.TitleBar{Height: room.Height, Right: room.Right}, 0)
			bar := int(max(room.Height, map[Platform]float32{Mac: 52, Windows: 48, Linux: 47}[p])) * scale
			differ := 0
			for y := 0; y < bar; y++ {
				for x := 0; x < pane*scale; x++ {
					if inPane.RGBAAt(headerPaneSidebar*scale+x, y) != alone.RGBAAt(x, y) {
						differ++
					}
				}
			}
			if differ > 0 {
				t.Errorf("%v dark %v: %d pixels of the bar over the pane differ from a bar across a window as wide as the pane", p, dark, differ)
			}
			if p == Windows {
				// WinUI's TitleBar draws no line along its bottom.
				continue
			}
			line, above := inPane.RGBAAt((headerPaneSidebar+100)*scale, bar-1), inPane.RGBAAt((headerPaneSidebar+100)*scale, bar-8)
			if line == above {
				t.Errorf("%v dark %v: the bar's bottom row reads %v as the bar above it, want its line", p, dark, line)
			}
			if side, sideAbove := inPane.RGBAAt(100*scale, bar-1), inPane.RGBAAt(100*scale, bar-8); side != sideAbove {
				t.Errorf("%v dark %v: the sidebar reads %v in the row of the bar's line and %v above it, want the line under the pane alone", p, dark, side, sideAbove)
			}
		}
	}
}

// TestScreenshotsHeaderBarInPane writes a window with a sidebar and a header bar over its content pane on every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names, once without a subtitle and once with one, the window controls' room as headerTestRoom gives it, as tools/appkit-reference's split header bar scenes draw AppKit's; and once more without a subtitle with the sidebar at the window's end, where the pane stops short of the window controls at the end on Windows and GNOME. Without it the test only builds the frames.
func TestScreenshotsHeaderBarInPane(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			for _, shot := range []struct {
				subtitle string
				atEnd    bool
			}{{"", false}, {"12 items", false}, {"", true}} {
				room := headerTestRoom(p, ui.LeftToRight)
				tester := ui.NewTester(headerPaneView(p, room, &headerTestState{subtitle: shot.subtitle, sidebar: headerPaneSidebar, sidebarAtEnd: shot.atEnd}), 720, 120)
				tester.SetScale(2)
				tester.SetDark(dark)
				tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
				tester.SetTitleBar(room)
				tester.Frame()
				tester.Frame()
				if !tester.HasText("Documents") || !tester.HasText("Edit") {
					t.Fatalf("%v: the bar's title or controls are missing", p)
				}
				mode := map[bool]string{false: "light", true: "dark"}[dark]
				name := "headerbar-pane"
				if shot.subtitle != "" {
					name += "-subtitle"
				}
				if shot.atEnd {
					name += "-sidebar-end"
				}
				writeShot(t, dir, fmt.Sprintf("breeze-%s-%s-%s.png", name, strings.ToLower(p.String()), mode), tester)
			}
		}
	}
}
