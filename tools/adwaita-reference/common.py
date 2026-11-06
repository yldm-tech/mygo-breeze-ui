# What the sets of scenes/ build their scenes with: GTK and libadwaita, and helpers that lay out widgets as Breeze UI's screenshot tests do.
import gi

gi.require_version("Gtk", "4.0")
gi.require_version("Adw", "1")
gi.require_version("Gdk", "4.0")
gi.require_version("Graphene", "1.0")
from gi.repository import Adw, Gdk, GLib, Gio, Graphene, Gtk  # noqa: E402, F401

# The window's frame, by which its content is offset from Breeze UI's screenshots: 5 pixels, 10 at 2x.
FRAME = 5


class Scene:
    """A window's content: its size, its root widget, the widgets whose bounds to print, and what to do once it is shown (open a popover, take the focus)."""

    def __init__(self, width, height, root, named=None, shown=None, popovers=None, open_popovers=None):
        self.width, self.height, self.root = width, height, root
        self.named = named or {}
        self.shown = shown
        # How long after the window is presented it is rendered, in milliseconds, and what to do once it is (let go of the mouse a set pressed).
        self.wait = 1500
        self.done = None
        self.popovers = popovers or []
        # What opens the popovers, once the window is up: by default each popover's popup.
        self.open_popovers = open_popovers or (lambda: [p.popup() for p in self.popovers])


def fixed(width, height, items, popovers=None, named=None):
    box = Gtk.Fixed()
    for widget, x, y in items:
        box.put(widget, x, y)
    return Scene(width, height, box, named=named, popovers=popovers)


def vbox(gap, *children, padding=0, width=-1):
    box = Gtk.Box(orientation=Gtk.Orientation.VERTICAL, spacing=gap)
    pad(box, padding)
    if width > 0:
        box.set_size_request(width, -1)
    for child in children:
        box.append(child)
    return box


def hbox(gap, *children, padding=0, align=Gtk.Align.START):
    box = Gtk.Box(orientation=Gtk.Orientation.HORIZONTAL, spacing=gap)
    pad(box, padding)
    for child in children:
        child.set_valign(align)
        box.append(child)
    return box


def pad(widget, padding):
    if isinstance(padding, int):
        padding = (padding, padding, padding, padding)
    top, end, bottom, start = padding
    widget.set_margin_top(top)
    widget.set_margin_end(end)
    widget.set_margin_bottom(bottom)
    widget.set_margin_start(start)
    return widget


def start(widget):
    """Keeps widget at its natural width in a vertical box, as a widget in a column aligned to the start."""
    widget.set_halign(Gtk.Align.START)
    return widget


def check(label, active, sensitive=True, inconsistent=False):
    c = Gtk.CheckButton(label=label, active=active)
    c.set_inconsistent(inconsistent)
    c.set_sensitive(sensitive)
    return c


def entry(text="", placeholder="", width=-1):
    e = Gtk.Entry(text=text, placeholder_text=placeholder)
    e.set_size_request(width, -1)
    return e


def strings(items):
    return Gtk.StringList.new(items)


def label_factory(css=None, xalign=0):
    factory = Gtk.SignalListItemFactory()

    def setup(_, item):
        label = Gtk.Label(xalign=xalign)
        if css:
            label.add_css_class(css)
        item.set_child(label)

    factory.connect("setup", setup)
    factory.connect("bind", lambda _, item: item.get_child().set_label(item.get_item().get_string()))
    return factory


def tree_model(roots, children):
    """A Gtk.TreeListModel of strings: roots and, for an item, children(item), all expanded."""

    def create(item):
        kids = children(item.get_string())
        return strings(kids) if kids else None

    return Gtk.TreeListModel.new(strings(roots), False, True, create)


def tree_factory(column=None):
    """A factory of rows with a Gtk.TreeExpander around a label; column, when set, labels another column of the row's item."""
    factory = Gtk.SignalListItemFactory()

    def setup(_, item):
        label = Gtk.Label(xalign=0)
        if column is None:
            expander = Gtk.TreeExpander(child=label)
            item.set_child(expander)
        else:
            item.set_child(label)

    def bind(_, item):
        row = item.get_item()
        text = row.get_item().get_string()
        if column is None:
            item.get_child().set_list_row(row)
            item.get_child().get_child().set_label(text)
        else:
            item.get_child().set_label(column(text))

    factory.connect("setup", setup)
    factory.connect("bind", bind)
    return factory


def descendants(widget):
    """Every widget inside widget, depth first."""
    child = widget.get_first_child()
    while child is not None:
        yield child
        yield from descendants(child)
        child = child.get_next_sibling()
