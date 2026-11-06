# collections: libadwaita's counterparts of the widgets of TestScreenshotsCollections (ui/breeze/widgets_collections_test.go).
from common import Gtk, Scene, vbox, hbox, start, entry, strings, label_factory, tree_model, tree_factory


def collections_set():
    """TestScreenshotsCollections's collectionsView: rich lists, a button and entries, an editable label, a grid view and a tree list. GNOME has no token field."""
    files = ["Notes.md", "Budget.xlsx", "Photo.heic", "Talk.key"]

    def rich_list(selected):
        sel = Gtk.SingleSelection(model=strings(files), selected=selected)
        view = Gtk.ListView(model=sel, factory=label_factory())
        view.add_css_class("rich-list")
        sw = Gtk.ScrolledWindow(child=view, hscrollbar_policy=Gtk.PolicyType.NEVER)
        sw.set_size_request(-1, 160)
        return sw, view

    files_list, focused = rich_list(1)
    other_list, _ = rich_list(2)
    left = vbox(12, files_list, other_list, width=240)

    text = entry(placeholder="Text input")
    text.set_hexpand(True)
    title = start(Gtk.EditableLabel(text="Notes.md"))
    grid_sel = Gtk.SingleSelection(model=strings([f[:4] for f in files]), selected=1)
    grid = Gtk.GridView(model=grid_sel, factory=label_factory(xalign=0), max_columns=4, min_columns=4)
    grid_sw = Gtk.ScrolledWindow(child=grid)
    grid_sw.set_size_request(-1, 90)
    outline_sel = Gtk.SingleSelection(model=tree_model(["src", "go.mod"], lambda k: ["main.go", "view.go"] if k == "src" else None), selected=0)
    outline = Gtk.ListView(model=outline_sel, factory=tree_factory())
    outline_sw = Gtk.ScrolledWindow(child=outline)
    outline_sw.set_size_request(-1, 110)
    right = vbox(12, hbox(8, Gtk.Button(label="Button"), text, align=Gtk.Align.CENTER), entry("Berlin"), entry(), title, grid_sw, outline_sw, width=300)
    root = hbox(16, left, right, padding=16)
    scene = Scene(600, 440, root, named={"list": files_list, "grid": grid_sw, "outline": outline_sw, "editable": title})
    scene.shown = lambda window: focused.grab_focus()
    return scene


SETS = [
    ("collections", "adwaita-collections", collections_set),
]
