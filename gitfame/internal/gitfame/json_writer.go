package gitfame

import (
	"encoding/json"
	"os"
)

type JSONWriter struct{}

func (JSONWriter JSONWriter) Write(statistics *[]Statistics) {
	data, err := json.Marshal(*statistics)

	if err != nil {
		panic(err)
	}

	if _, err = os.Stdout.Write(data); err != nil {
		return
	}
}
