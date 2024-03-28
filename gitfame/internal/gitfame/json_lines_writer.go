package gitfame

import (
	"encoding/json"
	"fmt"
	"os"
)

type JSONLinesWriter struct{}

func (JSONLinesWriter JSONLinesWriter) Write(statistics *[]Statistics) {
	for _, stats := range *statistics {
		data, err := json.Marshal(stats)

		if err != nil {
			panic(err)
		}

		if _, err = os.Stdout.Write(data); err != nil {
			panic(err)
		}

		fmt.Println()
	}
}
