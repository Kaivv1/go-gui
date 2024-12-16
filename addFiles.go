package main

import (
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func makeAddFilesSpace(storage *Storage, w fyne.Window) (fyne.CanvasObject, func()) {
	title := canvas.NewText("Добави заявки", color.White)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = TITLE_TEXT

	data, err := storage.GetFromStorage()
	if err != nil {
		data = &StorageStructure{}
	}
	DBTs := []string{}
	for _, DBT := range data.DBTs {
		fullStr := fmt.Sprintf("%d-%s", DBT.Num, DBT.DBT)
		DBTs = append(DBTs, fullStr)
	}

	DBTSelector := widget.NewSelect(DBTs, func(s string) {
		log.Printf("Selected: %s", s)
	})

	refreshSelector := func() {
		if len(DBTs) > 0 {
			return
		}
		data, err := storage.GetFromStorage()
		if err != nil {
			data = &StorageStructure{}
		}
		for _, DBT := range data.DBTs {
			fullStr := fmt.Sprintf("%d-%s", DBT.Num, DBT.DBT)
			DBTs = append(DBTs, fullStr)
		}
		DBTSelector.Options = DBTs
		DBTSelector.Refresh()
	}

	refreshSelector()

	input := widget.NewEntry()
	input.PlaceHolder = "Enter name"

	dropArea := widget.NewLabel("Drag file here")
	dropArea.Alignment = fyne.TextAlignCenter

	border := canvas.NewRectangle(color.RGBA{R: 0, G: 122, B: 255, A: 255})
	border.StrokeWidth = 2.0

	dropContainer := container.NewStack(border, container.NewCenter(dropArea))

	chooseFileButton := widget.NewButton("choose file manually", func() {})

	platformSelector := widget.NewSelect([]string{"Email", "Arhimed", "Hermes", "Regix"}, func(s string) {

	})

	addFileToFolderBtn := widget.NewButton("Add file to folder", func() {

	})

	grid := container.NewVBox(container.New(
		layout.NewGridLayout(2), DBTSelector, input, chooseFileButton, platformSelector), dropContainer, addFileToFolderBtn)

	return container.NewBorder(title, nil, nil, nil, grid), refreshSelector
}
