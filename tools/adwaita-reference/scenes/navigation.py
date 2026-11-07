# navigation: libadwaita's counterparts of the widgets of TestScreenshotsNavigation (ui/breeze/widgets_navigation_test.go).
from common import Adw, Gtk, Scene, vbox


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


SETS = [
    ("nav", "adwaita-nav", nav_set),
]
