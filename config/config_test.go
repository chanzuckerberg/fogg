package config

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/plugins"
	"github.com/stretchr/testify/assert"
	validator "gopkg.in/go-playground/validator.v9"
)

func TestParseDefaults(t *testing.T) {
	json := `
	{
		"defaults": {
			"aws_region_backend": "us-west-2",
			"aws_region_provider": "us-west-1",
			"aws_profile_backend": "czi",
			"aws_profile_provider": "czi",
			"infra_s3_bucket": "the-bucket",
			"infra_dynamo_db_table": "the-table",
			"project": "test-project",
			"terraform_version": "0.11.0"
		}
	}`
	r := ioutil.NopCloser(strings.NewReader(json))
	defer r.Close()
	c, e := ReadConfig(r)
	assert.NoError(t, e)
	assert.NotNil(t, c.Defaults)
	assert.Equal(t, "us-west-2", c.Defaults.AWSRegionBackend)
	assert.Equal(t, "us-west-1", c.Defaults.AWSRegionProvider)
	assert.Equal(t, "czi", c.Defaults.AWSProfileBackend)
	assert.Equal(t, "czi", c.Defaults.AWSProfileProvider)
	assert.Equal(t, "the-bucket", c.Defaults.InfraBucket)
	assert.Equal(t, "the-table", c.Defaults.InfraDynamoTable)
	assert.Equal(t, "test-project", c.Defaults.Project)
	assert.Equal(t, "0.11.0", c.Defaults.TerraformVersion)
	assert.Equal(t, true, c.Docker)
}

func TestParse(t *testing.T) {
	f, _ := os.Open("testdata/v1_full.json")
	defer f.Close()
	r := bufio.NewReader(f)
	c, e := ReadConfig(r)
	assert.Nil(t, e)
	assert.NotNil(t, c.Defaults)
	assert.Equal(t, int64(1), c.Defaults.AccountID)
	assert.Equal(t, "us-west-2", c.Defaults.AWSRegionBackend)
	assert.Equal(t, "us-west-1", c.Defaults.AWSRegionProvider)
	assert.Equal(t, "0.1.0", c.Defaults.AWSProviderVersion)
	assert.Equal(t, "czi", c.Defaults.AWSProfileBackend)
	assert.Equal(t, "czi", c.Defaults.AWSProfileProvider)
	assert.Equal(t, "the-bucket", c.Defaults.InfraBucket)
	assert.Equal(t, "the-table", c.Defaults.InfraDynamoTable)
	assert.Equal(t, "test-project", c.Defaults.Project)
	assert.Equal(t, "0.11.0", c.Defaults.TerraformVersion)

	assert.NotNil(t, c.Accounts)
	assert.Len(t, c.Accounts, 2)

	assert.NotNil(t, c.Envs)
	assert.Len(t, c.Envs, 1)
	assert.Len(t, c.Envs["stage"].Components, 2)
	env := c.Envs["stage"].Components["env"]
	assert.NotNil(t, env)
	assert.NotNil(t, env.ModuleSource)
	assert.Equal(t, "github.com/foo/bar", *env.ModuleSource)

	assert.NotNil(t, c.Modules)
}

func TestJsonFailure(t *testing.T) {
	json := `foo`
	r := ioutil.NopCloser(strings.NewReader(json))
	defer r.Close()
	c, e := ReadConfig(r)
	assert.Nil(t, c)
	assert.NotNil(t, e)
}

func TestValidation(t *testing.T) {
	json := `{}`
	r := ioutil.NopCloser(strings.NewReader(json))
	defer r.Close()
	c, e := ReadConfig(r)

	assert.NotNil(t, c)
	assert.Nil(t, e)

	e = c.Validate()
	assert.NotNil(t, e)

	_, ok := e.(*validator.InvalidValidationError)
	assert.False(t, ok)

	err, ok := e.(validator.ValidationErrors)
	assert.True(t, ok)
	assert.Len(t, err, 10)
}

func TestExtraVarsValidation(t *testing.T) {
	json := `
	{
		"defaults": {
			"aws_region_backend": "us-west-2",
			"account_id": 123456789,
			"aws_region_provider": "us-west-1",
			"aws_profile_backend": "czi",
			"aws_profile_provider": "czi",
			"aws_provider_version": "czi",
			"infra_s3_bucket": "the-bucket",
			"infra_dynamo_db_table": "the-table",
			"project": "test-project",
			"owner": "test@test.com",
			"terraform_version": "0.11.0"
		}
	}`
	r := ioutil.NopCloser(strings.NewReader(json))
	defer r.Close()
	c, e := ReadConfig(r)
	assert.Nil(t, e)

	e = c.Validate()
	assert.Nil(t, e)

	c.Defaults.ExtraVars = map[string]string{}
	c.Defaults.ExtraVars["env"] = "failme"
	e = c.Validate()
	assert.NotNil(t, e)
}

func TestInitConfig(t *testing.T) {
	c := InitConfig("proj", "reg", "buck", "table", "prof", "me@foo.example", "0.99.0")
	assert.Equal(t, "prof", c.Defaults.AWSProfileBackend)
	assert.Equal(t, "prof", c.Defaults.AWSProfileProvider)
	assert.Equal(t, "reg", c.Defaults.AWSRegionBackend)
	assert.Equal(t, "reg", c.Defaults.AWSRegionProvider)
	assert.Equal(t, "0.99.0", c.Defaults.AWSProviderVersion)
	assert.Equal(t, "buck", c.Defaults.InfraBucket)
	assert.Equal(t, "table", c.Defaults.InfraDynamoTable)
	assert.Equal(t, "me@foo.example", c.Defaults.Owner)
	assert.Equal(t, "proj", c.Defaults.Project)
	assert.Equal(t, "0.11.7", c.Defaults.TerraformVersion)
	assert.Equal(t, true, c.Docker)
}

func Test_detectVersion(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"implicit 1", args{[]byte(`{}`)}, 1, false},
		{"explicit 1", args{[]byte(`{"version": 1}`)}, 1, false},
		{"explicit 2", args{[]byte(`{"version": 2}`)}, 2, false},
		{"err", args{[]byte(`{`)}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := detectVersion(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("detectVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("detectVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func intptr(i int64) *int64 {
	return &i
}

func strptr(s string) *string {
	return &s
}

func boolptr(b bool) *bool {
	return &b
}

func TestUpgradeConfigVersion(t *testing.T) {
	a := assert.New(t)
	f, e := os.Open("testdata/v1_full.json")
	a.NoError(e)
	defer f.Close()
	r := bufio.NewReader(f)
	v1Full, e := ReadConfig(r)
	a.NoError(e)
	v2Full := &v2.Config{
		Version: 2,
		Docker:  false,
		Defaults: v2.Defaults{
			Common: v2.Common{
				Backend: v2.Backend{
					Bucket:      "the-bucket",
					DynamoTable: "the-table",
					Profile:     "czi",
					Region:      "us-west-2",
				},

				Providers: v2.Providers{
					AWS: &v2.AWSProvider{
						AccountID:         intptr(1),
						Profile:           strptr("czi"),
						Region:            strptr("us-west-1"),
						Version:           strptr("0.1.0"),
						AdditionalRegions: []string{"us-east-1"},
					},
				},
				Owner:            "default@example.com",
				Project:          "test-project",
				ExtraVars:        map[string]string{"foo": "bar"},
				TerraformVersion: "0.11.0",
			},
		},
		Tools: v2.Tools{
			TfLint: &v1.TfLint{
				Enabled: boolptr(true),
			},
			TravisCI: &v1.TravisCI{
				Enabled:        true,
				AWSIAMRoleName: "travis-role",
				TestBuckets:    13,
			},
		},
		Accounts: map[string]v2.Account{
			"foo": v2.Account{
				Common: v2.Common{
					Backend: v2.Backend{
						Bucket:      "foo-bucket",
						DynamoTable: "foo-table",
						Profile:     "czi-foo",
						Region:      "us-west-foo1",
					},

					Providers: v2.Providers{
						AWS: &v2.AWSProvider{
							AccountID:         intptr(2),
							Profile:           strptr("czi-foo"),
							Region:            strptr("us-west-foo1"),
							Version:           strptr("0.12.0"),
							AdditionalRegions: []string{"us-east-foo2"},
						},
					},
					Owner:            "foo@example.com",
					Project:          "foo-project",
					ExtraVars:        map[string]string{"foo": "foo"},
					TerraformVersion: "0.12.0",
				},
			},
			"bar": v2.Account{
				Common: v2.Common{
					Backend: v2.Backend{
						Bucket:      "bar-bucket",
						DynamoTable: "bar-table",
						Profile:     "czi-bar",
						Region:      "us-west-bar1",
					},

					Providers: v2.Providers{
						AWS: &v2.AWSProvider{
							AccountID:         intptr(3),
							Profile:           strptr("czi-bar"),
							Region:            strptr("us-west-bar1"),
							Version:           strptr("0.13.0"),
							AdditionalRegions: []string{"us-east-bar2"},
						},
					},
					Owner:            "bar@example.com",
					Project:          "bar-project",
					ExtraVars:        map[string]string{"foo": "bar"},
					TerraformVersion: "0.13.0",
				},
			},
		},

		Envs: map[string]v2.Env{
			"stage": v2.Env{
				Common: v2.Common{
					Backend: v2.Backend{
						Bucket:      "stage-bucket",
						DynamoTable: "stage-table",
						Profile:     "czi-stage",
						Region:      "us-west-stage1",
					},

					Providers: v2.Providers{
						AWS: &v2.AWSProvider{
							AccountID:         intptr(4),
							Profile:           strptr("czi-stage"),
							Region:            strptr("us-west-stage1"),
							Version:           strptr("0.14.0"),
							AdditionalRegions: []string{"us-east-stage2"},
						},
					},
					Owner:            "stage@example.com",
					Project:          "stage-project",
					ExtraVars:        map[string]string{"foo": "stage"},
					TerraformVersion: "0.14.0",
				},
				Components: map[string]v2.Component{
					"env": v2.Component{
						Common: v2.Common{
							Backend: v2.Backend{
								Bucket:      "env-bucket",
								DynamoTable: "env-table",
								Profile:     "czi-env",
								Region:      "us-west-env1",
							},

							Providers: v2.Providers{
								AWS: &v2.AWSProvider{
									AccountID:         intptr(5),
									Profile:           strptr("czi-env"),
									Region:            strptr("us-west-env1"),
									Version:           strptr("0.15.0"),
									AdditionalRegions: []string{"us-east-env2"},
								},
							},
							Owner:            "env@example.com",
							Project:          "env-project",
							ExtraVars:        map[string]string{"foo": "env"},
							TerraformVersion: "0.15.0",
						},
						ModuleSource: strptr("github.com/foo/bar"),
					},
					"helm": {},
				},
			},
		},

		Modules: map[string]v1.Module{
			"module1": v1.Module{TerraformVersion: strptr("0.15.0")},
		},
		Plugins: v1.Plugins{
			CustomPlugins: map[string]*plugins.CustomPlugin{
				"plugin": &plugins.CustomPlugin{
					URL:    "https://example.com/plugin.tgz",
					Format: plugins.TypePluginFormatTar,
					TarConfig: plugins.TarConfig{
						StripComponents: 7,
					},
				},
			},
			TerraformProviders: map[string]*plugins.CustomPlugin{
				"provider": &plugins.CustomPlugin{
					URL:    "https://example.com/provider.tgz",
					Format: plugins.TypePluginFormatTar,
					TarConfig: plugins.TarConfig{
						StripComponents: 7,
					},
				},
			},
		},
	}

	type args struct {
		c1 *v1.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *v2.Config
		wantErr bool
	}{
		{"v1 full", args{v1Full}, v2Full, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpgradeConfigVersion(tt.args.c1)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpgradeConfigVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			a.Equal(tt.want, got)
		})
	}
}
