package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "systray.h"
*/
import "C"

import (
	"log"
	"os"
	"runtime"
	"time"
)

// Arrange that main.main runs on main thread so that our calls into the Cocoa
// app all happen from the main thread.
func init() {
	runtime.LockOSThread()
}

func main() {
	// Start a goroutine for doing stuff that our Go application will do
	go func() {
		// This is just an example of some processing that happens outside of
		// the Cocoa app.
		for {
			log.Print("Waiting")
			time.Sleep(1 * time.Second)
		}
	}()
	// Start the Cocoa app (this blocks)
	C.StartApp()
}

func updateTitle(newTitle string) {
	C.updateTitle(C.CString(newTitle))
}

//export OnAction
func OnAction(actionChars *C.char) {
	action := C.GoString(actionChars)
	log.Printf("Got action: %s", action)
	switch action {
	case "dostuff":
		updateTitle("New Title")
	case "quit":
		log.Printf("Quitting")
		os.Exit(0)
	default:
		log.Printf("Got unexpected action: %s", action)
	}
}
