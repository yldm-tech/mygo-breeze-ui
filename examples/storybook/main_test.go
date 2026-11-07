package main

import (
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestNewBookDefaultsToSideBySide(t *testing.T) {
	b := newBook()
	if b.platform != 3 {
		t.Fatalf("new Storybook should default to %q, got %q", platformNames[3], platformNames[b.platform])
	}
}

// TestEveryStory builds every story on every platform choice, light and dark, and checks the toolbar names it. With BREEZE_SCREENSHOTS set, it writes a few views of the book there at 2x.
func TestEveryStory(t *testing.T) {
	b := newBook()
	tester := ui.NewTester(b.view, 1280, 820)
	for _, s := range b.stories {
		for platform := range platformNames {
			for dark := range 2 {
				b.selected, b.platform, b.dark = s.id(), platform, dark
				tester.Frame()
				if !tester.HasText(s.component) || !tester.HasText(s.name) {
					t.Fatalf("%s on %s, dark %d: the toolbar does not name the story", s.id(), platformNames[platform], dark)
				}
			}
		}
	}
	dir := os.Getenv("BREEZE_SCREENSHOTS")
	if dir == "" {
		return
	}
	for _, shot := range []struct {
		file           string
		story          string
		platform, dark int
	}{
		{"storybook-button-macos.png", "Button/Default", 0, 0},
		{"storybook-select-side-by-side.png", "Select/Default", 3, 0},
		{"storybook-toggles-side-by-side-dark.png", "ToggleGroup/Formatting", 3, 1},
		{"storybook-table-windows.png", "Table/Files", 1, 0},
	} {
		b.selected, b.platform, b.dark = shot.story, shot.platform, shot.dark
		shotTester := ui.NewTester(b.view, 1280, 820)
		shotTester.SetScale(2)
		shotTester.Frame()
		f, err := os.Create(filepath.Join(dir, shot.file))
		if err != nil {
			t.Fatal(err)
		}
		if err := png.Encode(f, shotTester.Image()); err != nil {
			t.Fatal(err)
		}
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

// TestControlsDriveTheStory checks that a control and the story share a value: editing the label control renames the button.
func TestControlsDriveTheStory(t *testing.T) {
	b := newBook()
	b.selected = "Button/Default"
	tester := ui.NewTester(b.view, 1280, 820)
	tester.Frame()
	*b.argsOf(b.story()).text["label"] = "Renamed"
	tester.Frame()
	if !tester.HasText("Renamed") {
		t.Fatal("the button did not take the label control's value")
	}
	if !tester.HasText(`breeze.Button(c, "Renamed")`) {
		t.Fatal("the code panel did not follow the label control")
	}
}

func TestStorybookChromeFollowsLocale(t *testing.T) {
	b := newBook()
	tester := ui.NewTester(b.view, 1280, 820)
	tester.SetLocale("zh-Hans")
	tester.Frame()
	for _, want := range []string{"组件", "交互式 story 和平台预览", "故事", "平台", "外观", "控件", "Go 代码", "复制"} {
		if !tester.HasText(want) {
			t.Fatalf("Chinese chrome does not contain %q", want)
		}
	}
	tester.SetLocale("en")
	tester.Frame()
	for _, want := range []string{"Components", "Interactive stories and platform previews", "Platform", "Appearance", "Controls", "Go code", "Copy"} {
		if !tester.HasText(want) {
			t.Fatalf("English chrome does not contain %q", want)
		}
	}
}
