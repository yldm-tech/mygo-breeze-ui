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

// InfoBar creates an inline message bar of a severity, kind, with a title and a message, as the platform shows a message inside a window: WinUI's InfoBar on Windows, libadwaita's AdwBanner on GNOME, and on macOS, which has no message bar, a section of a grouped form with a colored SF Symbol, as System Settings shows its inline notices (inferred). kind is a ui.Toast's Type, so that the same word colors a toast and a bar: "success", "warning", "error" ("danger" and "critical" too), and any other, as "" or "info", is informational; libadwaita's banners have no severity, so on GNOME every kind draws the same banner. A nil open makes a bar that cannot be closed; otherwise the bar shows while *open and its close button sets *open to false: WinUI's close button, a borderless xmark button on macOS, and on GNOME the banner's button labelled Dismiss. The bar is a status for assistive technology, named by its title and described by its message; the close button is named "Close" on macOS and Windows and "Dismiss" on GNOME. It stretches across a column, as the platforms' bars stretch. It is InfoBarActions without actions.
func InfoBar(c *ui.Context, kind, title, message string, open *bool) ui.Element {
	return InfoBarActions(c, kind, title, message, open, nil)
}

// InfoBarActions creates an InfoBar with actions: actions builds buttons, as Button, PrimaryButton or Link, which the bar shows where the platform puts a message bar's action. On Windows they are WinUI InfoBar's ActionButton, 16 DIPs after the message on its line while the title, the message and the actions fit across the bar, else 12 below the message at the start. On macOS they sit at the end of the bar before the close button, centered on the title's line, as SwiftUI lays a button out in that row of a grouped form (measured on macOS 27.2): a 24-point push button makes the row 4 points taller. On GNOME they take the place of the banner's button, as AdwBanner has one button (button-label, a PrimaryButton showing as its suggested action): the banner then has no Dismiss button, and an action closes it by setting *open to false. The buttons take their faces over the bar's fill, as the platforms lay their translucent faces over it. Assistive technology reads the title, the message, the actions and then the close button. A nil actions makes InfoBar.
func InfoBarActions(c *ui.Context, kind, title, message string, open *bool, actions func()) ui.Element {
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
		fluentInfoBar(c, bar, severity, title, message, closer, actions)
	case Linux:
		if actions != nil {
			closer = nil
		}
		adwBanner(c, bar, title, message, closer, actions)
	default:
		macInfoBar(c, bar, severity, title, message, closer, actions)
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

// macInfoBar draws the message bar on macOS, which has none: a row of a SwiftUI grouped Form, as measured on macOS 27.2: the section's fill (macFormSectionFill) with corners of 12, the row padded 10 points above and below and 11.5 at the start; a 13-point SF Symbol of the severity (info.circle.fill, checkmark.circle.fill, exclamationmark.triangle.fill or xmark.octagon.fill) in its color, the mark inside cut out to the fill, centered on the title's first line; 11.5 points after it the title in 13-point bold labelColor and the message in secondaryLabelColor, both over the fill (38 and 124 over 247, 222 and 157 over 37), 2 points apart (inferred); and at the end a borderless button showing xmark, 10.5 points across in a 16-point box, its ink 12 points from the end, in secondaryLabelColor, labelColor while pressed (inferred). That macOS shows such a row as a message bar is inferred. The actions, built over the fill (macOver), sit between the text and the close button, as SwiftUI lays Buttons out after a Spacer in that row (measured by tools/appkit-reference's statusActionsSet on macOS 27.2): their middle on the title line's, the row growing where they are taller than the line, so that a 24-point push button is 10 points from the row's top and the symbol, the text and the close button 4 points lower than without it; 6.25 points before the close button's box, or 10 from the row's end without one; and 8 apart and at least 8 after the text, as HStack spaces them (inferred: a message that fits beside a button in SwiftUI's render leaves 10).
func macInfoBar(c *ui.Context, bar ui.Element, severity int, title, message string, closer func(string) ui.Element, actions func()) {
	t := c.Theme()
	fill := macFormSectionFill(t)
	secondary := tint(fill, t.Dark, pick(t.Dark, float32(0.498), 0.549))
	line := lineHeight(c, t.FontSize, 700)
	bar.AlignItems(ui.Start).Gap(11.5).PaddingY(10).PaddingStart(11.5).PaddingEnd(9.25).Radius(12).Background(fill).TextColor(tint(fill, t.Dark, 0.847))
	bar.Children(func() {
		color := macInfoBarColor(t, severity)
		symbol := ui.Box(c).Size(13, line).Shrink(0).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			macInfoBarSymbol(p, ui.Rect{X: r.X, Y: r.Y + (r.H-13)/2, W: 13, H: 13}, severity, color, fill)
		})
		text := ui.Column(c).Grow(1).Shrink(1).Gap(2).Children(func() {
			ui.Text(c, title).Bold()
			if message != "" {
				ui.Text(c, message).TextColor(secondary)
			}
		})
		// How far the title's line moves down to keep its middle on the actions' middle, where they are taller than it.
		var drop float32
		if actions != nil {
			row := ui.Themed(c, macOver(t, fill)).Background(ui.Transparent).Row().Shrink(0).MinHeight(line).AlignItems(ui.Center).Gap(8).Children(actions)
			// The text keeps 8 before the actions, HStack's spacing, so that a message wraps where SwiftUI's does (the bar's gap of 11.5 less 3.5). The bar's gap less 5.25 leaves the 6.25 before the close button's box; without one the bar's end padding of 9.25 and 0.75 leave 10.
			row.MarginStart(-3.5)
			if closer != nil {
				row.MarginEnd(-5.25)
			} else {
				row.MarginEnd(0.75)
			}
			if _, h := prevSize(c, row); h > line {
				drop = (h - line) / 2
				symbol.Margin(drop, 0, 0, 0)
				text.Margin(drop, 0, 0, 0)
			}
		}
		if closer == nil {
			return
		}
		b := closer("Close")
		glyph := secondary
		if b.Pressed() {
			glyph = tint(fill, t.Dark, 0.847)
		}
		b.Size(16, line).Margin(drop, 0, 0, 0).Center().Children(func() {
			ui.Box(c).Size(16, 16).Draw(func(p *ui.Painter, r ui.Rect) { toastCross(p, r, 0.17, 1.3, glyph) })
		})
		macRing(b, 4)
	})
}

// macOver returns t with the colors macOS lays over the window laid over under instead: the theme of the controls of a message bar, whose push buttons' secondarySystemFill and labels show the bar's fill through them (a push button's face reads 228 over the section's 247 and 54 over 37, where over the window it reads 235 and 48, measured by tools/appkit-reference's statusActionsSet on macOS 27.2).
func macOver(t *ui.Theme, under ui.Color) *ui.Theme {
	o := *t
	o.Background = under
	o.Text = tint(under, t.Dark, 0.847)
	o.TextMuted = tint(under, t.Dark, pick(t.Dark, float32(0.498), 0.549))
	o.Border = macSeparator(under, t.Dark)
	o.Surface = macFace(under, t.Dark)
	o.SurfaceHover = o.Surface
	o.SurfacePressed = macPressed(under, t.Dark)
	return &o
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

// fluentInfoBar draws a WinUI InfoBar in the layout, as InfoBar.xaml lays it out and as fluentToast draws it floating: at least 48 DIPs high inside its border (InfoBarMinHeight on the template's Grid, inside the ContentRoot Border), so 50 here, where MinHeight includes the border, ControlCornerRadius, a 1-DIP CardStrokeColorDefault border (InfoBarBorderThickness) laid over the window, as the background stops inside it (measured from WinUI 3 3.3 on Windows 11 26H2 by tools/winui-reference: the edge reads 229 around every severity's bar in the light theme and 29 in the dark one, not tinted by the bar's color), InfoBarContentRootPadding 16 at the start; the severity icon (fluentInfoBarIcon) in a 16-DIP box with InfoBarIconMargin 0,16,14,16; the title in 14-DIP semibold and the message at regular weight, 16 before the close button (InfoBarPanelMargin); and the close button (fluentInfoBarClose). InfoBarPanel puts the message 12 after the title, both 14 from the top (InfoBarTitleHorizontalOrientationMargin, InfoBarMessageHorizontalOrientationMargin), the panel as high as its tallest child and the bar's MinHeight making the rest, while the two fit on one line, and otherwise, or when there is no message (InfoBarPanel::MeasureOverride lays a single child out vertically), the message 4 below the title, the panel padded 14 above and 18 below (InfoBarPanelVerticalOrientationPadding, InfoBarMessageVerticalOrientationMargin); Breeze UI decides it from the bar's width in the previous frame. In the layout the informational bar's background is SystemFillColorAttentionBackground, #80F6F6F6 and #08FFFFFF over the window, where the floating toast's is the solid one; the other severities' are the same solid colors (fluentToastColors). It has no shadow. The ActionButton is the panel's third child (InfoBar.xaml): across, 16 after the message and 8 from the top (InfoBarActionHorizontalOrientationMargin 16,8,0,0), and laid out vertically, 12 below the message (InfoBarActionVerticalOrientationMargin); the panel goes vertical when the title, the message and the actions with their margins are wider than it or the actions with their margin are taller than InfoBarMinHeight (InfoBarPanel::MeasureOverride), which Breeze UI decides from the actions' size in the previous frame. The actions are built over the bar's color (fluentOver), in a row at least 32 high, a Button's height, so that a Link lines up with the title as InfoBar's HyperlinkButton does, and 8 apart (inferred).
func fluentInfoBar(c *ui.Context, bar ui.Element, severity int, title, message string, closer func(string) ui.Element, actions func()) {
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
	bar.AlignItems(ui.Start).MinHeight(50).PaddingStart(16).Radius(4).Background(bg).
		Border(1, t.Background.Mix(black, pick(t.Dark, float32(0x0f)/255, 0x19/255.0))).TextColor(t.Text).FontSize(14)
	bar.Children(func() {
		ui.Box(c).Size(16, 16).Shrink(0).MarginY(16).MarginEnd(14).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			fluentInfoBarIcon(p, r, severity, iconFill, pick(t.Dark, white, ui.RGBA(0, 0, 0, 0.894)))
		})
		box := ui.Row(c).Grow(1).Shrink(1).AlignItems(ui.Start).MarginEnd(16)
		var head, body, act ui.Element
		box.Children(func() {
			head = ui.Text(c, title).FontWeight(600)
			if message != "" {
				body = ui.Text(c, message)
			}
			if actions != nil {
				act = ui.Themed(c, fluentOver(t, bg)).Background(ui.Transparent).Row().Shrink(0).MinHeight(32).AlignItems(ui.Center).Gap(8).Children(actions)
			}
		})
		// InfoBarPanel::MeasureOverride lays a single item, a title without a message or actions, out vertically too.
		items, across, tall := 1, tw, false
		if message != "" {
			items, across = items+1, across+12+mw
		}
		if actions != nil {
			aw, ah := prevSize(c, act)
			items, across, tall = items+1, across+16+aw, 8+ah > 48
		}
		if items > 1 && across <= panel && !tall {
			head.Margin(14, 0, 0, 0)
			body.Shrink(1).Margin(14, 0, 0, 0).MarginStart(12)
			act.Margin(8, 0, 0, 0).MarginStart(16)
		} else {
			box.Column().Padding(14, 0, 18, 0)
			body.Margin(4, 0, 0, 0)
			act.Margin(12, 0, 0, 0)
		}
		if closer != nil {
			fluentInfoBarClose(c, closer("Close"))
		}
	})
}

// fluentInfoBarIcon draws an InfoBar's severity icon in its 16-DIP box r, as WinUI 3 3.3 renders the glyphs of Segoe Fluent Icons on Windows 11 26H2 (measured by tools/winui-reference at 2x): a circle 15 DIPs across, flush with the box's start and a DIP below its top, in fill, and on it in mark, in strokes a DIP wide: for information an i, its dot 1.5 high from 3.5 above the circle's middle and its stem 3.5 high from half a DIP above it; for a warning the same turned over, its stem from 3 above the middle and its dot ending 3.5 below; for success a check mark from 3.25 left of the middle to its point 1 left and 2 below and up to 3.25 right and 2.25 above; for an error an × whose ink is 5 DIPs square. The marks are drawn from the render, as the font is not open.
func fluentInfoBarIcon(p *ui.Painter, r ui.Rect, severity int, fill, mark ui.Color) {
	circle := ui.Rect{X: r.X, Y: r.Y + 1, W: 15, H: 15}
	p.Fill(circle, fill, 7.5)
	cx, cy := circle.X+7.5, circle.Y+7.5
	switch severity {
	case fluentToastSuccess:
		var path ui.Path
		path.MoveTo(cx-3.25, cy-0.25).LineTo(cx-1, cy+2).LineTo(cx+3.25, cy-2.25)
		p.StrokePath(&path, 1, mark)
	case fluentToastError:
		toastCross(p, ui.Rect{X: cx - 2.25, Y: cy - 2.25, W: 4.5, H: 4.5}, 0, 1, mark)
	case fluentToastWarning:
		p.Fill(ui.Rect{X: cx - 0.5, Y: cy - 3, W: 1, H: 3.5}, mark, 0)
		p.Fill(ui.Rect{X: cx - 0.5, Y: cy + 2, W: 1, H: 1.5}, mark, 0)
	default:
		p.Fill(ui.Rect{X: cx - 0.5, Y: cy - 3.5, W: 1, H: 1.5}, mark, 0)
		p.Fill(ui.Rect{X: cx - 0.5, Y: cy - 0.5, W: 1, H: 3.5}, mark, 0)
	}
}

// prevWidth returns e's width in the previous frame, as a platform's layout that switches between arrangements by the width it is given knows it; for an element that frame lacked it returns the window's width and asks for another frame, which lays it out at its own.
func prevWidth(c *ui.Context, e ui.Element) float32 {
	// LayoutBounds is the box before transforms and zoom, in the DIPs the layout works in. A bar the previous frame had is never empty, as its padding and MinHeight make it high, so an empty box means a new bar and the frame is asked for once, not on every frame for a bar laid out zero wide.
	if b := e.LayoutBounds(); b.W > 0 || b.H > 0 {
		return b.W
	}
	c.Invalidate()
	w, _ := c.Size()
	return w
}

// prevSize returns e's size in the previous frame, as a layout that places parts by the size of what the caller builds in them needs it; for an element that frame lacked it returns 0, 0 and asks once for another frame, which lays it out by its own. An element laid out empty is asked for only once, not on every frame.
func prevSize(c *ui.Context, e ui.Element) (float32, float32) {
	b := e.LayoutBounds()
	if b.W == 0 && b.H == 0 {
		asked := ui.Local(e, "breeze-prev-size", func() bool { return false })
		if !*asked {
			*asked = true
			c.Invalidate()
		}
	}
	return b.W, b.H
}

// fluentInfoBarClose draws an InfoBar's close button (InfoBarCloseButtonStyle): 38 DIPs square (InfoBarCloseButtonSize) with a margin of 5, at the top, ControlCornerRadius, drawn with AppBarButton's brushes as InfoBar.xaml sets them: transparent, SubtleFillColorSecondary under the pointer and SubtleFillColorTertiary pressed; its Cancel symbol (InfoBarCloseButtonGlyphSize) in the text color, an × whose ink is 11 DIPs square in 1-DIP strokes, centered half a DIP above and before the button's middle (measured from WinUI 3 3.3 on Windows 11 26H2 by tools/winui-reference at 2x); and, as InfoBarCloseButtonStyle is based on DefaultButtonStyle, the system focus visual 3 DIPs around it (FocusVisualMargin -3, fluentFocusRing).
func fluentInfoBarClose(c *ui.Context, b ui.Element) {
	t := c.Theme()
	b.Size(38, 38).Margin(5).Shrink(0).Radius(4).Center().FocusRing(false)
	focused := b.FocusVisible()
	b.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if focused {
			fluentFocusRing(p, t, r, 3, 3, 3, 3, 4)
		}
	})
	hover, pressed := fluentSubtle(t)
	switch {
	case b.Pressed():
		b.Background(pressed)
	case b.Hovered():
		b.Background(hover)
	}
	b.Children(func() {
		ui.Box(c).Size(16, 16).Draw(func(p *ui.Painter, r ui.Rect) {
			toastCross(p, ui.Rect{X: r.X + r.W/2 - 5.5, Y: r.Y + r.H/2 - 5.5, W: 10, H: 10}, 0, 1, t.Text)
		})
	})
}

// adwBannerShares are the shares of currentColor a GtkButton's face lays over what is behind it, at rest, under the pointer and pressed, in the light theme (where currentColor is window_fg_color, 80% opaque) and the dark one: $button_color and its hover and active shades.
var adwBannerShares = [2][3]float32{{0.08, 0.12, 0.24}, {0.10, 0.15, 0.30}}

// adwBanner draws an AdwBanner: a full-width strip of --banner-color (#7d7d83) mixed 30% into window_bg_color, in window_fg_color, padded 6 (_toolbars.scss); its title in .heading, bold, wrapping, centered across its lines; and, while it can be closed, its button, a GtkButton at least 84 wide (BUTTON_HORZ_MIN_WIDTH), as adw-banner.c's allocate_content places them: without a button the title centered and padded a further 6 all around (HORZ_PADDING, VERT_PADDING); with one, the button at the end, centered on the strip's height, and the title centered across the whole strip while it fits with 36 to spare on each side of the button (HORZ_SPACING_CENTERED), else 6 from the start (HORZ_PADDING) and 6 before the button (HORZ_SPACING); and where the title does not fit beside the button at its natural width, the button below the title, centered, at least 160 wide (BUTTON_VERT_MIN_WIDTH), 9 below it (VERT_SPACING), the two padded 6 above, below and at the sides. Breeze UI picks the arrangement from the banner's width in the previous frame. A banner has only a title: the message follows it in the same label, at regular weight (inferred). libadwaita's banners have no severity; the .error, .warning and .success classes only recolor accent-colored parts, which a banner's title is not, so every kind draws this same banner. Actions take the button's place, built over the strip (adwOverBanner) and as wide as they are, 6 apart (inferred), placed by their width in the previous frame.
func adwBanner(c *ui.Context, bar ui.Element, title, message string, closer func(string) ui.Element, actions func()) {
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
	if closer == nil && actions == nil {
		bar.Justify(ui.Center).AlignItems(ui.Center).Children(func() { text().Padding(6) })
		return
	}
	const dismiss = "Dismiss"
	tw, _ := c.MeasureText(0, ui.Span{Text: title, Size: t.FontSize, Weight: 700})
	if message != "" {
		mw, _ := c.MeasureText(0, ui.Span{Text: " " + message, Size: t.FontSize})
		tw += mw
	}
	width := prevWidth(c, bar) - 12 // the strip's padding
	bar.Children(func() {
		label := text()
		var button ui.Element
		var bw float32
		if actions != nil {
			button = ui.Themed(c, adwOverBanner(t, bg, fg)).Background(ui.Transparent).Row().Shrink(0).AlignItems(ui.Center).Gap(6).Children(actions)
			bw, _ = prevSize(c, button)
		} else {
			bw, _ = c.MeasureText(0, ui.Span{Text: dismiss, Size: t.FontSize, Weight: 700})
			bw += 34 // the button's padding of 17 on each side
			button = closer(dismiss)
			adwButton(c, button, false)
			shares := pick(t.Dark, adwBannerShares[0], adwBannerShares[1])
			share := shares[0]
			switch {
			case button.Pressed():
				share = shares[2]
			case button.Hovered():
				share = shares[1]
			}
			face := bg.Mix(pick(t.Dark, adwInk, white), share)
			button.Background(face).TextColor(adwTextOn(t, face)).Children(func() { ui.Text(c, dismiss).SingleLine() })
		}
		// The minimum widths are the Dismiss button's; actions keep their own, so the arrangement and the title's margin go by the width the actions have.
		least := func(w float32) {
			if actions == nil {
				button.MinWidth(w)
			}
		}
		shown := bw
		if actions == nil {
			shown = max(bw, 84)
		}
		switch {
		case 6+tw+6+shown > width:
			bar.Column().AlignItems(ui.Center).Gap(9).Padding(12)
			least(min(max(bw, 160), width))
		case tw+(shown+36)*2 > width:
			bar.AlignItems(ui.Center).Gap(6).PaddingStart(12)
			label.Grow(1).TextAlign(ui.Start)
			least(84)
		default:
			// A margin as wide as the button and the gap keeps the title centered across the strip.
			bar.AlignItems(ui.Center).Gap(6)
			label.Grow(1).MarginStart(shown + 6)
			least(84)
		}
	})
}

// adwOverBanner returns t laid over a banner's strip, bg, its text in fg: the theme a banner's actions build in, so that a button's translucent face and its label take their colors over the strip, as the Dismiss button's do.
func adwOverBanner(t *ui.Theme, bg, fg ui.Color) *ui.Theme {
	o := *t
	o.Background, o.Text = bg, fg
	return &o
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

// macRingLight, macRingArc and macRingFull are the arc of a determinate spinning NSProgressIndicator for the default accent (#007aff), measured over the window: #0a82ff in Aqua, full or not, and in Dark Aqua #1769e6, and #177be5 once the ring is full.
var macRingLight, macRingArc, macRingFull = ui.Hex("#0a82ff"), ui.Hex("#1769e6"), ui.Hex("#177be5")

// macProgressRing draws an NSProgressIndicator of style spinning with isIndeterminate false, as measured by reference.swift on macOS 27.2: a track 5 points wide inside the edge, at both the regular size (32 points) and the small one (16, whose ring is then 3 points from its center), black at 5.1% in Aqua (242 over white, 235 along its edges) and white at 10.2% in Dark Aqua (53 over 30); and over it the arc of the value from 12 o'clock clockwise, as wide, with round ends, in the accent: for the default blue #0a82ff in Aqua and in Dark Aqua #1769e6, #177be5 once full, measured (macRingLight, macRingArc, macRingFull); for another accent the accent itself in Aqua and 10% darker in Dark Aqua, inferred. The track's darker edges and the arc's, a quarter pixel at 2x, are not drawn.
func macProgressRing(c *ui.Context, e ui.Element, value float64) {
	t := c.Theme()
	track := tint(t.Background, t.Dark, pick(t.Dark, float32(0.051), 0.102))
	arc := t.Accent
	if !t.Dark && t.Accent == macAccent {
		arc = macRingLight
	}
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
