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

type localLocation struct {
	Directory    string
	Files        []string
	CurrentIndex int
}

func NewLocalLocation(directory string) localLocation {
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	return localLocation{Directory: directory, CurrentIndex: -1}
}

func (l *localLocation) Scan() {
	fmt.Printf("Scanning directory %s\n", l.Directory)
	filepath.WalkDir(l.Directory, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && isLegalFile(path) {
			fmt.Printf("  - %s\n", filepath.Base(path))
			l.Files = append(l.Files, path)
		}
		return err
	})
}

func (l localLocation) Print() {
	fmt.Printf("Directory: %s\n", l.Directory)
	fmt.Println("Files:")
	l.PrintFiles(2)
}

func (l localLocation) PrintFiles(indent int) {
	for _, file := range l.Files {
		basename := filepath.Base(file)
		fmt.Printf("%s- %s\n", strings.Repeat(" ", indent), basename)
	}
}

func (l localLocation) getConfigPath() string {
	return l.Directory + "/.vaultlocation"
}

func (l localLocation) Save() {
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

func (l *localLocation) Load() {
	filepath := l.getConfigPath()
	fmt.Printf("Loading location %s\n", filepath)
	dataFile, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&l)
	if err != nil {
		dataFile.Close()
		log.Fatal(err)
	}
	dataFile.Close()
}

func (l localLocation) GetRandomFilePath() (string, error) {
	size := len(l.Files)
	if size == 0 {
		return "", errors.New("Location " + l.Directory + " has 0 files")
	}
	randomIndex := rand.Intn(size)
	return l.Files[randomIndex], nil
}
