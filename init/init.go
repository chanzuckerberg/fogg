package init

import (
	"github.com/chanzuckerberg/fogg/config"
	"github.com/spf13/afero"
)

const AWSProviderVersion = "4.34.0"

type FoggProject struct {
	Project, Region, Bucket, Table, Profile, Owner, AwsAccountID *string
}

// Init reads user console input and generates a fogg.yml file
func Init(fs afero.Fs, foggProject *FoggProject) error {
	config := config.InitConfig(
		foggProject.Project,
		foggProject.Region,
		foggProject.Bucket,
		foggProject.Table,
		foggProject.Profile,
		foggProject.Owner,
		foggProject.AwsAccountID,
		AWSProviderVersion,
	)
	e := config.Write(fs, "fogg.yml")
	return e
}
