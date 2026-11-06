namespace WinUIReference;

partial class App
{
    // StatusScenes are the status widgets of TestScreenshotsStatus (ui/breeze/widgets_status_test.go): statusView, statusTitlesView and statusActionsView.
    [Scenes]
    static IEnumerable<Scene> StatusScenes()
    {
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
    }
}
