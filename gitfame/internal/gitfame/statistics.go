package gitfame

import (
	"strconv"
	"strings"
)

type Statistics struct {
	Name    string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}

func GetStatistics(file string, blame []string, configuration *Configuration) (map[string][]string, map[string]int) {
	authors := make(map[string][]string)
	commits := make(map[string]int)
	nextHash := true

	var (
		iterations    int
		waitForAuthor bool
		lastHash      string
	)

	if len(blame) == 0 {
		hash, name := GitLog(file, configuration)
		commits[hash] = 0
		authors[name] = append(authors[name], hash)
	}

	for _, line := range blame {
		if nextHash {
			nextHash = false

			if iterations == 0 {
				split := strings.Split(line, " ")
				iterations, _ = strconv.Atoi(split[len(split)-1])
				commits[split[0]] += iterations
				waitForAuthor = true
				lastHash = split[0]
			}

			iterations--
		} else if line[0] != '\t' && waitForAuthor {
			split := strings.Split(line, " ")

			if !*configuration.UseCommitter && split[0] == "author" {
				name := line[len("author "):]
				authors[name] = append(authors[name], lastHash)
				waitForAuthor = false
			} else if *configuration.UseCommitter && split[0] == "committer" {
				name := line[len("committer "):]
				authors[name] = append(authors[name], lastHash)
				waitForAuthor = false
			}
		} else if line[0] == '\t' {
			nextHash = true
		}
	}

	return authors, commits
}

func TransformStatistics(commits, files map[string]int, authors map[string]map[string]struct{}, statistics []*Statistics) []*Statistics {
	for author, commitHashes := range authors {
		stats := &Statistics{
			Name:    author,
			Lines:   0,
			Commits: len(commitHashes),
			Files:   files[author],
		}

		statistics = append(statistics, stats)
	}

	for author, commitHashes := range authors {
		var lines int

		for hash := range commitHashes {
			lines += commits[hash]
		}

		for _, stats := range statistics {
			if strings.EqualFold(stats.Name, author) {
				stats.Lines += lines
				break
			}
		}
	}

	return statistics
}

func UpdateStatistics(commitHashes map[string]int, commitsUpdated, files map[string]int, commits map[string][]string, authors map[string]map[string]struct{}) {
	for hash, count := range commitHashes {
		commitsUpdated[hash] += count
	}

	for author, hashes := range commits {
		files[author]++

		for _, hash := range hashes {
			if _, ok := authors[author]; !ok {
				authors[author] = make(map[string]struct{})
			}

			authors[author][hash] = struct{}{}
		}
	}
}
