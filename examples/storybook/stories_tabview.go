package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, tabViewStories) }

// tabViewChosen is the chosen tab of the TabView story, which the platforms side by side share as they share its tabs.
var tabViewChosen int

// tabViewTitles returns the story's tabs, the titles its tabs control lists separated by commas.
func tabViewTitles(a *args) []string {
	var titles []string
	for _, s := range strings.Split(*a.text["tabs"], ",") {
		if s = strings.TrimSpace(s); s != "" {
			titles = append(titles, s)
		}
	}
	return titles
}

// tabViewStories are the stories of the widgets of widgets_tabview.go: document tabs.
func tabViewStories() []*story {
	return []*story{
		{component: "TabView", name: "Documents", category: "Navigation",
			controls: []control{textArg("tabs", "Documents, Downloads, Pictures"), toggleArg("closable", true), toggleArg("add", true), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				titles := tabViewTitles(a)
				var closeTab func(int)
				var addTab func()
				if *a.on["closable"] {
					closeTab = func(i int) {
						*a.text["tabs"] = strings.Join(slices.Delete(titles, i, i+1), ", ")
					}
				}
				if *a.on["add"] {
					addTab = func() {
						titles = append(titles, fmt.Sprintf("Untitled %d", len(titles)+1))
						*a.text["tabs"] = strings.Join(titles, ", ")
						tabViewChosen = len(titles) - 1
					}
				}
				ui.Column(c).Gap(12).FillWidth().MaxWidth(560).Children(func() {
					disabled(a, breeze.TabView(c, &tabViewChosen, titles, closeTab, addTab))
					if len(titles) > 0 {
						ui.Text(c, titles[min(tabViewChosen, len(titles)-1)]+" is open.")
					}
				})
			},
			code: func(a *args) string {
				closeTab, addTab := "nil", "nil"
				if *a.on["closable"] {
					closeTab = "func(i int) {\n\tapp.titles = slices.Delete(app.titles, i, i+1)\n}"
				}
				if *a.on["add"] {
					addTab = "func() {\n\tapp.titles = append(app.titles, \"Untitled\")\n\tapp.tab = len(app.titles) - 1\n}"
				}
				return "breeze.TabView(c, &app.tab, app.titles, " + closeTab + ", " + addTab + ")"
			}},
	}
}
