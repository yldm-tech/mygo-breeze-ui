package breeze

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func init() {
	for _, n := range []string{"Sidebar", "SidebarSection", "SidebarItem", "Breadcrumbs", "BackButton", "ForwardButton", "Toolbar"} {
		drawnWidgets[n] = true
	}
}

// navTestFolder is a folder icon for the sidebars of the tests.
var navTestFolder = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><path d="M1.5 3h4.2l1.6 1.6h7.2v8.4h-13z" fill="none" stroke="#000" stroke-width="1.3" stroke-linejoin="round"/></svg>`))

// navState is what the navigation widgets of the tests edit.
type navState struct {
	selected            string
	open                bool
	crumb               int
	sidebarChanged      int
	crumbChanged        int
	bold                bool
	shares              int
	router              *ui.Router
	disabledSidebar     bool
	disabledBreadcrumbs bool
}

func newNavState() *navState {
	r := ui.NewRouter("/notes")
	r.Push("/notes/42")
	return &navState{selected: "inbox", open: true, crumb: -1, router: r}
}

// navView lays out the navigation widgets: a sidebar of two sections at the start, and a toolbar above a path beside it.
func navView(p Platform, s *navState) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Row(c).Fill().AlignItems(ui.Stretch).Children(func() {
				sb := Sidebar(c, &s.selected, func() {
					SidebarSection(c, "Mailboxes", &s.open, func() {
						SidebarItem(c, "inbox", navTestFolder, "Inbox").Children(func() { Badge(c, "12") })
						SidebarItem(c, "sent", navTestFolder, "Sent")
						SidebarItem(c, "drafts", navTestFolder, "Drafts")
					})
					SidebarSection(c, "Smart Mailboxes", nil, func() {
						SidebarItem(c, "flagged", nil, "Flagged")
					})
				}).Width(200).Label("Folders").Disabled(s.disabledSidebar)
				if sb.Changed() {
					s.sidebarChanged++
				}
				ui.Column(c).Grow(1).Gap(12).Padding(0, 12).Children(func() {
					Toolbar(c, func() {
						BackButton(c, s.router)
						ForwardButton(c, s.router)
						Toggle(c, &s.bold, "Bold")
						ui.Spacer(c)
						if Button(c, "Share").Clicked() {
							s.shares++
						}
					}).Label("Format")
					bc := Breadcrumbs(c, []string{"Home", "Users", "Ada"}, &s.crumb).Label("Path").Disabled(s.disabledBreadcrumbs)
					if bc.Changed() {
						s.crumbChanged++
					}
				})
			})
		})
	}
}

// TestNavigationWidgetsKeepTheirBehavior checks on every platform, left to right and right to left, that the restyled sidebar chooses with clicks, the arrows and typing and is one stop of Tab, that its sections open and close, that a path's links choose with a click, Enter and Space while its last item does not, that the history buttons go back and forward, stop at the history's ends and list the pages before on a right click, and that a toolbar's controls click and take the focus with the arrows.
func TestNavigationWidgetsKeepTheirBehavior(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			name := fmt.Sprintf("%v %v", p, map[ui.Direction]string{ui.LeftToRight: "ltr", ui.RightToLeft: "rtl"}[dir])
			s := newNavState()
			tester := ui.NewTester(navView(p, s), 800, 420)
			tester.SetDirection(dir)
			tester.Frame()
			for _, text := range []string{"Mailboxes", "Inbox", "12", "Flagged", "Home", "Ada", "Bold"} {
				if !tester.HasText(text) {
					t.Fatalf("%s: %q is missing", name, text)
				}
			}

			// The sidebar.
			if err := tester.Click("Sent"); err != nil {
				t.Fatalf("%s: click Sent: %v", name, err)
			}
			tester.Frame()
			if s.selected != "sent" || s.sidebarChanged != 1 {
				t.Fatalf("%s: a click chose %q, changed %d", name, s.selected, s.sidebarChanged)
			}
			if !tester.Focused("Folders") {
				t.Fatalf("%s: the sidebar did not take the focus", name)
			}
			tester.Key(0, ui.KeyDown)
			tester.Frame()
			if s.selected != "drafts" {
				t.Fatalf("%s: Down chose %q, want drafts", name, s.selected)
			}
			tester.Key(0, ui.KeyEnd)
			tester.Frame()
			if s.selected != "flagged" {
				t.Fatalf("%s: End chose %q, want flagged", name, s.selected)
			}
			tester.Key(0, ui.KeyHome)
			tester.Frame()
			if s.selected != "inbox" {
				t.Fatalf("%s: Home chose %q, want inbox", name, s.selected)
			}
			tester.Key(0, ui.KeyF)
			tester.Frame()
			if s.selected != "flagged" {
				t.Fatalf("%s: typing F chose %q, want flagged", name, s.selected)
			}
			tester.Key(0, ui.KeyHome)
			tester.Frame()
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			if tester.Focused("Folders") || !tester.Focused("Back") {
				t.Fatalf("%s: Tab from the sidebar did not go to the toolbar's first control", name)
			}
			if err := tester.Click("Mailboxes"); err != nil {
				t.Fatalf("%s: click the section's title: %v", name, err)
			}
			tester.Frame()
			if s.open || tester.HasText("Drafts") {
				t.Fatalf("%s: a click on the title left the section open", name)
			}
			if err := tester.Click("Mailboxes"); err != nil {
				t.Fatalf("%s: click the section's title again: %v", name, err)
			}
			tester.Frame()
			if !s.open || !tester.HasText("Drafts") {
				t.Fatalf("%s: a second click left the section closed", name)
			}

			// The path.
			if err := tester.Click("Users"); err != nil {
				t.Fatalf("%s: click Users: %v", name, err)
			}
			tester.Frame()
			if s.crumb != 1 || s.crumbChanged != 1 {
				t.Fatalf("%s: a click on Users chose %d, changed %d", name, s.crumb, s.crumbChanged)
			}
			if !tester.Focused("Users") {
				t.Fatalf("%s: the link did not take the focus", name)
			}
			s.crumb = -1
			tester.Key(0, ui.KeyEnter)
			tester.Frame()
			if s.crumb != 1 || s.crumbChanged != 2 {
				t.Fatalf("%s: Enter on Users chose %d, changed %d", name, s.crumb, s.crumbChanged)
			}
			s.crumb = -1
			tester.Key(0, ui.KeySpace)
			tester.Frame()
			if s.crumb != 1 || s.crumbChanged != 3 {
				t.Fatalf("%s: Space on Users chose %d, changed %d", name, s.crumb, s.crumbChanged)
			}
			if err := tester.Click("Ada"); err != nil {
				t.Fatalf("%s: click Ada: %v", name, err)
			}
			tester.Frame()
			if s.crumb != 1 || s.crumbChanged != 3 {
				t.Fatalf("%s: a click on the last item chose %d, changed %d", name, s.crumb, s.crumbChanged)
			}

			// The history buttons.
			if err := tester.Click("Back"); err != nil {
				t.Fatalf("%s: click Back: %v", name, err)
			}
			tester.Frame()
			if got := s.router.Path(); got != "/notes" {
				t.Fatalf("%s: Back went to %q", name, got)
			}
			if err := tester.Click("Back"); err != nil {
				t.Fatalf("%s: click Back again: %v", name, err)
			}
			tester.Frame()
			if got := s.router.Path(); got != "/notes" {
				t.Fatalf("%s: Back at the history's start went to %q", name, got)
			}
			if err := tester.Click("Forward"); err != nil {
				t.Fatalf("%s: click Forward: %v", name, err)
			}
			tester.Frame()
			if got := s.router.Path(); got != "/notes/42" {
				t.Fatalf("%s: Forward went to %q", name, got)
			}
			if err := tester.Click("Forward"); err != nil {
				t.Fatalf("%s: click Forward again: %v", name, err)
			}
			tester.Frame()
			if got := s.router.Path(); got != "/notes/42" {
				t.Fatalf("%s: Forward at the history's end went to %q", name, got)
			}
			if err := tester.RightClick("Back"); err != nil {
				t.Fatalf("%s: right click Back: %v", name, err)
			}
			tester.Frame()
			if got := tester.Menu(); len(got) != 1 || got[0] != "/notes" {
				t.Fatalf("%s: Back's menu lists %v, want the page before", name, got)
			}
			if err := tester.ChooseMenuItem("/notes"); err != nil {
				t.Fatalf("%s: choose the page before: %v", name, err)
			}
			tester.Frame()
			if got := s.router.Path(); got != "/notes" {
				t.Fatalf("%s: Back's menu went to %q", name, got)
			}

			// The toolbar.
			if err := tester.Click("Bold"); err != nil {
				t.Fatalf("%s: click Bold: %v", name, err)
			}
			tester.Frame()
			if !s.bold {
				t.Fatalf("%s: the toolbar's toggle did not turn on", name)
			}
			next := ui.KeyRight
			if dir == ui.RightToLeft {
				next = ui.KeyLeft
			}
			tester.Key(0, next)
			tester.Frame()
			if !tester.Focused("Share") {
				t.Fatalf("%s: the arrow after Bold did not move the focus to Share", name)
			}
			tester.Key(0, ui.KeyEnter)
			tester.Frame()
			if s.shares != 1 {
				t.Fatalf("%s: Enter on Share clicked it %d times", name, s.shares)
			}
		}
	}
}

// TestNavigationWidgetsDisabled checks that a disabled sidebar and a disabled path refuse clicks on every platform.
func TestNavigationWidgetsDisabled(t *testing.T) {
	for _, p := range Platforms {
		s := newNavState()
		s.disabledSidebar, s.disabledBreadcrumbs = true, true
		tester := ui.NewTester(navView(p, s), 800, 420)
		tester.Frame()
		for _, text := range []string{"Sent", "Users"} {
			if err := tester.Click(text); err != nil {
				t.Fatalf("%v: click %q: %v", p, text, err)
			}
			tester.Frame()
		}
		if s.selected != "inbox" || s.crumb != -1 || s.sidebarChanged != 0 || s.crumbChanged != 0 {
			t.Fatalf("%v: disabled widgets chose %q and %d", p, s.selected, s.crumb)
		}
	}
}

// TestMacSidebarSelectionFollowsTheFocus checks that the chosen item of a sidebar on macOS is drawn in the emphasized accent color while the sidebar has the focus in an active window, and in unemphasizedSelectedContentBackgroundColor otherwise, as AppKit draws a source list.
func TestMacSidebarSelectionFollowsTheFocus(t *testing.T) {
	for _, dark := range []bool{false, true} {
		s := newNavState()
		tester := ui.NewTester(navView(Mac, s), 800, 420)
		tester.SetDark(dark)
		tester.Frame()
		at := func() ui.Color {
			r, ok := tester.Find("Inbox")
			if !ok {
				t.Fatal("Inbox is missing")
			}
			px := tester.Image().RGBAAt(int(r.X-3), int(r.Y+r.H/2))
			return ui.RGB(px.R, px.G, px.B)
		}
		gray := pick(dark, ui.Hex("#dcdcdc"), ui.Hex("#464646"))
		accent := pick(dark, ui.Hex("#0070f5"), ui.Hex("#007aff"))
		if got := at(); got != gray {
			t.Fatalf("dark %v: unfocused selection %v, want %v", dark, got, gray)
		}
		if err := tester.Click("Inbox"); err != nil {
			t.Fatal(err)
		}
		tester.Move(700, 400)
		tester.Frame()
		if got := at(); got != accent {
			t.Fatalf("dark %v: focused selection %v, want %v", dark, got, accent)
		}
		tester.SetFocused(false)
		tester.Frame()
		if got := at(); got != gray {
			t.Fatalf("dark %v: selection in an inactive window %v, want %v", dark, got, gray)
		}
	}
}

// TestNavigationToolbarOverflow checks on every platform that the controls of a toolbar too narrow for them go into its menu, from the last, and that choosing one there clicks it, as in ui.Toolbar.
func TestNavigationToolbarOverflow(t *testing.T) {
	for _, p := range Platforms {
		s := newNavState()
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Column(c).Width(160).Children(func() {
					Toolbar(c, func() {
						BackButton(c, s.router)
						ForwardButton(c, s.router)
						Toggle(c, &s.bold, "Bold")
						if Button(c, "Share").Clicked() {
							s.shares++
						}
					}).Label("Format")
				})
			})
		}, 400, 200)
		tester.Frame()
		if !tester.HasText("Back") || !tester.HasText("More") {
			t.Fatalf("%v: a narrow toolbar shows %v, want Back and its menu", p, tester.Texts())
		}
		if err := tester.Click("More"); err != nil {
			t.Fatalf("%v: click the toolbar's menu: %v", p, err)
		}
		tester.Frame()
		menu := strings.Join(tester.Menu(), ",")
		if !strings.Contains(menu, "Bold") || !strings.Contains(menu, "Share") || strings.Contains(menu, "Back") {
			t.Fatalf("%v: the toolbar's menu lists %q, want Bold and Share", p, menu)
		}
		if err := tester.ChooseMenuItem("Share"); err != nil {
			t.Fatalf("%v: choose Share: %v", p, err)
		}
		tester.Frame()
		if s.shares != 1 {
			t.Fatalf("%v: choosing Share in the menu clicked it %d times", p, s.shares)
		}
	}
}

// TestHistoryButtonsLookTheSameInAToolbar checks on every platform that a back button under the pointer is drawn the same inside a Toolbar as outside one, where package ui's toolbar would give it MyGo's hover fill: on macOS none, as an NSButton does not track the pointer.
func TestHistoryButtonsLookTheSameInAToolbar(t *testing.T) {
	for _, p := range Platforms {
		var seen [2][2]ui.Color
		for i, inBar := range []bool{false, true} {
			s := newNavState()
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					ui.Column(c).Fill().Children(func() {
						if inBar {
							Toolbar(c, func() { BackButton(c, s.router) })
						} else {
							ui.Row(c).Children(func() { BackButton(c, s.router) })
						}
					})
				})
			}, 200, 100)
			tester.Frame()
			r, ok := tester.Find("Back")
			if !ok {
				t.Fatalf("%v: Back is missing", p)
			}
			at := func() ui.Color {
				px := tester.Image().RGBAAt(int(r.X+4), int(r.Y+r.H/2))
				return ui.RGB(px.R, px.G, px.B)
			}
			seen[i][0] = at()
			tester.Move(r.X+r.W/2, r.Y+r.H/2)
			tester.Frame()
			seen[i][1] = at()
		}
		if seen[0] != seen[1] {
			t.Fatalf("%v: a back button outside a toolbar is %v and %v under the pointer, inside one %v and %v", p, seen[0][0], seen[0][1], seen[1][0], seen[1][1])
		}
		if p == Mac && seen[1][0] != seen[1][1] {
			t.Fatalf("macOS: a back button changes under the pointer, from %v to %v", seen[1][0], seen[1][1])
		}
		if p != Mac && seen[1][0] == seen[1][1] {
			t.Fatalf("%v: a back button does not change under the pointer", p)
		}
	}
}

// TestFluentSidebarSelectionIndicator checks that on Windows the chosen item of a sidebar has its accent pill at its start, at the left left to right and at the right right to left, and that the items not chosen have none.
func TestFluentSidebarSelectionIndicator(t *testing.T) {
	for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
		for _, dark := range []bool{false, true} {
			s := newNavState()
			var accent ui.Color
			view := navView(Windows, s)
			tester := ui.NewTester(func(c *ui.Context) {
				view(c)
				Preview(c, Windows, func() { accent = Theme(c, Windows).Accent })
			}, 800, 420)
			tester.SetDark(dark)
			tester.SetDirection(dir)
			tester.Frame()
			sidebar, _ := tester.Find("Folders")
			// pill reports whether the row of an item's text has the accent between the text and the sidebar's edge on its start side.
			pill := func(text string) bool {
				r, ok := tester.Find(text)
				if !ok {
					t.Fatalf("%q is missing", text)
				}
				img, y := tester.Image(), int(r.Y+r.H/2)
				from, to := int(sidebar.X), int(r.X)
				if dir == ui.RightToLeft {
					from, to = int(r.X+r.W), int(sidebar.X+sidebar.W)
				}
				for x := from; x < to; x++ {
					px := img.RGBAAt(x, y)
					if ui.RGB(px.R, px.G, px.B) == accent {
						return true
					}
				}
				return false
			}
			if !pill("Inbox") {
				t.Fatalf("dir %v, dark %v: the chosen item has no accent pill at its start", dir, dark)
			}
			if pill("Sent") {
				t.Fatalf("dir %v, dark %v: an item not chosen has an accent pill", dir, dark)
			}
		}
	}
}

// TestSidebarItemColors checks the fill of the chosen item on Windows (SubtleFillColorSecondary) and GNOME (currentColor at 10% over sidebar_bg_color), and that on Windows an item not chosen takes SubtleFillColorTertiary while pressed, as NavigationViewItemBackgroundPressed.
func TestSidebarItemColors(t *testing.T) {
	for _, p := range []Platform{Windows, Linux} {
		for _, dark := range []bool{false, true} {
			s := newNavState()
			var th *ui.Theme
			view := navView(p, s)
			tester := ui.NewTester(func(c *ui.Context) {
				view(c)
				Preview(c, p, func() { th = Theme(c, p) })
			}, 800, 420)
			tester.SetDark(dark)
			tester.Frame()
			// fill is the color just after an item's text, inside the item, away from the badge.
			fill := func(text string) ui.Color {
				r, ok := tester.Find(text)
				if !ok {
					t.Fatalf("%q is missing", text)
				}
				px := tester.Image().RGBAAt(int(r.X-2), int(r.Y+r.H/2))
				return ui.RGB(px.R, px.G, px.B)
			}
			want := adwInkAlpha(th, adwSidebarBg(th), 0.10)
			if p == Windows {
				want, _ = fluentSubtle(th)
			}
			if got := fill("Inbox"); got != want {
				t.Fatalf("%v dark %v: the chosen item is %v, want %v", p, dark, got, want)
			}
			if p != Windows {
				continue
			}
			r, _ := tester.Find("Sent")
			tester.Press(r.X+r.W/2, r.Y+r.H/2)
			tester.Frame()
			_, pressed := fluentSubtle(th)
			if got := fill("Sent"); got != pressed {
				t.Fatalf("Windows dark %v: a pressed item is %v, want %v", dark, got, pressed)
			}
			tester.Release(r.X+r.W/2, r.Y+r.H/2)
			tester.Frame()
		}
	}
}

// TestSidebarsSideBySideKeepTheirOwnFocus checks that of two sidebars in one window, as the platforms side by side in a gallery, only the one with the focus draws its chosen item in the accent color on macOS.
func TestSidebarsSideBySideKeepTheirOwnFocus(t *testing.T) {
	left, right := "a", "b"
	tester := ui.NewTester(func(c *ui.Context) {
		Preview(c, Mac, func() {
			Use(c)
			ui.Row(c).Fill().AlignItems(ui.Stretch).Children(func() {
				Sidebar(c, &left, func() {
					SidebarItem(c, "a", nil, "Left A")
					SidebarItem(c, "b", nil, "Left B")
				}).Width(200).Label("Left")
				Sidebar(c, &right, func() {
					SidebarItem(c, "a", nil, "Right A")
					SidebarItem(c, "b", nil, "Right B")
				}).Width(200).Label("Right")
			})
		})
	}, 400, 200)
	tester.Frame()
	if err := tester.Click("Left A"); err != nil {
		t.Fatal(err)
	}
	tester.Move(390, 190)
	tester.Frame()
	at := func(text string) ui.Color {
		r, _ := tester.Find(text)
		px := tester.Image().RGBAAt(int(r.X-3), int(r.Y+r.H/2))
		return ui.RGB(px.R, px.G, px.B)
	}
	if got := at("Left A"); got != ui.Hex("#0070f5") {
		t.Fatalf("the focused sidebar's chosen item is %v, want the accent", got)
	}
	if got := at("Right B"); got != ui.Hex("#dcdcdc") {
		t.Fatalf("the other sidebar's chosen item is %v, want gray", got)
	}
}

// TestNavigationWidgetsTakeEmptyInput checks on every platform that a path of no items or of an empty one, items and sections outside a Sidebar, an empty title and an empty sidebar build without panicking.
func TestNavigationWidgetsTakeEmptyInput(t *testing.T) {
	for _, p := range Platforms {
		chosen, selected := -1, ""
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				Breadcrumbs(c, nil, &chosen)
				Breadcrumbs(c, []string{""}, &chosen)
				Breadcrumbs(c, []string{"Only"}, &chosen)
				SidebarItem(c, "loose", nil, "Loose")
				SidebarSection(c, "", nil, func() { SidebarItem(c, "in", nil, "In") })
				Sidebar(c, &selected, func() {}).Width(100)
			})
		}, 400, 300)
		tester.Frame()
		if !tester.HasText("Only") || !tester.HasText("Loose") || !tester.HasText("In") {
			t.Fatalf("%v: shows %v", p, tester.Texts())
		}
		if err := tester.Click("Only"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		if chosen != -1 {
			t.Fatalf("%v: a click on a path's only item chose %d", p, chosen)
		}
	}
}

// TestScreenshotsNavigation writes the navigation widgets of every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names: the sidebar without the focus, then with it after a click, beside a toolbar and a path. Without it the test only builds the frames.
func TestScreenshotsNavigation(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			for _, focus := range []bool{false, true} {
				s := newNavState()
				s.selected = "sent"
				tester := ui.NewTester(navView(p, s), 560, 260)
				tester.SetScale(2)
				tester.SetDark(dark)
				tester.Frame()
				suffix := ""
				if focus {
					if err := tester.Click("Sent"); err != nil {
						t.Fatal(err)
					}
					tester.Move(540, 250)
					tester.Frame()
					suffix = "-focused"
				}
				for _, text := range []string{"Mailboxes", "Sent", "Ada", "Share"} {
					if !tester.HasText(text) {
						t.Fatalf("%v: %q is missing", p, text)
					}
				}
				writeShot(t, dir, fmt.Sprintf("breeze-nav-%s-%s%s.png", strings.ToLower(p.String()), mode, suffix), tester)
			}
		}
	}
}
