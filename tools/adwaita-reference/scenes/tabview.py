# tabview: libadwaita's counterpart of TestScreenshotsTabView (ui/breeze/widgets_tabview_test.go).
from common import Adw, Gtk, fixed, vbox


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


SETS = [
    ("tabview", "adwaita-tabview", tabview_set),
]
