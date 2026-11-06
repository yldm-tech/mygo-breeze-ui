import { createRoute, Link } from "@tanstack/react-router";
import { ArrowLeft, ArrowUpRight, Braces, Terminal } from "lucide-react";
import { RootRoute } from "./root";
import { useUiStore } from "../store";
import { Seo } from "../components/Seo";

export const CliMcpRoute = createRoute({
  getParentRoute: () => RootRoute,
  path: "/cli-mcp",
  component: CliMcpPage,
});

function CodeBlock({ children }: { children: string }) {
  return <pre className="overflow-x-auto rounded-xl border border-line bg-[#15171c] p-5 font-mono text-[13px] leading-[1.8] text-[#dce3e8]"><code>{children}</code></pre>;
}

function CliMcpPage() {
  const { locale } = useUiStore();
  const zh = locale === "zh";
  return (
    <main className="pb-[120px] pt-12 max-[680px]:pb-20 max-[680px]:pt-10">
      <Seo title={zh ? "MyGo Breeze CLI 与 MCP" : "MyGo Breeze CLI and MCP"} description={zh ? "使用 MyGo Breeze CLI 和只读 MCP 服务搜索组件、列出平台并生成 Go 代码。" : "Use the MyGo Breeze CLI and read-only MCP server to search components, list platforms and generate Go code."} socialTitle="MyGo Breeze CLI and MCP" socialDescription="Search components, list platforms and generate Go code with the MyGo Breeze CLI and MCP server." path="/cli-mcp" locale={locale} />
      <div className="mx-auto w-[calc(100%-48px)] max-w-[1160px] max-[680px]:w-[calc(100%-34px)]">
        <Link className="mb-12 inline-flex items-center gap-1.5 font-mono text-[11px] text-muted hover:text-ink" to="/">
          <ArrowLeft size={15} /> {zh ? "返回首页" : "Back home"}
        </Link>
        <div className="mb-14 max-w-[760px]">
          <p className="mb-5 font-mono text-[10px] uppercase tracking-[0.09em] text-mint">CLI · MCP</p>
          <h1 className="mb-6 text-[clamp(44px,6vw,76px)] font-bold leading-[1.04] tracking-[-0.075em]">{zh ? "CLI 与 MCP" : "CLI and MCP"}</h1>
          <p className="max-w-[700px] text-[16px] leading-[1.75] text-muted">
            {zh ? "从同一份 catalog.json 搜索组件、列出平台并生成 Go 代码。CLI 适合终端工作流，MCP 服务适合编程 Agent。" : "Search components, list platforms and generate Go code from the same catalog.json. Use the CLI in a terminal and the MCP server with a coding agent."}
          </p>
        </div>

        <div className="grid grid-cols-2 gap-4 max-[780px]:grid-cols-1">
          <article className="rounded-2xl border border-line bg-panel p-7 max-[680px]:p-5">
            <div className="mb-6 flex items-center gap-3"><span className="grid size-10 place-items-center rounded-xl bg-[#eaf7f2] text-mint"><Terminal size={19} /></span><div><p className="m-0 font-mono text-[10px] uppercase tracking-[0.08em] text-dim">CLI</p><h2 className="m-0 text-2xl tracking-[-0.04em]">breeze-ui</h2></div></div>
            <p className="mb-6 text-sm leading-[1.7] text-muted">{zh ? "列出组件和平台，或生成可以直接放进项目的代码片段。" : "List components and platforms, or generate snippets you can drop into a project."}</p>
            <CodeBlock>{`npx breeze-ui list\nnpx breeze-ui add Switch -o ui/switch_snippet.go`}</CodeBlock>
            <p className="mt-5 text-sm text-muted">{zh ? "也可以安装 Go 命令：" : "Or install the Go command:"}</p>
            <CodeBlock>{"go install github.com/yldm-tech/mygo-breeze-ui/cmd/breeze-ui@latest"}</CodeBlock>
            <a className="mt-5 inline-flex items-center gap-1.5 text-sm font-semibold text-mint hover:text-ink" href="https://github.com/yldm-tech/mygo-breeze-ui" target="_blank" rel="noreferrer">{zh ? "查看 CLI 源码" : "View the CLI source"}<ArrowUpRight size={15} /></a>
          </article>

          <article className="rounded-2xl border border-line bg-panel p-7 max-[680px]:p-5">
            <div className="mb-6 flex items-center gap-3"><span className="grid size-10 place-items-center rounded-xl bg-[#eaf7f2] text-mint"><Braces size={19} /></span><div><p className="m-0 font-mono text-[10px] uppercase tracking-[0.08em] text-dim">MCP</p><h2 className="m-0 text-2xl tracking-[-0.04em]">breeze-ui-mcp</h2></div></div>
            <p className="mb-6 text-sm leading-[1.7] text-muted">{zh ? "只读 MCP 服务把 Breeze UI 目录交给编程 Agent，不启动网络监听器。" : "The read-only MCP server gives the Breeze UI catalog to a coding agent without starting a network listener."}</p>
            <CodeBlock>{`bun install\nbun run packages/breeze-mcp/src/index.ts`}</CodeBlock>
            <p className="mt-5 text-sm leading-[1.7] text-muted">{zh ? "服务提供搜索组件、查看组件、列出平台和生成代码片段四个工具。" : "It exposes four tools: search components, inspect a component, list platforms and generate snippets."}</p>
            <CodeBlock>{"breeze-ui-mcp"}</CodeBlock>
            <a className="mt-5 inline-flex items-center gap-1.5 text-sm font-semibold text-mint hover:text-ink" href="https://github.com/yldm-tech/mygo-breeze-ui/tree/main/packages/breeze-mcp" target="_blank" rel="noreferrer">{zh ? "查看 MCP 源码" : "View the MCP source"}<ArrowUpRight size={15} /></a>
          </article>
        </div>

        <section className="mt-4 rounded-2xl border border-line bg-panel-2 p-7 max-[680px]:p-5">
          <div className="mb-6 flex items-center gap-3"><span className="grid size-9 place-items-center rounded-lg bg-white text-mint shadow-[inset_0_0_0_1px_#e2e4e9]"><Braces size={17} /></span><h2 className="m-0 text-2xl tracking-[-0.04em]">{zh ? "同一份目录，四种入口" : "One catalog, four entry points"}</h2></div>
          <div className="grid grid-cols-4 gap-3 max-[780px]:grid-cols-2 max-[480px]:grid-cols-1">
            {(zh ? [["搜索组件", "按名称、分类或说明查找"], ["查看组件", "读取 drawn 或 themed 外观"], ["列出平台", "读取三个平台的版本"], ["生成片段", "输出可直接使用的 Go 代码"]] : [["Search", "Find by name, category or text"], ["Inspect", "Read drawn or themed looks"], ["Platforms", "Read the three platform versions"], ["Snippets", "Generate usable Go code"]]).map(([title, text]) => <div className="rounded-xl border border-line bg-white p-4" key={title}><h3 className="m-0 mb-2 text-sm">{title}</h3><p className="m-0 text-xs leading-[1.6] text-muted">{text}</p></div>)}
          </div>
        </section>
      </div>
    </main>
  );
}
