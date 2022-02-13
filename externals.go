package main

import (
	"fmt"
	"os/exec"
)

func SetBackground(filepath string) error {
	fmt.Printf("Setting wallpaper to %s\n", filepath)
	cmd := exec.Command("wal", "-i", filepath)
	return cmd.Run()
}
