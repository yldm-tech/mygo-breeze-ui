package breeze

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestComponentsBuildAndRemainDiscoverable(t *testing.T) {
	name := ""
	open := true
	tester := ui.NewTester(func(c *ui.Context) {
		Column(c, P(2), Gap(2)).Children(func() {
			Button(c, "Save", ButtonPrimary)
			Badge(c, "Ready", BadgeSuccess)
			Input(c, &name, "Name")
			Card(c, func() { Text(c, "Details") })
			Dialog(c, &open, "Details", func() { Text(c, "Body") }, nil)
		})
	}, 320, 200)
	tester.Frame()
	if !tester.HasText("Save") || !tester.HasText("Ready") || !tester.HasText("Details") {
		t.Fatal("component text is missing")
	}
}

func TestButtonVariantsAndDisabledState(t *testing.T) {
	clicked := false
	tester := ui.NewTester(func(c *ui.Context) {
		b := ButtonWith(c, "Disabled", ButtonOptions{Variant: ButtonDanger, Disabled: true})
		if b.Clicked() {
			clicked = true
		}
	}, 320, 200)
	tester.Frame()
	if err := tester.Click("Disabled"); err != nil {
		t.Fatal(err)
	}
	if clicked {
		t.Fatal("disabled button accepted a click")
	}
}
