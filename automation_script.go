package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	// "fyne.io/fyne/v2/container"
	// "fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	excel "github.com/xuri/excelize/v2"
)

type Row struct {
	DBT   string
	Email string
	Num   int
}

func getDesktopPath() (string, error) {
	switch runtime.GOOS {
	case "windows":
		userProfile := os.Getenv("USERPROFILE")
		if userProfile == "" {
			return "", os.ErrNotExist
		}
		return filepath.Join(userProfile, "Desktop"), nil
	case "linux":
		home := os.Getenv("HOME")
		windowsDesktop := filepath.Join("/mnt/c/Users", filepath.Base(home), "Desktop")
		if _, err := os.Stat(windowsDesktop); err == nil {
			return windowsDesktop, nil
		}
	}
	return "", os.ErrNotExist
}

func execFolderScript(w fyne.Window) ([]Row, error) {
	desktopPath, err := getDesktopPath()
	if err != nil {
		return []Row{}, err
	}
	currDir, err := os.Getwd()
	if err != nil {
		return []Row{}, errors.New("cant get current dir")
	}
	filename := "DBT_s_imeili.xlsx"
	log.Printf("curr: %s", currDir)
	var excelFilePath string
	err = filepath.Walk(currDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.EqualFold(info.Name(), filename) {
			log.Printf("Found file: %s\n", info.Name())
			excelFilePath = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return []Row{}, fmt.Errorf("error walking path: %s", currDir)
	}
	excelFile, err := excel.OpenFile(excelFilePath)
	if err != nil {
		return []Row{}, fmt.Errorf("error opening %s: %s", filename, err.Error())
	}
	defer func() {
		if err := excelFile.Close(); err != nil {
			log.Printf("Error while closing file(%s) with error: %s\n", filename, err.Error())
		}
	}()
	sheets := excelFile.GetSheetList()
	neededSheet := sheets[0]
	rows, err := excelFile.GetRows(neededSheet)
	if err != nil {
		return []Row{}, fmt.Errorf("can't get rows on file: %s", filename)
	}
	DBT_FOLDER := filepath.Join(desktopPath, "ДБТ")
	if _, err = os.Stat(DBT_FOLDER); err == nil {
		log.Println("Folder exists")
		log.Println("Script is being canceled")
		return []Row{}, errors.New("папката вече съществува")
	} else if os.IsNotExist(err) {
		log.Println("Folder does not exist")
		log.Println("Folder is being created")
	} else {
		log.Println("Script is being canceled")
		return []Row{}, fmt.Errorf("проблем с намирането на папката")
	}
	err = os.MkdirAll(DBT_FOLDER, 0666)
	if err != nil {
		return []Row{}, fmt.Errorf("проблем със създаването на папката")
	}
	log.Println("ДБТ folder created")
	emailFilePath := filepath.Join(DBT_FOLDER, "ДБТ-имейли.txt")
	emailsFile, err := os.OpenFile(emailFilePath, os.O_RDWR|os.O_CREATE, 0666)
	emailsFileName := filepath.Base(emailFilePath)
	if err != nil {
		return []Row{}, fmt.Errorf("cannot open %s", emailsFileName)
	}
	defer emailsFile.Close()

	log.Println("ДБТ-имейли.txt created")
	rowsSlice := []Row{}
	for _, row := range rows {
		DBT := row[1]
		email := row[2]
		if DBT == "" || email == "" {
			continue
		}
		slice1 := strings.Split(email, "-")
		slice2 := strings.Split(slice1[1], "@")
		numStr := slice2[0]
		var num int
		fmt.Sscanf(numStr, "%d", &num)
		rowsSlice = append(rowsSlice, Row{
			DBT:   DBT,
			Email: email,
			Num:   num,
		})
	}
	sort.Slice(rowsSlice, func(i, j int) bool {
		return rowsSlice[i].Num < rowsSlice[j].Num
	})

	for _, row := range rowsSlice {
		if row.Email != "" {
			_, err := emailsFile.Write([]byte(fmt.Sprintf("%s - %s\n", row.DBT, row.Email)))
			if err != nil {
				log.Printf("Error while adding email: %s for DBT: %s\n", row.Email, row.DBT)
			}
			cityDBTPath := filepath.Join(DBT_FOLDER, fmt.Sprintf("%d-%s", row.Num, row.DBT))
			folder := filepath.Base(cityDBTPath)
			if err = os.Mkdir(cityDBTPath, os.ModePerm); err != nil {
				log.Printf("Error creating: %s", folder)
			} else {
				for _, platform := range PLATFORMS {
					platformPath := filepath.Join(cityDBTPath, platform)
					if err = os.Mkdir(platformPath, os.ModePerm); err != nil {
						log.Printf("Error creating platform %s for %s\n", platform, folder)
					}
					for _, request := range REQUEST_TYPES {
						requestPath := filepath.Join(platformPath, request)
						if err = os.Mkdir(requestPath, os.ModePerm); err != nil {
							log.Printf("Error creating request folder %s for %s at %s\n", request, platform, folder)
						}
					}
				}
				log.Printf("Folder: %s created", folder)
			}

		}
	}

	AZ_Folder_Path := filepath.Join(DBT_FOLDER, "АЗ")
	err = os.Mkdir(AZ_Folder_Path, os.ModePerm)
	if err != nil {
		log.Println("Error creating AZ folder")
	}
	for _, platform := range PLATFORMS {
		platformPath := filepath.Join(AZ_Folder_Path, platform)
		if err = os.Mkdir(platformPath, os.ModePerm); err != nil {
			log.Printf("Error creating platform %s for %s\n", platform, AZ_Folder_Path)
		}
		for _, request := range REQUEST_TYPES {
			requestPath := filepath.Join(platformPath, request)
			if err = os.Mkdir(requestPath, os.ModePerm); err != nil {
				log.Printf("Error creating request folder %s for %s at %s\n", request, platform, AZ_Folder_Path)
			}
		}
	}

	dialog.ShowInformation("Нотификация", "Папката е създадена успешно", w)
	return rowsSlice, nil
}
