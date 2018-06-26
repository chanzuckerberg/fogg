package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionString(t *testing.T) {
	s := versionString("0.1.0", "abcdef", true, false)
	assert.Equal(t, "0.1.0", s)

	s = versionString("0.1.0", "abcdef", false, false)
	assert.Equal(t, "0.1.0-abcdef", s)

	s = versionString("0.1.0", "abcdef", false, true)
	assert.Equal(t, "0.1.0-abcdef-dirty", s)

}
