package breeze

import "testing"

func TestCatalogIsStableAndUnique(t *testing.T) {
	items := Catalog()
	if len(items) < 10 {
		t.Fatalf("catalog has %d entries, want at least 10", len(items))
	}
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
	if item, ok := FindKind("component", "Button"); !ok || item.Category != "actions" {
		t.Fatalf("Button lookup = %#v, %v", item, ok)
	}
	if _, ok := FindKind("component", "missing"); ok {
		t.Fatal("missing component found")
	}
	if _, ok := FindKind("token", "Button"); ok {
		t.Fatal("component returned for token lookup")
	}
}

func TestCatalogCoversPublicComponents(t *testing.T) {
	want := []string{"Button", "Card", "Input", "TextArea", "Checkbox", "Switch", "Select", "Badge", "Progress", "Tabs", "Dialog"}
	for _, name := range want {
		if _, ok := FindKind("component", name); !ok {
			t.Errorf("component %q is missing from the catalog", name)
		}
	}
}
