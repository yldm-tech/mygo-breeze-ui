package main

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

// controlKind is what a story's control edits.
type controlKind uint8

const (
	textControl controlKind = iota
	toggleControl
	numberControl
	choiceControl
)

// control is one of a story's arguments, which the controls panel edits.
type control struct {
	name    string
	kind    controlKind
	text    string   // a text control's value, or a choice control's
	on      bool     // a toggle control's value
	num     float64  // a number control's value, between lo and hi
	lo, hi  float64  // a number control's range
	options []string // a choice control's options
}

func textArg(name, value string) control { return control{name: name, kind: textControl, text: value} }
func toggleArg(name string, on bool) control {
	return control{name: name, kind: toggleControl, on: on}
}
func numberArg(name string, v, lo, hi float64) control {
	return control{name: name, kind: numberControl, num: v, lo: lo, hi: hi}
}
func choiceArg(name, value string, options ...string) control {
	return control{name: name, kind: choiceControl, text: value, options: options}
}

// args holds the values of a story's controls. The controls panel and the story bind the same values, so a widget that changes its value, as a check box does, changes the control too.
type args struct {
	text map[string]*string
	on   map[string]*bool
	num  map[string]*float64
}

func newArgs(controls []control) *args {
	a := &args{text: map[string]*string{}, on: map[string]*bool{}, num: map[string]*float64{}}
	for _, ctl := range controls {
		switch ctl.kind {
		case textControl, choiceControl:
			v := ctl.text
			a.text[ctl.name] = &v
		case toggleControl:
			v := ctl.on
			a.on[ctl.name] = &v
		case numberControl:
			v := ctl.num
			a.num[ctl.name] = &v
		}
	}
	return a
}

// story is one state of a component, as a story of Storybook: the controls it takes, how it is built, and the Go that builds it.
type story struct {
	component, name, category string
	controls                  []control
	render                    func(c *ui.Context, a *args)
	code                      func(a *args) string
}

func (s *story) id() string { return s.component + "/" + s.name }

// categories lists the sidebar's sections in order.
var categories = []string{"Actions", "Choices", "Text", "Feedback", "Navigation", "Overlays", "Layout", "Data"}

// disabled applies the story's "disabled" control to e.
func disabled(a *args, e ui.Element) ui.Element {
	if on, ok := a.on["disabled"]; ok {
		e.Disabled(*on)
	}
	return e
}

func q(s string) string { return fmt.Sprintf("%q", s) }

// storyGroups are the functions of the stories_*.go files, each adding the stories of a group of widgets in its init function. Their stories follow the ones of this file in the sidebar, in the order of the files' names within each category.
var storyGroups []func() []*story

// stories are every story of the book, in sidebar order within their categories.
func stories() []*story {
	all := baseStories()
	for _, group := range storyGroups {
		all = append(all, group()...)
	}
	return all
}

// baseStories are the stories of the widgets of widgets.go and extras.go.
func baseStories() []*story {
	sizes := []string{"Small", "Medium", "Large"}
	fonts := []string{"Serif", "Sans", "Mono"}
	return []*story{
		// ---- Actions ----
		{component: "Button", name: "Default", category: "Actions",
			controls: []control{textArg("label", "Cancel"), toggleArg("disabled", false)},
			render:   func(c *ui.Context, a *args) { disabled(a, breeze.Button(c, *a.text["label"])) },
			code:     func(a *args) string { return "breeze.Button(c, " + q(*a.text["label"]) + ")" }},
		{component: "PrimaryButton", name: "Default", category: "Actions",
			controls: []control{textArg("label", "Save"), toggleArg("disabled", false)},
			render:   func(c *ui.Context, a *args) { disabled(a, breeze.PrimaryButton(c, *a.text["label"])) },
			code: func(a *args) string {
				return "if breeze.PrimaryButton(c, " + q(*a.text["label"]) + ").Clicked() {\n\tapp.save()\n}"
			}},
		{component: "Toggle", name: "Default", category: "Actions",
			controls: []control{textArg("label", "Bold"), toggleArg("on", true), toggleArg("disabled", false)},
			render:   func(c *ui.Context, a *args) { disabled(a, breeze.Toggle(c, a.on["on"], *a.text["label"])) },
			code:     func(a *args) string { return "breeze.Toggle(c, &app.bold, " + q(*a.text["label"]) + ")" }},
		{component: "ToggleGroup", name: "Formatting", category: "Actions",
			controls: []control{toggleArg("bold", true), toggleArg("italic", false), toggleArg("underline", false)},
			render: func(c *ui.Context, a *args) {
				breeze.ToggleGroup(c, func() {
					breeze.Toggle(c, a.on["bold"], "Bold")
					breeze.Toggle(c, a.on["italic"], "Italic")
					breeze.Toggle(c, a.on["underline"], "Underline")
				}).Label("Formatting")
			},
			code: func(*args) string {
				return "breeze.ToggleGroup(c, func() {\n\tbreeze.Toggle(c, &app.bold, \"Bold\")\n\tbreeze.Toggle(c, &app.italic, \"Italic\")\n\tbreeze.Toggle(c, &app.underline, \"Underline\")\n}).Label(\"Formatting\")"
			}},
		{component: "Link", name: "Default", category: "Actions",
			controls: []control{textArg("label", "Read the docs"), textArg("url", "https://github.com/yldm-tech/mygo-breeze-ui")},
			render:   func(c *ui.Context, a *args) { breeze.Link(c, *a.text["label"], *a.text["url"]) },
			code: func(a *args) string {
				return "breeze.Link(c, " + q(*a.text["label"]) + ", " + q(*a.text["url"]) + ")"
			}},
		{component: "MenuButton", name: "Native menu", category: "Actions",
			controls: []control{textArg("label", "Sort by")},
			render: func(c *ui.Context, a *args) {
				breeze.MenuButton(c, *a.text["label"], func(m *ui.Menu) {
					m.Item("Name")
					m.Item("Date")
					m.Item("Size")
				})
			},
			code: func(a *args) string {
				return "breeze.MenuButton(c, " + q(*a.text["label"]) + ", func(m *ui.Menu) {\n\tm.Item(\"Name\")\n\tm.Item(\"Date\")\n\tm.Item(\"Size\")\n})"
			}},
		{component: "PopoverMenuButton", name: "Drawn menu", category: "Actions",
			controls: []control{textArg("label", "More")},
			render: func(c *ui.Context, a *args) {
				breeze.PopoverMenuButton(c, *a.text["label"], func(m *ui.Menu) {
					m.Item("Rename")
					m.Item("Duplicate")
					m.Item("Delete")
				})
			},
			code: func(a *args) string {
				return "breeze.PopoverMenuButton(c, " + q(*a.text["label"]) + ", func(m *ui.Menu) {\n\tm.Item(\"Rename\")\n\tm.Item(\"Duplicate\")\n\tm.Item(\"Delete\")\n})"
			}},

		// ---- Choices ----
		{component: "Checkbox", name: "Default", category: "Choices",
			controls: []control{textArg("label", "Remember me"), toggleArg("checked", true), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.Checkbox(c, a.on["checked"], *a.text["label"]))
			},
			code: func(a *args) string { return "breeze.Checkbox(c, &app.remember, " + q(*a.text["label"]) + ")" }},
		{component: "Radio", name: "Group", category: "Choices",
			controls: []control{choiceArg("size", "Medium", sizes...), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.RadioGroup(c, func() {
					for _, s := range sizes {
						breeze.Radio(c, a.text["size"], s, s)
					}
				}).Label("Size"))
			},
			code: func(*args) string {
				return "breeze.RadioGroup(c, func() {\n\tfor _, s := range []string{\"Small\", \"Medium\", \"Large\"} {\n\t\tbreeze.Radio(c, &app.size, s, s)\n\t}\n}).Label(\"Size\")"
			}},
		{component: "RadioGroup", name: "Spacing", category: "Choices",
			controls: []control{choiceArg("size", "Medium", sizes...), toggleArg("row", false), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				g := breeze.RadioGroup(c, func() {
					for _, s := range sizes {
						breeze.Radio(c, a.text["size"], s, s)
					}
				}).Label("Size")
				if *a.on["row"] {
					g.Row()
				}
				disabled(a, g)
			},
			code: func(a *args) string {
				row := ""
				if *a.on["row"] {
					row = ".Row()"
				}
				return "breeze.RadioGroup(c, func() {\n\tfor _, s := range []string{\"Small\", \"Medium\", \"Large\"} {\n\t\tbreeze.Radio(c, &app.size, s, s)\n\t}\n}).Label(\"Size\")" + row
			}},
		{component: "Switch", name: "Default", category: "Choices",
			controls: []control{toggleArg("on", true), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				ui.Row(c).Gap(10).AlignItems(ui.Center).Children(func() {
					disabled(a, breeze.Switch(c, a.on["on"]).Label("Sync"))
					ui.Text(c, "Sync")
				})
			},
			code: func(*args) string { return "breeze.Switch(c, &app.sync).Label(\"Sync\")" }},
		{component: "Select", name: "Default", category: "Choices",
			controls: []control{choiceArg("size", "Medium", sizes...), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.Select(c, a.text["size"], sizes).Label("Size"))
			},
			code: func(*args) string {
				return "breeze.Select(c, &app.size, []string{\"Small\", \"Medium\", \"Large\"}).Label(\"Size\")"
			}},
		{component: "Combobox", name: "Default", category: "Choices",
			controls: []control{choiceArg("font", "Sans", fonts...)},
			render: func(c *ui.Context, a *args) {
				breeze.Combobox(c, a.text["font"], fonts).Label("Font").Width(220)
			},
			code: func(*args) string {
				return "breeze.Combobox(c, &app.font, []string{\"Serif\", \"Sans\", \"Mono\"}).Label(\"Font\")"
			}},
		{component: "Slider", name: "Default", category: "Choices",
			controls: []control{numberArg("value", 40, 0, 100), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.Slider(c, a.num["value"], 0, 100).Width(240).Label("Volume"))
			},
			code: func(*args) string { return "breeze.Slider(c, &app.volume, 0, 100).Label(\"Volume\")" }},
		{component: "Stepper", name: "Default", category: "Choices",
			controls: []control{numberArg("value", 3, 0, 10)},
			render: func(c *ui.Context, a *args) {
				ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
					ui.Text(c, fmt.Sprintf("%.0f copies", *a.num["value"]))
					breeze.Stepper(c, a.num["value"], 0, 10, 1).Label("Copies")
				})
			},
			code: func(*args) string { return "breeze.Stepper(c, &app.copies, 0, 10, 1).Label(\"Copies\")" }},
		{component: "Segmented", name: "Default", category: "Choices",
			controls: []control{choiceArg("view", "Week", "Day", "Week", "Month")},
			render: func(c *ui.Context, a *args) {
				i := indexOf([]string{"Day", "Week", "Month"}, *a.text["view"])
				if breeze.Segmented(c, &i, "Day", "Week", "Month").Changed() {
					*a.text["view"] = []string{"Day", "Week", "Month"}[i]
				}
			},
			code: func(*args) string { return "breeze.Segmented(c, &app.view, \"Day\", \"Week\", \"Month\")" }},

		// ---- Text ----
		{component: "TextInput", name: "Default", category: "Text",
			controls: []control{textArg("value", ""), textArg("placeholder", "Name"), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.TextInput(c, a.text["value"]).Placeholder(*a.text["placeholder"]).Width(260).Label("Name"))
			},
			code: func(a *args) string {
				return "breeze.TextInput(c, &app.name).Placeholder(" + q(*a.text["placeholder"]) + ")"
			}},
		{component: "SearchField", name: "Default", category: "Text",
			controls: []control{textArg("query", "")},
			render: func(c *ui.Context, a *args) {
				breeze.SearchField(c, a.text["query"]).Width(260).Label("Search")
			},
			code: func(*args) string { return "breeze.SearchField(c, &app.query)" }},
		{component: "TextArea", name: "Default", category: "Text",
			controls: []control{textArg("value", "Notes for the meeting"), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				disabled(a, breeze.TextArea(c, a.text["value"]).Size(300, 120).Label("Notes"))
			},
			code: func(*args) string { return "breeze.TextArea(c, &app.notes)" }},
		{component: "NumberInput", name: "Default", category: "Text",
			controls: []control{numberArg("value", 12, 0, 100)},
			render: func(c *ui.Context, a *args) {
				breeze.NumberInput(c, a.num["value"], 0, 100, 1).Label("Font size")
			},
			code: func(*args) string { return "breeze.NumberInput(c, &app.size, 0, 100, 1)" }},

		// ---- Feedback ----
		{component: "Progress", name: "Determinate", category: "Feedback",
			controls: []control{numberArg("value", 0.4, 0, 1)},
			render: func(c *ui.Context, a *args) {
				breeze.Progress(c, *a.num["value"]).Width(260)
			},
			code: func(*args) string { return "breeze.Progress(c, app.done)" }},
		{component: "Progress", name: "Indeterminate", category: "Feedback",
			render: func(c *ui.Context, a *args) { breeze.Progress(c, -1).Width(260) },
			code:   func(*args) string { return "breeze.Progress(c, -1)" }},
		{component: "Spinner", name: "Default", category: "Feedback",
			render: func(c *ui.Context, a *args) { breeze.Spinner(c) },
			code:   func(*args) string { return "breeze.Spinner(c)" }},
		{component: "Badge", name: "Count", category: "Feedback",
			controls: []control{textArg("text", "3")},
			render: func(c *ui.Context, a *args) {
				ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
					ui.Text(c, "Inbox")
					breeze.Badge(c, *a.text["text"])
				})
			},
			code: func(a *args) string { return "breeze.Badge(c, " + q(*a.text["text"]) + ")" }},

		// ---- Navigation ----
		{component: "Tabs", name: "Default", category: "Navigation",
			controls: []control{choiceArg("tab", "General", "General", "Advanced", "About")},
			render: func(c *ui.Context, a *args) {
				tabs := []string{"General", "Advanced", "About"}
				i := indexOf(tabs, *a.text["tab"])
				ui.Column(c).Gap(12).Width(360).Children(func() {
					if breeze.Tabs(c, &i, tabs...).Changed() {
						*a.text["tab"] = tabs[i]
					}
					ui.Text(c, tabs[i]+" settings")
				})
			},
			code: func(*args) string { return "breeze.Tabs(c, &app.tab, \"General\", \"Advanced\", \"About\")" }},

		// ---- Overlays ----
		{component: "AlertDialog", name: "Confirm", category: "Overlays",
			controls: []control{textArg("title", "Delete “Notes”?"), textArg("message", "You can’t undo this action."), toggleArg("open", false)},
			render: func(c *ui.Context, a *args) {
				if breeze.Button(c, "Delete…").Clicked() {
					*a.on["open"] = true
				}
				breeze.AlertDialog(c, a.on["open"], *a.text["title"], *a.text["message"], "Cancel", "Delete")
			},
			code: func(a *args) string {
				return "if breeze.AlertDialog(c, &app.asking, " + q(*a.text["title"]) + ", " + q(*a.text["message"]) + ", \"Cancel\", \"Delete\") == 1 {\n\tapp.delete()\n}"
			}},
		{component: "Modal", name: "Default", category: "Overlays",
			controls: []control{toggleArg("open", false)},
			render: func(c *ui.Context, a *args) {
				if breeze.Button(c, "Show dialog").Clicked() {
					*a.on["open"] = true
				}
				breeze.Modal(c, a.on["open"], func() {
					ui.Text(c, "Export").Bold()
					ui.Text(c, "Choose where the file goes.")
					ui.Row(c).Gap(8).Justify(ui.End).Children(func() {
						if breeze.Button(c, "Cancel").Clicked() {
							*a.on["open"] = false
						}
						if breeze.PrimaryButton(c, "Export").Clicked() {
							*a.on["open"] = false
						}
					})
				})
			},
			code: func(*args) string {
				return "breeze.Modal(c, &app.exporting, func() {\n\tui.Text(c, \"Export\").Bold()\n\t// ...\n})"
			}},
		{component: "Popover", name: "Default", category: "Overlays",
			controls: []control{toggleArg("open", false)},
			render: func(c *ui.Context, a *args) {
				anchor := breeze.Button(c, "Details")
				if anchor.Clicked() {
					*a.on["open"] = true
				}
				breeze.Popover(c, anchor, a.on["open"], func() {
					ui.Text(c, "Created yesterday").Bold()
					ui.Text(c, "42 KB, edited by Ada")
				})
			},
			code: func(*args) string {
				return "anchor := breeze.Button(c, \"Details\")\nbreeze.Popover(c, anchor, &app.details, func() {\n\tui.Text(c, \"Created yesterday\")\n})"
			}},
		{component: "Tooltip", name: "On a button", category: "Overlays",
			controls: []control{textArg("tip", "Saves the document")},
			render: func(c *ui.Context, a *args) {
				breeze.Button(c, "Hover me").Tooltip(*a.text["tip"])
			},
			code: func(a *args) string { return "breeze.Button(c, \"Save\").Tooltip(" + q(*a.text["tip"]) + ")" }},

		// ---- Layout ----
		{component: "Field", name: "Labeled", category: "Layout",
			controls: []control{textArg("label", "Name"), textArg("value", "")},
			render: func(c *ui.Context, a *args) {
				breeze.Field(c, *a.text["label"], func() {
					breeze.TextInput(c, a.text["value"]).Width(240)
				})
			},
			code: func(a *args) string {
				return "breeze.Field(c, " + q(*a.text["label"]) + ", func() {\n\tbreeze.TextInput(c, &app.name)\n})"
			}},
		{component: "Divider", name: "Default", category: "Layout",
			render: func(c *ui.Context, a *args) {
				ui.Column(c).Gap(12).Width(260).Children(func() {
					ui.Text(c, "Above")
					breeze.Divider(c)
					ui.Text(c, "Below")
				})
			},
			code: func(*args) string { return "breeze.Divider(c)" }},

		// ---- Data ----
		{component: "Table", name: "Files", category: "Data",
			render: func(c *ui.Context, a *args) {
				files := [][2]string{{"Notes.md", "4 KB"}, {"Budget.xlsx", "38 KB"}, {"Photo.heic", "2.1 MB"}, {"Talk.key", "12 MB"}}
				cols := []ui.TableColumn{{Title: "Name"}, {Title: "Size", Width: 90, Align: ui.End}}
				breeze.Table(c, tableState(c), cols, len(files), func(row, col int) {
					ui.Text(c, files[row][col]).SingleLine()
				}).Size(360, 170).Label("Files")
			},
			code: func(*args) string {
				return "cols := []ui.TableColumn{{Title: \"Name\"}, {Title: \"Size\", Width: 90, Align: ui.End}}\nbreeze.Table(c, &app.files, cols, len(files), func(row, col int) {\n\tui.Text(c, files[row][col])\n})"
			}},
		{component: "Tree", name: "Folders", category: "Data",
			controls: []control{toggleArg("src", true)},
			render: func(c *ui.Context, a *args) {
				breeze.Tree(c, func() {
					breeze.TreeItem(c, "src", a.on["src"], func() {
						breeze.TreeItem(c, "main.go", nil, nil)
						breeze.TreeItem(c, "view.go", nil, nil)
					})
					breeze.TreeItem(c, "go.mod", nil, nil)
				}).Width(260).Label("Files")
			},
			code: func(*args) string {
				return "breeze.Tree(c, func() {\n\tbreeze.TreeItem(c, \"src\", &app.srcOpen, func() {\n\t\tbreeze.TreeItem(c, \"main.go\", nil, nil)\n\t})\n})"
			}},
	}
}

func indexOf(list []string, s string) int {
	for i, v := range list {
		if v == s {
			return i
		}
	}
	return 0
}

// tables keeps the table story's list states, which must outlive a frame: one for each platform, as the platforms side by side show a table each and a ListState belongs to one list.
var tables = map[breeze.Platform]*ui.ListState{}

func tableState(c *ui.Context) *ui.ListState {
	p := breeze.PlatformOf(c)
	s, ok := tables[p]
	if !ok {
		row := 0
		s = &ui.ListState{Selected: &row}
		tables[p] = s
	}
	return s
}
