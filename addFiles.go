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

func AddBorder(w fyne.CanvasObject, c color.Color) *fyne.Container {
	border := canvas.NewRectangle(c)
	border.StrokeWidth = 2
	border.StrokeColor = c
	return container.NewBorder(border, nil, nil, nil, w)
}

func makeAddFilesSpace(storage *Storage, w fyne.Window) (fyne.CanvasObject, func()) {
	title := canvas.NewText("Добави заявки", color.White)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = TITLE_TEXT

	data, err := storage.GetFromStorage()
	if err != nil {
		data = &StorageStructure{}
	}
	DBTs := []string{"АЗ"}
	for _, DBT := range data.DBTs {
		fullStr := fmt.Sprintf("%d-%s", DBT.Num, DBT.DBT)
		DBTs = append(DBTs, fullStr)
	}

	DBTSelector := widget.NewSelect(DBTs, func(s string) {
	})
	DBTSelector.PlaceHolder = "Избери ДБТ"

	DBTSelectorContainer := AddBorder(DBTSelector, color.RGBA{R: 255, G: 0, B: 0, A: 0})

	refreshSelector := func() {
		if len(DBTs) > 1 {
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
	input.PlaceHolder = "Име"

	chooseFileButton := widget.NewButton("Избери файл", func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()

			filePath := reader.URI().Path()
			log.Printf("Selected file: %s", filePath)
		}, w)
		d.Resize(fyne.NewSize(950, 800))
		d.Show()
	})

	platformSelector := widget.NewSelect(PLATFORMS[:], func(s string) {

	})

	platformSelector.PlaceHolder = "Избери платформа"
	addFileToFolderBtn := widget.NewButton("Добави файла към папката", func() {
		log.Printf("DBT: %s", DBTSelector.Selected)
		log.Printf("Request type: %s", platformSelector.Selected)
		log.Printf("Input: %s", input.Text)
	})
	resetButton := widget.NewButton("Ресетни", func() {})

	actionSelector := widget.NewSelect(REQUEST_TYPES[:], func(s string) {

	})
	actionSelector.PlaceHolder = "Вид заявка"
	grid := container.NewVBox(container.NewBorder(
		nil, nil, chooseFileButton, container.NewVBox(actionSelector, resetButton, addFileToFolderBtn), container.NewVBox(DBTSelectorContainer, platformSelector, input)))

	return container.NewBorder(title, nil, nil, nil, grid), refreshSelector
}
