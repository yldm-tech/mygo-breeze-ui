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
		{component: "HeaderBar", name: "In a content pane", category: "Navigation",
			controls: []control{choiceArg("page", "All films", "All films", "Downloads"), textArg("subtitle", "Library"), toggleArg("start", true)},
			render: func(c *ui.Context, a *args) {
				ui.Row(c).MaxWidth(720).Height(220).AlignItems(ui.Stretch).Children(func() {
					breeze.Sidebar(c, a.text["page"], func() {
						breeze.SidebarItem(c, "All films", navigationFolderIcon, "All films")
						breeze.SidebarItem(c, "Downloads", navigationFolderIcon, "Downloads")
					}).Shrink(0).Width(200).Label("Library")
					ui.Column(c).Grow(1).Shrink(1).MinWidth(0).Children(func() {
						var start func()
						if *a.on["start"] {
							start = func() { breeze.IconButton(c, inputsIcons["folder"], "Folders") }
						}
						// Over the pane the bar keeps clear of only the window controls above it.
						breeze.HeaderBar(c, *a.text["page"], *a.text["subtitle"], start, func() {
							breeze.PrimaryButton(c, "Start")
						})
						ui.Column(c).Grow(1).Padding(16).Children(func() { breeze.Text(c, "Content") })
					})
				})
			},
			code: func(a *args) string {
				start := "nil"
				if *a.on["start"] {
					start = "func() {\n\t\t\tbreeze.IconButton(c, folderIcon, \"Folders\")\n\t\t}"
				}
				return "ui.Row(c).Fill().Children(func() {\n\tbreeze.Sidebar(c, &app.page, func() {\n\t\tbreeze.SidebarItem(c, \"All films\", folderIcon, \"All films\")\n\t\tbreeze.SidebarItem(c, \"Downloads\", folderIcon, \"Downloads\")\n\t}).Width(200).Label(\"Library\")\n\tui.Column(c).Grow(1).Children(func() {\n\t\t// Over the pane the bar keeps clear of only the window controls above it.\n\t\tbreeze.HeaderBar(c, app.page, " + q(*a.text["subtitle"]) + ", " + start + ", func() {\n\t\t\tbreeze.PrimaryButton(c, \"Start\")\n\t\t})\n\t\tbuildContent()\n\t})\n})"
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
