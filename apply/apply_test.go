package apply

import (
	"testing"

	"github.com/chanzuckerberg/fogg/plan"
	"github.com/stretchr/testify/assert"
)

func TestRemoveExtension(t *testing.T) {
	var x string

	x = removeExtension("foo")
	assert.Equal(t, "foo", x)

	x = removeExtension("foo.txt")
	assert.Equal(t, "foo", x)

	x = removeExtension("foo.txt.asdf")
	assert.Equal(t, "foo.txt", x)
}

func TestJoinEnvs(t *testing.T) {
	var m map[string]plan.Env
	var x string

	m = map[string]plan.Env{
		"foo": plan.Env{},
	}
	x = joinEnvs(m)
	assert.NotNil(t, x)
	assert.Equal(t, "foo", x)

	m = map[string]plan.Env{
		"foo": plan.Env{},
		"bar": plan.Env{},
	}
	x = joinEnvs(m)
	assert.NotNil(t, x)
	assert.Equal(t, "bar foo", x)

}
