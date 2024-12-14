//go:generate fyne bundle -o bundled.go assets/go-bear-2.png

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	// "fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	logo := resourceGoBear2Png
	a.SetIcon(logo)
	a.Settings().SetTheme(newCustomTheme())
	w := a.NewWindow("go-gui")
	w.Resize(fyne.NewSize(800, 800))
	w.SetIcon(logo)

	w.SetContent(MakeGui(logo, a))
	w.ShowAndRun()
}
