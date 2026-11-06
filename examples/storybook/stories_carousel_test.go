package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestCarouselStoriesFollowTheirControls checks on every platform that a click on a pip of the PipsPager story chooses its page, which the text below follows, and that the pages control changes the pips and the code; and that a click on a dot of the Carousel story shows its album.
func TestCarouselStoriesFollowTheirControls(t *testing.T) {
	for platform := range 3 {
		name := platformNames[platform]
		b := newBook()
		b.selected, b.platform = "PipsPager/Default", platform
		pagerChosen = 1
		tester := ui.NewTester(b.view, 1280, 820)
		a := b.argsOf(b.story())
		tester.Frame()
		tester.Frame()
		if err := tester.Click("Page 4"); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		tester.Frame()
		if pagerChosen != 3 || !tester.HasText("Page 4 of 5") {
			t.Fatalf("%s: a click on Page 4 chose %d", name, pagerChosen)
		}
		*a.num["pages"] = 3
		tester.Frame()
		tester.Frame()
		if _, ok := tester.Find("Page 4"); ok || pagerChosen != 2 || !tester.HasText("breeze.PipsPager(c, &app.page, 3)") {
			t.Errorf("%s: the pages control did not leave 3 pages (the page chosen %d)", name, pagerChosen)
		}

		b = newBook()
		b.selected, b.platform = "Carousel/Albums", platform
		carouselShown = 0
		tester = ui.NewTester(b.view, 1280, 820)
		tester.Frame()
		tester.Frame()
		if !tester.HasText("Mountains") {
			t.Fatalf("%s: the carousel story does not show its first album", name)
		}
		if err := tester.Click("Page 3"); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		for range 3 {
			tester.Advance(1e9)
			tester.Frame()
		}
		if carouselShown != 2 || !tester.HasText("Forests") {
			t.Errorf("%s: a click on the third dot showed album %d", name, carouselShown)
		}
	}
}
