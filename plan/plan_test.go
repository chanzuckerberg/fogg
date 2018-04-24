package plan

import (
	"testing"

	"github.com/ryanking/fogg/config"
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
