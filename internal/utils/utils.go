package utils

import (
	"fmt"
	"os"
)

func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateDir(pathDir string) error {
	err := os.Mkdir(pathDir, 0777)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
