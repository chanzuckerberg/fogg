package util

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/spf13/afero"
)

func Intptr(i int64) *int64 {
	return &i
}

func JSONNumberPtr(i int) *json.Number {
	j := json.Number(strconv.Itoa(i))
	return &j
}

func StrPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func ProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	return filepath.Join(basepath, "..")
}

func TestFile(name string) ([]byte, error) {
	// always look for yaml, json deprecated
	path := filepath.Join(ProjectRoot(), "testdata", name, "fogg.yml")
	return os.ReadFile(path)
}

func TestFs() (afero.Fs, string, error) {
	d, err := os.MkdirTemp("", "fogg")
	if err != nil {
		return nil, "", err
	}

	fs := afero.NewBasePathFs(afero.NewOsFs(), d)

	return fs, d, nil
}

func PwdFs() (afero.Fs, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return afero.NewBasePathFs(afero.NewOsFs(), pwd), nil
}
