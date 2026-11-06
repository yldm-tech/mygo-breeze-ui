package main

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, statusStories) }

// statusSearchIcon is the magnifying glass of the status page's story, drawn in the color the page gives it.
var statusSearchIcon = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><circle cx="6.5" cy="6.5" r="4.75" fill="none" stroke="#000" stroke-width="1.5"/><path d="M10 10l4.5 4.5" stroke="#000" stroke-width="1.5" stroke-linecap="round"/></svg>`))

// statusStories are the stories of the widgets of widgets_status.go.
func statusStories() []*story {
	return []*story{
		// ---- Feedback ----
		{component: "InfoBar", name: "Severities", category: "Feedback",
			controls: []control{choiceArg("kind", "warning", "info", "success", "warning", "error"), textArg("title", "Update available"), textArg("message", "Restart to finish installing."), toggleArg("closable", true), toggleArg("open", true)},
			render: func(c *ui.Context, a *args) {
				var open *bool
				if *a.on["closable"] {
					open = a.on["open"]
				}
				ui.Column(c).FillWidth().MaxWidth(560).Gap(12).AlignItems(ui.Stretch).Children(func() {
					breeze.InfoBar(c, *a.text["kind"], *a.text["title"], *a.text["message"], open)
					if open != nil && !*open && breeze.Button(c, "Show again").AlignSelf(ui.Start).Clicked() {
						*open = true
					}
				})
			},
			code: func(a *args) string {
				open := "nil"
				if *a.on["closable"] {
					open = "&app.showUpdate"
				}
				return "breeze.InfoBar(c, " + q(*a.text["kind"]) + ", " + q(*a.text["title"]) + ", " + q(*a.text["message"]) + ", " + open + ")"
			}},
		{component: "InfoBarActions", name: "With an action", category: "Feedback",
			controls: []control{choiceArg("kind", "warning", "info", "success", "warning", "error"), textArg("title", "Speech model missing"), textArg("message", "Download a model to transcribe."), textArg("action", "Open Settings"), toggleArg("primary", false), toggleArg("closable", true), toggleArg("open", true)},
			render: func(c *ui.Context, a *args) {
				var open *bool
				if *a.on["closable"] {
					open = a.on["open"]
				}
				ui.Column(c).FillWidth().MaxWidth(560).Gap(12).AlignItems(ui.Stretch).Children(func() {
					breeze.InfoBarActions(c, *a.text["kind"], *a.text["title"], *a.text["message"], open, func() {
						if *a.on["primary"] {
							breeze.PrimaryButton(c, *a.text["action"])
						} else {
							breeze.Button(c, *a.text["action"])
						}
					})
					if open != nil && !*open && breeze.Button(c, "Show again").AlignSelf(ui.Start).Clicked() {
						*open = true
					}
				})
			},
			code: func(a *args) string {
				open, button := "nil", "Button"
				if *a.on["closable"] {
					open = "&app.showNotice"
				}
				if *a.on["primary"] {
					button = "PrimaryButton"
				}
				return "breeze.InfoBarActions(c, " + q(*a.text["kind"]) + ", " + q(*a.text["title"]) + ", " + q(*a.text["message"]) + ", " + open + ", func() {\n\tif breeze." + button + "(c, " + q(*a.text["action"]) + ").Clicked() {\n\t\tapp.openSettings()\n\t}\n})"
			}},
		{component: "ProgressRing", name: "Determinate", category: "Feedback",
			controls: []control{numberArg("value", 0.6, 0, 1), toggleArg("small", false)},
			render: func(c *ui.Context, a *args) {
				r := breeze.ProgressRing(c, *a.num["value"]).Label("Upload")
				if *a.on["small"] {
					r.Size(16, 16)
				}
			},
			code: func(a *args) string {
				size := ""
				if *a.on["small"] {
					size = ".Size(16, 16)"
				}
				return fmt.Sprintf("breeze.ProgressRing(c, %.2f)%s.Label(\"Upload\")", *a.num["value"], size)
			}},
		{component: "StatusPage", name: "No results", category: "Feedback",
			controls: []control{textArg("title", "No Results"), textArg("description", "Check the spelling or try a new search."), toggleArg("icon", true), toggleArg("action", true)},
			render: func(c *ui.Context, a *args) {
				var icon *ui.SVG
				if *a.on["icon"] {
					icon = statusSearchIcon
				}
				var actions func()
				if *a.on["action"] {
					actions = func() { breeze.Button(c, "Clear Search") }
				}
				ui.Column(c).FillWidth().Height(400).Children(func() {
					breeze.StatusPage(c, icon, *a.text["title"], *a.text["description"], actions)
				})
			},
			code: func(a *args) string {
				icon, actions := "nil", "nil"
				if *a.on["icon"] {
					icon = "searchIcon"
				}
				if *a.on["action"] {
					actions = "func() {\n\tif breeze.Button(c, \"Clear Search\").Clicked() {\n\t\tapp.query = \"\"\n\t}\n}"
				}
				return "breeze.StatusPage(c, " + icon + ", " + q(*a.text["title"]) + ", " + q(*a.text["description"]) + ", " + actions + ")"
			}},
	}
}
