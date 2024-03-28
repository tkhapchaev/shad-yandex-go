package gitfame

import (
	"sort"
	"strings"
)

var (
	orderBy string
)

func Sort(statistics []Statistics) *[]Statistics {
	sort.Slice(statistics, func(i, j int) bool {
		var (
			firstStats, secondStats []int
		)

		switch orderBy {
		case "lines":
			firstStats = []int{statistics[i].Lines, statistics[i].Commits, statistics[i].Files}
			secondStats = []int{statistics[j].Lines, statistics[j].Commits, statistics[j].Files}
		case "commits":
			firstStats = []int{statistics[i].Commits, statistics[i].Lines, statistics[i].Files}
			secondStats = []int{statistics[j].Commits, statistics[j].Lines, statistics[j].Files}
		case "files":
			firstStats = []int{statistics[i].Files, statistics[i].Lines, statistics[i].Commits}
			secondStats = []int{statistics[j].Files, statistics[j].Lines, statistics[j].Commits}
		default:
			panic("invalid sorting key: " + orderBy + "\n")
		}

		if firstStats[0] != secondStats[0] {
			return firstStats[0] > secondStats[0]
		}

		if firstStats[1] != secondStats[1] {
			return firstStats[1] > secondStats[1]
		}

		if firstStats[2] != secondStats[2] {
			return firstStats[2] > secondStats[2]
		}

		return strings.ToLower(statistics[i].Name) < strings.ToLower(statistics[j].Name)
	})

	return &statistics
}
