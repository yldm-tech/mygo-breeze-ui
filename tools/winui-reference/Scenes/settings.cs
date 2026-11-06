namespace WinUIReference;

partial class App
{
    // SettingsScenes are the settings rows of TestScreenshotsSettings (ui/breeze/widgets_settings_test.go).
    [Scenes]
    static IEnumerable<Scene> SettingsScenes()
    {
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
    }
}
