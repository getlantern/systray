//+build !windows

package main

/*
#cgo linux pkg-config: gtk+-3.0 appindicator3-0.1
#cgo darwin CFLAGS: -DDARWIN -x objective-c -fobjc-arc
#cgo darwin LDFLAGS: -framework Cocoa -framework Webkit

#include "webview.h"
*/
import "C"

func configureWebview(title string, width, height int) {
	C.configureAppWindow(C.CString(title), C.int(width), C.int(height))
}

func showWebview(url string) {
	C.showAppWindow(C.CString(url))
}
