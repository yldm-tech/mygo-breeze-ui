// AppKit's own controls, for comparing Breeze UI's macOS look with them.
//
// It renders two sets of controls in Aqua and Dark Aqua at 2x, laid out where ui/breeze/screenshot_test.go places Breeze UI's (referenceView and referenceView2), writes them as PNGs and prints, for every control, its frame and the values along its middle row and column as runs (value x count): a gray pixel as its value, a colored one as #rrggbb, from which the sizes and fills in ui/breeze/mac.go were read. It draws into an sRGB bitmap, the color space MyGo writes, rather than the display's profile that bitmapImageRepForCachingDisplay picks, in which grays read the same but colors paler (systemGreen reads #65c466 in Display P3, #34c759 in sRGB).
//
// The window is offscreen and never key, so AppKit would draw it inactive, its accent controls (a checked box, an on switch, the default button, the chosen segment) gray. Its NSWindow subclass answers hasKeyAppearance (an NSWindow method outside the SDK's headers) with true, so they are drawn as in the active window, as Breeze UI draws them; -inactive leaves the window inactive. Materials (Liquid Glass, vibrancy), menus, popovers, tooltips and focus rings still do not render into the bitmap.
//
//	swift tools/appkit-reference/reference.swift [-inactive] <dir>
import AppKit
import SwiftUI

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

// A grouped Form of SwiftUI, as System Settings lays out its rows, padded 20 points around its sections, without the background it would lay over the items placed before it.
struct RowsForm<Content: View>: View {
    @ViewBuilder var content: () -> Content
    var body: some View { Form(content: content).formStyle(.grouped).scrollDisabled(true).scrollContentBackground(.hidden) }
}

// A SwiftUI view as an item whose top left corner is at x, y, w by h points.
func hosted<V: View>(_ name: String, _ view: V, x: CGFloat, y: CGFloat, w: CGFloat, h: CGFloat) -> Item {
    Item(name: name, view: NSHostingView(rootView: view), x: x, y: y, size: NSSize(width: w, height: h))
}

// rowsSet places SwiftUI's GroupBox and grouped Forms where ui/breeze/widgets_rows_test.go's rowsView places Breeze UI's cards and rows on macOS: a GroupBox of text at 16, 16, a section of a button row, a switch row with a subtitle and one without at 16, 57.5, a row with a value and a disabled switch row alone, and on the right a section under a header and a disclosure group of rows. The Forms are shifted by their padding of 20, so that their sections land where Breeze UI's do.
func rowsSet() -> [Item] {
    let card = GroupBox { Text("A card holds what belongs together.").frame(maxWidth: .infinity, alignment: .leading) }
    let settings = RowsForm {
        Section {
            LabeledContent { Button("Open") {} } label: { Text("Wallpaper"); Text("Choose a picture") }
            Toggle(isOn: .constant(true)) { Text("Dark mode"); Text("Use a dark appearance") }
            Toggle("Notifications", isOn: .constant(false))
        }
    }
    let storage = RowsForm { LabeledContent("Storage", value: "120 GB") }
    let bluetooth = RowsForm { Toggle("Bluetooth", isOn: .constant(true)).disabled(true) }
    let privacy = RowsForm {
        Section("Privacy") {
            Toggle(isOn: .constant(true)) { Text("Location"); Text("Apps may ask for it") }
            LabeledContent { EmptyView() } label: { Text("Camera"); Text("2 apps") }
        }
    }
    let network = RowsForm {
        DisclosureGroup("Network", isExpanded: .constant(true)) {
            LabeledContent { EmptyView() } label: { Text("Proxy"); Text("Automatic") }
            Toggle("Metered", isOn: .constant(false))
        }
    }
    return [
        hosted("group-box", card.frame(width: 300), x: 16, y: 16, w: 300, h: 26),
        hosted("form-rows", settings, x: -4, y: 37.5, w: 340, h: 200),
        hosted("form-row-value", storage, x: -4, y: 196.5, w: 340, h: 80),
        hosted("form-switch-disabled", bluetooth, x: -4, y: 249.5, w: 340, h: 77),
        hosted("form-section-header", privacy, x: 320, y: -4.5, w: 340, h: 190),
        hosted("form-disclosure-rows", network, x: 320, y: 143, w: 340, h: 166),
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

for (appearance, mode) in [(NSAppearance.Name.aqua, "light"), (.darkAqua, "dark")] {
    render(firstSet(), width: 360, height: 200, appearance: appearance, file: "appkit-\(mode).png")
    render(secondSet(), width: 520, height: 330, appearance: appearance, file: "appkit2-\(mode).png")
    render(rowsSet(), width: 672, height: 520, appearance: appearance, file: "appkit-rows-\(mode).png")
}
