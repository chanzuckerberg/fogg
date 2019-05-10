package apply_test

import (
	"flag"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/chanzuckerberg/fogg/apply"
	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/templates"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var updateGoldenFiles = flag.Bool("update", false, "when set, rewrite the golden files")

func TestIntegration(t *testing.T) {

	var testCases = []struct {
		fileName string
	}{
		{"empty"},
	}

	for _, tc := range testCases {
		t.Run(tc.fileName, func(t *testing.T) {
			a := assert.New(t)

			if *updateGoldenFiles {

				fs := afero.NewBasePathFs(afero.NewOsFs(), filepath.Join(util.ProjectRoot(), "testdata", tc.fileName))

				conf, e := config.FindAndReadConfig(fs, "fogg.json")
				a.NoError(e)
				fmt.Printf("conf %#v\n", conf)

				e = conf.Validate()
				a.NoError(e)

				e = apply.Apply(fs, conf, templates.Templates, true)
				a.NoError(e)
			}
			// else
			//	run apply and diff
		})
	}
}
