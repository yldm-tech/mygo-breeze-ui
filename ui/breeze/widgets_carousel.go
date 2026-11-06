package breeze

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// ---- pagers: a row of page indicators, and a carousel of pages above them ----
//
// macOS: AppKit and SwiftUI have no page indicator on macOS: the macOS 27.2 SDK marks SwiftUI's PageTabViewStyle and IndexViewStyle @available(macOS, unavailable), and AppKit's NSPageController, which swipes between pages, shows none. The closest that macOS draws is UIKit's UIPageControl in UIKit for Mac (Mac Catalyst, the Mac idiom), which SwiftUI's page TabView puts below its pages there; it was rendered at 2x on macOS 27.2 by a Mac Catalyst probe (UIPageControl alone and SwiftUI's TabView of the page style, light and dark), which printed its frames and its pixels: dots 6 points across in indicator views 7 points square and 19 apart, in a control 23 points high with 14 points at either end, the chosen dot white and the others white at 45% (115 over black, 185 over 128). The look is inferred from it: UIKit draws the dots white in either appearance, for content behind them, so Breeze UI draws them in black in Aqua.
//
// Windows 11: WinUI's PipsPager of microsoft/microsoft-ui-xaml main at b1b6306 (controls/dev/PipsPager): PipsPager.xaml (a horizontal StackPanel of the previous button, a ScrollViewer around an ItemsRepeater of pip buttons and the next button), PipsPager_themeresources.xaml (each pip 12 by 24 DIPs, its EA3B glyph at 4 DIPs and the chosen one's at 6, ControlStrongFillColorDefault, TextFillColorSecondary under the pointer and pressed, the pointer growing a pip to 6 and a press shrinking it to 4, the buttons' styles), PipsPager.idl and PipsPager.properties.cpp (MaxVisiblePips 5, the buttons collapsed by default, Orientation horizontal), PipsPager.cpp (the ScrollViewer as wide as MaxVisiblePips pips, the chosen pip brought to its middle, the names, the arrows moving the focus) and Strings/en-us/Resources.resw (Pager, Page). The carousel is WinUI's FlipView: FlipView_themeresources.xaml (the background, the corners, the previous and next buttons 16 by 38 DIPs at the sides with a margin of 1, the in-app acrylic behind CaretLeftSolid8 and CaretRightSolid8 at 8 DIPs, not Tab stops), FlipView_Partial.cpp and FlipView_Partial.h (the buttons shown while a mouse moves over the view and hidden 3 seconds after it stops, never at the first or last page; a wheel's notch flipping one page, 200 ms apart; the move animated to an adjacent page only), Selector_Partial.cpp (the keys) and docs/api-specs/PipsPager/PipsPager.md (a PipsPager centered 10 DIPs below a FlipView).
//
// GNOME: libadwaita 1.10.0's AdwCarouselIndicatorDots (adw-carousel-indicator-dots.c: DOTS_RADIUS 3 and DOTS_RADIUS_SELECTED 4, DOTS_OPACITY 0.3 and DOTS_OPACITY_SELECTED 0.9 of the widget's color, DOTS_SPACING 7, DOTS_MARGIN 6, the radius and the opacity following the carousel's position) and AdwCarousel (adw-carousel.c: the keys, the mouse wheel one page per notch and 150 ms apart, the spring of damping ratio 1, mass 0.5 and stiffness 500; adw-swipe-tracker.c: a touchpad's 400 pixels to a page, a drag of 16 pixels before it starts, and where a swipe ends), with the dots below the carousel as libadwaita's demo lays them out (adw-demo-page-carousel.ui: a margin of 6 around them).

// pipsMacDot, pipsMacPitch, pipsMacHeight and pipsMacInset are UIPageControl's sizes in UIKit for Mac, measured on macOS 27.2 by tools/appkit-reference/pagecontrol.swift: its indicator views are 7 points square, 19 apart, centered in a control 23 points high, 14 points from either end to the first and the last, and each draws a circle 6 points across in its middle (at 2x, 10 opaque pixels and an antialiased one at either side); each dot is the middle of a cell 19 points wide, so the cells leave 8 at either end (a control of 5 dots is 111 points wide, of 9 dots 187). pipsMacOthers is the opacity of the dots not chosen, 45% (white reads 115 over black, 185 over 128 and 131 over the dark window's 30).
const (
	pipsMacDot    = 6
	pipsMacPitch  = 19
	pipsMacHeight = 23
	pipsMacInset  = 8
	pipsMacOthers = 0.45
)

// pipsFluentWidth and pipsFluentHeight are a WinUI pip's button, PipsPagerHorizontalOrientationButtonWidth and Height; pipsFluentNormal and pipsFluentSelected the font sizes of its EA3B glyph, PipsPagerNormalGlyphFontSize and PipsPagerSelectedGlyphFontSize, drawn as a dot as wide as the font size, the glyph filling its em (inferred); pipsFluentVisible is MaxVisiblePips's default.
const (
	pipsFluentWidth    = 12
	pipsFluentHeight   = 24
	pipsFluentNormal   = 4
	pipsFluentSelected = 6
	pipsFluentVisible  = 5
)

// dotsAdw* are AdwCarouselIndicatorDots's values (adw-carousel-indicator-dots.c): the dots' radii and opacities, the room between two dots, and the margin around them.
const (
	dotsAdwRadius          = 3
	dotsAdwRadiusSelected  = 4
	dotsAdwOpacity         = 0.3
	dotsAdwOpacitySelected = 0.9
	dotsAdwSpacing         = 7
	dotsAdwMargin          = 6
)

// dotsAdwPitch is the room of one dot, 2 * DOTS_RADIUS_SELECTED + DOTS_SPACING.
const dotsAdwPitch = 2*dotsAdwRadiusSelected + dotsAdwSpacing

// carouselFlipButton* are FlipView's PreviousButtonHorizontal and NextButtonHorizontal: 16 by 38 DIPs with a Margin of 1, at the view's sides and centered down it (FlipView_themeresources.xaml); carouselFlipButtonsShown is FLIP_VIEW_BUTTONS_SHOW_DURATION_MS, how long they stay after the mouse last moved over the view, and carouselFlipWheelPause FLIP_VIEW_DISTINCT_SCROLL_WHEEL_DELAY_MS, the pause between the wheel's deltas that lets a notch flip again (FlipView_Partial.h).
const (
	carouselFlipButtonW      = 16
	carouselFlipButtonH      = 38
	carouselFlipButtonsShown = 3000 * time.Millisecond
	carouselFlipWheelPause   = 200 * time.Millisecond
	carouselFlipPagerGap     = 10
)

// carouselAdwWheelPause is AdwCarousel's SCROLL_TIMEOUT_DURATION, how long after a notch of the mouse wheel moved a page another one waits; carouselSwipeQuiet is adw-swipe-tracker.c's EVENT_HISTORY_THRESHOLD_MS, the time over which a swipe's velocity is taken, which Breeze UI also takes as the end of a touchpad's swipe that sends no momentum; carouselAdwTouchpadPage is TOUCHPAD_BASE_DISTANCE_H, how far a touchpad scrolls for a page; carouselAdwDragStart is DRAG_THRESHOLD_DISTANCE; carouselAdwTouchpadFling and carouselAdwDragFling are VELOCITY_THRESHOLD_TOUCHPAD and VELOCITY_THRESHOLD_TOUCH, in pixels per millisecond: a swipe ending faster goes on to the next page, a slower one to the nearest; carouselSwipeBias is SCROLL_ORIENTATION_FACTOR, by how much a touchpad's scroll must be more sideways than up or down to swipe.
const (
	carouselAdwWheelPause     = 150 * time.Millisecond
	carouselSwipeQuiet        = 150 * time.Millisecond
	carouselAdwTouchpadPage   = 400
	carouselAdwDragStart      = 16
	carouselAdwTouchpadFling  = 0.6
	carouselAdwDragFling      = 0.3
	carouselSwipeBias         = 1.2
	carouselAdwDotsMargin     = 6
	carouselAnimationKey      = "breeze-carousel"
	carouselPipsScrollKey     = "breeze-pips-scroll"
	carouselAdwDotsMotionKey  = "breeze-dots"
	carouselMacFluentDuration = 250 * time.Millisecond
)

// carouselSwipeDecelTouchpad and carouselSwipeDecelTouch are adw-swipe-tracker.c's DECELERATION_TOUCHPAD and DECELERATION_TOUCH, carouselSwipeCurve its VELOCITY_CURVE_THRESHOLD and carouselSwipeParabola its DECELERATION_PARABOLA_MULTIPLIER, which project where a quick swipe glides to (get_end_progress); carouselSwipeEpsilon is its EPSILON, within which a swipe begins at a page rather than between two (get_bounds).
const (
	carouselSwipeDecelTouchpad = 0.997
	carouselSwipeDecelTouch    = 0.998
	carouselSwipeCurve         = 2
	carouselSwipeParabola      = 0.35
	carouselSwipeEpsilon       = 0.005
)

// carouselAdwSpring is AdwCarousel's scroll-params, adw_spring_params_new (1, 0.5, 500): a damping ratio of 1, a mass of 0.5 and a stiffness of 500, so a damping of 2 * sqrt(0.5 * 500); velocity is the velocity the spring starts with, in pages a second, which a swipe hands on (carouselSwipeEnd).
func carouselAdwSpring(velocity float32) ui.Animation {
	return ui.Animation{Spring: ui.SpringOptions{Mass: 0.5, Stiffness: 500, Damping: 2 * math.Sqrt(0.5*500), Velocity: velocity}}
}

// pageName returns the name of page i's pip: WinUI's PipsPagerPageText and the page's number, "Page 3" (PipsPager::OnElementPrepared).
func pageName(c *ui.Context, i int) string {
	return breezeStrings().Format(c.Locale(), "pipsPagerPage", map[string]any{"page": i + 1})
}

// PipsPager creates a row of count page indicators, of which *selected is the index of the one chosen, drawn as the platform draws its page dots: a click on one chooses it, and the arrows choose the one before or after, Left and Right mirrored right to left, Up and Down as WinUI's PipsPager takes them. It is one Tab stop, none without pages; assistive technology sees a radio group named Pager whose radio buttons are named by their pages, Page 1, Page 2, as WinUI names its pips.
//   - macOS: UIPageControl as UIKit for Mac draws it (inferred, as macOS has no page indicator): 6-point dots 19 apart in a control 23 points high, the chosen dot black and the others black at 45%, white in Dark Aqua, and AppKit's focus ring around the control (inferred).
//   - Windows: WinUI's PipsPager: pips 12 by 24 DIPs, a 4-DIP dot of ControlStrongFillColorDefault, the chosen one 6 DIPs, a pip under the pointer 6 DIPs and pressed 4 in TextFillColorSecondary; at most 5 pips show, the row scrolling to keep the chosen one in its middle, as MaxVisiblePips's default; the system focus visual around the chosen pip. WinUI's previous and next buttons are collapsed by default and are not drawn.
//   - GNOME: AdwCarouselIndicatorDots: dots 15 pixels apart, of radius 3 at 30% of the text's color, the chosen one of radius 4 at 90%, 6 pixels from the edges, the dots growing and shrinking as the choice moves; the focus ring inside its edge (inferred, as libadwaita's dots take no focus).
//
// It returns the radio group, whose Changed reports a page chosen. Build it under a view showing *selected's page, or use Carousel, which holds its pages with one.
func PipsPager(c *ui.Context, selected *int, count int) ui.Element {
	return pipsPager(c, selected, count, false, -1)
}

// pipsPager builds PipsPager; in a carousel (inCarousel) it is no Tab stop, the carousel taking the keys, and on GNOME its dots follow position, the carousel's position in pages, rather than a motion of their own.
func pipsPager(c *ui.Context, selected *int, count int, inCarousel bool, position float32) ui.Element {
	count = max(count, 0)
	if count > 0 {
		*selected = max(0, min(*selected, count-1))
	}
	sel := *selected
	root := ui.Row(c).Shrink(0).Role(ui.RoleRadioGroup).Label(say(c, "pipsPagerName"))
	// Without pages there is no pip to take the focus, as in WinUI.
	if !inCarousel && count > 0 {
		root.Focusable()
	}
	chosen := -1
	pip := func(i int) ui.Element {
		e := ui.Box(c).Shrink(0).Role(ui.RoleRadio).Label(pageName(c, i)).Checked(i == sel)
		if e.Clicked() {
			chosen = i
		}
		if i == sel {
			root.ActiveDescendant(e)
		}
		return e
	}
	focused := !inCarousel && root.FocusVisible()
	switch PlatformOf(c) {
	case Windows:
		pipsFluent(c, root, count, sel, pip, focused)
	case Linux:
		pos := position
		if pos < 0 {
			pos = root.AnimateTo(carouselAdwDotsMotionKey, float32(sel), carouselAdwSpring(0))
		}
		dotsAdw(c, root, count, pos, pip, focused, inCarousel)
	default:
		pipsMac(c, root, count, sel, pip, focused)
	}
	if !inCarousel && count > 0 {
		back, forth := ui.KeyLeft, ui.KeyRight
		if c.Direction() == ui.RightToLeft {
			back, forth = forth, back
		}
		switch {
		case root.Shortcut(0, back), root.Shortcut(0, ui.KeyUp):
			chosen = max(sel-1, 0)
		case root.Shortcut(0, forth), root.Shortcut(0, ui.KeyDown):
			chosen = min(sel+1, count-1)
		}
	}
	if chosen >= 0 && chosen != *selected {
		*selected = chosen
		root.SetChanged()
		c.Invalidate()
	}
	return root
}

// ---- macOS: UIPageControl in UIKit for Mac ----

// pipsMacInk returns the color of UIPageControl's dots: white, the chosen one opaque and the others at 45%, measured over black and gray; in Aqua black in its place (inferred, as UIKit draws them white in either appearance).
func pipsMacInk(t *ui.Theme, chosen bool) ui.Color {
	ink := pick(t.Dark, black, white)
	if chosen {
		return ink
	}
	return ink.Alpha(pipsMacOthers)
}

// pipsMac draws the pager as UIPageControl in UIKit for Mac (see PipsPager): cells 19 points wide holding the dots in their middles, 8 points from the control's ends, 23 points high, and AppKit's focus ring around the control, round at its ends (inferred).
func pipsMac(c *ui.Context, root ui.Element, n, sel int, pip func(int) ui.Element, focused bool) {
	t := c.Theme()
	root.AlignItems(ui.Stretch).FocusRing(false)
	// UIPageControl of no pages fits in no room (sizeThatFits gives 0 by 0).
	if n > 0 {
		root.Height(pipsMacHeight).PaddingX(pipsMacInset)
	}
	root.Children(func() {
		for i := range n {
			ink := pipsMacInk(t, i == sel)
			pip(i).Width(pipsMacPitch).Draw(func(p *ui.Painter, r ui.Rect) {
				p.Fill(ui.Rect{X: r.X + (r.W-pipsMacDot)/2, Y: r.Y + (r.H-pipsMacDot)/2, W: pipsMacDot, H: pipsMacDot}, ink, pipsMacDot/2)
			})
		}
	})
	root.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if focused {
			macFocusRing(p, r, r.H/2)
		}
	})
}

// ---- Windows: WinUI's PipsPager ----

// pipsFluentInk returns a pip's glyph color, PipsPagerSelectionIndicatorForeground (ControlStrongFillColorDefault, #72000000 and #8BFFFFFF; the chosen pip's the same) or, under the pointer and pressed, TextFillColorSecondary (#9E000000 and #C5FFFFFF), as translucent colors over the window.
func pipsFluentInk(t *ui.Theme, hovered, pressed bool) ui.Color {
	if hovered || pressed {
		return pick(t.Dark, black.Alpha(0x9e/255.0), white.Alpha(0xc5/255.0))
	}
	return pick(t.Dark, black.Alpha(0x72/255.0), white.Alpha(0x8b/255.0))
}

// pipsFluentDot returns the diameter of a pip's dot: 6 DIPs chosen or under the pointer, 4 otherwise and while pressed, as PipsPagerButtonBaseStyle's PointerOver and Pressed states set the glyph's FontSize.
func pipsFluentDot(chosen, hovered, pressed bool) float32 {
	switch {
	case pressed:
		return pipsFluentNormal
	case hovered, chosen:
		return pipsFluentSelected
	}
	return pipsFluentNormal
}

// pipsFluentScroll returns how far WinUI's ScrollViewer scrolls the pips of a pager of n pages to bring the chosen pip sel into its middle (ScrollToCenterOfViewport: HorizontalAlignmentRatio 0.5), held between its start and its end: the viewport is as wide as MaxVisiblePips pips (CalculateScrollViewerSize, the chosen pip as wide as the others).
func pipsFluentScroll(n, sel int) float32 {
	view := float32(pipsFluentWidth * min(n, pipsFluentVisible))
	content := float32(pipsFluentWidth * n)
	at := float32(sel)*pipsFluentWidth + pipsFluentWidth/2 - view/2
	return max(0, min(at, content-view))
}

// pipsFluent draws the pager as WinUI's PipsPager (see PipsPager): a viewport as wide as at most 5 pips around the row of 12 by 24 pips, scrolled to keep the chosen pip in its middle, gliding there as WinUI's StartBringIntoView animates it (the curve inferred), and the system focus visual on the chosen pip's edge (FocusVisualMargin 0, ControlCornerRadius 4).
func pipsFluent(c *ui.Context, root ui.Element, n, sel int, pip func(int) ui.Element, focused bool) {
	t := c.Theme()
	rtl := c.Direction() == ui.RightToLeft
	view := float32(pipsFluentWidth * min(n, pipsFluentVisible))
	off := root.AnimateTo(carouselPipsScrollKey, pipsFluentScroll(n, sel), ui.Smooth)
	shift := -off
	if rtl {
		shift = off
	}
	// A PipsPager of no pages is as small as its empty ItemsRepeater, 0 by 0, its buttons collapsed.
	h := float32(pipsFluentHeight)
	if n == 0 {
		h = 0
	}
	root.Height(h).FocusRing(false)
	root.Children(func() {
		ui.Row(c).Width(view).Height(h).Shrink(0).ClipX().Children(func() {
			ui.Row(c).Width(float32(pipsFluentWidth*n)).Height(h).Shrink(0).Translate(shift, 0).Children(func() {
				for i := range n {
					e := pip(i).Size(pipsFluentWidth, pipsFluentHeight)
					hovered, pressed := e.Hovered(), e.Pressed()
					d, ink := pipsFluentDot(i == sel, hovered, pressed), pipsFluentInk(t, hovered, pressed)
					e.Draw(func(p *ui.Painter, r ui.Rect) {
						p.Fill(ui.Rect{X: r.X + (r.W-d)/2, Y: r.Y + (r.H-d)/2, W: d, H: d}, ink, d/2)
					})
				}
			})
		})
	})
	root.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if !focused {
			return
		}
		x := r.X + float32(sel)*pipsFluentWidth - off
		if rtl {
			x = r.X + view - float32(sel+1)*pipsFluentWidth + off
		}
		fluentFocusRing(p, t, ui.Rect{X: x, Y: r.Y, W: pipsFluentWidth, H: pipsFluentHeight}, 0, 0, 0, 0, 4)
	})
}

// ---- GNOME: AdwCarouselIndicatorDots ----

// dotsAdwWidth is AdwCarouselIndicatorDots's natural width for n pages: a dot's room for each and DOTS_MARGIN at either end.
func dotsAdwWidth(n int) float32 { return float32(dotsAdwPitch*n + 2*dotsAdwMargin) }

// dotsAdwHeight is its height: the chosen dot and DOTS_MARGIN above and below.
const dotsAdwHeight = 2*dotsAdwRadiusSelected + 2*dotsAdwMargin

// dotsAdwCenter returns the middle of dot slot s, from the left, of a row of n dots in a widget w pixels wide, as snapshot_dots places them: the dots' length taken as n rooms less one spacing, centered in the widget made a pixel narrower when the difference with n rooms is even, and each dot in the middle of its room; so the dots stand half a spacing toward the end of the widget's middle, as libadwaita 1.10.0 draws them.
func dotsAdwCenter(w float32, n, s int) float32 {
	length := float32(dotsAdwPitch*n - dotsAdwSpacing)
	full := float32(dotsAdwPitch * n)
	wl := int(w)
	if (wl-int(full))%2 == 0 {
		wl--
	}
	return (float32(wl)-length)/2 + float32(dotsAdwPitch)/2 + float32(dotsAdwPitch*s)
}

// dotsAdwShare returns how far dot i of n is chosen, from 0 to 1, at position, in pages, as snapshot_dots shares the progress between the dots on either side of a position between two pages.
func dotsAdwShare(n, i int, position float32) float32 {
	left := float32(1)
	for j := 0; j <= i && j < n; j++ {
		share := max(0, min(float32(j+1)-position, left))
		left -= share
		if j == i {
			return share
		}
	}
	return 0
}

// dotsAdw draws the pager as AdwCarouselIndicatorDots (see PipsPager) at position, in pages: each dot's radius from 3 to 4 and its opacity from 30% to 90% of the widget's color as it is chosen, mirrored right to left. Each page's cell, a dot's room 20 pixels high around its dot, takes a click, which libadwaita's dots do not (inferred); alone, the focus ring goes inside the widget's edge, round at its ends (inferred). In a carousel the dots are as wide as it, as GtkBox gives them its width.
func dotsAdw(c *ui.Context, root ui.Element, n int, position float32, pip func(int) ui.Element, focused, inCarousel bool) {
	t := c.Theme()
	rtl := c.Direction() == ui.RightToLeft
	root.Height(dotsAdwHeight).FocusRing(false)
	w := dotsAdwWidth(n)
	if inCarousel {
		root.AlignSelf(ui.Stretch)
		if b := root.LayoutBounds(); b.W > 0 {
			w = b.W
		}
	} else {
		root.Width(w)
	}
	slot := func(i int) int {
		if rtl {
			return n - 1 - i
		}
		return i
	}
	root.Children(func() {
		for i := range n {
			x := dotsAdwCenter(w, n, slot(i)) - float32(dotsAdwPitch)/2
			pip(i).Absolute().Left(x).Top(0).Size(dotsAdwPitch, dotsAdwHeight)
		}
	})
	root.Draw(func(p *ui.Painter, r ui.Rect) {
		if n < 2 {
			return
		}
		for i := range n {
			share := dotsAdwShare(n, i, position)
			radius := dotsAdwRadius + (dotsAdwRadiusSelected-dotsAdwRadius)*share
			opacity := dotsAdwOpacity + (dotsAdwOpacitySelected-dotsAdwOpacity)*share
			cx, cy := r.X+dotsAdwCenter(r.W, n, slot(i)), r.Y+float32(int(r.H)/2)
			p.Fill(ui.Rect{X: cx - radius, Y: cy - radius, W: 2 * radius, H: 2 * radius}, adwInkLayer(t, opacity), radius)
		}
	})
	root.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if focused {
			p.Stroke(r, t.Focus, r.H/2, 2)
		}
	})
}

// ---- the carousel ----

// carouselSample is a move of a swipe, as adw-swipe-tracker.c records it: how far, in DIPs, and when.
type carouselSample struct {
	delta float32
	at    time.Time
}

// carouselState is what a carousel keeps between frames.
type carouselState struct {
	// swiping is whether a touchpad's swipe or the pointer's drag moves the pages; pos is where they show then, in pages, between lo and hi, the pages it may reach (carouselSwipeBounds) from from, where the pages showed as it began, initial being the page nearest to from; dist is how far a page is, in DIPs, and fling the velocity, in DIPs a millisecond, past which the swipe goes on toward the next page, or -1 for a move to initial as it is (a notch of the wheel); touchpad tells a touchpad's swipe from a drag, which decelerate differently; samples are its last moves.
	swiping     bool
	pos         float32
	initial     int
	from        float32
	lo, hi      float32
	dist, fling float32
	touchpad    bool
	samples     []carouselSample
	lastSwipe   time.Time
	// endSwipe asks the next frame to end the swipe, or to go to initial; ended is the velocity a drag ended with, in DIPs a millisecond, when hasEnded; velocity is the velocity a swipe ended with, which the spring taking over starts with (see carouselSwipeEnd).
	endSwipe, hasEnded bool
	ended, velocity    float32
	// momentum is whether the touchpad's momentum after a swipe is the carousel's, to take and drop.
	momentum bool
	// dragged and draggedY are how far the pointer pressing the view went, sideways and down, before a drag began; dragging whether it drags the pages, and dragRejected whether it went up or down rather than sideways, which leaves the pages until it lets go.
	dragged, draggedY      float32
	dragging, dragRejected bool
	// wheelAt is the last delta of a wheel on Windows, and on GNOME when the wheel may move a page again; wheelDir is the direction of the last flip on Windows.
	wheelAt  time.Time
	wheelDir int
	// pointerAt is when the mouse last moved over the view.
	pointerAt time.Time
	// target is the page the pages last went to, and shown where they showed in the last frame, in pages; w and h the view's size in the last frame.
	target   int
	shown    float32
	w, h     float32
	disabled bool
}

// Carousel creates a view showing one of count pages, the one *selected names, which page builds, with the platform's page dots below it: page(i) builds page i's contents, and is called for the page shown and, while the view moves between pages, for those coming into view, each in a container of its own kept by its index. The arrows, Page Up and Page Down move to the page before or after, Home and End to the first and the last, as the platform's do; a click on a dot chooses its page. It is one Tab stop, none without pages, the view, named by the page it shows ("Page 2 of 4"), with the platform's focus ring around it; the dots below are no Tab stop.
//   - macOS: the pages as NSPageController's horizontal strip moves them, sliding over 0.25 seconds, a touchpad's swipe moving them with the fingers (inferred, as macOS has no carousel with page dots); UIPageControl's dots (PipsPager) centered right below (inferred).
//   - Windows: WinUI's FlipView: SolidBackgroundFillColorBase with corners of 4; while a mouse moves over it, and for 3 seconds after, its previous and next buttons at its sides, 16 by 38 DIPs, the in-app acrylic (#f9f9f9, #2c2c2c) with corners of 4 behind a filled caret in ControlStrongFillColorDefault (TextFillColorSecondary under the pointer and pressed, a press shrinking it to 0.875), the previous button hidden on the first page and the next on the last; a notch of the wheel, or a touchpad's swipe, flipping one page, again after a pause of 200 ms or the other way; the move animated to an adjacent page and immediate to another; Up and Down taken as Left and Right; the system focus visual around the view; and a PipsPager centered 10 DIPs below it, as WinUI's PipsPager documentation pairs them.
//   - GNOME: AdwCarousel: no background; a touchpad's swipe (400 pixels to a page) or a drag of the pointer moving the pages, which go on to the next page past a velocity or back to the nearest; a notch of the wheel moving a page, 150 ms apart; libadwaita's spring between pages; AdwCarouselIndicatorDots below it with a margin of 6, following the pages as they move; the focus ring inside its edge (inferred).
//
// It returns the carousel, whose Changed reports a page chosen. Give it a width, and a height if its pages should not set it: the pages fill the view.
func Carousel(c *ui.Context, selected *int, count int, page func(i int)) ui.Element {
	pl := PlatformOf(c)
	t := c.Theme()
	count = max(count, 0)
	if count > 0 {
		*selected = max(0, min(*selected, count-1))
	}
	rtl := c.Direction() == ui.RightToLeft
	now := c.Now()
	root := ui.Column(c).AlignItems(ui.Stretch)
	var area ui.Element
	var st *carouselState
	changed := false
	choose := func(i int) {
		i = max(0, min(i, count-1))
		if count > 0 && i != *selected {
			*selected = i
			changed = true
		}
	}
	root.Children(func() {
		area = ui.Column(c).Grow(1).AlignItems(ui.Stretch).FocusRing(false).Clip().Role(ui.RoleGroup)
		// Without pages there is nothing to take the focus, as a FlipView without items and an AdwCarousel without children have none.
		if count > 0 {
			area.Focusable()
		}
		st = ui.Local(area, "breeze-carousel-state", func() carouselState {
			return carouselState{target: *selected, shown: float32(*selected)}
		})
		if b := area.LayoutBounds(); b.W != st.w || b.H != st.h {
			st.w, st.h = b.W, b.H
			c.Invalidate()
		}
		carouselInput(c, pl, area, st, count, rtl)
		// The keys: Left and Right mirrored right to left, Page Up and Page Down, Home and End, as AdwCarousel's bindings and WinUI's Selector::HandleNavigationKey take them, and Up and Down on Windows, where a horizontal FlipView takes them too. FlipView moves from its SelectedIndex, the page the pages go to; AdwCarousel's navigate_to_direction from the page nearest to where they show (round (self->position)), so that a key pressed again while they move toward the next page moves to it again, and none at the first or the last.
		back, forth := ui.KeyLeft, ui.KeyRight
		if rtl {
			back, forth = forth, back
		}
		base := *selected
		if pl == Linux {
			base = carouselNearest(st.shown, count)
		}
		switch {
		case area.Shortcut(0, back), area.Shortcut(0, ui.KeyPageUp), pl == Windows && area.Shortcut(0, ui.KeyUp):
			if base > 0 {
				choose(base - 1)
			}
		case area.Shortcut(0, forth), area.Shortcut(0, ui.KeyPageDown), pl == Windows && area.Shortcut(0, ui.KeyDown):
			if base < count-1 {
				choose(base + 1)
			}
		case area.Shortcut(0, ui.KeyHome):
			choose(0)
		case area.Shortcut(0, ui.KeyEnd):
			choose(count - 1)
		}
		// A touchpad that sends no momentum ends a swipe by stopping: after 150 ms without a move.
		if st.swiping && !st.dragging {
			if d := now.Sub(st.lastSwipe); d >= carouselSwipeQuiet {
				st.swiping, st.endSwipe = false, true
			} else {
				c.After(carouselSwipeQuiet - d)
			}
		}
		if st.endSwipe {
			st.endSwipe = false
			choose(carouselSwipeEnd(st, count))
		}
		pos := carouselPosition(pl, area, st, *selected)
		st.shown = pos
		area.Children(func() {
			carouselStrip(c, st, count, pos, rtl, page)
			if pl == Windows {
				carouselFlipButtons(c, t, st, count, *selected, now, choose)
			}
		})
		if pl == Windows {
			// FlipViewBackground, SolidBackgroundFillColorBase, with ControlCornerRadius.
			area.Background(t.Background).Radius(4)
		}
		if count > 0 {
			area.Label(breezeStrings().Format(c.Locale(), "carouselPage", map[string]any{"page": *selected + 1, "count": count}))
		}
		before := *selected
		pager := pipsPager(c, selected, count, true, pos)
		if *selected != before {
			changed = true
		}
		switch pl {
		case Windows:
			pager.AlignSelf(ui.Center).Margin(carouselFlipPagerGap, 0, 0, 0)
		case Linux:
			pager.Margin(carouselAdwDotsMargin)
		default:
			pager.AlignSelf(ui.Center)
		}
	})
	focused := area.FocusVisible()
	area.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if !focused {
			return
		}
		switch pl {
		case Windows:
			fluentFocusRing(p, t, r, 0, 0, 0, 0, 0)
		case Linux:
			p.Stroke(r, t.Focus, 0, 2)
		}
	})
	root.Draw(func(p *ui.Painter, r ui.Rect) {
		if root.Valid() {
			st.disabled = root.IsDisabled()
		}
	}).DrawOver(func(p *ui.Painter, r ui.Rect) {
		if focused && pl == Mac {
			macFocusRing(p, ui.Rect{X: r.X, Y: r.Y, W: r.W, H: r.H - pipsMacHeight}, 0)
		}
	})
	if changed {
		root.SetChanged()
		c.Invalidate()
	}
	return root
}

// carouselPosition returns where the pages show this frame, in pages: under the fingers or the pointer while they swipe, else on the way to page sel, as the platform moves there. GNOME's spring goes on from a swipe's velocity; Windows animates the move only to an adjacent page, as FlipView::OnSelectedIndexChanged does, and jumps to another; macOS slides over 0.25 seconds (inferred: NSAnimationContext's default duration, which NSPageController's animator takes).
func carouselPosition(pl Platform, area ui.Element, st *carouselState, sel int) float32 {
	if st.swiping {
		return area.Jump(carouselAnimationKey, st.pos)
	}
	target := float32(sel)
	last := st.target
	st.target = sel
	switch pl {
	case Linux:
		v := st.velocity
		st.velocity = 0
		return area.AnimateTo(carouselAnimationKey, target, carouselAdwSpring(v))
	case Windows:
		if last-sel > 1 || sel-last > 1 {
			return area.Jump(carouselAnimationKey, target)
		}
		// ControlNormalAnimationDuration along ControlFastOutSlowInKeySpline, WinUI's own motion of its controls, stands for DirectManipulation's view change, whose curve its sources do not give (inferred).
		return area.AnimateTo(carouselAnimationKey, target, ui.Animation{Duration: carouselMacFluentDuration, Ease: ui.CubicBezier(0, 0, 0, 1)})
	}
	return area.AnimateTo(carouselAnimationKey, target, ui.Animation{Duration: carouselMacFluentDuration, Ease: ui.EaseInOut})
}

// carouselRange returns the pages whose boxes cross the view at pos, the first and the last: the one page at rest, or the two pos is between. Without the view's width of the last frame, which moving pages need, it is the nearest page alone.
func carouselRange(count int, pos, w float32) (lo, hi int) {
	if count == 0 {
		return 0, -1
	}
	pos = max(0, min(pos, float32(count-1)))
	lo, hi = int(math.Floor(float64(pos))), int(math.Ceil(float64(pos)))
	if w <= 0 {
		lo = int(math.Round(float64(pos)))
		hi = lo
	}
	return lo, hi
}

// carouselStrip builds the pages shown at pos (carouselRange), each in a container keyed by its index, side by side in a strip moved by how far pos is past the first of them, mirrored right to left. At rest the one page fills the view.
func carouselStrip(c *ui.Context, st *carouselState, count int, pos float32, rtl bool, page func(int)) {
	lo, hi := carouselRange(count, pos, st.w)
	moving := hi > lo
	shift := float32(0)
	if moving {
		shift = -(max(0, min(pos, float32(count-1))) - float32(lo)) * st.w
		if rtl {
			shift = -shift
		}
	}
	strip := ui.Row(c).Grow(1).AlignItems(ui.Stretch).Shrink(0).Translate(shift, 0)
	if moving {
		strip.Width(st.w * float32(hi-lo+1))
	}
	strip.Children(func() {
		for i := lo; i <= hi; i++ {
			e := ui.Column(c).Key(i).Shrink(0).AlignItems(ui.Stretch)
			if moving {
				e.Width(st.w)
			} else {
				e.Grow(1)
			}
			e.Children(func() { page(i) })
		}
	})
}

// carouselFlipButtons builds FlipView's previous and next buttons over the view's sides (see Carousel), built in every frame and hidden when not shown: while the mouse moved over the view in the last 3 seconds or is over a button, the previous one not on the first page and the next one not on the last. A click moves a page and keeps them for another 3 seconds, as OnPreviousButtonPartClick resets the fade-out timer.
func carouselFlipButtons(c *ui.Context, t *ui.Theme, st *carouselState, count, sel int, now time.Time, choose func(int)) {
	top := max(0, (st.h-carouselFlipButtonH)/2)
	shown := false
	if d := now.Sub(st.pointerAt); !st.pointerAt.IsZero() && d < carouselFlipButtonsShown {
		shown = true
		c.After(carouselFlipButtonsShown - d)
	}
	prev := ui.Box(c).Absolute().Start(1).Top(top).Size(carouselFlipButtonW, carouselFlipButtonH).Role(ui.RoleButton).Label(say(c, "carouselPrevious"))
	next := ui.Box(c).Absolute().End(1).Top(top).Size(carouselFlipButtonW, carouselFlipButtonH).Role(ui.RoleButton).Label(say(c, "carouselNext"))
	over := prev.Hovered() || next.Hovered()
	rtl := c.Direction() == ui.RightToLeft
	for k, b := range []ui.Element{prev, next} {
		back := k == 0
		b.Hidden(!(shown || over) || (back && sel <= 0) || (!back && sel >= count-1))
		hovered, pressed := b.Hovered(), b.Pressed()
		ink := pipsFluentInk(t, hovered, pressed)
		scale := float32(1)
		if pressed {
			scale = 0.875 // FlipViewButtonScalePressed
		}
		// The glyphs are MirroredWhenRightToLeft: the previous button's caret points to the start.
		left := back != rtl
		b.Draw(func(p *ui.Painter, r ui.Rect) {
			// FlipViewNextPreviousButtonBackground, AcrylicInAppFillColorDefaultBrush's fallback color, with ControlCornerRadius and FlipViewButtonBorderThemeThickness 0.
			p.Fill(r, pick(t.Dark, ui.Hex("#f9f9f9"), ui.Hex("#2c2c2c")), 4)
			fluentFlipCaret(p, r, left, scale, ink)
		})
		if b.Clicked() {
			if back {
				choose(sel - 1)
			} else {
				choose(sel + 1)
			}
			st.pointerAt = now
		}
	}
}

// fluentFlipCaret draws CaretLeftSolid8 (EDD9) or CaretRightSolid8 (EDDA) at FlipViewButtonFontSize 8 in the middle of r: a filled caret 6 DIPs wide and 8 high pointing left or right, CaretUpSolid8's measured 8 by 6 turned on its side (inferred), scaled by scale.
func fluentFlipCaret(p *ui.Painter, r ui.Rect, left bool, scale float32, ink ui.Color) {
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	w, h := 3*scale, 4*scale
	var tri ui.Path
	if left {
		tri.MoveTo(cx+w, cy-h).LineTo(cx-w, cy).LineTo(cx+w, cy+h).Close()
	} else {
		tri.MoveTo(cx-w, cy-h).LineTo(cx+w, cy).LineTo(cx-w, cy+h).Close()
	}
	p.FillPath(&tri, ink)
}

// carouselInput has the view take the scrolling and the pointer's moves over it as they come. On Windows a scroll flips a page as FlipView::OnPointerWheelChanged does (Ctrl left to zooming), and a move of the mouse shows the buttons. On GNOME a notch of a mouse's wheel moves a page, up or down too, as AdwCarousel's scroll_cb does, and a touchpad's swipe moves the pages with the fingers, as on macOS (inferred there, where a mouse's wheel moves none, NSPageController following only a touchpad); the swipe starts only more sideways than up or down. A drag of the pointer moves them on GNOME, as AdwSwipeTracker's allow-mouse-drag lets AdwCarousel, once it went 16 pixels.
func carouselInput(c *ui.Context, pl Platform, area ui.Element, st *carouselState, count int, rtl bool) {
	area.HandleInput(func(ev ui.InputEvent) bool {
		if st.disabled {
			return false
		}
		if ev.Kind == ui.InputPointerMove && pl == Windows {
			st.pointerAt = ev.Time
			c.Invalidate()
			return false
		}
		if ev.Kind != ui.InputScroll || count == 0 {
			return false
		}
		switch {
		case pl == Windows:
			return carouselFlip(st, ev, count)
		case ev.Momentum:
			// The fingers left the touchpad: the swipe ends, and the momentum after it is dropped.
			if st.swiping {
				st.swiping, st.endSwipe = false, true
				return true
			}
			return st.momentum
		case ev.Precise:
			d := ev.DX
			if !st.swiping && math.Abs(float64(d)) < math.Abs(float64(ev.DY))*carouselSwipeBias {
				st.momentum = false
				return false
			}
			if rtl {
				d = -d
			}
			dist := st.w
			if pl == Linux {
				dist = carouselAdwTouchpadPage
			}
			if dist <= 0 {
				return false
			}
			if !st.swiping {
				carouselSwipeBegin(st, true, dist, count)
			}
			carouselSwipeMove(st, d, ev.Time, count)
			return true
		case pl == Linux:
			if ev.Time.Before(st.wheelAt) {
				return false
			}
			// scroll_cb moves from the page nearest to where the pages show (get_page_at_position).
			from := carouselNearest(st.shown, count)
			dir := 0
			switch {
			case ev.DY != 0:
				dir = int(math.Copysign(1, float64(ev.DY)))
			case ev.DX != 0:
				dir = int(math.Copysign(1, float64(ev.DX)))
			default:
				return false
			}
			st.initial, st.fling, st.endSwipe = max(0, min(from+dir, count-1)), -1, true
			st.wheelAt = ev.Time.Add(carouselAdwWheelPause)
			return true
		}
		return false
	})
	if pl != Linux || st.disabled || count == 0 || st.w <= 0 {
		return
	}
	// Dragged reports the pointer pressing the view wherever it goes, where Pressed stops when it leaves the view: GTK's drag gesture keeps the pointer it took.
	dx, dy, held := area.Dragged()
	if held && !st.dragRejected {
		if rtl {
			dx = -dx
		}
		if !st.dragging {
			// drag_update_cb: a drag going more up or down than sideways at its first move, or as it passes 16 DIPs from where it began, is not the carousel's; one going sideways starts the swipe there, with its last move, the way before it being lost.
			first := st.dragged == 0 && st.draggedY == 0
			st.dragged += dx
			st.draggedY += dy
			vertical := math.Abs(float64(st.draggedY)) > math.Abs(float64(st.dragged))
			switch {
			case first && vertical:
				st.dragRejected = true
			case math.Hypot(float64(st.dragged), float64(st.draggedY)) >= carouselAdwDragStart:
				if vertical {
					st.dragRejected = true
					break
				}
				st.dragging = true
				carouselSwipeBegin(st, false, st.w, count)
				carouselSwipeMove(st, -dx, c.Now(), count)
			}
		} else {
			carouselSwipeMove(st, -dx, c.Now(), count)
		}
	}
	vx, _, ended := area.DragEnded()
	if ended || (!held && (st.dragging || st.dragged != 0 || st.draggedY != 0 || st.dragRejected)) {
		if st.dragging {
			if rtl {
				vx = -vx
			}
			st.swiping, st.endSwipe = false, true
			st.hasEnded, st.ended = ended, -vx/1000
		}
		st.dragging, st.dragged, st.draggedY, st.dragRejected = false, 0, 0, false
	}
}

// carouselFlip flips a page with a scroll on Windows, as FlipView::OnPointerWheelChanged does: a delta down or to the right moves to the next page, up or to the left to the one before, when it goes the other way than the last flip or comes after a pause of 200 ms, every delta restarting the pause; a delta that cannot flip is taken, so that the view around does not scroll either, but one at the first or the last page is left to it.
func carouselFlip(st *carouselState, ev ui.InputEvent, count int) bool {
	if ev.Mods&ui.Ctrl != 0 {
		return false
	}
	d := ev.DY
	if d == 0 {
		d = ev.DX
	}
	if d == 0 {
		return false
	}
	dir := int(math.Copysign(1, float64(d)))
	flip := dir != st.wheelDir || ev.Time.Sub(st.wheelAt) > carouselFlipWheelPause
	st.wheelAt = ev.Time
	if !flip {
		return true
	}
	to := st.target + dir
	if st.endSwipe {
		to = st.initial + dir
	}
	if to < 0 || to >= count {
		return false
	}
	st.wheelDir = dir
	st.initial, st.fling, st.endSwipe = to, -1, true
	return true
}

// carouselSwipeBegin starts a swipe from where the pages show, as begin_swipe_cb stops AdwCarousel's spring there: a touchpad's (touchpad) or a drag's, dist being the DIPs to a page.
func carouselSwipeBegin(st *carouselState, touchpad bool, dist float32, count int) {
	st.swiping, st.momentum, st.hasEnded = true, true, false
	st.from = max(0, min(st.shown, float32(count-1)))
	st.pos = st.from
	st.initial = carouselNearest(st.from, count)
	st.lo, st.hi = carouselSwipeBounds(st.from, count)
	st.dist, st.touchpad = dist, touchpad
	st.fling = carouselAdwDragFling
	if touchpad {
		st.fling = carouselAdwTouchpadFling
	}
	st.samples = st.samples[:0]
}

// carouselNearest returns the page nearest to pos, in pages, the first of two as near, as find_closest_point and get_closest_child_at take it.
func carouselNearest(pos float32, count int) int {
	return max(0, min(int(math.Ceil(float64(pos)-0.5)), count-1))
}

// carouselSwipeBounds returns the pages a swipe beginning at from, in pages, may reach without long swipes, as get_bounds gives them: from a page, the pages on either side of it; from between two pages, the page before the first and the page after the second; within the pages.
func carouselSwipeBounds(from float32, count int) (lo, hi float32) {
	prev, next := carouselNearest(from, count), carouselNearest(from, count)
	if math.Abs(float64(from)-float64(prev)) >= carouselSwipeEpsilon {
		prev, next = int(math.Floor(float64(from))), int(math.Ceil(float64(from)))
	}
	return float32(max(prev-1, 0)), float32(min(next+1, count-1))
}

// carouselSwipeMove moves a swipe by d DIPs toward the later pages at time at, within its bounds, as gesture_update does without overshoot, and records the move, keeping the last 150 ms as adw-swipe-tracker.c's history does.
func carouselSwipeMove(st *carouselState, d float32, at time.Time, count int) {
	st.pos = max(st.lo, min(st.pos+d/st.dist, st.hi))
	st.lastSwipe = at
	st.samples = append(st.samples, carouselSample{d, at})
	for len(st.samples) > 0 && at.Sub(st.samples[0].at) > carouselSwipeQuiet {
		st.samples = st.samples[1:]
	}
}

// carouselSwipeVelocity returns a swipe's velocity in DIPs a millisecond, as calculate_velocity takes it: the moves after the first of the history over the time it spans; a drag's is the pointer's as it let go. At a bound, a velocity beyond it is none, as AdwCarousel allows no overshoot.
func carouselSwipeVelocity(st *carouselState) float32 {
	v := st.ended
	if !st.hasEnded {
		v = 0
		if len(st.samples) >= 2 {
			if ms := float32(st.samples[len(st.samples)-1].at.Sub(st.samples[0].at)) / float32(time.Millisecond); ms > 0 {
				total := float32(0)
				for _, s := range st.samples[1:] {
					total += s.delta
				}
				v = total / ms
			}
		}
	}
	if (st.pos <= st.lo && v < 0) || (st.pos >= st.hi && v > 0) {
		return 0
	}
	return v
}

// carouselSwipeEnd returns the page a swipe that just ended goes to, as get_end_progress chooses it without long swipes: slower than its fling velocity, the page nearest to where it is; faster, where the pages would glide to, decelerating as a touchpad's swipe or a drag does (DECELERATION_TOUCHPAD 0.997, DECELERATION_TOUCH 0.998, past VELOCITY_CURVE_THRESHOLD 2 along DECELERATION_PARABOLA_MULTIPLIER's parabola), held within its bounds, then the page beyond it the swipe's way when the page before it is the one the swipe began at (find_point_for_projection), else the page nearest to it. A move of the wheel (fling -1) goes to its page as it is. The velocity, in DIPs a millisecond, is handed on to the spring as AdwCarousel's end_swipe_cb hands it on, the spring taking it as pages a second.
func carouselSwipeEnd(st *carouselState, count int) int {
	if st.fling < 0 {
		return st.initial
	}
	v := carouselSwipeVelocity(st)
	st.velocity = v
	speed := math.Abs(float64(v))
	if speed < float64(st.fling) {
		return max(int(st.lo), min(carouselNearest(st.pos, count), int(st.hi)))
	}
	pos := max(float64(st.lo), min(float64(st.pos)+math.Copysign(carouselSwipeGlide(speed, st.touchpad), float64(v)), float64(st.hi)))
	prev, next := int(math.Floor(pos)), int(math.Ceil(pos))
	switch {
	case v > 0 && prev == st.initial:
		return next
	case v < 0 && next == st.initial:
		return prev
	}
	return carouselNearest(float32(pos), count)
}

// carouselSwipeGlide returns how far, in pages, a swipe ending at speed, in DIPs a millisecond, glides on as get_end_progress projects it: a touchpad's (touchpad) or a drag's deceleration taken as a slope, linear up to VELOCITY_CURVE_THRESHOLD and along a parabola past it.
func carouselSwipeGlide(speed float64, touchpad bool) float64 {
	decel := carouselSwipeDecelTouch
	if touchpad {
		decel = carouselSwipeDecelTouchpad
	}
	slope := decel / (1 - decel) / 1000
	if speed <= carouselSwipeCurve {
		return speed * slope
	}
	c := slope / 2 / carouselSwipeParabola
	x := speed - carouselSwipeCurve + c
	return carouselSwipeParabola*x*x - carouselSwipeParabola*c*c + slope*carouselSwipeCurve
}
