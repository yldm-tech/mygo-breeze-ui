import AppKit

// FirstResponderDatePicker takes the window's focus as it joins the window, so that the offscreen render draws its chosen day as a focused NSDatePicker does.
class FirstResponderDatePicker: NSDatePicker {
    override func viewDidMoveToWindow() {
        super.viewDidMoveToWindow()
        window?.makeFirstResponder(self)
    }
}

// renderDatePicker places the date pickers where ui/breeze/widgets_datepicker_test.go's datePickerView places Breeze UI's CalendarView and DatePicker: an NSDatePicker of style clockAndCalendar (SwiftUI's graphical DatePicker on macOS) at 16, 16 showing November 2026 with the 18th chosen, the window's focus when focused, and an NSDatePicker of style textField with presentsCalendarOverlay (the field of SwiftUI's compact style) whose 21-point frame is at 361, 17, so that its bezel, drawn a point outside it, starts at 360, 16; both in American English with weeks starting on Sunday, in UTC. AppKit marks the real today, which TestScreenshotsDatePicker fixes as the 6th of November 2026. It renders them with render and prints the accessibility of the calendar's buttons, which AppKit leaves unnamed.
func renderDatePicker(appearance: NSAppearance.Name, focused: Bool, file: String) {
    var cal = Calendar(identifier: .gregorian)
    cal.locale = Locale(identifier: "en_US"); cal.firstWeekday = 1; cal.timeZone = TimeZone(identifier: "UTC")!
    let day = cal.date(from: DateComponents(year: 2026, month: 11, day: 18, hour: 10))!
    let calendar = focused ? FirstResponderDatePicker() : NSDatePicker()
    let field = NSDatePicker()
    for (p, style) in [(calendar, NSDatePicker.Style.clockAndCalendar), (field, .textField)] {
        p.datePickerStyle = style; p.datePickerElements = .yearMonthDay
        p.locale = cal.locale; p.calendar = cal; p.timeZone = cal.timeZone; p.dateValue = day
    }
    field.presentsCalendarOverlay = true
    render([Item(name: "calendar", view: calendar, x: 16, y: 16), Item(name: "date-field", view: field, x: 361, y: 17)], width: 560, height: 400, appearance: appearance, file: file)
    for child in calendar.accessibilityChildren() ?? [] {
        for button in (child as AnyObject).accessibilityChildren?() ?? [] {
            if let b = button as? NSAccessibilityElementProtocol {
                let o = button as AnyObject
                print("  button label=\"\(o.accessibilityLabel?() ?? "")\" title=\"\(o.accessibilityTitle?() ?? "")\" screen frame=\(b.accessibilityFrame())")
            }
        }
    }
}

// DatePickerScenes renders the graphical and field date pickers of TestScreenshotsDatePicker (datePickerView in ui/breeze/widgets_datepicker_test.go), unfocused and focused.
final class DatePickerScenes: SceneGroup {
    override class var file: String { #fileID }
    override class func renderScenes(appearance: NSAppearance.Name, mode: String) {
        renderDatePicker(appearance: appearance, focused: false, file: "appkit-datepicker-\(mode).png"); renderDatePicker(appearance: appearance, focused: true, file: "appkit-datepicker-focused-\(mode).png")
    }
}
