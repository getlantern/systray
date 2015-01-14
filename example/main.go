package main

import (
	"fmt"
	"log"
	"time"

	"github.com/getlantern/systray"
)

func main() {
	// Start a goroutine for doing stuff that our Go application will do
	go func() {
		// need some time to let app starts up
		time.Sleep(1 * time.Second)
		ch1 := systray.AddMenu("change", "Change Me", "Change Me")
		ch2 := systray.AddMenu("quit", "Quit", "Quit the whole app")
		// This is just an example of some processing that happens outside of
		// the Cocoa app.
		for {
			log.Print("Waiting")
			select {
			case _ = <-ch1:
				fmt.Println("clicked!")
			case _ = <-ch2:
				fmt.Println("quit!")
				systray.Quit()
				return
			}
		}
	}()
	// Start the Cocoa app (this blocks)
	systray.EnterLoop()
}
