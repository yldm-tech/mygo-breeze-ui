import AppKit
import SwiftUI

let args = CommandLine.arguments.dropFirst()
let inactive = args.contains("-inactive")
let out = args.last(where: { !$0.hasPrefix("-") }) ?? "."

// SceneGroup is a group of scenes, one subclass in each file of Scenes/, which sceneGroups finds at run time, so that adding a group touches no other file. A subclass returns its file (#fileID), which orders the groups, and renders its scenes in an appearance, mode naming it in the files it writes (light or dark).
class SceneGroup: NSObject {
    class var file: String { "" }
    class func renderScenes(appearance: NSAppearance.Name, mode: String) {}
}

// sceneGroups returns the subclasses of SceneGroup in this executable in the order of their files. It exits with 1 when there are none or two share a file.
func sceneGroups() -> [SceneGroup.Type] {
    var count: UInt32 = 0
    guard let image = class_getImageName(SceneGroup.self), let names = objc_copyClassNamesForImage(image, &count) else {
        FileHandle.standardError.write(Data("appkit-reference: cannot list the classes of the executable\n".utf8))
        exit(1)
    }
    defer { free(names) }
    var groups: [SceneGroup.Type] = []
    for i in 0..<Int(count) {
        // The names are this executable's classes only, and only SceneGroup's subclasses are cast: casting every class of the process (objc_copyClassList) to a Swift type crashes on some of the system's.
        if let cls = objc_lookUpClass(names[i]), class_getSuperclass(cls).map(ObjectIdentifier.init) == ObjectIdentifier(SceneGroup.self) {
            groups.append(cls as! SceneGroup.Type)
        }
    }
    groups.sort { $0.file < $1.file }
    let files = groups.map { $0.file }
    if groups.isEmpty || Set(files).count != files.count {
        FileHandle.standardError.write(Data("appkit-reference: want one SceneGroup per file of Scenes/, found \(files)\n".utf8))
        exit(1)
    }
    return groups
}

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

// Noop is the target of controls that draw disabled without an action.
class Noop: NSObject {
    @objc func act(_ sender: Any?) {}
}

let noop = Noop()

// A grouped Form of SwiftUI, as System Settings lays out its rows, padded 20 points around its sections, without the background it would lay over the items placed before it.
struct RowsForm<Content: View>: View {
    @ViewBuilder var content: () -> Content
    var body: some View { Form(content: content).formStyle(.grouped).scrollDisabled(true).scrollContentBackground(.hidden) }
}

// A SwiftUI view as an item whose top left corner is at x, y, w by h points.
func hosted<V: View>(_ name: String, _ view: V, x: CGFloat, y: CGFloat, w: CGFloat, h: CGFloat) -> Item {
    Item(name: name, view: NSHostingView(rootView: view), x: x, y: y, size: NSSize(width: w, height: h))
}

// A DisclosureGroup's label as ExpanderRow draws it on macOS: a title over a secondary subtitle, as LabeledContent lays out a row's label.
struct TitleSubtitle: View {
    let title: String, subtitle: String
    var body: some View {
        VStack(alignment: .leading, spacing: 2.5) {
            Text(title)
            if !subtitle.isEmpty { Text(subtitle).font(.subheadline).foregroundStyle(.secondary) }
        }
    }
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
    // SwiftUI's hosting views lay out their content on the run loop.
    RunLoop.current.run(until: Date().addingTimeInterval(0.3))
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

// renderToolbar renders a titled window's unified toolbar holding the items of items, ToolbarItems' unless given, the window's frame view at 2x into an sRGB bitmap, and prints the frames of the items' views and platters. It renders offscreen with the items' platters, 36 points high as a live capture shows them though lighter: the sizes inside them compare, and their fills as shares of the platter.
func renderToolbar(appearance: NSAppearance.Name, file: String, items: NSToolbarDelegate = toolbarItems, identifier: String = "inputs") {
    let win = KeyWindow(contentRect: NSRect(x: -4000, y: -4000, width: 520, height: 120), styleMask: [.titled, .closable, .fullSizeContentView], backing: .buffered, defer: false)
    win.appearance = NSAppearance(named: appearance)
    win.titleVisibility = .hidden
    let bar = NSToolbar(identifier: identifier); bar.delegate = items; bar.displayMode = .iconOnly
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
