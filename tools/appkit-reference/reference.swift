// AppKit's own controls, for comparing Breeze UI's macOS look with them.
//
// It renders two sets of controls in Aqua and Dark Aqua at 2x, laid out where ui/breeze/screenshot_test.go places Breeze UI's (referenceView and referenceView2), and the controls of TestScreenshotsInputs (inputsSet, with a window's toolbar holding a combo button and image items), writes them as PNGs and prints, for every control, its frame and the values along its middle row and column as runs (value x count): a gray pixel as its value, a colored one as #rrggbb, from which the sizes and fills in ui/breeze/mac.go were read. It draws into an sRGB bitmap, the color space MyGo writes, rather than the display's profile that bitmapImageRepForCachingDisplay picks, in which grays read the same but colors paler (systemGreen reads #65c466 in Display P3, #34c759 in sRGB).
//
// The window is offscreen and never key, so AppKit would draw it inactive, its accent controls (a checked box, an on switch, the default button, the chosen segment) gray. Its NSWindow subclass answers hasKeyAppearance (an NSWindow method outside the SDK's headers) with true, so they are drawn as in the active window, as Breeze UI draws them; -inactive leaves the window inactive. Materials (Liquid Glass, vibrancy), menus, popovers, tooltips and focus rings still do not render into the bitmap.
//
//	swift tools/appkit-reference/reference.swift [-inactive] <dir>
import AppKit

let app = NSApplication.shared
app.setActivationPolicy(.accessory)
let args = CommandLine.arguments.dropFirst()
let inactive = args.contains("-inactive")
let out = args.last(where: { $0 != "-inactive" }) ?? "."

// A window that draws as the key window, unless -inactive.
class KeyWindow: NSWindow {
    @objc func hasKeyAppearance() -> Bool { !inactive }
}

struct Item {
    let name: String
    let view: NSView
    let x: CGFloat // left edge, points
    let y: CGFloat // top edge, points from the top
    var size: NSSize? = nil
}

func runs(_ v: [String]) -> String {
    var s: [String] = []
    var i = 0
    while i < v.count {
        var j = i
        while j < v.count && v[j] == v[i] { j += 1 }
        s.append(j - i > 1 ? "\(v[i])x\(j - i)" : "\(v[i])")
        i = j
    }
    return s.joined(separator: " ")
}

func firstSet() -> [Item] {
    let cancel = NSButton(title: "Cancel", target: nil, action: nil); cancel.bezelStyle = .push
    let save = NSButton(title: "Save", target: nil, action: nil); save.bezelStyle = .push; save.keyEquivalent = "\r"
    let disabled = NSButton(title: "Disabled", target: nil, action: nil); disabled.bezelStyle = .push; disabled.isEnabled = false
    let c1 = NSButton(checkboxWithTitle: "Unchecked", target: nil, action: nil)
    let c2 = NSButton(checkboxWithTitle: "Checked", target: nil, action: nil); c2.state = .on
    let c3 = NSButton(checkboxWithTitle: "Disabled", target: nil, action: nil); c3.state = .on; c3.isEnabled = false
    let s1 = NSSwitch(), s2 = NSSwitch(), s3 = NSSwitch(); s2.state = .on; s3.isEnabled = false
    return [
        Item(name: "button", view: cancel, x: 16, y: 16), Item(name: "default-button", view: save, x: 110, y: 16), Item(name: "disabled-button", view: disabled, x: 200, y: 16),
        Item(name: "checkbox-off", view: c1, x: 16, y: 64), Item(name: "checkbox-on", view: c2, x: 130, y: 64), Item(name: "checkbox-disabled", view: c3, x: 240, y: 64),
        Item(name: "switch-off", view: s1, x: 16, y: 106), Item(name: "switch-on", view: s2, x: 100, y: 106), Item(name: "switch-disabled", view: s3, x: 184, y: 106),
    ]
}

func secondSet() -> [Item] {
    let r1 = NSButton(radioButtonWithTitle: "Off", target: nil, action: nil)
    let r2 = NSButton(radioButtonWithTitle: "On", target: nil, action: nil); r2.state = .on
    let field = NSTextField(string: ""); field.placeholderString = "Name"; field.isBezeled = true; field.bezelStyle = .squareBezel
    let search = NSSearchField(); search.placeholderString = "Search"
    let popup = NSPopUpButton(frame: .zero, pullsDown: false); popup.addItems(withTitles: ["Medium", "Small", "Large"])
    let toggle = NSButton(title: "Bold", target: nil, action: nil); toggle.bezelStyle = .push; toggle.setButtonType(.pushOnPushOff); toggle.state = .on
    let slider = NSSlider(value: 40, minValue: 0, maxValue: 100, target: nil, action: nil)
    let bar = NSProgressIndicator(); bar.style = .bar; bar.isIndeterminate = false; bar.minValue = 0; bar.maxValue = 100; bar.doubleValue = 40
    let spinner = NSProgressIndicator(); spinner.style = .spinning; spinner.isIndeterminate = true
    let segmented = NSSegmentedControl(labels: ["Day", "Week", "Month"], trackingMode: .selectOne, target: nil, action: nil); segmented.selectedSegment = 0
    let text = NSTextView(frame: NSRect(x: 0, y: 0, width: 160, height: 60)); text.string = "Notes"
    let scroll = NSScrollView(); scroll.borderType = .bezelBorder; scroll.documentView = text
    return [
        Item(name: "radio-off", view: r1, x: 16, y: 16), Item(name: "radio-on", view: r2, x: 100, y: 16),
        Item(name: "text-field", view: field, x: 16, y: 50, size: NSSize(width: 180, height: field.fittingSize.height)),
        Item(name: "search-field", view: search, x: 220, y: 50, size: NSSize(width: 180, height: search.fittingSize.height)),
        Item(name: "pop-up-button", view: popup, x: 16, y: 90), Item(name: "toggle", view: toggle, x: 220, y: 90),
        Item(name: "slider", view: slider, x: 16, y: 130, size: NSSize(width: 180, height: slider.fittingSize.height)),
        Item(name: "progress", view: bar, x: 220, y: 136, size: NSSize(width: 180, height: bar.fittingSize.height)),
        Item(name: "segmented", view: segmented, x: 16, y: 170),
        Item(name: "spinner", view: spinner, x: 220, y: 170, size: NSSize(width: 32, height: 32)),
        Item(name: "text-view", view: scroll, x: 16, y: 220, size: NSSize(width: 180, height: 70)),
    ]
}

// Noop is the target of controls that draw disabled without an action.
class Noop: NSObject {
    @objc func act(_ sender: Any?) {}
}

let noop = Noop()

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

func render(_ items: [Item], width W: CGFloat, height H: CGFloat, appearance: NSAppearance.Name, file: String) {
    let root = NSView(frame: NSRect(x: 0, y: 0, width: W, height: H))
    root.appearance = NSAppearance(named: appearance)
    let bg = NSBox(frame: root.bounds); bg.boxType = .custom; bg.fillColor = .windowBackgroundColor; bg.borderWidth = 0
    root.addSubview(bg)
    for it in items {
        let s = it.size ?? it.view.fittingSize
        it.view.frame = NSRect(x: it.x, y: H - it.y - s.height, width: s.width, height: s.height)
        root.addSubview(it.view)
    }
    let win = KeyWindow(contentRect: NSRect(x: -4000, y: -4000, width: W, height: H), styleMask: [.borderless], backing: .buffered, defer: false)
    win.contentView = root
    root.layoutSubtreeIfNeeded()
    root.display()
    // An sRGB bitmap at 2x: cacheDisplay draws into the bitmap's own color space.
    let rep = NSBitmapImageRep(bitmapDataPlanes: nil, pixelsWide: Int(W * 2), pixelsHigh: Int(H * 2), bitsPerSample: 8, samplesPerPixel: 4, hasAlpha: true, isPlanar: false, colorSpaceName: .calibratedRGB, bytesPerRow: 0, bitsPerPixel: 0)!.retagging(with: .sRGB)!
    rep.size = root.bounds.size
    root.cacheDisplay(in: root.bounds, to: rep)
    try! rep.representation(using: .png, properties: [:])!.write(to: URL(fileURLWithPath: out + "/" + file))
    func gray(_ x: Int, _ y: Int) -> String {
        var p = [Int](repeating: 0, count: 4); rep.getPixel(&p, atX: x, y: y)
        return p[0] == p[1] && p[1] == p[2] ? "\(p[0])" : String(format: "#%02x%02x%02x", p[0], p[1], p[2])
    }
    print("== \(file)")
    for it in items {
        let f = it.view.frame
        let x0 = Int(f.minX * 2), y0 = Int((H - f.maxY) * 2), w = Int(f.width * 2), h = Int(f.height * 2)
        print("\(it.name) x=\(f.minX) top=\(H - f.maxY) w=\(f.width) h=\(f.height)")
        print("  row: " + runs((max(0, x0 - 4)..<min(Int(W * 2), x0 + w + 4)).map { gray($0, y0 + h / 2) }))
        print("  col: " + runs((max(0, y0 - 4)..<min(Int(H * 2), y0 + h + 4)).map { gray(x0 + w / 2, $0) }))
    }
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

// renderToolbar renders a titled window's unified toolbar holding ToolbarItems' items, the window's frame view at 2x into an sRGB bitmap, and prints the frames of the items' views and platters. It renders offscreen with the items' platters, 36 points high as a live capture shows them though lighter: the sizes inside them compare, and their fills as shares of the platter.
func renderToolbar(appearance: NSAppearance.Name, file: String) {
    let win = KeyWindow(contentRect: NSRect(x: -4000, y: -4000, width: 520, height: 120), styleMask: [.titled, .closable, .fullSizeContentView], backing: .buffered, defer: false)
    win.appearance = NSAppearance(named: appearance)
    win.titleVisibility = .hidden
    let bar = NSToolbar(identifier: "inputs"); bar.delegate = toolbarItems; bar.displayMode = .iconOnly
    win.toolbar = bar; win.toolbarStyle = .unified
    win.orderFront(nil)
    let frame = win.contentView!.superview!
    frame.layoutSubtreeIfNeeded(); frame.display()
    let W = frame.bounds.width, H = frame.bounds.height
    let rep = NSBitmapImageRep(bitmapDataPlanes: nil, pixelsWide: Int(W * 2), pixelsHigh: Int(H * 2), bitsPerSample: 8, samplesPerPixel: 4, hasAlpha: true, isPlanar: false, colorSpaceName: .calibratedRGB, bytesPerRow: 0, bitsPerPixel: 0)!.retagging(with: .sRGB)!
    rep.size = frame.bounds.size
    frame.cacheDisplay(in: frame.bounds, to: rep)
    try! rep.representation(using: .png, properties: [:])!.write(to: URL(fileURLWithPath: out + "/" + file))
    print("== \(file)")
    func walk(_ v: NSView) {
        let name = String(describing: type(of: v))
        if name.contains("Platter") || v is NSComboButton || name == "NSToolbarButton" {
            let f = v.convert(v.bounds, to: frame)
            print("\(name) x=\(f.minX) top=\(H - f.maxY) w=\(f.width) h=\(f.height)")
        }
        v.subviews.forEach(walk)
    }
    walk(frame)
    win.orderOut(nil)
}

for (appearance, mode) in [(NSAppearance.Name.aqua, "light"), (.darkAqua, "dark")] {
    render(firstSet(), width: 360, height: 200, appearance: appearance, file: "appkit-\(mode).png")
    render(secondSet(), width: 520, height: 330, appearance: appearance, file: "appkit2-\(mode).png")
    render(inputsSet(), width: 520, height: 180, appearance: appearance, file: "appkit-inputs-\(mode).png")
    renderToolbar(appearance: appearance, file: "appkit-inputs-toolbar-\(mode).png")
}
