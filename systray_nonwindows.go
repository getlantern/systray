// +build !windows

package systray

/*
#cgo linux CFLAGS: -DWEBVIEW_GTK=1
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0 appindicator3-0.1
#cgo darwin CFLAGS: -DDARWIN -x objective-c -fobjc-arc
#cgo darwin LDFLAGS: -framework Cocoa -framework WebKit

#include "systray.h"
*/
import "C"

import (
	"unsafe"
)

func nativeLoop(title string, width int, height int) {
	C.nativeLoop(C.CString(title), C.int(width), C.int(height))
}

func quit() {
	C.quit()
}

// ShowAppWindow shows the given URL in the application window. Only works if
// configureAppWindow has been called first.
func ShowAppWindow(url string) {
	C.showAppWindow(C.CString(url))
}

// SetIcon sets the systray icon.
// iconBytes should be the content of .ico for windows and .ico/.jpg/.png
// for other platforms.
func SetIcon(iconBytes []byte) {
	cstr := (*C.char)(unsafe.Pointer(&iconBytes[0]))
	C.setIcon(cstr, (C.int)(len(iconBytes)), false)
}

// SetTitle sets the systray title, only available on Mac.
func SetTitle(title string) {
	C.setTitle(C.CString(title))
}

// SetTooltip sets the systray tooltip to display on mouse hover of the tray icon,
// only available on Mac and Windows.
func SetTooltip(tooltip string) {
	C.setTooltip(C.CString(tooltip))
}

func addOrUpdateMenuItem(item *MenuItem) {
	var disabled C.short
	if item.disabled {
		disabled = 1
	}
	var checked C.short
	if item.checked {
		checked = 1
	}
	if item.parent == nil {
		C.add_or_update_menu_item(
			C.int(item.id),
			C.CString(item.title),
			C.CString(item.tooltip),
			disabled,
			checked,
		)
	} else {
		C.add_or_update_submenu_item(
			C.int(item.parent.id),
			C.int(item.id),
			C.CString(item.title),
			C.CString(item.tooltip),
			disabled,
			checked,
		)
	}
}

func addSeparator(id int32) {
	C.add_separator(C.int(id))
}

func hideMenuItem(item *MenuItem) {
	C.hide_menu_item(
		C.int(item.id),
	)
}

func showMenuItem(item *MenuItem) {
	C.show_menu_item(
		C.int(item.id),
	)
}

//export systray_ready
func systray_ready() {
	systrayReady()
}

//export systray_on_exit
func systray_on_exit() {
	systrayExit()
}

//export systray_menu_item_selected
func systray_menu_item_selected(cID C.int) {
	systrayMenuItemSelected(int32(cID))
}
