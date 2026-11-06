package breeze

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// The values these tests pin were measured from libadwaita 1.10.0's rendering (GTK 4.24.1, Adwaita Sans) by tools/adwaita-reference, its inputs, rows and overcard sets, against TestScreenshotsInputs, TestScreenshotsRows and TestAdwControlsOverCard.

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

// TestAdwButtonOverCard checks that on GNOME a button shows what is behind it, as libadwaita's translucent $button_color does: in an AdwActionRow of a boxed list its face is 235 in the light theme and 72 in the dark one, the card's color mixed, where over the window it is 230 and 57.
func TestAdwButtonOverCard(t *testing.T) {
	for _, dark := range []bool{false, true} {
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Linux, func() {
				Use(c)
				ui.Column(c).Fill().Padding(16).Gap(16).Children(func() {
					ui.Column(c).Width(300).Children(func() {
						Card(c, func() {
							ActionRow(c, "Wallpaper", "Choose a picture", func() { Button(c, "Open") })
						})
					})
					ui.Row(c).Children(func() { Button(c, "Cancel") })
				})
			})
		}, 400, 200)
		tester.SetDark(dark)
		tester.Frame()
		face := func(name string) ui.Color {
			r, ok := tester.Find(name)
			if !ok {
				t.Fatalf("dark %v: %q is missing", dark, name)
			}
			// Inside the button's padding of 17, before its label, on its middle row.
			return adwPixel(tester, r.X+6, r.Y+r.H/2)
		}
		for _, k := range []struct {
			name string
			want ui.Color
		}{
			{"Open", pick(dark, ui.RGB(235, 235, 235), ui.RGB(72, 72, 75))},
			{"Cancel", pick(dark, ui.RGB(230, 230, 231), ui.RGB(57, 57, 60))},
		} {
			if got := face(k.name); adwColorDiff(got, k.want) > 2.0/255 {
				t.Errorf("dark %v: %q's face is %v, want libadwaita's %v", dark, k.name, got, k.want)
			}
		}
	}
}

// TestAdwControlsOverCard checks that on GNOME an entry and a split button show what is behind them, as libadwaita's translucent fills do, measured from libadwaita 1.10.0 (GTK 4.24.1) by tools/adwaita-reference's overcard set: in AdwActionRows over a boxed list's card a GtkPasswordEntry's fill is 235 in the light theme and 72 in the dark one and its placeholder, .dimmed, 132 and 172; an AdwSplitButton's face is 235 and 72 and its seam 217 and 91. Over the window, in the inputs set, the placeholder is 129 and 166, where MyGo's TextMuted read 140 and 156.
func TestAdwControlsOverCard(t *testing.T) {
	for _, dark := range []bool{false, true} {
		empty, alone := "", ""
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Linux, func() {
				Use(c)
				ui.Column(c).Fill().Padding(16).Gap(16).Children(func() {
					ui.Column(c).Width(300).Children(func() {
						Card(c, func() {
							ActionRow(c, "Share", "", func() { SplitButton(c, "Send", func(m *ui.Menu) { m.Item("Send Later") }) })
							ActionRow(c, "Secret", "", func() { PasswordField(c, &empty).Placeholder("Password").Label("Over the card").Width(160) })
						})
					})
					ui.Row(c).Children(func() { PasswordField(c, &alone).Placeholder("Password").Label("Over the window").Width(160) })
				})
			})
		}, 400, 220)
		tester.SetScale(2)
		tester.SetDark(dark)
		tester.Frame()
		img := tester.Image()
		find := func(name string) ui.Rect {
			r, ok := tester.Find(name)
			if !ok {
				t.Fatalf("dark %v: %q is missing", dark, name)
			}
			return r
		}
		at := func(x, y float32) ui.Color {
			px := img.RGBAAt(int(x*2), int(y*2))
			return ui.RGB(px.R, px.G, px.B)
		}
		// The color farthest from the fill in a box, the ink where it covers whole device pixels: the placeholder's glyphs, or the seam.
		ink := func(x0, y0, x1, y1 float32) ui.Color {
			best, bestLum := ui.Color{}, float32(-1)
			for y := int(y0 * 2); y < int(y1*2); y++ {
				for x := int(x0 * 2); x < int(x1*2); x++ {
					px := img.RGBAAt(x, y)
					c := ui.RGB(px.R, px.G, px.B)
					l := luminance(c)
					if !dark {
						l = 1 - l
					}
					if l > bestLum {
						best, bestLum = c, l
					}
				}
			}
			return best
		}
		check := func(what string, got, want ui.Color) {
			if adwColorDiff(got, want) > 2.0/255 {
				t.Errorf("dark %v: %s is %v, want libadwaita's %v", dark, what, got, want)
			}
		}
		card, window := find("Over the card"), find("Over the window")
		check("the entry's fill over the card", at(card.X+4, card.Y+card.H/2), pick(dark, ui.RGB(235, 235, 235), ui.RGB(72, 72, 75)))
		check("the placeholder over the card", ink(card.X+9, card.Y+8, card.X+card.W-9, card.Y+card.H-8), pick(dark, ui.RGB(132, 132, 135), ui.RGB(172, 172, 174)))
		check("the placeholder over the window", ink(window.X+9, window.Y+8, window.X+window.W-9, window.Y+window.H-8), pick(dark, ui.RGB(129, 129, 133), ui.RGB(166, 166, 167)))
		// "Send" is the label, inside the button's padding of 10.
		send, arrow := find("Send"), find(inputsMoreOptions)
		check("the split button's face over the card", at(send.X-4, send.Y+send.H/2), pick(dark, ui.RGB(235, 235, 235), ui.RGB(72, 72, 75)))
		check("the split button's seam over the card", ink(arrow.X-1, arrow.Y+arrow.H/2, arrow.X+2, arrow.Y+arrow.H/2+0.5), pick(dark, ui.RGB(217, 217, 217), ui.RGB(91, 91, 93)))
	}
}
