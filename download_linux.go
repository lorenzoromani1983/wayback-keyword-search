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

func saveFiles(array_ []string, waitgroup *sync.WaitGroup) {
	defer waitgroup.Done()
	for _, url := range array_ {
		if len(url) < 255 {
			content := getPage(url)
			if content != "page not available" {
				urlstring := strings.Replace(url, "/", "£", -1)

				dir, err := os.Getwd()
				if err != nil {
					log.Fatal(err)
				}

				pathToFile := dir + path + "/" + targetDomain + "/" + urlstring + ".txt"
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
	}
}

func createDir() {
	fmt.Print("Specify the target domain (only lowercase): ")
	fmt.Scanln(&targetDomain)
	fmt.Print("Specify timestamp in the format:'yyyymmdd' (also: 'yyyy' > download only a specific year; 'yyyymm' > year and date; '2' > everything): ")
	fmt.Scanln(&timeStamp)
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Saving data in:", path+"/"+targetDomain)
	err_ := os.Mkdir(path+"/"+targetDomain, 0777)
	if err_ != nil {
		log.Fatal(err_)
	}
}

func getHistory() []string {
	API_URL := "https://web.archive.org/cdx/search/cdx?url=*."
	query_url := API_URL + targetDomain
	client := http.Client{Timeout: time.Duration(50) * time.Second}
	response, err := client.Get(query_url)
	if err != nil {
		panic("Error: can't download history for the specified domain ")
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

func main() {
	createDir()
	history := getHistory() 
	fmt.Println("Number of pages to be downloaded: ")
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
	fmt.Println("Starting download...")
	for _, array := range arrayOfArrays {
		time.Sleep(time.Millisecond * 50)
		go saveFiles(array, &waitgroup)
	}
	waitgroup.Wait()
	fmt.Println("Download completed")
}
