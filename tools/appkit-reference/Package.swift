// swift-tools-version: 6.0
// AppKit's own controls, for comparing Breeze UI's macOS look with them: see Sources/appkit-reference/main.swift.
import PackageDescription

let package = Package(
    name: "appkit-reference",
    // The renders are of macOS 27's controls. SwiftPM records the deployment target as the SDK the executable was linked with, and AppKit draws for that SDK: built for macOS 14, the lowest the code compiles for, its window tab bar is 4 points shorter and its slider knob differs, as before macOS 26.
    platforms: [.macOS("27.0")],
    targets: [
        .executableTarget(name: "appkit-reference", swiftSettings: [.swiftLanguageMode(.v5)]),
    ]
)
