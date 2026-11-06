# rows: libadwaita's counterparts of the cards and rows of TestScreenshotsRows (ui/breeze/widgets_rows_test.go).
from common import Adw, Gtk, Scene, vbox, hbox, pad


def rows_set():
    """TestScreenshotsRows's rowsView: a .card of text padded 12, a boxed list of an action row with a button, a switch row with a subtitle and one without, an action row with a dimmed value and a disabled switch row each alone in a boxed list, a preferences group of rows for the fieldset, and an expander row with its nested rows for the accordion."""
    def boxed(*rows):
        box = Gtk.ListBox(selection_mode=Gtk.SelectionMode.NONE)
        box.add_css_class("boxed-list")
        for row in rows:
            box.append(row)
        return box

    def switch_row(title, active, subtitle="", sensitive=True):
        row = Adw.SwitchRow(title=title, subtitle=subtitle, active=active)
        row.set_sensitive(sensitive)
        return row

    text = Gtk.Label(label="A card holds what belongs together.", xalign=0, wrap=True)
    card = Gtk.Box()
    card.add_css_class("card")
    card.append(pad(text, 12))
    wallpaper = Adw.ActionRow(title="Wallpaper", subtitle="Choose a picture")
    open_button = Gtk.Button(label="Open", valign=Gtk.Align.CENTER)
    wallpaper.add_suffix(open_button)
    settings = boxed(wallpaper, switch_row("Dark mode", True, "Use a dark appearance"), switch_row("Notifications", False))
    storage = Adw.ActionRow(title="Storage")
    value = Gtk.Label(label="120 GB")
    value.add_css_class("dimmed")
    storage.add_suffix(value)
    left = vbox(16, card, settings, boxed(storage), boxed(switch_row("Bluetooth", True, sensitive=False)), width=300)

    privacy = Adw.PreferencesGroup(title="Privacy")
    privacy.add(switch_row("Location", True, "Apps may ask for it"))
    privacy.add(Adw.ActionRow(title="Camera", subtitle="2 apps"))
    network = Adw.ExpanderRow(title="Network", expanded=True)
    network.add_row(Adw.ActionRow(title="Proxy", subtitle="Automatic"))
    network.add_row(switch_row("Metered", False))
    accordion = boxed(network)
    right = vbox(16, privacy, accordion, width=300)
    # The rows ask to expand, which would widen the column past the test's 300 to the window's end.
    right.set_hexpand(False)
    root = hbox(24, left, right, padding=16)
    return Scene(672, 520, root, named={"card": card, "rows": settings, "storage": storage, "privacy": privacy, "network": network})


SETS = [
    ("rows", "adwaita-rows", rows_set),
]
