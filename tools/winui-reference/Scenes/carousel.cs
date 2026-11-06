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
    // CarouselScenes are the PipsPagers and the FlipView of TestScreenshotsCarousel (ui/breeze/widgets_carousel_test.go).
    [Scenes]
    static IEnumerable<Scene> CarouselScenes()
    {
        // widgets_carousel_test.go's carouselShotView on Windows (TestScreenshotsCarousel): a PipsPager of 5 pages with the 2nd chosen at 16, 16, one of 9 with the 7th chosen at 16, 56, which shows 5 pips (MaxVisiblePips) scrolled to keep the chosen one in the middle, and at 16, 96 a StackPanel 320 wide of a FlipView 160 high of four pages, each a Grid holding its slide's text centered, the 2nd shown, over a PipsPager of 4 centered 10 below it, as docs/api-specs/PipsPager/PipsPager.md pairs them, where Breeze UI's Carousel draws both. The FlipView's previous and next buttons show only once a mouse moves over it, which the scene's window never sees: winui-carousel-buttons shows them as a move would.
        yield return new Scene("winui-carousel", 352, 312, Canvas(352, 312, CarouselXaml), CarouselItems, WaitForPips);
        yield return new Scene("winui-carousel-buttons", 352, 312, Canvas(352, 312, CarouselXaml), CarouselItems, ShowFlipButtons);

        // The same right to left (TestScreenshotsCarousel's rtl-light shots), mirrored by an inner Canvas of FlowDirection RightToLeft as winui-datepicker-rtl is, without and with the FlipView's buttons.
        yield return new Scene("winui-carousel-rtl", 352, 312, Canvas(352, 312, $$"""
            <Canvas Width="352" Height="312" FlowDirection="RightToLeft">{{CarouselXaml}}</Canvas>
            """), CarouselItems, WaitForPips);
        yield return new Scene("winui-carousel-buttons-rtl", 352, 312, Canvas(352, 312, $$"""
            <Canvas Width="352" Height="312" FlowDirection="RightToLeft">{{CarouselXaml}}</Canvas>
            """), CarouselItems, ShowFlipButtons);

        // The carousel scene's controls in the states a still bitmap does not reach, set by VisualStateManager.GoToState as winui-states sets its controls': the pager of 5's 4th pip PointerOver, the FlipView's buttons shown with the next one PointerOver and the previous one Pressed, and the pager below it with its 3rd pip Pressed. Breeze UI's pip under the pointer, pressed pip and buttons under the pointer and pressed are compared with them one at a time, as one pointer reaches one of them. winui-carousel-disabled has the PipsPagers and the FlipView disabled (IsEnabled False), as Disabled(true) disables Breeze UI's pagers and carousel.
        yield return new Scene("winui-carousel-states", 352, 312, Canvas(352, 312, CarouselXaml), CarouselItems, HoldCarouselStates);
        yield return new Scene("winui-carousel-disabled", 352, 312, Canvas(352, 312, CarouselXaml), CarouselItems, DisableCarousel);
    }

    // CarouselXaml is the body of the carousel scenes (see winui-carousel).
    const string CarouselXaml = """
        <PipsPager x:Name="pager-five" Canvas.Left="16" Canvas.Top="16" NumberOfPages="5" SelectedPageIndex="1"/>
        <PipsPager x:Name="pager-nine" Canvas.Left="16" Canvas.Top="56" NumberOfPages="9"/>
        <StackPanel x:Name="carousel" Canvas.Left="16" Canvas.Top="96" Width="320">
          <FlipView x:Name="flipview" Height="160">
            <Grid><TextBlock Text="Slide 1" HorizontalAlignment="Center" VerticalAlignment="Center"/></Grid>
            <Grid><TextBlock Text="Slide 2" HorizontalAlignment="Center" VerticalAlignment="Center"/></Grid>
            <Grid><TextBlock Text="Slide 3" HorizontalAlignment="Center" VerticalAlignment="Center"/></Grid>
            <Grid><TextBlock Text="Slide 4" HorizontalAlignment="Center" VerticalAlignment="Center"/></Grid>
          </FlipView>
          <PipsPager x:Name="carousel-pager" HorizontalAlignment="Center" Margin="0,10,0,0" NumberOfPages="4" SelectedPageIndex="1"/>
        </StackPanel>
        """;

    // CarouselItems are the elements of the carousel scenes whose frames, parts and pixels the report gives.
    static readonly string[] CarouselItems = ["pager-five", "pager-nine", "carousel", "flipview", "carousel-pager"];

    // WaitForPips shows the FlipView's 2nd page, which XAML cannot choose before the pages are added (a Selector's SelectedIndex past its items is refused), chooses the 7th page of the pager of 9 once it is laid out, and gives the move to it and the PipsPagers' ScrollViewers the time to finish, which they animate (FlipView::OnSelectedIndexChanged animates a move to the adjacent page; ScrollToCenterOfViewport asks StartBringIntoView for an animation). A SelectedPageIndex set in XAML is brought into view only by OnApplyTemplate's UpdateSelectedPip, before the pager's first layout, where StartBringIntoView does nothing: rendered so, the pager of 9 shows its first 5 pips and not the chosen one, which a page chosen later, as an app's or a click's, does not.
    static async Task WaitForPips(FrameworkElement root, ElementTheme theme)
    {
        ((FlipView)root.FindName("flipview")).SelectedIndex = 1;
        ((PipsPager)root.FindName("pager-nine")).SelectedPageIndex = 6;
        await Task.Delay(1000);
    }

    // ShowFlipButtons shows the FlipView's PreviousButtonHorizontal and NextButtonHorizontal, as FlipView::ChangeVisualState does once a mouse moves over it (m_showNavigationButtons), the 2nd page being neither the first nor the last.
    static async Task ShowFlipButtons(FrameworkElement root, ElementTheme theme)
    {
        await WaitForPips(root, theme);
        var flip = (FlipView)root.FindName("flipview");
        foreach (var name in new[] { "PreviousButtonHorizontal", "NextButtonHorizontal" })
        {
            if (FindNamed(flip, name) is not FrameworkElement button)
            {
                throw new InvalidOperationException($"the FlipView has no {name}");
            }
            button.Visibility = Visibility.Visible;
        }
        await Task.Delay(300);
    }

    // HoldCarouselStates shows the FlipView's buttons and puts the parts winui-carousel-states names into their states, without transitions, so their final values show: a pip's PointerOver and Pressed states set its glyph's FontSize and Foreground (PipsPagerButtonBaseStyle), a FlipView button's set its caret's Foreground and its scale (FlipView_themeresources.xaml's PreviousButtonHorizontal and NextButtonHorizontal templates).
    static async Task HoldCarouselStates(FrameworkElement root, ElementTheme theme)
    {
        await ShowFlipButtons(root, theme);
        var flip = (FlipView)root.FindName("flipview");
        foreach (var (control, state) in new (Control?, string)[] {
            (PipAt((PipsPager)root.FindName("pager-five"), 3), "PointerOver"),
            (PipAt((PipsPager)root.FindName("carousel-pager"), 2), "Pressed"),
            (FindNamed(flip, "NextButtonHorizontal") as Control, "PointerOver"),
            (FindNamed(flip, "PreviousButtonHorizontal") as Control, "Pressed"),
        })
        {
            if (control == null)
            {
                throw new InvalidOperationException($"no control to put into {state}");
            }
            VisualStateManager.GoToState(control, state, false);
        }
        await Task.Delay(500);
    }

    // DisableCarousel disables the carousel scene's PipsPagers and FlipView once its pages are chosen, which puts their pips into PipsPagerButtonBaseStyle's Disabled state (ControlStrongFillColorDisabled).
    static async Task DisableCarousel(FrameworkElement root, ElementTheme theme)
    {
        await WaitForPips(root, theme);
        foreach (var name in new[] { "pager-five", "pager-nine", "flipview", "carousel-pager" })
        {
            ((Control)root.FindName(name)).IsEnabled = false;
        }
        await Task.Delay(500);
    }

    // PipAt returns the pip button of page index of pager, which its ItemsRepeater has prepared.
    static Control? PipAt(PipsPager pager, int index) =>
        Find<ItemsRepeater>(pager)?.TryGetElement(index) as Control;

    // FindNamed returns the element named name inside e's visual tree, its template parts included, which FindName on the scene's root does not reach.
    static FrameworkElement? FindNamed(DependencyObject? e, string name)
    {
        if (e == null)
        {
            return null;
        }
        if (e is FrameworkElement f && f.Name == name)
        {
            return f;
        }
        for (int i = 0; i < VisualTreeHelper.GetChildrenCount(e); i++)
        {
            if (FindNamed(VisualTreeHelper.GetChild(e, i), name) is FrameworkElement found)
            {
                return found;
            }
        }
        return null;
    }
}
