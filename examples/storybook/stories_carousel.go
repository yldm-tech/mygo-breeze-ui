package main

import (
	"fmt"
	"strconv"

	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, carouselStories) }

// pagerChosen and carouselShown are the pages the PipsPager and Carousel stories choose, which the platforms side by side share.
var pagerChosen, carouselShown = 1, 0

// carouselSlides are the titles of the Carousel story's pages, as a photo app's albums.
var carouselSlides = []string{"Mountains", "Lakes", "Forests", "Deserts", "Coasts", "Glaciers", "Islands", "Canyons"}

// storyPages returns the story's pages control as a count.
func storyPages(a *args) int { return int(*a.num["pages"] + 0.5) }

// carouselStories are the stories of the widgets of widgets_carousel.go: a row of page indicators, and a view of pages with them below.
func carouselStories() []*story {
	return []*story{
		// ---- Navigation ----
		{component: "PipsPager", name: "Default", category: "Navigation",
			controls: []control{numberArg("pages", 5, 1, 12), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				n := storyPages(a)
				ui.Column(c).Gap(12).AlignItems(ui.Start).Children(func() {
					disabled(a, breeze.PipsPager(c, &pagerChosen, n))
					ui.Text(c, fmt.Sprintf("Page %d of %d", pagerChosen+1, n))
				})
			},
			code: func(a *args) string { return "breeze.PipsPager(c, &app.page, " + strconv.Itoa(storyPages(a)) + ")" }},
		{component: "Carousel", name: "Albums", category: "Navigation",
			controls: []control{numberArg("pages", 4, 1, float64(len(carouselSlides))), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				n := storyPages(a)
				disabled(a, breeze.Carousel(c, &carouselShown, n, func(i int) {
					ui.Column(c).Height(180).Padding(12).Gap(4).Justify(ui.Center).AlignItems(ui.Center).Children(func() {
						breeze.Text(c, carouselSlides[i], breeze.TextTitle)
						breeze.Text(c, fmt.Sprintf("Album %d of %d", i+1, n), breeze.TextCaption|breeze.TextSecondary)
					})
				}).Width(360))
			},
			code: func(a *args) string {
				return "breeze.Carousel(c, &app.album, " + strconv.Itoa(storyPages(a)) + ", func(i int) {\n\tbreeze.Text(c, app.albums[i], breeze.TextTitle)\n}).Width(360)"
			}},
	}
}
