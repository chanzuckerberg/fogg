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
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

var updateGoldenFiles = flag.Bool("update", false, "when set, rewrite the golden files")

func TestIntegration(t *testing.T) {

	var testCases = []struct {
		fileName   string
		configFile string
	}{
		{"okta_provider_yaml", "fogg.yml"},
		{"github_provider_yaml", "fogg.yml"},
		{"bless_provider_yaml", "fogg.yml"},
		{"snowflake_provider_yaml", "fogg.yml"},
		{"v2_full_yaml", "fogg.yml"},
		{"v2_minimal_valid_yaml", "fogg.yml"},
		{"v2_no_aws_provider_yaml", "fogg.yml"},
		{"github_actions", "fogg.yml"},
		{"circleci", "fogg.yml"},
		{"remote_backend_yaml", "fogg.yml"},
	}

	for _, tc := range testCases {
		t.Run(tc.fileName, func(t *testing.T) {
			r := require.New(t)

			testdataFs := afero.NewBasePathFs(afero.NewOsFs(), filepath.Join(util.ProjectRoot(), "testdata", tc.fileName))

			if *updateGoldenFiles {
				// delete all files except fogg.json
				e := afero.Walk(testdataFs, ".", func(path string, info os.FileInfo, err error) error {
					if !info.IsDir() && !(path == tc.configFile) {
						return testdataFs.Remove(path)
					}
					return nil
				})
				r.NoError(e)

				conf, e := config.FindAndReadConfig(testdataFs, tc.configFile)
				r.NoError(e)
				fmt.Printf("conf %#v\n", conf)
				fmt.Println("READ CONFIG")

				w, e := conf.Validate()
				r.NoError(e)
				r.Len(w, 0)

				e = apply.Apply(testdataFs, conf, templates.Templates, true)
				r.NoError(e)
			} else {
				fileName := "fogg.yml"
				fs, _, e := util.TestFs()
				r.NoError(e)

				// copy tc.configFile into the tmp test dir (so that it doesn't show up as a diff)
				configContents, e := afero.ReadFile(testdataFs, fileName)
				if os.IsNotExist(e) { //If the error is related to the file being non-existent
					fileName = tc.configFile
					configContents, e = afero.ReadFile(testdataFs, fileName)
				}
				r.NoError(e)

				configMode, e := testdataFs.Stat(fileName)
				r.NoError(e)
				r.NoError(afero.WriteFile(fs, fileName, configContents, configMode.Mode()))

				conf, e := config.FindAndReadConfig(fs, fileName)
				r.NoError(e)
				fmt.Printf("conf %#v\n", conf)

				w, e := conf.Validate()
				r.NoError(e)
				r.Len(w, 0)

				e = apply.Apply(fs, conf, templates.Templates, true)
				r.NoError(e)

				r.NoError(afero.Walk(testdataFs, ".", func(path string, info os.FileInfo, err error) error {
					logrus.Debug("================================================")
					logrus.Debug(path)
					if !info.Mode().IsRegular() {
						logrus.Debug("dir or link")
					} else {
						i1, e1 := testdataFs.Stat(path)
						r.NotNil(i1)
						r.NoError(e1)

						i2, e2 := fs.Stat(path)
						r.NoError(e2)
						r.NotNil(i2)

						logrus.Debugf("i1 size: %d ii2 size %d", i1.Size(), i2.Size())
						r.Equalf(i1.Size(), i2.Size(), "file size: %s", path)
						// This (below) doesn't currently work for files created on a mac then tested on linux. :shrug:
						// r.Equalf(i1.Mode(), i2.Mode(), "file mode: %s, %o vs %o", path, i1.Mode(), i2.Mode())

						f1, e3 := afero.ReadFile(testdataFs, path)
						r.NoError(e3)
						f2, e4 := afero.ReadFile(fs, path)
						r.NoError(e4)

						logrus.Debugf("f1:\n%s\n\n---- ", f1)
						logrus.Debugf("f2:\n%s\n\n---- ", f2)

						r.Equal(f1, f2, path)
					}
					return nil
				}))
			}
		})
	}
}
