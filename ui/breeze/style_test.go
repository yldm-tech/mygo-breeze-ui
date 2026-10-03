package breeze

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestSpaceAndPalette(t *testing.T) {
	if got := Space(4); got != 16 {
		t.Fatalf("Space(4) = %v, want 16", got)
	}
	tester := ui.NewTester(func(c *ui.Context) {
		p := Palette(c)
		if p.Background.A == 0 || p.Text.A == 0 {
			t.Fatal("default palette contains transparent required colors")
		}
	}, 320, 200)
	tester.Frame()
}

func TestUtilityStylesRender(t *testing.T) {
	tester := ui.NewTester(func(c *ui.Context) {
		Row(c, Gap(2), P(3), Bg(Palette(c).Surface)).Children(func() {
			Text(c, "Ready", TextSize(TextLG), FontWeight(600))
		})
	}, 320, 200)
	tester.Frame()
	if !tester.HasText("Ready") {
		t.Fatal("styled text is missing")
	}
}
