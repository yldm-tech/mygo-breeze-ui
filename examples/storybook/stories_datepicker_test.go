package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestDatePickerStoriesChooseDays checks on every platform that the CalendarView story writes the day a click chooses into its date control, and that the DatePicker story shows its date, opens its calendar and writes the day chosen there.
func TestDatePickerStoriesChooseDays(t *testing.T) {
	for platform := range 3 {
		b := newBook()
		b.selected, b.platform = "CalendarView/Default", platform
		tester := ui.NewTester(b.view, 1280, 820)
		tester.SetLocale("en-US")
		a := b.argsOf(b.story())
		tester.Frame()
		if err := tester.Click("November 20, 2026"); err != nil {
			t.Fatalf("%s: %v", platformNames[platform], err)
		}
		tester.Frame()
		tester.Frame()
		if got := *a.text["date"]; got != "2026-11-20" {
			t.Errorf("%s: the calendar's click left the date control at %q", platformNames[platform], got)
		}
		if !tester.HasText("breeze.CalendarView(c, &app.day)") {
			t.Errorf("%s: the calendar story shows no code", platformNames[platform])
		}

		b = newBook()
		b.selected, b.platform = "DatePicker/Default", platform
		tester = ui.NewTester(b.view, 1280, 820)
		tester.SetLocale("en-US")
		a = b.argsOf(b.story())
		tester.Frame()
		if err := tester.Click("11/18/2026"); err != nil {
			t.Fatalf("%s: %v", platformNames[platform], err)
		}
		tester.Frame()
		tester.Frame()
		if err := tester.Click("November 3, 2026"); err != nil {
			t.Fatalf("%s: %v", platformNames[platform], err)
		}
		tester.Frame()
		tester.Frame()
		if got := *a.text["date"]; got != "2026-11-03" || !tester.HasText("11/3/2026") {
			t.Errorf("%s: choosing the 3rd left the date control at %q", platformNames[platform], got)
		}
	}
}
