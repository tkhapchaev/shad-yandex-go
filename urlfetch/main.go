//go:build !solution

package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	arguments := os.Args[1:]

	for _, url := range arguments {
		response, err := http.Get(url)

		if err != nil {
			os.Exit(1)
		}

		content, _ := io.ReadAll(response.Body)
		os.Stdout.Write(content)
	}
}
