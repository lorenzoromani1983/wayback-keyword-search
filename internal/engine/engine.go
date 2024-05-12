package engine

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const API_URL = "https://web.archive.org/cdx/search/cdx?url=*."
const BASE_URL = "http://web.archive.org/web/"

var blockByTypes = map[string]struct{}{
	"image/jpeg":               {},
	"image/png":                {},
	"image/gif":                {},
	"warc/revisit":             {},
	"text/css":                 {},
	"application/javascript":   {},
	"image/vnd.microsoft.icon": {},
}

func GetHistory(targetDomain string, timeStamp string) ([]string, error) {
	query_url := API_URL + targetDomain
	client := http.Client{Timeout: time.Duration(50) * time.Second}
	response, err := client.Get(query_url)
	if err != nil {
		return nil, fmt.Errorf("can't download history for the specified domain: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	content := string(body)
	lines := strings.Split(string(content), "\n")

	results := make([]string, 0, len(lines))

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

	return results, nil
}

func GetPage(ctx context.Context, url string) (string, error) {
	for n := 0; n <= 2; n++ {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return "", fmt.Errorf("error creating request: %w", err)
		}

		response, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("error making request: %w", err)
		}

		if response.StatusCode != http.StatusOK {
			continue
		}

		data, err := io.ReadAll(response.Body)
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}

		result := string(data)
		response.Body.Close()

		return result, nil
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
