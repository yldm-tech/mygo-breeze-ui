package main

import (
	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, iconLabelStories) }

// iconLabelStories are the stories of the widgets of widgets_iconlabel.go: buttons showing an icon before their label, drawn with the icons of the icon button stories.
func iconLabelStories() []*story {
	return []*story{
		// ---- Actions ----
		{component: "IconLabelButton", name: "Default", category: "Actions",
			controls: []control{choiceArg("icon", "share", "share", "trash", "folder"), textArg("label", "Export"), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.IconLabelButton(c, inputsIcons[*a.text["icon"]], *a.text["label"]))
			},
			code: func(a *args) string {
				return "if breeze.IconLabelButton(c, " + *a.text["icon"] + "Icon, " + q(*a.text["label"]) + ").Clicked() {\n\tapp.export()\n}"
			}},
		{component: "IconLabelButton", name: "In a toolbar", category: "Actions",
			render: func(c *ui.Context, a *args) {
				breeze.Toolbar(c, func() {
					breeze.IconLabelButton(c, inputsIcons["share"], "Share")
					breeze.IconButton(c, inputsIcons["trash"], "Delete")
					breeze.PrimaryIconLabelButton(c, inputsIcons["folder"], "Move")
				}).Width(360)
			},
			code: func(*args) string {
				return "breeze.Toolbar(c, func() {\n\tbreeze.IconLabelButton(c, shareIcon, \"Share\")\n\tbreeze.IconButton(c, trashIcon, \"Delete\")\n\tbreeze.PrimaryIconLabelButton(c, folderIcon, \"Move\")\n})"
			}},
		{component: "PrimaryIconLabelButton", name: "Default", category: "Actions",
			controls: []control{choiceArg("icon", "share", "share", "trash", "folder"), textArg("label", "Upload"), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.PrimaryIconLabelButton(c, inputsIcons[*a.text["icon"]], *a.text["label"]))
			},
			code: func(a *args) string {
				return "if breeze.PrimaryIconLabelButton(c, " + *a.text["icon"] + "Icon, " + q(*a.text["label"]) + ").Clicked() {\n\tapp.upload()\n}"
			}},
	}
}
