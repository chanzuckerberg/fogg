package v2

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"path/filepath"
	"reflect"

	v1 "github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

//ReadConfig take a byte array as input and outputs a json or yaml config file
func ReadConfig(fs afero.Fs, b []byte, configFile string) (*Config, error) {
	var e error
	c := &Config{
		Docker: false,
	}

	info, e := fs.Stat(configFile)
	if e != nil {
		return nil, errs.WrapUser(e, "unable to find file")
	}

	ext := filepath.Ext(info.Name())
	//Determines the file extension
	switch ext {
	case ".yml", ".yaml":
		e = yaml.Unmarshal(b, c)
	case ".json":
		e = json.Unmarshal(b, c)
	default:
		return nil, errs.NewUserf("File type %s is not supported", ext)
	}
	return c, e
}

type Config struct {
	Accounts map[string]Account   `json:"accounts,omitempty" yaml:"accounts,omitempty"`
	Defaults Defaults             `json:"defaults"  yaml:"defaults" validate:"required"`
	Docker   bool                 `json:"docker,omitempty" yaml:"docker,omitempty"`
	Envs     map[string]Env       `json:"envs,omitempty" yaml:"envs,omitempty"`
	Global   Component            `json:"global,omitempty" yaml:"global,omitempty"`
	Modules  map[string]v1.Module `json:"modules,omitempty" yaml:"modules,omitempty"`
	Plugins  v1.Plugins           `json:"plugins,omitempty" yaml:"plugins,omitempty"`
	Version  int                  `json:"version" yaml:"version" validate:"required,eq=2"`
}

type Common struct {
	Backend          *Backend          `json:"backend,omitempty" yaml:"backend,omitempty"`
	ExtraVars        map[string]string `json:"extra_vars,omitempty" yaml:"extra_vars,omitempty"`
	Owner            *string           `json:"owner,omitempty" yaml:"owner,omitempty" `
	Project          *string           `json:"project,omitempty"  yaml:"project,omitempty"`
	Providers        *Providers        `json:"providers,omitempty" yaml:"providers,omitempty" `
	TerraformVersion *string           `json:"terraform_version,omitempty" yaml:"terraform_version,omitempty"`
	Tools            *Tools            `json:"tools,omitempty" yaml:"tools,omitempty"`
}

type Defaults struct {
	Common `json:",inline" yaml:",inline"`
}

type Account struct {
	Common `json:",inline" yaml:",inline"`
}

type Tools struct {
	Atlantis *Atlantis    `json:"atlantis,omitempty" yaml:"atlantis,omitempty"`
	TravisCI *v1.TravisCI `json:"travis_ci,omitempty" yaml:"travis_ci,omitempty"`
	TfLint   *v1.TfLint   `json:"tflint,omitempty" yaml:"tflint,omitempty"`
}

type Atlantis struct {
	Enabled  *bool   `json:"enabled,omitempty"`
	RoleName *string `json:"role_name,omitempty"`
	RolePath *string `json:"role_path,omitempty"`
}

type Env struct {
	Common `json:",inline" yaml:",inline"`

	Components map[string]Component `json:"components,omitempty" yaml:"components,omitempty"`
}

type Component struct {
	Common `json:",inline" yaml:",inline"`

	EKS          *v1.EKSConfig     `json:"eks,omitempty" yaml:"eks,omitempty"`
	Kind         *v1.ComponentKind `json:"kind,omitempty" yaml:"kind,omitempty"`
	ModuleSource *string           `json:"module_source,omitempty" yaml:"module_source,omitempty"`
}

type Providers struct {
	AWS       *AWSProvider       `json:"aws,omitempty" yaml:"aws,omitempty"`
	Snowflake *SnowflakeProvider `json:"snowflake,omitempty" yaml:"snowflake,omitempty"`
	Bless     *BlessProvider     `json:"bless,omitempty" yaml:"bless,omitempty"`
	Okta      *OktaProvider      `json:"okta,omitempty" yaml:"okta,omitempty"`
}

// OktaProvider is an okta provider
type OktaProvider struct {
	// the okta provider is optional (above) but if supplied you must set an OrgName
	OrgName *string `json:"org_name,omitempty" yaml:"org_name,omitempty"`
	Version *string `json:"version,omitempty" yaml:"version,omitempty"`
}

// BlessProvider allows for terraform-provider-bless configuration
type BlessProvider struct {
	// the bless provider is optional (above) but if supplied you must set a region and aws_profile
	AdditionalRegions []string `json:"additional_regions,omitempty" yaml:"additional_regions,omitempty"`
	AWSProfile        *string  `json:"aws_profile,omitempty" yaml:"aws_profile,omitempty"`
	AWSRegion         *string  `json:"aws_region,omitempty" yaml:"aws_region,omitempty"`
	Version           *string  `json:"version,omitempty" yaml:"version,omitempty"`
}

type AWSProvider struct {
	// the aws provider is optional (above) but if supplied you must set account id and region
	AccountID         *json.Number `json:"account_id,omitempty" yaml:"account_id,omitempty"`
	AdditionalRegions []string     `json:"additional_regions,omitempty" yaml:"additional_regions,omitempty"`
	Profile           *string      `json:"profile,omitempty" yaml:"profile,omitempty"`
	Region            *string      `json:"region,omitempty" yaml:"region,omitempty"`
	Version           *string      `json:"version,omitempty" yaml:"version,omitempty"`
}

type SnowflakeProvider struct {
	Account *string `json:"account,omitempty" yaml:"account,omitempty"`
	Role    *string `json:"role,omitempty" yaml:"role,omitempty"`
	Region  *string `json:"region,omitempty" yaml:"region,omitempty"`
	Version *string `json:"version,omitempty" yaml:"version,omitempty"`
}

type Backend struct {
	AccountID   *string `json:"account_id,omitempty" yaml:"account_id,omitempty"`
	Bucket      *string `json:"bucket,omitempty" yaml:"bucket,omitempty"`
	DynamoTable *string `json:"dynamodb_table,omitempty" yaml:"dynamodb_table,omitempty"`
	Profile     *string `json:"profile,omitempty" yaml:"profile,omitempty"`
	Region      *string `json:"region,omitempty" yaml:"region,omitempty"`
}

// Generate is used for test/quick integration. There are supposedly ways to do this without polluting the public
// api, but givent that fogg isn't an api, it doesn't seem like a big deal
func (c *Config) Generate(r *rand.Rand, size int) reflect.Value {
	// TODO write this to be part of tests https://github.com/shiwano/submarine/blob/5c02c0cfcf05126454568ef9624550eb0d84f86c/server/battle/src/battle/util/util_test.go#L19

	fmt.Println("generate")
	conf := &Config{}

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	randString := func(r *rand.Rand, n int) string {
		b := make([]byte, n)
		for i := range b {
			b[i] = letterBytes[r.Intn(len(letterBytes))]
		}
		return string(b)
	}

	randNonEmptyString := func(r *rand.Rand, s int) string {
		return "asdf"
	}

	randStringPtr := func(r *rand.Rand, s int) *string {
		str := randString(r, s)
		return &str
	}

	randBoolPtr := func(r *rand.Rand, s int) *bool {
		b := r.Float32() > 0.5
		return &b
	}
	randIntPtr := func(r *rand.Rand, s int) *int {
		i := r.Intn(s)
		return &i
	}

	randStringMap := func(r *rand.Rand, s int) map[string]string {
		m := map[string]string{}

		for i := 0; i < s; i++ {
			m[randNonEmptyString(r, s)] = randString(r, s)
		}

		return map[string]string{}
	}

	randOktaProvider := func(r *rand.Rand, s int) *OktaProvider {
		if r.Float32() < 0.5 {
			return nil
		}
		return &OktaProvider{
			Version: randStringPtr(r, s),
			OrgName: randStringPtr(r, s),
		}
	}

	randBlessProvider := func(r *rand.Rand, s int) *BlessProvider {
		if r.Float32() < 0.5 {
			return nil
		}
		return &BlessProvider{
			Version:           randStringPtr(r, s),
			AWSRegion:         randStringPtr(r, s),
			AWSProfile:        randStringPtr(r, s),
			AdditionalRegions: []string{randString(r, s)},
		}
	}

	randAWSProvider := func(r *rand.Rand, s int) *AWSProvider {
		if r.Float32() < 0.5 {
			accountID := json.Number(randString(r, s))
			return &AWSProvider{
				AccountID: &accountID,
				Region:    randStringPtr(r, s),
				Profile:   randStringPtr(r, s),
				Version:   randStringPtr(r, s),
			}
		}
		return nil
	}

	randSnowflakeProvider := func(r *rand.Rand, s int) *SnowflakeProvider {
		if r.Float32() < 0.5 {
			return &SnowflakeProvider{
				Account: randStringPtr(r, size),
				Region:  randStringPtr(r, s),
				Role:    randStringPtr(r, s),
			}
		}
		return nil
	}

	randCommon := func(r *rand.Rand, s int) Common {
		c := Common{
			Backend: &Backend{
				Bucket: randStringPtr(r, s),
				Region: randStringPtr(r, s),
			},
			ExtraVars: randStringMap(r, s),
			Owner:     randStringPtr(r, s),
			Project:   randStringPtr(r, s),
			Providers: &Providers{
				AWS:       randAWSProvider(r, s),
				Snowflake: randSnowflakeProvider(r, s),
				Okta:      randOktaProvider(r, s),
				Bless:     randBlessProvider(r, s),
			},
			TerraformVersion: randStringPtr(r, s),
		}

		if r.Float32() < 0.5 {
			c.Tools = &Tools{}
			if r.Float32() < 0.5 {
				c.Tools.TravisCI = &v1.TravisCI{
					Enabled:     randBoolPtr(r, s),
					TestBuckets: randIntPtr(r, s),
				}
			}
			if r.Float32() < 0.5 {
				p := r.Float32() < 0.5
				c.Tools.TfLint = &v1.TfLint{
					Enabled: &p,
				}
			}
			if r.Float32() < 0.5 {
				t := true
				c.Tools.Atlantis = &Atlantis{
					Enabled:  &t,
					RolePath: randStringPtr(r, s),
					RoleName: randStringPtr(r, s),
				}
			}
		}

		return c
	}

	conf.Version = 2

	conf.Defaults = Defaults{
		Common: randCommon(r, size),
	}

	// tools

	conf.Accounts = map[string]Account{}
	acctN := r.Intn(size)

	for i := 0; i < acctN; i++ {
		acctName := randString(r, size)
		conf.Accounts[acctName] = Account{
			Common: randCommon(r, size),
		}

	}

	conf.Envs = map[string]Env{}
	envN := r.Intn(size)

	for i := 0; i < envN; i++ {
		envName := randString(r, size)
		e := Env{
			Common: randCommon(r, size),
		}
		e.Components = map[string]Component{}
		compN := r.Intn(size)

		for i := 0; i < compN; i++ {
			compName := randString(r, size)
			e.Components[compName] = Component{
				Common: randCommon(r, size),
			}
		}
		conf.Envs[envName] = e

	}

	return reflect.ValueOf(conf)
}
