# layout: libadwaita's counterparts of the widgets of TestScreenshotsLayout (ui/breeze/widgets_layout_test.go).
from common import Adw, Gtk, Scene, vbox, hbox, pad, start, check, entry


def layout_set():
    """TestScreenshotsLayout's layoutView: expanders, a boxed list of expander rows for the accordion, an entry beside its label, a preferences group for the fieldset, a check button group and a paned in a frame."""
    closed = Gtk.Expander(label="Closed", child=Gtk.Label(label="Hidden"))
    advanced = Gtk.Expander(label="Advanced", expanded=True, child=start(check("Verbose logging", False)))
    accordion = Gtk.ListBox(selection_mode=Gtk.SelectionMode.NONE)
    accordion.add_css_class("boxed-list")
    general = Adw.ExpanderRow(title="General")
    general.add_row(Adw.ActionRow(title="General settings"))
    privacy = Adw.ExpanderRow(title="Privacy", expanded=True)
    share = Adw.ActionRow(title="Share usage data")
    share.add_prefix(check("", True))
    privacy.add_row(share)
    accordion.append(general)
    accordion.append(privacy)
    left = vbox(12, start(closed), start(advanced), accordion, width=250)

    name_label = Gtk.Label(label="Name")
    name = entry("Ada")
    name.set_hexpand(True)
    name_row = hbox(12, name_label, name, align=Gtk.Align.CENTER)
    shipping = Adw.PreferencesGroup(title="Shipping")
    address = Adw.EntryRow(title="Address")
    shipping.add(address)
    group = vbox(0, start(check("Notifications", False, inconsistent=True)), pad(start(check("Mail", True)), (0, 0, 0, 24)), pad(start(check("Calendar", False)), (0, 0, 0, 24)))
    paned = Gtk.Paned(start_child=pad(Gtk.Label(label="Sidebar", xalign=0, yalign=0), 8), end_child=pad(Gtk.Label(label="Content", xalign=0, yalign=0), 8), position=120)
    paned.add_css_class("frame")
    paned.set_size_request(-1, 60)
    # The test's Form puts its Fieldset 24 below the field above it (ui.Form's gap of 12 and the 12 a Fieldset keeps from what is above it), as an AdwPreferencesPage spaces its groups 24 apart.
    form = vbox(24, name_row, shipping)
    right = vbox(16, form, group, paned, width=270)
    root = hbox(24, left, right, padding=16)
    return Scene(600, 470, root, named={"expander": closed, "accordion": accordion, "expanderrow": privacy, "name": name, "shipping": shipping, "group": group, "paned": paned})


SETS = [
    ("layout", "adwaita-layout", layout_set),
]
