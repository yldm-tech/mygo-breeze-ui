package main

import (
	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, settingsStories) }

// settingsStories are the stories of the widgets of widgets_settings.go: combo rows, entry rows, spin rows and expander rows.
func settingsStories() []*story {
	appearances := []string{"Light", "Dark", "Auto"}
	return []*story{
		// ---- Choices ----
		{component: "ComboRow", name: "Default", category: "Choices",
			controls: []control{textArg("title", "Appearance"), textArg("subtitle", "Light or dark"), choiceArg("selected", "Auto", appearances...), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.ComboRow(c, *a.text["title"], *a.text["subtitle"], a.text["selected"], appearances).FillWidth().MaxWidth(360))
			},
			code: func(a *args) string {
				return "breeze.ComboRow(c, " + q(*a.text["title"]) + ", " + rowsSubtitle(a) + ", &app.appearance, []string{\"Light\", \"Dark\", \"Auto\"})"
			}},

		// ---- Text ----
		{component: "EntryRow", name: "Default", category: "Text",
			controls: []control{textArg("title", "Name"), textArg("value", ""), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.EntryRow(c, *a.text["title"], a.text["value"]).FillWidth().MaxWidth(360))
			},
			code: func(a *args) string {
				return "breeze.EntryRow(c, " + q(*a.text["title"]) + ", &app.name)"
			}},
		{component: "SpinRow", name: "Default", category: "Text",
			controls: []control{textArg("title", "Font size"), textArg("subtitle", "Points"), numberArg("value", 12, 6, 72), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.SpinRow(c, *a.text["title"], *a.text["subtitle"], a.num["value"], 6, 72, 1).FillWidth().MaxWidth(360))
			},
			code: func(a *args) string {
				return "breeze.SpinRow(c, " + q(*a.text["title"]) + ", " + rowsSubtitle(a) + ", &app.size, 6, 72, 1)"
			}},

		// ---- Layout ----
		{component: "ExpanderRow", name: "Default", category: "Layout",
			controls: []control{textArg("title", "Network"), textArg("subtitle", "Wi-Fi and proxies"), toggleArg("open", true), toggleArg("metered", false)},
			render: func(c *ui.Context, a *args) {
				breeze.ExpanderRow(c, *a.text["title"], *a.text["subtitle"], a.on["open"], func() {
					breeze.ActionRow(c, "Proxy", "Automatic", nil)
					breeze.SwitchRow(c, "Metered", "", a.on["metered"])
				}).FillWidth().MaxWidth(360)
			},
			code: func(a *args) string {
				return "breeze.ExpanderRow(c, " + q(*a.text["title"]) + ", " + rowsSubtitle(a) + ", &app.network, func() {\n\tbreeze.ActionRow(c, \"Proxy\", \"Automatic\", nil)\n\tbreeze.SwitchRow(c, \"Metered\", \"\", &app.metered)\n})"
			}},
		{component: "ExpanderRow", name: "In a card", category: "Layout",
			controls: []control{choiceArg("appearance", "Auto", appearances...), textArg("name", "Ada"), numberArg("size", 12, 6, 72), toggleArg("open", false)},
			render: func(c *ui.Context, a *args) {
				breeze.Card(c, func() {
					breeze.ComboRow(c, "Appearance", "", a.text["appearance"], appearances)
					breeze.EntryRow(c, "Name", a.text["name"])
					breeze.SpinRow(c, "Font size", "", a.num["size"], 6, 72, 1)
					breeze.ExpanderRow(c, "Advanced", "", a.on["open"], func() {
						breeze.ActionRow(c, "Logs", "", nil)
					})
				}).FillWidth().MaxWidth(360)
			},
			code: func(*args) string {
				return "breeze.Card(c, func() {\n\tbreeze.ComboRow(c, \"Appearance\", \"\", &app.appearance, []string{\"Light\", \"Dark\", \"Auto\"})\n\tbreeze.EntryRow(c, \"Name\", &app.name)\n\tbreeze.SpinRow(c, \"Font size\", \"\", &app.size, 6, 72, 1)\n\tbreeze.ExpanderRow(c, \"Advanced\", \"\", &app.advanced, func() {\n\t\tbreeze.ActionRow(c, \"Logs\", \"\", nil)\n\t})\n})"
			}},
	}
}
