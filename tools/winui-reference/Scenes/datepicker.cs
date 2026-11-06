using System.Runtime.InteropServices;
using System.Runtime.InteropServices.WindowsRuntime;
using System.Text;
using Microsoft.UI.Dispatching;
using Microsoft.UI.Windowing;
using Microsoft.UI.Xaml;
using Microsoft.UI.Xaml.Controls;
using Microsoft.UI.Xaml.Controls.Primitives;
using Microsoft.UI.Xaml.Markup;
using Microsoft.UI.Xaml.Media;
using Microsoft.UI.Xaml.Media.Imaging;
using Microsoft.UI.Xaml.XamlTypeInfo;
using Windows.Foundation;
using Windows.Graphics;
using Windows.Graphics.Imaging;
using Windows.Storage.Streams;
using Windows.UI;

namespace WinUIReference;

partial class App
{
    // DatePickerScenes are the calendars and date pickers of TestScreenshotsDatePicker and TestScreenshotsFluentDatePickerOpen (ui/breeze/widgets_datepicker_test.go).
    [Scenes]
    static IEnumerable<Scene> DatePickerScenes()
    {
        // widgets_datepicker_test.go's datePickerView on Windows: a CalendarView at 16, 16 showing November 2026 with the 18th chosen, and at 360, 16 a CalendarDatePicker of the same day, closed, both in American English with weeks starting on Sunday. XAML cannot set a CalendarView's chosen dates or the month it shows, so the scene's hook adds the 18th to SelectedDates, shows its month with SetDisplayDate and sets the picker's Date. WinUI marks the machine's today, which TestScreenshotsDatePicker fixes as the 6th of November 2026.
        yield return new Scene("winui-datepicker", 560, 400, Canvas(560, 400, """
            <CalendarView x:Name="calendar" Canvas.Left="16" Canvas.Top="16" Language="en-US" FirstDayOfWeek="Sunday" SelectionMode="Single"/>
            <CalendarDatePicker x:Name="date-picker" Canvas.Left="360" Canvas.Top="16" Language="en-US" FirstDayOfWeek="Sunday"/>
            """), ["calendar", "date-picker"], async (root, theme) =>
        {
            var day = new DateTimeOffset(2026, 11, 18, 10, 0, 0, TimeSpan.Zero);
            var calendar = (CalendarView)root.FindName("calendar");
            calendar.SelectedDates.Add(day);
            calendar.SetDisplayDate(day);
            ((CalendarDatePicker)root.FindName("date-picker")).Date = day;
            await Task.Delay(800);
        });

        // The same right to left (TestScreenshotsDatePicker's rtl-light shot): an inner Canvas of FlowDirection RightToLeft, which measures Canvas.Left from its right side, mirrors the scene inside the root, whose own ScaleTransform a mirrored root would carry off the bitmap.
        yield return new Scene("winui-datepicker-rtl", 560, 400, Canvas(560, 400, """
            <Canvas Width="560" Height="400" FlowDirection="RightToLeft">
              <CalendarView x:Name="calendar" Canvas.Left="16" Canvas.Top="16" Language="en-US" FirstDayOfWeek="Sunday" SelectionMode="Single"/>
              <CalendarDatePicker x:Name="date-picker" Canvas.Left="360" Canvas.Top="16" Language="en-US" FirstDayOfWeek="Sunday"/>
            </Canvas>
            """), ["calendar", "date-picker"], async (root, theme) =>
        {
            var day = new DateTimeOffset(2026, 11, 18, 10, 0, 0, TimeSpan.Zero);
            var calendar = (CalendarView)root.FindName("calendar");
            calendar.SelectedDates.Add(day);
            calendar.SetDisplayDate(day);
            ((CalendarDatePicker)root.FindName("date-picker")).Date = day;
            await Task.Delay(800);
        });

        // widgets_datepicker_test.go's datePickerOpenView on Windows (TestScreenshotsFluentDatePickerOpen): a CalendarDatePicker of the 18th of November 2026 at 120, 16 with its calendar open, which OpenCalendar puts where WinUI's flyout shows it.
        yield return new Scene("winui-datepicker-open", 420, 440, Canvas(420, 440, """
            <CalendarDatePicker x:Name="date-picker" Canvas.Left="120" Canvas.Top="16" Language="en-US" FirstDayOfWeek="Sunday"/>
            """), ["date-picker", "flyout-calendar"], OpenCalendar);
    }

    // OpenCalendar sets the CalendarDatePicker named date-picker to the 18th of November 2026 and opens its calendar, to read where WinUI places the flyout's FlyoutPresenter, which Placement Bottom centers below the picker, and how its template dresses the presenter and the CalendarView, then closes it and puts a FlyoutPresenter of the same style, shadow and offset from the picker's bottom middle there in the scene, holding a CalendarView of the same style, corners, border and margin showing the same month and day, as RenderTargetBitmap does not draw popups and the scene's scale does not reach them.
    static async Task OpenCalendar(FrameworkElement root, ElementTheme theme)
    {
        var day = new DateTimeOffset(2026, 11, 18, 10, 0, 0, TimeSpan.Zero);
        var picker = (CalendarDatePicker)root.FindName("date-picker");
        picker.Date = day;
        await Task.Delay(500);
        // The template's Flyout (FlyoutBase.AttachedFlyout of its root) would open in a window of its own on the screen, out of the off-screen window; kept in the window's root it opens there, at the same place, as the root has room around the picker.
        if (VisualTreeHelper.GetChild(picker, 0) is FrameworkElement templateRoot && FlyoutBase.GetAttachedFlyout(templateRoot) is FlyoutBase attached)
        {
            attached.ShouldConstrainToRootBounds = true;
        }
        else
        {
            throw new InvalidOperationException("the picker's template has no attached flyout");
        }
        picker.IsCalendarOpen = true;
        await Task.Delay(1000);
        FlyoutPresenter? shown = null;
        Popup? holder = null;
        foreach (var popup in VisualTreeHelper.GetOpenPopupsForXamlRoot(root.XamlRoot))
        {
            if (shown == null && Find<FlyoutPresenter>(popup.Child) is FlyoutPresenter found)
            {
                shown = found;
                holder = popup;
            }
        }
        if (shown == null || holder == null || Find<CalendarView>(shown) is not CalendarView inner)
        {
            throw new InvalidOperationException("the calendar did not open");
        }
        // The popup is not scaled with the scene: its offset from the picker's bottom middle, in the window's DIPs, is the offset in the scene's.
        var target = picker.TransformToVisual(null).TransformPoint(new Point(picker.ActualWidth / 2, picker.ActualHeight));
        var opened = shown.TransformToVisual(null).TransformPoint(new Point(0, 0));
        double w = shown.ActualWidth, h = shown.ActualHeight;
        double dx = opened.X + w / 2 - target.X, dy = opened.Y - target.Y;
        var below = picker.TransformToVisual(root).TransformPoint(new Point(picker.ActualWidth / 2, picker.ActualHeight));
        var at = new Point(below.X - w / 2 + dx, below.Y + dy);
        var innerAt = inner.TransformToVisual(shown).TransformPoint(new Point(0, 0));
        var Ts = (Thickness t) => $"{t.Left},{t.Top},{t.Right},{t.Bottom}";
        var Cs = (CornerRadius r) => $"{r.TopLeft},{r.TopRight},{r.BottomRight},{r.BottomLeft}";
        var Bs = (Brush b) => b is SolidColorBrush s ? $"#{s.Color.A:x2}{s.Color.R:x2}{s.Color.G:x2}{s.Color.B:x2}" : b?.GetType().Name ?? "null";
        Report?.AppendLine($"opened flyout ({theme}): {w}x{h} at {dx},{dy} from the picker's bottom middle (window {opened.X},{opened.Y}); presenter padding {Ts(shown.Padding)} border {Ts(shown.BorderThickness)} {Bs(shown.BorderBrush)} corners {Cs(shown.CornerRadius)} background {Bs(shown.Background)} translation {shown.Translation} shadow {shown.Shadow?.GetType().Name ?? "none"} (popup's translation {holder.Translation} shadow {holder.Shadow?.GetType().Name ?? "none"}, its child's {holder.Child.Translation} {holder.Child.Shadow?.GetType().Name ?? "none"}) style {(shown.Style != null ? "set" : "none")}; calendar {inner.ActualWidth}x{inner.ActualHeight} at {innerAt.X},{innerAt.Y} margin {Ts(inner.Margin)} border {Ts(inner.BorderThickness)} {Bs(inner.BorderBrush)} corners {Cs(inner.CornerRadius)} background {Bs(inner.Background)}");
        picker.IsCalendarOpen = false;
        await Task.Delay(500);
        var calendar = new CalendarView
        {
            Name = "flyout-calendar",
            Language = "en-US",
            FirstDayOfWeek = Windows.Globalization.DayOfWeek.Sunday,
            SelectionMode = inner.SelectionMode,
            CornerRadius = inner.CornerRadius,
            BorderThickness = inner.BorderThickness,
            Margin = inner.Margin,
        };
        if (inner.Style != null)
        {
            calendar.Style = inner.Style;
        }
        calendar.SelectedDates.Add(day);
        calendar.SetDisplayDate(day);
        // Sized by its content, as in the popup: a width set on it would limit the CalendarView's, and WinUI would wrap the day names that do not fit their columns rather than widen the columns.
        var presenter = new FlyoutPresenter { Name = "flyout", Content = calendar, Translation = shown.Translation, Shadow = shown.Shadow };
        if (shown.Style != null)
        {
            presenter.Style = shown.Style;
        }
        Microsoft.UI.Xaml.Controls.Canvas.SetLeft(presenter, at.X);
        Microsoft.UI.Xaml.Controls.Canvas.SetTop(presenter, at.Y);
        ((Panel)root).Children.Add(presenter);
        await WhenLoaded(presenter);
        await Task.Delay(800);
    }
}
