This is a cross platfrom Go library to place an icon and menu in the notification area.
Supports Windows and Mac OSX, Linux coming soon.

## Install

```sh
go get github.com/getlantern/systray
```

```sh
sudo apt-get install libgtk-3-dev
```
if you don't have this package on Linux.


## Try

`cd` into `example` folder, and

```sh
go get
go build
./example # example.exe for Windows
```

Place your icon under `example` folder, and run `make_icon.bat` or `make_icon.sh`, whichever suit for you os.
Your icon should be .ico file under Windows, and .png for other platform.

## Credits

- https://github.com/xilp/systray
- https://github.com/cratonica/trayhost
