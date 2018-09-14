package plan

import (
	"bufio"
	"os"
	"testing"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/plugins"
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
func TestResolveRequired(t *testing.T) {
	resolved := resolveRequired("def", nil)
	assert.Equal(t, "def", resolved)

	over := "over"
	resolved = resolveRequired("def", &over)
	assert.Equal(t, "over", resolved)
}

func TestResolveAccounts(t *testing.T) {
	foo, bar := int64(123), int64(456)

	accounts := map[string]config.Account{
		"foo": {
			AccountID: &foo,
		},
		"bar": {
			AccountID: &bar,
		},
		"baz": {},
	}

	other := resolveAccounts(accounts)
	assert.NotNil(t, other)
	assert.Equal(t, map[string]int64{"bar": bar, "foo": foo}, other)
}

func TestResolveStringArray(t *testing.T) {
	def := []string{"foo"}
	override := []string{"bar"}

	result := resolveStringArray(def, override)
	assert.Len(t, result, 1)
	assert.Equal(t, "bar", result[0])

	override = nil

	result2 := resolveStringArray(def, override)
	assert.Len(t, result2, 1)
	assert.Equal(t, "foo", result2[0])

}

func TestPlanBasic(t *testing.T) {
	f, _ := os.Open("testdata/full.json")
	defer f.Close()
	r := bufio.NewReader(f)
	c, err := config.ReadConfig(r)
	assert.Nil(t, err)

	plan, e := Eval(c, true, false)
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
	assert.Equal(t, plan.Envs["staging"].TerraformVersion, "0.100.0")

	assert.NotNil(t, plan.Envs["staging"].Components)
	assert.Len(t, plan.Envs["staging"].Components, 3)

	assert.NotNil(t, plan.Envs["staging"])
	assert.NotNil(t, plan.Envs["staging"].Components["vpc"])
	log.Debugf("%#v\n", plan.Envs["staging"].Components["vpc"].ModuleSource)
	assert.NotNil(t, *plan.Envs["staging"].Components["vpc"].ModuleSource)
	assert.Equal(t, "github.com/terraform-aws-modules/terraform-aws-vpc?ref=v1.30.0", *plan.Envs["staging"].Components["vpc"].ModuleSource)

	assert.NotNil(t, plan.Envs["staging"].Components["comp1"])
	assert.Equal(t, "0.100.0", plan.Envs["staging"].Components["comp1"].TerraformVersion)

	assert.NotNil(t, plan.Plugins.CustomPlugins)
	assert.Len(t, plan.Plugins.CustomPlugins, 1)
	assert.NotNil(t, plan.Plugins.CustomPlugins["custom"])
	assert.Equal(t, plugins.TypePluginFormatZip, plan.Plugins.CustomPlugins["custom"].Format)
	assert.Equal(t, "https://example.com/custom.zip", plan.Plugins.CustomPlugins["custom"].URL)

	assert.NotNil(t, plan.Plugins.TerraformProviders)
	assert.Len(t, plan.Plugins.TerraformProviders, 1)
	assert.NotNil(t, plan.Plugins.TerraformProviders["provider"])
	assert.Equal(t, plugins.TypePluginFormatTar, plan.Plugins.TerraformProviders["provider"].Format)
	assert.Equal(t, "https://example.com/provider.tar.gz", plan.Plugins.TerraformProviders["provider"].URL)
}

func TestPlanNoPlugins(t *testing.T) {
	f, _ := os.Open("testdata/full.json")
	defer f.Close()
	r := bufio.NewReader(f)
	c, err := config.ReadConfig(r)
	assert.Nil(t, err)

	plan, e := Eval(c, true, true)
	assert.Nil(t, e)
	assert.NotNil(t, plan)

	assert.Nil(t, plan.Plugins.CustomPlugins)
	assert.Nil(t, plan.Plugins.TerraformProviders)
}

func TestExtraVarsComposition(t *testing.T) {
	f, _ := os.Open("testdata/full.json")
	defer f.Close()
	r := bufio.NewReader(f)
	c, err := config.ReadConfig(r)
	assert.Nil(t, err)

	plan, e := Eval(c, true, false)
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
			def := &config.TfLint{Enabled: r.def}
			over := &config.TfLint{Enabled: r.over}
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
			over := &config.TfLint{Enabled: r.over}
			result := resolveTfLintComponent(def, over)
			a.Equal(r.output, result.Enabled)
		})
	}
}
