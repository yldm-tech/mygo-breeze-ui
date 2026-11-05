package breeze

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

// idiomTestSidebarIcon and idiomTestShareIcon are 16-point template images for the image-only toolbar items of the tests.
var (
	idiomTestSidebarIcon = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><rect x="1.5" y="2.5" width="13" height="11" rx="2" fill="none" stroke="#000" stroke-width="1.3"/><path d="M6 2.5v11" stroke="#000" stroke-width="1.3"/></svg>`))
	idiomTestShareIcon   = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><path d="M8 1.5v9M4.8 4.5 8 1.5l3.2 3M5 7H3.5v7.5h9V7H11" fill="none" stroke="#000" stroke-width="1.3" stroke-linejoin="round"/></svg>`))
)

// idiomState is what idiomView's widgets edit.
type idiomState struct {
	router              *ui.Router
	wrap, plain, side   bool
	clicks              map[string]int
	city, font          string
	crumb, crumbChanged int
	chosen              string
}

func newIdiomState() *idiomState {
	r := ui.NewRouter("/notes")
	r.Push("/notes/42")
	return &idiomState{router: r, wrap: true, side: true, crumb: -1, font: "Serif", clicks: map[string]int{}}
}

// idiomView lays out the group's widgets on platform p: a toolbar holding the history buttons, a text button, two image-only buttons side by side, toggles on and off, an image-only toggle, a menu button, a popover menu button and a primary button; below it a path, an autocomplete, a combobox and the menu buttons outside a toolbar.
func idiomView(p Platform, s *idiomState) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			count := func(name string, e ui.Element) ui.Element {
				if e.Clicked() {
					s.clicks[name]++
				}
				return e
			}
			icon := func(label string, svg *ui.SVG) ui.Element {
				b := Button(c, "").Label(label)
				b.Children(func() { ui.Icon(c, svg).Size(16, 16) })
				return count(label, b)
			}
			menu := func(m *ui.Menu) {
				for _, it := range []string{"Duplicate", "Rename", "Delete"} {
					if m.Item(it).Chosen() {
						s.chosen = it
					}
				}
			}
			ui.Column(c).Fill().Padding(12).Gap(16).Children(func() {
				Toolbar(c, func() {
					BackButton(c, s.router)
					ForwardButton(c, s.router)
					count("Bold", Button(c, "Bold"))
					icon("Sidebar", idiomTestSidebarIcon)
					icon("Share icon", idiomTestShareIcon)
					Toggle(c, &s.wrap, "Wrap")
					Toggle(c, &s.plain, "Plain")
					side := Toggle(c, &s.side, "").Label("Inspector")
					side.Children(func() { ui.Icon(c, idiomTestSidebarIcon).Size(16, 16) })
					MenuButton(c, "Share", menu)
					PopoverMenuButton(c, "More", menu)
					ui.Spacer(c)
					count("Save", PrimaryButton(c, "Save"))
				}).Label("Format")
				bc := Breadcrumbs(c, []string{"Home", "Users", "Ada"}, &s.crumb).Label("Path")
				if bc.Changed() {
					s.crumbChanged++
				}
				ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
					Autocomplete(c, &s.city, []string{"Berlin", "Bern", "Boston", "Brno"}).Placeholder("City").Label("City").Width(200)
					Combobox(c, &s.font, []string{"Serif", "Sans", "Mono"}).Label("Font").Width(200)
					MenuButton(c, "Sort", menu)
					PopoverMenuButton(c, "Actions", menu)
				})
			})
		})
	}
}

// TestScreenshotsWave2Idioms writes the group's widgets on every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names: at rest, with the autocomplete's suggestions open, with the toolbar's popover menu open, with the pointer over the path's current item and pressing it, and with a toolbar item pressed. Without it the test only builds the frames.
func TestScreenshotsWave2Idioms(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			name := strings.ToLower(p.String())
			shot := func(what string, tester *ui.Tester) {
				writeShot(t, dir, fmt.Sprintf("idioms-%s-%s-%s.png", what, name, mode), tester)
			}
			fresh := func() *ui.Tester {
				tester := ui.NewTester(idiomView(p, newIdiomState()), 1000, 300)
				tester.SetScale(2)
				tester.SetDark(dark)
				tester.Frame()
				tester.Frame()
				for _, want := range []string{"Bold", "Wrap", "Share", "More", "Save", "Ada"} {
					if !tester.HasText(want) {
						t.Fatalf("%v: %q is missing: %v", p, want, tester.Texts())
					}
				}
				return tester
			}
			tester := fresh()
			shot("rest", tester)

			r, ok := tester.Find("City")
			if !ok {
				t.Fatalf("%v: the autocomplete is missing", p)
			}
			tester.ClickAt(r.X+r.W/2, r.Y+r.H/2)
			tester.Frame()
			tester.Type("B")
			tester.Frame()
			tester.Key(0, ui.KeyDown)
			tester.Frame()
			if !tester.HasText("Boston") {
				t.Fatalf("%v: the suggestions do not show", p)
			}
			shot("suggestions", tester)

			tester = fresh()
			if err := tester.Click("More"); err != nil {
				t.Fatalf("%v: click More: %v", p, err)
			}
			tester.Frame()
			if !tester.HasText("Rename") {
				t.Fatalf("%v: the popover menu does not show", p)
			}
			shot("menu", tester)

			tester = fresh()
			r, _ = tester.Find("Ada")
			tester.Move(r.X+r.W/2, r.Y+r.H/2)
			tester.Frame()
			shot("crumb-hover", tester)
			tester.Press(r.X+r.W/2, r.Y+r.H/2)
			tester.Frame()
			shot("crumb-pressed", tester)
			tester.Release(r.X+r.W/2, r.Y+r.H/2)
			tester.Frame()

			for _, item := range []string{"Bold", "Share icon"} {
				tester = fresh()
				r, ok = tester.Find(item)
				if !ok {
					t.Fatalf("%v: %q is missing", p, item)
				}
				tester.Press(r.X+r.W/2, r.Y+r.H/2)
				tester.Frame()
				shot("pressed-"+strings.ReplaceAll(strings.ToLower(item), " ", "-"), tester)
				tester.Release(r.X+r.W/2, r.Y+r.H/2)
				tester.Frame()
			}
		}
	}
}
