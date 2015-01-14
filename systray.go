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

//export callMe
func callMe(cname *C.char) {
	name := C.GoString(cname)
	fmt.Println(name)
	fmt.Println("callMe!!!")
}

var theCallMe = callMe

func AddMenu(name string, title string, tooltip string) chan bool {
	retChan := make(chan bool)
	/*menuItems[name] = menuItem{
		title,
		tooltip,
		retChan,
	}*/
	C.addMenu(
		C.CString(name),
		C.CString(title),
		C.CString(tooltip),
		unsafe.Pointer(&theCallMe),
	)
	return retChan
}

// Start the Cocoa app (this blocks)
func EnterLoop() {
	runtime.LockOSThread()
	C.nativeLoop()
}

func UpdateTitle(newTitle string) {
	C.updateTitle(C.CString(newTitle))
}
