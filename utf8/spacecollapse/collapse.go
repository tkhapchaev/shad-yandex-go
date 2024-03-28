//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func CollapseSpaces(input string) string {
	var result strings.Builder
	result.Grow(len(input))
	var previousRune rune

	for len(input) > 0 {
		currentRune, currentLength := utf8.DecodeRuneInString(input)

		if currentLength == 1 && currentRune == utf8.RuneError {
			result.WriteRune('\uFFFD')
		} else if !unicode.IsSpace(currentRune) {
			result.WriteRune(currentRune)
			previousRune = currentRune
		} else if previousRune != ' ' {
			result.WriteRune(' ')
			previousRune = ' '
		}

		input = input[currentLength:]
	}

	return result.String()
}
