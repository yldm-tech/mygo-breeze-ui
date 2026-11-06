import AppKit

// textSet: AppKit's text styles (NSFont.preferredFont(forTextStyle:)) in labels, two in secondaryLabelColor, where TestScreenshotsText's textView places Breeze UI's Text (textTops).
func textSet() -> [Item] {
    func label(_ s: String, _ style: NSFont.TextStyle, secondary: Bool = false) -> NSTextField {
        let l = NSTextField(labelWithString: s)
        l.font = NSFont.preferredFont(forTextStyle: style, options: [:])
        if secondary { l.textColor = .secondaryLabelColor }
        return l
    }
    let rows: [(String, NSFont.TextStyle, Bool, CGFloat)] = [
        ("Large Title", .largeTitle, false, 16), ("Title", .title1, false, 76), ("Subtitle", .title2, false, 122), ("Heading", .headline, false, 160),
        ("Body", .body, false, 188), ("Caption", .caption1, false, 214), ("Secondary body", .body, true, 240), ("Secondary caption", .caption1, true, 266),
    ]
    return rows.map { Item(name: $0.0, view: label($0.0, $0.1, secondary: $0.2), x: 16, y: $0.3) }
}

// TextScenes renders the text styles of TestScreenshotsText (textView in ui/breeze/widgets_text_test.go).
final class TextScenes: SceneGroup {
    override class var file: String { #fileID }
    override class func renderScenes(appearance: NSAppearance.Name, mode: String) {
        render(textSet(), width: 360, height: 300, appearance: appearance, file: "appkit-text-\(mode).png")
    }
}
