package plan

import (
	"testing"

	"github.com/ryanking/fogg/config"
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

func TestResolveOptional(t *testing.T) {
	def, override := "def", "override"
	var resolved *string

	resolved = resolveOptional(&def, &override)
	assert.NotNil(t, resolved)
	assert.Equal(t, *resolved, override)

	resolved = resolveOptional(&def, nil)
	assert.NotNil(t, resolved)
	assert.Equal(t, *resolved, def)

}

func TestResolveOtherAccounts(t *testing.T) {
	foo, bar := int64(123), int64(456)

	accounts := map[string]config.Account{
		"foo": config.Account{
			AccountId: &foo,
		},
		"bar": config.Account{
			AccountId: &bar,
		},
		"baz": config.Account{},
	}

	var other map[string]int64
	other = resolveOtherAccounts(accounts, "foo")
	assert.NotNil(t, other)
	assert.Equal(t, map[string]int64{"bar": bar}, other)
}

func TestResolveStringArray(t *testing.T) {
	def := &[]string{"foo"}
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
    "infra_bucket": "buck",
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
  }
}
`
	fs := afero.NewMemMapFs()
	afero.WriteFile(fs, "fogg.json", []byte(config), 0644)
	plan, _ := Plan(fs, "fogg.json")
	assert.NotNil(t, plan)
	assert.NotNil(t, plan.Accounts)
	assert.Len(t, plan.Accounts, 2)
}
