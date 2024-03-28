//go:build !solution

package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var (
	KeyToURL map[string]string
	URLToKey map[string]string
)

type URL struct {
	Value string `json:"url"`
}

type Response struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	key := r.URL.String()[4:]

	if KeyToURL[key] == "" {
		http.Error(w, "invalid key", http.StatusNotFound)

		return
	}

	url := KeyToURL[key]
	http.Redirect(w, r, url, http.StatusFound)
}

func Shorten(w http.ResponseWriter, r *http.Request) {
	var url URL

	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, "invalid JSON", 400)

		return
	}

	err := r.Body.Close()

	if err != nil {
		panic(err)
	}

	if URLToKey[url.Value] != "" {
		response := Response{URL: url.Value, Key: URLToKey[url.Value]}
		marshaled, _ := json.Marshal(response)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write(marshaled)

		return
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	key := strconv.Itoa(rand.Intn(1e9))

	KeyToURL[key] = url.Value
	URLToKey[url.Value] = key

	response := Response{URL: url.Value, Key: key}
	marshaled, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(marshaled)

	if err != nil {
		return
	}
}

func main() {
	port := flag.Int("port", 8000, "port string")
	flag.Parse()

	KeyToURL = make(map[string]string)
	URLToKey = make(map[string]string)

	http.HandleFunc("/shorten", Shorten)
	http.HandleFunc("/go/", Handle)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(*port), nil))
}
