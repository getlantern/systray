/*
Package systray is a cross platfrom Go library to place an icon and menu in the notification area.
Supports Windows and Mac OSX currently, Linux coming soon.
Methods can be called from any goroutine except Run(), which should be called at the very beginning of main() to lock at main thread.
*/
package systray

/*
#cgo linux pkg-config: gtk+-3.0 appindicator3-0.1
#cgo linux CFLAGS: -DLINUX
#cgo windows CFLAGS: -DWIN32 -DUNICODE -D_UNICODE
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

var (
	readyCh       = make(chan interface{})
	clickedCh     = make(chan interface{})
	menuItems     = make(map[string]chan interface{})
	menuItemsLock sync.RWMutex
)

// Run initializes GUI and starts event loop, then invoke the onReady callback.
// It blocks until systray.Quit() is called.
// Should be called at the very beginning of main() to lock at main thread.
func Run(onReady func()) {
	runtime.LockOSThread()
	go func() {
		<-readyCh
		onReady()
	}()

	C.nativeLoop()
}

// Quit the systray and whole app
func Quit() {
	C.quit()
}

// SetIcon sets the systray icon.
// iconBytes should be the content of .ico for windows and .png for other platforms.
func SetIcon(iconBytes []byte) {
	cstr := (*C.char)(unsafe.Pointer(&iconBytes[0]))
	C.setIcon(cstr, (C.int)(len(iconBytes)))
}

// SetTitle sets the systray title, only available on Mac.
func SetTitle(title string) {
	C.setTitle(C.CString(title))
}

// SetTitle sets the systray tooltip after the mouse stayed on tray icon for a while, only available on Mac.
func SetTooltip(tooltip string) {
	C.setTooltip(C.CString(tooltip))
}

// Add menu item with designated title and tooltip, returning a channel that
// notifies whenever that menu item has been clicked.
//
// Menu items are keyed to an id. If the same id is added twice, the 2nd one
// overwrites the first.
//
// AddMenuItem can be safely invoked from different goroutines.
func AddMenuItem(id string, title string, tooltip string) <-chan interface{} {
	retChan := make(chan interface{})
	menuItemsLock.Lock()
	menuItems[id] = retChan
	C.addMenuItem(
		C.CString(id),
		C.CString(title),
		C.CString(tooltip),
	)
	menuItemsLock.Unlock()
	return retChan
}

//export systray_ready
func systray_ready() {
	readyCh <- nil
}

//export systray_menu_item_selected
func systray_menu_item_selected(cId *C.char) {
	id := C.GoString(cId)
	menuItemsLock.RLock()
	ch := menuItems[id]
	menuItemsLock.RUnlock()
	ch <- nil
}
