package versioning

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
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

	modules, err := GetLocalModules(fs, "/Users/echanakira/Desktop/learning/shared-infra/terraform/envs/staging/golinks/")
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

//FIXME: Test cannot retrieve the repo
// func TestGetRegistryModuleFromGithub(t *testing.T) {
// 	r := require.New(t)

// 	repo := "github.com/terraform-aws-modules/terraform-aws-security-group?ref=v3.1.0"
// 	mod, err := GetFromGithub(repo)
// 	r.NoError(err)
// 	r.NotNil(mod)
// }

func TestDownloadRegistryFromGithub(t *testing.T) {
	r := require.New(t)
	pwd, _ := os.Getwd()
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

	path := "terraform-aws-modules/security-group/aws"
	mod, err := downloadModule(fs, path, "2.9.0")
	r.NoError(err)
	r.NotNil(mod)
}

func TestGetFromAlbHttp(t *testing.T) {
	r := require.New(t)
	pwd, _ := os.Getwd()
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

	path := "/Users/echanakira/Desktop/learning/shared-infra/terraform/modules/alb-http/"
	mods, err := getAllModules(fs, path)
	r.NoError(err)
	r.NotNil(mods)
}
