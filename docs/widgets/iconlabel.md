# Icon-label buttons

The widgets of [`ui/breeze/widgets_iconlabel.go`](../../ui/breeze/widgets_iconlabel.go), whose catalog entries are [`ui/breeze/catalog/iconlabel.json`](../../ui/breeze/catalog/iconlabel.json). What every widget shares, its platforms, themes and limits, is in [Breeze UI for MyGo](../breeze-ui.md).

## Looks

| Widget | macOS | Windows 11 | GNOME |
|---|---|---|---|
| `IconLabelButton`, `PrimaryIconLabelButton` | push button with a template image leading its title: the image 12 points in, the title 2 after it, 24 high; in a toolbar or header bar, a bordered item showing both, a 36-point capsule padded 12 that never joins an image item's platter, the primary one prominent | Button holding a 16-DIP icon and its label 8 apart, 32 high; in a toolbar, an AppBarButton with its label on the right, the icon 10 inside its face and the 12-DIP label 8 after it, the primary one an accent Button | a button holding an AdwButtonContent (`.image-text-button`): padded 9, the 16-pixel icon 6 from the bold label, 34 high; flat in a toolbar or header bar, where the suggested action stays raised |

## Sources

This group has no sources of its own yet: its row in the table says what each platform's look follows, and [`tools/appkit-reference`](../../tools/appkit-reference) renders AppKit's counterparts (`Scenes/iconlabel.swift`).

## Limits

Besides [those of every widget](../breeze-ui.md#limits):

- **Icon-label buttons.** `IconLabelButton` and `PrimaryIconLabelButton` are `Button` and `PrimaryButton` with the caller's SVG before the label, drawn 16 by 16 in the label's color, and the label as the button's name. On macOS they are compared at 2x with AppKit's push buttons and bordered NSToolbarItems showing a 16-point template image and a title, rendered by the icon-label scenes of [`appkit-reference`](../../tools/appkit-reference) where `TestScreenshotsIconLabel` places Breeze UI's: their widths, the image's place and the title 2 points after it match, and an item with an image and a title stays on a platter of its own beside an image item, 8 points away. The prominent item's accent does not render offscreen and is `PrimaryButton`'s in a toolbar; a disabled button fades at half opacity where AppKit grays the image and title over a lighter face (245, the ink 182). The Windows and GNOME looks are inferred from the sources rather than compared with renders: `AppBarButton_themeresources.xaml` of microsoft-ui-xaml's main branch (the LabelOnRight state's `AppBarButtonContentViewboxMargin` 12,16,0,10 and `AppBarButtonTextLabelOnRightMargin` 8,16,12,10, laid out beside the face, which is `AppBarButtonInnerBorderMargin` 2 inside) and the Button's content as a StackPanel of a FontIcon and a TextBlock with a Spacing of 8, the icon centered on the label's line (inferred); libadwaita 1.10.0's `_buttons.scss` (`.image-text-button`'s padding of 9 and `buttoncontent`'s spacing of 6, its bold label padded 2 at its end), `_toolbars.scss` (`%toolbar` makes `.image-text-button` flat) and `adw-button-content.c` (which gives the button that class).
