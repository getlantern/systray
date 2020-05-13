package main

import (
	"github.com/zserge/webview"
)

func main() {
	wv := webview.New(true)
	wv.SetTitle("Some Title")
	wv.SetSize(800, 600, webview.HintNone)
	wv.Navigate("https://www.getlantern.org")
	wv.Run()
}
