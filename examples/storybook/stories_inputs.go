package main

import (
	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, inputsStories) }

// inputsIcons are the icons the icon button stories offer, drawn in the color of the button's text.
var inputsIcons = map[string]*ui.SVG{
	"share":  ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><path d="M8 1.5v9M4.8 4.5 8 1.5l3.2 3M5 7H3.5v7.5h9V7H11" fill="none" stroke="#000" stroke-width="1.3" stroke-linejoin="round"/></svg>`)),
	"trash":  ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><path d="M2.5 4h11M6 4V2.5h4V4M4 4l.8 10h6.4L12 4M6.5 6.5v5M9.5 6.5v5" fill="none" stroke="#000" stroke-width="1.3" stroke-linejoin="round"/></svg>`)),
	"folder": navigationFolderIcon,
}

// inputsSaveMenu fills the split button stories' menu.
func inputsSaveMenu(m *ui.Menu) {
	m.Item("Save As…")
	m.Item("Save All")
}

// inputsStories are the stories of the widgets of widgets_inputs.go.
func inputsStories() []*story {
	return []*story{
		// ---- Actions ----
		{component: "SplitButton", name: "Default", category: "Actions",
			controls: []control{textArg("label", "Save"), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.SplitButton(c, *a.text["label"], inputsSaveMenu))
			},
			code: func(a *args) string {
				return "if breeze.SplitButton(c, " + q(*a.text["label"]) + ", func(m *ui.Menu) {\n\tif m.Item(\"Save As…\").Chosen() {\n\t\tapp.saveAs()\n\t}\n}).Clicked() {\n\tapp.save()\n}"
			}},
		{component: "SplitButton", name: "In a toolbar", category: "Actions",
			controls: []control{textArg("label", "New")},
			render: func(c *ui.Context, a *args) {
				breeze.Toolbar(c, func() {
					breeze.SplitButton(c, *a.text["label"], inputsSaveMenu)
					breeze.IconButton(c, inputsIcons["share"], "Share")
				}).Width(320)
			},
			code: func(a *args) string {
				return "breeze.Toolbar(c, func() {\n\tbreeze.SplitButton(c, " + q(*a.text["label"]) + ", app.newMenu)\n\tbreeze.IconButton(c, shareIcon, \"Share\")\n})"
			}},
		{component: "IconButton", name: "Default", category: "Actions",
			controls: []control{choiceArg("icon", "share", "share", "trash", "folder"), textArg("label", "Share"), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.IconButton(c, inputsIcons[*a.text["icon"]], *a.text["label"]))
			},
			code: func(a *args) string {
				return "breeze.IconButton(c, " + *a.text["icon"] + "Icon, " + q(*a.text["label"]) + ")"
			}},
		{component: "IconButton", name: "In a toolbar", category: "Actions",
			render: func(c *ui.Context, a *args) {
				breeze.Toolbar(c, func() {
					breeze.IconButton(c, inputsIcons["share"], "Share")
					breeze.IconButton(c, inputsIcons["trash"], "Delete")
					breeze.IconButton(c, inputsIcons["folder"], "Move")
				}).Width(320)
			},
			code: func(*args) string {
				return "breeze.Toolbar(c, func() {\n\tbreeze.IconButton(c, shareIcon, \"Share\")\n\tbreeze.IconButton(c, trashIcon, \"Delete\")\n\tbreeze.IconButton(c, folderIcon, \"Move\")\n})"
			}},

		// ---- Text ----
		{component: "PasswordField", name: "Default", category: "Text",
			controls: []control{textArg("password", "hunter2"), textArg("placeholder", "Password"), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.PasswordField(c, a.text["password"]).Placeholder(*a.text["placeholder"]).Label("Password").Width(240))
			},
			code: func(a *args) string {
				return "breeze.PasswordField(c, &app.password).Placeholder(" + q(*a.text["placeholder"]) + ").Label(\"Password\")"
			}},
	}
}
