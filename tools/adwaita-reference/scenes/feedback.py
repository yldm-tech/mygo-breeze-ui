# feedback: libadwaita's counterparts of the widgets of TestScreenshotsFeedback (ui/breeze/widgets_feedback_test.go).
from common import Adw, Gtk, Scene, vbox


def feedback_set():
    """TestScreenshotsFeedback's feedbackView: a search bar for the find bar, a button, level bars at their three levels, avatars and a toast with a button. GTK's tooltip is a surface that shows only under a pointer and is left out; GNOME has no rating control."""
    search = Gtk.SearchEntry(text="cat")
    search.set_hexpand(True)
    bar = Gtk.SearchBar(child=search, search_mode_enabled=True, show_close_button=True)

    def level(value):
        lb = Gtk.LevelBar(min_value=0, max_value=100, value=value)
        # Breeze UI's Meter: the warning level from 70 and the critical one from 90, as GtkLevelBar's offsets low and high.
        lb.add_offset_value(Gtk.LEVEL_BAR_OFFSET_LOW, 70)
        lb.add_offset_value(Gtk.LEVEL_BAR_OFFSET_HIGH, 90)
        lb.add_offset_value(Gtk.LEVEL_BAR_OFFSET_FULL, 100)
        lb.set_size_request(220, -1)
        return lb

    share = Gtk.Button(label="Share")
    levels = [level(40), level(75), level(95)]
    avatars = [Adw.Avatar(size=32, text=name, show_initials=True) for name in ["Ada Lovelace", "Grace Hopper", "Linus", "山田太郎", ""]]
    content = Gtk.Fixed()
    content.set_vexpand(True)
    content.put(share, 16, 16)
    for i, lb in enumerate(levels):
        # A meter's bar is centered in its 20-pixel row (the tests place the rows 28 apart).
        content.put(lb, 300, 16 + 28 * i + 6)
    for i, a in enumerate(avatars):
        content.put(a, 16 + 48 * i, 110)
    overlay = Adw.ToastOverlay(child=vbox(0, bar, content))
    toast = Adw.Toast(title="Note deleted", button_label="Undo", timeout=0)
    named = {"searchbar": bar, "share": share, "level40": levels[0], "level75": levels[1], "level95": levels[2], "avatar": avatars[0]}
    scene = Scene(560, 400, overlay, named=named)
    scene.shown = lambda window: overlay.add_toast(toast)
    return scene


SETS = [
    ("feedback", "adwaita-feedback", feedback_set),
]
