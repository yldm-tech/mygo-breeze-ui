package main

import (
	"strings"

	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, switcherStories) }

// switcherIcons are the icons of the ViewSwitcher stories' views, Library, Browse and Radio, drawn in the color of the items' text.
var switcherIcons = []*ui.SVG{
	ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><path d="M2.5 2.5h3v11h-3zM6.5 2.5h3v11h-3zM10.4 3.4l2.4-.6 2.4 10.2-2.4.6z" fill="none" stroke="#000" stroke-width="1.3" stroke-linejoin="round"/></svg>`)),
	ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><path d="M2.5 2.5h4v4h-4zM9.5 2.5h4v4h-4zM2.5 9.5h4v4h-4zM9.5 9.5h4v4h-4z" fill="none" stroke="#000" stroke-width="1.3" stroke-linejoin="round"/></svg>`)),
	ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><circle cx="8" cy="8" r="1.6"/><path d="M5.3 5.3a3.8 3.8 0 0 0 0 5.4M10.7 5.3a3.8 3.8 0 0 1 0 5.4M3.2 3.2a6.8 6.8 0 0 0 0 9.6M12.8 3.2a6.8 6.8 0 0 1 0 9.6" fill="none" stroke="#000" stroke-width="1.3" stroke-linecap="round"/></svg>`)),
}

// switcherLabels are the labels of the stories' views, and switcherIconNames the names their code gives the icons.
var (
	switcherLabels    = []string{"Library", "Browse", "Radio"}
	switcherIconNames = []string{"libraryIcon", "browseIcon", "radioIcon"}
)

// switcherChosen is the chosen view of the ViewSwitcher stories, which the platforms side by side share.
var switcherChosen = 1

// switcherItems returns the stories' views, with their icons when icons is true.
func switcherItems(icons bool) []breeze.ViewItem {
	items := make([]breeze.ViewItem, len(switcherLabels))
	for i, s := range switcherLabels {
		items[i].Label = s
		if icons {
			items[i].Icon = switcherIcons[i]
		}
	}
	return items
}

// switcherCode returns the Go building the stories' switcher, with the icons when icons is true.
func switcherCode(icons bool) string {
	var b strings.Builder
	b.WriteString("breeze.ViewSwitcher(c, &app.view,\n")
	for i, s := range switcherLabels {
		if icons {
			b.WriteString("\tbreeze.ViewItem{Icon: " + switcherIconNames[i] + ", Label: " + q(s) + "},\n")
		} else {
			b.WriteString("\tbreeze.ViewItem{Label: " + q(s) + "},\n")
		}
	}
	b.WriteString(")")
	return b.String()
}

// switcherStories are the stories of the widget of widgets_switcher.go: the control that switches an app's top-level views.
func switcherStories() []*story {
	return []*story{
		// ---- Navigation ----
		{component: "ViewSwitcher", name: "Default", category: "Navigation",
			controls: []control{toggleArg("icons", true), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				ui.Column(c).Gap(12).AlignItems(ui.Start).Children(func() {
					disabled(a, breeze.ViewSwitcher(c, &switcherChosen, switcherItems(*a.on["icons"])...))
					ui.Text(c, switcherLabels[switcherChosen]+" is shown.")
				})
			},
			code: func(a *args) string { return switcherCode(*a.on["icons"]) }},
		{component: "ViewSwitcher", name: "In a header bar", category: "Navigation",
			controls: []control{textArg("title", "Music")},
			render: func(c *ui.Context, a *args) {
				breeze.HeaderBar(c, *a.text["title"], "", nil, func() {
					breeze.ViewSwitcher(c, &switcherChosen, switcherItems(true)...)
				}).MaxWidth(560)
			},
			code: func(a *args) string {
				return "breeze.HeaderBar(c, " + q(*a.text["title"]) + ", \"\", nil, func() {\n\t" + strings.ReplaceAll(switcherCode(true), "\n", "\n\t") + "\n})"
			}},
	}
}
