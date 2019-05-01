package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponentKindGetOrDefault(t *testing.T) {
	ck := ComponentKindHelmTemplate
	assert.Equal(t, string(ck.GetOrDefault()), string(ComponentKindHelmTemplate))

	var nck *ComponentKind
	assert.Equal(t, string(nck.GetOrDefault()), string(ComponentKindTerraform))

	var zck ComponentKind
	assert.Equal(t, string(zck.GetOrDefault()), string(ComponentKindTerraform))
}
