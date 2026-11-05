package main

import (
	"strings"

	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, navigationStories) }

// navigationFolderIcon and navigationTrayIcon are the icons of the sidebar stories, drawn in the color of the text.
var (
	navigationFolderIcon = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><path d="M1.5 3h4.2l1.6 1.6h7.2v8.4h-13z" fill="none" stroke="#000" stroke-width="1.3" stroke-linejoin="round"/></svg>`))
	navigationTrayIcon   = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><path d="M1.5 8.5l2-5.5h9l2 5.5v4.5h-13zM1.5 8.5h4l1 1.5h3l1-1.5h4" fill="none" stroke="#000" stroke-width="1.3" stroke-linejoin="round"/></svg>`))
)

// navigationRouters keeps the history of the back and forward button stories, which must outlive a frame: one for each platform, as the platforms side by side show buttons each.
var navigationRouters = map[breeze.Platform]*ui.Router{}

func navigationRouter(c *ui.Context) *ui.Router {
	p := breeze.PlatformOf(c)
	r, ok := navigationRouters[p]
	if !ok {
		r = ui.NewRouter("/notes")
		r.Push("/notes/42")
		navigationRouters[p] = r
	}
	return r
}

// navigationHistory builds a story's back and forward buttons beside the page the router shows, and a button opening another page.
func navigationHistory(c *ui.Context, a *args) {
	r := navigationRouter(c)
	ui.Column(c).Gap(12).AlignItems(ui.Center).Children(func() {
		ui.Row(c).Gap(4).AlignItems(ui.Center).Children(func() {
			breeze.BackButton(c, r)
			breeze.ForwardButton(c, r)
			ui.Text(c, r.Path()).MarginStart(8)
		})
		if breeze.Button(c, "Open "+*a.text["page"]).Clicked() {
			r.Push(*a.text["page"])
		}
	})
}

// navigationStories are the stories of widgets_navigation.go.
func navigationStories() []*story {
	mailboxes := []string{"Inbox", "Sent", "Drafts", "Flagged"}
	return []*story{
		{component: "Sidebar", name: "Mailboxes", category: "Navigation",
			controls: []control{choiceArg("selected", "Inbox", mailboxes...), toggleArg("mailboxes", true), toggleArg("icons", true)},
			render: func(c *ui.Context, a *args) {
				icon := func(s *ui.SVG) *ui.SVG {
					if *a.on["icons"] {
						return s
					}
					return nil
				}
				breeze.Sidebar(c, a.text["selected"], func() {
					breeze.SidebarSection(c, "Mailboxes", a.on["mailboxes"], func() {
						breeze.SidebarItem(c, "Inbox", icon(navigationTrayIcon), "Inbox").Children(func() { breeze.Badge(c, "12") })
						breeze.SidebarItem(c, "Sent", icon(navigationFolderIcon), "Sent")
						breeze.SidebarItem(c, "Drafts", icon(navigationFolderIcon), "Drafts")
					})
					breeze.SidebarSection(c, "Smart Mailboxes", nil, func() {
						breeze.SidebarItem(c, "Flagged", icon(navigationFolderIcon), "Flagged")
					})
				}).Size(240, 300).Label("Mailboxes")
			},
			code: func(*args) string {
				return "breeze.Sidebar(c, &app.mailbox, func() {\n\tbreeze.SidebarSection(c, \"Mailboxes\", &app.mailboxes, func() {\n\t\tbreeze.SidebarItem(c, \"Inbox\", trayIcon, \"Inbox\").Children(func() {\n\t\t\tbreeze.Badge(c, \"12\")\n\t\t})\n\t\tbreeze.SidebarItem(c, \"Sent\", folderIcon, \"Sent\")\n\t\tbreeze.SidebarItem(c, \"Drafts\", folderIcon, \"Drafts\")\n\t})\n\tbreeze.SidebarSection(c, \"Smart Mailboxes\", nil, func() {\n\t\tbreeze.SidebarItem(c, \"Flagged\", folderIcon, \"Flagged\")\n\t})\n}).Width(240).Label(\"Mailboxes\")"
			}},
		{component: "SidebarSection", name: "Collapsible", category: "Navigation",
			controls: []control{textArg("title", "Favorites"), toggleArg("open", true), choiceArg("selected", "Desktop", "Desktop", "Documents")},
			render: func(c *ui.Context, a *args) {
				breeze.Sidebar(c, a.text["selected"], func() {
					breeze.SidebarSection(c, *a.text["title"], a.on["open"], func() {
						breeze.SidebarItem(c, "Desktop", navigationFolderIcon, "Desktop")
						breeze.SidebarItem(c, "Documents", navigationFolderIcon, "Documents")
					})
				}).Size(240, 200).Label("Places")
			},
			code: func(a *args) string {
				return "breeze.SidebarSection(c, " + q(*a.text["title"]) + ", &app.favorites, func() {\n\tbreeze.SidebarItem(c, \"Desktop\", folderIcon, \"Desktop\")\n\tbreeze.SidebarItem(c, \"Documents\", folderIcon, \"Documents\")\n})"
			}},
		{component: "SidebarItem", name: "With badge", category: "Navigation",
			controls: []control{textArg("label", "Inbox"), textArg("badge", "12"), toggleArg("chosen", true)},
			render: func(c *ui.Context, a *args) {
				selected := ""
				if *a.on["chosen"] {
					selected = "item"
				}
				breeze.Sidebar(c, &selected, func() {
					breeze.SidebarItem(c, "item", navigationTrayIcon, *a.text["label"]).Children(func() {
						if b := *a.text["badge"]; b != "" {
							breeze.Badge(c, b)
						}
					})
					breeze.SidebarItem(c, "other", navigationFolderIcon, "Archive")
				}).Size(240, 120).Label("Mailboxes")
				*a.on["chosen"] = selected == "item"
			},
			code: func(a *args) string {
				return "breeze.SidebarItem(c, \"inbox\", trayIcon, " + q(*a.text["label"]) + ").Children(func() {\n\tbreeze.Badge(c, " + q(*a.text["badge"]) + ")\n})"
			}},
		{component: "Breadcrumbs", name: "Path", category: "Navigation",
			controls: []control{textArg("path", "Macintosh HD/Users/ada/Documents"), textArg("opened", "")},
			render: func(c *ui.Context, a *args) {
				items := strings.Split(*a.text["path"], "/")
				chosen := -1
				ui.Column(c).Gap(12).AlignItems(ui.Center).Children(func() {
					if breeze.Breadcrumbs(c, items, &chosen).Label("Path").Changed() {
						*a.text["opened"] = strings.Join(items[:chosen+1], "/")
					}
					if o := *a.text["opened"]; o != "" {
						ui.Text(c, "Opened "+o).TextColor(c.Theme().TextMuted)
					}
				})
			},
			code: func(*args) string {
				return "path := []string{\"Macintosh HD\", \"Users\", \"ada\", \"Documents\"}\nif breeze.Breadcrumbs(c, path, &app.chosen).Label(\"Path\").Changed() {\n\tapp.open(path[:app.chosen+1])\n}"
			}},
		{component: "BackButton", name: "History", category: "Navigation",
			controls: []control{textArg("page", "/notes/7")},
			render:   navigationHistory,
			code:     func(*args) string { return "breeze.BackButton(c, app.router)" }},
		{component: "ForwardButton", name: "History", category: "Navigation",
			controls: []control{textArg("page", "/notes/7")},
			render:   navigationHistory,
			code:     func(*args) string { return "breeze.ForwardButton(c, app.router)" }},
		{component: "Toolbar", name: "Formatting", category: "Actions",
			controls: []control{toggleArg("bold", true), toggleArg("italic", false)},
			render: func(c *ui.Context, a *args) {
				r := navigationRouter(c)
				breeze.Toolbar(c, func() {
					breeze.BackButton(c, r)
					breeze.ForwardButton(c, r)
					breeze.Toggle(c, a.on["bold"], "Bold")
					breeze.Toggle(c, a.on["italic"], "Italic")
					ui.Spacer(c)
					breeze.Button(c, "Share")
				}).AlignSelf(ui.Stretch).Label("Formatting")
			},
			code: func(*args) string {
				return "breeze.Toolbar(c, func() {\n\tbreeze.BackButton(c, app.router)\n\tbreeze.ForwardButton(c, app.router)\n\tbreeze.Toggle(c, &app.bold, \"Bold\")\n\tbreeze.Toggle(c, &app.italic, \"Italic\")\n\tui.Spacer(c)\n\tbreeze.Button(c, \"Share\")\n}).Label(\"Formatting\")"
			}},
	}
}
