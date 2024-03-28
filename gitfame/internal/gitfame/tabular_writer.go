package gitfame

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
)

type TabularWriter struct{}

func (tabularWriter TabularWriter) Write(statistics *[]Statistics) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	if _, err := fmt.Fprintln(writer, "Name\tLines\tCommits\tFiles"); err != nil {
		return
	}

	for _, stats := range *statistics {
		var line string

		line += stats.Name + "\t"
		line += strconv.Itoa(stats.Lines) + "\t"
		line += strconv.Itoa(stats.Commits) + "\t"
		line += strconv.Itoa(stats.Files)

		if _, err := fmt.Fprintln(writer, line); err != nil {
			panic(err)
		}
	}

	err := writer.Flush()

	if err != nil {
		panic(err)
	}
}
