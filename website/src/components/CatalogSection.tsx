import { motion } from "motion/react";
import { useQuery } from "@tanstack/react-query";
import { getCatalog, type CatalogItem } from "../lib/catalog";
import { useUiStore, type CatalogFilter } from "../store";
import { useCopy } from "../copy";

const filterMap: Record<CatalogFilter, "all" | "component" | "utility" | "token"> = {
  all: "all",
  components: "component",
  utilities: "utility",
  tokens: "token",
};
function itemKind(item: CatalogItem, kind: string) {
  return item.kind ?? kind;
}

export function CatalogSection() {
  const { locale, filter, setFilter } = useUiStore();
  const t = useCopy(locale);
  const { data, isLoading } = useQuery({ queryKey: ["catalog"], queryFn: getCatalog });
  const selected = filterMap[filter];
  const groups = data
    ? (
        [
          ["component", data.components],
          ["utility", data.utilities],
          ["token", data.tokens],
        ] as const
      ).filter(([kind]) => selected === "all" || selected === kind)
    : [];
  const items = groups.flatMap(([kind, entries]) =>
    entries
      .slice(0, kind === "token" ? 8 : 12)
      .map((item) => ({ ...item, kind: itemKind(item, kind) })),
  );
  return (
    <section className="catalog-section" id="catalog">
      <div className="container">
        <div className="section-heading">
          <div>
            <p className="eyebrow">{t.catalogKicker}</p>
            <h2>{t.catalogTitle}</h2>
            <p>{t.catalogText}</p>
          </div>
          <div className="filters">
            {(["all", "components", "utilities", "tokens"] as const).map((value) => (
              <button
                key={value}
                className={`filter ${filter === value ? "active" : ""}`}
                onClick={() => setFilter(value)}
              >
                {
                  t[
                    value === "all"
                      ? "all"
                      : value === "components"
                        ? "component"
                        : value === "utilities"
                          ? "utility"
                          : "token"
                  ]
                }
              </button>
            ))}
          </div>
        </div>
        <div className="catalog-grid">
          {isLoading ? (
            <div className="catalog-loading">Loading catalog…</div>
          ) : (
            items.map((item, index) => (
              <motion.article
                className="catalog-card"
                key={`${item.kind}-${item.name}`}
                initial={{ opacity: 0, y: 12 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true, amount: 0.2 }}
                transition={{ delay: Math.min(index * 0.025, 0.25) }}
              >
                <div className="catalog-kind">{item.kind}</div>
                <h3>{item.name}</h3>
                <p>{item.description ?? item.value}</p>
                {item.example && <code>{item.example}</code>}
              </motion.article>
            ))
          )}
        </div>
      </div>
    </section>
  );
}
