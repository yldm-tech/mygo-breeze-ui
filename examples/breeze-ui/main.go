// Command breeze-ui is the Breeze UI gallery: the same view drawn with the controls of macOS, Windows 11 and GNOME side by side, or with the controls of the system it runs on.
package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

type gallery struct {
	mode     int // 0 side by side, 1 this system
	remember bool
	sync     bool
	name     string
	size     string
	volume   float64
	period   int
	saved    int
	density  string
	bold     bool
	query    string
	notes    string
	tab      int
}

func (g *gallery) view(c *ui.Context) {
	breeze.Use(c)
	ui.Column(c).Fill().Children(func() {
		ui.Row(c).Padding(12, 16).Gap(12).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "Breeze UI").FontSize(18).Bold().Grow(1)
			ui.Segmented(c, &g.mode, "Side by side", "This system")
		})
		if g.mode == 1 {
			ui.Scroll(c).Fill().Padding(16, 24).Children(func() { g.controls(c) })
			return
		}
		ui.Row(c).Fill().AlignItems(ui.Stretch).Children(func() {
			for _, p := range breeze.Platforms {
				breeze.Preview(c, p, func() {
					ui.Themed(c, breeze.Theme(c, p)).Grow(1).Basis(0).Padding(16, 20).Gap(14).Children(func() {
						ui.Text(c, p.String()).FontSize(15).Bold()
						g.controls(c)
					})
				})
			}
		})
	})
}

// controls builds the same controls whichever platform draws them; the state is shared, so a change in one column shows in the others.
func (g *gallery) controls(c *ui.Context) {
	ui.Column(c).Gap(12).Children(func() {
		breeze.Tabs(c, &g.tab, "General", "Advanced")
		ui.Row(c).Gap(8).Wrap().Children(func() {
			breeze.Button(c, "Cancel")
			if breeze.PrimaryButton(c, "Save").Clicked() {
				g.saved++
			}
			breeze.Toggle(c, &g.bold, "Bold")
			breeze.Button(c, "Disabled").Disabled(true)
		})
		breeze.Checkbox(c, &g.remember, "Remember me")
		breeze.RadioGroup(c, func() {
			ui.Row(c).Gap(12).Children(func() {
				breeze.Radio(c, &g.density, "compact", "Compact")
				breeze.Radio(c, &g.density, "roomy", "Roomy")
			})
		})
		ui.Row(c).Gap(10).AlignItems(ui.Center).Children(func() {
			breeze.Switch(c, &g.sync).Label("Sync")
			ui.Text(c, "Sync")
			ui.Spacer(c)
			breeze.Spinner(c)
		})
		breeze.TextInput(c, &g.name).Placeholder("Name")
		breeze.SearchField(c, &g.query)
		breeze.Select(c, &g.size, []string{"Small", "Medium", "Large"})
		breeze.Slider(c, &g.volume, 0, 100)
		breeze.Progress(c, g.volume/100)
		breeze.Segmented(c, &g.period, "Day", "Week", "Month")
		breeze.TextArea(c, &g.notes).Placeholder("Notes").MinHeight(64)
	})
}

func main() {
	g := &gallery{remember: true, size: "Medium", volume: 40, density: "compact", bold: true}
	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title:   "Breeze UI",
			Width:   1080,
			Height:  760,
			Content: ui.View(g.view),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
