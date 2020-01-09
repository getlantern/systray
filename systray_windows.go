// +build windows

package systray

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"os"
	"sync/atomic"

	"github.com/lxn/walk"
)

var (
	tmpDir     string
	mainWindow *walk.MainWindow
	webView    *walk.WebView
	notifyIcon *walk.NotifyIcon

	actions      = make(map[int32]*walk.Action)
	nextActionId int32

	okayToClose int32
)

func nativeLoop(title string, width int, height int) {
	var err error
	mainWindow, err = walk.NewMainWindow()
	if err != nil {
		fail("Unable to create main window", err)
	}
	mainWindow.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		// don't close app unless we're actually finished
		actuallyClose := atomic.LoadInt32(&okayToClose) == 1
		*canceled = !actuallyClose
		if !actuallyClose {
			mainWindow.SetVisible(false)
		}
	})
	layout := walk.NewVBoxLayout()
	if err := mainWindow.SetLayout(layout); err != nil {
		fail("Unable to set main layout", err)
	}
	notifyIcon, err = walk.NewNotifyIcon(mainWindow)
	if err != nil {
		fail("Unable to create notify icon", err)
	}
	if title != "" {
		webView, err = walk.NewWebView(mainWindow)
		if err != nil {
			fail("Unable to create web view", err)
		}
		if err := mainWindow.SetTitle(title); err != nil {
			fail("Unable to set main title", err)
		}
		if err := mainWindow.SetWidth(width); err != nil {
			fail("Unable to set width", err)
		}
		if err := mainWindow.SetHeight(height); err != nil {
			fail("Unable to set height", err)
		}
	}
	systrayReady()
	mainWindow.Run()
}

func quit() {
	atomic.StoreInt32(&okayToClose, 1)
	mainWindow.Close()
	notifyIcon.Dispose()
	systrayExit()
}

// SetIcon sets the systray icon.
// iconBytes should be the content of .ico for windows and .ico/.jpg/.png
// for other platforms.
func SetIcon(iconBytes []byte) {
	md5 := md5.Sum(iconBytes)
	filename := fmt.Sprintf("systray.%x.ico", md5)
	iconpath := filepath.Join(walk.Resources.RootDirPath(), filename)
	// First, try to find a previously loaded icon in walk cache
	icon, err := walk.Resources.Icon(filename)
	if err != nil {
		// Cache miss, load the icon
		err := ioutil.WriteFile(iconpath, iconBytes, 0644)
		if err != nil {
			fail("Unable to save icon to disk", err)
		}
		defer os.Remove(iconpath)
		icon, err = walk.Resources.Icon(filename)
		if err != nil {
			fail("Unable to load icon", err)
		}
	}
	err = notifyIcon.SetIcon(icon)
	if err != nil {
		fail("Unable to set systray icon", err)
	}
	err = notifyIcon.SetVisible(true)
	if err != nil {
		fail("Unable to make systray icon visible", err)
	}
}

// SetTemplateIcon sets the systray icon as a template icon (on macOS), falling back
// to a regular icon on other platforms.
// templateIconBytes and iconBytes should be the content of .ico for windows and
// .ico/.jpg/.png for other platforms.
func SetTemplateIcon(templateIconBytes []byte, regularIconBytes []byte) {
	SetIcon(regularIconBytes)
}

// SetTitle sets the systray title, only available on Mac.
func SetTitle(title string) {
	// not supported on Windows
}

// SetTooltip sets the systray tooltip to display on mouse hover of the tray icon,
// only available on Mac and Windows.
func SetTooltip(tooltip string) {
	if err := notifyIcon.SetToolTip(tooltip); err != nil {
		fail("Unable to set tooltip", err)
	}
}

// ShowAppWindow shows the given URL in the application window. Only works if
// configureAppWindow has been called first.
func ShowAppWindow(url string) {
	if webView == nil {
		return
	}
	webView.SetURL(url)
	mainWindow.SetVisible(true)
}

func addOrUpdateMenuItem(item *MenuItem) {
	action := actions[item.id]
	if action == nil {
		item.id = nextActionId
		action = walk.NewAction()
		action.Triggered().Attach(func() {
			select {
			case item.ClickedCh <- struct{}{}:
				// okay
			default:
				// no listener, ignore
			}
		})
		if err := notifyIcon.ContextMenu().Actions().Add(action); err != nil {
			fail("Unable to add menu item to systray", err)
		}
		actions[item.id] = action
		atomic.AddInt32(&nextActionId, 1)
	}
	err := action.SetText(item.title)
	if err != nil {
		fail("Unable to set menu item text", err)
	}
	err = action.SetChecked(item.checked)
	if err != nil {
		fail("Unable to set menu item checked", err)
	}
	err = action.SetEnabled(!item.Disabled())
	if err != nil {
		fail("Unable to set menu item enabled", err)
	}
}

// SetIcon sets the icon of a menu item. Only works on macOS and Windows.
// iconBytes should be the content of .ico/.jpg/.png
func (item *MenuItem) SetIcon(iconBytes []byte) {
	md5 := md5.Sum(iconBytes)
	filename := fmt.Sprintf("systray.%x.ico", md5)
	iconpath := filepath.Join(walk.Resources.RootDirPath(), filename)
	// First, try to find a previously loaded icon in walk cache
	icon, err := walk.Resources.Image(filename)
	if err != nil {
		// Cache miss, load the icon
		err := ioutil.WriteFile(iconpath, iconBytes, 0644)
		if err != nil {
			fail("Unable to save icon to disk", err)
		}
		defer os.Remove(iconpath)
		icon, err = walk.Resources.Image(filename)
		if err != nil {
			fail("Unable to load icon", err)
		}
	}
	actions[item.id].SetImage(icon)
}

// SetTemplateIcon sets the icon of a menu item as a template icon (on macOS). On Windows, it
// falls back to the regular icon bytes and on Linux it does nothing.
// templateIconBytes and regularIconBytes should be the content of .ico for windows and
// .ico/.jpg/.png for other platforms.
func (item *MenuItem) SetTemplateIcon(templateIconBytes []byte, regularIconBytes []byte) {
	item.SetIcon(regularIconBytes)
}

func addSeparator(id int32) {
	action := walk.NewSeparatorAction()
	if err := notifyIcon.ContextMenu().Actions().Add(action); err != nil {
		fail("Unable to add separator", err)
	}
}

func hideMenuItem(item *MenuItem) {
	actions[item.id].SetVisible(false)
}

func showMenuItem(item *MenuItem) {
	actions[item.id].SetVisible(true)
}

func fail(msg string, err error) {
	panic(fmt.Errorf("%v: %v", msg, err))
}
