package breeze

import (
	"math"

	"github.com/egoist/mygo/ui"
)

// ---- status: inline message bars, determinate progress rings and empty-state pages ----
//
// None of these is a widget of MyGo's; they are built from its primitives (ui.Row, ui.Column, ui.ButtonBase, ui.Icon, ui.RichText and painting) and drawn as each platform draws its own.
//
// macOS values are measured from AppKit and SwiftUI rendered at 2x into an sRGB bitmap in a window answering hasKeyAppearance, on macOS 27.2, by the status set of tools/appkit-reference/reference.swift: NSProgressIndicator of style spinning with isIndeterminate false, regular and small; SwiftUI's ContentUnavailableView in an NSHostingView; and, for the message bar, which AppKit does not have, a row of a SwiftUI grouped Form holding a Label with an SF Symbol and a borderless close button, the closest system look (System Settings shows its inline notices so). Windows values are WinUI 3's (microsoft/microsoft-ui-xaml main): InfoBar_themeresources.xaml and InfoBar.xaml, ProgressRing_themeresources.xaml, ProgressRing.xaml, ProgressRing.cpp and the Lottie of AnimatedVisuals/ProgressRingDeterminate.cpp, and Common_themeresources_any.xaml; WinUI has no empty-state control, so StatusPage follows WinUI's type ramp (SubtitleTextBlockStyle, BodyTextBlockStyle) and is inferred. GNOME values are libadwaita 1.10.0's: _toolbars.scss (banner), adw-banner.c and adw-banner.ui, _misc.scss (statuspage), adw-status-page.ui, _labels.scss (.title-1, .heading, .body) and _colors.scss (--dim-opacity); libadwaita has no determinate ring, so ProgressRing follows AdwSpinner's track and stroke width (adw-spinner-paintable.c) with GtkProgressBar's accent, inferred.

// InfoBar creates an inline message bar of a severity, kind, with a title and a message, as the platform shows a message inside a window: WinUI's InfoBar on Windows, libadwaita's AdwBanner on GNOME, and on macOS, which has no message bar, a section of a grouped form with a colored SF Symbol, as System Settings shows its inline notices (inferred). kind is a ui.Toast's Type, so that the same word colors a toast and a bar: "success", "warning", "error" ("danger" and "critical" too), and any other, as "" or "info", is informational; libadwaita's banners have no severity, so on GNOME every kind draws the same banner. A nil open makes a bar that cannot be closed; otherwise the bar shows while *open and its close button sets *open to false: WinUI's close button, a borderless xmark button on macOS, and on GNOME the banner's button labelled Dismiss. The bar is a status for assistive technology, named by its title and described by its message; the close button is named "Close" on macOS and Windows and "Dismiss" on GNOME. It stretches across a column, as the platforms' bars stretch.
func InfoBar(c *ui.Context, kind, title, message string, open *bool) ui.Element {
	if open != nil && !*open {
		return ui.Element{}
	}
	bar := ui.Row(c).AlignSelf(ui.Stretch).Role(ui.RoleStatus).Label(title)
	if message != "" {
		bar.Description(message)
	}
	var closer func(name string) ui.Element
	if open != nil {
		closer = func(name string) ui.Element {
			b := ui.ButtonBase(c).Label(name).Shrink(0)
			if b.Clicked() {
				*open = false
			}
			return b
		}
	}
	severity := fluentToastSeverity(kind)
	switch PlatformOf(c) {
	case Windows:
		fluentInfoBar(c, bar, severity, title, message, closer)
	case Linux:
		adwBanner(c, bar, title, message, closer)
	default:
		macInfoBar(c, bar, severity, title, message, closer)
	}
	return bar
}

// lineHeight returns the height of a line of text of size and weight, as a Text of the theme's font lays it out.
func lineHeight(c *ui.Context, size float32, weight int) float32 {
	_, h := c.MeasureText(0, ui.Span{Text: "Ag", Size: size, Weight: weight})
	return h
}

// macInfoBarColor returns the color of the SF Symbol of a severity, measured from SwiftUI on macOS 27.2: controlAccentColor (#007aff in both appearances, the theme's Accent), systemGreen (the theme's Success), systemYellow (#ffcc00, #ffd600 in Dark Aqua, as an NSLevelIndicator's warning fill) and systemRed (#ff383c, #ff4245: the theme's Danger).
func macInfoBarColor(t *ui.Theme, severity int) ui.Color {
	switch severity {
	case fluentToastSuccess:
		return t.Success
	case fluentToastWarning:
		return macMeterFill(t, 1)
	case fluentToastError:
		return t.Danger
	}
	return t.Accent
}

// macInfoBar draws the message bar on macOS, which has none: a row of a SwiftUI grouped Form, as measured on macOS 27.2: the section's fill (macFormSectionFill) with corners of 12, the row padded 10 points above and below and 11.5 at the start; a 13-point SF Symbol of the severity (info.circle.fill, checkmark.circle.fill, exclamationmark.triangle.fill or xmark.octagon.fill) in its color, the mark inside cut out to the fill, centered on the title's first line; 11.5 points after it the title in 13-point bold labelColor and the message in secondaryLabelColor, both over the fill (38 and 124 over 247, 222 and 157 over 37), 2 points apart (inferred); and at the end a borderless button showing xmark, 10.5 points across in a 16-point box, its ink 12 points from the end, in secondaryLabelColor, labelColor while pressed (inferred). That macOS shows such a row as a message bar is inferred.
func macInfoBar(c *ui.Context, bar ui.Element, severity int, title, message string, closer func(string) ui.Element) {
	t := c.Theme()
	fill := macFormSectionFill(t)
	secondary := tint(fill, t.Dark, pick(t.Dark, float32(0.498), 0.549))
	line := lineHeight(c, t.FontSize, 700)
	bar.AlignItems(ui.Start).Gap(11.5).Padding(10, 9.25, 10, 11.5).Radius(12).Background(fill).TextColor(tint(fill, t.Dark, 0.847))
	bar.Children(func() {
		color := macInfoBarColor(t, severity)
		ui.Box(c).Size(13, line).Shrink(0).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			macInfoBarSymbol(p, ui.Rect{X: r.X, Y: r.Y + (r.H-13)/2, W: 13, H: 13}, severity, color, fill)
		})
		ui.Column(c).Grow(1).Shrink(1).Gap(2).Children(func() {
			ui.Text(c, title).Bold()
			if message != "" {
				ui.Text(c, message).TextColor(secondary)
			}
		})
		if closer == nil {
			return
		}
		b := closer("Close")
		glyph := secondary
		if b.Pressed() {
			glyph = tint(fill, t.Dark, 0.847)
		}
		b.Size(16, line).Center().Children(func() {
			ui.Box(c).Size(16, 16).Draw(func(p *ui.Painter, r ui.Rect) { toastCross(p, r, 0.17, 1.3, glyph) })
		})
		macRing(b, 4)
	})
}

// macInfoBarSymbol draws the SF Symbol of a severity in r, 13 points: a circle, a triangle with rounded corners or an octagon in col, with its mark (i, a check mark, ! or ×) in cut, the fill behind it, as SwiftUI draws a .fill symbol in one color. The shapes are drawn as paths from the symbols' outlines.
func macInfoBarSymbol(p *ui.Painter, r ui.Rect, severity int, col, cut ui.Color) {
	cx, cy, s := r.X+r.W/2, r.Y+r.H/2, r.W/13
	switch severity {
	case fluentToastWarning:
		var tri ui.Path
		tri.MoveTo(cx, r.Y+1.2*s).LineTo(r.X+r.W-0.9*s, r.Y+r.H-1.1*s).LineTo(r.X+0.9*s, r.Y+r.H-1.1*s).Close()
		p.FillPath(&tri, col)
		p.StrokePath(&tri, 1.6*s, col)
		p.Fill(ui.Rect{X: cx - 0.7*s, Y: cy - 2.6*s, W: 1.4 * s, H: 4.2 * s}, cut, 0.7*s)
		p.Fill(ui.Rect{X: cx - 0.8*s, Y: cy + 2.5*s, W: 1.6 * s, H: 1.6 * s}, cut, 0.8*s)
		return
	case fluentToastError:
		var oct ui.Path
		for i := range 8 {
			a := (float64(i) + 0.5) * math.Pi / 4
			x, y := cx+6.5*s*float32(math.Cos(a)), cy+6.5*s*float32(math.Sin(a))
			if i == 0 {
				oct.MoveTo(x, y)
			} else {
				oct.LineTo(x, y)
			}
		}
		oct.Close()
		p.FillPath(&oct, col)
		toastCross(p, ui.Rect{X: cx - 2.6*s, Y: cy - 2.6*s, W: 5.2 * s, H: 5.2 * s}, 0, 1.4*s, cut)
		return
	}
	p.Fill(r, col, r.W/2)
	if severity == fluentToastSuccess {
		var check ui.Path
		check.MoveTo(cx-2.9*s, cy+0.2*s).LineTo(cx-0.9*s, cy+2.2*s).LineTo(cx+3*s, cy-2.2*s)
		p.StrokePath(&check, 1.4*s, cut)
		return
	}
	p.Fill(ui.Rect{X: cx - 0.8*s, Y: cy - 3.6*s, W: 1.6 * s, H: 1.6 * s}, cut, 0.8*s)
	p.Fill(ui.Rect{X: cx - 0.7*s, Y: cy - 1*s, W: 1.4 * s, H: 4.4 * s}, cut, 0.7*s)
}

// fluentInfoBar draws a WinUI InfoBar in the layout, as InfoBar.xaml lays it out and as fluentToast draws it floating: at least 48 DIPs high (InfoBarMinHeight), ControlCornerRadius, a 1-DIP CardStrokeColorDefault border (InfoBarBorderThickness), InfoBarContentRootPadding 16 at the start; the 16-DIP severity icon (fluentToastIcon) with InfoBarIconMargin 0,16,14,16; the title in 14-DIP semibold and the message at regular weight, 16 before the close button (InfoBarPanelMargin); and the close button (fluentInfoBarClose). InfoBarPanel puts the message 12 after the title, both 14 from the top (InfoBarTitleHorizontalOrientationMargin, InfoBarMessageHorizontalOrientationMargin), while the two fit on one line, and otherwise the message 4 below the title, the panel padded 14 above and 18 below (InfoBarPanelVerticalOrientationPadding, InfoBarMessageVerticalOrientationMargin); Breeze UI decides it from the bar's width in the previous frame. In the layout the informational bar's background is SystemFillColorAttentionBackground, #80F6F6F6 and #08FFFFFF over the window, where the floating toast's is the solid one; the other severities' are the same solid colors (fluentToastColors). It has no shadow.
func fluentInfoBar(c *ui.Context, bar ui.Element, severity int, title, message string, closer func(string) ui.Element) {
	t := c.Theme()
	bg, iconFill := fluentToastColors(t, fluentSystemAccent(c), severity)
	if severity == fluentToastInformational {
		bg = pick(t.Dark, t.Background.Mix(ui.Hex("#f6f6f6"), 0.502), t.Background.Mix(white, 0.031))
	}
	// The panel's width: the bar less its border, the start padding, the icon and its margin, the panel's margin and the close button with its margins.
	panel := prevWidth(c, bar) - 2 - 16 - 16 - 14 - 16
	if closer != nil {
		panel -= 48
	}
	tw, _ := c.MeasureText(0, ui.Span{Text: title, Size: 14, Weight: 600})
	mw, _ := c.MeasureText(0, ui.Span{Text: message, Size: 14})
	across := message == "" || tw+12+mw <= panel
	bar.AlignItems(ui.Start).MinHeight(48).PaddingStart(16).Radius(4).Background(bg).
		Border(1, pick(t.Dark, ui.RGBA(0, 0, 0, 0.059), ui.RGBA(0, 0, 0, 0.098))).TextColor(t.Text).FontSize(14)
	bar.Children(func() {
		ui.Box(c).Size(16, 16).Shrink(0).MarginY(16).MarginEnd(14).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			fluentToastIcon(p, r, severity, iconFill, pick(t.Dark, white, ui.RGBA(0, 0, 0, 0.894)))
		})
		if across {
			ui.Row(c).Grow(1).Shrink(1).AlignItems(ui.Start).MarginEnd(16).Padding(0, 0, 14, 0).Children(func() {
				ui.Text(c, title).FontWeight(600).Margin(14, 0, 0, 0)
				if message != "" {
					ui.Text(c, message).Shrink(1).Margin(14, 0, 0, 0).MarginStart(12)
				}
			})
		} else {
			ui.Column(c).Grow(1).Shrink(1).AlignItems(ui.Start).MarginEnd(16).Padding(14, 0, 18, 0).Children(func() {
				ui.Text(c, title).FontWeight(600)
				ui.Text(c, message).Margin(4, 0, 0, 0)
			})
		}
		if closer != nil {
			fluentInfoBarClose(c, closer("Close"))
		}
	})
}

// prevWidth returns e's width in the previous frame, as a platform's layout that switches between arrangements by the width it is given knows it; for an element that frame lacked it returns the window's width and asks for another frame, which lays it out at its own.
func prevWidth(c *ui.Context, e ui.Element) float32 {
	if w := e.Bounds().W; w > 0 {
		return w
	}
	c.Invalidate()
	w, _ := c.Size()
	return w
}

// fluentInfoBarClose draws an InfoBar's close button (InfoBarCloseButtonStyle): 38 DIPs square (InfoBarCloseButtonSize) with a margin of 5, at the top, ControlCornerRadius, drawn with AppBarButton's brushes as InfoBar.xaml sets them: transparent, SubtleFillColorSecondary under the pointer and SubtleFillColorTertiary pressed; its 16-DIP Cancel symbol (InfoBarCloseButtonGlyphSize) in the text color.
func fluentInfoBarClose(c *ui.Context, b ui.Element) {
	t := c.Theme()
	b.Size(38, 38).Margin(5).Shrink(0).Radius(4).Center()
	hover, pressed := fluentSubtle(t)
	switch {
	case b.Pressed():
		b.Background(pressed)
	case b.Hovered():
		b.Background(hover)
	}
	b.Children(func() {
		ui.Box(c).Size(16, 16).Draw(func(p *ui.Painter, r ui.Rect) { toastCross(p, r, 0.22, 1, t.Text) })
	})
}

// adwBannerShares are the shares of currentColor a GtkButton's face lays over what is behind it, at rest, under the pointer and pressed, in the light theme (where currentColor is window_fg_color, 80% opaque) and the dark one: $button_color and its hover and active shades.
var adwBannerShares = [2][3]float32{{0.08, 0.12, 0.24}, {0.10, 0.15, 0.30}}

// adwBanner draws an AdwBanner: a full-width strip of --banner-color (#7d7d83) mixed 30% into window_bg_color, in window_fg_color, padded 6 (_toolbars.scss); its title in .heading, bold, wrapping, centered across its lines; and, while it can be closed, its button, a GtkButton at least 84 wide (BUTTON_HORZ_MIN_WIDTH), as adw-banner.c's allocate_content places them: without a button the title centered and padded a further 6 all around (HORZ_PADDING, VERT_PADDING); with one, the button at the end, centered on the strip's height, and the title centered across the whole strip while it fits with 36 to spare on each side of the button (HORZ_SPACING_CENTERED), else 6 from the start (HORZ_PADDING) and 6 before the button (HORZ_SPACING); and where the title does not fit beside the button at its natural width, the button below the title, centered, at least 160 wide (BUTTON_VERT_MIN_WIDTH), 9 below it (VERT_SPACING), the two padded 6 above, below and at the sides. Breeze UI picks the arrangement from the banner's width in the previous frame. A banner has only a title: the message follows it in the same label, at regular weight (inferred). libadwaita's banners have no severity; the .error, .warning and .success classes only recolor accent-colored parts, which a banner's title is not, so every kind draws this same banner.
func adwBanner(c *ui.Context, bar ui.Element, title, message string, closer func(string) ui.Element) {
	t := c.Theme()
	bg := t.Background.Mix(ui.Hex("#7d7d83"), 0.3)
	fg := pick(t.Dark, bg.Mix(adwInk, 0.8), white)
	bar.Padding(6).Background(bg).TextColor(fg)
	text := func() ui.Element {
		return ui.RichText(c).Shrink(1).TextAlign(ui.Center).Children(func() {
			ui.Text(c, title).FontWeight(700)
			if message != "" {
				ui.Text(c, " "+message)
			}
		})
	}
	if closer == nil {
		bar.Justify(ui.Center).AlignItems(ui.Center).Children(func() { text().Padding(6) })
		return
	}
	const dismiss = "Dismiss"
	bw, _ := c.MeasureText(0, ui.Span{Text: dismiss, Size: t.FontSize, Weight: 700})
	bw += 34 // the button's padding of 17 on each side
	tw, _ := c.MeasureText(0, ui.Span{Text: title, Size: t.FontSize, Weight: 700})
	if message != "" {
		mw, _ := c.MeasureText(0, ui.Span{Text: " " + message, Size: t.FontSize})
		tw += mw
	}
	width := prevWidth(c, bar) - 12 // the strip's padding
	button := func() ui.Element {
		b := closer(dismiss)
		adwButton(c, b, false)
		shares := adwBannerShares[map[bool]int{false: 0, true: 1}[t.Dark]]
		share := shares[0]
		switch {
		case b.Pressed():
			share = shares[2]
		case b.Hovered():
			share = shares[1]
		}
		face := bg.Mix(pick(t.Dark, adwInk, white), share)
		return b.Background(face).TextColor(adwTextOn(t, face)).Children(func() { ui.Text(c, dismiss).SingleLine() })
	}
	switch {
	case 6+tw+6+max(bw, 84) > width:
		bar.Column().AlignItems(ui.Center).Gap(9).Padding(12).Children(func() {
			text()
			button().MinWidth(min(max(bw, 160), width))
		})
	case tw+(max(bw, 84)+36)*2 > width:
		bar.AlignItems(ui.Center).Gap(6).PaddingStart(12).Children(func() {
			text().Grow(1).TextAlign(ui.Start)
			button().MinWidth(84)
		})
	default:
		bar.AlignItems(ui.Center).Gap(6).Children(func() {
			// A space as wide as the button keeps the title centered across the strip.
			ui.Box(c).Width(max(bw, 84)).Shrink(0).Role(ui.RoleNone)
			text().Grow(1)
			button().MinWidth(84)
		})
	}
}

// ProgressRing creates a circular progress indicator filled to value between 0 and 1, as the platform draws a determinate ring: AppKit's NSProgressIndicator of style spinning with isIndeterminate false, WinUI's ProgressRing with IsIndeterminate false, and on GNOME, which has no determinate ring, a ring after AdwSpinner's (inferred). The ring fills clockwise from 12 o'clock in either reading direction, as a clock does. A negative value shows activity of unknown length: it is then Spinner. Assistive technology sees a progress indicator of value from 0 to 1, as for Progress; name it with Label. It is 32 points or DIPs across, as the platforms' regular rings are; Size makes it another size, as AppKit's small one, 16.
func ProgressRing(c *ui.Context, value float64) ui.Element {
	if value < 0 {
		return Spinner(c)
	}
	value = min(value, 1)
	e := ui.Box(c).Size(32, 32).Shrink(0).Role(ui.RoleProgress).Range(0, 1, value)
	switch PlatformOf(c) {
	case Windows:
		fluentProgressRing(c, e, value)
	case Linux:
		adwProgressRing(c, e, value)
	default:
		macProgressRing(c, e, value)
	}
	return e
}

// macRingArc and macRingFull are the arc of a determinate spinning NSProgressIndicator in Dark Aqua for the default accent, measured: #1769e6, and #177be5 once the ring is full, where Aqua draws the accent itself in both (#0a82ff measured over white for #007aff).
var macRingArc, macRingFull = ui.Hex("#1769e6"), ui.Hex("#177be5")

// macProgressRing draws an NSProgressIndicator of style spinning with isIndeterminate false, as measured by reference.swift on macOS 27.2: a track 5 points wide inside the edge, at both the regular size (32 points) and the small one (16, whose ring is then 3 points from its center), black at 5.1% in Aqua (242 over white, 235 along its edges) and white at 10.2% in Dark Aqua (53 over 30); and over it the arc of the value from 12 o'clock clockwise, as wide, with round ends, in the accent (in Dark Aqua #1769e6 for the default blue, #177be5 once full, measured, and the accent 10% darker for another, inferred). The track's darker edges and the arc's, a quarter pixel at 2x, are not drawn.
func macProgressRing(c *ui.Context, e ui.Element, value float64) {
	t := c.Theme()
	track := tint(t.Background, t.Dark, pick(t.Dark, float32(0.051), 0.102))
	arc := t.Accent
	if t.Dark {
		arc = t.Accent.Mix(black, 0.1)
		if t.Accent == macAccent {
			arc = macRingArc
			if value >= 1 {
				arc = macRingFull
			}
		}
	}
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		d := min(r.W, r.H)
		w := min(float32(5), d/3)
		cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, d/2-w/2
		var ring ui.Path
		ring.Circle(cx, cy, rad)
		p.StrokePath(&ring, w, track)
		if value > 0 {
			strokeArc(p, cx, cy, rad, 0, value*2*math.Pi, w, arc)
		}
	})
}

// fluentRingTrim returns how much of a determinate ProgressRing's ellipse is drawn at value, as the Lottie of ProgressRingDeterminate.cpp maps the player's progress, which ProgressRing.cpp sets to the value's share of the range (UpdateLottieProgress), to the ellipse's TrimEnd: hidden until 1/120 (its container's Scale steps from 0 to 1 there), then keyframes at 0.25, 0.5 and 0.75 equal to the progress, 0.9667 at 0.9833 and 1 at 0.9917, along cubic-bezier(0.167,0.167,0.833,0.833) between them.
func fluentRingTrim(value float64) float64 {
	keys := [][2]float64{{1.0 / 120, 0.0001}, {0.25, 0.25}, {0.5, 0.5}, {0.75, 0.75}, {0.983333349, 0.96666666}, {0.991666675, 1}}
	if value < keys[0][0] {
		return 0
	}
	ease := ui.CubicBezier(0.167, 0.167, 0.833, 0.833)
	for i := 1; i < len(keys); i++ {
		if value <= keys[i][0] {
			a, b := keys[i-1], keys[i]
			return a[1] + (b[1]-a[1])*float64(ease(float32((value-a[0])/(b[0]-a[0]))))
		}
	}
	return 1
}

// fluentProgressRing draws a ProgressRing with IsIndeterminate false, as its Lottie, ProgressRingDeterminate.cpp, draws it on the 32-DIP control (the style's Width and Height): an ellipse of radius 8 scaled 1.77 around the middle of a 32-unit canvas, so radius 14.16, stroked 1.5 scaled 1.77, so 2.655 DIPs, with round ends, in ProgressRingForegroundThemeBrush, AccentFillColorDefault (the theme's Accent), trimmed to the value (fluentRingTrim); under it the whole ellipse in ProgressRingBackgroundThemeBrush, ControlFillColorTransparent, which shows nothing. Stretch fill scales the canvas to another size. Where the trim starts is not in the source: 12 o'clock, clockwise, as the indeterminate ring's (fluentSpinner).
func fluentProgressRing(c *ui.Context, e ui.Element, value float64) {
	t := c.Theme()
	trim := fluentRingTrim(value)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		if trim <= 0 {
			return
		}
		s := min(r.W, r.H) / 32
		strokeArc(p, r.X+r.W/2, r.Y+r.H/2, 14.16*s, 0, trim*2*math.Pi, 2.655*s, t.Accent)
	})
}

// adwProgressRing draws a determinate ring on GNOME, which has none, after AdwSpinner (adwSpinner): the track a ring of the text color at 15%, stroked an eighth of the diameter wide as adw-spinner-paintable.c strokes it at the icon weight of 400 (calculate_line_width), 4 pixels at 32; and over it the value from 12 o'clock clockwise with round ends in the accent, as GtkProgressBar fills (adwProgress). The size, 32 pixels, the accent and where the arc starts are inferred.
func adwProgressRing(c *ui.Context, e ui.Element, value float64) {
	t := c.Theme()
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		d := min(r.W, r.H)
		w := d / 8
		cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, d/2-w/2
		var ring ui.Path
		ring.Circle(cx, cy, rad)
		p.StrokePath(&ring, w, t.Background.Mix(t.Text, 0.15))
		if value > 0 {
			strokeArc(p, cx, cy, rad, 0, value*2*math.Pi, w, t.Accent)
		}
	})
}

// StatusPage creates a page telling that a view is empty or what state it is in, with an icon, a title and a description, as the platform draws an empty state: libadwaita's AdwStatusPage on GNOME, SwiftUI's ContentUnavailableView on macOS, and on Windows, which has no such control, a page after WinUI's type ramp (inferred). fn builds actions below the text, as buttons, in a row; nil for none. A nil icon or an empty description leaves it out. The page grows to fill its parent and centers its content; for assistive technology it is a group named by its title and described by its description, and its title is a heading.
func StatusPage(c *ui.Context, icon *ui.SVG, title, description string, fn func()) ui.Element {
	t := c.Theme()
	page := ui.Column(c).Grow(1).AlignSelf(ui.Stretch).Justify(ui.Center).AlignItems(ui.Center).Role(ui.RoleGroup).Label(title)
	if description != "" {
		page.Description(description)
	}
	var look statusLook
	switch PlatformOf(c) {
	case Windows:
		look = fluentStatusLook(t)
	case Linux:
		look = adwStatusLook(t)
	default:
		look = macStatusLook(t)
	}
	page.Padding(look.padY, look.padX)
	page.Children(func() {
		if icon != nil {
			ui.Icon(c, icon).Size(look.icon, look.icon).Shrink(0).TextColor(look.iconColor).Margin(0, 0, look.iconGap, 0)
		}
		h := ui.Text(c, title).Role(ui.RoleHeading).Level(1).FontSize(look.titleSize).FontWeight(look.titleWeight).TextColor(look.titleColor).TextAlign(ui.Center)
		var d ui.Element
		if description != "" {
			d = ui.Text(c, description).FontSize(look.bodySize).TextColor(look.bodyColor).TextAlign(ui.Center).Margin(look.textGap, 0, 0, 0)
			if look.bodyLine > 0 {
				d.LineHeight(look.bodyLine)
			}
		}
		if look.maxWidth > 0 {
			h.MaxWidth(look.maxWidth)
			d.MaxWidth(look.maxWidth)
		}
		if fn != nil {
			ui.Row(c).Gap(look.actionGap).AlignItems(ui.Center).Justify(ui.Center).Wrap().Margin(look.actionsGap, 0, 0, 0).Children(fn)
		}
	})
	return page
}

// statusLook is how a platform lays out a StatusPage: the page's padding, the icon's size and color and the space after it, the title's and the description's sizes, weights and colors, the space between them and before the actions, the space between actions, and how wide the text may be, 0 for as wide as the page.
type statusLook struct {
	padY, padX                       float32
	icon, iconGap                    float32
	iconColor, titleColor, bodyColor ui.Color
	titleSize, bodySize, bodyLine    float32
	titleWeight                      int
	textGap, actionsGap, actionGap   float32
	maxWidth                         float32
}

// macStatusLook is ContentUnavailableView's, measured from SwiftUI on macOS 27.2 with a label of an SF Symbol and a title, a description and a button: the symbol's ink 35 points tall in tertiaryLabelColor (black at 25.9%, white at 25%: 189 over white, 86 over 30), the title 28 points below it in 15-point bold (title3, two points over the body size) and the description 7 below the title in 15-point regular, both in secondaryLabelColor and centered, and the actions 19 below the description. The gaps are between line boxes, read from the glyphs and SF's metrics. The page's padding of 20 and the 12 between actions are inferred.
func macStatusLook(t *ui.Theme) statusLook {
	size := t.FontSize + 2
	return statusLook{
		padY: 20, padX: 20,
		icon: 35, iconGap: 28, iconColor: tint(t.Background, t.Dark, pick(t.Dark, float32(0.259), 0.249)),
		titleSize: size, titleWeight: 700, titleColor: t.TextMuted,
		bodySize: size, bodyColor: t.TextMuted,
		textGap: 7, actionsGap: 19, actionGap: 12,
	}
}

// fluentStatusLook lays out an empty state on Windows, which has no control for one, after WinUI's type ramp: the title in SubtitleTextBlockStyle (20 DIPs, semibold, TextFillColorPrimary), the description in BodyTextBlockStyle (14 DIPs, regular) in TextFillColorSecondary, a 48-DIP icon in TextFillColorSecondary 16 above the title, the description 8 below it, the actions 24 below and 8 apart (ContentDialogButtonSpacing), the page padded 24 and the text at most 600 wide. All of it is inferred.
func fluentStatusLook(t *ui.Theme) statusLook {
	return statusLook{
		padY: 24, padX: 24,
		icon: 48, iconGap: 16, iconColor: t.TextMuted,
		titleSize: 20, titleWeight: 600, titleColor: t.Text,
		bodySize: 14, bodyColor: t.TextMuted,
		textGap: 8, actionsGap: 24, actionGap: 8,
		maxWidth: 600,
	}
}

// adwStatusLook is AdwStatusPage's (_misc.scss, adw-status-page.ui): its content box with a margin of 36 pixels above and below and 12 at the sides; the 128-pixel icon in currentColor at --dim-opacity (55%, the theme's TextMuted) 36 above the title; the title in .title-1, weight 800 at 181% of the font; the description 12 below it in .body, regular with a line height of 140%, both in window_fg_color, centered and at most 600 wide (AdwClamp's maximum-size); the actions, the page's child, 36 below the text. The 12 between actions are inferred.
func adwStatusLook(t *ui.Theme) statusLook {
	return statusLook{
		padY: 36, padX: 12,
		icon: 128, iconGap: 36, iconColor: t.TextMuted,
		titleSize: t.FontSize * 1.81, titleWeight: 800, titleColor: t.Text,
		bodySize: t.FontSize, bodyColor: t.Text, bodyLine: 1.4,
		textGap: 12, actionsGap: 36, actionGap: 12,
		maxWidth: 600,
	}
}
