package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestSidebarStories checks the sidebar stories on every platform: the add button of the Mailboxes section appends a mailbox to its control and shows it, the section staying open, and the Smart Mailboxes story shows its items with their images and chooses one with a click.
func TestSidebarStories(t *testing.T) {
	for platform := range 3 {
		b := newBook()
		b.platform = platform
		b.selected = "SidebarSectionWithActions/Mailboxes +"
		tester := ui.NewTester(b.view, 1280, 820)
		a := b.argsOf(b.story())
		tester.Frame()
		if err := tester.Click("Add Mailbox"); err != nil {
			t.Fatalf("%s: %v", platformNames[platform], err)
		}
		tester.Frame()
		if got := *a.text["mailboxes"]; got != "Inbox, Sent, Mailbox 3" || !*a.on["open"] || !tester.HasText("Mailbox 3") {
			t.Fatalf("%s: the add button left mailboxes %q, open %v", platformNames[platform], got, *a.on["open"])
		}
		b.selected = "SidebarItemImage/Smart Mailboxes"
		a = b.argsOf(b.story())
		tester.Frame()
		if err := tester.Click("Today"); err != nil {
			t.Fatalf("%s: %v", platformNames[platform], err)
		}
		tester.Frame()
		if got := *a.text["selected"]; got != "Today" || !tester.HasText("Unread") {
			t.Fatalf("%s: a click on Today chose %q", platformNames[platform], got)
		}
	}
}
