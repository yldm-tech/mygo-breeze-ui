import AppKit

// NewTabs answers newWindowForTab: at the end of the responder chain, as the app's delegate, which makes a tab bar show its + button.
class NewTabs: NSObject, NSApplicationDelegate {
    @objc func newWindowForTab(_ sender: Any?) {}
}

let newTabs = NewTabs()

// renderTabs joins three titled windows, Documents, Downloads and Pictures, with addTabbedWindow, Documents chosen, as TestScreenshotsTabView's first TabView holds them, and prints the frames of the window tab bar's parts (NSTabBar, its track, the NSTabButtons, the separators between them and the NSTabBarNewTabButton) and of each tab's close button, which shows only under the pointer, by its accessibility frame, then renders the window's frame view at 2x into an sRGB bitmap. The bar's glass does not render offscreen: with -live the window is put on the screen, ordered front without activating the app, and captured with screencapture -l, which needs the Screen Recording permission. The app is not active then, so the capture shows the titles dimmed. The + button appears once the tabs are laid out, and AppKit keeps the widths the tabs had without it for all but the last tab: the equal widths are those of a bar without the button.
func renderTabs(appearance: NSAppearance.Name, file: String) {
    app.delegate = newTabs
    let live = args.contains("-live")
    let windows = ["Documents", "Downloads", "Pictures"].map { title -> KeyWindow in
        let w = KeyWindow(contentRect: NSRect(x: live ? 200 : -4000, y: live ? 300 : -4000, width: 560, height: 120), styleMask: [.titled, .closable, .miniaturizable, .resizable], backing: .buffered, defer: false)
        w.title = title; w.tabbingMode = .preferred; w.tabbingIdentifier = "tabs"; w.isReleasedWhenClosed = false
        w.appearance = NSAppearance(named: appearance)
        return w
    }
    windows[0].orderFront(nil)
    windows[0].addTabbedWindow(windows[1], ordered: .above)
    windows[1].addTabbedWindow(windows[2], ordered: .above)
    windows[0].tabGroup?.selectedWindow = windows[0]
    // The tab bar lays out and animates its tabs in on the run loop.
    RunLoop.current.run(until: Date().addingTimeInterval(2))
    let frame = windows[0].contentView!.superview!
    frame.layoutSubtreeIfNeeded(); frame.display()
    let W = frame.bounds.width, H = frame.bounds.height
    let rep = NSBitmapImageRep(bitmapDataPlanes: nil, pixelsWide: Int(W * 2), pixelsHigh: Int(H * 2), bitsPerSample: 8, samplesPerPixel: 4, hasAlpha: true, isPlanar: false, colorSpaceName: .calibratedRGB, bytesPerRow: 0, bitsPerPixel: 0)!.retagging(with: .sRGB)!
    rep.size = frame.bounds.size
    frame.cacheDisplay(in: frame.bounds, to: rep)
    try! rep.representation(using: .png, properties: [:])!.write(to: URL(fileURLWithPath: out + "/" + file))
    print("== \(file)")
    let screen = windows[0].frame
    func walk(_ v: NSView) {
        let name = String(describing: type(of: v))
        if name == "NSTabBar" || name == "NSTabBarTrackView" || name == "NSTabButton" || name == "NSTabBarNewTabButton" || (name == "NSTabBarTrackInversionView" && v.frame.width < 2) {
            let f = v.convert(v.bounds, to: frame)
            print("\(name) x=\(f.minX) top=\(H - f.maxY) w=\(f.width) h=\(f.height) title=\(v.accessibilityTitle() ?? "") label=\(v.accessibilityLabel() ?? "")")
            if name == "NSTabButton" {
                for child in v.accessibilityChildren() ?? [] {
                    if let e = child as? NSAccessibilityElementProtocol {
                        let r = e.accessibilityFrame()
                        print("  close x=\(r.minX - screen.minX) top=\(screen.maxY - r.maxY) w=\(r.width) h=\(r.height)")
                    }
                }
            }
        }
        v.subviews.forEach(walk)
    }
    walk(frame)
    if live {
        let capture = Process()
        capture.executableURL = URL(fileURLWithPath: "/usr/sbin/screencapture")
        capture.arguments = ["-x", "-o", "-l\(windows[0].windowNumber)", out + "/" + file.replacingOccurrences(of: ".png", with: "-live.png")]
        try! capture.run(); capture.waitUntilExit()
    }
    windows.forEach { $0.orderOut(nil) }
    app.delegate = nil
}

// TabViewScenes renders the tab bar of tabbed windows for TestScreenshotsTabView (ui/breeze/widgets_tabview_test.go), captured on the screen too with -live.
final class TabViewScenes: SceneGroup {
    override class var file: String { #fileID }
    override class func renderScenes(appearance: NSAppearance.Name, mode: String) {
        renderTabs(appearance: appearance, file: "appkit-tabs-\(mode).png")
    }
}
