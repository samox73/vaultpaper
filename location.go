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
	"strings"
)

type location struct {
	Directory    string
	Files        []string
	CurrentIndex int
}

func NewLocation(directory string) location {
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	return location{Directory: directory, CurrentIndex: -1}
}

func (l *location) Scan() {
	fmt.Printf("Scanning directory %s\n", l.Directory)
	filepath.WalkDir(l.Directory, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && isLegalFile(path) {
			fmt.Printf("  - %s\n", filepath.Base(path))
			l.Files = append(l.Files, path)
		}
		return err
	})
}

func (l location) Print() {
	fmt.Printf("Directory: %s\n", l.Directory)
	fmt.Println("Files:")
	l.PrintFiles(2)
}

func (l location) PrintFiles(indent int) {
	for _, file := range l.Files {
		basename := filepath.Base(file)
		fmt.Printf("%s- %s\n", strings.Repeat(" ", indent), basename)
	}
}

func (l location) getConfigPath() string {
	return l.Directory + "/.vaultlocation"
}

func (l location) Save() {
	filepath := l.getConfigPath()
	dataFile, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	dataEncoder := gob.NewEncoder(dataFile)
	err = dataEncoder.Encode(l)
	if err != nil {
		log.Fatal(err)
	}
	dataFile.Close()
	fmt.Printf("Saved location to file '%s'\n", filepath)
}

func (l *location) Load() {
	filepath := l.getConfigPath()
	fmt.Printf("Loading location %s\n", filepath)
	dataFile, err := os.Open(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			l.Save()
			l.Load()
			return
		} else {
			log.Fatal(err)
		}
	}
	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&l)
	if err != nil {
		dataFile.Close()
		log.Fatal(err)
	}
	dataFile.Close()
}

func (l location) GetRandomFilePath() (string, error) {
	size := len(l.Files)
	if size == 0 {
		return "", errors.New("Location " + l.Directory + " has 0 files")
	}
	randomIndex := rand.Intn(size)
	return l.Files[randomIndex], nil
}
