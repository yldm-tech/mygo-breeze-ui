# values: libadwaita's counterparts of the widgets of TestScreenshotsValues (ui/breeze/widgets_values_test.go).
from common import Gdk, Gtk, Scene, vbox, hbox, start, entry


def values_set():
    """TestScreenshotsValues's valuesView: a button and an entry, a scale, scales with marks after them (one insensitive), entries for the date and the time, a spin button and color buttons beside a calendar. GTK has no range slider; GTK's date and time are entries."""
    def scale(value, top, step, marks, sensitive=True):
        s = Gtk.Scale.new_with_range(Gtk.Orientation.HORIZONTAL, 0, top, step)
        s.set_value(value)
        s.set_draw_value(False)
        s.set_sensitive(sensitive)
        if marks:
            v = 0
            while v <= top:
                s.add_mark(v, Gtk.PositionType.BOTTOM, None)
                v += step
        return s

    volume, quality, off = scale(40, 100, 1, False), scale(50, 100, 25, True), scale(3, 10, 1, True, False)
    spin = Gtk.SpinButton.new_with_range(0, 1000, 0.01)
    spin.set_digits(2)
    spin.set_value(12.5)

    def color(rgba, sensitive=True):
        b = Gtk.ColorDialogButton(dialog=Gtk.ColorDialog(with_alpha=True))
        c = Gdk.RGBA()
        c.parse(rgba)
        b.set_rgba(c)
        b.set_sensitive(sensitive)
        return b

    colors = [color("#3584e4"), color("#ffffff"), color("rgba(255,59,48,0.4)"), color("#3584e4", False)]
    left = vbox(14,
                hbox(8, Gtk.Button(label="Button"), entry(placeholder="Text", width=120), align=Gtk.Align.CENTER),
                volume, quality, off,
                start(entry("1/5/2026")), start(entry("5:45 PM")), start(spin),
                hbox(8, *colors, align=Gtk.Align.CENTER), width=260)
    calendar = Gtk.Calendar(year=2026, month=0, day=5)
    root = hbox(24, left, calendar, padding=16)
    return Scene(640, 460, root, named={"volume": volume, "quality": quality, "off": off, "spin": spin, "color": colors[0], "calendar": calendar})


SETS = [
    ("values", "adwaita-values", values_set),
]
