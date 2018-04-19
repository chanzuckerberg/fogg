package plan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoalesceStrings(t *testing.T) {
	var r *string
	r = coalesceStrings([]*string{})
	assert.Nil(t, r)

	r = coalesceStrings([]*string{nil})
	assert.Nil(t, r)

	in := "foo"
	r = coalesceStrings([]*string{&in})
	assert.NotNil(t, r)
	assert.Equal(t, "foo", *r)

	a, b, c := "a", "b", "c"
	r = coalesceStrings([]*string{&a, &b, &c})
	assert.NotNil(t, r)
	assert.Equal(t, "a", *r)

	d := "d"
	r = coalesceStrings([]*string{nil, &d})
	assert.NotNil(t, r)
	assert.Equal(t, "d", *r)
}
