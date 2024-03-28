//go:build !solution

package reverse

import (
	"strings"
	"unicode/utf8"
)

func Reverse(input string) string {
	var result strings.Builder
	result.Grow(len(input))

	for len(input) > 0 {
		currentRune, currentLength := utf8.DecodeLastRuneInString(input)

		if currentLength == 1 && currentRune == utf8.RuneError {
			result.WriteRune('\uFFFD')
		} else {
			result.WriteRune(currentRune)
		}

		input = input[:len(input)-currentLength]
	}

	return result.String()
}
