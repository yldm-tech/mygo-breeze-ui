// Command storybook shows every Breeze UI widget as Storybook shows components: stories in a sidebar, each drawn on a canvas as macOS, Windows 11 or GNOME draws it, or as all three side by side, light or dark, with controls that change its arguments and the Go that builds it. It runs as a window, or in a browser compiled to WebAssembly:
//
//	go run ./examples/storybook
//	GOOS=js GOARCH=wasm go build -o website/public/storybook/app.wasm ./examples/storybook
package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/egoist/mygo"
	_ "github.com/egoist/mygo/locales/zh-Hans"
	"github.com/egoist/mygo/ui"
	breeze "github.com/yldm-tech/mygo-breeze-ui/ui/breeze"
)

// Platform choices of the toolbar, in order; the last shows the three side by side.
var platformNames = []string{"macOS", "Windows", "GNOME", "Side by side"}

var storybookPreviewTransition = ui.ElementTransition{
	Duration: 220 * time.Millisecond,
	Ease:     ui.EaseInOut,
	Enter:    &ui.Motion{Opacity: 0, Scale: 0.97},
}

// Storybook owns translations for its chrome in these languages. Registering
// them with MyGo makes App.Locale follow navigator.languages in WebAssembly
// (and the user's preferred languages on desktop) instead of falling back to
// English before textFor can choose a translation.
var storybookLanguages = []string{"en", "zh-Hans", "ja", "ko", "es", "fr", "de"}

type storybookText struct {
	locale                                        storybookLocale
	components, subtitle, search, stories         string
	story, live, platform, appearance             string
	light, dark, sideBySide, controls, noControls string
	goCode, copy                                  string
}

type storybookLocale string

const (
	localeEN storybookLocale = "en"
	localeZH storybookLocale = "zh"
	localeJA storybookLocale = "ja"
	localeKO storybookLocale = "ko"
	localeES storybookLocale = "es"
	localeFR storybookLocale = "fr"
	localeDE storybookLocale = "de"
)

func localeFor(c *ui.Context) storybookLocale {
	raw := strings.ToLower(c.Locale())
	switch {
	case strings.HasPrefix(raw, "zh"):
		return localeZH
	case strings.HasPrefix(raw, "ja"):
		return localeJA
	case strings.HasPrefix(raw, "ko"):
		return localeKO
	case strings.HasPrefix(raw, "es"):
		return localeES
	case strings.HasPrefix(raw, "fr"):
		return localeFR
	case strings.HasPrefix(raw, "de"):
		return localeDE
	default:
		return localeEN
	}
}

func textFor(c *ui.Context) storybookText {
	l := localeFor(c)
	texts := map[storybookLocale]storybookText{
		localeZH: {locale: localeZH, components: "组件", subtitle: "交互式 story 和平台预览", search: "查找组件", stories: "故事", story: "故事", live: "实时", platform: "平台", appearance: "外观", light: "浅色", dark: "深色", sideBySide: "并排", controls: "控件", noControls: "这个 story 没有可调参数。", goCode: "Go 代码", copy: "复制"},
		localeJA: {locale: localeJA, components: "コンポーネント", subtitle: "インタラクティブなストーリーとプラットフォームプレビュー", search: "コンポーネントを検索", stories: "ストーリー", story: "ストーリー", live: "ライブ", platform: "プラットフォーム", appearance: "外観", light: "ライト", dark: "ダーク", sideBySide: "横並び", controls: "コントロール", noControls: "このストーリーには調整可能なパラメータがありません。", goCode: "Go コード", copy: "コピー"},
		localeKO: {locale: localeKO, components: "컴포넌트", subtitle: "인터랙티브 스토리 및 플랫폼 미리보기", search: "컴포넌트 찾기", stories: "스토리", story: "스토리", live: "라이브", platform: "플랫폼", appearance: "모양", light: "밝은 테마", dark: "어두운 테마", sideBySide: "나란히 보기", controls: "컨트롤", noControls: "이 스토리에는 조정 가능한 매개변수가 없습니다.", goCode: "Go 코드", copy: "복사"},
		localeES: {locale: localeES, components: "Componentes", subtitle: "Historias interactivas y vistas previas de plataforma", search: "Buscar un componente", stories: "Historias", story: "HISTORIA", live: "EN VIVO", platform: "Plataforma", appearance: "Apariencia", light: "Claro", dark: "Oscuro", sideBySide: "Lado a lado", controls: "Controles", noControls: "Esta historia no tiene controles.", goCode: "Código Go", copy: "Copiar"},
		localeFR: {locale: localeFR, components: "Composants", subtitle: "Histoires interactives et aperçus des plateformes", search: "Rechercher un composant", stories: "Histoires", story: "HISTOIRE", live: "EN DIRECT", platform: "Plateforme", appearance: "Apparence", light: "Clair", dark: "Sombre", sideBySide: "Côte à côte", controls: "Contrôles", noControls: "Cette histoire n’a aucun contrôle.", goCode: "Code Go", copy: "Copier"},
		localeDE: {locale: localeDE, components: "Komponenten", subtitle: "Interaktive Stories und Plattformvorschauen", search: "Komponente suchen", stories: "Stories", story: "STORY", live: "LIVE", platform: "Plattform", appearance: "Darstellung", light: "Hell", dark: "Dunkel", sideBySide: "Nebeneinander", controls: "Steuerelemente", noControls: "Diese Story hat keine Steuerelemente.", goCode: "Go-Code", copy: "Kopieren"},
		localeEN: {locale: localeEN, components: "Components", subtitle: "Interactive stories and platform previews", search: "Find a component", stories: "Stories", story: "STORY", live: "LIVE", platform: "Platform", appearance: "Appearance", light: "Light", dark: "Dark", sideBySide: "Side by side", controls: "Controls", noControls: "This story has no controls.", goCode: "Go code", copy: "Copy"},
	}
	if text, ok := texts[l]; ok {
		return text
	}
	return texts[localeEN]
}

func categoryLabel(category string, locale storybookLocale) string {
	labels := map[storybookLocale]map[string]string{
		localeZH: {"Actions": "操作", "Choices": "选择", "Text": "文本", "Feedback": "反馈", "Navigation": "导航", "Overlays": "浮层", "Layout": "布局", "Data": "数据"},
		localeJA: {"Actions": "操作", "Choices": "選択", "Text": "テキスト", "Feedback": "フィードバック", "Navigation": "ナビゲーション", "Overlays": "オーバーレイ", "Layout": "レイアウト", "Data": "データ"},
		localeKO: {"Actions": "동작", "Choices": "선택", "Text": "텍스트", "Feedback": "피드백", "Navigation": "탐색", "Overlays": "오버레이", "Layout": "레이아웃", "Data": "데이터"},
		localeES: {"Actions": "Acciones", "Choices": "Opciones", "Text": "Texto", "Feedback": "Comentarios", "Navigation": "Navegación", "Overlays": "Superposiciones", "Layout": "Diseño", "Data": "Datos"},
		localeFR: {"Actions": "Actions", "Choices": "Choix", "Text": "Texte", "Feedback": "Retour", "Navigation": "Navigation", "Overlays": "Superpositions", "Layout": "Mise en page", "Data": "Données"},
		localeDE: {"Actions": "Aktionen", "Choices": "Auswahl", "Text": "Text", "Feedback": "Feedback", "Navigation": "Navigation", "Overlays": "Overlays", "Layout": "Layout", "Data": "Daten"},
	}
	if label, ok := labels[locale][category]; ok {
		return label
	}
	return category
}

type navEntry struct {
	category string
	story    *story
}

type book struct {
	stories  []*story
	selected string // the chosen story's id
	platform int    // an index of platformNames
	dark     int    // 0 light, 1 dark
	query    string
	open     map[string]*bool // the sidebar sections shown
	args     map[string]*args // each story's argument values, kept as stories change
	navRows  []navEntry
	navList  ui.ListState
	navIndex int
	navID    string
}

func newBook() *book {
	b := &book{stories: stories(), platform: 3, open: map[string]*bool{}, args: map[string]*args{}}
	b.selected = b.stories[0].id()
	for _, cat := range categories {
		shown := true
		b.open[cat] = &shown
	}
	return b
}

func (b *book) story() *story {
	for _, s := range b.stories {
		if s.id() == b.selected {
			return s
		}
	}
	return b.stories[0]
}

func (b *book) argsOf(s *story) *args {
	a, ok := b.args[s.id()]
	if !ok {
		a = newArgs(s.controls)
		b.args[s.id()] = a
	}
	return a
}

func (b *book) view(c *ui.Context) {
	c.SetTheme(storybookTheme())
	ui.Row(c).Fill().AlignItems(ui.Stretch).Children(func() {
		b.sidebar(c)
		ui.Column(c).Grow(1).MinWidth(0).Children(func() {
			s := b.story()
			b.toolbar(c, s)
			b.canvas(c, s)
			b.panel(c, s)
		})
	})
}

// storybookTheme keeps the application chrome in the same visual language as
// the website: a soft gray canvas, white panels, ink text and Breeze mint
// accents. The platform previews still use their own themes below.
func storybookTheme() *ui.Theme {
	t := *ui.LightTheme()
	t.Background = ui.Hex("#f4f5f7")
	t.Surface = ui.Hex("#ffffff")
	t.SurfaceHover = ui.Hex("#f8f9fb")
	t.SurfacePressed = ui.Hex("#e2e4e9")
	t.Border = ui.Hex("#e2e4e9")
	t.Text = ui.Hex("#15171c")
	t.TextMuted = ui.Hex("#555c69")
	t.Accent = ui.Hex("#0b7a63")
	t.AccentHover = ui.Hex("#075a49")
	t.AccentPressed = ui.Hex("#064b3e")
	t.AccentText = ui.Hex("#ffffff")
	t.Selection = ui.Hex("#a6e1cd")
	t.Focus = ui.RGBA(11, 122, 99, 0.35)
	t.Radius = 10
	t.Spacing = 4
	t.FontSize = 14
	return &t
}

// sidebar lists the stories by category, filtered by the search field.
func (b *book) sidebar(c *ui.Context) {
	t := c.Theme()
	copy := textFor(c)
	ui.Column(c).Width(288).Shrink(0).Background(t.Surface).BorderWidth(0, 1, 0, 0).BorderColor(t.Border).Children(func() {
		ui.Column(c).Padding(18, 18, 14, 18).Gap(10).Children(func() {
			ui.Row(c).AlignItems(ui.Center).Children(func() {
				ui.Text(c, copy.components).FontSize(t.FontSize + 3).Bold().Grow(1)
				ui.Text(c, fmt.Sprintf("%d", len(b.stories))).Padding(4, 8).Radius(999).Background(t.SurfaceHover).TextColor(t.TextMuted).FontSize(t.FontSize - 2)
			})
			ui.Text(c, copy.subtitle).TextColor(t.TextMuted).FontSize(t.FontSize - 1)
			ui.SearchField(c, &b.query).Label(copy.search)
		})
		q := strings.ToLower(strings.TrimSpace(b.query))
		b.navRows = b.navRows[:0]
		for _, cat := range categories {
			var matches []*story
			for _, s := range b.stories {
				if s.category == cat && (q == "" || strings.Contains(strings.ToLower(s.component+" "+s.name), q)) {
					matches = append(matches, s)
				}
			}
			if len(matches) == 0 {
				continue
			}
			b.navRows = append(b.navRows, navEntry{category: cat})
			if shown := b.open[cat]; shown != nil && *shown {
				for _, s := range matches {
					b.navRows = append(b.navRows, navEntry{category: cat, story: s})
				}
			}
		}
		if b.navID != b.selected {
			b.navID = b.selected
			b.navIndex = -1
			for i, entry := range b.navRows {
				if entry.story != nil && entry.story.id() == b.selected {
					b.navIndex = i
					break
				}
			}
		}
		b.navList.Selected = &b.navIndex
		b.navList.RoundRows = true
		b.navList.Key = func(i int) any {
			entry := b.navRows[i]
			if entry.story == nil {
				return "category/" + entry.category
			}
			return entry.story.id()
		}
		b.navList.Header = func(i int) bool { return b.navRows[i].story == nil }
		b.navList.Label = func(i int) string {
			entry := b.navRows[i]
			if entry.story == nil {
				return categoryLabel(entry.category, copy.locale)
			}
			return entry.story.component + " " + entry.story.name
		}
		list := ui.List(c, &b.navList, len(b.navRows)).Rows(func(row ui.ListRow) {
			entry := b.navRows[row.Index]
			if entry.story == nil {
				head := ui.ButtonBase(row.Context).FillWidth().Padding(8, 10).Radius(8).Background(t.Surface).AlignItems(ui.Center)
				head.Children(func() {
					arrow := "▾"
					if shown := b.open[entry.category]; shown == nil || !*shown {
						arrow = "▸"
					}
					ui.Text(row.Context, arrow).TextColor(t.TextMuted).Width(16)
					ui.Text(row.Context, categoryLabel(entry.category, copy.locale)).TextColor(t.TextMuted).Bold().Grow(1)
				})
				if head.Clicked() {
					*b.open[entry.category] = !*b.open[entry.category]
				}
				return
			}
			item := ui.Row(row.Context).FillWidth().Gap(8).Padding(7, 10).Radius(8).AlignItems(ui.Center)
			if row.Selected() {
				item.Background(t.SurfacePressed).TextColor(t.Text)
			}
			item.Children(func() {
				ui.Text(row.Context, entry.story.component+" · "+entry.story.name).Grow(1).SingleLine()
			})
		})
		list.Grow(1).Label(copy.stories)
		if list.Changed() && b.navIndex >= 0 && b.navIndex < len(b.navRows) && b.navRows[b.navIndex].story != nil {
			b.selected = b.navRows[b.navIndex].story.id()
			b.navID = b.selected
		}
	})
}

// toolbar names the story and chooses the platform and the appearance.
func (b *book) toolbar(c *ui.Context, s *story) {
	t := c.Theme()
	copy := textFor(c)
	ui.Column(c).Padding(14, 22, 16, 22).Gap(12).Background(t.Surface).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
		ui.Row(c).Gap(14).AlignItems(ui.Center).Wrap().Children(func() {
			ui.Row(c).Grow(1).Gap(9).AlignItems(ui.Center).Children(func() {
				ui.Text(c, copy.story).Padding(4, 8).Radius(999).Background(t.Accent).TextColor(t.AccentText).FontSize(t.FontSize - 2).Bold()
				ui.Text(c, s.component).Bold().FontSize(t.FontSize + 4)
				ui.Text(c, "/").TextColor(t.Border)
				ui.Text(c, s.name).TextColor(t.TextMuted)
			})
			ui.Text(c, copy.live).Padding(4, 8).Radius(999).Background(ui.RGBA(11, 122, 99, 0.12)).TextColor(t.Accent).FontSize(t.FontSize - 2).Bold()
		})
		ui.Row(c).Gap(10).AlignItems(ui.Center).Wrap().Children(func() {
			ui.Row(c).Gap(8).Padding(5, 8).Radius(10).Background(t.SurfaceHover).Border(1, t.Border).AlignItems(ui.Center).Children(func() {
				ui.Text(c, copy.platform).TextColor(t.TextMuted).FontSize(t.FontSize - 1).Bold()
				platforms := platformNames
				platforms = []string{"macOS", "Windows", "GNOME", copy.sideBySide}
				ui.Segmented(c, &b.platform, platforms...)
			})
			ui.Row(c).Gap(8).Padding(5, 8).Radius(10).Background(t.SurfaceHover).Border(1, t.Border).AlignItems(ui.Center).Children(func() {
				ui.Text(c, copy.appearance).TextColor(t.TextMuted).FontSize(t.FontSize - 1).Bold()
				ui.Segmented(c, &b.dark, copy.light, copy.dark)
			})
		})
	})
}

// canvas draws the story as the chosen platform, or as each platform in a column of its own.
func (b *book) canvas(c *ui.Context, s *story) {
	t := c.Theme()
	a := b.argsOf(s)
	platforms := []breeze.Platform{breeze.Platforms[min(b.platform, 2)]}
	if b.platform == 3 {
		platforms = breeze.Platforms
	}
	appearance := ui.LightTheme()
	if b.dark == 1 {
		appearance = ui.DarkTheme()
	}
	ui.Row(c).Grow(1).MinHeight(260).AlignItems(ui.Stretch).Padding(20, 22).Gap(16).Background(t.Background).Children(func() {
		for i, p := range platforms {
			name := platformNames[i]
			if len(platforms) == 1 {
				name = platformNames[b.platform]
			}
			// The appearance's theme first, which breeze.Theme reads light or dark from, then the platform's.
			ui.Themed(c, appearance).Key(fmt.Sprintf("storybook-preview-%s-%s-%d-%d", s.id(), p.String(), b.platform, b.dark)).Transition(storybookPreviewTransition).Grow(1).Basis(0).MinWidth(0).Radius(16).Background(t.Surface).Border(1, t.Border).Shadow(0, 8, 24, 0, ui.RGBA(20, 24, 32, 0.06)).Clip().Children(func() {
				breeze.Preview(c, p, func() {
					ui.Themed(c, breeze.Theme(c, p)).Fill().Padding(20).Gap(16).Children(func() {
						if len(platforms) > 1 {
							ui.Text(c, name).Bold().TextColor(c.Theme().TextMuted)
						}
						ui.Column(c).Grow(1).Center().Gap(12).Children(func() { s.render(c, a) })
					})
				})
			})
		}
	})
}

// panel shows the story's controls and the Go that builds it.
func (b *book) panel(c *ui.Context, s *story) {
	t := c.Theme()
	copy := textFor(c)
	a := b.argsOf(s)
	ui.Row(c).Height(250).Shrink(0).AlignItems(ui.Stretch).Background(t.Surface).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
		ui.Scroll(c).Grow(1).Basis(0).Padding(18, 22).Gap(12).Children(func() {
			ui.Text(c, copy.controls).TextColor(t.Accent).Bold()
			if len(s.controls) == 0 {
				ui.Text(c, copy.noControls).TextColor(t.TextMuted)
			}
			ui.Form(c, func() {
				for _, ctl := range s.controls {
					ui.Field(c, ctl.name, func() { b.control(c, a, ctl) })
				}
			})
		})
		ui.Column(c).Grow(1).Basis(0).Padding(18, 22).Gap(12).Background(t.SurfaceHover).BorderWidth(0, 0, 0, 1).BorderColor(t.Border).Children(func() {
			ui.Row(c).AlignItems(ui.Center).Children(func() {
				ui.Text(c, copy.goCode).TextColor(t.Accent).Bold().Grow(1)
				if ui.Button(c, copy.copy).Clicked() {
					c.WriteClipboard(s.code(a))
				}
			})
			ui.Text(c, s.code(a)).Font("monospace").Selectable().TextColor(t.Text)
		})
	})
}

// control builds the editor of one of the story's arguments.
func (b *book) control(c *ui.Context, a *args, ctl control) {
	switch ctl.kind {
	case textControl:
		ui.TextInput(c, a.text[ctl.name]).Width(240)
	case toggleControl:
		ui.Switch(c, a.on[ctl.name])
	case numberControl:
		v := a.num[ctl.name]
		ui.Row(c).Gap(10).AlignItems(ui.Center).Children(func() {
			ui.Slider(c, v, ctl.lo, ctl.hi).Width(180)
			ui.Text(c, formatNumber(*v, ctl.hi)).TextColor(c.Theme().TextMuted)
		})
	case choiceControl:
		ui.Select(c, a.text[ctl.name], ctl.options)
	}
}

func formatNumber(v, hi float64) string {
	if hi <= 1 {
		return fmt.Sprintf("%.2f", v)
	}
	return fmt.Sprintf("%.0f", v)
}

func main() {
	mygo.App.SetLanguages(storybookLanguages...)
	if locale := browserLocale(); locale != "" {
		mygo.App.SetLocale(locale)
	}
	b := newBook()
	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title:   "Breeze UI Storybook",
			Width:   1280,
			Height:  820,
			Content: ui.View(b.view),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
