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
    // InputsScenes are the controls of TestScreenshotsInputs (ui/breeze/widgets_inputs_test.go).
    [Scenes]
    static IEnumerable<Scene> InputsScenes()
    {
        // inputsView of ui/breeze/widgets_inputs_test.go (TestScreenshotsInputs): two PasswordBoxes, empty with a placeholder and holding a password with its RevealButton showing (ShowRevealButton, as the window never takes the focus that shows it), a SplitButton, Buttons holding a 16-DIP FontIcon, the second disabled, and below them a CommandBar holding a SplitButton of SplitButtonCommandBarStyle and two icon-only AppBarButtons, as wide as they need, which puts them at its start as Breeze UI's Toolbar puts its items (a wider CommandBar puts its primary commands at its end).
        yield return new Scene("winui-inputs", 520, 180, Canvas(520, 180, """
            <PasswordBox x:Name="password-empty" Canvas.Left="16" Canvas.Top="16" Width="180" PlaceholderText="Password"/>
            <PasswordBox x:Name="password" Canvas.Left="220" Canvas.Top="16" Width="180" Password="hunter2"/>
            <SplitButton x:Name="split-button" Canvas.Left="16" Canvas.Top="64" Content="Save">
              <SplitButton.Flyout><MenuFlyout><MenuFlyoutItem Text="Save As…"/><MenuFlyoutItem Text="Save All"/></MenuFlyout></SplitButton.Flyout>
            </SplitButton>
            <Button x:Name="icon-button" Canvas.Left="160" Canvas.Top="64" AutomationProperties.Name="Share"><FontIcon Glyph="&#xE72D;" FontSize="16"/></Button>
            <Button x:Name="icon-button-disabled" Canvas.Left="220" Canvas.Top="64" AutomationProperties.Name="Delete" IsEnabled="False"><FontIcon Glyph="&#xE74D;" FontSize="16"/></Button>
            <CommandBar x:Name="toolbar" Canvas.Left="0" Canvas.Top="110" DefaultLabelPosition="Collapsed" OverflowButtonVisibility="Collapsed" IsDynamicOverflowEnabled="False">
              <AppBarElementContainer>
                <SplitButton x:Name="toolbar-split-button" Content="New" Style="{StaticResource SplitButtonCommandBarStyle}">
                  <SplitButton.Flyout><MenuFlyout><MenuFlyoutItem Text="Save As…"/></MenuFlyout></SplitButton.Flyout>
                </SplitButton>
              </AppBarElementContainer>
              <AppBarButton x:Name="toolbar-share" Label="Share" Icon="Share"/>
              <AppBarButton x:Name="toolbar-delete" Label="Delete" Icon="Delete"/>
            </CommandBar>
            """), ["password-empty", "password", "split-button", "icon-button", "icon-button-disabled", "toolbar", "toolbar-split-button", "toolbar-share", "toolbar-delete"], ShowRevealButton);
    }

    // ShowRevealButton puts the PasswordBox named password into its Focused and ButtonVisible states, which draw it as focused and show its RevealButton: WinUI shows the button once text is typed into the focused box, as TestScreenshotsInputs types into Breeze UI's, and the scene's window never takes the focus.
    static async Task ShowRevealButton(FrameworkElement root, ElementTheme theme)
    {
        var box = (Control)root.FindName("password");
        VisualStateManager.GoToState(box, "Focused", false);
        VisualStateManager.GoToState(box, "ButtonVisible", false);
        await Task.Delay(500);
    }
}
