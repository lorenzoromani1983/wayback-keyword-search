package download

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"wayback-keyword-search/internal/engine"
	"wayback-keyword-search/internal/utils"
)

type Task struct {
	sem chan struct{}
}

func New(maxWorkers int) *Task {
	return &Task{
		sem: make(chan struct{}, maxWorkers),
	}
}

func (t *Task) Run(ctx context.Context, wg *sync.WaitGroup, numWorker uint, rootDir string, url string) {
	t.sem <- struct{}{}
	defer func() { <-t.sem }()

	defer wg.Done()

	select {
	case <-ctx.Done():
		log.Printf("worker %d was done", numWorker)
	default:
		log.Printf("worker: %d, downloading: %s", numWorker+1, url)

		urlString_ := strings.Replace(url, "/", "_", -1)
		urlString__ := strings.Replace(urlString_, ":", "", -1)
		urlstring := strings.Replace(urlString__, "?", "§§", -1)
		fileNameCheck := urlstring + ".txt"

		pathToFile := filepath.Join(rootDir, fileNameCheck)
		if utils.PathExists(pathToFile) != false {
			log.Printf("skipping: %s", url)
			return
		}

		content, err := engine.GetPage(ctx, url)
		if err != nil {
			log.Printf("got err: %s", err)
		} else {
			file, err := os.Create(pathToFile)
			if err != nil {
				log.Printf("url: %s, got err: %s", url, err)
			} else {
				file.WriteString(content)
				file.Close()
				log.Printf("done: %s", url)
			}
		}
	}
}
