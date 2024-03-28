//go:build !solution

package varfmt

import (
	"fmt"
	"strconv"
	"strings"
)

func Sprintf(format string, args ...interface{}) string {
	var position int
	result := strings.Builder{}
	argsMap := make(map[interface{}]string)

	for i := 0; i < len(format); i++ {
		if format[i] != '{' {
			result.WriteString(string(format[i]))
		} else {
			j := i + 1

			for format[j] != '}' {
				j += 1
			}

			index := 0

			if j == i+1 {
				index = position
			} else {
				index, _ = strconv.Atoi(format[i+1 : j])
			}

			if argsMap[args[index]] == "" {
				arg := fmt.Sprint(args[index])
				result.WriteString(arg)
				argsMap[args[index]] = arg
			} else {
				result.WriteString(argsMap[args[index]])
			}

			position += 1
			i = j
		}
	}

	return result.String()
}
