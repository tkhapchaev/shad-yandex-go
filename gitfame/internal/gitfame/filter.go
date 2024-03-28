package gitfame

import (
	"path/filepath"
	"strings"

	"gitlab.com/manytask/itmo-go/public/gitfame/configs"
)

func Filter(files []string, configuration *Configuration) []string {
	if len(*configuration.Languages) != 0 {
		ApplyLanguages(*configuration.Languages, configuration.Extensions)
	}

	if len(*configuration.Extensions) != 0 {
		files = ApplyExtensions(files, *configuration.Extensions)
	}

	if len(*configuration.Exclusions) != 0 {
		files = ApplyExclusions(files, *configuration.Exclusions)
	}

	if len(*configuration.Restrictions) != 0 {
		files = ApplyRestrictions(files, *configuration.Restrictions)
	}

	return files
}

func ApplyLanguages(languages []string, extensions *[]string) {
	languageExtensions := *configs.GetLanguageExtensions()

	for _, language := range languages {
		for _, languageExtension := range languageExtensions {
			if strings.EqualFold(strings.ToLower(language), strings.ToLower(languageExtension.Name)) {
				*extensions = append(*extensions, languageExtension.Extensions...)
				break
			}
		}
	}
}

func ApplyExtensions(files, extensions []string) []string {
	var filesFiltered []string

	for _, file := range files {
		ok := false

		for _, extension := range extensions {
			if strings.EqualFold(extension, filepath.Ext(file)) {
				ok = true
			}
		}

		if ok {
			filesFiltered = append(filesFiltered, file)
		}
	}

	return filesFiltered
}

func ApplyExclusions(files, exclusions []string) []string {
	var filesFiltered []string

	for _, file := range files {
		ok := true

		for _, exclusion := range exclusions {
			matches, err := filepath.Match(exclusion, file)

			if err != nil {
				panic(err)
			}

			if matches {
				ok = false
			}
		}

		if ok {
			filesFiltered = append(filesFiltered, file)
		}
	}

	return filesFiltered
}

func ApplyRestrictions(files, restrictions []string) []string {
	var filesFiltered []string

	for _, file := range files {
		ok := false

		for _, restriction := range restrictions {
			matches, err := filepath.Match(restriction, file)

			if err != nil {
				panic(err)
			}

			if matches {
				ok = true
			}
		}

		if ok {
			filesFiltered = append(filesFiltered, file)
		}
	}

	return filesFiltered
}
