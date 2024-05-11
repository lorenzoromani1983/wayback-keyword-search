package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"wayback-keyword-search/internal/engine"
	"wayback-keyword-search/internal/utils"
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

	path, err := os.Getwd()
	if err != nil {
		log.Fatalf("%s", err)
	}

	pathDomain := filepath.Join(path, targetDomain)

	fmt.Println("Try to saving data in:", pathDomain)

	if utils.PathExists(pathDomain) == false {
		fmt.Println("Starting new download")
		err := utils.CreateDir(pathDomain)
		if err != nil {
			log.Fatal("%w", err)
		}
	} else {
		fmt.Println("Resuming download")
	}

	history, err := engine.GetHistory(targetDomain, timeStamp)
	if err != nil {
		log.Fatalf("%s", err)
	}

	historyLen := len(history)

	fmt.Printf("Number of pages saved by Archive: %d\n", historyLen)

	if len(history) == 0 {
		return
	}

	sem = make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup
	wg.Add(historyLen)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; i < historyLen; i++ {
		go downloader(ctx, &wg, pathDomain, history[i])
	}

	wg.Wait()

	fmt.Println("Download completed.")
}

func downloader(ctx context.Context, wg *sync.WaitGroup, basePathDir string, url string) {
	sem <- struct{}{}
	defer func() { <-sem }()

	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Printf("Worker downloading %s\n", url)

			urlString_ := strings.Replace(url, "/", "_", -1)
			urlString__ := strings.Replace(urlString_, ":", "", -1)
			urlstring := strings.Replace(urlString__, "?", "§§", -1)
			fileNameCheck := urlstring + ".txt"

			pathToFile := filepath.Join(basePathDir, fileNameCheck)

			if utils.PathExists(pathToFile) != false {
				fmt.Println("skipping:", url)

				return
			}

			content, err := engine.GetPage(url)
			if err != nil {
				log.Printf("got err: %s", err)
			} else {

				file, err := os.Create(pathToFile)
				if err != nil {
					log.Printf("url: %s, got err: %s", url, err)
				} else {
					file.WriteString(content)
					file.Close()

					fmt.Println("Done:", url)
				}
			}

			return
		}
	}
}

