package main

import (
	"fmt"
	"os/exec"
)

type pywalBackend struct{}

func NewPywalBackend() pywalBackend {
	return pywalBackend{}
}

func (backend pywalBackend) setWallpaper(filepath string) error {
	fmt.Printf("Pywal: setting wallpaper to %s\n", filepath)
	cmd := exec.Command("wal", "-i", filepath)
	return cmd.Run()
}
