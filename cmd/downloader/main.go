package main

import (
	"flag"
	"log"

	"wayback-keyword-search/internal/web_domain"
)

func main() {
	var targetDomain string
	var timeStamp string
	var maxWorkers int

	flag.StringVar(&targetDomain, "domain", "", "Specify the target domain (only lowercase)")
	flag.StringVar(&timeStamp, "timeStamp", "", "Specify timestamp in the format:'yyyymmdd' (also: 'yyyy' > download only a specific year; 'yyyymm' > year and month; '2' or '1' > everything for the years past 20** or 19**")
	flag.IntVar(&maxWorkers, "workers", 10, "Specify the max workers (default=10)")

	flag.Parse()

	if targetDomain == "" || timeStamp == "" {
		log.Fatalf("Please provide both domain and timestamp.")
	}

	webDomain := web_domain.New(targetDomain, timeStamp)
	err := webDomain.Init()
	if err != nil {
		log.Fatalf("%s", err)
	}

	err = webDomain.Download(maxWorkers)
	if err != nil {
		log.Fatalf("%s", err)
	}

	log.Print("finished downloading")
}
