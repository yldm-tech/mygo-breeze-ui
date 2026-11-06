# carousel: libadwaita's counterparts of TestScreenshotsCarousel, TestAdwCarouselDragAsRendered and TestAdwCarouselDotsWidth (ui/breeze/widgets_carousel_test.go).
import subprocess

import gi

gi.require_version("GdkX11", "4.0")
from gi.repository import GdkX11  # noqa: E402,F401 (GdkX11 gives the window's surface get_xid)

from common import FRAME, Adw, GLib, Gtk, fixed, vbox  # noqa: E402


def carousel_set(rtl=False, drag=False, insensitive=False):
    """TestScreenshotsCarousel's carouselShotView: AdwCarouselIndicatorDots of 5 pages with the 2nd chosen at 16, 16 and of 9 pages with the 7th chosen at 16, 56, each showing a carousel of empty pages that the scene keeps out of sight, 1 pixel square and transparent (the dots read its pages and position once it is laid out, and are measured again then), and at 16, 96 a vertical GtkBox 320 wide of an AdwCarousel 160 high of four pages, each a label of its slide centered, the 2nd shown, over AdwCarouselIndicatorDots with a margin of 6 around them, as libadwaita's demo (adw-demo-page-carousel.ui) lays them out, where Breeze UI's PipsPagers and Carousel stand. With rtl, the carouselrtl set, the default direction is right to left and each stands as far from the window's right edge. With drag, the carouseldrag and carouseldragrtl sets, xdotool presses the mouse in the carousel's middle and moves it 10 times by 10 pixels toward its start, the way to the next page, 100 ms apart, and the window is rendered with the mouse still pressed, as TestAdwCarouselDragAsRendered drags Breeze UI's; the carousel's position is printed. With insensitive, the carouselinsensitive set, the pagers and the carousel's box are insensitive, to show how libadwaita draws dots that are."""
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
    if insensitive:
        for widget in (five, nine, box):
            widget.set_sensitive(False)

    def at(width):
        """The left of a widget width pixels wide 16 from the window's start: its left side, or its right side right to left (GtkFixed places its children from the left in either direction)."""
        return 352 - 16 - width if rtl else 16

    # The dots' natural widths: 15 pixels a page and 6 at either end (adw_carousel_indicator_dots_measure).
    items = [(five_pages, 0, 0), (nine_pages, 0, 0)]
    items += [(five, at(87), 16), (nine, at(147), 56), (box, at(320), 96)]
    scene = fixed(352, 312, items, named={"five": five, "nine": nine, "carousel": slides, "dots": dots, "box": box})

    def shown(window):
        """Measures the dots again once the carousels are laid out: the dots measure their carousel's snap points, which AdwCarousel sets only as it allocates its pages, and they measure again only when its number of pages changes, so dots measured first stay 27 pixels wide (as for one page) and draw their dots past their sides."""
        GLib.timeout_add(300, lambda: ([d.queue_resize() for d in (five, nine, dots)], False)[1])
        if drag:
            GLib.timeout_add(400, lambda: (press_and_drag(window.get_surface().get_xid()), False)[1])

    def press_and_drag(xid):
        """Presses the mouse in the carousel's middle, in the window's coordinates (the content's and the frame), and moves it by 10 pixels 10 times, 100 ms apart, from another process so that GTK takes each move as it comes. Right to left, GTK maps the window off the screen's left edge, at -361 without a window manager, where the press would land on its resize edge: xdotool moves it to the screen's corner first."""
        x, y = at(320) + 160 + FRAME, 96 + 80 + FRAME
        step = 10 if rtl else -10
        args = ["xdotool", "windowmove", "--sync", str(xid), "0", "0", "sleep", "0.2", "mousemove", "--window", str(xid), str(x), str(y), "sleep", "0.2", "mousedown", "1"]
        for k in range(1, 11):
            args += ["sleep", "0.1", "mousemove", "--window", str(xid), str(x + step * k), str(y)]
        subprocess.Popen(args)

    def let_go():
        print(f"  position: {slides.get_position():g}")
        subprocess.run(["xdotool", "mouseup", "1"], check=True)

    scene.shown = shown
    if drag:
        scene.wait, scene.done = 2500, let_go
    return scene


def carousel_many_set():
    """TestAdwCarouselDotsWidth's view: at 16, 16 a vertical GtkBox with no width of its own, of an AdwCarousel of 25 pages, each a label of its number centered across the carousel (hexpand, else AdwCarousel's size_allocate makes its pages as wide as the widest label), the first shown, over AdwCarouselIndicatorDots with a margin of 6 around them, as in the carousel set. The dots ask for their natural width as their minimum (adw_carousel_indicator_dots_measure), 15 pixels a page and 6 at either end, so the box and the carousel are as wide as the dots and their margins, 399 pixels, rather than as the pages; they are measured again once the carousel is laid out, as in the carousel set."""
    slides = Adw.Carousel()
    for i in range(25):
        slides.append(Gtk.Label(label=f"{i + 1}", hexpand=True))
    dots = Adw.CarouselIndicatorDots(carousel=slides, margin_top=6, margin_bottom=6, margin_start=6, margin_end=6)
    box = vbox(0, slides, dots)
    scene = fixed(432, 96, [(box, 16, 16)], named={"carousel": slides, "dots": dots, "box": box})
    scene.shown = lambda window: GLib.timeout_add(300, lambda: (dots.queue_resize(), False)[1])
    return scene


SETS = [
    ("carousel", "adwaita-carousel", carousel_set),
    ("carouselrtl", "adwaita-carousel-rtl", lambda: carousel_set(True)),
    ("carouseldrag", "adwaita-carousel-drag", lambda: carousel_set(drag=True)),
    ("carouseldragrtl", "adwaita-carousel-drag-rtl", lambda: carousel_set(True, drag=True)),
    ("carouselinsensitive", "adwaita-carousel-insensitive", lambda: carousel_set(insensitive=True)),
    ("carouselmany", "adwaita-carousel-many", carousel_many_set),
]
