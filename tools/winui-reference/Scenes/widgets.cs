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
    // WidgetsScenes are the controls where ui/breeze/screenshot_test.go places Breeze UI's for Windows (TestScreenshots: referenceView, referenceView2, referenceView3 and the alert), and the same controls held in their PointerOver and Pressed visual states (TestScreenshotsWindowsStates).
    [Scenes]
    static IEnumerable<Scene> WidgetsScenes()
    {
        // referenceView: three buttons, three check boxes and three switches. Breeze UI's switch shows no On or Off content.
        yield return new Scene("winui", 360, 200, Canvas(360, 200, """
            <Button x:Name="button" Canvas.Left="16" Canvas.Top="16" Content="Cancel"/>
            <Button x:Name="accent-button" Canvas.Left="110" Canvas.Top="16" Content="Save" Style="{StaticResource AccentButtonStyle}"/>
            <Button x:Name="disabled-button" Canvas.Left="200" Canvas.Top="16" Content="Disabled" IsEnabled="False"/>
            <CheckBox x:Name="checkbox-off" Canvas.Left="16" Canvas.Top="64" Content="Unchecked"/>
            <CheckBox x:Name="checkbox-on" Canvas.Left="130" Canvas.Top="64" Content="Checked" IsChecked="True"/>
            <CheckBox x:Name="checkbox-disabled" Canvas.Left="240" Canvas.Top="64" Content="Disabled" IsChecked="True" IsEnabled="False"/>
            <ToggleSwitch x:Name="switch-off" Canvas.Left="16" Canvas.Top="106" OnContent="" OffContent=""/>
            <ToggleSwitch x:Name="switch-on" Canvas.Left="100" Canvas.Top="106" OnContent="" OffContent="" IsOn="True"/>
            <ToggleSwitch x:Name="switch-disabled" Canvas.Left="184" Canvas.Top="106" OnContent="" OffContent="" IsEnabled="False"/>
            """), ["button", "accent-button", "disabled-button", "checkbox-off", "checkbox-on", "checkbox-disabled", "switch-off", "switch-on", "switch-disabled"]);

        // referenceView2: radio buttons, fields, a ComboBox, a toggle, a slider, a progress bar, a SelectorBar (Segmented), a ProgressRing and a multi-line TextBox.
        yield return new Scene("winui2", 520, 330, Canvas(520, 330, """
            <RadioButton x:Name="radio-off" Canvas.Left="16" Canvas.Top="16" Content="Off" GroupName="size"/>
            <RadioButton x:Name="radio-on" Canvas.Left="100" Canvas.Top="16" Content="On" GroupName="size" IsChecked="True"/>
            <TextBox x:Name="text-field" Canvas.Left="16" Canvas.Top="50" Width="180" PlaceholderText="Name"/>
            <AutoSuggestBox x:Name="search-field" Canvas.Left="220" Canvas.Top="50" Width="180" PlaceholderText="Search" QueryIcon="Find"/>
            <ComboBox x:Name="combo-box" Canvas.Left="16" Canvas.Top="90" SelectedIndex="0">
              <x:String>Medium</x:String><x:String>Small</x:String><x:String>Large</x:String>
            </ComboBox>
            <ToggleButton x:Name="toggle" Canvas.Left="220" Canvas.Top="90" Content="Bold" IsChecked="True"/>
            <Slider x:Name="slider" Canvas.Left="16" Canvas.Top="130" Width="180" Minimum="0" Maximum="100" Value="40"/>
            <ProgressBar x:Name="progress" Canvas.Left="220" Canvas.Top="142" Width="180" Minimum="0" Maximum="100" Value="40"/>
            <SelectorBar x:Name="selector-bar" Canvas.Left="16" Canvas.Top="170">
              <SelectorBarItem Text="Day" IsSelected="True"/><SelectorBarItem Text="Week"/><SelectorBarItem Text="Month"/>
            </SelectorBar>
            <ProgressRing x:Name="progress-ring" Canvas.Left="220" Canvas.Top="170" IsActive="True"/>
            <TextBox x:Name="text-area" Canvas.Left="16" Canvas.Top="220" Width="180" Height="70" AcceptsReturn="True" TextWrapping="Wrap" Text="Notes"/>
            """), ["radio-off", "radio-on", "text-field", "search-field", "combo-box", "toggle", "slider", "progress", "selector-bar", "progress-ring", "text-area"]);

        // referenceView3: a hyperlink in text, a divider (DividerStrokeColorDefault), two ToggleButtons 8 apart (ToggleGroup), a DropDownButton (MenuButton) and a button with its Flyout open below it.
        yield return new Scene("winui3", 400, 240, Canvas(400, 240, """
            <TextBlock x:Name="link" Canvas.Left="16" Canvas.Top="16"><Hyperlink NavigateUri="https://example.com">Docs</Hyperlink></TextBlock>
            <Rectangle x:Name="divider" Canvas.Left="16" Canvas.Top="50" Width="200" Height="1" Fill="{ThemeResource DividerStrokeColorDefaultBrush}"/>
            <StackPanel x:Name="toggle-group" Canvas.Left="16" Canvas.Top="70" Orientation="Horizontal" Spacing="8">
              <ToggleButton Content="Bold" IsChecked="True"/><ToggleButton Content="Italic"/>
            </StackPanel>
            <DropDownButton x:Name="menu-button" Canvas.Left="16" Canvas.Top="120" Content="Sort">
              <DropDownButton.Flyout><MenuFlyout><MenuFlyoutItem Text="Name"/></MenuFlyout></DropDownButton.Flyout>
            </DropDownButton>
            <Button x:Name="info" Canvas.Left="200" Canvas.Top="70" Content="Info"/>
            """), ["link", "divider", "toggle-group", "menu-button", "info"], OpenFlyout);

        // The alert: a ContentDialog with a title, a message, the default (primary) button and the close button, over a window showing a line of text 16 DIPs in.
        yield return new Scene("winui-alert", 560, 340, $$"""
            <Grid {{Ns}} Width="560" Height="340" Background="{ThemeResource SolidBackgroundFillColorBaseBrush}">
              <TextBlock x:Name="content" Margin="16" Text="Window content"/>
              <ContentDialog x:Name="dialog" Title="Delete “Notes”?" Content="You can’t undo this action." PrimaryButtonText="Delete" CloseButtonText="Cancel" DefaultButton="Primary"/>
            </Grid>
            """, ["content", "dialog"], OpenDialog);

        // The controls of the first two sets held in their PointerOver and Pressed states, by VisualStateManager.GoToState, which the pointer reaches but a still bitmap does not.
        yield return new Scene("winui-states", 520, 260, Canvas(520, 260, """
            <Button x:Name="button-hover" Tag="PointerOver" Canvas.Left="16" Canvas.Top="16" Content="Cancel"/>
            <Button x:Name="button-pressed" Tag="Pressed" Canvas.Left="110" Canvas.Top="16" Content="Cancel"/>
            <Button x:Name="accent-hover" Tag="PointerOver" Canvas.Left="204" Canvas.Top="16" Content="Save" Style="{StaticResource AccentButtonStyle}"/>
            <Button x:Name="accent-pressed" Tag="Pressed" Canvas.Left="284" Canvas.Top="16" Content="Save" Style="{StaticResource AccentButtonStyle}"/>
            <CheckBox x:Name="checkbox-off-hover" Tag="UncheckedPointerOver" Canvas.Left="16" Canvas.Top="64" Content="Off"/>
            <CheckBox x:Name="checkbox-off-pressed" Tag="UncheckedPressed" Canvas.Left="140" Canvas.Top="64" Content="Off"/>
            <CheckBox x:Name="checkbox-on-hover" Tag="CheckedPointerOver" Canvas.Left="264" Canvas.Top="64" Content="On" IsChecked="True"/>
            <CheckBox x:Name="checkbox-on-pressed" Tag="CheckedPressed" Canvas.Left="388" Canvas.Top="64" Content="On" IsChecked="True"/>
            <ToggleSwitch x:Name="switch-off-hover" Tag="PointerOver" Canvas.Left="16" Canvas.Top="106" OnContent="" OffContent=""/>
            <ToggleSwitch x:Name="switch-off-pressed" Tag="Pressed" Canvas.Left="140" Canvas.Top="106" OnContent="" OffContent=""/>
            <ToggleSwitch x:Name="switch-on-hover" Tag="PointerOver" Canvas.Left="264" Canvas.Top="106" OnContent="" OffContent="" IsOn="True"/>
            <ToggleSwitch x:Name="switch-on-pressed" Tag="Pressed" Canvas.Left="388" Canvas.Top="106" OnContent="" OffContent="" IsOn="True"/>
            <RadioButton x:Name="radio-off-pressed" Tag="Pressed" Canvas.Left="16" Canvas.Top="150" Content="Off" GroupName="a"/>
            <RadioButton x:Name="radio-on-hover" Tag="PointerOver" Canvas.Left="140" Canvas.Top="150" Content="On" GroupName="b" IsChecked="True"/>
            <RadioButton x:Name="radio-on-pressed" Tag="Pressed" Canvas.Left="264" Canvas.Top="150" Content="On" GroupName="c" IsChecked="True"/>
            <Slider x:Name="slider-hover" Tag="PointerOver" Canvas.Left="16" Canvas.Top="194" Width="180" Minimum="0" Maximum="100" Value="40"/>
            <Slider x:Name="slider-pressed" Tag="Pressed" Canvas.Left="220" Canvas.Top="194" Width="180" Minimum="0" Maximum="100" Value="40"/>
            """), ["button-hover", "button-pressed", "accent-hover", "accent-pressed", "checkbox-off-hover", "checkbox-off-pressed", "checkbox-on-hover", "checkbox-on-pressed", "switch-off-hover", "switch-off-pressed", "switch-on-hover", "switch-on-pressed", "radio-off-pressed", "radio-on-hover", "radio-on-pressed", "slider-hover", "slider-pressed"], HoldStates);
    }

    // OpenFlyout opens the button's Flyout below it, aligned to its start, as Breeze UI's popover is, to read where WinUI places the FlyoutPresenter, then closes it and puts a FlyoutPresenter with the same content there in the scene, as RenderTargetBitmap does not draw popups and the scene's scale does not reach them.
    static async Task OpenFlyout(FrameworkElement root, ElementTheme theme)
    {
        var info = (Button)root.FindName("info");
        var flyout = new Flyout { Content = new TextBlock { Text = "Details" }, ShouldConstrainToRootBounds = true };
        flyout.ShowAt(info, new FlyoutShowOptions { Placement = FlyoutPlacementMode.BottomEdgeAlignedLeft, ShowMode = FlyoutShowMode.Transient });
        await Task.Delay(1000);
        FlyoutPresenter? shown = null;
        foreach (var popup in VisualTreeHelper.GetOpenPopupsForXamlRoot(root.XamlRoot))
        {
            shown ??= Find<FlyoutPresenter>(popup.Child);
        }
        if (shown == null)
        {
            throw new InvalidOperationException("the flyout did not open");
        }
        // The popup is not scaled with the scene: its offset from the button's bottom start, in the window's DIPs, is the offset in the scene's.
        var target = info.TransformToVisual(null).TransformPoint(new Point(0, info.ActualHeight));
        var opened = shown.TransformToVisual(null).TransformPoint(new Point(0, 0));
        var below = info.TransformToVisual(root).TransformPoint(new Point(0, info.ActualHeight));
        var at = new Point(below.X + opened.X - target.X, below.Y + opened.Y - target.Y);
        double w = shown.ActualWidth, h = shown.ActualHeight;
        flyout.Hide();
        var presenter = new FlyoutPresenter { Name = "flyout", Content = new TextBlock { Text = "Details" }, Width = w, Height = h };
        Microsoft.UI.Xaml.Controls.Canvas.SetLeft(presenter, at.X);
        Microsoft.UI.Xaml.Controls.Canvas.SetTop(presenter, at.Y);
        ((Panel)root).Children.Add(presenter);
        await WhenLoaded(presenter);
        await Task.Delay(500);
    }

    // OpenDialog shows the scene's ContentDialog in place, inside the scene, rather than in a popup, so that the scene's scale reaches it.
    static async Task OpenDialog(FrameworkElement root, ElementTheme theme)
    {
        var dialog = (ContentDialog)root.FindName("dialog");
        _ = dialog.ShowAsync(ContentDialogPlacement.InPlace);
        await Task.Delay(1500);
    }

    // HoldStates puts each control into the visual state its Tag names.
    static async Task HoldStates(FrameworkElement root, ElementTheme theme)
    {
        foreach (var child in ((Panel)root).Children)
        {
            if (child is Control control && control.Tag is string state)
            {
                VisualStateManager.GoToState(control, state, false);
                // A slider's thumb is a Thumb with visual states of its own, which scale its inner circle.
                if (Find<Thumb>(control) is Thumb thumb)
                {
                    VisualStateManager.GoToState(thumb, state, false);
                }
            }
        }
        await Task.Delay(800);
    }
}
