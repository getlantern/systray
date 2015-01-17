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

var (
	readyCh       = make(chan interface{})
	clickedCh     = make(chan interface{})
	menuItems     = make(map[string]chan interface{})
	menuItemsLock sync.RWMutex
)

// Run the Cocoa app (this blocks)
func Run(onReady func()) {
	go func() {
		<-readyCh
		onReady()
	}()

	runtime.LockOSThread()
	C.nativeLoop()
}

// Quit the Cocoa app
func Quit() {
	C.quit()
}

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

// Add menu item with designated title and tooltip, waiting on returned chan to get notified when menu item clicked.
// Add again with same menuId will override previous one.
// Can be invoked from different goroutines.
func AddMenu(menuId string, title string, tooltip string) chan interface{} {
	retChan := make(chan interface{})
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

//export systray_ready
func systray_ready() {
	readyCh <- nil
}

//export systray_menu_item_selected
func systray_menu_item_selected(cmenuId *C.char) {
	menuId := C.GoString(cmenuId)
	menuItemsLock.RLock()
	ch := menuItems[menuId]
	menuItemsLock.RUnlock()
	ch <- nil
}
