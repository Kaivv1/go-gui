package main

import (
	"fmt"
	"image/color"
	"io"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"

	// "fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

//	type FileToAdd struct {
//		DBT string
//		RequestType string
//		Name string
//		File
//	}

func makeAddFilesSpace(storage *Storage, w fyne.Window) (fyne.CanvasObject, func()) {
	var chosenFilePath string

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
	DBTSelector.Selected = "Избери ДБТ"

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

	chooseFileButton := &widget.Button{}

	chooseFileButton = widget.NewButton("Избери файл", func() {
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
			chooseFileButton.SetText("Има избран файл")

			chosenFilePath = filePath
			w.Content().Refresh()
		}, w)
		d.Resize(fyne.NewSize(950, 800))
		d.Show()
		d.SetOnClosed(func() {
			chooseFileButton.SetText("Избери файл")
			w.Content().Refresh()
		})
		w.Content().Refresh()
	})

	platformSelector := widget.NewSelect(PLATFORMS[:], func(s string) {
	})
	platformSelector.PlaceHolder = "Избери платформа"

	actionSelector := widget.NewSelect(REQUEST_TYPES[:], func(s string) {

	})
	actionSelector.PlaceHolder = "Вид заявка"

	resetButton := widget.NewButton("Ресетни", func() {
		chosenFilePath = ""
		platformSelector.Selected = ""
		actionSelector.Selected = ""
		DBTSelector.Selected = ""
		input.Text = ""
		chooseFileButton.SetText("Избери файл")

		DBTSelector.Refresh()
		actionSelector.Refresh()
		platformSelector.Refresh()
		input.Refresh()
	})

	addFileToFolderBtn := widget.NewButton("Добави файла към папката", func() {
		desktopPath, err := getDesktopPath()
		if err != nil {
			log.Println("Error getting error path when adding file to folder")
		}
		extension := filepath.Ext(chosenFilePath)

		newFilePath := filepath.Join(desktopPath, "ДБТ", DBTSelector.Selected, platformSelector.Selected, actionSelector.Selected, fmt.Sprintf("%s%s", input.Text, extension))

		chosenFile, err := os.Open(chosenFilePath)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		defer chosenFile.Close()

		newFile, err := os.Create(newFilePath)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, chosenFile)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		chooseFileButton.SetText("Избери файл")
		dialog.ShowInformation("Нотификация", "Файлът е добавен към папката", w)

		log.Printf("File path: %s\n", newFilePath)
		log.Printf("File: %s\n", filepath.Base(chosenFilePath))

	})
	grid := container.NewVBox(container.NewBorder(
		nil, nil, chooseFileButton, container.NewVBox(actionSelector, resetButton, addFileToFolderBtn), container.NewVBox(DBTSelector, platformSelector, input)))

	return container.NewBorder(title, nil, nil, nil, grid), refreshSelector
}
