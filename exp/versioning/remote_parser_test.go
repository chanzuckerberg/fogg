package versioning

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestModuleStruct(t *testing.T) {
	r := require.New(t)
	module := GetModule()
	r.NotNil(module)

}

func TestGetModuleVersion(t *testing.T) {
	r := require.New(t)
	module := GetModule()
	r.Equal("4.1.0", module.Version)

}

func TestFindModules(t *testing.T) {
	// repo := "github.com/chanzuckerberg/cztack//aws-params-reader-policy?ref=v0.15.1"
	// localModules := GetLocalModules(repo)
	// if localModules == nil{

	// }

	GetAWSModules()

	// registryModules := searchForModules(localModules, awsModules)

	// fmt.Println(registryModules)
}
