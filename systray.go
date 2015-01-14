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
	"fmt"
	"runtime"
	"unsafe"
)

/* type menuItem struct {
	title   string
	tooltip string
	channel chan bool
}

var menuItems map[string]menuItem = make(map[string]menuItem)

func callback(cname *C.char) {
	name := C.GoString(cname)
	menuItems[name].channel <- true
}

var theCallback = callback*/

func AddMenu(name string, title string, tooltip string) chan bool {
	retChan := make(chan bool)
	/*menuItems[name] = menuItem{
		title,
		tooltip,
		retChan,
	}*/
	callback := func(c chan bool) func() {
		return func() {
			fmt.Println("adsfasdf")
			c <- true
		}
	}(retChan)
	C.addMenu(
		C.CString(name),
		C.CString(title),
		C.CString(tooltip),
		unsafe.Pointer(&callback),
	)
	return retChan
}

// Start the Cocoa app (this blocks)
func EnterLoop() {
	C.nativeLoop()
}

func UpdateTitle(newTitle string) {
	C.updateTitle(C.CString(newTitle))
}

// Arrange that main.main runs on main thread so that our calls into the Cocoa
// app all happen from the main thread.
func init() {
	runtime.LockOSThread()
}

// export onAction
func onAction(name *C.char) {
	C.GoString(name)
}
