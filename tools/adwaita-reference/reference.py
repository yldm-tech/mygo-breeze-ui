# libadwaita's own controls, for comparing Breeze UI's GNOME look with them.
#
# It renders, with GTK 4 and libadwaita, the controls of Breeze UI's screenshot tests at the same places, light and dark, at 2x, and writes them as PNGs: libadwaita's counterparts of the widgets of each group's screenshot test (scenes/<group>.py), laid out in boxes with the paddings and gaps of those tests' rows and columns, so a control sits where Breeze UI's does when both are the same size. It prints every placed control's bounds in window DIPs, and writes popovers, which GTK draws in surfaces of their own, over the window at their places. Breeze UI follows libadwaita 1.10.0 on GTK 4.24 (GNOME 51); it was run with GTK 4.24.1, libadwaita 1.10.0, PyGObject 3.58.1 and the Adwaita Sans font on Arch Linux, under Xvfb:
#
#	xvfb-run -a -s "-screen 0 1600x1200x24" python3 tools/adwaita-reference/reference.py <dir> [set...]
#
# The sets are those of the modules of scenes/, each module's SETS a list of (key, file prefix, build); with keys given, it renders only those sets, and without any, it renders them all. The window's frame offsets its content by 5 pixels (10 at 2x) from Breeze UI's screenshots, and the window is 10 pixels wider and higher than the test's so that its content is the test's size. GNOME 51's font is Adwaita Sans 11 (font-name in gsettings-desktop-schemas 51.0); text differs from Breeze UI's in its glyphs and widths, which move the controls after a label, while the controls' sizes and colors compare. A window whose content is at most 360 pixels high makes the alert .short, as it is here.
import importlib
import sys
from pathlib import Path

from common import FRAME, Adw, Gdk, GLib, Graphene, Gtk

OUT = sys.argv[1] if len(sys.argv) > 1 else "."
ONLY = set(sys.argv[2:])
SCALE = 2


def load_sets():
    """The sets of the modules of scenes/, the modules in the order of their names and each one's sets in its own order. Two sets may not share a key."""
    sets, seen = [], {}
    for path in sorted((Path(__file__).parent / "scenes").glob("*.py")):
        if path.stem == "__init__":
            continue
        for key, prefix, build in importlib.import_module("scenes." + path.stem).SETS:
            if key in seen:
                raise SystemExit(f"set {key} is in scenes/{seen[key]}.py and scenes/{path.stem}.py")
            seen[key] = path.stem
            sets.append((key, prefix, build))
    return sets


def render(window, scene, path):
    """Snapshots the window, and the popovers over it at their places, at SCALE and writes it as a PNG."""
    width, height = scene.width + 2 * FRAME, scene.height + 2 * FRAME
    snapshot = Gtk.Snapshot()
    snapshot.scale(SCALE, SCALE)
    Gtk.WidgetPaintable.new(window).snapshot(snapshot, width, height)
    for popover in scene.popovers:
        surface = popover.get_surface()
        if surface is None or not popover.get_mapped():
            continue
        snapshot.save()
        snapshot.translate(Graphene.Point().init(surface.get_position_x(), surface.get_position_y()))
        Gtk.WidgetPaintable.new(popover).snapshot(snapshot, surface.get_width(), surface.get_height())
        snapshot.restore()
    node = snapshot.to_node()
    rect = Graphene.Rect().init(0, 0, width * SCALE, height * SCALE)
    window.get_renderer().render_texture(node, rect).save_to_png(path)
    print("wrote", path)
    for name, widget in scene.named.items():
        ok, bounds = widget.compute_bounds(window)
        if ok:
            # A widget in a popover is measured from the window's surface, its frame included: take the frame off, to print it in the content's DIPs as the others.
            x, y = bounds.get_x(), bounds.get_y()
            if widget.get_ancestor(Gtk.Popover) is not None:
                x, y = x - FRAME, y - FRAME
            print(f"  {name}: x {x:g} y {y:g} w {bounds.get_width():g} h {bounds.get_height():g}")


# Xvfb has no compositor, so GTK gives the window and its popovers the .solid-csd look: a popover's contents square, without their shadows, inside a 1-pixel border. These rules put back the stylesheet's own popover > contents and popover.background (libadwaita 1.10.0's gtk.css), as a desktop with a compositor draws them.
COMPOSITED_POPOVERS = """
.solid-csd popover { background: transparent; border: none; }
.solid-csd popover > contents { border-radius: 15px; background-clip: padding-box; box-shadow: 0 0 0 1px RGB(0 0 0/5%), 0 1px 5px 1px RGB(0 0 0/9%), 0 2px 14px 3px RGB(0 0 0/5%); }
"""


def main():
    Adw.init()
    css = Gtk.CssProvider()
    css.load_from_string(COMPOSITED_POPOVERS)
    Gtk.StyleContext.add_provider_for_display(Gdk.Display.get_default(), css, Gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
    sets = load_sets()
    jobs = []
    for dark in (False, True):
        mode = "dark" if dark else "light"
        for key, prefix, build in sets:
            if not ONLY or key in ONLY:
                jobs.append((dark, f"{prefix}-{mode}.png", build))

    loop = GLib.MainLoop()

    def next_job():
        if not jobs:
            loop.quit()
            return
        dark, file, build = jobs.pop(0)
        # A set may make the default direction right to left; the next starts left to right.
        Gtk.Widget.set_default_direction(Gtk.TextDirection.LTR)
        Adw.StyleManager.get_default().set_color_scheme(Adw.ColorScheme.FORCE_DARK if dark else Adw.ColorScheme.FORCE_LIGHT)
        scene = build()
        window = Adw.Window()
        # AdwWindow asks for at least 360 by 200 pixels, frame included; a scene less than 190 high (the inputs set) would get a taller window, which the snapshot below then squeezes into the scene's height.
        window.set_size_request(-1, -1)
        window.set_content(scene.root)
        window.set_default_size(scene.width + 2 * FRAME, scene.height + 2 * FRAME)
        window.present()
        if scene.shown:
            scene.shown(window)
        if scene.popovers:
            GLib.timeout_add(300, lambda: (scene.open_popovers(), False)[1])

        def shoot():
            render(window, scene, f"{OUT}/{file}")
            if scene.done:
                scene.done()
            for p in scene.popovers:
                p.popdown()
                p.unparent()
            window.destroy()
            GLib.timeout_add(200, lambda: (next_job(), False)[1])
            return False

        # Wait past the dialog's spring animation, the toast's and the spinner's first frames.
        GLib.timeout_add(scene.wait, shoot)

    GLib.idle_add(lambda: (next_job(), False)[1])
    loop.run()


if __name__ == "__main__":
    main()
