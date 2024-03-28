//go:build !solution

package main

import (
	"fmt"
	"os"
	"strconv"
	s "strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	countingMap := make(map[string]int)
	arguments := os.Args[1:]
	var strings []string

	for i := 0; i < len(arguments); i++ {
		data, err := os.ReadFile(arguments[i])
		check(err)

		strings = append(strings, s.Split(string(data), "\n")...)
	}

	for _, str := range strings {
		_, ok := countingMap[str]

		if ok {
			countingMap[str] += 1
		} else {
			countingMap[str] = 1
		}
	}

	for key, value := range countingMap {
		if value >= 2 {
			fmt.Println(strconv.Itoa(value) + "\t" + key)
		}
	}
}
