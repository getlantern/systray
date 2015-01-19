This is a cross platfrom Go library to place an icon and menu in the notification area.
Supports Windows, Mac OSX and Linux systems with gtk+.

## Install

```sh
go get github.com/getlantern/systray
```

## Try

`cd` into `example` folder, and

```sh
go build
./example # example.exe for Windows
```

Place your icon under `example` folder, and run `make_icon.bat` or `make_icon.sh`, whichever suit for you os.
Your icon should be .ico file under Windows, and .png for other platform.

## Credits

- https://github.com/xilp/systray
- https://github.com/cratonica/trayhost
