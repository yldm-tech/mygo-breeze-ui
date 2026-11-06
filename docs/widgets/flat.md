# Flat buttons

The widgets of [`ui/breeze/widgets_flat.go`](../../ui/breeze/widgets_flat.go), whose catalog entries are [`ui/breeze/catalog/flat.json`](../../ui/breeze/catalog/flat.json). What every widget shares, its platforms, themes and limits, is in [Breeze UI for MyGo](../breeze-ui.md).

## Looks

| Widget | macOS | Windows 11 | GNOME |
|---|---|---|---|
| `FlatButton` | accessoryBarAction button showing its bezel only under the pointer: no face at rest, the title in `secondaryLabelColor`; 24 points high, padded 12; under the pointer black or white at 7.8%, pressed 16.9%, the title in `labelColor` | Button of SubtleButtonStyle: transparent, `SubtleFillColorSecondary` under the pointer, `SubtleFillColorTertiary` pressed, 32 high | `.flat` text button: transparent, 7% under the pointer, 16% pressed, 34 high, bold |
| `FlatIconButton` | the same in a 24-point square, the icon in `secondaryLabelColor` at rest | Button of SubtleButtonStyle around a 16-DIP icon, 40 by 29 | `.flat` image button, 34 pixels square |
| `FlatToggle` | accessoryBar toggle in a 24-point square: off as `FlatIconButton`; on, black or white at 12.5% under a `labelColor` icon, the accent under the pointer | off as `FlatIconButton`; on, a checked ToggleButton's accent (inferred) | `.flat` image toggle: 10% while checked, 13% under the pointer, 19% pressed |

## Sources

`FlatButton`, `FlatIconButton` and `FlatToggle` (`widgets_flat.go`, with `TestScreenshotsFlat`) draw the buttons each platform has without a face at rest, for row accessories, close buttons, ⋯ menus and a player's controls. They are `ui.ButtonBase` (`ui.ToggleBase` for `FlatToggle`), so `Clicked`, `Menu`, `PopoverMenuButton`, the keys, the focus and `Disabled` are MyGo's; the icon buttons are named by their label and show it as the platform's tooltip, as `IconButton`. Their fills are translucent layers, as the platforms' are, so over a card, a sidebar or a dark part of the window (`ui.Themed` with a dark theme) they show what is behind them. Inside a `Toolbar` or a `HeaderBar` they are the bar's own items, `Button`, `IconButton` and `Toggle`, where those are the platform's flat bar items; where the bar's own would be raised they stay flat, as libadwaita's `.flat` and WinUI's `SubtleButtonStyle` keep them: a text button or toggle on GNOME, whose header bars and toolbars raise their text buttons, and a text button or any toggle in a `HeaderBar` on Windows, where those are push buttons and ToggleButtons. Their looks come from:

- **macOS:** NSButtons of bezel style accessoryBarAction (accessoryBar and pushOnPushOff for the toggle) with `showsBorderOnlyWhileMouseInside`, the image ones in 24-point square frames, rendered at 2x on macOS 27.2 by [`appkit-reference`](../../tools/appkit-reference)'s flat set (`appkit-flat-*.png`): at rest no face and the content in `secondaryLabelColor` (128 over white, 154 over the window in Dark Aqua), the bezel a push button's face (235 and 48) and pressed a pressed push button's (212 and 68) with the content in `labelColor`, a toggle on with its bezel hidden 12.5% (223 and 58), and with its bezel showing the accent (darker pressed, `#006ce1` and `#1e8bff`).
- **Windows 11:** `SubtleButtonStyle` of `Button_themeresources.xaml` (microsoft-ui-xaml main at 04a5ab5): `ButtonPadding` and `ControlCornerRadius` as a Button's, `SubtleFillColorTransparent`, `SubtleFillColorSecondary` and `SubtleFillColorTertiary` for its background and border, `TextFillColorPrimary`, `TextFillColorSecondary` pressed. WinUI has no subtle toggle; `FlatToggle` on is drawn as a checked ToggleButton.
- **GNOME:** libadwaita 1.10.0's `%button_basic_flat` in `_buttons.scss` and the `$hover_color`, `$active_color` and `$selected_*` shades of `_colors.scss`, which Breeze UI already draws for a toolbar's flat image buttons.

## Limits

Besides [those of every widget](../breeze-ui.md#limits):

- **Flat buttons.** AppKit does not put the pointer inside a button offscreen, so a flat button's look under the pointer on macOS is the bezel AppKit draws with `showsBorderOnlyWhileMouseInside` off, which is what the pointer shows (inferred); its press and its toggle on were rendered. An on toggle turns from gray to the accent under the pointer, as AppKit's accessoryBar toggle does. The Windows look follows `SubtleButtonStyle`'s resources and has not been compared with WinUI's render: [`tools/winui-reference`](../../tools/winui-reference/Program.cs) has a `winui-flat` scene for it that has not been run. WinUI has no subtle toggle, so a `FlatToggle` on is a checked ToggleButton, the accent (inferred). The GNOME look is the stylesheet's `.flat` button, compared with libadwaita's render only as a toolbar's flat image buttons. On every platform the content's colors are taken over the window's background, so over a card of another color the icon is a few levels off the platform's translucent ink.
