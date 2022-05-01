package main

import (
	"fmt"
	"os/exec"
)

type fehBackend struct{}

func NewFehBackend() fehBackend {
	return fehBackend{}
}

func (backend fehBackend) setWallpaper(filepath string) error {
	fmt.Printf("Feh: setting wallpaper to %s\n", filepath)
	cmd := exec.Command("feh", "--bg-fill", filepath)
	return cmd.Run()
}
