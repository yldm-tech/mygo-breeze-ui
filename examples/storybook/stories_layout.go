package main

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, layoutStories) }

// layoutSplit builds a split story: split lays out two panes, whose divider edits the story's "size" control.
func layoutSplit(c *ui.Context, a *args, split func(c *ui.Context, size *float32, first, second func()) ui.Element, first, second string) {
	size := float32(*a.num["size"])
	t := c.Theme()
	pane := func(label string) func() {
		return func() { ui.Text(c, label).Padding(10).TextColor(t.TextMuted) }
	}
	if split(c, &size, pane(first), pane(second)).FillWidth().MaxWidth(360).Height(180).Border(1, t.Border).Changed() {
		*a.num["size"] = float64(size)
	}
}

// layoutStories are the stories of the widgets of widgets_layout.go: disclosures, accordions, fieldsets, forms, check box groups and splits.
func layoutStories() []*story {
	return []*story{
		// ---- Choices ----
		{component: "CheckboxGroup", name: "Notifications", category: "Choices",
			controls: []control{textArg("label", "Notifications"), toggleArg("mail", true), toggleArg("calendar", false), toggleArg("reminders", false), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.CheckboxGroup(c, *a.text["label"], func() {
					breeze.Checkbox(c, a.on["mail"], "Mail")
					breeze.Checkbox(c, a.on["calendar"], "Calendar")
					breeze.Checkbox(c, a.on["reminders"], "Reminders")
				}))
			},
			code: func(a *args) string {
				return "breeze.CheckboxGroup(c, " + q(*a.text["label"]) + ", func() {\n\tbreeze.Checkbox(c, &app.mail, \"Mail\")\n\tbreeze.Checkbox(c, &app.calendar, \"Calendar\")\n\tbreeze.Checkbox(c, &app.reminders, \"Reminders\")\n})"
			}},

		// ---- Layout ----
		{component: "Collapsible", name: "Default", category: "Layout",
			controls: []control{textArg("label", "Advanced"), toggleArg("open", true), toggleArg("verbose", false), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.Collapsible(c, *a.text["label"], a.on["open"], func() {
					breeze.Checkbox(c, a.on["verbose"], "Verbose logging")
				}).FillWidth().MaxWidth(300))
			},
			code: func(a *args) string {
				return "breeze.Collapsible(c, " + q(*a.text["label"]) + ", &app.advanced, func() {\n\tbreeze.Checkbox(c, &app.verbose, \"Verbose logging\")\n})"
			}},
		{component: "Accordion", name: "Settings", category: "Layout",
			controls: []control{toggleArg("general", false), toggleArg("privacy", true), toggleArg("share", true)},
			render: func(c *ui.Context, a *args) {
				breeze.Accordion(c, func() {
					breeze.AccordionItem(c, "General", a.on["general"], func() {
						ui.Text(c, "Open at login and the default browser.")
					})
					breeze.AccordionItem(c, "Privacy", a.on["privacy"], func() {
						breeze.Checkbox(c, a.on["share"], "Share usage data")
					})
				}).FillWidth().MaxWidth(320)
			},
			code: func(*args) string {
				return "breeze.Accordion(c, func() {\n\tbreeze.AccordionItem(c, \"General\", &app.general, func() {\n\t\tui.Text(c, \"Open at login and the default browser.\")\n\t})\n\tbreeze.AccordionItem(c, \"Privacy\", &app.privacy, func() {\n\t\tbreeze.Checkbox(c, &app.share, \"Share usage data\")\n\t})\n})"
			}},
		{component: "AccordionItem", name: "One section", category: "Layout",
			controls: []control{textArg("title", "Privacy"), toggleArg("open", true), toggleArg("share", false)},
			render: func(c *ui.Context, a *args) {
				breeze.Accordion(c, func() {
					breeze.AccordionItem(c, *a.text["title"], a.on["open"], func() {
						breeze.Checkbox(c, a.on["share"], "Share usage data")
					})
				}).FillWidth().MaxWidth(320)
			},
			code: func(a *args) string {
				return "breeze.Accordion(c, func() {\n\tbreeze.AccordionItem(c, " + q(*a.text["title"]) + ", &app.privacy, func() {\n\t\tbreeze.Checkbox(c, &app.share, \"Share usage data\")\n\t})\n})"
			}},
		{component: "Fieldset", name: "Shipping", category: "Layout",
			controls: []control{textArg("legend", "Shipping"), textArg("address", ""), textArg("city", ""), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.Fieldset(c, *a.text["legend"], func() {
					breeze.Field(c, "Address", func() { breeze.TextInput(c, a.text["address"]).Placeholder("Street") })
					breeze.Field(c, "City", func() { breeze.TextInput(c, a.text["city"]) })
				}).FillWidth().MaxWidth(320))
			},
			code: func(a *args) string {
				return "breeze.Fieldset(c, " + q(*a.text["legend"]) + ", func() {\n\tbreeze.Field(c, \"Address\", func() { breeze.TextInput(c, &app.address).Placeholder(\"Street\") })\n\tbreeze.Field(c, \"City\", func() { breeze.TextInput(c, &app.city) })\n})"
			}},
		{component: "Form", name: "Profile", category: "Layout",
			controls: []control{textArg("name", "Ada Lovelace"), textArg("email", ""), toggleArg("newsletter", true)},
			render: func(c *ui.Context, a *args) {
				breeze.Form(c, func() {
					breeze.Field(c, "Name", func() { breeze.TextInput(c, a.text["name"]) })
					breeze.Field(c, "Email", func() { breeze.TextInput(c, a.text["email"]).Placeholder("ada@example.com") })
					breeze.Field(c, "Updates", func() { breeze.Checkbox(c, a.on["newsletter"], "Send me the newsletter") })
				}).FillWidth().MaxWidth(360)
			},
			code: func(*args) string {
				return "breeze.Form(c, func() {\n\tbreeze.Field(c, \"Name\", func() { breeze.TextInput(c, &app.name) })\n\tbreeze.Field(c, \"Email\", func() { breeze.TextInput(c, &app.email).Placeholder(\"ada@example.com\") })\n\tbreeze.Field(c, \"Updates\", func() { breeze.Checkbox(c, &app.newsletter, \"Send me the newsletter\") })\n})"
			}},
		{component: "Scroll", name: "Rows", category: "Layout",
			controls: []control{numberArg("rows", 30, 1, 100)},
			render: func(c *ui.Context, a *args) {
				t := c.Theme()
				breeze.Scroll(c).FillWidth().MaxWidth(320).Height(180).Border(1, t.Border).Label("Rows").Children(func() {
					for i := range int(*a.num["rows"]) {
						ui.Text(c, fmt.Sprintf("Row %d", i+1)).Padding(6, 10)
					}
				})
			},
			code: func(*args) string {
				return "breeze.Scroll(c).Height(180).Children(func() {\n\tfor _, row := range app.rows {\n\t\tui.Text(c, row)\n\t}\n})"
			}},
		{component: "ScrollHorizontal", name: "Tags", category: "Layout",
			controls: []control{numberArg("tags", 12, 1, 40)},
			render: func(c *ui.Context, a *args) {
				t := c.Theme()
				breeze.ScrollHorizontal(c).FillWidth().MaxWidth(320).Gap(8).Padding(8).Border(1, t.Border).Label("Tags").Children(func() {
					for i := range int(*a.num["tags"]) {
						breeze.Button(c, fmt.Sprintf("Tag %d", i+1)).Shrink(0)
					}
				})
			},
			code: func(*args) string {
				return "breeze.ScrollHorizontal(c).Gap(8).Children(func() {\n\tfor _, tag := range app.tags {\n\t\tbreeze.Button(c, tag)\n\t}\n})"
			}},
		{component: "ScrollBoth", name: "Canvas", category: "Layout",
			controls: []control{numberArg("width", 900, 100, 2000), numberArg("height", 600, 100, 2000)},
			render: func(c *ui.Context, a *args) {
				t := c.Theme()
				w, h := float32(*a.num["width"]), float32(*a.num["height"])
				breeze.ScrollBoth(c).FillWidth().MaxWidth(320).Height(200).Border(1, t.Border).Label("Canvas").Children(func() {
					ui.Box(c).Size(w, h).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
						// A grid every 50 DIPs, so the canvas shows how far it scrolled.
						for x := float32(0); x <= r.W; x += 50 {
							p.Fill(ui.Rect{X: r.X + x, Y: r.Y, W: 1, H: r.H}, t.Border, 0)
						}
						for y := float32(0); y <= r.H; y += 50 {
							p.Fill(ui.Rect{X: r.X, Y: r.Y + y, W: r.W, H: 1}, t.Border, 0)
						}
					})
				})
			},
			code: func(*args) string {
				return "breeze.ScrollBoth(c).Size(320, 200).Children(func() {\n\tui.Box(c).Size(app.width, app.height).Draw(app.drawCanvas)\n})"
			}},
		{component: "Split", name: "Sidebar", category: "Layout",
			controls: []control{numberArg("size", 120, 40, 280)},
			render:   func(c *ui.Context, a *args) { layoutSplit(c, a, breeze.Split, "Sidebar", "Content") },
			code:     func(*args) string { return "breeze.Split(c, &app.sidebar, app.files, app.editor).Fill()" }},
		{component: "SplitVertical", name: "Console", category: "Layout",
			controls: []control{numberArg("size", 100, 40, 140)},
			render:   func(c *ui.Context, a *args) { layoutSplit(c, a, breeze.SplitVertical, "Editor", "Console") },
			code: func(*args) string {
				return "breeze.SplitVertical(c, &app.editorHeight, app.editor, app.console).Fill()"
			}},
		{component: "SplitPercent", name: "Halves", category: "Layout",
			controls: []control{numberArg("size", 50, 0, 100)},
			render:   func(c *ui.Context, a *args) { layoutSplit(c, a, breeze.SplitPercent, "Original", "Translation") },
			code:     func(*args) string { return "breeze.SplitPercent(c, &app.share, app.original, app.translation).Fill()" }},
		{component: "SplitVerticalPercent", name: "Halves", category: "Layout",
			controls: []control{numberArg("size", 50, 0, 100)},
			render:   func(c *ui.Context, a *args) { layoutSplit(c, a, breeze.SplitVerticalPercent, "Preview", "Source") },
			code: func(*args) string {
				return "breeze.SplitVerticalPercent(c, &app.share, app.preview, app.source).Fill()"
			}},
	}
}
