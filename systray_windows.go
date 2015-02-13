// +build windows

package systray

import (
	"io/ioutil"
	"os"
	"syscall"
	"unsafe"
)

var (
	mod                      = syscall.NewLazyDLL("systray.dll")
	_nativeLoop              = mod.NewProc("nativeLoop")
	_quit                    = mod.NewProc("quit")
	_setIcon                 = mod.NewProc("setIcon")
	_setTitle                = mod.NewProc("setTitle")
	_setTooltip              = mod.NewProc("setTooltip")
	_add_or_update_menu_item = mod.NewProc("add_or_update_menu_item")

	iconFiles = make([]*os.File, 0)
)

func nativeLoop() {
	_nativeLoop.Call(
		syscall.NewCallback(systray_ready),
		syscall.NewCallback(systray_menu_item_selected))
}

func quit() {
	_quit.Call()
	for _, f := range iconFiles {
		err := os.Remove(f.Name())
		if err != nil {
			log.Debugf("Unable to delete temporary icon file %v: %v", f.Name(), err)
		}
	}
}

// SetIcon sets the systray icon.
// iconBytes should be the content of .ico for windows and .ico/.jpg/.png
// for other platforms.
func SetIcon(iconBytes []byte) {
	f, err := ioutil.TempFile("", "systray_temp_icon")
	if err != nil {
		log.Errorf("Unable to create temp icon: %v", err)
		return
	}
	defer f.Close()
	_, err = f.Write(iconBytes)
	if err != nil {
		log.Errorf("Unable to write icon to temp file %v: %v", f.Name(), f)
		return
	}
	f.Close()
	name, err := strPtr(f.Name())
	if err != nil {
		log.Errorf("Unable to convert name to string pointer: %v", err)
		return
	}
	_setIcon.Call(name)
}

// SetTitle sets the systray title, only available on Mac.
func SetTitle(title string) {
	t, err := strPtr(title)
	if err != nil {
		log.Errorf("Unable to convert title to string pointer: %v", err)
		return
	}
	_setTitle.Call(t)
}

// SetTitle sets the systray tooltip to display on mouse hover of the tray icon,
// only available on Mac.
func SetTooltip(tooltip string) {
	t, err := strPtr(tooltip)
	if err != nil {
		log.Errorf("Unable to convert tooltip to string pointer: %v", err)
		return
	}
	_setTooltip.Call(t)
}

func addOrUpdateMenuItem(item *MenuItem) {
	var disabled = 0
	if item.disabled {
		disabled = 1
	}
	var checked = 0
	if item.checked {
		checked = 1
	}
	title, err := strPtr(item.title)
	if err != nil {
		log.Errorf("Unable to convert title to string pointer: %v", err)
		return
	}
	tooltip, err := strPtr(item.tooltip)
	if err != nil {
		log.Errorf("Unable to convert tooltip to string pointer: %v", err)
		return
	}
	_add_or_update_menu_item.Call(
		uintptr(item.id),
		title,
		tooltip,
		uintptr(disabled),
		uintptr(checked),
	)
}

func strPtr(s string) (uintptr, error) {
	bp, err := syscall.BytePtrFromString(s)
	if err != nil {
		return 0, err
	}
	return uintptr(unsafe.Pointer(bp)), nil
}

func systray_ready() uintptr {
	systrayReady()
	return 0
}

func systray_menu_item_selected(id uintptr) uintptr {
	systrayMenuItemSelected(int32(id))
	return 0
}
