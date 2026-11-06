# headerbar: libadwaita's counterpart of TestScreenshotsHeaderBar (ui/breeze/widgets_headerbar_test.go).
from common import Adw, Gtk, Scene, vbox, pad


def headerbar_set(subtitle=""):
    """TestScreenshotsHeaderBar's headerView, the sets headerbar and headerbarsub: an AdwHeaderBar at the top of the window's content, outside a window's titlebar slot and an AdwToolbarView as Breeze UI's HeaderBar is drawn, its title an AdwWindowTitle reading Documents (with the subtitle 12 items in headerbarsub), go-previous-symbolic packed at its start and Edit and send-to-symbolic at its end, GTK's close button after them (decoration-layout :close, the window controls headerTestRoom gives room for on GNOME), then the window's content padded 16."""
    bar = Adw.HeaderBar(title_widget=Adw.WindowTitle(title="Documents", subtitle=subtitle), decoration_layout=":close")
    back = Gtk.Button(icon_name="go-previous-symbolic")
    edit = Gtk.Button(label="Edit")
    share = Gtk.Button(icon_name="send-to-symbolic")
    bar.pack_start(back)
    bar.pack_end(share)
    bar.pack_end(edit)
    content = Gtk.Label(label="Content", xalign=0)
    root = vbox(0, bar, pad(content, 16))
    return Scene(560, 120, root, named={"headerbar": bar, "back": back, "edit": edit, "share": share})


SETS = [
    ("headerbar", "adwaita-headerbar", headerbar_set),
    ("headerbarsub", "adwaita-headerbar-subtitle", lambda: headerbar_set("12 items")),
]
