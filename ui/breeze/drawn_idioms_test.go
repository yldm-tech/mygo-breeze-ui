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

// idiomAutocompleteRun is what a script of keys and clicks did to an autocomplete: its value, how many frames reported Changed and Submitted, the caret, whether the suggestions showed and whether the input had the focus, after each step.
type idiomAutocompleteRun struct {
	steps []string
}

// idiomRunAutocomplete runs the same script on an autocomplete built by build, on platform p in direction dir, and returns what each step left.
func idiomRunAutocomplete(t *testing.T, p Platform, dir ui.Direction, disabled bool, build func(c *ui.Context, value *string) ui.Element) idiomAutocompleteRun {
	t.Helper()
	value := ""
	changed, submitted, caret := 0, 0, 0
	tester := ui.NewTester(func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Column(c).Padding(16).Children(func() {
				in := build(c, &value).Placeholder("City").Label("City").Width(220).Disabled(disabled)
				if in.Changed() {
					changed++
				}
				if in.Submitted() {
					submitted++
				}
				_, caret = in.TextSelection()
				ui.Button(c, "Next")
			})
		})
	}, 480, 360)
	tester.SetDirection(dir)
	tester.Frame()
	var run idiomAutocompleteRun
	note := func(step string) {
		tester.Frame()
		run.steps = append(run.steps, fmt.Sprintf("%s: value %q changed %d submitted %d caret %d berlin %v boston %v focused %v", step, value, changed, submitted, caret, tester.HasText("Berlin"), tester.HasText("Boston"), tester.Focused("City")))
	}
	r, ok := tester.Find("City")
	if !ok {
		t.Fatalf("%v: the autocomplete is not named City", p)
	}
	tester.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	note("click")
	tester.Type("Ber")
	note("type Ber")
	tester.Key(0, ui.KeyEnter)
	note("Enter on none")
	tester.Key(0, ui.KeyDown)
	note("Down")
	tester.Key(0, ui.KeyDown)
	note("Down again")
	tester.Key(0, ui.KeyUp)
	note("Up")
	tester.Key(0, ui.KeyEnter)
	note("Enter on the first")
	for range 5 {
		tester.Key(0, ui.KeyBackspace)
	}
	tester.Type("o")
	note("retype Bo")
	tester.Key(0, ui.KeyEscape)
	note("Escape")
	tester.Type("s")
	note("type s")
	if err := tester.Click("Boston"); err == nil {
		note("click Boston")
	} else {
		note("no Boston to click")
	}
	tester.Key(0, ui.KeyTab)
	note("Tab away")
	return run
}

// TestIdiomAutocompleteKeepsMyGosBehavior runs one script of clicks and keys on ui.Autocomplete and on this package's, which is built on ui.ComboboxBase with the platform's menu, on every platform, left to right and right to left, enabled and disabled, and checks that each step leaves the same value, the same Changed and Submitted reports, the same caret, the same suggestions showing and the same focus: the input is named by Label, typing filters without highlighting the first suggestion (Enter on none is Submitted), the arrows move, Enter takes the suggestion with the caret after it, Escape closes the list, a click takes a suggestion, and Tab leaves.
func TestIdiomAutocompleteKeepsMyGosBehavior(t *testing.T) {
	suggestions := []string{"Berlin", "Bern", "Boston", "Brno"}
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			for _, disabled := range []bool{false, true} {
				want := idiomRunAutocomplete(t, p, dir, disabled, func(c *ui.Context, value *string) ui.Element { return ui.Autocomplete(c, value, suggestions) })
				got := idiomRunAutocomplete(t, p, dir, disabled, func(c *ui.Context, value *string) ui.Element { return Autocomplete(c, value, suggestions) })
				for i := range want.steps {
					if got.steps[i] != want.steps[i] {
						t.Fatalf("%v %v disabled %v: after %d steps\n got  %s\n want %s", p, dir, disabled, i+1, got.steps[i], want.steps[i])
					}
				}
				if !disabled && !strings.Contains(got.steps[1], "berlin true") {
					t.Fatalf("%v: typing showed no suggestions: %s", p, got.steps[1])
				}
				if !disabled && !strings.Contains(got.steps[2], "submitted 1") {
					t.Fatalf("%v: Enter without the arrows was not Submitted: %s", p, got.steps[2])
				}
			}
		}
	}
}

// TestIdiomAutocompleteDrawsThePlatformsMenu checks that the suggestions show in the platform's menu panel: macMenu's, fluentMenu's or adwMenu's background under the first suggestion's row, not MyGo's panel in the window's background.
func TestIdiomAutocompleteDrawsThePlatformsMenu(t *testing.T) {
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			city := ""
			var th *ui.Theme
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					th = c.Theme()
					ui.Column(c).Padding(16).Children(func() {
						Autocomplete(c, &city, []string{"Berlin", "Bern"}).Label("City").Width(220)
					})
				})
			}, 480, 360)
			tester.SetDark(dark)
			tester.Frame()
			r, _ := tester.Find("City")
			tester.ClickAt(r.X+r.W/2, r.Y+r.H/2)
			tester.Frame()
			tester.Type("Be")
			tester.Frame()
			row, ok := tester.Find("Bern")
			if !ok {
				t.Fatalf("%v: the suggestions do not show", p)
			}
			want := pick(th.Dark, ui.Hex("#f6f6f6"), ui.Hex("#2b2b2b"))
			switch p {
			case Windows:
				want = pick(th.Dark, ui.Hex("#f9f9f9"), ui.Hex("#2c2c2c"))
			case Linux:
				want = pick(th.Dark, white, ui.Hex("#36363a"))
			}
			px := tester.Image().RGBAAt(int(r.X+r.W-14), int(row.Y+row.H/2))
			if got := ui.RGB(px.R, px.G, px.B); got != want {
				t.Fatalf("%v dark %v: the suggestions' panel is %v, want %v", p, dark, got, want)
			}
		}
	}
}

// idiomToolbarState is what idiomToolbarView's controls edit.
type idiomToolbarState struct {
	wrap, side, disabled bool
	clicks               map[string]int
	chosen               string
	router               *ui.Router
}

// idiomToolbarView lays out a toolbar of every control this package draws as a toolbar item, with spacer between the two image-only buttons at its end when spaced.
func idiomToolbarView(p Platform, s *idiomToolbarState, spaced bool) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			count := func(name string, e ui.Element) {
				if e.Clicked() {
					s.clicks[name]++
				}
			}
			menu := func(m *ui.Menu) {
				for _, it := range []string{"Duplicate", "Rename"} {
					if m.Item(it).Chosen() {
						s.chosen = it
					}
				}
			}
			ui.Column(c).Fill().Padding(8).Children(func() {
				Toolbar(c, func() {
					BackButton(c, s.router)
					ForwardButton(c, s.router)
					count("Bold", Button(c, "Bold").Disabled(s.disabled))
					Toggle(c, &s.wrap, "Wrap")
					MenuButton(c, "Share", menu)
					PopoverMenuButton(c, "More", menu)
					count("Save", PrimaryButton(c, "Save"))
					side := Toggle(c, &s.side, "").Label("Inspector")
					side.Children(func() { ui.Icon(c, idiomTestSidebarIcon).Size(16, 16) })
					a := Button(c, "").Label("Sidebar")
					a.Children(func() { ui.Icon(c, idiomTestSidebarIcon).Size(16, 16) })
					count("Sidebar", a)
					if spaced {
						ui.Spacer(c)
					}
					b := Button(c, "").Label("Upload")
					b.Children(func() { ui.Icon(c, idiomTestShareIcon).Size(16, 16) })
					count("Upload", b)
				}).Label("Format")
			})
		})
	}
}

func newIdiomToolbarState() *idiomToolbarState {
	r := ui.NewRouter("/notes")
	r.Push("/notes/42")
	return &idiomToolbarState{clicks: map[string]int{}, router: r}
}

// TestIdiomToolbarItemsKeepTheirBehavior checks on every platform, left to right and right to left, that the controls drawn as toolbar items keep MyGo's behavior: a click clicks a button, toggles a toggle, opens a menu button's native menu and a popover menu button's drawn menu, whose item chosen is reported; the toolbar is one stop of Tab and the arrows move the focus among its controls, Enter clicking the one focused; the history buttons go back; a disabled button takes no click; and the controls are found by their labels.
func TestIdiomToolbarItemsKeepTheirBehavior(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			s := newIdiomToolbarState()
			tester := ui.NewTester(idiomToolbarView(p, s, false), 1100, 300)
			tester.SetDirection(dir)
			tester.Frame()
			tester.Frame()
			for _, name := range []string{"Back", "Forward", "Bold", "Wrap", "Share", "More", "Save", "Inspector", "Sidebar", "Upload"} {
				if _, ok := tester.Find(name); !ok {
					t.Fatalf("%v %v: %q is missing: %v", p, dir, name, tester.Texts())
				}
			}
			click := func(name string) {
				t.Helper()
				if err := tester.Click(name); err != nil {
					t.Fatalf("%v %v: click %q: %v", p, dir, name, err)
				}
				tester.Frame()
			}
			for _, name := range []string{"Bold", "Save", "Sidebar", "Upload"} {
				click(name)
				if s.clicks[name] != 1 {
					t.Fatalf("%v %v: a click on %q clicked it %d times", p, dir, name, s.clicks[name])
				}
			}
			click("Wrap")
			click("Inspector")
			if !s.wrap || !s.side {
				t.Fatalf("%v %v: the toggles are %v and %v after a click", p, dir, s.wrap, s.side)
			}
			click("Share")
			if got := strings.Join(tester.Menu(), ","); got != "Duplicate,Rename" {
				t.Fatalf("%v %v: the menu button's menu is %q", p, dir, got)
			}
			if err := tester.ChooseMenuItem("Rename"); err != nil {
				t.Fatalf("%v %v: choose Rename: %v", p, dir, err)
			}
			tester.Frame()
			if s.chosen != "Rename" {
				t.Fatalf("%v %v: choosing in the native menu gave %q", p, dir, s.chosen)
			}
			click("More")
			if !tester.HasText("Duplicate") {
				t.Fatalf("%v %v: the popover menu does not show", p, dir)
			}
			click("Duplicate")
			if s.chosen != "Duplicate" || tester.HasText("Duplicate") {
				t.Fatalf("%v %v: choosing in the drawn menu gave %q, shows %v", p, dir, s.chosen, tester.HasText("Duplicate"))
			}
			click("Back")
			if s.router.Path() != "/notes" {
				t.Fatalf("%v %v: Back went to %q", p, dir, s.router.Path())
			}
			// The keyboard: Tab into the toolbar, the arrows along it, Enter on a button.
			click("Bold")
			if !tester.Focused("Bold") {
				t.Fatalf("%v %v: a click left Bold without the focus", p, dir)
			}
			next := ui.KeyRight
			if dir == ui.RightToLeft {
				next = ui.KeyLeft
			}
			tester.Key(0, next)
			tester.Frame()
			if !tester.Focused("Wrap") {
				t.Fatalf("%v %v: the arrow did not move the focus from Bold to Wrap", p, dir)
			}
			tester.Key(0, ui.KeySpace)
			tester.Frame()
			if s.wrap {
				t.Fatalf("%v %v: Space did not turn the focused toggle off", p, dir)
			}
			tester.Key(0, next)
			tester.Key(0, next)
			tester.Key(0, next)
			tester.Frame()
			if !tester.Focused("Save") {
				t.Fatalf("%v %v: three arrows from Wrap did not reach Save", p, dir)
			}
			tester.Key(0, ui.KeyEnter)
			tester.Frame()
			if s.clicks["Save"] != 2 {
				t.Fatalf("%v %v: Enter on Save clicked it %d times in all", p, dir, s.clicks["Save"])
			}
			s.disabled = true
			tester.Frame()
			click("Bold")
			if s.clicks["Bold"] != 2 {
				t.Fatalf("%v %v: a disabled button took a click: %d", p, dir, s.clicks["Bold"])
			}
		}
	}
}

// TestIdiomMacToolbarJoinsImageItems checks that on macOS two image-only items side by side in a toolbar touch, sharing one platter, the history buttons too, left to right and right to left, and that a spacer between them keeps them apart, as do the other platforms, which never join their items.
func TestIdiomMacToolbarJoinsImageItems(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			for _, spaced := range []bool{false, true} {
				s := newIdiomToolbarState()
				tester := ui.NewTester(idiomToolbarView(p, s, spaced), 1100, 300)
				tester.SetDirection(dir)
				tester.Frame()
				tester.Frame()
				tester.Frame()
				gap := func(a, b string) float32 {
					ra, _ := tester.Find(a)
					rb, _ := tester.Find(b)
					return max(rb.X-(ra.X+ra.W), ra.X-(rb.X+rb.W))
				}
				joined := p == Mac && !spaced
				if g := gap("Sidebar", "Upload"); joined != (g < 0.5) {
					t.Fatalf("%v %v spaced %v: the image buttons are %v apart", p, dir, spaced, g)
				}
				if g := gap("Back", "Forward"); (p == Mac) != (g < 0.5) {
					t.Fatalf("%v %v: the history buttons are %v apart", p, dir, g)
				}
				// A toggle never joins: the image toggle stays apart from the image button after it.
				if g := gap("Inspector", "Sidebar"); g < 0.5 {
					t.Fatalf("%v %v: the image toggle joined the image button", p, dir)
				}
			}
		}
	}
}

// TestIdiomToolbarItemLooks checks the faces the platforms give a toolbar's items, light and dark: on macOS the selected toggle is filled with the accent and an item at rest is white (Aqua) or (69,68,68) (Dark Aqua) at its middle; on Windows an AppBarButton at rest shows the window behind it and a checked AppBarToggleButton the accent; on GNOME an image button at rest is flat and a text button raised.
func TestIdiomToolbarItemLooks(t *testing.T) {
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			s := newIdiomToolbarState()
			s.wrap = true
			var th *ui.Theme
			view := idiomToolbarView(p, s, true)
			tester := ui.NewTester(func(c *ui.Context) {
				view(c)
				th = Theme(c, p)
			}, 1100, 300)
			tester.SetDark(dark)
			tester.Frame()
			at := func(name string, dx float32) ui.Color {
				r, ok := tester.Find(name)
				if !ok {
					t.Fatalf("%v: %q is missing", p, name)
				}
				px := tester.Image().RGBAAt(int(r.X+dx), int(r.Y+3))
				return ui.RGB(px.R, px.G, px.B)
			}
			wrap, bold, sidebar := at("Wrap", 20), at("Bold", 20), at("Sidebar", 20)
			switch p {
			case Mac:
				if want := pick(th.Dark, white, ui.RGB(69, 68, 68)); navAbsDiff(bold.R, want.R) > 8 {
					t.Fatalf("macOS dark %v: a text item at rest is %v near its top, want about %v", dark, bold, want)
				}
				if wrap.B < 200 || wrap.R > 120 {
					t.Fatalf("macOS dark %v: the selected toggle is %v, not the accent", dark, wrap)
				}
			case Windows:
				if bold != th.Background {
					t.Fatalf("Windows dark %v: an AppBarButton at rest is %v, want the window %v", dark, bold, th.Background)
				}
				if wrap != th.Accent {
					t.Fatalf("Windows dark %v: a checked AppBarToggleButton is %v, want %v", dark, wrap, th.Accent)
				}
			case Linux:
				if sidebar != th.Background {
					t.Fatalf("Linux dark %v: an image button at rest is %v, want flat on %v", dark, sidebar, th.Background)
				}
				if bold == th.Background {
					t.Fatalf("Linux dark %v: a text button in a toolbar is flat", dark)
				}
			}
		}
	}
}

// TestIdiomMenuButtonsKeepTheirBehavior checks the menu buttons outside a toolbar on every platform, left to right and right to left: a click opens the native menu or the drawn one, Enter and Down open the drawn one from the keyboard, the arrows and Enter choose in it, Escape closes it giving the focus back, a disabled one opens nothing, and the drawn menu takes the platform's menu background.
func TestIdiomMenuButtonsKeepTheirBehavior(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			chosen, disabled := "", false
			var th *ui.Theme
			menu := func(m *ui.Menu) {
				for _, it := range []string{"Name", "Date", "Size"} {
					if m.Item(it).Chosen() {
						chosen = it
					}
				}
			}
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					th = c.Theme()
					ui.Row(c).Gap(12).Padding(16).Children(func() {
						MenuButton(c, "Sort", menu).Disabled(disabled)
						PopoverMenuButton(c, "Order", menu).Disabled(disabled)
					})
				})
			}, 600, 400)
			tester.SetDirection(dir)
			tester.Frame()
			if err := tester.Click("Sort"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			if got := strings.Join(tester.Menu(), ","); got != "Name,Date,Size" {
				t.Fatalf("%v %v: the native menu is %q", p, dir, got)
			}
			tester.CloseMenu()
			if err := tester.Click("Order"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			if !tester.HasText("Date") {
				t.Fatalf("%v %v: the drawn menu does not show", p, dir)
			}
			row, _ := tester.Find("Size")
			px := tester.Image().RGBAAt(int(row.X+row.W-6), int(row.Y+row.H/2))
			bg := pick(th.Dark, ui.Hex("#f6f6f6"), ui.Hex("#2b2b2b"))
			switch p {
			case Windows:
				bg = pick(th.Dark, ui.Hex("#f9f9f9"), ui.Hex("#2c2c2c"))
			case Linux:
				bg = pick(th.Dark, white, ui.Hex("#36363a"))
			}
			if got := ui.RGB(px.R, px.G, px.B); got != bg {
				t.Fatalf("%v %v: the drawn menu's panel is %v, want %v", p, dir, got, bg)
			}
			tester.Key(0, ui.KeyEscape)
			tester.Frame()
			if tester.HasText("Date") || !tester.Focused("Order") {
				t.Fatalf("%v %v: Escape left the menu %v, the focus on Order %v", p, dir, tester.HasText("Date"), tester.Focused("Order"))
			}
			tester.Key(0, ui.KeyEnter)
			tester.Frame()
			tester.Key(0, ui.KeyDown)
			tester.Frame()
			tester.Key(0, ui.KeyEnter)
			tester.Frame()
			if chosen != "Date" {
				t.Fatalf("%v %v: Enter, Down and Enter chose %q", p, dir, chosen)
			}
			disabled, chosen = true, ""
			tester.Frame()
			_ = tester.Click("Order")
			tester.Frame()
			if tester.HasText("Date") {
				t.Fatalf("%v %v: a disabled popover menu button opened its menu", p, dir)
			}
		}
	}
}

// TestIdiomBreadcrumbCurrentItem checks that on Windows the current item of a path turns TextFillColorSecondary under the pointer and TextFillColorTertiary pressed, as BreadcrumbBar draws it, and that on every platform a click on it chooses nothing and reports no change, while a click on a link still chooses it.
func TestIdiomBreadcrumbCurrentItem(t *testing.T) {
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			chosen, changed := -1, 0
			var th *ui.Theme
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					th = c.Theme()
					ui.Column(c).Padding(16).Children(func() {
						if Breadcrumbs(c, []string{"Home", "Users", "Ada"}, &chosen).Label("Path").Changed() {
							changed++
						}
					})
				})
			}, 480, 200)
			tester.SetDark(dark)
			tester.Frame()
			r, _ := tester.Find("Ada")
			ink := func() ui.Color {
				img, best, bestSum := tester.Image(), ui.Color{}, 1000
				for y := int(r.Y); y < int(r.Y+r.H); y++ {
					for x := int(r.X); x < int(r.X+r.W); x++ {
						px := img.RGBAAt(x, y)
						sum := int(px.R) + int(px.G) + int(px.B)
						if th.Dark {
							sum = 765 - sum
						}
						if sum < bestSum {
							best, bestSum = ui.RGB(px.R, px.G, px.B), sum
						}
					}
				}
				return best
			}
			rest := ink()
			tester.Move(r.X+r.W/2, r.Y+r.H/2)
			tester.Frame()
			hover := ink()
			tester.Press(r.X+r.W/2, r.Y+r.H/2)
			tester.Frame()
			pressed := ink()
			tester.Release(r.X+r.W/2, r.Y+r.H/2)
			tester.Frame()
			if p == Windows {
				fg := pick(th.Dark, black, white)
				if want := th.Background.Mix(fg, pick(th.Dark, float32(0.62), 0.773)); hover != want && hover != th.TextMuted {
					t.Fatalf("Windows dark %v: the current item under the pointer is %v, want %v", dark, hover, want)
				}
				if want := th.Background.Mix(fg, pick(th.Dark, float32(0.447), 0.529)); pressed != want {
					t.Fatalf("Windows dark %v: the current item pressed is %v, want %v", dark, pressed, want)
				}
			} else if hover != rest || pressed != rest {
				t.Fatalf("%v dark %v: the current item changes under the pointer: %v, %v, %v", p, dark, rest, hover, pressed)
			}
			if chosen != -1 || changed != 0 {
				t.Fatalf("%v: a click on the current item chose %d, changed %d", p, chosen, changed)
			}
			if err := tester.Click("Users"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			if chosen != 1 || changed != 1 {
				t.Fatalf("%v: a click on Users chose %d, changed %d", p, chosen, changed)
			}
		}
	}
}
