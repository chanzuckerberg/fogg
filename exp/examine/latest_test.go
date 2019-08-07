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

	localModules, err := GetLocalModules(fs, "../../testdata/version_detection/terraform/envs/staging/app/")
	r.NoError(err)
	r.NotNil(localModules)

	globalModules, err := LatestModuleVersions(fs, localModules)
	r.NoError(err)
	r.NotNil(globalModules)
}

//DISABLED
// func TestGetModuleVersion(t *testing.T) {
// 	r := require.New(t)
// 	module := GetModule()
// 	r.Equal("4.1.0", module.Version)
// }

//DISABLED
// func TestFindModules(t *testing.T) {
// repo := "github.com/chanzuckerberg/cztack//aws-params-reader-policy?ref=v0.15.1"
// localModules := GetLocalModules(repo)
// if localModules == nil{

// }

// registryModules := searchForModules(localModules, awsModules)

// fmt.Println(registryModules)
// }
