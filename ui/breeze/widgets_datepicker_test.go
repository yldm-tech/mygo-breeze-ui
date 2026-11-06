package breeze

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

func init() {
	drawnWidgets["CalendarView"] = true
	drawnWidgets["DatePicker"] = true
}

// datePickerToday is the day the tests' calendars mark as today, a Friday in the month they show.
var datePickerToday = time.Date(2026, time.November, 6, 9, 30, 0, 0, time.UTC)

// fixToday makes the calendars mark day as today until the test ends.
func fixToday(t *testing.T, day time.Time) {
	t.Helper()
	prev := calendarNow
	calendarNow = func(*ui.Context) time.Time { return day }
	t.Cleanup(func() { calendarNow = prev })
}

// dateTester tests view on p in a window w by h in American English, whose weeks start on Sunday, with the platform's text size.
func dateTester(p Platform, w, h int, view func(c *ui.Context)) *ui.Tester {
	tester := ui.NewTester(func(c *ui.Context) { Preview(c, p, func() { Use(c); view(c) }) }, w, h)
	tester.SetLocale("en-US")
	tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tester.Frame()
	return tester
}

// datePickerView lays out the date pickers for screenshots: CalendarView at 16, 16 showing November 2026 with the 18th chosen, and a closed DatePicker of the same day at 360, 16. tools/appkit-reference's datepicker scene, tools/winui-reference's winui-datepicker scene and tools/adwaita-reference's datepicker set place AppKit's, WinUI's and GTK's there.
func datePickerView(p Platform) func(c *ui.Context) {
	day := time.Date(2026, time.November, 18, 10, 0, 0, 0, time.UTC)
	picked := day
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Row(c).Fill().Padding(16).AlignItems(ui.Start).Children(func() {
				ui.Column(c).Width(344).Shrink(0).AlignItems(ui.Start).Children(func() {
					CalendarView(c, &day)
				})
				DatePicker(c, &picked)
			})
		})
	}
}

// TestScreenshotsDatePicker writes the date pickers of every platform at 2x, light and dark, and right to left in the light theme, to the directory BREEZE_SCREENSHOTS names, in American English with today the 6th of November 2026. Without it the test only builds the frames.
func TestScreenshotsDatePicker(t *testing.T) {
	fixToday(t, datePickerToday)
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, shot := range []struct {
			name string
			dark bool
			dir  ui.Direction
		}{{"light", false, ui.LeftToRight}, {"dark", true, ui.LeftToRight}, {"rtl-light", false, ui.RightToLeft}} {
			tester := ui.NewTester(datePickerView(p), 560, 400)
			tester.SetLocale("en-US")
			tester.SetScale(2)
			tester.SetDark(shot.dark)
			tester.SetDirection(shot.dir)
			tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tester.Frame()
			tester.Frame()
			for _, s := range []string{"18", "11/18/2026"} {
				if !tester.HasText(s) {
					t.Fatalf("%v %s: %q is missing", p, shot.name, s)
				}
			}
			writeShot(t, dir, fmt.Sprintf("breeze-datepicker-%s-%s.png", strings.ToLower(p.String()), shot.name), tester)
		}
	}
}

// calendarApp holds a calendar's date as a test keeps it, with how many times Changed reported a new one.
type calendarApp struct {
	date    time.Time
	changes int
}

// view builds a calendar after a button, which Tab starts from, and a button after it.
func (a *calendarApp) view(c *ui.Context) {
	ui.Column(c).Padding(16).Gap(12).AlignItems(ui.Start).Children(func() {
		Button(c, "Before")
		if CalendarView(c, &a.date).Changed() {
			a.changes++
		}
		Button(c, "After")
	})
}

// dayName is what assistive technology calls the day d of 2026, its date in American English.
func dayName(month time.Month, d int) string {
	return dateText("en-US", "yMMMMd", time.Date(2026, month, d, 0, 0, 0, 0, time.UTC))
}

// calendarTitleFor returns the text that names the month m of 2026 in a calendar on p in locale: NSDatePicker's abbreviated month and year, CalendarView's month and year, GtkCalendar's month alone.
func calendarTitleFor(p Platform, locale string, m time.Month) string {
	first := time.Date(2026, m, 1, 0, 0, 0, 0, time.UTC)
	switch p {
	case Windows:
		return dateText(locale, "yMMMM", first)
	case Linux:
		return i18nMonth(locale, m)
	}
	return dateText(locale, "yMMM", first)
}

// i18nMonth returns the name of month m in locale standing alone, as GtkCalendar's header shows it.
func i18nMonth(locale string, m time.Month) string { return i18n.Months(locale, i18n.Wide, true)[m-1] }

// TestCalendarViewChoosesDays checks on every platform, left to right and right to left, that a click chooses a day, which Changed reports, and that the keys move as the platform's calendar moves them: on macOS the arrows choose as they go, as NSDatePicker's do; on Windows and GNOME they move among the days and Enter chooses, as WinUI's CalendarView and GtkCalendar do; Page Down shows the next month on the same day, and Home and End go to its first and last days.
func TestCalendarViewChoosesDays(t *testing.T) {
	fixToday(t, datePickerToday)
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			app := &calendarApp{date: time.Date(2026, time.November, 18, 10, 30, 0, 0, time.UTC)}
			tester := dateTester(p, 640, 520, app.view)
			tester.SetDirection(dir)
			tester.Frame()
			if err := tester.Click(dayName(time.November, 20)); err != nil {
				t.Fatalf("%v %v: %v", p, dir, err)
			}
			tester.Frame()
			tester.Frame()
			if want := time.Date(2026, time.November, 20, 10, 30, 0, 0, time.UTC); !app.date.Equal(want) || app.changes != 1 {
				t.Fatalf("%v %v: a click on the 20th left %v with %d changes", p, dir, app.date, app.changes)
			}
			press := func(k ui.Key) {
				tester.Key(0, k)
				tester.Frame()
				tester.Frame()
			}
			// The day after, toward the end of the line: Right left to right, Left right to left.
			forth := ui.KeyRight
			if dir == ui.RightToLeft {
				forth = ui.KeyLeft
			}
			press(forth)
			moved := app.date.Day() == 21
			if moved != (p == Mac) {
				t.Errorf("%v %v: the arrow toward the line's end left the 21st chosen: %v", p, dir, moved)
			}
			press(ui.KeyEnter)
			if app.date.Day() != 21 {
				t.Errorf("%v %v: Enter after the arrow left %v", p, dir, app.date)
			}
			press(ui.KeyPageDown)
			if !tester.HasText(calendarTitleFor(p, "en-US", time.December)) {
				t.Errorf("%v %v: Page Down did not show December", p, dir)
			}
			press(ui.KeyEnd)
			press(ui.KeyEnter)
			if want := time.Date(2026, time.December, 31, 10, 30, 0, 0, time.UTC); !app.date.Equal(want) {
				t.Errorf("%v %v: Page Down, End and Enter chose %v, want %v", p, dir, app.date, want)
			}
			press(ui.KeyHome)
			press(ui.KeySpace)
			if app.date.Day() != 1 || app.date.Month() != time.December {
				t.Errorf("%v %v: Home and Space chose %v", p, dir, app.date)
			}
		}
	}
}

// TestCalendarViewHeaderButtons checks the header's buttons on every platform, named as the platform names them in breeze's catalog: on macOS ◀, ● and ▶ choose the same day a month before, today and a month after, as NSDatePicker's do; on Windows Previous and Next only show another month; on GNOME the month's and the year's arrows choose the same day there, as GtkCalendar's do, the last day of a shorter month.
func TestCalendarViewHeaderButtons(t *testing.T) {
	fixToday(t, datePickerToday)
	for _, p := range Platforms {
		app := &calendarApp{date: time.Date(2026, time.October, 31, 8, 0, 0, 0, time.UTC)}
		tester := dateTester(p, 640, 520, app.view)
		click := func(s string) {
			t.Helper()
			if err := tester.Click(s); err != nil {
				t.Fatalf("%v: %v", p, err)
			}
			tester.Frame()
			tester.Frame()
		}
		switch p {
		case Mac:
			click("Next Month")
			if app.date.Month() != time.November || app.date.Day() != 30 {
				t.Errorf("%v: Next Month chose %v, want the 30th of November", p, app.date)
			}
			click("Today")
			if !sameDay(app.date, datePickerToday) || app.date.Hour() != 8 {
				t.Errorf("%v: Today chose %v", p, app.date)
			}
			click("Previous Month")
			if app.date.Month() != time.October || app.date.Day() != 6 {
				t.Errorf("%v: Previous Month chose %v", p, app.date)
			}
		case Windows:
			click("Next")
			if app.date.Month() != time.October || !tester.HasText("November 2026") {
				t.Errorf("%v: Next left %v and did not show November", p, app.date)
			}
			click("Previous")
			click("Previous")
			if !tester.HasText("September 2026") || app.changes != 0 {
				t.Errorf("%v: Previous twice did not show September, or chose a day", p)
			}
		case Linux:
			click("Next Month")
			if app.date.Month() != time.November || app.date.Day() != 30 {
				t.Errorf("%v: Next Month chose %v, want the 30th of November", p, app.date)
			}
			click("Next Year")
			if app.date.Year() != 2027 || app.date.Month() != time.November {
				t.Errorf("%v: Next Year chose %v", p, app.date)
			}
			click("Previous Year")
			click("Previous Month")
			if app.date.Year() != 2026 || app.date.Month() != time.October || app.date.Day() != 30 {
				t.Errorf("%v: Previous Year and Previous Month chose %v", p, app.date)
			}
		}
	}
}

// TestCalendarViewTabStops checks on every platform where Tab stops: on macOS once, at the calendar, as NSDatePicker is one control; on Windows at Previous, Next and the days; on GNOME at the four arrows and the days.
func TestCalendarViewTabStops(t *testing.T) {
	fixToday(t, datePickerToday)
	stops := map[Platform][]string{
		Mac:     {dayName(time.November, 18)},
		Windows: {"Previous", "Next", dayName(time.November, 18)},
		Linux:   {"Previous Month", "Next Month", "Previous Year", "Next Year", dayName(time.November, 18)},
	}
	for _, p := range Platforms {
		app := &calendarApp{date: time.Date(2026, time.November, 18, 0, 0, 0, 0, time.UTC)}
		tester := dateTester(p, 640, 520, app.view)
		if err := tester.Click("Before"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		for _, want := range append(stops[p], "After") {
			tester.Key(0, ui.KeyTab)
			tester.Frame()
			if !tester.Focused(want) {
				t.Fatalf("%v: Tab did not go to %q", p, want)
			}
		}
	}
}

// TestCalendarViewFollowsTheLocale checks on every platform that the window's language names the month and puts the week's first day where its region does: in French the month is novembre and weeks start on Monday, so the 2nd of November, a Monday, comes before the 1st, a Sunday, on the line; in American English the 1st comes first.
func TestCalendarViewFollowsTheLocale(t *testing.T) {
	fixToday(t, datePickerToday)
	for _, p := range Platforms {
		for _, locale := range []string{"en-US", "fr-FR"} {
			app := &calendarApp{date: time.Date(2026, time.November, 18, 0, 0, 0, 0, time.UTC)}
			tester := dateTester(p, 640, 520, app.view)
			tester.SetLocale(locale)
			tester.Frame()
			title := calendarTitleFor(p, locale, time.November)
			if _, ok := tester.Find(title); !ok {
				t.Errorf("%v %s: no title %q", p, locale, title)
			}
			name := func(d int) string {
				return dateText(locale, "yMMMMd", time.Date(2026, time.November, d, 0, 0, 0, 0, time.UTC))
			}
			first, ok1 := tester.Find(name(1))
			second, ok2 := tester.Find(name(2))
			if !ok1 || !ok2 {
				t.Fatalf("%v %s: the days are not named %q and %q", p, locale, name(1), name(2))
			}
			if mondayFirst := second.X < first.X; mondayFirst != (locale == "fr-FR") {
				t.Errorf("%v %s: Monday the 2nd at %v, Sunday the 1st at %v", p, locale, second.X, first.X)
			}
		}
	}
}

// TestCalendarViewSizes pins the platforms' sizes: NSDatePicker's calendar 139 by 148 points with its first week's first day 4 points in and 38 down, 18.5 by 18 (measured on macOS 27.2); WinUI's CalendarView 300 by 351 DIPs with 40-DIP day items, the first 4 in and 97 down (its border, ScrollViewer and item margins, the 48-DIP header, its line and the 42-DIP row of names; CalendarView_themeresources.xaml); GtkCalendar's 34-pixel arrows and grid rows 27 apart, the first week starting a week before November's first, a Sunday (gtkcalendar.c).
func TestCalendarViewSizes(t *testing.T) {
	fixToday(t, datePickerToday)
	for _, p := range Platforms {
		day := time.Date(2026, time.November, 18, 0, 0, 0, 0, time.UTC)
		var box ui.Rect
		tester := dateTester(p, 640, 520, func(c *ui.Context) {
			ui.Row(c).Padding(16).AlignItems(ui.Start).Children(func() {
				box = CalendarView(c, &day).Bounds()
			})
		})
		near := func(what string, got, want float32) {
			t.Helper()
			if got < want-0.01 || got > want+0.01 {
				t.Errorf("%v: %s is %v, want %v", p, what, got, want)
			}
		}
		switch p {
		case Mac:
			near("the width", box.W, 139)
			near("the height", box.H, 148)
			one, _ := tester.Find(dayName(time.November, 1))
			near("the 1st's x", one.X-box.X, 4)
			near("the 1st's y", one.Y-box.Y, 38)
			near("a day's width", one.W, 18.5)
			near("a day's height", one.H, 18)
		case Windows:
			near("the width", box.W, 300)
			near("the height", box.H, 351)
			one, _ := tester.Find(dayName(time.November, 1))
			near("the 1st's x", one.X-box.X, 4)
			near("the 1st's y", one.Y-box.Y, 97)
			near("a day's width", one.W, 40)
			near("a day's height", one.H, 40)
		case Linux:
			arrow, _ := tester.Find("Previous Month")
			near("an arrow's width", arrow.W, 34)
			near("an arrow's height", arrow.H, 34)
			if _, ok := tester.Find(dayName(time.October, 25)); !ok {
				t.Errorf("%v: the first week does not start on the 25th of October", p)
			}
			one, _ := tester.Find(dayName(time.November, 1))
			eight, _ := tester.Find(dayName(time.November, 8))
			near("the rows' pitch", eight.Y-one.Y, 27)
			near("the 1st's x", one.X-box.X, 1+3)
		}
	}
}

// TestCalendarViewWithoutADate checks on every platform that a zero date shows today's month with no day chosen, and that a click chooses a day in the local time zone.
func TestCalendarViewWithoutADate(t *testing.T) {
	fixToday(t, datePickerToday)
	for _, p := range Platforms {
		app := &calendarApp{}
		tester := dateTester(p, 640, 520, app.view)
		if !tester.HasText(calendarTitleFor(p, "en-US", time.November)) {
			t.Fatalf("%v: a zero date does not show today's month", p)
		}
		if err := tester.Click(dayName(time.November, 3)); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		tester.Frame()
		if y, m, d := app.date.Date(); y != 2026 || m != time.November || d != 3 || app.date.Location() != time.Local {
			t.Errorf("%v: a click on the 3rd left %v", p, app.date)
		}
	}
}

// TestCalendarViewDisabled checks on every platform that a disabled calendar takes neither clicks nor keys.
func TestCalendarViewDisabled(t *testing.T) {
	fixToday(t, datePickerToday)
	for _, p := range Platforms {
		date := time.Date(2026, time.November, 18, 0, 0, 0, 0, time.UTC)
		want := date
		tester := dateTester(p, 640, 520, func(c *ui.Context) {
			CalendarView(c, &date).Disabled(true)
		})
		tester.Click(dayName(time.November, 20))
		tester.Frame()
		tester.Key(0, ui.KeyTab)
		tester.Key(0, ui.KeyRight)
		tester.Key(0, ui.KeyEnter)
		tester.Frame()
		tester.Frame()
		if !date.Equal(want) {
			t.Errorf("%v: a disabled calendar chose %v", p, date)
		}
	}
}

// datePickerApp holds a date picker's date as a test keeps it.
type datePickerApp struct {
	date    time.Time
	changes int
}

func (a *datePickerApp) view(c *ui.Context) {
	ui.Column(c).Padding(16).Gap(12).AlignItems(ui.Start).Children(func() {
		if DatePicker(c, &a.date).Changed() {
			a.changes++
		}
		Button(c, "After")
	})
}

// TestDatePickerOpensChoosesAndCloses checks on every platform that the button shows the date in the user's short format, opens the calendar, which takes the focus, that a click on a day chooses it, closes the calendar and gives the button the focus back, that Escape closes it without choosing, and that the keys of the calendar choose a day with Enter.
func TestDatePickerOpensChoosesAndCloses(t *testing.T) {
	fixToday(t, datePickerToday)
	for _, p := range Platforms {
		app := &datePickerApp{date: time.Date(2026, time.November, 18, 9, 0, 0, 0, time.UTC)}
		tester := dateTester(p, 640, 520, app.view)
		if err := tester.Click("11/18/2026"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		tester.Frame()
		if !tester.Focused(dayName(time.November, 18)) {
			t.Fatalf("%v: the open calendar does not have the focus", p)
		}
		if err := tester.Click(dayName(time.November, 24)); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		tester.Frame()
		if !tester.HasText("11/24/2026") || app.date.Day() != 24 || app.changes != 1 {
			t.Fatalf("%v: a click on the 24th left %v, %d changes", p, app.date, app.changes)
		}
		if _, open := tester.Find(dayName(time.November, 24)); open || !tester.Focused("11/24/2026") {
			t.Errorf("%v: the calendar stayed open or the button lost the focus", p)
		}
		tester.Key(0, ui.KeySpace)
		tester.Frame()
		tester.Frame()
		tester.Key(0, ui.KeyEscape)
		tester.Frame()
		tester.Frame()
		if _, open := tester.Find(dayName(time.November, 24)); open || app.date.Day() != 24 {
			t.Errorf("%v: Escape left the calendar open or chose %v", p, app.date)
		}
		tester.Key(0, ui.KeySpace)
		tester.Frame()
		tester.Frame()
		tester.Key(0, ui.KeyDown)
		tester.Frame()
		tester.Key(0, ui.KeyEnter)
		tester.Frame()
		tester.Frame()
		if app.date.Month() != time.December || app.date.Day() != 1 || !tester.HasText("12/1/2026") {
			t.Errorf("%v: Down and Enter chose %v", p, app.date)
		}
	}
}

// TestDatePickerPlaceholder checks on every platform that a zero date shows WinUI's placeholder, "select a date", and that choosing a day replaces it.
func TestDatePickerPlaceholder(t *testing.T) {
	fixToday(t, datePickerToday)
	for _, p := range Platforms {
		app := &datePickerApp{}
		tester := dateTester(p, 640, 520, app.view)
		if err := tester.Click("select a date"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		if err := tester.Click(dayName(time.November, 9)); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		tester.Frame()
		if !tester.HasText("11/9/2026") || tester.HasText("select a date") {
			t.Errorf("%v: choosing the 9th left %v", p, tester.Texts())
		}
	}
}

// TestDatePickerSizes pins the closed pickers' heights: NSDatePicker's textField bezel 23 points (its 21-point frame and the line a point outside it), CalendarDatePicker 32 DIPs, a GtkMenuButton 34 pixels.
func TestDatePickerSizes(t *testing.T) {
	fixToday(t, datePickerToday)
	for _, p := range Platforms {
		date := time.Date(2026, time.November, 18, 0, 0, 0, 0, time.UTC)
		var box ui.Rect
		dateTester(p, 640, 520, func(c *ui.Context) {
			ui.Row(c).Padding(16).AlignItems(ui.Start).Children(func() {
				box = DatePicker(c, &date).Bounds()
			})
		})
		if want := map[Platform]float32{Mac: 23, Windows: 32, Linux: 34}[p]; box.H != want {
			t.Errorf("%v: the picker is %v high, want %v", p, box.H, want)
		}
	}
}

// TestBreezeStringsTakeAppTranslations checks that the names Breeze UI gives the calendar's buttons come from its catalog in the namespace "breeze", which an app's registered catalog translates with keys "breeze.<key>": WinUI's Next button is "Sekva" in a window showing Esperanto when the app says so, and "Next" without a translation.
func TestBreezeStringsTakeAppTranslations(t *testing.T) {
	fixToday(t, datePickerToday)
	app, err := i18n.NewCatalog(fstest.MapFS{
		"en.json": {Data: []byte(`{"@@locale": "en"}`)},
		"eo.json": {Data: []byte(`{"@@locale": "eo", "breeze.calendarNext": "Sekva"}`)},
	}, "en")
	if err != nil {
		t.Fatal(err)
	}
	i18n.Register(app)
	date := time.Date(2026, time.November, 18, 0, 0, 0, 0, time.UTC)
	tester := dateTester(Windows, 640, 520, func(c *ui.Context) { CalendarView(c, &date) })
	if _, ok := tester.Find("Next"); !ok {
		t.Errorf("the Next button is not named Next in English")
	}
	tester.SetLocale("eo")
	tester.Frame()
	if _, ok := tester.Find("Sekva"); !ok {
		t.Errorf("the app's translation of breeze.calendarNext is not used")
	}
}

// TestCalendarViewColorsAsMeasured pins the fills of the chosen day and of today, read from the frame at 2x in the light theme: on macOS the chosen day's rectangle is unemphasizedSelectedContentBackgroundColor (#dcdcdc) while the calendar has no focus and selectedContentBackgroundColor (#0064e1) once it has, on a white field (measured on macOS 27.2); on Windows today's circle is AccentFillColorDefault (#005a9e) and the chosen day's ring the same accent, its inside the calendar's white; on GNOME the chosen day is accent_bg_color (#3584e4).
func TestCalendarViewColorsAsMeasured(t *testing.T) {
	fixToday(t, datePickerToday)
	for _, p := range Platforms {
		date := time.Date(2026, time.November, 18, 0, 0, 0, 0, time.UTC)
		tester := dateTester(p, 640, 520, func(c *ui.Context) {
			ui.Row(c).Padding(16).AlignItems(ui.Start).Children(func() { CalendarView(c, &date) })
		})
		tester.SetScale(2)
		tester.Frame()
		at := func(x, y float32) string {
			c := tester.Image().RGBAAt(int(x*2), int(y*2))
			return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
		}
		// Within 3 levels a channel, for the antialiasing of a ring's curve.
		want := func(what, got, want string) {
			t.Helper()
			g, w := ui.Hex(got), ui.Hex(want)
			for _, d := range []int{int(g.R) - int(w.R), int(g.G) - int(w.G), int(g.B) - int(w.B)} {
				if d < -3 || d > 3 {
					t.Errorf("%v: %s is %s, want %s", p, what, got, want)
					return
				}
			}
		}
		chosen, _ := tester.Find(dayName(time.November, 18))
		today, _ := tester.Find(dayName(time.November, 6))
		switch p {
		case Mac:
			want("the chosen day without the focus", at(chosen.X+1.5, chosen.Y+2), "#dcdcdc")
			want("the field", at(chosen.X+1.5, chosen.Y+17.5), "#ffffff")
			if err := tester.Click(dayName(time.November, 18)); err != nil {
				t.Fatal(err)
			}
			tester.Frame()
			tester.Frame()
			want("the chosen day with the focus", at(chosen.X+1.5, chosen.Y+2), "#0064e1")
		case Windows:
			want("today's circle", at(today.X+20, today.Y+4), "#005a9e")
			want("the chosen day's ring", at(chosen.X+0.75, chosen.Y+20), "#005a9e")
			want("inside the chosen day's ring", at(chosen.X+4, chosen.Y+20), "#ffffff")
		case Linux:
			want("the chosen day", at(chosen.X+2, chosen.Y+chosen.H/2), "#3584e4")
		}
	}
}

// shotPixel returns the pixel of tester's last frame, rendered at 2x, at x, y in DIPs.
func shotPixel(tester *ui.Tester, x, y float32) [3]int {
	c := tester.Image().RGBAAt(int(x*2), int(y*2))
	return [3]int{int(c.R), int(c.G), int(c.B)}
}

// samePixel reports whether a and b are within tolerance in every channel.
func samePixel(a, b [3]int, tolerance int) bool {
	for i := range a {
		if d := a[i] - b[i]; d < -tolerance || d > tolerance {
			return false
		}
	}
	return true
}

// TestMacCalendarAsMeasured pins NSDatePicker's calendar and field as rendered by tools/appkit-reference's datepicker scene on macOS 27.2, where Breeze UI's glyphs then land on AppKit's pixel for pixel: the numbers in 10 points (the height and width of a 10-point text of the same digits, without tabular digits) ending 2 points before their column's end, the day names ending there too and half a point below their 17-point row's middle, the title 3 points in, the field's date 3.5 points in; and the box's edges, half-point lines over the window, at half a point in Aqua (235 over white, the bottom's two lines 194 and 224) and a point in Dark Aqua (52 at the sides and 60 at the top over the window's 30, the bottom's lines 58 and 88 a point higher).
func TestMacCalendarAsMeasured(t *testing.T) {
	fixToday(t, datePickerToday)
	for _, dark := range []bool{false, true} {
		day := time.Date(2026, time.November, 18, 0, 0, 0, 0, time.UTC)
		var box, field, ref ui.Rect
		tester := dateTester(Mac, 640, 520, func(c *ui.Context) {
			ui.Row(c).Padding(16).Gap(16).AlignItems(ui.Start).Children(func() {
				box = CalendarView(c, &day).Bounds()
				field = DatePicker(c, &day).Bounds()
				ref = ui.Text(c, "11").SingleLine().FontSize(c.Theme().FontSize * 10 / 13).Bounds()
			})
		})
		tester.SetDark(dark)
		tester.SetScale(2)
		tester.Frame()
		tester.Frame()
		near := func(what string, got, want float32) {
			t.Helper()
			if got < want-0.01 || got > want+0.01 {
				t.Errorf("dark %v: %s is %v, want %v", dark, what, got, want)
			}
		}
		num, _ := tester.Find("11")
		cell, _ := tester.Find(dayName(time.November, 11))
		near("the 11th's height", num.H, ref.H)
		near("the 11th's width", num.W, ref.W)
		near("the 11th's end", cell.X+cell.W-(num.X+num.W), 2)
		su, _ := tester.Find("Su")
		near("Su's end", su.X+su.W-box.X, 4+18.5-2)
		near("Su's middle", su.Y+su.H/2-box.Y, 20+17.0/2+0.5)
		title, _ := tester.Find("Nov 2026")
		near("the title's start", title.X-box.X, 3)
		date, _ := tester.Find("11/18/2026")
		near("the field's date", date.X-field.X, 3.5)
		px := func(what string, x, y float32, want int) {
			t.Helper()
			if got := shotPixel(tester, box.X+x, box.Y+y); !samePixel(got, [3]int{want, want, want}, 2) {
				t.Errorf("dark %v: %s at %v, %v is %v, want %d", dark, what, x, y, got, want)
			}
		}
		if dark {
			px("the window at the start", 0.75, 70, 30)
			px("the start's line", 1.25, 70, 52)
			px("the end's line", 137.25, 70, 52)
			px("the top's line", 70, 1.25, 60)
			px("the bottom's first line", 70, 146.25, 58)
			px("the bottom's second line", 70, 146.75, 88)
			px("the window under the box", 70, 147.25, 30)
		} else {
			px("the start's line", 0.75, 70, 235)
			px("the end's line", 137.75, 70, 235)
			px("the top's line", 70, 0.75, 235)
			px("the bottom's first line", 70, 147.25, 194)
			px("the bottom's second line", 70, 147.75, 224)
		}
	}
}

// TestFluentCalendarStatesAsMeasured pins CalendarViewBaseItemChrome's states (CalendarView_themeresources.xaml, CalendarViewBaseItemChrome.cpp) and CalendarDatePicker's glyph: a day of another month pressed in the dark theme writes its number in TextFillColorTertiary, #87FFFFFF over the calendar; the chosen day under the pointer takes CalendarItemHoverBackground (SubtleFillColorSecondary, black at 3.5% over white) inside its ring, as GetItemBackgroundBrush gives it whether chosen or not; the picker's glyph is TextFillColorSecondary without a date and TextFillColorPrimary with one (the Selected state); and the CalendarView in the picker's flyout has corners of OverlayCornerRadius (8), its border along the flyout's round corner.
func TestFluentCalendarStatesAsMeasured(t *testing.T) {
	fixToday(t, datePickerToday)
	day := time.Date(2026, time.November, 18, 0, 0, 0, 0, time.UTC)
	tester := dateTester(Windows, 640, 520, func(c *ui.Context) {
		ui.Row(c).Padding(16).AlignItems(ui.Start).Children(func() { CalendarView(c, &day) })
	})
	tester.SetScale(2)
	tester.SetDark(true)
	tester.Frame()
	out, _ := tester.Find(dayName(time.December, 1))
	tester.Press(out.X+out.W/2, out.Y+out.H/2)
	tester.Frame()
	under := shotPixel(tester, out.X+0.5, out.Y+0.5)
	var want, brightest [3]int
	for i := range want {
		want[i] = under[i] + (255-under[i])*0x87/255
	}
	for y := out.Y + 8; y < out.Y+32; y += 0.5 {
		for x := out.X + 8; x < out.X+32; x += 0.5 {
			if p := shotPixel(tester, x, y); p[0] > brightest[0] {
				brightest = p
			}
		}
	}
	if !samePixel(brightest, want, 4) {
		t.Errorf("a pressed day of another month is written in %v, want TextFillColorTertiary %v over %v", brightest, want, under)
	}
	tester.Release(out.X+out.W/2, out.Y+out.H/2)
	tester.Frame()

	// The release chose the 1st of December: a calendar of its own for the pointer.
	day = time.Date(2026, time.November, 18, 0, 0, 0, 0, time.UTC)
	tester = dateTester(Windows, 640, 520, func(c *ui.Context) {
		ui.Row(c).Padding(16).AlignItems(ui.Start).Children(func() { CalendarView(c, &day) })
	})
	tester.SetScale(2)
	tester.Frame()
	chosen, _ := tester.Find(dayName(time.November, 18))
	tester.Move(chosen.X+chosen.W/2, chosen.Y+chosen.H/2)
	tester.Frame()
	if got, want := shotPixel(tester, chosen.X+20, chosen.Y+6), 246; !samePixel(got, [3]int{want, want, want}, 2) {
		t.Errorf("the chosen day under the pointer is %v inside its ring, want SubtleFillColorSecondary %d", got, want)
	}

	glyph := func(date time.Time) int {
		var box ui.Rect
		picker := dateTester(Windows, 640, 520, func(c *ui.Context) {
			ui.Row(c).Padding(16).AlignItems(ui.Start).Children(func() { box = DatePicker(c, &date).Bounds() })
		})
		picker.SetScale(2)
		picker.Frame()
		picker.Frame()
		darkest := 255
		for y := box.Y + 8; y < box.Y+24; y += 0.5 {
			for x := box.X + box.W - 26; x < box.X+box.W-6; x += 0.5 {
				darkest = min(darkest, shotPixel(picker, x, y)[0])
			}
		}
		return darkest
	}
	if dated, placeholder := glyph(day), glyph(time.Time{}); dated > 50 || placeholder < 70 {
		t.Errorf("the glyph's darkest pixel is %d with a date and %d without, want TextFillColorPrimary (under 50) and TextFillColorSecondary (over 70)", dated, placeholder)
	}

	opened := dateTester(Windows, 640, 520, func(c *ui.Context) {
		ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() { DatePicker(c, &day) })
	})
	opened.SetScale(2)
	opened.Frame()
	if err := opened.Click("11/18/2026"); err != nil {
		t.Fatal(err)
	}
	opened.Frame()
	opened.Frame()
	prev, ok := opened.Find("Previous")
	if !ok {
		t.Fatal("the flyout's calendar is not open")
	}
	// The calendar's start is 222 DIPs before Previous and its top 7 above (its border, Previous's margin of 6 and the header's start).
	x, y := prev.X-222, prev.Y-7
	if got := shotPixel(opened, x+2.6, y+2.6); got[0] > 250 {
		t.Errorf("the flyout's calendar has no border along its round corner: %v at 2.6, 2.6", got)
	}
}

// TestAdwCalendarAsStyled pins libadwaita's _calendar.scss: the focus ring of the label the keys are on (the focus-ring mixin with $focus-state ':focus', which GtkCalendar's own GTK_STATE_FLAG_FOCUSED sets) shows once a click has given the calendar the focus, 2 pixels of $focus_border_color inside the chosen day's accent; and the day numbers are in tabular digits (font-variant-numeric: tabular-nums), as wide as a text of the same digits so written.
func TestAdwCalendarAsStyled(t *testing.T) {
	fixToday(t, datePickerToday)
	day := time.Date(2026, time.November, 18, 0, 0, 0, 0, time.UTC)
	var ref ui.Rect
	tester := dateTester(Linux, 640, 520, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(16).AlignItems(ui.Start).Children(func() {
			CalendarView(c, &day)
			ref = ui.Text(c, "11").SingleLine().FontFeatures("tnum").Bounds()
		})
	})
	tester.SetScale(2)
	tester.Frame()
	if err := tester.Click(dayName(time.November, 20)); err != nil {
		t.Fatal(err)
	}
	tester.Frame()
	tester.Frame()
	label, _ := tester.Find(dayName(time.November, 20))
	accent := shotPixel(tester, label.X+label.W/2, label.Y+4)
	if ring := shotPixel(tester, label.X+1, label.Y+label.H/2); samePixel(ring, accent, 10) {
		t.Errorf("no focus ring after a click: %v at the label's edge, its accent %v", ring, accent)
	}
	num, _ := tester.Find("11")
	if num.W < ref.W-0.01 || num.W > ref.W+0.01 {
		t.Errorf("the 11th is %v wide, want %v, in tabular digits", num.W, ref.W)
	}
}
