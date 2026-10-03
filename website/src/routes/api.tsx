import { createRoute, Link } from "@tanstack/react-router";
import { motion } from "motion/react";
import { ArrowLeft, Braces, Check, Copy } from "lucide-react";
import { RootRoute } from "./root";
import { apiEndpoints } from "../lib/api";
import { useUiStore } from "../store";
import { useCopy } from "../copy";
import { useQuery } from "@tanstack/react-query";
import { getCatalog } from "../lib/catalog";
import { CatalogTable } from "../components/CatalogTable";

export const ApiRoute = createRoute({
  getParentRoute: () => RootRoute,
  path: "/api",
  component: ApiPage,
});
function ApiPage() {
  const { locale } = useUiStore();
  const t = useCopy(locale);
  const { data: catalog } = useQuery({ queryKey: ["catalog"], queryFn: getCatalog });
  return (
    <main className="px-0 pb-[120px] pt-12 max-[680px]:pb-20 max-[680px]:pt-12">
      <div className="mx-auto w-[calc(100%-48px)] max-w-[1160px] max-[680px]:w-[calc(100%-34px)]">
        <motion.div
          className="mb-14 max-w-[760px]"
          initial={{ opacity: 0, y: 15 }}
          animate={{ opacity: 1, y: 0 }}
        >
          <Link
            className="mb-10 inline-flex items-center gap-1.5 font-mono text-[10px] text-muted hover:text-ink"
            to="/"
          >
            <ArrowLeft size={15} /> {t.footer}
          </Link>
          <p className="mb-5 font-mono text-[10px] uppercase tracking-[0.09em] text-mint">
            {t.apiKicker}
          </p>
          <h1 className="mb-[26px] text-[clamp(44px,5.3vw,74px)] font-bold leading-[1.04] tracking-[-0.078em]">
            {locale === "zh" ? "Breeze UI API 文档" : "Breeze UI API reference"}
          </h1>
          <p className="max-w-[680px] text-[15px] text-muted">
            {locale === "zh"
              ? "用一套稳定、可预测的目录接口，把设计系统接入你的工具和工作流。"
              : "A stable, predictable catalog interface for connecting the design system to your tools and workflows."}
          </p>
          <div className="mt-6 flex flex-wrap gap-4 font-mono text-[10px] text-muted">
            <span className="inline-flex items-center gap-1.5">
              <Check size={14} /> versioned
            </span>
            <span className="inline-flex items-center gap-1.5">
              <Check size={14} /> read-only
            </span>
            <span className="inline-flex items-center gap-1.5">
              <Check size={14} /> typed JSON
            </span>
          </div>
        </motion.div>
        <div className="grid grid-cols-2 gap-2.5 max-[680px]:grid-cols-1">
          {apiEndpoints.map((endpoint, index) => (
            <motion.article
              className="rounded-lg border border-line bg-panel p-5"
              key={endpoint.path}
              initial={{ opacity: 0, y: 12 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.08 }}
            >
              <div className="flex items-center gap-2 border-b border-line pb-3 font-mono text-[10px]">
                <span
                  className={`rounded px-1.5 py-0.5 ${endpoint.method === "GET" ? "bg-[rgba(157,245,208,0.1)] text-mint" : "bg-[rgba(255,174,118,0.1)] text-orange"}`}
                >
                  {endpoint.method}
                </span>
                <code>{endpoint.path}</code>
                <button
                  className="ml-auto text-dim hover:text-mint"
                  title={t.copy}
                  onClick={() => navigator.clipboard?.writeText(endpoint.path)}
                >
                  <Copy size={14} />
                </button>
              </div>
              <h2 className="mb-2 mt-5 text-xl tracking-[-0.04em]">{endpoint.label}</h2>
              <p className="mb-5 text-xs text-muted">{endpoint.description}</p>
              <div className="mb-2 font-mono text-[9px] uppercase tracking-[0.08em] text-dim">
                {t.response}
              </div>
              <pre className="m-0 overflow-auto rounded-md bg-bg p-3 font-mono text-[10px] leading-[1.6] text-muted">
                {endpoint.response}
              </pre>
            </motion.article>
          ))}
        </div>
        <div className="mt-3.5 flex gap-3.5 rounded-lg border border-[rgba(157,245,208,0.22)] bg-[rgba(157,245,208,0.04)] p-5 text-mint">
          <Braces size={20} />
          <div>
            <h3 className="mb-1 text-sm text-ink">
              {locale === "zh" ? "和 Go 包保持一致" : "The same catalog, everywhere"}
            </h3>
            <p className="m-0 text-[11px] text-muted">
              {locale === "zh"
                ? "API、CLI 和 MCP 服务都从同一份 catalog.json 生成，避免文档和实现漂移。"
                : "The API, CLI, and MCP server share the same catalog.json, so docs and implementation stay in sync."}
            </p>
          </div>
        </div>
        {catalog && <CatalogTable catalog={catalog} />}
      </div>
    </main>
  );
}
