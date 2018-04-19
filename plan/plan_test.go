package plan

import (
	"testing"

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
