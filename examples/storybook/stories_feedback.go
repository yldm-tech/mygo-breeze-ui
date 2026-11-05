package main

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, feedbackStories) }

// feedbackText is what the find bar's story finds in.
const feedbackText = "The cat sat on the mat. A second cat came by, and the third cat stayed home."

// feedbackStories are the stories of the widgets of widgets_feedback.go.
func feedbackStories() []*story {
	return []*story{
		// ---- Text ----
		{component: "FindBar", name: "Default", category: "Text",
			controls: []control{toggleArg("open", true), textArg("query", "cat")},
			render: func(c *ui.Context, a *args) {
				query := a.text["query"]
				matches := 0
				if *query != "" {
					matches = strings.Count(strings.ToLower(feedbackText), strings.ToLower(*query))
				}
				current := feedbackFindMatch(c)
				ui.Column(c).FillWidth().MaxWidth(520).Gap(12).AlignItems(ui.Stretch).Children(func() {
					if !*a.on["open"] && breeze.Button(c, "Find…").AlignSelf(ui.Start).Clicked() {
						*a.on["open"] = true
					}
					breeze.FindBar(c, a.on["open"], query, matches, current)
					ui.Text(c, feedbackText)
					if matches > 0 {
						ui.Text(c, fmt.Sprintf("Showing match %d of %d", *current+1, matches)).TextColor(c.Theme().TextMuted)
					}
				})
			},
			code: func(*args) string {
				return "breeze.FindBar(c, &app.finding, &app.query, len(app.matches), &app.match)"
			}},

		// ---- Choices ----
		{component: "Rating", name: "Default", category: "Choices",
			controls: []control{numberArg("stars", 3, 0, 5), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				stars := int(*a.num["stars"] + 0.5)
				if disabled(a, breeze.Rating(c, &stars, 5).Label("Rating")).Changed() {
					*a.num["stars"] = float64(stars)
				}
			},
			code: func(*args) string { return "breeze.Rating(c, &app.stars, 5).Label(\"Rating\")" }},

		// ---- Feedback ----
		{component: "Meter", name: "Levels", category: "Feedback",
			controls: []control{numberArg("value", 60, 0, 100), numberArg("warning", 70, 0, 100), numberArg("critical", 90, 0, 100), toggleArg("disabled", false)},
			render: func(c *ui.Context, a *args) {
				levels := &ui.MeterLevels{Warning: *a.num["warning"], Critical: *a.num["critical"]}
				disabled(a, breeze.Meter(c, *a.num["value"], 0, 100, levels).Width(260).Label("Disk"))
			},
			code: func(a *args) string {
				return fmt.Sprintf("breeze.Meter(c, app.used, 0, 100, &ui.MeterLevels{Warning: %.0f, Critical: %.0f}).Label(\"Disk\")", *a.num["warning"], *a.num["critical"])
			}},
		{component: "Avatar", name: "Initials", category: "Feedback",
			controls: []control{textArg("name", "Ada Lovelace")},
			render: func(c *ui.Context, a *args) {
				ui.Row(c).Gap(10).AlignItems(ui.Center).Children(func() {
					breeze.Avatar(c, *a.text["name"], nil)
					ui.Text(c, *a.text["name"])
				})
			},
			code: func(a *args) string { return "breeze.Avatar(c, " + q(*a.text["name"]) + ", nil)" }},
		{component: "Toasts", name: "Default", category: "Feedback",
			controls: []control{textArg("message", "Note deleted"), choiceArg("type", "info", "info", "success", "warning", "error"), toggleArg("action", true)},
			render: func(c *ui.Context, a *args) {
				if breeze.Button(c, "Show toast").Clicked() {
					ts := ui.Toast{ID: "story", Title: *a.text["message"], Type: *a.text["type"]}
					if *a.on["action"] {
						ts.Action, ts.OnAction = "Undo", func() {}
					}
					c.AddToast(ts)
				}
				breeze.Toasts(c)
			},
			code: func(a *args) string {
				action := ""
				if *a.on["action"] {
					action = ", Action: \"Undo\", OnAction: app.undo"
				}
				return "breeze.Toasts(c) // once in the view\nc.AddToast(ui.Toast{Title: " + q(*a.text["message"]) + ", Type: " + q(*a.text["type"]) + action + "})"
			}},

		// ---- Overlays ----
		{component: "Tooltip", name: "Platform look", category: "Overlays",
			controls: []control{textArg("tip", "Share the document")},
			render: func(c *ui.Context, a *args) {
				breeze.Tooltip(c, breeze.Button(c, "Share"), *a.text["tip"])
			},
			code: func(a *args) string {
				return "breeze.Tooltip(c, breeze.Button(c, \"Share\"), " + q(*a.text["tip"]) + ")"
			}},
	}
}

// feedbackFindMatches keeps the find bar story's current match, which must outlive a frame: one for each platform, as the platforms side by side show a find bar each.
var feedbackFindMatches = map[breeze.Platform]*int{}

func feedbackFindMatch(c *ui.Context) *int {
	p := breeze.PlatformOf(c)
	m, ok := feedbackFindMatches[p]
	if !ok {
		m = new(int)
		feedbackFindMatches[p] = m
	}
	return m
}
