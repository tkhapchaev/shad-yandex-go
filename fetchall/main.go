//go:build !solution

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func fetchURL(url string, channel chan<- string) {
	start := time.Now()
	response, err := http.Get(url)

	if err != nil {
		channel <- fmt.Sprint("(!) Failed to fetch URL " + url)

		return
	}

	numberOfBytes, _ := io.Copy(ioutil.Discard, response.Body)
	channel <- fmt.Sprintf("%.3fs %10d\t%s", time.Since(start).Seconds(), numberOfBytes, url)
}

func main() {
	start := time.Now()
	arguments := os.Args[1:]
	channel := make(chan string)

	for _, url := range arguments {
		go fetchURL(url, channel)
	}

	for range arguments {
		fmt.Println(<-channel)
	}

	fmt.Printf("%.3fs elapsed\n", time.Since(start).Seconds())
}
