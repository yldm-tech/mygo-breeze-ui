# settings: libadwaita's counterparts of the rows of TestScreenshotsSettings (ui/breeze/widgets_settings_test.go), and the spin row of TestAdwPlusEvenWhenDisabled.
from common import Adw, Gtk, Scene, vbox, hbox, start, strings


def settings_set():
    """TestScreenshotsSettings's settingsView: in a 300-pixel column, a boxed list of an AdwComboRow with a subtitle, an AdwEntryRow holding a name, an empty one and an AdwSpinRow with a subtitle, then an AdwComboRow and a disabled AdwSpinRow each alone in a boxed list; beside them an open AdwExpanderRow with a subtitle holding an action row and a switch row, and a closed one, each in a boxed list."""
    def boxed(*rows):
        box = Gtk.ListBox(selection_mode=Gtk.SelectionMode.NONE)
        box.add_css_class("boxed-list")
        for row in rows:
            box.append(row)
        return box

    def combo(title, items, selected, subtitle=""):
        return Adw.ComboRow(title=title, subtitle=subtitle, model=strings(items), selected=selected)

    def entry(title, text):
        row = Adw.EntryRow(title=title)
        row.set_text(text)
        return row

    def spin(title, value, lower, upper, subtitle="", sensitive=True):
        row = Adw.SpinRow.new_with_range(lower, upper, 1)
        row.set_title(title)
        row.set_subtitle(subtitle)
        row.set_value(value)
        row.set_sensitive(sensitive)
        return row

    settings = boxed(combo("Appearance", ["Light", "Dark", "Auto"], 2, "Light or dark"), entry("Name", "Evan"), entry("Email", ""), spin("Font size", 12, 6, 72, "Points"))
    language = combo("Language", ["English", "Français"], 0)
    copies = spin("Copies", 3, 1, 9, sensitive=False)
    left = vbox(16, settings, boxed(language), boxed(copies), width=300)

    network = Adw.ExpanderRow(title="Network", subtitle="Wi-Fi and proxies", expanded=True)
    network.add_row(Adw.ActionRow(title="Proxy", subtitle="Automatic"))
    network.add_row(Adw.SwitchRow(title="Metered"))
    advanced = Adw.ExpanderRow(title="Advanced")
    advanced.add_row(Adw.ActionRow(title="Logs"))
    right = vbox(16, boxed(network), boxed(advanced), width=300)
    # The rows ask to expand, which would widen the columns past the test's 300 to the window's end.
    left.set_hexpand(False)
    right.set_hexpand(False)
    root = hbox(24, left, right, padding=16)
    return Scene(672, 480, root, named={"settings": settings, "language": language, "copies": copies, "network": network, "advanced": advanced})


def spin_limit_set():
    """The spin row of TestAdwPlusEvenWhenDisabled: in a boxed list 300 wide, 16 from the window's edges, an AdwSpinRow Copies at the upper limit of 1 to 9, its + button insensitive, to read the icon of a disabled button. Nothing keeps the focus, which the row would take as the window's only focusable widget."""
    box = Gtk.ListBox(selection_mode=Gtk.SelectionMode.NONE)
    box.add_css_class("boxed-list")
    row = Adw.SpinRow.new_with_range(1, 9, 1)
    row.set_title("Copies")
    row.set_value(9)
    box.append(row)
    root = vbox(0, box, padding=16, width=300)
    root.set_hexpand(False)
    return Scene(332, 120, start(root), named={"row": row}, shown=lambda window: window.set_focus(None))


SETS = [
    ("spinlimit", "adwaita-spinlimit", spin_limit_set),
    ("settings", "adwaita-settings", settings_set),
]
