package gitfame

import (
	"encoding/csv"
	"os"
	"strconv"
)

type CSVWriter struct{}

func (CSVWriter CSVWriter) Write(statistics *[]Statistics) {
	writer := csv.NewWriter(os.Stdout)
	err := writer.Write([]string{"Name", "Lines", "Commits", "Files"})

	if err != nil {
		panic(err)
	}

	for _, stats := range *statistics {
		line := []string{
			stats.Name,
			strconv.Itoa(stats.Lines),
			strconv.Itoa(stats.Commits),
			strconv.Itoa(stats.Files)}

		err = writer.Write(line)

		if err != nil {
			panic(err)
		}
	}

	writer.Flush()
}
