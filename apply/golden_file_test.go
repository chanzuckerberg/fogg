package apply_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
		fileName string
	}{
		{"auth0_provider_yaml"},
		{"okta_provider_yaml"},
		{"github_provider_yaml"},
		{"bless_provider_yaml"},
		{"snowflake_provider_yaml"},
		{"v2_full_yaml"},
		{"v2_minimal_valid_yaml"},
		{"v2_split_yaml"},
		{"v2_no_aws_provider_yaml"},
		{"github_actions"},
		{"github_actions_with_iam_role"},
		{"circleci"},
		{"tfe_provider_yaml"},
		{"remote_backend_yaml"},
		{"tfe_config"},
		{"v2_aws_default_tags"},
		{"v2_aws_ignore_tags"},
		{"v2_tf_registry_module"},
		{"v2_tf_registry_module_atlantis"},
		{"v2_tf_registry_module_atlantis_dup_module"},
		{"v2_integration_registry"},
		{"v2_github_actions_with_pre_commit"},
		{"generic_providers_yaml"},
	}

	for _, test := range testCases {
		tt := test

		t.Run(tt.fileName, func(t *testing.T) {
			r := require.New(t)

			testdataFs := afero.NewBasePathFs(afero.NewOsFs(), filepath.Join(util.ProjectRoot(), "testdata", tt.fileName))
			configFile := "fogg.yml"
			if *updateGoldenFiles {
				// delete all files except fogg.yml and conf.d, foo_modules directories
				e := afero.Walk(testdataFs, ".", func(path string, info os.FileInfo, err error) error {
					if !info.IsDir() && !(path == configFile) && !(strings.Contains(path, "fogg.d")) && !(strings.Contains(path, "foo_modules")) {
						return testdataFs.Remove(path)
					}
					return nil
				})
				r.NoError(e)

				conf, e := config.FindAndReadConfig(testdataFs, configFile)
				r.NoError(e)
				fmt.Printf("conf %#v\n", conf)
				fmt.Println("READ CONFIG")

				w, e := conf.Validate()
				r.NoError(e)
				r.Len(w, 0)

				e = apply.Apply(testdataFs, conf, templates.Templates, true)
				r.NoError(e)
			} else {
				fs, _, e := util.TestFs()
				r.NoError(e)

				// copy fogg.yml into the tmp test dir (so that it doesn't show up as a diff)
				configContents, e := afero.ReadFile(testdataFs, configFile)
				r.NoError(e)
				configMode, e := testdataFs.Stat(configFile)
				r.NoError(e)
				r.NoError(afero.WriteFile(fs, configFile, configContents, configMode.Mode()))
				// if fogg.d exists, copy all partial configs too
				confDir, e := testdataFs.Stat("fogg.d")
				fs.Mkdir("fogg.d", 0700)
				if e == nil && confDir.IsDir() {
					afero.Walk(testdataFs, "fogg.d", func(path string, info os.FileInfo, err error) error {
						if !info.IsDir() {
							partialConfigContents, e := afero.ReadFile(testdataFs, path)
							r.NoError(e)
							r.NoError(afero.WriteFile(fs, path, partialConfigContents, info.Mode()))
							return nil
						}
						return nil
					})
				}
				// if foo_modules exists, copy these too...
				fooModulesDir, e := testdataFs.Stat("foo_modules")
				fs.Mkdir("foo_modules", 0700)
				if e == nil && fooModulesDir.IsDir() {
					afero.Walk(testdataFs, "foo_modules", func(path string, info os.FileInfo, err error) error {
						if !info.IsDir() {
							moduleFileContents, e := afero.ReadFile(testdataFs, path)
							r.NoError(e)
							r.NoError(afero.WriteFile(fs, path, moduleFileContents, info.Mode()))
							return nil
						} else {
							fs.Mkdir(path, 0700)
						}
						return nil
					})
				}

				conf, e := config.FindAndReadConfig(fs, configFile)
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
