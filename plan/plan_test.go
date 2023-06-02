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

	logrus.Debugf("%#v\n", plan.Envs["staging"].Components["vpc"].ModuleSource)
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
