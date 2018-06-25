package plan

import (
	"testing"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestResolveRequired(t *testing.T) {
	var resolved string

	resolved = resolveRequired("def", nil)
	assert.Equal(t, "def", resolved)

	over := "over"
	resolved = resolveRequired("def", &over)
	assert.Equal(t, "over", resolved)
}

func TestResolveOtherAccounts(t *testing.T) {
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

	var other map[string]int64
	other = resolveOtherAccounts(accounts, "foo")
	assert.NotNil(t, other)
	assert.Equal(t, map[string]int64{"bar": bar}, other)
}

func TestResolveStringArray(t *testing.T) {
	def := []string{"foo"}
	override := &[]string{"bar"}

	result := resolveStringArray(def, override)
	assert.Len(t, result, 1)
	assert.Equal(t, "bar", result[0])

	override = nil

	result2 := resolveStringArray(def, override)
	assert.Len(t, result2, 1)
	assert.Equal(t, "foo", result2[0])

}

func TestPlanBasic(t *testing.T) {
	config := `
{
  "defaults": {
    "aws_region": "reg",
    "aws_profile": "prof",
    "infra_s3_bucket": "buck",
    "project": "proj",
    "terraform_version": "0.100.0",
    "owner": "foo@example.com"
  },
  "accounts": {
    "foo": {
      "account_id": 123
    },
    "bar": {
      "account_id": 456
    }
  },
  "modules": {
    "my_module": {}
  },
  "envs": {
    "staging":{
        "components": {
            "comp1": {},
            "comp2": {}
        }
    },
    "prod": {}
  }
}
`
	fs := afero.NewMemMapFs()
	afero.WriteFile(fs, "fogg.json", []byte(config), 0644)
	plan, e := Eval(fs, "fogg.json")
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
	assert.Len(t, plan.Envs["staging"].Components, 2)

	assert.NotNil(t, plan.Envs["staging"].Components["comp1"])
	assert.Equal(t, plan.Envs["staging"].Components["comp1"].TerraformVersion, "0.100.0")
}
