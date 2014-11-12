package main

import "C"

import (
	"log"
	"os"
)

//export OnAction
// Note - this has to be in a separate file from the one where we import
// systray.m, otherwise it doesn't compile.
func OnAction(actionChars *C.char) {
	action := C.GoString(actionChars)
	log.Printf("Got action: %s", action)
	switch action {
	case "quit":
		log.Printf("Quitting")
		os.Exit(0)
	default:
		log.Printf("Got unexpected action: %s", action)
	}
}
