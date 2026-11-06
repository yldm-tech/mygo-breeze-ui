package breeze

import (
	"embed"
	"io/fs"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// ---- date pickers: a month's calendar drawn per platform, and a button opening it ----
//
// macOS: NSDatePicker of style clockAndCalendar with the year, month and day, which SwiftUI's DatePicker of the graphical style draws on macOS (the two render alike, 139 by 148 points), and NSDatePicker of style textField with presentsCalendarOverlay for the closed picker, the field of SwiftUI's compact style without its stepper; both rendered at 2x in Aqua and Dark Aqua in a key window on macOS 27.2 by tools/appkit-reference/reference.swift (renderDatePicker), which prints the frames, the runs of pixels and the accessibility of the calendar's three buttons.
//
// Windows 11: WinUI's CalendarView and CalendarDatePicker of microsoft/microsoft-ui-xaml main at 5ed8dd3: CalendarView_themeresources.xaml (the brushes, CalendarViewDayItem's 40-DIP MinWidth and MinHeight and its Margin of 1, the header button's and the navigation buttons' margins, paddings and glyphs, the week day names' style, the 1-DIP line under the header, the month view's margins), CalendarViewBaseItemChrome.cpp (the round chrome of CalendarViewBaseItemRoundedChromeEnabled: today a filled circle, the chosen day a 1-DIP ring, both an inner ring), CalendarView_Partial_Interaction.cpp (the keys), CalendarDatePicker_themeresources.xaml and calendardatepicker_partial.cpp (the button, its glyph column, the flyout below it, the short date, the calendar closing as a day is chosen), and Microsoft.UI.Xaml.Common.rc for the strings.
//
// GNOME: GTK 4.24.1's gtkcalendar.c (the header's arrows and labels, the homogeneous grid of day names and six weeks, a week of the month before when the month starts on the first column, the keys) as libadwaita 1.10.0's _calendar.scss styles it, with _views.scss, _buttons.scss (%button_basic_flat, .image-button), _popovers.scss (a calendar in a popover loses its border and its background) and adwaita-icon-theme 51's pan-start-symbolic and pan-end-symbolic.

//go:embed i18n/*.json
var stringFiles embed.FS

// breezeStrings is Breeze UI's own catalog of the strings its widgets say, a library catalog of namespace "breeze" (i18n.NewLibraryCatalog) in English: an app translates or changes them with keys "breeze.<key>" in its own catalogs, as it changes MyGo's with "mygo.<key>", and they are said in the language the window shows.
var breezeStrings = sync.OnceValue(func() *i18n.Catalog {
	sub, err := fs.Sub(stringFiles, "i18n")
	if err != nil {
		panic(err)
	}
	cat, err := i18n.NewLibraryCatalog(sub, "en", "breeze")
	if err != nil {
		panic(err)
	}
	return cat
})

// say returns the string of key in breezeStrings in the window's language.
func say(c *ui.Context, key string) string { return breezeStrings().Format(c.Locale(), key, nil) }

// dateFormats returns the user's regional formats and the locale MyGo's own date controls write dates in, as its Context.formats puts them together: the window's language with the region of the formats (Region.FormatLocale), or where the system says nothing of them, as for a Tester, those CLDR gives the language.
func dateFormats(c *ui.Context) (i18n.Region, string) {
	lang := c.Locale()
	if lang == "" {
		lang = "en"
	}
	r := c.Region()
	if r.HourCycle != 12 && r.HourCycle != 24 {
		r = i18n.RegionFor(lang)
		r.Locale = ""
	}
	key := formatsKey{lang, r}
	if v, ok := formatLocales.Load(key); ok {
		return r, v.(string)
	}
	locale := r.FormatLocale(lang)
	formatLocales.Store(key, locale)
	return r, locale
}

type formatsKey struct {
	lang   string
	region i18n.Region
}

// formatLocales caches dateFormats's locales: the views ask for them every frame.
var formatLocales sync.Map // formatsKey -> string

// dateTexts caches the dates and numbers the date pickers write, by locale, pattern and day, as MyGo's controls keep theirs: a calendar writes some fifty every frame.
var dateTexts struct {
	sync.Mutex
	m map[dateTextKey]string
}

type dateTextKey struct {
	locale, skeleton string
	v                int64
}

// dateText writes the day of d in locale with a skeleton of i18n.DateOptions ("yMd", "yMMMMd"), or its day of the month alone for "d#", in the locale's digits.
func dateText(locale, skeleton string, d time.Time) string {
	y, m, day := d.Date()
	key := dateTextKey{locale, skeleton, time.Date(y, m, day, 0, 0, 0, 0, time.UTC).Unix() / 86400}
	dateTexts.Lock()
	defer dateTexts.Unlock()
	if s, ok := dateTexts.m[key]; ok {
		return s
	}
	var s string
	if skeleton == "d#" {
		s = i18n.FormatNumber(locale, day, i18n.NumberOptions{NoGrouping: true})
	} else {
		s = i18n.FormatDate(locale, time.Date(y, m, day, 12, 0, 0, 0, time.UTC), i18n.DateOptions{Skeleton: skeleton})
	}
	if dateTexts.m == nil || len(dateTexts.m) >= 2048 {
		dateTexts.m = map[dateTextKey]string{}
	}
	dateTexts.m[key] = s
	return s
}

// calendarNow returns the time whose day the calendars mark as today: the frame's time. Tests set it to a fixed day.
var calendarNow = func(c *ui.Context) time.Time { return c.Now() }

// dayOf returns the day of t at midnight in UTC, the calendars' days.
func dayOf(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// sameDay reports whether a and b fall on the same day.
func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

// addMonths returns d moved by n months, on the same day or the last of a shorter month, as AppKit's, WinUI's (Calendar.AddMonths) and GTK's (g_date_time_add_months) calendars move by months.
func addMonths(d time.Time, n int) time.Time {
	y, m, day := d.Date()
	first := time.Date(y, m+time.Month(n), 1, 0, 0, 0, 0, time.UTC)
	last := first.AddDate(0, 1, -1).Day()
	return time.Date(first.Year(), first.Month(), min(day, last), 0, 0, 0, 0, time.UTC)
}

// setDay sets the day of *date to d's, keeping its time of day and its location (the local time of a zero *date), and reports whether it changed, as MyGo's calendar sets it.
func setDay(date *time.Time, d time.Time) bool {
	loc := date.Location()
	if date.IsZero() {
		loc = time.Local
	}
	y, m, day := d.Date()
	h, mi, s := date.Clock()
	next := time.Date(y, m, day, h, mi, s, date.Nanosecond(), loc)
	if next.Equal(*date) {
		return false
	}
	*date = next
	return true
}

// calendarState is what a calendar keeps between frames: the day the keys are on, whose month it shows, the date it last saw, which the cursor follows when the app sets another, and whether a click on a day or a key that moves among the days has put the keys on a day since the calendar was built (GtkCalendar's focus_row and focus_col, -1 until then).
type calendarState struct {
	cursor, seen time.Time
	started      bool
	placed       bool
}

// calendarBuild is what a calendar's parts are built from in a frame.
type calendarBuild struct {
	c       *ui.Context
	t       *ui.Theme
	locale  string
	first   time.Weekday
	month   time.Time // the first of the month shown
	cursor  time.Time // the day the keys are on
	chosen  time.Time // *date's day; zero when *date is
	today   time.Time
	rtl     bool
	bare    bool       // inside a date picker's popover
	keys    ui.Element // the element that takes the keys and the focus
	focused bool       // keys has the focus, from the keyboard or not
	placed  bool       // a click on a day or a key among the days has put the keys on one (calendarState)
	// move shows d, and chooses it with choose; pick chooses a day clicked.
	move func(d time.Time, choose bool)
	pick func(d time.Time)
}

// start returns the first day of the six weeks shown: the first day of the week on or before the first of the month, or on GNOME, where it is the first itself, a week before it (gtkcalendar.c's calendar_compute_days: first_day 0 becomes 7).
func (b *calendarBuild) start(p Platform) time.Time {
	back := (int(b.month.Weekday()) - int(b.first) + 7) % 7
	if p == Linux && back == 0 {
		back = 7
	}
	return b.month.AddDate(0, 0, -back)
}

// CalendarView creates a month's calendar choosing *date, drawn as the platform draws its calendar, in the window's language with its weeks starting on the first day of the user's regional formats (Context.Region), as MyGo's own controls: the names of the months and the days, and the numbers, come from MyGo's i18n package. A zero *date shows today's month with no day chosen.
//
// A click on a day chooses it. The days are one Tab stop, after the header's buttons on Windows and GNOME; the arrows move the day the keys are on, Up and Down by weeks (Left and Right mirrored right to left), Page Up and Page Down by months, keeping the day or taking the last of a shorter month, Home and End to the first and the last day of the month, and Enter and Space choose the day the keys are on; the month shown is that day's. Assistive technology sees a table named by the month, of buttons named by their dates (Element.ActiveDescendant tells it which day the keys are on), checked where chosen.
//   - macOS: NSDatePicker of style clockAndCalendar, as SwiftUI's graphical DatePicker draws it: 139 by 148 points, a field (white, white at 4.9% over the window in Dark Aqua) inside a half-point line (black at 7.8%, white at 9.8%) whose bottom edge is darker; the month and the year abbreviated ("Nov 2026") in 13-point bold labelColor 3 points in; at the end the buttons ◀ (the month before), ● (today) and ▶ (the month after), 6, 7 and 6 points wide and 7 high, 8 apart, in secondaryLabelColor; the short day names ("Su") in 10-point bold secondaryLabelColor over a 1-point separator (black at 10%, white at 9.8%); six weeks in rows 18 points high, columns 18.5 wide from 4 points in, the numbers in 10 points with proportional digits at the end of their column, 2 before it; the days of the months around in tertiaryLabelColor, today's number in the accent (#0060e5, #26a0ff for the blue), the chosen day an 18 by 16 rectangle with corners of 2.5, selectedContentBackgroundColor with a white number while the calendar has the focus, unemphasizedSelectedContentBackgroundColor without it (measured on macOS 27.2). The whole calendar is the one Tab stop with AppKit's focus ring around it, its buttons not, as NSDatePicker is one control; every move chooses the day, as NSDatePicker's arrows and buttons do.
//   - Windows: CalendarView: 300 by 351 DIPs, ControlFillColorInputActive inside a ControlStrokeColorDefault border with corners of 4; the month and the year ("November 2026") in 14-DIP semibold 16 DIPs in, the Previous and Next buttons (32 by 35, the CaretUpSolid8 and CaretDownSolid8 glyphs, SubtleFillColorSecondary under the pointer and Tertiary pressed, the glyph 0.875 as large pressed) at the end; a 1-DIP line; the short day names in 12-DIP semibold, centered in 42-DIP columns; six weeks of 40-DIP day items with a margin of 1, their numbers in 14 DIPs, the days out of the month in TextFillColorSecondary; today a filled AccentFillColorDefault circle with TextOnAccentFillColorPrimary text, the chosen day a 1-DIP AccentFillColorDefault ring with AccentTextFillColorPrimary text (a TextOnAccentFillColorPrimary ring inside today's circle when today is chosen), a day under the pointer SubtleFillColorSecondary, pressed Tertiary; the system focus visual 3 DIPs around the day the keys are on. The month header button that zooms out to the year is the header's text here, and the Previous and Next buttons are Tab stops before the days; the arrows only move among the days, Enter or Space or a click choosing one.
//   - GNOME: GtkCalendar with libadwaita's style: view_bg_color inside a 1-pixel $border_color border; a header, its bottom border a line, of flat 34-pixel square image buttons with pan-start-symbolic and pan-end-symbolic around the month (the widest month's room) and, after free room, around the year (10 pixels wider than it); a grid padded 3 at the left and the bottom of rows 27 high (a label's 18 and its padding of 3 and margin of 3 above it) sharing the width in seven columns, a 3-pixel margin at each label's right; the abbreviated day names ("Sun"), the days in view_fg_color, those of the months around at 30%, today with a 2-pixel $border_color line along its label's bottom, the chosen day the accent with corners of 9 and accent_fg_color text; the focus ring inside the label the keys are on. The four buttons are Tab stops before the days; they and Ctrl with the arrows (Left and Right the month, Up and Down the year) choose the same day in another month or year, as GTK's do; the arrows only move among the days, Space or Enter or a click choosing one.
//
// It returns the calendar, whose Changed reports a day chosen. Calendar keeps MyGo's own, which takes Browse.
func CalendarView(c *ui.Context, date *time.Time) ui.Element {
	return calendarView(c, date, nil)
}

// calendarView builds CalendarView; with done, in a date picker's popover, it is drawn without its own frame where the platform draws it so, takes the focus as it appears, and calls done once a click, Enter or Space chooses a day.
func calendarView(c *ui.Context, date *time.Time, done func()) ui.Element {
	p := PlatformOf(c)
	region, locale := dateFormats(c)
	root := ui.Column(c).Shrink(0)
	st := ui.Local(root, "breeze-calendar", func() calendarState { return calendarState{} })
	now := dayOf(calendarNow(c))
	if !st.started || !st.seen.Equal(*date) {
		st.cursor, st.seen, st.started = now, *date, true
		if !date.IsZero() {
			st.cursor = dayOf(*date)
		}
	}
	b := &calendarBuild{c: c, t: c.Theme(), locale: locale, first: region.FirstWeekday, cursor: st.cursor, today: now, rtl: c.Direction() == ui.RightToLeft, bare: done != nil, placed: st.placed}
	b.month = time.Date(b.cursor.Year(), b.cursor.Month(), 1, 0, 0, 0, 0, time.UTC)
	if !date.IsZero() {
		b.chosen = dayOf(*date)
	}
	// What the input of this frame asks for, applied once the calendar is built, as it was built from the state the user saw.
	var moveTo, chooseDay *time.Time
	b.move = func(d time.Time, choose bool) {
		moveTo = &d
		if choose {
			chooseDay = &d
		}
	}
	finish, placed := false, false
	b.pick = func(d time.Time) {
		moveTo, chooseDay, finish, placed = &d, &d, true, true
	}
	// step moves the keys among the days, as the arrows, Page Up, Page Down, Home and End do, the header's buttons and Ctrl with the arrows not.
	step := func(d time.Time, choose bool) {
		b.move(d, choose)
		placed = true
	}
	switch p {
	case Windows:
		fluentCalendar(b, root)
	case Linux:
		adwCalendar(b, root)
	default:
		macCalendar(b, root)
	}
	keys := b.keys
	if done != nil {
		keys.AutoFocus()
	}
	// The keys: GtkCalendar moves the month with Ctrl+Left and Ctrl+Right and the year with Ctrl+Up and Ctrl+Down (Cmd on macOS), choosing the same day there; the others are MyGo's calendar's, WinUI's and NSDatePicker's.
	back, forth := ui.KeyLeft, ui.KeyRight
	if b.rtl {
		back, forth = forth, back
	}
	moves := p == Mac
	cur := b.cursor
	switch {
	case p == Linux && keys.Shortcut(ui.Ctrl, ui.KeyLeft):
		b.move(addMonths(cur, -1), true)
	case p == Linux && keys.Shortcut(ui.Ctrl, ui.KeyRight):
		b.move(addMonths(cur, 1), true)
	case p == Linux && keys.Shortcut(ui.Ctrl, ui.KeyUp):
		b.move(addMonths(cur, -12), true)
	case p == Linux && keys.Shortcut(ui.Ctrl, ui.KeyDown):
		b.move(addMonths(cur, 12), true)
	case keys.Shortcut(0, back):
		step(cur.AddDate(0, 0, -1), moves)
	case keys.Shortcut(0, forth):
		step(cur.AddDate(0, 0, 1), moves)
	case keys.Shortcut(0, ui.KeyUp):
		step(cur.AddDate(0, 0, -7), moves)
	case keys.Shortcut(0, ui.KeyDown):
		step(cur.AddDate(0, 0, 7), moves)
	case keys.Shortcut(0, ui.KeyPageUp):
		step(addMonths(cur, -1), moves)
	case keys.Shortcut(0, ui.KeyPageDown):
		step(addMonths(cur, 1), moves)
	case keys.Shortcut(0, ui.KeyHome):
		step(b.month, moves)
	case keys.Shortcut(0, ui.KeyEnd):
		step(b.month.AddDate(0, 1, -1), moves)
	case keys.Shortcut(0, ui.KeyEnter), keys.Shortcut(0, ui.KeySpace):
		b.pick(cur)
	}
	if moveTo != nil {
		st.cursor = *moveTo
		st.placed = st.placed || placed
		if chooseDay != nil && setDay(date, *chooseDay) {
			root.SetChanged()
		}
		st.seen = *date
		c.Invalidate()
	}
	if finish && done != nil {
		done()
	}
	return root
}

// calendarDays builds the days' grid into the element building it: the row of the weekdays' names, which weekday builds, and six weeks, rows that row builds (week -1 is the names' row) holding seven days that cell builds, returning the element a click chooses. For assistive technology keys is a table named by the month, its days buttons named by their dates, checked where chosen, the one the keys are on its active descendant.
func calendarDays(b *calendarBuild, p Platform, weekday func(name string), row func(week int, fn func()) ui.Element, cell func(d time.Time) ui.Element) {
	b.keys.Role(ui.RoleTable).Label(dateText(b.locale, "yMMMM", b.month))
	names := weekdayNames(b.locale, p)
	row(-1, func() {
		for i := range 7 {
			weekday(names[(int(b.first)+i)%7])
		}
	})
	start := b.start(p)
	for w := range 6 {
		row(w, func() {
			for d := range 7 {
				day := start.AddDate(0, 0, w*7+d)
				e := cell(day)
				e.Role(ui.RoleButton).Label(dateText(b.locale, "yMMMMd", day)).Checked(sameDay(day, b.chosen))
				if e.Clicked() {
					b.keys.Focus()
					b.pick(day)
				}
				if sameDay(day, b.cursor) {
					b.keys.ActiveDescendant(e)
				}
			}
		}).Role(ui.RoleRow)
	}
}

// weekdayNames returns the names of the days of the week in locale as the platform's calendar shows them, Sunday first: CLDR's short names ("Su") for NSDatePicker and for WinUI's DayOfWeekFormat {dayofweek.abbreviated(2)}, the abbreviated ones ("Sun") for GtkCalendar's %a.
func weekdayNames(locale string, p Platform) [7]string {
	if p == Linux {
		return i18n.Weekdays(locale, i18n.Abbreviated, true)
	}
	return i18n.Weekdays(locale, i18n.Short, true)
}

// ---- macOS: NSDatePicker of style clockAndCalendar ----

// macCalendarField returns the fill of NSDatePicker's calendar: white, and in Dark Aqua white at 4.9% over the window (41 over 30, measured on macOS 27.2); in a popover, the panel's own (macMenu's #f6f6f6 and #2b2b2b).
func macCalendarField(t *ui.Theme, bare bool) ui.Color {
	if bare {
		return pick(t.Dark, ui.Hex("#f6f6f6"), ui.Hex("#2b2b2b"))
	}
	return pick(t.Dark, white, t.Background.Mix(white, 0.049))
}

// macCalendarToday returns the color of today's number: #0060e5 in Aqua and #26a0ff in Dark Aqua for the blue accent (measured on macOS 27.2), the accent for another (inferred).
func macCalendarToday(t *ui.Theme) ui.Color {
	if t.Accent == macAccent {
		return pick(t.Dark, ui.Hex("#0060e5"), ui.Hex("#26a0ff"))
	}
	return t.Accent
}

// macCalendar builds NSDatePicker's calendar into root, CalendarView's macOS look (measured at 2x on macOS 27.2 by tools/appkit-reference's datepicker scene): 139 by 148 points; the box from half a point in, 137.5 by 146.5 (in Dark Aqua from a point in, 136.5 by 145), its edges half-point lines over the window (black at 7.8%, 235 over white; white at 9.8% in Dark Aqua, a lighter 13.3% along its top, 52 and 60 over the window's 30) around the field, and under its bottom two half-point lines of black at 24% and 12% (white at 12% and 26%); the header 20 points high, the title 3 points in, its buttons centered half a point below its middle, the last 6.5 points from the end; the day names' row 17 points high and the separator, a point, under it at 37, from 3 to 135 points in; the weeks from 38, rows 18 high, columns 18.5 wide from 4 points in, the day names and the numbers 2 points before their column's end, the names half a point below their row's middle; the chosen day's rectangle inset half a point at its column's start and a point at its row's top and bottom. The text colors are AppKit's label colors over the field: labelColor 84.7%, secondaryLabelColor 49.8% (54.9% in Dark Aqua), tertiaryLabelColor 26% (25%). In a popover the field and its lines are left out, the calendar lying on the panel (inferred).
func macCalendar(b *calendarBuild, root ui.Element) {
	c, t := b.c, b.t
	field := macCalendarField(t, b.bare)
	label := tint(field, t.Dark, 0.847)
	secondary := tint(field, t.Dark, pick(t.Dark, float32(0.498), 0.549))
	tertiary := tint(field, t.Dark, pick(t.Dark, float32(0.26), 0.25))
	sep := tint(field, t.Dark, pick(t.Dark, float32(0.10), 0.098))
	b.keys = root
	b.focused = root.Focused()
	ring := root.FocusVisible()
	root.Width(139).Height(148).Focusable().FocusRing(false)
	if !b.bare {
		edge := pick(t.Dark, black.Alpha(0.078), white.Alpha(0.098))
		top := pick(t.Dark, black.Alpha(0.078), white.Alpha(0.133))
		low1 := pick(t.Dark, black.Alpha(0.239), white.Alpha(0.124))
		low2 := pick(t.Dark, black.Alpha(0.122), white.Alpha(0.258))
		rtl := b.rtl
		// Dark Aqua draws the box half a point further in at the top and the sides and a point higher at the bottom.
		in, w, h := float32(0.5), float32(137.5), float32(146.5)
		if t.Dark {
			in, w, h = 1, 136.5, 145
		}
		root.Draw(func(p *ui.Painter, r ui.Rect) {
			box := ui.Rect{X: r.X + in, Y: r.Y + in, W: w, H: h}
			if rtl {
				box.X = r.X + r.W - in - box.W
			}
			// The lines lie over the window, the field inside them.
			p.Fill(ui.Rect{X: box.X + 0.5, Y: box.Y + 0.5, W: box.W - 1, H: box.H - 0.5}, field, 0)
			p.Fill(ui.Rect{X: box.X, Y: box.Y, W: 0.5, H: box.H}, edge, 0)
			p.Fill(ui.Rect{X: box.X + box.W - 0.5, Y: box.Y, W: 0.5, H: box.H}, edge, 0)
			p.Fill(ui.Rect{X: box.X + 0.5, Y: box.Y, W: box.W - 1, H: 0.5}, top, 0)
			p.Fill(ui.Rect{X: box.X, Y: box.Y + box.H, W: box.W, H: 0.5}, low1, 0)
			p.Fill(ui.Rect{X: box.X, Y: box.Y + box.H + 0.5, W: box.W, H: 0.5}, low2, 0)
		})
	}
	root.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if ring {
			macFocusRing(p, r, 0)
		}
	})
	sel := macListSelection(t, b.focused)
	root.Children(func() {
		ui.Row(c).Height(20).Shrink(0).PaddingStart(3).PaddingEnd(6.5).AlignItems(ui.Center).Children(func() {
			ui.Text(c, dateText(b.locale, "yMMM", b.month)).SingleLine().Grow(1).Shrink(1).MinWidth(0).FontSize(t.FontSize).FontWeight(700).TextColor(label)
			for i, key := range []string{"calendarPreviousMonth", "calendarToday", "calendarNextMonth"} {
				w := float32(6)
				if i == 1 {
					w = 7
				}
				btn := ui.Box(c).Size(w, 7).Shrink(0).Margin(1, 0, 0, 0).Role(ui.RoleButton).Label(say(c, key))
				if i > 0 {
					btn.MarginStart(8)
				}
				ink := secondary
				if btn.Pressed() {
					ink = label // inferred
				}
				kind, rtl := i, b.rtl
				btn.Draw(func(p *ui.Painter, r ui.Rect) { macCalendarButton(p, r, kind, rtl, ink) })
				if btn.Clicked() {
					switch i {
					case 0:
						b.move(addMonths(b.cursor, -1), true)
					case 1:
						b.move(b.today, true)
					default:
						b.move(addMonths(b.cursor, 1), true)
					}
				}
			}
		})
		calendarDays(b, Mac, func(name string) {
			ui.Row(c).Width(18.5).Height(17).Shrink(0).Justify(ui.End).AlignItems(ui.Center).Padding(1, 0, 0, 0).PaddingEnd(2).Role(ui.RoleColumnHeader).Children(func() {
				ui.Text(c, name).SingleLine().FontSize(t.FontSize * 10 / 13).FontWeight(700).TextColor(secondary)
			})
		}, func(week int, fn func()) ui.Element {
			r := ui.Row(c).Height(18).Shrink(0).PaddingStart(4).AlignItems(ui.Start).Children(fn)
			if week < 0 {
				// The separator, a point under the names' 17.
				rtl := b.rtl
				r.Draw(func(p *ui.Painter, rr ui.Rect) {
					x := rr.X + 3
					if rtl {
						x = rr.X + rr.W - 135
					}
					p.Fill(ui.Rect{X: x, Y: rr.Y + 17, W: 132, H: 1}, sep, 0)
				})
			}
			return r
		}, func(d time.Time) ui.Element {
			chosen := sameDay(d, b.chosen)
			fg := label
			switch {
			case chosen && b.focused:
				fg = white
			case sameDay(d, b.today):
				fg = macCalendarToday(t)
			case d.Month() != b.month.Month():
				fg = tertiary
			}
			e := ui.Row(c).Width(18.5).Height(18).Shrink(0).Justify(ui.End).AlignItems(ui.Center).PaddingEnd(2)
			if chosen {
				rtl := b.rtl
				e.Draw(func(p *ui.Painter, r ui.Rect) {
					x := r.X + 0.5
					if rtl {
						x = r.X
					}
					p.Fill(ui.Rect{X: x, Y: r.Y + 1, W: 18, H: 16}, sel, 2.5)
				})
			}
			e.Children(func() {
				ui.Text(c, dateText(b.locale, "d#", d)).SingleLine().FontSize(t.FontSize * 10 / 13).TextColor(fg)
			})
			return e
		})
	})
}

// macCalendarButton draws NSDatePicker's header button kind into its frame r: 0 the month before, a triangle pointing back, 5.5 points wide and 7 high; 1 today, a 7-point dot; 2 the month after, pointing forth (their accessibility frames and ink, measured on macOS 27.2), mirrored right to left (inferred).
func macCalendarButton(p *ui.Painter, r ui.Rect, kind int, rtl bool, ink ui.Color) {
	if kind == 1 {
		p.Fill(r, ink, r.W/2)
		return
	}
	back := kind == 0
	if rtl {
		back = !back
	}
	var tri ui.Path
	if back {
		tri.MoveTo(r.X+r.W, r.Y).LineTo(r.X+0.5, r.Y+r.H/2).LineTo(r.X+r.W, r.Y+r.H).Close()
	} else {
		tri.MoveTo(r.X, r.Y).LineTo(r.X+r.W-0.5, r.Y+r.H/2).LineTo(r.X, r.Y+r.H).Close()
	}
	p.FillPath(&tri, ink)
}

// ---- Windows: WinUI's CalendarView ----

// fluentCalendarBg returns CalendarViewBackground, ControlFillColorInputActive: white, and #1e1e1e at 70% in the dark theme, over the window.
func fluentCalendarBg(t *ui.Theme) ui.Color {
	return pick(t.Dark, white, t.Background.Mix(ui.Hex("#1e1e1e"), 0.702))
}

// fluentCalendar builds a CalendarView into root, CalendarView's Windows look (CalendarView_themeresources.xaml): 300 DIPs wide, BorderThickness 1 of CalendarViewBorderBrush (ControlStrokeColorDefault), corners of ControlCornerRadius (4), CalendarViewBackground; the header 48 high, its HeaderButton's text (FontSize 14, SemiBold) placed by its Margin 7,6,3,7 and Padding 8,7,8,8 on a 20-DIP line, the PreviousButton and NextButton (NavigationButtonStyle) with Margins 3,6,3,7 and 3,6,7,7, 32 wide (Padding 12 each side of the 8-DIP glyph) and stretched to 35 high; a 1-DIP Border of the border brush; WeekDayNames with a Margin of 2,2,2,0 holding seven star columns of CaptionTextBlockStyle (12) SemiBold names with a Margin of 1 and a Padding of 12, 42 high, in CalendarItemForeground; and the MonthViewScrollViewer with a Margin of 2 holding six weeks of CalendarViewDayItems, 40 by 40 with a Margin of 1. The text colors are WinUI's translucent ones over the calendar's background.
func fluentCalendar(b *calendarBuild, root ui.Element) {
	c, t := b.c, b.t
	bg := fluentCalendarBg(t)
	th := fluentOver(t, bg)
	// In CalendarDatePicker's flyout the template gives the CalendarView corners of OverlayCornerRadius (8).
	root.Width(300).Border(1, t.Border).Radius(pick(b.bare, float32(4), 8)).Background(bg)
	root.Children(func() {
		ui.Row(c).Height(48).Shrink(0).AlignItems(ui.Start).Children(func() {
			ui.Box(c).Grow(1).Shrink(1).MinWidth(0).Height(20).Margin(13, 0, 0, 0).MarginStart(15).MarginEnd(11).Justify(ui.Center).Children(func() {
				ui.Text(c, dateText(b.locale, "yMMMM", b.month)).SingleLine().FontSize(t.FontSize).FontWeight(600).TextColor(th.Text)
			})
			for i, key := range []string{"calendarPrevious", "calendarNext"} {
				btn := ui.ButtonBase(c).Size(32, 35).Shrink(0).Margin(6, 0, 0, 0).MarginStart(3).MarginEnd(3).Radius(4).Label(say(c, key))
				if i == 1 {
					btn.MarginEnd(7)
				}
				fluentCalendarNav(t, th, bg, btn, i == 0)
				if btn.Clicked() {
					b.move(addMonths(b.cursor, 2*i-1), false)
				}
			}
		})
		ui.Box(c).Height(1).Shrink(0).Background(t.Border)
		grid := ui.Column(c).Padding(2).Shrink(0).Focusable().FocusRing(false)
		b.keys = grid
		b.focused = grid.Focused()
		focusRing := grid.FocusVisible()
		grid.Children(func() {
			calendarDays(b, Windows, func(name string) {
				ui.Box(c).Size(42, 42).Shrink(0).Center().Role(ui.RoleColumnHeader).Children(func() {
					ui.Text(c, name).SingleLine().FontSize(t.FontSize * 12 / 14).FontWeight(600).TextColor(th.Text)
				})
			}, func(week int, fn func()) ui.Element {
				r := ui.Row(c).Shrink(0).Children(fn)
				if week < 0 {
					r.Margin(0, 0, 2, 0)
				}
				return r
			}, func(d time.Time) ui.Element {
				e := ui.Box(c).Size(40, 40).Margin(1).Shrink(0).Center()
				fluentCalendarDay(c, b, th, bg, e, d, focusRing && sameDay(d, b.cursor))
				return e
			})
		})
	})
}

// fluentCalendarNav draws CalendarView's PreviousButton (back) or NextButton btn over the calendar's background bg: transparent, CalendarViewNavigationButtonBackgroundPointerOver (SubtleFillColorSecondary) under the pointer and Pressed (SubtleFillColorTertiary) pressed, corners of ControlCornerRadius; the CaretUpSolid8 (EDDB) or CaretDownSolid8 (EDDC) glyph at 8 DIPs in the CalendarView's Foreground (TextFillColorPrimary), ControlStrongFillColorDefault under the pointer and pressed, scaled by CalendarViewNavigationButtonScalePressed (0.875) pressed, drawn as a filled caret 8 DIPs wide and 4 high (inferred from the glyphs' names); the system focus visual 2 DIPs around it (CalendarViewNavigationButtonFocusVisualMargin -2).
func fluentCalendarNav(t, th *ui.Theme, bg ui.Color, btn ui.Element, back bool) {
	hover, press := fluentSubtle(th)
	ink := th.Text
	strong := bg.Mix(pick(t.Dark, black, white), pick(t.Dark, float32(0x72)/255, 0x8b/255.0))
	scale := float32(1)
	switch {
	case btn.Pressed():
		btn.Background(press)
		ink, scale = strong, 0.875
	case btn.Hovered():
		btn.Background(hover)
		ink = strong
	}
	focused := btn.FocusVisible()
	btn.FocusRing(false).Draw(func(p *ui.Painter, r ui.Rect) {
		cx, cy := r.X+r.W/2, r.Y+r.H/2
		w, h := 4*scale, 2*scale
		var tri ui.Path
		if back {
			tri.MoveTo(cx-w, cy+h).LineTo(cx, cy-h).LineTo(cx+w, cy+h).Close()
		} else {
			tri.MoveTo(cx-w, cy-h).LineTo(cx, cy+h).LineTo(cx+w, cy-h).Close()
		}
		p.FillPath(&tri, ink)
	}).DrawOver(func(p *ui.Painter, r ui.Rect) {
		if focused {
			fluentFocusRing(p, t, r, 2, 2, 2, 2, 4)
		}
	})
}

// fluentCalendarDay draws CalendarViewDayItem e of day d, as CalendarViewBaseItemChrome draws it with CalendarViewBaseItemRoundedChromeEnabled: a circle as large as the item (radius half its 40 DIPs). Today: TodayBackground (AccentFillColorDefault; TodayHoverBackground, AccentFillColorSecondary, under the pointer, TodayPressedBackground, Tertiary, pressed) with TodayForeground (TextOnAccentFillColorPrimary), and when it is chosen too an inner 1-DIP ring of TodaySelectedInnerBorderBrush (TextOnAccentFillColorPrimary) 1 DIP in, today winning over the chosen day's ring. The chosen day: a ring of CalendarItemBorderThickness (1) in SelectedBorderBrush (AccentFillColorDefault), SelectedHoverBorderBrush (AccentFillColorSecondary) under the pointer and SelectedPressedBorderBrush (SubtleFillColorTertiary) pressed, its number in SelectedForeground (AccentTextFillColorPrimary; Tertiary pressed), and under the pointer and pressed the fills any day takes there. Another day: CalendarItemHoverBackground with HoverBorderBrush (both SubtleFillColorSecondary) under the pointer, Pressed and PressedBorderBrush (SubtleFillColorTertiary) pressed with PressedForeground (TextFillColorSecondary); its number in CalendarItemForeground (TextFillColorPrimary), or OutOfScopeForeground (TextFillColorSecondary) out of the month, Primary under the pointer and Tertiary pressed (#72000000, #87FFFFFF). The number, in DayItemFontSize 14, sits half a DIP above the middle: the item's Padding 0,0,0,4 under DayItemMargin 0,3,0,0 (inferred). The system focus visual goes 3 DIPs around the item the keys are on (FocusVisualMargin -3), round with it.
func fluentCalendarDay(c *ui.Context, b *calendarBuild, th *ui.Theme, bg ui.Color, e ui.Element, d time.Time, ringed bool) {
	t := b.t
	today, chosen, out := sameDay(d, b.today), sameDay(d, b.chosen), d.Month() != b.month.Month()
	hovered, pressed := e.Hovered(), e.Pressed()
	hover, press := fluentSubtle(th)
	var fill, ring, inner ui.Color
	fg := th.Text
	switch {
	case today:
		fill, fg = t.Accent, t.AccentText
		switch {
		case pressed:
			fill = bg.Mix(t.Accent, 0.8)
		case hovered:
			fill = bg.Mix(t.Accent, 0.9)
		}
		if chosen {
			inner = t.AccentText
		}
	case chosen:
		ring, fg = t.Accent, fluentLinkColor(c, false, pressed)
		// GetItemBackgroundBrush gives a chosen day the fills of any other under the pointer and pressed.
		switch {
		case pressed:
			fill, ring = press, press
		case hovered:
			fill, ring = hover, bg.Mix(t.Accent, 0.9)
		}
	case pressed:
		fill, ring, fg = press, press, th.TextMuted
		if out {
			fg = fluentTertiaryOver(t, bg)
		}
	case hovered:
		fill, ring = hover, hover
	case out:
		fg = th.TextMuted
	}
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		if fill.A > 0 {
			p.Fill(r, fill, r.W/2)
		}
		// Painter.Stroke draws inside the rectangle it is given.
		if ring.A > 0 {
			p.Stroke(r, ring, r.W/2, 1)
		}
		if inner.A > 0 {
			p.Stroke(ui.Rect{X: r.X + 1, Y: r.Y + 1, W: r.W - 2, H: r.H - 2}, inner, r.W/2-1, 1)
		}
	}).DrawOver(func(p *ui.Painter, r ui.Rect) {
		if ringed {
			fluentFocusRing(p, t, r, 3, 3, 3, 3, r.W/2)
		}
	})
	e.Children(func() {
		ui.Text(c, dateText(b.locale, "d#", d)).SingleLine().FontSize(t.FontSize).Margin(0, 0, 1, 0).TextColor(fg)
	})
}

// ---- GNOME: GtkCalendar as libadwaita styles it ----

// adwCalendar builds a GtkCalendar into root, CalendarView's GNOME look (gtkcalendar.c's widgets, _calendar.scss): calendar.view, view_bg_color and view_fg_color inside a 1-pixel $border_color border without corners, or in a popover none of these, the calendar lying on the popover (popover > contents > .view); the header box, a 1-pixel $border_color border along its bottom, holding the image buttons (34 pixels square: min-width and min-height 24 with padding 5, %button_basic_flat with border-radius 0: transparent, $hover_color 7% under the pointer, $active_color 16% pressed) with pan-start-symbolic and pan-end-symbolic (their -rtl variants right to left), around the month's name in a homogeneous stack as wide as the longest of the twelve, then free room (the next month's button hexpands, aligned at its start), and the buttons around the year, whose label is 10 pixels wider than its natural width; the calendar's text in tabular digits (font-variant-numeric: tabular-nums); the grid, padded 3 at the left and the bottom, homogeneous: each label a margin of 3 above and at the right, the day numbers padded 3 (rows 27 high for a label line of 18, Adwaita Sans 11's at 96 DPI, inferred). Its labels are the abbreviated day names (strftime's %a), then the days: in view_fg_color, the other months' (.other-month) at 30%, today (.today) with box-shadow inset 0 -2px $border_color, the chosen one (:selected) accent_bg_color with accent_fg_color and corners of $button_radius (9), and the focus ring (the focus-ring mixin on :focus: 2 pixels of $focus_border_color inside the label) on the one the keys are on while the calendar has the focus, from a click too.
func adwCalendar(b *calendarBuild, root ui.Element) {
	c, t := b.c, b.t
	under := adwListViewBg(t)
	if b.bare {
		under = pick(t.Dark, white, ui.Hex("#36363a")) // adwMenu's popover_bg_color
	} else {
		// The border lies outside the background (background-clip: padding-box), so its 15% of currentColor falls on the window's background: 220 in the light theme and 67 in the dark one in libadwaita 1.10.0's render, where the header's line over view_bg_color is 224 and 63.
		root.Background(under).Border(1, adwInkAlpha(t, t.Background, 0.15))
	}
	fg := adwTextOn(t, under)
	border := adwInkAlpha(t, under, 0.15)
	lineH := float32(18) * t.FontSize / 14.67
	months := i18n.Months(b.locale, i18n.Wide, true)
	root.Children(func() {
		ui.Row(c).Shrink(0).AlignItems(ui.Stretch).Children(func() {
			arrow := func(key string, back bool, delta int) {
				btn := ui.ButtonBase(c).Size(34, 34).Shrink(0).Label(say(c, key))
				adwCalendarArrow(t, btn, under, fg, back != b.rtl)
				if btn.Clicked() {
					b.move(addMonths(b.cursor, delta), true)
				}
			}
			arrow("calendarPreviousMonth", true, -1)
			ui.Column(c).Shrink(0).Justify(ui.Center).AlignItems(ui.Center).Children(func() {
				ui.Text(c, months[b.month.Month()-1]).SingleLine().TextColor(fg)
				// The stack is as wide as its widest page, a month's name.
				ui.Text(c, strings.Join(months[:], "\n")).Height(0).Invisible().Role(ui.RoleNone)
			})
			arrow("calendarNextMonth", false, 1)
			ui.Box(c).Grow(1).Role(ui.RoleNone)
			arrow("calendarPreviousYear", true, -12)
			ui.Box(c).Shrink(0).PaddingX(5).Center().Children(func() {
				ui.Text(c, i18n.FormatNumber(b.locale, b.month.Year(), i18n.NumberOptions{NoGrouping: true})).SingleLine().FontFeatures("tnum").TextColor(fg)
			})
			arrow("calendarNextYear", false, 12)
		})
		ui.Box(c).Height(1).Shrink(0).Background(border)
		grid := ui.Column(c).Shrink(0).Padding(0, 0, 3, 3).MinWidth(3 + 8*adwCalendarColumn(b)).Focusable().FocusRing(false)
		b.keys = grid
		b.focused = grid.Focused()
		// The focus-ring mixin is included with $focus-state ':focus', and GtkCalendar sets GTK_STATE_FLAG_FOCUSED itself on the label at focus_row and focus_col whatever moved the focus there, a click too: the ring shows while the calendar has the focus, not only from the keyboard. Those are -1 until a click on a day (calendar_select_and_focus_day) or an arrow (move_focus) sets them, so a calendar that has just taken the focus, by Tab or as its popover opens, shows no ring: libadwaita 1.10.0's render of the datepickernoarrow set shows none in the open popover, whose calendar has the focus.
		focusRing := b.focused && b.placed
		grid.Children(func() {
			calendarDays(b, Linux, func(name string) {
				ui.Box(c).Grow(1).Basis(0).Padding(3, 3, 0, 0).Role(ui.RoleColumnHeader).Children(func() {
					ui.Box(c).Grow(1).Center().Children(func() {
						ui.Text(c, name).SingleLine().TextColor(fg)
					})
				})
			}, func(week int, fn func()) ui.Element {
				return ui.Row(c).Height(lineH + 9).Shrink(0).AlignItems(ui.Stretch).Children(fn)
			}, func(d time.Time) ui.Element {
				var label ui.Element
				ui.Box(c).Grow(1).Basis(0).Padding(3, 3, 0, 0).Role(ui.RoleNone).Children(func() {
					label = ui.Box(c).Grow(1).Padding(3).Center()
					adwCalendarDay(c, b, under, fg, border, label, d, focusRing && sameDay(d, b.cursor))
				})
				return label
			})
		})
	})
}

// adwCalendarColumn returns the natural width of the widest label of GtkCalendar's grid, its margin of 3 included: a day name's text, or a day number's in tabular digits with its padding of 3 at either side, in GTK's whole pixels (gtk_label_measure rounds up). The grid is homogeneous, and GtkGridLayout counts the hidden week-number column in its request (grid_request_homogeneous gives every line the widest child's width, the empty one too, and grid_request_sum adds them all), so the grid asks for 3 and 8 such columns, which the 7 shown ones share: in libadwaita 1.10.0's render, with Adwaita Sans's "Wed" 35 wide with its margin, the grid is 283 wide and its columns 40.
func adwCalendarColumn(b *calendarBuild) float32 {
	c, t := b.c, b.t
	widest := float32(0)
	for _, name := range weekdayNames(b.locale, Linux) {
		w, _ := c.MeasureText(0, ui.Span{Text: name, Size: t.FontSize})
		widest = max(widest, float32(math.Ceil(float64(w)))+3)
	}
	start := b.start(Linux)
	for i := range 42 {
		w, _ := c.MeasureText(0, ui.Span{Text: dateText(b.locale, "d#", start.AddDate(0, 0, i)), Size: t.FontSize, Features: "tnum"})
		widest = max(widest, float32(math.Ceil(float64(w)))+9)
	}
	return widest
}

// adwCalendarArrow draws one of GtkCalendar's header buttons, btn, over under: flat and square (%button_basic_flat, border-radius 0), currentColor at 7% under the pointer and 16% pressed, its focus ring inside its edge, and the 16-pixel pan-start-symbolic (back) or pan-end-symbolic icon in fg: a 2-pixel round stroke along (9,4), (5,8), (9,12), or (7,4), (11,8), (7,12) (adwaita-icon-theme 51).
func adwCalendarArrow(t *ui.Theme, btn ui.Element, under, fg ui.Color, back bool) {
	switch {
	case btn.Pressed():
		btn.Background(adwInkAlpha(t, under, 0.16))
	case btn.Hovered():
		btn.Background(adwInkAlpha(t, under, 0.07))
	}
	adwFocusRing(t, btn, 0, nil)
	btn.Draw(func(p *ui.Painter, r ui.Rect) {
		x, y := r.X+r.W/2-8, r.Y+r.H/2-8
		var path ui.Path
		if back {
			path.MoveTo(x+9, y+4).LineTo(x+5, y+8).LineTo(x+9, y+12)
		} else {
			path.MoveTo(x+7, y+4).LineTo(x+11, y+8).LineTo(x+7, y+12)
		}
		p.StrokePath(&path, 2, fg)
	})
}

// adwCalendarRing returns the color of the focus ring over the opaque fill under: $focus_border_color, accent_color at 50%. In the dark theme with the default accent, accent_color lies outside sRGB, (128.6, 208.4, 310.8) (see adwTheme), and GTK 4.24 clamps the mix only once made, so over the chosen day's accent the ring reads (90, 170, 255) in libadwaita 1.10.0's render of tools/adwaita-reference's datepickerfocus set, where the theme's Focus, fitted to GTK's ring over an entry's and a button's faces, would give (91, 160, 245); the ring is that mix here. Otherwise it is the theme's Focus over under, (28, 114, 209) over the accent in the light theme, as in the render.
func adwCalendarRing(t *ui.Theme, under ui.Color) ui.Color {
	if !t.Dark || t.Accent != adwAccent {
		return t.Focus.Over(under)
	}
	mix := func(wide float32, u uint8) uint8 { return uint8(min(255, (wide+float32(u))/2+0.5)) }
	return ui.RGB(mix(128.6, under.R), mix(208.4, under.G), mix(310.8, under.B))
}

// adwCalendarDay draws GtkCalendar's day-number label e of day d over under (adwCalendar).
func adwCalendarDay(c *ui.Context, b *calendarBuild, under, fg, border ui.Color, e ui.Element, d time.Time, ringed bool) {
	t := b.t
	chosen := sameDay(d, b.chosen) && d.Month() == b.month.Month()
	today := sameDay(d, b.today) && d.Month() == b.month.Month()
	color := fg
	switch {
	case chosen:
		e.Background(t.Accent).Radius(9)
		color = t.AccentText
	case d.Month() != b.month.Month():
		color = adwInkAlpha(t, under, 0.3)
	}
	e.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if today && !chosen {
			p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 2, W: r.W, H: 2}, border, 0)
		}
		if ringed {
			p.Stroke(r, adwCalendarRing(t, pick(chosen, under, t.Accent)), pick(chosen, float32(0), 9), 2)
		}
	})
	e.Children(func() {
		ui.Text(c, dateText(b.locale, "d#", d)).SingleLine().FontFeatures("tnum").TextColor(color)
	})
}

// ---- the date picker: a button opening the calendar ----

// DatePicker creates a button showing *date in the user's regional formats, as the platform's date picker shows it (MyGo's short date, "11/18/2026" in the United States), that opens CalendarView in a popover below it; a zero *date shows WinUI's placeholder, "select a date", on every platform. A click on a day, or Enter or Space on the day the keys are on, chooses it and closes the popover, the focus going back to the button; Escape or a click outside closes it. The keys in the calendar are CalendarView's: on macOS the arrows choose as they move, and the popover stays open.
//   - macOS: SwiftUI's DatePicker of the compact style without its stepper, which is NSDatePicker of style textField with presentsCalendarOverlay: the date field's bezel (macDatePicker: 23 points high with its line, corners of 7), the date 3.5 points in (its tabular digits' ink from 5.5); the calendar overlay, the graphical calendar on a panel like a menu (macMenu) padded 6, below the field's start (inferred: the overlay is a window of its own that the renders do not show).
//   - Windows: CalendarDatePicker: a button's face 32 DIPs high (fluentDateTimeFace), the date (DateText, Padding 12,0,0,2) and the calendar glyph E787 at 12 DIPs in TextFillColorSecondary, TextFillColorPrimary at rest once a date is chosen, centered in a 32-DIP column at the end (drawn as a calendar's outline, inferred), the system focus visual 3 DIPs around it (FocusVisualMargin -3); its Flyout, Placement Bottom, centered below it with a FlyoutPresenter of Padding 0, BorderThickness 0, corners of OverlayCornerRadius (8) and the default shadow, holding the CalendarView, which closes once a day is chosen.
//   - GNOME: a GtkMenuButton showing the date, as GNOME's apps open a GtkCalendar from one: a raised button (adwButton) padded 5px 9px, the label 4 pixels before pan-down-symbolic, $button_checked_color (30%) while the popover is open; the popover (adwPopover) centered below it, holding the calendar without its border.
//
// Assistive technology sees a pop-up button whose value is the date, expanded while the calendar is open. It returns the button, whose Changed reports a day chosen.
func DatePicker(c *ui.Context, date *time.Time) ui.Element {
	p := PlatformOf(c)
	_, locale := dateFormats(c)
	text := say(c, "datePickerPlaceholder")
	if !date.IsZero() {
		text = dateText(locale, "yMd", *date)
	}
	b := ui.ButtonBase(c).Shrink(0)
	open := ui.Local(b, "breeze-datepicker-open", func() bool { return false })
	if b.Clicked() {
		*open = !*open
	}
	b.Role(ui.RolePopUpButton).Value(text).Expanded(*open)
	switch p {
	case Windows:
		fluentDatePicker(c, b, text, !date.IsZero())
	case Linux:
		adwDatePicker(c, b, text, *open)
	default:
		macDatePicker(c, b, false, nil)
		b.MinWidth(0).Padding(2, 4).PaddingStart(3.5).AlignItems(ui.Center).Children(func() {
			ui.Text(c, text).SingleLine().FontFeatures("tnum")
		})
	}
	before := *date
	panel := ui.PopoverBase(c, b, open, func(panel ui.Element) {
		Preview(c, p, func() {
			switch p {
			case Windows:
				t := c.Theme()
				panel.AttachTo(b, ui.AnchorBottom, ui.AnchorTop).Margin(4, 0, 0, 0).Radius(8).Clip().
					Shadow(0, 8, 17, 0, ui.RGBA(0, 0, 0, pick(t.Dark, float32(0.14), 0.26)))
			case Linux:
				// GtkPopover without its arrow (has-arrow false, as MyGo draws popovers) puts its contents right at the button's bottom edge, centered on it: in libadwaita 1.10.0's render (tools/adwaita-reference's datepickernoarrow set) the contents span 269 to 568 under a button spanning 360 to 476 and ending at 50, where they start. adwPopover's margin of 6 is a dropdown's.
				adwPopover(c, panel)
				panel.AttachTo(b, ui.AnchorBottom, ui.AnchorTop).Margin(0)
			default:
				macMenu(c, panel)
				panel.Padding(6)
			}
			calendarView(c, date, func() {
				*open = false
				b.Focus()
			})
		})
	})
	if panel.Dismissed() {
		b.Focus()
	}
	if !before.Equal(*date) {
		b.SetChanged()
	}
	return b
}

// fluentDatePicker draws CalendarDatePicker's button b showing text (DatePicker), a date when dated. The calendar glyph is CalendarDatePickerCalendarGlyphForeground (TextFillColorSecondary), and with a date the Selected state's CalendarDatePickerTextForegroundSelected (TextFillColorPrimary), which the PointerOver and Pressed states, started after it, cover with their TextFillColorSecondary (calendardatepicker_partial.cpp's UpdateVisualState).
func fluentDatePicker(c *ui.Context, b ui.Element, text string, dated bool) {
	t := c.Theme()
	hovered, pressed := b.Hovered(), b.Pressed()
	fluentDateTimeFace(c, b, hovered, pressed)
	b.MinWidth(0).Padding(0).PaddingStart(12).AlignItems(ui.Center).FocusRing(false)
	focused := b.FocusVisible()
	b.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if focused {
			fluentFocusRing(p, t, r, 3, 3, 3, 3, 4)
		}
	})
	glyph := t.TextMuted
	if dated && !hovered && !pressed {
		glyph = t.Text
	}
	b.Children(func() {
		ui.Text(c, text).SingleLine().Margin(0, 0, 2, 0)
		ui.Box(c).Size(32, 30).Shrink(0).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			// The Calendar glyph at 12 DIPs: an outline 11 by 10 with corners of 1.5 and the line under its top band, a DIP thick (inferred).
			x, y := r.X+r.W/2-5.5, r.Y+r.H/2-5
			p.Stroke(ui.Rect{X: x + 0.5, Y: y + 0.5, W: 10, H: 9}, glyph, 1.5, 1)
			p.Fill(ui.Rect{X: x + 1, Y: y + 3, W: 9, H: 1}, glyph, 0)
		})
	})
}

// adwDatePicker draws the GtkMenuButton b showing text (DatePicker), checked while its popover is open.
func adwDatePicker(c *ui.Context, b ui.Element, text string, open bool) {
	t := c.Theme()
	adwButton(c, b, false)
	b.Padding(5, 9).Gap(4)
	if open {
		b.Background(adwInkLayer(t, 0.30))
	}
	color := adwTextOn(t, adwInkAlpha(t, t.Background, 0.10))
	b.Children(func() {
		ui.Text(c, text).SingleLine()
		idiomAdwArrow(c, color)
	})
}
