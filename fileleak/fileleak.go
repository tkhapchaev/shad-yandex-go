//go:build !solution

package fileleak

import (
	"io/fs"
	"io/ioutil"
	"reflect"
)

type testingT interface {
	Errorf(msg string, args ...interface{})
	Cleanup(func())
}

func getFileInfos() []fs.FileInfo {
	files, err := ioutil.ReadDir("/proc/self/fd")

	if err != nil {
		panic(err)
	}

	return files
}

func VerifyNone(t testingT) {
	before := getFileInfos()

	t.Cleanup(func() {
		after := getFileInfos()

		if !reflect.DeepEqual(after, before) {
			t.Errorf("file leaks detected")
		}
	})
}
