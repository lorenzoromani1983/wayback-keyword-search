package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"wayback-keyword-search/internal/engine"
)

var sem chan struct{}

var path string
var targetDomain string
var timeStamp string
var maxWorkers int

func main() {
	flag.StringVar(&targetDomain, "domain", "", "Specify the target domain (only lowercase)")
	flag.StringVar(&timeStamp, "timeStamp", "", "Specify timestamp in the format:'yyyymmdd' (also: 'yyyy' > download only a specific year; 'yyyymm' > year and month; '2' or '1' > everything for the years past 20** or 19**")
	flag.IntVar(&maxWorkers, "workers", 10, "Specify the max workers (default=10)")

	flag.Parse()

	if targetDomain == "" || timeStamp == "" {
		fmt.Println("Please provide both domain and timestamp.")
		return
	}

	path, _ = os.Getwd()

	if pathExists(path+"/"+targetDomain) == false {
		setDir(createDir)
	} else {
		setDir(nil)
	}

	history := engine.GetHistory(targetDomain, timeStamp)
	historyLen := len(history)

	fmt.Printf("Number of pages saved by Archive: %d\n", historyLen)

	if len(history) == 0 {
		return
	}

	sem = make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup
	wg.Add(historyLen)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	for i := 0; i < historyLen; i++ {
		go downloader(ctx, &wg, history[i])
	}

	wg.Wait()

	fmt.Println("Download completed.")
}

func pathExists(Path string) bool {
	_, error := os.Stat(Path)
	if os.IsNotExist(error) {
		return false
	} else {
		return true
	}
}

func createDir() {
	fmt.Println("Saving data in:", path+"/"+targetDomain)
	err := os.Mkdir(path+"/"+targetDomain, 0777)
	if err != nil {
		log.Fatal(err)
	}
}

func downloader(ctx context.Context, wg *sync.WaitGroup, url string) {
	sem <- struct{}{}
	defer func() { <-sem }()

	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Printf("Worker downloading %s\n", url)

			urlString_ := strings.Replace(url, "/", "£", -1)
			urlString__ := strings.Replace(urlString_, ":", "!!!", -1)
			urlstring := strings.Replace(urlString__, "?", "§§", -1)
			fileNameCheck := urlstring + ".txt"
			pathToFile := path + "/" + targetDomain + "/" + fileNameCheck

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
