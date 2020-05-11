package main

import (
	"fmt"
	"github.com/lxn/walk"
)

func showWebviewOnWindows() {
	mainWindow, err := walk.NewMainWindow()
	if err != nil {
		panic(fmt.Sprintf("Failed to create main window: %v\n", err))
	}
	mainWindow.SetTitle("Webview")
	mainWindow.SetWidth(800)
	mainWindow.SetHeight(600)
	layout := walk.NewVBoxLayout()
	if err := mainWindow.SetLayout(layout); err != nil {
		panic(fmt.Sprintf("Failed to set layout: %v\n", err))
	}
	webView, err := walk.NewWebView(mainWindow)
	if err != nil {
		panic(fmt.Sprintf("Failed to create webview window: %v\n", err))
	}
	webView.SetURL("https://www.getlantern.org")
	mainWindow.SetVisible(true)
	mainWindow.Run()
}
