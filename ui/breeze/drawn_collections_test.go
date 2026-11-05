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

// collClick clicks the element showing s, failing the test without one, and moves the pointer away, so that what it clicked shows as it does not under the pointer.
func collClick(t *testing.T, tester *ui.Tester, s string) {
	t.Helper()
	if err := tester.Click(s); err != nil {
		t.Fatal(err)
	}
	tester.Frame()
	tester.Move(0, 0)
	tester.Frame()
}

func init() {
	for _, n := range []string{"EditableText", "GridView", "Outline", "OutlineTable", "Table", "Tree", "TreeItem"} {
		drawnWidgets[n] = true
	}
}

// collThemeOf returns the theme Breeze UI gives p's widgets, light or dark.
func collThemeOf(p Platform, dark bool) *ui.Theme {
	var th *ui.Theme
	probe := ui.NewTester(func(c *ui.Context) { th = Theme(c, p) }, 10, 10)
	probe.SetDark(dark)
	probe.Frame()
	return th
}

// collAt returns the color of the last frame at (x, y), in DIPs at a scale of 1.
func collAt(tester *ui.Tester, x, y float32) ui.Color {
	c := tester.Image().RGBAAt(int(x), int(y))
	return ui.RGB(c.R, c.G, c.B)
}

// collNear reports whether a and b differ by at most tol in each channel.
func collNear(a, b ui.Color, tol int) bool {
	d := func(x, y uint8) int { return max(int(x)-int(y), int(y)-int(x)) }
	return d(a.R, b.R) <= tol && d(a.G, b.G) <= tol && d(a.B, b.B) <= tol
}

// collChosenFill returns the fill of a chosen row of p's collections in a light theme, not under the pointer: AppKit's selection with or without the focus, WinUI's SubtleFillColorSecondary backplate, GNOME's accent at 25% over the view.
func collChosenFill(p Platform, th *ui.Theme, focused bool) ui.Color {
	switch p {
	case Windows:
		fill, _ := fluentSubtle(th)
		return fill
	case Linux:
		return adwListViewBg(th).Mix(th.Accent, 0.25)
	}
	return macListSelection(th, focused)
}

// collLuminanceRange returns the darkest and lightest luminance of the last frame in r.
func collLuminanceRange(tester *ui.Tester, r ui.Rect) (lo, hi float32) {
	lo, hi = 1, 0
	img := tester.Image()
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			c := img.RGBAAt(x, y)
			l := luminance(ui.RGB(c.R, c.G, c.B))
			lo, hi = min(lo, l), max(hi, l)
		}
	}
	return lo, hi
}

// TestTableKeepsItsBehavior checks the drawn table on every platform: the chosen row is the platform's (AppKit's gray selection without the focus and the accent one with it, WinUI's backplate, GNOME's tinted row, never MyGo's accent fill) with the text's own color, or white on AppKit's accent; a click chooses a row and gives the table the focus, Down moves the choice, Changed reports it, a click on a Sortable header sorts, dragging a header's edge resizes its column, F2 renames an editable text in the chosen row, rows under the pointer light up except on macOS, and a disabled table chooses nothing.
func TestTableKeepsItsBehavior(t *testing.T) {
	cols := []ui.TableColumn{{Title: "Name"}, {Title: "Kind", Width: 70}, {Title: "Size", Width: 80, Sortable: true}}
	for _, p := range Platforms {
		th := collThemeOf(p, false)
		names := append([]string(nil), collectionFiles...)
		sel, changes := 1, 0
		var sort ui.SortOrder
		state := &ui.ListState{Selected: &sel, Sort: &sort}
		off := 0
		offState := &ui.ListState{Selected: &off}
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Row(c).Fill().Gap(8).Padding(8).AlignItems(ui.Start).Children(func() {
					tb := Table(c, state, cols, len(names), func(row, col int) {
						if col == 0 {
							EditableText(c, &names[row])
							return
						}
						ui.Text(c, map[int]string{1: "file", 2: fmt.Sprintf("%d KB", row+1)}[col]).SingleLine()
					}).Size(340, 260).Label("Files")
					if tb.Changed() {
						changes++
					}
					Table(c, offState, cols[:1], 2, func(row, col int) { ui.Text(c, []string{"One", "Two"}[row]) }).Size(140, 120).Label("Off").Disabled(true)
				})
			})
		}, 520, 300)
		tester.Frame()
		// A point of a row in its first cell's padding, before its text.
		pad := func(name string) (float32, float32) {
			r, ok := tester.Find(name)
			if !ok {
				t.Fatalf("%v: %q is missing", p, name)
			}
			return r.X - 4, r.Y + r.H/2
		}
		x, y := pad("Budget.xlsx")
		if got, want := collAt(tester, x, y), collChosenFill(p, th, false); !collNear(got, want, 2) {
			t.Fatalf("%v: the chosen row without the focus is %v, want %v", p, got, want)
		}
		collClick(t, tester, "Photo.heic")
		if sel != 2 || changes != 1 || !tester.Focused("Files") {
			t.Fatalf("%v: a click chose row %d with %d changes, focused %v; want row 2, 1 change, focused", p, sel, changes, tester.Focused("Files"))
		}
		x, y = pad("Photo.heic")
		got := collAt(tester, x, y)
		if want := collChosenFill(p, th, true); !collNear(got, want, 2) || collNear(got, th.Accent, 8) && p != Mac {
			t.Fatalf("%v: the chosen row with the focus is %v, want %v", p, got, want)
		}
		text, _ := tester.Find("Photo.heic")
		if lo, _ := collLuminanceRange(tester, text); (p == Mac) != (lo > luminance(got)-0.05) {
			t.Fatalf("%v: the chosen row's text reaches a luminance of %.2f over its row's %.2f; only AppKit's is white", p, lo, luminance(got))
		}
		nx, ny := pad("Notes.md")
		rest := collAt(tester, nx, ny)
		notes, _ := tester.Find("Notes.md")
		tester.Move(notes.X+notes.W+40, notes.Y+notes.H/2)
		tester.Frame()
		if hovered := collAt(tester, nx, ny); (p == Mac) != (hovered == rest) {
			t.Fatalf("%v: a row under the pointer is %v, at rest %v; only AppKit's rows do not light up", p, hovered, rest)
		}
		tester.Move(0, 0)
		tester.Key(0, ui.KeyDown)
		tester.Frame()
		if sel != 3 {
			t.Fatalf("%v: Down chose row %d, want 3", p, sel)
		}
		tester.Key(0, ui.KeyF2)
		tester.Frame()
		tester.Type("Slides")
		tester.Frame()
		tester.Key(0, ui.KeyEnter)
		tester.Frame()
		if names[3] != "Slides.key" || !tester.Focused("Files") {
			t.Fatalf("%v: F2 renamed the row %q, focused %v; want Slides.key and the table focused", p, names[3], tester.Focused("Files"))
		}
		collClick(t, tester, "Size")
		if sort.Column != "Size" {
			t.Fatalf("%v: a click on the Size header sorted by %q", p, sort.Column)
		}
		kind, _ := tester.Find("Kind")
		// The edge at the end of the Kind column's header: its title is Space(2.5) into the 70-DIP cell.
		ex, ey := kind.X-th.Space(2.5)+70-2, kind.Y+kind.H/2
		tester.Press(ex, ey)
		tester.Move(ex-30, ey)
		tester.Release(ex-30, ey)
		tester.Frame()
		if w := state.Columns.Widths["Kind"]; w != 40 {
			t.Fatalf("%v: dragging the Kind header's edge 30 DIPs in left its width %v, want 40", p, w)
		}
		if err := tester.Click("Two"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		if off != 0 {
			t.Fatalf("%v: a click chose row %d of a disabled table", p, off)
		}
	}
}

// TestTableMirrorsRightToLeft checks a right-to-left table on every platform: a click still chooses, the chosen row is the platform's, and Windows's pill is at the row's start, on the right.
func TestTableMirrorsRightToLeft(t *testing.T) {
	cols := []ui.TableColumn{{Title: "Name"}, {Title: "Size", Width: 80}}
	for _, p := range Platforms {
		th := collThemeOf(p, false)
		sel := 0
		state := &ui.ListState{Selected: &sel}
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				Table(c, state, cols, len(collectionFiles), func(row, col int) {
					ui.Text(c, map[int]string{0: collectionFiles[row], 1: "file"}[col]).SingleLine()
				}).Size(300, 220).Label("Files")
			})
		}, 320, 240)
		tester.SetDirection(ui.RightToLeft)
		tester.Frame()
		collClick(t, tester, "Budget.xlsx")
		if sel != 1 {
			t.Fatalf("%v: a click chose row %d, want 1", p, sel)
		}
		files, _ := tester.Find("Files")
		row, _ := tester.Find("Budget.xlsx")
		y := row.Y + row.H/2
		// Past the end of the text, in its cell's padding on the right.
		if got, want := collAt(tester, row.X+row.W+4, y), collChosenFill(p, th, true); !collNear(got, want, 2) {
			t.Fatalf("%v: the chosen row right to left is %v, want %v", p, got, want)
		}
		if p == Windows {
			// The pill is 3 DIPs wide at the backplate's start, 4 in from the table's right edge; the left end holds none.
			end, start := collAt(tester, files.X+files.W-4-1.5, y), collAt(tester, files.X+4+1.5, y)
			if !collNear(end, th.Accent, 2) || collNear(start, th.Accent, 40) {
				t.Fatalf("%v: the pill is not at the right of a right-to-left table: end %v, start %v, accent %v", p, end, start, th.Accent)
			}
		}
	}
}

// TestOutlineKeepsItsBehavior checks the drawn outline on every platform, left to right and right to left: an Indent of zero becomes the platform's, the chosen row's fill reaches over MyGo's arrow box, its chevron there is not white without the focus and white on AppKit's accent with it, a click on the drawn chevron opens and closes the item, with Option (Alt) all inside, a click on a row chooses it and gives the outline the focus, Changed reports it, and the arrows close and open the item chosen.
func TestOutlineKeepsItsBehavior(t *testing.T) {
	for _, p := range Platforms {
		for _, rtl := range []bool{false, true} {
			th := collThemeOf(p, false)
			sel, changes := 0, 0
			var s ui.OutlineState[string]
			s.List.Selected = &sel
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					ui.Column(c).Fill().Padding(8).Children(func() {
						if Outline(c, &s, []string{"src", "go.mod"}, collChildren, func(k string) { ui.Text(c, k).SingleLine() }).Size(260, 240).Label("Files").Changed() {
							changes++
						}
					})
				})
			}, 300, 280)
			if rtl {
				tester.SetDirection(ui.RightToLeft)
			}
			tester.Frame()
			name := fmt.Sprintf("%v rtl=%v", p, rtl)
			if s.Indent != collIndent(p) {
				t.Fatalf("%s: the outline's Indent is %v, want %v", name, s.Indent, collIndent(p))
			}
			// MyGo's arrow box of the item src: Space(1) before the cover, which Windows pads 8 more at the top level.
			arrow := func() ui.Rect {
				r, _ := tester.Find("src")
				start := float32(0)
				if p == Windows {
					start = 8
				}
				box := ui.Rect{X: r.X - start - th.Space(1) - th.Space(4), Y: r.Y + r.H/2 - 8, W: 16, H: 16}
				if rtl {
					box.X = r.X + r.W + start + th.Space(1)
				}
				return box
			}
			box := arrow()
			// One DIP inside the arrow box's outer side: the row's fill reaches it.
			side := box.X + 1
			if rtl {
				side = box.X + box.W - 1
			}
			fill := collChosenFill(p, th, false)
			if got := collAt(tester, side, box.Y+8); !collNear(got, fill, 2) {
				t.Fatalf("%s: the chosen row at its arrow box is %v, want %v", name, got, fill)
			}
			// The chevron's 8 DIPs, inside the arrow box.
			inner := ui.Rect{X: box.X + 3, Y: box.Y + 3, W: 10, H: 10}
			if p == Mac {
				if _, hi := collLuminanceRange(tester, inner); hi > luminance(fill)+0.02 {
					t.Fatalf("%s: the chevron of a chosen row without the focus is white", name)
				}
			}
			tester.ClickAt(box.X+8, box.Y+8)
			tester.Frame()
			if !s.Open.Has("src") || !tester.HasText("main.go") {
				t.Fatalf("%s: a click on the chevron did not open src: %v", name, tester.Texts())
			}
			box = arrow()
			tester.ClickAt(box.X+8, box.Y+8)
			tester.Frame()
			if s.Open.Has("src") || tester.HasText("main.go") {
				t.Fatalf("%s: a second click on the chevron left src open", name)
			}
			box = arrow()
			tester.ClickAtWith(ui.Alt, box.X+8, box.Y+8)
			tester.Frame()
			if !s.Open.Has("src") || !s.Open.Has("ui") || !tester.HasText("view.go") {
				t.Fatalf("%s: an Option-click on the chevron did not open all inside: %v", name, tester.Texts())
			}
			collClick(t, tester, "go.mod")
			if sel != 4 || changes != 1 || !tester.Focused("Files") {
				t.Fatalf("%s: a click chose row %d with %d changes, focused %v; want row 4, 1 change, focused", name, sel, changes, tester.Focused("Files"))
			}
			collClick(t, tester, "src")
			if p == Mac {
				box := arrow()
				if _, hi := collLuminanceRange(tester, ui.Rect{X: box.X + 3, Y: box.Y + 3, W: 10, H: 10}); hi < 0.95 {
					t.Fatalf("%s: the chevron of the chosen row with the focus is not white (%.2f)", name, hi)
				}
			}
			back, forth := ui.KeyLeft, ui.KeyRight
			if rtl {
				back, forth = forth, back
			}
			tester.Key(0, back)
			tester.Frame()
			if s.Open.Has("src") {
				t.Fatalf("%s: the back arrow did not close src", name)
			}
			tester.Key(0, forth)
			tester.Frame()
			if !s.Open.Has("src") {
				t.Fatalf("%s: the forth arrow did not open src", name)
			}
		}
	}
}

// TestOutlineTableKeepsItsBehavior checks the drawn outline table on every platform, left to right and right to left: its chevron is drawn in MyGo's arrow box right before the first cell's content, where a click opens and closes the item, the chosen row's fill reaches over it, a click chooses a row, and Right opens it.
func TestOutlineTableKeepsItsBehavior(t *testing.T) {
	cols := []ui.TableColumn{{Title: "Name"}, {Title: "Size", Width: 60}}
	for _, p := range Platforms {
		for _, rtl := range []bool{false, true} {
			th := collThemeOf(p, false)
			sel := 1
			var s ui.OutlineState[string]
			s.List.Selected = &sel
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					OutlineTable(c, &s, cols, []string{"src", "go.mod"}, collChildren, func(k string, col int) {
						ui.Text(c, map[int]string{0: k, 1: "1 KB"}[col]).SingleLine()
					}).Size(300, 260).Label("Tree table")
				})
			}, 320, 280)
			if rtl {
				tester.SetDirection(ui.RightToLeft)
			}
			tester.Frame()
			name := fmt.Sprintf("%v rtl=%v", p, rtl)
			collClick(t, tester, "src")
			if sel != 0 || !tester.Focused("Tree table") {
				t.Fatalf("%s: a click chose row %d, focused %v", name, sel, tester.Focused("Tree table"))
			}
			r, _ := tester.Find("src")
			box := ui.Rect{X: r.X - th.Space(4), Y: r.Y + r.H/2 - 8, W: 16, H: 16}
			side := box.X + 1
			if rtl {
				box.X = r.X + r.W
				side = box.X + box.W - 1
			}
			if got, want := collAt(tester, side, box.Y+8), collChosenFill(p, th, true); !collNear(got, want, 2) {
				t.Fatalf("%s: the chosen row at its arrow box is %v, want %v", name, got, want)
			}
			tester.ClickAt(box.X+8, box.Y+8)
			tester.Frame()
			if !s.Open.Has("src") || !tester.HasText("main.go") {
				t.Fatalf("%s: a click on the chevron did not open src: %v", name, tester.Texts())
			}
			tester.ClickAt(box.X+8, box.Y+8)
			tester.Frame()
			if s.Open.Has("src") {
				t.Fatalf("%s: a second click on the chevron did not close src", name)
			}
			forth := ui.KeyRight
			if rtl {
				forth = ui.KeyLeft
			}
			tester.Key(0, forth)
			tester.Frame()
			if !s.Open.Has("src") {
				t.Fatalf("%s: the forth arrow did not open src", name)
			}
		}
	}
}

// TestGridViewKeepsItsBehavior checks the drawn grid on every platform: the chosen item is the platform's (AppKit's selection, WinUI's 2-DIP accent border, GNOME's tinted item) over MyGo's border, a click chooses an item and gives the grid the focus, the arrows, Home and End move the choice, Changed reports it, a grid choosing several draws every item chosen by Cmd-click and Shift-click, items under the pointer light up except on macOS, and a disabled grid chooses nothing.
func TestGridViewKeepsItsBehavior(t *testing.T) {
	for _, p := range Platforms {
		th := collThemeOf(p, false)
		item, changes := 1, 0
		grid := &ui.GridState{Selected: &item, Label: func(i int) string { return collColors[i] }}
		lead := -1
		var many ui.Selection[int]
		multi := &ui.GridState{Selected: &lead, Selection: &many, Label: func(i int) string { return "Many " + collColors[i] }}
		off := 0
		offGrid := &ui.GridState{Selected: &off, Label: func(i int) string { return "Off " + collColors[i] }}
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Column(c).Fill().Gap(8).Padding(8).Children(func() {
					if GridView(c, grid, len(collColors), 60, 40, func(i int) { ui.Text(c, strings.ToLower(collColors[i])) }).Height(60).Label("Grid").Changed() {
						changes++
					}
					GridView(c, multi, len(collColors), 60, 40, func(i int) { ui.Text(c, "m"+collColors[i]) }).Height(60).Label("Many")
					GridView(c, offGrid, len(collColors), 60, 40, func(i int) { ui.Text(c, "o"+collColors[i]) }).Height(60).Label("Off grid").Disabled(true)
				})
			})
		}, 360, 240)
		tester.Frame()
		// The middle of an item's top edge, inside MyGo's 2-DIP border.
		top := func(label string) ui.Color {
			r, ok := tester.Find(label)
			if !ok {
				t.Fatalf("%v: the item %q is missing", p, label)
			}
			return collAt(tester, r.X+r.W/2, r.Y+1)
		}
		chosen := func(focused bool) ui.Color {
			switch p {
			case Windows:
				return th.Accent
			case Linux:
				return adwListViewBg(th).Mix(th.Accent, 0.25)
			}
			return macListSelection(th, focused)
		}
		if got, want := top("Green"), chosen(false); !collNear(got, want, 2) {
			t.Fatalf("%v: the chosen item without the focus is %v at its edge, want %v", p, got, want)
		}
		collClick(t, tester, "Blue")
		if item != 2 || changes != 1 || !tester.Focused("Grid") {
			t.Fatalf("%v: a click chose item %d with %d changes, focused %v; want 2, 1 change, focused", p, item, changes, tester.Focused("Grid"))
		}
		if got, want := top("Blue"), chosen(true); !collNear(got, want, 2) {
			t.Fatalf("%v: the chosen item with the focus is %v at its edge, want %v", p, got, want)
		}
		for _, k := range []struct {
			key  ui.Key
			want int
		}{{ui.KeyRight, 3}, {ui.KeyHome, 0}, {ui.KeyEnd, 3}, {ui.KeyLeft, 2}} {
			tester.Key(0, k.key)
			tester.Frame()
			if item != k.want {
				t.Fatalf("%v: %v chose item %d, want %d", p, k.key, item, k.want)
			}
		}
		rest := top("Red")
		red, _ := tester.Find("Red")
		tester.Move(red.X+red.W/2, red.Y+red.H/2)
		tester.Frame()
		if hovered := top("Red"); (p == Mac) != (hovered == rest) {
			t.Fatalf("%v: an item under the pointer is %v, at rest %v; only macOS's do not light up", p, hovered, rest)
		}
		tester.Move(0, 0)
		tester.Frame()
		collClick(t, tester, "Many Red")
		if err := tester.ClickWith(ui.Cmd, "Many Gray"); err != nil {
			t.Fatal(err)
		}
		tester.Move(0, 0)
		tester.Frame()
		if many.Len() != 2 || !many.Has(0) || !many.Has(3) {
			t.Fatalf("%v: Cmd-click chose %d items", p, many.Len())
		}
		if a, b, other := top("Many Red"), top("Many Gray"), top("Many Green"); !collNear(a, chosen(true), 2) || !collNear(b, chosen(true), 2) || collNear(other, chosen(true), 2) {
			t.Fatalf("%v: the items chosen together are %v and %v, one not chosen %v; want %v", p, a, b, other, chosen(true))
		}
		if err := tester.ClickWith(ui.Shift, "Many Blue"); err != nil {
			t.Fatal(err)
		}
		tester.Move(0, 0)
		tester.Frame()
		if many.Len() != 2 || !many.Has(2) || !many.Has(3) || !collNear(top("Many Blue"), chosen(true), 2) {
			t.Fatalf("%v: Shift-click from Gray to Blue chose %d items", p, many.Len())
		}
		if err := tester.Click("Off Blue"); err != nil {
			t.Fatal(err)
		}
		tester.Frame()
		if off != 0 {
			t.Fatalf("%v: a click chose item %d of a disabled grid", p, off)
		}
	}
}

// TestTreeKeepsItsBehavior checks the drawn tree on every platform, left to right and right to left: its items are indented by the platform's step, the item marked Selected is filled as the platform fills a chosen row (AppKit's gray selection without the focus, its accent one with it), a click chooses an item and gives it the focus, Up moves it, the arrows close and open an item, and a click on an arrow opens and closes it.
func TestTreeKeepsItsBehavior(t *testing.T) {
	for _, p := range Platforms {
		for _, rtl := range []bool{false, true} {
			th := collThemeOf(p, false)
			library, music := true, true
			file := "Podcasts"
			leaf := func(c *ui.Context, label string) {
				item := TreeItem(c, label, nil, nil)
				if item.Clicked() {
					file = label
				}
				item.Selected(file == label)
			}
			tester := ui.NewTester(func(c *ui.Context) {
				Preview(c, p, func() {
					Use(c)
					ui.Column(c).Fill().Padding(8).Children(func() {
						Tree(c, func() {
							TreeItem(c, "Library", &library, func() {
								TreeItem(c, "Music", &music, func() { leaf(c, "Songs.m4a") })
								leaf(c, "Podcasts")
							})
							leaf(c, "Movies")
						}).Width(260).Label("Media")
					})
				})
			}, 300, 260)
			if rtl {
				tester.SetDirection(ui.RightToLeft)
			}
			tester.Frame()
			name := fmt.Sprintf("%v rtl=%v", p, rtl)
			find := func(s string) ui.Rect {
				r, ok := tester.Find(s)
				if !ok {
					t.Fatalf("%s: %q is missing", name, s)
				}
				return r
			}
			// Where a label starts, from the tree's start.
			start := func(s string) float32 {
				if r := find(s); rtl {
					return -(r.X + r.W)
				} else {
					return r.X
				}
			}
			step := collIndent(p)
			if d1, d2 := start("Music")-start("Library"), start("Songs.m4a")-start("Music"); d1 < step-0.5 || d1 > step+0.5 || d2 < step-0.5 || d2 > step+0.5 {
				t.Fatalf("%s: the levels are indented %v and %v, want %v", name, d1, d2, step)
			}
			gap := map[Platform]float32{Mac: 0, Windows: 12, Linux: 4}[p]
			// The middle of an item's arrow box, empty for a leaf.
			arrow := func(s string) (float32, float32) {
				r := find(s)
				if rtl {
					return r.X + r.W + gap + 8, r.Y + r.H/2
				}
				return r.X - gap - 8, r.Y + r.H/2
			}
			if ax, ay := arrow("Podcasts"); !collNear(collAt(tester, ax, ay), collChosenFill(p, th, false), 2) {
				got, want := collAt(tester, ax, ay), collChosenFill(p, th, false)
				t.Fatalf("%s: the item marked Selected is %v, want %v", name, got, want)
			}
			collClick(t, tester, "Movies")
			if file != "Movies" || !tester.Focused("Movies") {
				t.Fatalf("%s: a click chose %q, focused %v", name, file, tester.Focused("Movies"))
			}
			if ax, ay := arrow("Movies"); !collNear(collAt(tester, ax, ay), collChosenFill(p, th, true), 2) {
				got, want := collAt(tester, ax, ay), collChosenFill(p, th, true)
				t.Fatalf("%s: the chosen item with the focus is %v, want %v", name, got, want)
			}
			for _, want := range []string{"Podcasts", "Songs.m4a", "Music"} {
				tester.Key(0, ui.KeyUp)
				tester.Frame()
				if !tester.Focused(want) {
					t.Fatalf("%s: Up did not move the focus to %s", name, want)
				}
			}
			back, forth := ui.KeyLeft, ui.KeyRight
			if rtl {
				back, forth = forth, back
			}
			tester.Key(0, back)
			tester.Frame()
			if music || tester.HasText("Songs.m4a") {
				t.Fatalf("%s: the back arrow did not close Music", name)
			}
			tester.Key(0, forth)
			tester.Frame()
			if !music {
				t.Fatalf("%s: the forth arrow did not open Music", name)
			}
			tester.ClickAt(arrow("Music"))
			tester.Frame()
			if music {
				t.Fatalf("%s: a click on Music's arrow did not close it", name)
			}
		}
	}
}

// collEditing reports whether a text input has the keyboard focus, as an editable text being edited does.
func collEditing(tester *ui.Tester) bool {
	active, _ := tester.TextInputActive()
	return active
}

// TestEditableTextDrawsThePlatformsFrame checks the drawn editable text on every platform: Enter edits it, the input's edge is not MyGo's accent border but the platform's frame (a focused TextBox's edge above its accent line on Windows, the input's own fill on GNOME, AppKit's focus ring on macOS), Escape goes back, a double click edits it and Enter keeps what was typed, which Changed reports.
func TestEditableTextDrawsThePlatformsFrame(t *testing.T) {
	for _, p := range Platforms {
		th := collThemeOf(p, false)
		name, changed := "Draft.txt", 0
		tester := ui.NewTester(func(c *ui.Context) {
			Preview(c, p, func() {
				Use(c)
				ui.Column(c).Fill().Padding(20).Gap(12).AlignItems(ui.Start).Children(func() {
					if EditableText(c, &name).Label("Name").Width(200).Changed() {
						changed++
					}
					Button(c, "Other")
				})
			})
		}, 300, 160)
		tester.Frame()
		collClick(t, tester, "Draft.txt")
		tester.Key(0, ui.KeyEnter)
		tester.Frame()
		if !collEditing(tester) {
			t.Fatalf("%v: Enter did not edit the text", p)
		}
		box, _ := tester.Find("Name")
		mid := box.X + box.W/2
		// MyGo's input is the text's box grown by 3 DIPs above and below; its 1-DIP border is the row of pixels at its edge.
		edge, bottom := collAt(tester, mid, box.Y-2.5), collAt(tester, mid, box.Y+box.H+2.5)
		if collNear(edge, th.Accent, 8) {
			t.Fatalf("%v: the input's edge is MyGo's accent border %v", p, edge)
		}
		switch p {
		case Windows:
			if !collNear(edge, th.Border, 2) || !collNear(bottom, th.Accent, 2) {
				t.Fatalf("%v: the input's edge is %v and its bottom %v, want %v and the accent %v", p, edge, bottom, th.Border, th.Accent)
			}
		case Linux:
			if !collNear(edge, th.Background, 2) {
				t.Fatalf("%v: the input's edge is %v, want its fill %v", p, edge, th.Background)
			}
		default:
			if collNear(edge, th.Background, 8) {
				t.Fatalf("%v: the input shows no focus ring at its edge", p)
			}
		}
		tester.Key(0, ui.KeyEscape)
		tester.Frame()
		if name != "Draft.txt" || changed != 0 || collEditing(tester) {
			t.Fatalf("%v: Escape left %q, %d changes", p, name, changed)
		}
		r, _ := tester.Find("Draft.txt")
		tester.ClickAt(r.X+4, r.Y+r.H/2)
		tester.ClickAt(r.X+4, r.Y+r.H/2)
		tester.Frame()
		tester.Type("Plan")
		tester.Frame()
		tester.Key(0, ui.KeyEnter)
		tester.Frame()
		if name != "Plan.txt" || changed != 1 {
			t.Fatalf("%v: a double click and Enter gave %q with %d changes; want Plan.txt and 1", p, name, changed)
		}
	}
}
