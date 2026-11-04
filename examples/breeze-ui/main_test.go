package main

import (
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestGallery builds the side-by-side gallery and, when BREEZE_SCREENSHOTS names a directory, writes it there at 2x, light and dark.
func TestGallery(t *testing.T) {
	for _, dark := range []bool{false, true} {
		g := &gallery{remember: true, size: "Medium", volume: 40}
		tester := ui.NewTester(g.view, 1080, 520)
		tester.SetScale(2)
		tester.SetDark(dark)
		tester.Frame()
		for _, s := range []string{"macOS", "Windows", "Linux", "Remember me"} {
			if !tester.HasText(s) {
				t.Fatalf("gallery is missing %q", s)
			}
		}
		dir := os.Getenv("BREEZE_SCREENSHOTS")
		if dir == "" {
			continue
		}
		name := map[bool]string{false: "gallery-light.png", true: "gallery-dark.png"}[dark]
		f, err := os.Create(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		if err := png.Encode(f, tester.Image()); err != nil {
			t.Fatal(err)
		}
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}
	}
}
