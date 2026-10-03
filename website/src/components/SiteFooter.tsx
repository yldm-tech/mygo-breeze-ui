import { GitBranch } from "lucide-react";
import { Link } from "@tanstack/react-router";
import { useUiStore } from "../store";
import { useCopy } from "../copy";

export function SiteFooter() {
  const { locale } = useUiStore();
  const t = useCopy(locale);
  return (
    <footer className="site-footer container">
      <Link className="brand" to="/">
        <span className="brand-mark">B</span>
        <span>
          breeze<span className="brand-dot">.</span>ui
        </span>
      </Link>
      <span>{t.footer} · MIT License</span>
      <a href="https://github.com/yldm-tech/mygo-breeze-ui" target="_blank" rel="noreferrer">
        <GitBranch size={14} /> {t.source}
      </a>
    </footer>
  );
}
