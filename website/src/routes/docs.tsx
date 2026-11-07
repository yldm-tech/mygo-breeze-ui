import { useState } from "react";
import { Check, ChevronDown, ChevronRight, Copy, Info, Menu, Search, X } from "lucide-react";
import { createRoute, Link } from "@tanstack/react-router";
import { RootRoute } from "./root";
import { Seo } from "../components/Seo";
import { GoCodeBlock } from "../components/CodeHighlight";
import { useUiStore } from "../store";
import { categoryLabels, docsComponents, localizedCategory, localizedDescription } from "./components-docs";
import type { Locale } from "../store";

const afterCode = `func (app *settings) view(c *ui.Context) {
    breeze.Use(c) // 平台的颜色和尺寸，ui 包的组件也会跟着变
    ui.Column(c).Padding(16).Gap(12).Children(func() {
        breeze.Checkbox(c, &app.remember, "Remember me")
        breeze.Switch(c, &app.sync).Label("Sync")
        ui.Row(c).Gap(8).Children(func() {
            breeze.Button(c, "Cancel")
            if breeze.PrimaryButton(c, "Save").Clicked() {
                app.save()
            }
        })
    })
}`;

const beforeCode = `func (app *settings) view(c *ui.Context) {
    ui.Column(c).Padding(16).Gap(12).Children(func() {
        ui.Checkbox(c, &app.remember, "Remember me")
        ui.Switch(c, &app.sync).Label("Sync")
        ui.Row(c).Gap(8).Children(func() {
            ui.Button(c, "Cancel")
            if ui.PrimaryButton(c, "Save").Clicked() {
                app.save()
            }
        })
    })
}`;

const previewCode = `for _, p := range breeze.Platforms {
    breeze.Preview(c, p, func() {
        ui.Themed(c, breeze.Theme(c, p)).Grow(1).Padding(16).Gap(12).Children(func() {
            breeze.Button(c, "Cancel")
            breeze.Switch(c, &app.sync).Label("Sync")
        })
    })
}`;

function CopyButton({ value }: { value: string }) {
  const [copied, setCopied] = useState(false);
  return (
    <button
      type="button"
      className="inline-flex min-h-9 shrink-0 items-center rounded-lg bg-white px-3 text-[13px] text-[#555C69] shadow-[inset_0_0_0_1px_#E2E4E9] hover:text-[#0B7A63]"
      onClick={() => {
        void navigator.clipboard?.writeText(value);
        setCopied(true);
        window.setTimeout(() => setCopied(false), 1400);
      }}
    >
      {copied ? <Check size={14} /> : <Copy size={14} />}
      <span className="ml-1.5">{copied ? "已复制" : "复制"}</span>
    </button>
  );
}

function DocsSidebar({ open, onClose, locale }: { open: boolean; onClose: () => void; locale: Locale }) {
  const groups = [
    {
      title: "开始",
      items: [
        ["快速开始", "#top", true],
        ["安装", "#install", false],
        ["第一个视图", "#first-view", false],
        ["Storybook", "#storybook", false],
      ],
    },
    {
      title: "指南",
      items: [["平台与 Preview", "#platforms", false], ["主题", "#themes", false], ["已知限制", "#limits", false]],
    },
  ] as const;
  const categoryOrder = ["actions", "choices", "text", "layout", "navigation", "data", "feedback", "overlays"];
  const [query, setQuery] = useState("");
  const [expanded, setExpanded] = useState<Record<string, boolean>>({});
  const [activeSection, setActiveSection] = useState(() => (typeof window !== "undefined" ? window.location.hash.slice(1) || "top" : "top"));
  const visible = query.trim().toLowerCase();

  return (
    <aside
      id="docs-side"
      className={`${open ? "block" : "hidden"} w-full shrink-0 border-b border-[#E2E4E9] bg-white px-4 py-5 md:sticky md:top-16 md:block md:max-h-[calc(100vh-4rem)] md:w-[280px] md:overflow-y-auto md:border-b-0 md:border-r md:px-5 md:py-10`}
    >
      <nav className="flex flex-col gap-7 text-[13px]" aria-label="文档目录">
        {groups.map((group) => (
          <div key={group.title} className="flex flex-col gap-2">
            <p className="m-0 px-2 text-[12px] font-semibold uppercase tracking-[0.08em] text-[#8A90A0]">{group.title}</p>
            {group.items.map(([label, href, current]) => {
              const selected = href.startsWith("#") ? activeSection === href.slice(1) : current;
              const className = `flex min-h-[34px] items-center rounded-lg px-2.5 ${selected ? "bg-[#E3F3EE] font-semibold text-[#075A49]" : "text-[#3A404C] hover:bg-[#F8F9FB] hover:text-[#0B7A63]"}`;
              return <a key={label} className={className} href={href} onClick={() => { if (href.startsWith("#")) setActiveSection(href.slice(1)); onClose(); }} aria-current={selected ? "page" : undefined}>{label}</a>;
            })}
          </div>
        ))}
        <div className="flex flex-col gap-3">
          <p className="m-0 px-2 text-[12px] font-semibold uppercase tracking-[0.08em] text-[#8A90A0]">组件 · {docsComponents.length}</p>
          <label className="flex h-10 items-center gap-2 rounded-xl bg-white px-3 text-[#6B7280] shadow-[inset_0_0_0_1px_#E2E4E9]"><Search size={15} /><span className="sr-only">搜索组件</span><input className="min-w-0 flex-1 border-0 bg-transparent text-sm outline-none" value={query} onChange={(event) => setQuery(event.target.value)} placeholder="搜索组件" /></label>
          <div className="flex flex-col gap-1">
            {categoryOrder.map((category) => {
              const items = docsComponents.filter((item) => item.category === category && (!visible || `${item.name} ${localizedDescription(item, locale)}`.toLowerCase().includes(visible)));
              if (!items.length) return null;
              const isExpanded = Boolean(expanded[category] || visible);
              return <div key={category}>
                <button type="button" className="flex min-h-[36px] w-full items-center gap-2 rounded-lg px-2.5 text-left text-[#555C69] hover:bg-[#F8F9FB]" aria-expanded={isExpanded} onClick={() => setExpanded((value) => ({ ...value, [category]: !value[category] }))}>
                  {isExpanded ? <ChevronDown size={15} /> : <ChevronRight size={15} />}<span className="flex-1">{localizedCategory(category, locale)}</span><span className="text-xs text-[#8A90A0]">{docsComponents.filter((item) => item.category === category).length}</span>
                </button>
                {isExpanded && <div className="ml-5 flex flex-col gap-0.5 border-l border-[#E2E4E9] pl-2">{items.map((item) => <Link key={item.name} className="flex min-h-[32px] items-center rounded-lg px-2.5 text-[#3A404C] hover:bg-[#F8F9FB] hover:text-[#0B7A63]" to="/docs/components/$component" params={{ component: item.name }} onClick={onClose}>{item.name}</Link>)}</div>}
              </div>;
            })}
          </div>
        </div>
      </nav>
    </aside>
  );
}

function InlineCode({ children }: { children: React.ReactNode }) {
  return <code className="rounded-[5px] bg-[#F4F5F7] px-1.5 py-0.5 font-mono text-[14px] text-[#15171C]">{children}</code>;
}

function DocsPage() {
  const { locale } = useUiStore();
  const [menuOpen, setMenuOpen] = useState(false);
  const [diff, setDiff] = useState<"after" | "before">("after");

  return (
    <>
      <Seo
        title="Breeze UI 文档 · 快速开始"
        description="Breeze UI 快速开始、平台 Preview、主题和 Storybook 指南。"
        socialTitle="Breeze UI documentation · Quick start"
        socialDescription="Quick start, platform previews, themes and Storybook for Breeze UI."
        path="/docs"
        locale={locale}
      />
      <main className="min-h-screen bg-white text-[#15171C]" style={{ fontFamily: 'Inter, "Noto Sans SC", "PingFang SC", system-ui, sans-serif' }}>
        <div className="sticky top-16 z-20 flex items-center gap-2 border-b border-[#E2E4E9] bg-[#FAFBFC] px-4 py-2 md:hidden">
          <button
            type="button"
            className="inline-flex min-h-11 items-center gap-2 rounded-[10px] bg-white px-3.5 text-sm font-semibold shadow-[inset_0_0_0_1px_#E2E4E9]"
            onClick={() => setMenuOpen((value) => !value)}
            aria-expanded={menuOpen}
            aria-controls="docs-side"
          >
            {menuOpen ? <X size={18} /> : <Menu size={18} />}
            {menuOpen ? "收起目录" : "目录"}
          </button>
          <span className="min-w-0 overflow-hidden text-ellipsis whitespace-nowrap text-sm text-[#555C69]">开始 / 快速开始</span>
        </div>

        <div className="mx-auto flex max-w-[1440px] flex-wrap items-start">
          <DocsSidebar open={menuOpen} onClose={() => setMenuOpen(false)} locale={locale} />
          <article id="top" className="min-w-0 flex-[999_1_560px] px-4 pb-16 pt-6 md:px-12 md:pb-20 md:pt-10">
            <div className="mx-auto flex max-w-[780px] flex-col gap-12">
              <div className="flex flex-col gap-3.5">
                <nav className="flex gap-2 text-sm text-[#6B7280]" aria-label="面包屑">
                  <span>开始</span><span aria-hidden="true">/</span><span className="text-[#15171C]">快速开始</span>
                </nav>
                <h1 className="m-0 text-[40px] font-semibold leading-[1.15] tracking-[-0.02em] max-[680px]:text-[34px]">快速开始</h1>
                <p className="m-0 text-[18px] leading-[1.65] text-[#555C69] max-[680px]:text-[16px]">Breeze UI 让 MyGo 的原生 UI 组件按应用所在的系统绘制：Mac 上是 macOS 的控件，Windows 上是 Windows 11 的，Linux 上是 GNOME libadwaita 的。五分钟内把它接进你的应用。</p>
              </div>

              <ol className="m-0 grid list-none grid-cols-[repeat(auto-fit,minmax(min(100%,200px),1fr))] gap-3 p-0">
                {[
                  ["1", "安装 Go 包", "一条 go get，需要 MyGo v0.7.4+", "#install"],
                  ["2", "替换组件", "ui.Button → breeze.Button", "#first-view"],
                  ["3", "三平台对照", "用 Preview 并排检查", "#platforms"],
                ].map(([number, label, note, href]) => (
                  <li key={number} className="flex flex-col gap-1.5 rounded-[14px] bg-[#F8F9FB] p-[18px] shadow-[inset_0_0_0_1px_#E2E4E9]">
                    <span className="flex size-7 items-center justify-center rounded-full bg-[#0B7A63] text-[13px] font-bold text-white">{number}</span>
                    <a className="font-semibold text-[#15171C]" href={href}>{label}</a>
                    <span className="text-sm text-[#555C69]">{note}</span>
                  </li>
                ))}
              </ol>

              <section id="install" className="flex scroll-mt-24 flex-col gap-3.5">
                <h2 className="m-0 text-[26px] font-semibold tracking-[-0.01em]">安装</h2>
                <p className="m-0 text-[#3A404C]">在你的 MyGo 项目里添加依赖。Breeze UI 需要 MyGo v0.7.4 或更高版本。</p>
                <div className="flex items-center gap-3 rounded-xl bg-[#F8F9FB] p-1.5 pl-4 shadow-[inset_0_0_0_1px_#E2E4E9]">
                  <span className="font-mono text-sm text-[#8A90A0]">$</span>
                  <code className="min-w-0 flex-1 overflow-x-auto whitespace-nowrap font-mono text-sm">go get github.com/yldm-tech/mygo-breeze-ui</code>
                  <CopyButton value="go get github.com/yldm-tech/mygo-breeze-ui" />
                </div>
                <p className="m-0 text-[#3A404C]">想先看看效果？克隆仓库后运行组件画廊，三个平台会并排显示：</p>
                <GoCodeBlock code="go run ./examples/breeze-ui" className="m-0 overflow-x-auto rounded-xl bg-[#F8F9FB] p-3.5 px-4 font-mono text-sm shadow-[inset_0_0_0_1px_#E2E4E9]" />
              </section>

              <section id="first-view" className="flex scroll-mt-24 flex-col gap-3.5">
                <h2 className="m-0 text-[26px] font-semibold tracking-[-0.01em]">第一个视图</h2>
                <p className="m-0 text-[#3A404C]">在视图开头调用一次 <InlineCode>breeze.Use(c)</InlineCode>，然后把 <InlineCode>ui.</InlineCode> 换成 <InlineCode>breeze.</InlineCode>。组件返回 <InlineCode>ui.Element</InlineCode>、响应 <InlineCode>Clicked</InlineCode>，也能像 ui 的组件一样继续链式添加样式。</p>
                <div className="overflow-hidden rounded-xl shadow-[0_0_0_1px_#E2E4E9]">
                  <div className="flex items-center justify-between border-b border-[#E2E4E9] bg-[#FAFBFC] px-2 pl-1">
                    <div className="flex" role="tablist" aria-label="对比">
                      <button type="button" role="tab" aria-selected={diff === "after"} className={`min-h-11 border-0 bg-transparent px-3 text-[13px] font-semibold ${diff === "after" ? "text-[#15171C] shadow-[inset_0_-2px_#0B7A63]" : "text-[#6B7280]"}`} onClick={() => setDiff("after")}>使用 Breeze UI</button>
                      <button type="button" role="tab" aria-selected={diff === "before"} className={`min-h-11 border-0 bg-transparent px-3 text-[13px] font-semibold ${diff === "before" ? "text-[#15171C] shadow-[inset_0_-2px_#0B7A63]" : "text-[#6B7280]"}`} onClick={() => setDiff("before")}>原来的 ui</button>
                    </div>
                    <span className="text-xs text-[#555C69]">settings.go</span>
                  </div>
                  <GoCodeBlock code={diff === "after" ? afterCode : beforeCode} className="m-0 overflow-x-auto p-4 px-5 font-mono text-[13.5px] leading-[1.75] text-[#15171C]" />
                </div>
                <div className="flex gap-3 rounded-xl bg-[#EEF5FF] p-4 px-[18px] text-[15px] text-[#173B6B]">
                  <Info size={20} className="mt-0.5 shrink-0" />
                  <p className="m-0">指针、键盘、焦点和无障碍来自 MyGo 的组件基座（<code className="font-mono text-sm">ui.ButtonBase</code>、<code className="font-mono text-sm">ui.CheckboxBase</code>、<code className="font-mono text-sm">ui.SwitchBase</code>），所以替换之后交互行为不变。</p>
                </div>
              </section>

              <section id="platforms" className="flex scroll-mt-24 flex-col gap-3.5">
                <h2 className="m-0 text-[26px] font-semibold tracking-[-0.01em]">平台与 Preview</h2>
                <p className="m-0 text-[#3A404C]"><InlineCode>breeze.Current()</InlineCode> 是应用运行的平台：macOS 上是 <InlineCode>Mac</InlineCode>，Windows 上是 <InlineCode>Windows</InlineCode>，其他系统是 <InlineCode>Linux</InlineCode>。每个组件在构建时读取它。</p>
                <div className="overflow-x-auto rounded-xl shadow-[0_0_0_1px_#E2E4E9]">
                  <table className="w-full min-w-[520px] border-collapse text-sm">
                    <thead><tr className="bg-[#FAFBFC] text-left"><th className="border-b border-[#E2E4E9] px-4 py-3 font-semibold text-[#555C69]">函数</th><th className="border-b border-[#E2E4E9] px-4 py-3 font-semibold text-[#555C69]">作用</th></tr></thead>
                    <tbody>
                      {[
                        ["SetPlatform(breeze.Windows)", "让所有窗口按另一个平台绘制，适合截图。"],
                        ["SetPlatform(breeze.Auto)", "重新跟随系统。"],
                        ["Preview(c, p, body)", "只把 body 构建的部分按平台 p 绘制，在当前窗口里。"],
                      ].map(([name, description], index) => <tr key={name}><th scope="row" className={`${index < 2 ? "border-b border-[#EEF0F3]" : ""} whitespace-nowrap px-4 py-3.5 text-left font-mono font-medium text-[#0B7A63]`}>{name}</th><td className={`${index < 2 ? "border-b border-[#EEF0F3]" : ""} px-4 py-3.5 text-[#3A404C]`}>{description}</td></tr>)}
                    </tbody>
                  </table>
                </div>
                <p className="m-0 text-[#3A404C]">配合 <InlineCode>ui.Themed(c, breeze.Theme(c, p))</InlineCode>，就能像组件画廊那样把三个平台并排：</p>
                <GoCodeBlock code={previewCode} className="m-0 overflow-x-auto rounded-xl bg-[#F8F9FB] p-4 px-5 font-mono text-[13.5px] leading-[1.75] shadow-[inset_0_0_0_1px_#E2E4E9]" />
              </section>

              <section id="themes" className="flex scroll-mt-24 flex-col gap-3.5">
                <h2 className="m-0 text-[26px] font-semibold tracking-[-0.01em]">主题</h2>
                <p className="m-0 text-[#3A404C]"><InlineCode>breeze.Theme(c, p)</InlineCode> 返回平台 p 在当前外观（浅色或深色）下的主题：窗口背景、文字、分隔线和强调色。<InlineCode>breeze.Use(c)</InlineCode> 用 <InlineCode>c.SetTheme</InlineCode> 把它设给窗口。它跟随桌面偏好：</p>
                <div className="grid grid-cols-[repeat(auto-fit,minmax(min(100%,220px),1fr))] gap-3">
                  <div className="flex flex-col gap-2 rounded-[14px] p-[18px] shadow-[inset_0_0_0_1px_#E2E4E9]"><strong className="font-semibold">强调色</strong><span className="text-sm text-[#555C69]">用户选的强调色替换平台默认值。</span><div className="mt-1 flex gap-1.5"><span className="size-[18px] rounded-[5px] bg-[#007AFF]" /><span className="size-[18px] rounded-[5px] bg-[#0078D4]" /><span className="size-[18px] rounded-[5px] bg-[#3584E4]" /></div></div>
                  <div className="flex flex-col gap-2 rounded-[14px] p-[18px] shadow-[inset_0_0_0_1px_#E2E4E9]"><strong className="font-semibold">更高对比度</strong><span className="text-sm text-[#555C69]">分隔线和次要文字加深（深色下变亮），焦点环不透明。</span></div>
                  <div className="flex flex-col gap-2 rounded-[14px] p-[18px] shadow-[inset_0_0_0_1px_#E2E4E9]"><strong className="font-semibold">文字大小</strong><span className="text-sm text-[#555C69]">字号随 Windows 和 GNOME 的文字缩放变大。</span></div>
                </div>
              </section>

              <section id="storybook" className="flex scroll-mt-24 flex-col gap-3.5">
                <h2 className="m-0 text-[26px] font-semibold tracking-[-0.01em]">Storybook</h2>
                <div className="flex flex-wrap items-center gap-6 rounded-2xl bg-[#F8F9FB] p-6 shadow-[inset_0_0_0_1px_#E2E4E9]">
                  <div className="flex min-w-[min(100%,300px)] flex-1 flex-col gap-2.5"><p className="m-0 text-[#3A404C]">每个组件都有自己的 story：选 macOS、Windows 11、GNOME 或三者并排，切换浅色和深色，用控件调整参数，并查看构建它的 Go 代码。网页版由 WebAssembly 编译，MyGo 用 WebGL 2 在页面里绘制。</p><GoCodeBlock code="go run ./examples/storybook" className="m-0 overflow-x-auto rounded-[10px] bg-white p-3 px-3.5 font-mono text-[13px] shadow-[inset_0_0_0_1px_#E2E4E9]" /></div>
                  <a className="inline-flex min-h-11 items-center rounded-[10px] bg-[#15171C] px-[18px] text-sm font-semibold text-white" href="https://breeze-ui.yldm.ai/storybook/" target="_blank" rel="noreferrer">打开 Storybook →</a>
                </div>
              </section>

              <section id="limits" className="flex scroll-mt-24 flex-col gap-3.5">
                <h2 className="m-0 text-[26px] font-semibold tracking-[-0.01em]">已知限制</h2>
                <p className="m-0 text-[#3A404C]">完整列表见 <a className="text-[#1F5FBF] underline" href="https://github.com/yldm-tech/mygo-breeze-ui/blob/main/docs/breeze-ui.md#limits" target="_blank" rel="noreferrer">breeze-ui.md</a>，最常遇到的几条：</p>
                <ul className="m-0 flex list-disc flex-col gap-2 pl-5 text-[#3A404C]"><li><strong>尺寸不同。</strong>macOS 按钮 24 pt，Windows 32，GNOME 34。同一布局在 Mac 上更紧凑，请按内容定尺寸，不要写死高度。</li><li><strong>非活动窗口。</strong>MyGo v0.7.4 不告诉视图窗口是否激活，所以总是绘制激活窗口的样子。</li><li><strong>Liquid Glass。</strong>macOS 27 的菜单、弹出层、工具提示是玻璃材质的独立窗口，无法离屏渲染，是仿照绘制而非测量。</li><li><strong>图标与字体。</strong>不附带 SF Symbols、Segoe Fluent Icons 或 Adwaita 图标；勾选、箭头以路径绘制，文字使用系统字体。</li></ul>
              </section>

              <section id="next" className="flex scroll-mt-24 flex-col gap-3.5">
                <h2 className="m-0 text-[26px] font-semibold tracking-[-0.01em]">下一步</h2>
                <div className="grid grid-cols-[repeat(auto-fit,minmax(min(100%,240px),1fr))] gap-3"><Link to="/docs/components" className="flex flex-col gap-1.5 rounded-xl p-[18px] text-[#15171C] shadow-[inset_0_0_0_1px_#E2E4E9]"><span className="font-semibold">浏览组件 →</span><span className="text-sm text-[#555C69]">94 个组件，每个都写明在三个平台上的样子。</span></Link><Link to="/cli-mcp" className="flex flex-col gap-1.5 rounded-xl p-[18px] text-[#15171C] shadow-[inset_0_0_0_1px_#E2E4E9]"><span className="font-semibold">CLI 与 MCP →</span><span className="text-sm text-[#555C69]">生成代码片段，或把目录接入编程 Agent。</span></Link></div>
              </section>

              <p className="m-0 border-t border-[#E2E4E9] pt-6 text-[13px] text-[#6B7280]"><a className="text-[#1F5FBF]" href="https://github.com/yldm-tech/mygo-breeze-ui/blob/main/docs/breeze-ui.md" target="_blank" rel="noreferrer">在 GitHub 上编辑此页</a></p>
            </div>
          </article>
          <aside className="hidden w-[220px] shrink-0 px-0 pb-10 pr-6 pt-10 lg:block" aria-label="本页内容">
            <div className="sticky top-24 flex flex-col gap-1 text-[13px]">
              <span className="pb-1.5 font-semibold">本页内容</span>
              {[["安装", "#install"], ["第一个视图", "#first-view"], ["平台与 Preview", "#platforms"], ["主题", "#themes"], ["Storybook", "#storybook"], ["已知限制", "#limits"], ["下一步", "#next"]].map(([label, href], index) => <a key={label} className={`${index === 0 ? "text-[#0B7A63] shadow-[inset_2px_0_#0B7A63]" : "text-[#555C69] shadow-[inset_1px_0_#E2E4E9]"} px-0 py-1 pl-2.5`} href={href}>{label}</a>)}
            </div>
          </aside>
        </div>
      </main>
    </>
  );
}

export const DocsRoute = createRoute({
  getParentRoute: () => RootRoute,
  path: "/docs",
  component: DocsPage,
});
