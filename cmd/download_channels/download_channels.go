package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"wayback-keyword-search/internal/engine"
)

var path string
var targetDomain string
var timeStamp string

func main() {
	var wg sync.WaitGroup
	setDomain()
	if pathExists(path+"/"+targetDomain) == false {
		setDir(createDir)
	} else {
		setDir(nil)
	}
	history := engine.GetHistory(targetDomain, timeStamp)
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
		"Specify timestamp in the format:'yyyymmdd' (also: 'yyyy' > download only a specific year; 'yyyymm' > year and month; '2' or '1' > everything for the years past 20** or 19**): ",
	)
	fmt.Scanln(&timeStamp)
}

func createDir() {
	fmt.Println("Saving data in:", path+"/"+targetDomain)
	err := os.Mkdir(path+"/"+targetDomain, 0777)
	if err != nil {
		log.Fatal(err)
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
					content := engine.GetPage(url)
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