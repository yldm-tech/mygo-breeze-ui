package breeze

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/egoist/mygo/ui"
)

func init() {
	for _, n := range []string{"Collapsible", "Accordion", "AccordionItem", "Fieldset", "Form"} {
		drawnWidgets[n] = true
	}
}

// layoutTester builds view on platform p without motion, so that a collapsible opens in the frame it is clicked.
func layoutTester(p Platform, w, h int, view func(c *ui.Context)) *ui.Tester {
	tester := ui.NewTester(func(c *ui.Context) { Preview(c, p, func() { view(c) }) }, w, h)
	tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tester.Frame()
	return tester
}

// TestCollapsibleKeepsItsBehavior checks the disclosure on every platform: closed it hides its content, a click opens it and Changed reports it, Space closes it and Enter opens it again, and Tab goes into the content.
func TestCollapsibleKeepsItsBehavior(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			open, changes, inner := false, 0, 0
			tester := layoutTester(p, 400, 300, func(c *ui.Context) {
				ui.Column(c).Padding(8).Children(func() {
					if Collapsible(c, "Advanced", &open, func() {
						ui.Text(c, "Inside")
						if Button(c, "Inner").Clicked() {
							inner++
						}
					}).Changed() {
						changes++
					}
				})
			})
			tester.SetDirection(dir)
			tester.Frame()
			if tester.HasText("Inside") {
				t.Fatalf("%v: closed, it shows its content", p)
			}
			if err := tester.Click("Advanced"); err != nil {
				t.Fatalf("%v: click: %v", p, err)
			}
			tester.Frame()
			if !open || changes != 1 || !tester.HasText("Inside") {
				t.Fatalf("%v: a click left it open %v after %d changes, showing %v", p, open, changes, tester.Texts())
			}
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			if !tester.Focused("Inner") {
				t.Fatalf("%v: Tab did not go into the content", p)
			}
			tester.Key(0, ui.KeySpace)
			tester.Frame()
			if inner != 1 {
				t.Fatalf("%v: Space did not click the button inside", p)
			}
			tester.Key(ui.Shift, ui.KeyTab)
			tester.Frame()
			tester.Key(0, ui.KeySpace)
			tester.Frame()
			if open || tester.HasText("Inside") {
				t.Fatalf("%v: Space did not close it", p)
			}
			tester.Key(0, ui.KeyEnter)
			tester.Frame()
			if !open || changes != 3 {
				t.Fatalf("%v: Enter left it open %v after %d changes", p, open, changes)
			}
		}
	}
}

// TestAccordionKeepsItsBehavior checks the accordion on every platform: a click opens a section, Down moves the focus to the next header, Enter opens that one, and Home goes back to the first.
func TestAccordionKeepsItsBehavior(t *testing.T) {
	for _, p := range Platforms {
		general, privacy, changes := false, false, 0
		tester := layoutTester(p, 400, 400, func(c *ui.Context) {
			ui.Column(c).Padding(8).Children(func() {
				Accordion(c, func() {
					if AccordionItem(c, "General", &general, func() { ui.Text(c, "General settings") }).Changed() {
						changes++
					}
					AccordionItem(c, "Privacy", &privacy, func() { ui.Text(c, "Privacy settings") })
				})
			})
		})
		if err := tester.Click("General"); err != nil {
			t.Fatalf("%v: click: %v", p, err)
		}
		tester.Frame()
		if !general || changes != 1 || !tester.HasText("General settings") {
			t.Fatalf("%v: a click left General open %v after %d changes", p, general, changes)
		}
		tester.Key(0, ui.KeyDown)
		tester.Frame()
		if !tester.Focused("Privacy") {
			t.Fatalf("%v: Down did not move to the next header", p)
		}
		tester.Key(0, ui.KeyEnter)
		tester.Frame()
		if !privacy || !tester.HasText("Privacy settings") {
			t.Fatalf("%v: Enter did not open Privacy", p)
		}
		tester.Key(0, ui.KeyHome)
		tester.Frame()
		if !tester.Focused("General") {
			t.Fatalf("%v: Home did not move to the first header", p)
		}
	}
}

// TestFormsKeepTheirBehavior checks fieldsets, forms and check box groups on every platform: a click on a label in a fieldset focuses its control, a disabled fieldset refuses clicks, and a check box group checks all its boxes, then none.
func TestFormsKeepTheirBehavior(t *testing.T) {
	for _, p := range Platforms {
		name, address, clicks := "", "", 0
		mail, calendar := false, true
		tester := layoutTester(p, 500, 600, func(c *ui.Context) {
			ui.Column(c).Padding(8).Gap(8).Children(func() {
				Form(c, func() {
					Field(c, "Name", func() { TextInput(c, &name) })
					Fieldset(c, "Shipping", func() {
						Field(c, "Address", func() { TextInput(c, &address) })
					})
				})
				Fieldset(c, "Locked", func() {
					if Button(c, "Order").Clicked() {
						clicks++
					}
				}).Disabled(true)
				CheckboxGroup(c, "Notifications", func() {
					Checkbox(c, &mail, "Mail")
					Checkbox(c, &calendar, "Calendar")
				})
			})
		})
		if !tester.HasText("Shipping") || !tester.HasText("Locked") {
			t.Fatalf("%v: a legend is missing", p)
		}
		if err := tester.Click("Address"); err != nil {
			t.Fatalf("%v: click the label: %v", p, err)
		}
		tester.Frame()
		tester.Type("Main St")
		tester.Frame()
		if address != "Main St" {
			t.Fatalf("%v: typing after a click on the label gave %q", p, address)
		}
		if err := tester.Click("Order"); err != nil {
			t.Fatalf("%v: click the disabled button: %v", p, err)
		}
		tester.Frame()
		if clicks != 0 {
			t.Fatalf("%v: a disabled fieldset's button took a click", p)
		}
		if err := tester.Click("Notifications"); err != nil {
			t.Fatalf("%v: click the group: %v", p, err)
		}
		tester.Frame()
		if !mail || !calendar {
			t.Fatalf("%v: the group checked mail %v, calendar %v", p, mail, calendar)
		}
		if err := tester.Click("Notifications"); err != nil {
			t.Fatalf("%v: click the group again: %v", p, err)
		}
		tester.Frame()
		if mail || calendar {
			t.Fatalf("%v: a second click left mail %v, calendar %v", p, mail, calendar)
		}
	}
}

// TestFormLabelsFollowThePlatform checks where a form puts a field's label: beside the control on macOS and GNOME, above it on Windows.
func TestFormLabelsFollowThePlatform(t *testing.T) {
	for _, p := range Platforms {
		var name string
		var label, input ui.Rect
		tester := layoutTester(p, 500, 300, func(c *ui.Context) {
			ui.Column(c).Padding(8).Children(func() {
				Form(c, func() {
					Field(c, "Full name", func() { TextInput(c, &name).Label("Input") })
				})
			})
		})
		tester.Frame()
		label, _ = tester.Find("Full name")
		input, _ = tester.Find("Input")
		// The label names the field around it, so Find returns the field: the input starts at its top when the label is beside it, a line lower when above.
		above := input.Y > label.Y+10
		if want := p == Windows; above != want {
			t.Errorf("%v: label %v, input %v: above %v, want %v", p, label, input, above, want)
		}
	}
}

// TestSplitsKeepTheirBehavior checks the splits on every platform: Tab focuses the divider and the arrows move it.
func TestSplitsKeepTheirBehavior(t *testing.T) {
	for _, p := range Platforms {
		width, height, share := float32(120), float32(100), float32(40)
		tester := layoutTester(p, 600, 400, func(c *ui.Context) {
			ui.Column(c).Fill().Gap(8).Children(func() {
				Split(c, &width, func() { ui.Text(c, "Files") }, func() { ui.Text(c, "Editor") }).Height(100)
				SplitVertical(c, &height, func() { ui.Text(c, "Top") }, func() { ui.Text(c, "Bottom") }).Height(240)
			})
		})
		tester.Frame()
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		tester.Key(0, ui.KeyRight)
		tester.Frame()
		if width <= 120 {
			t.Fatalf("%v: Right left the split at %v", p, width)
		}
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		tester.Key(0, ui.KeyDown)
		tester.Frame()
		if height <= 100 {
			t.Fatalf("%v: Down left the vertical split at %v", p, height)
		}
		percent := layoutTester(p, 600, 200, func(c *ui.Context) {
			SplitPercent(c, &share, func() { ui.Text(c, "Left") }, func() { ui.Text(c, "Right") }).Fill()
		})
		percent.Frame()
		percent.Key(0, ui.KeyTab)
		percent.Frame()
		percent.Key(0, ui.KeyLeft)
		percent.Frame()
		if share >= 40 {
			t.Fatalf("%v: Left left the percent split at %v", p, share)
		}
	}
}

// TestLayoutDisabledRefusesInput checks on every platform that a disabled disclosure and a disabled accordion section neither open on a click nor take the focus from Tab.
func TestLayoutDisabledRefusesInput(t *testing.T) {
	for _, p := range Platforms {
		advanced, general := false, false
		tester := layoutTester(p, 400, 300, func(c *ui.Context) {
			ui.Column(c).Padding(8).Gap(8).Children(func() {
				Collapsible(c, "Advanced", &advanced, func() { ui.Text(c, "Inside") }).Disabled(true)
				Accordion(c, func() {
					AccordionItem(c, "General", &general, func() { ui.Text(c, "General settings") }).Disabled(true)
				})
			})
		})
		for _, s := range []string{"Advanced", "General"} {
			if err := tester.Click(s); err != nil {
				t.Fatalf("%v: click %q: %v", p, s, err)
			}
			tester.Frame()
		}
		if advanced || general || tester.HasText("Inside") || tester.HasText("General settings") {
			t.Fatalf("%v: a click opened a disabled section: disclosure %v, accordion %v", p, advanced, general)
		}
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		if tester.Focused("Advanced") || tester.Focused("General") {
			t.Fatalf("%v: Tab focused a disabled header", p)
		}
	}
}

// TestCollapsibleMirrors checks on every platform that right to left the disclosure's label is as far from the right edge as it is from the left one left to right, as each platform mirrors its disclosures.
func TestCollapsibleMirrors(t *testing.T) {
	const width, pad = 400, 8
	for _, p := range Platforms {
		var inset [2]float32
		for i, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			open := true
			tester := layoutTester(p, width, 300, func(c *ui.Context) {
				ui.Column(c).Padding(pad).Children(func() {
					Collapsible(c, "Advanced", &open, func() { ui.Text(c, "Inside") })
				})
			})
			tester.SetDirection(dir)
			tester.Frame()
			r, ok := tester.Find("Advanced")
			if !ok {
				t.Fatalf("%v %v: no label", p, dir)
			}
			inset[i] = r.X - pad
			if dir == ui.RightToLeft {
				inset[i] = width - pad - (r.X + r.W)
			}
		}
		if d := inset[0] - inset[1]; d < -0.5 || d > 0.5 {
			t.Errorf("%v: the label is %v in from the start left to right, %v right to left", p, inset[0], inset[1])
		}
	}
}

// TestFieldsetLinesUpInAForm checks on every platform that the controls of a fieldset's fields in a Form start and end where those of the form's other fields do, as ui.Fieldset promises, and that the form leaves no state behind.
func TestFieldsetLinesUpInAForm(t *testing.T) {
	for _, p := range Platforms {
		name, address := "", ""
		tester := layoutTester(p, 500, 300, func(c *ui.Context) {
			ui.Column(c).Padding(24).Children(func() {
				Form(c, func() {
					Field(c, "Name", func() { TextInput(c, &name).Label("Name input") })
					Fieldset(c, "Shipping", func() {
						Field(c, "Address", func() { TextInput(c, &address).Label("Address input") })
					})
				})
			})
		})
		tester.Frame()
		a, okA := tester.Find("Name input")
		b, okB := tester.Find("Address input")
		if !okA || !okB {
			t.Fatalf("%v: an input is missing", p)
		}
		if dx, dw := a.X-b.X, a.W-b.W; dx < -0.5 || dx > 0.5 || dw < -0.5 || dw > 0.5 {
			t.Errorf("%v: the form's input is at %v, the fieldset's at %v", p, a, b)
		}
		if !layoutStateEmpty(&layoutForms) {
			t.Errorf("%v: a Form left its state behind", p)
		}
	}
}

// TestAccordionCountsItsOwnItems checks that an accordion inside a section counts its items in its own, leaving the outer one's count as it was, and that no state outlives the build.
func TestAccordionCountsItsOwnItems(t *testing.T) {
	for _, p := range Platforms {
		a, b, x := true, true, false
		var built int
		tester := layoutTester(p, 400, 400, func(c *ui.Context) {
			Accordion(c, func() {
				AccordionItem(c, "A", &a, func() {
					Accordion(c, func() { AccordionItem(c, "X", &x, func() {}) })
				})
				AccordionItem(c, "B", &b, func() {
					if v, ok := accordions.Load(c); ok {
						built = v.(*layoutAccordion).built
					}
				})
			})
		})
		tester.Frame()
		if built != 2 {
			t.Errorf("%v: B is item %d of the outer accordion, want 2", p, built)
		}
		if !layoutStateEmpty(&accordions) {
			t.Errorf("%v: an Accordion left its state behind", p)
		}
	}
}

// layoutStateEmpty reports whether m holds no window's state.
func layoutStateEmpty(m *sync.Map) bool {
	empty := true
	m.Range(func(any, any) bool { empty = false; return false })
	return empty
}

// layoutView lays out the widgets of the layout group for screenshots: a closed and an open disclosure, an accordion with its second section open, a form with a fieldset, a check box group and a split.
func layoutView(p Platform) func(c *ui.Context) {
	closed, opened := false, true
	general, privacy := false, true
	name, address := "Ada", ""
	mail, calendar, photos := true, false, false
	split := float32(120)
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Row(c).Fill().Padding(16).Gap(24).AlignItems(ui.Start).Children(func() {
				ui.Column(c).Width(250).Gap(12).Children(func() {
					Collapsible(c, "Closed", &closed, func() { ui.Text(c, "Hidden") })
					Collapsible(c, "Advanced", &opened, func() {
						Checkbox(c, &photos, "Verbose logging")
					})
					Accordion(c, func() {
						AccordionItem(c, "General", &general, func() { ui.Text(c, "General settings") })
						AccordionItem(c, "Privacy", &privacy, func() { Checkbox(c, &mail, "Share usage data") })
					})
				})
				ui.Column(c).Width(270).Gap(16).Children(func() {
					Form(c, func() {
						Field(c, "Name", func() { TextInput(c, &name) })
						Fieldset(c, "Shipping", func() {
							Field(c, "Address", func() { TextInput(c, &address).Placeholder("Street") })
						})
					})
					CheckboxGroup(c, "Notifications", func() {
						Checkbox(c, &mail, "Mail")
						Checkbox(c, &calendar, "Calendar")
					})
					Split(c, &split, func() {
						ui.Text(c, "Sidebar").Padding(8)
					}, func() {
						ui.Text(c, "Content").Padding(8)
					}).Height(60).Border(1, c.Theme().Border)
				})
			})
		})
	}
}

// TestScreenshotsLayout writes the layout group of every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names. Without it the test only builds the frames.
func TestScreenshotsLayout(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			tester := ui.NewTester(layoutView(p), 600, 470)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tester.Frame()
			tester.Frame()
			for _, s := range []string{"Advanced", "Verbose logging", "Privacy", "Share usage data", "Shipping", "Notifications", "Sidebar"} {
				if !tester.HasText(s) {
					t.Fatalf("%v: %q is missing", p, s)
				}
			}
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			writeShot(t, dir, fmt.Sprintf("breeze-layout-%s-%s.png", strings.ToLower(p.String()), mode), tester)
		}
	}
}
