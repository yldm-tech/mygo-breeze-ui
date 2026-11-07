package breeze

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"sort"
)

// catalogFiles are the files the catalog is assembled from: catalog/platforms.json holds the version and the platforms, and each other file the components of one Go file, catalog/<group>.json those of widgets_<group>.go (or <group>.go for widgets.go and extras.go), so that a pull request adding widgets to a group touches only that group's files.
//
//go:embed catalog/*.json
var catalogFiles embed.FS

// catalogFile is one file of catalog/.
type catalogFile struct {
	Version    *int          `json:"version"`
	Platforms  []CatalogItem `json:"platforms"`
	Components []CatalogItem `json:"components"`
}

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

var loadedCatalog, catalogData = loadCatalog()

// loadCatalog merges the files of catalog/ in the order of their names into one catalog, its components in the byte order of their names, and returns it with its JSON, two-space indented, as the website and the MCP server assemble it too. The version must be given once, and no two components or platforms may share a name.
func loadCatalog() (catalog, []byte) {
	names, err := fs.Glob(catalogFiles, "catalog/*.json")
	if err != nil {
		panic(fmt.Sprintf("mygo-breeze-ui/ui/breeze: invalid embedded catalog: %v", err))
	}
	var c catalog
	versions := 0
	seen := map[string]bool{}
	for _, name := range names {
		data, err := catalogFiles.ReadFile(name)
		if err != nil {
			panic(fmt.Sprintf("mygo-breeze-ui/ui/breeze: invalid embedded catalog: %v", err))
		}
		var f catalogFile
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&f); err != nil {
			panic(fmt.Sprintf("mygo-breeze-ui/ui/breeze: invalid embedded catalog %s: %v", name, err))
		}
		if f.Version != nil {
			versions++
			c.Version = *f.Version
		}
		for kind, items := range map[string][]CatalogItem{"platform": f.Platforms, "component": f.Components} {
			for _, item := range items {
				if seen[kind+" "+item.Name] {
					panic(fmt.Sprintf("mygo-breeze-ui/ui/breeze: invalid embedded catalog %s: %s %q is listed twice", name, kind, item.Name))
				}
				seen[kind+" "+item.Name] = true
			}
		}
		c.Platforms = append(c.Platforms, f.Platforms...)
		c.Components = append(c.Components, f.Components...)
	}
	if versions != 1 {
		panic(fmt.Sprintf("mygo-breeze-ui/ui/breeze: invalid embedded catalog: %d files give a version, want 1", versions))
	}
	sort.Slice(c.Components, func(i, j int) bool { return c.Components[i].Name < c.Components[j].Name })
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(c); err != nil {
		panic(fmt.Sprintf("mygo-breeze-ui/ui/breeze: invalid embedded catalog: %v", err))
	}
	return c, buf.Bytes()
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

// CatalogJSON returns the catalog as JSON, the files of catalog/ merged, as a copy for documentation tools.
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
