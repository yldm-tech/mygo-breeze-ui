package main

import (
	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, themeStories) }

// themeAccents are the accents the ThemeOf story offers, by name: "System" is the zero color, the user's accent.
var themeAccents = map[string]string{"Purple": "#8e44ad", "Orange": "#e67e22", "Green": "#27ae60"}

// themeStories are the stories of ThemeOf in theme.go: a part of the canvas built in another appearance or an app's accent, beside the canvas's own controls.
func themeStories() []*story {
	views := []string{"Day", "Week"}
	return []*story{
		// ---- Layout ----
		{component: "ThemeOf", name: "Dark subtree", category: "Layout",
			controls: []control{toggleArg("dark", true), choiceArg("accent", "System", "System", "Purple", "Orange", "Green"), choiceArg("view", "Week", views...), toggleArg("on", true)},
			render: func(c *ui.Context, a *args) {
				var accent ui.Color
				if hex, ok := themeAccents[*a.text["accent"]]; ok {
					accent = ui.Hex(hex)
				}
				controls := func() {
					ui.Row(c).Gap(12).AlignItems(ui.Center).Wrap().Children(func() {
						breeze.Button(c, "Cancel")
						breeze.PrimaryButton(c, "Play")
						i := indexOf(views, *a.text["view"])
						if breeze.Segmented(c, &i, views...).Changed() {
							*a.text["view"] = views[i]
						}
						breeze.Switch(c, a.on["on"]).Label("Subtitles")
					})
					breeze.Link(c, "Release notes", "https://github.com/yldm-tech/mygo-breeze-ui")
				}
				controls()
				ui.Themed(c, breeze.ThemeOf(c, breeze.PlatformOf(c), *a.on["dark"], accent)).Padding(16).Gap(12).Radius(8).Children(controls)
			},
			code: func(a *args) string {
				dark, accent := "false", "ui.Color{}"
				if *a.on["dark"] {
					dark = "true"
				}
				if hex, ok := themeAccents[*a.text["accent"]]; ok {
					accent = "ui.Hex(" + q(hex) + ")"
				}
				return "ui.Themed(c, breeze.ThemeOf(c, breeze.PlatformOf(c), " + dark + ", " + accent + ")).Padding(16).Gap(12).Children(func() {\n\tbreeze.PrimaryButton(c, \"Play\")\n\tbreeze.Segmented(c, &app.view, \"Day\", \"Week\")\n\tbreeze.Switch(c, &app.subtitles).Label(\"Subtitles\")\n})"
			}},
	}
}
