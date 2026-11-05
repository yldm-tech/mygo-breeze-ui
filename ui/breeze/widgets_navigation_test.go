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

// TestNavigationWidgetsKeepTheirBehavior checks on every platform, left to right and right to left, that the restyled sidebar chooses with clicks, the arrows and typing, that its sections open and close, that a path's links choose with a click and Enter while its last item does not, that the history buttons go back and forward and stop at the history's start, and that a toolbar's controls click and take the focus with the arrows.
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
			if err := tester.Click("Ada"); err != nil {
				t.Fatalf("%s: click Ada: %v", name, err)
			}
			tester.Frame()
			if s.crumb != 1 || s.crumbChanged != 2 {
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
