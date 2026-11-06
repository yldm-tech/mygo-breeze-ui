package breeze

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestCatalogIsStableAndUnique(t *testing.T) {
	items := Catalog()
	for i := 1; i < len(items); i++ {
		if items[i-1].Name > items[i].Name {
			t.Fatalf("catalog is not sorted at %q and %q", items[i-1].Name, items[i].Name)
		}
		if items[i-1].Name == items[i].Name && items[i-1].Kind == items[i].Kind {
			t.Fatalf("duplicate %s %q", items[i].Kind, items[i].Name)
		}
	}
}

func TestCatalogLookup(t *testing.T) {
	if item, ok := FindKind("component", "Button"); !ok || item.Category != "actions" || item.Look != "drawn" {
		t.Fatalf("Button lookup = %#v, %v", item, ok)
	}
	if _, ok := FindKind("component", "missing"); ok {
		t.Fatal("missing component found")
	}
	if _, ok := FindKind("platform", "Button"); ok {
		t.Fatal("component returned for platform lookup")
	}
	for _, p := range Platforms {
		if _, ok := FindKind("platform", p.String()); !ok {
			t.Errorf("platform %v is missing from the catalog", p)
		}
	}
}

// widgetGroups maps each group of widgets to the exported functions its Go file declares, widgets.go and extras.go and the widgets_<group>.go files: a group is the file's name without .go and without a leading widgets_. The catalog (catalog/<group>.json) and the documentation (docs/widgets/<group>.md) are split by the same groups.
func widgetGroups(t *testing.T) map[string][]string {
	t.Helper()
	files, err := filepath.Glob("widgets*.go")
	if err != nil {
		t.Fatal(err)
	}
	files = append(files, "extras.go")
	groups := map[string][]string{}
	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(token.NewFileSet(), name, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		group := strings.TrimPrefix(strings.TrimSuffix(name, ".go"), "widgets_")
		groups[group] = []string{}
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv == nil && fn.Name.IsExported() {
				groups[group] = append(groups[group], fn.Name.Name)
			}
		}
		sort.Strings(groups[group])
	}
	return groups
}

// TestCatalogMatchesWidgets checks the catalog against the widget functions of extras.go and the widgets*.go files, group by group: catalog/<group>.json lists exactly the exported functions of the group's Go file, each with an example calling it, those Breeze UI draws itself marked drawn; only catalog/platforms.json gives the version and the platforms, and no file of catalog/ is without its Go file.
func TestCatalogMatchesWidgets(t *testing.T) {
	if _, err := os.Stat("catalog.json"); err == nil {
		t.Fatal("catalog.json is back: the catalog is the files of catalog/, one per group of widgets (a merge or rebase may have restored the old file)")
	}
	groups := widgetGroups(t)
	names, err := fs.Glob(catalogFiles, "catalog/*.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range names {
		data, err := catalogFiles.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		var f catalogFile
		if err := json.Unmarshal(data, &f); err != nil {
			t.Fatal(err)
		}
		group := strings.TrimSuffix(path.Base(name), ".json")
		if group == "platforms" {
			if f.Version == nil || len(f.Platforms) == 0 || len(f.Components) != 0 {
				t.Errorf("%s: want the version and the platforms and no components", name)
			}
			continue
		}
		if f.Version != nil || len(f.Platforms) != 0 {
			t.Errorf("%s: only catalog/platforms.json gives the version and the platforms", name)
		}
		widgets, ok := groups[group]
		if !ok {
			t.Errorf("%s lists components of no Go file: want widgets_%s.go", name, group)
			continue
		}
		delete(groups, group)
		var listed []string
		for _, item := range f.Components {
			listed = append(listed, item.Name)
			if !strings.Contains(item.Example, "breeze."+item.Name+"(") {
				t.Errorf("%s: example %q does not call breeze.%s", item.Name, item.Example, item.Name)
			}
			drawn := drawnWidgets[item.Name]
			if want := map[bool]string{true: "drawn", false: "themed"}[drawn]; item.Look != want {
				t.Errorf("%s: look %q, want %q", item.Name, item.Look, want)
			}
		}
		sort.Strings(listed)
		if strings.Join(widgets, ",") != strings.Join(listed, ",") {
			t.Errorf("the Go file of group %s has %v, %s lists %v", group, widgets, name, listed)
		}
	}
	for group, widgets := range groups {
		if len(widgets) > 0 {
			t.Errorf("group %s has %v, but there is no catalog/%s.json", group, widgets, group)
		}
	}
}

// drawnWidgets are the widgets Breeze UI draws as each platform does; the others are package ui's. The test file of each widgets_*.go adds its own in an init function.
var drawnWidgets = map[string]bool{
	"Button": true, "PrimaryButton": true, "Checkbox": true, "Switch": true, "Radio": true, "Toggle": true,
	"TextInput": true, "TextArea": true, "SearchField": true, "Select": true, "Slider": true,
	"Progress": true, "Spinner": true, "Segmented": true, "Tabs": true,
	"Combobox": true, "Modal": true, "AlertDialog": true,
	"Link": true, "Divider": true, "Badge": true, "MenuButton": true, "PopoverMenuButton": true,
	"Popover": true, "ToggleGroup": true, "Field": true, "Stepper": true,
}
