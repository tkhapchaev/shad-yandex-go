package gitfame

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

type Configuration struct {
	Repository   *string
	Revision     *string
	OrderBy      *string
	UseCommitter *bool
	Format       *string
	Extensions   *[]string
	Languages    *[]string
	Exclusions   *[]string
	Restrictions *[]string
	Writer       Writer
}

func GetConfiguration() *Configuration {
	var (
		flagRepository   = pflag.String("repository", ".", "Path to the Git repository")
		flagRevision     = pflag.String("revision", "HEAD", "Commit pointer (SHA-1 hash)")
		flagOrderBy      = pflag.String("order-by", "lines", "Order by key: lines, commits or files")
		flagUseCommitter = pflag.Bool("use-committer", false, "Use committer instead of author")
		flagFormat       = pflag.String("format", "tabular", "Output format: tabular, .csv, .json, json-lines")
		flagExtensions   = pflag.StringSlice("extensions", []string{}, "List of file extensions to include, comma-separated")
		flagLanguages    = pflag.StringSlice("languages", []string{}, "List of languages to include, comma-separated")
		flagExclude      = pflag.StringSlice("exclude", []string{}, "Set of glob patterns to exclude files")
		flagRestrictTo   = pflag.StringSlice("restrict-to", []string{}, "Set of glob patterns to restrict files to")
	)

	pflag.Parse()

	if _, err := os.Stat(*flagRepository); os.IsNotExist(err) {
		panic("invalid repository path: " + *flagRepository + "\n")
	}

	if *flagOrderBy != "lines" && *flagOrderBy != "commits" && *flagOrderBy != "files" {
		panic("invalid sorting key: " + *flagOrderBy + "\n")
	}

	if *flagFormat != "tabular" && *flagFormat != "csv" && *flagFormat != "json" && *flagFormat != "json-lines" {
		panic("invalid format: " + *flagFormat + "\n")
	}

	orderBy = *flagOrderBy
	var writer Writer

	if *flagFormat == "tabular" {
		writer = TabularWriter{}
	} else if *flagFormat == "csv" {
		writer = CSVWriter{}
	} else if *flagFormat == "json" {
		writer = JSONWriter{}
	} else if *flagFormat == "json-lines" {
		writer = JSONLinesWriter{}
	}

	for _, extension := range *flagExtensions {
		if !strings.HasPrefix(extension, ".") {
			_, _ = fmt.Fprintf(os.Stderr, "Invalid extension: %s", extension)
			os.Exit(1)
		}
	}

	return &Configuration{
		Repository:   flagRepository,
		Revision:     flagRevision,
		OrderBy:      flagOrderBy,
		UseCommitter: flagUseCommitter,
		Format:       flagFormat,
		Extensions:   flagExtensions,
		Languages:    flagLanguages,
		Exclusions:   flagExclude,
		Restrictions: flagRestrictTo,
		Writer:       writer,
	}
}
