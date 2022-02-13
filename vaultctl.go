package main

import (
	"fmt"
	"path/filepath"
)

type vaultctl struct {
	store store
}

func NewVaultctl() vaultctl {
	return vaultctl{store: NewStore()}
}

func (v vaultctl) Random() error {
	randomFile, err := v.store.ActiveLocation.GetRandomFilePath()
	if err != nil {
		return err
	}
	SetBackground(randomFile)
	return nil
}

func (v *vaultctl) AddFolder(path string) error {
	var err error
	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			return nil
		}
	}
	fmt.Printf("Adding new folder %s\n", path)
	v.store.AddLocation(path)
	return nil
}

func (v *vaultctl) Load() error {
	return v.store.Load()
}

func (v *vaultctl) Save() error {
	return v.store.Save()
}

func (v vaultctl) ListLocations() {
	fmt.Println("Current configured locations:")
	for _, location := range v.store.Locations {
		fmt.Printf("  - %s\n", location.Directory)
	}
}
