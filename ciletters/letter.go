//go:build !solution

package ciletters

import (
	"strings"
	"text/template"
)

func MakeLetter(n *Notification) (string, error) {
	var stringBuilder strings.Builder
	templateStr := `Your pipeline #{{ .Pipeline.ID}} {{if ne .Pipeline.Status "ok"}}has failed{{else}}passed{{end}}!
    Project:      {{ .Project.GroupID }}/{{.Project.ID }}
    Branch:       🌿 {{ .Branch }}
    Commit:       {{ slice .Commit.Hash 0 8 }} {{ .Commit.Message }}
    CommitAuthor: {{ .Commit.Author }}{{range $job := .Pipeline.FailedJobs}}
        Stage: {{$job.Stage}}, Job {{$job.Name}}{{range cmdLog $job.RunnerLog}}
            {{.}}{{end}}
{{end}}`

	notificationTemplate, err := template.New("email").Funcs(template.FuncMap{"cmdLog": func(str string) []string {
		result := strings.Split(str, "\n")

		if len(result) > 9 {
			result = result[9:]
		}

		return result
	},
	}).Parse(templateStr)

	if err != nil {
		return "", err
	}

	err = notificationTemplate.Execute(&stringBuilder, n)

	if err != nil {
		return "", err
	} else {
		return stringBuilder.String(), nil
	}
}
