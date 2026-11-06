# Breeze UI for MyGo

[![CI](https://github.com/yldm-tech/mygo-breeze-ui/actions/workflows/ci.yml/badge.svg)](https://github.com/yldm-tech/mygo-breeze-ui/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/yldm-tech/mygo-breeze-ui.svg)](https://pkg.go.dev/github.com/yldm-tech/mygo-breeze-ui)
[![License](https://img.shields.io/github/license/yldm-tech/mygo-breeze-ui)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/yldm-tech/mygo-breeze-ui)](https://github.com/yldm-tech/mygo-breeze-ui/stargazers)

[官网](https://yldm-tech.github.io/mygo-breeze-ui/) · [English](README.md) · [简体中文](README.zh-CN.md)

Breeze UI 让 [MyGo](https://github.com/egoist/mygo) 的原生 UI 控件按应用所在的系统绘制：在 Mac 上是 macOS 的控件，在 Windows 上是 Windows 11 的控件，在 Linux 上是 GNOME libadwaita 的控件。控件的名字和签名与 MyGo 的 `ui` 包完全一致，把 `ui.Button` 换成 `breeze.Button` 即可，其余代码不用改。

![同一个界面分别按 macOS、Windows 11 和 GNOME 绘制](docs/assets/gallery-light.png)

## 包含内容

- `ui/breeze`：按钮、切换按钮与切换组、复选框、单选按钮与单选按钮组、开关、文本框与搜索框（以及编辑 `ui.TextBuffer` 的文本框）、组合框、下拉选择、滑块、步进器、进度条、加载指示器、分段控件、标签页、可关闭标签并带新建按钮的文档标签栏、按平台字号层级排版的文字（标题、小标题、正文、说明文字和弱化的次要文字）、链接、分隔线、徽标、菜单按钮、弹出层、对话框和提示框、工具提示和 Toast 通知、折叠面板与手风琴、字段组与表单、卡片、设置行与开关行、下拉选择行、输入行、数值行与展开行、侧边栏、面包屑、后退与前进按钮、工具栏、隐藏标题栏窗口的标题栏（HeaderBar）、切换应用顶层视图的视图切换器（ViewSwitcher）、电平表、头像、评分、内联消息栏、环形进度、空状态页、带刻度的滑块、数字输入框、日期与时间输入框、按平台绘制月历（NSDatePicker、WinUI 的 CalendarView、GtkCalendar）并跟随用户语言与地区格式的日期选择器、颜色井、复选框组、分栏视图、列表、表格、树、大纲、网格、可编辑文本、自动补全与标记输入框、密码框、拆分按钮和图标按钮按各平台的样子绘制，对照 macOS 27、WinUI 3 和 libadwaita 1.10（GNOME 51），尺寸和颜色量自 AppKit 与 libadwaita 的实际渲染，或取自 WinUI 3 的主题资源；小数输入框、范围滑块、日历、颜色选择器、查找栏和滚动容器以同样的名字使用 `ui` 包的控件，配上平台主题。
- 三个平台的浅色和深色主题，跟随用户选择的强调色、对比度和文字大小。
- `breeze-ui` CLI：列出控件和平台，生成代码片段。
- `breeze-ui-mcp`：基于同一份目录的只读 MCP 服务。

## Go 包

```sh
go get github.com/yldm-tech/mygo-breeze-ui
```

```go
import (
    "github.com/egoist/mygo/ui"
    breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

func (app *settings) view(c *ui.Context) {
    breeze.Use(c) // the platform's colors and metrics
    ui.Column(c).Padding(16).Gap(12).Children(func() {
        breeze.Checkbox(c, &app.remember, "Remember me")
        breeze.Switch(c, &app.sync).Label("Sync")
        ui.Row(c).Gap(8).Children(func() {
            breeze.Button(c, "Cancel")
            if breeze.PrimaryButton(c, "Save").Clicked() {
                app.save()
            }
        })
    })
}
```

默认按当前系统选择平台，`breeze.SetPlatform` 可以指定其他平台；`breeze.Preview(c, p, body)` 把窗口的一部分按另一个平台绘制，示例 gallery 就是这样把三个平台并排显示的：

```sh
go run ./examples/breeze-ui
```

指针、键盘、焦点和无障碍都来自 MyGo 的无外观控件（`ui.ButtonBase`、`ui.CheckboxBase`、`ui.SwitchBase`）。每个控件在各平台上的样子、数值的来源和目前的限制见 [docs/breeze-ui.md](docs/breeze-ui.md)。

## Storybook

每个组件在三个平台、浅色和深色下的样子，带可调参数和对应的 Go 代码：**[在线 Storybook](https://yldm-tech.github.io/mygo-breeze-ui/storybook/)** 在浏览器里以 WebAssembly 运行，同一个程序也能在本机以窗口打开：

```sh
go run ./examples/storybook
```

## CLI

```sh
npx breeze-ui list
npx breeze-ui add Switch -o ui/switch_snippet.go
```

也可以安装 Go 命令：

```sh
go install github.com/yldm-tech/mygo-breeze-ui/cmd/breeze-ui@latest
```

## MCP 服务

```sh
bun install
bun run build
```

`breeze-ui-mcp` 提供只读工具，用来搜索组件、列出平台和生成代码片段。

## 开发

```sh
go test ./...
go run ./examples/breeze-ui
BREEZE_SCREENSHOTS=/tmp/shots go test ./ui/breeze -run TestScreenshots   # write each platform's controls as PNG
bun install
bun run test
bun run typecheck
bun run build
```

Pull request 标题和提交信息使用 [Conventional Commits](https://www.conventionalcommits.org/)（`feat(breeze): ...`、`fix: ...`）。Labels 工作流会检查标题，并按类型、大小、区域和平台给每个 pull request 打标签；Conflicts 工作流会给与目标分支冲突的 pull request 打标签。`main` 上的 CI 通过后，如果上一个 tag 之后的提交需要发版，Release 工作流就会打 tag 并发布 GitHub Release：`feat` 升次版本号，`fix`、`perf` 和 `revert` 升修订号，破坏性改动（`feat!:` 或 `BREAKING CHANGE:` 脚注）升主版本号，版本仍是 0.x 时只升次版本号。需要强制发版时可以手动运行 Release 工作流。

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=yldm-tech/mygo-breeze-ui&type=Date)](https://star-history.com/#yldm-tech/mygo-breeze-ui&Date)

## 许可证

[MIT](LICENSE)
