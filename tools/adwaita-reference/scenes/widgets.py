# widgets: the sets of TestScreenshots (ui/breeze/screenshot_test.go): referenceView, referenceView2, the alert and referenceView3.
from common import Adw, Gtk, Scene, fixed, pad, check, entry


def first_set():
    def switch(active, sensitive=True):
        s = Gtk.Switch(active=active)
        s.set_sensitive(sensitive)
        s.set_valign(Gtk.Align.START)
        return s

    save = Gtk.Button(label="Save")
    save.add_css_class("suggested-action")
    disabled = Gtk.Button(label="Disabled")
    disabled.set_sensitive(False)
    items = [
        (Gtk.Button(label="Cancel"), 16, 16), (save, 110, 16), (disabled, 200, 16),
        (check("Unchecked", False), 16, 64), (check("Checked", True), 130, 64), (check("Disabled", True, False), 240, 64),
        (switch(False), 16, 106), (switch(True), 100, 106), (switch(False, False), 184, 106),
    ]
    return fixed(360, 200, items, named={"cancel": items[0][0], "check": items[4][0], "switch": items[7][0]})


def second_set():
    off = Gtk.CheckButton(label="Off")
    on = Gtk.CheckButton(label="On", active=True)
    off.set_group(on)
    name = entry(placeholder="Name", width=180)
    search = Gtk.SearchEntry(placeholder_text="Search")
    search.set_size_request(180, -1)
    drop = Gtk.DropDown.new_from_strings(["Medium", "Small", "Large"])
    bold = Gtk.ToggleButton(label="Bold", active=True)
    scale = Gtk.Scale.new_with_range(Gtk.Orientation.HORIZONTAL, 0, 100, 1)
    scale.set_value(40)
    scale.set_draw_value(False)
    scale.set_size_request(180, -1)
    bar = Gtk.ProgressBar(fraction=0.4)
    bar.set_size_request(180, -1)
    group = Adw.ToggleGroup()
    for label in ["Day", "Week", "Month"]:
        group.add(Adw.Toggle(label=label, name=label))
    group.set_active_name("Day")
    # In a Gtk.Fixed the group gets its minimum width, which ellipsizes the labels; give it its natural one.
    group.set_size_request(group.measure(Gtk.Orientation.HORIZONTAL, -1)[1], -1)
    spinner = Adw.Spinner()
    text = Gtk.TextView()
    text.get_buffer().set_text("Notes")
    frame = Gtk.Frame(child=text)
    frame.set_size_request(180, 70)
    named = {"radio": on, "entry": name, "search": search, "dropdown": drop, "toggle": bold, "scale": scale, "progress": bar, "togglegroup": group, "spinner": spinner, "frame": frame}
    return fixed(520, 330, [
        (off, 16, 16), (on, 100, 16), (name, 16, 50), (search, 220, 50), (drop, 16, 90), (bold, 220, 90),
        (scale, 16, 130), (bar, 220, 142), (group, 16, 170), (spinner, 220, 170), (frame, 16, 220),
    ], named=named)


def third_set():
    link = Gtk.Label()
    link.set_markup('<a href="https://example.com">Docs</a>')
    sep = Gtk.Separator(orientation=Gtk.Orientation.HORIZONTAL)
    sep.set_size_request(200, -1)
    linked = Gtk.Box()
    linked.add_css_class("linked")
    linked.append(Gtk.ToggleButton(label="Bold", active=True))
    linked.append(Gtk.ToggleButton(label="Italic"))
    menu = Gtk.MenuButton(label="Sort")
    info = Gtk.Button(label="Info")
    popover = Gtk.Popover(child=Gtk.Label(label="Details"))
    popover.set_parent(info)
    popover.set_position(Gtk.PositionType.BOTTOM)
    popover.set_has_arrow(False)
    return fixed(400, 240, [(link, 16, 16), (sep, 16, 50), (linked, 16, 70), (menu, 16, 120), (info, 200, 70)], popovers=[popover], named={"link": link, "linked": linked, "menubutton": menu, "info": info, "popover": popover})


def alert_set():
    label = Gtk.Label(label="Window content", xalign=0, yalign=0)
    pad(label, (16, 0, 0, 16))
    scene = Scene(560, 340, label)

    def shown(window):
        dialog = Adw.AlertDialog(heading="Delete “Notes”?", body="You can’t undo this action.")
        dialog.add_response("cancel", "Cancel")
        dialog.add_response("delete", "Delete")
        dialog.set_response_appearance("delete", Adw.ResponseAppearance.SUGGESTED)
        dialog.set_default_response("delete")
        dialog.present(window)

    scene.shown = shown
    return scene


SETS = [
    ("first", "adwaita", first_set),
    ("second", "adwaita2", second_set),
    ("alert", "adwaita-alert", alert_set),
    ("third", "adwaita3", third_set),
]
