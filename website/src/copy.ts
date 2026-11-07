import type { Locale } from "./store";

export const copy = {
  en: {
    navComponents: "Components",
    navStorybook: "Storybook",
    navInstall: "Install",
    navGithub: "GitHub",
    language: "中文",
    eyebrow: "Native UI for Go, with a breeze.",
    title: "MyGo widgets that look like the desktop they run on.",
    intro:
      "Breeze UI draws MyGo's widgets as macOS 26, Windows 11 and GNOME draw their own controls. The names and signatures are those of MyGo's package ui, so ui.Button becomes breeze.Button and the rest of the view stays the same.",
    start: "Get started",
    explore: "Explore components",
    components: "components",
    platforms: "platforms",
    drawn: "drawn per platform",
    calmTitle: "One view, three native looks",
    calmText:
      "On a Mac the controls look like AppKit's, on Windows 11 like WinUI 3's, on GNOME like libadwaita's. MyGo draws them all; the view code is written once.",
    principle: "An app should look like it belongs on the desktop it runs on.",
    principleSource: "— Breeze UI principles",
    platformTitle: "One set of names, three platforms",
    platformText:
      "breeze.Use(c) gives the window the platform's theme. SetPlatform forces one platform, and Preview draws a subtree as another, for galleries side by side.",
    componentTitle: "Behavior from MyGo, looks from the platform",
    componentText:
      "Pointer, keyboard, focus and accessibility come from MyGo's widget bases. Button, PrimaryButton, Checkbox and Switch are drawn to match each platform; the other widgets are MyGo's own in the platform's theme.",
    toolingTitle: "Tooling that stays out of the way",
    toolingText:
      "Use the CLI to list components and platforms and generate snippets, or connect the MCP server to your coding agent.",
    aiKicker: "AI-native tooling",
    aiTitle: "A familiar home for your AI stack",
    aiText:
      "Bring native MyGo surfaces to the models and coding agents your team already uses. Breeze UI stays focused on the interface layer.",
    aiSource: "Powered by ai-logo",
    aiOpenAi: "Models and agents",
    aiClaude: "Thoughtful workflows",
    aiGemini: "Multimodal products",
    aiCursor: "Fast iteration",
    aiCodex: "Code-first tools",
    aiDeepSeek: "Open model workflows",
    aiDisclaimer: "Brand marks are shown for identification only and do not imply endorsement.",
    catalogKicker: "The catalog",
    catalogTitle: "A small, deliberate catalog",
    catalogText:
      "Components carry the names of MyGo's package ui. Drawn ones are drawn by Breeze UI to match the platform; themed ones are MyGo's widgets in the platform's theme.",
    all: "All",
    component: "Components",
    platform: "Platforms",
    installKicker: "Start here",
    installTitle: "One package. Three platforms.",
    installText:
      "Add Breeze UI to a MyGo app and call breeze.Use(c), or run the gallery to see every control as macOS, Windows 11 and GNOME draw it, side by side.",
    docs: "Read the docs",
    copy: "Copy",
    copied: "Copied",
    starTitle: "Build in the open",
    starText: "Breeze UI is MIT licensed and evolving with the MyGo ecosystem.",
    source: "Source on GitHub",
    footer: "Breeze UI for MyGo",
  },
  zh: {
    navComponents: "组件",
    navStorybook: "Storybook",
    navInstall: "安装",
    navGithub: "GitHub",
    language: "English",
    eyebrow: "为 Go 打造的原生 UI，轻盈如风。",
    title: "让 MyGo 组件和所在桌面的控件一个样。",
    intro:
      "Breeze UI 按 macOS 26、Windows 11 和 GNOME 绘制自家控件的方式来绘制 MyGo 的组件。名称和签名与 MyGo 的 ui 包一致，把 ui.Button 换成 breeze.Button，视图里的其他代码不用改。",
    start: "开始使用",
    explore: "浏览组件",
    components: "个组件",
    platforms: "个平台",
    drawn: "个按平台自绘",
    calmTitle: "一份视图，三种原生外观",
    calmText:
      "控件在 Mac 上看起来像 AppKit 的，在 Windows 11 上像 WinUI 3 的，在 GNOME 上像 libadwaita 的。它们都由 MyGo 绘制，视图代码只写一份。",
    principle: "应用应该看起来就属于它运行的那个桌面。",
    principleSource: "— Breeze UI 设计原则",
    platformTitle: "一套名称，三个平台",
    platformText:
      "breeze.Use(c) 为窗口套用当前平台的主题。SetPlatform 可以强制指定平台，Preview 可以把一段界面按另一个平台绘制，用来并排对比。",
    componentTitle: "行为来自 MyGo，外观来自平台",
    componentText:
      "指针、键盘、焦点和无障碍由 MyGo 的组件基座提供。Button、PrimaryButton、Checkbox 和 Switch 由 Breeze UI 按各平台绘制，其余组件是 MyGo 自带的组件，套用平台主题。",
    toolingTitle: "不会碍事的工具链",
    toolingText: "用 CLI 列出组件和平台、生成代码片段，或把 MCP 服务接入你的编程 Agent。",
    aiKicker: "AI 原生工具链",
    aiTitle: "为你的 AI 工具栈留出一层原生界面",
    aiText:
      "把 MyGo 的原生窗口带到团队已经在使用的模型和编程 Agent 旁边。Breeze UI 专注于界面这一层。",
    aiSource: "由 ai-logo 提供图标",
    aiOpenAi: "模型与 Agent",
    aiClaude: "更从容的工作流",
    aiGemini: "多模态产品",
    aiCursor: "快速迭代",
    aiCodex: "代码优先工具",
    aiDeepSeek: "开放模型工作流",
    aiDisclaimer: "品牌标识仅用于识别，不代表相关方背书。",
    catalogKicker: "组件目录",
    catalogTitle: "小而明确的组件目录",
    catalogText:
      "组件与 MyGo ui 包里的组件同名。drawn 表示由 Breeze UI 按平台绘制，themed 表示 MyGo 自带组件套用平台主题。",
    all: "全部",
    component: "组件",
    platform: "平台",
    installKicker: "从这里开始",
    installTitle: "一个包，三个平台。",
    installText:
      "把 Breeze UI 加进 MyGo 应用并调用 breeze.Use(c)，或运行示例画廊，并排查看每个控件在 macOS、Windows 11 和 GNOME 上的样子。",
    docs: "阅读文档",
    copy: "复制",
    copied: "已复制",
    starTitle: "在开源中一起构建",
    starText: "Breeze UI 使用 MIT 协议，和 MyGo 生态一起持续演进。",
    source: "在 GitHub 查看源码",
    footer: "Breeze UI for MyGo",
  },
} as const;

export function useCopy(locale: Locale) {
  return copy[locale];
}
