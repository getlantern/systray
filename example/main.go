package main

import (
	"fmt"
	"time"

	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	// Start a goroutine for doing stuff that our Go application will do
	go func() {
		time.Sleep(1 * time.Second)
		systray.SetIcon(iconArray)
		systray.SetTitle("Click Me!")
		ch := systray.WaitForSystrayClicked()
		_ = <-ch
		fmt.Println("systray clicked, set off now")

		go func() {
			systray.SetTitle("Awesome App")
			systray.SetTooltip("Pretty awesome")
			chQuit := systray.AddMenu("quit", "Quit", "Quit the whole app")
			_ = <-chQuit
			systray.Quit()
		}()

		// We can also manipulate systray in other goroutine
		go func() {
			time.Sleep(1 * time.Second)
			ch := systray.AddMenu("change", "Change Me", "Change Me")
			chUrl := systray.AddMenu("lantern", "Open Lantern.org", "my home")
			chQuit := systray.AddMenu("quit2", "Another Quit", "Quit the whole app")
			// This is just an example of some processing that happens outside of
			// the Cocoa app.
			for {
				select {
				case _ = <-ch:
					ch = systray.AddMenu("change", "I've Changed", "Catch Me")
				case _ = <-chUrl:
					open.Run("https://www.getlantern.org")
				case _ = <-chQuit:
					systray.Quit()
					return
				}
			}
		}()
	}()
	// Start the Cocoa app (this blocks)
	systray.EnterLoop()
}
