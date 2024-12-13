package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	// "fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()

	w := a.NewWindow("Window")

	w.Resize(fyne.NewSize(800, 800))
	w.SetContent(makeHeader())

	w.ShowAndRun()
}
