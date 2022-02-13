package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type vaultctl struct {
	store store
}

func NewVaultctl() vaultctl {
	return vaultctl{store: NewStore()}
}

func (v vaultctl) Random() {
	if v.store.ActiveLocation == nil {
		fmt.Println("The store has no active location configured")
		os.Exit(1)
	}
	randomFile, err := v.store.ActiveLocation.GetRandomFilePath()
	if err != nil {
		log.Fatal(err)
	}
	SetBackground(randomFile)
}

func (v *vaultctl) AddFolder(path string) {
	var err error
	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			log.Fatal(err)
		}
	}
	if contains(v.store.Locations, path) {
		fmt.Printf("Store already tracks folder %s\n", path)
		return 
	}
	fmt.Printf("Adding new folder %s\n", path)
	newLocation, err := v.store.AddLocation(path)
	if err != nil {
		log.Fatal(err)
	}
	newLocation.Scan()
	newLocation.Save()
}

func (v *vaultctl) Load() {
	v.store.Load()
}

func (v *vaultctl) Save() {
	v.store.Save()
}

func (v vaultctl) ListLocations(verbose bool) {
	fmt.Println("Current configured locations:")
	for _, location := range v.store.Locations {
		fmt.Printf("  *) %s\n", location.Directory)
		location.PrintFiles(4)
	}
}
