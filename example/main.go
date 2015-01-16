package main

import (
	"fmt"
	"log"
	"time"

	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	// Start a goroutine for doing stuff that our Go application will do
	go func() {
		// need some time to let app starts up
		time.Sleep(1 * time.Second)
		ch := systray.AddMenu("change", "Change Me", "Change Me")
		chUrl := systray.AddMenu("lantern", "Open Lantern.org", "my home")
		chQuit := systray.AddMenu("quit", "Quit", "Quit the whole app")
		// This is just an example of some processing that happens outside of
		// the Cocoa app.
		for {
			log.Print("Waiting")
			select {
			case _ = <-ch:
				fmt.Println("clicked!")
			case _ = <-chUrl:
				open.Run("https://www.getlantern.org")
			case _ = <-chQuit:
				fmt.Println("quit!")
				systray.Quit()
				return
			}
		}
	}()
	// Start the Cocoa app (this blocks)
	systray.EnterLoop()
}
