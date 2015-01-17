package main

import (
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

var (
	clickedOnce = int32(0)
)

func main() {
	systray.Run(onReady)
}

func onReady() {
	systray.SetIcon(iconArray)
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Pretty awesome")
	chQuit := systray.AddMenu("quit", "Quit", "Quit the whole app")
	go func() {
		_ = <-chQuit
		systray.Quit()
	}()

	// We can manipulate the systray in other goroutines
	go func() {
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
}
