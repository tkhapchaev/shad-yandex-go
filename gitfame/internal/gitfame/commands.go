package gitfame

import (
	"os/exec"
	"strings"
)

func GitLsTree(configuration *Configuration) []string {
	command := exec.Command("git", "ls-tree", "-r", "--name-only", *configuration.Revision)
	command.Dir = *configuration.Repository

	output, err := command.Output()

	if err != nil {
		panic(err)
	}

	files := strings.Split(string(output), "\n")

	if len(files) > 0 {
		files = files[:len(files)-1]
	}

	return Filter(files, configuration)
}

func GitLog(file string, configuration *Configuration) (string, string) {
	command := exec.Command("git", "log", "--pretty=format:%H %an", *configuration.Revision, "--", file)
	command.Dir = *configuration.Repository

	output, err := command.Output()

	if err != nil {
		panic(err)
	}

	var builder strings.Builder
	builder.Write(output)

	logs := strings.Split(builder.String(), "\n")
	hash := strings.Split(logs[0], " ")
	name := strings.Join(hash[1:], " ")

	return hash[0], name
}

func GitBlame(file string, configuration *Configuration) []string {
	command := exec.Command("git", "blame", "--porcelain", *configuration.Revision, file)
	command.Dir = *configuration.Repository

	output, err := command.Output()

	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(output), "\n")

	if len(lines) > 0 {
		lines = lines[:len(lines)-1]
	}

	return lines
}
