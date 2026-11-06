package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestTextStoryFollowsItsControls checks that the Text story's style control changes the step it draws and the Go it shows: a large title is taller than a caption on every platform, the body is ui.Text's call without a style, and the secondary control adds TextSecondary to the code.
func TestTextStoryFollowsItsControls(t *testing.T) {
	b := newBook()
	b.selected = "Text/Default"
	tester := ui.NewTester(b.view, 1280, 820)
	a := b.argsOf(b.story())
	for platform := range 3 {
		b.platform = platform
		heights := map[string]float32{}
		for _, style := range []string{"LargeTitle", "Caption"} {
			*a.text["style"] = style
			tester.Frame()
			r, ok := tester.Find("Settings")
			if !ok {
				t.Fatalf("%s: the story shows no text: %q", platformNames[platform], tester.Texts())
			}
			heights[style] = r.H
			if !tester.HasText("breeze.Text(c, \"Settings\", breeze.Text" + style + ")") {
				t.Errorf("%s: the code does not name breeze.Text%s", platformNames[platform], style)
			}
		}
		if heights["LargeTitle"] <= heights["Caption"] {
			t.Errorf("%s: the large title is %v high, the caption %v", platformNames[platform], heights["LargeTitle"], heights["Caption"])
		}
	}
	*a.on["secondary"] = true
	tester.Frame()
	if !tester.HasText("breeze.Text(c, \"Settings\", breeze.TextCaption|breeze.TextSecondary)") {
		t.Error("the code does not add TextSecondary")
	}
	*a.text["style"] = "Body"
	tester.Frame()
	if !tester.HasText("breeze.Text(c, \"Settings\", breeze.TextSecondary)") {
		t.Errorf("the code of the dimmed body is not Text with TextSecondary: %q", tester.Texts())
	}
	*a.on["secondary"] = false
	tester.Frame()
	if !tester.HasText("breeze.Text(c, \"Settings\")") {
		t.Errorf("the code of the body is not ui.Text's call: %q", tester.Texts())
	}
}
