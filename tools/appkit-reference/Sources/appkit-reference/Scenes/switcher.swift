import AppKit

// switcherIcons are the icons of TestScreenshotsViewSwitcher (switcherTestIcons in ui/breeze/widgets_switcher_test.go), the same SVGs, which NSImage reads, as template images 16 points square.
let switcherIcons: [NSImage] = [
    ##"<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><path d="M2.5 2.5h3v11h-3zM6.5 2.5h3v11h-3zM10.4 3.4l2.4-.6 2.4 10.2-2.4.6z" fill="none" stroke="#000" stroke-width="1.3" stroke-linejoin="round"/></svg>"##,
    ##"<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><path d="M2.5 2.5h4v4h-4zM9.5 2.5h4v4h-4zM2.5 9.5h4v4h-4zM9.5 9.5h4v4h-4z" fill="none" stroke="#000" stroke-width="1.3" stroke-linejoin="round"/></svg>"##,
    ##"<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><circle cx="8" cy="8" r="1.6"/><path d="M5.3 5.3a3.8 3.8 0 0 0 0 5.4M10.7 5.3a3.8 3.8 0 0 1 0 5.4M3.2 3.2a6.8 6.8 0 0 0 0 9.6M12.8 3.2a6.8 6.8 0 0 1 0 9.6" fill="none" stroke="#000" stroke-width="1.3" stroke-linecap="round"/></svg>"##,
].map { s in
    let image = NSImage(data: Data(s.utf8))!
    image.size = NSSize(width: 16, height: 16); image.isTemplate = true
    return image
}
let switcherLabels = ["Library", "Browse", "Radio"]

// SwitcherItems fills the toolbar of renderSwitcher as breeze.HeaderBar's end holds a ViewSwitcher in TestScreenshotsViewSwitcher: after a flexible space, an NSToolbarItemGroup of selectionMode selectOne holding the three icons, labeled Library, Browse and Radio, Browse chosen, as an app switches its views from its toolbar (Music, Finder's view switcher). The toolbar shows icons only, so the labels name the segments.
class SwitcherItems: NSObject, NSToolbarDelegate {
    func toolbarDefaultItemIdentifiers(_ t: NSToolbar) -> [NSToolbarItem.Identifier] { [.flexibleSpace, .init("views")] }
    func toolbarAllowedItemIdentifiers(_ t: NSToolbar) -> [NSToolbarItem.Identifier] { toolbarDefaultItemIdentifiers(t) }
    func toolbar(_ t: NSToolbar, itemForItemIdentifier id: NSToolbarItem.Identifier, willBeInsertedIntoToolbar f: Bool) -> NSToolbarItem? {
        let group = NSToolbarItemGroup(itemIdentifier: id, images: switcherIcons, selectionMode: .selectOne, labels: switcherLabels, target: noop, action: #selector(Noop.act(_:)))
        group.selectedIndex = 1
        return group
    }
}

let switcherItems = SwitcherItems()

// renderSwitcher renders TestScreenshotsViewSwitcher's switcherView on macOS: a titled window with a unified toolbar titled Music holding SwitcherItems' group at its end, and in the window's content 16 points in from its side and 16 below the 52-point bar an NSSegmentedControl of the same three segments, each an icon before its label, Browse chosen, where a ViewSwitcher alone draws Breeze UI's segmented control. It renders the window's frame view at 2x into an sRGB bitmap and prints the frames of the group's platter and segmented view and of the segmented control, and the values across the group's middle row and down the chosen segment's middle. The toolbar renders offscreen at its live sizes but its chosen segment's pill 26 points high where a live window draws it 28 high, 4 in from the platter's top and bottom: with -live the window is also put on the screen, ordered front without activating the app, and captured with screencapture -l, which needs the Screen Recording permission; the app is not active then, so the title shows dimmed, while the pill reads as in the key window (the offscreen render draws it alike with or without -inactive).
func renderSwitcher(appearance: NSAppearance.Name, file: String) {
    let live = args.contains("-live")
    let win = KeyWindow(contentRect: NSRect(x: live ? 200 : -4000, y: live ? 300 : -4000, width: 560, height: 140), styleMask: [.titled, .closable, .miniaturizable, .resizable, .fullSizeContentView], backing: .buffered, defer: false)
    win.appearance = NSAppearance(named: appearance)
    win.title = "Music"
    let bar = NSToolbar(identifier: "switcher"); bar.delegate = switcherItems; bar.displayMode = .iconOnly
    win.toolbar = bar; win.toolbarStyle = .unified
    let alone = NSSegmentedControl(labels: switcherLabels, trackingMode: .selectOne, target: noop, action: #selector(Noop.act(_:)))
    for (i, image) in switcherIcons.enumerated() { alone.setImage(image, forSegment: i) }
    alone.selectedSegment = 1
    let size = alone.fittingSize
    alone.frame = NSRect(x: 16, y: 140 - 68 - size.height, width: size.width, height: size.height)
    win.contentView!.addSubview(alone)
    win.orderFront(nil)
    RunLoop.current.run(until: Date().addingTimeInterval(live ? 1 : 0.3))
    let frame = win.contentView!.superview!
    frame.layoutSubtreeIfNeeded(); frame.display()
    let W = frame.bounds.width, H = frame.bounds.height
    let rep = NSBitmapImageRep(bitmapDataPlanes: nil, pixelsWide: Int(W * 2), pixelsHigh: Int(H * 2), bitsPerSample: 8, samplesPerPixel: 4, hasAlpha: true, isPlanar: false, colorSpaceName: .calibratedRGB, bytesPerRow: 0, bitsPerPixel: 0)!.retagging(with: .sRGB)!
    rep.size = frame.bounds.size
    frame.cacheDisplay(in: frame.bounds, to: rep)
    try! rep.representation(using: .png, properties: [:])!.write(to: URL(fileURLWithPath: out + "/" + file))
    print("== \(file)")
    func gray(_ x: Int, _ y: Int) -> String {
        var p = [Int](repeating: 0, count: 4); rep.getPixel(&p, atX: x, y: y)
        return p[0] == p[1] && p[1] == p[2] ? "\(p[0])" : String(format: "#%02x%02x%02x", p[0], p[1], p[2])
    }
    func walk(_ v: NSView) {
        let name = String(describing: type(of: v))
        if name == "NSToolbarPlatterView" || name == "NSToolbarItemGroupPickerExpandedView" || v === alone {
            let f = v.convert(v.bounds, to: frame)
            print("\(name) x=\(f.minX) top=\(H - f.maxY) w=\(f.width) h=\(f.height)")
            let x0 = Int(f.minX * 2), y0 = Int((H - f.maxY) * 2), w = Int(f.width * 2), h = Int(f.height * 2)
            print("  row: " + runs((max(0, x0 - 4)..<min(Int(W * 2), x0 + w + 4)).map { gray($0, y0 + h / 2) }))
            print("  col at the second segment: " + runs((max(0, y0 - 4)..<min(Int(H * 2), y0 + h + 4)).map { gray(x0 + w / 2, $0) }))
        }
        v.subviews.forEach(walk)
    }
    walk(frame)
    if live {
        let capture = Process()
        capture.executableURL = URL(fileURLWithPath: "/usr/sbin/screencapture")
        capture.arguments = ["-x", "-o", "-l\(win.windowNumber)", out + "/" + file.replacingOccurrences(of: ".png", with: "-live.png")]
        try! capture.run(); capture.waitUntilExit()
    }
    win.orderOut(nil)
}

// SwitcherScenes renders the view switchers of TestScreenshotsViewSwitcher (switcherView in ui/breeze/widgets_switcher_test.go), captured on the screen too with -live.
final class SwitcherScenes: SceneGroup {
    override class var file: String { #fileID }
    override class func renderScenes(appearance: NSAppearance.Name, mode: String) {
        renderSwitcher(appearance: appearance, file: "appkit-switcher-\(mode).png")
    }
}
