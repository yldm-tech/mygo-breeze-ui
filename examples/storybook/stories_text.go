package main

import (
	"strings"

	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, textStories) }

// textStyles are the steps of the type ramp, by the names the style control shows, with what the Ramp story shows in them, in its order.
var textStyles = []struct {
	name, sample string
	style        breeze.TextStyle
}{
	{"LargeTitle", "Welcome", breeze.TextLargeTitle},
	{"Title", "Settings", breeze.TextTitle},
	{"Subtitle", "Appearance", breeze.TextSubtitle},
	{"Heading", "Theme", breeze.TextHeading},
	{"Body", "Choose how windows look.", breeze.TextBody},
	{"Caption", "Applies to every window.", breeze.TextCaption},
}

// textStyleArg returns the step the style control names, with TextSecondary added when the secondary control is on, and the Go that passes it after the text: nothing for the body alone, which is Text's default.
func textStyleArg(a *args) (breeze.TextStyle, string) {
	style, code := breeze.TextBody, ""
	for _, s := range textStyles {
		if s.name == *a.text["style"] && s.style != breeze.TextBody {
			style, code = s.style, ", breeze.Text"+s.name
		}
	}
	if *a.on["secondary"] {
		style |= breeze.TextSecondary
		if code == "" {
			code = ", breeze.TextSecondary"
		} else {
			code += "|breeze.TextSecondary"
		}
	}
	return style, code
}

// textStories are the stories of the widgets of widgets_text.go.
func textStories() []*story {
	var names []string
	for _, s := range textStyles {
		names = append(names, s.name)
	}
	return []*story{
		{component: "Text", name: "Default", category: "Text",
			controls: []control{choiceArg("style", "Title", names...), toggleArg("secondary", false), textArg("text", "Settings")},
			render: func(c *ui.Context, a *args) {
				style, _ := textStyleArg(a)
				breeze.Text(c, *a.text["text"], style).MaxWidth(520)
			},
			code: func(a *args) string {
				_, style := textStyleArg(a)
				return "breeze.Text(c, " + q(*a.text["text"]) + style + ")"
			}},
		{component: "Text", name: "Ramp", category: "Text",
			controls: []control{toggleArg("secondary", false)},
			render: func(c *ui.Context, a *args) {
				ui.Column(c).Gap(8).AlignItems(ui.Start).Children(func() {
					for _, s := range textStyles {
						style := s.style
						if *a.on["secondary"] {
							style |= breeze.TextSecondary
						}
						breeze.Text(c, s.sample, style)
					}
				})
			},
			code: func(a *args) string {
				secondary := ""
				if *a.on["secondary"] {
					secondary = "|breeze.TextSecondary"
				}
				var lines []string
				for _, s := range textStyles {
					lines = append(lines, "breeze.Text(c, "+q(s.sample)+", breeze.Text"+s.name+secondary+")")
				}
				return strings.Join(lines, "\n")
			}},
	}
}
