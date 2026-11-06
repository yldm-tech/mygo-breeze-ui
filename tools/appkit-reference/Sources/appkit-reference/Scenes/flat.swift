import AppKit

// flatSet lays out the first row of TestScreenshotsFlat (flatView in ui/breeze/widgets_flat_test.go) as AppKit draws a button without a face at rest: NSButtons of bezel style accessoryBarAction (accessoryBar for the toggles, of type pushOnPushOff) that show their bezel only while the pointer is inside them (showsBorderOnlyWhileMouseInside), the image buttons in 24-point square frames. The pointer cannot be put inside a button offscreen, so the hovered one (trash, the second) is drawn with its bezel showing, as the pointer shows it; the row below holds the same buttons pressed (highlighted), which FlatButton draws only while the pointer presses one: an image button, a text button and the toggle on.
func flatSet() -> [Item] {
    func flat(_ symbol: String?, _ title: String, toggle: Bool = false, on: Bool = false, hover: Bool = false, pressed: Bool = false, enabled: Bool = true) -> NSButton {
        let b: NSButton
        if let symbol { b = NSButton(image: NSImage(systemSymbolName: symbol, accessibilityDescription: title)!, target: noop, action: #selector(Noop.act(_:))) } else { b = NSButton(title: title, target: noop, action: #selector(Noop.act(_:))) }
        b.bezelStyle = toggle ? .accessoryBar : .accessoryBarAction
        if toggle { b.setButtonType(.pushOnPushOff); b.state = on ? .on : .off }
        b.showsBorderOnlyWhileMouseInside = !(hover || pressed)
        b.isEnabled = enabled
        if pressed { b.highlight(true) }
        return b
    }
    let square = NSSize(width: 24, height: 24)
    return [
        Item(name: "flat-image", view: flat("square.and.arrow.up", "Share"), x: 16, y: 16, size: square),
        Item(name: "flat-image-hover", view: flat("trash", "Delete", hover: true), x: 64, y: 16, size: square),
        Item(name: "flat-image-disabled", view: flat("trash", "Delete", enabled: false), x: 112, y: 16, size: square),
        Item(name: "flat-text", view: flat(nil, "Skip"), x: 160, y: 16),
        Item(name: "flat-toggle-off", view: flat("sidebar.left", "Sidebar", toggle: true), x: 240, y: 16, size: square),
        Item(name: "flat-toggle-on", view: flat("sidebar.left", "Sidebar", toggle: true, on: true), x: 288, y: 16, size: square),
        Item(name: "flat-image-pressed", view: flat("square.and.arrow.up", "Share", pressed: true), x: 16, y: 60, size: square),
        Item(name: "flat-text-pressed", view: flat(nil, "Skip", pressed: true), x: 64, y: 60),
        Item(name: "flat-toggle-on-pressed", view: flat("sidebar.left", "Sidebar", toggle: true, on: true, pressed: true), x: 144, y: 60, size: square),
    ]
}

// FlatScenes renders the flat buttons of TestScreenshotsFlat (flatView in ui/breeze/widgets_flat_test.go).
final class FlatScenes: SceneGroup {
    override class var file: String { #fileID }
    override class func renderScenes(appearance: NSAppearance.Name, mode: String) {
        render(flatSet(), width: 360, height: 100, appearance: appearance, file: "appkit-flat-\(mode).png")
    }
}
