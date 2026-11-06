package main

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, sidebarStories) }

// sidebarPlusIcon is the add button of the sidebar section stories, drawn in the color of the text; sidebarFlaggedMark, sidebarTodayMark and sidebarUnreadMark are the colored marks of the SidebarItemImage story, rounded squares in their own colors with a white glyph, as an app's tiles.
var (
	sidebarPlusIcon    = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><path d="M8 3v10M3 8h10" fill="none" stroke="#000" stroke-width="1.5" stroke-linecap="round"/></svg>`))
	sidebarFlaggedMark = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><rect width="16" height="16" rx="4" fill="#ff9500"/><path d="M5 3.5v9M5 4h6l-1.5 2.25L11 8.5H5" fill="#fff" stroke="#fff" stroke-width="1.2" stroke-linejoin="round"/></svg>`))
	sidebarTodayMark   = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><rect width="16" height="16" rx="4" fill="#0a84ff"/><rect x="3.5" y="4.5" width="9" height="8" rx="1.5" fill="none" stroke="#fff" stroke-width="1.2"/><path d="M3.5 7h9" stroke="#fff" stroke-width="1.2"/></svg>`))
	sidebarUnreadMark  = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><rect width="16" height="16" rx="4" fill="#34c759"/><circle cx="8" cy="8" r="3" fill="#fff"/></svg>`))
)

// sidebarStories are the stories of SidebarItemImage and SidebarSectionWithActions in widgets_navigation.go.
func sidebarStories() []*story {
	return []*story{
		{component: "SidebarSectionWithActions", name: "Mailboxes +", category: "Navigation",
			controls: []control{textArg("mailboxes", "Inbox, Sent"), toggleArg("open", true), choiceArg("selected", "Inbox", "Inbox", "Sent")},
			render: func(c *ui.Context, a *args) {
				names := strings.Split(*a.text["mailboxes"], ", ")
				breeze.Sidebar(c, a.text["selected"], func() {
					breeze.SidebarSectionWithActions(c, "Mailboxes", a.on["open"], func() {
						if breeze.IconButton(c, sidebarPlusIcon, "Add Mailbox").Clicked() {
							*a.text["mailboxes"] += fmt.Sprintf(", Mailbox %d", len(names)+1)
						}
					}, func() {
						for _, n := range names {
							breeze.SidebarItem(c, n, navigationFolderIcon, n)
						}
					})
				}).Size(240, 300).Label("Mailboxes")
			},
			code: func(a *args) string {
				return "breeze.SidebarSectionWithActions(c, \"Mailboxes\", &app.mailboxesOpen, func() {\n\tif breeze.IconButton(c, plusIcon, \"Add Mailbox\").Clicked() {\n\t\tapp.addMailbox()\n\t}\n}, func() {\n\tfor _, m := range app.mailboxes {\n\t\tbreeze.SidebarItem(c, m, folderIcon, m)\n\t}\n})"
			}},
		{component: "SidebarItemImage", name: "Smart Mailboxes", category: "Navigation",
			controls: []control{choiceArg("selected", "Flagged", "Flagged", "Today", "Unread")},
			render: func(c *ui.Context, a *args) {
				breeze.Sidebar(c, a.text["selected"], func() {
					breeze.SidebarSection(c, "Smart Mailboxes", nil, func() {
						breeze.SidebarItemImage(c, "Flagged", sidebarFlaggedMark, "Flagged")
						breeze.SidebarItemImage(c, "Today", sidebarTodayMark, "Today")
						breeze.SidebarItemImage(c, "Unread", sidebarUnreadMark, "Unread").Children(func() { breeze.Badge(c, "3") })
					})
				}).Size(240, 200).Label("Smart Mailboxes")
			},
			code: func(a *args) string {
				return "breeze.SidebarSection(c, \"Smart Mailboxes\", nil, func() {\n\tbreeze.SidebarItemImage(c, \"Flagged\", flaggedMark, \"Flagged\")\n\tbreeze.SidebarItemImage(c, \"Today\", todayMark, \"Today\")\n\tbreeze.SidebarItemImage(c, \"Unread\", unreadMark, \"Unread\").Children(func() {\n\t\tbreeze.Badge(c, \"3\")\n\t})\n})"
			}},
	}
}
