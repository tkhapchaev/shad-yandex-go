package configs

import (
	_ "embed"
	"encoding/json"
)

type LanguageExtensions struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Extensions []string `json:"extensions"`
}

//go:embed language_extensions.json
var jsonData []byte

func GetLanguageExtensions() *[]LanguageExtensions {
	var languageExtensions []LanguageExtensions
	err := json.Unmarshal(jsonData, &languageExtensions)

	if err != nil {
		panic(err)
	}

	return &languageExtensions
}
