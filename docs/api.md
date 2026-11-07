# Breeze UI catalog API

The website's [API reference](https://breeze-ui.yldm.ai/api/) documents the read-only catalog contract used by the CLI and MCP server. The source of truth is the files of [`ui/breeze/catalog/`](../ui/breeze/catalog): `platforms.json` with the version and the platforms, and one file per group of widgets, `<group>.json` listing the widgets of `widgets_<group>.go` (`widgets.json` and `extras.json` those of `widgets.go` and `extras.go`). The Go package, the MCP server and the website merge them the same way into the catalog the API returns; the website validates its shape with Zod before rendering it, and a Go test checks that each file lists exactly the widgets of its Go file.

## Endpoints

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/api/catalog` | Return the complete versioned catalog. |
| `GET` | `/api/catalog?kind=component` | Filter entries by `component` or `platform`. |
| `GET` | `/api/health` | Return the catalog availability and version. |

## Catalog response

```json
{
  "version": 2,
  "platforms": [{ "name": "Windows", "value": "Windows 11", "description": "...", "example": "breeze.SetPlatform(breeze.Windows)" }],
  "components": [{ "name": "Button", "category": "actions", "look": "drawn", "description": "...", "example": "..." }]
}
```

A component's `look` is `drawn` when Breeze UI draws it as each platform does, and `themed` when it is package `ui`'s widget in the platform's theme. A platform's `value` is the version its look follows.

The API is read-only. The interactive page includes a Zod-validated query form and a searchable, paginated TanStack Table view of the catalog.
