import AppKit
import SwiftUI

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

// RowsScenes renders the cards and rows of TestScreenshotsRows (rowsView in ui/breeze/widgets_rows_test.go).
final class RowsScenes: SceneGroup {
    override class var file: String { #fileID }
    override class func renderScenes(appearance: NSAppearance.Name, mode: String) {
        render(rowsSet(), width: 672, height: 520, appearance: appearance, file: "appkit-rows-\(mode).png")
    }
}
