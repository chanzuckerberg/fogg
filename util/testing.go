package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

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

func ProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	return filepath.Join(basepath, "..")
}

func TestFile(name string) ([]byte, error) {
	// calculate the path to repository's root
	path := filepath.Join(ProjectRoot(), "testdata", name, "fogg.json")
	return ioutil.ReadFile(path)
}

//TestFileYml is created to support new yml file type
func TestFileYml(name string) ([]byte, error) {
	// calculate the path to repository's root
	path := filepath.Join(ProjectRoot(), "testdata", name, "fogg.yml")
	return ioutil.ReadFile(path)
}

func TestFs() (afero.Fs, string, error) {
	d, err := ioutil.TempDir("", "fogg")
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
