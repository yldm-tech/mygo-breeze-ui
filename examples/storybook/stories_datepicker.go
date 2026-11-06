package main

import (
	"time"

	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, datePickerStories) }

// datePickerStories are the stories of the widgets of widgets_datepicker.go: a month's calendar drawn per platform and the button opening it. Their date control, written as 2006-01-02, holds the day they choose.
func datePickerStories() []*story {
	return []*story{
		{component: "CalendarView", name: "Default", category: "Choices",
			controls: []control{textArg("date", "2026-11-18"), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				day := valuesDate(*a.text["date"])
				if disabled(a, breeze.CalendarView(c, &day)).Changed() {
					*a.text["date"] = day.Format(time.DateOnly)
				}
			},
			code: func(*args) string { return "breeze.CalendarView(c, &app.day)" }},
		{component: "DatePicker", name: "Default", category: "Text",
			controls: []control{textArg("date", "2026-11-18"), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				due := valuesDate(*a.text["date"])
				if disabled(a, breeze.DatePicker(c, &due).Label("Due")).Changed() {
					*a.text["date"] = due.Format(time.DateOnly)
				}
			},
			code: func(*args) string { return "breeze.DatePicker(c, &app.due).Label(\"Due\")" }},
	}
}
