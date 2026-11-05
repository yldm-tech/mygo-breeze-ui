package breeze

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

// collChildren returns the children of the items the outlines of these tests show: a folder of two files, and a nested folder in it.
func collChildren(k string) []string {
	switch k {
	case "src":
		return []string{"ui", "main.go"}
	case "ui":
		return []string{"view.go"}
	}
	return nil
}

// collColors are the items of the grids of these tests.
var collColors = []string{"Red", "Green", "Blue", "Gray"}

// collState is what collView's widgets edit.
type collState struct {
	table, outlineTable ui.ListState
	outline, nested     ui.OutlineState[string]
	grid                ui.GridState
	row, orow, trow     int
	item                int
	names               []string
	title               string
	srcOpen, uiOpen     bool
	file                string
}

func newCollState() *collState {
	s := &collState{row: 1, orow: 1, trow: 0, item: 1, names: append([]string(nil), collectionFiles...), title: "Readme.md", srcOpen: true, uiOpen: true, file: "Podcasts"}
	s.table.Selected, s.grid.Selected = &s.row, &s.item
	s.grid.Label = func(i int) string { return collColors[i] }
	s.outline.List.Selected, s.nested.List.Selected = &s.orow, &s.trow
	s.outline.Open.Add("src")
	s.nested.Open.Add("src")
	s.nested.Open.Add("ui")
	return s
}

// collView lays out the group's widgets: a table whose names are editable texts, an outline table, an outline, a grid, a tree whose Podcasts is chosen and an editable text.
func collView(p Platform, s *collState) func(c *ui.Context) {
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			ui.Row(c).Fill().Padding(12).Gap(12).AlignItems(ui.Start).Children(func() {
				ui.Column(c).Gap(12).Width(260).Children(func() {
					Table(c, &s.table, []ui.TableColumn{{Title: "Name"}, {Title: "Kind", Width: 70}}, len(s.names), func(row, col int) {
						if col == 0 {
							EditableText(c, &s.names[row])
							return
						}
						ui.Text(c, strings.TrimPrefix(s.names[row][strings.LastIndexByte(s.names[row], '.'):], ".")).SingleLine()
					}).Height(170).Label("Table")
					OutlineTable(c, &s.nested, []ui.TableColumn{{Title: "Name"}, {Title: "Size", Width: 60}}, []string{"src", "go.mod"}, collChildren, func(k string, col int) {
						ui.Text(c, map[int]string{0: k, 1: fmt.Sprintf("%d KB", len(k))}[col]).SingleLine()
					}).Height(200).Label("Outline table")
				})
				ui.Column(c).Gap(12).Width(220).Children(func() {
					Outline(c, &s.outline, []string{"src", "go.mod"}, collChildren, func(k string) { ui.Text(c, k).SingleLine() }).Height(130).Label("Outline")
					GridView(c, &s.grid, len(collColors), 46, 34, func(i int) { ui.Text(c, collColors[i]).SingleLine() }).Height(90).Label("Grid")
					Tree(c, func() {
						TreeItem(c, "Library", &s.srcOpen, func() {
							TreeItem(c, "Music", &s.uiOpen, func() { collTreeLeaf(c, s, "Songs.m4a") })
							collTreeLeaf(c, s, "Podcasts")
						})
						collTreeLeaf(c, s, "Movies")
					}).Label("Tree")
					EditableText(c, &s.title).Label("Title")
				})
			})
		})
	}
}

// collTreeLeaf builds a leaf of collView's tree, chosen while it is s.file, as a click chooses it.
func collTreeLeaf(c *ui.Context, s *collState, name string) {
	item := TreeItem(c, name, nil, nil)
	if item.Clicked() {
		s.file = name
	}
	item.Selected(s.file == name)
}

// TestScreenshotsWave2Collections writes the group's layout of every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names, once with each widget holding the focus: the table editing a name, the outline, the grid, the tree, and the editable text editing. Without it the test only builds the frames.
func TestScreenshotsWave2Collections(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	shots := []struct {
		name  string
		focus func(t *testing.T, tester *ui.Tester)
	}{
		{"table", func(t *testing.T, tester *ui.Tester) {
			collClick(t, tester, "Budget.xlsx")
			tester.Key(0, ui.KeyF2)
			tester.Frame()
		}},
		{"outline", func(t *testing.T, tester *ui.Tester) {
			r, _ := tester.Find("Outline")
			tester.ClickAt(r.X+r.W-20, r.Y+12)
			tester.Frame()
			tester.Key(0, ui.KeyDown)
			tester.Frame()
		}},
		{"grid", func(t *testing.T, tester *ui.Tester) { collClick(t, tester, "Blue") }},
		{"tree", func(t *testing.T, tester *ui.Tester) { collClick(t, tester, "Movies") }},
		{"edit", func(t *testing.T, tester *ui.Tester) {
			collClick(t, tester, "Readme.md")
			tester.Key(0, ui.KeyEnter)
			tester.Frame()
		}},
	}
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			for _, shot := range shots {
				s := newCollState()
				tester := ui.NewTester(collView(p, s), 520, 500)
				tester.SetScale(2)
				tester.SetDark(dark)
				tester.Frame()
				shot.focus(t, tester)
				tester.Move(0, 0)
				tester.Frame()
				for _, want := range []string{"Name", "main.go", "Blue", "Podcasts"} {
					if !tester.HasText(want) {
						t.Fatalf("%v %s: %q is missing", p, shot.name, want)
					}
				}
				writeShot(t, dir, fmt.Sprintf("breeze-wave2-collections-%s-%s-%s.png", shot.name, strings.ToLower(p.String()), mode), tester)
			}
		}
	}
}

// collClick clicks the element showing s, failing the test without one.
func collClick(t *testing.T, tester *ui.Tester, s string) {
	t.Helper()
	if err := tester.Click(s); err != nil {
		t.Fatal(err)
	}
	tester.Frame()
}
