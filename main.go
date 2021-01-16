package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
)

func main() {
	mainApp := app.New()
	window := mainApp.NewWindow("youtube-dl")
	window.Resize(fyne.Size{
		Width:  450,
		Height: 250,
	})
	window.CenterOnScreen()
	window.SetIcon(fyne.NewStaticResource("icon.png", icon))

	window.SetContent(createLayout(window))
	window.ShowAndRun()
}
