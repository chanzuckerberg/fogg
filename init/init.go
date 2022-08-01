package init

import (
	"github.com/chanzuckerberg/fogg/config"
	"github.com/spf13/afero"
)

const AWSProviderVersion = "2.47.0"

//Init reads user console input and generates a fogg.yml file
func Init(fs afero.Fs, project, region, bucket, table, awsProfile, owner string) error {
	config := config.InitConfig(project, region, bucket, table, awsProfile, owner, AWSProviderVersion)
	e := config.Write(fs, "fogg.yml")
	return e
}
