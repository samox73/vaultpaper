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
	if err := ctl.Load(); err != nil {
		log.Fatal(err)
	}

	random := flag.Bool("r", false, "use random picture of currenct location")
	list := flag.Bool("l", false, "list all locations")
	newLocation := flag.String("a", "", "add a folder as a new location")
	flag.Parse()

	if *random {
		ctl.Random()
	}
	if *newLocation != "" {
		ctl.AddFolder(*newLocation)
	}
	if *list {
		ctl.ListLocations()
	}

	if err := ctl.Save(); err != nil {
		log.Fatal(err)
	}
}
