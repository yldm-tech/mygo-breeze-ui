package main

import (
	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, flatStories) }

// flatIcons are the icons the flat button stories offer besides inputsIcons, drawn in the color of the button's text.
var flatIcons = map[string]*ui.SVG{
	"more":    ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><circle cx="3" cy="8" r="1.3"/><circle cx="8" cy="8" r="1.3"/><circle cx="13" cy="8" r="1.3"/></svg>`)),
	"close":   ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><path d="M4 4l8 8M12 4l-8 8" fill="none" stroke="#000" stroke-width="1.3" stroke-linecap="round"/></svg>`)),
	"sidebar": ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><rect x="1.5" y="2.5" width="13" height="11" rx="2" fill="none" stroke="#000" stroke-width="1.3"/><path d="M6 2.5v11" stroke="#000" stroke-width="1.3"/></svg>`)),
}

// flatIcon returns the story icon named name, of flatIcons or inputsIcons.
func flatIcon(name string) *ui.SVG {
	if icon, ok := flatIcons[name]; ok {
		return icon
	}
	return inputsIcons[name]
}

// flatCardMenu fills the menu of the ⋯ button story.
func flatCardMenu(m *ui.Menu) {
	m.Item("Play")
	m.Item("Show in Finder")
	m.Separator()
	m.Item("Remove from Library…")
}

// flatStories are the stories of the widgets of widgets_flat.go.
func flatStories() []*story {
	return []*story{
		{component: "FlatButton", name: "Default", category: "Actions",
			controls: []control{textArg("label", "Skip"), toggleArg("disabled", false)},
			render:   func(c *ui.Context, a *args) { disabled(a, breeze.FlatButton(c, *a.text["label"])) },
			code: func(a *args) string {
				return "if breeze.FlatButton(c, " + q(*a.text["label"]) + ").Clicked() {\n\tapp.skip()\n}"
			}},
		{component: "FlatIconButton", name: "Default", category: "Actions",
			controls: []control{choiceArg("icon", "close", "close", "more", "share", "trash", "folder"), textArg("label", "Close"), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.FlatIconButton(c, flatIcon(*a.text["icon"]), *a.text["label"]))
			},
			code: func(a *args) string {
				return "breeze.FlatIconButton(c, " + *a.text["icon"] + "Icon, " + q(*a.text["label"]) + ")"
			}},
		{component: "FlatIconButton", name: "With a menu", category: "Actions",
			render: func(c *ui.Context, a *args) {
				breeze.FlatIconButton(c, flatIcons["more"], "More").Menu(flatCardMenu)
			},
			code: func(*args) string {
				return "breeze.FlatIconButton(c, moreIcon, \"More\").Menu(func(m *ui.Menu) {\n\tif m.Item(\"Play\").Chosen() {\n\t\tapp.play()\n\t}\n})"
			}},
		{component: "FlatToggle", name: "Default", category: "Actions",
			controls: []control{toggleArg("on", true), textArg("label", "Panel"), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.FlatToggle(c, a.on["on"], flatIcons["sidebar"], *a.text["label"]))
			},
			code: func(a *args) string {
				return "breeze.FlatToggle(c, &app.showPanel, panelIcon, " + q(*a.text["label"]) + ")"
			}},
	}
}
