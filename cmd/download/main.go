package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"wayback-keyword-search/internal/engine"
)

var path string
var targetDomain string
var timeStamp string
var chunksize int

func main() {
	setDomain()
	if pathExists(path+"/"+targetDomain) == false {
		start(createDir)
	} else {
		start(nil)
	}
}

func setDomain() {
	path, _ = os.Getwd()
	fmt.Print("Specify the target domain (only lowercase): ")
	fmt.Scanln(&targetDomain)
	fmt.Print(
		"Specify timestamp in the format:'yyyymmdd' (also: 'yyyy' > download only a specific year; 'yyyymm' > year and month; '2' or '1' > everything for the years past 20** or 19**): ",
	)
	fmt.Scanln(&timeStamp)
}

func start(function func()) {
	var start_notice string
	if function != nil {
		function()
		start_notice = "Starting new download"
	} else {
		start_notice = "Resuming download"
	}
	history := engine.GetHistory(targetDomain, timeStamp)
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

func createDir() {
	fmt.Println("Saving data in:", path+"/"+targetDomain)
	err := os.Mkdir(path+"/"+targetDomain, 0777)
	if err != nil {
		log.Fatal(err)
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
				content, err := engine.GetPage(url)
				if err != nil {
					log.Printf("got err: %s", err)
				} else {
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

func pathExists(Path string) bool {
	_, error := os.Stat(Path)
	if os.IsNotExist(error) {
		return false
	} else {
		return true
	}
}
