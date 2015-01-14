package main

import "C"

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/getlantern/systray"
)

func main() {
	// Start a goroutine for doing stuff that our Go application will do
	go func() {
		ch1 := systray.AddMenu("change", "Change Me", "Change Me")
		ch2 := systray.AddMenu("quit", "Quit", "Quit the whole app")
		// This is just an example of some processing that happens outside of
		// the Cocoa app.
		for {
			log.Print("Waiting")
			time.Sleep(1 * time.Second)
			systray.UpdateTitle("New Title")
			clicked := systray.AddMenu("asdf", "New Title", "sadfds")
			if ret := <-clicked; ret == true {
				break
			}
		}
		ret := <-ch1
		ret = <-ch2
		fmt.Println(ret)
	}()
	// Start the Cocoa app (this blocks)
	systray.EnterLoop()
}

func OnAction(action string) {
	log.Printf("Got action: %s", action)
	switch action {
	case "dostuff":
		systray.UpdateTitle("New Title")
	case "quit":
		log.Printf("Quitting")
		os.Exit(0)
	default:
		log.Printf("Got unexpected action: %s", action)
	}
}
