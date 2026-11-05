package breeze

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// adwPixel returns the color of tester's frame at x, y, in window DIPs at scale 1.
func adwPixel(tester *ui.Tester, x, y float32) ui.Color {
	px := tester.Image().RGBAAt(int(x), int(y))
	return ui.RGB(px.R, px.G, px.B)
}

// TestAdwFocusRingInside checks that on GNOME a button and a check box show libadwaita's focus ring 2 pixels inside their edge, and none outside it, where MyGo's own ring would be.
func TestAdwFocusRingInside(t *testing.T) {
	for _, dark := range []bool{false, true} {
		on := false
		var th *ui.Theme
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Linux, func() {
				Use(c)
				th = c.Theme()
				ui.Column(c).Fill().Padding(20).Gap(20).AlignItems(ui.Start).Children(func() {
					Button(c, "Cancel")
					Checkbox(c, &on, "Mail")
				})
			})
		}, 300, 200)
		tester.SetDark(dark)
		tester.Frame()
		for _, label := range []string{"Cancel", "Mail"} {
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			if !tester.Focused(label) {
				t.Fatalf("dark %v: Tab does not focus %q", dark, label)
			}
			r, _ := tester.Find(label)
			// The control's left edge: a button's text is 17 in, a check box's row 4 before its 20-pixel box and the 4 between it and the label.
			left := r.X - 17
			if label == "Mail" {
				left = r.X - 28
			}
			y := r.Y + r.H/2
			// The ring is the focus color at its opacity over the button's face, or over the window behind the check box's row.
			face := th.Surface
			if label == "Mail" {
				face = th.Background
			}
			want := face.Mix(ui.RGB(th.Focus.R, th.Focus.G, th.Focus.B), float32(th.Focus.A)/255)
			if inside := adwPixel(tester, left+1, y); adwColorDiff(inside, want) > 0.03 {
				t.Fatalf("dark %v: %q's edge is %v, want the focus ring's %v", dark, label, inside, want)
			}
			if outside := adwPixel(tester, left-2.5, y); outside != th.Background {
				t.Fatalf("dark %v: %q shows %v outside its edge, want the window's %v: MyGo's ring is drawn", dark, label, outside, th.Background)
			}
		}
	}
}

// adwColorDiff returns how far apart two opaque colors are, from 0 to 1: the largest difference of their channels.
func adwColorDiff(a, b ui.Color) float32 {
	d := func(x, y uint8) float32 {
		if x > y {
			return float32(x-y) / 255
		}
		return float32(y-x) / 255
	}
	return max(d(a.R, b.R), d(a.G, b.G), d(a.B, b.B))
}

// TestAdwSidebarFollowsAdwSidebar checks that on GNOME a sidebar's rows are AdwSidebar's, 42 pixels high, and that a section's heading is in the sidebar's text color, not dimmed, 12 pixels in from its edge.
func TestAdwSidebarFollowsAdwSidebar(t *testing.T) {
	for _, dark := range []bool{false, true} {
		s := newNavState()
		var th *ui.Theme
		view := navView(Linux, s)
		tester := ui.NewTester(func(c *ui.Context) {
			view(c)
			Preview(c, Linux, func() { th = Theme(c, Linux) })
		}, 800, 420)
		tester.SetDark(dark)
		tester.Frame()
		r, ok := tester.Find("Inbox")
		if !ok {
			t.Fatal("the chosen item is missing")
		}
		fill := adwInkAlpha(th, adwSidebarBg(th), 0.10)
		x := r.X + r.W + 4
		top, bottom := r.Y, r.Y
		for top > 0 && adwPixel(tester, x, top-1) == fill {
			top--
		}
		for adwPixel(tester, x, bottom+1) == fill {
			bottom++
		}
		if h := bottom - top + 1; h < 41 || h > 42 {
			t.Fatalf("dark %v: the chosen row is %v high, want 42", dark, h)
		}
		head, _ := tester.Find("Mailboxes")
		if head.X < 11.5 || head.X > 12.5 {
			t.Fatalf("dark %v: the heading's text starts %v in, want 12", dark, head.X)
		}
		// The heading's strongest ink, the pixel farthest from the sidebar's background, is its text color.
		bg := adwSidebarBg(th)
		strongest := bg
		for y := head.Y; y < head.Y+head.H; y++ {
			for x := head.X; x < head.X+head.W; x++ {
				if px := adwPixel(tester, x, y); adwColorDiff(px, bg) > adwColorDiff(strongest, bg) {
					strongest = px
				}
			}
		}
		if want := adwSidebarFg(th); adwColorDiff(strongest, want) > 0.12 {
			t.Fatalf("dark %v: the heading's ink reaches %v, want the sidebar's text color %v", dark, strongest, want)
		}
	}
}

// TestAdwSidebarBadge checks that on GNOME a badge in a sidebar item is AdwViewSwitcherSidebar's count, a capsule 1.6em of its smaller text high, and a badge elsewhere stays indicatorbin's 13-pixel one.
func TestAdwSidebarBadge(t *testing.T) {
	selected := "a"
	var th *ui.Theme
	tester := ui.NewTester(func(c *ui.Context) {
		Preview(c, Linux, func() {
			Use(c)
			th = c.Theme()
			ui.Row(c).Fill().AlignItems(ui.Stretch).Children(func() {
				Sidebar(c, &selected, func() {
					SidebarItem(c, "b", nil, "Inbox").Children(func() { Badge(c, "12") })
				}).Width(200)
				ui.Column(c).Padding(20).AlignItems(ui.Start).Children(func() { Badge(c, "34") })
			})
		})
	}, 400, 200)
	tester.Frame()
	height := func(text string, bg ui.Color) float32 {
		r, ok := tester.Find(text)
		if !ok {
			t.Fatalf("%q is missing", text)
		}
		// Above and below the text's line, in the middle of the capsule, away from its round ends.
		x := r.X + r.W/2
		top, bottom := r.Y, r.Y+r.H-1
		for adwPixel(tester, x, top-1) == bg {
			top--
		}
		for adwPixel(tester, x, bottom+1) == bg {
			bottom++
		}
		return bottom - top + 1
	}
	if h, want := height("12", adwInkAlpha(th, adwSidebarBg(th), 0.4)), th.FontSize/1.2*1.6; h < want-1 || h > want+1 {
		t.Fatalf("the sidebar's badge is %v high, want %v", h, want)
	}
	if h := height("34", adwInkAlpha(th, th.Background, 0.4)); h < 12 || h > 16 {
		t.Fatalf("a badge outside a sidebar is %v high, want indicatorbin's 13 or its text's line", h)
	}
}

// TestAdwAlertOutline checks that on GNOME the alert shows libadwaita 1.10.0's 1-pixel outline of white at 7% inside its edge, over the dark dialog_bg_color.
func TestAdwAlertOutline(t *testing.T) {
	open := true
	tester := ui.NewTester(func(c *ui.Context) {
		Preview(c, Linux, func() {
			Use(c)
			AlertDialog(c, &open, "Delete “Notes”?", "You can’t undo this action.", "Cancel", "Delete")
		})
	}, 560, 340)
	tester.SetDark(true)
	tester.Frame()
	r, ok := tester.Find("You can’t undo this action.")
	if !ok {
		t.Fatal("the alert's message is missing")
	}
	bg := ui.Hex("#36363a")
	// The alert is 300 wide, centered.
	edge := float32(560-300) / 2
	y := r.Y + r.H/2
	if got, want := adwPixel(tester, edge, y), bg.Mix(white, 0.07); adwColorDiff(got, want) > 0.01 {
		t.Fatalf("the alert's edge is %v, want %v", got, want)
	}
	if got := adwPixel(tester, edge+2, y); got != bg {
		t.Fatalf("inside its edge the alert is %v, want %v", got, bg)
	}
}

// TestAdwThemeScrollbar checks the GNOME theme's scroll bar: an overlay indicator's slider, currentColor at 20% and 3 pixels wide.
func TestAdwThemeScrollbar(t *testing.T) {
	for _, dark := range []bool{false, true} {
		th := adwTheme(dark, ui.Color{})
		want := pick(dark, ui.RGBA(0, 0, 6, 0.16), ui.RGBA(255, 255, 255, 0.2))
		if th.Scrollbar != want || th.ScrollbarWidth != 3 {
			t.Fatalf("dark %v: the scroll bar is %v, %v wide; want %v, 3", dark, th.Scrollbar, th.ScrollbarWidth, want)
		}
	}
}
