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
)

var menuItems map[string]chan bool = make(map[string]chan bool)
var lock sync.RWMutex

//export callMe
func callMe(cname *C.char) {
	name := C.GoString(cname)
	lock.RLock()
	ch := menuItems[name]
	lock.RUnlock()
	ch <- true
}

func AddMenu(name string, title string, tooltip string) chan bool {
	retChan := make(chan bool)
	lock.Lock()
	menuItems[name] = retChan
	lock.Unlock()
	C.addMenu(
		C.CString(name),
		C.CString(title),
		C.CString(tooltip),
	)
	return retChan
}

// Start the Cocoa app (this blocks)
func EnterLoop() {
	runtime.LockOSThread()
	C.nativeLoop()
}

func Quit() {
	C.quit()
}
