package util

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
)

func Intptr(i int64) *int64 {
	return &i
}

func StrPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func TestFile(name string) ([]byte, error) {
	// calculate the path to repository's root
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	path := filepath.Join(basepath, "..", "testdata", name, "fogg.json")
	return ioutil.ReadFile(path)
}
