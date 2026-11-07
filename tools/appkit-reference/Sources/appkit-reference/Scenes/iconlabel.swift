import AppKit

// templateImage is a 16-point template image stroking the path draw builds, in the coordinates of an SVG's 16 by 16 viewBox, 1.3 points wide, as the icons of the Go tests draw (idiomTestShareIcon and inputsTestTrashIcon).
func templateImage(_ draw: @escaping (NSBezierPath) -> Void) -> NSImage {
    let img = NSImage(size: NSSize(width: 16, height: 16), flipped: true) { _ in
        let path = NSBezierPath(); path.lineWidth = 1.3; path.lineJoinStyle = .round
        draw(path)
        NSColor.black.setStroke(); path.stroke()
        return true
    }
    img.isTemplate = true
    return img
}

let shareImage = templateImage { p in
    p.move(to: NSPoint(x: 8, y: 1.5)); p.line(to: NSPoint(x: 8, y: 10.5))
    p.move(to: NSPoint(x: 4.8, y: 4.5)); p.line(to: NSPoint(x: 8, y: 1.5)); p.line(to: NSPoint(x: 11.2, y: 4.5))
    p.move(to: NSPoint(x: 5, y: 7)); p.line(to: NSPoint(x: 3.5, y: 7)); p.line(to: NSPoint(x: 3.5, y: 14.5)); p.line(to: NSPoint(x: 12.5, y: 14.5)); p.line(to: NSPoint(x: 12.5, y: 7)); p.line(to: NSPoint(x: 11, y: 7))
}

let trashImage = templateImage { p in
    p.move(to: NSPoint(x: 2.5, y: 4)); p.line(to: NSPoint(x: 13.5, y: 4))
    p.move(to: NSPoint(x: 6, y: 4)); p.line(to: NSPoint(x: 6, y: 2.5)); p.line(to: NSPoint(x: 10, y: 2.5)); p.line(to: NSPoint(x: 10, y: 4))
    p.move(to: NSPoint(x: 4, y: 4)); p.line(to: NSPoint(x: 4.8, y: 14)); p.line(to: NSPoint(x: 11.2, y: 14)); p.line(to: NSPoint(x: 12, y: 4))
    p.move(to: NSPoint(x: 6.5, y: 6.5)); p.line(to: NSPoint(x: 6.5, y: 11.5))
    p.move(to: NSPoint(x: 9.5, y: 6.5)); p.line(to: NSPoint(x: 9.5, y: 11.5))
}

// iconLabelSet lays out the push buttons of TestScreenshotsIconLabel (iconLabelView in ui/breeze/widgets_iconlabel_test.go): NSButtons of bezel style push showing a template image leading their titles (imagePosition .imageLeading), the second the default button and the third disabled.
func iconLabelSet() -> [Item] {
    func button(_ title: String, _ image: NSImage) -> NSButton {
        let b = NSButton(title: title, image: image, target: noop, action: #selector(Noop.act(_:))); b.bezelStyle = .push; b.imagePosition = .imageLeading
        return b
    }
    let share = button("Share", shareImage)
    let export = button("Export", shareImage); export.keyEquivalent = "\r"
    let delete = button("Delete", trashImage); delete.isEnabled = false
    return [Item(name: "icon-label-button", view: share, x: 16, y: 16), Item(name: "icon-label-default", view: export, x: 130, y: 16), Item(name: "icon-label-disabled", view: delete, x: 250, y: 16)]
}

// IconLabelItems fills the toolbar of the icon-label scene: bordered items showing an image and a title (an NSToolbarItem with both, which the iconOnly display mode shows side by side), an image item between them and the last prominent, as iconLabelView's Toolbar holds them.
class IconLabelItems: NSObject, NSToolbarDelegate {
    func toolbarDefaultItemIdentifiers(_ t: NSToolbar) -> [NSToolbarItem.Identifier] { [.init("share"), .init("delete"), .init("export")] }
    func toolbarAllowedItemIdentifiers(_ t: NSToolbar) -> [NSToolbarItem.Identifier] { toolbarDefaultItemIdentifiers(t) }
    func toolbar(_ t: NSToolbar, itemForItemIdentifier id: NSToolbarItem.Identifier, willBeInsertedIntoToolbar f: Bool) -> NSToolbarItem? {
        let it = NSToolbarItem(itemIdentifier: id)
        switch id.rawValue {
        case "share":
            it.image = shareImage; it.title = "Share"; it.label = "Share"
        case "delete":
            it.image = trashImage; it.label = "Delete"
        default:
            it.image = shareImage; it.title = "Export"; it.label = "Export"; it.style = .prominent
        }
        it.target = noop; it.action = #selector(Noop.act(_:)); it.isBordered = true
        return it
    }
}

let iconLabelItems = IconLabelItems()

// IconLabelScenes renders the buttons of TestScreenshotsIconLabel (iconLabelView in ui/breeze/widgets_iconlabel_test.go), alone and in a window's toolbar.
final class IconLabelScenes: SceneGroup {
    override class var file: String { #fileID }
    override class func renderScenes(appearance: NSAppearance.Name, mode: String) {
        render(iconLabelSet(), width: 520, height: 64, appearance: appearance, file: "appkit-iconlabel-\(mode).png")
        renderToolbar(appearance: appearance, file: "appkit-iconlabel-toolbar-\(mode).png", items: iconLabelItems, identifier: "iconlabel")
    }
}
