package breeze

import (
	"math"
	"strings"
	"unicode"
	"unicode/utf16"

	"github.com/egoist/mygo/ui"
)

// ---- feedback: tooltips, toasts, meters, avatars, ratings and the find bar ----
//
// macOS values are measured from AppKit's NSLevelIndicator and NSTextFinder's find bar rendered at 2x, as tools/appkit-reference/reference.swift renders the other controls, the colored ones drawn into an sRGB bitmap (the display's profile, which that script reads, gives the same grays but paler colors), or read from AppKit (NSFont.toolTipsFont, NSColor.systemYellow); AppKit draws tooltips, notifications and contact monograms out of a view's reach, so those are inferred. Windows values are WinUI 3's theme resources (ToolTip_themeresources.xaml, InfoBar_themeresources.xaml and InfoBar.xaml, ProgressBar_themeresources.xaml, RatingControl_themeresources.xaml and RatingControl.cpp, PersonPicture_themeresources.xaml, PersonPicture.xaml, PersonPicture.cpp and InitialsGenerator.cpp, Common_themeresources_any.xaml, AcrylicBrush_themeresources.xaml). GNOME values are libadwaita 1.7.6's stylesheet (_tooltip.scss, _misc.scss for toasts, _level-bar.scss, _avatar.scss, _toolbars.scss for the search bar, _colors.scss, _defaults.scss), adw-toast-widget.ui and .c, adw-avatar.c, and GTK 4.18's gtktooltip.c, gtktooltipwindow.ui and gtklevelbar.c; GNOME has no rating control, so its stars are GNOME Software's (gs-star-widget.c, gs-star-image.c, gs-review-dialog.ui and style.css of 48.0).

// Tooltip shows text by anchor while the pointer rests on it or the keyboard focus comes to it, in the platform's tooltip: a small panel below the element on macOS, WinUI's ToolTip above it on Windows, and GTK's tooltip below it on GNOME. It is built on ui.TooltipBase, so the timing, Escape, a disabled anchor showing none and the one-at-a-time rule are MyGo's. Assistive technology does not read the tip, so text becomes anchor's Description, as Element.Tooltip makes it; unlike Element.Tooltip, which keeps a Description given before it, Tooltip replaces one, as package ui offers no way to read it: give a Description of your own after Tooltip. It returns anchor, as Element.Tooltip does, whose counterpart it is.
func Tooltip(c *ui.Context, anchor ui.Element, text string) ui.Element {
	if text == "" || !anchor.Valid() {
		// ui.TooltipBase needs an anchor of the frame; Element.Tooltip on an empty Element does nothing too.
		return anchor
	}
	anchor.Description(text)
	p := PlatformOf(c)
	ui.TooltipBase(c, anchor, func(tip ui.Element) {
		switch p {
		case Windows:
			fluentTooltip(c, anchor, tip)
		case Linux:
			adwTooltip(c, anchor, tip)
		default:
			macTooltip(c, anchor, tip)
		}
		tip.Children(func() { ui.Text(c, text) })
	})
	return anchor
}

// macTooltip draws an AppKit tooltip: text in NSFont.toolTipsFont (11 points, read from AppKit), in a light or dark panel with a hairline edge, below the element. AppKit draws tooltips in a window of its own, which does not render offscreen: the panel's colors, padding, corners and shadow are inferred, following the menus of mac.go.
func macTooltip(c *ui.Context, anchor, tip ui.Element) {
	t := c.Theme()
	bg := pick(t.Dark, ui.Hex("#f6f6f6"), ui.Hex("#2b2b2b"))
	tip.AttachTo(anchor, ui.AnchorBottom, ui.AnchorTop).Margin(4, 0, 0, 0).
		MaxWidth(300).Padding(3, 7).Radius(5).Background(bg).Border(0.5, tint(bg, t.Dark, 0.18)).
		TextColor(t.Text).FontSize(11).Shadow(0, 2, 6, 0, ui.RGBA(0, 0, 0, pick(t.Dark, float32(0.12), 0.3)))
}

// fluentTooltip draws a ToolTip: ToolTipBorderPadding 9,6,9,8, ControlCornerRadius, a 1-DIP SurfaceStrokeColorFlyout border around the in-app acrylic's fallback color (the theme's Inverse), 12-DIP text in TextFillColorPrimary, at most ToolTipMaxWidth (320) wide, above the element. The 4-DIP gap and the shadow are inferred.
func fluentTooltip(c *ui.Context, anchor, tip ui.Element) {
	t := c.Theme()
	tip.AttachTo(anchor, ui.AnchorTop, ui.AnchorBottom).Margin(0, 0, 4, 0).
		MaxWidth(320).Padding(6, 9, 8, 9).Radius(4).Background(t.Inverse).TextColor(t.InverseText).FontSize(12).
		Border(1, pick(t.Dark, ui.RGBA(0, 0, 0, 0.059), ui.RGBA(0, 0, 0, 0.2))).
		Shadow(0, 4, 8, 0, ui.RGBA(0, 0, 0, pick(t.Dark, float32(0.14), 0.26)))
}

// adwTooltip draws a GTK tooltip as libadwaita styles it: the shade at 80% (RGB(0 0 6 / 80%), the theme's Inverse) inside a 1-pixel $tooltip_border_color (white at 10%), white text, padding 6px 10px, corners of 9 and no shadow. GTK puts it below the widget, 4 pixels from it (gtktooltip.c's anchor_rect_padding), and wraps its label at 50 characters (gtktooltipwindow.ui), about 25 em.
func adwTooltip(c *ui.Context, anchor, tip ui.Element) {
	t := c.Theme()
	tip.AttachTo(anchor, ui.AnchorBottom, ui.AnchorTop).Margin(4, 0, 0, 0).
		MaxWidth(t.FontSize*25+20).Padding(6, 10).Radius(9).Background(t.Inverse).TextColor(t.InverseText).
		Border(1, ui.RGBA(255, 255, 255, 0.1))
}

// Toasts draws the window's toasts, the ones Context.Toast, ToastAction and AddToast show, in the platform's look in place of MyGo's: GNOME's AdwToast, a dark pill at the bottom of the window; on Windows an InfoBar, colored by the toast's Type ("success", "warning" or "error"); on macOS, which has no toast, a panel like a notification banner. Call it once in the view, wherever, after Use. It is built on ui.ToastViewportBase and ui.ToastBase, so the toasts' time, their stopping under the pointer, Escape, the action and what screen readers read are MyGo's. It returns the toasts' viewport.
func Toasts(c *ui.Context) ui.Element {
	p := PlatformOf(c)
	return ui.ToastViewportBase(c, func(viewport ui.Element, toasts []ui.Toast) {
		switch p {
		case Windows:
			fluentToastViewport(c, viewport)
		case Linux:
			adwToastViewport(c, viewport)
		default:
			macToastViewport(c, viewport)
		}
		for _, ts := range toasts {
			parts := ui.ToastBase(c, ts)
			// Toasts fade in and out, as MyGo's own do; the platforms' motion is not measured.
			parts.Root.Transition(ui.ElementTransition{Position: true, Enter: &ui.Motion{Opacity: 0, Y: 8}, Exit: &ui.Motion{Opacity: 0}})
			switch p {
			case Windows:
				fluentToast(c, parts, ts)
			case Linux:
				adwToast(c, parts, ts)
			default:
				macToast(c, parts, ts)
			}
		}
	})
}

// macToastViewport holds the toasts at the bottom of the window, centered, as MyGo places its own; macOS shows no toasts in windows, so the place is inferred.
func macToastViewport(c *ui.Context, viewport ui.Element) {
	viewport.Padding(0, 16, 24, 16).AlignItems(ui.Center).Gap(8)
}

// macToast draws a toast as a panel like a macOS notification banner: the bold title over the description, a push button for the action beside them, in a rounded panel with a hairline edge and a soft shadow. macOS has no toast control and its banners do not render offscreen, so the corners, padding and colors are inferred, following the menus of mac.go.
func macToast(c *ui.Context, parts ui.ToastParts, ts ui.Toast) {
	t := c.Theme()
	w, _ := c.Size()
	bg := pick(t.Dark, ui.Hex("#f6f6f6"), ui.Hex("#2b2b2b"))
	parts.Root.Row().AlignItems(ui.Center).Gap(12).Padding(10, 14).MinHeight(44).MaxWidth(min(w-32, 420)).Radius(14).
		Background(bg).Border(0.5, tint(bg, t.Dark, 0.15)).TextColor(t.Text).
		Shadow(0, 8, 24, 0, ui.RGBA(0, 0, 0, pick(t.Dark, float32(0.18), 0.4)))
	parts.Root.Children(func() {
		ui.Column(c).Grow(1).Shrink(1).Gap(2).Children(func() {
			ui.Text(c, ts.Title).Bold()
			if ts.Description != "" {
				ui.Text(c, ts.Description).TextColor(t.TextMuted)
			}
		})
		if ts.Action != "" {
			b := parts.ActionButton()
			macButton(c, b, false)
			b.Shrink(0).Children(func() { ui.Text(c, ts.Action).SingleLine() })
		}
	})
}

// fluentToastViewport holds the toasts at the bottom of the window, centered. WinUI's InfoBar sits in the layout rather than over it, so the place and the 16-DIP margins are inferred.
func fluentToastViewport(c *ui.Context, viewport ui.Element) {
	viewport.Padding(0, 16, 24, 16).AlignItems(ui.Center).Gap(8)
}

// InfoBar severities, from a toast's Type.
const (
	fluentToastInformational = iota
	fluentToastSuccess
	fluentToastWarning
	fluentToastError
)

// fluentToastSeverity returns the InfoBar severity a toast's Type names: "success", "warning", "error" (or "danger", "critical"); any other is informational.
func fluentToastSeverity(typ string) int {
	switch strings.ToLower(typ) {
	case "success":
		return fluentToastSuccess
	case "warning":
		return fluentToastWarning
	case "error", "danger", "critical":
		return fluentToastError
	}
	return fluentToastInformational
}

// fluentToastColors returns an InfoBar's background and its icon's fill for a severity: SystemFillColorSolidAttentionBackground (the solid one, as the toast floats over content) or SystemFillColorSuccessBackground, CautionBackground and CriticalBackground behind the content; SystemFillColorAttention (the accent fill) or SystemFillColorSuccess, Caution and Critical for the icon.
func fluentToastColors(t *ui.Theme, severity int) (bg, icon ui.Color) {
	switch severity {
	case fluentToastSuccess:
		return pick(t.Dark, ui.Hex("#dff6dd"), ui.Hex("#393d1b")), pick(t.Dark, ui.Hex("#0f7b0f"), ui.Hex("#6ccb5f"))
	case fluentToastWarning:
		return pick(t.Dark, ui.Hex("#fff4ce"), ui.Hex("#433519")), pick(t.Dark, ui.Hex("#9d5d00"), ui.Hex("#fce100"))
	case fluentToastError:
		return pick(t.Dark, ui.Hex("#fde7e9"), ui.Hex("#442726")), pick(t.Dark, ui.Hex("#c42b1c"), ui.Hex("#ff99a4"))
	}
	return pick(t.Dark, ui.Hex("#f7f7f7"), ui.Hex("#2e2e2e")), t.Accent
}

// fluentToast draws a toast as an InfoBar: at least 48 DIPs high, ControlCornerRadius, a 1-DIP CardStrokeColorDefault border, InfoBarContentRootPadding 16 at the start; the 16-DIP severity icon with InfoBarIconMargin 0,16,14,16; the title semibold and the message 12 after it, both 14 DIPs and 14 from the top; the action button 16 after them and 8 from the top; and the close button, 38 DIPs with a margin of 5, drawn as an AppBarButton. The shadow, as the toast floats over the window, is a flyout's, inferred.
func fluentToast(c *ui.Context, parts ui.ToastParts, ts ui.Toast) {
	t := c.Theme()
	w, _ := c.Size()
	severity := fluentToastSeverity(ts.Type)
	bg, iconFill := fluentToastColors(t, severity)
	parts.Root.Row().AlignItems(ui.Start).MinHeight(48).MaxWidth(min(w-32, 640)).PaddingStart(16).Radius(4).
		Background(bg).Border(1, pick(t.Dark, ui.RGBA(0, 0, 0, 0.059), ui.RGBA(0, 0, 0, 0.098))).TextColor(t.Text).FontSize(14).
		Shadow(0, 8, 16, 0, ui.RGBA(0, 0, 0, pick(t.Dark, float32(0.14), 0.26)))
	parts.Root.Children(func() {
		ui.Box(c).Size(16, 16).Shrink(0).MarginY(16).MarginEnd(14).Draw(func(p *ui.Painter, r ui.Rect) {
			fluentToastIcon(p, r, severity, iconFill, pick(t.Dark, white, ui.RGBA(0, 0, 0, 0.894)))
		})
		ui.Row(c).Grow(1).Shrink(1).Wrap().AlignItems(ui.Start).MarginEnd(16).Padding(0, 0, 14, 0).Children(func() {
			ui.Text(c, ts.Title).FontWeight(600).Margin(14, 0, 0, 0)
			if ts.Description != "" {
				ui.Text(c, ts.Description).Margin(14, 0, 0, 0).MarginStart(12)
			}
			if ts.Action != "" {
				b := parts.ActionButton()
				fluentButton(c, b, false)
				b.Margin(8, 0, 0, 0).MarginStart(16).Children(func() { ui.Text(c, ts.Action).SingleLine() })
			}
		})
		dismiss := parts.CloseButton().Label("Close").Size(38, 38).Margin(5).Shrink(0).Radius(4).Center()
		hover, pressed := fluentSubtle(t)
		switch {
		case dismiss.Pressed():
			dismiss.Background(pressed)
		case dismiss.Hovered():
			dismiss.Background(hover)
		}
		dismiss.Children(func() {
			ui.Box(c).Size(16, 16).Draw(func(p *ui.Painter, r ui.Rect) { toastCross(p, r, 0.22, 1, t.Text) })
		})
	})
}

// fluentToastIcon draws an InfoBar's icon: the round InfoBarIconBackgroundGlyph in the severity's color, and on it in TextFillColorInverse the glyph of the severity (Segoe Fluent Icons' info, check mark, exclamation mark or cross), drawn here as paths.
func fluentToastIcon(p *ui.Painter, r ui.Rect, severity int, fill, mark ui.Color) {
	p.Fill(r, fill, r.W/2)
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	switch severity {
	case fluentToastSuccess:
		var path ui.Path
		path.MoveTo(cx-3.5, cy+0.2).LineTo(cx-1.2, cy+2.5).LineTo(cx+3.6, cy-2.4)
		p.StrokePath(&path, 1.3, mark)
	case fluentToastError:
		toastCross(p, ui.Rect{X: cx - 3, Y: cy - 3, W: 6, H: 6}, 0, 1.3, mark)
	case fluentToastWarning:
		p.Fill(ui.Rect{X: cx - 0.65, Y: cy - 4, W: 1.3, H: 5}, mark, 0.65)
		p.Fill(ui.Rect{X: cx - 0.8, Y: cy + 2.3, W: 1.6, H: 1.6}, mark, 0.8)
	default:
		p.Fill(ui.Rect{X: cx - 0.8, Y: cy - 4, W: 1.6, H: 1.6}, mark, 0.8)
		p.Fill(ui.Rect{X: cx - 0.65, Y: cy - 1.2, W: 1.3, H: 5.2}, mark, 0.65)
	}
}

// toastCross strokes an × in r, inset by a share in of its size on every side.
func toastCross(p *ui.Painter, r ui.Rect, in, width float32, col ui.Color) {
	x0, y0, x1, y1 := r.X+r.W*in, r.Y+r.H*in, r.X+r.W*(1-in), r.Y+r.H*(1-in)
	var path ui.Path
	path.MoveTo(x0, y0).LineTo(x1, y1).MoveTo(x1, y0).LineTo(x0, y1)
	p.StrokePath(&path, width, col)
}

// adwToastViewport holds the toasts as AdwToastOverlay does: at the bottom, centered, with toast's margin of 12 pixels and 24 below.
func adwToastViewport(c *ui.Context, viewport ui.Element) {
	viewport.Padding(0, 12, 24, 12).AlignItems(ui.Center).Gap(12)
}

// adwToast draws an AdwToast: $toast_bg_color (#505053) with white text, a pill (border-radius 150px) padded 6 pixels, 12 at the start, its parts 6 apart; the title in .heading, bold, with a margin of 0 6px; the action as a button (white at 10%, 15% under the pointer, 30% pressed, 34 pixels high, bold); and the close button, a 34-pixel .circular .flat button with window-close-symbolic (white at 7% under the pointer, 16% pressed), whose arms are 2 pixels wide and cut square at the edges of the square from (4,4) to (12,12) of its 16-pixel icon. Its two shadows are the stylesheet's. AdwToast has no description; one shows after the title, not bold, which is inferred.
func adwToast(c *ui.Context, parts ui.ToastParts, ts ui.Toast) {
	w, _ := c.Size()
	bg := ui.Hex("#505053")
	parts.Root.Row().AlignItems(ui.Center).Gap(6).Padding(6).PaddingStart(12).MinHeight(46).MaxWidth(w-24).Radius(23).
		Background(bg).TextColor(white).
		Shadow(0, 1, 3, 1, ui.RGBA(0, 0, 6, 0.07)).Shadow(0, 2, 6, 2, ui.RGBA(0, 0, 6, 0.03))
	parts.Root.Children(func() {
		ui.Row(c).Grow(1).Shrink(1).Gap(6).MarginX(6).AlignItems(ui.Center).Children(func() {
			ui.Text(c, ts.Title).FontWeight(700).SingleLine()
			if ts.Description != "" {
				ui.Text(c, ts.Description).SingleLine()
			}
		})
		if ts.Action != "" {
			b := parts.ActionButton()
			share := float32(0.10)
			switch {
			case b.Pressed():
				share = 0.30
			case b.Hovered():
				share = 0.15
			}
			b.MinHeight(34).Padding(5, 17).Radius(9).FontWeight(700).Justify(ui.Center).AlignItems(ui.Center).Shrink(0).
				Background(bg.Mix(white, share)).TextColor(white).Children(func() { ui.Text(c, ts.Action).SingleLine() })
		}
		dismiss := parts.CloseButton().Label("Dismiss").Size(34, 34).Radius(17).Shrink(0).Center()
		switch {
		case dismiss.Pressed():
			dismiss.Background(bg.Mix(white, 0.16))
		case dismiss.Hovered():
			dismiss.Background(bg.Mix(white, 0.07))
		}
		dismiss.Children(func() {
			// window-close-symbolic, a 16-pixel icon: the cross stroked 2 wide past the square its arms end in, clipped to it.
			ui.Box(c).Size(16, 16).Draw(func(p *ui.Painter, r ui.Rect) {
				p.Clip(ui.Rect{X: r.X + 4, Y: r.Y + 4, W: 8, H: 8}, 0, func() { toastCross(p, r, 0.1875, 2, white) })
			})
		})
	})
}

// meterLevel returns 0 for a value that is fine, 1 for one at the warning level and 2 for one at the critical level, as ui.Meter decides it: a Critical below Warning makes low values the bad ones.
func meterLevel(value float64, l *ui.MeterLevels) int {
	if l == nil || l.Warning == l.Critical {
		return 0
	}
	low := l.Critical < l.Warning
	switch {
	case !low && value >= l.Critical, low && value <= l.Critical:
		return 2
	case !low && value >= l.Warning, low && value <= l.Warning:
		return 1
	}
	return 0
}

// meterFraction returns where value is between lo and hi, from 0 to 1.
func meterFraction(value, lo, hi float64) float32 {
	if hi <= lo {
		return 0
	}
	return float32(min(max((value-lo)/(hi-lo), 0), 1))
}

// Meter creates a bar showing value between lo and hi, as the platform shows a level: AppKit's NSLevelIndicator of continuous capacity style in green, yellow or red, a WinUI ProgressBar in its normal, paused or error color, or a GtkLevelBar; levels, nil for none, give the warning and critical levels. It is ui.Meter's counterpart, whose role and value assistive technology sees.
func Meter(c *ui.Context, value, lo, hi float64, levels *ui.MeterLevels) ui.Element {
	m := ui.Meter(c, value, lo, hi, levels)
	level, frac := meterLevel(value, levels), meterFraction(value, lo, hi)
	switch PlatformOf(c) {
	case Windows:
		fluentMeter(c, m, frac, level)
	case Linux:
		adwMeter(c, m, frac, level, value >= hi)
	default:
		macMeter(c, m, frac, level)
	}
	return m
}

// macMeterFill returns an NSLevelIndicator's fill for a level, measured from AppKit drawn into an sRGB bitmap: systemGreen (the theme's Success), systemYellow (#ffcc00, #ffd600 in Dark Aqua) and systemRed (the theme's Danger).
func macMeterFill(t *ui.Theme, level int) ui.Color {
	switch level {
	case 2:
		return t.Danger
	case 1:
		return pick(t.Dark, ui.Hex("#ffcc00"), ui.Hex("#ffd600"))
	}
	return t.Success
}

// macMeter draws an NSLevelIndicator of continuous capacity style, regular size, as measured from AppKit: a bar 16 points high with corners of 2.5, a track of black at 7.1% (white at 10.2% in Dark Aqua) inside a half-point edge (black at 19.2%, white at 7.1%), filled from its start with the level's color inside a half-point edge 13% darker, square where the fill ends.
func macMeter(c *ui.Context, m ui.Element, frac float32, level int) {
	t := c.Theme()
	const radius = 2.5
	edge := tint(t.Background, t.Dark, pick(t.Dark, float32(0.192), 0.071))
	track := tint(t.Background, t.Dark, pick(t.Dark, float32(0.071), 0.102))
	fill := macMeterFill(t, level)
	rtl := c.Direction() == ui.RightToLeft
	m.Height(16).MinWidth(40).Radius(radius).Background(ui.Transparent).Border(0, ui.Color{})
	m.Draw(func(p *ui.Painter, r ui.Rect) {
		inner := ui.Rect{X: r.X + 0.5, Y: r.Y + 0.5, W: r.W - 1, H: r.H - 1}
		p.Fill(r, edge, radius)
		p.Fill(inner, track, radius-0.5)
		if frac <= 0 {
			return
		}
		w := r.W * frac
		bar := ui.Rect{X: r.X, Y: r.Y, W: w, H: r.H}
		if rtl {
			bar.X = r.X + r.W - w
		}
		p.Clip(r, radius, func() { p.Fill(bar, fill.Mix(black, 0.13), 0) })
		p.Clip(inner, radius-0.5, func() {
			in := ui.Rect{X: bar.X, Y: inner.Y, W: bar.W, H: inner.H}
			if frac < 1 {
				// Square where the fill ends: the edge is only along the track's.
				if rtl {
					in.W -= 0.5
				} else {
					in.X += 0.5
					in.W -= 0.5
				}
			}
			p.Fill(in, fill, 0)
		})
	})
}

// fluentMeter draws a level as WinUI draws a determinate ProgressBar, which it has instead of a meter: a 3-DIP bar with rounded ends (ProgressBarCornerRadius 1.5) over a 1-DIP ProgressBarBackground line (ControlStrongStrokeColorDefault), in ProgressBarForeground (the accent fill), ProgressBarPausedForegroundColor (SystemFillColorCaution) at the warning level or ProgressBarErrorForegroundColor (SystemFillColorCritical) at the critical one.
func fluentMeter(c *ui.Context, m ui.Element, frac float32, level int) {
	t := c.Theme()
	fill := t.Accent
	switch level {
	case 1:
		fill = pick(t.Dark, ui.Hex("#9d5d00"), ui.Hex("#fce100"))
	case 2:
		fill = pick(t.Dark, ui.Hex("#c42b1c"), ui.Hex("#ff99a4"))
	}
	line := t.Background.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0.447), 0.545))
	rtl := c.Direction() == ui.RightToLeft
	m.Height(3).Radius(0).Background(ui.Transparent).Border(0, ui.Color{})
	m.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X, Y: r.Y + 1, W: r.W, H: 1}, line, 0.5)
		w := r.W * frac
		x := r.X
		if rtl {
			x = r.X + r.W - w
		}
		p.Fill(ui.Rect{X: x, Y: r.Y, W: w, H: r.H}, fill, 1.5)
	})
}

// adwMeter draws a GtkLevelBar of continuous mode: an 8-pixel $trough_color trough with round ends, and the filled block round at both ends in accent_bg_color (the .high class), warning_bg_color (.low) at the warning level, error_bg_color at the critical one, which libadwaita has no class for (inferred), and success_bg_color (.full) once the value reaches the top, where GtkLevelBar's full offset is. The levels are MyGo's: GtkLevelBar's default offsets (0.25, 0.75 and 1 in gtklevelbar.c) would color a value .low up to a quarter and .full above three quarters, which a meter with levels of its own does not follow.
func adwMeter(c *ui.Context, m ui.Element, frac float32, level int, full bool) {
	t := c.Theme()
	fill := t.Accent
	switch {
	case level == 2:
		fill = t.Danger
	case level == 1:
		fill = t.Warning
	case full:
		fill = t.Success
	}
	rtl := c.Direction() == ui.RightToLeft
	m.Height(8).Radius(4).Background(adwTrough(t, false, false)).Border(0, ui.Color{})
	m.Draw(func(p *ui.Painter, r ui.Rect) {
		w := r.W * frac
		x := r.X
		if rtl {
			x = r.X + r.W - w
		}
		p.Fill(ui.Rect{X: x, Y: r.Y, W: w, H: r.H}, fill, r.H/2)
	})
}

// Avatar creates a picture of a person or a thing named name: the image, cropped to a circle, or initials, as the platform draws a contact without a picture: macOS's gray monogram, WinUI's PersonPicture, or libadwaita's AdwAvatar on a gradient it picks from the name; where the platform takes no initials from the name (an empty one, or on Windows one in a script PersonPicture takes none from, as Chinese), its contact glyph. It is ui.Avatar restyled, as large as MyGo's; MyGo's initials stay, transparent, under the platform's, which follow its size. It is ui.Avatar's counterpart: assistive technology sees an image named name.
func Avatar(c *ui.Context, name string, image *ui.Bitmap) ui.Element {
	t := c.Theme()
	p := PlatformOf(c)
	e := ui.Avatar(c, name, image)
	if image != nil {
		if p == Windows {
			fluentAvatarStroke(t, e)
		}
		return e
	}
	e.TextColor(ui.Transparent)
	switch p {
	case Windows:
		fluentAvatar(t, e, name)
	case Linux:
		adwAvatar(e, name)
	default:
		macAvatar(t, e, name)
	}
	return e
}

// avatarInitials returns the first letters of the first and last words of name, uppercased, as ui.Avatar takes them.
func avatarInitials(name string) string {
	var first, last rune
	inWord := false
	for _, r := range name {
		letter := unicode.IsLetter(r) || unicode.IsDigit(r)
		if letter && !inWord {
			if first == 0 {
				first = unicode.ToUpper(r)
			} else {
				last = unicode.ToUpper(r)
			}
		}
		inWord = letter
	}
	switch {
	case first == 0:
		return ""
	case last == 0:
		return string(first)
	}
	return string([]rune{first, last})
}

// avatarPerson returns a person's head and shoulders in the square of side s at x, y, the shape of a contact glyph: WinUI's Contact (U+E77B), libadwaita's avatar-default-symbolic, and the silhouette macOS shows for a contact without a name. The proportions are inferred from those glyphs.
func avatarPerson(x, y, s float32) *ui.Path {
	var path ui.Path
	path.Circle(x+s/2, y+s*0.28, s*0.22)
	path.MoveTo(x+s*0.08, y+s*0.96).
		CubeTo(x+s*0.08, y+s*0.7, x+s*0.28, y+s*0.58, x+s/2, y+s*0.58).
		CubeTo(x+s*0.72, y+s*0.58, x+s*0.92, y+s*0.7, x+s*0.92, y+s*0.96).Close()
	return &path
}

// drawAvatarText draws s centered in r at size, in weight and col.
func drawAvatarText(p *ui.Painter, r ui.Rect, s string, size float32, weight int, col ui.Color) {
	if s == "" {
		return
	}
	span := ui.Span{Text: s, Size: size, Weight: weight, Color: col}
	w, h := p.MeasureText(0, span)
	p.RichText(r.X+(r.W-w)/2, r.Y+(r.H-h)/2, 0, span)
}

// macAvatar draws a contact's monogram as macOS's Contacts shows one: white initials on a circle of systemGray, lighter at the top, or a white silhouette for a name without initials. AppKit has no avatar view, so the gradient, the initials' size and weight and the silhouette's size are inferred.
func macAvatar(t *ui.Theme, e ui.Element, name string) {
	gray := pick(t.Dark, ui.Hex("#8e8e93"), ui.Hex("#98989d")) // systemGray
	e.LinearGradient(ui.LinearGradient{Angle: 180, Stops: []ui.GradientStop{{Color: gray.Mix(white, 0.28), At: 0}, {Color: gray.Mix(black, 0.06), At: 1}}})
	in := avatarInitials(name)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		if in == "" {
			s := r.W * 0.6
			p.FillPath(avatarPerson(r.X+(r.W-s)/2, r.Y+(r.H-s)/2, s), white)
			return
		}
		drawAvatarText(p, r, in, r.W*0.42, 500, white)
	})
}

// fluentAvatarStroke draws PersonPicture's 1-DIP ellipse stroke, CardStrokeColorDefault, inside its edge.
func fluentAvatarStroke(t *ui.Theme, e ui.Element) {
	stroke := pick(t.Dark, ui.RGBA(0, 0, 0, 0.059), ui.RGBA(0, 0, 0, 0.098))
	e.DrawOver(func(p *ui.Painter, r ui.Rect) { p.Stroke(r, stroke, r.W/2, 1) })
}

// fluentAvatar draws a PersonPicture without a picture: the ellipse in PersonPictureEllipseFillThemeBrush, ControlAltFillColorQuarternary (black at 9.4%, white at 7.1% in the dark theme), with its stroke, and in TextFillColorPrimary at 42% of its width (PersonPicture.cpp) the initials semibold, or, for a name without them, the Contact glyph (U+E77B) of the NoPhotoOrInitials state, drawn here as a path.
func fluentAvatar(t *ui.Theme, e ui.Element, name string) {
	e.Background(pick(t.Dark, ui.RGBA(0, 0, 0, 0.094), ui.RGBA(255, 255, 255, 0.071)))
	fluentAvatarStroke(t, e)
	in := fluentAvatarInitials(name)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		if in == "" {
			s := r.W * 0.42
			p.StrokePath(avatarPerson(r.X+(r.W-s)/2, r.Y+(r.H-s)/2, s), 1, t.Text)
			return
		}
		drawAvatarText(p, r, in, r.W*0.42, 600, t.Text)
	})
}

// fluentAvatarInitials returns the initials WinUI's PersonPicture takes from a display name (InitialsGenerator.cpp's InitialsFromDisplayName): none unless the name's first three UTF-16 units include a Latin or Cyrillic character and none of the scripts it lists as glyphs or symbols (Arabic, the Indic scripts, Thai, Lao, CJK, Greek, Hebrew, Armenian, IPA); a final bracketed part, as in "John Smith (Work)", left out; then of the first and the last of the words between spaces, at most 25, the first character past ASCII punctuation with the combining marks after it, uppercased.
func fluentAvatarInitials(name string) string {
	units := utf16.Encode([]rune(name))
	latin := false
	for _, u := range units[:min(len(units), 3)] {
		switch fluentAvatarCharKind(rune(u)) {
		case 2:
			return ""
		case 1:
			latin = true
		}
	}
	if !latin {
		return ""
	}
	for _, d := range []string{"{}", "()", "[]"} {
		if strings.HasSuffix(name, d[1:]) {
			if i := strings.LastIndex(name, d[:1]); i >= 0 {
				name = name[:i]
				break
			}
		}
	}
	parts := strings.Split(name, " ")
	var words []string
	for _, w := range parts[:min(len(parts), 25)] {
		if w != "" {
			words = append(words, w)
		}
	}
	first := func(w string) string {
		rs := []rune(w)
		start := 0
		for start < len(rs) && (rs[start] >= 0x21 && rs[start] <= 0x2f || rs[start] >= 0x3a && rs[start] <= 0x40 || rs[start] >= 0x7b && rs[start] <= 0x7e) {
			start++
		}
		if start == len(rs) {
			start = 0
		}
		end := start + 1
		for end < len(rs) && rs[end] >= 0x300 && rs[end] <= 0x36f {
			end++
		}
		return string(rs[start:end])
	}
	switch len(words) {
	case 0:
		return ""
	case 1:
		return strings.ToUpper(first(words[0]))
	}
	return strings.ToUpper(first(words[0]) + first(words[len(words)-1]))
}

// fluentAvatarCharKind returns 1 for a UTF-16 unit InitialsGenerator.cpp counts as Latin (Basic Latin, Latin-1, the Latin Extended blocks, Cyrillic and its supplement, combining diacritical marks), 2 for one of the blocks it counts as glyphs or symbols, and 0 for any other.
func fluentAvatarCharKind(u rune) int {
	for _, b := range [][2]rune{
		{0x0250, 0x02af}, {0x0600, 0x06ff}, {0x0750, 0x077f}, {0x08a0, 0x08ff}, {0xfb50, 0xfdff}, {0xfe70, 0xfeff}, {0x0900, 0x097f}, {0xa8e0, 0xa8ff},
		{0x0980, 0x0dff}, {0x0e00, 0x0eff}, {0x4e00, 0x9fff}, {0x3400, 0x4dbf}, {0x2e80, 0x2eff}, {0x3000, 0x303f}, {0x31c0, 0x31ef}, {0x3200, 0x33ff},
		{0xf900, 0xfaff}, {0xfe30, 0xfe4f}, {0x0370, 0x03ff}, {0x0590, 0x05ff}, {0x0530, 0x058f},
	} {
		if u >= b[0] && u <= b[1] {
			return 2
		}
	}
	for _, b := range [][2]rune{{0x0001, 0x024f}, {0x2c60, 0x2c7f}, {0xa720, 0xa7ff}, {0xab30, 0xab6f}, {0x1e00, 0x1eff}, {0x0400, 0x052f}, {0x0300, 0x036f}} {
		if u >= b[0] && u <= b[1] {
			return 1
		}
	}
	return 0
}

// adwAvatarColors are _avatar.scss's $avatarcolorlist: the initials' color and the gradient's top and bottom, in the order of the color1 to color14 classes.
var adwAvatarColors = [14][3]ui.Color{
	{ui.Hex("#cfe1f5"), ui.Hex("#83b6ec"), ui.Hex("#337fdc")}, // blue
	{ui.Hex("#caeaf2"), ui.Hex("#7ad9f1"), ui.Hex("#0f9ac8")}, // cyan
	{ui.Hex("#cef8d8"), ui.Hex("#8de6b1"), ui.Hex("#29ae74")}, // green
	{ui.Hex("#e6f9d7"), ui.Hex("#b5e98a"), ui.Hex("#6ab85b")}, // lime
	{ui.Hex("#f9f4e1"), ui.Hex("#f8e359"), ui.Hex("#d29d09")}, // yellow
	{ui.Hex("#ffead1"), ui.Hex("#ffcb62"), ui.Hex("#d68400")}, // gold
	{ui.Hex("#ffe5c5"), ui.Hex("#ffa95a"), ui.Hex("#ed5b00")}, // orange
	{ui.Hex("#f8d2ce"), ui.Hex("#f78773"), ui.Hex("#e62d42")}, // raspberry
	{ui.Hex("#fac7de"), ui.Hex("#e973ab"), ui.Hex("#e33b6a")}, // magenta
	{ui.Hex("#e7c2e8"), ui.Hex("#cb78d4"), ui.Hex("#9945b5")}, // purple
	{ui.Hex("#d5d2f5"), ui.Hex("#9e91e8"), ui.Hex("#7a59ca")}, // violet
	{ui.Hex("#f2eade"), ui.Hex("#e3cf9c"), ui.Hex("#b08952")}, // beige
	{ui.Hex("#e5d6ca"), ui.Hex("#be916d"), ui.Hex("#785336")}, // brown
	{ui.Hex("#d8d7d3"), ui.Hex("#c0bfbc"), ui.Hex("#6e6d71")}, // gray
}

// adwAvatarColor returns the index into adwAvatarColors that adw-avatar.c picks for text: GLib's g_str_hash (djb2 over signed bytes) modulo the 14 colors. An empty name gets a random color in libadwaita; here the first.
func adwAvatarColor(text string) int {
	if text == "" {
		return 0
	}
	h := uint32(5381)
	for i := 0; i < len(text); i++ {
		h = h<<5 + h + uint32(int32(int8(text[i])))
	}
	return int(h % uint32(len(adwAvatarColors)))
}

// adwAvatarInitials returns the initials adw-avatar.c takes from text: uppercased and trimmed, its first character and the one after its last space.
func adwAvatarInitials(text string) string {
	s := []rune(strings.TrimSpace(strings.ToUpper(text)))
	if len(s) == 0 {
		return ""
	}
	out := []rune{s[0]}
	for i := len(s) - 1; i > 0; i-- {
		if s[i] == ' ' {
			if i+1 < len(s) {
				out = append(out, s[i+1])
			}
			break
		}
	}
	return string(out)
}

// adwAvatar draws an AdwAvatar without an image: the gradient of the color its name hashes to, from top to bottom, and the initials bold in the color's light tint, as large as fits the biggest square inside the circle less a padding of 0.4 times the size minus 5 pixels (adw-avatar.c's update_font_size); white below 25 pixels, the .contrasted class adw_avatar_set_size adds. Without a name it shows avatar-default-symbolic at half its size instead (update_visibility and adw_avatar_set_size).
func adwAvatar(e ui.Element, name string) {
	col := adwAvatarColors[adwAvatarColor(name)]
	e.LinearGradient(ui.LinearGradient{Angle: 180, Stops: []ui.GradientStop{{Color: col[1], At: 0}, {Color: col[2], At: 1}}})
	in := adwAvatarInitials(name)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		fg := col[0]
		if r.W < 25 {
			fg = white
		}
		if name == "" {
			s := float32(int(r.W) / 2)
			p.FillPath(avatarPerson(r.X+(r.W-s)/2, r.Y+(r.H-s)/2, s), fg)
			return
		}
		if in == "" {
			return
		}
		maxSize := r.W/1.4142 - max(r.W*0.4-5, 0)
		w, h := p.MeasureText(0, ui.Span{Text: in, Size: maxSize, Weight: 700})
		size := maxSize
		if w > 0 {
			size = min(max(h*(maxSize/w), 0), maxSize)
		}
		drawAvatarText(p, r, in, size, 700, fg)
	})
}

// ratingLook is how a platform draws a rating's stars: each box, the gap between them, and how far the star's inner points are from its center, as a share of the outer ones.
type ratingLook struct {
	box, gap, inner float32
}

// Rating creates a row of max stars of which *value are filled, as the platform draws a rating: AppKit's NSLevelIndicator of rating style, WinUI's RatingControl, and on GNOME, which has no rating control, GNOME Software's stars. A click on a star sets it, a click on the star set clears it, and the arrows, Home and End change it while it has the focus; Changed reports a new value. It is ui.Rating's counterpart, built in a column of its own: MyGo's stars are left transparent there, sized to the platform's, and the platform's are drawn over them, so the pointer, the keys, the focus and what assistive technology sees (a slider from 0 to max) are MyGo's. It returns ui.Rating's row, inside that column: Label, Disabled and Changed work on it as on ui.Rating's, but the column, not the row, takes the place in the parent's layout, so give a margin or a grow to an element around it.
func Rating(c *ui.Context, value *int, max int) ui.Element {
	p := PlatformOf(c)
	look := ratingLook{box: 11, gap: 1, inner: 0.4} // NSLevelIndicator: 12-point stars, measured, side by side
	switch p {
	case Windows:
		look = ratingLook{box: 16, gap: 8, inner: 0.45} // RatingControlFontSizeForRendering 32 drawn at half, RatingControlItemSpacing 8
	case Linux:
		look = ratingLook{box: 32, gap: 2, inner: 0.45} // gs-review-dialog.ui's icon-size 32, gs-star-widget.c's STAR_SPACING 2
	}
	// ui.Rating sizes its stars at 1.25 times the font size, half a unit of spacing apart, and paints them in Warning and TextMuted.
	th := *c.Theme()
	th.FontSize, th.Spacing = look.box/1.25, look.gap*2
	th.Warning, th.TextMuted = ui.Transparent, ui.Transparent
	var r ui.Element
	ui.Themed(c, &th).Background(ui.Transparent).AlignItems(ui.Start).Children(func() { r = ui.Rating(c, value, max) })
	t := c.Theme()
	// Under the pointer, unless the rating or what holds it is disabled, as ui.Rating shows its own hover.
	px, _, over := r.PointerPosition()
	over = over && r.Hovered()
	n := max
	rtl := c.Direction() == ui.RightToLeft
	// at returns the box of star i: from the start of the reading direction, as ui.Rating lays them out.
	at := func(rect ui.Rect, i int) ui.Rect {
		x := rect.X + float32(i)*(look.box+look.gap)
		if rtl {
			x = rect.X + rect.W - look.box - float32(i)*(look.box+look.gap)
		}
		return ui.Rect{X: x, Y: rect.Y + (rect.H-look.box)/2, W: look.box, H: look.box}
	}
	r.Draw(func(pt *ui.Painter, rect ui.Rect) {
		v := *value
		hover := -1
		if over && p == Windows {
			for i := range n {
				if b := at(rect, i); rect.X+px >= b.X-look.gap/2 && rect.X+px < b.X+b.W+look.gap/2 {
					hover = i
				}
			}
		}
		for i := range n {
			b := at(rect, i)
			switch p {
			case Windows:
				fluentRatingStar(pt, t, b, i < v, v == 0, hover, i, over, rect.X+px)
			case Linux:
				adwRatingStar(pt, t, b, look.inner, i < v)
			default:
				macRatingStar(pt, t, b, look.inner, i < v)
			}
		}
	})
	return r
}

// ratingStarPath returns a five-pointed star of outer radius rad around cx, cy, its inner points at inner times rad.
func ratingStarPath(cx, cy, rad, inner float32) *ui.Path {
	var path ui.Path
	for i := range 10 {
		d := rad
		if i%2 == 1 {
			d = rad * inner
		}
		a := float64(i)*math.Pi/5 - math.Pi/2
		x, y := cx+d*float32(math.Cos(a)), cy+d*float32(math.Sin(a))
		if i == 0 {
			path.MoveTo(x, y)
		} else {
			path.LineTo(x, y)
		}
	}
	path.Close()
	return &path
}

// macRatingStar draws a star of NSLevelIndicator's rating style, measured from AppKit: a 12-point star filled in secondaryLabelColor while set, and in systemFill (black or white at 9.8%) while not.
func macRatingStar(p *ui.Painter, t *ui.Theme, b ui.Rect, inner float32, set bool) {
	fill := tint(t.Background, t.Dark, macTroughAlpha)
	if set {
		fill = t.TextMuted
	}
	p.FillPath(ratingStarPath(b.X+b.W/2, b.Y+b.H/2+0.4, 6.2, inner), fill)
}

// fluentRatingStar draws a star of a RatingControl, whose outlined glyph (U+E734) in RatingControlUnselectedForeground, TextFillColorSecondary, lies under every star, and over it the filled glyph (U+E735) of the stars set in RatingControlSelectedForeground, the accent fill. Under the pointer the filled glyph runs up to the star it is on: in the accent while a rating is set (the PointerOverSet state), in RatingControlPointerOverPlaceholderForeground, ControlAltFillColorTertiary (black at 5.9%, white at 4.3% in the dark theme), while none is (PointerOverPlaceholder, as unset reads 0 here). The stars grow toward the pointer as RatingControl.cpp's scale animation does: max(0.8 − 0.0004·dx², 0.5) of the 32-DIP glyph, dx the distance from the star's center to the pointer.
func fluentRatingStar(p *ui.Painter, t *ui.Theme, b ui.Rect, set, unset bool, hover, i int, over bool, px float32) {
	scale := float32(0.5)
	cx := b.X + b.W/2
	if over {
		dx := cx - px
		scale = max(-0.0005*0.8*dx*dx+0.8, 0.5)
	}
	rad := 32 * scale / 2 * 0.94
	star := ratingStarPath(cx, b.Y+b.H/2+rad*0.08, rad, 0.45)
	if hover >= 0 {
		set = i <= hover
	}
	if set && !unset {
		p.FillPath(star, t.Accent)
		return
	}
	p.StrokePath(star, 1, t.TextMuted)
	if set {
		// Only under the pointer: an unset rating has no star set otherwise.
		p.FillPath(star, pick(t.Dark, ui.RGBA(0, 0, 0, 0.059), ui.RGBA(255, 255, 255, 0.043)))
	}
}

// adwRatingStar draws a star of GNOME Software's GsStarWidget: starred-symbolic, in @yellow_5 (#e5a50a) while set (star-image > image.starred) and at 20% opacity while not (image.non-starred), in star-image's color: the text color in the light theme and @yellow_1 (#f9f06b) in the dark one (style.css of 51.0, its prefers-color-scheme: dark block; style-dark.css of 48.0 had the same).
func adwRatingStar(p *ui.Painter, t *ui.Theme, b ui.Rect, inner float32, set bool) {
	fill := t.Background.Mix(pick(t.Dark, t.Text, ui.Hex("#f9f06b")), 0.2)
	if set {
		fill = ui.Hex("#e5a50a")
	}
	rad := b.W / 2 * 0.92
	p.FillPath(ratingStarPath(b.X+b.W/2, b.Y+b.H/2+rad*0.08, rad, inner), fill)
}

// FindBar creates a bar for finding text while *open, in the platform's place for it: NSTextFinder's find bar on macOS, a strip with a divider under it on Windows, which has no find bar, and GtkSearchBar's bar on GNOME. The field, the count, the previous and next buttons, Done, the keys (Enter, Shift+Enter, Cmd+G or F3, Escape) and what assistive technology hears are MyGo's, in the platform's theme: the bar is ui.FindBar's, restyled. While closed, it returns an element showing nothing. It is ui.FindBar's counterpart.
func FindBar(c *ui.Context, open *bool, query *string, matches int, current *int) ui.Element {
	shown := *open
	bar := ui.FindBar(c, open, query, matches, current)
	if !shown {
		return bar
	}
	switch PlatformOf(c) {
	case Windows:
		fluentFindBar(c, bar)
	case Linux:
		adwFindBar(c, bar)
	default:
		macFindBar(c, bar)
	}
	return bar
}

// macFindBar draws NSTextFinder's find bar, as measured from AppKit: 32 points high, the window's background, its controls 7 points from the edge and 9 to 11 apart, and along its bottom a 6-point shade (black or white from 0.4% to 5.9%) over a line of 13.3%.
func macFindBar(c *ui.Context, bar ui.Element) {
	t := c.Theme()
	bar.MinHeight(32).PaddingX(7).Gap(10).Background(t.Background).BorderWidth(0, 0, 1, 0).BorderColor(tint(t.Background, t.Dark, 0.133))
	from, to := tint(t.Background, t.Dark, 0.004), tint(t.Background, t.Dark, 0.059)
	bar.Draw(func(p *ui.Painter, r ui.Rect) {
		shade := ui.Rect{X: r.X, Y: r.Y + r.H - 7, W: r.W, H: 6}
		p.FillGradient(shade, ui.LinearGradient{Angle: 180, Stops: []ui.GradientStop{{Color: from, At: 0}, {Color: to, At: 1}}}, 0)
	})
}

// fluentFindBar lays out the find bar on Windows, which has no find bar control: the window's background with a DividerStrokeColorDefault line under it, padded and spaced as a CommandBar's content, 8 DIPs. The padding is inferred.
func fluentFindBar(c *ui.Context, bar ui.Element) {
	t := c.Theme()
	bar.MinHeight(48).Padding(8, 12).Gap(8).Background(t.Background).BorderWidth(0, 0, 1, 0).BorderColor(fluentDivider(t))
}

// adwFindBar draws a GtkSearchBar's bar (searchbar > revealer > box): padding 6px 6px 7px 6px, headerbar_bg_color (white, #2e2e32 in the dark theme) and an inset 1-pixel headerbar_shade_color along its bottom (the shade at 12%, 36% in the dark theme).
func adwFindBar(c *ui.Context, bar ui.Element) {
	t := c.Theme()
	bg := pick(t.Dark, white, ui.Hex("#2e2e32"))
	bar.Padding(6, 6, 7, 6).Gap(6).Background(bg).BorderWidth(0, 0, 1, 0).BorderColor(bg.Mix(adwInk, pick(t.Dark, float32(0.12), 0.36)))
}
