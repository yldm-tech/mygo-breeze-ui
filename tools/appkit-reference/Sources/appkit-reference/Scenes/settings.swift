import AppKit
import SwiftUI

// settingsSet places SwiftUI's grouped Forms where ui/breeze/widgets_settings_test.go's settingsView places Breeze UI's settings rows on macOS: a section of a Picker of the menu style with a subtitle, a TextField with a label holding a name, an empty one and a Stepper with a format and a subtitle at 16, 16; a Picker alone and a disabled Stepper alone below it; on the right a DisclosureGroup, open, whose label is a title over a subtitle, holding a row with a subtitle and a switch, and a closed one. The Forms are shifted by their padding of 20, so that their sections land where Breeze UI's do.
func settingsSet() -> [Item] {
    let settings = RowsForm {
        Section {
            Picker(selection: .constant("Auto")) { ForEach(["Light", "Dark", "Auto"], id: \.self) { Text($0).tag($0) } } label: { Text("Appearance"); Text("Light or dark") }
            TextField("Name", text: .constant("Evan"))
            TextField("Email", text: .constant(""))
            Stepper(value: .constant(12.0), in: 6...72, step: 1, format: .number) { Text("Font size"); Text("Points") }
        }
    }
    let language = RowsForm { Picker("Language", selection: .constant("English")) { ForEach(["English", "Français"], id: \.self) { Text($0).tag($0) } } }
    let copies = RowsForm { Stepper(value: .constant(3.0), in: 1...9, step: 1, format: .number) { Text("Copies") }.disabled(true) }
    let network = RowsForm {
        DisclosureGroup(isExpanded: .constant(true)) {
            LabeledContent { EmptyView() } label: { Text("Proxy"); Text("Automatic") }
            Toggle("Metered", isOn: .constant(false))
        } label: { TitleSubtitle(title: "Network", subtitle: "Wi-Fi and proxies") }
    }
    let advanced = RowsForm { DisclosureGroup(isExpanded: .constant(false)) { Text("Logs") } label: { TitleSubtitle(title: "Advanced", subtitle: "") } }
    return [
        hosted("form-settings", settings, x: -4, y: -4, w: 340, h: 230),
        hosted("form-language", language, x: -4, y: 193, w: 340, h: 80),
        hosted("form-copies", copies, x: -4, y: 246, w: 340, h: 80),
        hosted("form-network", network, x: 320, y: -4, w: 340, h: 170),
        hosted("form-advanced", advanced, x: 320, y: 129, w: 340, h: 90),
    ]
}

// SettingsScenes renders the settings rows of TestScreenshotsSettings (settingsView in ui/breeze/widgets_settings_test.go).
final class SettingsScenes: SceneGroup {
    override class var file: String { #fileID }
    override class func renderScenes(appearance: NSAppearance.Name, mode: String) {
        render(settingsSet(), width: 672, height: 480, appearance: appearance, file: "appkit-settings-\(mode).png")
    }
}
