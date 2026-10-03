// Package breeze is Breeze UI, a third-party design system for MyGo's native UI.
//
// It translates utility-first design tokens into typed Go helpers. It does
// not run CSS, parse a DOM, or require a browser: styles are applied directly
// to MyGo ui.Elements and components compose the existing accessible widgets.
package breeze

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
)

//go:embed catalog.json
var catalogData []byte

type catalog struct {
	Version    int           `json:"version"`
	Utilities  []CatalogItem `json:"utilities"`
	Components []CatalogItem `json:"components"`
	Tokens     []CatalogItem `json:"tokens"`
}

// CatalogItem describes a utility, component, or token for documentation and
// tooling. Name is stable and is used by the CLI and MCP server.
type CatalogItem struct {
	Name        string `json:"name"`
	Kind        string `json:"kind,omitempty"`
	Category    string `json:"category,omitempty"`
	Description string `json:"description"`
	Value       string `json:"value,omitempty"`
	Example     string `json:"example,omitempty"`
}

var loadedCatalog = loadCatalog()

func loadCatalog() catalog {
	var c catalog
	if err := json.Unmarshal(catalogData, &c); err != nil {
		panic(fmt.Sprintf("mygo/ui/breeze: invalid embedded catalog: %v", err))
	}
	return c
}

// Catalog returns a copy of all catalog entries in stable name order. The
// result is safe for callers to modify.
func Catalog() []CatalogItem {
	items := make([]CatalogItem, 0, len(loadedCatalog.Utilities)+len(loadedCatalog.Components)+len(loadedCatalog.Tokens))
	for _, item := range loadedCatalog.Utilities {
		item.Kind = "utility"
		items = append(items, item)
	}
	for _, item := range loadedCatalog.Components {
		item.Kind = "component"
		items = append(items, item)
	}
	for _, item := range loadedCatalog.Tokens {
		item.Kind = "token"
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Name == items[j].Name {
			return items[i].Kind < items[j].Kind
		}
		return items[i].Name < items[j].Name
	})
	return items
}

// CatalogJSON returns the source catalog as a copy for documentation tools.
func CatalogJSON() []byte { return append([]byte(nil), catalogData...) }

// Find returns the catalog item with name, or false when it does not exist.
func Find(name string) (CatalogItem, bool) {
	for _, item := range Catalog() {
		if item.Name == name {
			return item, true
		}
	}
	return CatalogItem{}, false
}

// FindKind returns an item only when name belongs to kind. Kind is one of
// "utility", "component", or "token".
func FindKind(kind, name string) (CatalogItem, bool) {
	item, ok := Find(name)
	return item, ok && item.Kind == kind
}
