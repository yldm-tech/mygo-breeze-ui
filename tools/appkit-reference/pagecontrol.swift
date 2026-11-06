// UIKit for Mac's UIPageControl, the page indicator that macOS draws, for comparing Breeze UI's macOS PipsPager and Carousel with it: AppKit has no page indicator and SwiftUI's PageTabViewStyle and IndexViewStyle are unavailable on macOS (the macOS 27.2 SDK marks them @available(macOS, unavailable)), while a Mac Catalyst app draws UIPageControl in the Mac idiom, which SwiftUI's TabView of the page style puts below its pages there.
//
// It renders, in light and dark at 2x, UIPageControls where TestScreenshotsCarousel (carouselShotView in ui/breeze/widgets_carousel_test.go) places Breeze UI's pagers: 5 pages with the 2nd chosen at 16, 16 and 9 pages with the 7th chosen at 16, 56, and 4 pages with the 2nd chosen centered right below the carousel's 320 by 160 view at 16, 96, over the window's background of Breeze UI's macOS theme (white, #1e1e1e), with "Slide 2" centered in the view in the system font of 13 points; then a SwiftUI TabView of the page style with four pages, 320 by 200, alone. It writes them as PNGs and prints the frames of the controls and of their dots, the interface idiom (5 is the Mac), and the values along each control's middle row as runs (value x count): a gray pixel as its value, a colored one as #rrggbb. UIKit draws the dots white in either appearance, the chosen one opaque and the others at 45%, for content behind them, so in the light window they do not show.
//
// It is a Mac Catalyst app, built for the macabi target and run from a bundle whose UIDeviceFamily is 6 (the Mac idiom):
//
//	SDK=$(xcrun --sdk macosx --show-sdk-path)
//	xcrun swiftc -target arm64-apple-ios27.2-macabi -sdk $SDK -Fsystem $SDK/System/iOSSupport/System/Library/Frameworks -L $SDK/System/iOSSupport/usr/lib tools/appkit-reference/pagecontrol.swift -o <bundle>/Contents/MacOS/pagecontrol
//	<bundle>/Contents/MacOS/pagecontrol <dir>
//
// with <bundle>/Contents/Info.plist naming the executable (CFBundleExecutable pagecontrol), a bundle identifier, CFBundlePackageType APPL, UIDeviceFamily [6], LSUIElement true and an empty UIApplicationSceneManifest, the bundle signed ad hoc (codesign -s - <bundle>).
import SwiftUI
import UIKit

class AppDelegate: UIResponder, UIApplicationDelegate {
    func application(_ application: UIApplication, didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]? = nil) -> Bool { true }
    func application(_ application: UIApplication, configurationForConnecting session: UISceneSession, options: UIScene.ConnectionOptions) -> UISceneConfiguration {
        let config = UISceneConfiguration(name: nil, sessionRole: session.role)
        config.delegateClass = SceneDelegate.self
        return config
    }
}

// SceneDelegate renders once UIKit has its scene, as a Mac Catalyst app must adopt the scene life cycle, and quits.
class SceneDelegate: UIResponder, UIWindowSceneDelegate {
    var window: UIWindow?
    func scene(_ scene: UIScene, willConnectTo session: UISceneSession, options: UIScene.ConnectionOptions) {
        DispatchQueue.main.async {
            renderAll()
            exit(0)
        }
    }
}

let out = CommandLine.arguments.dropFirst().last ?? "."

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

// snap renders view at 2x into a PNG named file and returns its pixels as RGBA bytes, a row after another.
func snap(_ view: UIView, _ file: String) -> (bytes: [UInt8], width: Int) {
    let format = UIGraphicsImageRendererFormat()
    format.scale = 2
    let image = UIGraphicsImageRenderer(bounds: view.bounds, format: format).image { ctx in view.layer.render(in: ctx.cgContext) }
    try! image.pngData()!.write(to: URL(fileURLWithPath: out + "/" + file))
    let cg = image.cgImage!
    let w = cg.width, h = cg.height
    var bytes = [UInt8](repeating: 0, count: w * h * 4)
    let ctx = CGContext(data: &bytes, width: w, height: h, bitsPerComponent: 8, bytesPerRow: w * 4, space: CGColorSpace(name: CGColorSpace.sRGB)!, bitmapInfo: CGImageAlphaInfo.premultipliedLast.rawValue)!
    ctx.draw(cg, in: CGRect(x: 0, y: 0, width: w, height: h))
    return (bytes, w)
}

func value(_ px: (bytes: [UInt8], width: Int), _ x: Int, _ y: Int) -> String {
    let i = (y * px.width + x) * 4
    let r = px.bytes[i], g = px.bytes[i + 1], b = px.bytes[i + 2]
    return r == g && g == b ? "\(r)" : String(format: "#%02x%02x%02x", r, g, b)
}

// pageControl returns a UIPageControl of n pages with page sel chosen, at its fitting size with its top left corner at x, y.
func pageControl(_ n: Int, _ sel: Int, x: CGFloat, y: CGFloat) -> UIPageControl {
    let pc = UIPageControl()
    pc.numberOfPages = n
    pc.currentPage = sel
    let size = pc.sizeThatFits(CGSize(width: 1000, height: 1000))
    pc.frame = CGRect(x: x, y: y, width: size.width, height: size.height)
    return pc
}

// renderPagers places the page controls where TestScreenshotsCarousel places Breeze UI's pagers and carousel, 352 by 312 points.
func renderPagers(dark: Bool) {
    let mode = dark ? "dark" : "light"
    let root = UIView(frame: CGRect(x: 0, y: 0, width: 352, height: 312))
    root.overrideUserInterfaceStyle = dark ? .dark : .light
    root.backgroundColor = dark ? UIColor(red: 0x1e / 255, green: 0x1e / 255, blue: 0x1e / 255, alpha: 1) : .white
    let slide = UILabel(frame: CGRect(x: 16, y: 96, width: 320, height: 160))
    slide.text = "Slide 2"
    slide.textAlignment = .center
    slide.font = .systemFont(ofSize: 13)
    root.addSubview(slide)
    let five = pageControl(5, 1, x: 16, y: 16), nine = pageControl(9, 6, x: 16, y: 56)
    let four = pageControl(4, 1, x: 0, y: 256)
    four.frame.origin.x = 16 + (320 - four.frame.width) / 2
    for pc in [five, nine, four] { root.addSubview(pc) }
    root.layoutIfNeeded()
    let px = snap(root, "uikit-pagecontrol-\(mode).png")
    print("== uikit-pagecontrol-\(mode).png idiom=\(UIDevice.current.userInterfaceIdiom.rawValue)")
    for (name, pc) in [("five", five), ("nine", nine), ("four", four)] {
        let f = pc.frame
        print("\(name) x=\(f.minX) top=\(f.minY) w=\(f.width) h=\(f.height)")
        func walk(_ v: UIView) {
            if String(describing: type(of: v)).contains("PageIndicatorView") {
                let d = v.convert(v.bounds, to: root)
                print("  dot x=\(d.minX) top=\(d.minY) w=\(d.width) h=\(d.height)")
            }
            v.subviews.forEach(walk)
        }
        walk(pc)
        let y = Int((f.midY) * 2)
        print("  row: " + runs((Int(f.minX * 2)..<Int(f.maxX * 2)).map { value(px, $0, y) }))
    }
}

// PagedTabView is SwiftUI's TabView of the page style, four pages, the 2nd shown, its index shown always.
struct PagedTabView: View {
    @State var page = 1
    var body: some View {
        TabView(selection: $page) {
            ForEach(0..<4) { i in Text("Slide \(i + 1)").tag(i) }
        }
        .tabViewStyle(.page(indexDisplayMode: .always))
        .frame(width: 320, height: 200)
    }
}

// renderTabView renders PagedTabView alone and prints where its UIPageControl stands.
func renderTabView(dark: Bool) {
    let mode = dark ? "dark" : "light"
    let host = UIHostingController(rootView: PagedTabView())
    let window = UIWindow(frame: CGRect(x: 0, y: 0, width: 320, height: 200))
    window.overrideUserInterfaceStyle = dark ? .dark : .light
    window.rootViewController = host
    window.isHidden = false
    host.view.frame = window.bounds
    host.view.backgroundColor = dark ? UIColor(red: 0x1e / 255, green: 0x1e / 255, blue: 0x1e / 255, alpha: 1) : .white
    RunLoop.main.run(until: Date().addingTimeInterval(0.5))
    host.view.layoutIfNeeded()
    _ = snap(host.view, "swiftui-pagetabview-\(mode).png")
    print("== swiftui-pagetabview-\(mode).png")
    func walk(_ v: UIView) {
        if v is UIPageControl {
            let f = v.convert(v.bounds, to: host.view)
            print("UIPageControl x=\(f.minX) top=\(f.minY) w=\(f.width) h=\(f.height)")
        }
        v.subviews.forEach(walk)
    }
    walk(host.view)
}

func renderAll() {
    for dark in [false, true] {
        renderPagers(dark: dark)
        renderTabView(dark: dark)
    }
}

UIApplicationMain(CommandLine.argc, CommandLine.unsafeArgv, nil, NSStringFromClass(AppDelegate.self))
