package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
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

func execFolderScript(a fyne.App) {

	desktopPath, err := getDesktopPath()
	if err != nil {
		log.Println(err)
	}

	currDir, err := os.Getwd()
	if err != nil {
		log.Println("Cant get current dir")
		return
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
		log.Printf("Error walking path: %s\n", currDir)
		return
	}
	excelFile, err := excel.OpenFile(excelFilePath)
	if err != nil {
		log.Printf("Error opening %s: %s", filename, err.Error())
		return
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
		log.Printf("Can't get rows on file: %s", filename)
		return
	}
	DBT_FOLDER := filepath.Join(desktopPath, "ДБТ")

	if _, err = os.Stat(DBT_FOLDER); err == nil {
		log.Println("Folder exists")
		log.Fatalln("Script is being canceled")
		return
	} else if os.IsNotExist(err) {
		log.Println("Folder does not exist")
		log.Println("Folder is being created")
	} else {
		log.Println("Error searching for folder")
		log.Fatalln("Script is being canceled")
		return
	}

	err = os.MkdirAll(DBT_FOLDER, 0666)
	if err != nil {
		log.Println("Error creating DBT folder")
		return
	}
	log.Println("ДБТ folder created")

	emailFilePath := filepath.Join(DBT_FOLDER, "ДБТ-имейли.txt")

	emailsFile, err := os.OpenFile(emailFilePath, os.O_RDWR|os.O_CREATE, 0666)
	emailsFileName := filepath.Base(emailFilePath)
	if err != nil {
		log.Fatalf("Cannot open %s\n", emailsFileName)
	}

	defer func() {
		if err := emailsFile.Close(); err != nil {
			log.Fatalf("Error closing %s\n", emailsFileName)
		}
	}()

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
				log.Printf("Folder: %s created", folder)
			}
		}
	}

	AZ_Folder_Path := filepath.Join(DBT_FOLDER, "AZ")
	err = os.Mkdir(AZ_Folder_Path, os.ModePerm)
	if err != nil {
		log.Println("Error creating AZ folder")
	}
	a.SendNotification(fyne.NewNotification("Ready", "Your folder is ready"))
}
