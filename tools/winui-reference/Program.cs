// WinUI 3's own controls, for comparing Breeze UI's Windows look with them.
//
// It renders WinUI's controls where ui/breeze/screenshot_test.go places Breeze UI's for Windows (referenceView, referenceView2, referenceView3 and the alert) and where ui/breeze/widgets_text_test.go places Breeze UI's Text (textView), a set of the same controls held in their PointerOver and Pressed visual states, and the controls of TestScreenshotsInputs (password boxes, split buttons and icon buttons, alone and in a CommandBar), the flat buttons of TestScreenshotsFlat (Buttons of SubtleButtonStyle and checked ToggleButtons, some held in their PointerOver and Pressed states), the InfoBars and determinate ProgressRings of widgets_status_test.go (statusView, statusTitlesView's InfoBars without a message and statusActionsView's with an ActionButton), and CommunityToolkit's SettingsCards and SettingsExpander where widgets_rows_test.go places Breeze UI's rows (rowsView), and SettingsCards holding a ComboBox, TextBoxes and NumberBoxes and SettingsExpanders where widgets_settings_test.go places Breeze UI's settings rows (settingsView), and wider ones that keep their content at the end, TabViews where widgets_tabview_test.go places Breeze UI's (tabViewView) and one whose tabs fill it (TestFluentTabViewAsMeasured) TitleBars where widgets_headerbar_test.go places Breeze UI's HeaderBar (headerView), SelectorBars where widgets_switcher_test.go places Breeze UI's ViewSwitcher (switcherView), and a CalendarView and a CalendarDatePicker where widgets_datepicker_test.go places Breeze UI's CalendarView and DatePicker (datePickerView), the same right to left, and a CalendarDatePicker with its calendar open (datePickerOpenView), in the light and the dark theme at 2x, writes them as PNGs (winui-light.png, winui2-dark.png, ...) and writes report.txt: for every control its frame in DIPs, the frames of its named template parts and text blocks, and the pixels along its middle row and column at 2x as runs (value x count), a gray pixel as its value and a colored one as #rrggbb, as tools/appkit-reference/reference.swift prints AppKit's. The accent is forced to Windows's default blue: SystemAccentColor #0078D4 and its shades (Light1-3, Dark1-3), the values Windows derives from it and ui/breeze/fluent.go cites, are set in the application's resources, so the result does not depend on the machine's accent. Text is WinUI's default, ContentControlThemeFontFamily, which is Segoe UI Variable on Windows 11.
//
// Each scene is shown in a borderless tool window far off the screen with WS_EX_NOACTIVATE and AppWindow.Show(false), so it never takes the focus or the keyboard, and the program exits by itself, within 55 seconds at most. The window's rasterization scale is the screen's, 1 on the machine it was written on, and RenderTargetBitmap.RenderAsync asked for a larger size only stretches the bitmap, so each scene is scaled by a ScaleTransform of 2 with its RasterizationScale set to 2, which makes WinUI draw its shapes and text at 2x; layout still rounds to whole DIPs, where at 200% it rounds to half DIPs. Popups are not drawn by RenderTargetBitmap and not reached by the scene's scale, so the ContentDialog is shown in place (ContentDialogPlacement.InPlace), which leaves out its smoke layer, and a flyout is opened once to read where WinUI places it, then closed and a FlyoutPresenter with the same content put there in the scene (CalendarDatePicker's own Flyout kept inside the window first, which would otherwise open in a window of its own on the screen). ThemeShadow's shadows do render; the window's backdrop does not. The visual states of the fifth set are set with VisualStateManager.GoToState without transitions, so their final values show.
//
// It is an unpackaged app (WindowsPackageType None) that carries the Windows App SDK with it (WindowsAppSDKSelfContained), built with the .NET 8 SDK or later on Windows 10 version 2004 (build 19041) or later, x64:
//
//	dotnet build -c Release -m:4 tools/winui-reference/WinUIReference.csproj
//	tools\winui-reference\bin\Release\net8.0-windows10.0.26100.0\win-x64\winui-reference.exe <dir> [scene,...]
//
// With a second argument it renders only the scenes named there, by their file names without the theme (winui-rows,winui-text).
//
// It must run in an interactive desktop session (not over SSH, which cannot show windows); started from a remote shell, a scheduled task with /IT runs it in the signed-in user's session:
//
//	schtasks /Create /TN winui-reference /TR "<exe> <dir>" /SC ONCE /ST 23:59 /IT /F
//	schtasks /Run /TN winui-reference
//	(wait for <dir>\done.txt)
//	schtasks /Delete /TN winui-reference /F
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

public static class Program
{
    [STAThread]
    static int Main(string[] args)
    {
        // Never outlive a minute, whatever happens.
        using var watchdog = new Timer(_ => Environment.Exit(3), null, TimeSpan.FromSeconds(55), Timeout.InfiniteTimeSpan);
        WinRT.ComWrappersSupport.InitializeComWrappers();
        Application.Start(p =>
        {
            SynchronizationContext.SetSynchronizationContext(new DispatcherQueueSynchronizationContext(DispatcherQueue.GetForCurrentThread()));
            _ = new App(args.Length > 0 ? args[0] : ".", args.Length > 1 ? args[1].Split(',') : []);
        });
        return App.ExitCode;
    }
}

// A scene: controls placed on a canvas W by H DIPs, each named for the report, and what to open over them.
record Scene(string File, int W, int H, string Xaml, string[] Items, Func<FrameworkElement, ElementTheme, Task>? Open = null);

public sealed class App : Application, IXamlMetadataProvider
{
    public static int ExitCode;
    const int Scale = 2;
    readonly string outDir;
    readonly string[] only;
    readonly StringBuilder report = new();
    // Report is the report of the running App, for the scenes' hooks to add what they read.
    static StringBuilder? Report;
    readonly List<IXamlMetadataProvider> providers = [new XamlControlsXamlMetaDataProvider()];

    // App renders the scenes into dir, only those whose file names are in only when it is not empty. The XAML type providers of the CommunityToolkit assemblies next to the program (SettingsControls and what it uses) answer for their types after WinUI's own.
    public App(string dir, string[] only)
    {
        outDir = dir;
        this.only = only;
        Report = report;
        foreach (var path in Directory.GetFiles(AppContext.BaseDirectory, "CommunityToolkit.WinUI*.dll"))
        {
            Type?[] types;
            try
            {
                types = System.Reflection.Assembly.LoadFrom(path).GetTypes();
            }
            catch (System.Reflection.ReflectionTypeLoadException ex)
            {
                types = ex.Types;
            }
            foreach (var t in types)
            {
                if (t != null && typeof(IXamlMetadataProvider).IsAssignableFrom(t) && t.IsClass && !t.IsAbstract && t.GetConstructor(Type.EmptyTypes) != null)
                {
                    providers.Add((IXamlMetadataProvider)Activator.CreateInstance(t)!);
                }
            }
        }
    }

    public IXamlType GetXamlType(Type type) => providers.Select(p => p.GetXamlType(type)).FirstOrDefault(x => x != null)!;
    public IXamlType GetXamlType(string fullName) => providers.Select(p => p.GetXamlType(fullName)).FirstOrDefault(x => x != null)!;
    public XmlnsDefinition[] GetXmlnsDefinitions() => providers.SelectMany(p => p.GetXmlnsDefinitions()).ToArray();

    // Windows's default accent and the shades Windows derives from it.
    static readonly (string Key, string Hex)[] Accent =
    [
        ("SystemAccentColor", "#0078D4"),
        ("SystemAccentColorLight1", "#429CE3"), ("SystemAccentColorLight2", "#76B9ED"), ("SystemAccentColorLight3", "#A6D8FF"),
        ("SystemAccentColorDark1", "#005A9E"), ("SystemAccentColorDark2", "#004275"), ("SystemAccentColorDark3", "#002642"),
    ];

    static Color Hex(string s) => Color.FromArgb(255, Convert.ToByte(s[1..3], 16), Convert.ToByte(s[3..5], 16), Convert.ToByte(s[5..7], 16));

    protected override async void OnLaunched(LaunchActivatedEventArgs e)
    {
        Resources.MergedDictionaries.Add(new XamlControlsResources());
        foreach (var (key, hex) in Accent)
        {
            Resources[key] = Hex(hex);
        }
        try
        {
            Directory.CreateDirectory(outDir);
            var xaml = System.Diagnostics.FileVersionInfo.GetVersionInfo(Path.Combine(AppContext.BaseDirectory, "Microsoft.UI.Xaml.dll"));
            var sdk = System.Diagnostics.FileVersionInfo.GetVersionInfo(typeof(Microsoft.Windows.ApplicationModel.WindowsAppRuntime.RuntimeInfo).Assembly.Location);
            report.AppendLine($"Windows App SDK {sdk.ProductVersion}, Microsoft.UI.Xaml.dll {xaml.FileVersion}, Windows {Environment.OSVersion.Version}");
            var window = MakeWindow();
            foreach (var scene in Scenes().Where(s => only.Length == 0 || only.Contains(s.File)))
            {
                foreach (var theme in new[] { ElementTheme.Light, ElementTheme.Dark })
                {
                    await Render(window, scene, theme);
                }
            }
        }
        catch (Exception ex)
        {
            report.AppendLine("error: " + ex);
            ExitCode = 1;
        }
        File.WriteAllText(Path.Combine(outDir, "report.txt"), report.ToString());
        File.WriteAllText(Path.Combine(outDir, "done.txt"), ExitCode.ToString());
        Exit();
    }

    const string Ns = "xmlns=\"http://schemas.microsoft.com/winfx/2006/xaml/presentation\" xmlns:x=\"http://schemas.microsoft.com/winfx/2006/xaml\" xmlns:tk=\"using:CommunityToolkit.WinUI.Controls\"";

    static string Canvas(int w, int h, string body) =>
        $"<Canvas {Ns} Width=\"{w}\" Height=\"{h}\" Background=\"{{ThemeResource SolidBackgroundFillColorBaseBrush}}\">{body}</Canvas>";

    static IEnumerable<Scene> Scenes()
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

        // The first row of flatView in ui/breeze/widgets_flat_test.go (TestScreenshotsFlat): Buttons of SubtleButtonStyle holding a 16-DIP FontIcon, the second held in PointerOver and the third disabled, one holding text, and for the toggles, which WinUI has no subtle style for, a subtle Button and a checked ToggleButton; the row below holds subtle Buttons in Pressed and a checked ToggleButton in CheckedPointerOver.
        yield return new Scene("winui-flat", 360, 100, Canvas(360, 100, """
            <Button x:Name="flat-image" Canvas.Left="16" Canvas.Top="16" Style="{StaticResource SubtleButtonStyle}" AutomationProperties.Name="Share"><FontIcon Glyph="&#xE72D;" FontSize="16"/></Button>
            <Button x:Name="flat-image-hover" Tag="PointerOver" Canvas.Left="64" Canvas.Top="16" Style="{StaticResource SubtleButtonStyle}" AutomationProperties.Name="Delete"><FontIcon Glyph="&#xE74D;" FontSize="16"/></Button>
            <Button x:Name="flat-image-disabled" Canvas.Left="112" Canvas.Top="16" Style="{StaticResource SubtleButtonStyle}" AutomationProperties.Name="Delete" IsEnabled="False"><FontIcon Glyph="&#xE74D;" FontSize="16"/></Button>
            <Button x:Name="flat-text" Canvas.Left="160" Canvas.Top="16" Style="{StaticResource SubtleButtonStyle}" Content="Skip"/>
            <Button x:Name="flat-toggle-off" Canvas.Left="240" Canvas.Top="16" Style="{StaticResource SubtleButtonStyle}" AutomationProperties.Name="Sidebar"><FontIcon Glyph="&#xE8A0;" FontSize="16"/></Button>
            <ToggleButton x:Name="flat-toggle-on" Canvas.Left="288" Canvas.Top="16" IsChecked="True" AutomationProperties.Name="Panel"><FontIcon Glyph="&#xE8A0;" FontSize="16"/></ToggleButton>
            <Button x:Name="flat-image-pressed" Tag="Pressed" Canvas.Left="16" Canvas.Top="60" Style="{StaticResource SubtleButtonStyle}" AutomationProperties.Name="Share"><FontIcon Glyph="&#xE72D;" FontSize="16"/></Button>
            <Button x:Name="flat-text-pressed" Tag="Pressed" Canvas.Left="64" Canvas.Top="60" Style="{StaticResource SubtleButtonStyle}" Content="Skip"/>
            <ToggleButton x:Name="flat-toggle-on-hover" Tag="CheckedPointerOver" Canvas.Left="144" Canvas.Top="60" IsChecked="True" AutomationProperties.Name="Panel"><FontIcon Glyph="&#xE8A0;" FontSize="16"/></ToggleButton>
            """), ["flat-image", "flat-image-hover", "flat-image-disabled", "flat-text", "flat-toggle-off", "flat-toggle-on", "flat-image-pressed", "flat-text-pressed", "flat-toggle-on-hover"], HoldStates);

        // widgets_status_test.go's statusView: four InfoBars, one of each severity, the first and the last closable, 16 DIPs from the window's edges and 68 apart, and determinate ProgressRings of 25, 60, 100 and 0 percent below them and one of 60 percent at 16 DIPs. Breeze UI's StatusPage has no WinUI counterpart.
        yield return new Scene("winui-status", 560, 360, Canvas(560, 360, """
            <InfoBar x:Name="bar-informational" Canvas.Left="16" Canvas.Top="16" Width="528" IsOpen="True" Severity="Informational" Title="Informational" Message="Some message about it."/>
            <InfoBar x:Name="bar-success" Canvas.Left="16" Canvas.Top="84" Width="528" IsOpen="True" IsClosable="False" Severity="Success" Title="Success" Message="Some message about it."/>
            <InfoBar x:Name="bar-warning" Canvas.Left="16" Canvas.Top="152" Width="528" IsOpen="True" IsClosable="False" Severity="Warning" Title="Warning" Message="Some message about it."/>
            <InfoBar x:Name="bar-error" Canvas.Left="16" Canvas.Top="220" Width="528" IsOpen="True" Severity="Error" Title="Error" Message="Some message about it."/>
            <ProgressRing x:Name="ring-25" Canvas.Left="16" Canvas.Top="300" IsIndeterminate="False" Value="25"/>
            <ProgressRing x:Name="ring-60" Canvas.Left="64" Canvas.Top="300" IsIndeterminate="False" Value="60"/>
            <ProgressRing x:Name="ring-100" Canvas.Left="112" Canvas.Top="300" IsIndeterminate="False" Value="100"/>
            <ProgressRing x:Name="ring-0" Canvas.Left="160" Canvas.Top="300" IsIndeterminate="False" Value="0"/>
            <ProgressRing x:Name="ring-small" Canvas.Left="208" Canvas.Top="308" Width="16" Height="16" IsIndeterminate="False" Value="60"/>
            """), ["bar-informational", "bar-success", "bar-warning", "bar-error", "ring-25", "ring-60", "ring-100", "ring-0", "ring-small"]);

        // widgets_status_test.go's statusTitlesView: the four InfoBars of statusView with a title and no message, which InfoBarPanel lays out as it lays a message below the title.
        yield return new Scene("winui-status-titles", 560, 300, Canvas(560, 300, """
            <InfoBar x:Name="title-informational" Canvas.Left="16" Canvas.Top="16" Width="528" IsOpen="True" Severity="Informational" Title="Informational"/>
            <InfoBar x:Name="title-success" Canvas.Left="16" Canvas.Top="84" Width="528" IsOpen="True" IsClosable="False" Severity="Success" Title="Success"/>
            <InfoBar x:Name="title-warning" Canvas.Left="16" Canvas.Top="152" Width="528" IsOpen="True" IsClosable="False" Severity="Warning" Title="Warning"/>
            <InfoBar x:Name="title-error" Canvas.Left="16" Canvas.Top="220" Width="528" IsOpen="True" Severity="Error" Title="Error"/>
            """), ["title-informational", "title-success", "title-warning", "title-error"]);

        // widgets_status_test.go's statusActionsView: InfoBars with an ActionButton, 16 DIPs from the window's edges and 68 apart: a closable warning with a Button, an informational bar with an accent Button, a closable success bar without a message with a Button, and a closable error bar 280 wide, which InfoBarPanel lays out vertically with the button below the message.
        yield return new Scene("winui-status-actions", 560, 380, Canvas(560, 380, """
            <InfoBar x:Name="actions-warning" Canvas.Left="16" Canvas.Top="16" Width="528" IsOpen="True" Severity="Warning" Title="Warning" Message="Some message about it."><InfoBar.ActionButton><Button Content="Open Settings"/></InfoBar.ActionButton></InfoBar>
            <InfoBar x:Name="actions-informational" Canvas.Left="16" Canvas.Top="84" Width="528" IsOpen="True" IsClosable="False" Severity="Informational" Title="Informational" Message="Some message about it."><InfoBar.ActionButton><Button Content="Restart" Style="{ThemeResource AccentButtonStyle}"/></InfoBar.ActionButton></InfoBar>
            <InfoBar x:Name="actions-success" Canvas.Left="16" Canvas.Top="152" Width="528" IsOpen="True" Severity="Success" Title="Success"><InfoBar.ActionButton><Button Content="Undo"/></InfoBar.ActionButton></InfoBar>
            <InfoBar x:Name="actions-error" Canvas.Left="16" Canvas.Top="220" Width="280" IsOpen="True" Severity="Error" Title="Error" Message="Some message about it."><InfoBar.ActionButton><Button Content="Retry"/></InfoBar.ActionButton></InfoBar>
            """), ["actions-warning", "actions-informational", "actions-success", "actions-error"]);

        // widgets_rows_test.go's rowsView on Windows, laid out as it lays Breeze UI's rows out, in two columns 300 wide, 16 from the window's edges and 24 apart, their parts 16 apart: a card of text (a Border with the card's resources, as an Expander's content and a SettingsCard draw it), a card of rows (SettingsCards with DefaultSettingsExpanderItemStyle's setters over the card's fill, padded 16 at the sides where the style's 58,8,44,8 makes room for an expander's icon and chevron, the first without its line), a SettingsCard holding a secondary text and a disabled one holding a ToggleSwitch; beside them a section header (BodyStrongTextBlockStyle) 10 above two SettingsCards 4 apart, as Windows Settings stacks them, and a SettingsExpander holding two items. The ToggleSwitches take the card's RightAlignedCompactToggleSwitchStyle. The cards are 300 wide, which SettingsCard's ContentAlignmentStates wrap (RightWrapped, between SettingsCardWrapNoIconThreshold 286 and SettingsCardWrapThreshold 476).
        yield return new Scene("winui-rows", 672, 600, Canvas(672, 600, """
            <StackPanel Canvas.Left="16" Canvas.Top="16" Width="300" Spacing="16">
              <Border x:Name="card-text" Padding="16" CornerRadius="{ThemeResource ControlCornerRadius}" BorderThickness="1" Background="{ThemeResource CardBackgroundFillColorDefaultBrush}" BorderBrush="{ThemeResource CardStrokeColorDefaultBrush}">
                <TextBlock Text="A card holds what belongs together."/>
              </Border>
              <Border x:Name="card-rows" CornerRadius="{ThemeResource ControlCornerRadius}" BorderThickness="1" Background="{ThemeResource CardBackgroundFillColorDefaultBrush}" BorderBrush="{ThemeResource CardStrokeColorDefaultBrush}">
                <StackPanel>
                  <tk:SettingsCard x:Name="row-wallpaper" Header="Wallpaper" Description="Choose a picture" MinHeight="52" Padding="16,8,16,8" CornerRadius="0" BorderThickness="0" Background="Transparent"><Button Content="Open"/></tk:SettingsCard>
                  <tk:SettingsCard x:Name="row-dark" Header="Dark mode" Description="Use a dark appearance" MinHeight="52" Padding="16,8,16,8" CornerRadius="0" BorderThickness="0,1,0,0" Background="Transparent"><ToggleSwitch IsOn="True"/></tk:SettingsCard>
                  <tk:SettingsCard x:Name="row-notifications" Header="Notifications" MinHeight="52" Padding="16,8,16,8" CornerRadius="0" BorderThickness="0,1,0,0" Background="Transparent"><ToggleSwitch/></tk:SettingsCard>
                </StackPanel>
              </Border>
              <tk:SettingsCard x:Name="row-storage" Header="Storage"><TextBlock Text="120 GB" Foreground="{ThemeResource TextFillColorSecondaryBrush}"/></tk:SettingsCard>
              <tk:SettingsCard x:Name="row-bluetooth" Header="Bluetooth" IsEnabled="False"><ToggleSwitch IsOn="True"/></tk:SettingsCard>
            </StackPanel>
            <StackPanel Canvas.Left="340" Canvas.Top="16" Width="300" Spacing="16">
              <StackPanel>
                <TextBlock x:Name="privacy" Text="Privacy" Style="{StaticResource BodyStrongTextBlockStyle}" Margin="0,0,0,10"/>
                <StackPanel x:Name="privacy-cards" Spacing="4">
                  <tk:SettingsCard x:Name="row-location" Header="Location" Description="Apps may ask for it"><ToggleSwitch IsOn="True"/></tk:SettingsCard>
                  <tk:SettingsCard x:Name="row-camera" Header="Camera" Description="2 apps"/>
                </StackPanel>
              </StackPanel>
              <tk:SettingsExpander x:Name="expander" Header="Network" IsExpanded="True">
                <tk:SettingsExpander.Items>
                  <tk:SettingsCard x:Name="row-proxy" Header="Proxy" Description="Automatic"/>
                  <tk:SettingsCard x:Name="row-metered" Header="Metered"><ToggleSwitch/></tk:SettingsCard>
                </tk:SettingsExpander.Items>
              </tk:SettingsExpander>
            </StackPanel>
            """), ["card-text", "card-rows", "row-wallpaper", "row-dark", "row-notifications", "row-storage", "row-bluetooth", "privacy", "privacy-cards", "row-location", "row-camera", "expander", "row-proxy", "row-metered"]);

        // SettingsCards wider than SettingsCardWrapThreshold (476), which keep their content at the end: one with a header, a description and a button, one with a ToggleSwitch that is on, and one with a header alone and a ToggleSwitch that is off, 4 apart as Windows Settings stacks them; Breeze UI's tests read the same layout from rows this wide.
        yield return new Scene("winui-rows-wide", 560, 260, Canvas(560, 260, """
            <StackPanel x:Name="wide-cards" Canvas.Left="16" Canvas.Top="16" Width="528" Spacing="4">
              <tk:SettingsCard x:Name="wide-action" Header="Wallpaper" Description="Choose a picture"><Button Content="Open"/></tk:SettingsCard>
              <tk:SettingsCard x:Name="wide-switch-on" Header="Dark mode" Description="Use a dark appearance"><ToggleSwitch IsOn="True"/></tk:SettingsCard>
              <tk:SettingsCard x:Name="wide-switch-off" Header="Notifications"><ToggleSwitch/></tk:SettingsCard>
            </StackPanel>
            """), ["wide-cards", "wide-action", "wide-switch-on", "wide-switch-off"]);

        // widgets_settings_test.go's settingsView on Windows (TestScreenshotsSettings), laid out as rowsView's scene is, in two columns 300 wide, 16 from the window's edges and 24 apart, their parts 16 apart: a card of rows (SettingsCards with DefaultSettingsExpanderItemStyle's setters over the card's fill, padded 16 at the sides, the first without its line) holding a ComboBox, a TextBox holding a name, an empty TextBox and a NumberBox with inline spin buttons, then a SettingsCard holding a ComboBox and a disabled one holding a NumberBox; beside them a SettingsExpander, open, with a description, holding two items, and a closed one. SettingsCard gives the ComboBox and the TextBox a MinWidth of 120; the NumberBox takes Breeze UI's 156, its 80-DIP text and the 76 DIPs of the spin buttons' column (NumberBox.xaml: UpSpinButton's Margin of 4 and DownSpinButton's 0,4,4,4 around two 32-DIP buttons).
        yield return new Scene("winui-settings", 672, 640, Canvas(672, 640, """
            <StackPanel Canvas.Left="16" Canvas.Top="16" Width="300" Spacing="16">
              <Border x:Name="settings-card" CornerRadius="{ThemeResource ControlCornerRadius}" BorderThickness="1" Background="{ThemeResource CardBackgroundFillColorDefaultBrush}" BorderBrush="{ThemeResource CardStrokeColorDefaultBrush}">
                <StackPanel>
                  <tk:SettingsCard x:Name="combo-row" Header="Appearance" Description="Light or dark" MinHeight="52" Padding="16,8,16,8" CornerRadius="0" BorderThickness="0" Background="Transparent">
                    <ComboBox SelectedIndex="2"><x:String>Light</x:String><x:String>Dark</x:String><x:String>Auto</x:String></ComboBox>
                  </tk:SettingsCard>
                  <tk:SettingsCard x:Name="entry-row" Header="Name" MinHeight="52" Padding="16,8,16,8" CornerRadius="0" BorderThickness="0,1,0,0" Background="Transparent"><TextBox Text="Evan" Width="120"/></tk:SettingsCard>
                  <tk:SettingsCard x:Name="entry-row-empty" Header="Email" MinHeight="52" Padding="16,8,16,8" CornerRadius="0" BorderThickness="0,1,0,0" Background="Transparent"><TextBox Width="120"/></tk:SettingsCard>
                  <tk:SettingsCard x:Name="spin-row" Header="Font size" Description="Points" MinHeight="52" Padding="16,8,16,8" CornerRadius="0" BorderThickness="0,1,0,0" Background="Transparent"><NumberBox Value="12" Minimum="6" Maximum="72" SpinButtonPlacementMode="Inline" Width="156"/></tk:SettingsCard>
                </StackPanel>
              </Border>
              <tk:SettingsCard x:Name="combo-row-alone" Header="Language">
                <ComboBox SelectedIndex="0"><x:String>English</x:String><x:String>Français</x:String></ComboBox>
              </tk:SettingsCard>
              <tk:SettingsCard x:Name="spin-row-disabled" Header="Copies" IsEnabled="False"><NumberBox Value="3" Minimum="1" Maximum="9" SpinButtonPlacementMode="Inline" Width="156"/></tk:SettingsCard>
            </StackPanel>
            <StackPanel Canvas.Left="340" Canvas.Top="16" Width="300" Spacing="16">
              <tk:SettingsExpander x:Name="expander-row" Header="Network" Description="Wi-Fi and proxies" IsExpanded="True">
                <tk:SettingsExpander.Items>
                  <tk:SettingsCard x:Name="expander-proxy" Header="Proxy" Description="Automatic"/>
                  <tk:SettingsCard x:Name="expander-metered" Header="Metered"><ToggleSwitch/></tk:SettingsCard>
                </tk:SettingsExpander.Items>
              </tk:SettingsExpander>
              <tk:SettingsExpander x:Name="expander-row-closed" Header="Advanced">
                <tk:SettingsExpander.Items>
                  <tk:SettingsCard Header="Logs"/>
                </tk:SettingsExpander.Items>
              </tk:SettingsExpander>
            </StackPanel>
            """), ["settings-card", "combo-row", "entry-row", "entry-row-empty", "spin-row", "combo-row-alone", "spin-row-disabled", "expander-row", "expander-proxy", "expander-metered", "expander-row-closed"]);

        // widgets_switcher_test.go's switcherView on Windows (TestScreenshotsViewSwitcher): a TitleBar along the top of the window titled Music with a SelectorBar as its RightHeader, where Breeze UI's HeaderBar holds a ViewSwitcher among its end controls, and a SelectorBar alone 16 DIPs in from the window's side and the 48-DIP bar's bottom, each of three SelectorBarItems with a FontIcon (FontIcon's own FontSize of 20, which the template scales by 0.8) and a Text, Browse selected. The TitleBar is 422 DIPs wide, the window's 560 less the 138 that MyGo's three 46-DIP caption buttons take (headerTestRoom), as in the headerbar scenes.
        yield return new Scene("winui-switcher", 560, 140, Canvas(560, 140, """
            <TitleBar x:Name="titlebar" Canvas.Left="0" Canvas.Top="0" Width="422" Title="Music">
              <TitleBar.RightHeader>
                <SelectorBar x:Name="switcher-bar">
                  <SelectorBarItem x:Name="bar-library" Text="Library"><SelectorBarItem.Icon><FontIcon Glyph="&#xE8F1;"/></SelectorBarItem.Icon></SelectorBarItem>
                  <SelectorBarItem x:Name="bar-browse" Text="Browse" IsSelected="True"><SelectorBarItem.Icon><FontIcon Glyph="&#xF0E2;"/></SelectorBarItem.Icon></SelectorBarItem>
                  <SelectorBarItem x:Name="bar-radio" Text="Radio"><SelectorBarItem.Icon><FontIcon Glyph="&#xEC05;"/></SelectorBarItem.Icon></SelectorBarItem>
                </SelectorBar>
              </TitleBar.RightHeader>
            </TitleBar>
            <SelectorBar x:Name="switcher" Canvas.Left="16" Canvas.Top="64">
              <SelectorBarItem x:Name="library" Text="Library"><SelectorBarItem.Icon><FontIcon Glyph="&#xE8F1;"/></SelectorBarItem.Icon></SelectorBarItem>
              <SelectorBarItem x:Name="browse" Text="Browse" IsSelected="True"><SelectorBarItem.Icon><FontIcon Glyph="&#xF0E2;"/></SelectorBarItem.Icon></SelectorBarItem>
              <SelectorBarItem x:Name="radio" Text="Radio"><SelectorBarItem.Icon><FontIcon Glyph="&#xEC05;"/></SelectorBarItem.Icon></SelectorBarItem>
            </SelectorBar>
            """), ["titlebar", "switcher-bar", "bar-library", "bar-browse", "bar-radio", "switcher", "library", "browse", "radio"]);

        // widgets_tabview_test.go's tabViewView on Windows: a TabView 528 wide at 16, 16 holding Documents (chosen), Downloads and Pictures with its add button, the chosen tab's text 12 below its 40-DIP strip, and at 16, 136 a TabView of General and Advanced (chosen) whose tabs do not close, without the add button; TabWidthMode Equal and CloseButtonOverlayMode Auto, the defaults.
        yield return new Scene("winui-tabview", 560, 260, Canvas(560, 260, """
            <TabView x:Name="tabview-docs" Canvas.Left="16" Canvas.Top="16" Width="528" SelectedIndex="0">
              <TabViewItem x:Name="tab-documents" Header="Documents"/>
              <TabViewItem x:Name="tab-downloads" Header="Downloads"/>
              <TabViewItem x:Name="tab-pictures" Header="Pictures"/>
            </TabView>
            <TextBlock x:Name="tab-content" Canvas.Left="16" Canvas.Top="68" Text="Showing Documents"/>
            <TabView x:Name="tabview-prefs" Canvas.Left="16" Canvas.Top="136" Width="528" SelectedIndex="1" IsAddTabButtonVisible="False">
              <TabViewItem x:Name="tab-general" Header="General" IsClosable="False"/>
              <TabViewItem x:Name="tab-advanced" Header="Advanced" IsClosable="False"/>
            </TabView>
            """), ["tabview-docs", "tab-documents", "tab-downloads", "tab-pictures", "tab-content", "tabview-prefs", "tab-general", "tab-advanced"]);

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

        // fluent_native_test.go's TestFluentTabViewAsMeasured, its third TabView: 528 wide at 16, 16, five tabs that do not close (A1 chosen), without the add button, so that the tabs fill the strip.
        yield return new Scene("winui-tabview-fill", 560, 80, Canvas(560, 80, """
            <TabView x:Name="tabview-fill" Canvas.Left="16" Canvas.Top="16" Width="528" SelectedIndex="0" IsAddTabButtonVisible="False">
              <TabViewItem x:Name="tab-a1" Header="A1" IsClosable="False"/>
              <TabViewItem Header="B2" IsClosable="False"/>
              <TabViewItem Header="C3" IsClosable="False"/>
              <TabViewItem Header="D4" IsClosable="False"/>
              <TabViewItem x:Name="tab-e5" Header="E5" IsClosable="False"/>
            </TabView>
            """), ["tabview-fill", "tab-a1", "tab-e5"]);

        // widgets_headerbar_test.go's headerView (TestScreenshotsHeaderBar): a TitleBar along the top of the window titled Documents, with the Subtitle 12 items in the second scene, its own back button (IsBackButtonVisible: PART_BackButton of TitleBarBackButtonStyle in its column, stretched to the bar less its Margin of 2, as Breeze UI draws an image button at the bar's start; a Button of that style in the LeftHeader is only as tall as its 16-DIP icon, since the header's presenter is centered), and as its RightHeader a Button reading Edit and a subtle Button 40 by 44 DIPs holding the Share glyph with TitleBarBackButtonStyle's Margin of 2, 4 apart, where Breeze UI's HeaderBar holds an image button, a push button and an image button, and below the bar the window's content, a TextBlock 16 DIPs in from the window's side and the bar's bottom. The TitleBar is 422 DIPs wide, the window's 560 less the 138 that MyGo's three 46-DIP caption buttons take (headerTestRoom), as the scene's borderless window has no caption buttons of its own.
        yield return new Scene("winui-headerbar", 560, 120, Canvas(560, 120, """
            <TitleBar x:Name="titlebar" Canvas.Left="0" Canvas.Top="0" Width="422" Title="Documents" IsBackButtonVisible="True">
              <TitleBar.RightHeader>
                <StackPanel Orientation="Horizontal" Spacing="4" VerticalAlignment="Stretch">
                  <Button x:Name="edit" Content="Edit" VerticalAlignment="Center"/>
                  <Button x:Name="share" Width="40" Height="44" Margin="2" Padding="0" BorderThickness="0" CornerRadius="{ThemeResource ControlCornerRadius}" Background="{ThemeResource SubtleFillColorTransparentBrush}" AutomationProperties.Name="Share"><FontIcon Glyph="&#xE72D;" FontSize="16"/></Button>
                </StackPanel>
              </TitleBar.RightHeader>
            </TitleBar>
            <TextBlock x:Name="content" Canvas.Left="16" Canvas.Top="64" Text="Content"/>
            """), ["titlebar", "edit", "share", "content"]);
        yield return new Scene("winui-headerbar-subtitle", 560, 120, Canvas(560, 120, """
            <TitleBar x:Name="titlebar" Canvas.Left="0" Canvas.Top="0" Width="422" Title="Documents" Subtitle="12 items" IsBackButtonVisible="True">
              <TitleBar.RightHeader>
                <StackPanel Orientation="Horizontal" Spacing="4" VerticalAlignment="Stretch">
                  <Button x:Name="edit" Content="Edit" VerticalAlignment="Center"/>
                  <Button x:Name="share" Width="40" Height="44" Margin="2" Padding="0" BorderThickness="0" CornerRadius="{ThemeResource ControlCornerRadius}" Background="{ThemeResource SubtleFillColorTransparentBrush}" AutomationProperties.Name="Share"><FontIcon Glyph="&#xE72D;" FontSize="16"/></Button>
                </StackPanel>
              </TitleBar.RightHeader>
            </TitleBar>
            <TextBlock x:Name="content" Canvas.Left="16" Canvas.Top="64" Text="Content"/>
            """), ["titlebar", "edit", "share", "content"]);

        // TestScreenshotsText's textView: the TextBlock styles of Breeze UI's Text steps where textTops places them, the last two in TextFillColorSecondaryBrush.
        yield return new Scene("winui-text", 360, 300, Canvas(360, 300, """
            <TextBlock x:Name="large-title" Canvas.Left="16" Canvas.Top="16" Text="Large Title" Style="{StaticResource TitleLargeTextBlockStyle}"/>
            <TextBlock x:Name="title" Canvas.Left="16" Canvas.Top="76" Text="Title" Style="{StaticResource TitleTextBlockStyle}"/>
            <TextBlock x:Name="subtitle" Canvas.Left="16" Canvas.Top="122" Text="Subtitle" Style="{StaticResource SubtitleTextBlockStyle}"/>
            <TextBlock x:Name="heading" Canvas.Left="16" Canvas.Top="160" Text="Heading" Style="{StaticResource BodyStrongTextBlockStyle}"/>
            <TextBlock x:Name="body" Canvas.Left="16" Canvas.Top="188" Text="Body" Style="{StaticResource BodyTextBlockStyle}"/>
            <TextBlock x:Name="caption" Canvas.Left="16" Canvas.Top="214" Text="Caption" Style="{StaticResource CaptionTextBlockStyle}"/>
            <TextBlock x:Name="secondary-body" Canvas.Left="16" Canvas.Top="240" Text="Secondary body" Style="{StaticResource BodyTextBlockStyle}" Foreground="{ThemeResource TextFillColorSecondaryBrush}"/>
            <TextBlock x:Name="secondary-caption" Canvas.Left="16" Canvas.Top="266" Text="Secondary caption" Style="{StaticResource CaptionTextBlockStyle}" Foreground="{ThemeResource TextFillColorSecondaryBrush}"/>
            """), ["large-title", "title", "subtitle", "heading", "body", "caption", "secondary-body", "secondary-caption"]);
    }

    // ShowRevealButton puts the PasswordBox named password into its Focused and ButtonVisible states, which draw it as focused and show its RevealButton: WinUI shows the button once text is typed into the focused box, as TestScreenshotsInputs types into Breeze UI's, and the scene's window never takes the focus.
    static async Task ShowRevealButton(FrameworkElement root, ElementTheme theme)
    {
        var box = (Control)root.FindName("password");
        VisualStateManager.GoToState(box, "Focused", false);
        VisualStateManager.GoToState(box, "ButtonVisible", false);
        await Task.Delay(500);
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

    static T? Find<T>(DependencyObject? e) where T : DependencyObject
    {
        if (e == null)
        {
            return null;
        }
        if (e is T found)
        {
            return found;
        }
        for (int i = 0; i < VisualTreeHelper.GetChildrenCount(e); i++)
        {
            if (Find<T>(VisualTreeHelper.GetChild(e, i)) is T t)
            {
                return t;
            }
        }
        return null;
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

    // MakeWindow shows an empty borderless tool window far off every screen without activating it, large enough for every scene at 2x.
    static Window MakeWindow()
    {
        var window = new Window();
        var hwnd = WinRT.Interop.WindowNative.GetWindowHandle(window);
        SetWindowLongPtr(hwnd, GwlExStyle, GetWindowLongPtr(hwnd, GwlExStyle) | WsExNoActivate | WsExToolWindow);
        var appWindow = window.AppWindow;
        if (appWindow.Presenter is OverlappedPresenter presenter)
        {
            presenter.SetBorderAndTitleBar(false, false);
            presenter.IsResizable = false;
            presenter.IsMaximizable = false;
            presenter.IsMinimizable = false;
        }
        appWindow.IsShownInSwitchers = false;
        appWindow.MoveAndResize(new RectInt32(-8000, -8000, 1400, 1300));
        appWindow.Show(false);
        return window;
    }

    // Render draws scene in theme at 2x. The window's rasterization scale is the screen's, 1 here, and RenderTargetBitmap.RenderAsync with a larger size only stretches that bitmap, so the scene is scaled by a ScaleTransform and its RasterizationScale set to match, which makes WinUI draw its shapes and text at 2x.
    async Task Render(Window window, Scene scene, ElementTheme theme)
    {
        var root = (FrameworkElement)XamlReader.Load(scene.Xaml);
        root.RequestedTheme = theme;
        root.HorizontalAlignment = HorizontalAlignment.Left;
        root.VerticalAlignment = VerticalAlignment.Top;
        root.RenderTransform = new ScaleTransform { ScaleX = Scale, ScaleY = Scale };
        root.RasterizationScale = Scale;
        int pw = scene.W * Scale, ph = scene.H * Scale;
        var host = new Grid { Width = pw, Height = ph, HorizontalAlignment = HorizontalAlignment.Left, VerticalAlignment = VerticalAlignment.Top };
        host.Children.Add(root);
        window.Content = host;
        await WhenLoaded(root);
        await Task.Delay(700);
        if (scene.Open != null)
        {
            await scene.Open(root, theme);
        }
        var mode = theme == ElementTheme.Dark ? "dark" : "light";
        var pixels = await Pixels(host, pw, ph);
        await WritePng(Path.Combine(outDir, $"{scene.File}-{mode}.png"), pixels, pw, ph);
        report.AppendLine($"== {scene.File}-{mode}.png (window rasterization scale {root.XamlRoot.RasterizationScale}, scene {root.RasterizationScale})");
        var items = scene.Items.ToList();
        if (root.FindName("flyout") != null)
        {
            items.Add("flyout");
        }
        foreach (var name in items)
        {
            if (root.FindName(name) is not FrameworkElement e)
            {
                report.AppendLine($"{name} not found");
                continue;
            }
            var p = e.TransformToVisual(root).TransformPoint(new Point(0, 0));
            report.AppendLine($"{name} x={p.X} top={p.Y} w={e.ActualWidth} h={e.ActualHeight}");
            int x0 = (int)Math.Round(p.X * Scale), y0 = (int)Math.Round(p.Y * Scale), w = (int)Math.Round(e.ActualWidth * Scale), h = (int)Math.Round(e.ActualHeight * Scale);
            report.AppendLine("  row: " + Runs(Enumerable.Range(Math.Max(0, x0 - 4), Math.Min(pw, x0 + w + 4) - Math.Max(0, x0 - 4)).Select(x => Value(pixels, pw, x, Math.Min(ph - 1, y0 + h / 2)))));
            report.AppendLine("  col: " + Runs(Enumerable.Range(Math.Max(0, y0 - 4), Math.Min(ph, y0 + h + 4) - Math.Max(0, y0 - 4)).Select(y => Value(pixels, pw, Math.Min(pw - 1, x0 + w / 2), y))));
            Parts(e, root, "  ", 0);
        }
        File.WriteAllText(Path.Combine(outDir, "report.txt"), report.ToString());
    }

    // Parts reports the named, visible template parts and the text blocks inside e, with their frames relative to root.
    void Parts(DependencyObject e, UIElement root, string indent, int depth)
    {
        int n = VisualTreeHelper.GetChildrenCount(e);
        for (int i = 0; i < n; i++)
        {
            var child = VisualTreeHelper.GetChild(e, i);
            if (child is FrameworkElement f && f.Visibility == Visibility.Visible && f.ActualWidth > 0 && f.ActualHeight > 0)
            {
                string label = f is TextBlock tb ? $"TextBlock \"{tb.Text}\" {tb.FontFamily?.Source} {tb.FontSize} {tb.FontWeight.Weight}" : f.Name;
                if (!string.IsNullOrEmpty(label))
                {
                    var p = f.TransformToVisual(root).TransformPoint(new Point(0, 0));
                    report.AppendLine($"{indent}{new string(' ', depth)}{label} ({f.GetType().Name}) x={p.X} top={p.Y} w={f.ActualWidth} h={f.ActualHeight} opacity={f.Opacity}");
                }
                Parts(child, root, indent, depth + 1);
            }
            else if (child is not FrameworkElement)
            {
                Parts(child, root, indent, depth + 1);
            }
        }
    }

    static Task WhenLoaded(FrameworkElement e)
    {
        if (e.IsLoaded)
        {
            return Task.CompletedTask;
        }
        var done = new TaskCompletionSource();
        e.Loaded += (_, _) => done.TrySetResult();
        return done.Task;
    }

    // Pixels renders e at w by h pixels: premultiplied BGRA rows.
    static async Task<byte[]> Pixels(UIElement e, int w, int h)
    {
        var bitmap = new RenderTargetBitmap();
        await bitmap.RenderAsync(e, w, h);
        if (bitmap.PixelWidth != w || bitmap.PixelHeight != h)
        {
            throw new InvalidOperationException($"rendered {bitmap.PixelWidth}x{bitmap.PixelHeight}, asked for {w}x{h}");
        }
        return (await bitmap.GetPixelsAsync()).ToArray();
    }

    static string Value(byte[] px, int w, int x, int y)
    {
        int i = (y * w + x) * 4;
        byte b = px[i], g = px[i + 1], r = px[i + 2];
        return r == g && g == b ? r.ToString() : $"#{r:x2}{g:x2}{b:x2}";
    }

    static string Runs(IEnumerable<string> values)
    {
        var v = values.ToList();
        var s = new List<string>();
        for (int i = 0; i < v.Count;)
        {
            int j = i;
            while (j < v.Count && v[j] == v[i])
            {
                j++;
            }
            s.Add(j - i > 1 ? $"{v[i]}x{j - i}" : v[i]);
            i = j;
        }
        return string.Join(" ", s);
    }

    static async Task WritePng(string path, byte[] pixels, int w, int h)
    {
        using var stream = new InMemoryRandomAccessStream();
        var encoder = await BitmapEncoder.CreateAsync(BitmapEncoder.PngEncoderId, stream);
        encoder.SetPixelData(BitmapPixelFormat.Bgra8, BitmapAlphaMode.Premultiplied, (uint)w, (uint)h, 96 * Scale, 96 * Scale, pixels);
        await encoder.FlushAsync();
        stream.Seek(0);
        using var file = File.Create(path);
        await stream.AsStreamForRead().CopyToAsync(file);
    }

    const int GwlExStyle = -20;
    const nint WsExNoActivate = 0x08000000, WsExToolWindow = 0x00000080;

    [DllImport("user32.dll", EntryPoint = "GetWindowLongPtrW")]
    static extern nint GetWindowLongPtr(nint hwnd, int index);

    [DllImport("user32.dll", EntryPoint = "SetWindowLongPtrW")]
    static extern nint SetWindowLongPtr(nint hwnd, int index, nint value);
}
