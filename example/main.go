package main

import (
	"fmt"
	"os"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	file, err := os.Create("out.txt")
	if err != nil {
		fmt.Println("Couldn't create out.txt")
		return
	}
	defer file.Close()
	onExit := func() {
		fmt.Println("Starting onExit")
		file.Write([]byte("On Exit...\n"))
		file.Sync()
		fmt.Println("Finished onExit")
	}
	// Should be called at the very beginning of main().
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Lantern")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuit.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	// We can manipulate the systray in other goroutines
	go func() {
		systray.SetIcon(icon.Data)
		systray.SetTitle("Awesome App")
		systray.SetTooltip("Pretty awesome棒棒嗒")
		mChange := systray.AddMenuItem("Change Me", "Change Me")
		mChecked := systray.AddMenuItem("Unchecked", "Check Me")
		mEnabled := systray.AddMenuItem("Enabled", "Enabled")
		systray.AddMenuItem("Ignored", "Ignored")
		mUrl := systray.AddMenuItem("Open Lantern.org", "my home")
		mQuit := systray.AddMenuItem("退出", "Quit the whole app")
		for {
			select {
			case <-mChange.ClickedCh:
				mChange.SetTitle("I've Changed")
			case <-mChecked.ClickedCh:
				if mChecked.Checked() {
					mChecked.Uncheck()
					mChecked.SetTitle("Unchecked")
				} else {
					mChecked.Check()
					mChecked.SetTitle("Checked")
				}
			case <-mEnabled.ClickedCh:
				mEnabled.SetTitle("Disabled")
				mEnabled.Disable()
			case <-mUrl.ClickedCh:
				open.Run("https://www.getlantern.org")
			case <-mQuit.ClickedCh:
				systray.Quit()
				fmt.Println("Quit2 now...")
				return
			}
		}
	}()
}
