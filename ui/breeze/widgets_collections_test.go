package breeze

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func init() {
	for _, n := range []string{"List", "Autocomplete", "TokenField"} {
		drawnWidgets[n] = true
	}
}

var collectionFiles = []string{"Notes.md", "Budget.xlsx", "Photo.heic", "Talk.key"}

// TestListKeepsItsBehavior checks the drawn list on every platform: a click chooses a row, Down moves the choice, the list reports Changed, its rows are list items, and rows built with Rows, without a row builder, are chosen as well.
func TestListKeepsItsBehavior(t *testing.T) {
	for _, p := range Platforms {
		sel, changes := 0, 0
		other := -1
		state := &ui.ListState{Selected: &sel}
		otherState := &ui.ListState{Selected: &other}
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				ui.Row(c).Fill().Gap(8).Children(func() {
					l := List(c, state, len(collectionFiles), func(i int) { ui.Text(c, collectionFiles[i]).SingleLine() }).Size(220, 220).Label("Files")
					if l.Changed() {
						changes++
					}
					List(c, otherState, 2).Size(160, 120).Label("Other").Rows(func(r ui.ListRow) {
						ui.Text(c, []string{"One", "Two"}[r.Index])
					})
				})
			})
		}, 480, 300)
		tester.Frame()
		if err := tester.Click("Budget.xlsx"); err != nil {
			t.Fatalf("%v: click a row: %v", p, err)
		}
		tester.Frame()
		if sel != 1 || changes != 1 {
			t.Fatalf("%v: a click chose row %d with %d changes, want row 1 and 1 change", p, sel, changes)
		}
		tester.Key(0, ui.KeyDown)
		tester.Frame()
		if sel != 2 {
			t.Fatalf("%v: Down chose row %d, want 2", p, sel)
		}
		if !tester.Focused("Files") {
			t.Fatalf("%v: the list does not have the focus after a click on a row", p)
		}
		if err := tester.Click("Two"); err != nil {
			t.Fatalf("%v: click a row built with Rows: %v", p, err)
		}
		tester.Frame()
		if other != 1 {
			t.Fatalf("%v: a click on a row built with Rows chose %d", p, other)
		}
	}
}

// TestListDrawsThePlatformsRows checks that the row MyGo fills with the accent is covered by the platform's row: on Windows and GNOME the chosen row's text keeps the text color rather than the accent's, and every platform paints the chosen row unlike the others.
func TestListDrawsThePlatformsRows(t *testing.T) {
	for _, p := range Platforms {
		sel := 1
		state := &ui.ListState{Selected: &sel}
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				List(c, state, len(collectionFiles), func(i int) { ui.Text(c, collectionFiles[i]).SingleLine() }).Size(240, 240).Label("Files")
			})
		}, 300, 300)
		tester.Frame()
		chosen, ok := tester.Find("Budget.xlsx")
		other, ok2 := tester.Find("Photo.heic")
		if !ok || !ok2 {
			t.Fatalf("%v: rows missing", p)
		}
		img := tester.Image()
		// A point of the row's background just before its text, in its padding.
		at := func(r ui.Rect) [3]uint8 {
			c := img.RGBAAt(int(r.X-4), int(r.Y+r.H/2))
			return [3]uint8{c.R, c.G, c.B}
		}
		if at(chosen) == at(other) {
			t.Fatalf("%v: the chosen row looks as the others: %v", p, at(chosen))
		}
		var theme *ui.Theme
		probe := ui.NewTester(func(c *ui.Context) { theme = Theme(c, p) }, 10, 10)
		probe.Frame()
		got := at(chosen)
		accent := [3]uint8{theme.Accent.R, theme.Accent.G, theme.Accent.B}
		if p != Mac && got == accent {
			t.Fatalf("%v: the chosen row is MyGo's accent fill %v", p, got)
		}
	}
}

// TestCollectionsKeepTheirBehavior checks the collections of package ui on every platform: a grid item is chosen by a click, an outline opens with Right, an outline table shows its columns, and an editable text is edited with Enter.
func TestCollectionsKeepTheirBehavior(t *testing.T) {
	children := func(k string) []string {
		if k == "src" {
			return []string{"main.go", "view.go"}
		}
		return nil
	}
	for _, p := range Platforms {
		item := -1
		grid := &ui.GridState{Selected: &item, Label: func(i int) string { return collectionFiles[i] }}
		row := 0
		var outline ui.OutlineState[string]
		outline.List.Selected = &row
		var table ui.OutlineState[string]
		name := "Draft.txt"
		renamed := false
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				ui.Column(c).Fill().Gap(8).Padding(8).Children(func() {
					GridView(c, grid, len(collectionFiles), 80, 40, func(i int) { ui.Text(c, collectionFiles[i]) }).Height(100).Label("Photos")
					Outline(c, &outline, []string{"src", "go.mod"}, children, func(k string) { ui.Text(c, k) }).Height(120).Label("Files")
					OutlineTable(c, &table, []ui.TableColumn{{Title: "Name"}, {Title: "Kind", Width: 80}}, []string{"src"}, children, func(k string, col int) {
						ui.Text(c, map[int]string{0: k, 1: "folder"}[col])
					}).Height(100).Label("Table")
					if EditableText(c, &name).Label("Name").Changed() {
						renamed = true
					}
				})
			})
		}, 480, 520)
		tester.Frame()
		if err := tester.Click("Photo.heic"); err != nil {
			t.Fatalf("%v: click a grid item: %v", p, err)
		}
		tester.Frame()
		if item != 2 {
			t.Fatalf("%v: the grid chose %d, want 2", p, item)
		}
		if err := tester.Click("src"); err != nil {
			t.Fatalf("%v: click the outline's root: %v", p, err)
		}
		tester.Frame()
		if tester.HasText("main.go") {
			t.Fatalf("%v: the outline showed a closed item's children", p)
		}
		tester.Key(0, ui.KeyRight)
		tester.Frame()
		if !tester.HasText("main.go") {
			t.Fatalf("%v: Right did not open the outline's item", p)
		}
		if !tester.HasText("Kind") || !tester.HasText("folder") {
			t.Fatalf("%v: the outline table's columns are missing", p)
		}
		if err := tester.Click("Draft.txt"); err != nil {
			t.Fatalf("%v: click the editable text: %v", p, err)
		}
		tester.Frame()
		tester.Key(0, ui.KeyEnter)
		tester.Frame()
		tester.Type("Plans")
		tester.Frame()
		tester.Key(0, ui.KeyEnter)
		tester.Frame()
		if name != "Plans.txt" || !renamed {
			t.Fatalf("%v: editing gave %q, changed %v; want Plans.txt", p, name, renamed)
		}
	}
}

// TestTextCollectionsKeepTheirBehavior checks the fields restyled from package ui's on every platform: an autocomplete takes a suggestion with Down and Enter, and a token field adds a token with a comma and takes out the last with Backspace.
func TestTextCollectionsKeepTheirBehavior(t *testing.T) {
	for _, p := range Platforms {
		city := ""
		tags := []string{"Design"}
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				ui.Column(c).Gap(12).Padding(8).Children(func() {
					Autocomplete(c, &city, []string{"Berlin", "Bern", "Boston"}).Label("City").Width(240)
					TokenField(c, &tags, []string{"Go", "Swift"}).Label("Tags").Width(300)
				})
			})
		}, 480, 400)
		tester.Frame()
		r, ok := tester.Find("City")
		if !ok {
			t.Fatalf("%v: the autocomplete is missing", p)
		}
		tester.ClickAt(r.X+r.W/2, r.Y+r.H/2)
		tester.Frame()
		tester.Type("Ber")
		tester.Frame()
		if !tester.HasText("Berlin") || !tester.HasText("Bern") {
			t.Fatalf("%v: the suggestions do not show: %v", p, tester.Texts())
		}
		tester.Key(0, ui.KeyDown)
		tester.Frame()
		tester.Key(0, ui.KeyEnter)
		tester.Frame()
		if city != "Berlin" {
			t.Fatalf("%v: Down and Enter took %q, want Berlin", p, city)
		}
		r, ok = tester.Find("Tags")
		if !ok {
			t.Fatalf("%v: the token field is missing", p)
		}
		tester.ClickAt(r.X+r.W-10, r.Y+r.H/2)
		tester.Frame()
		tester.Type("Go,")
		tester.Frame()
		if !slices.Equal(tags, []string{"Design", "Go"}) {
			t.Fatalf("%v: a comma gave tokens %v", p, tags)
		}
		tester.Key(0, ui.KeyBackspace)
		tester.Frame()
		if !slices.Equal(tags, []string{"Design"}) {
			t.Fatalf("%v: Backspace left tokens %v", p, tags)
		}
	}
}

// collectionsView lays out the group's widgets beside a button and a text input of the same platform, to compare their paddings, radii and text: a list with its second row chosen, the fields, an editable text, a grid and an outline.
func collectionsView(p Platform, s *collectionsState) func(c *ui.Context) {
	children := func(k string) []string {
		if k == "src" {
			return []string{"main.go", "view.go"}
		}
		return nil
	}
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Row(c).Fill().Padding(16).Gap(16).AlignItems(ui.Start).Children(func() {
				ui.Column(c).Gap(12).Width(240).Children(func() {
					List(c, &s.files, len(collectionFiles), func(i int) { ui.Text(c, collectionFiles[i]).SingleLine() }).Height(160).Label("Files")
					List(c, &s.unfocused, len(collectionFiles), func(i int) { ui.Text(c, collectionFiles[i]).SingleLine() }).Height(160).Label("Unfocused")
				})
				ui.Column(c).Gap(12).Width(300).Children(func() {
					ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
						Button(c, "Button")
						TextInput(c, &s.name).Placeholder("Text input").Grow(1)
					})
					Autocomplete(c, &s.city, []string{"Berlin", "Bern", "Boston"}).Placeholder("City").Label("City")
					TokenField(c, &s.tags, []string{"Go", "Swift"}).Label("Tags")
					EditableText(c, &s.title).Label("Title")
					GridView(c, &s.grid, 4, 60, 36, func(i int) { ui.Text(c, collectionFiles[i][:4]) }).Height(90).Label("Grid")
					Outline(c, &s.outline, []string{"src", "go.mod"}, children, func(k string) { ui.Text(c, k) }).Height(110).Label("Outline")
				})
			})
		})
	}
}

// collectionsState is what collectionsView's widgets edit.
type collectionsState struct {
	files, unfocused ui.ListState
	grid             ui.GridState
	outline          ui.OutlineState[string]
	sel, other, item int
	orow             int
	name, city       string
	title            string
	tags             []string
}

func newCollectionsState() *collectionsState {
	s := &collectionsState{sel: 1, other: 2, item: 1, city: "Berlin", title: "Notes.md", tags: []string{"Design", "Go"}}
	s.files.Selected, s.unfocused.Selected, s.grid.Selected = &s.sel, &s.other, &s.item
	s.outline.List.Selected = &s.orow
	s.outline.Open.Add("src")
	return s
}

// TestScreenshotsCollections writes the collections' reference layout of every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names: once with the first list focused, once with the autocomplete's suggestions open. Without it the test only builds the frames.
func TestScreenshotsCollections(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			name := strings.ToLower(p.String())
			s := newCollectionsState()
			tester := ui.NewTester(collectionsView(p, s), 600, 440)
			tester.SetScale(2)
			tester.SetDark(dark)
			tester.Frame()
			if err := tester.Click("Budget.xlsx"); err != nil {
				t.Fatalf("%v: click the list: %v", p, err)
			}
			tester.Frame()
			for _, want := range []string{"Budget.xlsx", "Design", "Notes.md", "main.go"} {
				if !tester.HasText(want) {
					t.Fatalf("%v: %q is missing", p, want)
				}
			}
			writeShot(t, dir, fmt.Sprintf("breeze-collections-%s-%s.png", name, mode), tester)
			r, ok := tester.Find("City")
			if !ok {
				t.Fatalf("%v: the autocomplete is missing", p)
			}
			s.city = ""
			tester.ClickAt(r.X+r.W/2, r.Y+r.H/2)
			tester.Frame()
			tester.Type("B")
			tester.Frame()
			tester.Key(0, ui.KeyDown)
			tester.Frame()
			if !tester.HasText("Boston") {
				t.Fatalf("%v: the suggestions do not show", p)
			}
			writeShot(t, dir, fmt.Sprintf("breeze-collections-open-%s-%s.png", name, mode), tester)
		}
	}
}
