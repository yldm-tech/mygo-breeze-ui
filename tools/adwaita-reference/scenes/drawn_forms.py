# drawn_forms: libadwaita's counterparts of the widgets of TestScreenshotsWave2Forms (ui/breeze/drawn_forms_test.go).
from common import Adw, Gtk, Scene, vbox, hbox, pad, start, check, entry


def forms_set():
    """TestScreenshotsWave2Forms's formsShotView: check button groups, spin buttons (one at its maximum), a label above an entry, a preferences group, entries for the time and the date, and a paned in a frame."""
    def group(head, children, inconsistent):
        heads = check(head, not inconsistent, inconsistent=inconsistent)
        rows = [start(heads)] + [pad(start(check(label, on)), (0, 0, 0, 24)) for label, on in children]
        return vbox(0, *rows)

    def spin(value):
        s = Gtk.SpinButton.new_with_range(0, 10, 1)
        s.set_value(value)
        return start(s)

    left = vbox(16, group("Notifications", [("Mail", True), ("Calendar", False), ("Photos", False)], True), group("Sync", [("Contacts", True), ("Notes", True)], False), spin(5), spin(10), width=260)
    name = vbox(6, start(Gtk.Label(label="Name")), entry("Ada"))
    shipping = Adw.PreferencesGroup(title="Shipping")
    shipping.add(Adw.EntryRow(title="Address"))
    paned = Gtk.Paned(start_child=pad(Gtk.Label(label="Sidebar", xalign=0, yalign=0), 8), end_child=pad(Gtk.Label(label="Content", xalign=0, yalign=0), 8), position=120)
    paned.add_css_class("frame")
    paned.set_size_request(-1, 60)
    right = vbox(16, name, shipping, hbox(12, entry("15:04"), entry("1/5/2026"), align=Gtk.Align.CENTER), paned, width=300)
    root = hbox(24, left, right, padding=16)
    return Scene(640, 440, root, named={"spin": left.get_last_child(), "shipping": shipping, "paned": paned})


SETS = [
    ("forms", "adwaita-forms", forms_set),
]
