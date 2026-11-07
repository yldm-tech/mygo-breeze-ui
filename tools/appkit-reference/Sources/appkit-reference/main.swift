// AppKit's own controls, for comparing Breeze UI's macOS look with them.
//
// It renders sets of controls in Aqua and Dark Aqua at 2x, each group of them (Scenes/<group>.swift, a SceneGroup) laid out where the screenshot test of a group of Breeze UI's widgets places Breeze UI's, writes them as PNGs and prints, for every control, its frame and the values along its middle row and column as runs (value x count): a gray pixel as its value, a colored one as #rrggbb, from which the sizes and fills in ui/breeze/mac.go were read. It draws into an sRGB bitmap, the color space MyGo writes, rather than the display's profile that bitmapImageRepForCachingDisplay picks, in which grays read the same but colors paler (systemGreen reads #65c466 in Display P3, #34c759 in sRGB). With -live, the scenes that read it also put their window on the screen and capture it.
//
// An NSTextField drawn by cacheDisplay lays its text color's alpha on twice (a label in 50% black reads 63 over white, in labelColor's 85% black 6, where NSString.draw in the same font reads 127 and 39, measured on macOS 27.2), so read text colors from the label colors themselves rather than from these renders.
//
// The window is offscreen and never key, so AppKit would draw it inactive, its accent controls (a checked box, an on switch, the default button, the chosen segment) gray. Its NSWindow subclass answers hasKeyAppearance (an NSWindow method outside the SDK's headers) with true, so they are drawn as in the active window, as Breeze UI draws them; -inactive leaves the window inactive. Materials (Liquid Glass, vibrancy), menus, popovers, tooltips and focus rings still do not render into the bitmap.
//
//	swift run -c release --package-path tools/appkit-reference appkit-reference [-inactive] [-live] <dir>
import AppKit

let app = NSApplication.shared
app.setActivationPolicy(.accessory)

// The groups render one appearance after the other, Aqua first.
let groups = sceneGroups()
for (appearance, mode) in [(NSAppearance.Name.aqua, "light"), (.darkAqua, "dark")] {
    for group in groups {
        group.renderScenes(appearance: appearance, mode: mode)
    }
}
