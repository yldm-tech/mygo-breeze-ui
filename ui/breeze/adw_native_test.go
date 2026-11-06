package breeze

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// The values these tests pin were measured from libadwaita 1.10.0's rendering (GTK 4.24.1, Adwaita Sans) by tools/adwaita-reference, its inputs, rows, overcard, settings, tabview, headerbar, forms and spinlimit sets, against TestScreenshotsInputs, TestScreenshotsRows, TestAdwControlsOverCard, TestScreenshotsSettings, TestScreenshotsTabView, TestScreenshotsHeaderBar, TestScreenshotsWave2Forms and TestAdwPlusEvenWhenDisabled.

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

// adwInkTop returns the first row, in device pixels of a frame at scale, of r (in DIPs) that holds ink darker than 200, or -1.
func adwInkTop(tester *ui.Tester, r ui.Rect, scale float32) int {
	img := tester.Image()
	for y := int(r.Y * scale); y < int((r.Y+r.H)*scale); y++ {
		for x := int(r.X * scale); x < int((r.X+r.W)*scale); x++ {
			if img.RGBAAt(x, y).R < 200 {
				return y
			}
		}
	}
	return -1
}

// TestAdwEntryRowTitleOverText checks on GNOME that a filled entry row draws its title in the box of the small text it is laid out as, its ink as far below the box's top as a row's subtitle's below its own: AdwEntryRow allocates the title label at the subtitle's size above the text (allocate_editable_area), and libadwaita 1.10.0 draws "Name" as high in its row as that (tools/adwaita-reference's settings set), where the title drawn as a small span on a line of the body's size sat about 2 pixels lower. It reads the ink from the frame against the elements' boxes, whatever the host's font.
func TestAdwEntryRowTitleOverText(t *testing.T) {
	const scale = 2
	shoot := func(build func(c *ui.Context)) *ui.Tester {
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Linux, func() {
				Use(c)
				ui.Column(c).Padding(16).Width(332).Children(func() { build(c) })
			})
		}, 400, 120)
		tester.SetScale(scale)
		tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		tester.Frame()
		tester.Frame()
		return tester
	}
	name := "Evan"
	entry := shoot(func(c *ui.Context) { EntryRow(c, "Name", &name) })
	row := shoot(func(c *ui.Context) { ActionRow(c, "Title", "Name", nil) })
	// The first element showing "Name" in the entry row is its small text, laid out above the field and drawn over.
	small, ok := entry.Find("Name")
	sub, ok2 := row.Find("Name")
	if !ok || !ok2 {
		t.Fatal("the title or the subtitle is missing")
	}
	got, want := adwInkTop(entry, small, scale)-int(small.Y*scale), adwInkTop(row, sub, scale)-int(sub.Y*scale)
	if got < want-1 || got > want+1 {
		t.Fatalf("the title's ink starts %d device pixels below its box's top, the subtitle's %d", got, want)
	}
}

// TestAdwSpinRowIconsOverTheirFace checks on GNOME, light and dark, that a spin row's − and + are drawn in the button's color over the button's own face, as libadwaita 1.10.0 renders them (tools/adwaita-reference's settings set): 47 over the light face of 235, where over the card they would read 51, and white over the dark face.
func TestAdwSpinRowIconsOverTheirFace(t *testing.T) {
	for _, dark := range []bool{false, true} {
		size := 12.0
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Linux, func() {
				Use(c)
				ui.Column(c).Padding(16).Width(332).Children(func() { SpinRow(c, "Font size", "Points", &size, 6, 72, 1) })
			})
		}, 400, 120)
		tester.SetScale(2)
		tester.SetDark(dark)
		tester.Frame()
		for _, name := range []string{"Decrease", "Increase"} {
			b, ok := tester.Find(name)
			if !ok {
				t.Fatalf("dark %v: %s is missing", dark, name)
			}
			// The middle of the bar across both icons, a whole device pixel inside its 2-pixel stroke.
			px := tester.Image().RGBAAt(int((b.X+b.W/2)*2), int((b.Y+b.H/2)*2))
			want := map[bool]uint8{false: 47, true: 255}[dark]
			if d := int(px.R) - int(want); d < -2 || d > 2 {
				t.Errorf("dark %v: %s's icon reads %d, want libadwaita's %d", dark, name, px.R, want)
			}
		}
	}
}

// TestAdwTabBarAsRendered checks on GNOME, light and dark, the tab bars of TestScreenshotsTabView against libadwaita 1.10.0's AdwTabBars in tools/adwaita-reference's tabview set: each bar 47 pixels high with headerbar_shade_color in its last row (224 over 255 in the light theme, 29 over 46 in the dark one) and the window below it (250, 34); the tabs share the bar in whole pixels as AdwTabBox does, the first a pixel narrower, at 22, 182 and 343, 155, 156 and 156 wide beside the add button, and at 22 and 282, 255 and 256 wide without it; the separator before Pictures fills the pixel column at 340 (224, 77) from 25 to 52.5, 3 inside the tabs' 34 pixels.
func TestAdwTabBarAsRendered(t *testing.T) {
	for _, dark := range []bool{false, true} {
		tester := ui.NewTester(tabViewView(Linux), 560, 260)
		tester.SetScale(2)
		tester.SetDark(dark)
		tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		tester.Frame()
		tester.Frame()
		want := map[string][2]float32{"Documents": {22, 155}, "Downloads": {182, 156}, "Pictures": {343, 156}, "General": {22, 255}, "Advanced": {282, 256}}
		for name, xw := range want {
			r, ok := tester.Find(name)
			if !ok {
				t.Fatalf("dark %v: the tab %s is missing", dark, name)
			}
			if r.X != xw[0] || r.W != xw[1] {
				t.Errorf("dark %v: the tab %s is %v wide at %v, want %v wide at %v", dark, name, r.W, r.X, xw[1], xw[0])
			}
		}
		img := tester.Image()
		at := func(x, y int) uint8 { return img.RGBAAt(x, y).R }
		bg, shade, window := map[bool]uint8{false: 255, true: 46}[dark], map[bool]uint8{false: 224, true: 29}[dark], map[bool]uint8{false: 250, true: 34}[dark]
		near := func(got, want uint8) bool { return int(got)-int(want) <= 2 && int(want)-int(got) <= 2 }
		// The bars start at 16 and 136; at x 100, 47 rows down, the window shows again.
		for _, top := range []int{16, 136} {
			for _, row := range []struct {
				y    int
				want uint8
			}{{2 * (top + 45), bg}, {2*(top+46) + 1, shade}, {2 * (top + 46), shade}, {2 * (top + 47), window}} {
				if got := at(200, row.y); !near(got, row.want) {
					t.Errorf("dark %v: the bar at %d reads %d at device row %d, want %d", dark, top, got, row.y, row.want)
				}
			}
		}
		sep := map[bool]uint8{false: 224, true: 77}[dark]
		if got := at(680, 2*38); !near(got, sep) || !near(at(681, 2*38), sep) || !near(at(682, 2*38), bg) {
			t.Errorf("dark %v: the separator before Pictures reads %d, %d, %d across device columns 680 to 682, want %d, %d, %d", dark, got, at(681, 2*38), at(682, 2*38), sep, sep, bg)
		}
	}
}

// adwInkSpan returns the first and the last device column of r (in DIPs) holding ink darker than 160 in a frame at scale, or -1, -1.
func adwInkSpan(tester *ui.Tester, r ui.Rect, scale float32) (int, int) {
	img := tester.Image()
	first, last := -1, -1
	for x := int(r.X * scale); x < int((r.X+r.W)*scale); x++ {
		for y := int(r.Y * scale); y < int((r.Y+r.H)*scale); y++ {
			if img.RGBAAt(x, y).R < 160 {
				if first < 0 {
					first = x
				}
				last = x
				break
			}
		}
	}
	return first, last
}

// TestAdwSpinDigitsTabular checks on GNOME that the spin row's value and NumberInput's are drawn in tabular digits, as GtkSpinButton is .numeric (font-variant-numeric: tabular-nums in libadwaita 1.10.0's _labels.scss): "11" spans as many device pixels as a label of it with the tnum feature, where libadwaita 1.10.0 draws the "1" of "12" 3 pixels further from the "2" than Adwaita Sans's proportional digits put it (tools/adwaita-reference's settings set). It compares the ink with the label's in the same font, so a font whose digits are tabular anyway passes as well.
func TestAdwSpinDigitsTabular(t *testing.T) {
	const scale = 2
	for _, number := range []bool{false, true} {
		v := 11.0
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Linux, func() {
				Use(c)
				ui.Column(c).Padding(16).Gap(16).Width(332).AlignItems(ui.Start).Children(func() {
					if number {
						NumberInput(c, &v, 0, 99, 1)
					} else {
						ui.Column(c).Width(300).Children(func() { SpinRow(c, "Copies", "", &v, 0, 99, 1) })
					}
					ui.Text(c, "11").FontFeatures("tnum").Label("tabular")
				})
			})
		}, 400, 200)
		tester.SetScale(scale)
		tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		tester.Frame()
		tester.Frame()
		dec, ok := tester.Find("Decrease")
		label, ok2 := tester.Find("tabular")
		if !ok || !ok2 {
			t.Fatal("the buttons or the label are missing")
		}
		// The text before the buttons: NumberInput's 80 wide right before Decrease, the spin row's 64 ending 2 (the button's margin) and 6 (the box's gap) before it.
		text := ui.Rect{X: dec.X - 80, Y: dec.Y, W: 80, H: dec.H}
		if !number {
			text = ui.Rect{X: dec.X - 8 - 64, Y: dec.Y, W: 64, H: dec.H}
		}
		a, b := adwInkSpan(tester, text, scale)
		la, lb := adwInkSpan(tester, label, scale)
		if a < 0 || la < 0 {
			t.Fatalf("number input %v: no ink in the text (%d) or the label (%d)", number, a, la)
		}
		if got, want := b-a, lb-la; got < want-1 || got > want+1 {
			t.Errorf("number input %v: \"11\" spans %d device pixels, a tabular label of it %d", number, got, want)
		}
	}
}

// TestAdwPlusEvenWhenDisabled checks on GNOME, light and dark, that the + of a disabled Increase is as light in its middle as along its arms, the spin row's at its upper limit and NumberInput's: libadwaita 1.10.0 draws the icon of a disabled button even (151 over the face of 245, 153 over 61 dark, in tools/adwaita-reference's spinlimit set; 175 over 230 in its forms set), where MyGo's half opacity, applied fill by fill, would darken the crossing of two bars drawn over each other.
func TestAdwPlusEvenWhenDisabled(t *testing.T) {
	for _, number := range []bool{false, true} {
		for _, dark := range []bool{false, true} {
			v := 9.0
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, Linux, func() {
					Use(c)
					ui.Column(c).Padding(16).Width(332).AlignItems(ui.Start).Children(func() {
						if number {
							NumberInput(c, &v, 1, 9, 1)
						} else {
							ui.Column(c).Width(300).Children(func() { SpinRow(c, "Copies", "", &v, 1, 9, 1) })
						}
					})
				})
			}, 400, 120)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tester.Frame()
			b, ok := tester.Find("Increase")
			if !ok {
				t.Fatalf("number input %v, dark %v: Increase is missing", number, dark)
			}
			// The middle of the icon, and a device pixel inside its bar across, 4 pixels to the left.
			img := tester.Image()
			cx, cy := int((b.X+b.W/2)*2), int((b.Y+b.H/2)*2)
			if !number {
				cx = int((b.X+b.W/2)*2) - 1
			}
			mid, arm := img.RGBAAt(cx, cy).R, img.RGBAAt(cx-8, cy).R
			if d := int(mid) - int(arm); d < -3 || d > 3 {
				t.Errorf("number input %v, dark %v: the disabled +'s middle reads %d, its arm %d", number, dark, mid, arm)
			}
		}
	}
}
