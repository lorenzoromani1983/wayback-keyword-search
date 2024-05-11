package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
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

	pathToDomain := filepath.Join(path, targetDomain)

	if pathExists(pathToDomain) == false {
		fmt.Println("Starting new download")
		createDir(pathToDomain)
	} else {
		fmt.Println("Resuming download")
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

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	for i := 0; i < historyLen; i++ {
		go downloader(ctx, &wg, history[i])
	}

	wg.Wait()

	fmt.Println("Download completed.")
}

func pathExists(path string) bool {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return false
    }

	return true
}

func createDir(pathDir string) {
	fmt.Println("Saving data in:", pathDir)
	err := os.Mkdir(pathDir, 0777)
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

			pathToFile := filepath.Join(path, targetDomain, fileNameCheck)

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
