import AppKit

// SlidePages is the delegate of renderCarousel's NSPageController in its custom mode: every page a view controller of one identifier whose view is a label, centered, showing the page's object, "Slide 1" to "Slide 4".
class SlidePages: NSObject, NSPageControllerDelegate {
    func pageController(_ pc: NSPageController, identifierFor object: Any) -> NSPageController.ObjectIdentifier { "slide" }
    func pageController(_ pc: NSPageController, viewControllerForIdentifier id: NSPageController.ObjectIdentifier) -> NSViewController {
        let vc = NSViewController()
        let label = NSTextField(labelWithString: "")
        label.alignment = .center
        label.translatesAutoresizingMaskIntoConstraints = false
        let view = NSView(frame: NSRect(x: 0, y: 0, width: 320, height: 160))
        view.addSubview(label)
        NSLayoutConstraint.activate([label.centerXAnchor.constraint(equalTo: view.centerXAnchor), label.centerYAnchor.constraint(equalTo: view.centerYAnchor)])
        vc.view = view
        return vc
    }
    func pageController(_ pc: NSPageController, prepare vc: NSViewController, with object: Any?) {
        (vc.view.subviews.first as? NSTextField)?.stringValue = object as? String ?? ""
    }
}

// carouselPages keeps renderCarousel's page controller and its delegate alive while it renders.
var carouselPages: [AnyObject] = []

// renderCarousel places what AppKit has of a carousel where ui/breeze/widgets_carousel_test.go's carouselShotView places Breeze UI's Carousel on macOS: an NSPageController of transitionStyle horizontalStrip, the controller that swipes between pages on macOS, its view 320 by 160 at 16, 96 showing the second of four pages, "Slide 2" centered. AppKit has no page indicator, and SwiftUI's PageTabViewStyle and IndexViewStyle are unavailable on macOS: the dots that Breeze UI draws below the pages and in its pagers at 16, 16 and 16, 56 are UIKit for Mac's UIPageControl, which pagecontrol.swift renders where the test places them. It renders the page controller's view with render.
func renderCarousel(appearance: NSAppearance.Name, file: String) {
    let pages = SlidePages()
    let pc = NSPageController()
    pc.view = NSView(frame: NSRect(x: 0, y: 0, width: 320, height: 160))
    pc.delegate = pages
    pc.transitionStyle = .horizontalStrip
    pc.arrangedObjects = ["Slide 1", "Slide 2", "Slide 3", "Slide 4"]
    pc.selectedIndex = 1
    carouselPages = [pc, pages]
    render([Item(name: "page-controller", view: pc.view, x: 16, y: 96, size: NSSize(width: 320, height: 160))], width: 352, height: 312, appearance: appearance, file: file)
    print("  selectedIndex=\(pc.selectedIndex) selectedViewController=\(pc.selectedViewController.map { String(describing: type(of: $0)) } ?? "none")")
}

// CarouselScenes renders an NSPageController's strip of pages where TestScreenshotsCarousel (carouselShotView in ui/breeze/widgets_carousel_test.go) places Breeze UI's Carousel; the page dots, which AppKit lacks, are UIKit for Mac's, which tools/appkit-reference/pagecontrol.swift, a Mac Catalyst app built on its own outside this package, renders.
final class CarouselScenes: SceneGroup {
    override class var file: String { #fileID }
    override class func renderScenes(appearance: NSAppearance.Name, mode: String) {
        renderCarousel(appearance: appearance, file: "appkit-carousel-\(mode).png")
    }
}
