import { Link } from "@tanstack/react-router";
import { GitBranch, Menu, X } from "lucide-react";
import { useUiStore } from "../store";
import { useCopy } from "../copy";

export function SiteHeader() {
  const { locale, setLocale, mobileMenu, toggleMobileMenu } = useUiStore();
  const t = useCopy(locale);
  return (
    <header className="relative z-20 mx-auto flex h-[76px] w-[calc(100%-48px)] max-w-[1240px] items-center gap-8 max-[680px]:h-[66px] max-[680px]:w-[calc(100%-34px)]">
      <Link
        className="flex items-center gap-2.5 text-[18px] font-extrabold tracking-[-0.06em]"
        to="/"
      >
        <span className="grid size-[29px] place-items-center rounded-[9px] bg-mint font-extrabold tracking-normal text-bg">
          B
        </span>
        <span>
          breeze<span className="text-mint">.</span>ui
        </span>
      </Link>
      <nav
        className={`${mobileMenu ? "flex" : "hidden"} absolute top-[58px] right-0 left-0 flex-col items-start gap-3.5 rounded-lg border border-line bg-panel p-[17px] text-xs font-bold text-muted md:static md:ml-auto md:flex md:flex-row md:items-center md:gap-[27px] md:border-0 md:bg-transparent md:p-0`}
      >
        <Link
          className="inline-flex items-center gap-1.5 transition-colors hover:text-ink"
          to="/"
          hash="catalog"
          onClick={() => mobileMenu && toggleMobileMenu()}
        >
          {t.navComponents}
        </Link>
        <Link
          className="inline-flex items-center gap-1.5 transition-colors hover:text-ink"
          to="/api"
          onClick={() => mobileMenu && toggleMobileMenu()}
        >
          {t.navApi}
        </Link>
        <Link
          className="inline-flex items-center gap-1.5 transition-colors hover:text-ink"
          to="/"
          hash="install"
          onClick={() => mobileMenu && toggleMobileMenu()}
        >
          {t.navInstall}
        </Link>
        <a
          className="inline-flex items-center gap-1.5 transition-colors hover:text-ink"
          href="https://github.com/yldm-tech/mygo-breeze-ui"
          target="_blank"
          rel="noreferrer"
        >
          {t.navGithub} <GitBranch size={14} />
        </a>
      </nav>
      <div className="ml-auto flex items-center gap-2 md:ml-0">
        <button
          className="rounded-[7px] border border-line bg-[rgba(7,17,31,0.65)] px-[11px] py-[7px] text-[11px] text-muted transition hover:border-mint hover:text-mint"
          onClick={() => setLocale(locale === "en" ? "zh" : "en")}
        >
          {t.language}
        </button>
        <button
          className="grid rounded-[7px] border border-line bg-[rgba(7,17,31,0.65)] p-1.5 text-muted md:hidden"
          onClick={toggleMobileMenu}
          aria-label="Toggle menu"
          aria-expanded={mobileMenu}
        >
          {mobileMenu ? <X size={18} /> : <Menu size={18} />}
        </button>
      </div>
    </header>
  );
}
