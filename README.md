Package systray is a cross platfrom Go library to place an icon and menu in the notification area.
Tested on Windows 8, Mac OSX, Ubuntu 14.10 and Debian 7.6.

## Usage
```go
func main() {
	// Should be called at the very beginning of main().
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Pretty awesome超级棒")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	// Sets the icon of a menu item. Only available on Mac.
	mQuit.SetIcon(icon.Data)
}

func onExit() {
	// clean up here
}
```
Menu item can be checked and / or disabled. Methods except `Run()` can be invoked from any goroutine. See demo code under `example` folder.

## Platform specific concerns

### Linux

```sh
sudo apt-get install libgtk-3-dev libappindicator3-dev
```
Checked menu item not implemented on Linux yet.

## Try

Under `example` folder.
Place tray icon under `icon`, and use `make_icon.bat` or `make_icon.sh`, whichever suit for your os, to convert the icon to byte array.
Your icon should be .ico file under Windows, whereas .ico, .jpg and .png is supported on other platform.

```sh
go get
go run main.go
```

## Building and the Console Window

By default, the binary created by `go build` will cause a console window to be opened on both Windows and macOS when run.

### Windows

To prevent launching a console window when running on Windows, add these command-line build flags:

```sh
go build -ldflags -H=windowsgui
```

### macOS

On macOS, you will need to create an application bundle to wrap the binary; simply folders with the following minimal structure and assets:

```
SystrayApp.app/
  Contents/
    Info.plist
  MacOS/
    go-executable
  Resources/
    SystrayApp.icns
```

Consult the [Official Apple Documentation here](https://developer.apple.com/library/archive/documentation/CoreFoundation/Conceptual/CFBundles/BundleTypes/BundleTypes.html#//apple_ref/doc/uid/10000123i-CH101-SW1).

## Credits

- https://github.com/xilp/systray
- https://github.com/cratonica/trayhost
