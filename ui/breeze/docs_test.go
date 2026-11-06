package breeze

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestWidgetDocs checks docs/widgets against the groups of widgets: each group has a page, docs/widgets/<group>.md, that names every widget of the group as `Name` or `Name(`, and every page is a group's.
func TestWidgetDocs(t *testing.T) {
	groups := widgetGroups(t)
	pages, err := filepath.Glob("../../docs/widgets/*.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, page := range pages {
		group := strings.TrimSuffix(filepath.Base(page), ".md")
		widgets, ok := groups[group]
		if !ok {
			t.Errorf("%s is the page of no group: want ui/breeze/widgets_%s.go", page, group)
			continue
		}
		delete(groups, group)
		data, err := os.ReadFile(page)
		if err != nil {
			t.Fatal(err)
		}
		text := string(data)
		for _, name := range widgets {
			if !strings.Contains(text, "`"+name+"`") && !strings.Contains(text, "`"+name+"(") {
				t.Errorf("%s does not name %s", page, name)
			}
		}
	}
	for group, widgets := range groups {
		if len(widgets) > 0 {
			t.Errorf("group %s has %v, but there is no docs/widgets/%s.md", group, widgets, group)
		}
	}
}
