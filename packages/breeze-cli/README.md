# breeze-ui-cli

The CLI for Breeze UI, which draws MyGo's native UI widgets as macOS, Windows 11 and GNOME draw their controls. It lists the widgets and the platforms they follow, and writes snippets that use them.

```sh
npx --yes --package breeze-ui-cli breeze-ui init
npx --yes --package breeze-ui-cli breeze-ui list -kind component
npx --yes --package breeze-ui-cli breeze-ui add Switch -o ui/switch_snippet.go
```

For a released binary, install the Go command directly:

```sh
go install github.com/yldm-tech/mygo-breeze-ui/cmd/breeze-ui@latest
```
