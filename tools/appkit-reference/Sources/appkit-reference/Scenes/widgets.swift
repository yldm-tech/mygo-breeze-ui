import AppKit

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

// WidgetsScenes renders AppKit's counterparts of the controls of TestScreenshots (ui/breeze/screenshot_test.go): referenceView's first set and referenceView2's second.
final class WidgetsScenes: SceneGroup {
    override class var file: String { #fileID }
    override class func renderScenes(appearance: NSAppearance.Name, mode: String) {
        render(firstSet(), width: 360, height: 200, appearance: appearance, file: "appkit-\(mode).png")
        render(secondSet(), width: 520, height: 330, appearance: appearance, file: "appkit2-\(mode).png")
    }
}
