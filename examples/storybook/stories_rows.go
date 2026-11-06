package main

import (
	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, rowsStories) }

// rowsSubtitle returns the Go of a row's subtitle argument.
func rowsSubtitle(a *args) string { return q(*a.text["subtitle"]) }

// rowsStories are the stories of the widgets of widgets_rows.go: cards, action rows and switch rows.
func rowsStories() []*story {
	return []*story{
		// ---- Choices ----
		{component: "SwitchRow", name: "Default", category: "Choices",
			controls: []control{textArg("title", "Dark mode"), textArg("subtitle", "Use a dark appearance"), toggleArg("on", true), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.SwitchRow(c, *a.text["title"], *a.text["subtitle"], a.on["on"]).FillWidth().MaxWidth(360))
			},
			code: func(a *args) string {
				return "breeze.SwitchRow(c, " + q(*a.text["title"]) + ", " + rowsSubtitle(a) + ", &app.dark)"
			}},

		// ---- Layout ----
		{component: "Card", name: "Text", category: "Layout",
			controls: []control{textArg("text", "A card holds what belongs together.")},
			render: func(c *ui.Context, a *args) {
				breeze.Card(c, func() {
					ui.Text(c, *a.text["text"])
				}).FillWidth().MaxWidth(360)
			},
			code: func(a *args) string {
				return "breeze.Card(c, func() {\n\tui.Text(c, " + q(*a.text["text"]) + ")\n})"
			}},
		{component: "Card", name: "Settings", category: "Layout",
			controls: []control{toggleArg("dark", true), toggleArg("notifications", false)},
			render: func(c *ui.Context, a *args) {
				breeze.Card(c, func() {
					breeze.ActionRow(c, "Wallpaper", "Choose a picture", func() { breeze.Button(c, "Open") })
					breeze.SwitchRow(c, "Dark mode", "Use a dark appearance", a.on["dark"])
					breeze.SwitchRow(c, "Notifications", "", a.on["notifications"])
				}).FillWidth().MaxWidth(360)
			},
			code: func(*args) string {
				return "breeze.Card(c, func() {\n\tbreeze.ActionRow(c, \"Wallpaper\", \"Choose a picture\", func() { breeze.Button(c, \"Open\") })\n\tbreeze.SwitchRow(c, \"Dark mode\", \"Use a dark appearance\", &app.dark)\n\tbreeze.SwitchRow(c, \"Notifications\", \"\", &app.notifications)\n})"
			}},
		{component: "ActionRow", name: "Default", category: "Layout",
			controls: []control{textArg("title", "Wallpaper"), textArg("subtitle", "Choose a picture"), textArg("button", "Open"), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.ActionRow(c, *a.text["title"], *a.text["subtitle"], func() {
					if *a.text["button"] != "" {
						breeze.Button(c, *a.text["button"])
					}
				}).FillWidth().MaxWidth(360))
			},
			code: func(a *args) string {
				if *a.text["button"] == "" {
					return "breeze.ActionRow(c, " + q(*a.text["title"]) + ", " + rowsSubtitle(a) + ", nil)"
				}
				return "breeze.ActionRow(c, " + q(*a.text["title"]) + ", " + rowsSubtitle(a) + ", func() {\n\tbreeze.Button(c, " + q(*a.text["button"]) + ")\n})"
			}},
		{component: "ActionRow", name: "In a fieldset", category: "Layout",
			controls: []control{textArg("legend", "Privacy"), toggleArg("location", true)},
			render: func(c *ui.Context, a *args) {
				breeze.Fieldset(c, *a.text["legend"], func() {
					breeze.SwitchRow(c, "Location", "Apps may ask for it", a.on["location"])
					breeze.ActionRow(c, "Camera", "2 apps", nil)
				}).FillWidth().MaxWidth(360)
			},
			code: func(a *args) string {
				return "breeze.Fieldset(c, " + q(*a.text["legend"]) + ", func() {\n\tbreeze.SwitchRow(c, \"Location\", \"Apps may ask for it\", &app.location)\n\tbreeze.ActionRow(c, \"Camera\", \"2 apps\", nil)\n})"
			}},
	}
}
