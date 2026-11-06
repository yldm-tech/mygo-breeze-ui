# drawn_collections: libadwaita's counterparts of the widgets of TestScreenshotsWave2Collections (ui/breeze/drawn_collections_test.go).
from common import Gtk, Scene, vbox, hbox, start, strings, label_factory, tree_model, tree_factory


def table_set():
    """TestScreenshotsWave2Collections's collView: column views for the table and the outline table, a tree list, a grid view, another tree list and an editable label."""
    files = ["Notes.md", "Budget.xlsx", "Photo.heic", "Talk.key"]

    def column_view(model, columns, selected, height):
        sel = Gtk.SingleSelection(model=model, selected=selected)
        view = Gtk.ColumnView(model=sel)
        for title, factory, width in columns:
            col = Gtk.ColumnViewColumn(title=title, factory=factory)
            if width:
                col.set_fixed_width(width)
            else:
                col.set_expand(True)
            view.append_column(col)
        sw = Gtk.ScrolledWindow(child=view, hscrollbar_policy=Gtk.PolicyType.NEVER)
        sw.set_size_request(-1, height)
        return sw

    def plain_factory(f):
        factory = Gtk.SignalListItemFactory()
        factory.connect("setup", lambda _, item: item.set_child(Gtk.Label(xalign=0)))
        factory.connect("bind", lambda _, item: item.get_child().set_label(f(item.get_item().get_string())))
        return factory

    table = column_view(strings(files), [("Name", label_factory(), 0), ("Kind", plain_factory(lambda s: s.rsplit(".", 1)[1]), 70)], 1, 170)

    def children(k):
        return {"src": ["ui", "main.go"], "ui": ["view.go"]}.get(k)

    outline_table = column_view(tree_model(["src", "go.mod"], children), [("Name", tree_factory(), 0), ("Size", tree_factory(lambda k: f"{len(k)} KB"), 60)], 0, 200)
    left = vbox(12, table, outline_table, width=260)

    outline_sel = Gtk.SingleSelection(model=tree_model(["src", "go.mod"], children), selected=1)
    outline = Gtk.ScrolledWindow(child=Gtk.ListView(model=outline_sel, factory=tree_factory()))
    outline.set_size_request(-1, 130)
    colors = ["Red", "Green", "Blue", "Gray"]
    grid = Gtk.ScrolledWindow(child=Gtk.GridView(model=Gtk.SingleSelection(model=strings(colors), selected=1), factory=label_factory(), max_columns=4, min_columns=4))
    grid.set_size_request(-1, 90)
    tree_children = {"Library": ["Music", "Podcasts"], "Music": ["Songs.m4a"]}
    tree_sel = Gtk.SingleSelection(model=tree_model(["Library", "Movies"], tree_children.get), selected=3)
    tree = Gtk.ListView(model=tree_sel, factory=tree_factory())
    right = vbox(12, outline, grid, tree, start(Gtk.EditableLabel(text="Readme.md")), width=220)
    root = hbox(12, left, right, padding=12)
    return Scene(520, 500, root, named={"table": table, "outlinetable": outline_table, "outline": outline, "grid": grid, "tree": tree})


SETS = [
    ("table", "adwaita-table", table_set),
]
