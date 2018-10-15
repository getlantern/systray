// +build !windows

package systray

/*
#cgo linux pkg-config: gtk+-3.0 appindicator3-0.1
#cgo darwin CFLAGS: -DDARWIN -x objective-c -fobjc-arc
#cgo darwin LDFLAGS: -framework Cocoa

#include "systray.h"
*/
import "C"

import (
	"unsafe"
)

func nativeLoop() {
	C.nativeLoop()
}

func quit() {
	C.quit()
}

// SetIcon sets the systray icon.
// iconBytes should be the content of .ico for windows and .ico/.jpg/.png
// for other platforms.
func SetIcon(iconBytes []byte) {
	cstr := (*C.char)(unsafe.Pointer(&iconBytes[0]))
	C.setIcon(cstr, (C.int)(len(iconBytes)))
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
	C.add_or_update_menu_item(
		C.int(item.id),
		C.CString(item.title),
		C.CString(item.tooltip),
		disabled,
		checked,
	)
}

// SetIcon sets the icon of a menu item. Only available on Mac.
// iconBytes should be the content of .ico/.jpg/.png
func (item *MenuItem) SetIcon(iconBytes []byte) {
	cstr := (*C.char)(unsafe.Pointer(&iconBytes[0]))
	C.setMenuItemIcon(cstr, (C.int)(len(iconBytes)), C.int(item.id))
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

// Text dialog things

var (
	eulaAcceptedCallback func()
	eulaDeclinedCallback func()
)

func _showTextDialog(rtfFile []byte, onAccepted, onDeclined func()) {
	eulaAcceptedCallback = onAccepted
	eulaDeclinedCallback = onDeclined

	cstr := (*C.char)(unsafe.Pointer(&rtfFile[0]))
	C.showTextDialog(cstr)
}

//export onTextDialogClosed
func onTextDialogClosed(accepted C.int) {
	eulaAccepted := int32(accepted) == 1

	if eulaAccepted && eulaAcceptedCallback != nil {
		eulaAcceptedCallback()
		eulaAcceptedCallback = nil
	}

	if !eulaAccepted && eulaDeclinedCallback != nil {
		eulaDeclinedCallback()
		eulaDeclinedCallback = nil
	}
}
