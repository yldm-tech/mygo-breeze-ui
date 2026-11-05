package breeze

import "github.com/egoist/mygo/ui"

// ---- text: the platforms' type ramps ----
//
// macOS values are AppKit's text styles, read on macOS 27.2 with NSFont.preferredFont(forTextStyle:): their sizes, their weights (the font's wght axis) and their line heights (NSLayoutManager.defaultLineHeight). Windows values are WinUI 3's TextBlock styles (TextBlock_themeresources.xaml of microsoft-ui-xaml main: CaptionTextBlockStyle, BodyTextBlockStyle, BodyStrongTextBlockStyle, SubtitleTextBlockStyle, TitleTextBlockStyle, TitleLargeTextBlockStyle, all BasedOn BaseTextBlockStyle). GNOME values are libadwaita 1.10.0's style classes (_labels.scss: .title-1, .title-2, .title-4, .heading, .caption and .dimmed; _colors.scss: --dim-opacity).

// TextStyle is a step of the platform's type ramp, for Text. TextSecondary added to a step with | dims it as the platform dims secondary text.
//
// The steps are the ones the three platforms share. A strong body, which WinUI calls BodyStrong, is TextHeading: macOS's headline is the body's 13 points in bold, WinUI's BodyStrongTextBlockStyle the body's 14 pixels in semibold and libadwaita's .heading the body's size in bold, so a separate step would draw the same text on every platform.
type TextStyle uint8

const (
	// TextBody is the platform's body text, the text ui.Text draws in the platform's theme: 13 points on macOS, 14 pixels on Windows 11 (BodyTextBlockStyle), the desktop's 11-point font on GNOME.
	TextBody TextStyle = iota
	// TextLargeTitle is the largest title: macOS's largeTitle (26 points, regular), WinUI's TitleLargeTextBlockStyle (40, semibold), libadwaita's .title-1 (181%, weight 800).
	TextLargeTitle
	// TextTitle is a page's or a dialog's title: macOS's title1 (22 points, regular), WinUI's TitleTextBlockStyle (28, semibold), libadwaita's .title-2 (136%, weight 800).
	TextTitle
	// TextSubtitle is the title of a part of a page: macOS's title2 (17 points, regular), WinUI's SubtitleTextBlockStyle (20, semibold), libadwaita's .title-4 (118%, bold).
	TextSubtitle
	// TextHeading is a heading at the body's size: macOS's headline (13 points, bold), WinUI's BodyStrongTextBlockStyle (14, semibold), libadwaita's .heading (bold).
	TextHeading
	// TextCaption is small text, as under a control: macOS's caption1 (10 points), WinUI's CaptionTextBlockStyle (12), libadwaita's .caption (82%, at a line height of 140%).
	TextCaption

	// TextSecondary dims the text it is added to: macOS's secondaryLabelColor, WinUI's TextFillColorSecondary, libadwaita's .dimmed (the text at 55% opacity, 90% with higher contrast). TextCaption|TextSecondary is a dimmed caption; TextSecondary alone, dimmed body text.
	TextSecondary TextStyle = 1 << 7
)

// textRamp is a step of a type ramp: its size as a multiple of the body's, its weight, and its line height as a multiple of its size, 0 for the font's own line spacing.
type textRamp struct {
	scale      float32
	weight     int
	lineHeight float32
}

// textRampOf returns style's step of p's type ramp, and false for TextBody or a style that is no step, which keep ui.Text's own. Every platform's body is the theme's FontSize, so the steps are kept as multiples of it, which also follow the desktop's text scale as the theme does.
//   - macOS: NSFont.preferredFont(forTextStyle:) on macOS 27.2 gives largeTitle 26 points, title1 22, title2 17, headline 13 at weight 700 (.SFNS-Bold), body 13 and caption1 10, all regular but headline. Their line heights (NSLayoutManager.defaultLineHeight: 32, 26, 22, 16, 16 and 13) are SF's own line spacing, its ascender, descender and leading, and none is set: MyGo lays out a line of every text, ui.Text's too, as the font's ascender and descender without the leading (15.3 for the body's 16, 30.6 for the large title's 32, measured on macOS 27.2), so a line height of AppKit's on the steps would leave the body, which must stay ui.Text's, the only step shorter than AppKit's by its leading. SF's tracking, which narrows large sizes and widens small ones, comes with the system font in Core Text, so no letter spacing is set either.
//   - Windows: TextBlock_themeresources.xaml sets FontSize (TitleLargeTextBlockFontSize 40, TitleTextBlockFontSize 28, SubtitleTextBlockFontSize 20, BodyTextBlockFontSize 14, CaptionTextBlockFontSize 12) and FontWeight (SemiBold from BaseTextBlockStyle, Normal on Body and Caption), and no LineHeight: LineStackingStrategy MaxHeight takes the font's own line spacing. It sets no CharacterSpacing.
//   - GNOME: _labels.scss sets .title-1 { font-weight: 800; font-size: 181% }, .title-2 { 800; 136% }, .title-4 { 700; 118% }, .heading { 700 }, .caption { 400; 82%; line-height: 140% }, and no letter-spacing. Subtitle is .title-4 rather than .title-3, which is .title-2's 136% at weight 700 and so no smaller than the title; .title-4 is the next smaller size, as Subtitle is on the other platforms. TextBody is a plain label, not .body: .body adds a line height of 140% for paragraphs, which ui.Text does not have.
func textRampOf(p Platform, style TextStyle) (textRamp, bool) {
	switch p {
	case Windows:
		switch style {
		case TextLargeTitle:
			return textRamp{40.0 / 14, 600, 0}, true
		case TextTitle:
			return textRamp{28.0 / 14, 600, 0}, true
		case TextSubtitle:
			return textRamp{20.0 / 14, 600, 0}, true
		case TextHeading:
			return textRamp{1, 600, 0}, true
		case TextCaption:
			return textRamp{12.0 / 14, 400, 0}, true
		}
	case Linux:
		switch style {
		case TextLargeTitle:
			return textRamp{1.81, 800, 0}, true
		case TextTitle:
			return textRamp{1.36, 800, 0}, true
		case TextSubtitle:
			return textRamp{1.18, 700, 0}, true
		case TextHeading:
			return textRamp{1, 700, 0}, true
		case TextCaption:
			return textRamp{0.82, 400, 1.4}, true
		}
	default:
		switch style {
		case TextLargeTitle:
			return textRamp{26.0 / 13, 400, 0}, true
		case TextTitle:
			return textRamp{22.0 / 13, 400, 0}, true
		case TextSubtitle:
			return textRamp{17.0 / 13, 400, 0}, true
		case TextHeading:
			return textRamp{1, 700, 0}, true
		case TextCaption:
			return textRamp{10.0 / 13, 400, 0}, true
		}
	}
	return textRamp{}, false
}

// Text creates a text in a step of the platform's type ramp, style, which TextSecondary added with | dims: breeze.Text(c, breeze.TextTitle, "Settings") draws macOS's title1, WinUI's TitleTextBlockStyle or libadwaita's .title-2. It is ui.Text with the step's size, weight and line height, in the theme's font, so it wraps, selects (Element.Selectable), reads right to left and is read by assistive technology as ui.Text is, and TextBody draws what ui.Text draws. A title stays plain text to assistive technology, as the platforms' title styles do: add Role(ui.RoleHeading) and Level to make it a heading.
//
// Unlike Breeze UI's other widgets it is not ui.Text's counterpart under the same signature, as package ui has no type ramp: it takes the step before the string.
func Text(c *ui.Context, style TextStyle, s string) ui.Element {
	e := ui.Text(c, s)
	p := PlatformOf(c)
	t := c.Theme()
	if r, ok := textRampOf(p, style&^TextSecondary); ok {
		base := t.FontSize
		if base <= 0 {
			base = 14 // MyGo's size for a theme that sets none
		}
		e.FontSize(base * r.scale).FontWeight(r.weight)
		if r.lineHeight > 0 {
			e.LineHeight(r.lineHeight)
		}
	}
	if style&TextSecondary != 0 {
		if p == Linux {
			// .dimmed lowers the label's opacity rather than its color (_labels.scss), to --dim-opacity: 55%, and 90% under prefers-contrast: more (_colors.scss).
			e.Opacity(pick(c.Preferences().HighContrast, float32(0.55), 0.9))
		} else {
			// secondaryLabelColor and TextFillColorSecondary are the theme's TextMuted (mac.go, fluent.go).
			e.TextColor(t.TextMuted)
		}
	}
	return e
}
