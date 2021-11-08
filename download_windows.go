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

func scrape(array_ []string, waitgroup *sync.WaitGroup) {
	defer waitgroup.Done()
	for _, url := range array_ {
		if len(url) < 255 {
			client := http.Client{Timeout: time.Duration(50) * time.Second}
			response, err := client.Get(url)
			if err != nil {
				fmt.Println("[!] Errore:", url)

				continue
			}
			body, _ := io.ReadAll(response.Body)
			content := string(body)

			if response.StatusCode == 200 {
				urlstring_ := strings.Replace(url, "/", "£", -1)
				urlstring := strings.Replace(urlstring_,":","$",-1)

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
			} else {
				fmt.Println("Pagina non disponibile:", url)
			}

			response.Body.Close()
			frame := time.Duration(rand.Intn(100))
			time.Sleep(time.Millisecond * frame)
		}
	}
}

func createDir() {
	fmt.Print("Specificare il dominio (caratteri minuscoli): ")
	fmt.Scanln(&targetDomain)
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Salvataggio dei dati in:", path+"/"+targetDomain)
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
		panic("Errore: impossibile scaricare lo storico per il domino specificato, controlla l'esattezza dell'URL e riprova o accertati di essere connesso a internet")
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
			wayback_url := BASE_URL + timestamp + "/" + url
			waybackurls = append(waybackurls, wayback_url)
		}
	}
	return waybackurls
}

func main() {
	createDir()
	history := getHistory()
	fmt.Println("Numero di pagine archiviate dalla WayBack Machine:")
	fmt.Println(len(history))
	chunksize := len(history) / 10
	arrayOfArrays := subslice(history, chunksize)
	runtime.GOMAXPROCS(runtime.NumCPU())
	var waitgroup sync.WaitGroup
	waitgroup.Add(len(arrayOfArrays))
	fmt.Println("Avvio il download...")
	for _, array := range arrayOfArrays {
		time.Sleep(time.Millisecond * 100)
		go scrape(array, &waitgroup)
	}
	waitgroup.Wait()
	fmt.Println("Download completato")
}
