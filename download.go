package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var path string
var targetDomain string
var timeStamp string
var chunksize int

func subslice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for {
		if len(slice) == 0 {
			break
		}
		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}
		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}
	return chunks
}

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

func saveFiles(array_ []string, waitgroup *sync.WaitGroup) {
	defer waitgroup.Done()
	for _, url := range array_ {
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
			savedpage := strings.Split(data[0], ")/")[1]
			url := targetDomain + "/"+ savedpage
			timestamp := string(data[1])
			if strings.HasPrefix(timestamp, timeStamp) {
				wayback_url := BASE_URL + timestamp + "/" + url
				waybackurls = append(waybackurls, wayback_url)
			}
		}
	}
	return waybackurls
}

func start(function func()) {
	var start_notice string
	if function != nil {
		function()
		start_notice = "Starting new download"
	} else {
		start_notice = "Resuming download"
	}
	history := getHistory()
	fmt.Println("Number of pages saved by Archive: ")
	fmt.Println(len(history))
	if len(history) >= 10 {
		chunksize = len(history) / 10
	} else {
		chunksize = len(history)
	}
	arrayOfArrays := subslice(history, chunksize)
	runtime.GOMAXPROCS(runtime.NumCPU())
	var waitgroup sync.WaitGroup
	waitgroup.Add(len(arrayOfArrays))
	fmt.Println(start_notice)
	for _, array := range arrayOfArrays {
		time.Sleep(time.Millisecond * 50)
		go saveFiles(array, &waitgroup)
	}
	waitgroup.Wait()
	fmt.Println("Download completed")
}

func main() {
	setDomain()
	if pathExists(path+"/"+targetDomain) == false {
		start(createDir)
	} else {
		start(nil)
	}
}
