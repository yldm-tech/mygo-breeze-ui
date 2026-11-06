package breeze

import (
	"math"
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

// TestAdwLinkedToggleFocusRing checks that on GNOME a toggle in a ToggleGroup shows its focus ring inside its edge, which the group's rounded clip would cut off outside it.
func TestAdwLinkedToggleFocusRing(t *testing.T) {
	bold, italic := false, false
	var th *ui.Theme
	tester := ui.NewTester(func(c *ui.Context) {
		Preview(c, Linux, func() {
			Use(c)
			th = c.Theme()
			ui.Column(c).Fill().Padding(20).AlignItems(ui.Start).Children(func() {
				ToggleGroup(c, func() {
					Toggle(c, &bold, "Bold")
					Toggle(c, &italic, "Italic")
				})
			})
		})
	}, 300, 120)
	tester.Frame()
	tester.Key(0, ui.KeyTab)
	tester.Frame()
	if !tester.Focused("Bold") {
		t.Fatal("Tab does not focus the first toggle")
	}
	r, _ := tester.Find("Bold")
	want := th.Surface.Mix(ui.RGB(th.Focus.R, th.Focus.G, th.Focus.B), float32(th.Focus.A)/255)
	// The toggle is 34 high, its text centered in it.
	top := r.Y + r.H/2 - 17
	if got := adwPixel(tester, r.X+r.W/2, top+1); adwColorDiff(got, want) > 0.03 {
		t.Fatalf("the focused toggle's top edge is %v, want the focus ring's %v", got, want)
	}
}

// TestAdwSidebarHeadingSpacing checks that on GNOME a section's heading sits where AdwSidebar puts it: the first one's text at the list's top padding of 6, 12 above its section's first row, and another's text 14 below the row above it.
func TestAdwSidebarHeadingSpacing(t *testing.T) {
	s := newNavState()
	tester := ui.NewTester(navView(Linux, s), 800, 420)
	tester.Frame()
	find := func(text string) ui.Rect {
		r, ok := tester.Find(text)
		if !ok {
			t.Fatalf("%q is missing", text)
		}
		return r
	}
	near := func(what string, got, want float32) {
		if got < want-0.5 || got > want+0.5 {
			t.Fatalf("%s is %v, want %v", what, got, want)
		}
	}
	// A row is 42 high with its label centered in it.
	rowTop := func(text string) float32 { r := find(text); return r.Y + r.H/2 - 21 }
	head, next := find("Mailboxes"), find("Smart Mailboxes")
	near("the first heading's top", head.Y, 6)
	near("the space below the first heading", rowTop("Inbox")-(head.Y+head.H), 12)
	near("the space above the second heading", next.Y-(rowTop("Drafts")+42), 14)
	near("the space below the second heading", rowTop("Flagged")-(next.Y+next.H), 12)
}

// TestAdwAlertResponsesFit checks that on GNOME an alert's responses stay inside it when their least width of 60 would not let them fit side by side, as libadwaita would stack them.
func TestAdwAlertResponsesFit(t *testing.T) {
	open := true
	tester := ui.NewTester(func(c *ui.Context) {
		Preview(c, Linux, func() {
			Use(c)
			AlertDialog(c, &open, "Save changes?", "Unsaved changes will be lost.", "Cancel", "Discard", "Later", "Save")
		})
	}, 600, 500)
	tester.Frame()
	// The alert is at most 372 wide, centered, its responses 24 in from its sides.
	right := float32(600+372)/2 - 24
	if r, ok := tester.Find("Save"); !ok || r.X+r.W > right {
		t.Fatalf("the last response's label ends at %v, past the alert's response area at %v", r.X+r.W, right)
	}
}

// adwRun returns where, along the row or column through x, y, the pixels of tester's frame stay color: the first and last of them, scanning from start to end (horizontally when across).
func adwRun(tester *ui.Tester, x, y float32, across bool, start, end float32, color ui.Color) (first, last float32, ok bool) {
	first, last = -1, -1
	for v := start; v <= end; v++ {
		px, py := x, y
		if across {
			px = v
		} else {
			py = v
		}
		if adwPixel(tester, px, py) == color {
			if first < 0 {
				first = v
			}
			last = v
		}
	}
	return first, last, first >= 0
}

// TestAdwScaleGeometry pins GtkScale as libadwaita 1.10.0 (GTK 4.24.1) renders it, measured by tools/adwaita-reference: 34 pixels high without marks, the trough 15 pixels from its top, and the slider's center 2 pixels in from the trough's start plus the rest of the trough times the value, rounded down; with marks after it 28 high, the trough 12 from its top and the marks 6 below the trough.
func TestAdwScaleGeometry(t *testing.T) {
	volume, quality := 40.0, 50.0
	var th *ui.Theme
	tester := ui.NewTester(func(c *ui.Context) {
		Preview(c, Linux, func() {
			Use(c)
			th = c.Theme()
			ui.Box(c).Fill().Children(func() {
				ui.Box(c).Absolute().Left(16).Top(20).Children(func() { Slider(c, &volume, 0, 100).Width(180) })
				ui.Box(c).Absolute().Left(16).Top(80).Children(func() { StepSlider(c, &quality, 0, 100, 25).Width(260) })
			})
		})
	}, 320, 140)
	tester.Frame()
	// The trough's filled part, at the trough's start: 12 in from the scale's edge.
	top, bottom, ok := adwRun(tester, 16+12+3, 20, false, 20, 60, th.Accent)
	if !ok || top != 20+15 || bottom != 20+18 {
		t.Fatalf("the scale's trough spans rows %v to %v, want 35 to 38", top, bottom)
	}
	// The slider is white, 20 across; at 40% of the 152 pixels it travels its center is 28+2+60 = 90.
	left, right, ok := adwRun(tester, 0, 20+17, true, 16, 200, white)
	if !ok || left+right+1 != 2*90 {
		t.Fatalf("the slider spans %v to %v, want it centered at 90", left, right)
	}
	top, bottom, ok = adwRun(tester, 16+12+3, 80, false, 80, 120, th.Accent)
	if !ok || top != 80+12 || bottom != 80+15 {
		t.Fatalf("the marked scale's trough spans rows %v to %v, want 92 to 95", top, bottom)
	}
	// The first mark: a pixel wide at the slider's place at 0, 28+2, from 6 below the trough for 6 pixels.
	top, bottom, ok = adwRun(tester, 16+12+2, 80, false, 80, 120, th.TextMuted)
	if !ok || top != 80+22 || bottom != 80+27 {
		t.Fatalf("the first mark spans rows %v to %v, want 102 to 107", top, bottom)
	}
}

// TestAdwSwitchSlider pins a GtkSwitch's slider as libadwaita 1.10.0 (GTK 4.24.1) renders it, measured by tools/adwaita-reference: off in the dark theme it is white at 80% over view_bg_color, #d2d2d2, checked it is white, and a disabled switch's slider casts no shadow.
func TestAdwSwitchSlider(t *testing.T) {
	off, disabled, on := false, false, true
	var th *ui.Theme
	tester := ui.NewTester(func(c *ui.Context) {
		Preview(c, Linux, func() {
			Use(c)
			th = c.Theme()
			ui.Row(c).Fill().Padding(20).Gap(20).AlignItems(ui.Start).Children(func() {
				Switch(c, &off)
				Switch(c, &disabled).Disabled(true)
				Switch(c, &on)
			})
		})
	}, 240, 80)
	tester.SetDark(true)
	tester.Frame()
	if got := adwPixel(tester, 20+13, 20+13); got != ui.Hex("#d2d2d2") {
		t.Fatalf("the dark slider is %v, want #d2d2d2", got)
	}
	// Below the slider, 2 pixels inside the track's bottom, the track shows as it does beside the slider: no shadow falls on it.
	x := float32(20 + 46 + 20)
	beside := adwPixel(tester, x+40, 20+13)
	if below := adwPixel(tester, x+13, 20+24); below != beside {
		t.Fatalf("below the disabled slider the track is %v, beside it %v: a shadow falls", below, beside)
	}
	// Checked, the slider is white in the dark theme too (switch:checked > slider), at the track's end.
	if got := adwPixel(tester, 2*(20+46)+20+46-13, 20+13); got != white {
		t.Fatalf("the checked dark slider is %v, want white", got)
	}
	_ = th
}

// TestAdwSelectWidth checks that on GNOME a drop-down is as wide as its choice and arrow need, as GtkDropDown sets no least width (a "Medium" drop-down is 101 pixels wide in Adwaita Sans, measured from libadwaita 1.10.0 by tools/adwaita-reference): padding of 10, the text, 6 and the 16-pixel arrow.
func TestAdwSelectWidth(t *testing.T) {
	choice := "S"
	var th *ui.Theme
	tester := ui.NewTester(func(c *ui.Context) {
		Preview(c, Linux, func() {
			Use(c)
			th = c.Theme()
			ui.Column(c).Fill().Padding(20).AlignItems(ui.Start).Children(func() { Select(c, &choice, []string{"S", "M"}) })
		})
	}, 300, 100)
	tester.Frame()
	r, ok := tester.Find("S")
	if !ok {
		t.Fatal("the choice is missing")
	}
	left, right, _ := adwRun(tester, 0, r.Y+r.H/2, true, 0, 299, th.Surface)
	want := 10 + r.W + 6 + 16 + 10
	if w := right - left + 1; w < want-1.5 || w > want+1.5 {
		t.Fatalf("the drop-down is %v wide, want %v", w, want)
	}
}

// TestAdwToastWithButtonWidth checks that on GNOME a toast with a button is 426 pixels wide however short its title, as libadwaita 1.10.0 (GTK 4.24.1) draws it (measured by tools/adwaita-reference), and one without a button only as wide as its content.
func TestAdwToastWithButtonWidth(t *testing.T) {
	for _, action := range []string{"Undo", ""} {
		added := false
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Linux, func() {
				Use(c)
				if !added {
					added = true
					c.AddToast(ui.Toast{Title: "Saved", Action: action, OnAction: func() {}, Timeout: -1})
				}
				Toasts(c)
			})
		}, 700, 300)
		tester.Frame()
		tester.Frame()
		r, ok := tester.Find("Saved")
		if !ok {
			t.Fatalf("action %q: the toast is missing", action)
		}
		left, right, _ := adwRun(tester, 0, r.Y+r.H/2, true, 0, 699, ui.Hex("#505053"))
		w := right - left + 1
		// The run of the toast's color leaves out the anti-aliased pixel at either round end.
		if action != "" && (w < 424 || w > 426) {
			t.Fatalf("the toast with a button is %v wide, want 426", w)
		}
		if action == "" && w > 200 {
			t.Fatalf("the toast without a button is %v wide, want its content's width", w)
		}
	}
}

// TestAdwSidebarRowSpan checks that on GNOME a sidebar's chosen row reaches 6 pixels from both its edges, as the line at the sidebar's end takes no room (an inset box-shadow): in a 200-pixel sidebar from 6 to 194, as libadwaita 1.10.0 (GTK 4.24.1) renders it (measured by tools/adwaita-reference), left to right and right to left.
func TestAdwSidebarRowSpan(t *testing.T) {
	for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
		s := newNavState()
		var th *ui.Theme
		view := navView(Linux, s)
		tester := ui.NewTester(func(c *ui.Context) {
			view(c)
			Preview(c, Linux, func() { th = Theme(c, Linux) })
		}, 800, 420)
		tester.SetDirection(dir)
		tester.Frame()
		r, _ := tester.Find("Inbox")
		fill := adwInkAlpha(th, adwSidebarBg(th), 0.10)
		left, right, ok := adwRun(tester, 0, r.Y+r.H/2, true, 0, 799, fill)
		want := float32(6)
		if dir == ui.RightToLeft {
			want = 800 - 200 + 6
		}
		if !ok || left != want || right != want+187 {
			t.Fatalf("dir %v: the chosen row spans %v to %v, want %v to %v", dir, left, right, want, want+187)
		}
	}
}

// TestAdwLinkedToggleSeam pins the line where two linked toggles meet: the second's fill laid over the first's, measured from libadwaita 1.10.0 (GTK 4.24.1) by tools/adwaita-reference as #afafb1 after a checked toggle in the light theme and #747476 in the dark one.
func TestAdwLinkedToggleSeam(t *testing.T) {
	for _, dark := range []bool{false, true} {
		bold, italic := true, false
		var th *ui.Theme
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Linux, func() {
				Use(c)
				th = c.Theme()
				ui.Column(c).Fill().Padding(20).AlignItems(ui.Start).Children(func() {
					ToggleGroup(c, func() {
						Toggle(c, &bold, "Bold")
						Toggle(c, &italic, "Italic")
					})
				})
			})
		}, 300, 100)
		tester.SetDark(dark)
		tester.Frame()
		r, _ := tester.Find("Bold")
		checked := adwInkAlpha(th, th.Background, 0.30)
		_, end, _ := adwRun(tester, 0, r.Y+r.H/2, true, 0, 299, checked)
		want := pick(dark, ui.Hex("#afafb1"), ui.Hex("#747476"))
		if got := adwPixel(tester, end+1, r.Y+r.H/2); adwColorDiff(got, want) > 2.0/255 {
			t.Fatalf("dark %v: the seam is %v, want %v", dark, got, want)
		}
	}
}

// TestAdwSpinButtonIcons checks that on GNOME a NumberInput's − and + icons are centered in their buttons past the 1-pixel border on the start side: the bar 12 pixels after the border, as libadwaita 1.10.0 (GTK 4.24.1) renders a GtkSpinButton (measured by tools/adwaita-reference).
func TestAdwSpinButtonIcons(t *testing.T) {
	value := 5.0
	var th *ui.Theme
	tester := ui.NewTester(func(c *ui.Context) {
		Preview(c, Linux, func() {
			Use(c)
			th = c.Theme()
			ui.Column(c).Fill().Padding(20).AlignItems(ui.Start).Children(func() { NumberInput(c, &value, 0, 10, 1) })
		})
	}, 300, 100)
	// At 2x, where half a pixel shows.
	tester.SetScale(2)
	tester.Frame()
	fill := adwInkAlpha(th, th.Background, 0.10)
	edge := adwInkAlpha(th, fill, 0.10)
	ink := adwInkAlpha(th, fill, 1)
	// The input is 34 high at the column's padding of 20; the text field before the buttons is MyGo's 80 wide. Positions here are the image's pixels.
	y := float32(2 * (20 + 17))
	border, _, ok := adwRun(tester, 0, y, true, 2*(20+40), 599, edge)
	if !ok {
		t.Fatal("the − button's border is missing")
	}
	bar, _, ok := adwRun(tester, 0, y, true, border, border+64, ink)
	if !ok || bar-border != 24 {
		t.Fatalf("the − bar starts %v pixels after the border at 2x, want 24", bar-border)
	}
}

// TestAdwColorWellWidth checks that on GNOME a color well is a GtkColorDialogButton's size: its swatch as wide as "Black" in the text, padded 5, and 34 high (49 by 34 in Adwaita Sans 11, measured from libadwaita 1.10.0 (GTK 4.24.1) by tools/adwaita-reference). It also checks a theme with 24-pixel text, where "Black" is far wider than in the default one, so that the width follows the text whatever the test font.
func TestAdwColorWellWidth(t *testing.T) {
	for _, size := range []float32{0, 24} {
		tint := ui.Hex("#3584e4")
		var th *ui.Theme
		var black float32
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, Linux, func() {
				Use(c)
				theme := *c.Theme()
				if size > 0 {
					theme.FontSize = size
				}
				ui.Themed(c, &theme).Fill().Padding(20).AlignItems(ui.Start).Children(func() {
					th = c.Theme()
					black, _ = c.MeasureText(0, ui.Span{Text: "Black"})
					ColorWell(c, &tint).Label("Tint")
				})
			})
		}, 300, 100)
		tester.Frame()
		left, right, _ := adwRun(tester, 0, 20+17, true, 0, 299, tint)
		if w, want := right-left+1, float32(math.Ceil(float64(black))); w < want-1 || w > want+1 {
			t.Fatalf("font size %v: the swatch is %v wide, want %v", size, w, want)
		}
		// Down the middle: the button's face above and below the swatch, the window around it.
		top, bottom := float32(-1), float32(-1)
		for y := float32(0); y < 100; y++ {
			if adwPixel(tester, (left+right)/2, y) != th.Background {
				if top < 0 {
					top = y
				}
				bottom = y
			}
		}
		if h := bottom - top + 1; h < 33 || h > 35 {
			t.Fatalf("font size %v: the button is %v high, want 34", size, h)
		}
	}
}
