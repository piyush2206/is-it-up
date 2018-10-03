// 1. Program to check status of 5 websites. Print "OK" if site is up and "NOT OK" if site is down.
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

// Enum values of tyoe Status
const (
	Ok    Status = "OK"
	NotOk Status = "NOT OK"
)

type (
	websiteStatus struct {
		// Url of the website
		URL string

		// Status indicates if the site is Up or otherwise
		Status Status
	}

	// Status will have values "OK" and "NOT OK" if website is Up or down respectively
	Status string
)

func main() {
	websiteUrls := []string{"https://www.github.com", "https://www.google.com"}
	chWebsiteStatus := checkWebsites(websiteUrls)
	for resStatus := range chWebsiteStatus {
		fmt.Println(resStatus)
	}
}

func (w *websiteStatus) String() string {
	return fmt.Sprintf("URL: %s | Status: %v", w.URL, w.Status)
}

func checkWebsites(websiteUrls []string) (chWebsiteStatus chan *websiteStatus) {
	chWebsiteStatus = make(chan *websiteStatus)

	var wg sync.WaitGroup
	wg.Add(len(websiteUrls))

	for _, url := range websiteUrls {
		go func(url string) {
			status := isUp(url)

			chWebsiteStatus <- &websiteStatus{
				URL:    url,
				Status: status,
			}
			wg.Done()
		}(url)
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
		status = Status(NotOk)
		return
	} else {
		defer response.Body.Close()
		_, err := ioutil.ReadAll(response.Body)
		if err != nil {
			status = NotOk
			return
		}
		status = Status(Ok)
	}
	return
}
