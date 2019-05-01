package config

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/spf13/afero"
)

func InitConfig(project, region, bucket, table, awsProfile, owner, awsProviderVersion string) *v1.Config {
	return &v1.Config{
		Defaults: v1.Defaults{
			AWSProfileBackend:  awsProfile,
			AWSProfileProvider: awsProfile,
			AWSProviderVersion: awsProviderVersion,
			AWSRegionBackend:   region,
			AWSRegionProvider:  region,
			ExtraVars:          map[string]string{},
			InfraBucket:        bucket,
			InfraDynamoTable:   table,
			Owner:              owner,
			Project:            project,
			TerraformVersion:   "0.11.7",
		},
		Accounts: map[string]v1.Account{},
		Docker:   true,
		Envs:     map[string]v1.Env{},
		Modules:  map[string]v1.Module{},
	}
}

func ReadConfig(f io.Reader) (*v1.Config, error) {
	c := &v1.Config{
		Docker: true,
	}
	b, e := ioutil.ReadAll(f)
	if e != nil {
		return nil, errs.WrapUser(e, "unable to read config")
	}
	e = json.Unmarshal(b, c)
	if e != nil {
		return nil, errs.WrapUser(e, "unable to parse json config file")
	}
	return c, nil
}

func FindAndReadConfig(fs afero.Fs, configFile string) (*v1.Config, error) {
	f, e := fs.Open(configFile)
	if e != nil {
		return nil, errs.WrapUser(e, "unable to open config file")
	}
	reader := io.ReadCloser(f)
	defer reader.Close()
	c, err2 := ReadConfig(reader)
	return c, err2
}
