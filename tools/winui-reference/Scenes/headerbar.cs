namespace WinUIReference;

partial class App
{
    // HeaderBarScenes are the TitleBars of TestScreenshotsHeaderBar (ui/breeze/widgets_headerbar_test.go).
    [Scenes]
    static IEnumerable<Scene> HeaderBarScenes()
    {
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
    }
}
