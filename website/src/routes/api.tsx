import { createRoute } from "@tanstack/react-router";
import { motion } from "motion/react";
import { ArrowLeft, Braces, Check, Copy } from "lucide-react";
import { RootRoute } from "./root";
import { apiEndpoints } from "../lib/api";
import { useUiStore } from "../store";
import { useCopy } from "../copy";

export const ApiRoute = createRoute({
  getParentRoute: () => RootRoute,
  path: "/api",
  component: ApiPage,
});
function ApiPage() {
  const { locale } = useUiStore();
  const t = useCopy(locale);
  return (
    <main className="api-docs">
      <div className="container">
        <motion.div
          className="docs-hero"
          initial={{ opacity: 0, y: 15 }}
          animate={{ opacity: 1, y: 0 }}
        >
          <a className="back-link" href="/#top">
            <ArrowLeft size={15} /> {t.footer}
          </a>
          <p className="eyebrow">{t.apiKicker}</p>
          <h1>{locale === "zh" ? "Breeze UI API 文档" : "Breeze UI API reference"}</h1>
          <p>
            {locale === "zh"
              ? "用一套稳定、可预测的目录接口，把设计系统接入你的工具和工作流。"
              : "A stable, predictable catalog interface for connecting the design system to your tools and workflows."}
          </p>
          <div className="docs-meta">
            <span>
              <Check size={14} /> versioned
            </span>
            <span>
              <Check size={14} /> read-only
            </span>
            <span>
              <Check size={14} /> typed JSON
            </span>
          </div>
        </motion.div>
        <div className="endpoint-list">
          {apiEndpoints.map((endpoint, index) => (
            <motion.article
              className="endpoint-card"
              key={endpoint.path}
              initial={{ opacity: 0, y: 12 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.08 }}
            >
              <div className="endpoint-top">
                <span className={`method method-${endpoint.method.toLowerCase()}`}>
                  {endpoint.method}
                </span>
                <code>{endpoint.path}</code>
                <button
                  title={t.copy}
                  onClick={() => navigator.clipboard?.writeText(endpoint.path)}
                >
                  <Copy size={14} />
                </button>
              </div>
              <h2>{endpoint.label}</h2>
              <p>{endpoint.description}</p>
              <div className="response-label">{t.response}</div>
              <pre>{endpoint.response}</pre>
            </motion.article>
          ))}
        </div>
        <div className="docs-note">
          <Braces size={20} />
          <div>
            <h3>{locale === "zh" ? "和 Go 包保持一致" : "The same catalog, everywhere"}</h3>
            <p>
              {locale === "zh"
                ? "API、CLI 和 MCP 服务都从同一份 catalog.json 生成，避免文档和实现漂移。"
                : "The API, CLI, and MCP server share the same catalog.json, so docs and implementation stay in sync."}
            </p>
          </div>
        </div>
      </div>
    </main>
  );
}
