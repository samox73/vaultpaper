package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type vaultctl struct {
	store store
}

func NewVaultctl() vaultctl {
	return vaultctl{store: NewStore()}
}

func (v *vaultctl) UseBackend(backend string) error {
	fmt.Printf("Setting backend to %s\n", backend)
	switch backend {
	case "pywal":
		v.store.Backend = NewPywalBackend()
	case "feh":
		v.store.Backend = NewFehBackend()
	default:
		return errors.New("Backend '" + backend + "' is not a valid backend")
	}
	return nil
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
	v.store.Backend.setWallpaper(randomFile)
}

func (v *vaultctl) AddLocation(uri string) {
	if uri == "" {
		log.Fatal(errors.New("cannot add empty path"))
	}
	if err := v.store.AddLocation(uri); err != nil {
		log.Fatal(err)
	}
}

func (v *vaultctl) Load() {
	v.store.Load()
}

func (v *vaultctl) Save() {
	v.store.Save()
}

func (v vaultctl) ListLocations(verbose bool) {
	fmt.Println("Current configured locations:")
	for _, location := range v.store.LocalLocations {
		fmt.Printf("  *) %s\n", location.Directory)
		location.PrintFiles(4)
	}
}
