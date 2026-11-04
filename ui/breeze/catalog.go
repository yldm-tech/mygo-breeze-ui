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
	Platforms  []CatalogItem `json:"platforms"`
	Components []CatalogItem `json:"components"`
}

// CatalogItem describes a component or a platform for documentation and tooling. Name is stable and is used by the CLI and the MCP server.
type CatalogItem struct {
	Name string `json:"name"`
	// Kind is "component" or "platform".
	Kind     string `json:"kind,omitempty"`
	Category string `json:"category,omitempty"`
	// Look tells how a component looks: "drawn" when Breeze UI draws it as each platform does, "themed" when it is package ui's widget in the platform's theme.
	Look        string `json:"look,omitempty"`
	Description string `json:"description"`
	// Value is the version a platform's look follows.
	Value   string `json:"value,omitempty"`
	Example string `json:"example,omitempty"`
}

var loadedCatalog = loadCatalog()

func loadCatalog() catalog {
	var c catalog
	if err := json.Unmarshal(catalogData, &c); err != nil {
		panic(fmt.Sprintf("mygo-breeze-ui/ui/breeze: invalid embedded catalog: %v", err))
	}
	return c
}

// Catalog returns a copy of all catalog entries in stable name order. The result is safe for callers to modify.
func Catalog() []CatalogItem {
	items := make([]CatalogItem, 0, len(loadedCatalog.Platforms)+len(loadedCatalog.Components))
	for _, item := range loadedCatalog.Components {
		item.Kind = "component"
		items = append(items, item)
	}
	for _, item := range loadedCatalog.Platforms {
		item.Kind = "platform"
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

// FindKind returns an item only when name belongs to kind, "component" or "platform".
func FindKind(kind, name string) (CatalogItem, bool) {
	for _, item := range Catalog() {
		if item.Kind == kind && item.Name == name {
			return item, true
		}
	}
	return CatalogItem{}, false
}
