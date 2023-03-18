package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"context"
	"strings"
	"sync"
	"time"
)

var path string
var targetDomain string
var timeStamp string

func getPage(url string) string {
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

func pathExists(Path string) bool {
	_, error := os.Stat(Path)
	if os.IsNotExist(error) {
		return false
	} else {
		return true
	}
}

func setDomain() {
	path, _ = os.Getwd()
	fmt.Print("Specify the target domain (only lowercase): ")
	fmt.Scanln(&targetDomain)
	fmt.Print(
		"Specify timestamp in the format:'yyyymmdd' (also: 'yyyy' > download only a specific year; 'yyyymm' > year and month; '2' > everything): ",
	)
	fmt.Scanln(&timeStamp)
}

func createDir() {
	fmt.Println("Saving data in:", path+"/"+targetDomain)
	err_ := os.Mkdir(path+"/"+targetDomain, 0777)
	if err_ != nil {
		log.Fatal(err_)
	}
}

func downloader(ctx context.Context, wg *sync.WaitGroup, c <-chan string, worker int) {
    defer wg.Done()
    for {
        select {
        case url, ok := <-c:
            if !ok {
                return
            }
            fmt.Printf("Worker %d downloading %s\n", worker, url)
            urlstring_ := strings.Replace(url, "/", "£", -1)
            urlstring__ := strings.Replace(urlstring_, ":", "!!!", -1)
            urlstring := strings.Replace(urlstring__, "?", "§§", -1)
            file_name_check := urlstring + ".txt"
            pathToFile := path + "/" + targetDomain + "/" + file_name_check
            if pathExists(pathToFile) == false {
                if len(url) < 255 {
                    content := getPage(url)
                    if content != "page not available" {
                        file, err := os.Create(pathToFile)
                        if err != nil {
                            fmt.Println(err)
                        }
                        file.WriteString(content)
                        file.Close()
                        fmt.Println("Done:", url)
                    }
                    frame := time.Duration(rand.Intn(100))
                    time.Sleep(time.Millisecond * frame)
                }
            } else {
                fmt.Println("Skipping:", url)
            }
        case <-ctx.Done():
            return
        }
    }
}

func getHistory() []string {
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
	var waybackurls []string
	for _, line := range lines {
		if len(line) > 0 {
			data := strings.Split(line, " ")
			savedpage := strings.Split(data[0], ")")[1]
			url := targetDomain + savedpage
			timestamp := string(data[1])
			if strings.HasPrefix(timestamp, timeStamp) {
				wayback_url := BASE_URL + timestamp + "/" + url
				waybackurls = append(waybackurls, wayback_url)
			}
		}
	}
	return waybackurls
}

func setDir(function func()) {
	var start_notice string
	if function != nil {
		function()
		start_notice = "Starting new download"
	} else {
		start_notice = "Resuming download"
	}
	fmt.Println(start_notice)
}

func main() {
    var wg sync.WaitGroup
	setDomain()
	if pathExists(path+"/"+targetDomain) == false {
		setDir(createDir)
	} else {
		setDir(nil)
	}
	history := getHistory()
	fmt.Println("Number of pages saved by Archive: ")
	fmt.Println(len(history))
	workers := 10
	channel := make(chan string)
        ctx, cancel := context.WithCancel(context.TODO())
        defer cancel()
        for i := 0; i < workers; i++ {
            wg.Add(1)
            go downloader(ctx, &wg, channel, i)
        }
        for _, url := range history {
            channel <- url
        }
        close(channel)
        wg.Wait()
        fmt.Println("Download completed. Press Enter to close window")
        fmt.Scanln()
}
