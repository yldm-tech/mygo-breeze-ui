# drawn_idioms: libadwaita's counterparts of the widgets of TestScreenshotsWave2Idioms (ui/breeze/drawn_idioms_test.go).
from common import Gio, Gtk, Scene, vbox, hbox, entry


def idioms_set(open_menu=False):
    """TestScreenshotsWave2Idioms's idiomView: a .toolbar of the history buttons, a text button, image buttons, toggles, an image toggle, menu buttons and a suggested button; below it entries and menu buttons outside a toolbar. GNOME has no path bar control. With open_menu, More's GtkPopoverMenu is open, as in the test's "menu" shot."""
    toolbar = Gtk.Box()
    toolbar.add_css_class("toolbar")

    def toggle(label=None, icon=None, active=False):
        t = Gtk.ToggleButton(active=active)
        if icon:
            t.set_icon_name(icon)
        else:
            t.set_label(label)
        return t

    # The menus' items act, so they are not drawn insensitive.
    actions = Gio.SimpleActionGroup()
    actions.add_action(Gio.SimpleAction(name="item"))

    def menu(label):
        model = Gio.Menu()
        for item in ["Duplicate", "Rename", "Delete"]:
            model.append(item, "menu.item")
        return Gtk.MenuButton(label=label, menu_model=model)

    save = Gtk.Button(label="Save")
    save.add_css_class("suggested-action")
    items = [
        Gtk.Button(icon_name="go-previous-symbolic"), Gtk.Button(icon_name="go-next-symbolic", sensitive=False),
        Gtk.Button(label="Bold"), Gtk.Button(icon_name="sidebar-show-symbolic"), Gtk.Button(icon_name="send-to-symbolic"),
        toggle("Wrap", active=True), toggle("Plain"), toggle(icon="sidebar-show-symbolic", active=True),
        menu("Share"), menu("More"), Gtk.Box(hexpand=True), save,
    ]
    for w in items:
        toolbar.append(w)
    city = entry(placeholder="City", width=200)
    font = entry("Serif", width=200)
    font.set_icon_from_icon_name(Gtk.EntryIconPosition.SECONDARY, "pan-down-symbolic")
    row = hbox(12, city, font, menu("Sort"), menu("Actions"), align=Gtk.Align.CENTER)
    root = vbox(16, toolbar, row, padding=12)
    root.insert_action_group("menu", actions)
    more = items[9]
    popovers = [more.get_popover()] if open_menu else []
    named = {"toolbar": toolbar, "back": items[0], "bold": items[2], "image": items[3], "wrap": items[5], "imagetoggle": items[7], "share": items[8], "save": save, "sort": row.get_last_child().get_prev_sibling()}
    if open_menu:
        named["menu"] = more.get_popover()
    # The menu button opens its popover, so it shows checked while the popover is open.
    return Scene(1000, 300, root, named=named, popovers=popovers, open_popovers=more.popup)


SETS = [
    ("idioms", "adwaita-idioms", idioms_set),
    ("menu", "adwaita-idioms-menu", lambda: idioms_set(True)),
]
