namespace WinUIReference;

partial class App
{
    // RowsScenes are CommunityToolkit's SettingsCards and SettingsExpander where TestScreenshotsRows (ui/breeze/widgets_rows_test.go) lays out Breeze UI's cards and rows.
    [Scenes]
    static IEnumerable<Scene> RowsScenes()
    {
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
    }
}
