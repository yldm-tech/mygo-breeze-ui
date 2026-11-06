import AppKit

// HeaderItems fills the toolbar of renderHeaderBar as an app's window title bar: a navigational back item before the title, then after a flexible space a bordered text item and a bordered image item, as breeze.HeaderBar's start and end hold them in TestScreenshotsHeaderBar.
class HeaderItems: NSObject, NSToolbarDelegate {
    var back = true
    func toolbarDefaultItemIdentifiers(_ t: NSToolbar) -> [NSToolbarItem.Identifier] { (back ? [.init("back")] : []) + [.flexibleSpace, .init("edit"), .init("share")] }
    func toolbarAllowedItemIdentifiers(_ t: NSToolbar) -> [NSToolbarItem.Identifier] { toolbarDefaultItemIdentifiers(t) }
    func toolbar(_ t: NSToolbar, itemForItemIdentifier id: NSToolbarItem.Identifier, willBeInsertedIntoToolbar f: Bool) -> NSToolbarItem? {
        let it = NSToolbarItem(itemIdentifier: id)
        switch id.rawValue {
        case "back":
            it.image = NSImage(systemSymbolName: "chevron.left", accessibilityDescription: "Back"); it.label = "Back"; it.isNavigational = true
        case "edit":
            it.title = "Edit"; it.label = "Edit"
        default:
            it.image = NSImage(systemSymbolName: "square.and.arrow.up", accessibilityDescription: "Share"); it.label = "Share"
        }
        it.target = noop; it.action = #selector(Noop.act(_:)); it.isBordered = true
        return it
    }
}

let headerItems = HeaderItems()
let headerItemsNoBack: HeaderItems = { let h = HeaderItems(); h.back = false; return h }()

// renderHeaderBar renders a titled window with a unified toolbar showing its title, and its subtitle when given, with or without the back item, the window's frame view at 2x into an sRGB bitmap, as breeze.HeaderBar draws the bar on macOS, and prints the frames of the traffic lights, the title's text fields with their fonts and colors, the toolbar's view and the items' platters.
func renderHeaderBar(appearance: NSAppearance.Name, subtitle: String, back: Bool = true, file: String) {
    let win = KeyWindow(contentRect: NSRect(x: -4000, y: -4000, width: 560, height: 120), styleMask: [.titled, .closable, .miniaturizable, .resizable, .fullSizeContentView], backing: .buffered, defer: false)
    win.appearance = NSAppearance(named: appearance)
    win.title = "Documents"; win.subtitle = subtitle
    let bar = NSToolbar(identifier: back ? "header" : "header-noback"); bar.delegate = back ? headerItems : headerItemsNoBack; bar.displayMode = .iconOnly
    win.toolbar = bar; win.toolbarStyle = .unified
    win.orderFront(nil)
    let frame = win.contentView!.superview!
    frame.layoutSubtreeIfNeeded(); frame.display()
    let W = frame.bounds.width, H = frame.bounds.height
    let rep = NSBitmapImageRep(bitmapDataPlanes: nil, pixelsWide: Int(W * 2), pixelsHigh: Int(H * 2), bitsPerSample: 8, samplesPerPixel: 4, hasAlpha: true, isPlanar: false, colorSpaceName: .calibratedRGB, bytesPerRow: 0, bitsPerPixel: 0)!.retagging(with: .sRGB)!
    rep.size = frame.bounds.size
    frame.cacheDisplay(in: frame.bounds, to: rep)
    try! rep.representation(using: .png, properties: [:])!.write(to: URL(fileURLWithPath: out + "/" + file))
    print("== \(file) window \(W)x\(H) contentLayoutRect top=\(H - win.contentLayoutRect.maxY)")
    for b in [NSWindow.ButtonType.closeButton, .miniaturizeButton, .zoomButton] {
        if let v = win.standardWindowButton(b) {
            let f = v.convert(v.bounds, to: frame)
            print("button \(b.rawValue) x=\(f.minX) top=\(H - f.maxY) w=\(f.width) h=\(f.height)")
        }
    }
    func walk(_ v: NSView) {
        let name = String(describing: type(of: v))
        let f = v.convert(v.bounds, to: frame)
        if let t = v as? NSTextField, !t.stringValue.isEmpty {
            let font = t.font!
            let color = t.textColor?.usingColorSpace(.sRGB)
            let size = t.attributedStringValue.size()
            print("text \(name) \"\(t.stringValue)\" x=\(f.minX) top=\(H - f.maxY) w=\(f.width) h=\(f.height) advance=\(size.width) font=\(font.fontName) \(font.pointSize) color=\(t.textColor?.description ?? "nil") srgb=\(color.map { "\($0.redComponent) \($0.greenComponent) \($0.blueComponent) \($0.alphaComponent)" } ?? "-")")
        } else if name.contains("Platter") || name == "NSToolbarButton" || name == "NSToolbarView" || name.contains("TitlebarView") || name.contains("ToolbarTitle") {
            print("\(name) x=\(f.minX) top=\(H - f.maxY) w=\(f.width) h=\(f.height)")
        }
        v.subviews.forEach(walk)
    }
    walk(frame)
    func gray(_ x: Int, _ y: Int) -> String {
        var p = [Int](repeating: 0, count: 4); rep.getPixel(&p, atX: x, y: y)
        return p[0] == p[1] && p[1] == p[2] ? "\(p[0])" : String(format: "#%02x%02x%02x", p[0], p[1], p[2])
    }
    // The bar's background down a column clear of the items, and across it under the title's middle.
    print("  col x=300: " + runs((0..<Int(H * 2)).map { gray(600, $0) }))
    print("  row y=26: " + runs((0..<Int(W * 2)).map { gray($0, 52) }))
    win.orderOut(nil)
}

// SplitHeaderItems fills the toolbar of renderHeaderBarSplit: a tracking separator over the split view's divider, then the items of HeaderItems after it, so that they and the title sit over the content pane, as breeze.HeaderBar's start and end hold them in TestScreenshotsHeaderBarInPane.
class SplitHeaderItems: HeaderItems {
    override func toolbarDefaultItemIdentifiers(_ t: NSToolbar) -> [NSToolbarItem.Identifier] { [.sidebarTrackingSeparator] + super.toolbarDefaultItemIdentifiers(t) }
}

let splitHeaderItems = SplitHeaderItems()
let splitHeaderItemsNoBack: SplitHeaderItems = { let h = SplitHeaderItems(); h.back = false; return h }()

// renderHeaderBarSplit renders a titled window with a unified toolbar over an NSSplitViewController, a 236-point sidebar beside the content, the toolbar's items and the title after a tracking separator over the divider, with or without the back item and the subtitle, the window's frame view at 2x into an sRGB bitmap, as breeze.HeaderBar draws a bar over a window's content pane on macOS, and prints the frames of the traffic lights, the divider, the title's text fields, the toolbar's view and the items' platters.
func renderHeaderBarSplit(appearance: NSAppearance.Name, subtitle: String, back: Bool = true, file: String) {
    let win = KeyWindow(contentRect: NSRect(x: -4000, y: -4000, width: 720, height: 120), styleMask: [.titled, .closable, .miniaturizable, .resizable, .fullSizeContentView], backing: .buffered, defer: false)
    win.appearance = NSAppearance(named: appearance)
    win.title = "Documents"; win.subtitle = subtitle
    let split = NSSplitViewController()
    let side = NSSplitViewItem(sidebarWithViewController: NSViewController()); side.viewController.view = NSView()
    side.minimumThickness = 236; side.maximumThickness = 236
    let content = NSSplitViewItem(viewController: NSViewController()); content.viewController.view = NSView()
    split.addSplitViewItem(side); split.addSplitViewItem(content)
    win.contentViewController = split
    win.setContentSize(NSSize(width: 720, height: 120))
    let bar = NSToolbar(identifier: back ? "split-header" : "split-header-noback"); bar.delegate = back ? splitHeaderItems : splitHeaderItemsNoBack; bar.displayMode = .iconOnly
    win.toolbar = bar; win.toolbarStyle = .unified
    win.orderFront(nil)
    split.splitView.setPosition(236, ofDividerAt: 0)
    let frame = win.contentView!.superview!
    frame.layoutSubtreeIfNeeded(); frame.display()
    let W = frame.bounds.width, H = frame.bounds.height
    let rep = NSBitmapImageRep(bitmapDataPlanes: nil, pixelsWide: Int(W * 2), pixelsHigh: Int(H * 2), bitsPerSample: 8, samplesPerPixel: 4, hasAlpha: true, isPlanar: false, colorSpaceName: .calibratedRGB, bytesPerRow: 0, bitsPerPixel: 0)!.retagging(with: .sRGB)!
    rep.size = frame.bounds.size
    frame.cacheDisplay(in: frame.bounds, to: rep)
    try! rep.representation(using: .png, properties: [:])!.write(to: URL(fileURLWithPath: out + "/" + file))
    let sideFrame = side.viewController.view.convert(side.viewController.view.bounds, to: frame)
    print("== \(file) window \(W)x\(H) sidebar x=\(sideFrame.minX) w=\(sideFrame.width)")
    for b in [NSWindow.ButtonType.closeButton, .miniaturizeButton, .zoomButton] {
        if let v = win.standardWindowButton(b) {
            let f = v.convert(v.bounds, to: frame)
            print("button \(b.rawValue) x=\(f.minX) top=\(H - f.maxY) w=\(f.width) h=\(f.height)")
        }
    }
    func walk(_ v: NSView) {
        let name = String(describing: type(of: v))
        let f = v.convert(v.bounds, to: frame)
        if let t = v as? NSTextField, !t.stringValue.isEmpty {
            let size = t.attributedStringValue.size()
            print("text \(name) \"\(t.stringValue)\" x=\(f.minX) top=\(H - f.maxY) w=\(f.width) h=\(f.height) advance=\(size.width) font=\(t.font!.fontName) \(t.font!.pointSize)")
        } else if name.contains("Platter") || name == "NSToolbarButton" || name == "NSToolbarView" || name.contains("TitlebarView") || name.contains("ToolbarTitle") || name.contains("Separator") {
            print("\(name) x=\(f.minX) top=\(H - f.maxY) w=\(f.width) h=\(f.height)")
        }
        v.subviews.forEach(walk)
    }
    walk(frame)
    func gray(_ x: Int, _ y: Int) -> String {
        var p = [Int](repeating: 0, count: 4); rep.getPixel(&p, atX: x, y: y)
        return p[0] == p[1] && p[1] == p[2] ? "\(p[0])" : String(format: "#%02x%02x%02x", p[0], p[1], p[2])
    }
    // The content pane's bar down a column clear of the items, and across the window under the title's middle.
    print("  col x=400: " + runs((0..<Int(H * 2)).map { gray(800, $0) }))
    print("  row y=26: " + runs((0..<Int(W * 2)).map { gray($0, 52) }))
    win.orderOut(nil)
}

// HeaderBarScenes renders the title bars of TestScreenshotsHeaderBar and TestScreenshotsHeaderBarInPane (ui/breeze/widgets_headerbar_test.go).
final class HeaderBarScenes: SceneGroup {
    override class var file: String { #fileID }
    override class func renderScenes(appearance: NSAppearance.Name, mode: String) {
        renderHeaderBar(appearance: appearance, subtitle: "", file: "appkit-headerbar-\(mode).png")
        renderHeaderBar(appearance: appearance, subtitle: "12 items", file: "appkit-headerbar-subtitle-\(mode).png")
        renderHeaderBar(appearance: appearance, subtitle: "", back: false, file: "appkit-headerbar-noback-\(mode).png")
        renderHeaderBarSplit(appearance: appearance, subtitle: "", file: "appkit-headerbar-split-\(mode).png")
        renderHeaderBarSplit(appearance: appearance, subtitle: "12 items", file: "appkit-headerbar-split-subtitle-\(mode).png")
        renderHeaderBarSplit(appearance: appearance, subtitle: "", back: false, file: "appkit-headerbar-split-noback-\(mode).png")
    }
}
