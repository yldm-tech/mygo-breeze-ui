import { useState } from "react";
import { ChevronDown, Globe2, Menu, X } from "lucide-react";
import { Link } from "@tanstack/react-router";
import { useUiStore } from "../store";
import { useCopy } from "../copy";

function BrandMark() {
  return (
    <svg className="block size-7" viewBox="0 0 64 64" aria-hidden="true">
      <rect x="4" y="6" width="36" height="28" rx="9" fill="#A6E1CD" />
      <rect x="14" y="17" width="36" height="28" rx="4" fill="#3FAE8C" />
      <rect x="24" y="28" width="36" height="28" rx="7" fill="#0B7A63" />
      <circle cx="31" cy="35" r="1.8" fill="#fff" />
      <circle cx="37" cy="35" r="1.8" fill="#fff" fillOpacity=".7" />
      <circle cx="43" cy="35" r="1.8" fill="#fff" fillOpacity=".45" />
    </svg>
  );
}

function GitHubMark() {
  return <svg viewBox="0 0 24 24" className="size-5" aria-hidden="true"><path fill="currentColor" d="M12 .297a12 12 0 0 0-3.79 23.387c.6.111.82-.261.82-.577v-2.234c-3.338.726-4.042-1.416-4.042-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.84 1.237 1.84 1.237 1.07 1.834 2.807 1.304 3.492.997.108-.775.418-1.305.762-1.605-2.665-.303-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.124-.303-.535-1.523.117-3.176 0 0 1.008-.322 3.3 1.23a11.5 11.5 0 0 1 3.003-.404c1.02.005 2.047.138 3.003.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.873.118 3.176.77.84 1.233 1.91 1.233 3.22 0 4.61-2.807 5.624-5.48 5.921.43.372.823 1.103.823 2.222v3.293c0 .319.216.694.825.576A12 12 0 0 0 12 .297Z" /></svg>;
}

export function SiteHeader() {
  const { locale, setLocale } = useUiStore();
  const t = useCopy(locale);
  const [open, setOpen] = useState(false);
  const close = () => setOpen(false);
  return (
    <header className="sticky top-0 z-30 border-b border-line bg-[rgba(244,245,247,0.88)] backdrop-blur-[16px]">
      <div className="mx-auto flex min-h-16 w-[min(1200px,calc(100%-48px))] items-center gap-2 md:gap-8 max-[680px]:w-[calc(100%-34px)]">
        <Link className="inline-flex shrink-0 items-center gap-2.5 whitespace-nowrap text-[17px] font-bold tracking-[-0.01em] text-ink" to="/" onClick={close}>
          <BrandMark />
          Breeze UI
        </Link>
        <nav className={`${open ? "flex" : "hidden"} absolute left-4 right-4 top-[58px] flex-col items-start gap-3 rounded-[14px] border border-line bg-white p-4 text-[14px] text-muted shadow-[0_18px_36px_rgba(20,24,32,0.1)] min-[681px]:static min-[681px]:flex min-[681px]:flex-1 min-[681px]:flex-row min-[681px]:flex-wrap min-[681px]:items-center min-[681px]:gap-x-6 min-[681px]:gap-y-1 min-[681px]:border-0 min-[681px]:bg-transparent min-[681px]:p-0 min-[681px]:text-[15px] min-[681px]:shadow-none`} aria-label="Breeze UI">
          <Link to="/" onClick={close}>{locale === "zh" ? "首页" : "Home"}</Link>
          <Link to="/api" onClick={close}>{locale === "zh" ? "文档" : "Docs"}</Link>
          <Link to="/" hash="catalog" onClick={close}>{locale === "zh" ? "组件" : "Components"}</Link>
          <Link to="/" hash="tooling" onClick={close}>Storybook</Link>
          <Link to="/cli-mcp" onClick={close}>{locale === "zh" ? "CLI 与 MCP" : "CLI & MCP"}</Link>
        </nav>
        <div className="ml-auto flex items-center gap-1">
          <button className="inline-flex min-h-10 items-center gap-1.5 whitespace-nowrap rounded-[10px] border-0 bg-transparent px-2.5 text-sm text-muted hover:bg-[#e9ecef] hover:text-ink max-[680px]:min-h-9 max-[680px]:px-2" type="button" onClick={() => setLocale(locale === "zh" ? "en" : "zh")} aria-label="Change language">
            <Globe2 size={15} />{locale === "zh" ? "简体中文" : "English"}<ChevronDown size={13} />
          </button>
          <a className="inline-flex size-11 items-center justify-center rounded-[10px] text-ink hover:bg-[#e9ecef] max-[680px]:size-9" href="https://github.com/yldm-tech/mygo-breeze-ui" target="_blank" rel="noreferrer" aria-label="GitHub"><GitHubMark /></a>
          <Link className="ml-1 inline-flex min-h-10 items-center rounded-[10px] bg-ink px-4 text-sm font-semibold text-white max-[680px]:hidden" to="/" hash="install">{locale === "zh" ? "开始使用" : "Get started"}</Link>
          <button className="hidden size-9 items-center justify-center rounded-[10px] border-0 bg-transparent text-muted hover:bg-[#e9ecef] max-[680px]:inline-flex" type="button" onClick={() => setOpen((value) => !value)} aria-label="Toggle menu" aria-expanded={open}>
            {open ? <X size={20} /> : <Menu size={20} />}
          </button>
        </div>
      </div>
    </header>
  );
}

export { BrandMark };
