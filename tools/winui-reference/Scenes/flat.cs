namespace WinUIReference;

partial class App
{
    // FlatScenes are the flat buttons of TestScreenshotsFlat (ui/breeze/widgets_flat_test.go).
    [Scenes]
    static IEnumerable<Scene> FlatScenes()
    {
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
    }
}
