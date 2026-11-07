namespace WinUIReference;

partial class App
{
    // SwitcherScenes are the SelectorBar of TestScreenshotsViewSwitcher (ui/breeze/widgets_switcher_test.go).
    [Scenes]
    static IEnumerable<Scene> SwitcherScenes()
    {
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
    }
}
