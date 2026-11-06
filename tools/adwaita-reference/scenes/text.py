# text: libadwaita's counterparts of the text of TestScreenshotsText (ui/breeze/widgets_text_test.go).
from common import Gtk, fixed


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
    ("text", "adwaita-text", text_set),
]
