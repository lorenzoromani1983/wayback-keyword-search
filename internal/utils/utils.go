package utils

import (
	"fmt"
	"net/url"
	"os"
	"strings"
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

func UrlToFileName(inputUrl string) (fileName string) {
	invalidChars := []string{"/", "?", ":", "@", "&", "=", "+", "$", ","}
	for _, char := range invalidChars {
		inputUrl = strings.ReplaceAll(inputUrl, char, "_")
	}

	fileName, _ = url.QueryUnescape(inputUrl)
	fileName = strings.ReplaceAll(fileName, " ", "")

	return
}