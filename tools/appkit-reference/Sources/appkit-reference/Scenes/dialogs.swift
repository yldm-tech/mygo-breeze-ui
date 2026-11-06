import AppKit

// renderDialogs shows, over a titled window of 560 by 340 points as TestScreenshotsDialogs's, an NSWindow sheet of 300 by 150 holding a push button with hasDestructiveAction and a default push button, then an NSAlert whose first button, Delete, has hasDestructiveAction, as breeze.Sheet and breeze.DestructiveAlertDialog draw them, and prints the frames of the sheet and the alert against the window's and the alert's buttons with their key equivalents. Sheets and alerts do not render offscreen: with -live the window is put on the screen and captured with screencapture -l, which shows the sheet over its window, and the capture is converted from the display's profile to sRGB and printed along its middle row and the sheet's middle column, from which the dimming, the sheet's fill and corners and the destructive button's face in ui/breeze/widgets_dialogs.go were read.
func renderDialogs(appearance: NSAppearance.Name, file: String) {
    let live = args.contains("-live")
    let win = KeyWindow(contentRect: NSRect(x: live ? 200 : -4000, y: live ? 200 : -4000, width: 560, height: 340), styleMask: [.titled, .closable], backing: .buffered, defer: false)
    win.appearance = NSAppearance(named: appearance)
    win.title = "Window content"
    win.orderFrontRegardless()
    RunLoop.current.run(until: Date().addingTimeInterval(0.5))
    print("== \(file)")
    func capture(_ name: String) {
        guard live else { return }
        let path = out + "/" + file.replacingOccurrences(of: ".png", with: "-\(name)-live.png")
        let p = Process()
        p.executableURL = URL(fileURLWithPath: "/usr/sbin/screencapture")
        p.arguments = ["-x", "-o", "-l\(win.windowNumber)", path]
        try! p.run(); p.waitUntilExit()
        let rep = (NSImageRep(contentsOfFile: path) as! NSBitmapImageRep).converting(to: .sRGB, renderingIntent: .default)!
        func gray(_ x: Int, _ y: Int) -> String {
            var px = [Int](repeating: 0, count: 4); rep.getPixel(&px, atX: x, y: y)
            return px[0] == px[1] && px[1] == px[2] ? "\(px[0])" : String(format: "#%02x%02x%02x", px[0], px[1], px[2])
        }
        print("  \(name) capture \(rep.pixelsWide)x\(rep.pixelsHigh)")
        print("  row: " + runs((0..<rep.pixelsWide).map { gray($0, rep.pixelsHigh / 2) }))
        print("  col: " + runs((0..<rep.pixelsHigh).map { gray(rep.pixelsWide / 2, $0) }))
    }
    let sheet = KeyWindow(contentRect: NSRect(x: 0, y: 0, width: 300, height: 150), styleMask: [.titled], backing: .buffered, defer: false)
    let delete = NSButton(title: "Delete", target: nil, action: nil); delete.bezelStyle = .push; delete.hasDestructiveAction = true
    delete.frame = NSRect(x: 20, y: 20, width: 100, height: 30); sheet.contentView!.addSubview(delete)
    let cancel = NSButton(title: "Cancel", target: nil, action: nil); cancel.bezelStyle = .push; cancel.keyEquivalent = "\r"
    cancel.frame = NSRect(x: 180, y: 20, width: 100, height: 30); sheet.contentView!.addSubview(cancel)
    win.beginSheet(sheet, completionHandler: nil)
    RunLoop.current.run(until: Date().addingTimeInterval(1.2))
    let w = win.frame, s = sheet.frame
    print("window \(w.width)x\(w.height) title bar \(w.height - win.contentLayoutRect.maxY) sheet x=\(s.minX - w.minX) top=\(w.maxY - s.maxY) w=\(s.width) h=\(s.height)")
    capture("sheet")
    win.endSheet(sheet)
    RunLoop.current.run(until: Date().addingTimeInterval(0.8))
    let alert = NSAlert()
    alert.messageText = "Delete “Notes”?"
    alert.informativeText = "You can’t undo this action."
    alert.addButton(withTitle: "Delete")
    alert.addButton(withTitle: "Cancel")
    alert.buttons[0].hasDestructiveAction = true
    alert.beginSheetModal(for: win, completionHandler: nil)
    RunLoop.current.run(until: Date().addingTimeInterval(1.2))
    let a = alert.window.frame
    print("alert x=\(a.minX - w.minX) top=\(w.maxY - a.maxY) w=\(a.width) h=\(a.height)")
    for b in alert.buttons {
        print("  button \"\(b.title)\" x=\(b.frame.minX) w=\(b.frame.width) h=\(b.frame.height) key=\(b.keyEquivalent.debugDescription) destructive=\(b.hasDestructiveAction)")
    }
    capture("alert")
    win.endSheet(alert.window)
    win.orderOut(nil)
}

// DialogsScenes renders a sheet and an alert with a destructive button over a window, for breeze.Sheet and breeze.DestructiveAlertDialog (ui/breeze/widgets_dialogs_test.go); only -live captures them.
final class DialogsScenes: SceneGroup {
    override class var file: String { #fileID }
    override class func renderScenes(appearance: NSAppearance.Name, mode: String) {
        renderDialogs(appearance: appearance, file: "appkit-dialogs-\(mode).png")
    }
}
