package breeze

import (
	"fmt"
	"math"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

func init() {
	for _, n := range []string{"StepSlider", "DateInput", "TimeInput", "ColorWell"} {
		drawnWidgets[n] = true
	}
}

func TestStepTicks(t *testing.T) {
	for _, c := range []struct {
		lo, hi, step float64
		want         []float32
	}{
		{0, 100, 25, []float32{0, 0.25, 0.5, 0.75, 1}},
		{0, 10, 3, []float32{0, 0.3, 0.6, 0.9}},
		{0, 100, 0, nil},
		{0, 100, 1, nil}, // 101 ticks, too many to tell apart
		{5, 5, 1, nil},
	} {
		got := stepTicks(c.lo, c.hi, c.step)
		if len(got) != len(c.want) {
			t.Fatalf("stepTicks(%v, %v, %v) = %v, want %v", c.lo, c.hi, c.step, got, c.want)
		}
		for i := range got {
			if math.Abs(float64(got[i]-c.want[i])) > 1e-6 {
				t.Fatalf("stepTicks(%v, %v, %v) = %v, want %v", c.lo, c.hi, c.step, got, c.want)
			}
		}
	}
}

// TestStepSliderKeepsItsBehavior checks on every platform, left to right and right to left, that a StepSlider snaps a click to a step, that the arrows move it by step, and that a disabled one does not move.
func TestStepSliderKeepsItsBehavior(t *testing.T) {
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			quality, off := 0.0, 50.0
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					ui.Column(c).Gap(16).Padding(16).Children(func() {
						StepSlider(c, &quality, 0, 100, 25).Width(200).Label("Quality")
						StepSlider(c, &off, 0, 100, 25).Width(200).Label("Off").Disabled(true)
					})
				})
			}, 300, 200)
			tester.SetDirection(dir)
			tester.Frame()
			r, ok := tester.Find("Quality")
			if !ok {
				t.Fatalf("%v: slider named Quality not found", p)
			}
			// A third of the way from lo, which snaps to the step nearest it.
			x := r.X + r.W*0.36
			if dir == ui.RightToLeft {
				x = r.X + r.W*0.64
			}
			tester.ClickAt(x, r.Y+r.H/2)
			tester.Frame()
			if quality != 25 {
				t.Fatalf("%v %v: a click a third along set %v, want 25", p, dir, quality)
			}
			forth := ui.KeyRight
			if dir == ui.RightToLeft {
				forth = ui.KeyLeft
			}
			tester.Key(0, forth)
			tester.Frame()
			if quality != 50 {
				t.Fatalf("%v %v: the arrow forth set %v, want 50", p, dir, quality)
			}
			end := r.X + r.W - 1
			if dir == ui.RightToLeft {
				end = r.X + 1
			}
			tester.ClickAt(end, r.Y+r.H/2)
			tester.Frame()
			if quality != 100 {
				t.Fatalf("%v %v: a click at the end set %v, want 100", p, dir, quality)
			}
			d, _ := tester.Find("Off")
			tester.ClickAt(d.X+d.W-1, d.Y+d.H/2)
			tester.Frame()
			if off != 50 {
				t.Fatalf("%v %v: a disabled slider moved to %v", p, dir, off)
			}
		}
	}
}

// TestDateAndTimeInputsKeepTheirBehavior checks on every platform that the date field opens MyGo's calendar and chooses a day from it with the pointer and the keys, and that the time field steps its hours with the arrows.
func TestDateAndTimeInputsKeepTheirBehavior(t *testing.T) {
	for _, p := range Platforms {
		due := time.Date(2026, 1, 5, 9, 30, 0, 0, time.UTC)
		alarm := time.Date(2026, 1, 5, 7, 30, 0, 0, time.UTC)
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				ui.Column(c).Gap(16).Padding(16).AlignItems(ui.Start).Children(func() {
					TimeInput(c, &alarm).Label("Alarm")
					DateInput(c, &due).Label("Due")
				})
			})
		}, 480, 520)
		tester.SetRegion(i18n.Region{Locale: "en-GB", FirstWeekday: time.Monday, HourCycle: 24})
		tester.Frame()
		if err := tester.Click("Due"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		if !tester.HasText("January 2026") {
			t.Fatalf("%v: no calendar after a click: %q", p, tester.Texts())
		}
		if err := tester.Click("15 January 2026"); err != nil {
			t.Fatalf("%v: %v (%q)", p, err, tester.Texts())
		}
		tester.Frame()
		if want := time.Date(2026, 1, 15, 9, 30, 0, 0, time.UTC); !due.Equal(want) || tester.HasText("January 2026") {
			t.Fatalf("%v: a click on the 15th: %v, calendar shown %v", p, due, tester.HasText("January 2026"))
		}
		// From the keyboard: Enter opens it, Right moves, Enter chooses.
		tester.Key(0, ui.KeyEnter)
		tester.Frame()
		tester.Key(0, ui.KeyRight)
		tester.Frame()
		tester.Key(0, ui.KeyEnter)
		tester.Frame()
		if due.Day() != 16 {
			t.Fatalf("%v: Right and Enter chose %v", p, due)
		}
		if err := tester.Click("07"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		tester.Key(0, ui.KeyUp)
		tester.Frame()
		if alarm.Hour() != 8 || alarm.Minute() != 30 {
			t.Fatalf("%v: Up on the hours set %v", p, alarm)
		}
	}
}

// TestColorWellKeepsItsBehavior checks on every platform that the restyled color well opens MyGo's picker, chooses a swatch and closes with Escape, and that a disabled one does not open.
func TestColorWellKeepsItsBehavior(t *testing.T) {
	for _, p := range Platforms {
		tint, other := ui.Hex("#3b82f6"), ui.Hex("#ffffff")
		changes := 0
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				ui.Row(c).Gap(16).Padding(16).Children(func() {
					if ColorWell(c, &tint).Label("Tint").Changed() {
						changes++
					}
					ColorWell(c, &other).Label("Other").Disabled(true)
				})
			})
		}, 480, 520)
		tester.Frame()
		if err := tester.Click("Other"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		if tester.HasText("Red") {
			t.Fatalf("%v: a disabled well opened its picker", p)
		}
		if err := tester.Click("Tint"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		if !tester.HasText("Red") {
			t.Fatalf("%v: the picker did not open: %q", p, tester.Texts())
		}
		if err := tester.Click("Green"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		if tint != ui.Hex("#22c55e") || changes != 1 {
			t.Fatalf("%v: chose %v with %d changes", p, tint, changes)
		}
		tester.Key(0, ui.KeyEscape)
		tester.Frame()
		if tester.HasText("Red") {
			t.Fatalf("%v: Escape left the picker open", p)
		}
	}
}

// TestValueWidgetsOfPackageUI checks the widgets of this group that are package ui's in the platform's theme: the range slider's knobs move with the keys, the decimal input steps with its button, the calendar chooses a day, and the picker a swatch.
func TestValueWidgetsOfPackageUI(t *testing.T) {
	for _, p := range Platforms {
		low, high := 20.0, 60.0
		amount := "1.50"
		day := time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC)
		color := ui.Hex("#3b82f6")
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Column(c).Gap(16).Padding(16).AlignItems(ui.Start).Children(func() {
					RangeSlider(c, &low, &high, 0, 100, 10).Width(220).Label("Price")
					DecimalInput(c, &amount, "0", "", "0.25").Label("Amount")
					Calendar(c, &day).Label("Day")
					ColorPicker(c, &color)
				})
			})
		}, 480, 1000)
		tester.Frame()
		tester.Key(0, ui.KeyTab)
		tester.Frame()
		tester.Key(0, ui.KeyRight)
		tester.Frame()
		if low != 30 || high != 60 {
			t.Fatalf("%v: Tab and Right set the range to %v..%v", p, low, high)
		}
		if err := tester.Click("Increase"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		if amount != "1.75" {
			t.Fatalf("%v: the decimal input stepped to %q", p, amount)
		}
		if err := tester.Click("January 20, 2026"); err != nil {
			t.Fatalf("%v: %v (%q)", p, err, tester.Texts())
		}
		tester.Frame()
		if day.Day() != 20 {
			t.Fatalf("%v: the calendar chose %v", p, day)
		}
		if err := tester.Click("Green"); err != nil {
			t.Fatalf("%v: %v", p, err)
		}
		tester.Frame()
		if color != ui.Hex("#22c55e") {
			t.Fatalf("%v: the picker chose %v", p, color)
		}
	}
}

// valuesView lays out the widgets of this group beside the platform's button, text field and slider, to compare their heights, corners and paddings with them: sliders with and without steps, a range slider, date, time and decimal fields, color wells of a dark, a light and a translucent color, and a calendar.
func valuesView(p Platform, s *valuesState) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Row(c).Fill().Padding(16).Gap(24).AlignItems(ui.Start).Children(func() {
				ui.Column(c).Gap(14).Width(260).Children(func() {
					ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
						Button(c, "Button")
						TextInput(c, &s.name).Placeholder("Text").Width(120)
					})
					Slider(c, &s.volume, 0, 100).Label("Volume")
					StepSlider(c, &s.quality, 0, 100, 25).Label("Quality")
					StepSlider(c, &s.off, 0, 10, 1).Label("Disabled").Disabled(true)
					RangeSlider(c, &s.low, &s.high, 0, 100, 0).Label("Price")
					ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
						DateInput(c, &s.due).Label("Due")
						TimeInput(c, &s.alarm).Label("Alarm")
					})
					DecimalInput(c, &s.amount, "0", "", "0.01").Label("Amount")
					ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
						ColorWell(c, &s.tint).Label("Tint")
						ColorWell(c, &s.paper).Label("Paper")
						ColorWell(c, &s.glass).Label("Glass")
						ColorWell(c, &s.tint).Label("Off").Disabled(true)
					})
				})
				Calendar(c, &s.due).Label("Calendar")
			})
		})
	}
}

// valuesState is what valuesView's widgets edit.
type valuesState struct {
	name                            string
	volume, quality, off, low, high float64
	due, alarm                      time.Time
	amount                          string
	tint, paper, glass              ui.Color
}

func newValuesState() *valuesState {
	return &valuesState{
		volume: 40, quality: 50, off: 3, low: 20, high: 70,
		due:    time.Date(2026, 1, 5, 9, 30, 0, 0, time.UTC),
		alarm:  time.Date(2026, 1, 5, 17, 45, 0, 0, time.UTC),
		amount: "12.50",
		tint:   ui.Hex("#3584e4"), paper: ui.Hex("#ffffff"), glass: ui.RGBA(255, 59, 48, 0.4),
	}
}

// TestScreenshotsValues writes valuesView on every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names, with the date field's calendar and a color well's picker open in two more. Without it the test only builds the frames.
func TestScreenshotsValues(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			name := strings.ToLower(p.String())
			shot := func(file string, open string) {
				tester := ui.NewTester(valuesView(p, newValuesState()), 640, 460)
				tester.SetScale(2)
				tester.SetDark(dark)
				tester.Frame()
				if open != "" {
					if err := tester.Click(open); err != nil {
						t.Fatalf("%v: %v", p, err)
					}
					tester.Frame()
				}
				tester.Frame()
				for _, s := range []string{"Button", "January 2026"} {
					if !tester.HasText(s) {
						t.Fatalf("%v: %q is missing", p, s)
					}
				}
				writeShot(t, dir, fmt.Sprintf(file, name, mode), tester)
			}
			shot("values-%s-%s.png", "")
			shot("values-date-%s-%s.png", "Due")
			shot("values-color-%s-%s.png", "Tint")
		}
	}
}
