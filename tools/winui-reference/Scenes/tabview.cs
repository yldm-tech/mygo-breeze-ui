namespace WinUIReference;

partial class App
{
    // TabViewScenes are the TabViews of TestScreenshotsTabView (ui/breeze/widgets_tabview_test.go) and of TestFluentTabViewAsMeasured (ui/breeze/fluent_native_test.go).
    [Scenes]
    static IEnumerable<Scene> TabViewScenes()
    {
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
    }
}
