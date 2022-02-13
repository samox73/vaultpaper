package main

import (
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
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
