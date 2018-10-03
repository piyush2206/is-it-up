// A simple GoLang program to check if provided wibsites URLs are Up.

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	flgs := readFlags()

	websiteURLs := readURLs(*flgs.uRLsPath)

	chWebsiteStatus := checkWebsites(flgs, websiteURLs)

	for resStatus := range chWebsiteStatus {
		fmt.Println(resStatus)
	}
}

type flags struct {
	uRLsPath    *string
	concurrency *int
}

func readFlags() (flgs *flags) {
	flgs = new(flags)

	flgs.uRLsPath = flag.String("u", "urls.csv", "website URLs file path")
	flgs.concurrency = flag.Int("c", 50, "concurrent requests")

	flag.Parse()
	return
}

// Enum values of tyoe Status
const (
	NotOk Status = iota // 0
	Ok                  // 1
)

type (
	websiteStatus struct {
		// Url of the website
		URL string

		// Status indicates if the site is Up or otherwise
		Status Status
	}

	// Status will have values "OK" and "NOT OK" if website is Up or down respectively
	Status int
)

func (s Status) String() (strStatus string) {
	switch s {
	case Ok:
		strStatus = "OK"
	case NotOk:
	default:
		strStatus = "NOT OK"
	}

	return
}

func (w *websiteStatus) String() string {
	return fmt.Sprintf("URL: %s | Status: %s", w.URL, w.Status)
}

func readURLs(urlsPath string) (websiteURLs []string) {
	reader, err := os.Open(urlsPath)
	csvReader := csv.NewReader(reader)

	URLs, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, URL := range URLs {
		websiteURLs = append(websiteURLs, URL...)
	}

	return
}

func checkWebsites(flgs *flags, websiteURLs []string) (chWebsiteStatus chan *websiteStatus) {
	chWebsiteURLs := make(chan string)
	chWebsiteStatus = make(chan *websiteStatus)

	var wg sync.WaitGroup
	wg.Add(len(websiteURLs))

	go func() {
		for _, url := range websiteURLs {
			chWebsiteURLs <- url
		}
		close(chWebsiteURLs)
	}()

	for i := 0; i < *flgs.concurrency; i++ {
		go func() {
			for url := range chWebsiteURLs {
				status := isUp(url)

				chWebsiteStatus <- &websiteStatus{
					URL:    url,
					Status: status,
				}
				wg.Done()
			}
		}()
	}

	go func() {
		wg.Wait()
		close(chWebsiteStatus)
	}()

	return
}

func isUp(url string) (status Status) {
	response, err := http.Get(url)
	if err != nil {
		status = NotOk
		return
	} else {
		defer response.Body.Close()
		_, err := ioutil.ReadAll(response.Body)
		if err != nil {
			status = NotOk
			return
		}
		status = Ok
	}
	return
}
