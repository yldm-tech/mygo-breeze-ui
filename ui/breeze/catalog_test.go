package breeze

import (
	"go/ast"
	"go/parser"
	"go/token"
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

// TestCatalogMatchesWidgets checks the catalog against the widget functions of extras.go and the widgets*.go files: every exported one is listed with an example calling it, and those Breeze UI draws itself are marked drawn.
func TestCatalogMatchesWidgets(t *testing.T) {
	var widgets []string
	files, err := filepath.Glob("widgets*.go")
	if err != nil {
		t.Fatal(err)
	}
	files = append(files, "extras.go")
	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(token.NewFileSet(), name, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv == nil && fn.Name.IsExported() {
				widgets = append(widgets, fn.Name.Name)
			}
		}
	}
	sort.Strings(widgets)
	var listed []string
	for _, item := range Catalog() {
		if item.Kind != "component" {
			continue
		}
		listed = append(listed, item.Name)
		if !strings.Contains(item.Example, "breeze."+item.Name+"(") {
			t.Errorf("%s: example %q does not call breeze.%s", item.Name, item.Example, item.Name)
		}
		drawn := drawnWidgets[item.Name]
		if want := map[bool]string{true: "drawn", false: "themed"}[drawn]; item.Look != want {
			t.Errorf("%s: look %q, want %q", item.Name, item.Look, want)
		}
	}
	if strings.Join(widgets, ",") != strings.Join(listed, ",") {
		t.Fatalf("widgets.go has %v, the catalog lists %v", widgets, listed)
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
