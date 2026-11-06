package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestTabViewStoryOpensAndClosesTabs checks on every platform that the TabView story keeps its tabs in its tabs control: the add button appends one and chooses it, a click on the chosen tab's close button removes it, and the closable and add controls take the buttons away and change the code.
func TestTabViewStoryOpensAndClosesTabs(t *testing.T) {
	names := map[int][2]string{0: {"Close", "New Tab"}, 1: {"Close Tab", "Add New Tab"}, 2: {"Close Tab", "New Tab"}}
	for platform := range 3 {
		b := newBook()
		b.selected, b.platform = "TabView/Documents", platform
		tabViewChosen = 0
		tester := ui.NewTester(b.view, 1280, 820)
		a := b.argsOf(b.story())
		tester.Frame()
		if err := tester.Click(names[platform][1]); err != nil {
			t.Fatalf("%s: %v", platformNames[platform], err)
		}
		tester.Frame()
		tester.Frame()
		if got := *a.text["tabs"]; got != "Documents, Downloads, Pictures, Untitled 4" || tabViewChosen != 3 || !tester.HasText("Untitled 4 is open.") {
			t.Fatalf("%s: the add button left tabs %q, %d chosen", platformNames[platform], got, tabViewChosen)
		}
		tab, _ := tester.Find("Untitled 4")
		tester.Move(tab.X+tab.W/2, tab.Y+tab.H/2)
		tester.Frame()
		x := tab.X + tab.W - 20
		if platform == 0 {
			x = tab.X + 13
		}
		tester.ClickAt(x, tab.Y+tab.H/2)
		tester.Frame()
		tester.Frame()
		if got := *a.text["tabs"]; got != "Documents, Downloads, Pictures" || tabViewChosen != 2 {
			t.Fatalf("%s: the close button left tabs %q, %d chosen", platformNames[platform], got, tabViewChosen)
		}
		*a.on["closable"], *a.on["add"] = false, false
		tester.Frame()
		for _, s := range names[platform] {
			if _, ok := tester.Find(s); ok {
				t.Errorf("%s: %q shows with its control off", platformNames[platform], s)
			}
		}
		if !tester.HasText("breeze.TabView(c, &app.tab, app.titles, nil, nil)") {
			t.Errorf("%s: the code does not follow the controls", platformNames[platform])
		}
	}
}
