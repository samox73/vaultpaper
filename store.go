package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
)

type store struct {
	Locations      []location
	ActiveLocation *location
	FileName       string
	Backend        backend
}

func NewStore() store {
	gob.Register(fehBackend{})
	gob.Register(pywalBackend{})
	return store{FileName: "config", Backend: NewPywalBackend()}
}

func (s *store) AddLocation(path string) (*location, error) {
	if path == "" {
		return nil, errors.New("cannot add empty path")
	}
	location := NewLocation(path)
	s.Locations = append(s.Locations, location)
	if len(s.Locations) == 1 {
		s.ActiveLocation = &s.Locations[0]
		fmt.Printf("Active location of store: %s\n", path)
	}
	return &s.Locations[len(s.Locations)-1], nil
}

func (s store) CreateConfig() {
	configDir := GetStoreDir()
	filePath := configDir + "/" + s.FileName
	fmt.Printf("Saving store to %s\n", filePath)
	dataFile, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	dataEncoder := gob.NewEncoder(dataFile)
	if err := dataEncoder.Encode(s); err != nil {
		log.Fatal(err)
		return
	}
	dataFile.Close()
}

func (s *store) LoadStoreFile(filePath string) *os.File {
	fmt.Printf("Loading store file %s\n", filePath)
	dataFile, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Error: store config not found, creating default config")
			s.CreateConfig()
			return s.LoadStoreFile(filePath)
		} else {
			log.Fatal(err)
		}
	}
	return dataFile
}

func (s store) Save() {
	configDir := GetStoreDir()
	if !exists(configDir) {
		err := os.Mkdir(configDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	s.CreateConfig()
}

func (s *store) Load() {
	configDir := GetStoreDir()
	if !exists(configDir) {
		s.LoadDefault()
		return
	}
	filePath := configDir + "/" + s.FileName
	dataFile := s.LoadStoreFile(filePath)
	dataDecoder := gob.NewDecoder(dataFile)
	if err := dataDecoder.Decode(&s); err != nil {
		dataFile.Close()
		log.Fatal(err)
	}
	dataFile.Close()
	for _, location := range s.Locations {
		location.Load()
	}

}

func (s *store) LoadDefault() {
	fmt.Println("Inizializing data directory")
	configDir := GetStoreDir()
	if err := os.Mkdir(configDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}
