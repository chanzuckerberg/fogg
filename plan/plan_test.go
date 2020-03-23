package plan

import (
	"encoding/json"
	"fmt"
	"testing"

	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	formatter := &logrus.TextFormatter{
		DisableTimestamp: true,
	}
	logrus.SetFormatter(formatter)
}

// func TestResolveRequired(t *testing.T) {
// 	resolved := resolveRequired("def", nil)
// 	assert.Equal(t, "def", resolved)

// 	over := "over"
// 	resolved = resolveRequired("def", &over)
// 	assert.Equal(t, "over", resolved)
// }

func TestResolveAccounts(t *testing.T) {
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
	assert.NotNil(t, other)
	var nilJsonNumber *json.Number
	assert.Equal(t, map[string]*json.Number{"bar": &bar, "foo": &foo, "baz": nilJsonNumber}, other)
}

func TestPlanBasicV2Yaml(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	b, e := util.TestFile("v2_full_yaml")
	assert.NoError(t, e)

	fs, _, err := util.TestFs()
	a.NoError(err)
	err = afero.WriteFile(fs, "fogg.yml", b, 0644)
	a.NoError(err)
	c2, err := v2.ReadConfig(fs, b, "fogg.yml")
	assert.Nil(t, err)

	w, err := c2.Validate()
	a.NoError(err)
	a.Len(w, 0)

	plan, e := Eval(c2)
	r.NoError(e)
	assert.NotNil(t, plan)
	assert.NotNil(t, plan.Accounts)
	assert.Len(t, plan.Accounts, 2)

	assert.NotNil(t, plan.Modules)
	assert.Len(t, plan.Modules, 1)
	assert.Equal(t, "0.100.0", plan.Modules["my_module"].TerraformVersion)

	assert.NotNil(t, plan.Envs)
	assert.Len(t, plan.Envs, 2)

	assert.NotNil(t, plan.Envs["staging"])

	assert.NotNil(t, plan.Envs["staging"].Components)
	assert.Len(t, plan.Envs["staging"].Components, 4)

	assert.NotNil(t, plan.Envs["staging"])
	assert.NotNil(t, plan.Envs["staging"].Components["vpc"])
	logrus.Debugf("%#v\n", plan.Envs["staging"].Components["vpc"].ModuleSource)
	assert.NotNil(t, *plan.Envs["staging"].Components["vpc"].ModuleSource)
	assert.Equal(t, "github.com/terraform-aws-modules/terraform-aws-vpc?ref=v1.30.0", *plan.Envs["staging"].Components["vpc"].ModuleSource)

	assert.NotNil(t, plan.Envs["staging"].Components["comp1"])
	assert.Equal(t, "0.100.0", plan.Envs["staging"].Components["comp1"].TerraformVersion)

	assert.NotNil(t, plan.Envs["staging"].Components["comp_helm_template"])
	assert.Equal(t, "k8s", plan.Envs["staging"].Components["comp_helm_template"].EKS.ClusterName)

	assert.NotNil(t, plan.Envs["prod"])
	assert.NotNil(t, plan.Envs["prod"].Components["hero"])
	assert.NotNil(t, plan.Envs["prod"].Components["hero"].Providers)
	assert.NotNil(t, plan.Envs["prod"].Components["hero"].Providers.Heroku)

	assert.NotNil(t, plan.Envs["prod"])
	assert.NotNil(t, plan.Envs["prod"].Components["datadog"])
	assert.NotNil(t, plan.Envs["prod"].Components["datadog"].Providers)
	assert.NotNil(t, plan.Envs["prod"].Components["datadog"].Providers.Datadog)

	// accts inherit defaults
	assert.Equal(t, "bar1", plan.Accounts["foo"].ExtraVars["foo"])
	// envs overwrite defaults
	assert.Equal(t, "bar2", plan.Envs["staging"].Components["comp1"].ExtraVars["foo"])
	// component overwrite env
	assert.Equal(t, "bar3", plan.Envs["staging"].Components["vpc"].ExtraVars["foo"])

	assert.Equal(t, "terraform/proj/accounts/bar.tfstate", plan.Accounts["bar"].Backend.S3.KeyPath)
	assert.Equal(t, "terraform/proj/accounts/foo.tfstate", plan.Accounts["foo"].Backend.S3.KeyPath)

	r.Len(plan.Accounts["foo"].Accounts, 1)
	r.NotNil(plan.Accounts["foo"].Accounts["bar"])
	r.NotNil(plan.Accounts["foo"].Accounts["bar"].Backend)
	fmt.Println("accts")
	r.Equal(BackendKindS3, plan.Accounts["foo"].Accounts["bar"].Backend.Kind)
	r.NotNil(plan.Accounts["foo"].Accounts["bar"].Backend.S3)
	// assert.Equal(t, "terraform/proj/accounts/bar.tfstate", plan.Accounts["foo"].Accounts["bar"].Backend.S3.KeyPath)

	r.Len(plan.Accounts["foo"].Accounts, 1)
}

func TestResolveEKSConfig(t *testing.T) {
	a := assert.New(t)
	a.Equal("", resolveEKSConfig(nil, nil).ClusterName)
	a.Equal("a", resolveEKSConfig(&v2.EKSConfig{ClusterName: "a"}, nil).ClusterName)
	a.Equal("b", resolveEKSConfig(&v2.EKSConfig{ClusterName: "a"}, &v2.EKSConfig{ClusterName: "b"}).ClusterName)
}

func TestRemoteBackendPlan(t *testing.T) {
	r := require.New(t)

	b, e := util.TestFile("remote_backend_yaml")
	r.NoError(e)

	fs, _, err := util.TestFs()
	r.NoError(err)
	err = afero.WriteFile(fs, "fogg.yml", b, 0644)
	r.NoError(err)
	c2, err := v2.ReadConfig(fs, b, "fogg.yml")
	r.Nil(err)

	fmt.Println("c2")

	w, err := c2.Validate()
	r.NoError(err)
	r.Len(w, 0)

	plan, e := Eval(c2)
	r.NoError(e)
	r.NotNil(plan)

	r.NotNil(plan.Global)
	r.NotNil(plan.Global.Backend.Kind)
	r.Equal(plan.Global.Backend.Kind, BackendKindRemote)
}
