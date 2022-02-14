package main

import (
	"log"
	"os"
	"strings"
)

var legalTypes []string = []string{"png", "jpg", "jpeg"}

func isLegalFile(path string) bool {
	for _, typeSuffix := range legalTypes {
		if strings.HasSuffix(path, "."+typeSuffix) {
			return true
		}
	}
	return false
}

// exists returns whether the given file or directory exists
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	log.Fatal(err)
	return false
}

func contains(locations []location, path string) bool {
	for _, location := range locations {
		if location.Directory == path {
			return true
		}
	}
	return false
}

func GetStoreDir() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dirname := ".vault"
	return homedir + "/" + dirname
}
