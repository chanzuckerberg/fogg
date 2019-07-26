package versioning

import(
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)


func TestModuleStruct(t *testing.T){
	a := assert.New(t)
	module := GetModule();

	a.NotNil(module)
	fmt.Println(module)
}

func TestGetModuleVersion(t *testing.T){
	a := assert.New(t)
	module := GetModule();
	
	a.Equal("4.1.0",module.Version)
}


func TestFindModules(t *testing.T){
	// repo := "github.com/chanzuckerberg/cztack//aws-params-reader-policy?ref=v0.15.1"
	// localModules := GetLocalModules(repo)
	// if localModules == nil{

	// }

	GetAWSModules()

	// registryModules := searchForModules(localModules, awsModules)
	
	// fmt.Println(registryModules)
}