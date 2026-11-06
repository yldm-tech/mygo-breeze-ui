package breeze

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func init() {
	drawnWidgets["TabView"] = true
}

// tabViewApp is a window's documents as a test keeps them: the labels, the chosen one, and what the close and add callbacks did.
type tabViewApp struct {
	labels            []string
	selected          int
	closed, adds      int
	closable, addable bool
}

func (a *tabViewApp) close(i int) {
	a.labels = append(a.labels[:i:i], a.labels[i+1:]...)
	a.closed++
}

func (a *tabViewApp) add() {
	a.adds++
	a.labels = append(a.labels, fmt.Sprintf("New %d", a.adds))
	a.selected = len(a.labels) - 1
}

// view builds the tab view with the chosen tab's content and a button after it, which takes the focus away from the strip.
func (a *tabViewApp) view(c *ui.Context) {
	var closeFn func(int)
	var addFn func()
	if a.closable {
		closeFn = a.close
	}
	if a.addable {
		addFn = a.add
	}
	ui.Column(c).Padding(16).Gap(12).Children(func() {
		TabView(c, &a.selected, a.labels, closeFn, addFn)
		if len(a.labels) > 0 {
			ui.Text(c, "Showing "+a.labels[a.selected])
		}
		Button(c, "Inside")
	})
}

// tabViewCloseAt returns where the close button of the tab in box r is: at the end of WinUI's and GNOME's tabs (4 in, 32 and 24 wide) and at the leading edge of macOS's (4 in, 18 wide), mirrored right to left.
func tabViewCloseAt(p Platform, r ui.Rect, rtl bool) (x, y float32) {
	from := map[Platform]float32{Mac: 4 + 9, Windows: 4 + 16, Linux: 4 + 12}[p]
	atEnd := p != Mac
	if atEnd == rtl {
		return r.X + from, r.Y + r.H/2
	}
	return r.X + r.W - from, r.Y + r.H/2
}

// TestTabViewChoosesClosesAndAdds checks the tab view on every platform, left to right and right to left: a click chooses a tab; a click on the chosen tab's close button closes it, the tab after it chosen in its place; closing the last tab chooses the one before; a click on the add button adds a tab, which the callback chooses; and the content follows.
func TestTabViewChoosesClosesAndAdds(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			rtl := dir == ui.RightToLeft
			app := &tabViewApp{labels: []string{"One", "Two", "Three"}, closable: true, addable: true}
			tester := layoutTester(p, 640, 300, app.view)
			tester.SetDirection(dir)
			tester.Frame()
			if err := tester.Click("Two"); err != nil {
				t.Fatalf("%v %v: %v", p, dir, err)
			}
			tester.Frame()
			if app.selected != 1 || !tester.HasText("Showing Two") {
				t.Fatalf("%v %v: a click on Two chose %d", p, dir, app.selected)
			}
			two, _ := tester.Find("Two")
			// The pointer over the tab shows its close button on macOS, as AppKit shows it only under the pointer.
			tester.Move(two.X+two.W/2, two.Y+two.H/2)
			tester.Frame()
			tester.ClickAt(tabViewCloseAt(p, two, rtl))
			tester.Frame()
			tester.Frame()
			if app.closed != 1 || strings.Join(app.labels, ",") != "One,Three" || app.selected != 1 || !tester.HasText("Showing Three") {
				t.Fatalf("%v %v: closing the chosen tab left %v, %d chosen, %d closed", p, dir, app.labels, app.selected, app.closed)
			}
			three, _ := tester.Find("Three")
			tester.Move(three.X+three.W/2, three.Y+three.H/2)
			tester.Frame()
			tester.ClickAt(tabViewCloseAt(p, three, rtl))
			tester.Frame()
			tester.Frame()
			if strings.Join(app.labels, ",") != "One" || app.selected != 0 {
				t.Fatalf("%v %v: closing the last, chosen tab left %v, %d chosen", p, dir, app.labels, app.selected)
			}
			if err := tester.Click(tabViewAdd(p)); err != nil {
				t.Fatalf("%v %v: %v", p, dir, err)
			}
			tester.Frame()
			tester.Frame()
			if strings.Join(app.labels, ",") != "One,New 1" || app.selected != 1 || !tester.HasText("Showing New 1") {
				t.Fatalf("%v %v: the add button left %v, %d chosen", p, dir, app.labels, app.selected)
			}
		}
	}
}

// TestTabViewCloseButtonsShowAsThePlatformShowsThem checks where a tab's close button takes a click on every platform: on Windows every closable tab shows one (CloseButtonOverlayMode Auto), so a click there closes an unchosen tab; on macOS (under the pointer only) and GNOME (the chosen tab, or under the pointer) a tab the pointer was not over takes the click there itself and is chosen.
func TestTabViewCloseButtonsShowAsThePlatformShowsThem(t *testing.T) {
	for _, p := range Platforms {
		app := &tabViewApp{labels: []string{"One", "Two", "Three"}, closable: true}
		tester := layoutTester(p, 640, 300, app.view)
		tester.Frame()
		three, _ := tester.Find("Three")
		tester.ClickAt(tabViewCloseAt(p, three, false))
		tester.Frame()
		tester.Frame()
		if p == Windows {
			if app.closed != 1 || app.selected != 0 || strings.Join(app.labels, ",") != "One,Two" {
				t.Errorf("%v: a click on an unchosen tab's close button left %v, %d chosen", p, app.labels, app.selected)
			}
			continue
		}
		if app.closed != 0 || app.selected != 2 {
			t.Errorf("%v: a click where an unchosen tab's hidden close button is closed %d tabs and chose %d", p, app.closed, app.selected)
		}
	}
}

// TestTabViewKeys checks the platform's keys on every platform, with the focus on a button after the strip and in the strip: Ctrl+Tab and Ctrl+Shift+Tab switch tabs, wrapping around, without moving the focus; Ctrl+F4 closes the chosen tab on Windows only; on GNOME Ctrl+Page Down and Up switch too, Alt+digits choose a tab, and Ctrl+Home and Ctrl+End choose the first and the last only from the strip. In the strip the arrows choose tabs, as MyGo's tabs do.
func TestTabViewKeys(t *testing.T) {
	for _, p := range Platforms {
		app := &tabViewApp{labels: []string{"One", "Two", "Three"}, closable: true}
		tester := layoutTester(p, 640, 300, app.view)
		if err := tester.Click("Inside"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		press := func(mods ui.Modifiers, k ui.Key) {
			tester.Key(mods, k)
			tester.Frame()
			tester.Frame()
		}
		press(ui.Ctrl, ui.KeyTab)
		if app.selected != 1 {
			t.Errorf("%v: Ctrl+Tab chose %d", p, app.selected)
		}
		press(ui.Ctrl|ui.Shift, ui.KeyTab)
		press(ui.Ctrl|ui.Shift, ui.KeyTab)
		if app.selected != 2 {
			t.Errorf("%v: Ctrl+Shift+Tab twice from the second tab chose %d, want the last", p, app.selected)
		}
		if !tester.Focused("Inside") {
			t.Errorf("%v: the keys moved the focus off the button", p)
		}
		press(ui.Ctrl, ui.KeyF4)
		if closed := app.closed == 1 && strings.Join(app.labels, ",") == "One,Two" && app.selected == 1; closed != (p == Windows) {
			t.Errorf("%v: Ctrl+F4 left %v, %d chosen", p, app.labels, app.selected)
		}
		app.labels, app.selected = []string{"One", "Two", "Three", "Four"}, 0
		tester.Frame()
		press(ui.Ctrl, ui.KeyPageDown)
		press(ui.Alt, ui.Key4)
		if want := map[bool]int{true: 3, false: 0}[p == Linux]; app.selected != want {
			t.Errorf("%v: Ctrl+Page Down and Alt+4 chose %d, want %d", p, app.selected, want)
		}
		app.selected = 1
		tester.Frame()
		press(ui.Ctrl, ui.KeyEnd)
		if app.selected != 1 {
			t.Errorf("%v: Ctrl+End from the button chose %d", p, app.selected)
		}
		// Into the strip: Shift+Tab from the button goes past the add-less strip to the chosen tab.
		press(ui.Shift, ui.KeyTab)
		if !tester.Focused("Two") {
			t.Fatalf("%v: Shift+Tab from the button did not focus the chosen tab", p)
		}
		// On GNOME Ctrl+End chooses the last tab, which takes the focus; elsewhere it does nothing.
		press(ui.Ctrl, ui.KeyEnd)
		linux := p == Linux
		if want, focus := map[bool]int{true: 3, false: 1}[linux], map[bool]string{true: "Four", false: "Two"}[linux]; app.selected != want || !tester.Focused(focus) {
			t.Errorf("%v: Ctrl+End in the strip chose %d, want %d with the focus on %s", p, app.selected, want, focus)
		}
		press(0, ui.KeyRight)
		if want, focus := map[bool]int{true: 0, false: 2}[linux], map[bool]string{true: "One", false: "Three"}[linux]; app.selected != want || !tester.Focused(focus) {
			t.Errorf("%v: Right in the strip chose %d, want %d with the focus on %s", p, app.selected, want, focus)
		}
		press(ui.Ctrl, ui.KeyTab)
		if want := map[bool]int{true: 1, false: 3}[linux]; app.selected != want {
			t.Errorf("%v: Ctrl+Tab in the strip chose %d, want %d", p, app.selected, want)
		}
	}
}

// TestTabViewFocusAndNames checks on every platform that Tab stops once in the tabs, at the chosen one, then at the add button, named as the platform names it, and not at the close buttons, which are named as the platform names them; without close and add, neither shows.
func TestTabViewFocusAndNames(t *testing.T) {
	for _, p := range Platforms {
		app := &tabViewApp{labels: []string{"One", "Two"}, selected: 1, closable: true, addable: true}
		tester := layoutTester(p, 640, 300, func(c *ui.Context) {
			ui.Column(c).Children(func() {
				Button(c, "Before")
				app.view(c)
			})
		})
		if err := tester.Click("Before"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		for _, want := range []string{"Two", tabViewAdd(p), "Inside"} {
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			if !tester.Focused(want) {
				t.Fatalf("%v: Tab did not go to %q", p, want)
			}
		}
		if _, ok := tester.Find(tabViewClose(p)); !ok {
			t.Errorf("%v: no close button is named %q", p, tabViewClose(p))
		}
		bare := &tabViewApp{labels: []string{"One", "Two"}}
		tester = layoutTester(p, 640, 300, bare.view)
		tester.Frame()
		for _, s := range []string{tabViewClose(p), tabViewAdd(p)} {
			if _, ok := tester.Find(s); ok {
				t.Errorf("%v: %q shows without its callback", p, s)
			}
		}
		if _, ok := tester.Find("Two"); !ok {
			t.Errorf("%v: the tabs are missing", p)
		}
	}
}

// TestTabViewDisabled checks on every platform that a disabled tab view takes neither clicks nor its keys.
func TestTabViewDisabled(t *testing.T) {
	for _, p := range Platforms {
		app := &tabViewApp{labels: []string{"One", "Two"}, closable: true, addable: true}
		tester := layoutTester(p, 640, 300, func(c *ui.Context) {
			TabView(c, &app.selected, app.labels, app.close, app.add).Disabled(true)
			Button(c, "Inside")
		})
		tester.Frame()
		for _, s := range []string{"Two", tabViewAdd(p)} {
			if err := tester.Click(s); err != nil {
				t.Fatalf("%v: %v", p, err)
			}
			tester.Frame()
		}
		if err := tester.Click("Inside"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		tester.Key(ui.Ctrl, ui.KeyTab)
		tester.Frame()
		tester.Frame()
		if app.selected != 0 || app.adds != 0 || app.closed != 0 {
			t.Errorf("%v: a disabled tab view chose %d, added %d, closed %d", p, app.selected, app.adds, app.closed)
		}
	}
}

// TestTabViewSharesTheWidth checks the tabs' sizes on every platform: they are as high as the platform's (24 points, 32 and 34 DIPs) and equal in width, on GNOME in whole pixels with the first ones a pixel narrower where the width does not divide evenly, as AdwTabBox shares it; on Windows each is at most 240 wide and the add button follows the last 7 DIPs after it (the footer's 4 and the container's 3), where on macOS and GNOME the tabs fill the bar and the button is at its end, 8 points from it and 6 pixels (the end action's 5 and the bar's 1).
func TestTabViewSharesTheWidth(t *testing.T) {
	for _, p := range Platforms {
		for _, width := range []int{400, 1200} {
			app := &tabViewApp{labels: []string{"One", "Two", "Three"}, closable: true, addable: true}
			var strip ui.Rect
			tester := layoutTester(p, width, 300, func(c *ui.Context) {
				strip = TabView(c, &app.selected, app.labels, app.close, app.add).Bounds()
			})
			tester.Frame()
			var tabs [3]ui.Rect
			for i, s := range app.labels {
				tabs[i], _ = tester.Find(s)
			}
			plus, _ := tester.Find(tabViewAdd(p))
			if want := map[Platform]float32{Mac: 24, Windows: 32, Linux: 34}[p]; tabs[0].H != want {
				t.Errorf("%v %d: a tab is %v high, want %v", p, width, tabs[0].H, want)
			}
			for i := 1; i < 3; i++ {
				d := tabs[i].W - tabs[i-1].W
				if p == Linux {
					if tabs[i].W != float32(int(tabs[i].W)) || d < 0 || d > 1 {
						t.Errorf("%v %d: tabs %v and %v are not whole pixels wide, the later one as wide or a pixel wider", p, width, tabs[i-1], tabs[i])
					}
				} else if d < -0.01 || d > 0.01 {
					t.Errorf("%v %d: tabs %v and %v differ in width", p, width, tabs[0], tabs[i])
				}
			}
			after := plus.X - (tabs[2].X + tabs[2].W)
			end := strip.X + strip.W - (plus.X + plus.W)
			switch p {
			case Windows:
				if tabs[0].W > 240.5 || after < 6.5 || after > 7.5 {
					t.Errorf("%v %d: tabs %v wide, the add button %v after them", p, width, tabs[0].W, after)
				}
			default:
				want := map[Platform]float32{Mac: 8, Linux: 6}[p]
				if end < want-0.5 || end > want+0.5 {
					t.Errorf("%v %d: the add button ends %v before the strip's end, want %v", p, width, end, want)
				}
			}
		}
	}
}

// tabViewView lays out the tab views for screenshots: three closable tabs with the add button, the first chosen, over its content, and below two tabs without close buttons or the add button, the second chosen. tools/winui-reference's winui-tabview scene and tools/adwaita-reference's tabview set place WinUI's and libadwaita's there.
func tabViewView(p Platform) func(c *ui.Context) {
	docs := &tabViewApp{labels: []string{"Documents", "Downloads", "Pictures"}, closable: true, addable: true}
	prefs := &tabViewApp{labels: []string{"General", "Advanced"}, selected: 1}
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Column(c).Fill().Padding(16).Children(func() {
				ui.Column(c).Height(120).Gap(12).Children(func() {
					TabView(c, &docs.selected, docs.labels, docs.close, docs.add)
					ui.Text(c, "Showing "+docs.labels[docs.selected])
				})
				TabView(c, &prefs.selected, prefs.labels, nil, nil)
			})
		})
	}
}

// TestScreenshotsTabView writes the tab views of every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names. Without it the test only builds the frames.
func TestScreenshotsTabView(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			tester := ui.NewTester(tabViewView(p), 560, 260)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tester.Frame()
			tester.Frame()
			for _, s := range []string{"Documents", "Pictures", "Showing Documents", "Advanced"} {
				if !tester.HasText(s) {
					t.Fatalf("%v: %q is missing", p, s)
				}
			}
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			writeShot(t, dir, fmt.Sprintf("breeze-tabview-%s-%s.png", strings.ToLower(p.String()), mode), tester)
		}
	}
}

// TestTabViewLeavesKeysItDoesNotUse checks on every platform that the tab view takes a key only when it acts on it, as WinUI's accelerators report it unhandled and AdwTabView's shortcuts propagate it: Ctrl+Tab and Ctrl+Shift+Tab with a single tab, Ctrl+F4 when the tabs do not close, Ctrl+Page Down with a single tab on GNOME, and an Alt+digit there for a tab that is not there or is the chosen one, all reach a shortcut of the window built after it.
func TestTabViewLeavesKeysItDoesNotUse(t *testing.T) {
	type press struct {
		mods ui.Modifiers
		key  ui.Key
	}
	for _, p := range Platforms {
		presses := []press{{ui.Ctrl, ui.KeyTab}, {ui.Ctrl | ui.Shift, ui.KeyTab}, {ui.Ctrl, ui.KeyF4}, {ui.Ctrl, ui.KeyPageDown}, {ui.Alt, ui.Key5}, {ui.Alt, ui.Key1}}
		for _, labels := range [][]string{{"Only"}, {"One", "Two", "Three"}} {
			selected := 0
			got := map[press]bool{}
			tester := layoutTester(p, 640, 300, func(c *ui.Context) {
				TabView(c, &selected, labels, nil, nil)
				Button(c, "Inside")
				for _, k := range presses {
					if c.Shortcut(k.mods, k.key) {
						got[k] = true
					}
				}
			})
			if err := tester.Click("Inside"); err != nil {
				t.Fatalf("%v: %v", p, err)
			}
			tester.Frame()
			for _, k := range presses {
				// With three tabs the switching keys act; Alt+5 and Alt+1 (the chosen tab) do not on GNOME.
				acts := len(labels) > 1 && (k.key == ui.KeyTab || (p == Linux && k.key == ui.KeyPageDown))
				tester.Key(k.mods, k.key)
				tester.Frame()
				tester.Frame()
				if got[k] == acts {
					t.Errorf("%v %d tabs: the window's shortcut got %v+%v: %v", p, len(labels), k.mods, k.key, got[k])
				}
				selected = 0
				tester.Frame()
			}
		}
	}
}

// TestTabViewFocusFollowsTheKeys checks on every platform that with the focus on a tab the keys switching tabs take it to the tab they choose, as AdwTabBox gives the newly chosen tab the focus and as MyGo's tabs follow the arrows, and on Windows that Ctrl+F4 closing the focused last tab leaves the focus on the tab chosen in its place, as WinUI's TabView moves it to an adjacent tab.
func TestTabViewFocusFollowsTheKeys(t *testing.T) {
	for _, p := range Platforms {
		app := &tabViewApp{labels: []string{"One", "Two", "Three"}, closable: true}
		tester := layoutTester(p, 640, 300, app.view)
		if err := tester.Click("Inside"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		tester.Key(ui.Shift, ui.KeyTab)
		tester.Frame()
		if !tester.Focused("One") {
			t.Fatalf("%v: Shift+Tab did not focus the chosen tab", p)
		}
		tester.Key(ui.Ctrl, ui.KeyTab)
		tester.Frame()
		tester.Frame()
		if app.selected != 1 || !tester.Focused("Two") {
			t.Errorf("%v: Ctrl+Tab in the strip chose %d and left the focus behind", p, app.selected)
		}
		tester.Key(ui.Ctrl|ui.Shift, ui.KeyTab)
		tester.Key(ui.Ctrl|ui.Shift, ui.KeyTab)
		tester.Frame()
		tester.Frame()
		tester.Frame()
		if app.selected != 2 || !tester.Focused("Three") {
			t.Errorf("%v: Ctrl+Shift+Tab twice in the strip chose %d and left the focus behind", p, app.selected)
		}
		if p != Windows {
			continue
		}
		tester.Key(ui.Ctrl, ui.KeyF4)
		tester.Frame()
		tester.Frame()
		if strings.Join(app.labels, ",") != "One,Two" || app.selected != 1 || !tester.Focused("Two") {
			t.Errorf("%v: Ctrl+F4 on the focused last tab left %v, %d chosen, the focus elsewhere", p, app.labels, app.selected)
		}
	}
}

// TestTabViewDisabledLeavesKeysToAnother checks on every platform that a disabled tab view does not take the window's keys from an enabled one built after it.
func TestTabViewDisabledLeavesKeysToAnother(t *testing.T) {
	for _, p := range Platforms {
		first, second := 0, 0
		tester := layoutTester(p, 640, 300, func(c *ui.Context) {
			TabView(c, &first, []string{"One", "Two"}, nil, nil).Disabled(true)
			TabView(c, &second, []string{"Three", "Four"}, nil, nil)
			Button(c, "Inside")
		})
		if err := tester.Click("Inside"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		tester.Key(ui.Ctrl, ui.KeyTab)
		tester.Frame()
		tester.Frame()
		if first != 0 || second != 1 {
			t.Errorf("%v: Ctrl+Tab chose %d in the disabled tab view and %d in the other", p, first, second)
		}
	}
}

// TestTabViewCentersTitlesWithoutClose checks on macOS and GNOME, which center a tab's title, that the title of a tab without a close button is centered on the tab, as AdwTab centers it across the whole tab and AppKit's tab bar on its tab, reading the title's ink from the frame rather than the host font's widths.
func TestTabViewCentersTitlesWithoutClose(t *testing.T) {
	for _, p := range []Platform{Mac, Linux} {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			selected := 1
			tester := layoutTester(p, 640, 200, func(c *ui.Context) {
				ui.Column(c).Padding(16).Children(func() {
					TabView(c, &selected, []string{"General", "Advanced"}, nil, nil)
				})
			})
			tester.SetDirection(dir)
			tester.Frame()
			img := tester.Image()
			for _, s := range []string{"General", "Advanced"} {
				r, _ := tester.Find(s)
				y := int(r.Y + r.H/2)
				bg := img.RGBAAt(int(r.X+r.W/2), int(r.Y+2))
				lo, hi := -1, -1
				for x := int(r.X) + 6; x < int(r.X+r.W)-6; x++ {
					for dy := -4; dy <= 4; dy++ {
						c := img.RGBAAt(x, y+dy)
						if d := int(c.G) - int(bg.G); d > 48 || d < -48 {
							if lo < 0 {
								lo = x
							}
							hi = x
							break
						}
					}
				}
				if lo < 0 {
					t.Fatalf("%v %v: no ink in %q", p, dir, s)
				}
				if off := float32(lo+hi+1)/2 - (r.X + r.W/2); off < -2 || off > 2 {
					t.Errorf("%v %v: the title of %q is %v off the tab's center", p, dir, s, off)
				}
			}
		}
	}
}

// TestTabViewCloseButtonInsets checks on every platform, left to right and right to left, where the first tab's close button is, chosen and not: 4 points in from the tab's leading edge on macOS (its accessibility frame on macOS 27.2, the chosen tab's rim inside those 4), 4 DIPs in from the end of the chosen tab on Windows (TabViewItemHeaderPaddingWithCloseButton) and 5 from another's, whose TabContainer keeps a transparent 1-DIP border where the chosen one's reaches past the tab, and 4 pixels in from the end on GNOME (the tab's padding).
func TestTabViewCloseButtonInsets(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			for _, chosen := range []int{0, 1} {
				app := &tabViewApp{labels: []string{"One", "Two"}, selected: chosen, closable: true}
				tester := layoutTester(p, 640, 300, app.view)
				tester.SetDirection(dir)
				tester.Frame()
				tab, _ := tester.Find("One")
				x, _ := tester.Find(tabViewClose(p))
				fromLeft, fromRight := x.X-tab.X, tab.X+tab.W-(x.X+x.W)
				in := fromRight
				if (p == Mac) != (dir == ui.RightToLeft) {
					in = fromLeft
				}
				want := float32(4)
				if p == Windows && chosen != 0 {
					want = 5
				}
				if in < want-0.01 || in > want+0.01 {
					t.Errorf("%v %v, tab %d chosen: the first tab's close button is %v inside it, want %v", p, dir, chosen, in, want)
				}
			}
		}
	}
}
