//go:build !solution

package gzep

import (
	"compress/gzip"
	"io"
	"sync"
)

var pool = sync.Pool{
	New: func() any {
		return gzip.NewWriter(nil)
	},
}

func Encode(data []byte, w io.Writer) error {
	writer, _ := pool.Get().(*gzip.Writer)
	defer pool.Put(writer)
	writer.Reset(w)

	defer func() {
		err := writer.Close()

		if err != nil {
			panic(err)
		}
	}()

	if _, err := writer.Write(data); err != nil {
		return err
	}

	return writer.Flush()
}
