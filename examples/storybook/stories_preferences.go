package main

import (
	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, preferencesStories) }

// preferencesTile builds the tile the preferences stories put before a row's title: a rounded square of systemBlue, 30 points, 20 DIPs on Windows, where SettingsCard's HeaderIcon is 20.
func preferencesTile(c *ui.Context) {
	s := float32(30)
	if breeze.PlatformOf(c) == breeze.Windows {
		s = 20
	}
	ui.Box(c).Size(s, s).Shrink(0).Radius(s * 7 / 30).Background(ui.Hex("#0088ff"))
}

// preferencesStories are the stories of the widgets of widgets_preferences.go: preferences groups, prefix rows, content rows and radio rows.
func preferencesStories() []*story {
	engines := []string{"whisper", "cloud", "off"}
	return []*story{
		// ---- Choices ----
		{component: "RadioRow", name: "In a group", category: "Choices",
			controls: []control{choiceArg("engine", "whisper", engines...), toggleArg("cloud disabled", false)},
			render: func(c *ui.Context, a *args) {
				breeze.PreferencesGroup(c, "Engine", "", nil, func() {
					breeze.RadioRow(c, a.text["engine"], "whisper", "Whisper", "On this Mac", nil)
					breeze.RadioRow(c, a.text["engine"], "cloud", "Cloud", "Sent to a service", nil).Disabled(*a.on["cloud disabled"])
					breeze.RadioRow(c, a.text["engine"], "off", "Off", "", nil)
				}).FillWidth().MaxWidth(360)
			},
			code: func(a *args) string {
				return "breeze.PreferencesGroup(c, \"Engine\", \"\", nil, func() {\n\tbreeze.RadioRow(c, &app.engine, \"whisper\", \"Whisper\", \"On this Mac\", nil)\n\tbreeze.RadioRow(c, &app.engine, \"cloud\", \"Cloud\", \"Sent to a service\", nil)\n\tbreeze.RadioRow(c, &app.engine, \"off\", \"Off\", \"\", nil)\n})"
			}},

		// ---- Layout ----
		{component: "PreferencesGroup", name: "Default", category: "Layout",
			controls: []control{textArg("title", "Speech"), textArg("description", "Models are stored on this Mac."), textArg("button", "Detect"), toggleArg("auto", true)},
			render: func(c *ui.Context, a *args) {
				var suffix func()
				if *a.text["button"] != "" {
					suffix = func() { breeze.Button(c, *a.text["button"]) }
				}
				breeze.PreferencesGroup(c, *a.text["title"], *a.text["description"], suffix, func() {
					breeze.ActionRow(c, "Language", "", func() { ui.Text(c, "Auto").TextColor(c.Theme().TextMuted) })
					breeze.SwitchRow(c, "Detect on launch", "", a.on["auto"])
				}).FillWidth().MaxWidth(360)
			},
			code: func(a *args) string {
				suffix := "nil"
				if *a.text["button"] != "" {
					suffix = "func() {\n\tbreeze.Button(c, " + q(*a.text["button"]) + ")\n}"
				}
				return "breeze.PreferencesGroup(c, " + q(*a.text["title"]) + ", " + q(*a.text["description"]) + ", " + suffix + ", func() {\n\tbreeze.ActionRow(c, \"Language\", \"\", func() { ui.Text(c, \"Auto\") })\n\tbreeze.SwitchRow(c, \"Detect on launch\", \"\", &app.auto)\n})"
			}},
		{component: "PrefixRow", name: "Default", category: "Layout",
			controls: []control{textArg("title", "Whisper"), textArg("subtitle", "Base · 142 MB"), textArg("button", "Remove"), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				var suffix func()
				if *a.text["button"] != "" {
					suffix = func() { breeze.Button(c, *a.text["button"]) }
				}
				disabled(a, breeze.PrefixRow(c, *a.text["title"], *a.text["subtitle"], func() { preferencesTile(c) }, suffix).FillWidth().MaxWidth(360))
			},
			code: func(a *args) string {
				suffix := "nil"
				if *a.text["button"] != "" {
					suffix = "func() {\n\tbreeze.Button(c, " + q(*a.text["button"]) + ")\n}"
				}
				return "breeze.PrefixRow(c, " + q(*a.text["title"]) + ", " + rowsSubtitle(a) + ", func() {\n\tui.Icon(c, app.modelIcon)\n}, " + suffix + ")"
			}},
		{component: "ContentRow", name: "Default", category: "Layout",
			controls: []control{textArg("label", "Prompt"), textArg("prompt", "Translate the subtitles.")},
			render: func(c *ui.Context, a *args) {
				breeze.Card(c, func() {
					breeze.ActionRow(c, "Model", "", func() { ui.Text(c, "GPT").TextColor(c.Theme().TextMuted) })
					breeze.ContentRow(c, *a.text["label"], func() { breeze.TextArea(c, a.text["prompt"]).Height(80) })
				}).FillWidth().MaxWidth(360)
			},
			code: func(a *args) string {
				return "breeze.ContentRow(c, " + q(*a.text["label"]) + ", func() {\n\tbreeze.TextArea(c, &app.prompt).Height(80)\n})"
			}},
	}
}
