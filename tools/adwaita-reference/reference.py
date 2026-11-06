# libadwaita's own controls, for comparing Breeze UI's GNOME look with them.
#
# It renders, with GTK 4 and libadwaita, the controls of Breeze UI's screenshot tests at the same places, light and dark, at 2x, and writes them as PNGs: the sets of ui/breeze/screenshot_test.go (referenceView, referenceView2, the alert and referenceView3) and libadwaita's counterparts of the widgets of TestScreenshotsFeedback, TestScreenshotsLayout, TestScreenshotsNavigation, TestScreenshotsValues, TestScreenshotsCollections, TestScreenshotsRows, TestScreenshotsSettings, TestScreenshotsTabView, TestScreenshotsHeaderBar and TestScreenshotsWave2Forms, Wave2Idioms, Wave2Collections, TestScreenshotsInputs and TestScreenshotsText, laid out in boxes with the paddings and gaps of those tests' rows and columns, so a control sits where Breeze UI's does when both are the same size. It prints every placed control's bounds in window DIPs, and writes popovers, which GTK draws in surfaces of their own, over the window at their places. Breeze UI follows libadwaita 1.10.0 on GTK 4.24 (GNOME 51); it was run with GTK 4.24.1, libadwaita 1.10.0, PyGObject 3.58.1 and the Adwaita Sans font on Arch Linux, under Xvfb:
#
#	xvfb-run -a -s "-screen 0 1600x1200x24" python3 tools/adwaita-reference/reference.py <dir> [set...]
#
# The sets are first, second, alert, third, feedback, layout, nav, values, collections, forms, idioms, menu (idioms with a popover menu open), table, inputs, status, statuspage, rows, overcard, spinlimit, settings, tabview, text, headerbar and headerbarsub; without any, it renders them all. The window's frame offsets its content by 5 pixels (10 at 2x) from Breeze UI's screenshots, and the window is 10 pixels wider and higher than the test's so that its content is the test's size. GNOME 51's font is Adwaita Sans 11 (font-name in gsettings-desktop-schemas 51.0); text differs from Breeze UI's in its glyphs and widths, which move the controls after a label, while the controls' sizes and colors compare. A window whose content is at most 360 pixels high makes the alert .short, as it is here.
import sys

import gi

gi.require_version("Gtk", "4.0")
gi.require_version("Adw", "1")
gi.require_version("Gdk", "4.0")
gi.require_version("Graphene", "1.0")
from gi.repository import Adw, Gdk, GLib, Gio, Graphene, Gtk  # noqa: E402

OUT = sys.argv[1] if len(sys.argv) > 1 else "."
ONLY = set(sys.argv[2:])
SCALE = 2
FRAME = 5


class Scene:
    """A window's content: its size, its root widget, the widgets whose bounds to print, and what to do once it is shown (open a popover, take the focus)."""

    def __init__(self, width, height, root, named=None, shown=None, popovers=None, open_popovers=None):
        self.width, self.height, self.root = width, height, root
        self.named = named or {}
        self.shown = shown
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


# ---- the sets of ui/breeze/screenshot_test.go ----


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


# ---- libadwaita's counterparts of the other screenshot tests' widgets ----


def feedback_set():
    """TestScreenshotsFeedback's feedbackView: a search bar for the find bar, a button, level bars at their three levels, avatars and a toast with a button. GTK's tooltip is a surface that shows only under a pointer and is left out; GNOME has no rating control."""
    search = Gtk.SearchEntry(text="cat")
    search.set_hexpand(True)
    bar = Gtk.SearchBar(child=search, search_mode_enabled=True, show_close_button=True)

    def level(value):
        lb = Gtk.LevelBar(min_value=0, max_value=100, value=value)
        # Breeze UI's Meter: the warning level from 70 and the critical one from 90, as GtkLevelBar's offsets low and high.
        lb.add_offset_value(Gtk.LEVEL_BAR_OFFSET_LOW, 70)
        lb.add_offset_value(Gtk.LEVEL_BAR_OFFSET_HIGH, 90)
        lb.add_offset_value(Gtk.LEVEL_BAR_OFFSET_FULL, 100)
        lb.set_size_request(220, -1)
        return lb

    share = Gtk.Button(label="Share")
    levels = [level(40), level(75), level(95)]
    avatars = [Adw.Avatar(size=32, text=name, show_initials=True) for name in ["Ada Lovelace", "Grace Hopper", "Linus", "山田太郎", ""]]
    content = Gtk.Fixed()
    content.set_vexpand(True)
    content.put(share, 16, 16)
    for i, lb in enumerate(levels):
        # A meter's bar is centered in its 20-pixel row (the tests place the rows 28 apart).
        content.put(lb, 300, 16 + 28 * i + 6)
    for i, a in enumerate(avatars):
        content.put(a, 16 + 48 * i, 110)
    overlay = Adw.ToastOverlay(child=vbox(0, bar, content))
    toast = Adw.Toast(title="Note deleted", button_label="Undo", timeout=0)
    named = {"searchbar": bar, "share": share, "level40": levels[0], "level75": levels[1], "level95": levels[2], "avatar": avatars[0]}
    scene = Scene(560, 400, overlay, named=named)
    scene.shown = lambda window: overlay.add_toast(toast)
    return scene


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


def nav_set():
    """TestScreenshotsNavigation's navView: an AdwSidebar of two sections with Sent chosen, and a toolbar of the history buttons, a toggle and a button. Inbox's count is left out: AdwSidebar styles a count only in AdwViewSwitcherSidebar (view-switcher-sidebar .indicator), which has no sections. GNOME has no path bar control; Nautilus draws its own."""
    sidebar = Adw.Sidebar()
    mail = Adw.SidebarSection(title="Mailboxes")
    for title in ["Inbox", "Sent", "Drafts"]:
        mail.append(Adw.SidebarItem(title=title, icon_name="folder-symbolic"))
    smart = Adw.SidebarSection(title="Smart Mailboxes")
    smart.append(Adw.SidebarItem(title="Flagged"))
    sidebar.append(mail)
    sidebar.append(smart)
    sidebar.set_selected(1)
    # As in an AdwNavigationSplitView's sidebar: sidebar_bg_color and the line at its end. The sidebar and the content both ask for their share of the window, so the sidebar is 200 wide, as the test's.
    sidebar.add_css_class("sidebar-pane")
    sidebar.set_size_request(200, -1)

    back = Gtk.Button(icon_name="go-previous-symbolic")
    forward = Gtk.Button(icon_name="go-next-symbolic", sensitive=False)
    bold = Gtk.ToggleButton(label="Bold")
    spacer = Gtk.Box(hexpand=True)
    share = Gtk.Button(label="Share")
    toolbar = Gtk.Box()
    toolbar.add_css_class("toolbar")
    for w in [back, forward, bold, spacer, share]:
        toolbar.append(w)
    content = vbox(12, toolbar, padding=(0, 12, 0, 12))
    content.set_hexpand(True)
    content.set_size_request(560 - 200 - 24, -1)
    root = Gtk.Box()
    root.append(sidebar)
    root.append(content)
    return Scene(560, 260, root, named={"sidebar": sidebar, "toolbar": toolbar, "back": back, "bold": bold, "share": share})


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


def forms_set():
    """TestScreenshotsWave2Forms's formsShotView: check button groups, spin buttons (one at its maximum), a label above an entry, a preferences group, entries for the time and the date, and a paned in a frame."""
    def group(head, children, inconsistent):
        heads = check(head, not inconsistent, inconsistent=inconsistent)
        rows = [start(heads)] + [pad(start(check(label, on)), (0, 0, 0, 24)) for label, on in children]
        return vbox(0, *rows)

    def spin(value):
        s = Gtk.SpinButton.new_with_range(0, 10, 1)
        s.set_value(value)
        return start(s)

    left = vbox(16, group("Notifications", [("Mail", True), ("Calendar", False), ("Photos", False)], True), group("Sync", [("Contacts", True), ("Notes", True)], False), spin(5), spin(10), width=260)
    name = vbox(6, start(Gtk.Label(label="Name")), entry("Ada"))
    shipping = Adw.PreferencesGroup(title="Shipping")
    shipping.add(Adw.EntryRow(title="Address"))
    paned = Gtk.Paned(start_child=pad(Gtk.Label(label="Sidebar", xalign=0, yalign=0), 8), end_child=pad(Gtk.Label(label="Content", xalign=0, yalign=0), 8), position=120)
    paned.add_css_class("frame")
    paned.set_size_request(-1, 60)
    right = vbox(16, name, shipping, hbox(12, entry("15:04"), entry("1/5/2026"), align=Gtk.Align.CENTER), paned, width=300)
    root = hbox(24, left, right, padding=16)
    return Scene(640, 440, root, named={"spin": left.get_last_child(), "shipping": shipping, "paned": paned})


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


def status_set():
    """TestScreenshotsStatus's statusView: four banners, 16 from the window's edges and 68 apart, the first and the last with a Dismiss button (banners have no severity, so the four differ only in their titles), and AdwSpinners where the determinate rings are, as libadwaita has no determinate ring. A banner has only a title: Breeze UI's message follows it in the same label at regular weight, which markup gives the bold title here."""
    items, named = [], {}
    for i, title in enumerate(["Informational", "Success", "Warning", "Error"]):
        banner = Adw.Banner(title=title + ' <span weight="normal">Some message about it.</span>', use_markup=True, revealed=True)
        if i in (0, 3):
            banner.set_button_label("Dismiss")
        banner.set_size_request(528, -1)
        items.append((banner, 16, 16 + 68 * i))
        named["banner-" + title.lower()] = banner
    for i in range(4):
        spinner = Adw.Spinner()
        spinner.set_size_request(32, 32)
        items.append((spinner, 16 + 48 * i, 300))
    small = Adw.Spinner()
    small.set_size_request(16, 16)
    items.append((small, 208, 308))
    return fixed(560, 360, items, named=named)


def status_page_set():
    """TestScreenshotsStatus's statusPageView: an AdwStatusPage filling the window with system-search-symbolic, a title, a description and a button."""
    button = Gtk.Button(label="Clear Search", halign=Gtk.Align.CENTER)
    page = Adw.StatusPage(icon_name="system-search-symbolic", title="No Results", description="Check the spelling or try a new search.", child=button)
    named = {"statuspage": page, "button": button}
    # The page's icon, title and description, which it does not expose, by their types and style classes (adw-status-page.ui).
    for w in descendants(page):
        if isinstance(w, Gtk.Image):
            named["icon"] = w
        elif isinstance(w, Gtk.Label) and w.has_css_class("title"):
            named["title"] = w
        elif isinstance(w, Gtk.Label) and w.has_css_class("description"):
            named["description"] = w
    return Scene(560, 420, page, named=named)


def descendants(widget):
    """Every widget inside widget, depth first."""
    child = widget.get_first_child()
    while child is not None:
        yield child
        yield from descendants(child)
        child = child.get_next_sibling()


def rows_set():
    """TestScreenshotsRows's rowsView: a .card of text padded 12, a boxed list of an action row with a button, a switch row with a subtitle and one without, an action row with a dimmed value and a disabled switch row each alone in a boxed list, a preferences group of rows for the fieldset, and an expander row with its nested rows for the accordion."""
    def boxed(*rows):
        box = Gtk.ListBox(selection_mode=Gtk.SelectionMode.NONE)
        box.add_css_class("boxed-list")
        for row in rows:
            box.append(row)
        return box

    def switch_row(title, active, subtitle="", sensitive=True):
        row = Adw.SwitchRow(title=title, subtitle=subtitle, active=active)
        row.set_sensitive(sensitive)
        return row

    text = Gtk.Label(label="A card holds what belongs together.", xalign=0, wrap=True)
    card = Gtk.Box()
    card.add_css_class("card")
    card.append(pad(text, 12))
    wallpaper = Adw.ActionRow(title="Wallpaper", subtitle="Choose a picture")
    open_button = Gtk.Button(label="Open", valign=Gtk.Align.CENTER)
    wallpaper.add_suffix(open_button)
    settings = boxed(wallpaper, switch_row("Dark mode", True, "Use a dark appearance"), switch_row("Notifications", False))
    storage = Adw.ActionRow(title="Storage")
    value = Gtk.Label(label="120 GB")
    value.add_css_class("dimmed")
    storage.add_suffix(value)
    left = vbox(16, card, settings, boxed(storage), boxed(switch_row("Bluetooth", True, sensitive=False)), width=300)

    privacy = Adw.PreferencesGroup(title="Privacy")
    privacy.add(switch_row("Location", True, "Apps may ask for it"))
    privacy.add(Adw.ActionRow(title="Camera", subtitle="2 apps"))
    network = Adw.ExpanderRow(title="Network", expanded=True)
    network.add_row(Adw.ActionRow(title="Proxy", subtitle="Automatic"))
    network.add_row(switch_row("Metered", False))
    accordion = boxed(network)
    right = vbox(16, privacy, accordion, width=300)
    # The rows ask to expand, which would widen the column past the test's 300 to the window's end.
    right.set_hexpand(False)
    root = hbox(24, left, right, padding=16)
    return Scene(672, 520, root, named={"card": card, "rows": settings, "storage": storage, "privacy": privacy, "network": network})


def settings_set():
    """TestScreenshotsSettings's settingsView: in a 300-pixel column, a boxed list of an AdwComboRow with a subtitle, an AdwEntryRow holding a name, an empty one and an AdwSpinRow with a subtitle, then an AdwComboRow and a disabled AdwSpinRow each alone in a boxed list; beside them an open AdwExpanderRow with a subtitle holding an action row and a switch row, and a closed one, each in a boxed list."""
    def boxed(*rows):
        box = Gtk.ListBox(selection_mode=Gtk.SelectionMode.NONE)
        box.add_css_class("boxed-list")
        for row in rows:
            box.append(row)
        return box

    def combo(title, items, selected, subtitle=""):
        return Adw.ComboRow(title=title, subtitle=subtitle, model=strings(items), selected=selected)

    def entry(title, text):
        row = Adw.EntryRow(title=title)
        row.set_text(text)
        return row

    def spin(title, value, lower, upper, subtitle="", sensitive=True):
        row = Adw.SpinRow.new_with_range(lower, upper, 1)
        row.set_title(title)
        row.set_subtitle(subtitle)
        row.set_value(value)
        row.set_sensitive(sensitive)
        return row

    settings = boxed(combo("Appearance", ["Light", "Dark", "Auto"], 2, "Light or dark"), entry("Name", "Evan"), entry("Email", ""), spin("Font size", 12, 6, 72, "Points"))
    language = combo("Language", ["English", "Français"], 0)
    copies = spin("Copies", 3, 1, 9, sensitive=False)
    left = vbox(16, settings, boxed(language), boxed(copies), width=300)

    network = Adw.ExpanderRow(title="Network", subtitle="Wi-Fi and proxies", expanded=True)
    network.add_row(Adw.ActionRow(title="Proxy", subtitle="Automatic"))
    network.add_row(Adw.SwitchRow(title="Metered"))
    advanced = Adw.ExpanderRow(title="Advanced")
    advanced.add_row(Adw.ActionRow(title="Logs"))
    right = vbox(16, boxed(network), boxed(advanced), width=300)
    # The rows ask to expand, which would widen the columns past the test's 300 to the window's end.
    left.set_hexpand(False)
    right.set_hexpand(False)
    root = hbox(24, left, right, padding=16)
    return Scene(672, 480, root, named={"settings": settings, "language": language, "copies": copies, "network": network, "advanced": advanced})


def switcher_set():
    """TestScreenshotsViewSwitcher's switcherView: an AdwHeaderBar at the top of the window's content titled Music with an AdwViewSwitcher of the wide policy packed at its end, where Breeze UI's HeaderBar holds a ViewSwitcher among its end controls (libadwaita apps put it in place of the title, as the header bar's title-widget, which Breeze UI's HeaderBar does not take), GTK's close button after it (decoration-layout :close), and below the bar a wide AdwViewSwitcher alone, 16 pixels in from the window's side and 16 below the bar; each switches an AdwViewStack of Library, Browse and Radio with icons, Browse visible. The stacks are hidden and take no room: only the switchers show."""
    def switcher():
        stack = Adw.ViewStack()
        for title, icon in (("Library", "folder-music-symbolic"), ("Browse", "view-grid-symbolic"), ("Radio", "network-wireless-symbolic")):
            stack.add_titled_with_icon(Gtk.Label(label=title), title.lower(), title, icon)
        stack.set_visible_child_name("browse")
        stack.set_visible(False)
        return Adw.ViewSwitcher(stack=stack, policy=Adw.ViewSwitcherPolicy.WIDE), stack

    bar_switcher, bar_stack = switcher()
    bar = Adw.HeaderBar(title_widget=Adw.WindowTitle(title="Music", subtitle=""), decoration_layout=":close")
    bar.pack_end(bar_switcher)
    alone, alone_stack = switcher()
    root = vbox(0, bar, vbox(0, start(alone), alone_stack, bar_stack, padding=16))
    return Scene(560, 140, root, named={"headerbar": bar, "bar-switcher": bar_switcher, "switcher": alone})


def tabview_set():
    """TestScreenshotsTabView's tabViewView: an AdwTabBar 528 wide at 16, 16 over an AdwTabView of Documents (chosen), Downloads and Pictures, a flat button with adw-tab-new-symbolic as its end action, the chosen tab's text 12 below the 47-pixel bar, and at 16, 136 a bar of General and Advanced (chosen) without the end action. The pages' close buttons show on the chosen tab and under the pointer, as AdwTab shows them; the second bar's tabs close too, as AdwTabBar has no tabs that do not, while Breeze UI's tab view built without close shows none. The views take no room: only their bars show."""
    def tab_bar(titles, chosen, end_action):
        view = Adw.TabView()
        pages = [view.append(Gtk.Label(label=title)) for title in titles]
        for page, title in zip(pages, titles):
            page.set_title(title)
        view.set_selected_page(pages[chosen])
        view.set_visible(False)
        bar = Adw.TabBar(view=view, autohide=False)
        if end_action:
            new_tab = Gtk.Button(icon_name="adw-tab-new-symbolic", tooltip_text="New Tab")
            new_tab.add_css_class("flat")
            bar.set_end_action_widget(new_tab)
        bar.set_size_request(528, -1)
        return vbox(0, bar, view), bar

    docs, docs_bar = tab_bar(["Documents", "Downloads", "Pictures"], 0, True)
    content = Gtk.Label(label="Showing Documents")
    prefs, prefs_bar = tab_bar(["General", "Advanced"], 1, False)
    return fixed(560, 260, [(docs, 16, 16), (content, 16, 75), (prefs, 16, 136)], named={"docs": docs_bar, "content": content, "prefs": prefs_bar})


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


def spin_limit_set():
    """The spin row of TestAdwPlusEvenWhenDisabled: in a boxed list 300 wide, 16 from the window's edges, an AdwSpinRow Copies at the upper limit of 1 to 9, its + button insensitive, to read the icon of a disabled button. Nothing keeps the focus, which the row would take as the window's only focusable widget."""
    box = Gtk.ListBox(selection_mode=Gtk.SelectionMode.NONE)
    box.add_css_class("boxed-list")
    row = Adw.SpinRow.new_with_range(1, 9, 1)
    row.set_title("Copies")
    row.set_value(9)
    box.append(row)
    root = vbox(0, box, padding=16, width=300)
    root.set_hexpand(False)
    return Scene(332, 120, start(root), named={"row": row}, shown=lambda window: window.set_focus(None))


def text_set():
    """TestScreenshotsText's textView: labels in the style classes of Breeze UI's Text steps (.title-1, .title-2, .title-4, .heading, a plain label, .caption), the last two .dimmed, where textTops places them."""
    rows = [("Large Title", ["title-1"], 16), ("Title", ["title-2"], 76), ("Subtitle", ["title-4"], 122), ("Heading", ["heading"], 160), ("Body", [], 188), ("Caption", ["caption"], 214), ("Secondary body", ["dimmed"], 240), ("Secondary caption", ["caption", "dimmed"], 266)]
    items, named = [], {}
    for text, classes, top in rows:
        label = Gtk.Label(label=text, css_classes=classes)
        items.append((label, 16, top))
        named[text.lower().replace(" ", "-")] = label
    return fixed(360, 300, items, named=named)


SETS = [
    ("first", "adwaita", first_set), ("second", "adwaita2", second_set), ("alert", "adwaita-alert", alert_set), ("third", "adwaita3", third_set),
    ("feedback", "adwaita-feedback", feedback_set), ("layout", "adwaita-layout", layout_set), ("nav", "adwaita-nav", nav_set),
    ("values", "adwaita-values", values_set), ("collections", "adwaita-collections", collections_set), ("forms", "adwaita-forms", forms_set),
    ("idioms", "adwaita-idioms", idioms_set), ("menu", "adwaita-idioms-menu", lambda: idioms_set(True)), ("table", "adwaita-table", table_set),
    ("inputs", "adwaita-inputs", inputs_set),
    ("status", "adwaita-status", status_set), ("statuspage", "adwaita-statuspage", status_page_set),
    ("rows", "adwaita-rows", rows_set), ("overcard", "adwaita-overcard", over_card_set), ("spinlimit", "adwaita-spinlimit", spin_limit_set),
    ("settings", "adwaita-settings", settings_set), ("tabview", "adwaita-tabview", tabview_set),
    ("switcher", "adwaita-switcher", switcher_set),
    ("text", "adwaita-text", text_set),
]


def headerbar_set(subtitle=""):
    """TestScreenshotsHeaderBar's headerView, the sets headerbar and headerbarsub: an AdwHeaderBar at the top of the window's content, outside a window's titlebar slot and an AdwToolbarView as Breeze UI's HeaderBar is drawn, its title an AdwWindowTitle reading Documents (with the subtitle 12 items in headerbarsub), go-previous-symbolic packed at its start and Edit and send-to-symbolic at its end, GTK's close button after them (decoration-layout :close, the window controls headerTestRoom gives room for on GNOME), then the window's content padded 16."""
    bar = Adw.HeaderBar(title_widget=Adw.WindowTitle(title="Documents", subtitle=subtitle), decoration_layout=":close")
    back = Gtk.Button(icon_name="go-previous-symbolic")
    edit = Gtk.Button(label="Edit")
    share = Gtk.Button(icon_name="send-to-symbolic")
    bar.pack_start(back)
    bar.pack_end(share)
    bar.pack_end(edit)
    content = Gtk.Label(label="Content", xalign=0)
    root = vbox(0, bar, pad(content, 16))
    return Scene(560, 120, root, named={"headerbar": bar, "back": back, "edit": edit, "share": share})


SETS += [("headerbar", "adwaita-headerbar", headerbar_set), ("headerbarsub", "adwaita-headerbar-subtitle", lambda: headerbar_set("12 items"))]


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
            print(f"  {name}: x {bounds.get_x():g} y {bounds.get_y():g} w {bounds.get_width():g} h {bounds.get_height():g}")


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
    jobs = []
    for dark in (False, True):
        mode = "dark" if dark else "light"
        for key, prefix, build in SETS:
            if not ONLY or key in ONLY:
                jobs.append((dark, f"{prefix}-{mode}.png", build))

    loop = GLib.MainLoop()

    def next_job():
        if not jobs:
            loop.quit()
            return
        dark, file, build = jobs.pop(0)
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
            for p in scene.popovers:
                p.popdown()
                p.unparent()
            window.destroy()
            GLib.timeout_add(200, lambda: (next_job(), False)[1])
            return False

        # Wait past the dialog's spring animation, the toast's and the spinner's first frames.
        GLib.timeout_add(1500, shoot)

    GLib.idle_add(lambda: (next_job(), False)[1])
    loop.run()


main()
