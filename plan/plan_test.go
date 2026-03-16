package plan

import (
	"encoding/json"
	"testing"

	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/chanzuckerberg/go-misc/ptr"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	formatter := &logrus.TextFormatter{
		DisableTimestamp: true,
	}
	logrus.SetFormatter(formatter)
}

func TestResolveAccounts(t *testing.T) {
	r := require.New(t)
	foo, bar := json.Number("123"), json.Number("456")

	accounts := map[string]v2.Account{
		"foo": {
			Common: v2.Common{
				Providers: &v2.Providers{
					AWS: &v2.AWSProvider{
						AccountID: &foo,
					},
				},
			},
		},
		"bar": {
			Common: v2.Common{
				Providers: &v2.Providers{
					AWS: &v2.AWSProvider{
						AccountID: &bar,
					},
				},
			},
		},
		"baz": {},
	}

	other := resolveAccounts(accounts)
	r.NotNil(other)
	r.Equal(map[string]*json.Number{"bar": &bar, "foo": &foo}, other)
}

func TestPlanBasicV2Yaml(t *testing.T) {
	r := require.New(t)

	b, e := util.TestFile("v2_full_yaml")
	r.NoError(e)

	fs, _, err := util.TestFs()
	r.NoError(err)
	err = afero.WriteFile(fs, "fogg.yml", b, 0644)
	r.NoError(err)
	c2, err := v2.ReadConfig(fs, b, "fogg.yml")
	r.Nil(err)

	w, err := c2.Validate()
	r.NoError(err)
	r.Len(w, 0)

	plan, e := Eval(c2)
	r.NoError(e)
	r.NotNil(plan)
	r.NotNil(plan.Accounts)
	r.Len(plan.Accounts, 2)

	r.NotNil(plan.Modules)
	r.Len(plan.Modules, 1)
	r.Equal("0.100.0", plan.Modules["my_module"].TerraformVersion)

	r.NotNil(plan.Envs)
	r.Len(plan.Envs, 2)

	r.NotNil(plan.Envs["staging"])

	r.NotNil(plan.Envs["staging"].Components)
	r.Len(plan.Envs["staging"].Components, 4)

	r.NotNil(plan.Envs["staging"])
	r.NotNil(plan.Envs["staging"].Components["vpc"])
	r.NotNil(plan.Envs["staging"].Components["k8s-comp"])
	r.NotNil(plan.Envs["staging"].Components["k8s-comp"].ProviderConfiguration.Kubernetes)
	r.NotNil(plan.Envs["staging"].Components["k8s-comp"].ProviderConfiguration.Kubernetes.ClusterComponentName)

	r.NotNil(*plan.Envs["staging"].Components["vpc"].ModuleSource)
	r.Equal("github.com/terraform-aws-modules/terraform-aws-vpc?ref=v1.30.0", *plan.Envs["staging"].Components["vpc"].ModuleSource)
	r.Nil(plan.Envs["staging"].Components["vpc"].ModuleName)

	r.NotNil(plan.Envs["staging"].Components["comp1"])
	r.Equal("0.100.0", plan.Envs["staging"].Components["comp1"].TerraformVersion)

	r.NotNil(plan.Envs["prod"])
	r.NotNil(plan.Envs["prod"].Components["hero"])
	r.NotNil(plan.Envs["prod"].Components["hero"].ProviderConfiguration)
	r.NotNil(plan.Envs["prod"].Components["hero"].ProviderConfiguration.Heroku)

	r.NotNil(plan.Envs["prod"])
	r.NotNil(plan.Envs["prod"].Components["datadog"])
	r.NotNil(plan.Envs["prod"].Components["datadog"].ProviderConfiguration)
	r.NotNil(plan.Envs["prod"].Components["datadog"].ProviderConfiguration.Datadog)

	r.NotNil(plan.Envs["prod"])
	r.NotNil(plan.Envs["prod"].Components["vpc"])
	r.NotNil(*plan.Envs["prod"].Components["vpc"].ModuleSource)
	r.Equal("github.com/terraform-aws-modules/terraform-aws-vpc?ref=v1.30.0", *plan.Envs["staging"].Components["vpc"].ModuleSource)
	r.Equal(*plan.Envs["prod"].Components["vpc"].ModuleName, "prod-vpc")

	// accts inherit defaults
	r.Equal("bar1", plan.Accounts["foo"].ExtraVars["foo"])
	// envs overwrite defaults
	r.Equal("bar2", plan.Envs["staging"].Components["comp1"].ExtraVars["foo"])
	// component overwrite env
	r.Equal("bar3", plan.Envs["staging"].Components["vpc"].ExtraVars["foo"])

	r.Equal("terraform/proj/accounts/bar.tfstate", plan.Accounts["bar"].Backend.S3.KeyPath)
	r.Equal("terraform/proj/accounts/foo.tfstate", plan.Accounts["foo"].Backend.S3.KeyPath)

	r.Len(plan.Accounts["bar"].AccountBackends, 2)
	r.NotNil(plan.Accounts["bar"].AccountBackends["foo"])
	r.Equal(BackendKindS3, plan.Accounts["bar"].AccountBackends["foo"].Kind)
	r.NotNil(plan.Accounts["bar"].AccountBackends["foo"].S3)
	r.Equal("terraform/proj/accounts/foo.tfstate", plan.Accounts["bar"].AccountBackends["foo"].S3.KeyPath)

	r.NotNil(plan.GitHubActionsCI)
	r.True(plan.GitHubActionsCI.Enabled)
}

func TestRemoteBackendPlan(t *testing.T) {
	r := require.New(t)

	plan := buildPlan(t, "remote_backend_yaml")

	r.NotNil(plan.Global)
	r.NotNil(plan.Global.Backend.Kind)
	r.Equal(plan.Global.Backend.Kind, BackendKindRemote)
}

func TestComponentKindNotTerraform(t *testing.T) {
	r := require.New(t)

	plan := buildPlan(t, "component_kind")
	backends := plan.Envs["env1"].Components["foo"].ComponentBackends

	// only contains 1 backend, for itself, and not for "bar"
	r.Len(backends, 1)
	r.Contains(backends, "foo")
}

func buildPlan(t *testing.T, testfile string) *Plan {
	r := require.New(t)
	b, e := util.TestFile(testfile)
	r.NoError(e)

	fs, _, err := util.TestFs()
	r.NoError(err)
	err = afero.WriteFile(fs, "fogg.yml", b, 0644)
	r.NoError(err)
	c2, err := v2.ReadConfig(fs, b, "fogg.yml")
	r.Nil(err)

	w, err := c2.Validate()
	r.NoError(err)
	r.Len(w, 0)

	plan, e := Eval(c2)
	r.NoError(e)
	r.NotNil(plan)

	return plan
}

func TestTfeProvider(t *testing.T) {
	r := require.New(t)

	plan := buildPlan(t, "tfe_provider_yaml")

	enabled := func(c ComponentCommon) {
		r.NotNil(c)
		r.NotNil(c.ProviderConfiguration.Tfe)
		r.True(c.ProviderConfiguration.Tfe.Enabled)
	}

	disabled := func(c ComponentCommon) {
		r.NotNil(c)
		r.Nil(c.ProviderConfiguration.Tfe)
	}

	enabled(plan.Global.ComponentCommon)
	enabled(plan.Accounts["foo"].ComponentCommon)
	disabled(plan.Envs["bar"].Components["bam"].ComponentCommon)
}

func TestSentryProvider(t *testing.T) {
	r := require.New(t)

	plan := buildPlan(t, "v2_full_yaml")

	enabled := func(c ComponentCommon) {
		r.NotNil(c)
		r.NotNil(c.ProviderConfiguration.Sentry)
		r.True(c.ProviderConfiguration.Sentry.Enabled)
	}

	disabled := func(c ComponentCommon) {
		r.NotNil(c)
		r.Nil(c.ProviderConfiguration.Sentry)
	}

	disabled(plan.Global.ComponentCommon)
	disabled(plan.Accounts["foo"].ComponentCommon)
	enabled(plan.Envs["prod"].Components["sentry"].ComponentCommon)
}

func TestOktaProvider(t *testing.T) {
	r := require.New(t)

	plan := buildPlan(t, "v2_full_yaml")

	enabled := func(c ComponentCommon) {
		r.NotNil(c)
		r.NotNil(c.ProviderConfiguration.Okta)
		r.Equal(c.ProviderConfiguration.Okta.BaseURL, ptr.String("https://foo.okta.com/"))
	}

	disabled := func(c ComponentCommon) {
		r.NotNil(c)
		r.Nil(c.ProviderConfiguration.Sentry)
	}

	disabled(plan.Global.ComponentCommon)
	disabled(plan.Accounts["foo"].ComponentCommon)
	enabled(plan.Envs["prod"].Components["okta"].ComponentCommon)
}

func TestGrafanaProvider(t *testing.T) {
	r := require.New(t)

	plan := buildPlan(t, "v2_full_yaml")

	enabled := func(c ComponentCommon) {
		r.NotNil(c)
		r.NotNil(c.ProviderConfiguration.Grafana)
		r.NotNil(c.ProviderVersions["grafana"])
	}

	// disabled := func(c ComponentCommon) {
	// 	r.NotNil(c)
	// 	r.Nil(c.ProviderConfiguration.Sentry)
	// }

	// disabled(plan.Global.ComponentCommon)
	// disabled(plan.Accounts["foo"].ComponentCommon)
	enabled(plan.Envs["prod"].Components["hero"].ComponentCommon)
}

func TestResolveCustomProviderConfig(t *testing.T) {
	roleArn := "arn:aws:iam::123456789:role/tfe-si"
	awsPlan := &AWSProvider{
		AccountID: json.Number("123456789"),
		Region:    "us-west-2",
		RoleArn:   &roleArn,
	}

	t.Run("nil config", func(t *testing.T) {
		r := require.New(t)
		result, err := resolveCustomProviderConfig(nil, customProviderTemplateContext{AWS: awsPlan})
		r.NoError(err)
		r.Nil(result)
	})

	t.Run("no templates", func(t *testing.T) {
		r := require.New(t)
		config := map[string]any{"region": "eu-west-1"}
		result, err := resolveCustomProviderConfig(config, customProviderTemplateContext{AWS: awsPlan})
		r.NoError(err)
		r.Equal("eu-west-1", result["region"])
	})

	t.Run("template with AWS region", func(t *testing.T) {
		r := require.New(t)
		config := map[string]any{"region": "{{ .AWS.Region }}"}
		result, err := resolveCustomProviderConfig(config, customProviderTemplateContext{AWS: awsPlan})
		r.NoError(err)
		r.Equal("us-west-2", result["region"])
	})

	t.Run("template with AWS RoleArn", func(t *testing.T) {
		r := require.New(t)
		config := map[string]any{
			"assume_role": map[string]any{
				"role_arn": "{{ .AWS.RoleArn }}",
			},
		}
		result, err := resolveCustomProviderConfig(config, customProviderTemplateContext{AWS: awsPlan})
		r.NoError(err)
		nested := result["assume_role"].(map[string]any)
		r.Equal("arn:aws:iam::123456789:role/tfe-si", nested["role_arn"])
	})

	t.Run("nil AWS provider", func(t *testing.T) {
		r := require.New(t)
		config := map[string]any{"region": "{{ .AWS.Region }}"}
		_, err := resolveCustomProviderConfig(config, customProviderTemplateContext{AWS: nil})
		r.Error(err)
	})

	t.Run("non-template values pass through", func(t *testing.T) {
		r := require.New(t)
		config := map[string]any{
			"count":   42,
			"enabled": true,
			"items":   []any{"a", "b"},
		}
		result, err := resolveCustomProviderConfig(config, customProviderTemplateContext{AWS: awsPlan})
		r.NoError(err)
		r.Equal(42, result["count"])
		r.Equal(true, result["enabled"])
		r.Equal([]any{"a", "b"}, result["items"])
	})

	t.Run("template with AWS AccountID", func(t *testing.T) {
		r := require.New(t)
		config := map[string]any{"account": "{{ .AWS.AccountID }}"}
		result, err := resolveCustomProviderConfig(config, customProviderTemplateContext{AWS: awsPlan})
		r.NoError(err)
		r.Equal("123456789", result["account"])
	})

	t.Run("mixed templates and literals in nested map", func(t *testing.T) {
		r := require.New(t)
		config := map[string]any{
			"region": "eu-west-1",
			"assume_role": map[string]any{
				"role_arn":     "{{ .AWS.RoleArn }}",
				"session_name": "my-session",
			},
		}
		result, err := resolveCustomProviderConfig(config, customProviderTemplateContext{AWS: awsPlan})
		r.NoError(err)
		r.Equal("eu-west-1", result["region"])
		nested := result["assume_role"].(map[string]any)
		r.Equal("arn:aws:iam::123456789:role/tfe-si", nested["role_arn"])
		r.Equal("my-session", nested["session_name"])
	})

	t.Run("template in list elements", func(t *testing.T) {
		r := require.New(t)
		config := map[string]any{
			"regions": []any{"{{ .AWS.Region }}", "eu-west-1"},
		}
		result, err := resolveCustomProviderConfig(config, customProviderTemplateContext{AWS: awsPlan})
		r.NoError(err)
		r.Equal([]any{"us-west-2", "eu-west-1"}, result["regions"])
	})

	t.Run("deeply nested template", func(t *testing.T) {
		r := require.New(t)
		config := map[string]any{
			"outer": map[string]any{
				"inner": map[string]any{
					"value": "{{ .AWS.Region }}",
				},
			},
		}
		result, err := resolveCustomProviderConfig(config, customProviderTemplateContext{AWS: awsPlan})
		r.NoError(err)
		outer := result["outer"].(map[string]any)
		inner := outer["inner"].(map[string]any)
		r.Equal("us-west-2", inner["value"])
	})

	t.Run("empty config returns empty map", func(t *testing.T) {
		r := require.New(t)
		config := map[string]any{}
		result, err := resolveCustomProviderConfig(config, customProviderTemplateContext{AWS: awsPlan})
		r.NoError(err)
		r.NotNil(result)
		r.Len(result, 0)
	})

	t.Run("empty map values preserved", func(t *testing.T) {
		r := require.New(t)
		config := map[string]any{
			"features": map[string]any{},
		}
		result, err := resolveCustomProviderConfig(config, customProviderTemplateContext{AWS: awsPlan})
		r.NoError(err)
		features := result["features"].(map[string]any)
		r.Len(features, 0)
	})

	t.Run("invalid template syntax returns error", func(t *testing.T) {
		r := require.New(t)
		config := map[string]any{"region": "{{ .Invalid"}
		_, err := resolveCustomProviderConfig(config, customProviderTemplateContext{AWS: awsPlan})
		r.Error(err)
	})
}
