# Inputs

The widgets of [`ui/breeze/widgets_inputs.go`](../../ui/breeze/widgets_inputs.go), whose catalog entries are [`ui/breeze/catalog/inputs.json`](../../ui/breeze/catalog/inputs.json). What every widget shares, its platforms, themes and limits, is in [Breeze UI for MyGo](../breeze-ui.md).

## Looks

| Widget | macOS | Windows 11 | GNOME |
|---|---|---|---|
| `PasswordField` | NSSecureTextField: the text field's bezel, nothing to reveal the text (its caps-lock indicator is not drawn) | PasswordBox: a TextBox with the 28 by 24 RevealButton at its end while it holds text typed since the focus came, revealing while pressed or with Alt+F8 held | GtkPasswordEntry: an entry with Show Text in its context menu, no peek icon (GTK's default) |
| `SplitButton` | NSComboButton of style split: a push button's face, a divider 1 by 14 points, a 16-point part with a 9 by 5 chevron; in a toolbar, the title in a 25-point pill on the item's platter and a 40-point arrow part | SplitButton: two button faces with a 1-DIP divider, their border over the faces, a 35-DIP part with ChevronDownSmall; in a toolbar, SplitButtonCommandBarStyle, the label 13 in, a 23-DIP part with a ChevronDown glyph | AdwSplitButton: raised buttons meeting at a seam, a 24-pixel arrow part; in a toolbar, flat with a separator 6 in from the edges |
| `IconButton` | bordered push button, 12 points either side of the image; in a toolbar, an image item (joined with one beside it) | Button with a 16-DIP icon, 40 by 29; in a toolbar, an icon-only AppBarButton | `.image-button`, 34 pixels square, raised; flat in a toolbar |

## Sources

They were measured with the other widgets: see [Sources](../breeze-ui.md#sources).

## Limits

Besides [those of every widget](../breeze-ui.md#limits):

- **Password fields.** The bullets are MyGo's (U+2022 in the field's font), smaller than AppKit's 5.5-point dots 8 points apart (measured from NSSecureTextField's render) and WinUI's PasswordChar (U+25CF); GTK's are its invisible character, U+25CF where the font has it (`find_invisible_char` in `gtktext.c`). NSSecureTextField's and GtkPasswordEntry's caps-lock indicator is not drawn, as MyGo does not tell a view whether Caps Lock is on. WinUI's RevealButton also reveals the text while assistive technology toggles it (the Toggle pattern); Breeze UI's reveals it while the pointer presses it or Alt+F8 is held. On GNOME the peek icon stays off, as GtkPasswordEntry's default; an app wanting AdwPasswordEntryRow's has to draw it.
- **Split buttons.** Both parts are Tab stops, as in GTK, where WinUI's SplitButton takes the focus once and opens its menu with Alt+Down or F4 (`SplitButton.cpp`), and the focus ring goes around the part that has it. Styles given to what `SplitButton` returns apply to the main part only, but for `Disabled`, which the arrow part takes from it in the next frame; give margins and sizes to an element around it. WinUI's FlyoutOpen state (both parts pressed while the menu shows) is not drawn: the system's menu does not tell the view it shows.

## What MyGo would need

Besides [what every widget would need](../breeze-ui.md#what-mygo-would-need):

- **Inputs.** The state of Caps Lock (a modifier bit or `Context.CapsLock`), for the caps-lock indicators of NSSecureTextField and GtkPasswordEntry; the character a password field hides its text with (`Element.PasswordChar`), for AppKit's and WinUI's larger dots; a way to take a button out of the Tab order (`Element.Focusable(false)`, as MyGo's own search field does for its clear button), so that a split button is one stop; and `Element.OpenMenu` on a menu button, so that the main part of a split button opens the menu with Alt+Down, as WinUI's does.
