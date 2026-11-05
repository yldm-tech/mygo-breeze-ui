package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, valuesStories) }

// valuesStories are the stories of the widgets of widgets_values.go: sliders with steps and ranges, dates, times, decimals and colors. Their dates, times and colors are text controls, which the widgets write back as they change them.
func valuesStories() []*story {
	steps := []string{"10", "20", "25", "50"}
	return []*story{
		// ---- Choices ----
		{component: "StepSlider", name: "Default", category: "Choices",
			controls: []control{numberArg("value", 50, 0, 100), choiceArg("step", "25", steps...), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				step, _ := strconv.ParseFloat(*a.text["step"], 64)
				disabled(a, breeze.StepSlider(c, a.num["value"], 0, 100, step).Width(240).Label("Quality"))
			},
			code: func(a *args) string {
				return "breeze.StepSlider(c, &app.quality, 0, 100, " + *a.text["step"] + ").Label(\"Quality\")"
			}},
		{component: "RangeSlider", name: "Price", category: "Choices",
			controls: []control{numberArg("low", 20, 0, 100), numberArg("high", 70, 0, 100), choiceArg("step", "0", "0", "10", "25")},
			render: func(c *ui.Context, a *args) {
				step, _ := strconv.ParseFloat(*a.text["step"], 64)
				breeze.RangeSlider(c, a.num["low"], a.num["high"], 0, 100, step).Width(240).Label("Price")
			},
			code: func(a *args) string {
				return "breeze.RangeSlider(c, &app.min, &app.max, 0, 100, " + *a.text["step"] + ").Label(\"Price\")"
			}},
		{component: "ColorWell", name: "Default", category: "Choices",
			controls: []control{textArg("color", "#3584e4"), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				color := valuesColor(*a.text["color"])
				if disabled(a, breeze.ColorWell(c, &color).Label("Tint")).Changed() {
					*a.text["color"] = valuesHex(color)
				}
			},
			code: func(*args) string { return "breeze.ColorWell(c, &app.tint).Label(\"Tint\")" }},
		{component: "ColorPicker", name: "Default", category: "Choices",
			controls: []control{textArg("color", "#3584e4")},
			render: func(c *ui.Context, a *args) {
				color := valuesColor(*a.text["color"])
				if breeze.ColorPicker(c, &color).Width(260).Changed() {
					*a.text["color"] = valuesHex(color)
				}
			},
			code: func(*args) string { return "breeze.ColorPicker(c, &app.tint)" }},
		{component: "Calendar", name: "Default", category: "Choices",
			controls: []control{textArg("date", "2026-01-05")},
			render: func(c *ui.Context, a *args) {
				day := valuesDate(*a.text["date"])
				if breeze.Calendar(c, &day).Label("Day").Changed() {
					*a.text["date"] = day.Format(time.DateOnly)
				}
			},
			code: func(*args) string { return "breeze.Calendar(c, &app.day)" }},

		// ---- Text ----
		{component: "DateInput", name: "Default", category: "Text",
			controls: []control{textArg("date", "2026-01-05"), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				due := valuesDate(*a.text["date"])
				if disabled(a, breeze.DateInput(c, &due).Label("Due")).Changed() {
					*a.text["date"] = due.Format(time.DateOnly)
				}
			},
			code: func(*args) string { return "breeze.DateInput(c, &app.due).Label(\"Due\")" }},
		{component: "TimeInput", name: "Default", category: "Text",
			controls: []control{textArg("time", "17:45"), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				alarm := valuesTime(*a.text["time"])
				if disabled(a, breeze.TimeInput(c, &alarm).Label("Alarm")).Changed() {
					*a.text["time"] = alarm.Format("15:04")
				}
			},
			code: func(*args) string { return "breeze.TimeInput(c, &app.alarm).Label(\"Alarm\")" }},
		{component: "DecimalInput", name: "Amount", category: "Text",
			controls: []control{textArg("value", "12.50"), choiceArg("step", "0.01", "0.01", "0.05", "1")},
			render: func(c *ui.Context, a *args) {
				breeze.DecimalInput(c, a.text["value"], "0", "", *a.text["step"]).Label("Amount")
			},
			code: func(a *args) string {
				return "breeze.DecimalInput(c, &app.amount, \"0\", \"\", " + q(*a.text["step"]) + ").Label(\"Amount\")"
			}},
	}
}

// valuesDate reads a date control written as 2006-01-02, or the first of January 2026 when it is not one.
func valuesDate(s string) time.Time {
	d, err := time.Parse(time.DateOnly, strings.TrimSpace(s))
	if err != nil {
		return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	return d
}

// valuesTime reads a time control written as 15:04, or noon when it is not one.
func valuesTime(s string) time.Time {
	tm, err := time.Parse("15:04", strings.TrimSpace(s))
	if err != nil {
		return time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	}
	return time.Date(2026, 1, 1, tm.Hour(), tm.Minute(), 0, 0, time.UTC)
}

// valuesColor reads a color control written as #rrggbb or #rrggbbaa, or gray when it is not one: ui.Hex panics on text that is not a color, which a control being typed in often is.
func valuesColor(s string) ui.Color {
	h := strings.TrimPrefix(strings.TrimSpace(s), "#")
	v, err := strconv.ParseUint(h, 16, 32)
	switch {
	case err != nil:
	case len(h) == 6:
		return ui.Color{R: uint8(v >> 16), G: uint8(v >> 8), B: uint8(v), A: 255}
	case len(h) == 8:
		return ui.Color{R: uint8(v >> 24), G: uint8(v >> 16), B: uint8(v >> 8), A: uint8(v)}
	}
	return ui.RGB(128, 128, 128)
}

// valuesHex writes a color as the color control reads it.
func valuesHex(c ui.Color) string {
	if c.A == 255 {
		return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
	}
	return fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, c.A)
}
