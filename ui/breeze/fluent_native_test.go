package breeze

import (
	"fmt"
	"image"
	"image/color"
	"testing"

	"github.com/egoist/mygo/ui"
)

// The values these tests check were measured from WinUI 3 3.3 (Windows App SDK 2.6.5, and CommunityToolkit.WinUI.Controls.SettingsControls 8.2.251219 for the SettingsCards) on Windows 11 26H2 by tools/winui-reference, at 2x, in the scenes winui-inputs, winui-status, winui-status-titles, winui-rows, winui-rows-wide and winui-settings. They read positions from Find and element bounds and colors from pixels away from text, so that they do not depend on the host's font.

// fluentNativeTester renders view on Windows at 2x, light or dark, after two frames.
func fluentNativeTester(w, h int, dark bool, view func(c *ui.Context)) *ui.Tester {
	tester := ui.NewTester(func(c *ui.Context) { Preview(c, Windows, func() { Use(c); view(c) }) }, w, h)
	tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tester.SetScale(2)
	tester.SetDark(dark)
	tester.Frame()
	tester.Frame()
	return tester
}

// fluentInkBox returns the box, in pixels at 2x, of the pixels of r (in DIPs) that hit says are ink, and whether there were any.
func fluentInkBox(img *image.RGBA, r ui.Rect, hit func(x, y int) bool) (image.Rectangle, bool) {
	box, found := image.Rectangle{}, false
	for y := int(2 * r.Y); y < int(2*(r.Y+r.H)); y++ {
		for x := int(2 * r.X); x < int(2*(r.X+r.W)); x++ {
			if !hit(x, y) {
				continue
			}
			px := image.Rect(x, y, x+1, y+1)
			if !found {
				box, found = px, true
			} else {
				box = box.Union(px)
			}
		}
	}
	return box, found
}

// gray reads the red channel of the pixel at x, y, which is the gray level of a gray pixel.
func gray(img *image.RGBA, x, y int) int { return int(img.RGBAAt(x, y).R) }

// TestFluentRevealGlyphAsMeasured checks the PasswordBox's RevealButton glyph, U+F78D, against WinUI's render: an eye without a lower lid, so that the column through its middle crosses three strokes (the lid, then the pupil's top and bottom), its ink 12 DIPs across and from 4 DIPs above the button's middle to 4 below it.
func TestFluentRevealGlyphAsMeasured(t *testing.T) {
	pw := ""
	tester := fluentNativeTester(400, 120, false, func(c *ui.Context) {
		ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() { PasswordField(c, &pw).Label("Password").Width(200) })
	})
	if err := tester.Click("Password"); err != nil {
		t.Fatal(err)
	}
	tester.Type("abc")
	tester.Frame()
	b, ok := tester.Find("Show Password")
	if !ok {
		t.Fatal("no reveal button")
	}
	img := tester.Image()
	ink := func(x, y int) bool { return gray(img, x, y) < 200 }
	box, ok := fluentInkBox(img, b, ink)
	cx, cy := 2*(b.X+b.W/2), 2*(b.Y+b.H/2)
	if !ok || box.Dx() < 23 || box.Dx() > 25 || float32(box.Min.Y) < cy-8.5 || float32(box.Max.Y) > cy+8.5 || float32(box.Max.Y) < cy+7 {
		t.Errorf("the glyph's ink is %v at 2x, want 24 pixels across, from 8 above %v to 8 below", box, cy)
	}
	if mid := float32(box.Min.X+box.Max.X) / 2; mid < cx-0.75 || mid > cx+0.75 {
		t.Errorf("the glyph's ink is centered at %v, want %v", mid, cx)
	}
	strokes, in := 0, false
	for y := box.Min.Y; y < box.Max.Y; y++ {
		if on := ink(int(cx), y); on && !in {
			strokes++
			in = true
		} else if !on {
			in = false
		}
	}
	if strokes != 3 {
		t.Errorf("the column through the glyph's middle crosses %d strokes, want 3: the lid and the pupil, no lower lid", strokes)
	}
}

// TestFluentPlaceholderAsMeasured checks an empty PasswordBox's placeholder against WinUI's render: TextFillColorSecondary over the box's fill, its ink reading 95 over 251 in the light theme and 207 over 45 in the dark one, where over the window it would read 92 and 204.
func TestFluentPlaceholderAsMeasured(t *testing.T) {
	for _, dark := range []bool{false, true} {
		pw := ""
		var f ui.Rect
		tester := fluentNativeTester(400, 120, dark, func(c *ui.Context) {
			ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() { f = PasswordField(c, &pw).Placeholder("Password").Width(180).Bounds() })
		})
		img := tester.Image()
		ink, fill := 255, gray(img, int(2*f.X)+6, int(2*f.Y)+6)
		if dark {
			ink = 0
		}
		for y := int(2 * f.Y); y < int(2*(f.Y+f.H)); y++ {
			for x := int(2 * f.X); x < int(2*(f.X+f.W)); x++ {
				if v := gray(img, x, y); (!dark && v < ink) || (dark && v > ink) {
					ink = v
				}
			}
		}
		want := map[bool][2]int{false: {95, 251}, true: {207, 45}}[dark]
		if ink != want[0] || fill != want[1] {
			t.Errorf("dark %v: the placeholder's ink reads %d over %d, want %d over %d", dark, ink, fill, want[0], want[1])
		}
	}
}

// TestFluentSplitButtonEdgesAsMeasured checks that a SplitButton's border reads over its parts' faces, as its border Grids lie over its background Grids: the top edge 236 and the bottom 211 over the face of 251 in the light theme, 65 and 60 over 45 in the dark one, where a Button's edges, over the window, read 229 and 204, 53 and 48.
func TestFluentSplitButtonEdgesAsMeasured(t *testing.T) {
	for _, dark := range []bool{false, true} {
		tester := fluentNativeTester(400, 120, dark, func(c *ui.Context) {
			ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() { SplitButton(c, "Save", inputsSaveMenu(new(string))) })
		})
		a, _ := tester.Find(inputsMoreOptions)
		img := tester.Image()
		x := int(2 * (a.X + a.W/2))
		top, bottom := gray(img, x, int(2*a.Y)), gray(img, x, int(2*(a.Y+a.H))-1)
		want := map[bool][2]int{false: {236, 211}, true: {65, 60}}[dark]
		if top < want[0]-1 || top > want[0]+1 || bottom < want[1]-1 || bottom > want[1]+1 {
			t.Errorf("dark %v: the arrow part's edges read %d at the top and %d at the bottom, want %v", dark, top, bottom, want)
		}
	}
}

// TestFluentCommandBarSplitButtonAsMeasured checks a SplitButton of SplitButtonCommandBarStyle in a Toolbar against WinUI's render: the label 13 DIPs into the main part, which ends 5 after it, and the chevron's ink centered 8 DIPs into the 23-DIP arrow part, its TextBlock at the start of the part's content.
func TestFluentCommandBarSplitButtonAsMeasured(t *testing.T) {
	var main ui.Rect
	tester := fluentNativeTester(400, 120, false, func(c *ui.Context) {
		ui.Column(c).AlignItems(ui.Start).Children(func() {
			Toolbar(c, func() { main = SplitButton(c, "New", inputsSaveMenu(new(string))).Bounds() })
		})
	})
	label, _ := tester.Find("New")
	a, _ := tester.Find(inputsMoreOptions)
	if d := label.X - main.X; d < 12.75 || d > 13.25 {
		t.Errorf("the label is %v into the main part %v, want 13", d, main)
	}
	if d := main.X + main.W - (label.X + label.W); d < 4.75 || d > 5.25 {
		t.Errorf("the main part %v ends %v after the label %v, want 5", main, d, label)
	}
	img := tester.Image()
	box, ok := fluentInkBox(img, a, func(x, y int) bool { return gray(img, x, y) < 150 })
	if mid := float32(box.Min.X+box.Max.X) / 4; !ok || mid < a.X+7.5 || mid > a.X+8.5 {
		t.Errorf("the chevron's ink %v is centered at %v in the arrow part %v, want 8 into it", box, mid, a)
	}
	// U+E972's ink from half a DIP above the part's middle to 2.5 below it, 18 pixels of it at 2x (WinUI's render: rows 267 to 272 around the middle at 268, the ink summing to 18.0 pixels), its stroke two thirds of a DIP.
	if mid := int(2 * (a.Y + a.H/2)); box.Min.Y != mid-1 || box.Max.Y != mid+5 {
		t.Errorf("the chevron's ink rows %d to %d at 2x, want %d to %d", box.Min.Y, box.Max.Y, mid-1, mid+5)
	}
	bg, ink := float64(gray(img, int(2*a.X)+1, int(2*a.Y)+1)), float64(26)
	area := 0.0
	for y := int(2 * a.Y); y < int(2*(a.Y+a.H)); y++ {
		for x := int(2 * a.X); x < int(2*(a.X+a.W)); x++ {
			area += (bg - float64(gray(img, x, y))) / (bg - ink)
		}
	}
	if area < 15 || area > 22 {
		t.Errorf("the chevron's ink sums to %.1f pixels at 2x, want 18", area)
	}
}

// TestFluentInfoBarAsMeasured checks an InfoBar against WinUI's render, light and dark: its border is CardStrokeColorDefault over the window, 229 and 29, around an informational and a success bar alike; the severity icon's circle is 15 DIPs across, flush with its 16-DIP box's start and a DIP below its top, and its marks' ink, from the circle's middle, from 3.5 above to 3 below for the i, 3 above to 3.5 below for the !, and 3.5 either side and 2.5 above and below for the check mark (half way from the circle's color to the mark's); the close button's × is 11 DIPs square, centered half a DIP above and before the button's middle.
func TestFluentInfoBarAsMeasured(t *testing.T) {
	for _, dark := range []bool{false, true} {
		open := true
		tester := fluentNativeTester(560, 200, dark, func(c *ui.Context) {
			ui.Column(c).Padding(16).Gap(16).Children(func() {
				InfoBar(c, "info", "Informational", "Some message about it.", &open)
				InfoBar(c, "success", "Success", "Some message about it.", nil)
				InfoBar(c, "warning", "Warning", "Some message about it.", nil)
			})
		})
		img := tester.Image()
		edge := map[bool]int{false: 229, true: 29}[dark]
		// The marks' ink boxes in DIPs from the circle's middle: start, top, end, bottom.
		marks := map[string][4]float32{"Informational": {-0.5, -3.5, 0.5, 3}, "Success": {-3.5, -2.5, 3.5, 2.5}, "Warning": {-0.5, -3, 0.5, 3.5}}
		for _, title := range []string{"Informational", "Success", "Warning"} {
			bar, _ := tester.Find(title)
			x := int(2 * (bar.X + bar.W/2))
			if got := img.RGBAAt(x, int(2*bar.Y)); absDiff(got.R, uint8(edge)) > 1 || absDiff(got.G, uint8(edge)) > 1 || absDiff(got.B, uint8(edge)) > 1 {
				t.Errorf("dark %v, %s: the top edge reads %v, want gray %d", dark, title, got, edge)
			}
			// The icon's box: 1 DIP of border and InfoBarContentRootPadding 16 at the start, InfoBarIconMargin 16 above.
			icon := ui.Rect{X: bar.X + 17 - 1, Y: bar.Y + 17 - 1, W: 18, H: 18}
			box, ok := fluentInkBox(img, icon, func(x, y int) bool {
				p := img.RGBAAt(x, y)
				return max(p.R, p.G, p.B)-min(p.R, p.G, p.B) > 40
			})
			// The warning bar's own background is colored too much for this.
			if want := image.Rect(int(2*(bar.X+17)), int(2*(bar.Y+18)), int(2*(bar.X+32)), int(2*(bar.Y+33))); title != "Warning" && (!ok || box != want) {
				t.Errorf("dark %v, %s: the icon's circle covers %v at 2x, want %v", dark, title, box, want)
			}
			// The mark: white in the light theme, black at 89.4% over the circle in the dark one.
			cx, cy := bar.X+17+7.5, bar.Y+18+7.5
			fill := img.RGBAAt(int(2*(bar.X+17))+4, int(2*cy))
			mark := func(v uint8) int { return int(pick(dark, 255, float32(v)*0.106)) }
			far := func(p color.RGBA) int { return absDiff(p.R, fill.R) + absDiff(p.G, fill.G) + absDiff(p.B, fill.B) }
			half := (absDiff(uint8(mark(fill.R)), fill.R) + absDiff(uint8(mark(fill.G)), fill.G) + absDiff(uint8(mark(fill.B)), fill.B)) / 2
			got, ok := fluentInkBox(img, ui.Rect{X: cx - 5, Y: cy - 5, W: 10, H: 10}, func(x, y int) bool { return far(img.RGBAAt(x, y)) > half })
			m := marks[title]
			if want := image.Rect(int(2*(cx+m[0])), int(2*(cy+m[1])), int(2*(cx+m[2])), int(2*(cy+m[3]))); !ok || got != want {
				t.Errorf("dark %v, %s: the icon's mark covers %v at 2x, want %v", dark, title, got, want)
			}
		}
		b, _ := tester.Find("Close")
		box, ok := fluentInkBox(img, b, func(x, y int) bool {
			return absDiff(img.RGBAAt(x, y).R, img.RGBAAt(int(2*b.X)+2, int(2*b.Y)+2).R) > 120
		})
		cx, cy := 2*(b.X+b.W/2)-1, 2*(b.Y+b.H/2)-1
		if !ok || box.Dx() != 22 || box.Dy() != 22 || float32(box.Min.X+box.Max.X)/2 != cx || float32(box.Min.Y+box.Max.Y)/2 != cy {
			t.Errorf("dark %v: the close glyph's ink is %v at 2x, want 22 pixels square centered on %v, %v", dark, box, cx, cy)
		}
	}
}

// TestSettingsCardWraps checks a SettingsCard's ContentAlignmentStates as WinUI renders them, for an action row, a switch row and the items of an AccordionItem: narrower than 476 DIPs the content goes below the header and the description, 8 below them, at the card's start, 17 from its edge (16 inside the border); wider, it stays at the end, 17 from the card's edge, centered on the card. The switch's On and Off share a column as wide as the wider, so its track does not move as it toggles. An AccordionItem's rows, a SettingsExpander's items, put their header 58 DIPs in and their content 44 before their end.
func TestSettingsCardWraps(t *testing.T) {
	for _, width := range []float32{300, 528} {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			name := fmt.Sprintf("width %v %v", width, dir)
			on, off := true, false
			tester := fluentNativeTester(600, 600, false, func(c *ui.Context) {
				ui.Column(c).Padding(16).Gap(16).Width(width + 32).Children(func() {
					ActionRow(c, "Wallpaper", "Choose a picture", func() { Button(c, "Open") })
					SwitchRow(c, "Dark mode", "Use a dark appearance", &on)
					SwitchRow(c, "Notifications", "", &off)
				})
			})
			tester.SetDirection(dir)
			tester.Frame()
			tester.Frame()
			row, _ := tester.Find("Wallpaper")
			sub, _ := tester.Find("Choose a picture")
			open, _ := tester.Find("Open")
			// The button's label, 12 DIPs into the 32-DIP button (ButtonPadding 11 inside the border).
			start := open.X - 12 - row.X
			end := row.X + row.W - (open.X + open.W + 12)
			if dir == ui.RightToLeft {
				start, end = end, start
			}
			if width < 476 {
				// The 32-DIP button ends 17 above the row's end, its top 8 below the description.
				top := row.Y + row.H - 17 - 32
				if start < 16.5 || start > 17.5 || top < sub.Y+sub.H+8-0.25 || top > sub.Y+sub.H+8+0.25 || open.Y < top || open.Y+open.H > top+32 {
					t.Errorf("%s: the button %v is not 8 below the description %v at the start of the row %v", name, open, sub, row)
				}
			} else if end < 16.5 || end > 17.5 || open.Y+open.H/2 < row.Y+row.H/2-1 || open.Y+open.H/2 > row.Y+row.H/2+1 {
				t.Errorf("%s: the button %v is not 17 from the end of the row %v and centered on it", name, open, row)
			}
			onText, _ := tester.Find("On")
			offText, _ := tester.Find("Off")
			if onText.W != offText.W {
				t.Errorf("%s: On is %v wide and Off %v: they do not share a column", name, onText.W, offText.W)
			}
		}
	}
	open, metered := true, false
	tester := fluentNativeTester(600, 400, false, func(c *ui.Context) {
		ui.Column(c).Padding(16).Width(560).Children(func() {
			Accordion(c, func() {
				AccordionItem(c, "Network", &open, func() {
					ActionRow(c, "Proxy", "Automatic", nil)
					SwitchRow(c, "Metered", "", &metered)
				})
			})
		})
	})
	row, _ := tester.Find("Proxy")
	sub, _ := tester.Find("Automatic")
	if d := sub.X - row.X; d < 57.5 || d > 58.5 {
		t.Errorf("an AccordionItem's row puts its text %v in, want 58", d)
	}
	sw, _ := tester.Find("Metered")
	offText, _ := tester.Find("Off")
	// The track, 40 wide, 12 after the Off column, ends 44 before the row's end.
	if d := sw.X + sw.W - (offText.X + offText.W + 12 + 40); d < 43.5 || d > 44.5 {
		t.Errorf("an AccordionItem's switch row ends its switch %v before its end, want 44", d)
	}
}

// TestFluentSettingsRowsOverTheCard checks what a card of settings rows holds against WinUI's render, light and dark: SettingsCards on a card lay their translucent colors over the card's fill, not over the window. The line between two rows is CardStrokeColorDefault over the card, 236 and 39; a TextBox reads 254 and 55 inside a border of 236 and 58 with a bottom edge of 139 and 159, a ComboBox 254 and 55, and the description TextFillColorSecondary over the card, 95 and 207 (where over the window they read 229 and 29, 251 and 45, 229 and 48, 134 and 154, 92 and 204). A disabled TextBox's or NumberBox's bottom edge is its border's color, without the strong line (measured from WinUI 3 3.3 with CommunityToolkit's SettingsControls 8.2.251219 by tools/winui-reference's winui-settings). The controls are found from the rows' bounds: in a card 300 wide they sit 8 above the row's end at its start, 16 in, and alone 17 inside the card's edge.
func TestFluentSettingsRowsOverTheCard(t *testing.T) {
	for _, dark := range []bool{false, true} {
		appearance, name, email, size, copies := "Auto", "Evan", "", 12.0, 3.0
		var combo, entry, spin, disabled, disabledEntry ui.Rect
		tester := fluentNativeTester(400, 600, dark, func(c *ui.Context) {
			ui.Column(c).Padding(16).Width(332).Gap(16).Children(func() {
				Card(c, func() {
					combo = ComboRow(c, "Appearance", "Light or dark", &appearance, []string{"Light", "Dark", "Auto"}).Bounds()
					entry = EntryRow(c, "Name", &name).Bounds()
					spin = SpinRow(c, "Font size", "Points", &size, 6, 72, 1).Bounds()
				})
				disabled = SpinRow(c, "Copies", "", &copies, 1, 9, 1).Disabled(true).Bounds()
				disabledEntry = EntryRow(c, "Email", &email).Disabled(true).Bounds()
			})
		})
		img := tester.Image()
		at := func(x, y float32) int { return gray(img, int(2*x), int(2*y)) }
		pickv := func(l, d int) int {
			if dark {
				return d
			}
			return l
		}
		check := func(what string, got, wantv int) {
			t.Helper()
			if got != wantv {
				t.Errorf("dark %v: %s reads %d, want %d", dark, what, got, wantv)
			}
		}
		check("the line above the entry row", at(entry.X+150, entry.Y), pickv(236, 39))
		check("the line above the spin row", at(spin.X+150, spin.Y), pickv(236, 39))
		// The TextBox: 120 wide, its top 40 above the row's end.
		fx, fy := entry.X+16, entry.Y+entry.H-8-32
		check("the TextBox's fill", at(fx+100, fy+4), pickv(254, 55))
		check("the TextBox's top edge", gray(img, int(2*(fx+100)), int(2*fy+0.5)), pickv(236, 58))
		check("the TextBox's bottom edge", gray(img, int(2*(fx+100)), int(2*(fy+32)+0.5)-1), pickv(139, 159))
		check("the ComboBox's face", at(combo.X+16+60, combo.Y+combo.H-8-32+4), pickv(254, 55))
		check("the NumberBox's fill", at(spin.X+16+60, spin.Y+spin.H-8-32+4), pickv(254, 55))
		// "Points", below "Font size", which has no descenders to reach into its box.
		sub, _ := tester.Find("Points")
		ink := pickv(255, 0)
		for y := int(2 * sub.Y); y < int(2*(sub.Y+sub.H)); y++ {
			for x := int(2 * sub.X); x < int(2*(sub.X+sub.W)); x++ {
				if v := gray(img, x, y); (!dark && v < ink) || (dark && v > ink) {
					ink = v
				}
			}
		}
		check("the description's ink", ink, pickv(95, 207))
		// The disabled NumberBox and TextBox alone: 17 inside the card, their bottom 17 above the card's end. The rows' heights follow the font's line, so the box's edges are found at the nearest pixel.
		near := func(v float32) int { return int(2*v + 0.5) }
		for _, d := range []struct {
			name string
			row  ui.Rect
		}{{"NumberBox", disabled}, {"TextBox", disabledEntry}} {
			bx, by := d.row.X+17+60, d.row.Y+d.row.H-17-32
			if top, bottom := gray(img, int(2*bx), near(by)), gray(img, int(2*bx), near(by+32)-1); top != bottom {
				t.Errorf("dark %v: the disabled %s's edges read %d at the top and %d at the bottom, want the same", dark, d.name, top, bottom)
			}
		}
	}
}
