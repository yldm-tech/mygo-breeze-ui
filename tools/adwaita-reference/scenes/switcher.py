# switcher: libadwaita's counterpart of TestScreenshotsViewSwitcher (ui/breeze/widgets_switcher_test.go).
from common import Adw, Gtk, Scene, vbox, start


def switcher_set():
    """TestScreenshotsViewSwitcher's switcherView: an AdwHeaderBar at the top of the window's content titled Music with an AdwViewSwitcher of the wide policy packed at its end, where Breeze UI's HeaderBar holds a ViewSwitcher among its end controls (libadwaita apps put it in place of the title, as the header bar's title-widget, which Breeze UI's HeaderBar does not take), GTK's close button after it (decoration-layout :close), and below the bar a wide AdwViewSwitcher alone, 16 pixels in from the window's side and 16 below the bar; each switches an AdwViewStack of Library, Browse and Radio with icons, Browse visible. The stacks are hidden and take no room: only the switchers show."""
    def switcher():
        stack = Adw.ViewStack()
        for title, icon in (("Library", "folder-music-symbolic"), ("Browse", "view-grid-symbolic"), ("Radio", "network-wireless-symbolic")):
            stack.add_titled_with_icon(Gtk.Label(label=title), title.lower(), title, icon)
        stack.set_visible_child_name("browse")
        stack.set_visible(False)
        return Adw.ViewSwitcher(stack=stack, policy=Adw.ViewSwitcherPolicy.WIDE), stack

    bar_switcher, bar_stack = switcher()
    bar = Adw.HeaderBar(title_widget=Adw.WindowTitle(title="Music", subtitle=""), decoration_layout=":close")
    bar.pack_end(bar_switcher)
    alone, alone_stack = switcher()
    root = vbox(0, bar, vbox(0, start(alone), alone_stack, bar_stack, padding=16))
    return Scene(560, 140, root, named={"headerbar": bar, "bar-switcher": bar_switcher, "switcher": alone})


SETS = [
    ("switcher", "adwaita-switcher", switcher_set),
]
