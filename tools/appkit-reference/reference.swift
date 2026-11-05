// AppKit's own controls, for comparing Breeze UI's macOS look with them.
//
// It renders two sets of controls in Aqua and Dark Aqua at 2x, laid out where ui/breeze/screenshot_test.go places Breeze UI's (referenceView and referenceView2), and the text styles where ui/breeze/widgets_text_test.go places Breeze UI's Text (textView), writes them as PNGs and prints, for every control, its frame and the values along its middle row and column as runs (value x count): a gray pixel as its value, a colored one as #rrggbb, from which the sizes and fills in ui/breeze/mac.go were read. It draws into an sRGB bitmap, the color space MyGo writes, rather than the display's profile that bitmapImageRepForCachingDisplay picks, in which grays read the same but colors paler (systemGreen reads #65c466 in Display P3, #34c759 in sRGB).
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

// textSet: AppKit's text styles (NSFont.preferredFont(forTextStyle:)) in labels, two in secondaryLabelColor, where TestScreenshotsText's textView places Breeze UI's Text (textTops).
func textSet() -> [Item] {
    func label(_ s: String, _ style: NSFont.TextStyle, secondary: Bool = false) -> NSTextField {
        let l = NSTextField(labelWithString: s)
        l.font = NSFont.preferredFont(forTextStyle: style, options: [:])
        if secondary { l.textColor = .secondaryLabelColor }
        return l
    }
    let rows: [(String, NSFont.TextStyle, Bool, CGFloat)] = [
        ("Large Title", .largeTitle, false, 16), ("Title", .title1, false, 76), ("Subtitle", .title2, false, 122), ("Heading", .headline, false, 160),
        ("Body", .body, false, 188), ("Caption", .caption1, false, 214), ("Secondary body", .body, true, 240), ("Secondary caption", .caption1, true, 266),
    ]
    return rows.map { Item(name: $0.0, view: label($0.0, $0.1, secondary: $0.2), x: 16, y: $0.3) }
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

for (appearance, mode) in [(NSAppearance.Name.aqua, "light"), (.darkAqua, "dark")] {
    render(firstSet(), width: 360, height: 200, appearance: appearance, file: "appkit-\(mode).png")
    render(secondSet(), width: 520, height: 330, appearance: appearance, file: "appkit2-\(mode).png")
    render(textSet(), width: 360, height: 300, appearance: appearance, file: "appkit-text-\(mode).png")
}
