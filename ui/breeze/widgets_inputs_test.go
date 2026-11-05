package breeze

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func init() {
	for _, n := range []string{"PasswordField", "SplitButton", "IconButton"} {
		drawnWidgets[n] = true
	}
}

// inputsTestTrashIcon is a trash can for the icon buttons of the tests.
var inputsTestTrashIcon = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><path d="M2.5 4h11M6 4V2.5h4V4M4 4l.8 10h6.4L12 4M6.5 6.5v5M9.5 6.5v5" fill="none" stroke="#000" stroke-width="1.3" stroke-linejoin="round"/></svg>`))

// inputsSaveMenu builds the menu of the split buttons of the tests, setting *chosen to the item chosen.
func inputsSaveMenu(chosen *string) func(m *ui.Menu) {
	return func(m *ui.Menu) {
		for _, it := range []string{"Save As…", "Save All"} {
			if m.Item(it).Chosen() {
				*chosen = it
			}
		}
	}
}

// TestPasswordFieldHidesItsText checks on every platform that the field edits its value, submits on Enter, is named by Label, and hides its text: Copy takes nothing from it.
func TestPasswordFieldHidesItsText(t *testing.T) {
	for _, p := range Platforms {
		pw, submitted := "", 0
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Column(c).Padding(16).Gap(12).AlignItems(ui.Start).Children(func() {
					if PasswordField(c, &pw).Label("Password").Width(200).Submitted() {
						submitted++
					}
					Button(c, "Next")
				})
			})
		}, 400, 200)
		tester.Frame()
		if err := tester.Click("Password"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Type("hunter2")
		tester.Frame()
		if pw != "hunter2" || !tester.Focused("Password") {
			t.Fatalf("%v: typing gave %q, focused %v", p, pw, tester.Focused("Password"))
		}
		tester.Key(ui.Cmd, ui.KeyA)
		tester.SetClipboard("")
		tester.Key(ui.Cmd, ui.KeyC)
		if got := tester.Clipboard(); got != "" {
			t.Errorf("%v: Copy took %q from the hidden text", p, got)
		}
		tester.Key(0, ui.KeyEnter)
		tester.Frame()
		if submitted != 1 {
			t.Errorf("%v: Enter submitted %d times", p, submitted)
		}
		if f, ok := tester.Find("Password"); !ok || f.W != 200 {
			t.Errorf("%v: the field is %v, %v", p, f, ok)
		}
	}
}

// TestPasswordBoxRevealButton checks Windows's PasswordBox, left to right and right to left: the RevealButton shows at the field's end once text is typed into it while it has the focus, not for text it held when it got the focus, nor in a field 5 em wide or less; it is not a Tab stop; the text shows while it is pressed and hides once it is released.
func TestPasswordBoxRevealButton(t *testing.T) {
	for _, rtl := range []bool{false, true} {
		pw, short := "", ""
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Windows, func() {
				Use(c)
				ui.Column(c).Padding(16).Gap(12).AlignItems(ui.Start).Children(func() {
					PasswordField(c, &pw).Label("Password").Width(200)
					Button(c, "Next")
					PasswordField(c, &short).Label("PIN").Width(70)
				})
			})
		}, 400, 240)
		if rtl {
			tester.SetDirection(ui.RightToLeft)
		}
		tester.Frame()
		if err := tester.Click("Password"); err != nil {
			t.Fatal(err)
		}
		if tester.HasText("Show Password") {
			t.Fatalf("rtl %v: the reveal button shows in an empty field", rtl)
		}
		tester.Type("abc")
		tester.Frame()
		f, _ := tester.Find("Password")
		b, ok := tester.Find("Show Password")
		if !ok {
			t.Fatalf("rtl %v: no reveal button once text was typed", rtl)
		}
		// 26 by 24, 4 DIPs in from the field's end and centered on it.
		inset := f.X + f.W - (b.X + b.W)
		if rtl {
			inset = b.X - f.X
		}
		if b.W != 26 || b.H != 24 || inset < 3.5 || inset > 4.5 || b.Y+b.H/2-(f.Y+f.H/2) > 0.5 || f.Y+f.H/2-(b.Y+b.H/2) > 0.5 {
			t.Errorf("rtl %v: the reveal button is %v in the field %v", rtl, b, f)
		}
		tester.Key(ui.Cmd, ui.KeyA)
		tester.SetClipboard("")
		tester.Press(b.X+b.W/2, b.Y+b.H/2)
		if !tester.Focused("Password") {
			t.Errorf("rtl %v: pressing the reveal button took the focus from the field", rtl)
		}
		tester.Key(ui.Cmd, ui.KeyC)
		if got := tester.Clipboard(); got != "abc" {
			t.Errorf("rtl %v: the text did not show while the button was pressed: Copy took %q", rtl, got)
		}
		tester.Release(b.X+b.W/2, b.Y+b.H/2)
		tester.SetClipboard("")
		tester.Key(ui.Cmd, ui.KeyC)
		if got := tester.Clipboard(); got != "" {
			t.Errorf("rtl %v: the text stayed shown after the release: Copy took %q", rtl, got)
		}
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		if !tester.Focused("Next") || tester.HasText("Show Password") {
			t.Errorf("rtl %v: Tab went to the reveal button or left it showing", rtl)
		}
		tester.Key(ui.Shift, ui.KeyTab)
		tester.Frame()
		if !tester.Focused("Password") || tester.HasText("Show Password") {
			t.Errorf("rtl %v: the reveal button shows for text the field held as it got the focus", rtl)
		}
		tester.Key(ui.Cmd, ui.KeyA)
		tester.Key(0, ui.KeyBackspace)
		tester.Type("x")
		tester.Frame()
		if pw != "x" || !tester.HasText("Show Password") {
			t.Errorf("rtl %v: emptying the field and typing again did not show the button (%q)", rtl, pw)
		}
		if !rtl {
			// A password longer than the field scrolls, its end stopping 6 DIPs before the button's 30-DIP column (TextControlThemePadding's end inside the border).
			tester.Type(strings.Repeat("x", 40))
			tester.Frame()
			b, _ = tester.Find("Show Password")
			if caret, ok := tester.TextCaret(); !ok || caret.X+caret.W > b.X+b.W+4-30-6+0.5 {
				t.Errorf("the caret at the end of a long password is at %v, past the reveal button's column at %v", caret, b.X+b.W+4-30)
			}
		}
		if err := tester.Click("PIN"); err != nil {
			t.Fatal(err)
		}
		tester.Type("1234")
		tester.Frame()
		if short != "1234" || tester.HasText("Show Password") {
			t.Errorf("rtl %v: a field 5 em wide shows the reveal button", rtl)
		}
	}
}

// TestPasswordBoxAltF8 checks the keys of Windows's PasswordBox: Alt+F8 reveals the text until F8 is released, and F8 alone, or other keys, do nothing and go on to the field.
func TestPasswordBoxAltF8(t *testing.T) {
	st := &fluentPeek{}
	keys := fluentPeekKeys(st)
	if keys(ui.InputEvent{Kind: ui.InputKeyDown, Key: ui.KeyF8}) || st.key {
		t.Fatal("F8 without Alt revealed the text")
	}
	if !keys(ui.InputEvent{Kind: ui.InputKeyDown, Key: ui.KeyF8, Mods: ui.Alt}) || !st.key {
		t.Fatal("Alt+F8 did not reveal the text")
	}
	if keys(ui.InputEvent{Kind: ui.InputKeyDown, Key: ui.KeyA, Mods: ui.Alt}) || keys(ui.InputEvent{Kind: ui.InputText, Text: "a"}) {
		t.Fatal("the handler took a key that is not its own")
	}
	if !keys(ui.InputEvent{Kind: ui.InputKeyUp, Key: ui.KeyF8}) || st.key {
		t.Fatal("releasing F8 left the text shown")
	}
	if keys(ui.InputEvent{Kind: ui.InputKeyUp, Key: ui.KeyF8}) {
		t.Fatal("a second release of F8 was taken")
	}
}

// TestPasswordEntryShowText checks GNOME's password entry: its context menu ends with Show Text, which shows the text and, chosen again, hides it; macOS has no such item.
func TestPasswordEntryShowText(t *testing.T) {
	for _, p := range []Platform{Linux, Mac} {
		pw := ""
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Column(c).Padding(16).Children(func() { PasswordField(c, &pw).Label("Password").Width(200) })
			})
		}, 400, 200)
		tester.Frame()
		if err := tester.Click("Password"); err != nil {
			t.Fatal(err)
		}
		tester.Type("hunter2")
		if err := tester.RightClick("Password"); err != nil {
			t.Fatal(err)
		}
		menu := tester.Menu()
		if p == Mac {
			if strings.Contains(strings.Join(menu, ","), "Show Text") {
				t.Errorf("macOS: the context menu has Show Text: %q", menu)
			}
			continue
		}
		if len(menu) < 3 || menu[len(menu)-1] != "Show Text" || menu[len(menu)-2] != "-" {
			t.Fatalf("GNOME: the context menu is %q, want the editing items, a separator and Show Text", menu)
		}
		copied := func() string {
			tester.Key(ui.Cmd, ui.KeyA)
			tester.SetClipboard("")
			tester.Key(ui.Cmd, ui.KeyC)
			return tester.Clipboard()
		}
		if err := tester.ChooseMenuItem("Show Text"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		if got := copied(); got != "hunter2" {
			t.Errorf("GNOME: Show Text left the text hidden: Copy took %q", got)
		}
		if err := tester.RightClick("Password"); err != nil {
			t.Fatal(err)
		}
		if err := tester.ChooseMenuItem("Show Text"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		if got := copied(); got != "" {
			t.Errorf("GNOME: Show Text chosen again left the text shown: Copy took %q", got)
		}
	}
}

// TestSplitButtonParts checks the split button on every platform, left to right and right to left, alone and in a Toolbar: the other part is after the main one in the reading direction and named More Options; Tab reaches the main part, where Enter clicks it, then Tab (the arrows in a toolbar) the other, where Down opens the menu build fills, whose item chosen reports Chosen; a click on the main part clicks it, and a click on the other opens the menu.
func TestSplitButtonParts(t *testing.T) {
	for _, p := range Platforms {
		for _, bar := range []bool{false, true} {
			for _, rtl := range []bool{false, true} {
				name := fmt.Sprintf("%v toolbar %v rtl %v", p, bar, rtl)
				clicks, chosen := 0, ""
				tester := ui.NewTester(func(c *ui.Context) {
					Preview(c, p, func() {
						Use(c)
						ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() {
							build := func() {
								if SplitButton(c, "Save", inputsSaveMenu(&chosen)).Clicked() {
									clicks++
								}
							}
							if bar {
								Toolbar(c, build).Label("Edit")
							} else {
								build()
							}
						})
					})
				}, 400, 200)
				if rtl {
					tester.SetDirection(ui.RightToLeft)
				}
				tester.Frame()
				tester.Frame()
				s, ok1 := tester.Find("Save")
				m, ok2 := tester.Find(inputsMoreOptions)
				if !ok1 || !ok2 {
					t.Fatalf("%s: parts missing: %v %v", name, ok1, ok2)
				}
				after := m.X >= s.X+s.W-1.5
				if rtl {
					after = m.X+m.W <= s.X+1.5
				}
				if !after || m.Y > s.Y+s.H/2 || m.Y+m.H < s.Y+s.H/2 {
					t.Errorf("%s: the arrow part %v is not after the label %v", name, m, s)
				}
				// Tab from nothing focused goes to the main part (in a toolbar, the toolbar's stop, its first control).
				tester.Key(0, ui.KeyTab)
				tester.Frame()
				if !tester.Focused("Save") {
					t.Fatalf("%s: Tab did not focus the main part", name)
				}
				tester.Key(0, ui.KeyEnter)
				tester.Frame()
				if clicks != 1 || tester.Menu() != nil {
					t.Fatalf("%s: Enter on the main part clicked %d times, menu %q", name, clicks, tester.Menu())
				}
				next := ui.KeyTab
				if bar {
					next = map[bool]ui.Key{false: ui.KeyRight, true: ui.KeyLeft}[rtl]
				}
				tester.Key(0, next)
				tester.Frame()
				if !tester.Focused(inputsMoreOptions) {
					t.Fatalf("%s: the keys did not reach the arrow part", name)
				}
				tester.Key(0, ui.KeyDown)
				if got := strings.Join(tester.Menu(), ","); got != "Save As…,Save All" {
					t.Fatalf("%s: Down on the arrow part opened %q", name, got)
				}
				if err := tester.ChooseMenuItem("Save All"); err != nil {
					t.Fatal(err)
				}
				tester.Frame()
				if chosen != "Save All" || clicks != 1 {
					t.Fatalf("%s: the item chosen is %q, clicks %d", name, chosen, clicks)
				}
				if err := tester.Click("Save"); err != nil {
					t.Fatal(err)
				}
				tester.Frame()
				if clicks != 2 || tester.Menu() != nil {
					t.Fatalf("%s: a click on the main part clicked %d times in all, menu %q", name, clicks, tester.Menu())
				}
				if err := tester.Click(inputsMoreOptions); err != nil {
					t.Fatal(err)
				}
				if got := strings.Join(tester.Menu(), ","); got != "Save As…,Save All" || clicks != 2 {
					t.Fatalf("%s: a click on the arrow opened %q, clicks %d", name, got, clicks)
				}
				tester.CloseMenu()
			}
		}
	}
}

// TestSplitButtonSizes checks the parts' sizes on every platform, alone and in a Toolbar: macOS's 24-point NSComboButton with its 17-point arrow part (the divider and 16 points) and, in a toolbar, the 25-point pill and the 40-point arrow part measured from AppKit's render; WinUI's 32-DIP SplitButton with its 35-DIP secondary part and the 36-DIP SplitButtonCommandBarStyle with its 23-DIP one; libadwaita's 34-pixel buttons with the 24-pixel arrow, raised and flat.
func TestSplitButtonSizes(t *testing.T) {
	type size struct{ h, arrowW float32 }
	want := map[Platform][2]size{
		Mac:     {{24, 17}, {25, 40}},
		Windows: {{32, 35}, {36, 23}},
		Linux:   {{34, 24}, {34, 24}},
	}
	for _, p := range Platforms {
		for i, bar := range []bool{false, true} {
			name := fmt.Sprintf("%v toolbar %v", p, bar)
			var main ui.Rect
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() {
						build := func() { main = SplitButton(c, "Save", inputsSaveMenu(new(string))).Bounds() }
						if bar {
							Toolbar(c, build)
						} else {
							build()
						}
					})
				})
			}, 400, 200)
			tester.Frame()
			tester.Frame()
			a, ok := tester.Find(inputsMoreOptions)
			w := want[p][i]
			if !ok || main.H != w.h || a.H != w.h || a.W != w.arrowW {
				t.Errorf("%s: the main part is %v and the arrow part %v, want both %v high and the arrow %v wide", name, main, a, w.h, w.arrowW)
			}
		}
	}
}

// TestMacComboButtonPixels checks macOS's NSComboButton against AppKit's render (appkit-inputs-light.png of tools/appkit-reference): the chevron's ink is centered on the divider and the arrow part together, as AppKit's spans the pixels 148 to 165 of 140 to 173 at 2x; disabled, the face fades to 245 as a disabled push button's does.
func TestMacComboButtonPixels(t *testing.T) {
	for _, off := range []bool{false, true} {
		var main ui.Rect
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Mac, func() {
				Use(c)
				ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() {
					main = SplitButton(c, "Save", inputsSaveMenu(new(string))).Disabled(off).Bounds()
				})
			})
		}, 400, 200)
		tester.SetScale(2)
		tester.Frame()
		tester.Frame()
		tester.Frame()
		img := tester.Image()
		a, _ := tester.Find(inputsMoreOptions)
		y := int(2 * (a.Y + a.H/2))
		face := img.RGBAAt(int(2*main.X)+6, y)
		if wantFace := map[bool]uint8{false: 235, true: 245}[off]; absDiff(face.R, wantFace) > 1 {
			t.Errorf("disabled %v: the face reads %v, want %d", off, face, wantFace)
		}
		if off {
			continue
		}
		// The chevron's ink across the arrow part, past the divider's two pixels.
		first, last := -1, -1
		for x := int(2*a.X) + 2; x < int(2*(a.X+a.W)); x++ {
			for dy := -6; dy <= 6; dy++ {
				if img.RGBAAt(x, y+dy).R < face.R-40 {
					if first < 0 {
						first = x
					}
					last = x
				}
			}
		}
		mid := float32(first+last+1) / 2
		if center := 2 * (a.X + a.W/2); first < 0 || mid < center-0.5 || mid > center+0.5 {
			t.Errorf("the chevron's ink spans the pixels %d to %d, centered on %v, want %v", first, last, mid, center)
		}
	}
}

// TestFluentSplitButtonChevron checks the ink of WinUI's SplitButton's chevron, light and dark: TextFillColorSecondary over the face at rest, and TextFillColorPrimary with the pointer over the secondary part (SecondaryPointerOver), read at the chevron's darkest pixel in the light theme and its brightest in the dark one.
func TestFluentSplitButtonChevron(t *testing.T) {
	for _, dark := range []bool{false, true} {
		var th *ui.Theme
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Windows, func() {
				Use(c)
				th = c.Theme()
				ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() { SplitButton(c, "Save", inputsSaveMenu(new(string))) })
			})
		}, 400, 200)
		tester.SetScale(2)
		tester.SetDark(dark)
		tester.Frame()
		tester.Frame()
		a, _ := tester.Find(inputsMoreOptions)
		ink := func() ui.Color {
			img := tester.Image()
			best := img.RGBAAt(int(2*a.X)+4, int(2*(a.Y+a.H/2)))
			for y := int(2*a.Y) + 4; y < int(2*(a.Y+a.H))-4; y++ {
				for x := int(2*a.X) + 4; x < int(2*(a.X+a.W))-4; x++ {
					if px := img.RGBAAt(x, y); (px.R < best.R) != dark && px.R != best.R {
						best = px
					}
				}
			}
			return ui.RGB(best.R, best.G, best.B)
		}
		check := func(state string, want ui.Color) {
			if got := ink(); max(absDiff(got.R, want.R), absDiff(got.G, want.G), absDiff(got.B, want.B)) > 3 {
				t.Errorf("dark %v, %s: the chevron reads %v, want %v", dark, state, got, want)
			}
		}
		check("at rest", fluentSecondaryOver(th, fluentFace(th, false, false)))
		tester.Move(a.X+a.W/2, a.Y+a.H/2)
		tester.Frame()
		check("under the pointer", th.Text)
	}
}

// TestSplitButtonDisabled checks on every platform that Disabled on what SplitButton returns disables both parts: neither the click nor the menu happens.
func TestSplitButtonDisabled(t *testing.T) {
	for _, p := range Platforms {
		clicks, chosen := 0, ""
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() {
					if SplitButton(c, "Save", inputsSaveMenu(&chosen)).Disabled(true).Clicked() {
						clicks++
					}
				})
			})
		}, 400, 200)
		tester.Frame()
		for _, s := range []string{"Save", inputsMoreOptions} {
			if err := tester.Click(s); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
		}
		if clicks != 0 || tester.Menu() != nil {
			t.Errorf("%v: the disabled split button clicked %d times, menu %q", p, clicks, tester.Menu())
		}
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		if tester.Focused("Save") || tester.Focused(inputsMoreOptions) {
			t.Errorf("%v: Tab focused a part of the disabled split button", p)
		}
	}
}

// TestIconButton checks the icon button on every platform, alone and in a Toolbar: its size is the platform's, label names it for clicks and the keyboard, its tooltip shows label once the pointer rests on it (not while disabled), and the icon is drawn in the button's text color.
func TestIconButton(t *testing.T) {
	sizes := map[Platform][2][2]float32{
		Mac:     {{40, 24}, {36, 36}},
		Windows: {{40, 32}, {64, 36}},
		Linux:   {{34, 34}, {34, 34}},
	}
	for _, p := range Platforms {
		for i, bar := range []bool{false, true} {
			name := fmt.Sprintf("%v toolbar %v", p, bar)
			clicks := 0
			var ink ui.Color
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					ui.Column(c).Padding(40).AlignItems(ui.Start).Children(func() {
						build := func() {
							if IconButton(c, idiomTestShareIcon, "Share").Clicked() {
								clicks++
							}
							IconButton(c, inputsTestTrashIcon, "Delete").Disabled(true)
						}
						if bar {
							Toolbar(c, build)
						} else {
							ui.Row(c).Gap(8).Children(build)
						}
					})
					ink = inputsIconInk(c, p, bar)
				})
			}, 400, 300)
			tester.SetScale(2)
			tester.Frame()
			tester.Frame()
			b, ok := tester.Find("Share")
			if !ok || b.W != sizes[p][i][0] || b.H != sizes[p][i][1] {
				t.Fatalf("%s: the button is %v, %v, want %v", name, b, ok, sizes[p][i])
			}
			// The share icon's upright stroke, 1.3 wide along x = 8 of its 16, read where it covers two whole pixels at 2x.
			x, y := b.X+(b.W-16)/2, b.Y+(b.H-16)/2
			got := tester.Image().RGBAAt(int(2*(x+8))-1, int(2*(y+6)))
			if d := max(absDiff(got.R, ink.R), absDiff(got.G, ink.G), absDiff(got.B, ink.B)); d > 24 {
				t.Errorf("%s: the icon reads %v, want the button's text color %v", name, got, ink)
			}
			tester.Move(b.X+b.W/2, b.Y+b.H/2)
			tester.Advance(700 * time.Millisecond)
			tester.Frame()
			if n := strings.Count(strings.Join(tester.Texts(), "\n"), "Share"); n != 2 {
				t.Errorf("%s: the tooltip did not show the label once the pointer rested: %q", name, tester.Texts())
			}
			d, _ := tester.Find("Delete")
			tester.Move(d.X+d.W/2, d.Y+d.H/2)
			tester.Advance(time.Second)
			tester.Frame()
			if n := strings.Count(strings.Join(tester.Texts(), "\n"), "Delete"); n != 1 {
				t.Errorf("%s: the disabled button showed its tooltip", name)
			}
			if err := tester.Click("Share"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			if clicks != 1 {
				t.Fatalf("%s: a click clicked %d times", name, clicks)
			}
			tester.Click("Share")
			tester.Frame()
			tester.Key(0, ui.KeySpace)
			tester.Frame()
			if clicks != 3 {
				t.Errorf("%s: Space on the focused button clicked %d times in all", name, clicks)
			}
		}
	}
}

// inputsIconInk is the color the platform's button gives its icon at rest: the label's.
func inputsIconInk(c *ui.Context, p Platform, bar bool) ui.Color {
	t := c.Theme()
	switch {
	case p == Windows:
		return t.Text
	case p == Linux && bar:
		return adwTextOn(t, t.Background)
	case p == Linux:
		return adwTextOn(t, t.Surface)
	case bar:
		return idiomMacInk(t, false)
	}
	return tint(macFace(t.Background, t.Dark), t.Dark, 0.847)
}

func absDiff(a, b uint8) int {
	if a > b {
		return int(a - b)
	}
	return int(b - a)
}

// inputsState is what inputsView's widgets edit.
type inputsState struct {
	empty, secret, chosen string
}

// inputsView lays out the group's widgets on platform p as tools/appkit-reference/reference.swift places AppKit's (inputsSet): an empty password field with its placeholder and one holding a password, a split button, an icon button and a disabled one; below them, a toolbar holding a split button and two icon buttons side by side.
func inputsView(p Platform, s *inputsState) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Box(c).Fill().Children(func() {
				at := func(x, y float32, build func()) {
					ui.Box(c).Absolute().Left(x).Top(y).Children(build)
				}
				at(16, 16, func() { PasswordField(c, &s.empty).Placeholder("Password").Width(180) })
				at(220, 16, func() { PasswordField(c, &s.secret).Label("Secret").Width(180) })
				at(16, 64, func() { SplitButton(c, "Save", inputsSaveMenu(&s.chosen)) })
				at(160, 64, func() { IconButton(c, idiomTestShareIcon, "Share") })
				at(220, 64, func() { IconButton(c, inputsTestTrashIcon, "Delete").Disabled(true) })
				ui.Box(c).Absolute().Left(0).Top(110).Width(520).Children(func() {
					Toolbar(c, func() {
						SplitButton(c, "New", inputsSaveMenu(&s.chosen))
						IconButton(c, idiomTestShareIcon, "Share")
						IconButton(c, inputsTestTrashIcon, "Delete")
					})
				})
			})
		})
	}
}

// TestScreenshotsInputs writes inputsView on every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names, the second password field focused with a password typed into it, which shows Windows's reveal button. Without it the test only builds the frames.
func TestScreenshotsInputs(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			s := &inputsState{}
			tester := ui.NewTester(inputsView(p, s), 520, 180)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.Frame()
			tester.Frame()
			if err := tester.Click("Secret"); err != nil {
				t.Fatalf("%v: %v", p, err)
			}
			tester.Type("hunter2")
			tester.Frame()
			for _, text := range []string{"Save", "New", inputsMoreOptions, "Share", "Delete"} {
				if _, ok := tester.Find(text); !ok {
					t.Fatalf("%v: %q is missing", p, text)
				}
			}
			if _, ok := tester.Find("Show Password"); ok != (p == Windows) {
				t.Fatalf("%v: the reveal button shows: %v", p, ok)
			}
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			writeShot(t, dir, fmt.Sprintf("breeze-inputs-%s-%s.png", strings.ToLower(p.String()), mode), tester)
		}
	}
}
