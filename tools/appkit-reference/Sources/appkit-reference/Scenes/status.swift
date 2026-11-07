import AppKit
import SwiftUI

// StatusBar is the closest macOS has to an inline message bar, which breeze.InfoBar draws there: a row of a grouped Form holding a Label of an SF Symbol in the severity's color over a bold title and a secondary message, and a borderless button showing xmark. The Form insets its section 20 points, which statusSet takes back by placing the view 20 points up and to the left.
struct StatusBar: View {
    let symbol: String, color: Color, title: String, closable: Bool
    var message = "Some message about it.", action = "", prominent = false
    var body: some View {
        Form {
            Section {
                HStack(alignment: .firstTextBaseline) {
                    Label { VStack(alignment: .leading, spacing: 2) { Text(title).bold(); if !message.isEmpty { Text(message).foregroundStyle(.secondary) } } } icon: { Image(systemName: symbol).foregroundStyle(color) }
                    Spacer()
                    if !action.isEmpty {
                        if prominent { Button(action) {}.buttonStyle(.borderedProminent) } else { Button(action) {} }
                    }
                    if closable { Button {} label: { Image(systemName: "xmark") }.buttonStyle(.borderless) }
                }
            }
        }.formStyle(.grouped).scrollDisabled(true).scrollContentBackground(.hidden)
    }
}

// statusSet is TestScreenshotsStatus's statusView: four message bars, one of each severity, the first and the last closable, 16 from the window's edges and 68 apart, and determinate spinning progress indicators of 0.25, 0.6, 1 and 0 at the regular size and one of 0.6 at the small size below them.
func statusSet() -> [Item] {
    var items: [Item] = []
    let bars: [(String, Color, String)] = [("info.circle.fill", .accentColor, "Informational"), ("checkmark.circle.fill", .green, "Success"), ("exclamationmark.triangle.fill", .yellow, "Warning"), ("xmark.octagon.fill", .red, "Error")]
    for (i, b) in bars.enumerated() {
        let view = NSHostingView(rootView: StatusBar(symbol: b.0, color: b.1, title: b.2, closable: i == 0 || i == 3))
        items.append(Item(name: "bar-\(b.2.lowercased())", view: view, x: -4, y: 16 + CGFloat(i) * 68 - 20, size: NSSize(width: 568, height: 94)))
    }
    for (i, v) in [0.25, 0.6, 1, 0].enumerated() {
        let ring = NSProgressIndicator(); ring.style = .spinning; ring.isIndeterminate = false; ring.minValue = 0; ring.maxValue = 1; ring.doubleValue = v
        items.append(Item(name: "ring-\(v)", view: ring, x: 16 + CGFloat(i) * 48, y: 300, size: NSSize(width: 32, height: 32)))
    }
    let small = NSProgressIndicator(); small.style = .spinning; small.isIndeterminate = false; small.controlSize = .small; small.minValue = 0; small.maxValue = 1; small.doubleValue = 0.6
    items.append(Item(name: "ring-small", view: small, x: 208, y: 308, size: NSSize(width: 16, height: 16)))
    return items
}

// statusActionsSet is TestScreenshotsStatus's statusActionsView: message bars with actions, InfoBarActions' buttons, 16 from the window's edges and 68 apart: a closable warning with a push button, an informational bar with a default (prominent) button, a closable success bar without a message with a push button, and a closable error bar 280 points wide with a push button.
func statusActionsSet() -> [Item] {
    let bars: [(String, Color, String, String, String, Bool, Bool, CGFloat)] = [
        ("exclamationmark.triangle.fill", .yellow, "Warning", "Some message about it.", "Open Settings", true, false, 528),
        ("info.circle.fill", .accentColor, "Informational", "Some message about it.", "Restart", false, true, 528),
        ("checkmark.circle.fill", .green, "Success", "", "Undo", true, false, 528),
        ("xmark.octagon.fill", .red, "Error", "Some message about it.", "Retry", true, false, 280),
    ]
    return bars.enumerated().map { i, b in
        let view = NSHostingView(rootView: StatusBar(symbol: b.0, color: b.1, title: b.2, closable: b.5, message: b.3, action: b.4, prominent: b.6))
        return Item(name: "actions-\(b.2.lowercased())", view: view, x: -4, y: 16 + CGFloat(i) * 68 - 20, size: NSSize(width: b.7 + 40, height: i == 3 ? 140 : 94))
    }
}

// statusPageSet is TestScreenshotsStatus's statusPageView: SwiftUI's ContentUnavailableView filling the window, with a magnifying glass, a title, a description and a button.
func statusPageSet() -> [Item] {
    let page = NSHostingView(rootView: ContentUnavailableView {
        Label("No Results", systemImage: "magnifyingglass")
    } description: {
        Text("Check the spelling or try a new search.")
    } actions: {
        Button("Clear Search") {}
    })
    return [Item(name: "content-unavailable", view: page, x: 0, y: 0, size: NSSize(width: 560, height: 420))]
}

// StatusScenes renders the status widgets of TestScreenshotsStatus (statusView, statusPageView and statusActionsView in ui/breeze/widgets_status_test.go).
final class StatusScenes: SceneGroup {
    override class var file: String { #fileID }
    override class func renderScenes(appearance: NSAppearance.Name, mode: String) {
        render(statusSet(), width: 560, height: 360, appearance: appearance, file: "appkit-status-\(mode).png")
        render(statusPageSet(), width: 560, height: 420, appearance: appearance, file: "appkit-statuspage-\(mode).png")
        render(statusActionsSet(), width: 560, height: 380, appearance: appearance, file: "appkit-statusactions-\(mode).png")
    }
}
