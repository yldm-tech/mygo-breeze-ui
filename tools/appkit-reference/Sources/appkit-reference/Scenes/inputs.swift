import AppKit

// inputsSet lays out the controls of TestScreenshotsInputs (inputsView in ui/breeze/widgets_inputs_test.go): two NSSecureTextFields, empty with a placeholder and holding a password, an NSComboButton of style split, and bordered push buttons showing only an image, the second disabled.
func inputsSet() -> [Item] {
    let empty = NSSecureTextField(string: ""); empty.placeholderString = "Password"
    let filled = NSSecureTextField(string: "hunter2")
    let menu = NSMenu(); menu.addItem(withTitle: "Save As…", action: nil, keyEquivalent: ""); menu.addItem(withTitle: "Save All", action: nil, keyEquivalent: "")
    // An NSComboButton without an action draws its title disabled.
    let combo = NSComboButton(title: "Save", menu: menu, target: noop, action: #selector(Noop.act(_:))); combo.style = .split
    let share = NSButton(image: NSImage(systemSymbolName: "square.and.arrow.up", accessibilityDescription: "Share")!, target: nil, action: nil); share.bezelStyle = .push
    let trash = NSButton(image: NSImage(systemSymbolName: "trash", accessibilityDescription: "Delete")!, target: nil, action: nil); trash.bezelStyle = .push; trash.isEnabled = false
    return [
        Item(name: "secure-field-empty", view: empty, x: 16, y: 16, size: NSSize(width: 180, height: empty.fittingSize.height)),
        Item(name: "secure-field", view: filled, x: 220, y: 16, size: NSSize(width: 180, height: filled.fittingSize.height)),
        Item(name: "combo-button", view: combo, x: 16, y: 64), Item(name: "image-button", view: share, x: 160, y: 64), Item(name: "image-button-disabled", view: trash, x: 220, y: 64),
    ]
}

// ToolbarItems fills the toolbar of renderToolbar: an NSComboButton as an item's view, then two bordered image items, which share one platter.
class ToolbarItems: NSObject, NSToolbarDelegate {
    func toolbarDefaultItemIdentifiers(_ t: NSToolbar) -> [NSToolbarItem.Identifier] { [.init("combo"), .init("share"), .init("trash")] }
    func toolbarAllowedItemIdentifiers(_ t: NSToolbar) -> [NSToolbarItem.Identifier] { toolbarDefaultItemIdentifiers(t) }
    func toolbar(_ t: NSToolbar, itemForItemIdentifier id: NSToolbarItem.Identifier, willBeInsertedIntoToolbar f: Bool) -> NSToolbarItem? {
        let it = NSToolbarItem(itemIdentifier: id)
        switch id.rawValue {
        case "combo":
            let menu = NSMenu(); menu.addItem(withTitle: "Save As…", action: nil, keyEquivalent: "")
            it.view = NSComboButton(title: "New", menu: menu, target: noop, action: #selector(Noop.act(_:))); it.label = "New"
        case "share":
            it.image = NSImage(systemSymbolName: "square.and.arrow.up", accessibilityDescription: "Share"); it.label = "Share"
        default:
            it.image = NSImage(systemSymbolName: "trash", accessibilityDescription: "Delete"); it.label = "Delete"
        }
        if it.view == nil { it.target = noop; it.action = #selector(Noop.act(_:)); it.isBordered = true }
        return it
    }
}

let toolbarItems = ToolbarItems()

// InputsScenes renders the controls of TestScreenshotsInputs (inputsView in ui/breeze/widgets_inputs_test.go), alone and in a window's toolbar.
final class InputsScenes: SceneGroup {
    override class var file: String { #fileID }
    override class func renderScenes(appearance: NSAppearance.Name, mode: String) {
        render(inputsSet(), width: 520, height: 180, appearance: appearance, file: "appkit-inputs-\(mode).png")
        renderToolbar(appearance: appearance, file: "appkit-inputs-toolbar-\(mode).png")
    }
}
