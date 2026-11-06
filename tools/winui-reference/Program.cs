// WinUI 3's own controls, for comparing Breeze UI's Windows look with them.
//
// It renders WinUI's controls where Breeze UI's screenshot tests place Breeze UI's for Windows, each group of scenes (Scenes/<group>.cs, a method of App marked [Scenes]) where its group's test places them, in the light and the dark theme at 2x, writes them as PNGs (winui-light.png, winui2-dark.png, ...) and writes report.txt: for every control its frame in DIPs, the frames of its named template parts and text blocks, and the pixels along its middle row and column at 2x as runs (value x count), a gray pixel as its value and a colored one as #rrggbb, as tools/appkit-reference prints AppKit's. The accent is forced to Windows's default blue: SystemAccentColor #0078D4 and its shades (Light1-3, Dark1-3), the values Windows derives from it and ui/breeze/fluent.go cites, are set in the application's resources, so the result does not depend on the machine's accent. Text is WinUI's default, ContentControlThemeFontFamily, which is Segoe UI Variable on Windows 11.
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

// ScenesAttribute marks a static method of App that returns the scenes of one group, alone in its file of Scenes/; File, the file it is declared in, orders the groups.
[AttributeUsage(AttributeTargets.Method)]
sealed class ScenesAttribute([System.Runtime.CompilerServices.CallerFilePath] string path = "") : System.Attribute
{
    public string File { get; } = path;
}

public sealed partial class App : Application, IXamlMetadataProvider
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

    // Scenes are the scenes of every group, the methods marked [Scenes] in the files of Scenes/, one in each, in the order of the files. It throws when there are none or two share a file.
    static IEnumerable<Scene> Scenes()
    {
        var groups = typeof(App).GetMethods(System.Reflection.BindingFlags.NonPublic | System.Reflection.BindingFlags.Static | System.Reflection.BindingFlags.DeclaredOnly)
            .Select(m => (Method: m, Mark: (ScenesAttribute?)System.Attribute.GetCustomAttribute(m, typeof(ScenesAttribute))))
            .Where(g => g.Mark != null)
            .OrderBy(g => g.Mark!.File, StringComparer.Ordinal)
            .ToList();
        if (groups.Count == 0)
        {
            throw new InvalidOperationException("no method is marked [Scenes]");
        }
        foreach (var shared in groups.GroupBy(g => g.Mark!.File).Where(f => f.Count() > 1))
        {
            throw new InvalidOperationException($"{shared.Key} has more than one method marked [Scenes]: {string.Join(", ", shared.Select(g => g.Method.Name))}");
        }
        return groups.SelectMany(g => (IEnumerable<Scene>)g.Method.Invoke(null, null)!);
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
