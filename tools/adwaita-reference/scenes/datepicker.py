# datepicker: libadwaita's counterparts of TestScreenshotsDatePicker, TestAdwDatePickerPopover and TestAdwCalendarAsRendered (ui/breeze/widgets_datepicker_test.go, adw_native_test.go).
from common import Gdk, GLib, Graphene, Gtk, fixed, descendants


def datepicker_set(open_popover=False, arrow=True):
    """TestScreenshotsDatePicker's datePickerView: a GtkCalendar at 16, 16 showing November 2026 with the 18th chosen, and at 360, 16 a GtkMenuButton showing the same day as American English writes it short, whose popover holds a calendar, closed; in the sets datepickeropen and datepickernoarrow, those of TestAdwDatePickerPopover, the popover is open, with GtkPopover's arrow as the menu button gives it and without it (has-arrow false), as MyGo's popovers are drawn. Run it with LC_ALL=en_US.UTF-8, so that GTK names the months and days in English and starts the weeks on Sunday; GTK marks the machine's today, which TestScreenshotsDatePicker fixes as the 6th of November 2026."""
    calendar = Gtk.Calendar(year=2026, month=10, day=18)
    inner = Gtk.Calendar(year=2026, month=10, day=18)
    popover = Gtk.Popover(child=inner, has_arrow=arrow)
    button = Gtk.MenuButton(label="11/18/2026", popover=popover)
    named = {"calendar": calendar, "date-picker": button}
    if open_popover:
        named.update({"popover-calendar": inner})
    scene = fixed(560, 400, [(calendar, 16, 16), (button, 360, 16)], named=named, popovers=[popover] if open_popover else None)
    # The menu button opens its popover, so it shows checked while the popover is open.
    scene.open_popovers = button.popup
    return scene


def calendar_day(calendar, day):
    """The label of the calendar's day of the month shown."""
    return next(w for w in descendants(calendar.get_last_child()) if isinstance(w, Gtk.Label) and w.get_label() == str(day) and not w.has_css_class("other-month"))


def click_day(calendar, day):
    """Clicks the calendar's day: its GtkGestureClick's pressed signal at the label's middle, which GtkCalendar handles as a click."""
    label = calendar_day(calendar, day)
    ok, point = label.compute_point(calendar, Graphene.Point().init(label.get_width() / 2, label.get_height() / 2))
    for controller in calendar.observe_controllers():
        if isinstance(controller, Gtk.GestureClick):
            controller.emit("pressed", 1, point.x, point.y)


def print_day(name, calendar, day):
    """Prints whether the calendar has the focus and its day's state flags, where GTK_STATE_FLAG_FOCUSED (32) shows the focus ring, before the line of the set's PNG."""
    print(f"  {name}: the calendar has the focus {calendar.has_focus()}, the {day}th's flags {int(calendar_day(calendar, day).get_state_flags())}")


def datepicker_focus_set(keys=(), away=False):
    """The datepicker set after a click on the calendar's 20th, the click of TestAdwCalendarAsStyled and TestAdwCalendarAsRendered: GtkCalendar chooses the day, takes the focus and sets GTK_STATE_FLAG_FOCUSED on its label, which shows libadwaita's focus ring. The window does not show focus rings for the keyboard (focus-visible false), as after a click. In the set datepickerkeys, the keys given are then pressed on the calendar (its GtkEventControllerKey's key-pressed signal): Right puts the ring on the 21st, a day not chosen. In the set datepickeraway, the menu button then takes the focus: GtkCalendar sets the flag on the label whether it has the focus or not, so the ring stays."""
    scene = datepicker_set()
    calendar, button = scene.named["calendar"], scene.named["date-picker"]

    def click(window):
        click_day(calendar, 20)
        for key in keys:
            for controller in calendar.observe_controllers():
                if isinstance(controller, Gtk.EventControllerKey):
                    controller.emit("key-pressed", key, 0, 0)
        if away:
            button.grab_focus()
        window.set_focus_visible(False)
        print_day("after the click", calendar, 20)

    scene.shown = lambda window: GLib.timeout_add(300, lambda: (click(window), False)[1])
    return scene


def datepicker_reopen_set():
    """The datepickernoarrow set after a click on the popover's calendar's 20th, the popover then closed, as an app closes it once a day is chosen, and opened again: the calendar keeps the label it focused (GtkCalendar's focus_row and focus_col), so it shows the ring as the popover opens again, where the popover opened the first time showed none. Before the click it prints where Tab goes from the calendar, which has the focus as the popover opens: GtkCalendar is focusable itself and comes before its buttons."""
    scene = datepicker_set(True, False)
    button, inner = scene.named["date-picker"], scene.named["popover-calendar"]
    popover = button.get_popover()

    def reopen():
        popover.popup()
        GLib.timeout_add(200, lambda: (print_day("opened again", inner, 20), False)[1])

    def opened():
        print_day("opened", inner, 20)
        # The Tab order from the calendar, which has the focus: the popover's contents take Tab as the window would.
        window = popover.get_root()
        order = []
        for _ in range(4):
            popover.child_focus(Gtk.DirectionType.TAB_FORWARD)
            focus = window.get_focus()
            order.append(focus.get_tooltip_text() or type(focus).__name__)
        print("  Tab from the calendar: " + ", ".join(order))
        click_day(inner, 20)
        popover.popdown()
        GLib.timeout_add(200, lambda: (reopen(), False)[1])

    scene.open_popovers = lambda: (button.popup(), GLib.timeout_add(200, lambda: (opened(), False)[1]))
    return scene


def datepicker_rtl_set():
    """TestScreenshotsDatePicker's rtl-light shot: the datepicker set right to left, the default direction as a right-to-left language sets it (a widget's own direction would not reach the menu button's arrow), the calendar's right edge at 544 and the button's at 200, where the shot puts Breeze UI's."""
    Gtk.Widget.set_default_direction(Gtk.TextDirection.RTL)
    calendar = Gtk.Calendar(year=2026, month=10, day=18)
    button = Gtk.MenuButton(label="11/18/2026", popover=Gtk.Popover(child=Gtk.Calendar(year=2026, month=10, day=18)))
    width = calendar.measure(Gtk.Orientation.HORIZONTAL, -1)[1]
    return fixed(560, 400, [(calendar, 544 - width, 16), (button, 200 - button.measure(Gtk.Orientation.HORIZONTAL, -1)[1], 16)], named={"calendar": calendar, "date-picker": button})


SETS = [
    ("datepicker", "adwaita-datepicker", datepicker_set),
    ("datepickeropen", "adwaita-datepicker-open", lambda: datepicker_set(True)),
    ("datepickernoarrow", "adwaita-datepicker-noarrow", lambda: datepicker_set(True, False)),
    ("datepickerfocus", "adwaita-datepicker-focus", datepicker_focus_set),
    ("datepickerkeys", "adwaita-datepicker-keys", lambda: datepicker_focus_set([Gdk.KEY_Right])),
    ("datepickeraway", "adwaita-datepicker-away", lambda: datepicker_focus_set(away=True)),
    ("datepickerreopen", "adwaita-datepicker-reopen", datepicker_reopen_set),
    ("datepickerrtl", "adwaita-datepicker-rtl", datepicker_rtl_set),
]
