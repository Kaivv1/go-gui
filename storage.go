package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Storage struct {
	Path string
}

type StorageStructure struct {
	DBTs []Row
}

func NewStorage(path string) (*Storage, error) {
	s := &Storage{
		Path: path,
	}
	err := s.DoesStorageExist()
	if err != nil {
		return &Storage{}, err
	}
	return s, nil
}

func (s *Storage) DoesStorageExist() error {
	currDir, err := os.Getwd()
	storageFullPath := filepath.Join(currDir, s.Path)
	if _, err = os.Stat(storageFullPath); os.IsNotExist(err) {
		_, err = os.Create(storageFullPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) WriteToStorage(args *StorageStructure) error {
	data := &StorageStructure{
		DBTs: args.DBTs,
	}
	dataInBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = os.WriteFile(s.Path, dataInBytes, 0666)
	if err != nil {
		return errors.New("Cannot write to storage file")
	}
	return nil
}

func (s *Storage) GetFromStorage() (*StorageStructure, error) {
	dataInBytes, err := os.ReadFile(s.Path)
	if err != nil {
		return &StorageStructure{}, errors.New("Cannot read storage file")
	}
	storageStructure := &StorageStructure{}

	err = json.Unmarshal(dataInBytes, storageStructure)
	if err != nil {
		return &StorageStructure{}, err
	}

	return storageStructure, nil
}
