package gitfame

import (
	"gitlab.com/manytask/itmo-go/public/gitfame/internal/gitfame"
)

func Run() {
	var (
		statistics []*gitfame.Statistics
		toSort     []gitfame.Statistics
	)

	configuration := gitfame.GetConfiguration()
	gitFiles := gitfame.GitLsTree(configuration)

	authors := make(map[string]map[string]struct{})
	commits := make(map[string]int)
	files := make(map[string]int)

	for _, file := range gitFiles {
		blame := gitfame.GitBlame(file, configuration)
		authorsUpdated, commitsUpdated := gitfame.GetStatistics(file, blame, configuration)
		gitfame.UpdateStatistics(commitsUpdated, commits, files, authorsUpdated, authors)
	}

	statistics = gitfame.TransformStatistics(commits, files, authors, statistics)

	for _, stats := range statistics {
		toSort = append(toSort, *stats)
	}

	configuration.Writer.Write(gitfame.Sort(toSort))
}
