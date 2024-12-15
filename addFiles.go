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

	return container.NewBorder(title, nil, nil, nil, container.NewVBox(container.New(layout.NewGridLayout(2), DBTSelector, input))), refreshSelector
}
