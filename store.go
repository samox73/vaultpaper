package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strings"
)

type store struct {
	RedditLocations []redditLocation
	LocalLocations  []localLocation
	ActiveLocation  location
	FileName        string
	Backend         backend
}

func NewStore() store {
	gob.Register(fehBackend{})
	gob.Register(pywalBackend{})
	gob.Register(redditLocation{})
	gob.Register(localLocation{})
	return store{FileName: "config", Backend: NewPywalBackend()}
}

func (s *store) AddLocation(uri string) error {
	if strings.HasPrefix(uri, "r/") {
		return s.AddRedditLocation(uri)
	} else {
		return s.AddLocalLocation(uri)
	}
}

func (s *store) AddRedditLocation(sub string) error {
	fmt.Printf("Adding subreddit '%s'", sub)
	if subredditIsPresent(s.RedditLocations, sub) {
		return &LocationPresentError{sub}
	}
	s.RedditLocations = append(s.RedditLocations, NewRedditLocation(sub))
	if s.ActiveLocation == nil {
		s.ActiveLocation = &s.RedditLocations[0]
		fmt.Printf("Active location of store is now '%s'\n", sub)
	}
	return nil
}

func (s *store) AddLocalLocation(path string) error {
	path = mapAbsPath(path)
	if localLocationIsPresent(s.LocalLocations, path) {
		return &LocationPresentError{path}
	}
	fmt.Printf("Adding local location '%s'\n", path)
	location := NewLocalLocation(path)
	s.LocalLocations = append(s.LocalLocations, location)
	if s.ActiveLocation == nil {
		s.ActiveLocation = &s.LocalLocations[0]
		fmt.Printf("Active location of store is now '%s'\n", path)
	}
	newLocation := &s.LocalLocations[len(s.LocalLocations)-1]
	newLocation.Scan()
	newLocation.Save()
	return nil
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
	for _, location := range s.LocalLocations {
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

type LocationPresentError struct {
	uri string
}

func (e *LocationPresentError) Error() string {
	return fmt.Sprintf("Location '%s' is already present", e.uri)
}
