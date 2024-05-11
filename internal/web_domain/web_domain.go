package web_domain

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"wayback-keyword-search/internal/engine"
	"wayback-keyword-search/internal/task"
	"wayback-keyword-search/internal/utils"
)

type Domain struct {
	targetDomain string
	timeStamp    string

	domainDir string
}

func New(targetDomain string, timeStamp string) *Domain {
	return &Domain{
		targetDomain: targetDomain,
		timeStamp:    timeStamp,
	}
}

func (d *Domain) Init() error {
	rootDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	d.domainDir = filepath.Join(rootDir, d.targetDomain)

	log.Printf("try to saving data in: %s", d.domainDir)

	if utils.PathExists(d.domainDir) == false {
		log.Print("starting new download")
		err := utils.CreateDir(d.domainDir)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	} else {
		log.Print("resuming download")
	}

	return nil
}

func (d *Domain) Download(maxWorkers int) error {
	log.Printf("retrieving information for %s, please wait.", d.targetDomain)
	history, err := engine.GetHistory(d.targetDomain, d.timeStamp)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	countTasks := len(history)

	log.Printf("number of pages saved by archive: %d", countTasks)

	if len(history) == 0 {
		return errors.New("empty list tasks")
	}

	task := task.New(maxWorkers)

	var wg sync.WaitGroup
	wg.Add(countTasks)

	for i := 0; i < countTasks; i++ {
		go task.Download(&wg, uint(i), d.domainDir, history[i])
	}

	wg.Wait()

	log.Printf("task completed: %d", countTasks)

	return nil
}
