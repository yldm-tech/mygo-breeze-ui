import { Link } from "@tanstack/react-router";
import { BrandMark } from "./SiteHeader";
import { useUiStore } from "../store";

export function SiteFooter() {
  const { locale } = useUiStore();
  const zh = locale === "zh";
  return (
    <footer className="site-footer">
      <div className="footer-inner">
        <div className="footer-brand">
          <span><BrandMark />Breeze UI</span>
          <small>{zh ? "MyGo 的原生 UI 组件。MIT License。" : "Native UI widgets for MyGo. MIT License."}</small>
        </div>
        <nav className="footer-nav" aria-label={zh ? "文档" : "Docs"}>
          <strong>{zh ? "文档" : "Docs"}</strong>
          <Link to="/docs">{zh ? "快速开始" : "Quick start"}</Link>
          <Link to="/docs" hash="platforms">{zh ? "平台与主题" : "Platforms & themes"}</Link>
          <Link to="/" hash="catalog">{zh ? "组件" : "Components"}</Link>
          <Link to="/cli-mcp">{zh ? "CLI 与 MCP" : "CLI & MCP"}</Link>
        </nav>
        <nav className="footer-nav" aria-label={zh ? "资源" : "Resources"}>
          <strong>{zh ? "资源" : "Resources"}</strong>
          <a href="https://breeze-ui.yldm.ai/storybook/" target="_blank" rel="noreferrer">Storybook</a>
          <a href="https://pkg.go.dev/github.com/yldm-tech/mygo-breeze-ui" target="_blank" rel="noreferrer">pkg.go.dev</a>
          <a href="https://github.com/egoist/mygo" target="_blank" rel="noreferrer">MyGo</a>
        </nav>
        <nav className="footer-nav" aria-label={zh ? "社区" : "Community"}>
          <strong>{zh ? "社区" : "Community"}</strong>
          <a href="https://github.com/yldm-tech/mygo-breeze-ui" target="_blank" rel="noreferrer">GitHub</a>
          <a href="https://github.com/yldm-tech/mygo-breeze-ui/issues" target="_blank" rel="noreferrer">{zh ? "问题反馈" : "Issues"}</a>
        </nav>
      </div>
    </footer>
  );
}
