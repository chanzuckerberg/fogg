package examine

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

//TODO: Move fs to versioning.go
func TestGetLocalModules(t *testing.T) {
	r := require.New(t)
	pwd, err := os.Getwd()
	r.NoError(err)
	pwd = "../../" + pwd
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

	modules, err := GetLocalModules(fs, "../../testdata/version_detection/terraform/envs/staging/app/")
	r.NoError(err)
	r.NotNil(modules)
}
func TestGetCztackModuleFromGithub(t *testing.T) {
	r := require.New(t)
	pwd, _ := os.Getwd()
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

	repo := "github.com/chanzuckerberg/cztack//aws-params-reader-policy?ref=v0.15.1"
	mod, err := GetFromGithub(fs, repo)
	r.NoError(err)
	r.NotNil(mod)
}

func TestDownloadModuleFromRegistry(t *testing.T) {
	r := require.New(t)
	pwd, _ := os.Getwd()
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

	path := "terraform-aws-modules/security-group/aws"
	mod, err := downloadModule(fs, path, "2.9.0")
	r.NoError(err)
	r.NotNil(mod)
}

func TestGetLocalRegistryModuleFromRegistry(t *testing.T) {
	r := require.New(t)
	pwd, _ := os.Getwd()
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

	path := "../../testdata/version_detection/terraform/modules/test-component-2"
	mods, err := getAllModules(fs, path)
	r.NoError(err)
	r.NotNil(mods)
}
