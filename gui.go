package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	// "fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func makeHeader(logo *fyne.StaticResource) fyne.CanvasObject {
	createDbtBtn := widget.NewButtonWithIcon("Create DBT", theme.ContentAddIcon(), func() {
		execFolderScript()
	})
	createDbtBtn.IconPlacement = widget.ButtonIconTrailingText

	icon := canvas.NewImageFromResource(logo)
	icon.SetMinSize(fyne.NewSize(100, 100))
	icon.FillMode = canvas.ImageFillContain

	addBtn := widget.NewButton("Add", func() {

	})

	return container.NewBorder(
		icon,
		nil,
		createDbtBtn,
		addBtn,
	)
}

func MakeGui(logo *fyne.StaticResource) fyne.CanvasObject {
	header := makeHeader(logo)

	return container.NewVBox(header)
}
