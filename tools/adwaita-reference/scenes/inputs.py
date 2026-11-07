# inputs: libadwaita's counterparts of the widgets of TestScreenshotsInputs (ui/breeze/widgets_inputs_test.go), and the card of TestAdwControlsOverCard.
from common import Adw, Gio, Gtk, Scene, fixed, vbox, start


def inputs_set():
    """TestScreenshotsInputs's inputsView: GtkPasswordEntrys, empty with a placeholder and holding a password with the focus, an AdwSplitButton, image buttons, the second insensitive, and below them a .toolbar of 520 pixels holding an AdwSplitButton and two image buttons, flat there."""
    model = Gio.Menu()
    for item in ["Save As…", "Save All"]:
        model.append(item, "menu.item")
    actions = Gio.SimpleActionGroup()
    actions.add_action(Gio.SimpleAction(name="item"))
    empty = Gtk.PasswordEntry(placeholder_text="Password")
    empty.set_size_request(180, -1)
    secret = Gtk.PasswordEntry(text="hunter2")
    secret.set_size_request(180, -1)
    trash = Gtk.Button(icon_name="user-trash-symbolic", sensitive=False)
    toolbar = Gtk.Box()
    toolbar.add_css_class("toolbar")
    toolbar.set_size_request(520, -1)
    for w in [Adw.SplitButton(label="New", menu_model=model), Gtk.Button(icon_name="send-to-symbolic"), Gtk.Button(icon_name="user-trash-symbolic")]:
        toolbar.append(w)
    split = Adw.SplitButton(label="Save", menu_model=model)
    share = Gtk.Button(icon_name="send-to-symbolic")
    items = [(empty, 16, 16), (secret, 220, 16), (split, 16, 64), (share, 160, 64), (trash, 220, 64), (toolbar, 0, 110)]
    scene = fixed(520, 180, items, named={"password": secret, "split": split, "share": share, "toolbar": toolbar})
    scene.root.insert_action_group("menu", actions)
    scene.shown = lambda window: secret.grab_focus()
    return scene


def over_card_set():
    """The card of TestAdwControlsOverCard: a boxed list 300 wide, 16 from the window's edges, of an AdwActionRow holding an AdwSplitButton and one holding an empty GtkPasswordEntry with a placeholder, 160 wide, to read what the translucent faces of libadwaita's controls show over a card."""
    box = Gtk.ListBox(selection_mode=Gtk.SelectionMode.NONE)
    box.add_css_class("boxed-list")
    share = Adw.ActionRow(title="Share")
    split = Adw.SplitButton(label="Send", valign=Gtk.Align.CENTER, menu_model=Gio.Menu())
    share.add_suffix(split)
    secret = Adw.ActionRow(title="Secret")
    entry = Gtk.PasswordEntry(placeholder_text="Password", valign=Gtk.Align.CENTER)
    entry.set_size_request(160, -1)
    secret.add_suffix(entry)
    box.append(share)
    box.append(secret)
    root = vbox(0, box, padding=16, width=300)
    root.set_hexpand(False)
    return Scene(332, 160, start(root), named={"split": split, "password": entry})


SETS = [
    ("inputs", "adwaita-inputs", inputs_set),
    ("overcard", "adwaita-overcard", over_card_set),
]
