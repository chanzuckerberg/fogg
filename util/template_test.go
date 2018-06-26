package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDict(t *testing.T) {
	m := make(map[string]string)
	m["foo"] = "bar"
	r := dict(m)
	assert.NotNil(t, r)
	assert.IsType(t, map[string]interface{}{}, r)
	assert.Equal(t, "bar", r["foo"])
}
