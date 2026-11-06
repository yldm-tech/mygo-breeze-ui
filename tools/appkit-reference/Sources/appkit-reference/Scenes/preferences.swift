import AppKit
import SwiftUI

// A tile as the preferences scene's PrefixRows hold one: a 30-point rounded square of systemBlue.
struct Tile: View {
    var body: some View { RoundedRectangle(cornerRadius: 7).fill(Color.blue).frame(width: 30, height: 30) }
}

// preferencesSet places SwiftUI's grouped Forms where ui/breeze/widgets_preferences_test.go's prefsView places Breeze UI's preferences widgets on macOS: at 16, 16 a section whose header holds a title and a Detect button and whose footer is a description, holding a row whose label is a Label with a tile, a title and a subtitle and holds a Remove button, a row with a value and a row whose label is a Label with a tile and a title; below it a Picker of the radio-group style of three options, two with a subtitle; on the right a section under a header holding a row with a value and a VStack row of a title over a TextEditor 60 points high. The Forms are shifted by their padding of 20, so that their sections land where Breeze UI's do.
func preferencesSet() -> [Item] {
    let speech = RowsForm {
        Section {
            LabeledContent { Button("Remove") {} } label: { Label { Text("Whisper"); Text("Base · 142 MB") } icon: { Tile() } }
            LabeledContent("Language", value: "Auto")
            LabeledContent { EmptyView() } label: { Label { Text("Tiny") } icon: { Tile() } }
        } header: {
            HStack { Text("Speech"); Spacer(); Button("Detect") {} }
        } footer: { Text("Models are stored on this Mac.") }
    }
    let engine = RowsForm {
        Picker("Engine", selection: .constant("whisper")) {
            TitleSubtitle(title: "Whisper", subtitle: "On this Mac").tag("whisper")
            TitleSubtitle(title: "Cloud", subtitle: "Sent to a service").tag("cloud")
            TitleSubtitle(title: "Off", subtitle: "").tag("off")
        }.pickerStyle(.radioGroup).labelsHidden()
    }
    let translation = RowsForm {
        Section("Translation") {
            LabeledContent("Model", value: "GPT")
            VStack(alignment: .leading) { Text("Prompt"); TextEditor(text: .constant("Translate the subtitles.")).frame(height: 60) }
        }
    }
    return [
        hosted("form-preferences-group", speech, x: -4, y: -4, w: 340, h: 260),
        hosted("form-radio-rows", engine, x: -4, y: 217, w: 340, h: 200),
        hosted("form-content-row", translation, x: 320, y: -4.5, w: 340, h: 240),
    ]
}

// PreferencesScenes renders the preferences widgets of TestScreenshotsPreferences (prefsView in ui/breeze/widgets_preferences_test.go).
final class PreferencesScenes: SceneGroup {
    override class var file: String { #fileID }
    override class func renderScenes(appearance: NSAppearance.Name, mode: String) {
        render(preferencesSet(), width: 672, height: 480, appearance: appearance, file: "appkit-preferences-\(mode).png")
    }
}
