//go:generate fyne bundle -o bundled.go assets/go-bear-2.png

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	// "fyne.io/fyne/v2/widget"
)

const (
	TITLE_TEXT = 24
	BASE_TEXT  = 16
)

var PLATFORMS = [4]string{"Имейл", "Архимед", "Хермес", "Regix"}

var REQUEST_TYPES = [3]string{"Предоставяне", "Променяне", "Прекратяване"}

func main() {
	a := app.NewWithID("com.example.myapp")

	logo := resourceGoBear2Png
	a.SetIcon(logo)
	a.Settings().SetTheme(newCustomTheme())
	w := a.NewWindow("go-gui")
	w.Resize(fyne.NewSize(1000, 900))
	w.SetIcon(logo)
	storage, err := NewStorage("storage.json")
	if err != nil {
		dialog.ShowError(err, w)
	}

	w.SetContent(MakeGui(logo, w, storage))
	w.ShowAndRun()
}
