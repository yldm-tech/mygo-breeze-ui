import catalog from "../../ui/breeze/catalog.json"

type Language = "en" | "zh"
type CatalogItem = { name: string; description?: string; example?: string; value?: string; kind?: string }

const components = catalog.components as CatalogItem[]
const utilities = catalog.utilities as CatalogItem[]
const tokens = catalog.tokens as CatalogItem[]
let language: Language = "en"
let filter = "all"

const copy = {
  en: {
    navComponents: "Components", navInstall: "Install", navDocs: "Docs", navGithub: "GitHub",
    eyebrow: "Native UI for Go, with a breeze.", title: "Build native interfaces with a familiar utility vocabulary.",
    intro: "Breeze UI brings Tailwind-inspired tokens, utilities, and accessible components to MyGo's GPU-rendered native UI. No DOM. No CSS runtime. Just typed Go.",
    start: "Get started", explore: "Explore components", statComponents: "components", statUtilities: "utilities", statRuntime: "runtime dependencies",
    sectionTitle: "A calm layer over native UI", sectionText: "Keep the speed and control of native rendering while giving your team a consistent design language.",
    utilityTitle: "Utility-first, Go-native", utilityText: "Compose typed styles with familiar spacing, layout, color, and typography helpers.",
    componentTitle: "Components with behavior included", componentText: "Buttons, fields, dialogs, tabs, and more delegate focus, keyboard, and accessibility to MyGo's native primitives.",
    toolingTitle: "Tooling that stays out of the way", toolingText: "Use the CLI to discover the catalog, or connect the MCP server to your coding workflow.",
    catalogTitle: "A small, deliberate catalog", catalogText: "Start with the primitives you use every day. Each item is typed, composable, and documented.", all: "All", components: "Components", utilities: "Utilities", tokens: "Tokens",
    installTitle: "Start in one command", installText: "Breeze UI is a Go package. Add it to an existing MyGo app or explore the included gallery.", copy: "Copy", copied: "Copied", docs: "Read the docs", starTitle: "Follow the project", starText: "Breeze UI is open source and evolving in the open.", footer: "Breeze UI for MyGo", source: "Source on GitHub", language: "中文"
  },
  zh: {
    navComponents: "组件", navInstall: "安装", navDocs: "文档", navGithub: "GitHub",
    eyebrow: "为 Go 打造的原生 UI。轻盈如风。", title: "用熟悉的工具样式，构建真正的原生界面。",
    intro: "Breeze UI 将 Tailwind 风格的设计令牌、工具样式和可访问组件带到 MyGo 的 GPU 原生 UI 中。不需要 DOM，不需要运行时 CSS，只需要类型安全的 Go。",
    start: "开始使用", explore: "浏览组件", statComponents: "个组件", statUtilities: "个工具样式", statRuntime: "运行时依赖",
    sectionTitle: "原生 UI 之上的平静一层", sectionText: "保留原生渲染的速度与控制力，同时让团队拥有统一的设计语言。",
    utilityTitle: "Go 原生的 Utility-first", utilityText: "用熟悉的间距、布局、颜色和排版工具组合出类型安全的样式。",
    componentTitle: "自带行为的组件", componentText: "按钮、输入框、对话框、标签页等组件将焦点、键盘和可访问性交给 MyGo 原生基元。",
    toolingTitle: "不会碍事的工具链", toolingText: "用 CLI 浏览组件目录，或把 MCP 服务接入你的编码工作流。",
    catalogTitle: "小而明确的组件目录", catalogText: "从每天都会用到的基础组件开始。每个条目都类型安全、可组合并有文档。", all: "全部", components: "组件", utilities: "工具样式", tokens: "令牌",
    installTitle: "一条命令开始", installText: "Breeze UI 是 Go 包。把它加进现有 MyGo 应用，或运行仓库里的组件展示。", copy: "复制", copied: "已复制", docs: "阅读文档", starTitle: "关注项目", starText: "Breeze UI 开源开发，欢迎一起完善。", footer: "Breeze UI for MyGo", source: "在 GitHub 查看源码", language: "English"
  }
} as const

const t = () => copy[language]
const esc = (value: string) => value.replace(/[&<>"']/g, (char) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" }[char]!))
const icon = (name: string) => `<span class="icon icon-${name}" aria-hidden="true"></span>`

function itemCards() {
  const sets: Array<[string, CatalogItem[]]> = filter === "all" ? [["component", components], ["utility", utilities], ["token", tokens]] : filter === "components" ? [["component", components]] : filter === "utilities" ? [["utility", utilities]] : [["token", tokens]]
  return sets.flatMap(([kind, items]) => items.slice(0, kind === "component" ? 12 : kind === "utility" ? 12 : 8).map((item) => `<article class="catalog-card"><div class="catalog-kind">${kind}</div><h3>${esc(item.name)}</h3><p>${esc(item.description ?? item.value ?? "")}</p>${item.example ? `<code>${esc(item.example)}</code>` : ""}</article>`)).join("")
}

function render() {
  const text = t()
  document.documentElement.lang = language === "zh" ? "zh-CN" : "en"
  document.querySelector<HTMLDivElement>("#app")!.innerHTML = `
    <header class="site-header"><a class="brand" href="#top"><span class="brand-mark">B</span><span>breeze<span class="brand-dot">.</span>ui</span></a><nav><a href="#catalog">${text.navComponents}</a><a href="#install">${text.navInstall}</a><a href="https://github.com/yldm-tech/mygo-breeze-ui/tree/main/docs" target="_blank">${text.navDocs}</a><a href="https://github.com/yldm-tech/mygo-breeze-ui" target="_blank">${text.navGithub} ${icon("arrow-up-right")}</a></nav><button class="language" data-action="language">${text.language}</button><button class="menu-button" data-action="menu" aria-label="Menu">☰</button></header>
    <main id="top">
      <section class="hero container"><div class="hero-copy"><p class="eyebrow"><span class="pulse"></span>${text.eyebrow}</p><h1>${text.title}</h1><p class="hero-intro">${text.intro}</p><div class="hero-actions"><a class="button button-primary" href="#install">${text.start} ${icon("arrow-right")}</a><a class="button button-ghost" href="#catalog">${text.explore}</a></div><div class="stats"><div><strong>${components.length}</strong><span>${text.statComponents}</span></div><div><strong>${utilities.length}</strong><span>${text.statUtilities}</span></div><div><strong>0</strong><span>${text.statRuntime}</span></div></div></div><div class="hero-code"><div class="window-bar"><span></span><span></span><span></span><small>main.go</small></div><pre><span class="code-keyword">package</span> main

<span class="code-keyword">import</span> (
  <span class="code-string">"github.com/egoist/mygo/ui"</span>
  breeze <span class="code-string">"github.com/yldm-tech/mygo-breeze-ui/ui/breeze"</span>
)

<span class="code-keyword">func</span> view(c *ui.Context) {
  breeze.<span class="code-function">Card</span>(c, <span class="code-keyword">func</span>() {
    breeze.<span class="code-function">Text</span>(c, <span class="code-string">"Welcome home"</span>)
    <span class="code-keyword">if</span> breeze.<span class="code-function">Button</span>(c, <span class="code-string">"Continue"</span>, breeze.ButtonPrimary).<span class="code-function">Clicked</span>() {
      <span class="code-comment">// handle the action</span>
    }
  })
}</pre><div class="code-status"><span class="status-dot"></span> compiled with Go <span class="status-spacer"></span> native renderer</div></div></section>
      <section class="signal"><div class="container signal-grid"><div><span class="signal-number">01</span><h2>${text.sectionTitle}</h2><p>${text.sectionText}</p></div><div class="signal-line"></div><div class="signal-note"><span class="quote-mark">“</span><p>Design tokens should make the right thing easy.</p><small>— Breeze UI principles</small></div></div></section>
      <section class="features container"><article class="feature"><div class="feature-icon mint">${icon("spark")}</div><h3>${text.utilityTitle}</h3><p>${text.utilityText}</p><code>breeze.Row(c, breeze.Gap(4), breeze.P(6))</code></article><article class="feature"><div class="feature-icon orange">${icon("layers")}</div><h3>${text.componentTitle}</h3><p>${text.componentText}</p><code>breeze.Dialog(c, &amp;open, "Hello", body, footer)</code></article><article class="feature"><div class="feature-icon violet">${icon("terminal")}</div><h3>${text.toolingTitle}</h3><p>${text.toolingText}</p><code>breeze-ui add Button -o button.go</code></article></section>
      <section class="catalog-section" id="catalog"><div class="container"><div class="section-heading"><div><p class="eyebrow">${text.navComponents}</p><h2>${text.catalogTitle}</h2><p>${text.catalogText}</p></div><div class="filters"><button class="filter ${filter === "all" ? "active" : ""}" data-filter="all">${text.all}</button><button class="filter ${filter === "components" ? "active" : ""}" data-filter="components">${text.components}</button><button class="filter ${filter === "utilities" ? "active" : ""}" data-filter="utilities">${text.utilities}</button><button class="filter ${filter === "tokens" ? "active" : ""}" data-filter="tokens">${text.tokens}</button></div></div><div class="catalog-grid">${itemCards()}</div></div></section>
      <section class="install-section" id="install"><div class="container install-grid"><div><p class="eyebrow">${text.navInstall}</p><h2>${text.installTitle}</h2><p>${text.installText}</p><a class="text-link" href="https://github.com/yldm-tech/mygo-breeze-ui/blob/main/docs/breeze-ui.md" target="_blank">${text.docs} ${icon("arrow-right")}</a></div><div class="install-card"><div class="install-tabs"><span class="active">Go</span><span>CLI</span><span>MCP</span></div><div class="install-command"><code>go get github.com/yldm-tech/mygo-breeze-ui</code><button data-action="copy" data-copy="go get github.com/yldm-tech/mygo-breeze-ui">${text.copy}</button></div><div class="install-command muted-command"><code>go run ./examples/breeze-ui</code><button data-action="copy" data-copy="go run ./examples/breeze-ui">${text.copy}</button></div></div></div></section>
      <section class="star-section"><div class="container star-card"><div><p class="eyebrow">GitHub</p><h2>${text.starTitle}</h2><p>${text.starText}</p><a class="button button-primary" href="https://github.com/yldm-tech/mygo-breeze-ui" target="_blank">${text.source} ${icon("arrow-up-right")}</a></div><div class="star-visual"><div class="star-grid"></div><div class="star-path"></div><span class="star-label">★ 0</span></div></div></section>
    </main><footer class="site-footer container"><a class="brand" href="#top"><span class="brand-mark">B</span><span>breeze<span class="brand-dot">.</span>ui</span></a><span>${text.footer} · MIT License</span><a href="https://github.com/yldm-tech/mygo-breeze-ui" target="_blank">${text.source} ${icon("arrow-up-right")}</a></footer>`
  bind()
}

function bind() {
  document.querySelectorAll<HTMLElement>("[data-filter]").forEach((button) => button.addEventListener("click", () => { filter = button.dataset.filter ?? "all"; render() }))
  document.querySelector<HTMLElement>("[data-action=language]")?.addEventListener("click", () => { language = language === "en" ? "zh" : "en"; render() })
  document.querySelector<HTMLElement>("[data-action=menu]")?.addEventListener("click", () => document.querySelector("nav")?.classList.toggle("open"))
  document.querySelectorAll<HTMLElement>("[data-action=copy]").forEach((button) => button.addEventListener("click", async () => { await navigator.clipboard?.writeText(button.dataset.copy ?? ""); button.textContent = t().copied; setTimeout(() => { button.textContent = t().copy }, 1400) }))
}

render()
