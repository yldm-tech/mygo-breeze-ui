import { useMemo, useState } from "react";
import { ArrowLeft, Check, ChevronDown, ChevronRight, Copy, Search, Sparkles } from "lucide-react";
import { createRoute, Link } from "@tanstack/react-router";
import { RootRoute } from "./root";
import { Seo } from "../components/Seo";
import { GoCodeBlock } from "../components/CodeHighlight";
import { catalog, type CatalogComponent } from "../lib/catalog";
import { useUiStore } from "../store";
import type { Locale } from "../store";
import translations from "../locales/components.json";

const translationData = translations as unknown as Record<string, { descriptions: Record<string, string>; demo: Record<string, any> }> & { categories: Record<string, Record<string, string>> };
export const categoryLabels: Record<string, string> = translations.categories.zh;
export function localizedCategory(category: string, locale: Locale) {
  const labels = translationData.categories[locale] ?? translationData.categories.en ?? {};
  return labels[category] ?? category;
}
function localeTranslations(locale: Locale) {
  return (translationData[locale] as { descriptions: Record<string, string>; demo: Record<string, any> } | undefined) ?? translationData.en as { descriptions: Record<string, string>; demo: Record<string, any> };
}

const platformMeta = [
  { key: "mac", label: "macOS", accent: "#007AFF", font: "-apple-system, BlinkMacSystemFont, sans-serif" },
  { key: "win", label: "Windows 11", accent: "#0078D4", font: "'Segoe UI Variable', 'Segoe UI', sans-serif" },
  { key: "gnome", label: "GNOME", accent: "#3584E4", font: "'Adwaita Sans', Cantarell, sans-serif" },
] as const;

// Keep the docs gallery in lockstep with the 94-component reference catalog.
const referenceComponentNames = "Accordion,AccordionItem,ActionRow,AlertDialog,Autocomplete,Avatar,BackButton,Badge,Breadcrumbs,Button,Calendar,CalendarView,Card,Checkbox,CheckboxGroup,Collapsible,ColorPicker,ColorWell,ComboRow,Combobox,DateInput,DatePicker,DecimalInput,Divider,EditableText,EntryRow,ExpanderRow,Field,Fieldset,FindBar,Form,ForwardButton,GridView,HeaderBar,IconButton,InfoBar,Link,List,MenuButton,Meter,Modal,NumberInput,Outline,OutlineTable,PasswordField,Popover,PopoverMenuButton,PrimaryButton,Progress,ProgressRing,Radio,RadioGroup,RangeSlider,Rating,Scroll,ScrollBoth,ScrollHorizontal,SearchField,Segmented,Select,Sidebar,SidebarItem,SidebarSection,Slider,SpinRow,Spinner,Split,SplitButton,SplitPercent,SplitVertical,SplitVerticalPercent,StatusPage,StepSlider,Stepper,Switch,SwitchRow,TabView,Table,Tabs,Text,TextArea,TextAreaBuffer,TextInput,TextInputBuffer,TimeInput,Toasts,Toggle,ToggleGroup,TokenField,Toolbar,Tooltip,Tree,TreeItem,ViewSwitcher".split(",");
export const docsComponents = catalog.components.filter((item) => referenceComponentNames.includes(item.name));

export function localizedDescription(component: CatalogComponent, locale: Locale) {
  const descriptions = localeTranslations(locale).descriptions;
  return descriptions[component.name] ?? component.description;
}

function findComponent(name: string) {
  return docsComponents.find((item) => item.name.toLowerCase() === decodeURIComponent(name).toLowerCase());
}

function CopyExample({ value }: { value: string }) {
  const [copied, setCopied] = useState(false);
  return (
    <button
      type="button"
      className="inline-flex min-h-9 items-center gap-1.5 rounded-lg bg-white px-3 text-xs font-semibold text-[#555C69] shadow-[inset_0_0_0_1px_#E2E4E9] hover:text-[#0B7A63]"
      onClick={() => {
        void navigator.clipboard?.writeText(value);
        setCopied(true);
        window.setTimeout(() => setCopied(false), 1300);
      }}
    >
      {copied ? <Check size={14} /> : <Copy size={14} />}
      {copied ? "已复制" : "复制代码"}
    </button>
  );
}

function ComponentDemo({ component, platform, locale }: { component: CatalogComponent; platform: (typeof platformMeta)[number]; locale: Locale }) {
  const copy = localeTranslations(locale).demo as any;
  const [on, setOn] = useState(true);
  const [checked, setChecked] = useState(true);
  const [selected, setSelected] = useState("all");
  const [open, setOpen] = useState(false);
  const [value, setValue] = useState(62);
  const name = component.name;
  const isWin = platform.key === "win";
  const isGnome = platform.key === "gnome";
  const radius = isWin ? "rounded-md" : isGnome ? "rounded-xl" : "rounded-lg";
  const button = `${radius} min-h-9 px-4 text-sm font-semibold transition active:scale-[.98]`;
  const surface = `${radius} border border-[#E2E4E9] bg-white px-4 py-3 shadow-[0_12px_28px_rgba(20,24,32,0.08)]`;

  if (["Switch", "SwitchRow", "Toggle"].includes(name)) {
    return <button type="button" className={`${surface} flex min-w-[220px] items-center justify-between gap-6 text-left`} onClick={() => setOn((v) => !v)}><span><strong className="block text-sm">{name}</strong><span className="text-xs text-[#6B7280]">{on ? copy.enabled : copy.disabled}</span></span><span className="relative block h-6 w-12 rounded-full" style={{ background: on ? platform.accent : "#D5D9DF" }}><span className="absolute top-1 size-4 rounded-full bg-white shadow-sm transition-all" style={{ left: on ? "26px" : "4px" }} /></span></button>;
  }
  if (["Checkbox", "CheckboxGroup", "Radio", "RadioGroup", "RadioRow"].includes(name)) {
    return <div className={`${surface} flex min-w-[240px] flex-col gap-3 text-sm`}><label className="flex cursor-pointer items-center gap-2"><input className="size-4 accent-[#0B7A63]" type={name.includes("Radio") ? "radio" : "checkbox"} checked={checked} onChange={(event) => setChecked(event.target.checked)} /> {copy.rememberMe}</label><label className="flex cursor-pointer items-center gap-2"><input className="size-4 accent-[#0B7A63]" type={name.includes("Radio") ? "radio" : "checkbox"} checked={!checked} onChange={(event) => setChecked(!event.target.checked)} /> {copy.sendNotifications}</label></div>;
  }
  if (["Slider", "RangeSlider", "StepSlider", "Progress", "Meter", "ProgressRing"].includes(name)) {
    return <div className={`${surface} flex min-w-[260px] flex-col gap-3`}><div className="flex items-center justify-between text-sm"><strong>{name}</strong><span className="font-mono text-xs text-[#555C69]">{value}%</span></div><input className="w-full accent-[#0B7A63]" type="range" min="0" max="100" value={value} onChange={(event) => setValue(Number(event.target.value))} /><div className="h-2 overflow-hidden rounded-full bg-[#E5E7EB]"><div className="h-full rounded-full" style={{ width: `${value}%`, background: platform.accent }} /></div></div>;
  }
  if (["Modal", "AlertDialog", "Sheet", "Popover", "PopoverMenuButton", "Tooltip"].includes(name)) {
    return <div className="relative flex min-h-[140px] min-w-[260px] flex-col items-center justify-center gap-3"><button type="button" className={button} style={{ background: platform.accent, color: "#fff" }} onClick={() => setOpen((v) => !v)}>{open ? copy.close : copy.open.replace("{name}", name)}</button>{open && <div className={`${surface} absolute bottom-0 left-1/2 z-10 w-[min(280px,90%)] -translate-x-1/2 translate-y-1/2`}><strong className="block text-sm">{name}</strong><p className="m-0 mt-1 text-xs text-[#555C69]">{copy.interactivePreview}</p><button type="button" className="mt-3 text-xs font-semibold" style={{ color: platform.accent }} onClick={() => setOpen(false)}>{copy.done}</button></div>}</div>;
  }
  if (["Tabs", "Segmented", "ToggleGroup", "ViewSwitcher"].includes(name)) {
    const tabs = ["all", "recent", "pinned"] as const;
    return <div className={`${surface} flex min-w-[280px] flex-wrap gap-1 p-1`}>{tabs.map((item) => <button key={item} type="button" className={`${radius} px-3 py-2 text-xs font-semibold`} style={{ background: selected === item ? platform.accent : "transparent", color: selected === item ? "#fff" : "#555C69" }} onClick={() => setSelected(item)}>{copy.tabs[item]}</button>)}</div>;
  }
  if (["TextInput", "TextInputBuffer", "TextArea", "TextAreaBuffer", "SearchField", "DateInput", "DatePicker", "TimeInput", "NumberInput", "DecimalInput", "PasswordField", "EntryRow", "Autocomplete", "Combobox", "Select", "TokenField"].includes(name)) {
    return <div className={`${surface} flex min-w-[280px] flex-col gap-2`}><label className="text-xs font-semibold text-[#555C69]">{name} <span className="font-normal">{copy.value}</span></label>{name.includes("Area") ? <textarea className="min-h-16 resize-none rounded-md border border-[#D5D9DF] p-2 text-sm outline-none focus:border-[#0B7A63]" defaultValue="Breeze UI" /> : <input className="h-9 rounded-md border border-[#D5D9DF] px-3 text-sm outline-none focus:border-[#0B7A63]" placeholder={copy.typeInto.replace("{name}", name)} />}</div>;
  }
  if (["Calendar", "CalendarView"].includes(name)) {
    return <div className={`${surface} min-w-[280px]`}><div className="mb-3 flex items-center justify-between text-sm font-semibold"><span>{copy.month}</span><span className="text-xs" style={{ color: platform.accent }}>{copy.today}</span></div><div className="grid grid-cols-7 gap-1 text-center text-xs text-[#555C69]">{[...copy.weekdays, ...Array.from({ length: 28 }, (_, i) => String(i + 1))].map((day, index) => <button key={`${day}-${index}`} type="button" className="grid aspect-square place-items-center rounded-md hover:bg-[#EEF5FF]" style={index === 13 ? { background: platform.accent, color: "#fff" } : undefined}>{day}</button>)}</div></div>;
  }
  if (["Button", "PrimaryButton", "DestructiveButton", "Link", "IconButton", "BackButton", "ForwardButton", "MenuButton"].includes(name)) {
    return <div className="flex flex-wrap items-center justify-center gap-2"><button type="button" className={button} style={{ background: name.includes("Destructive") ? "#C9363E" : platform.accent, color: "#fff" }}>{name}</button><button type="button" className={`${button} border border-[#D5D9DF] bg-white text-[#3A404C]`}>{copy.cancel}</button></div>;
  }
  return <div className={`${surface} flex min-w-[260px] flex-col gap-3`}><div className="flex items-center justify-between"><strong className="text-sm">{name}</strong><Sparkles size={16} style={{ color: platform.accent }} /></div><div className="h-10 rounded-md bg-[#F4F5F7]" /><div className="flex gap-2"><span className="h-2 w-20 rounded-full bg-[#D5D9DF]" /><span className="h-2 w-12 rounded-full bg-[#E6E9ED]" /></div></div>;
}

function ComponentCard({ item, locale }: { item: CatalogComponent; locale: Locale }) {
  return <Link to="/docs/components/$component" params={{ component: item.name }} className="group flex min-h-[198px] flex-col rounded-2xl bg-[#F8F9FB] p-5 shadow-[inset_0_0_0_1px_#E2E4E9] transition hover:-translate-y-0.5 hover:bg-white hover:shadow-[0_14px_30px_rgba(20,24,32,0.1)]"><div className="mb-4 flex items-center justify-between gap-3"><span className="rounded-full bg-white px-2.5 py-1 font-mono text-[11px] text-[#555C69] shadow-[inset_0_0_0_1px_#E2E4E9]">{localizedCategory(item.category, locale)}</span><span className={`size-2 rounded-full ${item.look === "drawn" ? "bg-[#0B7A63]" : "bg-[#B4540A]"}`} title={item.look === "drawn" ? "按平台自绘" : "套用平台主题"} /></div><strong className="text-lg tracking-[-0.02em] group-hover:text-[#0B7A63]">{item.name}</strong><p className="mt-2 line-clamp-3 text-sm leading-[1.55] text-[#555C69]">{localizedDescription(item, locale)}</p><span className="mt-auto pt-4 text-xs font-semibold text-[#0B7A63]">查看用法与演示 →</span></Link>;
}

function ComponentDocsSidebar({ active, locale }: { active?: string; locale: Locale }) {
  const categoryOrder = ["actions", "choices", "text", "layout", "navigation", "data", "feedback", "overlays"];
  const [query, setQuery] = useState("");
  const activeCategory = docsComponents.find((item) => item.name === active)?.category;
  const [expanded, setExpanded] = useState<Record<string, boolean>>(() => (activeCategory ? { [activeCategory]: true } : {}));
  const visible = query.trim().toLowerCase();
  const link = (label: string, href: string) => {
    const selected = label === active || (!active && href === "/docs/components");
    const className = `flex min-h-[32px] items-center rounded-lg px-2.5 ${selected ? "bg-[#E3F3EE] font-semibold text-[#075A49]" : "text-[#3A404C] hover:bg-[#F8F9FB] hover:text-[#0B7A63]"}`;
    if (href === "/docs/components") return <Link key={`${label}-${href}`} to="/docs/components" className={className} aria-current={selected ? "page" : undefined}>{label}</Link>;
    if (href.startsWith("/docs/components/")) {
      const component = href.slice("/docs/components/".length);
      return <Link key={`${label}-${href}`} to="/docs/components/$component" params={{ component }} className={className} aria-current={selected ? "page" : undefined}>{label}</Link>;
    }
    return <a key={`${label}-${href}`} className={className} href={href} aria-current={selected ? "page" : undefined}>{label}</a>;
  };
  const sectionLink = "flex min-h-[34px] items-center rounded-lg px-2.5 text-[#3A404C] hover:bg-[#F8F9FB] hover:text-[#0B7A63]";
  return <aside className="block w-full shrink-0 border-b border-[#E2E4E9] bg-white px-4 py-5 lg:w-[280px] lg:border-b-0 lg:border-r lg:px-4 lg:py-7"><nav className="flex max-h-[calc(100vh-112px)] flex-col gap-5 overflow-y-auto pr-1 text-[13px] lg:sticky lg:top-24" aria-label="文档目录">
    <div className="flex flex-col gap-1">
      <p className="m-0 px-2 text-[12px] font-semibold uppercase tracking-[0.08em] text-[#8A90A0]">开始</p>
      <Link className={sectionLink} to="/docs" hash="top">快速开始</Link>
      <Link className={sectionLink} to="/docs" hash="install">安装</Link>
      <Link className={sectionLink} to="/docs" hash="first-view">第一个视图</Link>
      <Link className={sectionLink} to="/docs" hash="storybook">Storybook</Link>
    </div>
    <div className="flex flex-col gap-1">
      <p className="m-0 px-2 text-[12px] font-semibold uppercase tracking-[0.08em] text-[#8A90A0]">指南</p>
      <Link className={sectionLink} to="/docs" hash="platforms">平台与 Preview</Link>
      <Link className={sectionLink} to="/docs" hash="themes">主题</Link>
      <Link className={sectionLink} to="/docs" hash="limits">已知限制</Link>
    </div>
    <div>
      <p className="mb-2 px-2 text-[12px] font-semibold uppercase tracking-[0.08em] text-[#8A90A0]">组件 · {docsComponents.length}</p>
      <label className="mb-2 flex h-11 items-center gap-2 rounded-xl bg-white px-3.5 text-[#6B7280] shadow-[inset_0_0_0_1px_#E2E4E9,0_2px_8px_rgba(20,24,32,0.05)]"><Search size={16} /><span className="sr-only">搜索组件</span><input className="min-w-0 flex-1 border-0 bg-transparent text-sm outline-none" value={query} onChange={(event) => setQuery(event.target.value)} placeholder="搜索组件" /></label>
    </div>
    <div className="mt-1 flex flex-col gap-1">{categoryOrder.map((category) => {
      const items = docsComponents.filter((item) => item.category === category && (!visible || `${item.name} ${localizedDescription(item, locale)}`.toLowerCase().includes(visible)));
      if (!items.length) return null;
      const isExpanded = Boolean(expanded[category] || visible);
      return <div key={category}>
        <button type="button" className="flex min-h-[36px] w-full items-center gap-2 rounded-lg px-2.5 text-left text-[#555C69] hover:bg-[#F8F9FB]" aria-expanded={isExpanded} onClick={() => setExpanded((value) => ({ ...value, [category]: !value[category] }))}>
          {isExpanded ? <ChevronDown size={15} /> : <ChevronRight size={15} />}<span className="flex-1">{localizedCategory(category, locale)}</span><span className="text-xs text-[#8A90A0]">{docsComponents.filter((item) => item.category === category).length}</span>
        </button>
        {isExpanded && <div className="ml-5 flex flex-col gap-0.5 border-l border-[#E2E4E9] pl-2">{items.map((item) => link(item.name, `/docs/components/${item.name}`))}</div>}
      </div>;
    })}</div>
  </nav></aside>;
}

function ComponentsIndexPage() {
  const { locale } = useUiStore();
  const [query, setQuery] = useState("");
  const [category, setCategory] = useState("all");
  const categories = ["all", ...Array.from(new Set(docsComponents.map((item) => item.category)))];
  const filtered = useMemo(() => docsComponents.filter((item) => (!query || `${item.name} ${localizedDescription(item, locale)}`.toLowerCase().includes(query.toLowerCase())) && (category === "all" || item.category === category)), [category, locale, query]);
  return <><Seo title="Breeze UI 文档 · 组件" description="Breeze UI 94 个原生组件的用法、平台差异和交互演示。" socialTitle="Breeze UI components" socialDescription="Usage, platform differences and interactive previews for all Breeze UI components." path="/docs/components" locale={locale} /><main className="min-h-screen bg-white"><div className="mx-auto grid w-[min(1440px,100%)] grid-cols-[280px_minmax(0,1fr)] max-[900px]:grid-cols-1"><ComponentDocsSidebar locale={locale} /><div className="min-w-0 px-12 py-12 max-[680px]:px-[17px] max-[680px]:py-9"><div className="mb-10 max-w-[780px]"><Link className="mb-7 inline-flex items-center gap-1.5 text-sm text-[#555C69] hover:text-[#0B7A63]" to="/docs"><ArrowLeft size={15} /> 文档首页</Link><p className="mb-3 font-mono text-xs uppercase tracking-[0.09em] text-[#0B7A63]">组件 · {docsComponents.length}</p><h1 className="m-0 text-[42px] font-semibold leading-[1.1] tracking-[-0.04em] max-[680px]:text-[34px]">每个组件，都有用法和演示</h1><p className="mt-4 text-lg leading-[1.65] text-[#555C69]">从按钮、输入框到日期选择器和设置行。每个页面都列出 Go 调用、三平台外观、交互状态和实现差异。</p></div><div className="mb-7 flex flex-wrap items-center gap-3"><label className="flex h-11 min-w-[260px] flex-1 items-center gap-2 rounded-xl bg-[#F4F5F7] px-3.5 text-[#6B7280] shadow-[inset_0_0_0_1px_#E2E4E9]"><Search size={16} /><span className="sr-only">搜索组件</span><input className="min-w-0 flex-1 border-0 bg-transparent text-sm outline-none" value={query} onChange={(event) => setQuery(event.target.value)} placeholder="搜索组件，例如 Calendar" /></label><span className="text-sm text-[#6B7280]">显示 {filtered.length} / {docsComponents.length}</span></div><div className="mb-8 flex gap-2 overflow-x-auto pb-1">{categories.map((item) => <button key={item} type="button" className={`whitespace-nowrap rounded-full px-3 py-2 text-sm font-semibold ${category === item ? "bg-[#15171C] text-white" : "bg-[#F4F5F7] text-[#555C69]"}`} onClick={() => setCategory(item)}>{item === "all" ? "总览" : localizedCategory(item, locale)}</button>)}</div><div className="grid grid-cols-4 gap-3 max-[1100px]:grid-cols-3 max-[760px]:grid-cols-2 max-[520px]:grid-cols-1">{filtered.map((item) => <ComponentCard item={item} locale={locale} key={item.name} />)}</div></div></div></main></>;
}

function ComponentDocPage() {
  const { locale } = useUiStore();
  const { component: componentParam } = ComponentDocRoute.useParams();
  const component = findComponent(componentParam);
  const [platformKey, setPlatformKey] = useState<(typeof platformMeta)[number]["key"]>("mac");
  const platform = platformMeta.find((item) => item.key === platformKey) ?? platformMeta[0];
  if (!component) return <main className="mx-auto min-h-screen w-[min(900px,calc(100%-48px))] py-20"><h1 className="text-3xl font-semibold">找不到这个组件</h1><Link className="mt-5 inline-flex text-[#0B7A63]" to="/docs/components">返回组件目录</Link></main>;
  const description = localizedDescription(component, locale);
  const related = docsComponents.filter((item) => item.category === component.category && item.name !== component.name).slice(0, 4);
  return <><Seo title={`Breeze UI · ${component.name}`} description={`${component.name} 的 Go 用法、三平台差异和交互演示。`} socialTitle={`Breeze UI ${component.name}`} socialDescription={`${component.name} usage and platform preview`} path={`/docs/components/${component.name}`} locale={locale} /><main className="min-h-screen bg-white"><div className="mx-auto grid w-[min(1440px,100%)] grid-cols-[280px_minmax(0,780px)_220px] gap-8 py-12 max-[1100px]:grid-cols-1 max-[680px]:w-[calc(100%-34px)] max-[680px]:py-9"><ComponentDocsSidebar active={component.name} locale={locale} /><article><div className="mb-8 flex items-center gap-2 text-sm text-[#6B7280]"><Link className="hover:text-[#0B7A63]" to="/docs">文档</Link><span>/</span><span>组件</span><span>/</span><span className="text-[#15171C]">{component.name}</span></div><div className="mb-8 flex flex-wrap items-start justify-between gap-5"><div><h1 className="m-0 text-[44px] font-semibold leading-[1.1] tracking-[-0.04em] max-[680px]:text-[35px]">{component.name}</h1><p className="mt-3 max-w-[700px] text-lg leading-[1.65] text-[#555C69]">{description}</p></div><span className="rounded-full bg-[#E3F3EE] px-3 py-1.5 text-xs font-semibold text-[#075A49]">{component.look === "drawn" ? "按平台自绘" : "套用平台主题"}</span></div><section className="mb-12 rounded-2xl bg-[#F8F9FB] p-5 shadow-[inset_0_0_0_1px_#E2E4E9]" aria-labelledby="component-demo-heading"><div className="mb-4 flex flex-wrap items-center justify-between gap-3"><div><h2 id="component-demo-heading" className="m-0 text-xl font-semibold">交互演示</h2><p className="m-0 mt-1 text-sm text-[#6B7280]">点击控件查看 {component.name} 的状态变化。</p></div><div className="flex rounded-xl bg-white p-1 shadow-[inset_0_0_0_1px_#E2E4E9]">{platformMeta.map((item) => <button key={item.key} type="button" className={`rounded-lg px-3 py-2 text-xs font-semibold ${platformKey === item.key ? "text-white" : "text-[#555C69]"}`} style={platformKey === item.key ? { background: item.accent } : undefined} onClick={() => setPlatformKey(item.key)}>{item.label}</button>)}</div></div><div className="flex min-h-[240px] items-center justify-center rounded-xl border border-[#E2E4E9] bg-white p-6" style={{ fontFamily: platform.font }}><ComponentDemo component={component} platform={platform} locale={locale} /></div><p className="m-0 mt-3 text-xs text-[#6B7280]">当前预览：{platform.label} · {component.look === "drawn" ? "控件由 Breeze UI 按平台绘制" : "MyGo 组件套用平台主题"}</p></section><section className="mb-12" id="usage"><div className="mb-4 flex items-center justify-between gap-3"><div><h2 className="m-0 text-[26px] font-semibold">用法</h2><p className="m-0 mt-1 text-sm text-[#555C69]">API 名称和签名与 MyGo 的 ui 包保持一致。</p></div><CopyExample value={component.example} /></div><GoCodeBlock code={component.example} className="m-0 overflow-x-auto rounded-2xl bg-[#15171C] p-5 font-mono text-[13px] leading-[1.75] text-[#E8EEF2]" /></section><section className="mb-12" id="platforms"><h2 className="m-0 text-[26px] font-semibold">三平台对照</h2><p className="mt-2 text-[#3A404C]">同一段视图代码会根据窗口平台呈现不同的尺寸、间距、颜色和交互反馈。</p><div className="mt-5 overflow-x-auto rounded-xl shadow-[0_0_0_1px_#E2E4E9]"><table className="w-full min-w-[660px] border-collapse text-sm"><thead><tr className="bg-[#FAFBFC] text-left"><th className="border-b border-[#E2E4E9] px-4 py-3 font-semibold text-[#555C69]">平台</th><th className="border-b border-[#E2E4E9] px-4 py-3 font-semibold text-[#555C69]">外观</th><th className="border-b border-[#E2E4E9] px-4 py-3 font-semibold text-[#555C69]">怎么检查</th></tr></thead><tbody>{platformMeta.map((item) => <tr key={item.key}><th className="border-b border-[#EEF0F3] px-4 py-3.5 text-left font-semibold" style={{ color: item.accent }}>{item.label}</th><td className="border-b border-[#EEF0F3] px-4 py-3.5 text-[#3A404C]">{item.key === "mac" ? "紧凑的原生控件、较小的按钮和圆润的系统间距。" : item.key === "win" ? "WinUI 3 的 32px 控件、强调色边框和键盘焦点。" : "libadwaita 的 34px 控件、卡片列表和更宽的触控目标。"}</td><td className="border-b border-[#EEF0F3] px-4 py-3.5 text-[#3A404C]"><code className="font-mono text-xs">breeze.SetPlatform(breeze.{item.key === "mac" ? "Mac" : item.key === "win" ? "Windows" : "Linux"})</code></td></tr>)}</tbody></table></div></section><section className="mb-12" id="details"><h2 className="m-0 text-[26px] font-semibold">实现说明</h2><p className="mt-2 text-[#3A404C]">{description}</p><ul className="mt-4 flex list-disc flex-col gap-2 pl-5 text-[#3A404C]"><li>组件返回 <code className="font-mono text-sm">ui.Element</code>，可以继续链式调用布局和样式方法。</li><li>交互事件沿用 MyGo 基座：按钮响应 <code className="font-mono text-sm">Clicked()</code>，输入和选择控件绑定指针值。</li><li>目录标记为 <strong>{component.look === "drawn" ? "drawn" : "themed"}</strong>，可在截图和 Storybook 中核对实际平台外观。</li></ul></section>{related.length > 0 && <section className="border-t border-[#E2E4E9] pt-8"><h2 className="m-0 text-[22px] font-semibold">同类组件</h2><div className="mt-4 grid grid-cols-2 gap-3 max-[560px]:grid-cols-1">{related.map((item) => <Link key={item.name} to="/docs/components/$component" params={{ component: item.name }} className="rounded-xl p-4 shadow-[inset_0_0_0_1px_#E2E4E9] hover:text-[#0B7A63]"><strong>{item.name}</strong><span className="mt-1 block text-xs text-[#6B7280]">{localizedCategory(item.category, locale)}</span></Link>)}</div></section>}</article><aside className="hidden lg:block"><div className="sticky top-24 flex flex-col gap-1 text-[13px]"><span className="pb-2 font-semibold">本页内容</span><a className="border-l-2 border-[#0B7A63] py-1 pl-2.5 text-[#0B7A63]" href="#usage">用法</a><a className="border-l border-[#E2E4E9] py-1 pl-2.5 text-[#555C69]" href="#platforms">三平台对照</a><a className="border-l border-[#E2E4E9] py-1 pl-2.5 text-[#555C69]" href="#details">实现说明</a></div></aside></div></main></>;
}

export const ComponentsRoute = createRoute({ getParentRoute: () => RootRoute, path: "/docs/components", component: ComponentsIndexPage });
export const ComponentDocRoute = createRoute({ getParentRoute: () => RootRoute, path: "/docs/components/$component", component: ComponentDocPage });
