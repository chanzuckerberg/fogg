package config

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/config/v2"
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

func FindAndReadConfig(fs afero.Fs, configFile string) (*v2.Config, error) {
	f, err := fs.Open(configFile)
	if err != nil {
		return nil, errs.WrapUser(err, "unable to open config file")
	}
	reader := io.ReadCloser(f)
	defer reader.Close()

	c, err := ReadConfig(reader)
	if err != nil {
		return nil, err
	}

	c2, err := UpgradeConfigVersion(c)

	return c2, err
}

type ver struct {
	Version int `json:"version"`
}

func detectVersion(b []byte) (int, error) {
	v := &ver{}
	err := json.Unmarshal(b, v)
	if err != nil {
		return 0, err
	}
	if v.Version == 0 {
		return 1, nil
	}
	return v.Version, nil
}

// UpgradeConfigVersion will convert a v1.Config to a v2.Config. Note that due to some semantic changes, this a lossy
//  conversion
func UpgradeConfigVersion(c1 *v1.Config) (*v2.Config, error) {
	c2 := &v2.Config{
		Version:  2,
		Docker:   c1.Docker,
		Accounts: map[string]v2.Account{},
		Envs:     map[string]v2.Env{},
	}

	def1 := c1.Defaults
	c2.Defaults = v2.Defaults{
		Common: v2.Common{
			Backend: v2.Backend{
				Bucket:      def1.InfraBucket,
				DynamoTable: def1.InfraDynamoTable,
				Profile:     def1.AWSProfileBackend,
				Region:      def1.AWSRegionBackend,
			},
			ExtraVars: def1.ExtraVars,
			Providers: v2.Providers{
				AWS: &v2.AWSProvider{
					AccountID:         &def1.AccountID,
					AdditionalRegions: def1.AWSRegions,
					Profile:           &def1.AWSProfileProvider,
					Region:            &def1.AWSRegionProvider,
					Version:           &def1.AWSProviderVersion,
				},
			},
			Owner:            def1.Owner,
			Project:          def1.Project,
			TerraformVersion: def1.TerraformVersion,
		},
	}

	c2.Tools = v2.Tools{}
	if def1.TfLint != nil {
		c2.Tools.TfLint = def1.TfLint
	}

	if c1.TravisCI != nil {
		c2.Tools.TravisCI = c1.TravisCI
	}

	var ptrstr = func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}

	for acctName, acct := range c1.Accounts {
		c2.Accounts[acctName] = v2.Account{
			Common: v2.Common{
				Backend: v2.Backend{
					Bucket:      ptrstr(acct.InfraBucket),
					DynamoTable: ptrstr(acct.InfraDynamoTable),
					Profile:     ptrstr(acct.AWSProfileBackend),
					Region:      ptrstr(acct.AWSRegionBackend),
				},
				ExtraVars: acct.ExtraVars,
				Providers: v2.Providers{
					AWS: &v2.AWSProvider{
						AccountID:         acct.AccountID,
						AdditionalRegions: acct.AWSRegions,
						Profile:           acct.AWSProfileProvider,
						Region:            acct.AWSRegionProvider,
						Version:           acct.AWSProviderVersion,
					},
				},
				Owner:            ptrstr(acct.Owner),
				Project:          ptrstr(acct.Project),
				TerraformVersion: ptrstr(acct.TerraformVersion),
			},
		}
	}

	for envName, env := range c1.Envs {
		env2 := v2.Env{
			Common: v2.Common{
				Backend: v2.Backend{
					Bucket:      ptrstr(env.InfraBucket),
					DynamoTable: ptrstr(env.InfraDynamoTable),
					Profile:     ptrstr(env.AWSProfileBackend),
					Region:      ptrstr(env.AWSRegionBackend),
				},
				ExtraVars: env.ExtraVars,
				Providers: v2.Providers{
					AWS: &v2.AWSProvider{
						AccountID:         env.AccountID,
						AdditionalRegions: env.AWSRegions,
						Profile:           env.AWSProfileProvider,
						Region:            env.AWSRegionProvider,
						Version:           env.AWSProviderVersion,
					},
				},
				Owner:            ptrstr(env.Owner),
				Project:          ptrstr(env.Project),
				TerraformVersion: ptrstr(env.TerraformVersion),
			},
		}
		env2.Components = map[string]v2.Component{}

		for componentName, component := range env.Components {
			c2 := v2.Component{
				Common: v2.Common{
					Backend: v2.Backend{
						Bucket:      ptrstr(component.InfraBucket),
						DynamoTable: ptrstr(component.InfraDynamoTable),
						Profile:     ptrstr(component.AWSProfileBackend),
						Region:      ptrstr(component.AWSRegionBackend),
					},
					ExtraVars: component.ExtraVars,
					Providers: v2.Providers{
						// see below
					},
					Owner:            ptrstr(component.Owner),
					Project:          ptrstr(component.Project),
					TerraformVersion: ptrstr(component.TerraformVersion),
				},
				EKS:          component.EKS,
				Kind:         component.Kind,
				ModuleSource: component.ModuleSource,
			}

			if component.AccountID != nil || component.AWSRegions != nil || component.AWSProfileProvider != nil || component.AWSRegionProvider != nil || component.TerraformVersion != nil {
				c2.Providers.AWS = &v2.AWSProvider{
					AccountID:         component.AccountID,
					AdditionalRegions: component.AWSRegions,
					Profile:           component.AWSProfileProvider,
					Region:            component.AWSRegionProvider,
					Version:           component.AWSProviderVersion,
				}
			}
			env2.Components[componentName] = c2
		}
		c2.Envs[envName] = env2
	}

	// we don't bother doing a deep copy bc we assume these structures are immutable
	c2.Modules = c1.Modules
	c2.Plugins = c1.Plugins

	return c2, nil
}
