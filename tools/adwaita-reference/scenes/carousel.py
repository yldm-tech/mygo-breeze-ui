# carousel: libadwaita's counterparts of TestScreenshotsCarousel (ui/breeze/widgets_carousel_test.go).
from common import Adw, Gtk, fixed, vbox


def carousel_set(rtl=False):
    """TestScreenshotsCarousel's carouselShotView: AdwCarouselIndicatorDots of 5 pages with the 2nd chosen at 16, 16 and of 9 pages with the 7th chosen at 16, 56, each showing a carousel of empty pages that the scene keeps out of sight, 1 pixel square and transparent (the dots read its pages and position once it is laid out), and at 16, 96 a vertical GtkBox 320 wide of an AdwCarousel 160 high of four pages, each a label of its slide centered, the 2nd shown, over AdwCarouselIndicatorDots with a margin of 6 around them, as libadwaita's demo (adw-demo-page-carousel.ui) lays them out, where Breeze UI's PipsPagers and Carousel stand. With rtl, the carouselrtl set, the default direction is right to left and each stands as far from the window's right edge."""
    if rtl:
        Gtk.Widget.set_default_direction(Gtk.TextDirection.RTL)

    def carousel_of(pages, chosen, page):
        carousel = Adw.Carousel()
        for i in range(pages):
            carousel.append(page(i))
        carousel.scroll_to(carousel.get_nth_page(chosen), False)
        return carousel

    def unseen(pages, chosen):
        carousel = carousel_of(pages, chosen, lambda i: Gtk.Box())
        carousel.set_size_request(1, 1)
        carousel.set_opacity(0)
        return carousel

    five_pages, nine_pages = unseen(5, 1), unseen(9, 6)
    five = Adw.CarouselIndicatorDots(carousel=five_pages)
    nine = Adw.CarouselIndicatorDots(carousel=nine_pages)
    slides = carousel_of(4, 1, lambda i: Gtk.Label(label=f"Slide {i + 1}", hexpand=True, vexpand=True))
    slides.set_size_request(320, 160)
    dots = Adw.CarouselIndicatorDots(carousel=slides, margin_top=6, margin_bottom=6, margin_start=6, margin_end=6)
    box = vbox(0, slides, dots, width=320)

    def at(width):
        """The left of a widget width pixels wide 16 from the window's start: its left side, or its right side right to left (GtkFixed places its children from the left in either direction)."""
        return 352 - 16 - width if rtl else 16

    # The dots' natural widths: 15 pixels a page and 6 at either end (adw_carousel_indicator_dots_measure).
    items = [(five_pages, 0, 0), (nine_pages, 0, 0)]
    items += [(five, at(87), 16), (nine, at(147), 56), (box, at(320), 96)]
    return fixed(352, 312, items, named={"five": five, "nine": nine, "carousel": slides, "dots": dots, "box": box})


SETS = [
    ("carousel", "adwaita-carousel", carousel_set),
    ("carouselrtl", "adwaita-carousel-rtl", lambda: carousel_set(True)),
]
