package systray

/*
#cgo linux pkg-config: gtk+-3.0
#cgo linux CFLAGS: -DLINUX
#cgo windows CFLAGS: -DWIN32
#cgo darwin CFLAGS: -DDARWIN -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa

#include "systray.h"
*/
import "C"
import (
	"runtime"
	"sync"
	"unsafe"
)

func SetIcon(iconBytes []byte) {
	cstr := (*C.char)(unsafe.Pointer(&iconBytes[0]))
	C.setIcon(cstr, (C.int)(len(iconBytes)))
}

func SetTitle(title string) {
	C.setTitle(C.CString(title))
}

func SetTooltip(tooltip string) {
	C.setTooltip(C.CString(tooltip))
}

// Waiting on returned chan to get notified when systray clicked.
// Only valid if no menu item added.
func WaitForSystrayClicked() <-chan bool {
	return systrayClickedChan
}

// Add menu item with designated title and tooltip, waiting on returned chan to get notified when menu item clicked.
// Add again with same menuId will override previous one.
// Can be invoked from different goroutines.
func AddMenu(menuId string, title string, tooltip string) chan bool {
	retChan := make(chan bool)
	menuItemsLock.Lock()
	menuItems[menuId] = retChan
	C.addMenu(
		C.CString(menuId),
		C.CString(title),
		C.CString(tooltip),
	)
	menuItemsLock.Unlock()
	return retChan
}

// Start the Cocoa app (this blocks)
func EnterLoop() {
	runtime.LockOSThread()
	C.nativeLoop()
}

// Quit the Cocoa app
func Quit() {
	C.quit()
}

var systrayClickedChan chan bool = make(chan bool)
var menuItems map[string]chan bool = make(map[string]chan bool)
var menuItemsLock sync.RWMutex

//export systray_clicked_call_back
func systray_clicked_call_back() {
	systrayClickedChan <- true
}

//export systray_menu_item_call_back
func systray_menu_item_call_back(cmenuId *C.char) {
	menuId := C.GoString(cmenuId)
	menuItemsLock.RLock()
	ch := menuItems[menuId]
	menuItemsLock.RUnlock()
	ch <- true
}
