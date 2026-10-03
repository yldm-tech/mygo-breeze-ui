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
    <section
      className="border-t border-line bg-[rgba(9,22,36,0.72)] px-0 py-[92px] pb-[105px] max-[680px]:py-[70px] max-[680px]:pb-20"
      id="catalog"
    >
      <div className="mx-auto w-[calc(100%-48px)] max-w-[1160px] max-[680px]:w-[calc(100%-34px)]">
        <div className="mb-9 flex items-end justify-between gap-7 max-[680px]:mb-[27px] max-[680px]:block">
          <div>
            <p className="mb-5 font-mono text-[10px] uppercase tracking-[0.09em] text-mint">
              {t.catalogKicker}
            </p>
            <h2 className="mb-[9px] text-[clamp(29px,3.5vw,44px)] leading-[1.1] tracking-[-0.065em]">
              {t.catalogTitle}
            </h2>
            <p className="m-0 text-[13px] text-muted">{t.catalogText}</p>
          </div>
          <div className="flex shrink-0 gap-1.5 max-[680px]:mt-[22px] max-[680px]:overflow-auto max-[680px]:pb-0.5">
            {(["all", "components", "utilities", "tokens"] as const).map((value) => (
              <button
                key={value}
                className={`rounded-full border px-3 py-[7px] text-[10px] transition ${filter === value ? "border-mint bg-mint text-bg" : "border-line bg-transparent text-muted hover:border-muted"}`}
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
        <div className="grid grid-cols-4 gap-2.5 max-[940px]:grid-cols-3 max-[680px]:grid-cols-2 max-[680px]:gap-2">
          {isLoading ? (
            <div className="col-span-full text-muted">Loading catalog…</div>
          ) : (
            items.map((item, index) => (
              <motion.article
                className="min-h-[164px] min-w-0 rounded-lg border border-[#173547] bg-panel p-[18px] transition hover:-translate-y-0.5 hover:border-[#466576] max-[680px]:min-h-[160px] max-[680px]:p-3.5"
                key={`${item.kind}-${item.name}`}
                initial={{ opacity: 0, y: 12 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true, amount: 0.2 }}
                transition={{ delay: Math.min(index * 0.025, 0.25) }}
              >
                <div className="mb-[19px] font-mono text-[8px] uppercase tracking-[0.08em] text-mint max-[680px]:mb-[15px]">
                  {item.kind}
                </div>
                <h3 className="mb-2 text-base tracking-[-0.025em]">{item.name}</h3>
                <p className="m-0 text-xs leading-[1.75] text-muted max-[680px]:text-[9px]">
                  {item.description ?? item.value}
                </p>
                {item.example && (
                  <code className="mt-2 block max-w-full overflow-hidden text-ellipsis whitespace-nowrap font-mono text-[9px] text-dim">
                    {item.example}
                  </code>
                )}
              </motion.article>
            ))
          )}
        </div>
      </div>
    </section>
  );
}
