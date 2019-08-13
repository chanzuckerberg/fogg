package examine

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestCompareLocalAndGlobal(t *testing.T) {
	r := require.New(t)
	pwd, err := os.Getwd()
	r.NoError(err)
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

	config, err := GetLocalModules(fs, "../../testdata/version_detection/terraform/envs/staging/app/")
	r.NoError(err)
	r.NotNil(config)

	globalModules, err := LatestModuleVersions(fs, config)
	r.NoError(err)
	r.NotNil(globalModules)
}

func TestCreateGitUrl(t *testing.T) {
	r := require.New(t)
	pwd, err := os.Getwd()
	r.NoError(err)
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

	config, err := GetLocalModules(fs, "../../testdata/version_detection/terraform/envs/staging/app/")
	r.NoError(err)
	r.NotNil(config)

	url, err := createGitUrl(config.Modules[1])
	r.NoError(err)
	r.NotEmpty(url)
}
