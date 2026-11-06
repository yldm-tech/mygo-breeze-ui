# status: libadwaita's counterparts of the widgets of TestScreenshotsStatus (ui/breeze/widgets_status_test.go).
from common import Adw, Gtk, Scene, fixed, descendants


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


def status_actions_set():
    """TestScreenshotsStatus's statusActionsView: banners whose button is the bar's action (button-label), 16 from the window's edges and 68 apart: Open Settings, Restart as a suggested action (button-style, libadwaita 1.7), Undo on a banner without a message, and Retry on a banner 280 wide, which puts its button below the title. Breeze UI's message follows the title in the same label at regular weight, as in status_set."""
    items, named = [], {}
    bars = [("Warning", "Some message about it.", "Open Settings", False, 528), ("Informational", "Some message about it.", "Restart", True, 528), ("Success", "", "Undo", False, 528), ("Error", "Some message about it.", "Retry", False, 280)]
    for i, (title, message, action, suggested, width) in enumerate(bars):
        markup = title + (' <span weight="normal">' + message + "</span>" if message else "")
        banner = Adw.Banner(title=markup, use_markup=True, revealed=True, button_label=action)
        if suggested:
            banner.set_button_style(Adw.BannerButtonStyle.SUGGESTED)
        banner.set_size_request(width, -1)
        items.append((banner, 16, 16 + 68 * i))
        named["actions-" + title.lower()] = banner
    return fixed(560, 380, items, named=named)


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


SETS = [
    ("status", "adwaita-status", status_set),
    ("statuspage", "adwaita-statuspage", status_page_set),
    ("statusactions", "adwaita-statusactions", status_actions_set),
]
