package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
)

type store struct {
	Locations      []location
	ActiveLocation location
	FileName       string
}

func NewStore() store {
	return store{FileName: "config"}
}

func (s *store) AddLocation(path string) error {
	if path == "" {
		return errors.New("cannot add empty path")
	}
	location := NewLocation(path)
	if len(s.Locations) == 0 {
		s.ActiveLocation = location
		fmt.Printf("Active location of store: %s\n", path)
	}
	s.Locations = append(s.Locations, location)
	return nil
}

func GetStoreDir() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dirname := ".vault"
	return homedir + "/" + dirname, nil
}

func (s store) Save() error {
	configDir, err := GetStoreDir()
	if err != nil {
		return err
	}
	configDirExists, err := exists(configDir)
	if err != nil {
		return err
	}
	if !configDirExists {
		if err := os.Mkdir(configDir, os.ModePerm); err != nil {
			return err
		}
	}
	filePath := configDir + "/" + s.FileName
	fmt.Printf("Saving store to %s\n", filePath)
	dataFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	dataEncoder := gob.NewEncoder(dataFile)
	if err := dataEncoder.Encode(s); err != nil {
		return err
	}
	dataFile.Close()
	return nil
}

func (s *store) Load() error {
	configDir, err := GetStoreDir()
	if err != nil {
		return err
	}
	configDirExists, err := exists(configDir)
	if err != nil {
		return err
	}
	if !configDirExists {
		if err := s.LoadDefault(); err != nil {
			return err
		}
		return nil
	}
	filePath := configDir + "/" + s.FileName
	fmt.Printf("Loading data file %s\n", filePath)
	dataFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&s)
	dataFile.Close()
	fmt.Printf("Loaded store from file '%s'\n", filePath)
	return err
}
func (s *store) LoadDefault() error {
	fmt.Println("Inizializing data directory")
	configDir, err := GetStoreDir()
	if err != nil {
		return err
	}
	if err := os.Mkdir(configDir, os.ModePerm); err != nil {
		return err
	}
	return nil

}
