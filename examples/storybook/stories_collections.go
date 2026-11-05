package main

import (
	"strings"

	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, collectionsStories) }

// collectionsFiles are the rows the collections' stories show.
var collectionsFiles = []string{"Notes.md", "Budget.xlsx", "Photo.heic", "Talk.key", "Draft.txt"}

// collectionsChildren returns the children of the outline stories' folders, nil for a file.
func collectionsChildren(k string) []string {
	switch k {
	case "src":
		return []string{"src/main.go", "src/view.go"}
	case "docs":
		return []string{"docs/guide.md"}
	}
	return nil
}

// collectionsName returns the last part of a path of the outline stories.
func collectionsName(k string) string { return k[strings.LastIndexByte(k, '/')+1:] }

// collectionsStories are the stories of the widgets of widgets_collections.go.
func collectionsStories() []*story {
	cities := []string{"Amsterdam", "Berlin", "Bern", "Boston", "Brussels", "Lisbon"}
	return []*story{
		// ---- Text ----
		{component: "Autocomplete", name: "Cities", category: "Text",
			controls: []control{textArg("city", ""), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.Autocomplete(c, a.text["city"], cities).Placeholder("City").Width(260).Label("City"))
			},
			code: func(*args) string {
				return "breeze.Autocomplete(c, &app.city, []string{\"Amsterdam\", \"Berlin\", \"Bern\", \"Boston\"}).Label(\"City\")"
			}},
		{component: "TokenField", name: "Tags", category: "Text",
			controls: []control{textArg("tags", "Design, Go")},
			render: func(c *ui.Context, a *args) {
				tags := collectionsSplit(*a.text["tags"])
				if breeze.TokenField(c, &tags, []string{"Design", "Go", "Swift", "Rust", "Linux"}).Width(300).Label("Tags").Changed() {
					*a.text["tags"] = strings.Join(tags, ", ")
				}
			},
			code: func(*args) string {
				return "breeze.TokenField(c, &app.tags, []string{\"Design\", \"Go\", \"Swift\"}).Label(\"Tags\")"
			}},
		{component: "EditableText", name: "Rename", category: "Text",
			controls: []control{textArg("name", "Notes.md")},
			render: func(c *ui.Context, a *args) {
				ui.Column(c).Gap(6).Children(func() {
					breeze.EditableText(c, a.text["name"]).Label("File name")
					ui.Text(c, "Double-click or press Enter to rename").TextColor(c.Theme().TextMuted)
				})
			},
			code: func(*args) string {
				return "if breeze.EditableText(c, &app.name).Changed() {\n\tapp.rename()\n}"
			}},

		// ---- Data ----
		{component: "List", name: "Files", category: "Data",
			controls: []control{choiceArg("selected", "Budget.xlsx", collectionsFiles...)},
			render: func(c *ui.Context, a *args) {
				s := collectionsListState(c)
				i := indexOf(collectionsFiles, *a.text["selected"])
				s.Selected = &i
				if breeze.List(c, s, len(collectionsFiles), func(row int) {
					ui.Text(c, collectionsFiles[row]).SingleLine()
				}).Size(260, 200).Label("Files").Changed() {
					*a.text["selected"] = collectionsFiles[i]
				}
			},
			code: func(*args) string {
				return "app.files.Selected = &app.file\nbreeze.List(c, &app.files, len(names), func(i int) {\n\tui.Text(c, names[i])\n}).Label(\"Files\")"
			}},
		{component: "GridView", name: "Items", category: "Data",
			controls: []control{choiceArg("selected", "Photo.heic", collectionsFiles...)},
			render: func(c *ui.Context, a *args) {
				s := collectionsGridState(c)
				i := indexOf(collectionsFiles, *a.text["selected"])
				s.Selected = &i
				if breeze.GridView(c, s, len(collectionsFiles), 80, 56, func(item int) {
					ui.Column(c).Grow(1).Center().Children(func() { ui.Text(c, collectionsFiles[item]).SingleLine() })
				}).Size(300, 150).Label("Items").Changed() {
					*a.text["selected"] = collectionsFiles[i]
				}
			},
			code: func(*args) string {
				return "app.items.Selected = &app.item\nbreeze.GridView(c, &app.items, len(names), 80, 56, func(i int) {\n\tui.Text(c, names[i])\n})"
			}},
		{component: "Outline", name: "Folders", category: "Data",
			controls: []control{toggleArg("src open", true)},
			render: func(c *ui.Context, a *args) {
				s := collectionsOutlineState(c)
				collectionsSyncOpen(s, a.on["src open"], func() {
					breeze.Outline(c, s, []string{"src", "docs", "go.mod"}, collectionsChildren, func(k string) {
						ui.Text(c, collectionsName(k)).SingleLine()
					}).Size(260, 180).Label("Files")
				})
			},
			code: func(*args) string {
				return "app.files.List.Selected = &app.file\nbreeze.Outline(c, &app.files, []string{\"src\", \"docs\"}, app.children, func(path string) {\n\tui.Text(c, filepath.Base(path))\n})"
			}},
		{component: "OutlineTable", name: "Files", category: "Data",
			controls: []control{toggleArg("src open", true)},
			render: func(c *ui.Context, a *args) {
				s := collectionsOutlineTableState(c)
				cols := []ui.TableColumn{{Title: "Name"}, {Title: "Kind", Width: 80}}
				collectionsSyncOpen(s, a.on["src open"], func() {
					breeze.OutlineTable(c, s, cols, []string{"src", "docs", "go.mod"}, collectionsChildren, func(k string, col int) {
						if col == 0 {
							ui.Text(c, collectionsName(k)).SingleLine()
							return
						}
						kind := "File"
						if collectionsChildren(k) != nil {
							kind = "Folder"
						}
						ui.Text(c, kind).SingleLine()
					}).Size(290, 200).Label("Files")
				})
			},
			code: func(*args) string {
				return "cols := []ui.TableColumn{{Title: \"Name\"}, {Title: \"Kind\", Width: 90}}\nbreeze.OutlineTable(c, &app.files, cols, roots, app.children, func(path string, col int) {\n\tui.Text(c, app.cell(path, col))\n})"
			}},
	}
}

// collectionsSplit returns the tags of a comma-separated control's text.
func collectionsSplit(s string) []string {
	var tags []string
	for _, t := range strings.Split(s, ",") {
		if t = strings.TrimSpace(t); t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}

// collectionsSyncOpen opens the outline's "src" as the control says, builds it, and keeps the control as the user opens or closes it.
func collectionsSyncOpen(s *ui.OutlineState[string], open *bool, build func()) {
	if *open {
		s.Open.Add("src")
	} else {
		s.Open.Remove("src")
	}
	build()
	*open = s.Open.Has("src")
}

// The collections' stories keep a state for each platform, which must outlive a frame, as the platforms side by side show a widget each and a state belongs to one widget.
var (
	collectionsLists         = map[breeze.Platform]*ui.ListState{}
	collectionsGrids         = map[breeze.Platform]*ui.GridState{}
	collectionsOutlines      = map[breeze.Platform]*ui.OutlineState[string]{}
	collectionsOutlineTables = map[breeze.Platform]*ui.OutlineState[string]{}
)

func collectionsListState(c *ui.Context) *ui.ListState {
	p := breeze.PlatformOf(c)
	if collectionsLists[p] == nil {
		collectionsLists[p] = &ui.ListState{}
	}
	return collectionsLists[p]
}

func collectionsGridState(c *ui.Context) *ui.GridState {
	p := breeze.PlatformOf(c)
	if collectionsGrids[p] == nil {
		collectionsGrids[p] = &ui.GridState{Label: func(i int) string { return collectionsFiles[i] }}
	}
	return collectionsGrids[p]
}

// collectionsOutline returns states[p], made on the first call with a row chosen.
func collectionsOutline(states map[breeze.Platform]*ui.OutlineState[string], c *ui.Context) *ui.OutlineState[string] {
	p := breeze.PlatformOf(c)
	if states[p] == nil {
		row := 0
		states[p] = &ui.OutlineState[string]{}
		states[p].List.Selected = &row
	}
	return states[p]
}

func collectionsOutlineState(c *ui.Context) *ui.OutlineState[string] {
	return collectionsOutline(collectionsOutlines, c)
}

func collectionsOutlineTableState(c *ui.Context) *ui.OutlineState[string] {
	return collectionsOutline(collectionsOutlineTables, c)
}
