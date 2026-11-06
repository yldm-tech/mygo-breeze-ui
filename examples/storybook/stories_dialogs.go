package main

import (
	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, dialogsStories) }

// dialogsStories are the stories of the widgets of widgets_dialogs.go.
func dialogsStories() []*story {
	return []*story{
		// ---- Actions ----
		{component: "DestructiveButton", name: "Default", category: "Actions",
			controls: []control{textArg("label", "Delete"), toggleArg("disabled", false)},
			render:   func(c *ui.Context, a *args) { disabled(a, breeze.DestructiveButton(c, *a.text["label"])) },
			code: func(a *args) string {
				return "if breeze.DestructiveButton(c, " + q(*a.text["label"]) + ").Clicked() {\n\tapp.confirmDelete = true\n}"
			}},

		// ---- Overlays ----
		{component: "DestructiveAlertDialog", name: "Confirm", category: "Overlays",
			controls: []control{textArg("title", "Delete “Notes”?"), textArg("message", "You can’t undo this action."), toggleArg("open", false)},
			render: func(c *ui.Context, a *args) {
				if breeze.DestructiveButton(c, "Delete…").Clicked() {
					*a.on["open"] = true
				}
				breeze.DestructiveAlertDialog(c, a.on["open"], *a.text["title"], *a.text["message"], 1, "Cancel", "Delete")
			},
			code: func(a *args) string {
				return "if breeze.DestructiveAlertDialog(c, &app.asking, " + q(*a.text["title"]) + ", " + q(*a.text["message"]) + ", 1, \"Cancel\", \"Delete\") == 1 {\n\tapp.delete()\n}"
			}},
		{component: "Sheet", name: "Form", category: "Overlays",
			controls: []control{textArg("name", "Notes"), toggleArg("open", false)},
			render: func(c *ui.Context, a *args) {
				if breeze.Button(c, "Rename…").Clicked() {
					*a.on["open"] = true
				}
				breeze.Sheet(c, a.on["open"], func() {
					breeze.Text(c, "Rename", breeze.TextHeading)
					breeze.TextInput(c, a.text["name"]).Label("Name").Width(240)
					ui.Row(c).Gap(8).Justify(ui.End).Children(func() {
						if breeze.Button(c, "Cancel").Clicked() {
							*a.on["open"] = false
						}
						if breeze.PrimaryButton(c, "Rename").Clicked() {
							*a.on["open"] = false
						}
					})
				})
			},
			code: func(*args) string {
				return "breeze.Sheet(c, &app.renaming, func() {\n\tbreeze.Text(c, \"Rename\", breeze.TextHeading)\n\tbreeze.TextInput(c, &app.name).Label(\"Name\")\n\t// ...\n})"
			}},
	}
}
