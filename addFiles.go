package main

import (
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"

	// "fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// type FileToAdd struct {
// 	DBT string
// 	RequestType string
// 	Name string
// 	File
// }

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
	})
	DBTSelector.PlaceHolder = "Choose DBT"
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

	chooseFileButton := widget.NewButton("choose file manually", func() {
		dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()

			// Get file path or URI
			filePath := reader.URI().Path()
			log.Printf("Selected file: %s", filePath)
		}, w).Show()
	})

	platformSelector := widget.NewSelect([]string{"Email", "Arhimed", "Hermes", "Regix"}, func(s string) {

	})
	platformSelector.PlaceHolder = "Choose request type"
	addFileToFolderBtn := widget.NewButton("Add file to folder", func() {
		log.Printf("DBT: %s", DBTSelector.Selected)
		log.Printf("Request type: %s", platformSelector.Selected)
		log.Printf("Input: %s", input.Text)
	})
	resetButton := widget.NewButton("Reset", func() {})

	actionSelector := widget.NewSelect([]string{"Предоставяне", "Променяне", "Прекратяване"}, func(s string) {

	})
	grid := container.NewVBox(container.NewBorder(
		nil, nil, chooseFileButton, container.NewVBox(actionSelector, addFileToFolderBtn, resetButton), container.NewVBox(DBTSelector, platformSelector, input)))

	return container.NewBorder(title, nil, nil, nil, grid), refreshSelector
}
