Package systray is a cross platfrom Go library to place an icon and menu in the notification area.
Supports Windows and Mac OSX currently, Linux coming soon.

## Usage
```go
func main() {
	// Should be called at the very beginning of main().
	systray.Run(onReady)
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Pretty awesome超级棒")
	chQuit := systray.AddMenuItem("quit", "Quit", "Quit the whole app")
}
```


## Platform specific concerns

### Linux

```sh
sudo apt-get install libgtk-3-dev libappindicator3-dev
```

### Windows

Install [MinGW-W64](http://sourceforge.net/projects/mingw-w64) as it has up to date SDK headers we require.


## Try

Under `example` folder, place your tray icon there, and run `make_icon.bat` or `make_icon.sh`, whichever suit for you os.
Your icon should be .ico file under Windows, and .png for other platform.

```sh
go get
go build
./example # example.exe for Windows
```

## Credits

- https://github.com/xilp/systray
- https://github.com/cratonica/trayhost
