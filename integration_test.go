package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/chanzuckerberg/fogg/apply"
	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/templates"
	"github.com/phayes/permbits"
	"github.com/spf13/afero"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.DebugLevel)
}
func TestIntegration(t *testing.T) {
	assert.Equal(t, 1, 1)
	files, _ := ioutil.ReadDir("integration-tests")
	for _, test := range files {
		t.Run(test.Name(), func(t *testing.T) {
			log.Debug(test.Name())
			path := filepath.Join("integration-tests", test.Name())
			outputFs := afero.NewMemMapFs()
			configPath := filepath.Join(path, "fogg.json")
			f, e := os.Open(configPath)
			assert.Nil(t, e)
			reader := io.ReadCloser(f)
			defer reader.Close()
			c, e := config.ReadConfig(reader)
			assert.Nil(t, e)
			apply.Apply(outputFs, c, templates.Templates, false)
			expectedFs := afero.NewBasePathFs(afero.NewOsFs(), filepath.Join("integration-tests", test.Name(), "expected"))
			assertEqualFses(t, expectedFs, outputFs)
		})
	}
}

func fileList(fileInfos []os.FileInfo) []string {
	var r []string
	for _, i := range fileInfos {
		r = append(r, i.Name())
	}
	return r
}

func assertEqualFses(t *testing.T, expectedFs afero.Fs, actualFs afero.Fs) {
	expectedDirList, e := afero.ReadDir(expectedFs, ".")
	assert.Nil(t, e)
	actualDirList, e := afero.ReadDir(actualFs, ".")
	assert.Nil(t, e)

	assert.Equal(t, fileList(expectedDirList), fileList(actualDirList))

	for i := range expectedDirList {
		expected := expectedDirList[i]
		actual := actualDirList[i]

		assert.Equal(t, expected.Name(), actual.Name())

		if expected.Name() != ".fogg-version" {
			assert.Equal(t, expected.IsDir(), actual.IsDir())
			if expected.IsDir() {
				// TODO make sure directory modes are the same
				// assert.Equalf(t, ex.Mode(), out.Mode(), "mode for %s, expected %s, actual %s",
				// 	expected.Name(), permbits.FileMode(ex.Mode()).String(), permbits.FileMode(out.Mode()).String())

				assertEqualFses(t, afero.NewBasePathFs(expectedFs, expected.Name()), afero.NewBasePathFs(actualFs, expected.Name()))
			} else {
				expectedBytes, e := afero.ReadFile(expectedFs, expected.Name())
				assert.Nil(t, e)
				actualBytes, e := afero.ReadFile(actualFs, expected.Name())
				assert.Nil(t, e)
				if string(actualBytes) != string(expectedBytes) {
					t.Errorf("Result not as expected:\n%v", diff.LineDiff(string(expectedBytes), string(actualBytes)))
				}

				assert.Equalf(t, expected.Mode(), actual.Mode(), "mode for %s, expected %s, actual %s",
					expectedFs.Name(), permbits.FileMode(expected.Mode()),
					permbits.FileMode(actual.Mode()))
			}
		}
	}
}
