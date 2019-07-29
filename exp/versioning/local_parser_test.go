package versioning

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetLocalModules(t *testing.T) {
	r := require.New(t)
	modules := GetLocalModules("/Users/echanakira/Desktop/learning/shared-infra/terraform/envs/staging/golinks/")
	r.NotNil(modules)
}

func TestGetCztackModuleFromGithub(t *testing.T) {
	r := require.New(t)

	repo := "github.com/chanzuckerberg/cztack//aws-params-reader-policy?ref=v0.15.1"
	mod, err := getFromGithub(repo)
	r.NoError(err)
	r.NotNil(mod)
}

//FIXME: Broken Test
// func TestGetRegistryModuleFromGithub(t *testing.T) {
// 	r := require.New(t)

// 	repo := "github.com/terraform-aws-modules/terraform-aws-security-group?ref=v3.1.0"
// 	mod, err := getFromGithub(repo)
// 	r.NoError(err)
// 	r.NotNil(mod)
// }

func TestDownloadRegistryFromGithub(t *testing.T) {
	r := require.New(t)

	path := "terraform-aws-modules/security-group/aws"
	mod, err := downloadModule(path, "2.9.0")
	r.NoError(err)
	r.NotNil(mod)
}

func TestGetFromAlbHttp(t *testing.T) {
	r := require.New(t)

	path := "/Users/echanakira/Desktop/learning/shared-infra/terraform/modules/alb-http/"
	mods, err := retrieveAllDependencies(path)
	r.NoError(err)
	r.NotNil(mods)
}
