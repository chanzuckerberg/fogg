package plan

import (
	"bufio"
	"os"
	"testing"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/stretchr/testify/assert"
)

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
	f, _ := os.Open("fixtures/full.json")
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

	assert.NotNil(t, plan.Envs["staging"].Components["cloud-env"])
	assert.Equal(t, "git@github.com:chanzuckerberg/shared-infra//terraform/modules/aws-env?ref=v0.10.0", *plan.Envs["staging"].Components["cloud-env"].ModuleSource)

	assert.NotNil(t, plan.Envs["staging"].Components["comp1"])
	assert.Equal(t, "0.100.0", plan.Envs["staging"].Components["comp1"].TerraformVersion)
}
