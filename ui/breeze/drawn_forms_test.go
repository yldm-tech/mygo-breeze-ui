package breeze

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// formsShotState is what formsShotView's widgets edit.
type formsShotState struct {
	mail, calendar, photos, all1, all2 bool
	copies, full                       float64
	name, address                      string
	alarm, due                         time.Time
	split                              float32
}

func newFormsShotState() *formsShotState {
	return &formsShotState{
		mail: true, all1: true, all2: true,
		copies: 5, full: 10,
		name:  "Ada",
		alarm: time.Date(2026, 1, 5, 15, 4, 0, 0, time.UTC), due: time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
		split: 120,
	}
}

// formsShotView lays out the widgets of the forms group for screenshots: a mixed and a checked check box group, two number inputs (one at its maximum), a field outside a form, a fieldset, a time and a date input, and a split.
func formsShotView(p Platform, s *formsShotState) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Row(c).Fill().Padding(16).Gap(24).AlignItems(ui.Start).Children(func() {
				ui.Column(c).Width(260).Gap(16).Children(func() {
					CheckboxGroup(c, "Notifications", func() {
						Checkbox(c, &s.mail, "Mail")
						Checkbox(c, &s.calendar, "Calendar")
						Checkbox(c, &s.photos, "Photos")
					})
					CheckboxGroup(c, "Sync", func() {
						Checkbox(c, &s.all1, "Contacts")
						Checkbox(c, &s.all2, "Notes")
					})
					NumberInput(c, &s.copies, 0, 10, 1)
					NumberInput(c, &s.full, 0, 10, 1)
				})
				ui.Column(c).Width(300).Gap(16).Children(func() {
					Field(c, "Name", func() { TextInput(c, &s.name) })
					Fieldset(c, "Shipping", func() {
						Field(c, "Address", func() { TextInput(c, &s.address).Placeholder("Street") })
					})
					ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
						TimeInput(c, &s.alarm).Label("Alarm")
						DateInput(c, &s.due).Label("Due")
					})
					Split(c, &s.split, func() {
						ui.Text(c, "Sidebar").Padding(8)
					}, func() {
						ui.Text(c, "Content").Padding(8)
					}).Height(60).Border(1, c.Theme().Border)
				})
			})
		})
	}
}

// TestScreenshotsWave2Forms writes the forms group of every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names, with the pointer over the split's divider. Without it the test only builds the frames.
func TestScreenshotsWave2Forms(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			s := newFormsShotState()
			tester := ui.NewTester(formsShotView(p, s), 640, 440)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tester.SetRegion(i18n.Region{Locale: "en-US", FirstWeekday: time.Sunday, HourCycle: 24})
			tester.Frame()
			tester.Frame()
			for _, want := range []string{"Notifications", "Mail", "Shipping", "Name", "Sidebar"} {
				if !tester.HasText(want) {
					t.Fatalf("%v: %q is missing", p, want)
				}
			}
			side, ok := tester.Find("Sidebar")
			if !ok {
				t.Fatalf("%v: no sidebar pane", p)
			}
			// The divider is s.split DIPs after the first pane's start, where its text starts.
			tester.Move(side.X+s.split+0.5, side.Y+side.H/2)
			tester.Frame()
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			writeShot(t, dir, fmt.Sprintf("breeze-forms-%s-%s.png", strings.ToLower(p.String()), mode), tester)
		}
	}
}
