import { createRoute } from "@tanstack/react-router";
import { motion } from "motion/react";
import { useQuery } from "@tanstack/react-query";
import { Check, Clipboard, Code2, ExternalLink, Layers3, Search, Sparkles, Terminal, ArrowRight } from "lucide-react";
import { useMemo, useState } from "react";
import type { ReactNode } from "react";
import { RootRoute } from "./root";
import { getCatalog, type CatalogEntry } from "../lib/catalog";
import { useUiStore } from "../store";
import { Seo } from "../components/Seo";

export const HomeRoute = createRoute({ getParentRoute: () => RootRoute, path: "/", component: HomePage });

type Text = {
  h1a: string; h1b: string; intro: string; quick: string; storybook: string; copied: string; stat: string; demoCaption: string; setPlatform: string; preview: string; toolStorybook: string; toolCli: string; toolMcp: string;
  platformTitle: string; platformText: string; galleryTitle: string; galleryText: string; galleryCaption: string; platformMethod: string; componentTitle: string; componentText: string;
  themeTitle: string; themeText: string; accent: string; accentText: string; contrast: string; contrastText: string; size: string; sizeText: string;
  toolingTitle: string; toolingText: string; catalogTitle: string; catalogText: string; ctaTitle: string; ctaText: string; github: string;
  docs: string; all: string; search: string; shown: string; empty: string; byPlatform: string; light: string; dark: string;
};
const texts: Record<"zh" | "en", Text> = {
  zh: {
    h1a: "让 MyGo 组件", h1b: "长得像它所在的桌面。",
    intro: "Breeze UI 用 macOS、Windows 11 和 GNOME 各自的原生样式来绘制 MyGo 组件。名称和签名与 ui 包完全一致，把 ui.Button 换成 breeze.Button 就行，其余代码一行都不用改。",
    quick: "快速开始", storybook: "打开 Storybook", copied: "已复制", stat: "{n} 个组件，其中 {drawn} 个按平台自绘，其余套用平台主题。", demoCaption: "同一段视图代码 · 按 AppKit 绘制",
    platformTitle: "三个平台的对照", galleryTitle: "组件画廊截图", galleryText: "组件画廊 examples/breeze-ui 的截图：同一个视图，用 Preview 分别按 macOS、Windows 和 Linux 绘制。", galleryCaption: "运行 go run ./examples/breeze-ui，可以在本机看到同样的画面。", platformMethod: "", platformText: "每个平台的尺寸与颜色，都取自 AppKit、WinUI 3 和 libadwaita 的实际渲染或主题资源，再与原生控件逐像素对比。",
    componentTitle: "代码示例", componentText: "指针、键盘、焦点和无障碍由 MyGo 的组件基座提供，Breeze UI 只负责让它们看起来属于这个桌面。",
    themeTitle: "跟着用户的桌面偏好走。", themeText: "和 MyGo 的默认主题一样，Breeze UI 读取系统设置，不需要你额外写一行代码。",
    accent: "强调色", accentText: "用户选的强调色替换平台默认值；Windows 上按 WinUI 的方式取深浅。", contrast: "更高对比度", contrastText: "分隔线和次要文字加深，焦点环变为不透明。", size: "文字大小", sizeText: "字号随 Windows 和 GNOME 的文字缩放一起变大。",
    toolingTitle: "CLI、MCP 和 Storybook", toolingText: "Go 包、CLI、MCP 服务和网站读取同一份 catalog.json，名称和示例始终一致。",
    catalogTitle: "组件", catalogText: "从按钮、开关到日期选择器、设置行和标题栏。drawn 由 Breeze UI 按平台自绘；themed 是 MyGo 自带组件套用平台主题。",
    ctaTitle: "一个包，三个平台。", ctaText: "把 Breeze UI 加进 MyGo 应用并调用 breeze.Use(c)，或者运行组件画廊，并排看每个控件在 macOS、Windows 11 和 GNOME 上的样子。", github: "在 GitHub 上查看",
    docs: "阅读文档", setPlatform: "截图时强制一个平台，设回 Auto 即跟随系统。", preview: "只把一段界面按另一个平台绘制，用来并排对比三个平台。", toolStorybook: "每个组件都有可运行的交互示例，可以查看三种平台外观。", toolCli: "搜索组件、列出平台、生成代码片段，把目录接入你的工作流。", toolMcp: "把 Breeze UI 的组件目录直接交给你的编程 Agent。", all: "全部", search: "搜索组件，例如 Date", shown: "显示 {n} / 94", empty: "没有找到匹配的组件。", byPlatform: "按平台自绘", light: "浅色", dark: "深色",
  },
  en: {
    h1a: "MyGo widgets that look like", h1b: "the desktop they run on.",
    intro: "Breeze UI draws MyGo widgets the way macOS, Windows 11 and GNOME draw their own controls. Names and signatures match MyGo's ui package: swap ui.Button for breeze.Button and keep the rest of the view.",
    quick: "Quick start", storybook: "Open Storybook", copied: "Copied", stat: "{n} components, {drawn} drawn per platform and the rest themed.", demoCaption: "Same view code · drawn as AppKit",
    platformTitle: "How each platform is measured", galleryTitle: "Gallery screenshot", galleryText: "A screenshot of the gallery in examples/breeze-ui: one view, drawn with Preview as macOS, Windows and Linux.", galleryCaption: "Run go run ./examples/breeze-ui to see the same window on your machine.", platformMethod: "", platformText: "Sizes and colors come from AppKit, WinUI 3 and libadwaita rendering or theme resources, compared pixel by pixel with native controls.",
    componentTitle: "Example", componentText: "Pointer, keyboard, focus and accessibility come from MyGo's widget bases. Breeze UI only makes them look at home on the desktop.",
    themeTitle: "Follows the user's desktop.", themeText: "Like MyGo's default theme, Breeze UI reads system preferences. No extra code is needed.",
    accent: "Accent color", accentText: "The user's accent replaces the platform default; Windows shades follow WinUI.", contrast: "Higher contrast", contrastText: "Separators and secondary text get stronger, and the focus ring becomes opaque.", size: "Text size", sizeText: "Text scales with Windows and GNOME accessibility settings.",
    toolingTitle: "CLI, MCP and Storybook", toolingText: "The Go package, CLI, MCP server and website read the same catalog.json, so names and examples stay in sync.",
    catalogTitle: "Components", catalogText: "From buttons and switches to date pickers, settings rows and header bars. Drawn ones are rendered by Breeze UI per platform; themed uses MyGo's widget with the platform theme.",
    ctaTitle: "One package. Three platforms.", ctaText: "Add Breeze UI to a MyGo app and call breeze.Use(c), or run the gallery to see every control as macOS, Windows 11 and GNOME draw it side by side.", github: "View on GitHub",
    docs: "Read the docs", setPlatform: "Fix every window to one platform for screenshots, then return to Auto to follow the system.", preview: "Draw just a subtree as another platform to compare the three side by side.", toolStorybook: "Runnable examples for every component, with each platform look visible.", toolCli: "Search components, list platforms and generate snippets from the same catalog.", toolMcp: "Give the Breeze UI catalog directly to your coding agent.", all: "All", search: "Search components, e.g. Date", shown: "Showing {n} / 94", empty: "No component matches that search.", byPlatform: "drawn per platform", light: "Light", dark: "Dark",
  },
};

function Mark() {
  return <svg viewBox="0 0 64 64" width="22" height="22" aria-hidden="true"><rect x="4" y="6" width="36" height="28" rx="9" fill="#A6E1CD"/><rect x="14" y="17" width="36" height="28" rx="4" fill="#3FAE8C"/><rect x="24" y="28" width="36" height="28" rx="7" fill="#0B7A63"/><circle cx="31" cy="35" r="1.8" fill="#fff"/><circle cx="37" cy="35" r="1.8" fill="#fff" fillOpacity=".7"/><circle cx="43" cy="35" r="1.8" fill="#fff" fillOpacity=".45"/></svg>;
}

function HeroDemo({ t }: { t: Text }) {
  const [platform, setPlatform] = useState("macOS");
  const [dark, setDark] = useState(false);
  const platformClass = platform === "Windows" ? "windows" : platform === "GNOME" ? "gnome" : "macos";
  return (
    <div className={`demo-stage demo-stage-${platformClass}`}>
      <div className="demo-controls">
        <div className="segmented" role="radiogroup" aria-label="Platform">
          {["macOS", "Windows", "GNOME"].map((name) => <button key={name} type="button" aria-checked={platform === name} onClick={() => setPlatform(name)}>{name === "Windows" ? "Windows 11" : name}</button>)}
        </div>
        <div className="segmented" role="radiogroup" aria-label="Appearance"><button type="button" aria-checked={!dark} onClick={() => setDark(false)}>{t.light}</button><button type="button" aria-checked={dark} onClick={() => setDark(true)}>{t.dark}</button></div>
      </div>
      <motion.div className={`demo-window demo-window-${platformClass} ${dark ? "dark" : ""}`} layout>
        <div className="demo-titlebar"><div className="traffic"><span style={{ background: "#ff5f57" }} /><span style={{ background: "#febc2e" }} /><span style={{ background: "#28c840" }} /></div><span>Settings</span>{platform !== "macOS" && <b className="demo-close" aria-hidden="true">×</b>}</div>
        <div className="demo-settings">
          <div className="setting-row"><span>Name</span><div className="setting-input">Breeze</div></div>
          <div className="setting-row"><span>Sync</span><span className="fake-toggle" /></div>
          <div className="setting-row"><span>Volume</span><div className="setting-slider"><i /></div></div>
          <label className="setting-check"><span className="check-box">✓</span>Remember me</label>
          <div className="setting-actions"><span className="fake-control">Cancel</span><span className="fake-control primary">Save</span></div>
        </div>
      </motion.div>
      <div className="demo-caption">{t.demoCaption}</div>
    </div>
  );
}

function Preview({ index }: { index: number }) {
  if (index % 4 === 0) return <div className="preview-buttons"><span>Cancel</span><span>Save</span></div>;
  if (index % 4 === 1) return <div className="preview-choice"><span className="preview-toggle" /><span className="preview-check">✓</span><span className="preview-radio" /></div>;
  if (index % 4 === 2) return <div className="preview-inputs"><div className="preview-input">Breeze<span style={{ marginLeft: "auto", width: 2, height: 14, background: "#0b7a63" }} /></div><div className="preview-input muted">⌕&nbsp; Search</div></div>;
  return <div className="preview-list"><span /><span /><span /></div>;
}

// Keep the homepage catalog aligned with the 94 component entries used by the
// reference page. The repository catalog also contains internal/demo entries.
const referenceCatalog = "Accordion:l,AccordionItem:l,ActionRow:l,AlertDialog:o,Autocomplete:t,Avatar:d,BackButton:n,Badge:f,Breadcrumbs:n,Button:a,Calendar:c*,CalendarView:c,Card:l,Checkbox:c,CheckboxGroup:c,Collapsible:l,ColorPicker:c*,ColorWell:c,ComboRow:c,Combobox:c,DateInput:t,DatePicker:t,DecimalInput:t*,Divider:l,EditableText:t,EntryRow:t,ExpanderRow:l,Field:l,Fieldset:l,FindBar:t*,Form:l,ForwardButton:n,GridView:d,HeaderBar:n,IconButton:a,InfoBar:f,Link:a,List:d,MenuButton:a,Meter:f,Modal:o,NumberInput:t,Outline:d,OutlineTable:d,PasswordField:t,Popover:o,PopoverMenuButton:a,PrimaryButton:a,Progress:f,ProgressRing:f,Radio:c,RadioGroup:c,RangeSlider:c*,Rating:c,Scroll:l*,ScrollBoth:l*,ScrollHorizontal:l*,SearchField:t,Segmented:c,Select:c,Sidebar:n,SidebarItem:n,SidebarSection:n,Slider:c,SpinRow:t,Spinner:f,Split:l,SplitButton:a,SplitPercent:l,SplitVertical:l,SplitVerticalPercent:l,StatusPage:f,StepSlider:c,Stepper:c,Switch:c,SwitchRow:c,TabView:n,Table:d,Tabs:n,Text:t,TextArea:t,TextAreaBuffer:t,TextInput:t,TextInputBuffer:t,TimeInput:t,Toasts:f,Toggle:a,ToggleGroup:a,TokenField:t,Toolbar:a,Tooltip:o,Tree:d,TreeItem:d,ViewSwitcher:n";
const referenceCategory: Record<string, string> = { a: "actions", c: "choices", t: "text", l: "layout", n: "navigation", d: "data", f: "feedback", o: "overlays" };
const referenceNames = referenceCatalog.split(",").map((item) => {
  const [name = "", rawKind = "l"] = item.split(":");
  return { name, category: referenceCategory[rawKind.charAt(0)] ?? "layout", look: rawKind.includes("*") ? "themed" : "drawn" } as const;
});

function Catalog({ t, entries, locale }: { t: Text; entries: CatalogEntry[]; locale: "zh" | "en" }) {
  const [query, setQuery] = useState("");
  const [category, setCategory] = useState("all");
  const categoryData: readonly (readonly [string, string, number])[] = locale === "zh" ? [
    ["all", "全部", 94], ["actions", "操作", 10], ["choices", "选择", 19], ["text", "文本输入", 18], ["layout", "布局", 17], ["navigation", "导航", 10], ["data", "数据", 8], ["feedback", "反馈", 8], ["overlays", "浮层", 4],
  ] as const : [
    ["all", "All", 94], ["actions", "Actions", 10], ["choices", "Choices", 19], ["text", "Text input", 18], ["layout", "Layout", 17], ["navigation", "Navigation", 10], ["data", "Data", 8], ["feedback", "Feedback", 8], ["overlays", "Overlays", 4],
  ];
  const overview = locale === "zh" ? [
    ["操作", 10, "Button, IconButton, Link +7"], ["选择", 19, "Calendar, CalendarView, Checkbox +16"], ["文本输入", 18, "Autocomplete, DateInput, DatePicker +15"], ["布局", 17, "Accordion, AccordionItem, ActionRow …"], ["导航", 10, "BackButton, Breadcrumbs, HeaderBar …"], ["数据", 8, "Avatar, Badge, DataTable …"], ["反馈", 8, "Alert, Progress, Toast …"], ["浮层", 4, "Dialog, Popover, Tooltip …"],
  ] : [
    ["Actions", 10, "Button, IconButton, Link +7"], ["Choices", 19, "Calendar, CalendarView, Checkbox +16"], ["Text input", 18, "Autocomplete, DateInput, DatePicker +15"], ["Layout", 17, "Accordion, AccordionItem, ActionRow …"], ["Navigation", 10, "BackButton, Breadcrumbs, HeaderBar …"], ["Data", 8, "Avatar, Badge, DataTable …"], ["Feedback", 8, "Alert, Progress, Toast …"], ["Overlays", 4, "Dialog, Popover, Tooltip …"],
  ];
  const filtered = useMemo(() => entries.filter((item) => {
    const q = query.trim().toLowerCase();
    return item.kind === "component" && (!q || `${item.name} ${item.description}`.toLowerCase().includes(q)) && (category === "all" || item.category === category);
  }), [entries, query, category]);
  const searching = query.trim().length > 0 || category !== "all";
  return <section className="section-band catalog-band" id="catalog"><div className="page-section" style={{ display: "flex", flexDirection: "column", gap: 32 }}>
    <div className="section-heading"><h2>{t.catalogTitle}</h2><p>{t.catalogText}</p></div>
    <div className="catalog-controls"><label className="search-box"><Search size={16} /><input type="search" value={query} onChange={(e) => setQuery(e.target.value)} placeholder={t.search} aria-label={t.search} /></label><div className="catalog-cats" role="radiogroup" aria-label={t.all}>{categoryData.map(([value, label, count]) => <button key={value} type="button" className={category === value ? "active" : ""} onClick={() => setCategory(value)}>{label} <span>{count}</span></button>)}</div></div>
    {!searching ? <div className="catalog-grid catalog-overview">{overview.map(([label, count, examples], index) => <motion.button type="button" className="catalog-card" key={label} initial={{ opacity: 0, y: 16 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true, amount: .2 }} transition={{ delay: Math.min(index * .03, .24) }}><span className="catalog-preview"><Preview index={index} /></span><span className="catalog-body"><span className="catalog-meta"><strong>{label}</strong><span>{count}</span></span><span className="catalog-description">{examples}</span></span></motion.button>)}</div> : <div className="catalog-grid">{filtered.map((item, index) => <motion.button type="button" className="catalog-card" key={`${item.kind}-${item.name}`} initial={{ opacity: 0, y: 16 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true, amount: .2 }} transition={{ delay: Math.min(index * .03, .24) }}><span className="catalog-preview"><Preview index={index} /></span><span className="catalog-body"><span className="catalog-meta"><span>{item.look}</span><span>{item.category}</span></span><strong>{item.name}</strong><span className="catalog-description">{item.description}</span></span></motion.button>)}</div>}
    {searching && filtered.length === 0 ? <p className="catalog-more">{t.empty}</p> : <p className="catalog-more">{searching ? t.shown.replace("{n}", String(filtered.length)) : ""} · <a href="https://breeze-ui.yldm.ai/storybook/" target="_blank" rel="noreferrer">{t.storybook} <ArrowRight size={14} style={{ verticalAlign: -2 }} /></a></p>}
  </div></section>;
}

function CodeSection({ t }: { t: Text }) {
  const [tab, setTab] = useState<"view" | "gallery">("view");
  return <section className="page-section" id="platforms"><div className="code-layout"><div className="code-copy"><div className="section-heading"><h2>{t.componentTitle}</h2><p>{t.componentText}</p></div><div className="code-rows"><div className="code-row"><code>breeze.Use(c)</code><span>{t.themeText}</span></div><div className="code-row"><code>SetPlatform</code><span>{t.setPlatform}</span></div><div className="code-row"><code>Preview(c, p, body)</code><span>{t.preview}</span></div></div></div><div className="code-card"><div className="code-tabs"><div role="tablist" aria-label="Go"><button type="button" className={tab === "view" ? "active" : ""} onClick={() => setTab("view")}>View</button><button type="button" className={tab === "gallery" ? "active" : ""} onClick={() => setTab("gallery")}>Gallery</button></div><span>Go</span></div>{tab === "view" ? <pre><span className="kw">import</span> ({`\n`}    <span className="str">"github.com/egoist/mygo/ui"</span>{`\n`}    breeze <span className="str">"github.com/yldm-tech/mygo-breeze-ui/ui/breeze"</span>{`\n`}){`\n\n`}<span className="kw">func</span> (app *settings) <span className="fn">view</span>(c *ui.Context) {'{'}{`\n`}    breeze.<span className="fn">Use</span>(c) <span className="comment">// 平台的颜色和尺寸</span>{`\n`}    breeze.<span className="fn">Checkbox</span>(c, &amp;app.remember, <span className="str">"Remember me"</span>){`\n`}    breeze.<span className="fn">Switch</span>(c, &amp;app.sync).<span className="fn">Label</span>(<span className="str">"Sync"</span>){`\n`}    <span className="kw">if</span> breeze.<span className="fn">PrimaryButton</span>(c, <span className="str">"Save"</span>).<span className="fn">Clicked</span>() {'{'}{`\n`}        app.<span className="fn">save</span>(){`\n`}    {'}'}{`\n`}{'}'}</pre> : <pre><span className="comment">// 三个平台并排，像组件画廊那样</span>{`\n`}<span className="kw">for</span> _, p := <span className="kw">range</span> breeze.Platforms {'{'}{`\n`}    breeze.<span className="fn">Preview</span>(c, p, <span className="kw">func</span>() {'{'}{`\n`}        breeze.<span className="fn">Switch</span>(c, &amp;app.sync).<span className="fn">Label</span>(<span className="str">"Sync"</span>){`\n`}    {'}'}){`\n`}{'}'}{`\n\n`}breeze.<span className="fn">SetPlatform</span>(breeze.Windows){`\n`}breeze.<span className="fn">SetPlatform</span>(breeze.Auto)</pre>}</div></div></section>;
}

function HomePage() {
  const { locale } = useUiStore();
  const t = texts[locale];
  const { data } = useQuery({ queryKey: ["catalog"], queryFn: getCatalog });
  const entries: CatalogEntry[] = referenceNames.map((item) => {
    const source = data?.components.find((candidate) => candidate.name === item.name);
    return {
      kind: "component",
      name: item.name,
      category: item.category,
      look: item.look,
      description: source?.description ?? `${item.name} component rendered per platform.`,
      example: source?.example ?? item.name,
    };
  });
  const [copied, setCopied] = useState(false);
  const [accent, setAccent] = useState("#0b7a63");
  const [galleryDark, setGalleryDark] = useState(false);
  const command = "go get github.com/yldm-tech/mygo-breeze-ui";
  const copyCommand = async () => { await navigator.clipboard?.writeText(command); setCopied(true); window.setTimeout(() => setCopied(false), 1400); };
  const stat = t.stat.replace("{drawn}", locale === "zh" ? "86" : "86").replace("{n}", "94");
  return <div className="site-shell">
    <Seo title={locale === "zh" ? "MyGo Breeze — MyGo 原生 UI 组件" : "MyGo Breeze — Native UI widgets for MyGo"} description={locale === "zh" ? "MyGo Breeze 按 macOS、Windows 11 和 GNOME 的原生样式绘制 MyGo 组件。" : "MyGo Breeze draws MyGo widgets to match macOS, Windows 11 and GNOME."} socialTitle="MyGo Breeze — Native UI for MyGo" socialDescription="Native UI widgets for MyGo, drawn to match macOS, Windows 11 and GNOME." path="/" locale={locale} />
    <main>
      <section className="hero" id="top"><div className="hero-grid"><span className="gust" /><span className="gust" /><span className="gust" /><span className="gust" /></div><motion.div className="hero-copy rise"><h1><span>{t.h1a}</span> <span>{t.h1b}</span></h1></motion.div><div className="hero-content"><motion.p className="hero-intro rise" transition={{ delay: .08 }}>{t.intro}</motion.p><motion.div className="hero-buttons rise" transition={{ delay: .16 }}><a className="btn-primary" href="#install">{t.quick}<ArrowRight size={16} /></a><a className="btn-secondary" href="https://breeze-ui.yldm.ai/storybook/" target="_blank" rel="noreferrer">{t.storybook}</a></motion.div><motion.div className="command-chip rise" transition={{ delay: .24 }}><span>$</span><code>{command}</code><button className="command-copy" type="button" onClick={copyCommand} aria-label={command}>{copied ? <Check size={16} color="#0b7a63" /> : <Clipboard size={16} />}</button></motion.div><motion.p className="hero-statline rise" transition={{ delay: .32 }}>{stat}</motion.p></div><motion.div className="hero-demo rise" transition={{ delay: .2}}><HeroDemo t={t} /></motion.div></section>
      <div className="ticker-wrap" aria-hidden="true"><div className="ticker-row">{(entries.length ? entries : [{ name: "Accordion" }, { name: "Button" }, { name: "Calendar" }, { name: "Checkbox" }, { name: "DatePicker" }, { name: "Dialog" }, { name: "HeaderBar" }, { name: "ProgressBar" }, { name: "Slider" }, { name: "Switch" }, { name: "AccordionItem" }, { name: "AlertDialog" }, { name: "Autocomplete" }, { name: "ComboBox" }, { name: "Menu" }, { name: "Popover" }, { name: "RadioButton" }, { name: "SearchField" }, { name: "TabView" }, { name: "Tooltip" }]).slice(0, 20).map((item, index) => <span className="ticker-pill" key={`a-${item.name}-${index}`}>{item.name}</span>)}</div><div className="ticker-row ticker-reverse">{(entries.length ? entries : [{ name: "Accordion" }, { name: "Button" }, { name: "Calendar" }, { name: "Checkbox" }, { name: "DatePicker" }, { name: "Dialog" }, { name: "HeaderBar" }, { name: "ProgressBar" }, { name: "Slider" }, { name: "Switch" }, { name: "AccordionItem" }, { name: "AlertDialog" }, { name: "Autocomplete" }, { name: "ComboBox" }, { name: "Menu" }, { name: "Popover" }, { name: "RadioButton" }, { name: "SearchField" }, { name: "TabView" }, { name: "Tooltip" }]).slice(20, 40).map((item, index) => <span className="ticker-pill" key={`b-${item.name}-${index}`}>{item.name}</span>)}</div></div>
      <section className="section-band" id="platforms-overview"><div className="page-section"><div className="section-heading" style={{ marginBottom: 48 }}><h2>{t.platformTitle}</h2><p>{t.platformText}</p></div><div className="platform-grid">{[
        { name: "macOS", version: "macOS 27", method: "与 AppKit 在 macOS 27.2 上的渲染逐像素对比。", button: "24 pt", track: "54 × 24", checkbox: "16 pt", accent: "#007aff" },
        { name: "Windows 11", version: "WinUI 3", method: "取自 WinUI 3 主题资源，并与 WinUI 3.3 在 Windows 11 26H2 上的渲染对比。", button: "32 px", track: "40 × 20", checkbox: "20 px", accent: "#0078d4" },
        { name: "GNOME", version: "libadwaita 1.10", method: "与 libadwaita 1.10.0（GTK 4.24.1）的渲染逐像素对比。", button: "34 px", track: "46 × 26", checkbox: "20 px", accent: "#3584e4" },
      ].map((p, i) => <motion.article className="platform-card" key={p.name} initial={{ opacity: 0, y: 18 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true, amount: .2 }} transition={{ delay: i * .08 }}><div className="platform-head"><h3>{p.name}</h3><span className="pill">{p.version}</span></div><p>{p.method}</p><dl><dt>{locale === "zh" ? "按钮高度" : "Button height"}</dt><dd>{p.button}</dd><dt>{locale === "zh" ? "Switch 轨道" : "Switch track"}</dt><dd>{p.track}</dd><dt>Checkbox</dt><dd>{p.checkbox}</dd><dt>{locale === "zh" ? "默认强调色" : "System accent"}</dt><dd><span className="swatch" style={{ background: p.accent }} /> {p.accent}</dd></dl></motion.article>)}</div></div></section>
      <section className="page-section gallery-section"><div className="gallery-header"><div className="section-heading"><h2>{t.galleryTitle}</h2><p>{t.galleryText}</p></div><div className="gallery-toggle"><button className="active" type="button">Side by side</button><button type="button" onClick={() => setGalleryDark((value) => !value)}>This system</button></div></div><figure className="gallery-figure"><div className="gallery-images">{[
        { name: "macOS", className: "mac", radius: 12, barHeight: 32 },
        { name: "Windows 11", className: "win", radius: 8, barHeight: 32 },
        { name: "GNOME", className: "gnome", radius: 14, barHeight: 38 },
      ].map((card, index) => <div className={`gallery-card gallery-${card.className}${galleryDark ? " dark" : ""}`} key={card.name} style={{ borderRadius: card.radius }}><div className="gallery-bar" style={{ height: card.barHeight }}><span>{card.name}</span><b aria-hidden="true">×</b></div><div className="gallery-crop"><img src={`${import.meta.env.BASE_URL}${galleryDark ? "gallery-dark.png" : "gallery-light.png"}`} style={{ marginLeft: `${-100 * index}%` }} alt={`Breeze UI gallery: ${card.name}`} /></div></div>)}</div><figcaption className="gallery-caption">{t.galleryCaption}</figcaption></figure></section>
      <CodeSection t={t} />
      <section className="section-band"><div className="page-section theme-layout"><div className="theme-copy"><div className="section-heading"><h2>{t.themeTitle}</h2><p>{t.themeText}</p></div><div className="theme-cells"><div><h3>{t.accent}</h3><p>{t.accentText}</p></div><div><h3>{t.contrast}</h3><p>{t.contrastText}</p></div><div><h3>{t.size}</h3><p>{t.sizeText}</p></div></div></div><div className="theme-demo" style={{ ["--accent" as string]: accent }}><div className="theme-top"><span>System accent</span><div className="accent-buttons">{["#0b7a63", "#0078d4", "#8a3fd1", "#d14d72", "#b4540a"].map((value) => <button key={value} type="button" className={accent === value ? "active" : ""} style={{ ["--accent" as string]: value }} onClick={() => setAccent(value)} aria-label={`Set accent ${value}`} />)}</div></div><div className="theme-cells"><div className="theme-cell"><span>macOS</span><span className="theme-switch" /><span className="theme-button">Save</span></div><div className="theme-cell"><span>Windows</span><span className="theme-switch" /><span className="theme-button">Save</span></div><div className="theme-cell"><span>GNOME</span><span className="theme-switch" /><span className="theme-button">Save</span></div></div></div></div></section>
      <section className="page-section" id="tooling"><div className="section-heading" style={{ marginBottom: 48 }}><h2>{t.toolingTitle}</h2><p>{t.toolingText}</p></div><div className="tool-grid"><ToolCard icon={<Sparkles size={20} />} title="Storybook" text={t.toolStorybook} code="go run ./examples/storybook" href="https://breeze-ui.yldm.ai/storybook/" link={t.storybook} /><ToolCard icon={<Terminal size={20} />} title="breeze-ui CLI" text={t.toolCli} code="npx --yes --package breeze-ui-cli breeze-ui add Switch" href="https://github.com/yldm-tech/mygo-breeze-ui" link={t.docs} /><ToolCard icon={<Code2 size={20} />} title="MCP server" text={t.toolMcp} code="breeze-ui-mcp" href="https://github.com/yldm-tech/mygo-breeze-ui" link="MCP 文档 →" /></div></section>
      <Catalog t={t} entries={entries} locale={locale} />
      <section className="page-section cta-wrap" id="install"><div className="cta-card"><div className="cta-copy"><h2>{t.ctaTitle}</h2><p>{t.ctaText}</p><code>go run ./examples/breeze-ui</code></div><div className="hero-buttons"><a className="btn-dark" href="https://github.com/yldm-tech/mygo-breeze-ui" target="_blank" rel="noreferrer">{t.quick}</a><a className="btn-outline-dark" href="https://github.com/yldm-tech/mygo-breeze-ui" target="_blank" rel="noreferrer">{t.github}<ExternalLink size={16} /></a></div></div></section>
    </main>
  </div>;
}

function ToolCard({ icon, title, text, code, href, link }: { icon: ReactNode; title: string; text: string; code: string; href: string; link: string }) {
  return <motion.article className="tool-card" initial={{ opacity: 0, y: 18 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true, amount: .2 }}><div style={{ color: "#0b7a63" }}>{icon}</div><h3>{title}</h3><p>{text}</p><code>{code}</code><a className="tool-link" href={href} target="_blank" rel="noreferrer">{link} <ArrowRight size={15} /></a></motion.article>;
}
