// Command storybook shows every Breeze UI widget as Storybook shows components: stories in a sidebar, each drawn on a canvas as macOS, Windows 11 or GNOME draws it, or as all three side by side, light or dark, with controls that change its arguments and the Go that builds it. It runs as a window, or in a browser compiled to WebAssembly:
//
//	go run ./examples/storybook
//	GOOS=js GOARCH=wasm go build -o website/public/storybook/app.wasm ./examples/storybook
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

// Platform choices of the toolbar, in order; the last shows the three side by side.
var platformNames = []string{"macOS", "Windows", "GNOME", "Side by side"}

type navEntry struct {
	category string
	story    *story
}

type book struct {
	stories  []*story
	selected string // the chosen story's id
	platform int    // an index of platformNames
	dark     int    // 0 light, 1 dark
	query    string
	open     map[string]*bool // the sidebar sections shown
	args     map[string]*args // each story's argument values, kept as stories change
	navRows  []navEntry
	navList  ui.ListState
	navIndex int
	navID    string
}

func newBook() *book {
	b := &book{stories: stories(), open: map[string]*bool{}, args: map[string]*args{}}
	b.selected = b.stories[0].id()
	for _, cat := range categories {
		shown := true
		b.open[cat] = &shown
	}
	return b
}

func (b *book) story() *story {
	for _, s := range b.stories {
		if s.id() == b.selected {
			return s
		}
	}
	return b.stories[0]
}

func (b *book) argsOf(s *story) *args {
	a, ok := b.args[s.id()]
	if !ok {
		a = newArgs(s.controls)
		b.args[s.id()] = a
	}
	return a
}

func (b *book) view(c *ui.Context) {
	c.SetTheme(storybookTheme())
	ui.Row(c).Fill().AlignItems(ui.Stretch).Children(func() {
		b.sidebar(c)
		ui.Column(c).Grow(1).MinWidth(0).Children(func() {
			s := b.story()
			b.toolbar(c, s)
			b.canvas(c, s)
			b.panel(c, s)
		})
	})
}

// storybookTheme keeps the application chrome in the same visual language as
// the website: a soft gray canvas, white panels, ink text and Breeze mint
// accents. The platform previews still use their own themes below.
func storybookTheme() *ui.Theme {
	t := *ui.LightTheme()
	t.Background = ui.Hex("#f4f5f7")
	t.Surface = ui.Hex("#ffffff")
	t.SurfaceHover = ui.Hex("#f8f9fb")
	t.SurfacePressed = ui.Hex("#e2e4e9")
	t.Border = ui.Hex("#e2e4e9")
	t.Text = ui.Hex("#15171c")
	t.TextMuted = ui.Hex("#555c69")
	t.Accent = ui.Hex("#0b7a63")
	t.AccentHover = ui.Hex("#075a49")
	t.AccentPressed = ui.Hex("#064b3e")
	t.AccentText = ui.Hex("#ffffff")
	t.Selection = ui.Hex("#a6e1cd")
	t.Focus = ui.RGBA(11, 122, 99, 0.35)
	t.Radius = 10
	t.Spacing = 4
	t.FontSize = 14
	return &t
}

func storybookLogo(c *ui.Context) ui.Element {
	return ui.Box(c).Size(30, 30).Role(ui.RoleImage).Label("Breeze UI").Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X + 2, Y: r.Y + 3, W: 18, H: 14}, ui.Hex("#a6e1cd"), 4.5)
		p.Fill(ui.Rect{X: r.X + 7, Y: r.Y + 8.5, W: 18, H: 14}, ui.Hex("#3fae8c"), 2)
		p.Fill(ui.Rect{X: r.X + 12, Y: r.Y + 14, W: 18, H: 14}, ui.Hex("#0b7a63"), 3.5)
		for _, x := range []float32{15.5, 18.5, 21.5} {
			p.Fill(ui.Rect{X: r.X + x - 0.9, Y: r.Y + 17.5 - 0.9, W: 1.8, H: 1.8}, ui.RGB(255, 255, 255), 0.9)
		}
	})
}

// sidebar lists the stories by category, filtered by the search field.
func (b *book) sidebar(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).Width(288).Shrink(0).Background(t.Surface).BorderWidth(0, 1, 0, 0).BorderColor(t.Border).Children(func() {
		ui.Column(c).Padding(22, 18, 14, 18).Gap(12).Children(func() {
			ui.Row(c).Gap(10).AlignItems(ui.Center).Children(func() {
				storybookLogo(c)
				ui.Column(c).Gap(0).Children(func() {
					ui.Text(c, "Breeze UI").FontSize(t.FontSize + 3).Bold()
					ui.Text(c, "Storybook").TextColor(t.TextMuted).FontSize(t.FontSize - 1)
				})
			})
			ui.Text(c, fmt.Sprintf("%d components · live previews", len(b.stories))).TextColor(t.TextMuted).FontSize(t.FontSize - 1)
			ui.SearchField(c, &b.query).Label("Find a component")
		})
		q := strings.ToLower(strings.TrimSpace(b.query))
		b.navRows = b.navRows[:0]
		for _, cat := range categories {
			var matches []*story
			for _, s := range b.stories {
				if s.category == cat && (q == "" || strings.Contains(strings.ToLower(s.component+" "+s.name), q)) {
					matches = append(matches, s)
				}
			}
			if len(matches) == 0 {
				continue
			}
			b.navRows = append(b.navRows, navEntry{category: cat})
			if shown := b.open[cat]; shown != nil && *shown {
				for _, s := range matches {
					b.navRows = append(b.navRows, navEntry{category: cat, story: s})
				}
			}
		}
		if b.navID != b.selected {
			b.navID = b.selected
			b.navIndex = -1
			for i, entry := range b.navRows {
				if entry.story != nil && entry.story.id() == b.selected {
					b.navIndex = i
					break
				}
			}
		}
		b.navList.Selected = &b.navIndex
		b.navList.RoundRows = true
		b.navList.Key = func(i int) any {
			entry := b.navRows[i]
			if entry.story == nil {
				return "category/" + entry.category
			}
			return entry.story.id()
		}
		b.navList.Header = func(i int) bool { return b.navRows[i].story == nil }
		b.navList.Label = func(i int) string {
			entry := b.navRows[i]
			if entry.story == nil {
				return entry.category
			}
			return entry.story.component + " " + entry.story.name
		}
		list := ui.List(c, &b.navList, len(b.navRows)).Rows(func(row ui.ListRow) {
			entry := b.navRows[row.Index]
			if entry.story == nil {
				head := ui.ButtonBase(row.Context).FillWidth().Padding(8, 10).Radius(8).AlignItems(ui.Center)
				head.Children(func() {
					arrow := "▾"
					if shown := b.open[entry.category]; shown == nil || !*shown {
						arrow = "▸"
					}
					ui.Text(row.Context, arrow).TextColor(t.TextMuted).Width(16)
					ui.Text(row.Context, entry.category).TextColor(t.TextMuted).Bold().Grow(1)
				})
				if head.Clicked() {
					*b.open[entry.category] = !*b.open[entry.category]
				}
				return
			}
			item := ui.Row(row.Context).FillWidth().Gap(8).Padding(7, 10).Radius(8).AlignItems(ui.Center)
			if row.Selected() {
				item.Background(t.SurfacePressed).TextColor(t.Text)
			}
			item.Children(func() {
				ui.Text(row.Context, entry.story.component+" · "+entry.story.name).Grow(1).SingleLine()
			})
		})
		list.Grow(1).Label("Stories")
		if list.Changed() && b.navIndex >= 0 && b.navIndex < len(b.navRows) && b.navRows[b.navIndex].story != nil {
			b.selected = b.navRows[b.navIndex].story.id()
			b.navID = b.selected
		}
	})
}

// toolbar names the story and chooses the platform and the appearance.
func (b *book) toolbar(c *ui.Context, s *story) {
	t := c.Theme()
	ui.Column(c).Padding(14, 22, 16, 22).Gap(12).Background(t.Surface).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
		ui.Row(c).Gap(14).AlignItems(ui.Center).Wrap().Children(func() {
			ui.Row(c).Grow(1).Gap(9).AlignItems(ui.Center).Children(func() {
				ui.Text(c, "STORY").Padding(4, 8).Radius(999).Background(t.Accent).TextColor(t.AccentText).FontSize(t.FontSize - 2).Bold()
				ui.Text(c, s.component).Bold().FontSize(t.FontSize + 4)
				ui.Text(c, "/").TextColor(t.Border)
				ui.Text(c, s.name).TextColor(t.TextMuted)
			})
			ui.Text(c, "LIVE").Padding(4, 8).Radius(999).Background(ui.RGBA(11, 122, 99, 0.12)).TextColor(t.Accent).FontSize(t.FontSize - 2).Bold()
		})
		ui.Row(c).Gap(10).AlignItems(ui.Center).Wrap().Children(func() {
			ui.Row(c).Gap(8).Padding(5, 8).Radius(10).Background(t.SurfaceHover).Border(1, t.Border).AlignItems(ui.Center).Children(func() {
				ui.Text(c, "Platform").TextColor(t.TextMuted).FontSize(t.FontSize - 1).Bold()
				ui.Segmented(c, &b.platform, platformNames...)
			})
			ui.Row(c).Gap(8).Padding(5, 8).Radius(10).Background(t.SurfaceHover).Border(1, t.Border).AlignItems(ui.Center).Children(func() {
				ui.Text(c, "Appearance").TextColor(t.TextMuted).FontSize(t.FontSize - 1).Bold()
				ui.Segmented(c, &b.dark, "Light", "Dark")
			})
		})
	})
}

// canvas draws the story as the chosen platform, or as each platform in a column of its own.
func (b *book) canvas(c *ui.Context, s *story) {
	t := c.Theme()
	a := b.argsOf(s)
	platforms := []breeze.Platform{breeze.Platforms[min(b.platform, 2)]}
	if b.platform == 3 {
		platforms = breeze.Platforms
	}
	appearance := ui.LightTheme()
	if b.dark == 1 {
		appearance = ui.DarkTheme()
	}
	ui.Row(c).Grow(1).MinHeight(260).AlignItems(ui.Stretch).Padding(20, 22).Gap(16).Background(t.Background).Children(func() {
		for i, p := range platforms {
			name := platformNames[i]
			if len(platforms) == 1 {
				name = platformNames[b.platform]
			}
			// The appearance's theme first, which breeze.Theme reads light or dark from, then the platform's.
			ui.Themed(c, appearance).Grow(1).Basis(0).MinWidth(0).Radius(16).Background(t.Surface).Border(1, t.Border).Shadow(0, 8, 24, 0, ui.RGBA(20, 24, 32, 0.06)).Clip().Children(func() {
				breeze.Preview(c, p, func() {
					ui.Themed(c, breeze.Theme(c, p)).Fill().Padding(20).Gap(16).Children(func() {
						if len(platforms) > 1 {
							ui.Text(c, name).Bold().TextColor(c.Theme().TextMuted)
						}
						ui.Column(c).Grow(1).Center().Gap(12).Children(func() { s.render(c, a) })
					})
				})
			})
		}
	})
}

// panel shows the story's controls and the Go that builds it.
func (b *book) panel(c *ui.Context, s *story) {
	t := c.Theme()
	a := b.argsOf(s)
	ui.Row(c).Height(250).Shrink(0).AlignItems(ui.Stretch).Background(t.Surface).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
		ui.Scroll(c).Grow(1).Basis(0).Padding(18, 22).Gap(12).Children(func() {
			ui.Text(c, "Controls").TextColor(t.Accent).Bold()
			if len(s.controls) == 0 {
				ui.Text(c, "This story has no controls.").TextColor(t.TextMuted)
			}
			ui.Form(c, func() {
				for _, ctl := range s.controls {
					ui.Field(c, ctl.name, func() { b.control(c, a, ctl) })
				}
			})
		})
		ui.Column(c).Grow(1).Basis(0).Padding(18, 22).Gap(12).Background(t.SurfaceHover).BorderWidth(0, 0, 0, 1).BorderColor(t.Border).Children(func() {
			ui.Row(c).AlignItems(ui.Center).Children(func() {
				ui.Text(c, "Go code").TextColor(t.Accent).Bold().Grow(1)
				if ui.Button(c, "Copy").Clicked() {
					c.WriteClipboard(s.code(a))
				}
			})
			ui.Text(c, s.code(a)).Font("monospace").Selectable().TextColor(t.Text)
		})
	})
}

// control builds the editor of one of the story's arguments.
func (b *book) control(c *ui.Context, a *args, ctl control) {
	switch ctl.kind {
	case textControl:
		ui.TextInput(c, a.text[ctl.name]).Width(240)
	case toggleControl:
		ui.Switch(c, a.on[ctl.name])
	case numberControl:
		v := a.num[ctl.name]
		ui.Row(c).Gap(10).AlignItems(ui.Center).Children(func() {
			ui.Slider(c, v, ctl.lo, ctl.hi).Width(180)
			ui.Text(c, formatNumber(*v, ctl.hi)).TextColor(c.Theme().TextMuted)
		})
	case choiceControl:
		ui.Select(c, a.text[ctl.name], ctl.options)
	}
}

func formatNumber(v, hi float64) string {
	if hi <= 1 {
		return fmt.Sprintf("%.2f", v)
	}
	return fmt.Sprintf("%.0f", v)
}

func main() {
	b := newBook()
	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title:   "Breeze UI Storybook",
			Width:   1280,
			Height:  820,
			Content: ui.View(b.view),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
