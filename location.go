package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"path/filepath"
)

type location struct {
	Directory string
	Files     []string
}

func NewLocation(directory string) location {
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	return location{Directory: directory}
}

func (l *location) Scan() {
	fmt.Printf("Scanning directory %s\n", l.Directory)
	filepath.WalkDir(l.Directory, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && isLegalFile(path) {
			fmt.Printf("Found file %s (%s)\n", path, err)
			l.Files = append(l.Files, path)
		}
		return err
	})
}

func (l location) Print() {
	fmt.Printf("Directory: %s\n", l.Directory)
	fmt.Println("Files:")
	for _, file := range l.Files {
		fmt.Printf("  - %s\n", file)
	}
}

func (l location) Save(filename string) error {
	if filename == "" {
		return errors.New("empty filename")
	}
	dataFile, err := os.Create(l.Directory + "/" + filename)
	if err != nil {
		return err
	}
	dataEncoder := gob.NewEncoder(dataFile)
	err = dataEncoder.Encode(l)
	dataFile.Close()
	fmt.Printf("Saved location to file '%s'\n", filename)
	return err
}

func (l *location) Load(filename string) error {
	if filename == "" {
		return errors.New("empty filename")
	}
	fmt.Printf("Loading location %s\n", filename)
	dataFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&l)
	dataFile.Close()
	fmt.Printf("Loaded location from file '%s'\n", filename)
	return err
}

func (l location) GetRandomFilePath() (string, error) {
	size := len(l.Files)
	if size == 0 {
		return "", errors.New("Location " + l.Directory + " has 0 files")
	}
	randomIndex := rand.Intn(size)
	return l.Files[randomIndex], nil
}