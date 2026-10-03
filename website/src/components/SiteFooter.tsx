import { GitBranch } from "lucide-react";
import { Link } from "@tanstack/react-router";
import { useUiStore } from "../store";
import { useCopy } from "../copy";

export function SiteFooter() {
  const { locale } = useUiStore();
  const t = useCopy(locale);
  return (
    <footer className="mx-auto flex w-[calc(100%-48px)] max-w-[1160px] items-center gap-2 border-t border-line py-[23px] text-[10px] text-dim max-[680px]:w-[calc(100%-34px)] max-[680px]:flex-wrap">
      <Link
        className="mr-auto flex items-center gap-2.5 text-[14px] font-extrabold tracking-[-0.06em] text-ink max-[680px]:w-full"
        to="/"
      >
        <span className="grid size-6 place-items-center rounded-[7px] bg-mint text-[11px] font-extrabold tracking-normal text-bg">
          B
        </span>
        <span>
          breeze<span className="text-mint">.</span>ui
        </span>
      </Link>
      <span>{t.footer} · MIT License · ai-logo icons</span>
      <a
        className="ml-5 inline-flex items-center gap-1.5 text-muted hover:text-ink max-[680px]:ml-auto"
        href="https://github.com/yldm-tech/mygo-breeze-ui"
        target="_blank"
        rel="noreferrer"
      >
        <GitBranch size={14} /> {t.source}
      </a>
    </footer>
  );
}
