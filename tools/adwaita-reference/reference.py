# libadwaita's own controls, for comparing Breeze UI's GNOME look with them.
#
# It renders, with GTK 4 and libadwaita, the controls of ui/breeze/screenshot_test.go at the same places (referenceView, referenceView2, the alert and referenceView3), light and dark, at 2x, and writes them as PNGs. Breeze UI follows libadwaita 1.10.0 on GTK 4.24 (GNOME 51), so run it where those are installed, under a display; a GNOME 51 runtime such as the org.gnome.Sdk//51 Flatpak SDK is expected to carry them (inferred, not tried). Where libadwaita 1.7 is installed (Debian 13: python3-gi gir1.2-gtk-4.0 gir1.2-adw-1 xvfb xauth libgl1-mesa-dri) the controls 1.10 left alone compare too, but not the popover's ring, the dialogs' shadows, the alert's outline and body's line height, the spinner's stroke or the check mark:
#
#	xvfb-run -s "-screen 0 1600x1200x24" python3 tools/adwaita-reference/reference.py <dir>
#
# The window's frame offsets its content by 5 pixels (10 at 2x) from Breeze UI's screenshots. GNOME 51's font is Adwaita Sans 11 (font-name in gsettings-desktop-schemas 51.0), as GNOME 48's was, which Debian 13 does not package; Inter, on which it is based (fonts-inter, with gtk-font-name=Inter 11 in ~/.config/gtk-4.0/settings.ini), comes closest, so text differs a little while the controls' sizes and colors compare. A window at most 360 pixels high makes the alert .short, as it is here.
import sys

import gi

gi.require_version("Gtk", "4.0")
gi.require_version("Adw", "1")
gi.require_version("Graphene", "1.0")
from gi.repository import Adw, GLib, Graphene, Gtk  # noqa: E402

OUT = sys.argv[1] if len(sys.argv) > 1 else "."
SCALE = 2


def first_set():
    def check(label, active, sensitive=True):
        c = Gtk.CheckButton(label=label, active=active)
        c.set_sensitive(sensitive)
        return c

    def switch(active, sensitive=True):
        s = Gtk.Switch(active=active)
        s.set_sensitive(sensitive)
        s.set_valign(Gtk.Align.START)
        return s

    save = Gtk.Button(label="Save")
    save.add_css_class("suggested-action")
    disabled = Gtk.Button(label="Disabled")
    disabled.set_sensitive(False)
    return 360, 200, [
        (Gtk.Button(label="Cancel"), 16, 16), (save, 110, 16), (disabled, 200, 16),
        (check("Unchecked", False), 16, 64), (check("Checked", True), 130, 64), (check("Disabled", True, False), 240, 64),
        (switch(False), 16, 106), (switch(True), 100, 106), (switch(False, False), 184, 106),
    ]


def second_set():
    off = Gtk.CheckButton(label="Off")
    on = Gtk.CheckButton(label="On", active=True)
    off.set_group(on)
    name = Gtk.Entry(placeholder_text="Name")
    name.set_size_request(180, -1)
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
    return 520, 330, [
        (off, 16, 16), (on, 100, 16), (name, 16, 50), (search, 220, 50), (drop, 16, 90), (bold, 220, 90),
        (scale, 16, 130), (bar, 220, 142), (group, 16, 170), (spinner, 220, 170), (frame, 16, 220),
    ]


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
    return 400, 240, [(link, 16, 16), (sep, 16, 50), (linked, 16, 70), (menu, 16, 120), (info, 200, 70)], popover


def render(window, width, height, path):
    """Snapshots the window at SCALE and writes it as a PNG."""
    paintable = Gtk.WidgetPaintable.new(window)
    snapshot = Gtk.Snapshot()
    snapshot.scale(SCALE, SCALE)
    paintable.snapshot(snapshot, width, height)
    node = snapshot.to_node()
    rect = Graphene.Rect().init(0, 0, width * SCALE, height * SCALE)
    window.get_renderer().render_texture(node, rect).save_to_png(path)
    print("wrote", path)


def main():
    Adw.init()
    jobs = []
    for dark in (False, True):
        mode = "dark" if dark else "light"
        jobs.append((dark, f"adwaita-{mode}.png", first_set, False))
        jobs.append((dark, f"adwaita2-{mode}.png", second_set, False))
        jobs.append((dark, f"adwaita-alert-{mode}.png", None, True))
        jobs.append((dark, f"adwaita3-{mode}.png", third_set, False))

    loop = GLib.MainLoop()

    def next_job():
        if not jobs:
            loop.quit()
            return
        dark, file, build, alert = jobs.pop(0)
        Adw.StyleManager.get_default().set_color_scheme(Adw.ColorScheme.FORCE_DARK if dark else Adw.ColorScheme.FORCE_LIGHT)
        window = Adw.Window()
        popover = None
        if alert:
            width, height = 560, 340
            label = Gtk.Label(label="Window content", xalign=0, yalign=0)
            label.set_margin_top(16)
            label.set_margin_start(16)
            window.set_content(label)
        else:
            built = build()
            width, height, items = built[:3]
            popover = built[3] if len(built) > 3 else None
            fixed = Gtk.Fixed()
            for widget, x, y in items:
                fixed.put(widget, x, y)
            window.set_content(fixed)
        window.set_default_size(width, height)
        window.present()
        if not alert and popover is not None:
            GLib.timeout_add(300, lambda: (popover.popup(), False)[1])
        if alert:
            dialog = Adw.AlertDialog(heading="Delete “Notes”?", body="You can’t undo this action.")
            dialog.add_response("cancel", "Cancel")
            dialog.add_response("delete", "Delete")
            dialog.set_response_appearance("delete", Adw.ResponseAppearance.SUGGESTED)
            dialog.set_default_response("delete")
            dialog.present(window)

        def shoot():
            render(window, width, height, f"{OUT}/{file}")
            window.destroy()
            GLib.timeout_add(200, lambda: (next_job(), False)[1])
            return False

        # Wait past the dialog's spring animation and the spinner's first frames.
        GLib.timeout_add(1500, shoot)

    GLib.idle_add(lambda: (next_job(), False)[1])
    loop.run()


main()
