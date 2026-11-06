package breeze

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func init() {
	for _, n := range []string{"SidebarItemImage", "SidebarSectionWithActions"} {
		drawnWidgets[n] = true
	}
}

// sidebarTestPlus is the add button's icon of the tests, drawn in the color of the text.
var sidebarTestPlus = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><path d="M8 3v10M3 8h10" fill="none" stroke="#000" stroke-width="1.5" stroke-linecap="round"/></svg>`))

// sidebarTestRed and sidebarTestBlue are colored marks for SidebarItemImage: rounded squares, red and blue, the blue one with a white glyph, as an app's tiles.
var (
	sidebarTestRed  = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><rect width="16" height="16" rx="4" fill="#ff0000"/></svg>`))
	sidebarTestBlue = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><rect width="16" height="16" rx="4" fill="#1e64dc"/><path d="M5 4h6v8l-3-2.5-3 2.5z" fill="#fff"/></svg>`))
)

// sidebarState is what the sidebar of the tests edits.
type sidebarState struct {
	selected      string
	open          bool
	added         int
	addedSmart    int
	sidebarChange int
}

// sidebarView lays out a sidebar of two sections with actions, the first holding items with tinted icons under an add button, the second, which does not close, items with colored images under its own add button, as Mail's Mailboxes and Smart Mailboxes.
func sidebarView(p Platform, s *sidebarState) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Row(c).Fill().AlignItems(ui.Stretch).Children(func() {
				sb := Sidebar(c, &s.selected, func() {
					SidebarSectionWithActions(c, "Mailboxes", &s.open, func() {
						if IconButton(c, sidebarTestPlus, "Add Mailbox").Clicked() {
							s.added++
						}
					}, func() {
						SidebarItem(c, "inbox", navTestFolder, "Inbox").Children(func() { Badge(c, "12") })
						SidebarItem(c, "sent", navTestFolder, "Sent")
					})
					SidebarSectionWithActions(c, "Smart Mailboxes", nil, func() {
						if IconButton(c, sidebarTestPlus, "Add Smart Mailbox").Clicked() {
							s.addedSmart++
						}
					}, func() {
						SidebarItemImage(c, "flagged", sidebarTestRed, "Flagged")
						SidebarItemImage(c, "today", sidebarTestBlue, "Today")
					})
				}).Width(220).Label("Folders")
				if sb.Changed() {
					s.sidebarChange++
				}
				ui.Column(c).Grow(1)
			})
		})
	}
}

// TestScreenshotsSidebarImagesAndActions writes the sidebar of sidebarView on every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names: Flagged chosen without the focus, then with it after a click. Without it the test only builds the frames.
func TestScreenshotsSidebarImagesAndActions(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			for _, focus := range []bool{false, true} {
				s := &sidebarState{selected: "flagged", open: true}
				tester := ui.NewTester(sidebarView(p, s), 280, 300)
				tester.SetScale(2)
				tester.SetDark(dark)
				tester.Frame()
				suffix := ""
				if focus {
					if err := tester.Click("Flagged"); err != nil {
						t.Fatal(err)
					}
					tester.Move(270, 290)
					tester.Frame()
					suffix = "-focused"
				}
				for _, text := range []string{"Mailboxes", "Inbox", "Smart Mailboxes", "Flagged", "Today"} {
					if !tester.HasText(text) {
						t.Fatalf("%v: %q is missing", p, text)
					}
				}
				writeShot(t, dir, fmt.Sprintf("breeze-sidebar-%s-%s%s.png", strings.ToLower(p.String()), mode, suffix), tester)
			}
		}
	}
}

// sidebarInk reports whether the region of tester's last frame, at scale 2, between the sidebar's start edge and the text of an item (its end edge and the text right to left) holds a pixel of color want: where the item's icon is.
func sidebarInk(t *testing.T, tester *ui.Tester, rtl bool, text string, want ui.Color) bool {
	t.Helper()
	sidebar, _ := tester.Find("Folders")
	r, ok := tester.Find(text)
	if !ok {
		t.Fatalf("%q is missing", text)
	}
	from, to := sidebar.X, r.X-1
	if rtl {
		from, to = r.X+r.W+1, sidebar.X+sidebar.W
	}
	img := tester.Image()
	for y := int(r.Y * 2); y < int((r.Y+r.H)*2); y++ {
		for x := int(from * 2); x < int(to*2); x++ {
			px := img.RGBAAt(x, y)
			if ui.RGB(px.R, px.G, px.B) == want {
				return true
			}
		}
	}
	return false
}

// TestSidebarIconTint checks the color of a sidebar item's icon, light and dark: on macOS the accent, as AppKit draws a source list's template images (#0070f5 in Aqua and #007aff in Dark Aqua, measured by tools/appkit-reference's renderSidebar on macOS 27.2), on an item not chosen and on the chosen one without the focus, and white on the chosen one while the sidebar has the focus; on Windows and GNOME the text's color, never the accent.
func TestSidebarIconTint(t *testing.T) {
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			s := newNavState()
			var th *ui.Theme
			view := navView(p, s)
			tester := ui.NewTester(func(c *ui.Context) {
				view(c)
				Preview(c, p, func() { th = Theme(c, p) })
			}, 800, 420)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.Frame()
			name := fmt.Sprintf("%v dark %v", p, dark)
			switch p {
			case Mac:
				accent := pick(dark, ui.Hex("#0070f5"), ui.Hex("#007aff"))
				if !sidebarInk(t, tester, false, "Sent", accent) {
					t.Fatalf("%s: the icon of an item not chosen is not %v", name, accent)
				}
				if !sidebarInk(t, tester, false, "Inbox", accent) {
					t.Fatalf("%s: the icon of the chosen item without the focus is not %v", name, accent)
				}
				if sidebarInk(t, tester, false, "Inbox", white) {
					t.Fatalf("%s: the icon of the chosen item without the focus is white", name)
				}
				if err := tester.Click("Inbox"); err != nil {
					t.Fatal(err)
				}
				tester.Move(700, 400)
				tester.Frame()
				if !sidebarInk(t, tester, false, "Inbox", white) {
					t.Fatalf("%s: the icon of the chosen item with the focus is not white", name)
				}
				if !sidebarInk(t, tester, false, "Sent", accent) {
					t.Fatalf("%s: the icon of an item not chosen is not %v with the focus", name, accent)
				}
			default:
				ink := th.Text
				if p == Linux {
					ink = adwSidebarFg(th)
				}
				if !sidebarInk(t, tester, false, "Sent", ink) || sidebarInk(t, tester, false, "Sent", th.Accent) {
					t.Fatalf("%s: the icon of an item is not in %v", name, ink)
				}
			}
		}
	}
}

// sidebarImageView lays out a sidebar of two items, Alpha, chosen, and Beta: Alpha shows sidebarTestRed, tinted as SidebarItem's icon, or with image, in its own colors with SidebarItemImage.
func sidebarImageView(p Platform, selected *string, image bool) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Row(c).Fill().AlignItems(ui.Stretch).Children(func() {
				Sidebar(c, selected, func() {
					if image {
						SidebarItemImage(c, "a", sidebarTestRed, "Alpha")
					} else {
						SidebarItem(c, "a", sidebarTestRed, "Alpha")
					}
					SidebarItem(c, "b", nil, "Beta")
				}).Width(200).Label("Folders")
				ui.Column(c).Grow(1)
			})
		})
	}
}

// TestSidebarItemImageKeepsItsColors checks on every platform, light and dark, left to right and right to left, that SidebarItemImage shows a red SVG red on the chosen row while the sidebar has the focus from the keyboard, where macOS fills the row with the accent and whitens an icon, 16 points square where SidebarItem shows its icon, and that the row is otherwise drawn as SidebarItem's, its focus ring, selection and Windows's pill included: the two frames differ only inside the image's square. It also checks that a click on the image chooses the item and that typing the first letter of its label does.
func TestSidebarItemImageKeepsItsColors(t *testing.T) {
	red := ui.Hex("#ff0000")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
				name := fmt.Sprintf("%v dark %v %v", p, dark, map[ui.Direction]string{ui.LeftToRight: "ltr", ui.RightToLeft: "rtl"}[dir])
				var testers [2]*ui.Tester
				selected := [2]string{"a", "a"}
				for i, image := range []bool{false, true} {
					tester := ui.NewTester(sidebarImageView(p, &selected[i], image), 320, 160)
					tester.SetScale(2)
					tester.SetDark(dark)
					tester.SetDirection(dir)
					tester.Frame()
					tester.Key(0, ui.KeyTab)
					tester.Frame()
					if !tester.Focused("Folders") {
						t.Fatalf("%s: Tab did not give the sidebar the focus", name)
					}
					testers[i] = tester
				}
				tinted, shown := testers[0].Image(), testers[1].Image()
				// The image's square: the red pixels of the frame with the image.
				x0, y0, x1, y1 := shown.Bounds().Max.X, shown.Bounds().Max.Y, -1, -1
				for y := 0; y < shown.Bounds().Max.Y; y++ {
					for x := 0; x < shown.Bounds().Max.X; x++ {
						px := shown.RGBAAt(x, y)
						if ui.RGB(px.R, px.G, px.B) == red {
							x0, y0, x1, y1 = min(x0, x), min(y0, y), max(x1, x), max(y1, y)
						}
					}
				}
				if x1-x0+1 != 32 || y1-y0+1 != 32 {
					t.Fatalf("%s: the image's red is %d by %d pixels at 2x, want 32 by 32", name, x1-x0+1, y1-y0+1)
				}
				alpha, _ := testers[1].Find("Alpha")
				if cy := float32(y0+y1+1) / 4; cy-(alpha.Y+alpha.H/2) > 0.5 || alpha.Y+alpha.H/2-cy > 0.5 {
					t.Fatalf("%s: the image is centered at %v, the label at %v", name, cy, alpha.Y+alpha.H/2)
				}
				if px := tinted.RGBAAt((x0+x1)/2, (y0+y1)/2); ui.RGB(px.R, px.G, px.B) == red {
					t.Fatalf("%s: SidebarItem's icon is red: the test does not tell the two apart", name)
				}
				for y := 0; y < shown.Bounds().Max.Y; y++ {
					for x := 0; x < shown.Bounds().Max.X; x++ {
						if x >= x0-1 && x <= x1+1 && y >= y0-1 && y <= y1+1 {
							continue
						}
						if a, b := tinted.RGBAAt(x, y), shown.RGBAAt(x, y); a != b {
							t.Fatalf("%s: at %d,%d SidebarItem's frame is %v and SidebarItemImage's %v, outside the image", name, x, y, a, b)
						}
					}
				}

				tester := testers[1]
				tester.Key(0, ui.KeyDown)
				tester.Frame()
				if selected[1] != "b" {
					t.Fatalf("%s: Down chose %q, want b", name, selected[1])
				}
				tester.Key(0, ui.KeyA)
				tester.Frame()
				if selected[1] != "a" {
					t.Fatalf("%s: typing A chose %q, want a", name, selected[1])
				}
				selected[1] = "b"
				tester.Frame()
				tester.ClickAt(float32(x0+x1)/4, float32(y0+y1)/4)
				tester.Frame()
				if selected[1] != "a" {
					t.Fatalf("%s: a click on the image chose %q, want a", name, selected[1])
				}
			}
		}
	}
}

// TestSidebarItemImageShowsABitmap checks on every platform that SidebarItemImage scales a bitmap to its 16-point square: a 64-pixel green square shows 32 pixels square at 2x, green.
func TestSidebarItemImageShowsABitmap(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 64, 64))
	draw.Draw(src, src.Bounds(), &image.Uniform{C: color.RGBA{0, 255, 0, 255}}, image.Point{}, draw.Src)
	bitmap := ui.NewBitmap(src)
	green := ui.Hex("#00ff00")
	for _, p := range Platforms {
		selected := "b"
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				Sidebar(c, &selected, func() {
					SidebarItemImage(c, "a", bitmap, "Alpha")
					SidebarItem(c, "b", nil, "Beta")
				}).Width(200).Label("Folders")
			})
		}, 320, 160)
		tester.SetScale(2)
		tester.Frame()
		img := tester.Image()
		n := 0
		for y := 0; y < img.Bounds().Max.Y; y++ {
			for x := 0; x < img.Bounds().Max.X; x++ {
				px := img.RGBAAt(x, y)
				if ui.RGB(px.R, px.G, px.B) == green {
					n++
				}
			}
		}
		if n < 30*30 || n > 32*32 {
			t.Fatalf("%v: the bitmap shows %d green pixels at 2x, want a 32-pixel square", p, n)
		}
	}
}

// TestSidebarSectionActions checks on every platform, left to right and right to left, that a section's actions click without showing or hiding its items or choosing one, that a click on the free part of the title still shows and hides them, that each action is a Tab stop after the sidebar, which Enter clicks, and that the actions are centered on the title, 12 from the sidebar's end in a section that does not close (2 points inside the selection's inset on macOS, 12 from the pane's edge on Windows, at the heading's end on GNOME).
func TestSidebarSectionActions(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			name := fmt.Sprintf("%v %v", p, map[ui.Direction]string{ui.LeftToRight: "ltr", ui.RightToLeft: "rtl"}[dir])
			s := &sidebarState{selected: "inbox", open: true}
			tester := ui.NewTester(sidebarView(p, s), 280, 320)
			tester.SetDirection(dir)
			tester.Frame()
			sidebar, _ := tester.Find("Folders")
			for _, pair := range [][2]string{{"Mailboxes", "Add Mailbox"}, {"Smart Mailboxes", "Add Smart Mailbox"}} {
				title, _ := tester.Find(pair[0])
				action, ok := tester.Find(pair[1])
				if !ok {
					t.Fatalf("%s: %q is missing", name, pair[1])
				}
				if d := (action.Y + action.H/2) - (title.Y + title.H/2); d > 0.5 || d < -0.5 {
					t.Fatalf("%s: %q is centered %v below its title", name, pair[1], d)
				}
			}
			smart, _ := tester.Find("Add Smart Mailbox")
			if dir == ui.LeftToRight && smart.X+smart.W != sidebar.X+sidebar.W-12 || dir == ui.RightToLeft && smart.X != sidebar.X+12 {
				t.Fatalf("%s: the action of a section that does not close spans %v to %v in a sidebar of %v to %v", name, smart.X, smart.X+smart.W, sidebar.X, sidebar.X+sidebar.W)
			}
			// The first item stays clear of the action above it.
			if flagged, _ := tester.Find("Flagged"); flagged.Y < smart.Y+smart.H {
				t.Fatalf("%s: Flagged at %v is under the action, which ends at %v", name, flagged.Y, smart.Y+smart.H)
			}

			for _, label := range []string{"Add Mailbox", "Add Smart Mailbox"} {
				if err := tester.Click(label); err != nil {
					t.Fatalf("%s: click %s: %v", name, label, err)
				}
				tester.Frame()
			}
			if s.added != 1 || s.addedSmart != 1 || !s.open || s.selected != "inbox" || s.sidebarChange != 0 || !tester.HasText("Sent") {
				t.Fatalf("%s: the actions clicked %d and %d times, left the section open %v, chose %q", name, s.added, s.addedSmart, s.open, s.selected)
			}
			if err := tester.Click("Mailboxes"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			if s.open || tester.HasText("Sent") {
				t.Fatalf("%s: a click on the title left the section open", name)
			}
			if err := tester.Click("Mailboxes"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			if !s.open {
				t.Fatalf("%s: a second click on the title left the section closed", name)
			}

			if err := tester.Click("Sent"); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			if !tester.Focused("Add Mailbox") {
				t.Fatalf("%s: Tab from the sidebar did not go to the first action", name)
			}
			tester.Key(0, ui.KeyEnter)
			tester.Frame()
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			if !tester.Focused("Add Smart Mailbox") {
				t.Fatalf("%s: Tab from the first action did not go to the second", name)
			}
			if s.added != 2 || !s.open {
				t.Fatalf("%s: Enter on the action clicked it to %d, left the section open %v", name, s.added, s.open)
			}
		}
	}
}

// TestSidebarSectionActionsKeepTheTitle checks on every platform, light and dark, that a section's title and items are where SidebarSection puts them, its title drawn the same, when its actions are no taller than the title's row; and that taller ones grow the row as much above as below, the items following as much lower.
func TestSidebarSectionActionsKeepTheTitle(t *testing.T) {
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			name := fmt.Sprintf("%v dark %v", p, dark)
			// view builds a sidebar of a section whose actions are none (0), an 8-point square (1) or IconButton (2).
			view := func(kind int) func(c *ui.Context) {
				return func(c *ui.Context) {
					Preview(c, p, func() {
						Use(c)
						selected, open := "x", true
						Sidebar(c, &selected, func() {
							items := func() { SidebarItem(c, "inbox", navTestFolder, "Inbox") }
							switch kind {
							case 0:
								SidebarSection(c, "Mailboxes", &open, items)
							case 1:
								SidebarSectionWithActions(c, "Mailboxes", &open, func() { ui.Box(c).Size(8, 8).Background(ui.Hex("#ff0000")).Label("Mark") }, items)
							default:
								SidebarSectionWithActions(c, "Mailboxes", &open, func() { IconButton(c, sidebarTestPlus, "Add Mailbox") }, items)
							}
						}).Size(220, 200).Label("Folders")
					})
				}
			}
			var titles, items [3]ui.Rect
			var testers [3]*ui.Tester
			for kind := range 3 {
				tester := ui.NewTester(view(kind), 240, 200)
				tester.SetScale(2)
				tester.SetDark(dark)
				tester.Frame()
				titles[kind], _ = tester.Find("Mailboxes")
				items[kind], _ = tester.Find("Inbox")
				testers[kind] = tester
			}
			if titles[1] != titles[0] || items[1] != items[0] {
				t.Fatalf("%s: with a small action the title is at %v and the item at %v, without %v and %v", name, titles[1], items[1], titles[0], items[0])
			}
			// The title's pixels, from the sidebar's start to the small action.
			mark, _ := testers[1].Find("Mark")
			a, b := testers[0].Image(), testers[1].Image()
			for y := int(titles[0].Y * 2); y < int((titles[0].Y+titles[0].H)*2); y++ {
				for x := 0; x < int(mark.X*2)-2; x++ {
					if a.RGBAAt(x, y) != b.RGBAAt(x, y) {
						t.Fatalf("%s: the title's pixel at %d,%d is %v with an action, %v without", name, x, y, b.RGBAAt(x, y), a.RGBAAt(x, y))
					}
				}
			}
			button, _ := testers[2].Find("Add Mailbox")
			grow := titles[2].Y - titles[0].Y
			if grow < 0 {
				t.Fatalf("%s: a tall action moved the title up by %v", name, -grow)
			}
			if d := items[2].Y - items[0].Y - 2*grow; d > 0.01 || d < -0.01 {
				t.Fatalf("%s: the title moved %v lower and the item %v, want twice as much", name, grow, items[2].Y-items[0].Y)
			}
			if items[2].Y < button.Y+button.H {
				t.Fatalf("%s: the item at %v is under the action, which ends at %v", name, items[2].Y, button.Y+button.H)
			}
		}
	}
}

// TestSidebarSectionActionsSettleWhenHidden checks on every platform that a hidden sidebar whose section has actions, which then have no box, does not ask for frames forever, and that it lays them out when shown.
func TestSidebarSectionActionsSettleWhenHidden(t *testing.T) {
	for _, p := range Platforms {
		hidden := true
		s := &sidebarState{selected: "inbox", open: true}
		view := sidebarView(p, s)
		tester := ui.NewTester(func(c *ui.Context) {
			ui.Column(c).Fill().Hidden(hidden).Children(func() { view(c) })
		}, 280, 320)
		tester.Frame()
		if tester.FramePending() {
			t.Fatalf("%v: a hidden sidebar with section actions still asks for a frame", p)
		}
		hidden = false
		tester.Frame()
		title, _ := tester.Find("Mailboxes")
		action, ok := tester.Find("Add Mailbox")
		if !ok || tester.FramePending() {
			t.Fatalf("%v: the shown sidebar has its action %v and asks for a frame %v", p, ok, tester.FramePending())
		}
		if d := (action.Y + action.H/2) - (title.Y + title.H/2); d > 0.5 || d < -0.5 {
			t.Fatalf("%v: once shown, the action is centered %v below its title", p, d)
		}
	}
}
