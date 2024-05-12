package download

import (
	"context"
	"log"
	"os"
	"sync"

	"path/filepath"

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

func (t *Task) Run(ctx context.Context, wg *sync.WaitGroup, numWorker uint, rootDir string, inputURL string) {
	t.sem <- struct{}{}
	defer func() { <-t.sem }()

	defer wg.Done()

	select {
	case <-ctx.Done():
		log.Printf("worker %d was done", numWorker)
	default:
		log.Printf("worker: %d, downloading: %s", numWorker+1, inputURL)

		fileName := utils.UrlToFileName(inputURL) + ".txt"

		pathToFile := filepath.Join(rootDir, fileName)
		if utils.PathExists(pathToFile) != false {
			log.Printf("skipping: %s", inputURL)
			return
		}

		content, err := engine.GetPage(ctx, inputURL)
		if err != nil {
			log.Printf("error encountered: %s while retrieving content from URL: %s", err, inputURL)
			return
		}

		file, err := os.Create(pathToFile)
		if err != nil {
			log.Printf("error creating file: %s", err)
			return
		}

		file.WriteString(content)
		file.Close()

		log.Printf("done: %s", pathToFile)
	}
}
