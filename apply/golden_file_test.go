package apply_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/chanzuckerberg/fogg/apply"
	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/templates"
	"github.com/chanzuckerberg/fogg/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var updateGoldenFiles = flag.Bool("update", false, "when set, rewrite the golden files")

func TestIntegration(t *testing.T) {

	var testCases = []struct {
		fileName string
	}{
		{"v1_full"},
	}

	for _, tc := range testCases {
		t.Run(tc.fileName, func(t *testing.T) {
			a := assert.New(t)

			testdataFs := afero.NewBasePathFs(afero.NewOsFs(), filepath.Join(util.ProjectRoot(), "testdata", tc.fileName))

			if *updateGoldenFiles {
				// delete all files except fogg.json
				e := afero.Walk(testdataFs, ".", func(path string, info os.FileInfo, err error) error {
					if !info.IsDir() && path != "fogg.json" {
						return testdataFs.Remove(path)
					}
					return nil
				})
				a.NoError(e)

				conf, e := config.FindAndReadConfig(testdataFs, "fogg.json")
				a.NoError(e)
				fmt.Printf("conf %#v\n", conf)

				e = conf.Validate()
				a.NoError(e)

				e = apply.Apply(testdataFs, conf, templates.Templates, true)
				a.NoError(e)
			} else {

				fs, _, e := util.TestFs()
				a.NoError(e)

				// copy fogg.json into the tmp test dir (so that it doesn't show up as a diff)
				configContents, e := afero.ReadFile(testdataFs, "fogg.json")
				a.NoError(e)
				configMode, e := testdataFs.Stat("fogg.json")
				a.NoError(e)
				afero.WriteFile(fs, "fogg.json", configContents, configMode.Mode())

				conf, e := config.FindAndReadConfig(fs, "fogg.json")
				a.NoError(e)
				fmt.Printf("conf %#v\n", conf)

				e = conf.Validate()
				a.NoError(e)

				e = apply.Apply(fs, conf, templates.Templates, true)
				a.NoError(e)

				afero.Walk(testdataFs, ".", func(path string, info os.FileInfo, err error) error {
					log.Debug("================================================")
					log.Debug(path)
					if !info.Mode().IsRegular() {
						log.Debug("dir or link")
					} else {
						i1, e := testdataFs.Stat(path)
						a.NoError(e)
						i2, e := fs.Stat(path)
						a.NoError(e)
						a.Equal(i1.Size(), i2.Size())
						a.Equal(i1.Mode(), i2.Mode())

						f1, e := afero.ReadFile(testdataFs, path)
						a.NoError(e)
						f2, e := afero.ReadFile(fs, path)
						a.NoError(e)

						log.Debugf("f1:\n%s\n\n---- ", f1)
						log.Debugf("f2:\n%s\n\n---- ", f2)

						a.Equal(f1, f2)
					}
					return nil
				})
			}
		})
	}
}
