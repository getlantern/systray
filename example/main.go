package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	systray.Run(onReady)
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Pretty awesome棒")
	chQuit := systray.AddMenuItem("quit", "Quit", "Quit the whole app")
	go func() {
		<-chQuit
		systray.Quit()
		fmt.Println("Quit now...")
	}()

	// We can manipulate the systray in other goroutines
	go func() {
		ch := systray.AddMenuItem("change", "Change Me", "Change Me")
		chUrl := systray.AddMenuItem("lantern", "Open Lantern.org", "my home")
		chQuit := systray.AddMenuItem("quit2", "退出", "Quit the whole app")
		// This is just an example of some processing that happens outside of
		// the Cocoa app.
		for {
			select {
			case <-ch:
				ch = systray.AddMenuItem("change", "I've Changed", "Catch Me")
			case <-chUrl:
				open.Run("https://www.getlantern.org")
			case <-chQuit:
				systray.Quit()
				fmt.Println("Quit2 now...")
				return
			}
		}
	}()
}
