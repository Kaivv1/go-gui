package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func makeHeader() fyne.CanvasObject {
	return container.NewBorder(
		container.NewHBox(
			widget.NewButton("Създай ДБТ", func() {
				execFolderScript()
			}),
		),
		container.NewHBox(
			widget.NewLabel("Icon"),
		),
		container.NewHBox(
			widget.NewButton("Добави", func() {}),
		),
		nil,
	)
}

// func МakeGui() fyne.CanvasObject {
// 	return container.NewVBox(makeHeader())
// }
