package breeze

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// The values these tests pin were measured from libadwaita 1.10.0's rendering (GTK 4.24.1, Adwaita Sans) by tools/adwaita-reference, its inputs and rows sets, against TestScreenshotsInputs and TestScreenshotsRows.

// TestAdwPasswordFocusRingDark checks the focus ring of a GtkPasswordEntry in the dark theme, which libadwaita 1.10.0 renders (92, 132, 185) over the entry's fill: GTK mixes the dark theme's accent_color, outside sRGB, at 50% without clamping it, bluer than the clamped #81d0ff at 50% (93, 132, 158).
func TestAdwPasswordFocusRingDark(t *testing.T) {
	secret := ""
	tester := ui.NewTester(func(c *ui.Context) {
		Preview(c, Linux, func() {
			Use(c)
			ui.Box(c).Fill().Padding(16).AlignItems(ui.Start).Children(func() {
				PasswordField(c, &secret).Label("Secret").Width(180)
			})
		})
	}, 300, 100)
	tester.SetScale(2)
	tester.SetDark(true)
	tester.Frame()
	if err := tester.Click("Secret"); err != nil {
		t.Fatal(err)
	}
	tester.Type("hunter2")
	tester.Frame()
	r, ok := tester.Find("Secret")
	if !ok {
		t.Fatal("the password field is missing")
	}
	// The ring's 2 pixels lie inside the entry's edge: at 2x, the 4 device pixels from its left edge, on its middle row, away from the corners and the text.
	px := tester.Image().RGBAAt(int(r.X*2)+2, int((r.Y+r.H/2)*2))
	got := ui.RGB(px.R, px.G, px.B)
	if want := ui.RGB(92, 132, 185); adwColorDiff(got, want) > 2.0/255 {
		t.Fatalf("the focused password field's ring is %v, want libadwaita's %v", got, want)
	}
}

// TestAdwPreferencesGroupHeader checks that on GNOME a Fieldset of rows lays its rows out as an AdwPreferencesGroup with a title and no description: a header 34 pixels high (box.single-line) and the boxed list 6 below it, its first row 40 below the group's top, as libadwaita 1.10.0 lays out the rows set's "Privacy" group (group at y 16, its list at y 56), whatever the legend's font.
func TestAdwPreferencesGroupHeader(t *testing.T) {
	location := true
	tester := ui.NewTester(func(c *ui.Context) {
		Preview(c, Linux, func() {
			Use(c)
			ui.Column(c).Fill().Padding(16).Children(func() {
				ui.Column(c).Width(300).Children(func() {
					Fieldset(c, "Privacy", func() {
						SwitchRow(c, "Location", "Apps may ask for it", &location)
						ActionRow(c, "Camera", "2 apps", nil)
					})
				})
			})
		})
	}, 400, 260)
	tester.Frame()
	row, ok := tester.Find("Location")
	if !ok {
		t.Fatal("the switch row is missing")
	}
	if got := row.Y - 16; got < 39.5 || got > 40.5 {
		t.Fatalf("the group's first row is %v below its top, want libadwaita's 40", got)
	}
}

// TestAdwNestedRowLines checks the lines between the rows an AccordionItem holds on GNOME, as between AdwExpanderRow's nested rows: card_shade_color over the nested list, 229 in the light theme and 27 in the dark one, where the line below the header, over the card, is 237 and 33.
func TestAdwNestedRowLines(t *testing.T) {
	for _, dark := range []bool{false, true} {
		open, metered := true, false
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Linux, func() {
				Use(c)
				ui.Column(c).Fill().Padding(16).Children(func() {
					ui.Column(c).Width(300).Children(func() {
						Accordion(c, func() {
							AccordionItem(c, "Network", &open, func() {
								ActionRow(c, "Proxy", "Automatic", nil)
								SwitchRow(c, "Metered", "", &metered)
							})
						})
					})
				})
			})
		}, 400, 300)
		tester.SetScale(2)
		tester.SetDark(dark)
		tester.Frame()
		line := func(name string) ui.Color {
			r, ok := tester.Find(name)
			if !ok {
				t.Fatalf("dark %v: %q is missing", dark, name)
			}
			// The line is the row's top pixel, 2 device pixels at 2x, read in its middle, near the row's start, left of its title.
			px := tester.Image().RGBAAt(int(r.X*2)+8, int(r.Y*2+1))
			return ui.RGB(px.R, px.G, px.B)
		}
		header, between := line("Proxy"), line("Metered")
		wantHeader, wantBetween := pick(dark, ui.RGB(237, 237, 237), ui.RGB(33, 33, 37)), pick(dark, ui.RGB(229, 229, 229), ui.RGB(27, 27, 31))
		if adwColorDiff(header, wantHeader) > 2.0/255 {
			t.Errorf("dark %v: the line below the header is %v, want libadwaita's %v", dark, header, wantHeader)
		}
		if adwColorDiff(between, wantBetween) > 2.0/255 {
			t.Errorf("dark %v: the line between the nested rows is %v, want libadwaita's %v", dark, between, wantBetween)
		}
	}
}
