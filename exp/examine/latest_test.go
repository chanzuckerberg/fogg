package examine

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestCompareLocalAndGlobal(t *testing.T) {
	t.Skip("This hits rate limits in CI. See https://travis-ci.com/chanzuckerberg/fogg/builds/123534108.")
	r := require.New(t)
	pwd, err := os.Getwd()
	r.NoError(err)
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

	module, err := GetLocalModules(fs, "../../testdata/version_detection/terraform/envs/staging/app/")
	r.NoError(err)
	r.NotNil(module)

	globalModules, err := LatestModuleVersions(fs, module)
	r.NoError(err)
	r.NotNil(globalModules)
}

func TestCreateGitUrl(t *testing.T) {
	t.Skip("This hits rate limits in CI. See https://travis-ci.com/chanzuckerberg/fogg/builds/123534108.")
	r := require.New(t)
	pwd, err := os.Getwd()
	r.NoError(err)
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

	module, err := GetLocalModules(fs, "../../testdata/version_detection/terraform/envs/staging/app/")
	r.NoError(err)
	r.NotNil(module)

	url, err := createGitURL(module.ModuleCalls["parameters-policy"])
	r.NoError(err)
	r.NotEmpty(url)
}
