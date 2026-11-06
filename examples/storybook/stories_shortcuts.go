package main

import (
	"strings"

	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func init() { storyGroups = append(storyGroups, shortcutStories) }

// shortcutKeys are the keys the key control offers, by the names of their ui.Key constants.
var shortcutKeys = []struct {
	name string
	key  ui.Key
}{
	{"T", ui.KeyT},
	{"K", ui.KeyK},
	{"Left", ui.KeyLeft},
	{"Space", ui.KeySpace},
	{"Enter", ui.KeyEnter},
	{"Escape", ui.KeyEscape},
	{"Comma", ui.KeyComma},
	{"F5", ui.KeyF5},
}

// shortcutArg returns the modifiers and the key the story's controls choose, and the Go that passes them.
func shortcutArg(a *args) (ui.Modifiers, ui.Key, string) {
	var mods ui.Modifiers
	var names []string
	for _, m := range []struct {
		control, code string
		mod           ui.Modifiers
	}{{"cmd", "ui.Cmd", ui.Cmd}, {"alt", "ui.Alt", ui.Alt}, {"shift", "ui.Shift", ui.Shift}} {
		if *a.on[m.control] {
			mods |= m.mod
			names = append(names, m.code)
		}
	}
	code := "0"
	if len(names) > 0 {
		code = strings.Join(names, "|")
	}
	key := shortcutKeys[0]
	for _, k := range shortcutKeys {
		if k.name == *a.text["key"] {
			key = k
		}
	}
	return mods, key.key, code + ", ui.Key" + key.name
}

// shortcutControls are the controls choosing a shortcut: its modifiers and its key.
func shortcutControls() []control {
	var names []string
	for _, k := range shortcutKeys {
		names = append(names, k.name)
	}
	return []control{toggleArg("cmd", true), toggleArg("alt", false), toggleArg("shift", true), choiceArg("key", "T", names...)}
}

// shortcutStories are the stories of the widgets of widgets_shortcuts.go: keyboard shortcuts as the platforms write and draw them.
func shortcutStories() []*story {
	return []*story{
		{component: "ShortcutLabel", name: "Default", category: "Text",
			controls: shortcutControls(),
			render: func(c *ui.Context, a *args) {
				mods, key, _ := shortcutArg(a)
				breeze.ShortcutLabel(c, mods, key)
			},
			code: func(a *args) string {
				_, _, code := shortcutArg(a)
				return "breeze.ShortcutLabel(c, " + code + ")"
			}},
		{component: "ShortcutLabel", name: "Cheat sheet", category: "Text",
			render: func(c *ui.Context, a *args) {
				ui.Column(c).Gap(10).Width(320).Children(func() {
					for _, s := range []struct {
						label string
						mods  ui.Modifiers
						key   ui.Key
					}{
						{"Play or pause", 0, ui.KeySpace},
						{"Command palette", ui.Cmd, ui.KeyK},
						{"Show translation", ui.Cmd | ui.Shift, ui.KeyT},
						{"Previous sentence", ui.Alt | ui.Cmd, ui.KeyLeft},
					} {
						ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
							breeze.Text(c, s.label).Grow(1)
							breeze.ShortcutLabel(c, s.mods, s.key)
						})
					}
					ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
						breeze.Text(c, "Speed").Grow(1)
						breeze.Keycaps(c, "[", "]")
					})
				})
			},
			code: func(*args) string {
				return "ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {\n\tbreeze.Text(c, \"Show translation\").Grow(1)\n\tbreeze.ShortcutLabel(c, ui.Cmd|ui.Shift, ui.KeyT)\n})"
			}},
		{component: "Keycaps", name: "Chord", category: "Text",
			controls: []control{textArg("keys", "[ ]")},
			render:   func(c *ui.Context, a *args) { breeze.Keycaps(c, strings.Fields(*a.text["keys"])...) },
			code: func(a *args) string {
				var keys []string
				for _, k := range strings.Fields(*a.text["keys"]) {
					keys = append(keys, q(k))
				}
				return "breeze.Keycaps(c, " + strings.Join(keys, ", ") + ")"
			}},
		{component: "ShortcutText", name: "Menu row", category: "Text",
			controls: append(shortcutControls(), textArg("label", "Show translation")),
			render: func(c *ui.Context, a *args) {
				mods, key, _ := shortcutArg(a)
				ui.Row(c).Gap(12).Width(260).AlignItems(ui.Center).Children(func() {
					breeze.Text(c, *a.text["label"]).Grow(1)
					breeze.Text(c, breeze.ShortcutText(c, mods, key), breeze.TextSecondary)
				})
			},
			code: func(a *args) string {
				_, _, code := shortcutArg(a)
				return "ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {\n\tbreeze.Text(c, " + q(*a.text["label"]) + ").Grow(1)\n\tbreeze.Text(c, breeze.ShortcutText(c, " + code + "), breeze.TextSecondary)\n})"
			}},
	}
}
