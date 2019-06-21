package v2

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"

	v1 "github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/errs"
)

func ReadConfig(b []byte) (*Config, error) {
	c := &Config{
		Docker: false,
	}
	e := json.Unmarshal(b, c)

	return c, errs.WrapUser(e, "unable to parse config")
}

type Config struct {
	Accounts map[string]Account   `json:"accounts,omitempty"`
	Defaults Defaults             `json:"defaults" validate:"required"`
	Docker   bool                 `json:"docker,omitempty"`
	Envs     map[string]Env       `json:"envs,omitempty"`
	Global   Component            `json:"global,omitempty"`
	Modules  map[string]v1.Module `json:"modules,omitempty"`
	Plugins  v1.Plugins           `json:"plugins,omitempty"`
	Version  int                  `json:"version" validate:"required,eq=2"`
}

type Common struct {
	Backend          *Backend          `json:"backend,omitempty"`
	ExtraVars        map[string]string `json:"extra_vars,omitempty"`
	Owner            *string           `json:"owner,omitempty" `
	Project          *string           `json:"project,omitempty" `
	Providers        *Providers        `json:"providers,omitempty" `
	TerraformVersion *string           `json:"terraform_version,omitempty"`
	Tools            *Tools            `json:"tools,omitempty"`
}

type Defaults struct {
	Common
}

type Account struct {
	Common
}

type Tools struct {
	Atlantis *Atlantis    `json:"atlantis,omitempty"`
	TravisCI *v1.TravisCI `json:"travis_ci,omitempty"`
	TfLint   *v1.TfLint   `json:"tflint,omitempty"`
}

type Atlantis struct {
	Enabled  *bool   `json:"enabled,omitempty"`
	RoleName *string `json:"role_name,omitempty"`
	RolePath *string `json:"role_path,omitempty"`
}

type Env struct {
	Common

	Components map[string]Component `json:"components,omitempty"`
}

type Component struct {
	Common

	EKS          *v1.EKSConfig     `json:"eks,omitempty"`
	Kind         *v1.ComponentKind `json:"kind,omitempty"`
	ModuleSource *string           `json:"module_source,omitempty"`
}

type Providers struct {
	AWS       *AWSProvider       `json:"aws,omitempty"`
	Snowflake *SnowflakeProvider `json:"snowflake,omitempty"`
	Bless     *BlessProvider     `json:"bless,omitempty"`
	Okta      *OktaProvider      `json:"okta,omitempty"`
}

// OktaProvider is an okta provider
type OktaProvider struct {
	// the okta provider is optional (above) but if supplied you must set an OrgName
	OrgName *string `json:"org_name,omitempty"`
	Version *string `json:"version,omitempty"`
}

// BlessProvider allows for terraform-provider-bless configuration
type BlessProvider struct {
	// the bless provider is optional (above) but if supplied you must set a region and aws_profile
	AdditionalRegions []string `json:"additional_regions,omitempty"`
	AWSProfile        *string  `json:"aws_profile,omitempty"`
	AWSRegion         *string  `json:"aws_region,omitempty"`
	Version           *string  `json:"version,omitempty"`
}

type AWSProvider struct {
	// the aws provider is optional (above) but if supplied you must set account id and region
	AccountID         *json.Number `json:"account_id,omitempty"`
	AdditionalRegions []string     `json:"additional_regions,omitempty"`
	Profile           *string      `json:"profile,omitempty"`
	Region            *string      `json:"region,omitempty"`
	Version           *string      `json:"version,omitempty"`
}

type SnowflakeProvider struct {
	Account *string `json:"account,omitempty"`
	Role    *string `json:"role,omitempty"`
	Region  *string `json:"region,omitempty"`
	Version *string `json:"version,omitempty"`
}

type Backend struct {
	AccountID   *string `json:"account_id,omitempty"`
	Bucket      *string `json:"bucket,omitempty"`
	DynamoTable *string `json:"dynamodb_table,omitempty"`
	Profile     *string `json:"profile,omitempty"`
	Region      *string `json:"region,omitempty"`
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
					Enabled:     r.Float32() < 0.5,
					TestBuckets: r.Intn(size),
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
