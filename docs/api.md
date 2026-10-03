# Breeze UI catalog API

The website's [API reference](https://yldm-tech.github.io/mygo-breeze-ui/api/) documents the read-only catalog contract used by the CLI and MCP server. The source of truth is [`ui/breeze/catalog.json`](../ui/breeze/catalog.json); the website validates the same shape with Zod before rendering it.

## Endpoints

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/api/catalog` | Return the complete versioned catalog. |
| `GET` | `/api/catalog?kind=component` | Filter entries by `component`, `utility`, or `token`. |
| `GET` | `/api/health` | Return the catalog availability and version. |

## Catalog response

```json
{
  "version": 1,
  "components": [{ "name": "Button", "description": "...", "example": "..." }],
  "utilities": [{ "name": "row", "description": "...", "example": "..." }],
  "tokens": [{ "name": "mint", "value": "#9df5d0" }]
}
```

The API is intentionally read-only. The interactive page includes a Zod-validated query form and a searchable, paginated TanStack Table view of the catalog.
