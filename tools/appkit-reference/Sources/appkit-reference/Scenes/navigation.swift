import AppKit

// SourceList fills the outline view of renderSidebar as a sidebar: a group row, Mailboxes, holding Inbox, Sent and Drafts, each a cell of a template SF Symbol and a title, and Flagged, whose image is not a template, a red square, as breeze.SidebarItemImage shows a colored image.
class SourceList: NSObject, NSOutlineViewDataSource, NSOutlineViewDelegate {
    let group = "Mailboxes"
    let items = [("Inbox", "tray"), ("Sent", "paperplane"), ("Drafts", "doc"), ("Flagged", "")]
    func outlineView(_ o: NSOutlineView, numberOfChildrenOfItem item: Any?) -> Int { item == nil ? 1 : items.count }
    func outlineView(_ o: NSOutlineView, child i: Int, ofItem item: Any?) -> Any { item == nil ? group : items[i].0 }
    func outlineView(_ o: NSOutlineView, isItemExpandable item: Any) -> Bool { (item as? String) == group }
    func outlineView(_ o: NSOutlineView, isGroupItem item: Any) -> Bool { (item as? String) == group }
    func outlineView(_ o: NSOutlineView, shouldSelectItem item: Any) -> Bool { (item as? String) != group }
    func outlineView(_ o: NSOutlineView, viewFor col: NSTableColumn?, item: Any) -> NSView? {
        let name = item as! String
        let cell = NSTableCellView()
        let text = NSTextField(labelWithString: name)
        text.translatesAutoresizingMaskIntoConstraints = false
        cell.addSubview(text); cell.textField = text
        if name == group {
            NSLayoutConstraint.activate([text.leadingAnchor.constraint(equalTo: cell.leadingAnchor, constant: 2), text.centerYAnchor.constraint(equalTo: cell.centerYAnchor)])
            return cell
        }
        let symbol = items.first(where: { $0.0 == name })!.1
        let image: NSImage
        if symbol.isEmpty {
            image = NSImage(size: NSSize(width: 16, height: 16), flipped: false) { r in NSColor(srgbRed: 1, green: 0, blue: 0, alpha: 1).setFill(); NSBezierPath(roundedRect: r, xRadius: 4, yRadius: 4).fill(); return true }
        } else {
            image = NSImage(systemSymbolName: symbol, accessibilityDescription: name)!
        }
        let view = NSImageView(image: image)
        view.translatesAutoresizingMaskIntoConstraints = false
        cell.addSubview(view); cell.imageView = view
        NSLayoutConstraint.activate([
            view.leadingAnchor.constraint(equalTo: cell.leadingAnchor), view.widthAnchor.constraint(equalToConstant: 23), view.centerYAnchor.constraint(equalTo: cell.centerYAnchor),
            text.leadingAnchor.constraint(equalTo: view.trailingAnchor, constant: 3), text.centerYAnchor.constraint(equalTo: cell.centerYAnchor),
        ])
        return cell
    }
}

let sourceList = SourceList()

// renderSidebar renders an NSOutlineView of style sourceList at the medium row size, as breeze.Sidebar draws it on macOS, 200 points wide, the item chosen (Sent unless given), the outline view the window's first responder when focused, so that AppKit draws the emphasized selection, and prints the frames of the rows and their image views and, for every row, the pixels across the middle of its image view: whether the template images take the text color or the accent, and whether the colored one keeps its own color on the emphasized selection.
func renderSidebar(appearance: NSAppearance.Name, focused: Bool, chosen: String = "Sent", file: String) {
    let W: CGFloat = 200, H: CGFloat = 200
    let outline = NSOutlineView()
    outline.style = .sourceList
    outline.rowSizeStyle = .medium
    outline.headerView = nil
    let col = NSTableColumn(identifier: .init("main")); col.width = W
    outline.addTableColumn(col); outline.outlineTableColumn = col
    outline.dataSource = sourceList; outline.delegate = sourceList
    outline.frame = NSRect(x: 0, y: 0, width: W, height: H)
    let root = NSView(frame: NSRect(x: 0, y: 0, width: W, height: H))
    root.appearance = NSAppearance(named: appearance)
    let scroll = NSScrollView(frame: root.bounds); scroll.documentView = outline; scroll.drawsBackground = false
    root.addSubview(scroll)
    let win = KeyWindow(contentRect: NSRect(x: -4000, y: -4000, width: W, height: H), styleMask: [.borderless], backing: .buffered, defer: false)
    win.contentView = root
    outline.reloadData(); outline.expandItem(nil, expandChildren: true)
    outline.selectRowIndexes(IndexSet(integer: outline.row(forItem: chosen)), byExtendingSelection: false)
    if focused { win.makeFirstResponder(outline) } else { win.makeFirstResponder(nil) }
    root.layoutSubtreeIfNeeded()
    RunLoop.current.run(until: Date().addingTimeInterval(0.3))
    root.layoutSubtreeIfNeeded(); root.display()
    let rep = NSBitmapImageRep(bitmapDataPlanes: nil, pixelsWide: Int(W * 2), pixelsHigh: Int(H * 2), bitsPerSample: 8, samplesPerPixel: 4, hasAlpha: true, isPlanar: false, colorSpaceName: .calibratedRGB, bytesPerRow: 0, bitsPerPixel: 0)!.retagging(with: .sRGB)!
    rep.size = root.bounds.size
    root.cacheDisplay(in: root.bounds, to: rep)
    try! rep.representation(using: .png, properties: [:])!.write(to: URL(fileURLWithPath: out + "/" + file))
    func gray(_ x: Int, _ y: Int) -> String {
        var p = [Int](repeating: 0, count: 4); rep.getPixel(&p, atX: x, y: y)
        return p[0] == p[1] && p[1] == p[2] ? "\(p[0])" : String(format: "#%02x%02x%02x", p[0], p[1], p[2])
    }
    print("== \(file) first responder \(win.firstResponder === outline)")
    for row in 0..<outline.numberOfRows {
        guard let rv = outline.rowView(atRow: row, makeIfNecessary: false), let cell = outline.view(atColumn: 0, row: row, makeIfNecessary: false) as? NSTableCellView else { continue }
        let rf = rv.convert(rv.bounds, to: root)
        print("row \(row) \"\(cell.textField?.stringValue ?? "")\" x=\(rf.minX) top=\(H - rf.maxY) w=\(rf.width) h=\(rf.height) selected=\(rv.isSelected) emphasized=\(rv.isEmphasized) tint=\(cell.imageView?.contentTintColor?.description ?? "nil")")
        if let iv = cell.imageView {
            let f = iv.convert(iv.bounds, to: root)
            print("  image x=\(f.minX) top=\(H - f.maxY) w=\(f.width) h=\(f.height)")
            print("  row: " + runs((Int(f.minX * 2)..<Int(f.maxX * 2)).map { gray($0, Int((H - f.midY) * 2)) }))
        }
    }
}

// NavigationScenes renders a source list for the icons of Sidebar's items, which TestSidebarIconTint (ui/breeze/widgets_sidebar_test.go) compares with: without and with the focus, a template image chosen and then a colored one.
final class NavigationScenes: SceneGroup {
    override class var file: String { #fileID }
    override class func renderScenes(appearance: NSAppearance.Name, mode: String) {
        renderSidebar(appearance: appearance, focused: false, file: "appkit-sidebar-\(mode).png")
        renderSidebar(appearance: appearance, focused: true, file: "appkit-sidebar-\(mode)-focused.png")
        renderSidebar(appearance: appearance, focused: true, chosen: "Flagged", file: "appkit-sidebar-image-\(mode)-focused.png")
    }
}
