import { useMemo, useState } from "react";
import {
  createColumnHelper,
  createFilteredRowModel,
  createPaginatedRowModel,
  columnFilteringFeature,
  globalFilteringFeature,
  rowPaginationFeature,
  tableFeatures,
  useTable,
  type PaginationState,
} from "@tanstack/react-table";
import { ChevronLeft, ChevronRight, Search } from "lucide-react";
import type { Catalog, CatalogItem } from "../lib/catalog";

type CatalogRow = CatalogItem & { kind: "component" | "utility" | "token" };
const features = tableFeatures({
  columnFilteringFeature,
  globalFilteringFeature,
  rowPaginationFeature,
  filteredRowModel: createFilteredRowModel(),
  paginatedRowModel: createPaginatedRowModel(),
});
const columnHelper = createColumnHelper<typeof features, CatalogRow>();
const columns = columnHelper.columns([
  columnHelper.accessor("name", {
    header: "Name",
    cell: ({ getValue }) => <code className="font-mono text-xs text-mint">{getValue()}</code>,
  }),
  columnHelper.accessor("kind", {
    header: "Kind",
    cell: ({ getValue }) => (
      <span className="rounded-full bg-[rgba(157,245,208,0.08)] px-2 py-1 font-mono text-[9px] uppercase text-mint">
        {getValue()}
      </span>
    ),
  }),
  columnHelper.accessor("description", {
    header: "Description",
    cell: ({ getValue, row }) => (
      <span className="text-muted">{getValue() ?? row.original.value}</span>
    ),
  }),
  columnHelper.accessor("example", {
    header: "Example",
    cell: ({ getValue }) => (
      <code className="font-mono text-[11px] text-dim">{getValue() ?? "—"}</code>
    ),
  }),
]);

export function CatalogTable({ catalog }: { catalog: Catalog }) {
  const [globalFilter, setGlobalFilter] = useState("");
  const [pagination, setPagination] = useState<PaginationState>({ pageIndex: 0, pageSize: 8 });
  const data = useMemo<CatalogRow[]>(
    () => [
      ...catalog.components.map((item) => ({ ...item, kind: "component" as const })),
      ...catalog.utilities.map((item) => ({ ...item, kind: "utility" as const })),
      ...catalog.tokens.map((item) => ({ ...item, kind: "token" as const })),
    ],
    [catalog],
  );
  const table = useTable({
    features,
    data,
    columns,
    state: { globalFilter, pagination },
    onGlobalFilterChange: setGlobalFilter,
    onPaginationChange: setPagination,
  });

  return (
    <section className="mt-20" aria-labelledby="catalog-reference-title">
      <div className="mb-7 flex items-end justify-between gap-7 max-[680px]:block">
        <div>
          <p className="mb-5 font-mono text-[10px] uppercase tracking-[0.09em] text-mint">
            Schema explorer
          </p>
          <h2
            className="mb-3.5 text-[clamp(29px,3.5vw,44px)] leading-[1.1] tracking-[-0.065em]"
            id="catalog-reference-title"
          >
            Browse the typed catalog
          </h2>
          <p className="m-0 max-w-[650px] text-[13px] text-muted">
            Every entry is sourced from the same validated catalog used by the Go package, CLI, and
            MCP server.
          </p>
        </div>
        <label className="mt-4 flex shrink-0 items-center gap-2 rounded-md border border-line bg-panel px-3 py-2 text-muted">
          <Search size={15} aria-hidden="true" />
          <span className="sr-only">Search catalog</span>
          <input
            className="w-[190px] bg-transparent text-xs outline-none placeholder:text-dim max-[680px]:w-full"
            value={globalFilter}
            onChange={(event) => setGlobalFilter(event.target.value)}
            placeholder="Search catalog"
          />
        </label>
      </div>
      <div className="overflow-x-auto rounded-lg border border-line bg-panel">
        <table className="w-full min-w-[680px] border-collapse text-left text-xs">
          <thead>
            {table.getHeaderGroups().map((group) => (
              <tr key={group.id}>
                {group.headers.map((header) => (
                  <th
                    className="border-b border-line px-4 py-3 font-mono text-[10px] uppercase tracking-[0.08em] text-dim"
                    key={header.id}
                  >
                    {header.isPlaceholder ? null : <table.FlexRender header={header} />}
                  </th>
                ))}
              </tr>
            ))}
          </thead>
          <tbody>
            {table.getRowModel().rows.length ? (
              table.getRowModel().rows.map((row) => (
                <tr key={row.id}>
                  {row.getAllCells().map((cell) => (
                    <td
                      className="border-b border-[rgba(27,53,69,0.65)] px-4 py-3.5 text-muted"
                      key={cell.id}
                    >
                      <table.FlexRender cell={cell} />
                    </td>
                  ))}
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan={columns.length} className="px-4 py-8 text-center text-muted">
                  No catalog entries match that search.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
      <div className="mt-3 flex items-center justify-between gap-3 font-mono text-[10px] text-dim">
        <span className="truncate">
          {table.getFilteredRowModel().rows.length} entries · page {pagination.pageIndex + 1} of{" "}
          {table.getPageCount() || 1}
        </span>
        <div className="flex gap-1.5">
          <button
            className="rounded border border-line p-1.5 text-muted enabled:hover:border-mint enabled:hover:text-mint disabled:cursor-not-allowed disabled:opacity-40"
            type="button"
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
            aria-label="Previous page"
          >
            <ChevronLeft size={15} />
          </button>
          <button
            className="rounded border border-line p-1.5 text-muted enabled:hover:border-mint enabled:hover:text-mint disabled:cursor-not-allowed disabled:opacity-40"
            type="button"
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
            aria-label="Next page"
          >
            <ChevronRight size={15} />
          </button>
        </div>
      </div>
    </section>
  );
}
