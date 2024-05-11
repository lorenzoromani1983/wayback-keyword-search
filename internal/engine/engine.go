package engine

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var blockByTypes = map[string]struct{}{
	"image/jpeg":               {},
	"image/png":                {},
	"image/gif":                {},
	"warc/revisit":             {},
	"text/css":                 {},
	"application/javascript":   {},
	"image/vnd.microsoft.icon": {},
}

func GetHistory(targetDomain string, timeStamp string) []string {
	API_URL := "https://web.archive.org/cdx/search/cdx?url=*."
	query_url := API_URL + targetDomain
	client := http.Client{Timeout: time.Duration(50) * time.Second}
	response, err := client.Get(query_url)
	if err != nil {
		panic("can't download history for the specified domain ")
	}

	body, _ := io.ReadAll(response.Body)
	content := string(body)
	lines := strings.Split(string(content), "\n")
	BASE_URL := "http://web.archive.org/web/"

	var results []string

	for _, line := range lines {

		if len(line) == 0 {
			continue
		}
		data := strings.Split(line, " ")

		if len(data) != 7 || data[4] != "200" || IsValueExists(data[3], blockByTypes) {
			continue
		}

		if strings.Contains(data[0], ")/") == true {
			savedpage := strings.Split(data[0], ")/")[1]
			url := targetDomain + "/" + savedpage
			timestamp := string(data[1])
			if strings.HasPrefix(timestamp, timeStamp) {
				wayback_url := BASE_URL + timestamp + "/" + url
				results = append(results, wayback_url)
			}
		}
	}

	return results
}

func GetPage(url string) (string, error) {

	for n := 0; n <= 2; n++ {
		client := http.Client{Timeout: time.Duration(10) * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}

		if resp.StatusCode == http.StatusOK {
			var response string
			responseByte, err := io.ReadAll(resp.Body)
			if err == nil {
				response = string(responseByte)
			}
			resp.Body.Close()
			return response, nil
		}

		continue
	}

	return "", errors.New("without any response data")
}

func IsValueExists(target string, list map[string]struct{}) bool {
	if _, ok := list[target]; ok {
		return true
	} else {
		return false
	}
}
