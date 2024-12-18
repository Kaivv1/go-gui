package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"

	// "fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func makeHeader(logo *fyne.StaticResource, w fyne.Window, storage *Storage, refresh func()) fyne.CanvasObject {
	createDbtBtn := widget.NewButtonWithIcon("Генерирай ДБТ Папка", theme.ContentAddIcon(), func() {
		data, err := execFolderScript(w)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if err = storage.WriteToStorage(&StorageStructure{
			DBTs: data,
		}); err != nil {
			dialog.ShowError(err, w)
			return
		}
		refresh()
	})

	createDbtBtn.IconPlacement = widget.ButtonIconTrailingText
	createDbtBtn.Resize(fyne.NewSize(150, 50))
	icon := canvas.NewImageFromResource(logo)
	icon.SetMinSize(fyne.NewSize(70, 70))
	icon.FillMode = canvas.ImageFillContain

	return container.NewHBox(
		icon,
		layout.NewSpacer(),
		container.NewVBox(
			layout.NewSpacer(),
			createDbtBtn,
			layout.NewSpacer(),
		),
	)
}

func MakeGui(logo *fyne.StaticResource, w fyne.Window, storage *Storage) fyne.CanvasObject {
	addFiles, refreshSelector := makeAddFilesSpace(storage, w)
	header := makeHeader(logo, w, storage, refreshSelector)

	return container.NewBorder(header, nil, nil, nil, addFiles)
}
