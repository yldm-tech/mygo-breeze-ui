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
	for _, n := range []string{"FlatButton", "FlatIconButton", "FlatToggle"} {
		drawnWidgets[n] = true
	}
}

// flatDarkTheme is platform p's dark theme, which a view gives a dark part of a light window, as a player's controls over its video.
func flatDarkTheme(p Platform) *ui.Theme {
	switch p {
	case Windows:
		return fluentTheme(true, ui.Color{})
	case Linux:
		return adwTheme(true, ui.Color{})
	}
	return macTheme(true, ui.Color{})
}

// flatLayers returns the translucent fills platform p lays over what is behind a flat button under the pointer and pressed, and over a flat toggle on at rest, in a light or dark theme, as the platforms give them: on macOS black or white at 7.8%, 16.9% and 12.5% (measured from AppKit on macOS 27.2: 235, 212 and 223 over white); on Windows SubtleFillColorSecondary and SubtleFillColorTertiary (#09000000 and #06000000, #0FFFFFFF and #0AFFFFFF; the toggle on is the accent, not a layer); on GNOME currentColor at 7%, 16% and 10%, window_fg_color being 80% opaque in the light theme.
func flatLayers(p Platform, dark bool) (hover, pressed, on ui.Color) {
	ink := pick(dark, black, white)
	switch p {
	case Windows:
		if dark {
			return white.Alpha(0x0f / 255.0), white.Alpha(0x0a / 255.0), ui.Color{}
		}
		return black.Alpha(0x09 / 255.0), black.Alpha(0x06 / 255.0), ui.Color{}
	case Linux:
		if dark {
			return white.Alpha(0.07), white.Alpha(0.16), white.Alpha(0.10)
		}
		return adwInk.Alpha(0.07 * 0.8), adwInk.Alpha(0.16 * 0.8), adwInk.Alpha(0.10 * 0.8)
	}
	return ink.Alpha(0.078), ink.Alpha(0.169), ink.Alpha(0.125)
}

// flatClose reports whether two colors are within tol of each other in every channel.
func flatClose(got, want ui.Color, tol int) bool {
	return max(absDiff(got.R, want.R), absDiff(got.G, want.G), absDiff(got.B, want.B)) <= tol
}

// flatPixel reads the tester's image at x, y in points, at the scale of 2.
func flatPixel(tester *ui.Tester, x, y float32) ui.Color {
	px := tester.Image().RGBAAt(int(2*x), int(2*y))
	return ui.RGB(px.R, px.G, px.B)
}

// TestFlatButtonSizes checks the flat buttons' sizes on every platform, alone and in a Toolbar: macOS's 24-point accessory bar button, an image one in a 24-point square; a Button of WinUI's SubtleButtonStyle, 40 by 29 around an icon and 32 high around a label; libadwaita's 34-pixel .flat buttons; in a toolbar the toolbar's own items, IconButton's and Toggle's sizes there.
func TestFlatButtonSizes(t *testing.T) {
	type size struct{ icon, text [2]float32 }
	want := map[Platform][2]size{
		Mac:     {{[2]float32{24, 24}, [2]float32{0, 24}}, {[2]float32{36, 36}, [2]float32{0, 36}}},
		Windows: {{[2]float32{40, 29}, [2]float32{0, 32}}, {[2]float32{64, 36}, [2]float32{0, 36}}},
		Linux:   {{[2]float32{34, 34}, [2]float32{0, 34}}, {[2]float32{34, 34}, [2]float32{0, 34}}},
	}
	for _, p := range Platforms {
		for i, bar := range []bool{false, true} {
			name := fmt.Sprintf("%v toolbar %v", p, bar)
			on := false
			var text ui.Rect
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() {
						build := func() {
							FlatIconButton(c, idiomTestShareIcon, "Share")
							FlatToggle(c, &on, idiomTestSidebarIcon, "Sidebar")
							text = FlatButton(c, "Skip").Bounds()
						}
						if bar {
							Toolbar(c, build)
						} else {
							ui.Row(c).Gap(8).Children(build)
						}
					})
				})
			}, 400, 200)
			tester.Frame()
			tester.Frame()
			w := want[p][i]
			for _, n := range []string{"Share", "Sidebar"} {
				if b, ok := tester.Find(n); !ok || b.W != w.icon[0] || b.H != w.icon[1] {
					t.Errorf("%s: %s is %v, %v, want %v", name, n, b, ok, w.icon)
				}
			}
			if text.H != w.text[1] {
				t.Errorf("%s: Skip is %v, want %v high", name, text, w.text[1])
			}
		}
	}
	// macOS pads a title 12 points either side, as a push button: AppKit's "Skip" is 51 points wide where its title is 27.
	var flat, push ui.Rect
	tester := ui.NewTester(func(c *ui.Context) {
		Preview(c, Mac, func() {
			Use(c)
			ui.Row(c).Padding(16).AlignItems(ui.Start).Gap(8).Children(func() {
				flat = FlatButton(c, "Skip").Bounds()
				push = Button(c, "Skip").Bounds()
			})
		})
	}, 400, 200)
	tester.Frame()
	tester.Frame()
	if flat.W != push.W {
		t.Errorf("macOS: the flat button is %v wide, a push button around the same title %v", flat.W, push.W)
	}
}

// TestFlatButtonBehavior checks on every platform that the flat buttons are buttons: a click and Space click them, Tab reaches them in order and skips a disabled one, a disabled one neither clicks nor shows its tooltip, an icon button's and a toggle's label shows as their tooltip once the pointer rests on them, Menu opens the system's menu from a flat icon button, and a flat toggle turns its value over on a click and on Space.
func TestFlatButtonBehavior(t *testing.T) {
	for _, p := range Platforms {
		clicks, skips, chosen, on := 0, 0, "", false
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Row(c).Padding(40).Gap(8).AlignItems(ui.Start).Children(func() {
					if FlatIconButton(c, idiomTestShareIcon, "Share").Clicked() {
						clicks++
					}
					FlatIconButton(c, inputsTestTrashIcon, "Delete").Disabled(true)
					FlatIconButton(c, idiomTestShareIcon, "More").Menu(inputsSaveMenu(&chosen))
					if FlatButton(c, "Skip").Clicked() {
						skips++
					}
					FlatToggle(c, &on, idiomTestSidebarIcon, "Sidebar")
				})
			})
		}, 500, 300)
		tester.Frame()
		tester.Frame()
		if err := tester.Click("Share"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		tester.Key(0, ui.KeySpace)
		tester.Frame()
		if clicks != 2 {
			t.Errorf("%v: a click and Space clicked %d times", p, clicks)
		}
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		if !tester.Focused("More") {
			t.Errorf("%v: Tab from Share did not skip the disabled button to More", p)
		}
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		tester.Key(0, ui.KeyEnter)
		tester.Frame()
		if !tester.Focused("Skip") || skips != 1 {
			t.Errorf("%v: Tab and Enter reached Skip %v and clicked it %d times", p, tester.Focused("Skip"), skips)
		}
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		tester.Key(0, ui.KeySpace)
		tester.Frame()
		if !tester.Focused("Sidebar") || !on {
			t.Errorf("%v: Space on the focused toggle left it %v", p, on)
		}
		if err := tester.Click("Sidebar"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		if on {
			t.Errorf("%v: a click on the toggle did not turn it off", p)
		}
		if err := tester.Click("Delete"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		if clicks != 2 {
			t.Errorf("%v: a click on the disabled button clicked", p)
		}
		for _, n := range []string{"Share", "Sidebar"} {
			b, _ := tester.Find(n)
			tester.Move(b.X+b.W/2, b.Y+b.H/2)
			tester.Advance(700 * time.Millisecond)
			tester.Frame()
			if got := strings.Count(strings.Join(tester.Texts(), "\n"), n); got != 2 {
				t.Errorf("%v: the tooltip did not show %q once the pointer rested: %q", p, n, tester.Texts())
			}
		}
		d, _ := tester.Find("Delete")
		tester.Move(d.X+d.W/2, d.Y+d.H/2)
		tester.Advance(time.Second)
		tester.Frame()
		if got := strings.Count(strings.Join(tester.Texts(), "\n"), "Delete"); got != 1 {
			t.Errorf("%v: the disabled button showed its tooltip", p)
		}
		if err := tester.Click("More"); err != nil {
			t.Fatal(err)
		}
		if got := strings.Join(tester.Menu(), ","); got != "Save As…,Save All" {
			t.Fatalf("%v: a click on More opened %q", p, got)
		}
		if err := tester.ChooseMenuItem("Save All"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		if chosen != "Save All" {
			t.Errorf("%v: the menu chose %q", p, chosen)
		}
	}
}

// TestFlatButtonFaces checks the faces on every platform, light and dark, over the window and over a colored card: at rest the button shows what is behind it, under the pointer and pressed the platform's translucent fill laid over it, and a flat toggle on its fill at rest (the accent on Windows). On macOS the faces over the window are AppKit's as measured on macOS 27.2 (appkit-flat-*.png): 255, 235 under the pointer, 212 pressed and 223 for a toggle on in Aqua; 30, 48, 68 and 58 in Dark Aqua.
func TestFlatButtonFaces(t *testing.T) {
	card := ui.RGB(120, 160, 200)
	measured := map[bool][4]uint8{false: {255, 235, 212, 223}, true: {30, 48, 68, 58}}
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			var th *ui.Theme
			on := true
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					th = c.Theme()
					ui.Column(c).Padding(16).Gap(16).AlignItems(ui.Start).Children(func() {
						ui.Row(c).Gap(8).Children(func() {
							FlatIconButton(c, idiomTestShareIcon, "Window")
							FlatToggle(c, &on, idiomTestSidebarIcon, "Window on")
						})
						ui.Row(c).Gap(8).Padding(8).Background(card).Children(func() {
							FlatIconButton(c, idiomTestShareIcon, "Card")
							FlatToggle(c, &on, idiomTestSidebarIcon, "Card on")
						})
					})
				})
			}, 400, 200)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.Frame()
			tester.Frame()
			hover, pressed, onLayer := flatLayers(p, dark)
			for _, under := range []struct {
				name string
				bg   ui.Color
			}{{"Window", th.Background}, {"Card", card}} {
				name := fmt.Sprintf("%v dark %v over the %s", p, dark, strings.ToLower(under.name))
				b, _ := tester.Find(under.name)
				// A point of the face clear of the icon and of the focus ring inside GNOME's edge.
				x, y := b.X+3, b.Y+b.H/2
				check := func(state string, want ui.Color) {
					if got := flatPixel(tester, x, y); !flatClose(got, want, 2) {
						t.Errorf("%s, %s: the face reads %v, want %v", name, state, got, want)
					}
				}
				check("at rest", under.bg)
				cx, cy := b.X+b.W/2, b.Y+b.H/2
				tester.Move(cx, cy)
				tester.Frame()
				check("under the pointer", hover.Over(under.bg))
				tester.Press(cx, cy)
				tester.Frame()
				check("pressed", pressed.Over(under.bg))
				tester.Release(cx, cy)
				tester.Move(390, 190)
				tester.Frame()
				tg, _ := tester.Find(under.name + " on")
				want := th.Accent
				if p != Windows {
					want = onLayer.Over(under.bg)
				}
				if got := flatPixel(tester, tg.X+3, tg.Y+tg.H/2); !flatClose(got, want, 2) {
					t.Errorf("%s: the toggle on reads %v, want %v", name, got, want)
				}
				if p != Mac || under.name != "Window" {
					continue
				}
				// AppKit's own values over the window.
				m := measured[dark]
				tester.Move(cx, cy)
				tester.Frame()
				hovered := flatPixel(tester, x, y)
				tester.Press(cx, cy)
				tester.Frame()
				held := flatPixel(tester, x, y)
				tester.Release(cx, cy)
				tester.Move(390, 190)
				tester.Frame()
				onFace := flatPixel(tester, tg.X+3, tg.Y+tg.H/2)
				for i, got := range []ui.Color{flatPixel(tester, x, y), hovered, held, onFace} {
					if absDiff(got.R, m[i]) > 1 {
						t.Errorf("%s: state %d reads %v, AppKit's %d", name, i, got, m[i])
					}
				}
			}
		}
	}
}

// TestMacFlatButtonInk checks the color of a macOS flat button's image against AppKit's render (appkit-flat-*.png): secondaryLabelColor at rest (128 over white, 154 over the window in Dark Aqua), labelColor over the bezel under the pointer (36 and 223) and pressed (32 and 226); a toggle on away from the pointer labelColor over its fill (34 and 225), and under the pointer white over the accent. It reads the share icon's upright stroke, 1.3 points wide along x = 8 of its 16, where it covers two whole pixels at 2x.
func TestMacFlatButtonInk(t *testing.T) {
	want := map[bool][4]uint8{false: {128, 36, 32, 34}, true: {154, 223, 226, 225}}
	for _, dark := range []bool{false, true} {
		on := true
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Mac, func() {
				Use(c)
				ui.Row(c).Padding(16).Gap(8).AlignItems(ui.Start).Children(func() {
					FlatIconButton(c, idiomTestShareIcon, "Share")
					FlatToggle(c, &on, idiomTestShareIcon, "Shared")
				})
			})
		}, 300, 100)
		tester.SetScale(2)
		tester.SetDark(dark)
		tester.Frame()
		tester.Frame()
		ink := func(name string) uint8 {
			b, _ := tester.Find(name)
			x, y := b.X+(b.W-16)/2, b.Y+(b.H-16)/2
			return tester.Image().RGBAAt(int(2*(x+8))-1, int(2*(y+6))).R
		}
		b, _ := tester.Find("Share")
		cx, cy := b.X+b.W/2, b.Y+b.H/2
		rest := ink("Share")
		tester.Move(cx, cy)
		tester.Frame()
		hovered := ink("Share")
		tester.Press(cx, cy)
		tester.Frame()
		pressed := ink("Share")
		tester.Release(cx, cy)
		tester.Move(290, 90)
		tester.Frame()
		toggled := ink("Shared")
		for i, got := range []uint8{rest, hovered, pressed, toggled} {
			if absDiff(got, want[dark][i]) > 3 {
				t.Errorf("dark %v: state %d of the image reads %d, AppKit's %d", dark, i, got, want[dark][i])
			}
		}
		g, _ := tester.Find("Shared")
		tester.Move(g.X+g.W/2, g.Y+g.H/2)
		tester.Frame()
		if got := ink("Shared"); got < 245 {
			t.Errorf("dark %v: the toggle on under the pointer draws its image at %d, want white over the accent", dark, got)
		}
		if face := flatPixel(tester, g.X+3, g.Y+g.H/2); !flatClose(face, macAccent, 2) {
			t.Errorf("dark %v: the toggle on under the pointer reads %v, want the accent", dark, face)
		}
	}
}

// TestFlatButtonInDarkPart checks on every platform that a flat button and a flat toggle in a dark part of a light window, as a player's controls over its video, take that part's dark theme: the button under the pointer lays the dark theme's fill (white) over the part's background, and the toggle on is the dark theme's fill or accent.
func TestFlatButtonInDarkPart(t *testing.T) {
	for _, p := range Platforms {
		th := flatDarkTheme(p)
		on := true
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() {
					ui.Themed(c, th).Padding(8).Children(func() {
						ui.Row(c).Gap(8).Children(func() {
							FlatIconButton(c, idiomTestShareIcon, "Share")
							FlatToggle(c, &on, idiomTestSidebarIcon, "Panel")
						})
					})
				})
			})
		}, 300, 120)
		tester.SetScale(2)
		tester.Frame()
		tester.Frame()
		hover, _, onLayer := flatLayers(p, true)
		b, _ := tester.Find("Share")
		if got := flatPixel(tester, b.X+3, b.Y+b.H/2); !flatClose(got, th.Background, 2) {
			t.Errorf("%v: the button at rest reads %v, want the dark part's %v", p, got, th.Background)
		}
		tester.Move(b.X+b.W/2, b.Y+b.H/2)
		tester.Frame()
		if got, want := flatPixel(tester, b.X+3, b.Y+b.H/2), hover.Over(th.Background); !flatClose(got, want, 2) {
			t.Errorf("%v: the button under the pointer reads %v, want %v", p, got, want)
		}
		tester.Move(290, 110)
		tester.Frame()
		g, _ := tester.Find("Panel")
		want := th.Accent
		if p != Windows {
			want = onLayer.Over(th.Background)
		}
		if got := flatPixel(tester, g.X+3, g.Y+g.H/2); !flatClose(got, want, 2) {
			t.Errorf("%v: the toggle on reads %v, want %v", p, got, want)
		}
	}
}

// TestFlatButtonFocusRing checks on every platform that the platform's focus ring shows around a flat button that has the focus from the keyboard, where nothing was drawn before: AppKit's ring 1.5 points outside its edge, WinUI's focus visual 2 DIPs outside it (FocusVisualMargin -3), libadwaita's ring a pixel inside it.
func TestFlatButtonFocusRing(t *testing.T) {
	inset := map[Platform]float32{Mac: -1.5, Windows: -2, Linux: 1}
	for _, p := range Platforms {
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Row(c).Padding(16).AlignItems(ui.Start).Children(func() { FlatIconButton(c, idiomTestShareIcon, "Share") })
			})
		}, 300, 100)
		tester.SetScale(2)
		tester.Frame()
		b, _ := tester.Find("Share")
		x, y := b.X+inset[p], b.Y+b.H/2
		before := flatPixel(tester, x, y)
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		if !tester.Focused("Share") {
			t.Fatalf("%v: Tab did not focus the button", p)
		}
		if got := flatPixel(tester, x, y); flatClose(got, before, 20) {
			t.Errorf("%v: no focus ring at %v, %v: it reads %v as before", p, x, y, got)
		}
	}
}

// flatState is what flatView's toggles edit.
type flatState struct{ off, on, darkOn bool }

// flatView lays out the flat buttons on platform p as tools/appkit-reference/reference.swift places AppKit's (flatSet): an image button, a second one under the pointer (the screenshot test moves it there), a disabled one, a text button, and an image toggle off and on, 48 points apart from 16 in; below them a dark part, as a player's controls over its video, holding the same buttons in the platform's dark theme.
func flatView(p Platform, s *flatState) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Box(c).Fill().Children(func() {
				at := func(x, y float32, build func()) {
					ui.Box(c).Absolute().Left(x).Top(y).Children(build)
				}
				row := func(y float32, toggle *bool) {
					at(16, y, func() { FlatIconButton(c, idiomTestShareIcon, "Share") })
					at(64, y, func() { FlatIconButton(c, inputsTestTrashIcon, "Delete") })
					at(112, y, func() { FlatIconButton(c, inputsTestTrashIcon, "Deleted").Disabled(true) })
					at(160, y, func() { FlatButton(c, "Skip") })
					at(240, y, func() { FlatToggle(c, &s.off, idiomTestSidebarIcon, "Sidebar") })
					at(288, y, func() { FlatToggle(c, toggle, idiomTestSidebarIcon, "Panel") })
				}
				row(16, &s.on)
				ui.Themed(c, flatDarkTheme(p)).Absolute().Left(0).Top(60).Width(360).Height(60).Children(func() { row(10, &s.darkOn) })
			})
		})
	}
}

// TestScreenshotsFlat writes flatView on every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names, the pointer over the second button of the first row. Without it the test only builds the frames.
func TestScreenshotsFlat(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			s := &flatState{on: true, darkOn: true}
			tester := ui.NewTester(flatView(p, s), 360, 120)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.Frame()
			tester.Frame()
			for _, text := range []string{"Share", "Delete", "Deleted", "Skip", "Sidebar", "Panel"} {
				if _, ok := tester.Find(text); !ok {
					t.Fatalf("%v: %q is missing", p, text)
				}
			}
			d, _ := tester.Find("Delete")
			tester.Move(d.X+d.W/2, d.Y+d.H/2)
			tester.Frame()
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			writeShot(t, dir, fmt.Sprintf("breeze-flat-%s-%s.png", strings.ToLower(p.String()), mode), tester)
		}
	}
}

// TestFlatToggleText checks on every platform that a flat toggle without an icon shows its label as its text, at the height of a flat button, and turns its value over on a click.
func TestFlatToggleText(t *testing.T) {
	for _, p := range Platforms {
		on := false
		var toggle, button ui.Rect
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Row(c).Padding(16).Gap(8).AlignItems(ui.Start).Children(func() {
					toggle = FlatToggle(c, &on, nil, "Wrap").Bounds()
					button = FlatButton(c, "Skip").Bounds()
				})
			})
		}, 300, 100)
		tester.Frame()
		tester.Frame()
		if !tester.HasText("Wrap") || toggle.H != button.H {
			t.Errorf("%v: the toggle shows its label %v, %v high, a flat button %v", p, tester.HasText("Wrap"), toggle.H, button.H)
		}
		if err := tester.Click("Wrap"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		if !on {
			t.Errorf("%v: a click left the toggle off", p)
		}
	}
}
