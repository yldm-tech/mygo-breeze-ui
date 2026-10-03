// Breeze UI is a small native gallery for the Breeze UI design system for MyGo.
package main

import (
	"fmt"
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

type gallery struct {
	page    string
	count   int
	checked bool
	name    string
	open    bool
}

func (g *gallery) view(c *ui.Context) {
	p := breeze.Palette(c)
	breeze.Row(c, breeze.Fill, breeze.Bg(p.Background)).Children(func() {
		breeze.Column(c, breeze.W(190), breeze.Fill, breeze.P(4), breeze.Gap(2), breeze.Bg(p.Surface)).Children(func() {
			breeze.Text(c, "Breeze UI", breeze.TextSize(breeze.TextLG), breeze.FontWeight(700), breeze.Muted(c)).Padding(breeze.Space(1), breeze.Space(2))
			for _, page := range []string{"Overview", "Forms", "Overlays"} {
				item := breeze.Row(c, breeze.Py(2), breeze.Px(2), breeze.Gap(2), breeze.Rounded(breeze.RadiusMD), breeze.Shrink(0)).Focusable()
				if g.page == page {
					item.Background(p.Accent).TextColor(p.AccentText)
				} else if item.Hovered() {
					item.Background(p.Hover)
				}
				if item.Clicked() {
					g.page = page
				}
				item.Children(func() { breeze.Text(c, page) })
			}
		})
		breeze.Scroll(c, breeze.Fill, breeze.Grow(1), breeze.P(8), breeze.Gap(5)).Children(func() {
			breeze.Text(c, g.page, breeze.TextSize(breeze.Text2XL), breeze.FontWeight(700))
			switch g.page {
			case "Forms":
				g.forms(c)
			case "Overlays":
				g.overlays(c)
			default:
				g.overview(c)
			}
		})
	})
}

func (g *gallery) overview(c *ui.Context) {
	p := breeze.Palette(c)
	breeze.Text(c, "Typed utility styles and accessible components, drawn by MyGo without HTML or CSS.", breeze.Muted(c))
	breeze.Row(c, breeze.Gap(4), breeze.Wrap).Children(func() {
		breeze.Card(c, func() {
			breeze.Text(c, "Counter", breeze.Muted(c))
			breeze.Text(c, fmt.Sprint(g.count), breeze.TextSize(breeze.Text2XL), breeze.FontWeight(700))
			breeze.Row(c, breeze.Gap(2)).Children(func() {
				if breeze.Button(c, "−", breeze.ButtonSecondary).Clicked() {
					g.count--
				}
				if breeze.Button(c, "Increment", breeze.ButtonPrimary).Clicked() {
					g.count++
				}
			})
		})
		breeze.Card(c, func() {
			breeze.Text(c, "Status", breeze.Muted(c))
			breeze.Badge(c, "Ready", breeze.BadgeSuccess)
			breeze.Text(c, "The design tokens follow the active appearance.", breeze.Muted(c))
			_ = p
		})
	})
}

func (g *gallery) forms(c *ui.Context) {
	breeze.Card(c, func() {
		breeze.Text(c, "Profile", breeze.TextSize(breeze.TextLG), breeze.FontWeight(600))
		breeze.Input(c, &g.name, "Your name")
		breeze.Checkbox(c, &g.checked, "Keep me signed in")
		if breeze.Button(c, "Save", breeze.ButtonPrimary).Clicked() {
			g.open = true
		}
	})
}

func (g *gallery) overlays(c *ui.Context) {
	breeze.Card(c, func() {
		breeze.Text(c, "Dialog and status", breeze.TextSize(breeze.TextLG), breeze.FontWeight(600))
		if breeze.Button(c, "Open dialog", breeze.ButtonSecondary).Clicked() {
			g.open = true
		}
		if g.open {
			breeze.Dialog(c, &g.open, "Hello", func() {
				breeze.Text(c, "This dialog uses MyGo's accessible modal base.", breeze.Muted(c))
			}, func() {
				if breeze.Button(c, "Close", breeze.ButtonPrimary).Clicked() {
					g.open = false
				}
			})
		}
	})
}

func main() {
	g := &gallery{page: "Overview"}
	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title: "Breeze UI for MyGo",
			Width: 920, Height: 620,
			Content: ui.View(g.view),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
