package breeze

import "github.com/egoist/mygo/ui"

// IconLabelButton creates a push button showing icon before label, as the platform draws a button with an image and a title: on macOS an NSButton push button with a template image leading its title, 24 points high, the image 12 points in and its title 2 points after it (measured from AppKit's rendering on macOS 27.2 by tools/appkit-reference: a push button showing a 16-point image and "Export" is 82 points wide where one showing only "Export" is 64); on Windows a Button whose content is a horizontal StackPanel of a 16-DIP FontIcon and a TextBlock 8 apart, inside its ButtonPadding and border, 32 high; on GNOME a GtkButton holding an AdwButtonContent (.image-text-button), 34 pixels high, padded 9 at either side, the 16-pixel icon 6 from the bold label, which keeps 2 more at its end (_buttons.scss of libadwaita 1.10.0). Inside a Toolbar or a HeaderBar it is the toolbar's item with an image and a title: on macOS a bordered NSToolbarItem showing both, a 36-point capsule padded 12 with the image 2 before the title, which never joins an image item beside it into one platter (measured as above: 82 points wide, as wide as the push button); on Windows an AppBarButton with its label on the right (DefaultLabelPosition Right): the icon 10 DIPs inside its 36-DIP face, the 12-DIP label 8 after it and 10 before the face's end (AppBarButtonContentViewboxMargin 12 and AppBarButtonTextLabelOnRightMargin 8,16,12,10, less AppBarButtonInnerBorderMargin's 2); on GNOME a flat image-text button, as %toolbar makes it. In a HeaderBar on Windows, whose TitleBar gives no look to the controls of its headers, it is drawn as outside a toolbar.
//
// The icon is the caller's SVG drawn 16 by 16 in the color of the button's text, as IconButton draws it, which ui.Icon takes from its ancestors whatever colors the SVG has: a template image takes the title's color in AppKit, a FontIcon the Button's Foreground in WinUI, and a symbolic icon currentColor in GTK. label is the button's name for assistive technology and the tester, as Label gives it. Clicked, the focus, the keys and Disabled are ui.ButtonBase's.
func IconLabelButton(c *ui.Context, icon *ui.SVG, label string) ui.Element {
	return iconLabelButton(c, icon, label, false).Label(label)
}

// PrimaryIconLabelButton creates IconLabelButton in the accent color for the main action, as PrimaryButton draws it: the default button of macOS (a prominent toolbar item, an accent capsule, inside a Toolbar), the accent button of Windows, the suggested action of GNOME, which stays raised in a toolbar. The icon takes the label's color over the accent, white with the platforms' default accents.
func PrimaryIconLabelButton(c *ui.Context, icon *ui.SVG, label string) ui.Element {
	return iconLabelButton(c, icon, label, true).Label(label)
}

// iconLabelLook spaces b, a button holding an icon and a label that the platform's look of button (or of a toolbar item, when bar) was just given, as IconLabelButton describes it.
func iconLabelLook(c *ui.Context, b ui.Element, bar, primary bool) {
	switch PlatformOf(c) {
	case Windows:
		// The StackPanel's Spacing of 8 is the Button's Gap already. An AppBarButton's 1-DIP border takes room inside its face, which WinUI's content, laid out beside the face rather than inside it, does not: 9 and the border put the icon 10 inside the face.
		if bar && !primary {
			b.PaddingStart(9).PaddingEnd(9)
		}
	case Linux:
		// .image-text-button's padding, buttoncontent's border-spacing and the label's padding-right (padding-left right to left) of 2.
		b.Padding(5, 9).PaddingEnd(11).Gap(6).FontWeight(700)
	default:
		b.Gap(2)
	}
}
