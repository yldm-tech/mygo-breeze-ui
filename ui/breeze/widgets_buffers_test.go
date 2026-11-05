package breeze

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func init() {
	for _, n := range []string{"TextInputBuffer", "TextAreaBuffer"} {
		drawnWidgets[n] = true
	}
}

// bufferField builds the text field of kind (singleLine or multiLine) under test: breeze's TextInputBuffer or TextAreaBuffer when buffer is set, else breeze's TextInput or TextArea, both showing text.
func bufferField(c *ui.Context, kind fieldKind, buffer bool, text *string, doc *ui.TextBuffer) ui.Element {
	switch {
	case kind == multiLine && buffer:
		return TextAreaBuffer(c, doc)
	case kind == multiLine:
		return TextArea(c, text)
	case buffer:
		return TextInputBuffer(c, doc)
	default:
		return TextInput(c, text)
	}
}

// TestBufferInputsLookAsTextFields checks that TextInputBuffer and TextAreaBuffer draw exactly what TextInput and TextArea draw on every platform, light and dark, empty with a placeholder, holding text, focused and disabled: the two frames are the same pixel for pixel. It also checks that they are drawn by Breeze UI, not left as package ui's.
func TestBufferInputsLookAsTextFields(t *testing.T) {
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			for _, kind := range []fieldKind{singleLine, multiLine} {
				for _, state := range []string{"placeholder", "text", "focused", "disabled"} {
					shot := func(buffer, plain bool) *image.RGBA {
						text := "Notes for the meeting"
						if state == "placeholder" {
							text = ""
						}
						doc := ui.NewTextBuffer(text)
						tester := ui.NewTester(func(c *ui.Context) {
							Preview(c, p, func() {
								Use(c)
								ui.Column(c).Fill().Padding(20).AlignItems(ui.Start).Children(func() {
									var f ui.Element
									switch {
									case plain && kind == multiLine:
										f = ui.TextAreaBuffer(c, doc)
									case plain:
										f = ui.TextInputBuffer(c, doc)
									default:
										f = bufferField(c, kind, buffer, &text, doc)
									}
									f.Placeholder("Name").Width(220).Label("probe")
									if kind == multiLine {
										f.Height(90)
									}
									if state == "focused" {
										f.AutoFocus()
									}
									f.Disabled(state == "disabled")
								})
							})
						}, 300, 150)
						tester.SetScale(2)
						tester.SetDark(dark)
						tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
						tester.Frame()
						tester.Frame()
						return tester.Image()
					}
					name := fmt.Sprintf("%v dark=%v kind=%d %s", p, dark, kind, state)
					want, got := shot(false, false), shot(true, false)
					if !bytes.Equal(want.Pix, got.Pix) {
						t.Errorf("%s: the buffer input does not draw as the text field (%d pixels differ)", name, bufferPixelsDiffer(want, got))
					}
					if bytes.Equal(got.Pix, shot(false, true).Pix) {
						t.Errorf("%s: the buffer input draws as package ui's", name)
					}
				}
			}
		}
	}
}

// bufferPixelsDiffer counts the pixels a and b differ in.
func bufferPixelsDiffer(a, b *image.RGBA) int {
	n := 0
	for i := 0; i+3 < len(a.Pix) && i+3 < len(b.Pix); i += 4 {
		if !bytes.Equal(a.Pix[i:i+4], b.Pix[i:i+4]) {
			n++
		}
	}
	return n
}

// TestBufferInputsKeepUIBehavior checks that TextInputBuffer and TextAreaBuffer edit their buffer as ui's do on every platform: typing, undo and redo, selecting all and replacing it, an input method's composition, which stays out of the buffer until it is committed, and, in the area, a new line.
func TestBufferInputsKeepUIBehavior(t *testing.T) {
	for _, p := range Platforms {
		for _, kind := range []fieldKind{singleLine, multiLine} {
			doc := ui.NewTextBuffer("Hello")
			tester := formsTester(p, ui.LeftToRight, 400, 200, func(c *ui.Context) {
				ui.Column(c).Fill().Padding(20).Children(func() {
					bufferField(c, kind, true, nil, doc).Label("probe")
				})
			})
			name := fmt.Sprintf("%v kind=%d", p, kind)
			r := formsFind(t, tester, "probe")
			tester.ClickAt(r.X+r.W/2, r.Y+r.H/2)
			if !tester.Focused("probe") {
				t.Fatalf("%s: a click does not focus the input", name)
			}
			tester.Key(0, ui.KeyEnd)
			tester.Type(" world")
			if got := doc.String(); got != "Hello world" {
				t.Fatalf("%s: typing gives %q", name, got)
			}
			tester.Key(ui.Cmd, ui.KeyZ)
			if got := doc.String(); got != "Hello" {
				t.Fatalf("%s: undo gives %q", name, got)
			}
			tester.Key(ui.Cmd|ui.Shift, ui.KeyZ)
			if got := doc.String(); got != "Hello world" {
				t.Fatalf("%s: redo gives %q", name, got)
			}
			tester.Key(ui.Cmd, ui.KeyA)
			tester.Type("X")
			if got := doc.String(); got != "X" {
				t.Fatalf("%s: typing over the selection gives %q", name, got)
			}
			tester.Compose("ni", 2)
			if _, active := tester.TextCaret(); !active {
				t.Fatalf("%s: no input method composes in the input", name)
			}
			if got := doc.String(); got != "X" {
				t.Fatalf("%s: the composition went into the buffer: %q", name, got)
			}
			tester.Compose("", 0)
			tester.Type("你")
			if got := doc.String(); got != "X你" {
				t.Fatalf("%s: committing the composition gives %q", name, got)
			}
			if kind == multiLine {
				tester.Key(0, ui.KeyEnter)
				tester.Type("two")
				if got := doc.String(); got != "X你\ntwo" || doc.LineCount() != 2 {
					t.Fatalf("%s: Enter and typing give %q", name, got)
				}
			}
		}
	}
}

// TestBufferInputsKeepUILayout checks that TextInputBuffer and TextAreaBuffer take the place TextInput and TextArea take, and the width package ui's buffer inputs take, on every platform, left to right and right to left: grown in a row, stretched by a column, and, for the area, grown down a column.
func TestBufferInputsKeepUILayout(t *testing.T) {
	type build func(c *ui.Context, doc *ui.TextBuffer, text *string) ui.Element
	for _, kind := range []fieldKind{singleLine, multiLine} {
		mine := func(c *ui.Context, doc *ui.TextBuffer, _ *string) ui.Element {
			return bufferField(c, kind, true, nil, doc)
		}
		field := func(c *ui.Context, _ *ui.TextBuffer, text *string) ui.Element {
			return bufferField(c, kind, false, text, nil)
		}
		plain := func(c *ui.Context, doc *ui.TextBuffer, _ *string) ui.Element {
			if kind == multiLine {
				return ui.TextAreaBuffer(c, doc)
			}
			return ui.TextInputBuffer(c, doc)
		}
		layouts := map[string]func(c *ui.Context, e func() ui.Element){
			"grown in a row": func(c *ui.Context, e func() ui.Element) {
				ui.Row(c).Width(700).Padding(10).Gap(10).Children(func() {
					ui.Text(c, "before")
					e().Grow(1).Label("probe")
					ui.Text(c, "after")
				})
			},
			"stretched by a column": func(c *ui.Context, e func() ui.Element) {
				ui.Column(c).Fill().Padding(16).Gap(8).Children(func() {
					ui.Text(c, "above")
					e().Label("probe")
				})
			},
			"grown down a column": func(c *ui.Context, e func() ui.Element) {
				ui.Column(c).Width(300).Height(400).Padding(16).Gap(8).Children(func() {
					ui.Text(c, "above")
					e().Grow(1).Label("probe")
					ui.Text(c, "below")
				})
			},
		}
		for _, p := range Platforms {
			for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
				for name, layout := range layouts {
					if name == "grown down a column" && kind != multiLine {
						continue
					}
					box := func(b build) ui.Rect {
						doc, text := ui.NewTextBuffer("Hello"), "Hello"
						tester := formsTester(p, dir, 760, 500, func(c *ui.Context) {
							layout(c, func() ui.Element { return b(c, doc, &text) })
						})
						return formsFind(t, tester, "probe")
					}
					what := fmt.Sprintf("%v %v kind=%d %s", p, dir, kind, name)
					got, want, base := box(mine), box(field), box(plain)
					if got != want {
						t.Errorf("%s: at %v, the text field at %v", what, got, want)
					}
					if d := got.W - base.W; d > 1 || d < -1 || got.X-base.X > 1 || base.X-got.X > 1 {
						t.Errorf("%s: at %v, package ui's at %v", what, got, base)
					}
					if name == "grown down a column" && (got.H-base.H > 1 || base.H-got.H > 1) {
						t.Errorf("%s: %v high, package ui's %v", what, got.H, base.H)
					}
				}
			}
		}
	}
}

// scrollKind names a scroll container: breeze's or package ui's Scroll, ScrollHorizontal or ScrollBoth.
type scrollKind string

// newScroll builds the scroll container of kind, breeze's when mine is set, with content overflowing it the ways it scrolls.
func newScroll(c *ui.Context, kind scrollKind, mine bool) ui.Element {
	var s ui.Element
	switch {
	case kind == "Scroll" && mine:
		s = Scroll(c)
	case kind == "Scroll":
		s = ui.Scroll(c)
	case kind == "ScrollHorizontal" && mine:
		s = ScrollHorizontal(c)
	case kind == "ScrollHorizontal":
		s = ui.ScrollHorizontal(c)
	case mine:
		s = ScrollBoth(c)
	default:
		s = ui.ScrollBoth(c)
	}
	return s.Children(func() {
		switch kind {
		case "Scroll":
			for i := range 30 {
				ui.Text(c, fmt.Sprintf("Row %d", i)).Height(30)
			}
		case "ScrollHorizontal":
			for i := range 30 {
				ui.Text(c, fmt.Sprintf("Item %d", i)).Width(90).Shrink(0)
			}
		default:
			ui.Column(c).Width(900).Shrink(0).Children(func() {
				for i := range 30 {
					ui.Text(c, fmt.Sprintf("Line %d of a canvas wider and taller than its view", i)).Height(30).SingleLine()
				}
			})
		}
	})
}

var scrollKinds = []scrollKind{"Scroll", "ScrollHorizontal", "ScrollBoth"}

// TestScrollKeepsUIBehaviorAndLayout checks that Scroll, ScrollHorizontal and ScrollBoth are package ui's on every platform, left to right and right to left: they take its place grown in a row and stretched by a column, and the wheel and the keys scroll them as far as they scroll ui's.
func TestScrollKeepsUIBehaviorAndLayout(t *testing.T) {
	type result struct {
		row, column ui.Rect
		wheel, keys [2]float32 // X and Y of the container's ScrollState
	}
	for _, p := range Platforms {
		for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			for _, kind := range scrollKinds {
				run := func(mine bool) result {
					var res result
					var wheel, keys ui.ScrollState
					tester := formsTester(p, dir, 760, 500, func(c *ui.Context) {
						ui.Row(c).Width(700).Height(200).Padding(10).Gap(10).Children(func() {
							ui.Text(c, "before")
							newScroll(c, kind, mine).Grow(1).Label("probe")
							ui.Text(c, "after")
						})
					})
					res.row = formsFind(t, tester, "probe")
					tester = formsTester(p, dir, 760, 500, func(c *ui.Context) {
						ui.Column(c).Fill().Padding(16).Gap(8).Children(func() {
							ui.Text(c, "above")
							newScroll(c, kind, mine).Height(200).Label("probe").TrackScroll(&wheel)
						})
					})
					res.column = formsFind(t, tester, "probe")
					r := res.column
					tester.Move(r.X+r.W/2, r.Y+r.H/2)
					// A wheel's dx moves the view to the right, which right to left is back toward the content's start (scrollBy in input.go).
					dx := map[ui.Direction]float32{ui.LeftToRight: 120, ui.RightToLeft: -120}[dir]
					tester.Wheel(r.X+r.W/2, r.Y+r.H/2, dx, 120)
					tester.Advance(2e9)
					tester.Frame()
					res.wheel = [2]float32{wheel.X, wheel.Y}
					tester = formsTester(p, dir, 760, 500, func(c *ui.Context) {
						ui.Column(c).Fill().Padding(16).Gap(8).Children(func() {
							newScroll(c, kind, mine).Height(200).Label("probe").TrackScroll(&keys)
						})
					})
					r = formsFind(t, tester, "probe")
					tester.Move(r.X+r.W/2, r.Y+r.H/2)
					tester.Key(0, ui.KeyPageDown)
					// Right to left, Left scrolls toward the content's end, as its start is at the right.
					tester.Key(0, map[ui.Direction]ui.Key{ui.LeftToRight: ui.KeyRight, ui.RightToLeft: ui.KeyLeft}[dir])
					tester.Frame()
					res.keys = [2]float32{keys.X, keys.Y}
					return res
				}
				what := fmt.Sprintf("%v %v %s", p, dir, kind)
				got, want := run(true), run(false)
				if got.row != want.row || got.column != want.column {
					t.Errorf("%s: at %v and %v, package ui's at %v and %v", what, got.row, got.column, want.row, want.column)
				}
				if got.wheel != want.wheel {
					t.Errorf("%s: the wheel scrolls to %v, package ui's to %v", what, got.wheel, want.wheel)
				}
				if got.keys != want.keys {
					t.Errorf("%s: the keys scroll to %v, package ui's to %v", what, got.keys, want.keys)
				}
				if got.wheel == [2]float32{} || got.keys == [2]float32{} {
					t.Errorf("%s: the wheel scrolled to %v and the keys to %v: the test expects both to scroll", what, got.wheel, got.keys)
				}
			}
		}
	}
}

// TestScrollBarTakesThePlatformsColor checks that the thumb MyGo paints in Scroll, while the pointer is over it, is the platform theme's Scrollbar color over the window's background, at the side MyGo puts it: the end of the container left to right, its start right to left. MyGo's thumb is ScrollbarWidth wide (6 when the theme leaves it 0), 3 DIPs in from the edge (scrollBars in paint.go).
func TestScrollBarTakesThePlatformsColor(t *testing.T) {
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			for _, dir := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
				var th ui.Theme
				tester := ui.NewTester(func(c *ui.Context) {
					Preview(c, p, func() {
						Use(c)
						th = *c.Theme()
						ui.Column(c).Fill().Padding(16).Children(func() {
							newScroll(c, "Scroll", true).Height(240).Label("probe")
						})
					})
				}, 400, 300)
				tester.SetDark(dark)
				tester.SetDirection(dir)
				tester.Frame()
				r := formsFind(t, tester, "probe")
				tester.Move(r.X+r.W/2, r.Y+r.H/2)
				tester.Frame()
				width := th.ScrollbarWidth
				if width <= 0 {
					width = 6
				}
				x := r.X + r.W - 3 - width/2
				if dir == ui.RightToLeft {
					x = r.X + 3 + width/2
				}
				y := r.Y + 2 + 12 // within the thumb, which is at least 24 long and starts 2 DIPs down the track at the top
				want := th.Background.Mix(ui.RGBA(th.Scrollbar.R, th.Scrollbar.G, th.Scrollbar.B, 1), float32(th.Scrollbar.A)/255)
				got := tester.Image().RGBAAt(int(x), int(y))
				if !formsNear(got, want, 3) {
					t.Errorf("%v dark=%v %v: the thumb at %v,%v is %v, want the theme's Scrollbar %v over the background, %v", p, dark, dir, x, y, got, th.Scrollbar, want)
				}
			}
		}
	}
}

// buffersShotView lays out the buffer inputs beside the text fields they look like, focused, with a placeholder and disabled, and the three scroll containers, for screenshots. With plain it builds package ui's buffer inputs and scroll containers in their place, as an app had to before Breeze UI wrapped them.
func buffersShotView(p Platform, plain bool, doc, notes, empty, off *ui.TextBuffer) func(c *ui.Context) {
	name, area := "Ada Lovelace", "Notes for the meeting"
	input := func(c *ui.Context, b *ui.TextBuffer) ui.Element {
		if plain {
			return ui.TextInputBuffer(c, b)
		}
		return TextInputBuffer(c, b)
	}
	textArea := func(c *ui.Context, b *ui.TextBuffer) ui.Element {
		if plain {
			return ui.TextAreaBuffer(c, b)
		}
		return TextAreaBuffer(c, b)
	}
	return func(c *ui.Context) {
		Preview(c, p, func() {
			Use(c)
			t := c.Theme()
			ui.Row(c).Fill().Padding(16).Gap(24).AlignItems(ui.Start).Children(func() {
				ui.Column(c).Width(260).Gap(10).Children(func() {
					ui.Text(c, "TextInput").TextColor(t.TextMuted)
					TextInput(c, &name)
					ui.Text(c, "TextInputBuffer").TextColor(t.TextMuted)
					input(c, doc).AutoFocus()
					input(c, empty).Placeholder("Command")
					input(c, off).Disabled(true)
					ui.Text(c, "TextArea").TextColor(t.TextMuted)
					TextArea(c, &area).Height(70)
					ui.Text(c, "TextAreaBuffer").TextColor(t.TextMuted)
					textArea(c, notes).Height(70)
				})
				ui.Column(c).Width(260).Gap(10).Children(func() {
					for _, kind := range scrollKinds {
						ui.Text(c, string(kind)).TextColor(t.TextMuted)
						var s ui.Element
						if plain {
							s = newScroll(c, kind, false)
						} else {
							s = newScroll(c, kind, true)
						}
						s.Height(map[scrollKind]float32{"Scroll": 110, "ScrollHorizontal": 40, "ScrollBoth": 110}[kind]).Border(1, t.Border).Label(string(kind) + " view")
					}
				})
			})
		})
	}
}

// TestScreenshotsBuffers writes the buffer inputs and the scroll containers of every platform, light and dark, at 2x to the directory BREEZE_SCREENSHOTS names, once with the pointer over each scroll container so that MyGo shows its bars. Without it the test only builds the frames.
func TestScreenshotsBuffers(t *testing.T) {
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	for _, p := range Platforms {
		for _, dark := range []bool{false, true} {
			mode := map[bool]string{false: "light", true: "dark"}[dark]
			for _, kind := range scrollKinds {
				doc, notes := ui.NewTextBuffer("Ada Lovelace"), ui.NewTextBuffer("Notes for the meeting")
				empty, off := ui.NewTextBuffer(""), ui.NewTextBuffer("Disabled")
				tester := ui.NewTester(buffersShotView(p, false, doc, notes, empty, off), 600, 500)
				tester.SetScale(2)
				tester.SetDark(dark)
				tester.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
				tester.Frame()
				tester.Frame()
				for _, want := range []string{"TextInputBuffer", "TextAreaBuffer", "Row 0", "Item 0"} {
					if !tester.HasText(want) {
						t.Fatalf("%v: %q is missing", p, want)
					}
				}
				r := formsFind(t, tester, string(kind)+" view")
				tester.Move(r.X+r.W/2, r.Y+r.H/2)
				tester.Frame()
				writeShot(t, dir, fmt.Sprintf("breeze-buffers-%s-%s-%s.png", strings.ToLower(p.String()), mode, strings.ToLower(string(kind))), tester)
			}
		}
	}
}
