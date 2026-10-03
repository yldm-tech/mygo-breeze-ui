import { Link } from "@tanstack/react-router";
import { GitBranch, Menu, X } from "lucide-react";
import { useUiStore } from "../store";
import { useCopy } from "../copy";

export function SiteHeader() {
  const { locale, setLocale, mobileMenu, toggleMobileMenu } = useUiStore();
  const t = useCopy(locale);
  return (
    <header className="site-header">
      <Link className="brand" to="/">
        <span className="brand-mark">B</span>
        <span>
          breeze<span className="brand-dot">.</span>ui
        </span>
      </Link>
      <nav className={mobileMenu ? "main-nav open" : "main-nav"}>
        <a href="#catalog" onClick={() => mobileMenu && toggleMobileMenu()}>
          {t.navComponents}
        </a>
        <Link to="/api" onClick={() => mobileMenu && toggleMobileMenu()}>
          {t.navApi}
        </Link>
        <a href="#install" onClick={() => mobileMenu && toggleMobileMenu()}>
          {t.navInstall}
        </a>
        <a href="https://github.com/yldm-tech/mygo-breeze-ui" target="_blank" rel="noreferrer">
          {t.navGithub} <GitBranch size={14} />
        </a>
      </nav>
      <div className="header-actions">
        <button
          className="language-button"
          onClick={() => setLocale(locale === "en" ? "zh" : "en")}
        >
          {t.language}
        </button>
        <button className="menu-button" onClick={toggleMobileMenu} aria-label="Toggle menu">
          {mobileMenu ? <X size={18} /> : <Menu size={18} />}
        </button>
      </div>
    </header>
  );
}
