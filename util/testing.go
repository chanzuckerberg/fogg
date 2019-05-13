package util

import (
	"io/ioutil"

	"github.com/spf13/afero"
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

func TestFs() (afero.Fs, string, error) {
	d, err := ioutil.TempDir("", "fogg")
	if err != nil {
		return nil, "", err
	}

	fs := afero.NewBasePathFs(afero.NewOsFs(), d)

	return fs, d, nil
}
