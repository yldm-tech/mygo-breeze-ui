namespace WinUIReference;

partial class App
{
    // TextScenes are the TextBlock styles of TestScreenshotsText (ui/breeze/widgets_text_test.go).
    [Scenes]
    static IEnumerable<Scene> TextScenes()
    {
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
}
