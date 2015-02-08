package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	// Should be called at the very beginning of main().
	systray.Run(onReady)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Pretty awesome超级棒")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuit.Ch
		systray.Quit()
		fmt.Println("Quit now...")
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
			case <-mChange.Ch:
				mChange.Title = "I've Changed"
				mChange.Update()
			case <-mChecked.Ch:
				mChecked.Checked = !mChecked.Checked
				if mChecked.Checked {
					mChecked.Title = "Checked"
				} else {
					mChecked.Title = "Unchecked"
				}
				mChecked.Update()
			case <-mEnabled.Ch:
				mEnabled.Disabled = !mEnabled.Disabled
				mEnabled.Title = "Disabled"
				mEnabled.Update()
			case <-mUrl.Ch:
				open.Run("https://www.getlantern.org")
			case <-mQuit.Ch:
				systray.Quit()
				fmt.Println("Quit2 now...")
				return
			}
		}
	}()
}
