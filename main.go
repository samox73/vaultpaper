package main

import (
	"flag"
	"log"
	"math/rand"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	rand.Seed(time.Now().UnixNano())

	ctl := NewVaultctl()
	ctl.Load()

	random := flag.Bool("r", false, "use random picture of currenct location")
	list := flag.Bool("l", false, "list all locations")
	verbose := flag.Bool("v", false, "verbose")
	newLocation := flag.String("a", "", "add a folder as a new location")
	useBackend := flag.String("b", "", "use a specific backend, possible values include:\n  - pywal\n  - feh")
	flag.Parse()

	if *random {
		ctl.Random()
	} else if *newLocation != "" {
		ctl.AddFolder(*newLocation)
	} else if *list {
		ctl.ListLocations(*verbose)
	} else if *useBackend != "" {
		ctl.UseBackend(*useBackend)
	}

	ctl.Save()
}
