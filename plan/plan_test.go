package plan

import (
	"testing"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/util"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.DebugLevel)
	formatter := &log.TextFormatter{
		DisableTimestamp: true,
	}
	log.SetFormatter(formatter)
}

// func TestResolveRequired(t *testing.T) {
// 	resolved := resolveRequired("def", nil)
// 	assert.Equal(t, "def", resolved)

// 	over := "over"
// 	resolved = resolveRequired("def", &over)
// 	assert.Equal(t, "over", resolved)
// }

func TestResolveAccounts(t *testing.T) {
	foo, bar := int64(123), int64(456)

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
	assert.Equal(t, map[string]int64{"bar": bar, "foo": foo}, other)
}

func TestPlanBasicV1(t *testing.T) {
	a := assert.New(t)
	b, e := util.TestFile("v1_full_plan")
	a.NoError(e)
	c, err := v1.ReadConfig(b)
	assert.Nil(t, err)

	c2, err := config.UpgradeConfigVersion(c)
	a.NoError(err)

	plan, e := Eval(c2)
	assert.Nil(t, e)
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
	log.Debugf("%#v\n", plan.Envs["staging"].Components["vpc"].ModuleSource)
	assert.NotNil(t, *plan.Envs["staging"].Components["vpc"].ModuleSource)
	assert.Equal(t, "github.com/terraform-aws-modules/terraform-aws-vpc?ref=v1.30.0", *plan.Envs["staging"].Components["vpc"].ModuleSource)

	assert.NotNil(t, plan.Envs["staging"].Components["comp1"])
	assert.Equal(t, "0.100.0", plan.Envs["staging"].Components["comp1"].TerraformVersion)

	assert.NotNil(t, plan.Envs["staging"].Components["comp_helm_template"])
	assert.Equal(t, "k8s", plan.Envs["staging"].Components["comp_helm_template"].EKS.ClusterName)
}

func TestPlanBasicV2(t *testing.T) {
	a := assert.New(t)

	b, e := util.TestFile("v2_full")
	assert.NoError(t, e)

	c2, err := v2.ReadConfig(b)
	assert.Nil(t, err)

	w, err := c2.Validate()
	a.NoError(err)
	a.Len(w, 0)

	plan, e := Eval(c2)
	assert.Nil(t, e)
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
	log.Debugf("%#v\n", plan.Envs["staging"].Components["vpc"].ModuleSource)
	assert.NotNil(t, *plan.Envs["staging"].Components["vpc"].ModuleSource)
	assert.Equal(t, "github.com/terraform-aws-modules/terraform-aws-vpc?ref=v1.30.0", *plan.Envs["staging"].Components["vpc"].ModuleSource)

	assert.NotNil(t, plan.Envs["staging"].Components["comp1"])
	assert.Equal(t, "0.100.0", plan.Envs["staging"].Components["comp1"].TerraformVersion)

	assert.NotNil(t, plan.Envs["staging"].Components["comp_helm_template"])
	assert.Equal(t, "k8s", plan.Envs["staging"].Components["comp_helm_template"].EKS.ClusterName)
}

func TestExtraVarsCompositionV2(t *testing.T) {
	a := assert.New(t)
	b, e := util.TestFile("v1_full_plan")
	a.NoError(e)
	c, err := v1.ReadConfig(b)
	assert.Nil(t, err)

	c2, err := config.UpgradeConfigVersion(c)
	a.NoError(err)

	plan, e := Eval(c2)
	assert.Nil(t, e)
	assert.NotNil(t, plan)

	// accts inherit defaults
	assert.Equal(t, "bar1", plan.Accounts["foo"].ExtraVars["foo"])
	// envs overwrite defaults
	assert.Equal(t, "bar2", plan.Envs["staging"].Components["comp1"].ExtraVars["foo"])
	// component overwrite env
	assert.Equal(t, "bar3", plan.Envs["staging"].Components["vpc"].ExtraVars["foo"])

}

func TestResolveTfLint(test *testing.T) {
	a := assert.New(test)
	t := true
	f := false

	data := []struct {
		def    *bool
		over   *bool
		output bool
	}{
		{nil, nil, false},
		{nil, &t, true},
		{nil, &f, false},
		{&t, nil, true},
		{&t, &t, true},
		{&t, &f, false},
		{&f, nil, false},
		{&f, &t, true},
		{&f, &f, false},
	}
	for _, r := range data {
		test.Run("", func(t *testing.T) {
			def := &v1.TfLint{Enabled: r.def}
			over := &v1.TfLint{Enabled: r.over}
			result := resolveTfLint(def, over)
			a.Equal(r.output, result.Enabled)
		})
	}
}

func TestResolveTfLintComponent(test *testing.T) {
	a := assert.New(test)
	t := true
	f := false

	data := []struct {
		def    bool
		over   *bool
		output bool
	}{
		{t, nil, true},
		{t, &t, true},
		{t, &f, false},
		{f, nil, false},
		{f, &t, true},
		{f, &f, false},
	}
	for _, r := range data {
		test.Run("", func(t *testing.T) {
			def := TfLint{Enabled: r.def}
			over := &v1.TfLint{Enabled: r.over}
			result := resolveTfLintComponent(def, over)
			a.Equal(r.output, result.Enabled)
		})
	}
}

func TestResolveEKSConfig(t *testing.T) {
	a := assert.New(t)
	a.Equal("", resolveEKSConfig(nil, nil).ClusterName)
	a.Equal("a", resolveEKSConfig(&v1.EKSConfig{ClusterName: "a"}, nil).ClusterName)
	a.Equal("b", resolveEKSConfig(&v1.EKSConfig{ClusterName: "a"}, &v1.EKSConfig{ClusterName: "b"}).ClusterName)
}
