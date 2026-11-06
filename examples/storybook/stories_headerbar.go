package main

import (
	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, headerBarStories) }

// headerBarStories are the stories of the widget of widgets_headerbar.go: the title bar a view draws in a window whose title bar is hidden. The storybook's window keeps its own title bar, so c.TitleBar() is zero here and the bar keeps clear of no window controls.
func headerBarStories() []*story {
	return []*story{
		// ---- Navigation ----
		{component: "HeaderBar", name: "Default", category: "Navigation",
			controls: []control{textArg("title", "Documents"), textArg("subtitle", "12 items"), toggleArg("start", true), toggleArg("end", true)},
			render: func(c *ui.Context, a *args) {
				var start, end func()
				if *a.on["start"] {
					start = func() { breeze.IconButton(c, inputsIcons["folder"], "Folders") }
				}
				if *a.on["end"] {
					end = func() {
						breeze.Button(c, "Edit")
						breeze.IconButton(c, inputsIcons["share"], "Share")
					}
				}
				breeze.HeaderBar(c, *a.text["title"], *a.text["subtitle"], start, end).MaxWidth(560)
			},
			code: func(a *args) string {
				start, end := "nil", "nil"
				if *a.on["start"] {
					start = "func() {\n\tbreeze.IconButton(c, folderIcon, \"Folders\")\n}"
				}
				if *a.on["end"] {
					end = "func() {\n\tbreeze.Button(c, \"Edit\")\n\tbreeze.IconButton(c, shareIcon, \"Share\")\n}"
				}
				return "breeze.HeaderBar(c, " + q(*a.text["title"]) + ", " + q(*a.text["subtitle"]) + ", " + start + ", " + end + ")"
			}},
		{component: "HeaderBar", name: "With history", category: "Navigation",
			controls: []control{textArg("title", "Notes")},
			render: func(c *ui.Context, a *args) {
				r := navigationRouter(c)
				breeze.HeaderBar(c, *a.text["title"], r.Path(), func() {
					breeze.BackButton(c, r)
					breeze.ForwardButton(c, r)
				}, func() {
					if breeze.Button(c, "Open").Clicked() {
						r.Push("/notes/42")
					}
				}).MaxWidth(560)
			},
			code: func(a *args) string {
				return "breeze.HeaderBar(c, " + q(*a.text["title"]) + ", app.router.Path(), func() {\n\tbreeze.BackButton(c, app.router)\n\tbreeze.ForwardButton(c, app.router)\n}, func() {\n\tif breeze.Button(c, \"Open\").Clicked() {\n\t\tapp.router.Push(\"/notes/42\")\n\t}\n})"
			}},
	}
}
