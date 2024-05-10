package engine

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

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
		if len(line) > 0 {
			data := strings.Split(line, " ")
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
	}

	return results
}

func GetPage(url string) string {
	var response string

	for n := 0; n <= 2; n++ {
		client := http.Client{Timeout: time.Duration(10) * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			fmt.Println("Get Error:", url)
		} else {
			if resp.StatusCode == 200 {
				response_, _ := io.ReadAll(resp.Body)
				response = string(response_)
				resp.Body.Close()
				break
			} else {
				fmt.Println("Page not available:", url)
				response = "page not available"
				resp.Body.Close()
				break
			}
		}
	}

	return response
}
