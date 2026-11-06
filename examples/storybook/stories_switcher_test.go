package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestViewSwitcherStoryFollowsItsControls checks on every platform that a click on a view of the ViewSwitcher story chooses it, which the text below follows, and that the icons control takes the icons out of the code; and that the header bar story shows the switcher's views beside its title.
func TestViewSwitcherStoryFollowsItsControls(t *testing.T) {
	for platform := range 3 {
		b := newBook()
		b.selected, b.platform = "ViewSwitcher/Default", platform
		switcherChosen = 1
		tester := ui.NewTester(b.view, 1280, 820)
		a := b.argsOf(b.story())
		tester.Frame()
		tester.Frame()
		if err := tester.Click("Radio"); err != nil {
			t.Fatalf("%s: %v", platformNames[platform], err)
		}
		tester.Frame()
		if switcherChosen != 2 || !tester.HasText("Radio is shown.") {
			t.Fatalf("%s: a click on Radio chose %d", platformNames[platform], switcherChosen)
		}
		if !tester.HasText("breeze.ViewItem{Icon: radioIcon, Label: \"Radio\"}") {
			t.Errorf("%s: the code does not give the views their icons", platformNames[platform])
		}
		*a.on["icons"] = false
		tester.Frame()
		if !tester.HasText("breeze.ViewItem{Label: \"Radio\"}") || tester.HasText("radioIcon") {
			t.Errorf("%s: the code does not follow the icons control", platformNames[platform])
		}
		b.selected = "ViewSwitcher/In a header bar"
		tester.Frame()
		tester.Frame()
		for _, s := range []string{"Music", "Library", "Browse"} {
			if _, ok := tester.Find(s); !ok {
				t.Errorf("%s: the header bar story shows no %q", platformNames[platform], s)
			}
		}
	}
}
