import { Link } from "@tanstack/react-router";
import { BrandMark } from "./SiteHeader";
import { useUiStore } from "../store";

export function SiteFooter() {
  const { locale } = useUiStore();
  const zh = locale === "zh";
  return (
    <footer className="border-t border-line bg-white">
      <div className="mx-auto grid w-[min(1200px,calc(100%-48px))] grid-cols-4 gap-8 py-12 text-sm max-[680px]:w-[calc(100%-34px)] max-[640px]:grid-cols-2 max-[640px]:gap-x-5 max-[640px]:gap-y-[26px] max-[640px]:py-9">
        <div className="flex flex-col gap-2.5">
          <span className="flex items-center gap-2 text-[15px] font-bold"><BrandMark />Breeze UI</span>
          <small className="text-muted">{zh ? "MyGo 的原生 UI 组件。MIT License。" : "Native UI widgets for MyGo. MIT License."}</small>
        </div>
        <nav className="flex flex-col gap-2" aria-label={zh ? "文档" : "Docs"}>
          <strong className="font-semibold">{zh ? "文档" : "Docs"}</strong>
          <Link className="text-muted hover:text-mint" to="/docs">{zh ? "快速开始" : "Quick start"}</Link>
          <Link className="text-muted hover:text-mint" to="/docs" hash="platforms">{zh ? "平台与主题" : "Platforms & themes"}</Link>
          <Link className="text-muted hover:text-mint" to="/docs/components">{zh ? "组件" : "Components"}</Link>
          <Link className="text-muted hover:text-mint" to="/cli-mcp">{zh ? "CLI 与 MCP" : "CLI & MCP"}</Link>
        </nav>
        <nav className="flex flex-col gap-2" aria-label={zh ? "资源" : "Resources"}>
          <strong className="font-semibold">{zh ? "资源" : "Resources"}</strong>
          <a className="text-muted hover:text-mint" href="https://breeze-ui.yldm.ai/storybook/" target="_blank" rel="noreferrer">Storybook</a>
          <a className="text-muted hover:text-mint" href="https://pkg.go.dev/github.com/yldm-tech/mygo-breeze-ui" target="_blank" rel="noreferrer">pkg.go.dev</a>
          <a className="text-muted hover:text-mint" href="https://github.com/egoist/mygo" target="_blank" rel="noreferrer">MyGo</a>
        </nav>
        <nav className="flex flex-col gap-2" aria-label={zh ? "社区" : "Community"}>
          <strong className="font-semibold">{zh ? "社区" : "Community"}</strong>
          <a className="text-muted hover:text-mint" href="https://github.com/yldm-tech/mygo-breeze-ui" target="_blank" rel="noreferrer">GitHub</a>
          <a className="text-muted hover:text-mint" href="https://github.com/yldm-tech/mygo-breeze-ui/issues" target="_blank" rel="noreferrer">{zh ? "问题反馈" : "Issues"}</a>
        </nav>
      </div>
    </footer>
  );
}
